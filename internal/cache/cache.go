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

	// Close performs cleanup operations
	Close()
}

// NewMemoryCache creates a new memory cache instance
func NewMemoryCache(ttl time.Duration) Cache {
	return New(ttl)
}
