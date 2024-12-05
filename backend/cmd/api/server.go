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
	for message := range app.rules.channel {
		notif := data.Notification{
			Title:       "Rule passed!",
			Description: fmt.Sprintf("Message %v sent to %v", message.Payload, message.To),
			Level:       data.NotificationLevelSuccess,
		}
		app.logger.Debug("handleRuleRequests", "step", "sending notification", "note", "sending notification")
		ids, err := app.models.Notifications.InsertForAll(&notif)
		if err != nil {
			app.logger.Error("handleRuleRequests", "step", "sending notification", "error", err)
		} else {
			app.logger.Debug("handleRuleRequests", "step", "sending notification", "ids", ids)
		}

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

		if err = app.sendValue(url, body); err != nil {
			app.logger.Error("handleRuleRequests request", "url", url, "error", err)
		}
	}
}
