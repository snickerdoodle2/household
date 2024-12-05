package main

import (
	"inzynierka/internal/data"
	"net/http"
	"time"

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

func (app *App) requestAllNotifsHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	notifs := []data.UserNotification{
		{
			Notification: data.Notification{
				ID:          uuid.New(),
				Level:       data.NotificationLevelError,
				Title:       "Error",
				Description: "Example error notification",
				CreatedAt:   time.Now(),
			},
			Read:  false,
			Users: []uuid.UUID{user.ID},
		},
		{
			Notification: data.Notification{
				ID:          uuid.New(),
				Level:       data.NotificationLevelSuccess,
				Title:       "Success",
				Description: "Example success notification",
				CreatedAt:   time.Now(),
			},
			Read:  false,
			Users: []uuid.UUID{user.ID},
		},
		{
			Notification: data.Notification{
				ID:          uuid.New(),
				Level:       data.NotificationLevelInfo,
				Title:       "Info",
				Description: "Example info notification",
				CreatedAt:   time.Now(),
			},
			Read:  false,
			Users: []uuid.UUID{user.ID},
		},
		{
			Notification: data.Notification{
				ID:          uuid.New(),
				Level:       data.NotificationLevelWarning,
				Title:       "Warning",
				Description: "Example warning notification",
				CreatedAt:   time.Now(),
			},
			Read:  false,
			Users: []uuid.UUID{user.ID},
		},
	}

	for _, notif := range notifs {
		app.notificationBroker.Publish(notif)
	}
}
