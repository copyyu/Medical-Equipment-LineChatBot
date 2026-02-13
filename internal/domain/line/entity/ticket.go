package entity

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// TicketStatus represents the status of a ticket
type TicketStatus string

const (
	TicketStatusInProgress      TicketStatus = "in_progress"           // กำลังดำเนินการ
	TicketStatusCompleted       TicketStatus = "return_equipment_back" // ส่งคืนเครื่องแล้ว (งานเสร็จ)
	TicketStatusSendToOutsource TicketStatus = "send_to_outsource"     // ส่งซ่อมภายนอก
)

// GetStatusText returns Thai text for ticket status
func (s TicketStatus) GetStatusText() string {
	switch s {
	case TicketStatusInProgress:
		return "กำลังดำเนินการ"
	case TicketStatusCompleted:
		return "ส่งคืนเครื่องแล้ว"
	case TicketStatusSendToOutsource:
		return "ส่งซ่อมภายนอก"
	default:
		return "ไม่ทราบสถานะ"
	}
}

// GetColor returns hex color for ticket status
func (s TicketStatus) GetColor() string {
	switch s {
	case TicketStatusInProgress:
		return "#42A5F5" // Blue
	case TicketStatusCompleted:
		return "#66BB6A" // Green
	case TicketStatusSendToOutsource:
		return "#FFA726" // Orange
	default:
		return "#78909C" // Grey
	}
}

type TicketPriority string

const (
	PriorityLow    TicketPriority = "low"    // ต่ำ
	PriorityMedium TicketPriority = "medium" // ปานกลาง
	PriorityHigh   TicketPriority = "high"   // สูง
	PriorityUrgent TicketPriority = "urgent" // เร่งด่วน
)

func (p TicketPriority) GetPriorityText() string {
	switch p {
	case PriorityLow:
		return "ต่ำ"
	case PriorityMedium:
		return "ปานกลาง"
	case PriorityHigh:
		return "สูง"
	case PriorityUrgent:
		return "เร่งด่วน"
	default:
		return "ไม่ระบุ"
	}
}

// GetColor returns hex color for ticket priority
func (p TicketPriority) GetColor() string {
	switch p {
	case PriorityLow:
		return "#78909C" // Grey
	case PriorityMedium:
		return "#42A5F5" // Blue
	case PriorityHigh:
		return "#FFA726" // Orange
	case PriorityUrgent:
		return "#EF5350" // Red
	default:
		return "#78909C" // Grey
	}
}

// Ticket represents a repair/issue ticket in the system
type Ticket struct {
	// Core Identification
	ID       uint   `gorm:"primaryKey" json:"id"`
	TicketNo string `gorm:"size:50;uniqueIndex;not null" json:"ticket_no"`

	// Ticket Content

	Description *string        `gorm:"type:text" json:"description"`
	CategoryID  uint           `gorm:"not null;index" json:"category_id"`
	Priority    TicketPriority `gorm:"size:20;default:'medium'" json:"priority"`

	// Equipment Information
	EquipmentID   *uint   `gorm:"index" json:"equipment_id"`
	EquipmentName *string `gorm:"size:300" json:"equipment_name"`
	Location      *string `gorm:"size:300" json:"location"`

	// Reporter Information
	ReporterID       *uint   `gorm:"index" json:"reporter_id"`
	ReporterName     string  `gorm:"size:200" json:"reporter_name"`
	ReporterLineID   *string `gorm:"size:100" json:"reporter_line_id"`
	DepartmentID     *uint   `gorm:"index" json:"department_id"`
	ContactInfo      *string `gorm:"size:300" json:"contact_info"`
	ReporterPhotoURL *string `gorm:"size:500" json:"reporter_photo_url"`

	// Status Field
	Status TicketStatus `gorm:"size:50;default:'in_progress'" json:"status"`

	// Timestamps
	ReportedAt  time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"reported_at"`
	StartedAt   *time.Time     `json:"started_at"`
	CompletedAt *time.Time     `json:"completed_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Category   TicketCategory  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Equipment  *Equipment      `gorm:"foreignKey:EquipmentID" json:"equipment,omitempty"`
	Department *Department     `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	Reporter   *Admin          `gorm:"foreignKey:ReporterID" json:"reporter,omitempty"`
	Histories  []TicketHistory `gorm:"foreignKey:TicketID" json:"histories,omitempty"`
}

// TableName returns the database table name for Ticket
func (Ticket) TableName() string {
	return "tickets"
}

// GetDurationHours returns the duration in hours from reported to completed (or now)
func (t *Ticket) GetDurationHours() float64 {
	endTime := time.Now()
	if t.CompletedAt != nil {
		endTime = *t.CompletedAt
	}
	return endTime.Sub(t.ReportedAt).Hours()
}

// ticketCounter is used to generate sequential ticket numbers
var ticketCounter uint64 = 0

// GenerateTicketNumber generates a unique ticket number with format REQ-YYYY-XXXXX
func GenerateTicketNumber() string {
	year := time.Now().Year()
	ticketCounter++
	return fmt.Sprintf("REQ-%d-%05d", year, ticketCounter)
}
