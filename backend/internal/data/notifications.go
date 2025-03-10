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
	NotificationLevelSuccess NotificationLevel = "success"
	NotificationLevelWarning NotificationLevel = "warning"
	NotificationLevelInfo    NotificationLevel = "info"
)

type Notification struct {
	ID          uuid.UUID         `json:"id"`
	Level       NotificationLevel `json:"level"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	CreatedAt   time.Time         `json:"created_at"`
}

type UserNotification struct {
	Notification
	Read  bool        `json:"read"`
	Users []uuid.UUID `json:"-"`
}

type NotificationModel struct {
	DB *pgxpool.Pool
}

func (m *NotificationModel) insert(notification *Notification) error {
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
	err := m.insert(notification)
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

func (m *NotificationModel) InsertForAll(notification *Notification) ([]uuid.UUID, error) {
	err := m.insert(notification)
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO user_notifications (notification_id, user_id)
    SELECT $1, id FROM users
    RETURNING user_id
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.Query(ctx, query, notification.ID)
	if err != nil {
		return nil, err
	}

	ids := make([]uuid.UUID, 0)

	for rows.Next() {
		var id uuid.UUID

		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	return ids, nil
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

func (m *NotificationModel) MarkAllAsRead(userId uuid.UUID) error {
	query := `
    UPDATE user_notifications
    SET read = true
    WHERE NOT read AND user_id = $1
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.DB.Exec(ctx, query, userId)

	return err
}

func (m *NotificationModel) GetUnread(userId uuid.UUID) ([]*UserNotification, error) {
	query := `
    SELECT id, level, title, description, created_at, read FROM notifications
    INNER JOIN user_notifications ON notifications.id = user_notifications.notification_id
    WHERE user_id = $1 AND NOT read
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}

	notifications := make([]*UserNotification, 0)

	for rows.Next() {
		notif := UserNotification{Users: []uuid.UUID{userId}}
		err = rows.Scan(&notif.ID, &notif.Level, &notif.Title, &notif.Description, &notif.CreatedAt, &notif.Read)
		if err != nil {
			return nil, err
		}

		notifications = append(notifications, &notif)
	}

	return notifications, nil
}
