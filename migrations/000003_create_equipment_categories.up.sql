CREATE TABLE IF NOT EXISTS equipment_categories (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(300) NOT NULL,
    ecri_risk       VARCHAR(20)  NOT NULL DEFAULT 'MEDIUM',
    classification  VARCHAR(150) DEFAULT '',
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_equipment_categories_name ON equipment_categories (name) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_equipment_categories_deleted_at ON equipment_categories (deleted_at);
