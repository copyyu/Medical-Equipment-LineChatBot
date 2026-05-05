CREATE TABLE IF NOT EXISTS ticket_categories (
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    name_en    VARCHAR(100) DEFAULT '',
    color      VARCHAR(20)  NOT NULL DEFAULT '#78909C',
    icon       VARCHAR(50)  DEFAULT '',
    is_active  BOOLEAN      NOT NULL DEFAULT TRUE,
    sort_order INT          NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_ticket_categories_name ON ticket_categories (name) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_ticket_categories_deleted_at  ON ticket_categories (deleted_at);

CREATE TABLE IF NOT EXISTS tickets (
    id                 BIGSERIAL PRIMARY KEY,
    ticket_no          VARCHAR(50)  NOT NULL,
    description        TEXT,
    category_id        BIGINT       NOT NULL REFERENCES ticket_categories(id) ON DELETE RESTRICT,
    priority           VARCHAR(20)  NOT NULL DEFAULT 'medium',
    equipment_id       BIGINT       REFERENCES equipments(id) ON DELETE SET NULL,
    equipment_name     VARCHAR(300),
    location           VARCHAR(300),
    reporter_id        BIGINT,
    reporter_name      VARCHAR(200) NOT NULL DEFAULT '',
    reporter_line_id   VARCHAR(100),
    department_id      BIGINT       REFERENCES departments(id) ON DELETE SET NULL,
    contact_info       VARCHAR(300),
    reporter_photo_url VARCHAR(500),
    status             VARCHAR(50)  NOT NULL DEFAULT 'in_process',
    reported_at        TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    started_at         TIMESTAMPTZ,
    completed_at       TIMESTAMPTZ,
    created_at         TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at         TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_tickets_ticket_no   ON tickets (ticket_no) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_tickets_category_id        ON tickets (category_id);
CREATE INDEX IF NOT EXISTS idx_tickets_equipment_id       ON tickets (equipment_id);
CREATE INDEX IF NOT EXISTS idx_tickets_reporter_id        ON tickets (reporter_id);
CREATE INDEX IF NOT EXISTS idx_tickets_department_id      ON tickets (department_id);
CREATE INDEX IF NOT EXISTS idx_tickets_deleted_at         ON tickets (deleted_at);

CREATE TABLE IF NOT EXISTS ticket_histories (
    id         BIGSERIAL PRIMARY KEY,
    ticket_id  BIGINT              NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    admin_id   BIGINT,
    changed_by VARCHAR(255)        DEFAULT '',
    action     VARCHAR(50)         NOT NULL,
    field      VARCHAR(100),
    old_value  TEXT,
    new_value  TEXT,
    note       TEXT,
    is_system  BOOLEAN             NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_ticket_histories_ticket_id  ON ticket_histories (ticket_id);
CREATE INDEX IF NOT EXISTS idx_ticket_histories_admin_id   ON ticket_histories (admin_id);
CREATE INDEX IF NOT EXISTS idx_ticket_histories_deleted_at ON ticket_histories (deleted_at);
