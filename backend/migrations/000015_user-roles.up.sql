CREATE TYPE user_role AS ENUM (
    'admin',
    'user'
    );

ALTER TABLE users ADD COLUMN IF NOT EXISTS role user_role NOT NULL DEFAULT 'user';
UPDATE users SET role = 'admin' WHERE username = 'admin';