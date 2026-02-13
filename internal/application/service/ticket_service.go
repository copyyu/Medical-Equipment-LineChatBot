package service

import (
	"fmt"
	"log"
	"medical-webhook/internal/domain/line/entity"
	"medical-webhook/internal/domain/line/repository"
	"medical-webhook/internal/infrastructure/line/templates"
)

type TicketNotificationService struct {
	lineRepo   repository.LineRepository
	ticketRepo repository.TicketRepository
}

func NewTicketNotificationService(
	lineRepo repository.LineRepository,
	ticketRepo repository.TicketRepository,
) *TicketNotificationService {
	return &TicketNotificationService{
		lineRepo:   lineRepo,
		ticketRepo: ticketRepo,
	}
}

// NotifyStatusChange sends a Flex Message notification to the reporter when ticket status changes
func (s *TicketNotificationService) NotifyStatusChange(ticketID uint, oldStatus, newStatus entity.TicketStatus, note string) error {
	// Get ticket with full details
	ticket, err := s.ticketRepo.FindTicketByID(ticketID)
	if err != nil || ticket == nil {
		return fmt.Errorf("failed to find ticket: %w", err)
	}

	// Check if ticket has LINE user ID
	if ticket.ReporterLineID == nil || *ticket.ReporterLineID == "" {
		log.Printf("⚠️ Ticket %s has no LINE user ID, skip notification", ticket.TicketNo)
		return nil
	}

	// Send Flex Message with status transition
	flexMsg := templates.GetTicketStatusChangedFlex(ticket, oldStatus, newStatus, note)
	err = s.lineRepo.PushFlexMessage(*ticket.ReporterLineID, "🔔 อัปเดตสถานะ Ticket", flexMsg)
	if err != nil {
		log.Printf("❌ Failed to send status change flex notification for ticket %s: %v", ticket.TicketNo, err)
		return err
	}

	log.Printf("✅ Sent status change notification to user %s for ticket %s (%s → %s)",
		*ticket.ReporterLineID, ticket.TicketNo, oldStatus.GetStatusText(), newStatus.GetStatusText())
	return nil
}
