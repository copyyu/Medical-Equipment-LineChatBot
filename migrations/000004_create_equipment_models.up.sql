CREATE TABLE IF NOT EXISTS equipment_models (
    id                      BIGSERIAL PRIMARY KEY,
    brand_id                BIGINT       NOT NULL REFERENCES brands(id) ON DELETE RESTRICT,
    category_id             BIGINT       NOT NULL REFERENCES equipment_categories(id) ON DELETE RESTRICT,
    model_name              VARCHAR(250) NOT NULL,
    default_life_expectancy DECIMAL(10,2) NOT NULL DEFAULT 10,
    specifications          TEXT         DEFAULT '',
    created_at              TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at              TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_equipment_models_brand_id    ON equipment_models (brand_id);
CREATE INDEX IF NOT EXISTS idx_equipment_models_category_id ON equipment_models (category_id);
CREATE INDEX IF NOT EXISTS idx_equipment_models_deleted_at  ON equipment_models (deleted_at);
