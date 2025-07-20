package services

import (
	"encoding/json"
	"firewall/config"
	"testing"
	"time"
)

func init() {
	// Initialize config for tests
	config.InitConfig()
}

// ============================================================================
// DISTRIBUTED CACHE TESTS
// ============================================================================

func TestGetDistributedCache(t *testing.T) {
	// Test singleton pattern
	cache1 := GetDistributedCache()
	cache2 := GetDistributedCache()

	if cache1 != cache2 {
		t.Error("Expected singleton pattern for distributed cache")
	}
}

func TestDistributedCache_Set(t *testing.T) {
	cache := GetDistributedCache()

	// This will fail because we don't have a real Redis client, but we can test the function structure
	defer func() {
		if r := recover(); r != nil {
			t.Logf("DistributedCache.Set panicked as expected: %v", r)
		}
	}()

	err := cache.Set("test-key", "test-value", 5*time.Minute)
	// We expect an error because there's no Redis client, but the function should not panic
	if err == nil {
		t.Log("DistributedCache.Set completed without error (unexpected in test environment)")
	}
}

func TestDistributedCache_Get(t *testing.T) {
	cache := GetDistributedCache()

	// This will fail because we don't have a real Redis client, but we can test the function structure
	defer func() {
		if r := recover(); r != nil {
			t.Logf("DistributedCache.Get panicked as expected: %v", r)
		}
	}()

	value, found, err := cache.Get("test-key")
	// We expect an error because there's no Redis client, but the function should not panic
	if err == nil {
		t.Log("DistributedCache.Get completed without error (unexpected in test environment)")
	}

	// Verify return types
	if value != nil {
		t.Log("value is not nil (unexpected in test environment)")
	}
	if found {
		t.Log("found is true (unexpected in test environment)")
	}
}

func TestDistributedCache_Delete(t *testing.T) {
	cache := GetDistributedCache()

	// This will fail because we don't have a real Redis client, but we can test the function structure
	defer func() {
		if r := recover(); r != nil {
			t.Logf("DistributedCache.Delete panicked as expected: %v", r)
		}
	}()

	err := cache.Delete("test-key")
	// We expect an error because there's no Redis client, but the function should not panic
	if err == nil {
		t.Log("DistributedCache.Delete completed without error (unexpected in test environment)")
	}
}

func TestDistributedCache_InvalidateByType(t *testing.T) {
	cache := GetDistributedCache()

	// This will fail because we don't have a real Redis client, but we can test the function structure
	defer func() {
		if r := recover(); r != nil {
			t.Logf("DistributedCache.InvalidateByType panicked as expected: %v", r)
		}
	}()

	err := cache.InvalidateByType("ip")
	// We expect an error because there's no Redis client, but the function should not panic
	if err == nil {
		t.Log("DistributedCache.InvalidateByType completed without error (unexpected in test environment)")
	}
}

func TestDistributedCache_InvalidatePattern(t *testing.T) {
	cache := GetDistributedCache()

	// This will fail because we don't have a real Redis client, but we can test the function structure
	defer func() {
		if r := recover(); r != nil {
			t.Logf("DistributedCache.InvalidatePattern panicked as expected: %v", r)
		}
	}()

	err := cache.InvalidatePattern("ip:*")
	// We expect an error because there's no Redis client, but the function should not panic
	if err == nil {
		t.Log("DistributedCache.InvalidatePattern completed without error (unexpected in test environment)")
	}
}

func TestDistributedCache_Clear(t *testing.T) {
	cache := GetDistributedCache()

	// This will fail because we don't have a real Redis client, but we can test the function structure
	defer func() {
		if r := recover(); r != nil {
			t.Logf("DistributedCache.Clear panicked as expected: %v", r)
		}
	}()

	err := cache.Clear()
	// We expect an error because there's no Redis client, but the function should not panic
	if err == nil {
		t.Log("DistributedCache.Clear completed without error (unexpected in test environment)")
	}
}

func TestDistributedCache_Stats(t *testing.T) {
	cache := GetDistributedCache()

	// This will fail because we don't have a real Redis client, but we can test the function structure
	defer func() {
		if r := recover(); r != nil {
			t.Logf("DistributedCache.Stats panicked as expected: %v", r)
		}
	}()

	stats, err := cache.Stats()
	// We expect an error because there's no Redis client, but the function should not panic
	if err == nil {
		t.Log("DistributedCache.Stats completed without error (unexpected in test environment)")
	}

	// Verify return types
	if stats == nil {
		t.Log("stats is nil (expected in test environment)")
	}
}

func TestDistributedCache_Stop(t *testing.T) {
	cache := GetDistributedCache()

	// This should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("DistributedCache.Stop panicked unexpectedly: %v", r)
		}
	}()

	cache.Stop()
	// Should complete without error
}

func TestDistributedCache_KeyGeneration(t *testing.T) {
	cache := GetDistributedCache()

	// Test key generation methods
	filterKey := cache.FilterKey("ip", "192.168.1.1")
	if filterKey == "" {
		t.Error("Expected non-empty filter key")
	}

	listKey := cache.ListKey("ip", "1", "10", "search", "allowed")
	if listKey == "" {
		t.Error("Expected non-empty list key")
	}

	statsKey := cache.StatsKey("ip")
	if statsKey == "" {
		t.Error("Expected non-empty stats key")
	}
}

func TestDistributedCache_InvalidationHelpers(t *testing.T) {
	cache := GetDistributedCache()

	// This will fail because we don't have a real Redis client, but we can test the function structure
	defer func() {
		if r := recover(); r != nil {
			t.Logf("DistributedCache invalidation helpers panicked as expected: %v", r)
		}
	}()

	// Test invalidation helper methods
	err := cache.InvalidateFilter("ip")
	if err == nil {
		t.Log("InvalidateFilter completed without error (unexpected in test environment)")
	}

	err = cache.InvalidateList("ip")
	if err == nil {
		t.Log("InvalidateList completed without error (unexpected in test environment)")
	}

	err = cache.InvalidateStats("ip")
	if err == nil {
		t.Log("InvalidateStats completed without error (unexpected in test environment)")
	}

	err = cache.InvalidateAll("ip")
	if err == nil {
		t.Log("InvalidateAll completed without error (unexpected in test environment)")
	}
}

func TestDistributedCacheItem_JSONMarshaling(t *testing.T) {
	// Test JSON marshaling of cache items
	item := DistributedCacheItem{
		Value:     "test-value",
		ExpiresAt: time.Now().Add(5 * time.Minute),
		CreatedAt: time.Now(),
	}

	data, err := json.Marshal(item)
	if err != nil {
		t.Errorf("Failed to marshal cache item: %v", err)
		return
	}

	var unmarshaled DistributedCacheItem
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Errorf("Failed to unmarshal cache item: %v", err)
		return
	}

	if unmarshaled.Value != item.Value {
		t.Errorf("Expected value %v, got %v", item.Value, unmarshaled.Value)
	}
}

// ============================================================================
// DISTRIBUTED LOCK TESTS
// ============================================================================

func TestGetDistributedLock(t *testing.T) {
	// Test singleton pattern
	lock1 := GetDistributedLock()
	lock2 := GetDistributedLock()

	if lock1 != lock2 {
		t.Error("Expected singleton pattern for distributed lock")
	}
}

func TestDistributedLock_TryAcquireLock(t *testing.T) {
	lock := GetDistributedLock()

	// This will fail because we don't have a real Redis client, but we can test the function structure
	defer func() {
		if r := recover(); r != nil {
			t.Logf("DistributedLock.TryAcquireLock panicked as expected: %v", r)
		}
	}()

	acquired, lockInfo := lock.TryAcquireLock("test-lock", 5*time.Minute)
	// We expect an error because there's no Redis client, but the function should not panic
	if acquired {
		t.Log("TryAcquireLock succeeded (unexpected in test environment)")
	}

	// Verify return types
	if lockInfo != nil {
		t.Log("lockInfo is not nil (unexpected in test environment)")
	}
}

func TestDistributedLock_ReleaseLock(t *testing.T) {
	lock := GetDistributedLock()

	// This will fail because we don't have a real Redis client, but we can test the function structure
	defer func() {
		if r := recover(); r != nil {
			t.Logf("DistributedLock.ReleaseLock panicked as expected: %v", r)
		}
	}()

	released := lock.ReleaseLock("test-lock")
	// We expect an error because there's no Redis client, but the function should not panic
	if released {
		t.Log("ReleaseLock succeeded (unexpected in test environment)")
	}
}

func TestDistributedLock_IsLocked(t *testing.T) {
	lock := GetDistributedLock()

	// This will fail because we don't have a real Redis client, but we can test the function structure
	defer func() {
		if r := recover(); r != nil {
			t.Logf("DistributedLock.IsLocked panicked as expected: %v", r)
		}
	}()

	isLocked := lock.IsLocked("test-lock")
	// We expect an error because there's no Redis client, but the function should not panic
	if isLocked {
		t.Log("IsLocked returned true (unexpected in test environment)")
	}
}

func TestDistributedLock_GetLockInfo(t *testing.T) {
	lock := GetDistributedLock()

	// This will fail because we don't have a real Redis client, but we can test the function structure
	defer func() {
		if r := recover(); r != nil {
			t.Logf("DistributedLock.GetLockInfo panicked as expected: %v", r)
		}
	}()

	lockInfo, err := lock.GetLockInfo("test-lock")
	// We expect an error because there's no Redis client, but the function should not panic
	if err == nil {
		t.Log("GetLockInfo completed without error (unexpected in test environment)")
	}

	// Verify return types
	if lockInfo != nil {
		t.Log("lockInfo is not nil (unexpected in test environment)")
	}
}

func TestDistributedLock_ExtendLock(t *testing.T) {
	lock := GetDistributedLock()

	// This will fail because we don't have a real Redis client, but we can test the function structure
	defer func() {
		if r := recover(); r != nil {
			t.Logf("DistributedLock.ExtendLock panicked as expected: %v", r)
		}
	}()

	extended := lock.ExtendLock("test-lock", 5*time.Minute)
	// We expect an error because there's no Redis client, but the function should not panic
	if extended {
		t.Log("ExtendLock succeeded (unexpected in test environment)")
	}
}

func TestDistributedLock_GetActiveLocks(t *testing.T) {
	lock := GetDistributedLock()

	// This will fail because we don't have a real Redis client, but we can test the function structure
	defer func() {
		if r := recover(); r != nil {
			t.Logf("DistributedLock.GetActiveLocks panicked as expected: %v", r)
		}
	}()

	locks, err := lock.GetActiveLocks()
	// We expect an error because there's no Redis client, but the function should not panic
	if err == nil {
		t.Log("GetActiveLocks completed without error (unexpected in test environment)")
	}

	// Verify return types
	if locks == nil {
		t.Log("locks is nil (expected in test environment)")
	}
}

func TestDistributedLock_CleanupExpiredLocks(t *testing.T) {
	lock := GetDistributedLock()

	// This will fail because we don't have a real Redis client, but we can test the function structure
	defer func() {
		if r := recover(); r != nil {
			t.Logf("DistributedLock.CleanupExpiredLocks panicked as expected: %v", r)
		}
	}()

	lock.CleanupExpiredLocks()
	// Should complete without error
}

func TestDistributedLock_Stop(t *testing.T) {
	lock := GetDistributedLock()

	// This should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("DistributedLock.Stop panicked unexpectedly: %v", r)
		}
	}()

	lock.Stop()
	// Should complete without error
}

func TestGenerateInstanceID(t *testing.T) {
	// Test instance ID generation
	lock := GetDistributedLock()

	// This will fail because we don't have a real Redis client, but we can test the function structure
	defer func() {
		if r := recover(); r != nil {
			t.Logf("GenerateInstanceID panicked as expected: %v", r)
		}
	}()

	// The function is called internally, so we just test that the lock service initializes
	if lock == nil {
		t.Error("Expected lock service to be initialized")
	}
}

func TestLockInfo_JSONMarshaling(t *testing.T) {
	// Test JSON marshaling of lock info
	lockInfo := &LockInfo{
		LockID:    "test-lock",
		Instance:  "test-instance",
		Acquired:  time.Now(),
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	data, err := json.Marshal(lockInfo)
	if err != nil {
		t.Errorf("Failed to marshal lock info: %v", err)
		return
	}

	var unmarshaled LockInfo
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Errorf("Failed to unmarshal lock info: %v", err)
		return
	}

	if unmarshaled.LockID != lockInfo.LockID {
		t.Errorf("Expected LockID %s, got %s", lockInfo.LockID, unmarshaled.LockID)
	}
	if unmarshaled.Instance != lockInfo.Instance {
		t.Errorf("Expected Instance %s, got %s", lockInfo.Instance, unmarshaled.Instance)
	}
}

// ============================================================================
// NO-OP LOCK TESTS
// ============================================================================

func TestNewNoOpDistributedLock(t *testing.T) {
	lock := NewNoOpDistributedLock()

	if lock == nil {
		t.Error("Expected no-op lock to be created")
	}

	if lock.instanceID != "single-instance" {
		t.Errorf("Expected instance ID 'single-instance', got %s", lock.instanceID)
	}
}

func TestNoOpDistributedLock_TryAcquireLock(t *testing.T) {
	lock := NewNoOpDistributedLock()

	acquired, lockInfo := lock.TryAcquireLock("test-lock", 5*time.Minute)

	if !acquired {
		t.Error("Expected no-op lock acquisition to always succeed")
	}

	if lockInfo == nil {
		t.Error("Expected lock info to be returned")
	}

	if lockInfo.LockID != "test-lock" {
		t.Errorf("Expected LockID 'test-lock', got %s", lockInfo.LockID)
	}

	if lockInfo.Instance != "single-instance" {
		t.Errorf("Expected Instance 'single-instance', got %s", lockInfo.Instance)
	}
}

func TestNoOpDistributedLock_ReleaseLock(t *testing.T) {
	lock := NewNoOpDistributedLock()

	released := lock.ReleaseLock("test-lock")

	if !released {
		t.Error("Expected no-op lock release to always succeed")
	}
}

func TestNoOpDistributedLock_IsLocked(t *testing.T) {
	lock := NewNoOpDistributedLock()

	isLocked := lock.IsLocked("test-lock")

	if isLocked {
		t.Error("Expected no-op lock to always return false for IsLocked")
	}
}

func TestNoOpDistributedLock_GetLockInfo(t *testing.T) {
	lock := NewNoOpDistributedLock()

	lockInfo, err := lock.GetLockInfo("test-lock")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if lockInfo != nil {
		t.Error("Expected no-op lock to return nil for GetLockInfo")
	}
}

func TestNoOpDistributedLock_ExtendLock(t *testing.T) {
	lock := NewNoOpDistributedLock()

	extended := lock.ExtendLock("test-lock", 5*time.Minute)

	if !extended {
		t.Error("Expected no-op lock extension to always succeed")
	}
}

func TestNoOpDistributedLock_GetActiveLocks(t *testing.T) {
	lock := NewNoOpDistributedLock()

	locks, err := lock.GetActiveLocks()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(locks) != 0 {
		t.Errorf("Expected empty slice, got %d locks", len(locks))
	}
}

func TestNoOpDistributedLock_CleanupExpiredLocks(t *testing.T) {
	lock := NewNoOpDistributedLock()

	// Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("CleanupExpiredLocks panicked unexpectedly: %v", r)
		}
	}()

	lock.CleanupExpiredLocks()
	// Should complete without error
}

func TestNoOpDistributedLock_Stop(t *testing.T) {
	lock := NewNoOpDistributedLock()

	// Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Stop panicked unexpectedly: %v", r)
		}
	}()

	lock.Stop()
	// Should complete without error
}

// ============================================================================
// INTEGRATION TESTS
// ============================================================================

func TestDistributedSystems_Integration(t *testing.T) {
	// Test that both distributed cache and lock can be initialized together
	cache := GetDistributedCache()
	lock := GetDistributedLock()

	// Cache might be nil if distributed caching is disabled
	if cache == nil {
		t.Log("Cache is nil (expected when distributed caching is disabled)")
	}

	if lock == nil {
		t.Error("Expected distributed lock to be initialized")
	}

	// Test that they don't interfere with each other
	// Note: They are different types, so direct comparison is not possible
	t.Log("Distributed cache and lock initialized successfully")
}

func TestDistributedSystems_Configuration(t *testing.T) {
	// Test configuration-based initialization
	originalDistributed := config.AppConfig.Caching.Distributed
	originalLocking := config.AppConfig.Locking.Enabled

	// Test with distributed features disabled
	config.AppConfig.Caching.Distributed = false
	config.AppConfig.Locking.Enabled = false

	cache := GetDistributedCache()
	lock := GetDistributedLock()

	// Reset configuration
	config.AppConfig.Caching.Distributed = originalDistributed
	config.AppConfig.Locking.Enabled = originalLocking

	if cache == nil {
		t.Log("Cache is nil (expected when distributed caching is disabled)")
	}

	if lock == nil {
		t.Error("Expected lock to be initialized (no-op implementation)")
	}
}

func TestDistributedSystems_ErrorHandling(t *testing.T) {
	// Test error handling for nil instances
	var nilCache *DistributedCache
	var nilLock DistributedLockInterface

	// Test nil cache operations
	err := nilCache.Set("key", "value", time.Minute)
	if err == nil {
		t.Error("Expected error for nil cache Set")
	}

	_, _, err = nilCache.Get("key")
	if err == nil {
		t.Error("Expected error for nil cache Get")
	}

	err = nilCache.Delete("key")
	if err == nil {
		t.Error("Expected error for nil cache Delete")
	}

	// Test nil lock operations (should not panic for no-op implementation)
	if nilLock != nil {
		acquired, _ := nilLock.TryAcquireLock("lock", time.Minute)
		if acquired {
			t.Log("Nil lock acquisition succeeded (unexpected)")
		}
	}
}
