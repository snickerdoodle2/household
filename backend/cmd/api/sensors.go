package main

import (
	"errors"
	"inzynierka/internal/data"
	"inzynierka/internal/data/validator"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func (app *App) listSensorsHandler(w http.ResponseWriter, r *http.Request) {
	sensors, err := app.models.Sensors.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"data": sensors}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) getSensorHandler(w http.ResponseWriter, r *http.Request) {
	sensorIdStr := chi.URLParam(r, "id")
	sensorId, err := uuid.Parse(sensorIdStr)

	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": "not a valid uuid"}, nil)
		return
	}

	sensor, err := app.models.Sensors.Get(sensorId)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"sensor": sensor}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) getSensorValueHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": "not a valid uuid"}, nil)
		return
	}

	listener, ok := app.listeners[id]
	if !ok {
		app.errorResponse(w, r, http.StatusInternalServerError, nil)
	}

	// TODO: Change broker depending on sensor id
	conn, err := app.upgrader.Upgrade(w, r, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	defer conn.Close()

	endCh := make(chan struct{})
	go func() {
		// TODO: Handle closing ws from server (this function never ends)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				endCh <- struct{}{}
				return
			}
		}
	}()

	msgCh := listener.GetBroker().Subscribe()

	for msg := range msgCh {
		conn.WriteMessage(websocket.TextMessage, msg)
	}

	// TODO: Send messages to client
}

func (app *App) createSensorHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string          `json:"name"`
		URI         string          `json:"uri"`
		Type        data.SensorType `json:"type"`
		RefreshRate int             `json:"refresh_rate"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	sensor := &data.Sensor{
		Name:        input.Name,
		URI:         input.URI,
		Type:        input.Type,
		RefreshRate: input.RefreshRate,
	}

	v := validator.New()

	if data.ValidateSensor(v, sensor); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Sensors.Insert(sensor)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"data": sensor}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) updateSensorHandler(w http.ResponseWriter, r *http.Request) {
	sensorIdStr := chi.URLParam(r, "id")
	sensorId, err := uuid.Parse(sensorIdStr)

	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": "not a valid uuid"}, nil)
		return
	}

	sensor, err := app.models.Sensors.Get(sensorId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Name        *string          `json:"name"`
		URI         *string          `json:"uri"`
		Type        *data.SensorType `json:"type"`
		RefreshRate *int             `json:"refresh_rate"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if input.Name != nil {
		sensor.Name = *input.Name
	}

	if input.URI != nil {
		sensor.URI = *&sensor.URI
	}

	if input.Type != nil {
		sensor.Type = *input.Type
	}

	if input.RefreshRate != nil {
		sensor.RefreshRate = *input.RefreshRate
	}

	v := validator.New()
	if data.ValidateSensor(v, sensor); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Sensors.Update(sensor)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"sensor": sensor}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) deleteSensorHandler(w http.ResponseWriter, r *http.Request) {
	sensorIdStr := chi.URLParam(r, "id")
	sensorId, err := uuid.Parse(sensorIdStr)

	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": "not a valid uuid"}, nil)
		return
	}

	app.stopSensorListener(sensorId)

	err = app.models.Sensors.Delete(sensorId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "sensor successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
