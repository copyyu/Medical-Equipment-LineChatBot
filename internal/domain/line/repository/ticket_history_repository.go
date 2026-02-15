package repository

import (
	"medical-webhook/internal/domain/line/entity"
	"time"
)

// StatusChangeLogEntry is a flattened view of a ticket_history row joined with ticket + admin
type StatusChangeLogEntry struct {
	ID            uint
	TicketID      uint
	TicketNo      string
	EquipmentName string
	AdminName     string
	FromStatus    string
	ToStatus      string
	Note          string
	ChangedAt     time.Time
}

// StatusChangeLogStats holds aggregate stats for activity logs
type StatusChangeLogStats struct {
	TotalChanges      int64
	TodayChanges      int64
	WeekChanges       int64
	MostChangedStatus string
}

type TicketHistoryRepository interface {
	CreateTicketHistory(history *entity.TicketHistory) error

	// GetStatusChangeLogs returns paginated & filtered status change entries
	GetStatusChangeLogs(page, limit int, search, fromStatus, toStatus, startDate, endDate string) ([]StatusChangeLogEntry, int64, error)

	// GetStatusChangeLogStats returns aggregate stats for activity logs
	GetStatusChangeLogStats() (*StatusChangeLogStats, error)
}
