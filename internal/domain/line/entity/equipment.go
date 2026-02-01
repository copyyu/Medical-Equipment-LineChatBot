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

type Equipment struct {
	ID           uint    `gorm:"primaryKey" json:"id"`
	IDCode       string  `gorm:"size:100;uniqueIndex" json:"id_code"`
	SerialNo     *string `gorm:"size:150" json:"serial_no"`
	ModelID      uint    `gorm:"not null;index" json:"model_id"`
	DepartmentID uint    `gorm:"not null;index" json:"department_id"`
	AssessmentID *string `gorm:"size:100" json:"assessment_id"`

	Status AssetStatus `gorm:"size:50;default:active" json:"status"` // Asset Status

	ReceiveDate   *time.Time `json:"receive_date"`                                       // Receive Date
	PurchasePrice float64    `gorm:"type:decimal(15,2);default:0" json:"purchase_price"` // Purchase price

	EquipmentAge          float64    `gorm:"type:decimal(10,2);default:0" json:"equipment_age"`          // Equipment Age (ปี)
	ComputeDate           *time.Time `json:"compute_date"`                                               // Compute Date
	LifeExpectancy        float64    `gorm:"type:decimal(10,2);default:10" json:"life_expectancy"`       // Life Expect (ปี)
	RemainLife            float64    `gorm:"type:decimal(10,2);default:0" json:"remain_life"`            // Remain Life (ปี)
	UsefulLifetimePercent float64    `gorm:"type:decimal(5,2);default:0" json:"useful_lifetime_percent"` // % of useful lifetime
	ReplacementYear       *int       `json:"replacement_year"`                                           // Replacement Year

	Technology      *float64 `gorm:"type:decimal(5,2)" json:"technology"`       // Technology
	UsageStatistics *float64 `gorm:"type:decimal(5,2)" json:"usage_statistics"` // Usage Statistics
	Efficiency      *float64 `gorm:"type:decimal(5,2)" json:"efficiency"`       // Efficiency
	Others          *string  `gorm:"type:text" json:"others"`                   // Others

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
