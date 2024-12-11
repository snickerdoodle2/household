ALTER TABLE rules
DROP CONSTRAINT rules_valid_sensor_id_fkey,
ALTER COLUMN valid_payload DROP NOT NULL;

ALTER TABLE rules
RENAME COLUMN valid_sensor_id TO valid_target_id;

ALTER TABLE rules
RENAME COLUMN valid_payload TO valid_target_payload;

ALTER TABLE rules
ADD COLUMN valid_target_type varchar(255);

UPDATE rules
SET valid_target_type = 'sensor';

ALTER TABLE rules
ALTER COLUMN valid_target_type SET NOT NULL,
ALTER COLUMN valid_target_id SET NOT NULL;