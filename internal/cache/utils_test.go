package cache

import (
	"errors"
	"testing"
	"time"
)

func TestKeyBuilder(t *testing.T) {
	kb := NewKeyBuilder("test")

	tests := []struct {
		name     string
		parts    []string
		expected string
	}{
		{
			name:     "no parts",
			parts:    []string{},
			expected: "test",
		},
		{
			name:     "single part",
			parts:    []string{"123"},
			expected: "test:123",
		},
		{
			name:     "multiple parts",
			parts:    []string{"user", "123", "profile"},
			expected: "test:user:123:profile",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := kb.Build(tt.parts...)
			if result != tt.expected {
				t.Errorf("Build() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestKeyBuilder_Helpers(t *testing.T) {
	kb := NewKeyBuilder("program")

	tests := []struct {
		name     string
		method   func() string
		expected string
	}{
		{
			name:     "BuildWithID",
			method:   func() string { return kb.BuildWithID("123") },
			expected: "program:123",
		},
		{
			name:     "BuildList",
			method:   func() string { return kb.BuildList() },
			expected: "program:list",
		},
		{
			name:     "BuildSearch",
			method:   func() string { return kb.BuildSearch("golang") },
			expected: "program:search:golang",
		},
		{
			name:     "BuildCategory",
			method:   func() string { return kb.BuildCategory("tech") },
			expected: "program:category:tech",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.method()
			if result != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}

func TestCacheWrapper_GetOrSet(t *testing.T) {
	cache := NewMemoryCache(time.Minute)
	defer cache.Close()

	wrapper := NewCacheWrapper(cache, "test")

	// Test cache miss and set
	callCount := 0
	provider := func() (any, error) {
		callCount++
		return "test_value", nil
	}

	result, err := wrapper.GetOrSet("key1", provider)
	if err != nil {
		t.Errorf("GetOrSet() error = %v", err)
	}
	if result != "test_value" {
		t.Errorf("GetOrSet() = %v, want %v", result, "test_value")
	}
	if callCount != 1 {
		t.Errorf("Provider called %d times, want 1", callCount)
	}

	// Test cache hit
	result, err = wrapper.GetOrSet("key1", provider)
	if err != nil {
		t.Errorf("GetOrSet() error = %v", err)
	}
	if result != "test_value" {
		t.Errorf("GetOrSet() = %v, want %v", result, "test_value")
	}
	if callCount != 1 {
		t.Errorf("Provider called %d times, want 1 (should be cached)", callCount)
	}
}

func TestCacheWrapper_GetOrSet_ProviderError(t *testing.T) {
	cache := NewMemoryCache(time.Minute)
	defer cache.Close()

	wrapper := NewCacheWrapper(cache, "test")

	provider := func() (any, error) {
		return nil, errors.New("provider error")
	}

	result, err := wrapper.GetOrSet("key1", provider)
	if err == nil {
		t.Error("Expected error from provider")
	}
	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}

	// Verify nothing was cached
	if _, found := cache.Get("key1"); found {
		t.Error("Expected nothing to be cached on error")
	}
}

func TestCacheWrapper_HelperMethods(t *testing.T) {
	cache := NewMemoryCache(time.Minute)
	defer cache.Close()

	wrapper := NewCacheWrapper(cache, "program")

	provider := func() (any, error) {
		return "test_data", nil
	}

	tests := []struct {
		name     string
		method   func() (any, error)
		cacheKey string
	}{
		{
			name:     "GetByID",
			method:   func() (any, error) { return wrapper.GetByID("123", provider) },
			cacheKey: "program:123",
		},
		{
			name:     "GetList",
			method:   func() (any, error) { return wrapper.GetList(provider) },
			cacheKey: "program:list",
		},
		{
			name:     "GetSearch",
			method:   func() (any, error) { return wrapper.GetSearch("golang", provider) },
			cacheKey: "program:search:golang",
		},
		{
			name:     "GetByCategory",
			method:   func() (any, error) { return wrapper.GetByCategory("tech", provider) },
			cacheKey: "program:category:tech",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.method()
			if err != nil {
				t.Errorf("%s error = %v", tt.name, err)
			}
			if result != "test_data" {
				t.Errorf("%s = %v, want %v", tt.name, result, "test_data")
			}

			// Verify it was cached with the correct key
			if cached, found := cache.Get(tt.cacheKey); !found || cached != "test_data" {
				t.Errorf("Data not properly cached with key %s", tt.cacheKey)
			}
		})
	}
}

func TestCacheWrapper_Invalidation(t *testing.T) {
	cache := NewMemoryCache(time.Minute)
	defer cache.Close()

	wrapper := NewCacheWrapper(cache, "program")

	// Set up some test data
	cache.Set("program:123", "data1")
	cache.Set("program:list", "data2")
	cache.Set("program:search:golang", "data3")
	cache.Set("program:search:python", "data4")
	cache.Set("program:category:tech", "data5")
	cache.Set("program:category:science", "data6")

	tests := []struct {
		name           string
		invalidate     func()
		shouldExist    []string
		shouldNotExist []string
	}{
		{
			name:           "InvalidateID",
			invalidate:     func() { wrapper.InvalidateID("123") },
			shouldExist:    []string{"program:list", "program:search:golang"},
			shouldNotExist: []string{"program:123"},
		},
		{
			name:           "InvalidateList",
			invalidate:     func() { wrapper.InvalidateList() },
			shouldExist:    []string{"program:search:golang", "program:category:tech"},
			shouldNotExist: []string{"program:list"},
		},
		{
			name:           "InvalidateSearch",
			invalidate:     func() { wrapper.InvalidateSearch() },
			shouldExist:    []string{"program:category:tech"},
			shouldNotExist: []string{"program:search:golang", "program:search:python"},
		},
		{
			name:           "InvalidateCategory",
			invalidate:     func() { wrapper.InvalidateCategory() },
			shouldExist:    []string{},
			shouldNotExist: []string{"program:category:tech", "program:category:science"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset cache for each test
			cache.Clear()
			cache.Set("program:123", "data1")
			cache.Set("program:list", "data2")
			cache.Set("program:search:golang", "data3")
			cache.Set("program:search:python", "data4")
			cache.Set("program:category:tech", "data5")
			cache.Set("program:category:science", "data6")

			tt.invalidate()

			for _, key := range tt.shouldExist {
				if _, found := cache.Get(key); !found {
					t.Errorf("Expected key %s to exist after %s", key, tt.name)
				}
			}

			for _, key := range tt.shouldNotExist {
				if _, found := cache.Get(key); found {
					t.Errorf("Expected key %s to not exist after %s", key, tt.name)
				}
			}
		})
	}
}

func TestCacheWrapper_InvalidateAll(t *testing.T) {
	cache := NewMemoryCache(time.Minute)
	defer cache.Close()

	wrapper := NewCacheWrapper(cache, "program")

	// Set up test data with different prefixes
	cache.Set("program:123", "data1")
	cache.Set("program:list", "data2")
	cache.Set("user:456", "data3") // Different prefix

	wrapper.InvalidateAll()

	// Program entries should be gone
	if _, found := cache.Get("program:123"); found {
		t.Error("Expected program:123 to be invalidated")
	}
	if _, found := cache.Get("program:list"); found {
		t.Error("Expected program:list to be invalidated")
	}

	// User entry should remain
	if _, found := cache.Get("user:456"); !found {
		t.Error("Expected user:456 to remain")
	}
}

func TestCacheManager(t *testing.T) {
	config := CacheConfig{
		DefaultTTL: 15 * time.Minute,
		ShortTTL:   5 * time.Minute,
		LongTTL:    1 * time.Hour,
	}

	manager := NewCacheManager(config)
	defer manager.Close()

	// Test that all caches are accessible
	if manager.Default() == nil {
		t.Error("Default cache should not be nil")
	}
	if manager.Short() == nil {
		t.Error("Short cache should not be nil")
	}
	if manager.Long() == nil {
		t.Error("Long cache should not be nil")
	}

	// Test TTL values
	if manager.Default().TTL() != config.DefaultTTL {
		t.Errorf("Default cache TTL = %v, want %v", manager.Default().TTL(), config.DefaultTTL)
	}
	if manager.Short().TTL() != config.ShortTTL {
		t.Errorf("Short cache TTL = %v, want %v", manager.Short().TTL(), config.ShortTTL)
	}
	if manager.Long().TTL() != config.LongTTL {
		t.Errorf("Long cache TTL = %v, want %v", manager.Long().TTL(), config.LongTTL)
	}
}

func TestCacheKeyHelpers(t *testing.T) {
	tests := []struct {
		name     string
		function func() string
		expected string
	}{
		{
			name:     "CacheKey",
			function: func() string { return CacheKey("user", "123", "profile") },
			expected: "user:123:profile",
		},
		{
			name:     "ProgramKey",
			function: func() string { return ProgramKey("123") },
			expected: "program:123",
		},
		{
			name:     "ProgramsListKey",
			function: func() string { return ProgramsListKey() },
			expected: "programs:list",
		},
		{
			name:     "ProgramsSearchKey",
			function: func() string { return ProgramsSearchKey("golang") },
			expected: "programs:search:golang",
		},
		{
			name:     "ProgramsCategoryKey",
			function: func() string { return ProgramsCategoryKey("tech") },
			expected: "programs:category:tech",
		},
		{
			name:     "CategoryKey",
			function: func() string { return CategoryKey("123") },
			expected: "category:123",
		},
		{
			name:     "CategoriesListKey",
			function: func() string { return CategoriesListKey() },
			expected: "categories:list",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.function()
			if result != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}

func TestSafeTypeAssert(t *testing.T) {
	// Test successful assertion
	result, err := SafeTypeAssert[string]("hello")
	if err != nil {
		t.Errorf("SafeTypeAssert() error = %v", err)
	}
	if result != "hello" {
		t.Errorf("SafeTypeAssert() = %v, want %v", result, "hello")
	}

	// Test failed assertion
	_, err = SafeTypeAssert[string](123)
	if err == nil {
		t.Error("Expected error for type assertion failure")
	}

	// Test nil value
	_, err = SafeTypeAssert[string](nil)
	if err == nil {
		t.Error("Expected error for nil value")
	}
}

func TestCacheMetrics(t *testing.T) {
	metrics := &CacheMetrics{}

	// Initial state
	if metrics.HitRate() != 0 {
		t.Errorf("Initial hit rate = %v, want 0", metrics.HitRate())
	}
	if metrics.Total() != 0 {
		t.Errorf("Initial total = %v, want 0", metrics.Total())
	}

	// Record some operations
	metrics.RecordHit()
	metrics.RecordHit()
	metrics.RecordMiss()

	if metrics.Hits != 2 {
		t.Errorf("Hits = %v, want 2", metrics.Hits)
	}
	if metrics.Misses != 1 {
		t.Errorf("Misses = %v, want 1", metrics.Misses)
	}
	if metrics.Total() != 3 {
		t.Errorf("Total = %v, want 3", metrics.Total())
	}

	expectedHitRate := float64(2) / float64(3)
	if metrics.HitRate() != expectedHitRate {
		t.Errorf("Hit rate = %v, want %v", metrics.HitRate(), expectedHitRate)
	}

	// Reset
	metrics.Reset()
	if metrics.Hits != 0 || metrics.Misses != 0 {
		t.Error("Reset() should zero all counters")
	}
}

func TestDefaultCacheConfig(t *testing.T) {
	config := DefaultCacheConfig()

	if config.DefaultTTL != 15*time.Minute {
		t.Errorf("DefaultTTL = %v, want %v", config.DefaultTTL, 15*time.Minute)
	}
	if config.ShortTTL != 5*time.Minute {
		t.Errorf("ShortTTL = %v, want %v", config.ShortTTL, 5*time.Minute)
	}
	if config.LongTTL != 1*time.Hour {
		t.Errorf("LongTTL = %v, want %v", config.LongTTL, 1*time.Hour)
	}
}
