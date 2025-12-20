package entity

import (
	"time"

	"gorm.io/gorm"
)

type Brand struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null;uniqueIndex" json:"name"` // AIRCAST
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Models []EquipmentModel `gorm:"foreignKey:BrandID" json:"models,omitempty"`
}

func (Brand) TableName() string {
	return "brands"
}
