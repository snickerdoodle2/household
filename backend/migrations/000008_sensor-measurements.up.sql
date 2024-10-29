CREATE TABLE IF NOT EXISTS sensor_measurements (
    sensor_id uuid NOT NULL REFERENCES sensors(id),
    measured_at timestamptz(1) NOT NULL,
    measured_value real NOT NULL,
    PRIMARY KEY (sensor_id, measured_at)
)