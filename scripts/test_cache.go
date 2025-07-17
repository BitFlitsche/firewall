package main

import (
	"firewall/services"
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== Testing In-Memory Cache with TTL and Auto-Invalidation ===")

	// Initialize cache
	cache := services.GetCache()

	// Test 1: Basic caching with TTL
	fmt.Println("\n1. Testing basic caching with TTL...")

	// Set a value with 5-second TTL
	cache.Set("test:key1", "test_value_1", 5*time.Second)
	fmt.Println("Set 'test:key1' with 5-second TTL")

	// Get the value immediately
	if value, exists := cache.Get("test:key1"); exists {
		fmt.Printf("✓ Retrieved: %v\n", value)
	} else {
		fmt.Println("✗ Failed to retrieve value")
	}

	// Test 2: TTL expiration
	fmt.Println("\n2. Testing TTL expiration...")

	// Set a value with 2-second TTL
	cache.Set("test:key2", "test_value_2", 2*time.Second)
	fmt.Println("Set 'test:key2' with 2-second TTL")

	// Wait 3 seconds
	fmt.Println("Waiting 3 seconds for expiration...")
	time.Sleep(3 * time.Second)

	// Try to get the expired value
	if value, exists := cache.Get("test:key2"); exists {
		fmt.Printf("✗ Value still exists: %v (should be expired)\n", value)
	} else {
		fmt.Println("✓ Value correctly expired")
	}

	// Test 3: Cache invalidation by type
	fmt.Println("\n3. Testing cache invalidation by type...")

	// Set multiple values for different types
	cache.Set("ip:filter:192.168.1.1", "denied", 10*time.Minute)
	cache.Set("ip:list:page1", "ip_list_data", 10*time.Minute)
	cache.Set("email:filter:test@example.com", "allowed", 10*time.Minute)
	cache.Set("email:list:page1", "email_list_data", 10*time.Minute)

	fmt.Println("Set cache items for 'ip' and 'email' types")

	// Invalidate all IP-related cache
	cache.InvalidateAll("ip")
	fmt.Println("Invalidated all 'ip' cache")

	// Check what's left
	if _, exists := cache.Get("ip:filter:192.168.1.1"); exists {
		fmt.Println("✗ IP filter cache still exists (should be invalidated)")
	} else {
		fmt.Println("✓ IP filter cache correctly invalidated")
	}

	if _, exists := cache.Get("email:filter:test@example.com"); exists {
		fmt.Println("✓ Email filter cache still exists (correctly not invalidated)")
	} else {
		fmt.Println("✗ Email filter cache incorrectly invalidated")
	}

	// Test 4: Pattern-based invalidation
	fmt.Println("\n4. Testing pattern-based invalidation...")

	// Set values with specific patterns
	cache.Set("filter:ip:192.168.1.1", "denied", 10*time.Minute)
	cache.Set("filter:ip:192.168.1.2", "allowed", 10*time.Minute)
	cache.Set("list:ip:page1", "list_data", 10*time.Minute)

	fmt.Println("Set cache items with 'filter:ip' and 'list:ip' patterns")

	// Invalidate only filter patterns
	cache.InvalidatePattern("filter:ip")
	fmt.Println("Invalidated 'filter:ip' pattern")

	// Check results
	if _, exists := cache.Get("filter:ip:192.168.1.1"); exists {
		fmt.Println("✗ Filter cache still exists (should be invalidated)")
	} else {
		fmt.Println("✓ Filter cache correctly invalidated")
	}

	if _, exists := cache.Get("list:ip:page1"); exists {
		fmt.Println("✓ List cache still exists (correctly not invalidated)")
	} else {
		fmt.Println("✗ List cache incorrectly invalidated")
	}

	// Test 5: Cache statistics
	fmt.Println("\n5. Testing cache statistics...")

	// Add some test data
	cache.Set("stats:test1", "value1", 1*time.Minute)
	cache.Set("stats:test2", "value2", 1*time.Minute)

	stats := cache.Stats()
	fmt.Printf("Cache Statistics:\n")
	for key, value := range stats {
		fmt.Printf("  %s: %v\n", key, value)
	}

	// Test 6: Event-driven invalidation simulation
	fmt.Println("\n6. Testing event-driven invalidation...")

	// Simulate IP update event
	fmt.Println("Simulating IP update event...")
	services.PublishEvent("ip", "updated", map[string]string{"address": "192.168.1.1"})

	// Wait a moment for event processing
	time.Sleep(100 * time.Millisecond)

	// Check if IP cache was invalidated
	if _, exists := cache.Get("ip:filter:192.168.1.1"); exists {
		fmt.Println("✗ IP cache still exists after update event")
	} else {
		fmt.Println("✓ IP cache correctly invalidated by update event")
	}

	fmt.Println("\n=== Cache Test Completed Successfully! ===")

	// Stop the cache service
	cache.Stop()
}
