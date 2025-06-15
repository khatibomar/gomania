package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/khatibomar/gomania/internal/cache"
	"github.com/khatibomar/gomania/internal/database"
)

// ErrAlreadyExists is returned when trying to create a resource that already exists.
type ErrAlreadyExists struct {
	Message string
}

func (e *ErrAlreadyExists) Error() string {
	return e.Message
}

// IsErrAlreadyExists checks if an error is of type ErrAlreadyExists.
func IsErrAlreadyExists(err error) bool {
	var target *ErrAlreadyExists
	return errors.As(err, &target)
}

// ErrNotFound is returned when a resource is not found.
var ErrNotFound = errors.New("resource not found")

type ProgramService struct {
	db        *pgxpool.Pool
	q         *database.Queries
	logger    *slog.Logger
	validator *validator.Validate
	cache     cache.Cache
}

type CreateProgramRequest struct {
	Title       string    `json:"title" validate:"required,min=3,max=100"`
	Description string    `json:"description" validate:"omitempty,max=1000"`
	CategoryID  uuid.UUID `json:"category_id" validate:"required"`
	Language    string    `json:"language" validate:"required,min=2,max=50"`
	Duration    int       `json:"duration" validate:"required,gt=0"`
}

type UpdateProgramRequest struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Title       string    `json:"title" validate:"required,min=3,max=100"`
	Description string    `json:"description" validate:"omitempty,max=1000"`
	CategoryID  uuid.UUID `json:"category_id" validate:"required"`
	Language    string    `json:"language" validate:"required,min=2,max=50"`
	Duration    int       `json:"duration" validate:"required,gt=0"`
}

type SearchRequest struct {
	Query string `json:"query" validate:"required,min=1,max=100"`
}

type CategoryRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
}

func NewProgramService(db *pgxpool.Pool, logger *slog.Logger) *ProgramService {
	return NewProgramServiceWithCache(db, logger, 15*time.Minute)
}

func NewProgramServiceWithCache(db *pgxpool.Pool, logger *slog.Logger, cacheTTL time.Duration) *ProgramService {
	return &ProgramService{
		db:        db,
		q:         database.New(db),
		logger:    logger,
		validator: validator.New(),
		cache:     cache.NewMemoryCache(cacheTTL),
	}
}

func (s *ProgramService) CreateProgram(ctx context.Context, req CreateProgramRequest) (*database.CreateProgramRow, error) {
	if err := s.validator.Struct(req); err != nil {
		s.logger.Error("Invalid create program request", "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	s.logger.Info("Creating new program", "title", req.Title)

	program, err := s.q.CreateProgram(ctx, database.CreateProgramParams{
		Title:       req.Title,
		Description: pgtype.Text{String: req.Description, Valid: req.Description != ""},
		CategoryID:  pgtype.UUID{Bytes: req.CategoryID, Valid: true},
		Language:    pgtype.Text{String: req.Language, Valid: req.Language != ""},
		Duration:    pgtype.Int4{Int32: int32(req.Duration), Valid: req.Duration > 0},
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // 23505 is unique_violation
			s.logger.Warn("Attempted to create a program that already exists", "title", req.Title, "error", err)
			return nil, &ErrAlreadyExists{Message: fmt.Sprintf("program with title '%s' already exists or conflicts with an existing one", req.Title)}
		}
		s.logger.Error("Failed to create program", "title", req.Title, "error", err)
		return nil, fmt.Errorf("failed to create program: %w", err)
	}

	// Invalidate relevant cache entries after successful creation
	s.cache.Delete(cache.ProgramsListKey())
	s.cache.InvalidatePattern(cache.KeyPatternProgramsSearch)          // Clears all search results
	s.cache.Delete(cache.ProgramsCategoryKey(req.CategoryID.String())) // Clears specific category list

	s.logger.Info("Program created successfully", "title", req.Title, "id", program.ID)
	return &program, nil
}

func (s *ProgramService) GetProgram(ctx context.Context, id uuid.UUID) (*database.GetProgramRow, error) {
	cacheKey := cache.ProgramKey(id.String())

	if cached, found := s.cache.Get(cacheKey); found {
		s.logger.Debug("Program found in cache", "id", id)
		if program, ok := cached.(*database.GetProgramRow); ok {
			return program, nil
		}
		// Invalid type in cache, remove it
		s.logger.Warn("Invalid cache entry type for program, removing", "id", id)
		s.cache.Delete(cacheKey)
	}

	pgUUID := pgtype.UUID{Bytes: id, Valid: true}
	program, err := s.q.GetProgram(ctx, pgUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.Info("Program not found in DB", "id", id)
			return nil, fmt.Errorf("%w: program with ID '%s' not found", ErrNotFound, id.String())
		}
		s.logger.Error("Failed to get program from DB", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get program: %w", err)
	}

	s.cache.Set(cacheKey, &program)
	s.logger.Debug("Program cached", "id", id)

	return &program, nil
}

func (s *ProgramService) UpdateProgram(ctx context.Context, req UpdateProgramRequest) (*database.UpdateProgramRow, error) {
	if err := s.validator.Struct(req); err != nil {
		s.logger.Error("Invalid update program request", "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get the program's current state, especially its CategoryID, before updating.
	// This uses the cached GetProgram method.
	existingProgram, err := s.GetProgram(ctx, req.ID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			s.logger.Info("Program not found for update", "id", req.ID)
			// Return a specific error indicating the resource to update was not found
			return nil, fmt.Errorf("%w: program with ID '%s' not found for update", ErrNotFound, req.ID.String())
		}
		s.logger.Error("Failed to get program details before update", "id", req.ID, "error", err)
		return nil, fmt.Errorf("failed to get program details before update: %w", err)
	}
	oldProgramCategoryID := existingProgram.CategoryID // This is pgtype.UUID

	s.logger.Info("Updating program", "id", req.ID, "title", req.Title)

	updatedProgramData, err := s.q.UpdateProgram(ctx, database.UpdateProgramParams{
		ID:          pgtype.UUID{Bytes: req.ID, Valid: true},
		Title:       req.Title,
		Description: pgtype.Text{String: req.Description, Valid: req.Description != ""},
		CategoryID:  pgtype.UUID{Bytes: req.CategoryID, Valid: true},
		Language:    pgtype.Text{String: req.Language, Valid: req.Language != ""},
		Duration:    pgtype.Int4{Int32: int32(req.Duration), Valid: req.Duration > 0},
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			s.logger.Warn("Attempted to update program to a conflicting state", "id", req.ID, "title", req.Title, "error", err)
			return nil, &ErrAlreadyExists{Message: fmt.Sprintf("cannot update program, title '%s' may already exist or conflict", req.Title)}
		}
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.Info("Program not found during DB update operation (race condition or already deleted)", "id", req.ID)
			return nil, fmt.Errorf("%w: program with ID '%s' not found for update", ErrNotFound, req.ID.String())
		}
		s.logger.Error("Failed to update program in DB", "id", req.ID, "error", err)
		return nil, fmt.Errorf("failed to update program: %w", err)
	}

	// Invalidate relevant cache entries after successful update
	s.cache.Delete(cache.ProgramKey(req.ID.String()))         // Specific program cache
	s.cache.Delete(cache.ProgramsListKey())                   // General list of programs
	s.cache.InvalidatePattern(cache.KeyPatternProgramsSearch) // All search results

	// Invalidate cache for the new/current category
	s.cache.Delete(cache.ProgramsCategoryKey(req.CategoryID.String()))

	// If the category ID changed, invalidate the cache for the old category as well
	if oldProgramCategoryID.Valid {
		oldCatUUID, convErr := uuid.FromBytes(oldProgramCategoryID.Bytes[:])
		if convErr == nil {
			if oldCatUUID != req.CategoryID {
				s.logger.Debug("Program category changed, invalidating old category cache", "old_category_id", oldCatUUID.String())
				s.cache.Delete(cache.ProgramsCategoryKey(oldCatUUID.String()))
			}
		} else {
			s.logger.Error("Failed to convert old category pgtype.UUID to uuid.UUID for comparison", "error", convErr)
		}
	}

	s.logger.Info("Program updated successfully", "id", updatedProgramData.ID, "title", updatedProgramData.Title)
	return &updatedProgramData, nil
}

func (s *ProgramService) DeleteProgram(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("Deleting program", "id", id)

	var categoryIDToDeleteFromCache pgtype.UUID
	// Fetch program details to get CategoryID for targeted cache invalidation.
	programToDelete, err := s.GetProgram(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			s.logger.Info("Program not found before deletion attempt (already deleted or never existed)", "id", id)
			// Program is already gone, or never existed. DB delete will confirm.
			// categoryIDToDeleteFromCache will remain invalid/zero.
		} else {
			// For other errors fetching program details, return the error.
			s.logger.Error("Failed to get program details before deletion", "id", id, "error", err)
			return fmt.Errorf("failed to get program details before deletion: %w", err)
		}
	} else if programToDelete != nil {
		categoryIDToDeleteFromCache = programToDelete.CategoryID
	}

	pgUUID := pgtype.UUID{Bytes: id, Valid: true}
	dbErr := s.q.DeleteProgram(ctx, pgUUID)
	if dbErr != nil {
		if errors.Is(dbErr, pgx.ErrNoRows) {
			s.logger.Info("Program not found in DB during deletion (confirmed already deleted)", "id", id)
			// Fall through to cache invalidation, as the state is "item does not exist".
		} else {
			s.logger.Error("Failed to delete program from DB", "id", id, "error", dbErr)
			return fmt.Errorf("failed to delete program: %w", dbErr)
		}
	}

	// Invalidate relevant cache entries
	s.cache.Delete(cache.ProgramKey(id.String()))
	s.cache.Delete(cache.ProgramsListKey())
	s.cache.InvalidatePattern(cache.KeyPatternProgramsSearch)

	if categoryIDToDeleteFromCache.Valid {
		categoryUUID, convErr := uuid.FromBytes(categoryIDToDeleteFromCache.Bytes[:])
		if convErr == nil {
			s.logger.Debug("Invalidating category cache for deleted program", "category_id", categoryUUID.String())
			s.cache.Delete(cache.ProgramsCategoryKey(categoryUUID.String()))
		} else {
			s.logger.Error("Failed to convert category pgtype.UUID to uuid.UUID for cache key on delete", "error", convErr)
		}
	}

	s.logger.Info("Program deleted successfully", "id", id)
	return nil
}

func (s *ProgramService) ListPrograms(ctx context.Context) ([]database.ListProgramsRow, error) {
	cacheKey := cache.ProgramsListKey()

	if cached, found := s.cache.Get(cacheKey); found {
		s.logger.Debug("Programs list found in cache")
		if programs, ok := cached.([]database.ListProgramsRow); ok {
			return programs, nil
		}
		// Invalid type in cache, remove it
		s.logger.Warn("Invalid cache entry type for programs list, removing")
		s.cache.Delete(cacheKey)
	}

	s.logger.Info("Listing all programs")
	programs, err := s.q.ListPrograms(ctx)
	if err != nil {
		s.logger.Error("Failed to list programs", "error", err)
		return nil, fmt.Errorf("failed to list programs: %w", err)
	}

	s.cache.Set(cacheKey, programs)
	s.logger.Debug("Programs list cached")

	s.logger.Info("Successfully listed programs", "count", len(programs))
	return programs, nil
}

func (s *ProgramService) SearchPrograms(ctx context.Context, req SearchRequest) ([]database.SearchProgramsRow, error) {
	if err := s.validator.Struct(req); err != nil {
		s.logger.Error("Invalid search request", "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	cacheKey := cache.ProgramsSearchKey(req.Query)

	if cached, found := s.cache.Get(cacheKey); found {
		s.logger.Debug("Search results found in cache", "query", req.Query)
		if programs, ok := cached.([]database.SearchProgramsRow); ok {
			return programs, nil
		}
		// Invalid type in cache, remove it
		s.logger.Warn("Invalid cache entry type for search results, removing", "query", req.Query)
		s.cache.Delete(cacheKey)
	}

	s.logger.Info("Searching programs", "query", req.Query)

	programs, err := s.q.SearchPrograms(ctx, pgtype.Text{String: req.Query, Valid: true})
	if err != nil {
		s.logger.Error("Failed to search programs", "query", req.Query, "error", err)
		return nil, fmt.Errorf("failed to search programs: %w", err)
	}

	s.cache.Set(cacheKey, programs)
	s.logger.Debug("Search results cached", "query", req.Query)

	s.logger.Info("Search completed", "query", req.Query, "found", len(programs))
	return programs, nil
}

func (s *ProgramService) GetProgramsByCategory(ctx context.Context, categoryID uuid.UUID) ([]database.GetProgramsByCategoryRow, error) {
	cacheKey := cache.ProgramsCategoryKey(categoryID.String())

	if cached, found := s.cache.Get(cacheKey); found {
		s.logger.Debug("Programs by category found in cache", "category_id", categoryID)
		if programs, ok := cached.([]database.GetProgramsByCategoryRow); ok {
			return programs, nil
		}
		// Invalid type in cache, remove it
		s.logger.Warn("Invalid cache entry type for programs by category, removing", "category_id", categoryID)
		s.cache.Delete(cacheKey)
	}

	pgUUID := pgtype.UUID{Bytes: categoryID, Valid: true}
	s.logger.Info("Getting programs by category", "category_id", categoryID)
	programs, err := s.q.GetProgramsByCategory(ctx, pgUUID)
	if err != nil {
		s.logger.Error("Failed to get programs by category", "category_id", categoryID, "error", err)
		return nil, fmt.Errorf("failed to get programs by category: %w", err)
	}

	s.cache.Set(cacheKey, programs)
	s.logger.Debug("Programs by category cached", "category_id", categoryID)

	s.logger.Info("Successfully fetched programs by category", "category_id", categoryID, "count", len(programs))
	return programs, nil
}

// Category management
func (s *ProgramService) CreateCategory(ctx context.Context, req CategoryRequest) (*database.CreateCategoryRow, error) {
	if err := s.validator.Struct(req); err != nil {
		s.logger.Error("Invalid category request", "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	s.logger.Info("Creating new category", "name", req.Name)

	category, err := s.q.CreateCategory(ctx, req.Name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // 23505 is unique_violation
			s.logger.Warn("Attempted to create a category that already exists", "name", req.Name, "error", err)
			return nil, &ErrAlreadyExists{Message: fmt.Sprintf("category with name '%s' already exists", req.Name)}
		}
		s.logger.Error("Failed to create category", "name", req.Name, "error", err)
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	s.cache.Delete(cache.CategoriesListKey())

	s.logger.Info("Category created successfully", "name", req.Name, "id", category.ID)
	return &category, nil
}

func (s *ProgramService) GetCategories(ctx context.Context) ([]database.GetCategoriesRow, error) {
	cacheKey := cache.CategoriesListKey()

	if cached, found := s.cache.Get(cacheKey); found {
		s.logger.Debug("Categories found in cache")
		if categories, ok := cached.([]database.GetCategoriesRow); ok {
			return categories, nil
		}
		// Invalid type in cache, remove it
		s.logger.Warn("Invalid cache entry type for categories list, removing")
		s.cache.Delete(cacheKey)
	}

	s.logger.Info("Getting all categories")
	categories, err := s.q.GetCategories(ctx)
	if err != nil {
		s.logger.Error("Failed to get categories", "error", err)
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	s.cache.Set(cacheKey, categories)
	s.logger.Debug("Categories cached")

	s.logger.Info("Successfully fetched categories", "count", len(categories))
	return categories, nil
}

// GetCacheStats returns current cache statistics
func (s *ProgramService) GetCacheStats() cache.Stats {
	hits, misses, evicted := s.cache.Metrics()
	return cache.Stats{
		TotalItems: s.cache.Size(),
		TTL:        s.cache.TTL(),
		Hits:       hits,
		Misses:     misses,
		Evicted:    evicted,
	}
}

// ClearCache clears all cache entries
func (s *ProgramService) ClearCache() {
	s.logger.Info("Clearing all cache entries")
	s.cache.Clear()
}

// InvalidateProgramCache invalidates all program-related cache entries
// This includes individual programs, lists of programs, search results, and programs-by-category lists.
func (s *ProgramService) InvalidateProgramCache() {
	s.logger.Info("Invalidating all program-related cache")
	s.cache.InvalidatePattern(cache.KeyPatternPrograms) // e.g., "program:*" for individual items
	s.cache.InvalidatePattern("programs:*")             // e.g., "programs:list", "programs:search:*", "programs:category:*"
}

// InvalidateCategoryCache invalidates all category-entity-related cache entries (e.g., list of all categories)
// This does NOT invalidate program lists grouped by category. For that, use InvalidateProgramCache or more targeted invalidations.
func (s *ProgramService) InvalidateCategoryCache() {
	s.logger.Info("Invalidating category entity cache")
	s.cache.InvalidatePattern(cache.KeyPatternCategories) // e.g., "categories:*"
}

// SetCacheTTL updates the cache TTL (affects new entries only)
func (s *ProgramService) SetCacheTTL(ttl time.Duration) {
	s.cache.SetTTL(ttl)
	s.logger.Info("Cache TTL updated", "new_ttl", ttl)
}
