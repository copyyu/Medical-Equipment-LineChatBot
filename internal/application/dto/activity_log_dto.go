package dto

import "time"

// ActivityLogListRequest represents query params for activity log listing
type ActivityLogListRequest struct {
	Page       int    `query:"page"`
	Limit      int    `query:"limit"`
	Search     string `query:"search"`
	FromStatus string `query:"from_status"`
	ToStatus   string `query:"to_status"`
	StartDate  string `query:"start_date"`
	EndDate    string `query:"end_date"`
}

// ActivityLogItem represents a single activity log entry
type ActivityLogItem struct {
	ID            uint      `json:"id"`
	TicketID      uint      `json:"ticket_id"`
	TicketNo      string    `json:"ticket_no"`
	EquipmentName string    `json:"equipment_name"`
	UserName      string    `json:"user_name"`
	FromStatus    string    `json:"from_status"`
	ToStatus      string    `json:"to_status"`
	Note          string    `json:"note"`
	ChangedAt     time.Time `json:"changed_at"`
}

// ActivityLogListResponse represents the paginated activity log response
type ActivityLogListResponse struct {
	Data       []ActivityLogItem `json:"data"`
	Pagination Pagination        `json:"pagination"`
}

// ActivityLogStatsResponse represents aggregate stats
type ActivityLogStatsResponse struct {
	TotalChanges      int64  `json:"total_changes"`
	TodayChanges      int64  `json:"today_changes"`
	WeekChanges       int64  `json:"week_changes"`
	MostChangedStatus string `json:"most_changed_status"`
}
