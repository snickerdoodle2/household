ALTER TABLE sensors ADD COLUMN active boolean DEFAULT false;

ALTER TABLE sensors 
ADD CONSTRAINT sensors_active_or_refresh_rate_check 
CHECK (active = true OR refresh_rate IS NOT NULL);