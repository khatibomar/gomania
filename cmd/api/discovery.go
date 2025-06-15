package main

import (
	"net/http"

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
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"search": response}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
