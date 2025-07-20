package services

import (
	"bufio"
	"firewall/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
)

// StopForumSpamImportService handles importing data from StopForumSpam
type StopForumSpamImportService struct {
	db *gorm.DB
}

// NewStopForumSpamImportService creates a new StopForumSpam import service
func NewStopForumSpamImportService(db *gorm.DB) *StopForumSpamImportService {
	return &StopForumSpamImportService{
		db: db,
	}
}

// ImportToxicCIDRs imports toxic IP addresses in CIDR format from StopForumSpam
func (s *StopForumSpamImportService) ImportToxicCIDRs() error {
	log.Println("Starting StopForumSpam toxic CIDR import...")

	// URL for the toxic IP CIDR list
	url := "https://www.stopforumspam.com/downloads/toxic_ip_cidr.txt"

	// Download the file
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download StopForumSpam data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download StopForumSpam data: HTTP %d", resp.StatusCode)
	}

	// Read and parse the file
	cidrs, err := s.parseToxicCIDRFile(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to parse StopForumSpam data: %w", err)
	}

	log.Printf("Found %d toxic CIDR ranges to import", len(cidrs))

	// Start a transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	// Delete existing StopForumSpam toxic CIDR entries (but preserve manual entries)
	if err := tx.Where("source = ?", "stopforumspam_toxic_cidr").Delete(&models.IP{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete existing StopForumSpam entries: %w", err)
	}

	// Import new entries
	var importedCount int
	for _, cidr := range cidrs {
		ip := models.IP{
			Address: cidr,
			Status:  "denied",
			IsCIDR:  true,
			Source:  "stopforumspam_toxic_cidr",
		}

		if err := tx.Create(&ip).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create IP entry for %s: %w", cidr, err)
		}
		importedCount++
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Successfully imported %d toxic CIDR ranges from StopForumSpam", importedCount)

	// Publish event for async processing
	PublishEvent("ip", "imported", map[string]interface{}{
		"source": "stopforumspam_toxic_cidr",
		"count":  importedCount,
	})

	return nil
}

// parseToxicCIDRFile parses the StopForumSpam toxic CIDR file
func (s *StopForumSpamImportService) parseToxicCIDRFile(reader io.Reader) ([]string, error) {
	var cidrs []string
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Validate CIDR format
		if s.isValidCIDR(line) {
			cidrs = append(cidrs, line)
		} else {
			log.Printf("Warning: Invalid CIDR format: %s", line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return cidrs, nil
}

// isValidCIDR checks if a string is a valid CIDR notation
func (s *StopForumSpamImportService) isValidCIDR(cidr string) bool {
	// Basic CIDR validation
	parts := strings.Split(cidr, "/")
	if len(parts) != 2 {
		return false
	}

	// Check if it's a valid IP address
	if !s.isValidIP(parts[0]) {
		return false
	}

	// Check if prefix length is valid
	prefixLen := 0
	if _, err := fmt.Sscanf(parts[1], "%d", &prefixLen); err != nil {
		return false
	}

	// Valid prefix lengths: 0-32 for IPv4, 0-128 for IPv6
	// For now, we'll assume IPv4 (0-32)
	if prefixLen < 0 || prefixLen > 32 {
		return false
	}

	return true
}

// isValidIP checks if a string is a valid IP address
func (s *StopForumSpamImportService) isValidIP(ip string) bool {
	// Basic IP validation - this is a simplified version
	// In production, you might want to use a more robust IP validation library
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}

	for _, part := range parts {
		var num int
		if _, err := fmt.Sscanf(part, "%d", &num); err != nil {
			return false
		}
		if num < 0 || num > 255 {
			return false
		}
	}

	return true
}

// GetStopForumSpamImportStats returns statistics about StopForumSpam imports
func (s *StopForumSpamImportService) GetStopForumSpamImportStats() (map[string]interface{}, error) {
	var count int64
	if err := s.db.Model(&models.IP{}).Where("source = ?", "stopforumspam_toxic_cidr").Count(&count).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_stopforumspam_cidrs": count,
		"last_import":               time.Now().Format(time.RFC3339), // This would be better stored in a separate table
	}, nil
}

// GetStopForumSpamImportStatus returns the current status of StopForumSpam imports
func (s *StopForumSpamImportService) GetStopForumSpamImportStatus() (map[string]interface{}, error) {
	var count int64
	if err := s.db.Model(&models.IP{}).Where("source = ?", "stopforumspam_toxic_cidr").Count(&count).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"import_enabled": true,
		"is_running":     false, // This would be tracked in a separate table
		"total_imported": count,
		"last_import":    time.Now().Format(time.RFC3339), // This would be stored in a separate table
	}, nil
}
