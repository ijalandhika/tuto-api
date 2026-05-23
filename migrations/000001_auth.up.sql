-- migrate:up

CREATE TYPE actor_type AS ENUM ('parent', 'kid');

CREATE TABLE parents (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    email            VARCHAR     NOT NULL UNIQUE,
    password_hash    VARCHAR     NOT NULL,
    display_name     VARCHAR     NOT NULL,
    marketing_opt_in BOOLEAN     NOT NULL DEFAULT false,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE auth_sessions (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_type   actor_type  NOT NULL,
    actor_id     UUID        NOT NULL,
    token_hash   VARCHAR     NOT NULL UNIQUE,
    device_id    VARCHAR,
    expires_at   TIMESTAMPTZ NOT NULL,
    revoked_at   TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE parent_gate_attempts (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_id    UUID        NOT NULL REFERENCES parents(id) ON DELETE CASCADE,
    a            INT         NOT NULL,
    b            INT         NOT NULL,
    submitted    INT         NOT NULL,
    succeeded    BOOLEAN     NOT NULL DEFAULT false,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_auth_sessions_actor ON auth_sessions(actor_type, actor_id);
CREATE INDEX idx_auth_sessions_token ON auth_sessions(token_hash);
CREATE INDEX idx_parent_gate_parent  ON parent_gate_attempts(parent_id);

