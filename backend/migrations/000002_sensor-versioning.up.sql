ALTER TABLE sensors
ADD COLUMN created_at timestamptz(0) NOT NULL DEFAULT now(),
ADD COLUMN version integer NOT NULL DEFAULT 1;