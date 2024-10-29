CREATE TABLE IF NOT EXISTS rules (
    id uuid PRIMARY KEY,
    name varchar(255) NOT NULL,
    description text,
    internal json NOT NULL,
    valid_sensor_id uuid NOT NULL REFERENCES sensors(id),
    valid_payload json NOT NULL,
    created_at timestamptz(0) NOT NULL DEFAULT now(),
    version integer NOT NULL DEFAULT 1
);