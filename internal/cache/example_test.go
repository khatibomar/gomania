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

	// Verify all items are cached
	if _, found := c.Get("user:1"); found {
		fmt.Println("Users cached")
	}
	if _, found := c.Get("product:100"); found {
		fmt.Println("Products cached")
	}

	// Invalidate all user entries
	c.InvalidatePattern("user:*")

	// Check what remains
	if _, found := c.Get("user:1"); !found {
		fmt.Println("Users invalidated")
	}
	if _, found := c.Get("product:100"); found {
		fmt.Println("Products still cached")
	}

	// Output:
	// Users cached
	// Products cached
	// Users invalidated
	// Products still cached
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
