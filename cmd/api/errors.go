package main

import (
	"errors"
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) badRequestErrorResponse(w http.ResponseWriter, r *http.Request, err error, message string) {
	app.logError(r, err)

	app.errorResponse(w, r, http.StatusBadRequest, message)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.logError(r, errors.New(message))
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, message string) {
	app.logError(r, errors.New(message)) // Log the conflict as an error
	app.errorResponse(w, r, http.StatusConflict, message)
}
