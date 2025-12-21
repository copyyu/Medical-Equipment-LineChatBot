package entity

import (
	"time"

	"gorm.io/gorm"
)

type ECRIRiskLevel string

const (
	RiskHigh   ECRIRiskLevel = "HIGH"
	RiskMedium ECRIRiskLevel = "MEDIUM"
	RiskLow    ECRIRiskLevel = "LOW"
)

type EquipmentCategory struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Name           string         `gorm:"size:300;not null;uniqueIndex" json:"name"`
	ECRIRisk       ECRIRiskLevel  `gorm:"size:20;default:'MEDIUM'" json:"ecri_risk"`
	Classification string         `gorm:"size:150" json:"classification"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Models []EquipmentModel `gorm:"foreignKey:CategoryID" json:"models,omitempty"`
}

func (EquipmentCategory) TableName() string {
	return "equipment_categories"
}
