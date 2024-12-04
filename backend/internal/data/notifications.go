package data

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

type UserNotification struct {
	Notification
	Read bool
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

	return m.DB.QueryRow(ctx, query, args...).Scan(&notification.CreatedAt)
}

func (m *NotificationModel) InsertForUsers(notification *Notification, users []uuid.UUID) error {
	err := m.Insert(notification)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = m.DB.CopyFrom(ctx,
		pgx.Identifier{"user_notifications"},
		[]string{"notification_id", "user_id"},
		pgx.CopyFromSlice(len(users), func(i int) ([]any, error) {
			return []any{notification.ID, users[i]}, nil
		}),
	)

	return err
}

func (m *NotificationModel) InsertUser(notificationId, userId uuid.UUID) error {
	query := `
    INSERT INTO user_notifications (notification_id, user_id)
    VALUES ($1, $2)
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.DB.Exec(ctx, query, notificationId, userId)

	// TODO: Handle user already notified
	return err
}

func (m *NotificationModel) MarkAsRead(notificationId, userId uuid.UUID) error {
	query := `
    UPDATE user_notifications
    SET read = true
    WHERE notification_id = $1 AND user_id = $2
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.DB.Exec(ctx, query, notificationId, userId)

	// TODO: Handle user already read this
	return err
}

func (m *NotificationModel) GetUnread(userId uuid.UUID) error {
	return nil
}
