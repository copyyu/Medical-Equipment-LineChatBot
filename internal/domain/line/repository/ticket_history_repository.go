package repository

import "medical-webhook/internal/domain/line/entity"

type TicketHistoryRepository interface {
	CreateTicketHistory(history *entity.TicketHistory) error
}
