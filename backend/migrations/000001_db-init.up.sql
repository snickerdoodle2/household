CREATE TYPE sensor_type AS ENUM (
    'binary_switch',
    'binary_sensor',
    'decimal_switch',
    'decimal_sensor',
    'button'
    );

CREATE TABLE IF NOT EXISTS sensors (
    id uuid PRIMARY KEY,
    name varchar(255) NOT NULL,
    uri text NOT NULL,
    sensor_type sensor_type NOT NULL,
    refresh_rate int NOT NULL DEFAULT 5
)