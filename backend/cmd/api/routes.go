package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *App) routes() http.Handler {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/healthcheck", app.healthcheckHandler)

		r.Get("/sensor", app.listSensorsHandler)
		r.Get("/sensor/{id}", app.getSensorHandler)
		r.Post("/sensor", app.createSensorHandler)

		r.NotFound(app.notFoundResponse)
		r.MethodNotAllowed(app.methodNotAllowed)
	})

	return r
}
