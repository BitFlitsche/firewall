package services

import (
	"context"
	"firewall/config"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

// Global sync lock to prevent conflicts between full and incremental sync
var (
	syncLock          sync.Mutex
	isFullSyncRunning int32 // atomic flag for full sync status
)

// ScheduledSync handles periodic sync operations
type ScheduledSync struct {
	ctx             context.Context
	cancel          context.CancelFunc
	wg              sync.WaitGroup
	incrementalSync *IncrementalSync
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
			ctx:             ctx,
			cancel:          cancel,
			incrementalSync: NewIncrementalSync(),
		}
		scheduledSync.start()
	})
	return scheduledSync
}

// IsFullSyncRunning returns true if a full sync is currently in progress
func IsFullSyncRunning() bool {
	return atomic.LoadInt32(&isFullSyncRunning) == 1
}

// SetFullSyncRunning sets the full sync status
func SetFullSyncRunning(running bool) {
	if running {
		atomic.StoreInt32(&isFullSyncRunning, 1)
	} else {
		atomic.StoreInt32(&isFullSyncRunning, 0)
	}
}

// start begins the scheduled sync operations
func (ss *ScheduledSync) start() {
	// Only run incremental sync every 30 seconds
	// Full sync should be run manually when needed
	ss.wg.Add(1)
	go ss.runIncrementalSync(30 * time.Second)

	log.Println("Scheduled sync service started (incremental only)")
}

// runFullSync runs full data sync at specified intervals
func (ss *ScheduledSync) runFullSync(interval time.Duration) {
	defer ss.wg.Done()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Run initial sync
	log.Println("Running initial full sync...")
	if err := ss.incrementalSync.ForceFullSync(); err != nil {
		log.Printf("Initial sync failed: %v", err)
	}

	for {
		select {
		case <-ticker.C:
			log.Println("Running scheduled full sync...")
			if err := ss.incrementalSync.ForceFullSync(); err != nil {
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

	// Get distributed lock service
	distributedLock := GetDistributedLock()

	for {
		select {
		case <-ticker.C:
			// Check if full sync is running before starting incremental sync
			if IsFullSyncRunning() {
				log.Println("Skipping incremental sync - full sync in progress")
				continue
			}

			// Try to acquire distributed lock for incremental sync
			lockName := "incremental_sync"
			lockTTL := config.AppConfig.Locking.IncrementalTTL

			acquired, lockInfo := distributedLock.TryAcquireLock(lockName, lockTTL)
			if !acquired {
				log.Println("Skipping incremental sync - another instance is running it")
				continue
			}

			log.Printf("Running incremental sync (lock: %s, instance: %s)...", lockInfo.LockID, lockInfo.Instance)

			// Ensure lock is released after sync operation
			defer func() {
				distributedLock.ReleaseLock(lockName)
				log.Printf("Released incremental sync lock")
			}()

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
	// Use the new incremental sync that only syncs changed records
	return ss.incrementalSync.SyncIncrementalAll()
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
	return ss.incrementalSync.ForceFullSync()
}
