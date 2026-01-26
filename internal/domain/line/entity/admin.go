package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AdminStatus enum
type AdminStatus string

const (
	AdminStatusActive   AdminStatus = "active"
	AdminStatusInactive AdminStatus = "inactive"
)

type AdminRole string

const (
	RoleAdmin      AdminRole = "admin"
	RoleSuperAdmin AdminRole = "super_admin"
	RoleStaff      AdminRole = "staff"
)

// Admin model
type Admin struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Username     string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	FullName     string    `gorm:"type:varchar(255);not null" json:"full_name"`
	Role         string    `gorm:"type:varchar(50);default:'admin';not null" json:"role"`

	// Status        AdminStatus    `gorm:"type:varchar(20);default:'active'" json:"status"`
	LastLoginAt *time.Time     `gorm:"type:timestamp" json:"last_login_at"`
	CreatedAt   time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Sessions []AdminSession `gorm:"foreignKey:AdminID;constraint:OnDelete:CASCADE" json:"-"`
}

// TableName overrides the table name
func (Admin) TableName() string {
	return "admin"
}
