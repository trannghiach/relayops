ALTER TABLE delivery_attempts
DROP CONSTRAINT IF EXISTS chk_delivery_attempts_status;

ALTER TABLE jobs
DROP CONSTRAINT IF EXISTS chk_jobs_channel;

ALTER TABLE jobs
DROP CONSTRAINT IF EXISTS chk_jobs_status;

ALTER TABLE events
DROP CONSTRAINT IF EXISTS chk_events_status;