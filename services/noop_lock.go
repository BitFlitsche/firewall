package services

import (
	"time"
)

// NoOpDistributedLock provides a no-op implementation when distributed locking is disabled
type NoOpDistributedLock struct {
	instanceID string
}

// NewNoOpDistributedLock creates a new no-op distributed lock
func NewNoOpDistributedLock() *NoOpDistributedLock {
	return &NoOpDistributedLock{
		instanceID: "single-instance",
	}
}

// TryAcquireLock always succeeds for no-op implementation
func (dl *NoOpDistributedLock) TryAcquireLock(lockName string, ttl time.Duration) (bool, *LockInfo) {
	return true, &LockInfo{
		LockID:    lockName,
		Instance:  dl.instanceID,
		Acquired:  time.Now(),
		ExpiresAt: time.Now().Add(ttl),
	}
}

// ReleaseLock always succeeds for no-op implementation
func (dl *NoOpDistributedLock) ReleaseLock(lockName string) bool {
	return true
}

// IsLocked always returns false for no-op implementation (no locks are held)
func (dl *NoOpDistributedLock) IsLocked(lockName string) bool {
	return false
}

// GetLockInfo returns nil for no-op implementation
func (dl *NoOpDistributedLock) GetLockInfo(lockName string) (*LockInfo, error) {
	return nil, nil
}

// ExtendLock always succeeds for no-op implementation
func (dl *NoOpDistributedLock) ExtendLock(lockName string, ttl time.Duration) bool {
	return true
}

// GetActiveLocks returns empty slice for no-op implementation
func (dl *NoOpDistributedLock) GetActiveLocks() ([]*LockInfo, error) {
	return []*LockInfo{}, nil
}

// CleanupExpiredLocks does nothing for no-op implementation
func (dl *NoOpDistributedLock) CleanupExpiredLocks() {
	// No-op
}

// Stop does nothing for no-op implementation
func (dl *NoOpDistributedLock) Stop() {
	// No-op
}
