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

	app.writeJSON(w, http.StatusOK, envelope{"data": sensors}, nil)

}
