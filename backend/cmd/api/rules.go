package main

import (
	"errors"
	"inzynierka/internal/data/rule"
	"inzynierka/internal/data/validator"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (app *App) createRuleHandler(w http.ResponseWriter, r *http.Request) {
	var ruleIn rule.Rule
	err := app.readJSON(w, r, &ruleIn)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	app.logger.Info("rule from request", ruleIn)

	v := validator.New()

	if rule.ValidateRule(v, &ruleIn); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Rules.Insert(&ruleIn)
	if err != nil {
		switch {
		case errors.Is(err, rule.ErrNonExistingTo):
			v.AddError("on_valid.id", "referencing non existing device")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"data": ruleIn}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) getRuleHandler(w http.ResponseWriter, r *http.Request) {
	ruleIdStr := chi.URLParam(r, "id")
	ruleId, err := uuid.Parse(ruleIdStr)

	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": "not a valid uuid"}, nil)
		return
	}

	ruleS, err := app.models.Rules.Get(ruleId)

	if err != nil {
		switch {
		case errors.Is(err, rule.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"rule": ruleS}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
