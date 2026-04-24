DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM pg_catalog.pg_roles
        WHERE rolname = 'relayops_readonly'
    ) THEN
        REVOKE SELECT ON challenge_flags FROM relayops_readonly;
        REVOKE SELECT ON dead_letters FROM relayops_readonly;
        REVOKE SELECT ON delivery_attempts FROM relayops_readonly;
        REVOKE SELECT ON jobs FROM relayops_readonly;
        REVOKE SELECT ON events FROM relayops_readonly;

        REVOKE USAGE ON SCHEMA public FROM relayops_readonly;

        EXECUTE format(
            'REVOKE CONNECT ON DATABASE %I FROM relayops_readonly',
            current_database()
        );
    END IF;
END
$$;

DROP TABLE IF EXISTS challenge_flags;