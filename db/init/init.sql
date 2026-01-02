CREATE DATABASE "smart-home";

DO
$$
    BEGIN
        IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'user') THEN
            CREATE USER "user" WITH PASSWORD 'password';
        END IF;
    END
$$;

\c "smart-home"

CREATE
    EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users
(
    id         UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    name       VARCHAR(256),
    email      VARCHAR(256) NOT NULL UNIQUE,
    password   TEXT         NOT NULL,
    confirmed  BOOL                  DEFAULT FALSE,
    created_at timestamptz  NOT NULL DEFAULT now(),
    updated_at timestamptz  NOT NULL DEFAULT now(),
    deleted_at timestamptz  NOT NULL DEFAULT NULL
);

CREATE TABLE sessions
(
    id          UUID PRIMARY KEY     DEFAULT uuid_generate_v4(),
    user_id     UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    last_active timestamptz NOT NULL DEFAULT NOW(),
    created_at  timestamptz NOT NULL DEFAULT NOW(),
    expires_at  timestamptz          DEFAULT NULL
);

CREATE TABLE manufactured_devices
(
    id          UUID PRIMARY KEY     DEFAULT uuid_generate_v4(),
    secret      TEXT        NOT NULL,
    mac_address TEXT        NOT NULL UNIQUE,
    registered  BOOLEAN              DEFAULT FALSE,
    created_at  timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE devices
(
    id          UUID PRIMARY KEY REFERENCES manufactured_devices (id),
    user_id     UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name        TEXT,
    description TEXT,
    room        TEXT,
    type        TEXT,
    status_info JSONB,
    custom_data JSONB,
    device_info JSONB,
    last_seen   timestamptz DEFAULT NOW(),
    created_at  timestamptz DEFAULT NOW(),
    updated_at  timestamptz DEFAULT NOW(),
    deleted_at  timestamptz DEFAULT NULL
);

CREATE TABLE capabilities
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id   UUID NOT NULL REFERENCES devices (id) ON DELETE CASCADE,
    type        TEXT NOT NULL CHECK (type IN
                                     ('devices.capabilities.on_off', 'devices.capabilities.color_setting',
                                      'devices.capabilities.mode', 'devices.capabilities.range',
                                      'devices.capabilities.toggle')),
    retrievable BOOL NOT NULL    DEFAULT True,
    reportable  BOOL NOT NULL    DEFAULT False,
    parameters  JSONB,
    state       JSONB
);

CREATE TABLE properties
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_id   UUID NOT NULL REFERENCES devices (id) ON DELETE CASCADE,
    type        TEXT NOT NULL CHECK (type IN ('devices.properties.float', 'devices.properties.event')),
    retrievable BOOL NOT NULL    DEFAULT True,
    reportable  BOOL NOT NULL    DEFAULT False,
    parameters  JSONB,
    state       JSONB
);

CREATE TABLE oauth_clients
(
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    client_secret TEXT,
    redirect_uri  TEXT,
    scope         TEXT,
    name          TEXT,
    enabled       BOOL             DEFAULT True,
    created_at    timestamptz      DEFAULT NOW(),
    updated_at    timestamptz      DEFAULT NOW(),
    deleted_at    timestamptz      DEFAULT null
);

GRANT CONNECT ON DATABASE "smart-home" TO "user";

GRANT USAGE ON SCHEMA public TO "user";

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "user";
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO "user";

ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO "user";
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO "user";
