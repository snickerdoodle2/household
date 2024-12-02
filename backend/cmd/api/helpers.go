package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"inzynierka/internal/data"
	"inzynierka/internal/data/validator"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type envelope map[string]any

func (app *App) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *App) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func (app *App) readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

func (app *App) readCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)

	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}

func (app *App) readInt(qs url.Values, key string, defaulValue int, v *validator.Validator) int {
	s := qs.Get(key)

	if s == "" {
		return defaulValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be valid integer")
		return defaulValue
	}

	return i
}

// function creates a new listener, adds it to app module and depending on sensor active flag starts it or just starts broker
func (app *App) setupSensorListener(sensor *data.Sensor) {
	listener := app.createAndAddSensorListener(sensor)
	if sensor.Active {
		go listener.Broker.Start()
	} else {
		go listener.Start()
	}
}

// creates and adds a sensor listener to map in app module and returns pointer to it
func (app *App) createAndAddSensorListener(sensor *data.Sensor) (listener *data.Listener[float64]) {
	onNewValue := func(value float64) {
		measuserment := data.SensorMeasurement{
			SensorID:      sensor.ID,
			MeasuredAt:    time.Now(),
			MeasuredValue: value,
		}

		if err := app.models.SensorMeasurements.Insert(&measuserment); err != nil {
			app.logger.Error("Writing measurement to DB", "error", err)
		}
	}

	l := data.NewListener[float64](sensor, onNewValue)
	app.listeners[sensor.ID] = l
	return l
}

func (app *App) stopAndDeleteSensorListener(sensorId uuid.UUID) {
	if l, ok := app.listeners[sensorId]; ok {
		l.GetStopCh() <- struct{}{}
	}

	delete(app.listeners, sensorId)
}

func (app *App) startRule(rule *data.Rule) {
	ch := make(chan struct{}, 2)
	app.rules.stopChannels[rule.ID] = ch
	go rule.Run(app.listeners, app.rules.channel, ch, &app.models.SensorMeasurements)
}

func (app *App) stopRule(ruleId uuid.UUID) {
	if stopCh, ok := app.rules.stopChannels[ruleId]; ok {
		stopCh <- struct{}{}
	}

	delete(app.rules.stopChannels, ruleId)
}

func (app *App) sendValue(url string, body *bytes.Buffer) error {
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		app.logger.Error("handleRuleRequests request creation", "error", err.Error())
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	_, err = app.client.Do(req)
	if err != nil {
		app.logger.Error("handleRuleRequests request", "error", err.Error())
	}

	return nil
}
