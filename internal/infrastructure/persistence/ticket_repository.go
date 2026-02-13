package persistence

import (
	"fmt"
	"medical-webhook/internal/domain/line/entity"
	"time"

	"gorm.io/gorm"
)

type TicketRepository struct {
	db *gorm.DB
}

// NewTicketRepository creates a new ticket repository
func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) CreateTicket(ticket *entity.Ticket) error {
	return r.db.Create(ticket).Error
}

func (r *TicketRepository) FindTicketByNo(ticketNo string) (*entity.Ticket, error) {
	var ticket entity.Ticket
	err := r.db.Preload("Category").
		Preload("Equipment.Model.Brand").
		Preload("Equipment.Department").
		Preload("Department").
		Preload("Reporter").
		Preload("Histories").
		Where("ticket_no = ?", ticketNo).
		First(&ticket).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &ticket, err
}

func (r *TicketRepository) FindTicketByID(id uint) (*entity.Ticket, error) {
	var ticket entity.Ticket
	err := r.db.Preload("Category").
		Preload("Equipment").
		Preload("Department").
		Preload("Reporter").
		Preload("Histories.Admin").
		First(&ticket, id).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &ticket, err
}

func (r *TicketRepository) UpdateTicket(ticket *entity.Ticket) error {
	return r.db.Save(ticket).Error
}

func (r *TicketRepository) UpdateTicketStatus(ticketID uint, newStatus entity.TicketStatus, note string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get current ticket
		var ticket entity.Ticket
		if err := tx.First(&ticket, ticketID).Error; err != nil {
			return err
		}

		oldStatus := ticket.Status

		ticket.Status = newStatus

		// Update timestamps
		now := time.Now()
		if newStatus == entity.TicketStatusInProgress && ticket.StartedAt == nil {
			ticket.StartedAt = &now
		} else if newStatus == entity.TicketStatusCompleted && ticket.CompletedAt == nil {
			ticket.CompletedAt = &now
		}

		// Save ticket
		if err := tx.Save(&ticket).Error; err != nil {
			return err
		}

		// Create history for Status Change
		if oldStatus != newStatus {
			history := &entity.TicketHistory{
				TicketID: ticketID,
				Action:   entity.ActionStatusChanged,
				Field:    stringPtr("status"),
				OldValue: stringPtr(string(oldStatus)),
				NewValue: stringPtr(string(newStatus)),
				Note:     stringPtr(note),
				IsSystem: true,
			}
			if err := tx.Create(history).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *TicketRepository) GetTicketsByLineUserID(lineUserID string) ([]entity.Ticket, error) {
	var tickets []entity.Ticket
	err := r.db.Preload("Category").
		Preload("Equipment").
		Where("reporter_line_id = ?", lineUserID).
		Order("created_at DESC").
		Find(&tickets).Error
	return tickets, err
}

// Helper function
func stringPtr(s string) *string {
	return &s
}

func (r *TicketRepository) GetLatestTicketNumber(year int) (string, error) {
	var ticketNo string
	prefix := fmt.Sprintf("REQ-%d-%%", year)

	err := r.db.Model(&entity.Ticket{}).
		Select("ticket_no").
		Where("ticket_no LIKE ?", prefix).
		Order("ticket_no DESC").
		Limit(1).
		Pluck("ticket_no", &ticketNo).Error

	if err == gorm.ErrRecordNotFound {
		return "", nil
	}
	return ticketNo, err
}

func (r *TicketRepository) GetAllTickets(page, limit int, status, priority, search, sortBy, sortDir string) ([]entity.Ticket, int64, error) {
	var tickets []entity.Ticket
	var total int64

	offset := (page - 1) * limit
	query := r.db.Model(&entity.Ticket{}).
		Preload("Category").
		Preload("Equipment.Model").
		Preload("Department").
		Preload("Reporter")

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("ticket_no LIKE ? OR description LIKE ?", searchPattern, searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if sortBy != "" {
		order := sortBy
		if sortDir != "" {
			order += " " + sortDir
		}
		query = query.Order(order)
	} else {
		query = query.Order("created_at DESC")
	}

	err := query.Offset(offset).Limit(limit).Find(&tickets).Error
	return tickets, total, err
}

func (r *TicketRepository) GetTicketStats() (total, inProgress, completed, sendToOutsource int64, err error) {
	type Result struct {
		Status entity.TicketStatus
		Count  int64
	}
	var results []Result

	// Count per status
	err = r.db.Model(&entity.Ticket{}).
		Select("status, count(*) as count").
		Group("status").
		Scan(&results).Error

	if err != nil {
		return 0, 0, 0, 0, err
	}

	for _, res := range results {
		total += res.Count
		switch res.Status {
		case entity.TicketStatusInProgress:
			inProgress = res.Count
		case entity.TicketStatusCompleted:
			completed = res.Count
		case entity.TicketStatusSendToOutsource:
			sendToOutsource = res.Count
		}
	}

	return total, inProgress, completed, sendToOutsource, nil
}

// GetRecentTickets returns limited number of recent tickets
func (r *TicketRepository) GetRecentTickets(limit int) ([]entity.Ticket, error) {
	var tickets []entity.Ticket
	err := r.db.Preload("Category").
		Preload("Equipment.Model.Brand").
		Preload("Equipment.Department").
		Preload("Department").
		Preload("Reporter").
		Order("created_at DESC").
		Limit(limit).
		Find(&tickets).Error
	return tickets, err
}

// FindPendingTicketByEquipmentAndUser finds existing pending/in_progress ticket for equipment by LINE user
func (r *TicketRepository) FindPendingTicketByEquipmentAndUser(equipmentID uint, lineUserID string) (*entity.Ticket, error) {
	var ticket entity.Ticket
	err := r.db.Preload("Category").
		Where("equipment_id = ?", equipmentID).
		Where("reporter_line_id = ?", lineUserID).
		Where("status IN ?", []entity.TicketStatus{entity.TicketStatusInProgress}).
		First(&ticket).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &ticket, err
}
