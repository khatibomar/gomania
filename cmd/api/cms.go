package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/khatibomar/gomania/internal/service"
)

func (app *application) createProgramHandler(w http.ResponseWriter, r *http.Request) {
	var req service.CreateProgramRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	program, err := app.programService.CreateProgram(r.Context(), req)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusCreated, envelope{"program": program}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listProgramsHandler(w http.ResponseWriter, r *http.Request) {
	programs, err := app.programService.ListPrograms(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"programs": programs}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getProgramHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	program, err := app.programService.GetProgram(r.Context(), id)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"program": program}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateProgramHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var req service.UpdateProgramRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	req.ID = id

	program, err := app.programService.UpdateProgram(r.Context(), req)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"program": program}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteProgramHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.programService.DeleteProgram(r.Context(), id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

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

// Category handlers
func (app *application) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var req service.CategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	category, err := app.programService.CreateCategory(r.Context(), req)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusCreated, envelope{"category": category}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := app.programService.GetCategories(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"categories": categories}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getProgramsByCategoryHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	programs, err := app.programService.GetProgramsByCategory(r.Context(), id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"programs": programs}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
