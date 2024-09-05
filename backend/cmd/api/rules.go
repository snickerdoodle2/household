package main

import (
	"errors"
	"inzynierka/internal/data"
	"inzynierka/internal/data/validator"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (app *App) createRuleHandler(w http.ResponseWriter, r *http.Request) {
	var rule data.Rule
	err := app.readJSON(w, r, &rule)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	app.logger.Info("rule from request", rule)

	v := validator.New()

	if data.ValidateRule(v, &rule); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Rules.Insert(&rule)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNonExistingTo):
			v.AddError("on_valid.id", "referencing non existing device")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.startRule(&rule)

	err = app.writeJSON(w, http.StatusCreated, envelope{"data": rule}, nil)
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

	rule, err := app.models.Rules.Get(ruleId)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"rule": rule}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) listRulesHandler(w http.ResponseWriter, r *http.Request) {
	rule, err := app.models.Rules.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"data": rule}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) updateRuleHanlder(w http.ResponseWriter, r *http.Request) {
	ruleIdStr := chi.URLParam(r, "id")
	ruleId, err := uuid.Parse(ruleIdStr)

	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": "not a valid uuid"}, nil)
		return
	}

	rule, err := app.models.Rules.Get(ruleId)

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
		Name        *string                 `json:"name"`
		Description *string                 `json:"description"`
		Internal    *map[string]interface{} `json:"internal"`
		OnValid     *data.ValidRuleAction   `json:"on_valid"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if input.Name != nil {
		rule.Name = *input.Name
	}

	if input.Description != nil {
		rule.Description = *input.Description
	}

	if input.Internal != nil {
		internal, err := data.UnmarshalInternalRuleJSON(*input.Internal)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		rule.Internal = internal
	}

	if input.OnValid != nil {
		rule.OnValid = *input.OnValid
	}

	v := validator.New()
	if data.ValidateRule(v, rule); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Rules.Update(rule)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrNonExistingTo):
			v.AddError("on_valid.id", "referencing non existing device")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"data": rule}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) deleteRuleHandler(w http.ResponseWriter, r *http.Request) {
	ruleIdStr := chi.URLParam(r, "id")
	ruleId, err := uuid.Parse(ruleIdStr)

	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": "not a valid uuid"}, nil)
		return
	}

	app.stopRule(ruleId)

	err = app.models.Rules.Delete(ruleId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "rule successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
