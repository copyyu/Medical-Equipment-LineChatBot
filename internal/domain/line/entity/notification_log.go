package entity

import (
	"time"

	"gorm.io/gorm"
)

type NotificationStatus string

const (
	NotificationStatusSent   NotificationStatus = "SENT"
	NotificationStatusFailed NotificationStatus = "FAILED"
)

type NotificationLog struct {
	ID          uint               `gorm:"primaryKey" json:"id"`
	EquipmentID uint               `gorm:"not null;index" json:"equipment_id"`
	NotifyRound string             `gorm:"size:20" json:"notify_round"` // JUNE, AUGUST
	Message     string             `gorm:"type:text;not null" json:"message"`
	Status      NotificationStatus `gorm:"size:20;default:'SENT'" json:"status"`
	SentAt      time.Time          `gorm:"not null;index" json:"sent_at"`
	ErrorMsg    *string            `gorm:"type:text" json:"error_msg"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	DeletedAt   gorm.DeletedAt     `gorm:"index" json:"-"`

	// Relations
	Equipment Equipment `gorm:"foreignKey:EquipmentID" json:"equipment,omitempty"`
}

func (NotificationLog) TableName() string {
	return "notification_logs"
}
