ALTER TABLE rules
DROP COLUMN valid_sensor_id,
DROP COLUMN valid_payload,
ADD COLUMN valid_target_type varchar(255) NOT NULL,
ADD COLUMN valid_target_id uuid NOT NULL,
ADD COLUMN valid_target_payload json;
