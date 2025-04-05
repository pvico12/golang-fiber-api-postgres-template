DO $$
DECLARE
    db_name TEXT := '<db_name>';
    role_name TEXT := '<role_name>';
    role_password TEXT := '<your-password>';
BEGIN
    -- Create the role
    EXECUTE format('CREATE ROLE %I WITH LOGIN PASSWORD %L', role_name, role_password);

    -- Grant database connection
    EXECUTE format('GRANT CONNECT ON DATABASE %I TO %I', db_name, role_name);

    -- Grant schema permissions
    EXECUTE format('GRANT ALL ON SCHEMA public TO %I', role_name);
    EXECUTE format('GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO %I', role_name);
    EXECUTE format('GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO %I', role_name);

    -- Set default privileges for future objects
    EXECUTE format('ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO %I', role_name);
    EXECUTE format('ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO %I', role_name);
END $$;