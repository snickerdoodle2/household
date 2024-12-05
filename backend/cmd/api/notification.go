package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (app *App) readNotificationHandler(w http.ResponseWriter, r *http.Request) {
	notifId, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": "not a valid uuid"}, nil)
		return
	}
	user := app.contextGetUser(r)
	app.models.Notifications.MarkAsRead(notifId, user.ID)
}
