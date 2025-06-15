package main

import (
	"expvar"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// general
	mux.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)
	mux.Handle("GET /debug/vars", expvar.Handler())

	// CMS Programs
	mux.HandleFunc("POST /v1/cms/programs", app.createProgramHandler)
	mux.HandleFunc("GET /v1/cms/programs", app.listProgramsHandler)
	mux.HandleFunc("GET /v1/cms/programs/{id}", app.getProgramHandler)
	mux.HandleFunc("PUT /v1/cms/programs/{id}", app.updateProgramHandler)
	mux.HandleFunc("DELETE /v1/cms/programs/{id}", app.deleteProgramHandler)

	// CMS Categories
	mux.HandleFunc("POST /v1/cms/categories", app.createCategoryHandler)
	mux.HandleFunc("GET /v1/cms/categories", app.listCategoriesHandler)
	mux.HandleFunc("GET /v1/cms/categories/{id}/programs", app.getProgramsByCategoryHandler)

	// discovery
	mux.HandleFunc("GET /v1/programs", app.discoveryHandler)

	return app.logRequest(app.recoverPanic(app.enableCORS(mux)))
}
