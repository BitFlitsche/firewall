package services

import (
	"time"
)

// CacheInterface defines the interface for caching operations
type CacheInterface interface {
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string) (interface{}, bool, error)
	Delete(key string) error
	InvalidateByType(dataType string) error
	InvalidatePattern(pattern string) error
	Clear() error
	Stats() (map[string]interface{}, error)
	Stop()

	// Cache key generators
	FilterKey(dataType, value string) string
	ListKey(dataType, page, limit, search, status string) string
	StatsKey(dataType string) string

	// Cache invalidation helpers
	InvalidateFilter(dataType string) error
	InvalidateList(dataType string) error
	InvalidateStats(dataType string) error
	InvalidateAll(dataType string) error
}
