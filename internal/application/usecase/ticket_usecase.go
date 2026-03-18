package usecase

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"medical-webhook/internal/application/dto"
	"medical-webhook/internal/application/service"
	"medical-webhook/internal/domain/event"
	"medical-webhook/internal/domain/line/entity"
	"medical-webhook/internal/domain/line/repository"
	"medical-webhook/internal/infrastructure/line/templates"
	"medical-webhook/internal/utils/ptr"
)

// TicketUseCase handles ticket-related business logic
type TicketUseCase struct {
	lineRepo      repository.LineRepository
	equipmentRepo repository.EquipmentRepository
	ticketRepo    repository.TicketRepository
	categoryRepo  repository.TicketCategoryRepository
	historyRepo   repository.TicketHistoryRepository
	notifyService *service.TicketNotificationService
	eventBus      event.EventBus
}

// NewTicketUseCase creates a new ticket use case
func NewTicketUseCase(
	lineRepo repository.LineRepository,
	equipmentRepo repository.EquipmentRepository,
	ticketRepo repository.TicketRepository,
	categoryRepo repository.TicketCategoryRepository,
	historyRepo repository.TicketHistoryRepository,
	notifyService *service.TicketNotificationService,
	eventBus event.EventBus,
) *TicketUseCase {
	return &TicketUseCase{
		lineRepo:      lineRepo,
		equipmentRepo: equipmentRepo,
		ticketRepo:    ticketRepo,
		categoryRepo:  categoryRepo,
		historyRepo:   historyRepo,
		notifyService: notifyService,
		eventBus:      eventBus,
	}
}

// ErrDuplicateTicket is returned when user already has a pending ticket for this equipment
var ErrDuplicateTicket = fmt.Errorf("duplicate ticket exists")

// GetTicketList returns paginated ticket list
func (uc *TicketUseCase) GetTicketList(ctx context.Context, req dto.TicketListRequest) (*dto.TicketListResponse, error) {
	tickets, total, err := uc.ticketRepo.GetAllTickets(
		req.Page,
		req.Limit,
		req.Status,
		req.Priority,
		req.Search,
		req.SortBy,
		req.SortDir,
	)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	var ticketDTOs []dto.TicketItemResponse
	for _, t := range tickets {
		item := dto.TicketItemResponse{
			ID:       t.ID,
			TicketNo: t.TicketNo,

			ReporterName:     t.ReporterName,
			ReporterPhotoURL: t.ReporterPhotoURL,
			ReportedAt:       t.ReportedAt,
			CompletedAt:      t.CompletedAt,
			CreatedAt:        t.CreatedAt,
			Description:      t.Description,
			Priority:         string(t.Priority),
			PriorityText:     t.Priority.GetPriorityText(),
			Status:           string(t.Status),
			StatusText:       t.Status.GetStatusText(),
		}

		if t.Category.ID != 0 {
			item.CategoryID = t.Category.ID
			item.CategoryName = t.Category.Name
		}

		if t.Equipment != nil {
			item.EquipmentName = &t.Equipment.Model.ModelName
			item.EquipmentIDCode = &t.Equipment.IDCode
		}

		if t.Department != nil {
			item.DepartmentName = &t.Department.Name
		}

		ticketDTOs = append(ticketDTOs, item)
	}

	return &dto.TicketListResponse{
		Data: ticketDTOs,
		Pagination: dto.Pagination{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetTicketByID returns ticket detail
func (uc *TicketUseCase) GetTicketByID(ctx context.Context, id uint) (*dto.TicketDetailResponse, error) {
	ticket, err := uc.ticketRepo.FindTicketByID(id)
	if err != nil {
		return nil, err
	}
	if ticket == nil {
		return nil, nil
	}

	resp := &dto.TicketDetailResponse{
		ID:       ticket.ID,
		TicketNo: ticket.TicketNo,

		ReporterName:     ticket.ReporterName,
		ReporterLineID:   ticket.ReporterLineID,
		ReporterPhotoURL: ticket.ReporterPhotoURL,
		DepartmentID:     ticket.DepartmentID,
		ContactInfo:      ticket.ContactInfo,
		ReportedAt:       ticket.ReportedAt,
		StartedAt:        ticket.StartedAt,
		CompletedAt:      ticket.CompletedAt,
		CreatedAt:        ticket.CreatedAt,
		UpdatedAt:        ticket.UpdatedAt,
		Priority:         string(ticket.Priority),
		PriorityText:     ticket.Priority.GetPriorityText(),
		Status:           string(ticket.Status),
		StatusText:       ticket.Status.GetStatusText(),
		Description:      ticket.Description,
		CategoryID:       ticket.CategoryID,
	}

	if ticket.CompletedAt != nil {
		hours := ticket.GetDurationHours()
		resp.DurationHours = &hours
	}

	if ticket.Category.ID != 0 {
		resp.CategoryName = ticket.Category.Name
	}

	if ticket.Equipment != nil {
		resp.EquipmentID = &ticket.Equipment.ID
		resp.EquipmentName = &ticket.Equipment.Model.ModelName
		resp.EquipmentIDCode = &ticket.Equipment.IDCode
		// resp.Location = &ticket.Equipment.Location // Location field not found in Equipment
	}

	if ticket.Department != nil {
		resp.DepartmentName = &ticket.Department.Name
	}

	for _, h := range ticket.Histories {
		historyDTO := dto.TicketHistoryDTO{
			ID:        h.ID,
			Action:    string(h.Action),
			Field:     h.Field,
			OldValue:  h.OldValue,
			NewValue:  h.NewValue,
			Note:      h.Note,
			IsSystem:  h.IsSystem,
			CreatedAt: h.CreatedAt,
		}
		if h.Admin != nil {
			historyDTO.AdminName = &h.Admin.FullName
		}
		resp.Histories = append(resp.Histories, historyDTO)
	}

	return resp, nil
}

// UpdateTicket updates ticket details including status, priority, and description
func (uc *TicketUseCase) UpdateTicket(ctx context.Context, id uint, req dto.UpdateTicketRequest) error {
	ticket, err := uc.ticketRepo.FindTicketByID(id)
	if err != nil {
		return err
	}
	if ticket == nil {
		return fmt.Errorf("ticket not found")
	}

	var changes []entity.TicketHistory

	// Update Priority
	if req.Priority != nil && string(*req.Priority) != string(ticket.Priority) {
		oldVal := string(ticket.Priority)
		ticket.Priority = entity.TicketPriority(*req.Priority)
		changes = append(changes, entity.TicketHistory{
			TicketID: ticket.ID,
			Action:   entity.ActionUpdated,
			Field:    ptr.StringPtr("priority"),
			OldValue: ptr.StringPtr(oldVal),
			NewValue: req.Priority,
			Note:     ptr.StringPtr(req.Note),
			IsSystem: false,
		})
	}

	// Update Description
	if req.Description != nil && (ticket.Description == nil || *req.Description != *ticket.Description) {
		oldVal := ""
		if ticket.Description != nil {
			oldVal = *ticket.Description
		}
		ticket.Description = req.Description
		changes = append(changes, entity.TicketHistory{
			TicketID: ticket.ID,
			Action:   entity.ActionUpdated,
			Field:    ptr.StringPtr("description"),
			OldValue: ptr.StringPtr(oldVal),
			NewValue: req.Description,
			Note:     ptr.StringPtr(req.Note),
			IsSystem: false,
		})
	}

	// Update Status
	if req.Status != nil && string(*req.Status) != string(ticket.Status) {
		oldVal := string(ticket.Status)
		newStatus := entity.TicketStatus(*req.Status)
		ticket.Status = newStatus

		// Update timestamps
		now := time.Now()
		if newStatus == entity.TicketStatusInProcess && ticket.StartedAt == nil {
			ticket.StartedAt = &now
		} else if newStatus == entity.TicketStatusCompleted && ticket.CompletedAt == nil {
			ticket.CompletedAt = &now
		}

		changes = append(changes, entity.TicketHistory{
			TicketID: ticket.ID,
			Action:   entity.ActionStatusChanged,
			Field:    ptr.StringPtr("status"),
			OldValue: ptr.StringPtr(oldVal),
			NewValue: req.Status,
			Note:     ptr.StringPtr(req.Note),
			IsSystem: false, // or true if we consider this system action? But it is triggered by user/admin
		})
	}

	// Track if status changed for notification
	statusChanged := false
	var oldStatusValue entity.TicketStatus
	var newStatusValue entity.TicketStatus
	var statusNote string
	for _, change := range changes {
		if change.Action == entity.ActionStatusChanged {
			statusChanged = true
			oldStatusValue = entity.TicketStatus(*change.OldValue)
			newStatusValue = entity.TicketStatus(*change.NewValue)
			if change.Note != nil {
				statusNote = *change.Note
			}
			break
		}
	}

	// Update ticket in DB
	if err := uc.ticketRepo.UpdateTicket(ticket); err != nil {
		return err
	}

	// Save history records
	for i := range changes {
		changes[i].ChangedBy = req.ChangedBy
		_ = uc.historyRepo.CreateTicketHistory(&changes[i])
	}

	// Send LINE notification if status changed (fire-and-forget)
	if statusChanged && uc.notifyService != nil {
		go func() {
			if err := uc.notifyService.NotifyStatusChange(ticket.ID, oldStatusValue, newStatusValue, statusNote); err != nil {
				log.Printf("Failed to send status change notification for ticket %s: %v", ticket.TicketNo, err)
			}
		}()
	}

	// Publish ticket updated event
	if uc.eventBus != nil {
		go func() {
			publishErr := uc.eventBus.Publish(context.Background(), event.NewEvent(event.TicketUpdated, map[string]interface{}{
				"ticket_id": ticket.ID,
				"ticket_no": ticket.TicketNo,
				"status":    string(ticket.Status),
				"priority":  string(ticket.Priority),
			}))
			if publishErr != nil {
				log.Printf("Failed to publish ticket.updated event: %v", publishErr)
			}
		}()
	}

	return nil
}

// GetTicketStats returns ticket statistics
func (uc *TicketUseCase) GetTicketStats(ctx context.Context) (*dto.TicketStatsResponse, error) {
	total, inProgress, completed, sendToOutsource, err := uc.ticketRepo.GetTicketStats()
	if err != nil {
		return nil, err
	}

	return &dto.TicketStatsResponse{
		Total:           total,
		InProcess:       inProgress,
		Completed:       completed,
		SendToOutsource: sendToOutsource,
	}, nil
}

func (uc *TicketUseCase) GetTicketCategories(ctx context.Context) ([]entity.TicketCategory, error) {
	return uc.categoryRepo.GetTicketCategories()
}

// CreateTicketFromLINE creates a ticket from LINE report
func (uc *TicketUseCase) CreateTicketFromLINE(
	serialOrCode string,
	description string,
	lineUserID string,
	lineDisplayName string,
	linePhotoURL string,
	categoryID uint,
) (*entity.Ticket, error) {
	// Find equipment first
	equipment, err := uc.equipmentRepo.FindBySerialOrCode(serialOrCode)
	if err != nil || equipment == nil {
		return nil, fmt.Errorf("equipment not found: %s", serialOrCode)
	}

	// Check for existing pending/in_progress ticket for this equipment by this user
	existingTicket, err := uc.ticketRepo.FindPendingTicketByEquipmentAndUser(equipment.ID, lineUserID)
	if err != nil {
		log.Printf("❌ Failed to check existing ticket: %v", err)
	}
	if existingTicket != nil {
		log.Printf("⚠️ User %s already has pending ticket %s for equipment %d", lineUserID, existingTicket.TicketNo, equipment.ID)
		return existingTicket, ErrDuplicateTicket
	}

	// Generate ticket number from DB to avoid duplicates
	ticketNo, err := uc.generateTicketNumberFromDB()
	if err != nil {
		return nil, fmt.Errorf("failed to generate ticket number: %w", err)
	}

	equipmentName := "อุปกรณ์"
	if equipment.Model.ModelName != "" {
		equipmentName = equipment.Model.ModelName
	}

	// Prepare title

	// Use provided categoryID or find/create default
	if categoryID == 0 {
		category, err := uc.categoryRepo.FindCategoryByName("แจ้งซ่อม")
		if err == nil && category != nil {
			categoryID = category.ID
		} else {
			// Create default category if not found
			newCategory := &entity.TicketCategory{
				Name:      "แจ้งซ่อม",
				NameEn:    "Repair",
				Color:     "#EF5350",
				Icon:      "🔧",
				IsActive:  true,
				SortOrder: 1,
			}
			if err := uc.categoryRepo.CreateCategory(newCategory); err == nil {
				categoryID = newCategory.ID
				log.Printf("Created default category 'แจ้งซ่อม' (ID: %d)", categoryID)
			} else {
				log.Printf("Failed to create default category: %v, using fallback ID 1", err)
				categoryID = 1
			}
		}
	}

	// Create ticket (LINE users don't have Admin accounts, use ReporterLineID and ReporterName)
	ticket := &entity.Ticket{
		TicketNo:         ticketNo,
		Description:      &description,
		CategoryID:       categoryID,
		Priority:         entity.PriorityMedium,
		EquipmentID:      &equipment.ID,
		EquipmentName:    &equipmentName,
		ReporterName:     lineDisplayName,
		ReporterLineID:   &lineUserID,
		ReporterPhotoURL: &linePhotoURL,
		DepartmentID:     &equipment.DepartmentID,
		Status:           entity.TicketStatusInProcess,
		ReportedAt:       time.Now(),
	}

	err = uc.ticketRepo.CreateTicket(ticket)
	if err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	// Create initial history (no AdminID for LINE users)
	history := &entity.TicketHistory{
		TicketID: ticket.ID,
		Action:   entity.ActionCreated,
		Note:     ptr.StringPtr("สร้างจาก LINE โดย " + lineDisplayName),
		IsSystem: true,
	}
	_ = uc.historyRepo.CreateTicketHistory(history)

	// Publish ticket created event
	if uc.eventBus != nil {
		go func() {
			publishErr := uc.eventBus.Publish(context.Background(), event.NewEvent(event.TicketCreated, map[string]interface{}{
				"ticket_id": ticket.ID,
				"ticket_no": ticketNo,
				"equipment": serialOrCode,
				"reporter":  lineDisplayName,
				"status":    string(ticket.Status),
			}))
			if publishErr != nil {
				log.Printf("Failed to publish ticket.created event: %v", publishErr)
			}
		}()
	}

	log.Printf("Created ticket %s for equipment %s", ticketNo, serialOrCode)
	return ticket, nil
}

// GetTicketByNo finds ticket by ticket number
func (uc *TicketUseCase) GetTicketByNo(ticketNo string, lineUserID string) (*entity.Ticket, error) {
	ticket, err := uc.ticketRepo.FindTicketByNo(ticketNo)
	if err != nil {
		return nil, fmt.Errorf("failed to find ticket: %w", err)
	}

	if ticket == nil {
		return nil, nil
	}

	// Check if ticket belongs to this user
	if ticket.ReporterLineID == nil || *ticket.ReporterLineID != lineUserID {
		return nil, fmt.Errorf("unauthorized access")
	}

	return ticket, nil
}

// GetUserTickets gets all tickets for a LINE user
func (uc *TicketUseCase) GetUserTickets(lineUserID string) ([]entity.Ticket, error) {
	return uc.ticketRepo.GetTicketsByLineUserID(lineUserID)
}

// SendTicketCreatedMessage sends ticket created flex message
func (uc *TicketUseCase) SendTicketCreatedMessage(replyToken string, ticket *entity.Ticket) error {
	flexMsg := templates.GetTicketCreatedFlex(ticket)
	return uc.lineRepo.ReplyFlexMessage(replyToken, "สร้าง Ticket สำเร็จ", flexMsg)
}

// SendTicketStatusMessage sends ticket status flex message
func (uc *TicketUseCase) SendTicketStatusMessage(replyToken string, ticket *entity.Ticket) error {
	flexMsg := templates.GetTicketStatusFlex(ticket)
	return uc.lineRepo.ReplyFlexMessage(replyToken, "สถานะ Ticket", flexMsg)
}

// SendMyTicketsMessage sends user's tickets as carousel
func (uc *TicketUseCase) SendMyTicketsMessage(replyToken string, tickets []entity.Ticket) error {
	if len(tickets) == 0 {
		return uc.lineRepo.ReplyMessage(replyToken, "ℹ️ คุณยังไม่มีรายการแจ้งปัญหา")
	}

	flexMsg := templates.GetMyTicketsFlex(tickets)
	return uc.lineRepo.ReplyFlexMessage(replyToken, "รายการ Ticket ของคุณ", flexMsg)
}

// generateTicketNumberFromDB generates a unique ticket number using DB sequence
func (uc *TicketUseCase) generateTicketNumberFromDB() (string, error) {
	year := time.Now().Year()
	latestNo, err := uc.ticketRepo.GetLatestTicketNumber(year)
	if err != nil {
		return "", err
	}

	var nextNum int
	if latestNo == "" {
		nextNum = 1
	} else {
		// Parse number from format REQ-YYYY-XXXXX
		var parsedNum int
		_, err := fmt.Sscanf(latestNo, "REQ-%d-%d", new(int), &parsedNum)
		if err != nil {
			nextNum = 1
		} else {
			nextNum = parsedNum + 1
		}
	}

	return fmt.Sprintf("REQ-%d-%05d", year, nextNum), nil
}

// GetTicketsByEquipmentID finds all tickets for an equipment
func (uc *TicketUseCase) GetTicketsByEquipmentID(equipmentID uint) ([]entity.Ticket, error) {
	return uc.ticketRepo.GetTicketsByEquipmentID(equipmentID)
}
