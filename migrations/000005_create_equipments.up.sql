CREATE TABLE IF NOT EXISTS equipments (
    id                    BIGSERIAL PRIMARY KEY,
    id_code               VARCHAR(100) NOT NULL,

    -- Basic Info
    asset_type_name       VARCHAR(200),
    asset_name            VARCHAR(300),
    asset_id              VARCHAR(100),
    serial_no             VARCHAR(150),
    ecri_code             VARCHAR(100),

    -- Relations
    model_id              BIGINT  NOT NULL REFERENCES equipment_models(id) ON DELETE RESTRICT,
    department_id         BIGINT  NOT NULL REFERENCES departments(id) ON DELETE RESTRICT,

    -- Status
    status                VARCHAR(50)  NOT NULL DEFAULT 'active',
    asset_status_internal VARCHAR(100),
    rental_status         VARCHAR(100),
    borrow_status         VARCHAR(100),

    -- Location
    building              VARCHAR(200),
    floor                 VARCHAR(100),
    room                  VARCHAR(100),
    phone_no              VARCHAR(50),

    -- Business
    business_name         VARCHAR(200),
    item_no               VARCHAR(100),
    sku_no                VARCHAR(100),

    -- Dates
    receive_date          TIMESTAMPTZ,
    purchase_date         TIMESTAMPTZ,
    registration_date     TIMESTAMPTZ,
    purchase_price        DECIMAL(15,2) NOT NULL DEFAULT 0,

    -- Lifecycle
    life_expectancy       DECIMAL(10,2) NOT NULL DEFAULT 10,
    equipment_age         DECIMAL(10,2) NOT NULL DEFAULT 0,
    remain_life           DECIMAL(10,2) NOT NULL DEFAULT 0,
    replacement_year      INT,

    -- Warranty
    warranty_period        VARCHAR(100),
    warranty_start_date    TIMESTAMPTZ,
    warranty_end_date      TIMESTAMPTZ,
    warranty_pm            VARCHAR(200),
    warranty_cal           VARCHAR(200),

    -- PM & Calibration
    last_pm_date           TIMESTAMPTZ,
    last_cal_date          TIMESTAMPTZ,
    pm_period              VARCHAR(100),
    cal_period             VARCHAR(100),
    vendor_pm              VARCHAR(200),
    vendor_cal             VARCHAR(200),

    -- Power & Technical
    power_consumption      VARCHAR(100),

    -- Procurement
    supplier               VARCHAR(200),
    ownership              VARCHAR(200),
    po_no                  VARCHAR(100),
    contract_no            VARCHAR(100),
    invoice_no             VARCHAR(100),
    document_no            VARCHAR(100),
    tor_no                 VARCHAR(100),
    manufacturing_country  VARCHAR(100),

    -- Financial
    revenue_per_month      DECIMAL(15,2),

    -- Misc
    remark                 TEXT,
    approved_by            VARCHAR(200),
    nsmart_item_code       VARCHAR(100),
    updated_by             VARCHAR(200),

    -- Timestamps
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_equipments_id_code    ON equipments (id_code) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_equipments_model_id          ON equipments (model_id);
CREATE INDEX IF NOT EXISTS idx_equipments_department_id     ON equipments (department_id);
CREATE INDEX IF NOT EXISTS idx_equipments_deleted_at        ON equipments (deleted_at);
