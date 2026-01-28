package dto

import "medical-webhook/internal/domain/line/entity"

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

// MapEquipmentToListItem converts entity.Equipment to EquipmentListItem
func MapEquipmentToListItem(e entity.Equipment) EquipmentListItem {
	// Get name from model
	name := ""
	if e.Model.ModelName != "" {
		name = e.Model.ModelName
	} else {
		name = e.IDCode
	}

	// Get category from model
	category := ""
	if e.Model.Category.Name != "" {
		category = e.Model.Category.Name
	}

	// Get location from department
	location := ""
	if e.Department.Name != "" {
		location = e.Department.Name
	}

	// Calculate expiry date based on replacement year
	expiry := ""
	isExpiring := false
	if e.ReplacementYear != nil {
		expiry = formatYear(*e.ReplacementYear)
		// Check if expiring within 1 year
		if e.RemainLife <= 1 {
			isExpiring = true
		}
	} else if e.RemainLife > 0 {
		// Calculate based on remain life
		currentYear := 2026 // Current year
		expiryYear := currentYear + int(e.RemainLife)
		expiry = formatYear(expiryYear)
		if e.RemainLife <= 1 {
			isExpiring = true
		}
	}

	// Get last check date from latest maintenance record
	lastCheck := ""
	if len(e.MaintenanceRecords) > 0 {
		lastCheck = e.MaintenanceRecords[0].MaintenanceDate.Format("2006-01-02")
	} else if e.ComputeDate != nil {
		lastCheck = e.ComputeDate.Format("2006-01-02")
	}

	// Map asset status to frontend status
	status := mapAssetStatusToFrontend(e.Status)

	return EquipmentListItem{
		ID:         e.IDCode,
		Name:       name,
		Category:   category,
		Status:     status,
		Location:   location,
		LastCheck:  lastCheck,
		Expiry:     expiry,
		IsExpiring: isExpiring,
	}
}

// mapAssetStatusToFrontend maps backend AssetStatus to frontend EquipmentStatus
// Return raw status values to match frontend expectations
func mapAssetStatusToFrontend(status entity.AssetStatus) string {
	// Return the status value directly (active, defective, etc.)
	if status == "" {
		return "active"
	}
	return string(status)
}

func formatYear(year int) string {
	return string(rune('0'+year/1000)) + string(rune('0'+(year/100)%10)) + string(rune('0'+(year/10)%10)) + string(rune('0'+year%10))
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

// MapEquipmentToDetailResponse converts entity.Equipment to EquipmentDetailResponse
func MapEquipmentToDetailResponse(e entity.Equipment) EquipmentDetailResponse {
	item := MapEquipmentToListItem(e)

	serialNo := ""
	if e.SerialNo != nil {
		serialNo = *e.SerialNo
	}

	brand := ""
	if e.Model.Brand.Name != "" {
		brand = e.Model.Brand.Name
	}

	return EquipmentDetailResponse{
		ID:           item.ID,
		Name:         item.Name,
		Category:     item.Category,
		Status:       item.Status,
		Location:     item.Location,
		LastCheck:    item.LastCheck,
		Expiry:       item.Expiry,
		IsExpiring:   item.IsExpiring,
		SerialNo:     serialNo,
		Brand:        brand,
		DepartmentID: e.DepartmentID,
	}
}
