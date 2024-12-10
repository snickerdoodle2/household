ALTER TABLE rules
DROP COLUMN valid_target_type,
DROP COLUMN valid_target_id,
DROP COLUMN valid_target_payload,
ADD COLUMN valid_sensor_id uuid NOT NULL REFERENCES sensors(id),
ADD COLUMN valid_payload json NOT NULL;