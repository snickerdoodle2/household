package main

import (
	"inzynierka/internal/data"
	"net/http"
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
