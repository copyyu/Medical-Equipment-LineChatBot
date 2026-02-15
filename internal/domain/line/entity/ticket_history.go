package entity

import (
	"time"

	"gorm.io/gorm"
)

type TicketHistoryAction string

const (
	ActionCreated       TicketHistoryAction = "created"
	ActionStatusChanged TicketHistoryAction = "status_changed"
	ActionUpdated       TicketHistoryAction = "updated"
	ActionCommented     TicketHistoryAction = "commented"
	ActionAttached      TicketHistoryAction = "attached"
	ActionCancelled     TicketHistoryAction = "cancelled"
)

type TicketHistory struct {
	ID        uint                `gorm:"primaryKey" json:"id"`
	TicketID  uint                `gorm:"not null;index" json:"ticket_id"`
	AdminID   *uint               `gorm:"index" json:"admin_id"`
	ChangedBy string              `gorm:"size:255" json:"changed_by"` // Admin username who made the change
	Action    TicketHistoryAction `gorm:"size:50;not null" json:"action"`
	Field     *string             `gorm:"size:100" json:"field"`
	OldValue  *string             `gorm:"type:text" json:"old_value"`
	NewValue  *string             `gorm:"type:text" json:"new_value"`
	Note      *string             `gorm:"type:text" json:"note"`
	IsSystem  bool                `gorm:"default:false" json:"is_system"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Ticket Ticket `gorm:"foreignKey:TicketID" json:"ticket,omitempty"`
	Admin  *Admin `gorm:"foreignKey:AdminID" json:"admin,omitempty"`
}

func (TicketHistory) TableName() string {
	return "ticket_histories"
}
