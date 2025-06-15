package cache

import "time"

// Cache defines the interface for cache implementations
type Cache interface {
	// Set stores an item in the cache
	Set(key string, value any)

	// Get retrieves an item from the cache
	Get(key string) (any, bool)

	// Delete removes an item from the cache
	Delete(key string)

	// Clear removes all items from the cache
	Clear()

	// InvalidatePattern removes all cache entries that match a pattern
	InvalidatePattern(pattern string)

	// Size returns the number of items in the cache
	Size() int

	// TTL returns the current TTL setting
	TTL() time.Duration

	// SetTTL updates the cache TTL (affects new entries only)
	SetTTL(ttl time.Duration)

	// Keys returns all cache keys (useful for debugging)
	Keys() []string

	// Close performs cleanup operations
	Close()

	// Metrics returns cache performance metrics (hits, misses, evicted)
	Metrics() (int64, int64, int64)

	// ResetMetrics resets all performance counters
	ResetMetrics()
}

// Stats represents cache statistics
type Stats struct {
	TotalItems int           `json:"total_items"`
	TTL        time.Duration `json:"ttl"`
	Hits       int64         `json:"hits"`
	Misses     int64         `json:"misses"`
	Evicted    int64         `json:"evicted"`
}

// NewMemoryCache creates a new memory cache instance
func NewMemoryCache(ttl time.Duration) Cache {
	return New(ttl)
}
