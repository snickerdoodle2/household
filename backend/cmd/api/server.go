package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func (app *App) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	sensors, err := app.models.Sensors.GetAll()
	if err != nil {
		return err
	}

	for _, sensor := range sensors {
		app.startSensorListener(sensor)
	}

	app.logger.Info("starting server", "addr", srv.Addr)

	return srv.ListenAndServe()
}
