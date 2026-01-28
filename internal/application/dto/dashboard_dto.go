package dto

// DashboardSummaryResponse represents the main dashboard statistics
type DashboardSummaryResponse struct {
	TotalEquipment    int64               `json:"total_equipment"`
	RentalEquipment   int64               `json:"rental_equipment"`
	NearExpiry        int64               `json:"near_expiry"`
	TotalMaintenance  int64               `json:"total_maintenance"`
	AssetStatusCounts []AssetStatusCount  `json:"asset_status_counts"`
	JobStatusCounts   []JobStatusCount    `json:"job_status_counts"`
	RecentJobs        []RecentJobResponse `json:"recent_jobs"`
}

// AssetStatusCount represents equipment count by status
type AssetStatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

// JobStatusCount represents maintenance job count by status
type JobStatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

// RecentJobResponse represents a recent maintenance job
type RecentJobResponse struct {
	ID            string `json:"id"`
	EquipmentName string `json:"equipment_name"`
	Status        string `json:"status"`
	Assignee      string `json:"assignee"`
	UpdatedAt     string `json:"updated_at"`
}
