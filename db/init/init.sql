CREATE DATABASE "smart-home";

DO
$$
    BEGIN
        IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'user') THEN
            CREATE USER "user" WITH PASSWORD 'password';
        END IF;
    END
$$;


ALTER ROLE "user" SET timezone TO 'Europe/Moscow';
ALTER ROLE root SET timezone TO 'Europe/Moscow';

\c "smart-home"

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users
(
    id         UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    name       VARCHAR(256),
    email      VARCHAR(256) NOT NULL UNIQUE,
    password   TEXT         NOT NULL,
    confirmed  BOOL                  DEFAULT FALSE,
    created_at timestamptz  NOT NULL DEFAULT now(),
    updated_at timestamptz  NOT NULL DEFAULT now(),
    deleted_at timestamptz           DEFAULT NULL
);

CREATE TABLE oauth_sessions
(
    id          UUID PRIMARY KEY     DEFAULT uuid_generate_v4(),
    user_id     UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    last_active timestamptz NOT NULL DEFAULT NOW(),
    created_at  timestamptz NOT NULL DEFAULT NOW(),
    expires_at  timestamptz          DEFAULT NULL
);

CREATE TABLE panel_sessions
(
    id            UUID PRIMARY KEY     DEFAULT uuid_generate_v4(),
    user_id       UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token_type    VARCHAR(25),
    access_token  TEXT        NOT NULL,
    refresh_token TEXT        NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT NOW(),
    expires_at    timestamptz          DEFAULT NULL
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
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    device_uid  TEXT NOT NULL,
    mac_address TEXT NOT NULL,
    user_id     UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name        TEXT,
    description TEXT,
    room        TEXT,
    type        TEXT,
    status_info JSONB,
    custom_data JSONB,
    device_info JSONB,
    last_seen   timestamptz      DEFAULT NOW(),
    created_at  timestamptz      DEFAULT NOW(),
    updated_at  timestamptz      DEFAULT NOW(),
    deleted_at  timestamptz      DEFAULT NULL
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

CREATE TYPE llm_message_status AS ENUM ('pending', 'completed', 'failed');
CREATE TABLE llm_chats
(
    id         UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    user_id    UUID             NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    model      VARCHAR(64)      NOT NULL,
    title      TEXT             NOT NULL DEFAULT '',
    created_at timestamptz      NOT NULL DEFAULT NOW(),
    updated_at timestamptz      NOT NULL DEFAULT NOW(),
    deleted_at timestamptz               DEFAULT NULL
);

CREATE TABLE llm_chat_message
(
    id            UUID PRIMARY KEY   NOT NULL DEFAULT uuid_generate_v4(),
    chat_id       UUID               NOT NULL REFERENCES llm_chats (id) ON DELETE CASCADE,
    role          TEXT               NOT NULL CHECK (role IN ('user', 'assistant', 'system', 'tool')),
    model_name    VARCHAR(64)        NOT NULL,

    input_tokens  INTEGER                     DEFAULT 0,
    output_tokens INTEGER                     DEFAULT 0,
    content       TEXT                        DEFAULT '',

    tool_call_id  VARCHAR(255)                DEFAULT NULL,
    tool_name     VARCHAR(255)                DEFAULT NULL,
    tool_args     JSONB                       DEFAULT NULL,
    tool_result   JSONB                       DEFAULT NULL,

    status        llm_message_status NOT NULL DEFAULT 'completed',
    created_at    timestamptz        NOT NULL DEFAULT NOW(),
    updated_at    timestamptz        NOT NULL DEFAULT NOW(),
    deleted_at    timestamptz                 DEFAULT NULL
);
CREATE TABLE user_settings
(
    id            UUID PRIMARY KEY     DEFAULT uuid_generate_v4(),
    user_id       UUID UNIQUE NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    language      VARCHAR(10)          DEFAULT 'ru',
    timezone      VARCHAR(64)          DEFAULT 'Europe/Moscow',
    llm_config    JSONB                DEFAULT '{}'::jsonb,
    notifications JSONB                DEFAULT '{}'::jsonb,
    preferences   JSONB                DEFAULT '{}'::jsonb,
    created_at    timestamptz NOT NULL DEFAULT NOW(),
    updated_at    timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE user_llm_keys
(
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id       UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    provider      TEXT NOT NULL,
    key_encrypted TEXT NOT NULL,
    is_active     BOOLEAN          DEFAULT TRUE,
    last_used_at  timestamptz,
    created_at    timestamptz      DEFAULT NOW(),
    UNIQUE (user_id, provider)
);

GRANT CONNECT ON DATABASE "smart-home" TO "user";

GRANT USAGE ON SCHEMA public TO "user";

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "user";
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO "user";

ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO "user";
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO "user";
