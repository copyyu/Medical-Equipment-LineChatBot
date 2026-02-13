package entity

import (
	"time"

	"gorm.io/gorm"
)

type TicketCategory struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `gorm:"size:100;not null;uniqueIndex" json:"name"`
	NameEn    string `gorm:"size:100" json:"name_en"`
	Color     string `gorm:"size:20;default:'#78909C'" json:"color"`
	Icon      string `gorm:"size:50" json:"icon"`
	IsActive  bool   `gorm:"default:true" json:"is_active"`
	SortOrder int    `gorm:"default:0" json:"sort_order"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Tickets []Ticket `gorm:"foreignKey:CategoryID" json:"tickets,omitempty"`
}

func (TicketCategory) TableName() string {
	return "ticket_categories"
}
