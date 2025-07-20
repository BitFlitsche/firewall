package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"firewall/config"
	"firewall/models"
)

// SpamhausASNRecord represents a single ASN record from Spamhaus
type SpamhausASNRecord struct {
	ASN    int    `json:"asn"`
	RIR    string `json:"rir"`
	Domain string `json:"domain"`
	CC     string `json:"cc"`
	ASName string `json:"asname"`
}

// ImportSpamhausASNDrop imports ASN data from Spamhaus ASN-DROP list
func ImportSpamhausASNDrop() error {
	// Get configured URL
	importURL := config.AppConfig.Spamhaus.ImportURL
	if importURL == "" {
		importURL = "https://www.spamhaus.org/drop/asndrop.json"
	}

	// Fetch data from Spamhaus
	resp, err := http.Get(importURL)
	if err != nil {
		return fmt.Errorf("failed to fetch Spamhaus data from %s: %w", importURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Spamhaus API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the JSONL response (each line is a separate JSON object)
	lines := strings.Split(string(body), "\n")
	var records []SpamhausASNRecord
	var skippedLines int

	fmt.Printf("Processing %d lines from Spamhaus response...\n", len(lines))

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Skip metadata lines that start with {"type":
		if strings.HasPrefix(line, `{"type":`) {
			skippedLines++
			continue
		}

		var record SpamhausASNRecord
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			// Skip invalid lines but continue processing
			linePreview := line
			if len(line) > 100 {
				linePreview = line[:100] + "..."
			}
			fmt.Printf("Warning: Skipping invalid line %d: %s\n", i+1, linePreview)
			skippedLines++
			continue
		}

		// Only add records that have valid ASN data
		if record.ASN > 0 {
			records = append(records, record)
		} else {
			skippedLines++
		}
	}

	fmt.Printf("Parsed %d valid ASN records, skipped %d lines\n", len(records), skippedLines)

	// Get database connection
	db := config.DB
	if db == nil {
		return fmt.Errorf("database connection not available")
	}

	// Begin transaction
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Delete existing Spamhaus records (but preserve manual entries)
	if err := tx.Where("source = ?", "spamhaus").Delete(&models.ASN{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete existing Spamhaus records: %w", err)
	}

	// Import new records
	var importedCount int
	var skippedCount int
	for _, record := range records {
		// Convert ASN number to string format
		asnString := fmt.Sprintf("AS%d", record.ASN)

		// Check if this ASN already exists as a manual entry
		var existingASN models.ASN
		if err := tx.Where("asn = ? AND source = ?", asnString, "manual").First(&existingASN).Error; err == nil {
			// ASN exists as manual entry, skip it
			skippedCount++
			continue
		}

		// Create ASN record
		asnRecord := models.ASN{
			ASN:     asnString,
			RIR:     record.RIR,
			Domain:  record.Domain,
			Country: strings.ToUpper(record.CC),
			Name:    record.ASName,
			Status:  "denied", // Spamhaus ASN-DROP list contains only denied ASNs
			Source:  "spamhaus",
		}

		if err := tx.Create(&asnRecord).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create ASN record %s: %w", asnString, err)
		}
		importedCount++
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Sync to Elasticsearch
	if err := SyncAllASNs(); err != nil {
		return fmt.Errorf("failed to sync ASNs to Elasticsearch: %w", err)
	}

	// Update sync tracker
	if err := updateSyncTracker("asns"); err != nil {
		return fmt.Errorf("failed to update sync tracker: %w", err)
	}

	fmt.Printf("Successfully imported %d ASN records from Spamhaus ASN-DROP list (skipped %d manual entries)\n", importedCount, skippedCount)
	return nil
}

// updateSyncTracker updates the last sync timestamp for ASNs
func updateSyncTracker(dataType string) error {
	db := config.DB
	if db == nil {
		return fmt.Errorf("database connection not available")
	}

	var tracker models.SyncTracker
	result := db.Where("data_type = ?", dataType).First(&tracker)

	if result.Error != nil {
		// Create new tracker if it doesn't exist
		tracker = models.SyncTracker{
			DataType: dataType,
			LastSync: time.Now(),
		}
		return db.Create(&tracker).Error
	}

	// Update existing tracker
	tracker.LastSync = time.Now()
	return db.Save(&tracker).Error
}

// GetSpamhausImportStats returns statistics about the Spamhaus import
func GetSpamhausImportStats() (map[string]interface{}, error) {
	db := config.DB
	if db == nil {
		return nil, fmt.Errorf("database connection not available")
	}

	var stats struct {
		TotalSpamhausASNs int64
		LastSync          *time.Time
	}

	// Count Spamhaus ASNs
	if err := db.Model(&models.ASN{}).Where("source = ?", "spamhaus").Count(&stats.TotalSpamhausASNs).Error; err != nil {
		return nil, fmt.Errorf("failed to count Spamhaus ASNs: %w", err)
	}

	// Get last sync time
	var tracker models.SyncTracker
	if err := db.Where("data_type = ?", "asns").First(&tracker).Error; err == nil {
		stats.LastSync = &tracker.LastSync
	}

	return map[string]interface{}{
		"total_spamhaus_asns": stats.TotalSpamhausASNs,
		"last_sync":           stats.LastSync,
	}, nil
}
