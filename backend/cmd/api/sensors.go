package main

import (
	"inzynierka/internal/data"
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
	// TODO: add input validating + db insert

	err := app.writeJSON(w, http.StatusNotImplemented, envelope{"status": "todo"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
