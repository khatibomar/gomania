package cache

// Common cache key patterns
const (
	KeyPatternProgramsSearch = "programs:search:*"
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
