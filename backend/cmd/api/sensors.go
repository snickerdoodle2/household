package main

import (
	"inzynierka/internal/data"
	"inzynierka/internal/data/validator"
	"net/http"

	"github.com/google/uuid"
)

func (app *App) listSensorsHandler(w http.ResponseWriter, r *http.Request) {
	sensors := [...]data.Sensor{
		{
			ID:          uuid.New(),
			Name:        "Temperature",
			URI:         "localhost:9999",
			Type:        data.DecimalSensor,
			RefreshRate: 5,
		},
		{
			ID:          uuid.New(),
			Name:        "Entry doors",
			URI:         "localhost:11111",
			Type:        data.BinarySensor,
			RefreshRate: 10,
		},
	}

	err := app.writeJSON(w, http.StatusOK, envelope{"data": sensors}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) getSensorHandler(w http.ResponseWriter, r *http.Request) {
	// idk co tutaj

	err := app.writeJSON(w, http.StatusNotImplemented, envelope{"status": "todo"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) createSensorHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		URI         string `json:"uri"`
		Type        string `json:"type"`
		RefreshRate int    `json:"refresh_rate"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	sensor := &data.Sensor{
		Name:        input.Name,
		URI:         input.URI,
		Type:        data.SensorType(input.Type),
		RefreshRate: input.RefreshRate,
	}

	v := validator.New()

	if data.ValidateSensor(v, sensor); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.writeJSON(w, http.StatusNotImplemented, envelope{"data": sensor}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) updateSensorHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: uuid from url + db update

	err := app.writeJSON(w, http.StatusNotImplemented, envelope{"status": "todo"}, nil)
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
