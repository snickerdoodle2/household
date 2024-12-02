package main

import (
	"errors"
	"inzynierka/internal/data"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

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
