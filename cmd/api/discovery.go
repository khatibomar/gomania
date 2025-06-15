package main

import (
	"net/http"
	"strconv"

	"github.com/khatibomar/gomania/internal/service"
)

func (app *application) discoveryHandler(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("q")
	if searchQuery != "" {
		app.searchProgramsHandler(w, r, searchQuery)
		return
	}

	app.listProgramsHandler(w, r)
}

func (app *application) searchProgramsHandler(w http.ResponseWriter, r *http.Request, query string) {
	req := service.SearchRequest{
		Query: query,
	}

	programs, err := app.programService.SearchPrograms(r.Context(), req)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response := map[string]any{
		"query":   query,
		"results": programs,
		"count":   len(programs),
		"sources": map[string]any{
			"local": map[string]any{
				"count": len(programs),
			},
		},
	}

	// If no local results found, search external sources
	if len(programs) == 0 {
		app.logger.Info("No local results found, searching external sources", "query", query)

		externalResults, err := app.sourcesManager.SearchAllSources(r.Context(), query, 10)
		if err != nil {
			app.logger.Error("Failed to search external sources", "query", query, "error", err)
		} else {
			response["sources"].(map[string]any)["external"] = externalResults

			// Count total external results
			totalExternal := 0
			for _, podcasts := range externalResults {
				totalExternal += len(podcasts)
			}
			response["external_count"] = totalExternal
		}
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"search": response}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) searchExternalSourcesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		app.badRequestErrorResponse(w, r, nil, "search query is required")
		return
	}

	sourceName := r.URL.Query().Get("source")
	if sourceName == "" {
		app.badRequestErrorResponse(w, r, nil, "source parameter is required")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 10 // default
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	app.logger.Info("Searching external source", "source", sourceName, "query", query, "limit", limit)

	podcasts, err := app.sourcesManager.SearchBySource(r.Context(), sourceName, query, limit)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response := map[string]any{
		"query":   query,
		"source":  sourceName,
		"results": podcasts,
		"count":   len(podcasts),
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"external_search": response}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listExternalSourcesHandler(w http.ResponseWriter, r *http.Request) {
	sources := app.sourcesManager.GetAvailableSources()

	response := map[string]any{
		"sources": sources,
		"count":   len(sources),
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"external_sources": response}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
