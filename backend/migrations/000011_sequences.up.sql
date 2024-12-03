CREATE TABLE IF NOT EXISTS sequences (
    id uuid PRIMARY KEY,
    name varchar(255) NOT NULL,
    description text,
    actions json NOT NULL,
    created_at timestamptz(0) NOT NULL DEFAULT now(),
    version integer NOT NULL DEFAULT 1
);