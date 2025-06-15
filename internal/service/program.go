package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/khatibomar/gomania/internal/database"
)

type ProgramService struct {
	db     *pgxpool.Pool
	q      *database.Queries
	logger *slog.Logger
}

type CreateProgramRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CategoryID  uuid.UUID `json:"category_id"`
	Language    string    `json:"language"`
	Duration    int       `json:"duration"`
}

type UpdateProgramRequest struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CategoryID  uuid.UUID `json:"category_id"`
	Language    string    `json:"language"`
	Duration    int       `json:"duration"`
}

type SearchRequest struct {
	Query string `json:"query"`
}

type CategoryRequest struct {
	Name string `json:"name"`
}

func NewProgramService(db *pgxpool.Pool, logger *slog.Logger) *ProgramService {
	return &ProgramService{
		db:     db,
		q:      database.New(db),
		logger: logger,
	}
}

func (s *ProgramService) CreateProgram(ctx context.Context, req CreateProgramRequest) (*database.CreateProgramRow, error) {
	s.logger.Info("Creating new program", "title", req.Title)

	program, err := s.q.CreateProgram(ctx, database.CreateProgramParams{
		Title:       req.Title,
		Description: pgtype.Text{String: req.Description, Valid: req.Description != ""},
		CategoryID:  pgtype.UUID{Bytes: req.CategoryID, Valid: true},
		Language:    pgtype.Text{String: req.Language, Valid: req.Language != ""},
		Duration:    pgtype.Int4{Int32: int32(req.Duration), Valid: req.Duration > 0},
	})

	if err != nil {
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
		return nil, fmt.Errorf("failed to get program: %w", err)
	}
	return &program, nil
}

func (s *ProgramService) UpdateProgram(ctx context.Context, req UpdateProgramRequest) (*database.UpdateProgramRow, error) {
	program, err := s.q.UpdateProgram(ctx, database.UpdateProgramParams{
		ID:          pgtype.UUID{Bytes: req.ID, Valid: true},
		Title:       req.Title,
		Description: pgtype.Text{String: req.Description, Valid: req.Description != ""},
		CategoryID:  pgtype.UUID{Bytes: req.CategoryID, Valid: true},
		Language:    pgtype.Text{String: req.Language, Valid: req.Language != ""},
		Duration:    pgtype.Int4{Int32: int32(req.Duration), Valid: req.Duration > 0},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to update program: %w", err)
	}

	return &program, nil
}

func (s *ProgramService) DeleteProgram(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("Deleting program", "id", id)

	pgUUID := pgtype.UUID{Bytes: id, Valid: true}
	err := s.q.DeleteProgram(ctx, pgUUID)
	if err != nil {
		s.logger.Error("Failed to delete program", "id", id, "error", err)
		return fmt.Errorf("failed to delete program: %w", err)
	}

	s.logger.Info("Program deleted successfully", "id", id)
	return nil
}

func (s *ProgramService) ListPrograms(ctx context.Context) ([]database.ListProgramsRow, error) {
	programs, err := s.q.ListPrograms(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list programs: %w", err)
	}
	return programs, nil
}

func (s *ProgramService) SearchPrograms(ctx context.Context, req SearchRequest) ([]database.SearchProgramsRow, error) {
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
	programs, err := s.q.GetProgramsByCategory(ctx, pgUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get programs by category: %w", err)
	}
	return programs, nil
}

// Category management
func (s *ProgramService) CreateCategory(ctx context.Context, req CategoryRequest) (*database.CreateCategoryRow, error) {
	s.logger.Info("Creating new category", "name", req.Name)

	category, err := s.q.CreateCategory(ctx, req.Name)
	if err != nil {
		s.logger.Error("Failed to create category", "name", req.Name, "error", err)
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	s.logger.Info("Category created successfully", "name", req.Name, "id", category.ID)
	return &category, nil
}

func (s *ProgramService) GetCategories(ctx context.Context) ([]database.GetCategoriesRow, error) {
	categories, err := s.q.GetCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	return categories, nil
}
