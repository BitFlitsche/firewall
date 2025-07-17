package services

import (
	"fmt"
	"time"
)

// ExampleCacheUsage demonstrates how to use the cache service
func ExampleCacheUsage() {
	cache := GetCache()

	// Example 1: Cache filter results with 5-minute TTL
	filterKey := cache.FilterKey("ip", "192.168.1.1")
	if cached, exists := cache.Get(filterKey); exists {
		// Use cached result
		fmt.Printf("Using cached filter result: %v\n", cached)
	} else {
		// Perform expensive filter operation
		result := performExpensiveFilter("192.168.1.1")

		// Cache the result for 5 minutes
		cache.Set(filterKey, result, 5*time.Minute)
		fmt.Printf("Cached new filter result: %v\n", result)
	}

	// Example 2: Cache list data with 2-minute TTL
	listKey := cache.ListKey("ip", "1", "10", "", "all")
	if cached, exists := cache.Get(listKey); exists {
		// Use cached list
		fmt.Printf("Using cached list: %v\n", cached)
	} else {
		// Perform expensive database query
		result := performExpensiveListQuery("ip", 1, 10, "", "all")

		// Cache the result for 2 minutes
		cache.Set(listKey, result, 2*time.Minute)
		fmt.Printf("Cached new list: %v\n", result)
	}

	// Example 3: Cache system stats with 30-second TTL
	statsKey := cache.StatsKey("system")
	if cached, exists := cache.Get(statsKey); exists {
		// Use cached stats
		fmt.Printf("Using cached stats: %v\n", cached)
	} else {
		// Collect system stats
		result := collectSystemStats()

		// Cache the result for 30 seconds
		cache.Set(statsKey, result, 30*time.Second)
		fmt.Printf("Cached new stats: %v\n", result)
	}
}

// Example controller function with caching
func ExampleControllerWithCache() {
	cache := GetCache()

	// Cache key for this specific request
	cacheKey := fmt.Sprintf("api:ips:page:%d:limit:%d:search:%s:status:%s",
		1, 10, "192.168", "all")

	// Try to get from cache first
	if cached, exists := cache.Get(cacheKey); exists {
		// Return cached response
		fmt.Printf("Cache hit: %v\n", cached)
		return
	}

	// Cache miss - perform database query
	result := performDatabaseQuery()

	// Cache the result for 2 minutes
	cache.Set(cacheKey, result, 2*time.Minute)

	// Return the result
	fmt.Printf("Cache miss - stored new result: %v\n", result)
}

// Example of cache invalidation in update operations
func ExampleCacheInvalidation() {
	cache := GetCache()

	// After updating an IP address
	// This would be called in your update controller
	cache.InvalidateAll("ip")    // Invalidate all IP-related cache
	cache.InvalidateFilter("ip") // Invalidate IP filter cache
	cache.InvalidateList("ip")   // Invalidate IP list cache
	cache.InvalidateStats("ip")  // Invalidate IP stats cache

	// Or invalidate specific patterns
	cache.InvalidatePattern("ip:192.168.1.1") // Specific IP
	cache.InvalidatePattern("filter:ip")      // All IP filters
}

// Helper functions for examples
func performExpensiveFilter(ip string) interface{} {
	// Simulate expensive filter operation
	return map[string]interface{}{
		"result": "denied",
		"reason": "ip_blocked",
		"ip":     ip,
	}
}

func performExpensiveListQuery(dataType string, page, limit int, search, status string) interface{} {
	// Simulate expensive database query
	return map[string]interface{}{
		"data":      []interface{}{"item1", "item2", "item3"},
		"total":     100,
		"page":      page,
		"limit":     limit,
		"data_type": dataType,
		"search":    search,
		"status":    status,
	}
}

func performDatabaseQuery() interface{} {
	// Simulate database query
	return map[string]interface{}{
		"items": []interface{}{"192.168.1.1", "192.168.1.2"},
		"total": 2,
	}
}

func collectSystemStats() interface{} {
	// Simulate system stats collection
	return map[string]interface{}{
		"cpu_percent":    15.5,
		"memory_percent": 45.2,
		"uptime":         3600,
	}
}
