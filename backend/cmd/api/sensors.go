package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"inzynierka/internal/data"
	"inzynierka/internal/data/validator"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (app *App) listSensorsHandler(w http.ResponseWriter, r *http.Request) {
	sensors, err := app.models.Sensors.GetAllInfo()
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
		Active      bool            `json:"active"`
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
		Active:      input.Active,
	}

	if sensor.Active {
		app.initSensor(*sensor)
	} else {
		app.insertAndHandleListener(sensor, w, r)
	}
}

func (app *App) insertAndHandleListener(sensor *data.Sensor, w http.ResponseWriter, r *http.Request) {
	v := validator.New()

	if data.ValidateSensor(v, sensor); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err := app.models.Sensors.Insert(sensor)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateUri):
			v.AddError("uri", "a sensor with this URI already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	listener := app.createAndAddSensorListener(sensor)
	if !sensor.Active {
		listener.Start()
	}
}

// function to initialize active (for now) sensor.
// Sends init request to the sensor and adds it to initBuffer.
// Sensor should send init-ack request to init-ack endpoint to be removed from initBuffer and further processed
func (app *App) initSensor(sensor data.Sensor) error {
	sensorEndpoint := fmt.Sprintf("http://%v/init", sensor.URI)
	initAckEndpoint := "/api/v1/sensor/measurements" // TODO: find a way to get this from the app
	measurementsEndpoint := "/api/v1/sensor/init-ack"
	idToken, err := uuid.NewRandom()
	app.initBuffer[idToken] = sensor

	var requestBody struct {
		IdToken              uuid.UUID `json:"id-token"`
		ServerUri            string    `json:"server-uri"`
		InitAckEndpoint      string    `json:"init-ack-endpoint"`
		MeasurementsEndpoint string    `json:"measurements-endpoint"`
	}

	requestBody.IdToken = idToken
	requestBody.ServerUri = fmt.Sprintf("%s:%d", app.config.host, app.config.port)
	requestBody.InitAckEndpoint = initAckEndpoint
	requestBody.MeasurementsEndpoint = measurementsEndpoint

	if err != nil {
		return err
	}

	client := &http.Client{}
	client.Timeout = 5 * time.Second

	body := new(bytes.Buffer)
	err = json.NewEncoder(body).Encode(requestBody)
	if err != nil {
		app.logger.Error("sendInitRequest marshall", "error", err.Error())
		return err
	}

	req, err := http.NewRequest(http.MethodPut, sensorEndpoint, body)
	if err != nil {
		app.logger.Error("handleRuleRequests request creation", "error", err.Error())
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = client.Do(req)
	if err != nil {
		app.logger.Error("handleInitRequest request", "error", err.Error())
		return err
	}

	return nil
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
		Active      *bool            `json:"active"`
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
		sensor.URI = *input.URI
	}

	if input.Type != nil {
		sensor.Type = *input.Type
	}

	if input.RefreshRate != nil {
		sensor.RefreshRate = *input.RefreshRate
	}

	if input.Active != nil {
		sensor.Active = *input.Active
	}

	v := validator.New()
	if data.ValidateSensor(v, sensor); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Sensors.Update(sensor)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateUri):
			v.AddError("uri", "a sensor with this URI already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.stopAndDeleteSensorListener(sensorId)
	listener := app.createAndAddSensorListener(sensor)
	if !sensor.Active {
		listener.Start()
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

	app.stopAndDeleteSensorListener(sensorId)

	err = app.models.Sensors.DeleteSensorAndMeasurements(sensorId)
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

func (app *App) activeSensorHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: different sensor identification - address:port impossible
	uri := r.RemoteAddr
	app.logger.Info("received active sensor measurement", "uri", uri)

	var requestBody struct {
		MessageType string  `json:"message-type"`
		SensorType  string  `json:"sensor-type"`
		Value       float64 `json:"value"`
	}

	err := app.readJSON(w, r, &requestBody)
	if err != nil {
		app.logger.Warn("active sensor error", "request body", err.Error())
		return
	}
	r.Body.Close()

	id, err := app.models.Sensors.GetIdByUriAndType(uri, requestBody.SensorType)

	if err != nil {
		app.logger.Warn("error getting active sensor ID", "error", err.Error())
		return
	}

	app.listeners[id].Broker.Publish([]float64{requestBody.Value})

	measurement := data.SensorMeasurement{
		SensorID:      id,
		MeasuredAt:    time.Now(),
		MeasuredValue: requestBody.Value,
	}

	err = app.models.SensorMeasurements.Insert(&measurement)
	if err != nil {
		app.logger.Error("Writing measurement to DB", "error", err)
	}

	w.WriteHeader(http.StatusAccepted)
}

func (app *App) initAckHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Debug("init-ack received")
	var requestBodyData struct {
		IdToken              uuid.UUID `json:"id-token"`
		ServerUri            string    `json:"server-uri"`
		InitAckEndpoint      string    `json:"init-ack-endpoint"`
		MeasurementsEndpoint string    `json:"measurements-endpoint"`
	}

	err := app.readJSON(w, r, &requestBodyData)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	sensor, ok := app.initBuffer[requestBodyData.IdToken]
	if !ok {
		app.logger.Warn("init-ack received for unknown sensor", "id-token", requestBodyData.IdToken)
		return
	}

	sensor.IdToken = requestBodyData.IdToken

	app.insertAndHandleListener(&sensor, w, r)

	delete(app.initBuffer, requestBodyData.IdToken)

}
