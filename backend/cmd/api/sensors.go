package main

import (
	"errors"
	"inzynierka/internal/data"
	"inzynierka/internal/data/validator"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
	// TODO: uuid from url + db drop

	err := app.writeJSON(w, http.StatusNotImplemented, envelope{"status": "todo"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
