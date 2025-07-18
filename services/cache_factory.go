package services

import (
	"firewall/config"
	"log"
	"sync"
	"time"
)

// CacheFactory provides a unified interface for caching operations
type CacheFactory struct {
	cache CacheInterface
}

var (
	cacheFactoryInstance *CacheFactory
	cacheFactoryOnce     sync.Once
)

// GetCacheFactory returns the singleton cache factory instance
func GetCacheFactory() *CacheFactory {
	cacheFactoryOnce.Do(func() {
		var cache CacheInterface

		// Check if distributed caching is enabled
		if config.AppConfig.Caching.Distributed {
			distributedCache := GetDistributedCache()
			if distributedCache != nil {
				cache = distributedCache
				log.Println("Using distributed cache (Redis)")
			} else {
				// Fallback to in-memory cache if distributed cache fails to initialize
				cache = GetCache()
				log.Println("Falling back to in-memory cache")
			}
		} else {
			// Use in-memory cache for single instance
			cache = GetCache()
			log.Println("Using in-memory cache")
		}

		cacheFactoryInstance = &CacheFactory{
			cache: cache,
		}
	})
	return cacheFactoryInstance
}

// Set stores a value in the cache with TTL
func (cf *CacheFactory) Set(key string, value interface{}, ttl time.Duration) error {
	return cf.cache.Set(key, value, ttl)
}

// Get retrieves a value from the cache
func (cf *CacheFactory) Get(key string) (interface{}, bool, error) {
	return cf.cache.Get(key)
}

// Delete removes a specific key from cache
func (cf *CacheFactory) Delete(key string) error {
	return cf.cache.Delete(key)
}

// InvalidateByType removes all cached items for a specific data type
func (cf *CacheFactory) InvalidateByType(dataType string) error {
	return cf.cache.InvalidateByType(dataType)
}

// InvalidatePattern removes all cached items matching a pattern
func (cf *CacheFactory) InvalidatePattern(pattern string) error {
	return cf.cache.InvalidatePattern(pattern)
}

// Clear removes all items from cache
func (cf *CacheFactory) Clear() error {
	return cf.cache.Clear()
}

// Stats returns cache statistics
func (cf *CacheFactory) Stats() (map[string]interface{}, error) {
	return cf.cache.Stats()
}

// Stop gracefully stops the cache service
func (cf *CacheFactory) Stop() {
	cf.cache.Stop()
}

// Cache key generators
func (cf *CacheFactory) FilterKey(dataType, value string) string {
	return cf.cache.FilterKey(dataType, value)
}

func (cf *CacheFactory) ListKey(dataType, page, limit, search, status string) string {
	return cf.cache.ListKey(dataType, page, limit, search, status)
}

func (cf *CacheFactory) StatsKey(dataType string) string {
	return cf.cache.StatsKey(dataType)
}

// Cache invalidation helpers
func (cf *CacheFactory) InvalidateFilter(dataType string) error {
	return cf.cache.InvalidateFilter(dataType)
}

func (cf *CacheFactory) InvalidateList(dataType string) error {
	return cf.cache.InvalidateList(dataType)
}

func (cf *CacheFactory) InvalidateStats(dataType string) error {
	return cf.cache.InvalidateStats(dataType)
}

func (cf *CacheFactory) InvalidateAll(dataType string) error {
	return cf.cache.InvalidateAll(dataType)
}

// GetCacheType returns the type of cache being used
func (cf *CacheFactory) GetCacheType() string {
	if config.AppConfig.Caching.Distributed {
		return "distributed"
	}
	return "in-memory"
}
