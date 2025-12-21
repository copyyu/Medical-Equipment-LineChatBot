package dto

import "time"

type EquipmentReplacementAlertDTO struct {
	EquipmentID     uint    `json:"equipment_id"`
	IDCode          string  `json:"id_code"`
	SerialNo        string  `json:"serial_no"`
	BrandName       string  `json:"brand_name"`
	ModelName       string  `json:"model_name"`
	DepartmentName  string  `json:"department_name"`
	ReplacementYear int     `json:"replacement_year"`
	MonthsRemaining int     `json:"months_remaining"`
	PurchasePrice   float64 `json:"purchase_price"`
	NotifyRound     string  `json:"notify_round"`
}

type NotificationSummaryDTO struct {
	TotalEquipments  int        `json:"total_equipments"`
	JuneAlerts       int        `json:"june_alerts"`
	AugustAlerts     int        `json:"august_alerts"`
	LastNotification *time.Time `json:"last_notification"`
}

type NotificationSettingDTO struct {
	IsEnabled bool `json:"is_enabled"`
}
