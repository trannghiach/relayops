CREATE TABLE events (
    id UUID PRIMARY KEY,
    event_type VARCHAR(100) NOT NULL,
    source VARCHAR(100) NOT NULL,
    tenant_id VARCHAR(100),
    aggregate_type VARCHAR(100),
    aggregate_id VARCHAR(100),
    payload JSONB NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    idempotency_key VARCHAR(255),
    status VARCHAR(32) NOT NULL DEFAULT 'received',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE jobs (
    id UUID PRIMARY KEY,
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    job_type VARCHAR(50) NOT NULL,
    channel VARCHAR(50) NOT NULL,
    status VARCHAR(32) NOT NULL,
    recipient VARCHAR(255),
    target_url TEXT,
    payload JSONB NOT NULL,
    dedupe_key VARCHAR(255),
    scheduled_at TIMESTAMPTZ,
    available_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    attempts INT NOT NULL DEFAULT 0,
    max_attempts INT NOT NULL DEFAULT 5,
    last_error TEXT,
    claimed_at TIMESTAMPTZ,
    worker_id VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE delivery_attempts (
    id UUID PRIMARY KEY,
    job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    attempt_no INT NOT NULL,
    provider VARCHAR(50) NOT NULL,
    status VARCHAR(32) NOT NULL,
    response_code INT,
    response_body TEXT,
    error_message TEXT,
    started_at TIMESTAMPTZ NOT NULL,
    finished_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE dead_letters (
    id UUID PRIMARY KEY,
    job_id UUID NOT NULL UNIQUE REFERENCES jobs(id) ON DELETE CASCADE,
    reason TEXT NOT NULL,
    payload_snapshot JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY,
    actor_type VARCHAR(32) NOT NULL,
    actor_id VARCHAR(100),
    entity_type VARCHAR(32) NOT NULL,
    entity_id UUID NOT NULL,
    action VARCHAR(50) NOT NULL,
    details JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_events_event_type ON events(event_type);
CREATE INDEX idx_events_source ON events(source);
CREATE INDEX idx_events_created_at ON events(created_at DESC);
CREATE UNIQUE INDEX uq_events_idempotency_key ON events(idempotency_key)
WHERE idempotency_key IS NOT NULL;

CREATE INDEX idx_jobs_event_id ON jobs(event_id);
CREATE INDEX idx_jobs_status_available_at ON jobs(status, available_at);
CREATE INDEX idx_jobs_channel_status ON jobs(channel, status);
CREATE INDEX idx_jobs_created_at ON jobs(created_at DESC);
CREATE UNIQUE INDEX uq_jobs_dedupe_key ON jobs(dedupe_key)
WHERE dedupe_key IS NOT NULL;

CREATE INDEX idx_delivery_attempts_job_id ON delivery_attempts(job_id);
CREATE INDEX idx_delivery_attempts_started_at ON delivery_attempts(started_at DESC);

CREATE INDEX idx_dead_letters_created_at ON dead_letters(created_at DESC);

CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);