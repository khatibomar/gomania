package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
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

type ProgramService struct {
	db        *pgxpool.Pool
	q         *database.Queries
	logger    *slog.Logger
	validator *validator.Validate
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
	return &ProgramService{
		db:        db,
		q:         database.New(db),
		logger:    logger,
		validator: validator.New(),
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
			// Assuming the unique constraint is on the title or a combination involving it.
			s.logger.Warn("Attempted to create a program that already exists", "title", req.Title, "error", err)
			return nil, &ErrAlreadyExists{Message: fmt.Sprintf("program with title '%s' already exists or conflicts with an existing one", req.Title)}
		}
		s.logger.Error("Failed to create program", "title", req.Title, "error", err)
		return nil, fmt.Errorf("failed to create program: %w", err)
	}

	s.logger.Info("Program created successfully", "title", req.Title, "id", program.ID)
	return &program, nil
}

func (s *ProgramService) GetProgram(ctx context.Context, id uuid.UUID) (*database.GetProgramRow, error) {
	pgUUID := pgtype.UUID{Bytes: id, Valid: true}
	program, err := s.q.GetProgram(ctx, pgUUID)
	if err != nil {
		// Consider adding specific error handling for "not found" if pgx returns a distinct error for it.
		// For now, a generic error is returned.
		if errors.Is(err, pgx.ErrNoRows) { // Assuming pgx.ErrNoRows is the correct error for not found
			s.logger.Info("Program not found", "id", id)
			return nil, fmt.Errorf("program with ID '%s' not found", id) // Or a custom ErrNotFound
		}
		s.logger.Error("Failed to get program", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get program: %w", err)
	}
	return &program, nil
}

func (s *ProgramService) UpdateProgram(ctx context.Context, req UpdateProgramRequest) (*database.UpdateProgramRow, error) {
	if err := s.validator.Struct(req); err != nil {
		s.logger.Error("Invalid update program request", "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	s.logger.Info("Updating program", "id", req.ID, "title", req.Title)

	program, err := s.q.UpdateProgram(ctx, database.UpdateProgramParams{
		ID:          pgtype.UUID{Bytes: req.ID, Valid: true},
		Title:       req.Title,
		Description: pgtype.Text{String: req.Description, Valid: req.Description != ""},
		CategoryID:  pgtype.UUID{Bytes: req.CategoryID, Valid: true},
		Language:    pgtype.Text{String: req.Language, Valid: req.Language != ""},
		Duration:    pgtype.Int4{Int32: int32(req.Duration), Valid: req.Duration > 0},
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation, e.g., if title is updated to an existing one
			s.logger.Warn("Attempted to update program to a conflicting state", "id", req.ID, "title", req.Title, "error", err)
			return nil, &ErrAlreadyExists{Message: fmt.Sprintf("cannot update program, title '%s' may already exist or conflict", req.Title)}
		}
		// Check if the error is due to the program not being found for update
		// This depends on how your DB/sqlc query handles updates on non-existent rows.
		// If it returns pgx.ErrNoRows or similar, you can handle it here.
		// For example:
		// if errors.Is(err, pgx.ErrNoRows) {
		// 	s.logger.Info("Program not found for update", "id", req.ID)
		// 	return nil, fmt.Errorf("program with ID '%s' not found for update", req.ID) // Or a custom ErrNotFound
		// }
		s.logger.Error("Failed to update program", "id", req.ID, "error", err)
		return nil, fmt.Errorf("failed to update program: %w", err)
	}

	s.logger.Info("Program updated successfully", "id", program.ID, "title", program.Title)
	return &program, nil
}

func (s *ProgramService) DeleteProgram(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("Deleting program", "id", id)

	pgUUID := pgtype.UUID{Bytes: id, Valid: true}
	err := s.q.DeleteProgram(ctx, pgUUID)
	if err != nil {
		// Consider checking if the error is because the program was not found,
		// if your DeleteProgram query or DB behavior allows distinguishing this.
		// For example, if DeleteProgram returns an error or a specific result when no rows are affected.
		s.logger.Error("Failed to delete program", "id", id, "error", err)
		return fmt.Errorf("failed to delete program: %w", err)
	}

	s.logger.Info("Program deleted successfully", "id", id)
	return nil
}

func (s *ProgramService) ListPrograms(ctx context.Context) ([]database.ListProgramsRow, error) {
	s.logger.Info("Listing all programs")
	programs, err := s.q.ListPrograms(ctx)
	if err != nil {
		s.logger.Error("Failed to list programs", "error", err)
		return nil, fmt.Errorf("failed to list programs: %w", err)
	}
	s.logger.Info("Successfully listed programs", "count", len(programs))
	return programs, nil
}

func (s *ProgramService) SearchPrograms(ctx context.Context, req SearchRequest) ([]database.SearchProgramsRow, error) {
	if err := s.validator.Struct(req); err != nil {
		s.logger.Error("Invalid search request", "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	s.logger.Info("Searching programs", "query", req.Query)

	programs, err := s.q.SearchPrograms(ctx, pgtype.Text{String: req.Query, Valid: true})
	if err != nil {
		s.logger.Error("Failed to search programs", "query", req.Query, "error", err)
		return nil, fmt.Errorf("failed to search programs: %w", err)
	}

	s.logger.Info("Search completed", "query", req.Query, "found", len(programs))
	return programs, nil
}

func (s *ProgramService) GetProgramsByCategory(ctx context.Context, categoryID uuid.UUID) ([]database.GetProgramsByCategoryRow, error) {
	pgUUID := pgtype.UUID{Bytes: categoryID, Valid: true}
	s.logger.Info("Getting programs by category", "category_id", categoryID)
	programs, err := s.q.GetProgramsByCategory(ctx, pgUUID)
	if err != nil {
		s.logger.Error("Failed to get programs by category", "category_id", categoryID, "error", err)
		return nil, fmt.Errorf("failed to get programs by category: %w", err)
	}
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
			// Assuming the unique constraint is on the name.
			s.logger.Warn("Attempted to create a category that already exists", "name", req.Name, "error", err)
			return nil, &ErrAlreadyExists{Message: fmt.Sprintf("category with name '%s' already exists", req.Name)}
		}
		s.logger.Error("Failed to create category", "name", req.Name, "error", err)
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	s.logger.Info("Category created successfully", "name", req.Name, "id", category.ID)
	return &category, nil
}

func (s *ProgramService) GetCategories(ctx context.Context) ([]database.GetCategoriesRow, error) {
	s.logger.Info("Getting all categories")
	categories, err := s.q.GetCategories(ctx)
	if err != nil {
		s.logger.Error("Failed to get categories", "error", err)
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	s.logger.Info("Successfully fetched categories", "count", len(categories))
	return categories, nil
}
