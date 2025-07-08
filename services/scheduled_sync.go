package services

import (
	"context"
	"log"
	"sync"
	"time"
)

// ScheduledSync handles periodic sync operations
type ScheduledSync struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

var (
	scheduledSync *ScheduledSync
	syncOnce      sync.Once
)

// GetScheduledSync returns the singleton scheduled sync service
func GetScheduledSync() *ScheduledSync {
	syncOnce.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		scheduledSync = &ScheduledSync{
			ctx:    ctx,
			cancel: cancel,
		}
		scheduledSync.start()
	})
	return scheduledSync
}

// start begins the scheduled sync operations
func (ss *ScheduledSync) start() {
	// Start full sync every 5 minutes
	ss.wg.Add(1)
	go ss.runFullSync(5 * time.Minute)

	// Start incremental sync every 30 seconds
	ss.wg.Add(1)
	go ss.runIncrementalSync(30 * time.Second)

	log.Println("Scheduled sync service started")
}

// runFullSync runs full data sync at specified intervals
func (ss *ScheduledSync) runFullSync(interval time.Duration) {
	defer ss.wg.Done()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Run initial sync
	log.Println("Running initial full sync...")
	if err := SyncAllData(); err != nil {
		log.Printf("Initial sync failed: %v", err)
	}

	for {
		select {
		case <-ticker.C:
			log.Println("Running scheduled full sync...")
			if err := SyncAllData(); err != nil {
				log.Printf("Scheduled full sync failed: %v", err)
			} else {
				log.Println("Scheduled full sync completed successfully")
			}
		case <-ss.ctx.Done():
			log.Println("Full sync service stopped")
			return
		}
	}
}

// runIncrementalSync runs incremental sync at specified intervals
func (ss *ScheduledSync) runIncrementalSync(interval time.Duration) {
	defer ss.wg.Done()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Running incremental sync...")
			if err := ss.runIncrementalSyncOperation(); err != nil {
				log.Printf("Incremental sync failed: %v", err)
			} else {
				log.Println("Incremental sync completed successfully")
			}
		case <-ss.ctx.Done():
			log.Println("Incremental sync service stopped")
			return
		}
	}
}

// runIncrementalSyncOperation performs incremental sync
func (ss *ScheduledSync) runIncrementalSyncOperation() error {
	// For now, just run a full sync
	// In a production system, you'd track last sync time and only sync new/modified records
	return SyncAllData()
}

// Stop gracefully stops the scheduled sync service
func (ss *ScheduledSync) Stop() {
	ss.cancel()
	ss.wg.Wait()
	log.Println("Scheduled sync service stopped")
}

// ForceSync triggers an immediate sync
func (ss *ScheduledSync) ForceSync() error {
	log.Println("Forcing immediate sync...")
	return SyncAllData()
}
