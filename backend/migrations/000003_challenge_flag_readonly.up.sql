CREATE TABLE IF NOT EXISTS challenge_flags (
    name TEXT PRIMARY KEY,
    value TEXT NOT NULL
);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_catalog.pg_roles
        WHERE rolname = 'relayops_readonly'
    ) THEN
        RAISE EXCEPTION 'missing database role: relayops_readonly. Create it before running migrations.';
    END IF;
END
$$;

DO $$
BEGIN
    EXECUTE format(
        'GRANT CONNECT ON DATABASE %I TO relayops_readonly',
        current_database()
    );
END
$$;

GRANT USAGE ON SCHEMA public TO relayops_readonly;

GRANT SELECT ON events TO relayops_readonly;
GRANT SELECT ON jobs TO relayops_readonly;
GRANT SELECT ON delivery_attempts TO relayops_readonly;
GRANT SELECT ON dead_letters TO relayops_readonly;
GRANT SELECT ON challenge_flags TO relayops_readonly;

REVOKE CREATE ON SCHEMA public FROM relayops_readonly;