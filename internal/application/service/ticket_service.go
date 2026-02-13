package service

import (
	"fmt"
	"log"
	"medical-webhook/internal/domain/line/entity"
	"medical-webhook/internal/domain/line/model"
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

// NotifyStatusChange sends notification when ticket status changes
func (s *TicketNotificationService) NotifyStatusChange(ticketID uint, newStatus entity.TicketStatus, note string) error {
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

	statusText := newStatus.GetStatusText()
	statusEmoji := getStatusEmoji(newStatus)

	// Determine equipment name
	equipmentName := "-"
	if ticket.EquipmentName != nil {
		equipmentName = *ticket.EquipmentName
	}

	// Build message text
	messageText := fmt.Sprintf(
		"%s สถานะ Ticket อัปเดต\n\n"+
			"หมายเลข: %s\n"+
			"สถานะใหม่: %s\n"+
			"อุปกรณ์: %s",
		statusEmoji,
		ticket.TicketNo,
		statusText,
		equipmentName,
	)

	if note != "" {
		messageText += fmt.Sprintf("\n\n📝 หมายเหตุ: %s", note)
	}

	messageText += fmt.Sprintf(
		"\n\n💬 ดูรายละเอียดเพิ่มเติม กด 'ติดตามสถานะ' แล้วพิมพ์: %s",
		ticket.TicketNo,
	)

	err = s.lineRepo.PushMessage(&model.OutgoingMessage{
		To:   *ticket.ReporterLineID,
		Text: messageText,
	})
	if err != nil {
		log.Printf("Failed to send status notification: %v", err)
		return err
	}

	flexMsg := templates.GetTicketStatusFlex(ticket)
	err = s.lineRepo.PushFlexMessage(*ticket.ReporterLineID, "สถานะ Ticket", flexMsg)
	if err != nil {
		log.Printf("Failed to send flex notification: %v", err)
	}

	log.Printf("Sent status notification to user %s for ticket %s", *ticket.ReporterLineID, ticket.TicketNo)
	return nil
}

// NotifyTicketCompleted sends notification when ticket is completed
func (s *TicketNotificationService) NotifyTicketCompleted(ticketID uint) error {
	ticket, err := s.ticketRepo.FindTicketByID(ticketID)
	if err != nil || ticket == nil {
		return fmt.Errorf("failed to find ticket: %w", err)
	}

	if ticket.ReporterLineID == nil {
		return nil
	}

	// Calculate duration
	durationText := "N/A"
	if ticket.CompletedAt != nil {
		hours := ticket.GetDurationHours()
		if hours < 24 {
			durationText = fmt.Sprintf("%.1f ชั่วโมง", hours)
		} else {
			days := hours / 24
			durationText = fmt.Sprintf("%.1f วัน", days)
		}
	}

	equipmentName := "-"
	if ticket.EquipmentName != nil {
		equipmentName = *ticket.EquipmentName
	}

	messageText := fmt.Sprintf(
		"Ticket เสร็จสิ้นแล้ว!\n\n"+
			"หมายเลข: %s\n"+
			"อุปกรณ์: %s\n"+
			"ระยะเวลาดำเนินการ: %s\n\n"+
			"ขอบคุณที่ใช้บริการค่ะ 🙏",
		ticket.TicketNo,
		equipmentName,
		durationText,
	)

	return s.lineRepo.PushMessage(&model.OutgoingMessage{
		To:   *ticket.ReporterLineID,
		Text: messageText,
	})
}

// Helper function to get emoji for status
func getStatusEmoji(status entity.TicketStatus) string {
	switch status {
	case entity.TicketStatusInProgress:
		return "🔧"
	case entity.TicketStatusCompleted:
		return "✅"
	case entity.TicketStatusSendToOutsource:
		return "📤"
	default:
		return "📋"
	}
}
