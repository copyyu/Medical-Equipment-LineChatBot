package entity

import (
	"time"

	"gorm.io/gorm"
)

type MaintenanceType string

const (
	MaintenancePM MaintenanceType = "PM" // Preventive Maintenance (บำรุงรักษา)
)

type MaintenanceRecord struct {
	ID              uint            `gorm:"primaryKey" json:"id"`
	EquipmentID     uint            `gorm:"not null;index" json:"equipment_id"`
	MaintenanceType MaintenanceType `gorm:"size:10;not null" json:"maintenance_type"` // CM, PM
	MaintenanceDate time.Time       `gorm:"not null;index" json:"maintenance_date"`
	Cost            float64         `gorm:"type:decimal(15,2);default:0" json:"cost"` // ใช้สำหรับคำนวณ Total of Cost
	Description     string          `gorm:"type:text" json:"description"`
	Technician      string          `gorm:"size:100" json:"technician"`

	// Status field removed as it is moved to Ticket entity

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Equipment Equipment `gorm:"foreignKey:EquipmentID" json:"equipment,omitempty"`
}

func (MaintenanceRecord) TableName() string {
	return "maintenance_records"
}
