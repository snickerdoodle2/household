package main

import (
	"errors"
	"inzynierka/internal/data/rule"
	"inzynierka/internal/data/validator"
	"net/http"
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
