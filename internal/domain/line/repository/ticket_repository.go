package repository

import "medical-webhook/internal/domain/line/entity"

type TicketRepository interface {
	CreateTicket(ticket *entity.Ticket) error
	FindTicketByNo(ticketNo string) (*entity.Ticket, error)
	FindTicketByID(id uint) (*entity.Ticket, error)
	UpdateTicket(ticket *entity.Ticket) error
	UpdateTicketStatus(ticketID uint, newStatus entity.TicketStatus, note string) error
	GetTicketsByLineUserID(lineUserID string) ([]entity.Ticket, error)
	GetLatestTicketNumber(year int) (string, error)
	GetAllTickets(page, limit int, status, priority, search, sortBy, sortDir string) ([]entity.Ticket, int64, error)
	GetTicketStats() (total, inProgress, completed, sendToOutsource int64, err error)
	GetRecentTickets(limit int) ([]entity.Ticket, error)
	FindPendingTicketByEquipmentAndUser(equipmentID uint, lineUserID string) (*entity.Ticket, error)
	GetTicketsByEquipmentID(equipmentID uint) ([]entity.Ticket, error)
}
