package data

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationLevel string

const (
	NotificationLevelError   NotificationLevel = "error"
	NotificationLevelUrgent  NotificationLevel = "urgent"
	NotificationLevelWarning NotificationLevel = "warning"
	NotificationLevelInfo    NotificationLevel = "info"
)

type Notification struct {
	ID          uuid.UUID
	Level       NotificationLevel
	Title       string
	Description string
	CreatedAt   time.Time
}

type NotificationModel struct {
	DB *pgxpool.Pool
}

func (m *NotificationModel) Insert(notification *Notification) error {
	query := `
    INSERT INTO notifications (id, level, title, description)
    VALUES ($1, $2, $3, $4)
    RETURNING created_at
    `

	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	notification.ID = id

	args := []any{notification.ID, notification.Level, notification.Title, notification.Description}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return m.DB.QueryRow(ctx, query, args...).Scan(notification.CreatedAt)
}
