DELETE FROM rules
WHERE valid_target_type = 'sequence';

ALTER TABLE rules
    DROP COLUMN valid_target_type;

ALTER TABLE rules
    RENAME COLUMN valid_target_id TO valid_sensor_id;

ALTER TABLE rules
    RENAME COLUMN valid_target_payload TO valid_payload;

ALTER TABLE rules
    ADD CONSTRAINT rules_valid_sensor_id_fkey
        FOREIGN KEY (valid_sensor_id) REFERENCES sensors(id),
    ALTER COLUMN valid_sensor_id SET NOT NULL,
    ALTER COLUMN valid_payload SET NOT NULL;