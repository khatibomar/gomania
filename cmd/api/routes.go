package main

import (
	"expvar"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)
	mux.Handle("GET /debug/vars", expvar.Handler())

	return app.logRequest(app.recoverPanic(app.enableCORS(mux)))
}
