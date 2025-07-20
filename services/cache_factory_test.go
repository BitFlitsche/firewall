package services

import (
	"firewall/config"
	"testing"
	"time"
)

func init() {
	// Initialize config for tests
	config.InitConfig()
}

func TestCacheFactory_Singleton(t *testing.T) {
	// Test that GetCacheFactory returns the same instance
	factory1 := GetCacheFactory()
	factory2 := GetCacheFactory()

	if factory1 != factory2 {
		t.Error("GetCacheFactory should return the same instance")
	}
}

func TestCacheFactory_SetAndGet(t *testing.T) {
	factory := GetCacheFactory()
	defer factory.Clear()

	// Test basic set and get
	err := factory.Set("test-key", "test-value", 1*time.Hour)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	value, exists, err := factory.Get("test-key")
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
	value, exists, err = factory.Get("non-existent")
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

func TestCacheFactory_Delete(t *testing.T) {
	factory := GetCacheFactory()
	defer factory.Clear()

	// Set item
	err := factory.Set("delete-key", "delete-value", 1*time.Hour)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	// Verify it exists
	value, exists, err := factory.Get("delete-key")
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if !exists {
		t.Error("Expected key to exist")
	}

	// Delete item
	err = factory.Delete("delete-key")
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	// Verify it's gone
	value, exists, err = factory.Get("delete-key")
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
	err = factory.Delete("non-existent")
	if err != nil {
		t.Errorf("Delete of non-existent key should not error: %v", err)
	}
}

func TestCacheFactory_InvalidateByType(t *testing.T) {
	factory := GetCacheFactory()
	defer factory.Clear()

	// Set items with different types
	testData := map[string]string{
		"ip:192.168.1.1":       "ip-value-1",
		"ip:192.168.1.2":       "ip-value-2",
		"email:test@test.com":  "email-value-1",
		"email:admin@test.com": "email-value-2",
		"user:admin":           "user-value-1",
	}

	for key, value := range testData {
		err := factory.Set(key, value, 1*time.Hour)
		if err != nil {
			t.Errorf("Set failed for %s: %v", key, err)
		}
	}

	// Verify all items exist
	for key, expectedValue := range testData {
		value, exists, err := factory.Get(key)
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
	err := factory.InvalidateByType("ip")
	if err != nil {
		t.Errorf("InvalidateByType failed: %v", err)
	}

	// Verify IP items are gone
	for key := range testData {
		if key[:2] == "ip" {
			value, exists, err := factory.Get(key)
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
			value, exists, err := factory.Get(key)
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

func TestCacheFactory_InvalidatePattern(t *testing.T) {
	factory := GetCacheFactory()
	defer factory.Clear()

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
		err := factory.Set(key, value, 1*time.Hour)
		if err != nil {
			t.Errorf("Set failed for %s: %v", key, err)
		}
	}

	// Invalidate items containing "stats"
	err := factory.InvalidatePattern("stats")
	if err != nil {
		t.Errorf("InvalidatePattern failed: %v", err)
	}

	// Verify stats items are gone
	for key := range testData {
		if key[:5] == "stats" {
			value, exists, err := factory.Get(key)
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
			value, exists, err := factory.Get(key)
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

func TestCacheFactory_Clear(t *testing.T) {
	factory := GetCacheFactory()
	defer factory.Clear()

	// Set multiple items
	testData := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	for key, value := range testData {
		err := factory.Set(key, value, 1*time.Hour)
		if err != nil {
			t.Errorf("Set failed for %s: %v", key, err)
		}
	}

	// Verify all items exist
	for key, expectedValue := range testData {
		value, exists, err := factory.Get(key)
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
	err := factory.Clear()
	if err != nil {
		t.Errorf("Clear failed: %v", err)
	}

	// Verify all items are gone
	for key := range testData {
		value, exists, err := factory.Get(key)
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

func TestCacheFactory_Stats(t *testing.T) {
	factory := GetCacheFactory()
	defer factory.Clear()

	// Get initial stats
	stats, err := factory.Stats()
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
		err := factory.Set(key, value, 1*time.Hour)
		if err != nil {
			t.Errorf("Set failed for %s: %v", key, err)
		}
	}

	// Get stats after adding items
	stats, err = factory.Stats()
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
}

func TestCacheFactory_KeyGenerators(t *testing.T) {
	factory := GetCacheFactory()

	// Test FilterKey
	filterKey := factory.FilterKey("ip", "192.168.1.1")
	expectedFilterKey := "ip:filter:192.168.1.1"
	if filterKey != expectedFilterKey {
		t.Errorf("Expected FilterKey to return %s, got %s", expectedFilterKey, filterKey)
	}

	// Test ListKey
	listKey := factory.ListKey("ip", "1", "10", "search", "active")
	expectedListKey := "ip:list:1:10:search:active"
	if listKey != expectedListKey {
		t.Errorf("Expected ListKey to return %s, got %s", expectedListKey, listKey)
	}

	// Test StatsKey
	statsKey := factory.StatsKey("ip")
	expectedStatsKey := "ip:stats"
	if statsKey != expectedStatsKey {
		t.Errorf("Expected StatsKey to return %s, got %s", expectedStatsKey, statsKey)
	}
}

func TestCacheFactory_InvalidationHelpers(t *testing.T) {
	factory := GetCacheFactory()
	defer factory.Clear()

	// Set up test data
	testData := map[string]string{
		"ip:filter:192.168.1.1":      "filter-value",
		"ip:list:1:10:search:active": "list-value",
		"ip:stats":                   "stats-value",
		"email:filter:test@test.com": "email-filter-value",
	}

	for key, value := range testData {
		err := factory.Set(key, value, 1*time.Hour)
		if err != nil {
			t.Errorf("Set failed for %s: %v", key, err)
		}
	}

	// Test InvalidateFilter
	err := factory.InvalidateFilter("ip")
	if err != nil {
		t.Errorf("InvalidateFilter failed: %v", err)
	}

	// Verify filter items are gone
	_, exists, err := factory.Get("ip:filter:192.168.1.1")
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if exists {
		t.Error("Expected filter key to not exist after invalidation")
	}

	// Test InvalidateList
	err = factory.InvalidateList("ip")
	if err != nil {
		t.Errorf("InvalidateList failed: %v", err)
	}

	// Verify list items are gone
	_, exists, err = factory.Get("ip:list:1:10:search:active")
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if exists {
		t.Error("Expected list key to not exist after invalidation")
	}

	// Test InvalidateStats
	err = factory.InvalidateStats("ip")
	if err != nil {
		t.Errorf("InvalidateStats failed: %v", err)
	}

	// Verify stats items are gone
	_, exists, err = factory.Get("ip:stats")
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if exists {
		t.Error("Expected stats key to not exist after invalidation")
	}

	// Test InvalidateAll
	err = factory.InvalidateAll("ip")
	if err != nil {
		t.Errorf("InvalidateAll failed: %v", err)
	}

	// Verify all IP items are gone
	for key := range testData {
		if key[:2] == "ip" {
			_, exists, err := factory.Get(key)
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
			value, exists, err := factory.Get(key)
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

func TestCacheFactory_GetCacheType(t *testing.T) {
	factory := GetCacheFactory()

	cacheType := factory.GetCacheType()

	// Should return either "in-memory" or "distributed" based on config
	if cacheType != "in-memory" && cacheType != "distributed" {
		t.Errorf("Expected cache type to be 'in-memory' or 'distributed', got %s", cacheType)
	}
}

func TestCacheFactory_Stop(t *testing.T) {
	factory := GetCacheFactory()
	defer factory.Clear()

	// Set some items
	err := factory.Set("test-key", "test-value", 1*time.Hour)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	// Stop the factory
	factory.Stop()

	// Verify we can still get items (Stop only stops cleanup goroutine)
	value, exists, err := factory.Get("test-key")
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
