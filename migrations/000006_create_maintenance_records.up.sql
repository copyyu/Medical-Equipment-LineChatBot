CREATE TABLE IF NOT EXISTS maintenance_records (
    id               BIGSERIAL PRIMARY KEY,
    equipment_id     BIGINT       NOT NULL REFERENCES equipments(id) ON DELETE CASCADE,
    maintenance_type VARCHAR(10)  NOT NULL,
    maintenance_date TIMESTAMPTZ  NOT NULL,
    cost             DECIMAL(15,2) NOT NULL DEFAULT 0,
    description      TEXT         DEFAULT '',
    technician       VARCHAR(100) DEFAULT '',
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at       TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_maintenance_records_equipment_id     ON maintenance_records (equipment_id);
CREATE INDEX IF NOT EXISTS idx_maintenance_records_maintenance_date ON maintenance_records (maintenance_date);
CREATE INDEX IF NOT EXISTS idx_maintenance_records_deleted_at       ON maintenance_records (deleted_at);
