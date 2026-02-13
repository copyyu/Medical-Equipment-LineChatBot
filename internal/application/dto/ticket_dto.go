package dto

import "time"

// TicketListRequest represents the request for listing tickets with pagination
type TicketListRequest struct {
	Page     int    `query:"page"`
	Limit    int    `query:"limit"`
	Status   string `query:"status"`
	Priority string `query:"priority"`
	Search   string `query:"search"`
	SortBy   string `query:"sort_by"`
	SortDir  string `query:"sort_dir"`
}

// TicketListResponse represents paginated ticket list response
type TicketListResponse struct {
	Data       []TicketItemResponse `json:"data"`
	Pagination Pagination           `json:"pagination"`
}

// Pagination represents pagination
type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// TicketItemResponse represents a single ticket in list
type TicketItemResponse struct {
	ID       uint   `json:"id"`
	TicketNo string `json:"ticket_no"`

	Description  *string `json:"description"`
	CategoryID   uint    `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Priority     string  `json:"priority"`
	PriorityText string  `json:"priority_text"`
	Status       string  `json:"status"`
	StatusText   string  `json:"status_text"`

	EquipmentName    *string    `json:"equipment_name"`
	EquipmentIDCode  *string    `json:"equipment_id_code"`
	ReporterName     string     `json:"reporter_name"`
	ReporterPhotoURL *string    `json:"reporter_photo_url"`
	DepartmentName   *string    `json:"department_name"`
	ReportedAt       time.Time  `json:"reported_at"`
	CompletedAt      *time.Time `json:"completed_at"`
	CreatedAt        time.Time  `json:"created_at"`
}

// TicketDetailResponse represents full ticket detail
type TicketDetailResponse struct {
	ID       uint   `json:"id"`
	TicketNo string `json:"ticket_no"`

	Description  *string `json:"description"`
	CategoryID   uint    `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Priority     string  `json:"priority"`
	PriorityText string  `json:"priority_text"`
	Status       string  `json:"status"`
	StatusText   string  `json:"status_text"`

	EquipmentID      *uint              `json:"equipment_id"`
	EquipmentName    *string            `json:"equipment_name"`
	EquipmentIDCode  *string            `json:"equipment_id_code"`
	Location         *string            `json:"location"`
	ReporterName     string             `json:"reporter_name"`
	ReporterLineID   *string            `json:"reporter_line_id"`
	ReporterPhotoURL *string            `json:"reporter_photo_url"`
	DepartmentID     *uint              `json:"department_id"`
	DepartmentName   *string            `json:"department_name"`
	ContactInfo      *string            `json:"contact_info"`
	ReportedAt       time.Time          `json:"reported_at"`
	StartedAt        *time.Time         `json:"started_at"`
	CompletedAt      *time.Time         `json:"completed_at"`
	DurationHours    *float64           `json:"duration_hours"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
	Histories        []TicketHistoryDTO `json:"histories"`
}

// TicketHistoryDTO represents a history record
type TicketHistoryDTO struct {
	ID        uint      `json:"id"`
	Action    string    `json:"action"`
	Field     *string   `json:"field"`
	OldValue  *string   `json:"old_value"`
	NewValue  *string   `json:"new_value"`
	Note      *string   `json:"note"`
	AdminName *string   `json:"admin_name"`
	IsSystem  bool      `json:"is_system"`
	CreatedAt time.Time `json:"created_at"`
}

// UpdateTicketRequest represents ticket update request
type UpdateTicketRequest struct {
	CategoryID  *uint   `json:"category_id"`
	Priority    *string `json:"priority" validate:"omitempty,oneof=low medium high urgent"`
	Status      *string `json:"status" validate:"omitempty,oneof=in_progress return_equipment_back send_to_outsource"`
	Description *string `json:"description"`
	Note        string  `json:"note"` // For history log
}

// TicketStatsResponse represents ticket statistics
type TicketStatsResponse struct {
	Total           int64 `json:"total"`
	InProgress      int64 `json:"in_progress"`
	Completed       int64 `json:"completed"`
	SendToOutsource int64 `json:"send_to_outsource"`
}
