package services

import (
	"context"
	"firewall/config"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// DistributedLockInterface defines the interface for distributed locking
type DistributedLockInterface interface {
	TryAcquireLock(lockName string, ttl time.Duration) (bool, *LockInfo)
	ReleaseLock(lockName string) bool
	IsLocked(lockName string) bool
	GetLockInfo(lockName string) (*LockInfo, error)
	ExtendLock(lockName string, ttl time.Duration) bool
	GetActiveLocks() ([]*LockInfo, error)
	CleanupExpiredLocks()
	Stop()
}

// DistributedLock provides distributed locking using Redis
type DistributedLock struct {
	client *redis.Client
	ctx    context.Context
}

// LockInfo represents lock metadata
type LockInfo struct {
	LockID    string    `json:"lock_id"`
	Instance  string    `json:"instance"`
	Acquired  time.Time `json:"acquired"`
	ExpiresAt time.Time `json:"expires_at"`
}

var (
	distributedLock DistributedLockInterface
	lockOnce        sync.Once
	instanceID      string
)

// GetDistributedLock returns the singleton distributed lock service
func GetDistributedLock() DistributedLockInterface {
	lockOnce.Do(func() {
		// Check if distributed locking is enabled
		if !config.AppConfig.Locking.Enabled {
			log.Println("Distributed locking disabled - using no-op implementation")
			distributedLock = NewNoOpDistributedLock()
			return
		}

		// Generate unique instance ID
		instanceID = generateInstanceID()

		// Initialize Redis client
		redisClient := redis.NewClient(&redis.Options{
			Addr:     config.AppConfig.Redis.GetRedisAddr(),
			Password: config.AppConfig.Redis.Password,
			DB:       config.AppConfig.Redis.DB,
		})

		distributedLock = &DistributedLock{
			client: redisClient,
			ctx:    context.Background(),
		}

		log.Printf("Distributed lock service initialized with instance ID: %s", instanceID)
	})
	return distributedLock
}

// generateInstanceID creates a unique instance identifier
func generateInstanceID() string {
	hostname, _ := os.Hostname()
	pid := os.Getpid()
	return fmt.Sprintf("%s-%d-%d", hostname, pid, time.Now().Unix())
}

// TryAcquireLock attempts to acquire a distributed lock
func (dl *DistributedLock) TryAcquireLock(lockName string, ttl time.Duration) (bool, *LockInfo) {
	lockKey := fmt.Sprintf("lock:%s", lockName)
	lockValue := fmt.Sprintf("%s:%s", instanceID, time.Now().Format(time.RFC3339))

	// Try to acquire lock using Redis SET with NX and EX
	result := dl.client.Set(dl.ctx, lockKey, lockValue, ttl)

	if result.Err() != nil {
		log.Printf("Error acquiring lock %s: %v", lockName, result.Err())
		return false, nil
	}

	// Check if lock was acquired (SET NX returns OK only if key didn't exist)
	if result.Val() == "OK" {
		lockInfo := &LockInfo{
			LockID:    lockName,
			Instance:  instanceID,
			Acquired:  time.Now(),
			ExpiresAt: time.Now().Add(ttl),
		}

		log.Printf("Lock acquired: %s by instance %s", lockName, instanceID)
		return true, lockInfo
	}

	// Lock was not acquired
	return false, nil
}

// ReleaseLock releases a distributed lock (only if owned by this instance)
func (dl *DistributedLock) ReleaseLock(lockName string) bool {
	lockKey := fmt.Sprintf("lock:%s", lockName)

	// Use Lua script for atomic check-and-delete
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`

	lockValue := fmt.Sprintf("%s:%s", instanceID, time.Now().Format(time.RFC3339))
	result := dl.client.Eval(dl.ctx, script, []string{lockKey}, lockValue)

	if result.Err() != nil {
		log.Printf("Error releasing lock %s: %v", lockName, result.Err())
		return false
	}

	if result.Val().(int64) == 1 {
		log.Printf("Lock released: %s by instance %s", lockName, instanceID)
		return true
	}

	log.Printf("Lock not released: %s (not owned by this instance)", lockName)
	return false
}

// IsLocked checks if a lock is currently held
func (dl *DistributedLock) IsLocked(lockName string) bool {
	lockKey := fmt.Sprintf("lock:%s", lockName)
	result := dl.client.Exists(dl.ctx, lockKey)

	if result.Err() != nil {
		log.Printf("Error checking lock %s: %v", lockName, result.Err())
		return false
	}

	return result.Val() > 0
}

// GetLockInfo returns information about a lock
func (dl *DistributedLock) GetLockInfo(lockName string) (*LockInfo, error) {
	lockKey := fmt.Sprintf("lock:%s", lockName)
	result := dl.client.Get(dl.ctx, lockKey)

	if result.Err() != nil {
		if result.Err() == redis.Nil {
			return nil, nil // Lock doesn't exist
		}
		return nil, result.Err()
	}

	// Parse lock value (instanceID:timestamp)
	lockValue := result.Val()
	parts := strings.Split(lockValue, ":")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid lock value format: %s", lockValue)
	}

	instance := parts[0]
	timestamp, err := time.Parse(time.RFC3339, strings.Join(parts[1:], ":"))
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp in lock value: %v", err)
	}

	// Get TTL to calculate expiration
	ttlResult := dl.client.TTL(dl.ctx, lockKey)
	if ttlResult.Err() != nil {
		return nil, ttlResult.Err()
	}

	expiresAt := time.Now().Add(ttlResult.Val())

	return &LockInfo{
		LockID:    lockName,
		Instance:  instance,
		Acquired:  timestamp,
		ExpiresAt: expiresAt,
	}, nil
}

// ExtendLock extends the TTL of a lock (only if owned by this instance)
func (dl *DistributedLock) ExtendLock(lockName string, ttl time.Duration) bool {
	lockKey := fmt.Sprintf("lock:%s", lockName)
	lockValue := fmt.Sprintf("%s:%s", instanceID, time.Now().Format(time.RFC3339))

	// Use Lua script for atomic check-and-expire
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("expire", KEYS[1], ARGV[2])
		else
			return 0
		end
	`

	result := dl.client.Eval(dl.ctx, script, []string{lockKey}, lockValue, int(ttl.Seconds()))

	if result.Err() != nil {
		log.Printf("Error extending lock %s: %v", lockName, result.Err())
		return false
	}

	if result.Val().(int64) == 1 {
		log.Printf("Lock extended: %s by instance %s", lockName, instanceID)
		return true
	}

	return false
}

// GetActiveLocks returns all active locks
func (dl *DistributedLock) GetActiveLocks() ([]*LockInfo, error) {
	pattern := "lock:*"
	result := dl.client.Keys(dl.ctx, pattern)

	if result.Err() != nil {
		return nil, result.Err()
	}

	var locks []*LockInfo
	for _, key := range result.Val() {
		lockName := strings.TrimPrefix(key, "lock:")
		lockInfo, err := dl.GetLockInfo(lockName)
		if err != nil {
			log.Printf("Error getting lock info for %s: %v", lockName, err)
			continue
		}
		if lockInfo != nil {
			locks = append(locks, lockInfo)
		}
	}

	return locks, nil
}

// CleanupExpiredLocks removes expired locks (safety cleanup)
func (dl *DistributedLock) CleanupExpiredLocks() {
	pattern := "lock:*"
	result := dl.client.Keys(dl.ctx, pattern)

	if result.Err() != nil {
		log.Printf("Error getting lock keys: %v", result.Err())
		return
	}

	for _, key := range result.Val() {
		// Check TTL
		ttlResult := dl.client.TTL(dl.ctx, key)
		if ttlResult.Err() != nil {
			continue
		}

		// If TTL is -1 (no expiration) or -2 (key doesn't exist), skip
		if ttlResult.Val() <= 0 {
			dl.client.Del(dl.ctx, key)
			log.Printf("Cleaned up expired lock: %s", key)
		}
	}
}

// Stop gracefully stops the distributed lock service
func (dl *DistributedLock) Stop() {
	if dl.client != nil {
		dl.client.Close()
	}
	log.Println("Distributed lock service stopped")
}
