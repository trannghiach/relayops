DROP INDEX IF EXISTS idx_audit_logs_created_at;
DROP INDEX IF EXISTS idx_audit_logs_entity;

DROP INDEX IF EXISTS idx_dead_letters_created_at;

DROP INDEX IF EXISTS idx_delivery_attempts_started_at;
DROP INDEX IF EXISTS idx_delivery_attempts_job_id;

DROP INDEX IF EXISTS uq_jobs_dedupe_key;
DROP INDEX IF EXISTS idx_jobs_created_at;
DROP INDEX IF EXISTS idx_jobs_channel_status;
DROP INDEX IF EXISTS idx_jobs_status_available_at;
DROP INDEX IF EXISTS idx_jobs_event_id;

DROP INDEX IF EXISTS uq_events_idempotency_key;
DROP INDEX IF EXISTS idx_events_created_at;
DROP INDEX IF EXISTS idx_events_source;
DROP INDEX IF EXISTS idx_events_event_type;

DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS dead_letters;
DROP TABLE IF EXISTS delivery_attempts;
DROP TABLE IF EXISTS jobs;
DROP TABLE IF EXISTS events;