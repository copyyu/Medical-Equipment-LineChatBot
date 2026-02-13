package persistence

import (
	"medical-webhook/internal/domain/line/entity"

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
