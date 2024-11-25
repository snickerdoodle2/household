ALTER TABLE sensors
ADD COLUMN created_at timestamp(0) with time zone NOT NULL DEFAULT now(),
ADD COLUMN version integer NOT NULL DEFAULT 1;