CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY,
    username citext UNIQUE NOT NULL,
    display_name text NOT NULL,
    password_hash bytea NOT NULL,
    created_at timestamptz(0) NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);