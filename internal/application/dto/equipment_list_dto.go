package dto

import (
	"time"
)

// EquipmentListRequest represents the request parameters for equipment list
type EquipmentListRequest struct {
	Page    int    `query:"page"`
	Limit   int    `query:"limit"`
	Status  string `query:"status"`
	Search  string `query:"search"`
	SortBy  string `query:"sort_by"`
	SortDir string `query:"sort_dir"`
}

// EquipmentListItem represents a single equipment item for the table
type EquipmentListItem struct {
	ID         string  `json:"id"`          // ID Code
	Name       string  `json:"name"`        // Asset name or Model name
	Category   string  `json:"category"`    // Equipment category
	Status     string  `json:"status"`      // Asset status
	Location   string  `json:"location"`    // Department/location
	LastCheck  string  `json:"last_check"`  // Last PM date
	Expiry     string  `json:"expiry"`      // Based on replacement year or remain life
	IsExpiring bool    `json:"is_expiring"` // Flag for expiring soon
	RemainLife float64 `json:"remain_life"`
}

// EquipmentListResponse represents the paginated equipment list response
type EquipmentListResponse struct {
	Data       []EquipmentListItem `json:"data"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	TotalPages int                 `json:"total_pages"`
}

type CreateEquipmentRequest struct {
	IDCode            string  `json:"id_code" binding:"required"`
	SerialNo          string  `json:"serial_no" binding:"required"`
	Department        string  `json:"department" binding:"required"`
	Brand             string  `json:"brand" binding:"required"`
	Model             string  `json:"model" binding:"required"`
	Category          string  `json:"category" binding:"required"`
	AssetTypeName     string  `json:"asset_type_name"`
	ECRICode          string  `json:"ecri_code"`
	AssetName         string  `json:"asset_name"`
	AssetID           string  `json:"asset_id"`
	ReceiveDate       string  `json:"receive_date" binding:"required"`
	PurchaseDate      string  `json:"purchase_date"`
	PurchasePrice     float64 `json:"purchase_price" binding:"required"`
	LifeExpectancy    float64 `json:"life_expectancy"`
	Building          string  `json:"building"`
	Floor             string  `json:"floor"`
	Room              string  `json:"room"`
	WarrantyPeriod    string  `json:"warranty_period"`
	WarrantyStartDate string  `json:"warranty_start_date"`
	WarrantyEndDate   string  `json:"warranty_end_date"`
	Remark            string  `json:"remark"`
}

// EquipmentResponse - ใช้ snake_case ตาม entity
type EquipmentResponse struct {
	ID                   uint               `json:"id"`
	IDCode               string             `json:"id_code"`
	SerialNo             *string            `json:"serial_no"`
	ECRICode             *string            `json:"ecri_code"`
	Status               string             `json:"status"`
	AssetTypeName        *string            `json:"asset_type_name"`
	AssetName            *string            `json:"asset_name"`
	AssetID              *string            `json:"asset_id"`
	AssetStatusInternal  *string            `json:"asset_status_internal"`
	RentalStatus         *string            `json:"rental_status"`
	BorrowStatus         *string            `json:"borrow_status"`
	Building             *string            `json:"building"`
	Floor                *string            `json:"floor"`
	Room                 *string            `json:"room"`
	PhoneNo              *string            `json:"phone_no"`
	BusinessName         *string            `json:"business_name"`
	ItemNo               *string            `json:"item_no"`
	SKUNo                *string            `json:"sku_no"`
	ReceiveDate          *time.Time         `json:"receive_date"`
	PurchaseDate         *time.Time         `json:"purchase_date"`
	RegistrationDate     *time.Time         `json:"registration_date"`
	PurchasePrice        float64            `json:"purchase_price"`
	EquipmentAge         float64            `json:"equipment_age"`
	LifeExpectancy       float64            `json:"life_expectancy"`
	RemainLife           float64            `json:"remain_life"`
	ReplacementYear      *int               `json:"replacement_year"`
	WarrantyPeriod       *string            `json:"warranty_period"`
	WarrantyStartDate    *time.Time         `json:"warranty_start_date"`
	WarrantyEndDate      *time.Time         `json:"warranty_end_date"`
	WarrantyPM           *string            `json:"warranty_pm"`
	WarrantyCal          *string            `json:"warranty_cal"`
	LastPMDate           *time.Time         `json:"last_pm_date"`
	LastCalDate          *time.Time         `json:"last_cal_date"`
	PMPeriod             *string            `json:"pm_period"`
	CalPeriod            *string            `json:"cal_period"`
	VendorPM             *string            `json:"vendor_pm"`
	VendorCal            *string            `json:"vendor_cal"`
	PowerConsumption     *string            `json:"power_consumption"`
	Supplier             *string            `json:"supplier"`
	Ownership            *string            `json:"ownership"`
	PoNo                 *string            `json:"po_no"`
	ContractNo           *string            `json:"contract_no"`
	InvoiceNo            *string            `json:"invoice_no"`
	DocumentNo           *string            `json:"document_no"`
	TorNo                *string            `json:"tor_no"`
	ManufacturingCountry *string            `json:"manufacturing_country"`
	RevenuePerMonth      *float64           `json:"revenue_per_month"`
	Remark               *string            `json:"remark"`
	ApprovedBy           *string            `json:"approved_by"`
	NsmartItemCode       *string            `json:"nsmart_item_code"`
	UpdatedBy            *string            `json:"updated_by"`
	Model                *EquipmentModelDTO `json:"model,omitempty"`
	Department           *DepartmentDTO     `json:"department,omitempty"`
	CreatedAt            time.Time          `json:"created_at"`
	UpdatedAt            time.Time          `json:"updated_at"`
}

type EquipmentModelDTO struct {
	ID                    uint         `json:"id"`
	ModelName             string       `json:"model_name"`
	DefaultLifeExpectancy float64      `json:"default_life_expectancy"`
	Brand                 *BrandDTO    `json:"brand,omitempty"`
	Category              *CategoryDTO `json:"category,omitempty"`
}

// EquipmentUpdateRequest represents the request body for updating equipment
type EquipmentUpdateRequest struct {
	Status     string `json:"status"`      // Asset Status (active, defective, etc.)
	Location   string `json:"location"`    // Department name
	ExpiryDate string `json:"expiry_date"` // Expiry date (YYYY-MM-DD) - will calculate RemainLife
}

// EquipmentDetailResponse represents a single equipment detail for GET by ID
type EquipmentDetailResponse struct {
	ID           string `json:"id"`            // ID Code
	Name         string `json:"name"`          // Asset name or Model name
	Category     string `json:"category"`      // Equipment category
	Status       string `json:"status"`        // Asset status
	Location     string `json:"location"`      // Department/location
	LastCheck    string `json:"last_check"`    // Last PM date
	Expiry       string `json:"expiry"`        // Expiry year
	IsExpiring   bool   `json:"is_expiring"`   // Flag for expiring soon
	SerialNo     string `json:"serial_no"`     // Serial number
	Brand        string `json:"brand"`         // Brand name
	DepartmentID uint   `json:"department_id"` // Department ID for updates
}
