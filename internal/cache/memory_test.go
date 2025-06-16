package cache

import (
	"fmt"
	"testing"
	"time"
)

func TestMemoryCache_SetAndGet(t *testing.T) {
	cache := New(time.Minute)
	defer cache.Close()

	// Test setting and getting a value
	cache.Set("key1", "value1")

	value, found := cache.Get("key1")
	if !found {
		t.Error("Expected to find key1")
	}
	if value != "value1" {
		t.Errorf("Expected 'value1', got %v", value)
	}
}

func TestMemoryCache_GetNonExistent(t *testing.T) {
	cache := New(time.Minute)
	defer cache.Close()

	value, found := cache.Get("nonexistent")
	if found {
		t.Error("Expected not to find nonexistent key")
	}
	if value != nil {
		t.Errorf("Expected nil value, got %v", value)
	}
}

func TestMemoryCache_Delete(t *testing.T) {
	cache := New(time.Minute)
	defer cache.Close()

	cache.Set("key1", "value1")
	cache.Delete("key1")

	_, found := cache.Get("key1")
	if found {
		t.Error("Expected key1 to be deleted")
	}
}

func TestMemoryCache_Clear(t *testing.T) {
	cache := New(time.Minute)
	defer cache.Close()

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	// Verify items exist
	_, found1 := cache.Get("key1")
	_, found2 := cache.Get("key2")
	if !found1 || !found2 {
		t.Error("Expected both keys to exist before clear")
	}

	cache.Clear()

	// Verify items are gone
	_, found1 = cache.Get("key1")
	_, found2 = cache.Get("key2")
	if found1 || found2 {
		t.Error("Expected keys to be cleared")
	}
}

func TestMemoryCache_Expiration(t *testing.T) {
	cache := New(100 * time.Millisecond)
	defer cache.Close()

	cache.Set("key1", "value1")

	// Should be available immediately
	_, found := cache.Get("key1")
	if !found {
		t.Error("Expected to find key1 before expiration")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	_, found = cache.Get("key1")
	if found {
		t.Error("Expected key1 to be expired")
	}
}

func TestMemoryCache_InvalidatePattern(t *testing.T) {
	cache := New(time.Minute)
	defer cache.Close()

	cache.Set("user:123", "data1")
	cache.Set("user:456", "data2")
	cache.Set("product:789", "data3")

	cache.InvalidatePattern("user:*")

	_, found := cache.Get("user:123")
	if found {
		t.Error("Expected user:123 to be invalidated")
	}

	_, found = cache.Get("user:456")
	if found {
		t.Error("Expected user:456 to be invalidated")
	}

	_, found = cache.Get("product:789")
	if !found {
		t.Error("Expected product:789 to remain")
	}
}

func TestMemoryCache_ConcurrentAccess(t *testing.T) {
	cache := New(time.Minute)
	defer cache.Close()

	// Test concurrent writes and reads
	done := make(chan bool, 100)

	// Start multiple goroutines writing to cache
	for i := range 50 {
		go func(id int) {
			cache.Set(fmt.Sprintf("key%d", id), fmt.Sprintf("value%d", id))
			done <- true
		}(i)
	}

	// Start multiple goroutines reading from cache
	for i := range 50 {
		go func(id int) {
			cache.Get(fmt.Sprintf("key%d", id))
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for range 100 {
		<-done
	}

	// Verify some data was written by checking a few keys
	foundAny := false
	for i := range 10 {
		if _, found := cache.Get(fmt.Sprintf("key%d", i)); found {
			foundAny = true
			break
		}
	}
	if !foundAny {
		t.Error("Expected cache to contain some data after concurrent writes")
	}
}

func TestMemoryCache_DifferentTypes(t *testing.T) {
	cache := New(time.Minute)
	defer cache.Close()

	// Test different data types
	cache.Set("string", "test")
	cache.Set("int", 42)
	cache.Set("slice", []string{"a", "b", "c"})
	cache.Set("map", map[string]int{"a": 1, "b": 2})

	// Verify string
	value, found := cache.Get("string")
	if !found || value.(string) != "test" {
		t.Error("String value not stored/retrieved correctly")
	}

	// Verify int
	value, found = cache.Get("int")
	if !found || value.(int) != 42 {
		t.Error("Int value not stored/retrieved correctly")
	}

	// Verify slice
	value, found = cache.Get("slice")
	if !found {
		t.Error("Slice value not found")
	}
	slice := value.([]string)
	if len(slice) != 3 || slice[0] != "a" {
		t.Error("Slice value not stored/retrieved correctly")
	}

	// Verify map
	value, found = cache.Get("map")
	if !found {
		t.Error("Map value not found")
	}
	m := value.(map[string]int)
	if m["a"] != 1 || m["b"] != 2 {
		t.Error("Map value not stored/retrieved correctly")
	}
}

// Benchmark tests
func BenchmarkMemoryCache_Set(b *testing.B) {
	cache := New(time.Minute)
	defer cache.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}
}

func BenchmarkMemoryCache_Get(b *testing.B) {
	cache := New(time.Minute)
	defer cache.Close()

	// Pre-populate cache
	for i := range 1000 {
		cache.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(fmt.Sprintf("key%d", i%1000))
	}
}

func BenchmarkMemoryCache_ConcurrentSet(b *testing.B) {
	cache := New(time.Minute)
	defer cache.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
			i++
		}
	})
}

func BenchmarkMemoryCache_ConcurrentGet(b *testing.B) {
	cache := New(time.Minute)
	defer cache.Close()

	// Pre-populate cache
	for i := range 1000 {
		cache.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.Get(fmt.Sprintf("key%d", i%1000))
			i++
		}
	})
}
