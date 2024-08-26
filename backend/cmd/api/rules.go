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

func (app *App) listRulesHandler(w http.ResponseWriter, r *http.Request) {
	rules, err := app.models.Rules.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"data": rules}, nil)
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

	ruleS, err := app.models.Rules.Get(ruleId)
	_ = ruleS

	if err != nil {
		switch {
		case errors.Is(err, rule.ErrRecordNotFound):
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
		OnValid     *rule.ValidRuleAction   `json:"on_valid"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if input.Name != nil {
		ruleS.Name = *input.Name
	}

	if input.Description != nil {
		ruleS.Description = *input.Description
	}

	if input.Internal != nil {
		internal, err := rule.UnmarshalInternalRuleJSON(*input.Internal)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		ruleS.Internal = internal
	}

	if input.OnValid != nil {
		ruleS.OnValid = *input.OnValid
	}

	v := validator.New()
	if rule.ValidateRule(v, ruleS); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Rules.Update(ruleS)

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

	err = app.writeJSON(w, http.StatusCreated, envelope{"data": ruleS}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
