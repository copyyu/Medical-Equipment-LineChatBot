package entity

import (
	"time"

	"github.com/google/uuid"
)

// AdminSession model
type AdminSession struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AdminID   uuid.UUID `gorm:"type:uuid;not null;index" json:"admin_id"`
	Token     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"token"`
	IPAddress string    `gorm:"type:varchar(45)" json:"ip_address"`
	ExpiresAt time.Time `gorm:"type:timestamp;not null;index" json:"expires_at"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relations
	Admin Admin `gorm:"foreignKey:AdminID;constraint:OnDelete:CASCADE" json:"admin,omitempty"`
}

// TableName overrides the table name
func (AdminSession) TableName() string {
	return "admin_session"
}
