package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"inzynierka/internal/data"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type actionData struct {
	Body  bytes.Buffer
	Url   string
	Delay time.Duration
}

func (app *App) createSequenceHandler(w http.ResponseWriter, r *http.Request) {
	var sequence data.Sequence

	err := app.readJSON(w, r, &sequence)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.Sequences.Insert(&sequence)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"data": sequence}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *App) listSequencesHandler(w http.ResponseWriter, r *http.Request) {
	sequencesInfo, err := app.models.Sequences.GetAllInfo()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"data": sequencesInfo}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) getSequenceHandler(w http.ResponseWriter, r *http.Request) {
	sequenceIdStr := chi.URLParam(r, "id")
	sequenceId, err := uuid.Parse(sequenceIdStr)

	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": "not a valid uuid"}, nil)
		return
	}

	sequence, err := app.models.Sequences.Get(sequenceId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"sequence": sequence}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) updateSequenceHandler(w http.ResponseWriter, r *http.Request) {
	sequenceIdStr := chi.URLParam(r, "id")
	sequenceId, err := uuid.Parse(sequenceIdStr)

	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": "not a valid uuid"}, nil)
		return
	}

	sequence, err := app.models.Sequences.Get(sequenceId)
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
		Name        *string                `json:"name"`
		Description *string                `json:"description"`
		Actions     *[]data.SequenceAction `json:"actions"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if input.Name != nil {
		sequence.Name = *input.Name
	}
	if input.Description != nil {
		sequence.Description = *input.Description
	}
	if input.Actions != nil {
		sequence.Actions = *input.Actions
	}

	err = app.models.Sequences.Update(sequence)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"data": sequence}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) deleteSequenceHandler(w http.ResponseWriter, r *http.Request) {
	sequenceIdStr := chi.URLParam(r, "id")
	sequenceId, err := uuid.Parse(sequenceIdStr)

	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": "not a valid uuid"}, nil)
		return
	}

	err = app.models.Sequences.Delete(sequenceId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "sequence succesfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) startSequenceHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Debug("Start sequence request")
	sequenceIdStr := chi.URLParam(r, "id")
	sequenceId, err := uuid.Parse(sequenceIdStr)

	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": "not a valid uuid"}, nil)
		return
	}

	sequence, err := app.models.Sequences.Get(sequenceId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	preparedData, err := app.prepareActionData(sequence.Actions)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"success": "sequence execution has started"}, nil)

	go app.executeSequence(preparedData)
}

func (app *App) executeSequence(preparedData []actionData) error {
	for _, actionData := range preparedData {
		time.Sleep(actionData.Delay)

		app.sendValue(actionData.Url, &actionData.Body)
	}

	return nil
}

func (app *App) prepareActionData(actions []data.SequenceAction) ([]actionData, error) {
	var preparedData []actionData

	for _, action := range actions {
		var data actionData

		uri, err := app.models.Sensors.GetUri(action.Target)
		if err != nil {
			return nil, err
		}

		url := fmt.Sprintf("http://%s/value", uri)

		data.Url = url
		data.Delay = time.Duration(action.MsDelay) * time.Millisecond

		var payload struct {
			Value float32 `json:"value"`
		}

		payload.Value = action.Value
		body := new(bytes.Buffer)
		err = json.NewEncoder(body).Encode(payload)
		if err != nil {
			return nil, err
		}
		data.Body = *body

		preparedData = append(preparedData, data)
	}

	return preparedData, nil
}
