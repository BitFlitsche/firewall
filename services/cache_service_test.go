package services

import (
	"fmt"
	"testing"
	"time"
)

func TestCache_SetAndGet(t *testing.T) {
	cache := GetCache()
	defer cache.Clear()

	// Test basic set and get
	err := cache.Set("test-key", "test-value", 1*time.Hour)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	value, exists, err := cache.Get("test-key")
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if !exists {
		t.Error("Expected key to exist")
	}
	if value != "test-value" {
		t.Errorf("Expected 'test-value', got %v", value)
	}

	// Test non-existent key
	value, exists, err = cache.Get("non-existent")
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if exists {
		t.Error("Expected key to not exist")
	}
	if value != nil {
		t.Errorf("Expected nil value, got %v", value)
	}
}

func TestCache_Expiration(t *testing.T) {
	cache := GetCache()
	defer cache.Clear()

	// Set item with short TTL
	err := cache.Set("expire-key", "expire-value", 10*time.Millisecond)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	// Should exist immediately
	value, exists, err := cache.Get("expire-key")
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if !exists {
		t.Error("Expected key to exist")
	}
	if value != "expire-value" {
		t.Errorf("Expected 'expire-value', got %v", value)
	}

	// Wait for expiration
	time.Sleep(20 * time.Millisecond)

	// Should not exist after expiration
	value, exists, err = cache.Get("expire-key")
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if exists {
		t.Error("Expected key to not exist after expiration")
	}
	if value != nil {
		t.Errorf("Expected nil value after expiration, got %v", value)
	}
}

func TestCache_Delete(t *testing.T) {
	cache := GetCache()
	defer cache.Clear()

	// Set item
	err := cache.Set("delete-key", "delete-value", 1*time.Hour)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	// Verify it exists
	value, exists, err := cache.Get("delete-key")
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if !exists {
		t.Error("Expected key to exist")
	}

	// Delete item
	err = cache.Delete("delete-key")
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	// Verify it's gone
	value, exists, err = cache.Get("delete-key")
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if exists {
		t.Error("Expected key to not exist after deletion")
	}
	if value != nil {
		t.Errorf("Expected nil value after deletion, got %v", value)
	}

	// Delete non-existent key should not error
	err = cache.Delete("non-existent")
	if err != nil {
		t.Errorf("Delete of non-existent key should not error: %v", err)
	}
}

func TestCache_InvalidateByType(t *testing.T) {
	cache := GetCache()
	defer cache.Clear()

	// Set items with different types
	testData := map[string]string{
		"ip:192.168.1.1":       "ip-value-1",
		"ip:192.168.1.2":       "ip-value-2",
		"email:test@test.com":  "email-value-1",
		"email:admin@test.com": "email-value-2",
		"user:admin":           "user-value-1",
	}

	for key, value := range testData {
		err := cache.Set(key, value, 1*time.Hour)
		if err != nil {
			t.Errorf("Set failed for %s: %v", key, err)
		}
	}

	// Verify all items exist
	for key, expectedValue := range testData {
		value, exists, err := cache.Get(key)
		if err != nil {
			t.Errorf("Get failed for %s: %v", key, err)
		}
		if !exists {
			t.Errorf("Expected key %s to exist", key)
		}
		if value != expectedValue {
			t.Errorf("Expected %s for key %s, got %v", expectedValue, key, value)
		}
	}

	// Invalidate IP type
	err := cache.InvalidateByType("ip")
	if err != nil {
		t.Errorf("InvalidateByType failed: %v", err)
	}

	// Verify IP items are gone
	for key := range testData {
		if key[:2] == "ip" {
			value, exists, err := cache.Get(key)
			if err != nil {
				t.Errorf("Get failed for %s: %v", key, err)
			}
			if exists {
				t.Errorf("Expected key %s to not exist after invalidation", key)
			}
			if value != nil {
				t.Errorf("Expected nil value for %s after invalidation, got %v", key, value)
			}
		}
	}

	// Verify other items still exist
	for key, expectedValue := range testData {
		if key[:2] != "ip" {
			value, exists, err := cache.Get(key)
			if err != nil {
				t.Errorf("Get failed for %s: %v", key, err)
			}
			if !exists {
				t.Errorf("Expected key %s to still exist", key)
			}
			if value != expectedValue {
				t.Errorf("Expected %s for key %s, got %v", expectedValue, key, value)
			}
		}
	}
}

func TestCache_InvalidatePattern(t *testing.T) {
	cache := GetCache()
	defer cache.Clear()

	// Set items with different patterns
	testData := map[string]string{
		"ip:192.168.1.1":      "ip-value-1",
		"ip:192.168.1.2":      "ip-value-2",
		"email:test@test.com": "email-value-1",
		"user:admin":          "user-value-1",
		"stats:ip":            "stats-value-1",
		"stats:email":         "stats-value-2",
	}

	for key, value := range testData {
		err := cache.Set(key, value, 1*time.Hour)
		if err != nil {
			t.Errorf("Set failed for %s: %v", key, err)
		}
	}

	// Invalidate items containing "stats"
	err := cache.InvalidatePattern("stats")
	if err != nil {
		t.Errorf("InvalidatePattern failed: %v", err)
	}

	// Verify stats items are gone
	for key := range testData {
		if key[:5] == "stats" {
			value, exists, err := cache.Get(key)
			if err != nil {
				t.Errorf("Get failed for %s: %v", key, err)
			}
			if exists {
				t.Errorf("Expected key %s to not exist after pattern invalidation", key)
			}
			if value != nil {
				t.Errorf("Expected nil value for %s after pattern invalidation, got %v", key, value)
			}
		}
	}

	// Verify other items still exist
	for key, expectedValue := range testData {
		if key[:5] != "stats" {
			value, exists, err := cache.Get(key)
			if err != nil {
				t.Errorf("Get failed for %s: %v", key, err)
			}
			if !exists {
				t.Errorf("Expected key %s to still exist", key)
			}
			if value != expectedValue {
				t.Errorf("Expected %s for key %s, got %v", expectedValue, key, value)
			}
		}
	}
}

func TestCache_Clear(t *testing.T) {
	cache := GetCache()
	defer cache.Clear()

	// Set multiple items
	testData := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	for key, value := range testData {
		err := cache.Set(key, value, 1*time.Hour)
		if err != nil {
			t.Errorf("Set failed for %s: %v", key, err)
		}
	}

	// Verify all items exist
	for key, expectedValue := range testData {
		value, exists, err := cache.Get(key)
		if err != nil {
			t.Errorf("Get failed for %s: %v", key, err)
		}
		if !exists {
			t.Errorf("Expected key %s to exist", key)
		}
		if value != expectedValue {
			t.Errorf("Expected %s for key %s, got %v", expectedValue, key, value)
		}
	}

	// Clear cache
	err := cache.Clear()
	if err != nil {
		t.Errorf("Clear failed: %v", err)
	}

	// Verify all items are gone
	for key := range testData {
		value, exists, err := cache.Get(key)
		if err != nil {
			t.Errorf("Get failed for %s: %v", key, err)
		}
		if exists {
			t.Errorf("Expected key %s to not exist after clear", key)
		}
		if value != nil {
			t.Errorf("Expected nil value for %s after clear, got %v", key, value)
		}
	}
}

func TestCache_Stats(t *testing.T) {
	cache := GetCache()
	defer cache.Clear()

	// Get initial stats
	stats, err := cache.Stats()
	if err != nil {
		t.Errorf("Stats failed: %v", err)
	}

	initialItems := stats["items"].(int)
	if initialItems != 0 {
		t.Errorf("Expected 0 items initially, got %d", initialItems)
	}

	// Add some items
	testData := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	for key, value := range testData {
		err := cache.Set(key, value, 1*time.Hour)
		if err != nil {
			t.Errorf("Set failed for %s: %v", key, err)
		}
	}

	// Get stats after adding items
	stats, err = cache.Stats()
	if err != nil {
		t.Errorf("Stats failed: %v", err)
	}

	items := stats["items"].(int)
	validItems := stats["valid_items"].(int)
	expiredItems := stats["expired_items"].(int)
	cacheType := stats["type"].(string)

	if items != 3 {
		t.Errorf("Expected 3 items, got %d", items)
	}
	if validItems != 3 {
		t.Errorf("Expected 3 valid items, got %d", validItems)
	}
	if expiredItems != 0 {
		t.Errorf("Expected 0 expired items, got %d", expiredItems)
	}
	if cacheType != "in-memory" {
		t.Errorf("Expected 'in-memory' type, got %s", cacheType)
	}

	// Add expired item
	err = cache.Set("expired-key", "expired-value", 1*time.Millisecond)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	// Wait for expiration
	time.Sleep(10 * time.Millisecond)

	// Get stats after adding expired item
	stats, err = cache.Stats()
	if err != nil {
		t.Errorf("Stats failed: %v", err)
	}

	items = stats["items"].(int)
	validItems = stats["valid_items"].(int)
	expiredItems = stats["expired_items"].(int)

	// Note: The cleanup might have already removed the expired item
	if items < 3 || items > 4 {
		t.Errorf("Expected 3-4 items, got %d", items)
	}
}

func TestCache_KeyGenerators(t *testing.T) {
	cache := GetCache()

	// Test FilterKey
	filterKey := cache.FilterKey("ip", "192.168.1.1")
	expectedFilterKey := "ip:filter:192.168.1.1"
	if filterKey != expectedFilterKey {
		t.Errorf("Expected FilterKey to return %s, got %s", expectedFilterKey, filterKey)
	}

	// Test ListKey
	listKey := cache.ListKey("ip", "1", "10", "search", "active")
	expectedListKey := "ip:list:1:10:search:active"
	if listKey != expectedListKey {
		t.Errorf("Expected ListKey to return %s, got %s", expectedListKey, listKey)
	}

	// Test StatsKey
	statsKey := cache.StatsKey("ip")
	expectedStatsKey := "ip:stats"
	if statsKey != expectedStatsKey {
		t.Errorf("Expected StatsKey to return %s, got %s", expectedStatsKey, statsKey)
	}
}

func TestCache_InvalidationHelpers(t *testing.T) {
	cache := GetCache()
	defer cache.Clear()

	// Set up test data
	testData := map[string]string{
		"ip:filter:192.168.1.1":      "filter-value",
		"ip:list:1:10:search:active": "list-value",
		"ip:stats":                   "stats-value",
		"email:filter:test@test.com": "email-filter-value",
	}

	for key, value := range testData {
		err := cache.Set(key, value, 1*time.Hour)
		if err != nil {
			t.Errorf("Set failed for %s: %v", key, err)
		}
	}

	// Test InvalidateFilter
	err := cache.InvalidateFilter("ip")
	if err != nil {
		t.Errorf("InvalidateFilter failed: %v", err)
	}

	// Verify filter items are gone
	_, exists, err := cache.Get("ip:filter:192.168.1.1")
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if exists {
		t.Error("Expected filter key to not exist after invalidation")
	}

	// Test InvalidateList
	err = cache.InvalidateList("ip")
	if err != nil {
		t.Errorf("InvalidateList failed: %v", err)
	}

	// Verify list items are gone
	_, exists, err = cache.Get("ip:list:1:10:search:active")
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if exists {
		t.Error("Expected list key to not exist after invalidation")
	}

	// Test InvalidateStats
	err = cache.InvalidateStats("ip")
	if err != nil {
		t.Errorf("InvalidateStats failed: %v", err)
	}

	// Verify stats items are gone
	_, exists, err = cache.Get("ip:stats")
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if exists {
		t.Error("Expected stats key to not exist after invalidation")
	}

	// Test InvalidateAll
	err = cache.InvalidateAll("ip")
	if err != nil {
		t.Errorf("InvalidateAll failed: %v", err)
	}

	// Verify all IP items are gone
	for key := range testData {
		if key[:2] == "ip" {
			_, exists, err := cache.Get(key)
			if err != nil {
				t.Errorf("Get failed for %s: %v", key, err)
			}
			if exists {
				t.Errorf("Expected key %s to not exist after InvalidateAll", key)
			}
		}
	}

	// Verify non-IP items still exist
	for key, expectedValue := range testData {
		if key[:5] == "email" {
			value, exists, err := cache.Get(key)
			if err != nil {
				t.Errorf("Get failed for %s: %v", key, err)
			}
			if !exists {
				t.Errorf("Expected key %s to still exist", key)
			}
			if value != expectedValue {
				t.Errorf("Expected %s for key %s, got %v", expectedValue, key, value)
			}
		}
	}
}

func TestCache_Concurrency(t *testing.T) {
	cache := GetCache()
	defer cache.Clear()

	// Test concurrent reads and writes
	done := make(chan bool)
	numGoroutines := 10
	numOperations := 100

	// Start multiple goroutines doing concurrent operations
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("key-%d-%d", id, j)
				value := fmt.Sprintf("value-%d-%d", id, j)

				// Set
				err := cache.Set(key, value, 1*time.Hour)
				if err != nil {
					t.Errorf("Set failed in goroutine %d: %v", id, err)
				}

				// Get
				retrievedValue, exists, err := cache.Get(key)
				if err != nil {
					t.Errorf("Get failed in goroutine %d: %v", id, err)
				}
				if !exists {
					t.Errorf("Expected key %s to exist in goroutine %d", key, id)
				}
				if retrievedValue != value {
					t.Errorf("Expected %s for key %s in goroutine %d, got %v", value, key, id, retrievedValue)
				}
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Verify all items were set correctly
	for i := 0; i < numGoroutines; i++ {
		for j := 0; j < numOperations; j++ {
			key := fmt.Sprintf("key-%d-%d", i, j)
			expectedValue := fmt.Sprintf("value-%d-%d", i, j)

			value, exists, err := cache.Get(key)
			if err != nil {
				t.Errorf("Get failed for %s: %v", key, err)
			}
			if !exists {
				t.Errorf("Expected key %s to exist", key)
			}
			if value != expectedValue {
				t.Errorf("Expected %s for key %s, got %v", expectedValue, key, value)
			}
		}
	}
}

func TestCache_Stop(t *testing.T) {
	cache := GetCache()
	defer cache.Clear()

	// Set some items
	err := cache.Set("test-key", "test-value", 1*time.Hour)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	// Stop the cache
	cache.Stop()

	// Verify we can still get items (Stop only stops cleanup goroutine)
	value, exists, err := cache.Get("test-key")
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if !exists {
		t.Error("Expected key to exist after Stop")
	}
	if value != "test-value" {
		t.Errorf("Expected 'test-value', got %v", value)
	}
}
