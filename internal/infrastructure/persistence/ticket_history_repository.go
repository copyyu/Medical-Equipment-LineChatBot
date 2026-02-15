package persistence

import (
	"medical-webhook/internal/domain/line/entity"
	"medical-webhook/internal/domain/line/repository"
	"time"

	"gorm.io/gorm"
)

// TicketHistoryRepository implements repository.TicketHistoryRepository
type TicketHistoryRepository struct {
	db *gorm.DB
}

// NewTicketHistoryRepository creates a new ticket history repository
func NewTicketHistoryRepository(db *gorm.DB) *TicketHistoryRepository {
	return &TicketHistoryRepository{db: db}
}

func (r *TicketHistoryRepository) CreateTicketHistory(history *entity.TicketHistory) error {
	return r.db.Create(history).Error
}

// GetStatusChangeLogs returns paginated & filtered status change log entries
func (r *TicketHistoryRepository) GetStatusChangeLogs(
	page, limit int,
	search, fromStatus, toStatus, startDate, endDate string,
) ([]repository.StatusChangeLogEntry, int64, error) {

	// Base query: only status_changed actions
	query := r.db.Model(&entity.TicketHistory{}).
		Select(`
			ticket_histories.id,
			ticket_histories.ticket_id,
			tickets.ticket_no,
			COALESCE(tickets.equipment_name, '') as equipment_name,
			COALESCE(ticket_histories.changed_by, 'System') as admin_name,
			COALESCE(ticket_histories.old_value, '') as from_status,
			COALESCE(ticket_histories.new_value, '') as to_status,
			COALESCE(ticket_histories.note, '') as note,
			ticket_histories.created_at as changed_at
		`).
		Joins("JOIN tickets ON tickets.id = ticket_histories.ticket_id AND tickets.deleted_at IS NULL").
		Where("ticket_histories.action = ?", string(entity.ActionStatusChanged)).
		Where("ticket_histories.deleted_at IS NULL")

	// Filter: search (ticket_no, equipment_name, admin full_name)
	if search != "" {
		pattern := "%" + search + "%"
		query = query.Where(
			"(tickets.ticket_no ILIKE ? OR tickets.equipment_name ILIKE ? OR ticket_histories.changed_by ILIKE ?)",
			pattern, pattern, pattern,
		)
	}

	// Filter: from status
	if fromStatus != "" {
		query = query.Where("ticket_histories.old_value = ?", fromStatus)
	}

	// Filter: to status
	if toStatus != "" {
		query = query.Where("ticket_histories.new_value = ?", toStatus)
	}

	// Filter: date range
	if startDate != "" {
		t, err := time.Parse("2006-01-02", startDate)
		if err == nil {
			query = query.Where("ticket_histories.created_at >= ?", t)
		}
	}
	if endDate != "" {
		t, err := time.Parse("2006-01-02", endDate)
		if err == nil {
			// end of the day
			end := t.Add(24*time.Hour - time.Nanosecond)
			query = query.Where("ticket_histories.created_at <= ?", end)
		}
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Paginate + sort
	offset := (page - 1) * limit
	var entries []repository.StatusChangeLogEntry
	err := query.
		Order("ticket_histories.created_at DESC").
		Offset(offset).
		Limit(limit).
		Scan(&entries).Error

	if err != nil {
		return nil, 0, err
	}

	return entries, total, nil
}

// GetStatusChangeLogStats returns aggregate stats of status change logs
func (r *TicketHistoryRepository) GetStatusChangeLogStats() (*repository.StatusChangeLogStats, error) {
	stats := &repository.StatusChangeLogStats{}

	baseWhere := "action = ? AND deleted_at IS NULL"
	action := string(entity.ActionStatusChanged)

	// Total
	if err := r.db.Model(&entity.TicketHistory{}).
		Where(baseWhere, action).
		Count(&stats.TotalChanges).Error; err != nil {
		return nil, err
	}

	// Today
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if err := r.db.Model(&entity.TicketHistory{}).
		Where(baseWhere+" AND created_at >= ?", action, todayStart).
		Count(&stats.TodayChanges).Error; err != nil {
		return nil, err
	}

	// Last 7 days
	weekStart := todayStart.AddDate(0, 0, -7)
	if err := r.db.Model(&entity.TicketHistory{}).
		Where(baseWhere+" AND created_at >= ?", action, weekStart).
		Count(&stats.WeekChanges).Error; err != nil {
		return nil, err
	}

	// Most changed "to" status
	type StatusCount struct {
		NewValue string
		Count    int64
	}
	var result StatusCount
	err := r.db.Model(&entity.TicketHistory{}).
		Select("new_value, count(*) as count").
		Where(baseWhere, action).
		Group("new_value").
		Order("count DESC").
		Limit(1).
		Scan(&result).Error
	if err == nil && result.NewValue != "" {
		stats.MostChangedStatus = result.NewValue
	}

	return stats, nil
}
