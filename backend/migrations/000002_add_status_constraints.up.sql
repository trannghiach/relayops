ALTER TABLE events
ADD CONSTRAINT chk_events_status
CHECK (status IN ('received', 'planned', 'failed'));

ALTER TABLE jobs
ADD CONSTRAINT chk_jobs_status
CHECK (status IN (
    'pending',
    'processing',
    'succeeded',
    'failed',
    'dead_lettered',
    'cancelled'
));

ALTER TABLE jobs
ADD CONSTRAINT chk_jobs_channel
CHECK (channel IN ('email', 'webhook'));

ALTER TABLE delivery_attempts
ADD CONSTRAINT chk_delivery_attempts_status
CHECK (status IN ('succeeded', 'failed', 'timeout'));