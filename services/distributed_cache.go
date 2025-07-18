package services

import (
	"context"
	"encoding/json"
	"firewall/config"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// DistributedCache provides Redis-based distributed caching
type DistributedCache struct {
	client *redis.Client
	ctx    context.Context
}

// DistributedCacheItem represents a cached item with expiration
type DistributedCacheItem struct {
	Value     interface{} `json:"value"`
	ExpiresAt time.Time   `json:"expires_at"`
	CreatedAt time.Time   `json:"created_at"`
}

var (
	distributedCache     *DistributedCache
	distributedCacheOnce sync.Once
)

// GetDistributedCache returns the singleton distributed cache service
func GetDistributedCache() *DistributedCache {
	distributedCacheOnce.Do(func() {
		// Check if distributed caching is enabled
		if !config.AppConfig.Caching.Distributed {
			log.Println("Distributed caching disabled - using in-memory cache")
			return
		}

		// Initialize Redis client
		redisClient := redis.NewClient(&redis.Options{
			Addr:     config.AppConfig.Redis.GetRedisAddr(),
			Password: config.AppConfig.Redis.Password,
			DB:       config.AppConfig.Redis.DB,
		})

		distributedCache = &DistributedCache{
			client: redisClient,
			ctx:    context.Background(),
		}

		log.Println("Distributed cache service initialized")
	})
	return distributedCache
}

// Set stores a value in the distributed cache with TTL
func (dc *DistributedCache) Set(key string, value interface{}, ttl time.Duration) error {
	if dc == nil {
		return fmt.Errorf("distributed cache not initialized")
	}

	cacheItem := DistributedCacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
		CreatedAt: time.Now(),
	}

	// Serialize the cache item
	data, err := json.Marshal(cacheItem)
	if err != nil {
		return fmt.Errorf("failed to serialize cache item: %v", err)
	}

	// Store in Redis with TTL
	err = dc.client.Set(dc.ctx, key, data, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache item: %v", err)
	}

	return nil
}

// Get retrieves a value from the distributed cache
func (dc *DistributedCache) Get(key string) (interface{}, bool, error) {
	if dc == nil {
		return nil, false, fmt.Errorf("distributed cache not initialized")
	}

	// Get from Redis
	data, err := dc.client.Get(dc.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, false, nil // Key not found
		}
		return nil, false, fmt.Errorf("failed to get cache item: %v", err)
	}

	// Deserialize the cache item
	var cacheItem DistributedCacheItem
	if err := json.Unmarshal([]byte(data), &cacheItem); err != nil {
		return nil, false, fmt.Errorf("failed to deserialize cache item: %v", err)
	}

	// Check if expired
	if time.Now().After(cacheItem.ExpiresAt) {
		// Remove expired item
		dc.client.Del(dc.ctx, key)
		return nil, false, nil
	}

	return cacheItem.Value, true, nil
}

// Delete removes a specific key from cache
func (dc *DistributedCache) Delete(key string) error {
	if dc == nil {
		return fmt.Errorf("distributed cache not initialized")
	}

	return dc.client.Del(dc.ctx, key).Err()
}

// InvalidateByType removes all cached items for a specific data type
func (dc *DistributedCache) InvalidateByType(dataType string) error {
	if dc == nil {
		return fmt.Errorf("distributed cache not initialized")
	}

	pattern := dataType + ":*"
	keys, err := dc.client.Keys(dc.ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get keys for pattern %s: %v", pattern, err)
	}

	if len(keys) > 0 {
		err = dc.client.Del(dc.ctx, keys...).Err()
		if err != nil {
			return fmt.Errorf("failed to delete keys: %v", err)
		}
		log.Printf("Distributed Cache: Invalidated %d items for type '%s'", len(keys), dataType)
	}

	return nil
}

// InvalidatePattern removes all cached items matching a pattern
func (dc *DistributedCache) InvalidatePattern(pattern string) error {
	if dc == nil {
		return fmt.Errorf("distributed cache not initialized")
	}

	keys, err := dc.client.Keys(dc.ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get keys for pattern %s: %v", pattern, err)
	}

	if len(keys) > 0 {
		err = dc.client.Del(dc.ctx, keys...).Err()
		if err != nil {
			return fmt.Errorf("failed to delete keys: %v", err)
		}
		log.Printf("Distributed Cache: Invalidated %d items matching pattern '%s'", len(keys), pattern)
	}

	return nil
}

// Clear removes all items from cache
func (dc *DistributedCache) Clear() error {
	if dc == nil {
		return fmt.Errorf("distributed cache not initialized")
	}

	// Get all keys
	keys, err := dc.client.Keys(dc.ctx, "*").Result()
	if err != nil {
		return fmt.Errorf("failed to get all keys: %v", err)
	}

	if len(keys) > 0 {
		err = dc.client.Del(dc.ctx, keys...).Err()
		if err != nil {
			return fmt.Errorf("failed to delete all keys: %v", err)
		}
		log.Printf("Distributed Cache: Cleared all %d items", len(keys))
	}

	return nil
}

// Stats returns cache statistics
func (dc *DistributedCache) Stats() (map[string]interface{}, error) {
	if dc == nil {
		return nil, fmt.Errorf("distributed cache not initialized")
	}

	// Get all keys
	keys, err := dc.client.Keys(dc.ctx, "*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys: %v", err)
	}

	total := len(keys)
	expired := 0
	now := time.Now()

	// Check for expired items
	for _, key := range keys {
		data, err := dc.client.Get(dc.ctx, key).Result()
		if err != nil {
			continue
		}

		var cacheItem DistributedCacheItem
		if err := json.Unmarshal([]byte(data), &cacheItem); err != nil {
			continue
		}

		if now.After(cacheItem.ExpiresAt) {
			expired++
		}
	}

	// Calculate memory usage (rough estimate)
	memoryUsage := total * 200 // Rough estimate per item in distributed cache

	return map[string]interface{}{
		"items":         total,
		"expired_items": expired,
		"valid_items":   total - expired,
		"memory_usage":  memoryUsage,
		"type":          "distributed",
	}, nil
}

// Stop gracefully stops the distributed cache service
func (dc *DistributedCache) Stop() {
	if dc != nil && dc.client != nil {
		dc.client.Close()
		log.Println("Distributed cache service stopped")
	}
}

// Cache key generators (same as in-memory cache for compatibility)
func (dc *DistributedCache) FilterKey(dataType, value string) string {
	return dataType + ":filter:" + value
}

func (dc *DistributedCache) ListKey(dataType, page, limit, search, status string) string {
	return dataType + ":list:" + page + ":" + limit + ":" + search + ":" + status
}

func (dc *DistributedCache) StatsKey(dataType string) string {
	return dataType + ":stats"
}

// Cache invalidation helpers
func (dc *DistributedCache) InvalidateFilter(dataType string) error {
	return dc.InvalidatePattern(dataType + ":filter:")
}

func (dc *DistributedCache) InvalidateList(dataType string) error {
	return dc.InvalidatePattern(dataType + ":list:")
}

func (dc *DistributedCache) InvalidateStats(dataType string) error {
	return dc.Delete(dataType + ":stats")
}

// InvalidateAll invalidates all cache for a data type
func (dc *DistributedCache) InvalidateAll(dataType string) error {
	return dc.InvalidateByType(dataType)
}
