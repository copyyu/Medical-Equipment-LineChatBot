package entity

import (
	"time"

	"gorm.io/gorm"
)

type NotificationSetting struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	IsEnabled bool           `gorm:"default:true" json:"is_enabled"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (NotificationSetting) TableName() string {
	return "notification_settings"
}
