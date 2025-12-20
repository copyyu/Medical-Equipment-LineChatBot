package entity

import (
	"time"

	"gorm.io/gorm"
)

// Department - แผนก/หน่วยงาน
type Department struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null;uniqueIndex" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Equipments []Equipment `gorm:"foreignKey:DepartmentID" json:"equipments,omitempty"`
}

func (Department) TableName() string {
	return "departments"
}
