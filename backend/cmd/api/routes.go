package main

import (
	"inzynierka/internal/data"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *App) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(app.recoverPanic, cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: app.config.cors.trustedOrigins,
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}), app.authenticate)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/healthcheck", app.healthcheckHandler)

		r.Get("/sensor/measurements", app.upgradeSensorWebsocket)
		r.Post("/sensor/measurements", app.activeSensorHandler)

		r.Post("/sensor/init-ack", app.initAckHandler)
		r.Post("/sensor/re-init/{id}", app.reInitSensorHandler)

		r.Route("/", func(r chi.Router) {
			r.Use(app.requireAuthenticated)

			r.Get("/sensor", app.listSensorsHandler)
			r.Get("/sensor/{id}", app.getSensorHandler)
			r.Put("/sensor/{id}/value", app.setSensorValue)

			r.Post("/sensor", app.requireRole(data.UserRoleAdmin, http.HandlerFunc(app.createSensorHandler)))
			r.Put("/sensor/{id}", app.requireRole(data.UserRoleAdmin, http.HandlerFunc(app.updateSensorHandler)))
			r.Delete("/sensor/{id}", app.requireRole(data.UserRoleAdmin, http.HandlerFunc(app.deleteSensorHandler)))

			r.Get("/rule", app.listRulesHandler)
			r.Get("/rule/{id}", app.getRuleHandler)

			r.Post("/rule", app.requireRole(data.UserRoleAdmin, http.HandlerFunc(app.createRuleHandler)))
			r.Put("/rule/{id}", app.requireRole(data.UserRoleAdmin, http.HandlerFunc(app.updateRuleHanlder)))
			r.Delete("/rule/{id}", app.requireRole(data.UserRoleAdmin, http.HandlerFunc(app.deleteRuleHandler)))

			r.Get("/sequence", app.listSequencesHandler)
			r.Get("/sequence/{id}", app.getSequenceHandler)
			r.Post("/sequence/{id}/start", app.startSequenceHandler)

			r.Post("/sequence", app.requireRole(data.UserRoleAdmin, http.HandlerFunc(app.createSequenceHandler)))
			r.Put("/sequence/{id}", app.requireRole(data.UserRoleAdmin, http.HandlerFunc(app.updateSequenceHandler)))
			r.Delete("/sequence/{id}", app.requireRole(data.UserRoleAdmin, http.HandlerFunc(app.deleteSequenceHandler)))

			r.Put("/notification/{id}", app.readNotificationHandler)
			r.Put("/notification", app.readAllNotificationHandler)
			r.Post("/notification/debug", app.requestAllNotifsHandler)

			// TODO: make sure only person who can change user data is THE user (or admin)
			r.Get("/user/me", app.getCurrentUserHandler)

			r.Get("/user", app.getAllUsersHandler)
			r.Get("/user/{username}", app.getUserHandler)
			r.Post("/user", app.requireRole(data.UserRoleAdmin, http.HandlerFunc(app.createUserHandler)))
			r.Put("/user/{username}", app.updateUserHandler)
			r.Delete("/user/{username}", app.requireRole(data.UserRoleAdmin, http.HandlerFunc(app.deleteUserHandler)))

		})
		r.Post("/login", app.loginHandler)
		r.Post("/logout", app.logoutHandler)

		r.NotFound(app.notFoundResponse)
		r.MethodNotAllowed(app.methodNotAllowed)
	})

	r.NotFound(app.spaHandler)

	return r
}
