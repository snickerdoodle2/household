package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (app *App) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", app.config.host, app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     app.logger.StandardLog(),
	}

	sensors, err := app.models.Sensors.GetAll()
	if err != nil {
		return err
	}

	for _, sensor := range sensors {
		listener := app.createAndAddSensorListener(sensor)
		if !sensor.Active {
			go listener.Start()
		}
	}

	rules, err := app.models.Rules.GetAll()
	if err != nil {
		return err
	}

	for _, rule := range rules {
		app.startRule(rule)
	}

	go app.handleRuleRequests()

	app.logger.Info("starting server", "addr", srv.Addr)

	return srv.ListenAndServe()
}

func (app *App) handleRuleRequests() {
	client := &http.Client{}
	client.Timeout = 5 * time.Second
	for message := range app.rules.channel {
		uri, err := app.models.Sensors.GetUri(message.To)
		if err != nil {
			app.logger.Error("handleRuleRequests query", "error", err.Error(), "uuid", message.To)
			continue
		}
		url := fmt.Sprintf("http://%s/value", uri)

		body := new(bytes.Buffer)
		err = json.NewEncoder(body).Encode(message.Payload)
		if err != nil {
			app.logger.Error("handleRuleRequests marshall", "error", err.Error())
			continue
		}

		req, err := http.NewRequest(http.MethodPut, url, body)
		if err != nil {
			app.logger.Error("handleRuleRequests request creation", "error", err.Error())
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		_, err = client.Do(req)
		if err != nil {
			app.logger.Error("handleRuleRequests request", "error", err.Error())
			continue
		}
	}
}
