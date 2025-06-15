package cache

import (
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

// Entry represents a cached item with expiration time
type Entry struct {
	Data      any
	ExpiresAt time.Time
}

// IsExpired checks if the cache entry has expired
func (e *Entry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// Memory provides thread-safe in-memory caching with TTL support
type Memory struct {
	mu    sync.RWMutex
	items map[string]*Entry
	ttl   time.Duration
	stop  chan struct{}
	// Metrics
	hits    int64
	misses  int64
	evicted int64
}

// New creates a new memory cache instance with the specified TTL
func New(ttl time.Duration) *Memory {
	cache := &Memory{
		items: make(map[string]*Entry),
		ttl:   ttl,
		stop:  make(chan struct{}),
	}

	// Start background cleanup goroutine
	go cache.cleanup()

	return cache
}

// Set stores an item in the cache
func (c *Memory) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = &Entry{
		Data:      value,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

// Get retrieves an item from the cache
func (c *Memory) Get(key string) (any, bool) {
	c.mu.RLock()
	entry, exists := c.items[key]
	if !exists {
		c.mu.RUnlock()
		return nil, false
	}

	if entry.IsExpired() {
		c.mu.RUnlock()
		// Upgrade to write lock to remove expired entry
		c.mu.Lock()
		// Double-check after acquiring write lock
		if entry, exists := c.items[key]; exists && entry.IsExpired() {
			delete(c.items, key)
			atomic.AddInt64(&c.evicted, 1)
		}
		c.mu.Unlock()
		atomic.AddInt64(&c.misses, 1)
		return nil, false
	}

	data := entry.Data
	c.mu.RUnlock()
	atomic.AddInt64(&c.hits, 1)
	return data, true
}

// Delete removes an item from the cache
func (c *Memory) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Clear removes all items from the cache
func (c *Memory) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*Entry)
}

// InvalidatePattern removes all cache entries that match a pattern
func (c *Memory) InvalidatePattern(pattern string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key := range c.items {
		if match, _ := filepath.Match(pattern, key); match {
			delete(c.items, key)
		}
	}
}

// Size returns the number of items in the cache
func (c *Memory) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// TTL returns the current TTL setting
func (c *Memory) TTL() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ttl
}

// SetTTL updates the cache TTL (affects new entries only)
func (c *Memory) SetTTL(ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ttl = ttl
}

// Keys returns all cache keys (useful for debugging)
func (c *Memory) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.items))
	for key := range c.items {
		keys = append(keys, key)
	}
	return keys
}

// Metrics returns cache performance metrics
func (c *Memory) Metrics() (hits, misses, evicted int64) {
	return atomic.LoadInt64(&c.hits), atomic.LoadInt64(&c.misses), atomic.LoadInt64(&c.evicted)
}

// ResetMetrics resets all performance counters
func (c *Memory) ResetMetrics() {
	atomic.StoreInt64(&c.hits, 0)
	atomic.StoreInt64(&c.misses, 0)
	atomic.StoreInt64(&c.evicted, 0)
}

// Close stops the cleanup goroutine
func (c *Memory) Close() {
	close(c.stop)
}

// cleanup runs periodically to remove expired entries
func (c *Memory) cleanup() {
	ticker := time.NewTicker(time.Minute * 5) // Clean up every 5 minutes
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			evictedCount := int64(0)
			for key, entry := range c.items {
				if entry.IsExpired() {
					delete(c.items, key)
					evictedCount++
				}
			}
			if evictedCount > 0 {
				atomic.AddInt64(&c.evicted, evictedCount)
			}
			c.mu.Unlock()
		case <-c.stop:
			return
		}
	}
}
