package dto

import (
	"time"
)

// EquipmentListRequest represents the request parameters for equipment list
type EquipmentListRequest struct {
	Page    int    `query:"page"`
	Limit   int    `query:"limit"`
	Status  string `query:"status"`
	Search  string `query:"search"`
	SortBy  string `query:"sort_by"`
	SortDir string `query:"sort_dir"`
}

// EquipmentListItem represents a single equipment item for the table
type EquipmentListItem struct {
	ID         string `json:"id"`          // ID Code
	Name       string `json:"name"`        // Model name or description
	Category   string `json:"category"`    // Equipment category
	Status     string `json:"status"`      // Asset status
	Location   string `json:"location"`    // Department/location
	LastCheck  string `json:"last_check"`  // Last maintenance date
	Expiry     string `json:"expiry"`      // Based on replacement year or remain life
	IsExpiring bool   `json:"is_expiring"` // Flag for expiring soon
}

// EquipmentListResponse represents the paginated equipment list response
type EquipmentListResponse struct {
	Data       []EquipmentListItem `json:"data"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	TotalPages int                 `json:"total_pages"`
}

type CreateEquipmentRequest struct {
	IDCode                string   `json:"id_code" binding:"required"`
	SerialNo              string   `json:"serial_no" binding:"required"`
	AssessmentID          string   `json:"assessment_id"`
	Department            string   `json:"department" binding:"required"`
	Brand                 string   `json:"brand" binding:"required"`
	Model                 string   `json:"model" binding:"required"`
	Category              string   `json:"category" binding:"required"`
	ReceiveDate           string   `json:"receive_date" binding:"required"` // Format: YYYY-MM-DD
	PurchasePrice         float64  `json:"purchase_price" binding:"required"`
	EquipmentAge          float64  `json:"equipment_age"`
	ComputeDate           string   `json:"compute_date"` // Format: YYYY-MM-DD
	LifeExpectancy        float64  `json:"life_expectancy"`
	RemainLife            float64  `json:"remain_life"`
	UsefulLifetimePercent float64  `json:"useful_lifetime_percent"`
	ReplacementYear       int      `json:"replacement_year"`
	Technology            *float64 `json:"technology"`
	UsageStatistics       *float64 `json:"usage_statistics"`
	Efficiency            *float64 `json:"efficiency"`
	Others                string   `json:"others"`
}

// EquipmentResponse - ใช้ snake_case ตาม entity
type EquipmentResponse struct {
	ID                    uint               `json:"id"`
	IDCode                string             `json:"id_code"`
	SerialNo              *string            `json:"serial_no"`
	AssessmentID          *string            `json:"assessment_id"`
	Status                string             `json:"status"`
	ReceiveDate           *time.Time         `json:"receive_date"`
	PurchasePrice         float64            `json:"purchase_price"`
	EquipmentAge          float64            `json:"equipment_age"`
	ComputeDate           *time.Time         `json:"compute_date"`
	LifeExpectancy        float64            `json:"life_expectancy"`
	RemainLife            float64            `json:"remain_life"`
	UsefulLifetimePercent float64            `json:"useful_lifetime_percent"`
	ReplacementYear       *int               `json:"replacement_year"`
	Technology            *float64           `json:"technology"`
	UsageStatistics       *float64           `json:"usage_statistics"`
	Efficiency            *float64           `json:"efficiency"`
	Others                *string            `json:"others"`
	Model                 *EquipmentModelDTO `json:"model,omitempty"`
	Department            *DepartmentDTO     `json:"department,omitempty"`
	CreatedAt             time.Time          `json:"created_at"`
	UpdatedAt             time.Time          `json:"updated_at"`
}

type EquipmentModelDTO struct {
	ID                    uint         `json:"id"`
	ModelName             string       `json:"model_name"`
	DefaultLifeExpectancy float64      `json:"default_life_expectancy"`
	Brand                 *BrandDTO    `json:"brand,omitempty"`
	Category              *CategoryDTO `json:"category,omitempty"`
}

// EquipmentUpdateRequest represents the request body for updating equipment
type EquipmentUpdateRequest struct {
	Status      string `json:"status"`       // Asset Status (active, defective, etc.)
	Location    string `json:"location"`     // Department name
	ComputeDate string `json:"compute_date"` // Last check date (YYYY-MM-DD)
}

// EquipmentDetailResponse represents a single equipment detail for GET by ID
type EquipmentDetailResponse struct {
	ID           string `json:"id"`            // ID Code
	Name         string `json:"name"`          // Model name
	Category     string `json:"category"`      // Equipment category
	Status       string `json:"status"`        // Asset status
	Location     string `json:"location"`      // Department/location
	LastCheck    string `json:"last_check"`    // Last maintenance date
	Expiry       string `json:"expiry"`        // Expiry year
	IsExpiring   bool   `json:"is_expiring"`   // Flag for expiring soon
	SerialNo     string `json:"serial_no"`     // Serial number
	Brand        string `json:"brand"`         // Brand name
	DepartmentID uint   `json:"department_id"` // Department ID for updates
}
