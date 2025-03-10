package main

import (
	"errors"
	"inzynierka/internal/data"
	"inzynierka/internal/data/validator"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *App) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string        `json:"name"`
		Username string        `json:"username"`
		Role     data.UserRole `json:"role"`
		Password string        `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:     input.Name,
		Role:     input.Role,
		Username: input.Username,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateUsername):
			v.AddError("username", "user with this username already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// NOTE: for now you can update only Name
func (app *App) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	user, err := app.models.Users.GetByUsername(username)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	curUser := app.contextGetUser(r)
	if !(curUser.Role == data.UserRoleAdmin || curUser.ID == user.ID) {
		app.authenticationRequiredResponse(w, r)
		return
	}
	var input struct {
		Name     *string        `json:"name"`
		Role     *data.UserRole `json:"role"`
		Password *string        `json:"password"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if curUser.Role != data.UserRoleAdmin {
		input.Role = nil
	}

	if input.Name != nil {
		user.Name = *input.Name
	}

	if input.Role != nil {
		user.Role = *input.Role
	}

	if input.Password != nil {
		err = user.Password.Set(*input.Password)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err = app.models.Users.Update(user); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *App) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	err := app.models.Users.DeleteByUsername(username)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "user successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) getCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	err := app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []*data.User
	curUser := app.contextGetUser(r)
	if curUser.Role == data.UserRoleAdmin {
		usersTMP, err := app.models.Users.GetAllUsers()
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		users = usersTMP
	} else {
		users = []*data.User{curUser}
	}

	err := app.writeJSON(w, http.StatusOK, envelope{"data": users}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *App) getUserHandler(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	user, err := app.models.Users.GetByUsername(username)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
	}

	curUser := app.contextGetUser(r)
	if !(curUser.Role == data.UserRoleAdmin || curUser.ID == user.ID) {
		app.logger.Debug("updateUserHandler", "role", curUser.Role)
		app.authenticationRequiredResponse(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
