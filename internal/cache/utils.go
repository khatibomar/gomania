package cache

import (
	"fmt"
	"time"
)

// KeyBuilder provides a consistent way to build cache keys
type KeyBuilder struct {
	prefix string
}

// NewKeyBuilder creates a new key builder with the given prefix
func NewKeyBuilder(prefix string) *KeyBuilder {
	return &KeyBuilder{prefix: prefix}
}

// Build creates a cache key with the given parts
func (kb *KeyBuilder) Build(parts ...string) string {
	if len(parts) == 0 {
		return kb.prefix
	}

	key := kb.prefix
	for _, part := range parts {
		key += ":" + part
	}
	return key
}

// BuildWithID creates a cache key for a specific ID
func (kb *KeyBuilder) BuildWithID(id string) string {
	return kb.Build(id)
}

// BuildList creates a cache key for list operations
func (kb *KeyBuilder) BuildList() string {
	return kb.Build("list")
}

// BuildSearch creates a cache key for search operations
func (kb *KeyBuilder) BuildSearch(query string) string {
	return kb.Build("search", query)
}

// BuildCategory creates a cache key for category operations
func (kb *KeyBuilder) BuildCategory(categoryID string) string {
	return kb.Build("category", categoryID)
}

// CacheWrapper provides higher-level caching operations
type CacheWrapper struct {
	cache Cache
	kb    *KeyBuilder
}

// NewCacheWrapper creates a new cache wrapper with the given cache and key prefix
func NewCacheWrapper(cache Cache, keyPrefix string) *CacheWrapper {
	return &CacheWrapper{
		cache: cache,
		kb:    NewKeyBuilder(keyPrefix),
	}
}

// GetOrSet retrieves a value from cache or sets it using the provider function
func (cw *CacheWrapper) GetOrSet(key string, provider func() (any, error)) (any, error) {
	// Try to get from cache first
	if value, found := cw.cache.Get(key); found {
		return value, nil
	}

	// Cache miss - use provider function
	value, err := provider()
	if err != nil {
		return nil, err
	}

	// Cache the result
	cw.cache.Set(key, value)
	return value, nil
}

// GetByID retrieves an item by ID from cache or uses the provider function
func (cw *CacheWrapper) GetByID(id string, provider func() (any, error)) (any, error) {
	key := cw.kb.BuildWithID(id)
	return cw.GetOrSet(key, provider)
}

// GetList retrieves a list from cache or uses the provider function
func (cw *CacheWrapper) GetList(provider func() (any, error)) (any, error) {
	key := cw.kb.BuildList()
	return cw.GetOrSet(key, provider)
}

// GetSearch retrieves search results from cache or uses the provider function
func (cw *CacheWrapper) GetSearch(query string, provider func() (any, error)) (any, error) {
	key := cw.kb.BuildSearch(query)
	return cw.GetOrSet(key, provider)
}

// GetByCategory retrieves items by category from cache or uses the provider function
func (cw *CacheWrapper) GetByCategory(categoryID string, provider func() (any, error)) (any, error) {
	key := cw.kb.BuildCategory(categoryID)
	return cw.GetOrSet(key, provider)
}

// InvalidateID invalidates a specific item by ID
func (cw *CacheWrapper) InvalidateID(id string) {
	key := cw.kb.BuildWithID(id)
	cw.cache.Delete(key)
}

// InvalidateList invalidates the list cache
func (cw *CacheWrapper) InvalidateList() {
	key := cw.kb.BuildList()
	cw.cache.Delete(key)
}

// InvalidateSearch invalidates all search results
func (cw *CacheWrapper) InvalidateSearch() {
	pattern := cw.kb.Build("search", "*")
	cw.cache.InvalidatePattern(pattern)
}

// InvalidateCategory invalidates all category-related cache entries
func (cw *CacheWrapper) InvalidateCategory() {
	pattern := cw.kb.Build("category", "*")
	cw.cache.InvalidatePattern(pattern)
}

// InvalidateAll invalidates all cache entries with this wrapper's prefix
func (cw *CacheWrapper) InvalidateAll() {
	pattern := cw.kb.Build("*")
	cw.cache.InvalidatePattern(pattern)
}

// CacheConfig holds configuration for different cache scenarios
type CacheConfig struct {
	DefaultTTL time.Duration
	ShortTTL   time.Duration
	LongTTL    time.Duration
}

// DefaultCacheConfig returns a sensible default cache configuration
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{
		DefaultTTL: 15 * time.Minute,
		ShortTTL:   5 * time.Minute,
		LongTTL:    1 * time.Hour,
	}
}

// CacheManager manages multiple cache instances with different TTLs
type CacheManager struct {
	defaultCache Cache
	shortCache   Cache
	longCache    Cache
}

// NewCacheManager creates a new cache manager with different TTL caches
func NewCacheManager(config CacheConfig) *CacheManager {
	return &CacheManager{
		defaultCache: NewMemoryCache(config.DefaultTTL),
		shortCache:   NewMemoryCache(config.ShortTTL),
		longCache:    NewMemoryCache(config.LongTTL),
	}
}

// Default returns the default cache instance
func (cm *CacheManager) Default() Cache {
	return cm.defaultCache
}

// Short returns the short TTL cache instance
func (cm *CacheManager) Short() Cache {
	return cm.shortCache
}

// Long returns the long TTL cache instance
func (cm *CacheManager) Long() Cache {
	return cm.longCache
}

// Close closes all cache instances
func (cm *CacheManager) Close() {
	cm.defaultCache.Close()
	cm.shortCache.Close()
	cm.longCache.Close()
}

// Common cache key patterns
const (
	KeyPatternAll              = "*"
	KeyPatternUsers            = "user:*"
	KeyPatternPrograms         = "program:*"
	KeyPatternProgramsSearch   = "programs:search:*"
	KeyPatternProgramsCategory = "programs:category:*"
	KeyPatternCategories       = "categories:*"
)

// Helper functions for common cache operations

// CacheKey builds a cache key from components
func CacheKey(components ...string) string {
	if len(components) == 0 {
		return ""
	}

	key := components[0]
	for i := 1; i < len(components); i++ {
		key += ":" + components[i]
	}
	return key
}

// ProgramKey builds a cache key for a program
func ProgramKey(id string) string {
	return CacheKey("program", id)
}

// ProgramsListKey builds a cache key for programs list
func ProgramsListKey() string {
	return CacheKey("programs", "list")
}

// ProgramsSearchKey builds a cache key for program search
func ProgramsSearchKey(query string) string {
	return CacheKey("programs", "search", query)
}

// ProgramsCategoryKey builds a cache key for programs by category
func ProgramsCategoryKey(categoryID string) string {
	return CacheKey("programs", "category", categoryID)
}

// CategoryKey builds a cache key for a category
func CategoryKey(id string) string {
	return CacheKey("category", id)
}

// CategoriesListKey builds a cache key for categories list
func CategoriesListKey() string {
	return CacheKey("categories", "list")
}

// SafeTypeAssert safely performs type assertion with error handling
func SafeTypeAssert[T any](value any) (T, error) {
	var zero T
	if value == nil {
		return zero, fmt.Errorf("value is nil")
	}

	result, ok := value.(T)
	if !ok {
		return zero, fmt.Errorf("type assertion failed: expected %T, got %T", zero, value)
	}

	return result, nil
}

// CacheMetrics provides basic metrics for cache operations
type CacheMetrics struct {
	Hits   int64
	Misses int64
}

// HitRate calculates the cache hit rate
func (m *CacheMetrics) HitRate() float64 {
	total := m.Hits + m.Misses
	if total == 0 {
		return 0
	}
	return float64(m.Hits) / float64(total)
}

// Total returns the total number of cache operations
func (m *CacheMetrics) Total() int64 {
	return m.Hits + m.Misses
}

// RecordHit increments the hit counter
func (m *CacheMetrics) RecordHit() {
	m.Hits++
}

// RecordMiss increments the miss counter
func (m *CacheMetrics) RecordMiss() {
	m.Misses++
}

// Reset resets all counters
func (m *CacheMetrics) Reset() {
	m.Hits = 0
	m.Misses = 0
}
