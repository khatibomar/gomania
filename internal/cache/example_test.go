package cache_test

import (
	"fmt"
	"log"
	"time"

	"github.com/khatibomar/gomania/internal/cache"
)

func Example_basic() {
	// Create a new cache with 5-minute TTL
	c := cache.NewMemoryCache(5 * time.Minute)
	defer c.Close()

	// Store some data
	c.Set("user:123", map[string]any{
		"name":  "John Doe",
		"email": "john@example.com",
		"age":   30,
	})

	// Retrieve data
	if data, found := c.Get("user:123"); found {
		user := data.(map[string]any)
		fmt.Printf("User: %s (%s)\n", user["name"], user["email"])
	}

	// Output: User: John Doe (john@example.com)
}

func Example_patterns() {
	c := cache.NewMemoryCache(time.Hour)
	defer c.Close()

	// Store multiple related items
	c.Set("user:1", "Alice")
	c.Set("user:2", "Bob")
	c.Set("user:3", "Charlie")
	c.Set("product:100", "Laptop")
	c.Set("product:101", "Mouse")

	fmt.Printf("Cache size: %d\n", c.Size())

	// Invalidate all user entries
	c.InvalidatePattern("user:*")

	fmt.Printf("Cache size after invalidating users: %d\n", c.Size())

	// Check remaining items
	if _, found := c.Get("product:100"); found {
		fmt.Println("Product still cached")
	}

	// Output:
	// Cache size: 5
	// Cache size after invalidating users: 2
	// Product still cached
}

func Example_expiration() {
	// Create cache with very short TTL for demonstration
	c := cache.NewMemoryCache(100 * time.Millisecond)
	defer c.Close()

	c.Set("temp_data", "This will expire soon")

	// Check immediately
	if _, found := c.Get("temp_data"); found {
		fmt.Println("Data found immediately")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	if _, found := c.Get("temp_data"); !found {
		fmt.Println("Data expired")
	}

	// Output:
	// Data found immediately
	// Data expired
}

func Example_stats() {
	c := cache.NewMemoryCache(time.Minute)
	defer c.Close()

	// Add some data
	for i := range 5 {
		c.Set(fmt.Sprintf("key:%d", i), fmt.Sprintf("value:%d", i))
	}

	// Get statistics
	stats := cache.Stats{
		TotalItems: c.Size(),
		TTL:        c.TTL(),
	}

	fmt.Printf("Items in cache: %d\n", stats.TotalItems)
	fmt.Printf("TTL: %v\n", stats.TTL)

	// Output:
	// Items in cache: 5
	// TTL: 1m0s
}

func Example_programService() {
	// This example shows how the cache integrates with a service
	c := cache.NewMemoryCache(15 * time.Minute)
	defer c.Close()

	// Simulate service layer usage
	getProgramFromCache := func(id string) (any, bool) {
		cacheKey := fmt.Sprintf("program:%s", id)
		return c.Get(cacheKey)
	}

	setProgramInCache := func(id string, program any) {
		cacheKey := fmt.Sprintf("program:%s", id)
		c.Set(cacheKey, program)
	}

	invalidateProgramCache := func() {
		c.InvalidatePattern("program:*")
		c.InvalidatePattern("programs:*")
		log.Println("Program cache invalidated")
	}

	// Usage
	programID := "123"
	program := map[string]string{
		"title":    "Go Programming",
		"language": "Go",
	}

	// Cache miss - would typically fetch from database
	if _, found := getProgramFromCache(programID); !found {
		fmt.Println("Cache miss - fetching from database")
		setProgramInCache(programID, program)
	}

	// Cache hit
	if cached, found := getProgramFromCache(programID); found {
		fmt.Println("Cache hit - returning cached data")
		_ = cached // Use cached data
	}

	// Invalidate on update
	invalidateProgramCache()

	// Output:
	// Cache miss - fetching from database
	// Cache hit - returning cached data
}
