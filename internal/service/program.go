package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/khatibomar/gomania/internal/database"
	"github.com/khatibomar/gomania/internal/sources"
	"github.com/khatibomar/gomania/internal/sources/itunes"
)

type ProgramService struct {
	db             *pgxpool.Pool
	q              *database.Queries
	sourcesManager *sources.Manager
	logger         *slog.Logger
}

type CreateProgramRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Category    string     `json:"category"`
	Language    string     `json:"language"`
	Duration    int        `json:"duration"`
	PublishedAt *time.Time `json:"published_at"`
}

type UpdateProgramRequest struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Category    string     `json:"category"`
	Language    string     `json:"language"`
	Duration    int        `json:"duration"`
	PublishedAt *time.Time `json:"published_at"`
}

type SearchRequest struct {
	Query            string `json:"query"`
	SearchExternal   bool   `json:"search_external"`
	ImportIfNotFound bool   `json:"import_if_not_found"`
}

func NewProgramService(db *pgxpool.Pool, logger *slog.Logger) *ProgramService {
	sourcesManager := sources.NewManager()
	sourcesManager.RegisterClient(itunes.NewClient())

	return &ProgramService{
		db:             db,
		q:              database.New(db),
		sourcesManager: sourcesManager,
		logger:         logger,
	}
}

func (s *ProgramService) CreateProgram(ctx context.Context, req CreateProgramRequest) (*database.Program, error) {
	s.logger.Info("Creating new program", "title", req.Title)

	var publishedAt pgtype.Timestamp
	if req.PublishedAt != nil {
		publishedAt = pgtype.Timestamp{Time: *req.PublishedAt, Valid: true}
	}

	program, err := s.q.CreateProgram(ctx, database.CreateProgramParams{
		Title:       req.Title,
		Description: pgtype.Text{String: req.Description, Valid: req.Description != ""},
		Category:    pgtype.Text{String: req.Category, Valid: req.Category != ""},
		Language:    pgtype.Text{String: req.Language, Valid: req.Language != ""},
		Duration:    pgtype.Int4{Int32: int32(req.Duration), Valid: req.Duration > 0},
		PublishedAt: publishedAt,
		Source:      "local",
	})

	if err != nil {
		s.logger.Error("Failed to create program", "title", req.Title, "error", err)
		return nil, fmt.Errorf("failed to create program: %w", err)
	}

	s.logger.Info("Program created successfully", "title", req.Title, "id", program.ID)
	return &program, nil
}

func (s *ProgramService) GetProgram(ctx context.Context, id uuid.UUID) (*database.Program, error) {
	pgUUID := pgtype.UUID{Bytes: id, Valid: true}
	program, err := s.q.GetProgram(ctx, pgUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get program: %w", err)
	}
	return &program, nil
}

func (s *ProgramService) UpdateProgram(ctx context.Context, req UpdateProgramRequest) (*database.Program, error) {
	var publishedAt pgtype.Timestamp
	if req.PublishedAt != nil {
		publishedAt = pgtype.Timestamp{Time: *req.PublishedAt, Valid: true}
	}

	program, err := s.q.UpdateProgram(ctx, database.UpdateProgramParams{
		ID:          pgtype.UUID{Bytes: req.ID, Valid: true},
		Title:       req.Title,
		Description: pgtype.Text{String: req.Description, Valid: req.Description != ""},
		Category:    pgtype.Text{String: req.Category, Valid: req.Category != ""},
		Language:    pgtype.Text{String: req.Language, Valid: req.Language != ""},
		Duration:    pgtype.Int4{Int32: int32(req.Duration), Valid: req.Duration > 0},
		PublishedAt: publishedAt,
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

func (s *ProgramService) ListPrograms(ctx context.Context) ([]database.Program, error) {
	programs, err := s.q.ListPrograms(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list programs: %w", err)
	}
	return programs, nil
}

func (s *ProgramService) SearchPrograms(ctx context.Context, req SearchRequest) ([]database.Program, error) {
	s.logger.Info("Searching programs", "query", req.Query, "external", req.SearchExternal, "import", req.ImportIfNotFound)

	// First search local database
	localPrograms, err := s.q.SearchPrograms(ctx, pgtype.Text{String: req.Query, Valid: true})
	if err != nil {
		s.logger.Error("Failed to search local programs", "query", req.Query, "error", err)
		return nil, fmt.Errorf("failed to search local programs: %w", err)
	}

	s.logger.Info("Local search completed", "query", req.Query, "found", len(localPrograms))

	// If we have local results or not searching external, return local results
	if len(localPrograms) > 0 || !req.SearchExternal {
		return localPrograms, nil
	}

	s.logger.Info("Searching external sources", "query", req.Query)

	// Search external sources if no local results and external search is enabled
	externalResults, err := s.sourcesManager.SearchAllSources(ctx, req.Query, 10)
	if err != nil {
		s.logger.Error("Failed to search external sources", "query", req.Query, "error", err)
		return nil, fmt.Errorf("failed to search external sources: %w", err)
	}

	// If import is enabled, import the results
	if req.ImportIfNotFound {
		s.logger.Info("Importing external results", "query", req.Query)
		for sourceName, podcasts := range externalResults {
			for _, podcast := range podcasts {
				_, err := s.importFromExternalSource(ctx, podcast, sourceName)
				if err != nil {
					s.logger.Warn("Failed to import podcast", "source", sourceName, "podcast", podcast.Title, "error", err)
					continue
				}
			}
		}

		// Return fresh search results after import
		return s.q.SearchPrograms(ctx, pgtype.Text{String: req.Query, Valid: true})
	}

	return localPrograms, nil
}

func (s *ProgramService) ImportFromExternalSource(ctx context.Context, sourceName, externalID string) (*database.Program, error) {
	client, exists := s.sourcesManager.GetClient(sourceName)
	if !exists {
		return nil, fmt.Errorf("source '%s' not supported", sourceName)
	}

	searchResults, err := client.SearchPodcasts(externalID, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to search %s: %w", sourceName, err)
	}

	if len(searchResults) == 0 {
		return nil, fmt.Errorf("item not found in %s", sourceName)
	}

	result := searchResults[0]
	return s.importFromExternalSource(ctx, result, sourceName)
}

func (s *ProgramService) importFromExternalSource(ctx context.Context, podcast sources.Podcast, sourceName string) (*database.Program, error) {
	s.logger.Info("Importing podcast from external source", "source", sourceName, "title", podcast.Title, "external_id", podcast.ExternalID)

	existingProgram, err := s.q.GetProgramByExternalID(ctx, database.GetProgramByExternalIDParams{
		SourceName: sourceName,
		ExternalID: podcast.ExternalID,
	})
	if err == nil {
		s.logger.Info("Podcast already exists, skipping import", "source", sourceName, "title", podcast.Title, "id", existingProgram.ID)
		return &existingProgram, nil
	}

	var publishedAt pgtype.Timestamp
	if podcast.PublishedAt != nil {
		publishedAt = pgtype.Timestamp{Time: *podcast.PublishedAt, Valid: true}
	}

	// Create program
	program, err := s.q.CreateProgram(ctx, database.CreateProgramParams{
		Title:       podcast.Title,
		Description: pgtype.Text{String: podcast.Description, Valid: podcast.Description != ""},
		Category:    pgtype.Text{String: podcast.Genre, Valid: podcast.Genre != ""},
		Language:    pgtype.Text{String: podcast.Country, Valid: podcast.Country != ""},
		Duration:    pgtype.Int4{Int32: int32(podcast.Duration), Valid: podcast.Duration > 0},
		PublishedAt: publishedAt,
		Source:      sourceName,
	})

	if err != nil {
		s.logger.Error("Failed to create program from external source", "source", sourceName, "title", podcast.Title, "error", err)
		return nil, fmt.Errorf("failed to create program from %s: %w", sourceName, err)
	}

	_, err = s.q.CreateExternalSource(ctx, database.CreateExternalSourceParams{
		ProgramID:  program.ID,
		SourceName: sourceName,
		ExternalID: podcast.ExternalID,
	})

	if err != nil {
		s.logger.Error("Failed to create external source reference", "program_id", program.ID, "source", sourceName, "external_id", podcast.ExternalID, "error", err)
		return nil, fmt.Errorf("failed to create external source reference: %w", err)
	}

	s.logger.Info("Successfully imported podcast", "source", sourceName, "title", podcast.Title, "id", program.ID)
	return &program, nil
}
