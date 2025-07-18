package services

import (
	"firewall/config"
	"firewall/models"
	"fmt"
	"log"
	"time"
)

// IncrementalSync handles incremental sync operations
type IncrementalSync struct {
	lastSyncTimes map[string]time.Time
}

// NewIncrementalSync creates a new incremental sync instance
func NewIncrementalSync() *IncrementalSync {
	return &IncrementalSync{
		lastSyncTimes: make(map[string]time.Time),
	}
}

// getLastSyncTime gets the last sync time for a data type
func (is *IncrementalSync) getLastSyncTime(dataType string) time.Time {
	var tracker models.SyncTracker
	if err := config.DB.Where("data_type = ?", dataType).First(&tracker).Error; err != nil {
		// If no record exists, return a very old time to sync everything
		return time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	return tracker.LastSync
}

// updateLastSyncTime updates the last sync time for a data type
func (is *IncrementalSync) updateLastSyncTime(dataType string) error {
	now := time.Now()

	// Try to update existing record, if not exists create new one
	result := config.DB.Where("data_type = ?", dataType).Updates(&models.SyncTracker{
		DataType: dataType,
		LastSync: now,
	})

	if result.RowsAffected == 0 {
		// No record was updated, create new one
		return config.DB.Create(&models.SyncTracker{
			DataType: dataType,
			LastSync: now,
		}).Error
	}

	return result.Error
}

// SyncIncrementalIPs syncs only IPs modified since last sync
func (is *IncrementalSync) SyncIncrementalIPs() error {
	lastSync := is.getLastSyncTime("ips")

	var ips []models.IP
	query := config.DB.Where("updated_at > ? OR created_at > ?", lastSync, lastSync)
	if err := query.Find(&ips).Error; err != nil {
		return err
	}

	if len(ips) == 0 {
		log.Println("No IPs to sync incrementally")
		return nil
	}

	syncedCount := 0
	for _, ip := range ips {
		if err := IndexIPAddress(ip); err != nil {
			log.Printf("Error syncing IP %s: %v", ip.Address, err)
		} else {
			syncedCount++
		}
	}

	if syncedCount > 0 {
		if err := is.updateLastSyncTime("ips"); err != nil {
			log.Printf("Error updating IP sync time: %v", err)
		}
		log.Printf("Incrementally synced %d IPs to Elasticsearch", syncedCount)
	}

	return nil
}

// SyncIncrementalEmails syncs only emails modified since last sync
func (is *IncrementalSync) SyncIncrementalEmails() error {
	lastSync := is.getLastSyncTime("emails")

	var emails []models.Email
	query := config.DB.Where("updated_at > ? OR created_at > ?", lastSync, lastSync)
	if err := query.Find(&emails).Error; err != nil {
		return err
	}

	if len(emails) == 0 {
		log.Println("No emails to sync incrementally")
		return nil
	}

	syncedCount := 0
	for _, email := range emails {
		if err := IndexEmail(email); err != nil {
			log.Printf("Error syncing email %s: %v", email.Address, err)
		} else {
			syncedCount++
		}
	}

	if syncedCount > 0 {
		if err := is.updateLastSyncTime("emails"); err != nil {
			log.Printf("Error updating email sync time: %v", err)
		}
		log.Printf("Incrementally synced %d emails to Elasticsearch", syncedCount)
	}

	return nil
}

// SyncIncrementalUserAgents syncs only user agents modified since last sync
func (is *IncrementalSync) SyncIncrementalUserAgents() error {
	lastSync := is.getLastSyncTime("user_agents")

	var userAgents []models.UserAgent
	query := config.DB.Where("updated_at > ? OR created_at > ?", lastSync, lastSync)
	if err := query.Find(&userAgents).Error; err != nil {
		return err
	}

	if len(userAgents) == 0 {
		log.Println("No user agents to sync incrementally")
		return nil
	}

	syncedCount := 0
	for _, ua := range userAgents {
		if err := IndexUserAgent(ua); err != nil {
			log.Printf("Error syncing user agent %s: %v", ua.UserAgent, err)
		} else {
			syncedCount++
		}
	}

	if syncedCount > 0 {
		if err := is.updateLastSyncTime("user_agents"); err != nil {
			log.Printf("Error updating user agent sync time: %v", err)
		}
		log.Printf("Incrementally synced %d user agents to Elasticsearch", syncedCount)
	}

	return nil
}

// SyncIncrementalCountries syncs only countries modified since last sync
func (is *IncrementalSync) SyncIncrementalCountries() error {
	lastSync := is.getLastSyncTime("countries")

	var countries []models.Country
	query := config.DB.Where("updated_at > ? OR created_at > ?", lastSync, lastSync)
	if err := query.Find(&countries).Error; err != nil {
		return err
	}

	if len(countries) == 0 {
		log.Println("No countries to sync incrementally")
		return nil
	}

	syncedCount := 0
	for _, country := range countries {
		if err := IndexCountry(country); err != nil {
			log.Printf("Error syncing country %s: %v", country.Code, err)
		} else {
			syncedCount++
		}
	}

	if syncedCount > 0 {
		if err := is.updateLastSyncTime("countries"); err != nil {
			log.Printf("Error updating country sync time: %v", err)
		}
		log.Printf("Incrementally synced %d countries to Elasticsearch", syncedCount)
	}

	return nil
}

// SyncIncrementalCharsetRules syncs only charset rules modified since last sync
func (is *IncrementalSync) SyncIncrementalCharsetRules() error {
	lastSync := is.getLastSyncTime("charsets")

	var charsets []models.CharsetRule
	query := config.DB.Where("updated_at > ? OR created_at > ?", lastSync, lastSync)
	if err := query.Find(&charsets).Error; err != nil {
		return err
	}

	if len(charsets) == 0 {
		log.Println("No charset rules to sync incrementally")
		return nil
	}

	syncedCount := 0
	for _, charset := range charsets {
		if err := IndexCharsetRule(charset); err != nil {
			log.Printf("Error syncing charset rule %s: %v", charset.Charset, err)
		} else {
			syncedCount++
		}
	}

	if syncedCount > 0 {
		if err := is.updateLastSyncTime("charsets"); err != nil {
			log.Printf("Error updating charset sync time: %v", err)
		}
		log.Printf("Incrementally synced %d charset rules to Elasticsearch", syncedCount)
	}

	return nil
}

// SyncIncrementalUsernameRules syncs only username rules modified since last sync
func (is *IncrementalSync) SyncIncrementalUsernameRules() error {
	lastSync := is.getLastSyncTime("usernames")

	var usernames []models.UsernameRule
	query := config.DB.Where("updated_at > ? OR created_at > ?", lastSync, lastSync)
	if err := query.Find(&usernames).Error; err != nil {
		return err
	}

	if len(usernames) == 0 {
		log.Println("No username rules to sync incrementally")
		return nil
	}

	syncedCount := 0
	for _, username := range usernames {
		if err := IndexUsernameRule(username); err != nil {
			log.Printf("Error syncing username rule %s: %v", username.Username, err)
		} else {
			syncedCount++
		}
	}

	if syncedCount > 0 {
		if err := is.updateLastSyncTime("usernames"); err != nil {
			log.Printf("Error updating username sync time: %v", err)
		}
		log.Printf("Incrementally synced %d username rules to Elasticsearch", syncedCount)
	}

	return nil
}

// SyncIncrementalAll syncs all data types incrementally
func (is *IncrementalSync) SyncIncrementalAll() error {
	// Check if full sync is running before starting incremental sync
	if IsFullSyncRunning() {
		log.Println("Skipping incremental sync - full sync in progress")
		return nil
	}

	log.Println("Starting incremental sync to Elasticsearch...")

	if err := is.SyncIncrementalIPs(); err != nil {
		log.Printf("Error in incremental IP sync: %v", err)
	}

	if err := is.SyncIncrementalEmails(); err != nil {
		log.Printf("Error in incremental email sync: %v", err)
	}

	if err := is.SyncIncrementalUserAgents(); err != nil {
		log.Printf("Error in incremental user agent sync: %v", err)
	}

	if err := is.SyncIncrementalCountries(); err != nil {
		log.Printf("Error in incremental country sync: %v", err)
	}

	if err := is.SyncIncrementalCharsetRules(); err != nil {
		log.Printf("Error in incremental charset rule sync: %v", err)
	}

	if err := is.SyncIncrementalUsernameRules(); err != nil {
		log.Printf("Error in incremental username rule sync: %v", err)
	}

	log.Println("Incremental sync completed")
	return nil
}

// ForceFullSync performs a full sync and updates all sync timestamps
func (is *IncrementalSync) ForceFullSync() error {
	log.Println("Forcing full sync to Elasticsearch...")

	// Get distributed lock service
	distributedLock := GetDistributedLock()

	// Try to acquire distributed lock for full sync
	lockName := "full_sync"
	lockTTL := config.AppConfig.Locking.FullSyncTTL

	acquired, lockInfo := distributedLock.TryAcquireLock(lockName, lockTTL)
	if !acquired {
		return fmt.Errorf("full sync already in progress by another instance")
	}

	log.Printf("Acquired full sync lock (instance: %s)", lockInfo.Instance)

	// Ensure lock is released after sync operation
	defer func() {
		distributedLock.ReleaseLock(lockName)
		log.Printf("Released full sync lock")
	}()

	// Set full sync running flag to prevent incremental sync conflicts
	SetFullSyncRunning(true)
	defer SetFullSyncRunning(false)

	// Perform full sync
	if err := SyncAllData(); err != nil {
		return err
	}

	// Update all sync timestamps
	dataTypes := []string{"ips", "emails", "user_agents", "countries", "charsets", "usernames"}

	for _, dataType := range dataTypes {
		if err := is.updateLastSyncTime(dataType); err != nil {
			log.Printf("Error updating sync time for %s: %v", dataType, err)
		}
	}

	log.Println("Full sync completed and timestamps updated")
	return nil
}
