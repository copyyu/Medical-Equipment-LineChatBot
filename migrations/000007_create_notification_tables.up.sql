CREATE TABLE IF NOT EXISTS notification_settings (
    id          BIGSERIAL PRIMARY KEY,
    is_enabled  BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_notification_settings_deleted_at ON notification_settings (deleted_at);

CREATE TABLE IF NOT EXISTS notification_logs (
    id           BIGSERIAL PRIMARY KEY,
    equipment_id BIGINT             NOT NULL REFERENCES equipments(id) ON DELETE CASCADE,
    notify_round VARCHAR(20)        DEFAULT '',
    message      TEXT               NOT NULL,
    status       VARCHAR(20)        NOT NULL DEFAULT 'SENT',
    sent_at      TIMESTAMPTZ        NOT NULL,
    error_msg    TEXT,
    created_at   TIMESTAMPTZ        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMPTZ        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_notification_logs_equipment_id ON notification_logs (equipment_id);
CREATE INDEX IF NOT EXISTS idx_notification_logs_sent_at      ON notification_logs (sent_at);
CREATE INDEX IF NOT EXISTS idx_notification_logs_deleted_at   ON notification_logs (deleted_at);
