CREATE TYPE notification_level AS ENUM (
    'error',
    'urgent',
    'warning',
    'info'
    );

CREATE TABLE IF NOT EXISTS notifications (
    id uuid PRIMARY KEY,
    level notification_level NOT NULL,
    title text NOT NULL,
    description text NOT NULL,
    created_at timestamptz(0) DEFAULT now(),
    CHECK ( char_length(title) > 0 AND char_length(description) > 0 )
);

CREATE TABLE IF NOT EXISTS user_notifications (
    notification_id uuid references notifications(id) ON DELETE CASCADE,
    user_id uuid NOT NULL references users(id) ON DELETE CASCADE,
    read bool DEFAULT FALSE,
    PRIMARY KEY (notification_id, user_id)
)