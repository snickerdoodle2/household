package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *App) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(app.logRequest, app.recoverPanic, cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: app.config.cors.trustedOrigins,
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/healthcheck", app.healthcheckHandler)

		r.Get("/sensor", app.listSensorsHandler)
		r.Get("/sensor/{id}", app.getSensorHandler)
		r.Get("/sensor/{id}/value", app.getSensorValueHandler)
		r.Post("/sensor", app.createSensorHandler)
		r.Put("/sensor/{id}", app.updateSensorHandler)
		r.Delete("/sensor/{id}", app.deleteSensorHandler)

		r.NotFound(app.notFoundResponse)
		r.MethodNotAllowed(app.methodNotAllowed)
	})

	r.NotFound(app.spaHandler)

	return r
}
