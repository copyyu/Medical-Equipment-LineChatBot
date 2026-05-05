-- Enable uuid-ossp extension for gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS admin (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    username      VARCHAR(100) NOT NULL,
    email         VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name     VARCHAR(255) NOT NULL,
    role          VARCHAR(50)  NOT NULL DEFAULT 'admin',
    last_login_at TIMESTAMPTZ,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_admin_username ON admin (username) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_admin_email    ON admin (email)    WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_admin_deleted_at      ON admin (deleted_at);

CREATE TABLE IF NOT EXISTS admin_session (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    admin_id   UUID         NOT NULL REFERENCES admin(id) ON DELETE CASCADE,
    token      VARCHAR(255) NOT NULL,
    ip_address VARCHAR(45)  DEFAULT '',
    expires_at TIMESTAMPTZ  NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_admin_session_token ON admin_session (token);
CREATE INDEX IF NOT EXISTS idx_admin_session_admin_id     ON admin_session (admin_id);
CREATE INDEX IF NOT EXISTS idx_admin_session_expires_at   ON admin_session (expires_at);
