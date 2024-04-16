package main

import (
	"fmt"
	"inzynierka/internal/listener"
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
		switch {
		case sensor.Type == "binary_switch" || sensor.Type == "binary_sensor" || sensor.Type == "button":
			l := listener.New[bool](sensor)
			go l.Start()
			app.listeners[sensor.ID] = l
		case sensor.Type == "decimal_switch" || sensor.Type == "decimal_sensor":
			l := listener.New[float64](sensor)
			go l.Start()
			app.listeners[sensor.ID] = l
		}
	}

	app.logger.Info("starting server", "addr", srv.Addr)

	return srv.ListenAndServe()
}
