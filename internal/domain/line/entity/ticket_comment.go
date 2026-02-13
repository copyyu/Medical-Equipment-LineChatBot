package entity

import (
	"time"

	"gorm.io/gorm"
)

type TicketComment struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	TicketID   uint   `gorm:"not null;index" json:"ticket_id"`
	AdminID    *uint  `gorm:"index" json:"admin_id"`
	Content    string `gorm:"type:text;not null" json:"content"`
	IsSystem   bool   `gorm:"default:false" json:"is_system"`
	IsInternal bool   `gorm:"default:false" json:"is_internal"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Ticket Ticket `gorm:"foreignKey:TicketID" json:"ticket,omitempty"`
	Admin  *Admin `gorm:"foreignKey:AdminID" json:"admin,omitempty"`
}

func (TicketComment) TableName() string {
	return "ticket_comments"
}
