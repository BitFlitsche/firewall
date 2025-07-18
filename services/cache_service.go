package services

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"
)

// CacheItem represents a cached item with expiration
type CacheItem struct {
	Value     interface{}
	ExpiresAt time.Time
	CreatedAt time.Time
}

// Cache provides thread-safe in-memory caching with TTL and automatic eviction
type Cache struct {
	data   map[string]CacheItem
	mu     sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

var (
	cacheInstance *Cache
	cacheOnce     sync.Once
)

// GetCache returns the singleton cache instance
func GetCache() *Cache {
	cacheOnce.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		cacheInstance = &Cache{
			data:   make(map[string]CacheItem),
			ctx:    ctx,
			cancel: cancel,
		}
		cacheInstance.startCleanup()
	})
	return cacheInstance
}

// Set stores a value in the cache with TTL
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = CacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
		CreatedAt: time.Now(),
	}
	return nil
}

// Get retrieves a value from the cache
func (c *Cache) Get(key string) (interface{}, bool, error) {
	c.mu.RLock()
	item, exists := c.data[key]
	c.mu.RUnlock()

	if !exists {
		return nil, false, nil
	}

	// Check if expired
	if time.Now().After(item.ExpiresAt) {
		c.mu.Lock()
		delete(c.data, key)
		c.mu.Unlock()
		return nil, false, nil
	}

	return item.Value, true, nil
}

// Delete removes a specific key from cache
func (c *Cache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	return nil
}

// InvalidateByType removes all cached items for a specific data type
func (c *Cache) InvalidateByType(dataType string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	count := 0
	for key := range c.data {
		if strings.HasPrefix(key, dataType+":") {
			delete(c.data, key)
			count++
		}
	}

	if count > 0 {
		log.Printf("Cache: Invalidated %d items for type '%s'", count, dataType)
	}
	return nil
}

// InvalidatePattern removes all cached items matching a pattern
func (c *Cache) InvalidatePattern(pattern string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	count := 0
	for key := range c.data {
		if strings.Contains(key, pattern) {
			delete(c.data, key)
			count++
		}
	}

	if count > 0 {
		log.Printf("Cache: Invalidated %d items matching pattern '%s'", count, pattern)
	}
	return nil
}

// Clear removes all items from cache
func (c *Cache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	count := len(c.data)
	c.data = make(map[string]CacheItem)
	log.Printf("Cache: Cleared all %d items", count)
	return nil
}

// Stats returns cache statistics
func (c *Cache) Stats() (map[string]interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	total := len(c.data)
	expired := 0
	now := time.Now()

	for _, item := range c.data {
		if now.After(item.ExpiresAt) {
			expired++
		}
	}

	// Calculate memory usage in bytes (rough estimate)
	memoryUsage := total * 100 // Rough estimate per item

	return map[string]interface{}{
		"items":         total,
		"expired_items": expired,
		"valid_items":   total - expired,
		"memory_usage":  memoryUsage,
		"type":          "in-memory",
	}, nil
}

// startCleanup begins the background cleanup process
func (c *Cache) startCleanup() {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		ticker := time.NewTicker(5 * time.Minute) // Cleanup every 5 minutes
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c.cleanup()
			case <-c.ctx.Done():
				return
			}
		}
	}()
}

// cleanup removes expired items from cache
func (c *Cache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	removed := 0

	for key, item := range c.data {
		if now.After(item.ExpiresAt) {
			delete(c.data, key)
			removed++
		}
	}

	if removed > 0 {
		log.Printf("Cache: Cleaned up %d expired items", removed)
	}
}

// Stop gracefully stops the cache cleanup
func (c *Cache) Stop() {
	c.cancel()
	c.wg.Wait()
	log.Println("Cache service stopped")
}

// Cache key generators for different data types
func (c *Cache) FilterKey(dataType, value string) string {
	return dataType + ":filter:" + value
}

func (c *Cache) ListKey(dataType, page, limit, search, status string) string {
	return dataType + ":list:" + page + ":" + limit + ":" + search + ":" + status
}

func (c *Cache) StatsKey(dataType string) string {
	return dataType + ":stats"
}

// Cache invalidation helpers
func (c *Cache) InvalidateFilter(dataType string) error {
	return c.InvalidatePattern(dataType + ":filter:")
}

func (c *Cache) InvalidateList(dataType string) error {
	return c.InvalidatePattern(dataType + ":list:")
}

func (c *Cache) InvalidateStats(dataType string) error {
	return c.Delete(dataType + ":stats")
}

// InvalidateAll invalidates all cache for a data type
func (c *Cache) InvalidateAll(dataType string) error {
	return c.InvalidateByType(dataType)
}
