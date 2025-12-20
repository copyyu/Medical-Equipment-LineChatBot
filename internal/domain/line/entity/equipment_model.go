package entity

import (
	"time"

	"gorm.io/gorm"
)

type EquipmentModel struct {
	ID                    uint           `gorm:"primaryKey" json:"id"`
	BrandID               uint           `gorm:"not null;index" json:"brand_id"`
	CategoryID            uint           `gorm:"not null;index" json:"category_id"`
	ModelName             string         `gorm:"size:200;not null" json:"model_name"`
	DefaultLifeExpectancy float64        `gorm:"default:10" json:"default_life_expectancy"` // อายุการใช้งานมาตรฐาน (ปี)
	Specifications        string         `gorm:"type:text" json:"specifications"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Brand      Brand             `gorm:"foreignKey:BrandID" json:"brand,omitempty"`
	Category   EquipmentCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Equipments []Equipment       `gorm:"foreignKey:ModelID" json:"equipments,omitempty"`
}

func (EquipmentModel) TableName() string {
	return "equipment_models"
}
