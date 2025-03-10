package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"inzynierka/internal/data"
	"net/http"
	"time"
)

func (app *App) parseSettings() error {
	// TODO: ADD PARSING
	app.settings.MeasurementsAmount = 32
	return nil
}

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
		app.setupSensorListener(sensor)
	}

	go app.notificationBroker.Start()

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
	// reading from channel and handling rule requests
	for message := range app.rules.channel {
		switch message.TargetType {
		case data.SensorTarget:
			_ = app.sendNotificationToAll("Rule passed!", fmt.Sprintf("Sent message %v to sensor %v", message.Payload, message.TargetId), data.NotificationLevelSuccess)

			uri, err := app.models.Sensors.GetUri(message.TargetId)
			if err != nil {
				app.logger.Error("handleRuleRequests query", "error", err.Error(), "uuid", message.TargetId)
				continue
			}

			url := fmt.Sprintf("http://%s/value", uri)

			body := new(bytes.Buffer)
			err = json.NewEncoder(body).Encode(message.Payload)
			if err != nil {
				app.logger.Error("handleRuleRequests marshall", "error", err.Error())
				continue
			}

			if err = app.sendValue(url, body); err != nil {
				app.logger.Error("handleRuleRequests request", "url", url, "error", err)
			}

		case data.SequenceTarget:
			_ = app.sendNotificationToAll("Rule passed!", fmt.Sprintf("Starting sequence: %v", message.TargetId), data.NotificationLevelSuccess)

			sequence, err := app.models.Sequences.Get(message.TargetId)
			if err != nil {
				app.logger.Error("handleRuleRequests query", "error", err.Error(), "uuid", message.TargetId)
				continue
			}

			preparedData, err := app.prepareActionData(sequence.Actions)
			if err != nil {
				app.logger.Error("handleRuleRequests prepareActionData", "error", err.Error())
				continue
			}
			go app.executeSequence(preparedData)
		}
	}
}
