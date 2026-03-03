package entity

import (
	"time"

	"gorm.io/gorm"
)

// AssetStatus represents the status of equipment
type AssetStatus string

const (
	AssetStatusActive            AssetStatus = "active"               // ใช้งานอยู่
	AssetStatusDefective         AssetStatus = "defective"            // ชำรุด
	AssetStatusWaitDecom         AssetStatus = "wait_decom"           // รอปลดระวาง
	AssetStatusDecommission      AssetStatus = "decommission"         // ปลดระวางแล้ว
	AssetStatusActiveReadyToSell AssetStatus = "active_ready_to_sell" // พร้อมขาย
	AssetStatusMissing           AssetStatus = "missing"              // สูญหาย
	AssetStatusPlanToReplace     AssetStatus = "plan_to_replace"      // รอเปลี่ยนใหม่
)

// GetStatusText returns Thai text for asset status
func (s AssetStatus) GetStatusText() string {
	switch s {
	case AssetStatusActive:
		return "ใช้งานอยู่"
	case AssetStatusDefective:
		return "ชำรุด"
	case AssetStatusWaitDecom:
		return "รอปลดระวาง"
	case AssetStatusDecommission:
		return "ปลดระวางแล้ว"
	case AssetStatusActiveReadyToSell:
		return "พร้อมขาย"
	case AssetStatusMissing:
		return "สูญหาย"
	case AssetStatusPlanToReplace:
		return "รอเปลี่ยนใหม่"
	default:
		return "ไม่ทราบสถานะ"
	}
}

// GetColor returns hex color for asset status
func (s AssetStatus) GetColor() string {
	switch s {
	case AssetStatusActive:
		return "#4CAF50" // Green
	case AssetStatusDefective:
		return "#EF5350" // Red
	case AssetStatusWaitDecom:
		return "#FFA726" // Orange
	case AssetStatusDecommission:
		return "#78909C" // Grey
	case AssetStatusActiveReadyToSell:
		return "#42A5F5" // Blue
	case AssetStatusMissing:
		return "#E53935" // Dark Red
	case AssetStatusPlanToReplace:
		return "#AB47BC" // Purple
	default:
		return "#78909C" // Grey
	}
}

type Equipment struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	IDCode string `gorm:"size:100;uniqueIndex" json:"id_code"`

	// === Basic Info ===
	AssetTypeName *string `gorm:"size:200" json:"asset_type_name"`
	AssetName     *string `gorm:"size:300" json:"asset_name"`
	AssetID       *string `gorm:"size:100" json:"asset_id"`
	SerialNo      *string `gorm:"size:150" json:"serial_no"`
	ECRICode      *string `gorm:"size:100" json:"ecri_code"`

	// === Relations (resolved via Model/Department) ===
	ModelID      uint `gorm:"not null;index" json:"model_id"`
	DepartmentID uint `gorm:"not null;index" json:"department_id"`

	// === Status ===
	Status              AssetStatus `gorm:"size:50;default:active" json:"status"`
	AssetStatusInternal *string     `gorm:"size:100" json:"asset_status_internal"`
	RentalStatus        *string     `gorm:"size:100" json:"rental_status"`
	BorrowStatus        *string     `gorm:"size:100" json:"borrow_status"`

	// === Location ===
	Building *string `gorm:"size:200" json:"building"`
	Floor    *string `gorm:"size:100" json:"floor"`
	Room     *string `gorm:"size:100" json:"room"`
	PhoneNo  *string `gorm:"size:50" json:"phone_no"`

	// === Business ===
	BusinessName *string `gorm:"size:200" json:"business_name"`
	ItemNo       *string `gorm:"size:100" json:"item_no"`
	SKUNo        *string `gorm:"size:100" json:"sku_no"`

	// === Dates ===
	ReceiveDate      *time.Time `json:"receive_date"`
	PurchaseDate     *time.Time `json:"purchase_date"`
	RegistrationDate *time.Time `json:"registration_date"`
	PurchasePrice    float64    `gorm:"type:decimal(15,2);default:0" json:"purchase_price"`

	// === Lifecycle (LifeExpectancy from Excel, rest computed) ===
	LifeExpectancy  float64 `gorm:"type:decimal(10,2);default:10" json:"life_expectancy"`
	EquipmentAge    float64 `gorm:"type:decimal(10,2);default:0" json:"equipment_age"` // COMPUTED: now - ReceiveDate
	RemainLife      float64 `gorm:"type:decimal(10,2);default:0" json:"remain_life"`   // COMPUTED: LifeExpectancy - EquipmentAge
	ReplacementYear *int    `json:"replacement_year"`                                  // COMPUTED: ReceiveDate.Year + LifeExpectancy

	// === Warranty ===
	WarrantyPeriod    *string    `gorm:"size:100" json:"warranty_period"`
	WarrantyStartDate *time.Time `json:"warranty_start_date"`
	WarrantyEndDate   *time.Time `json:"warranty_end_date"`
	WarrantyPM        *string    `gorm:"size:200" json:"warranty_pm"`
	WarrantyCal       *string    `gorm:"size:200" json:"warranty_cal"`

	// === PM & Calibration ===
	LastPMDate  *time.Time `json:"last_pm_date"`
	LastCalDate *time.Time `json:"last_cal_date"`
	PMPeriod    *string    `gorm:"size:100" json:"pm_period"`
	CalPeriod   *string    `gorm:"size:100" json:"cal_period"`
	VendorPM    *string    `gorm:"size:200" json:"vendor_pm"`
	VendorCal   *string    `gorm:"size:200" json:"vendor_cal"`

	// === Power & Technical ===
	PowerConsumption *string `gorm:"size:100" json:"power_consumption"`

	// === Procurement ===
	Supplier             *string `gorm:"size:200" json:"supplier"`
	Ownership            *string `gorm:"size:200" json:"ownership"`
	PoNo                 *string `gorm:"size:100" json:"po_no"`
	ContractNo           *string `gorm:"size:100" json:"contract_no"`
	InvoiceNo            *string `gorm:"size:100" json:"invoice_no"`
	DocumentNo           *string `gorm:"size:100" json:"document_no"`
	TorNo                *string `gorm:"size:100" json:"tor_no"`
	ManufacturingCountry *string `gorm:"size:100" json:"manufacturing_country"`

	// === Financial ===
	RevenuePerMonth *float64 `gorm:"type:decimal(15,2)" json:"revenue_per_month"`

	// === Misc ===
	Remark         *string `gorm:"type:text" json:"remark"`
	ApprovedBy     *string `gorm:"size:200" json:"approved_by"`
	NsmartItemCode *string `gorm:"size:100" json:"nsmart_item_code"`
	UpdatedBy      *string `gorm:"size:200" json:"updated_by"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Model              EquipmentModel      `gorm:"foreignKey:ModelID" json:"model,omitempty"`
	Department         Department          `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	MaintenanceRecords []MaintenanceRecord `gorm:"foreignKey:EquipmentID" json:"maintenance_records,omitempty"`
}

func (Equipment) TableName() string {
	return "equipments"
}
