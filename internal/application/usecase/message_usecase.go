package usecase

import (
	"context"
	"fmt"
	"log"
	"strings"

	"medical-webhook/internal/application/service"
	"medical-webhook/internal/domain/constants"
	"medical-webhook/internal/domain/line/entity"
	"medical-webhook/internal/domain/line/model"
	"medical-webhook/internal/domain/line/repository"
	"medical-webhook/internal/infrastructure/client"
	"medical-webhook/internal/infrastructure/line/templates"
	"medical-webhook/internal/infrastructure/session"
)

// MessageUseCase orchestrates message processing flow
type MessageUseCase struct {
	lineRepo       repository.LineRepository
	equipmentRepo  repository.EquipmentRepository
	departmentRepo repository.DepartmentRepository
	ocrClient      *client.OCRClient
	sessionStore   *session.SessionStore
	messageService *service.MessageService
	ticketUseCase  *TicketUseCase
	baseURL        string
}

// NewMessageUseCase creates a new message use case
func NewMessageUseCase(
	lineRepo repository.LineRepository,
	equipmentRepo repository.EquipmentRepository,
	departmentRepo repository.DepartmentRepository,
	ocrClient *client.OCRClient,
	sessionStore *session.SessionStore,
	messageService *service.MessageService,
	ticketUseCase *TicketUseCase,
	baseURL string,
) *MessageUseCase {
	return &MessageUseCase{
		lineRepo:       lineRepo,
		equipmentRepo:  equipmentRepo,
		departmentRepo: departmentRepo,
		ocrClient:      ocrClient,
		sessionStore:   sessionStore,
		messageService: messageService,
		ticketUseCase:  ticketUseCase,
		baseURL:        baseURL,
	}
}

// HandleTextMessage handles incoming text message from Rich Menu or direct input.
// It routes messages to appropriate handlers based on Rich Menu commands or session state.
func (uc *MessageUseCase) HandleTextMessage(msg *model.IncomingMessage) error {
	log.Printf("📝 Processing text: %s", msg.Text)
	text := strings.TrimSpace(msg.Text)

	// Try to handle as Rich Menu command first
	if handled, err := uc.handleRichMenuCommand(msg, text); handled {
		return err
	}

	// Otherwise, handle as user input based on session
	return uc.handleUserInput(msg, text)
}

// handleRichMenuCommand handles Rich Menu button commands.
// Returns (true, error) if command was handled, (false, nil) if not a Rich Menu command.
func (uc *MessageUseCase) handleRichMenuCommand(msg *model.IncomingMessage, text string) (bool, error) {
	switch {
	case strings.Contains(text, "แจ้งปัญหา") || strings.Contains(text, "เช็คสถานะ"):
		// แสดง sub-menu ให้เลือก: แจ้งปัญหา / ดูเครื่องใกล้หมดอายุ
		return true, uc.lineRepo.ReplyFlexMessage(msg.ReplyToken, "เลือกบริการ", templates.GetReportMenuFlex())

	case strings.Contains(text, "ติดตามสถานะ"):
		// Directly show user's tickets
		tickets, err := uc.ticketUseCase.GetUserTickets(msg.UserID)
		if err != nil {
			log.Printf("❌ GetUserTickets error: %v", err)
			return true, uc.lineRepo.ReplyMessage(msg.ReplyToken, "❌ ไม่สามารถดึงข้อมูลได้ กรุณาลองใหม่ค่ะ")
		}
		if len(tickets) == 0 {
			return true, uc.lineRepo.ReplyMessage(msg.ReplyToken, "📋 คุณยังไม่มีรายการแจ้งปัญหาค่ะ\n\nหากต้องการแจ้งปัญหา กรุณากดเมนู \"แจ้งปัญหา / เช็คสถานะ\" ค่ะ")
		}
		return true, uc.lineRepo.ReplyFlexMessage(msg.ReplyToken, "รายการแจ้งปัญหาของคุณ", templates.GetMyTicketsFlex(tickets))

	case strings.Contains(text, "เปลี่ยนเครื่อง"):
		return true, uc.lineRepo.ReplyFlexMessage(msg.ReplyToken, "แจ้งเปลี่ยนเครื่อง", uc.messageService.GetEquipmentChangeFlex("https://www.google.com/"))

	case strings.Contains(text, "ติดต่อ"):
		return true, uc.lineRepo.ReplyFlexMessage(msg.ReplyToken, "ติดต่อเจ้าหน้าที่", uc.messageService.GetContactStaffFlex())

	case text == "เมนู" || strings.ToLower(text) == "menu":
		return true, uc.lineRepo.ReplyMessage(msg.ReplyToken, constants.MsgSelectMenu)

	default:
		return false, nil
	}
}

// handleUserInput handles user text input based on current session mode.
func (uc *MessageUseCase) handleUserInput(msg *model.IncomingMessage, text string) error {
	sess := uc.sessionStore.Get(msg.UserID)
	if sess == nil || sess.Mode == session.ModeNone {
		return uc.lineRepo.ReplyMessage(msg.ReplyToken, constants.MsgSelectMenuFirst)
	}

	switch sess.Mode {
	case session.ModeReportProblem:
		return uc.handleReportProblemInput(msg, text)
	case session.ModeTrackStatus:
		return uc.handleTrackStatusInput(msg, text)
	case session.ModeInputIssueDesc:
		return uc.handleInputIssueDescInput(msg, text)
	case session.ModeSelectDeptForExpiry:
		return uc.handleSelectDeptForExpiryInput(msg, text)
	default:
		return uc.lineRepo.ReplyMessage(msg.ReplyToken, constants.MsgSelectMenuFirst)
	}
}

// handleReportProblemInput handles user input when in report problem mode.
func (uc *MessageUseCase) handleReportProblemInput(msg *model.IncomingMessage, text string) error {
	// Validate and sanitize input
	sanitizedText, isValid := ValidateAndSanitizeSerial(text)
	if !isValid {
		return uc.lineRepo.ReplyMessage(msg.ReplyToken, constants.MsgRequestPhoto)
	}

	// Try to find equipment by id_code or serial_no
	equipment, err := uc.equipmentRepo.FindBySerialOrCode(sanitizedText)
	if err != nil {
		log.Printf("❌ DB lookup error: %v", err)
		return uc.lineRepo.ReplyMessage(msg.ReplyToken, constants.MsgDBLookupFailed)
	}

	if equipment != nil {
		log.Printf("✅ Found equipment by text query: %s", sanitizedText)
		uc.sessionStore.Delete(msg.UserID)
		return uc.lineRepo.ReplyFlexMessage(msg.ReplyToken, "เลือกการดำเนินการ", templates.GetActionMenuFlex(sanitizedText))
	}

	log.Printf("⚠️ Equipment not found for text: %s", sanitizedText)
	return uc.lineRepo.ReplyFlexMessage(msg.ReplyToken, "ไม่พบข้อมูล", templates.GetOCRNotFoundFlex(sanitizedText))
}

// handleTrackStatusInput handles user input when in track status mode.
func (uc *MessageUseCase) handleTrackStatusInput(msg *model.IncomingMessage, text string) error {
	sanitizedText := SanitizeInput(text)

	// Check if input looks like a ticket number (TK- or REQ-)
	upperText := strings.ToUpper(sanitizedText)
	if strings.HasPrefix(upperText, "TK-") || strings.HasPrefix(upperText, "REQ-") {
		// Look up specific ticket by number
		ticket, err := uc.ticketUseCase.GetTicketByNo(upperText, msg.UserID)
		if err != nil {
			log.Printf("❌ Ticket lookup error: %v", err)
			uc.sessionStore.Delete(msg.UserID)
			return uc.lineRepo.ReplyMessage(msg.ReplyToken, "❌ ไม่พบ Ticket หมายเลข: "+sanitizedText+"\nกรุณาตรวจสอบหมายเลขอีกครั้งค่ะ")
		}
		uc.sessionStore.Delete(msg.UserID)
		return uc.lineRepo.ReplyFlexMessage(msg.ReplyToken, "สถานะ Ticket", templates.GetTicketStatusFlex(ticket))
	}

	// Otherwise, try to show all user's tickets
	tickets, err := uc.ticketUseCase.GetUserTickets(msg.UserID)
	if err != nil {
		log.Printf("❌ GetUserTickets error: %v", err)
		uc.sessionStore.Delete(msg.UserID)
		return uc.lineRepo.ReplyMessage(msg.ReplyToken, "❌ ไม่สามารถดึงข้อมูลได้ กรุณาลองใหม่ค่ะ")
	}

	uc.sessionStore.Delete(msg.UserID)

	if len(tickets) == 0 {
		return uc.lineRepo.ReplyMessage(msg.ReplyToken, "📋 คุณยังไม่มีรายการแจ้งปัญหาค่ะ")
	}

	return uc.lineRepo.ReplyFlexMessage(msg.ReplyToken, "รายการแจ้งปัญหาของคุณ", templates.GetMyTicketsFlex(tickets))
}

// handleInputIssueDescInput handles user input when waiting for issue description.
func (uc *MessageUseCase) handleInputIssueDescInput(msg *model.IncomingMessage, text string) error {
	session := uc.sessionStore.Get(msg.UserID)
	if session == nil || session.SerialNumber == "" {
		uc.sessionStore.Delete(msg.UserID)
		return uc.lineRepo.ReplyMessage(msg.ReplyToken, constants.MsgSelectMenuFirst)
	}

	// Sanitize input
	description := SanitizeInput(text)

	// Get user profile for ticket creation
	displayName := ""
	photoURL := ""

	var profile *model.UserProfile
	var err error

	switch msg.SourceType {
	case "group":
		profile, err = uc.lineRepo.GetGroupMemberProfile(msg.GroupID, msg.UserID)
	case "room":
		profile, err = uc.lineRepo.GetRoomMemberProfile(msg.GroupID, msg.UserID)
	default:
		profile, err = uc.lineRepo.GetProfile(msg.UserID)
	}

	if err != nil {
		log.Printf("❌ Failed to get user profile: %v", err)
		displayName = "LINE User"
	} else if profile != nil {
		displayName = profile.DisplayName
		photoURL = profile.PictureURL
	} else {
		displayName = "LINE User"
	}

	// Create ticket with category from session
	ticket, err := uc.ticketUseCase.CreateTicketFromLINE(
		session.SerialNumber,
		description,
		msg.UserID,
		displayName,
		photoURL,
		session.CategoryID,
	)

	// Clear session first
	uc.sessionStore.Delete(msg.UserID)

	if err != nil {
		// Check if it's a duplicate ticket error
		if err == ErrDuplicateTicket && ticket != nil {
			log.Printf("⚠️ Duplicate ticket found: %s", ticket.TicketNo)
			return uc.lineRepo.ReplyFlexMessage(msg.ReplyToken, "พบรายการซ้ำ", templates.GetDuplicateTicketFlex(ticket.TicketNo, session.SerialNumber, ticket.Status.GetStatusText()))
		}
		log.Printf("❌ Failed to create ticket: %v", err)
		return uc.lineRepo.ReplyMessage(msg.ReplyToken, constants.MsgIssueReportFailed)
	}

	// Show success with ticket info
	return uc.lineRepo.ReplyFlexMessage(msg.ReplyToken, "สร้าง Ticket สำเร็จ", templates.GetTicketCreatedFlex(ticket))
}

// HandleImageMessage handles incoming image message - processes OCR
func (uc *MessageUseCase) HandleImageMessage(msg *model.IncomingMessage) error {
	log.Printf("🖼️ Processing image from user: %s, imageID: %s", msg.UserID, msg.ImageID)

	// Check if user has selected menu first
	sess := uc.sessionStore.Get(msg.UserID)
	if sess == nil || sess.Mode != session.ModeReportProblem {
		// User hasn't pressed "แจ้งปัญหา" menu first
		return uc.lineRepo.ReplyMessage(msg.ReplyToken, constants.MsgPleaseSelectReport)
	}

	// Check if OCR client is configured
	if uc.ocrClient == nil {
		log.Println("⚠️ OCR client not configured, using default response")
		return uc.lineRepo.ReplyMessage(msg.ReplyToken, constants.MsgImageReceived)
	}

	// Step 1: Download image from LINE
	imageBytes, err := uc.lineRepo.GetImageContent(msg.ImageID)
	if err != nil {
		log.Printf("❌ Failed to download image: %v", err)
		return uc.lineRepo.ReplyFlexMessage(msg.ReplyToken, "เกิดข้อผิดพลาด", templates.GetOCRErrorFlex())
	}

	// Step 2: Send to OCR API
	ocrResult, err := uc.ocrClient.ProcessImage(imageBytes, msg.ImageID+".jpg")
	if err != nil {
		log.Printf("❌ OCR processing failed: %v", err)
		return uc.lineRepo.ReplyFlexMessage(msg.ReplyToken, "เกิดข้อผิดพลาด", templates.GetOCRErrorFlex())
	}

	// Step 3: Get best text from OCR result
	detectedText := uc.ocrClient.GetDetectedCode(ocrResult)
	if detectedText == "" {
		log.Println("⚠️ OCR detected no text")
		return uc.lineRepo.ReplyFlexMessage(msg.ReplyToken, "อ่านรูปไม่สำเร็จ", templates.GetOCRErrorFlex())
	}

	log.Printf("📝 OCR detected: %s", detectedText)

	// Step 4: Check if equipment exists in DB
	equipment, err := uc.equipmentRepo.FindBySerialOrCode(detectedText)
	if err != nil {
		log.Printf("❌ DB lookup failed: %v", err)
		return uc.lineRepo.ReplyFlexMessage(msg.ReplyToken, "เกิดข้อผิดพลาด", templates.GetOCRErrorFlex())
	}

	if equipment == nil {
		log.Printf("⚠️ Equipment not found: %s", detectedText)
		return uc.lineRepo.ReplyFlexMessage(msg.ReplyToken, "ไม่พบในฐานระบบ", templates.GetOCRNotFoundFlex(detectedText))
	}

	// Step 5: Store session for confirmation
	uc.sessionStore.Set(msg.UserID, &session.OCRSession{
		SerialNumber: detectedText,
	})

	// Step 6: Send confirmation Flex Message
	return uc.lineRepo.ReplyFlexMessage(msg.ReplyToken, "ยืนยันหมายเลข", templates.GetOCRConfirmationFlex(detectedText, ""))
}

// HandleLocationMessage handles incoming location message
func (uc *MessageUseCase) HandleLocationMessage(msg *model.IncomingMessage) error {
	log.Printf("📍 Processing location from user: %s", msg.UserID)
	return uc.lineRepo.ReplyMessage(msg.ReplyToken, constants.MsgLocationReceived)
}

// SendWelcomeMessage sends welcome message to new follower
func (uc *MessageUseCase) SendWelcomeMessage(userID string) error {
	log.Printf("👋 Sending welcome to user: %s", userID)
	return uc.lineRepo.PushMessage(&model.OutgoingMessage{
		To:   userID,
		Text: constants.MsgWelcome,
	})
}

// handleSelectDeptForExpiryInput handles user text input to search department by name.
func (uc *MessageUseCase) handleSelectDeptForExpiryInput(msg *model.IncomingMessage, text string) error {
	ctx := context.Background()
	sanitizedText := SanitizeInput(text)

	// ค้นหาแผนกที่ตรงกับ keyword
	departments, err := uc.departmentRepo.SearchByNameLike(ctx, sanitizedText, 10)
	if err != nil {
		log.Printf("❌ SearchByNameLike error: %v", err)
		return uc.lineRepo.ReplyMessage(msg.ReplyToken, "❌ ไม่สามารถค้นหาแผนกได้ กรุณาลองใหม่ค่ะ")
	}

	if len(departments) == 0 {
		// ไม่เจอแผนก — ไม่ลบ session ให้ลองพิมพ์ใหม่ได้
		return uc.lineRepo.ReplyMessage(msg.ReplyToken, constants.MsgDeptNotFound)
	}

	if len(departments) == 1 {
		// เจอแผนกเดียว → แสดงเครื่องใกล้หมดอายุเลย
		uc.sessionStore.Delete(msg.UserID)
		dept := departments[0]

		expired, err := uc.equipmentRepo.FindExpiredByDepartment(ctx, dept.ID, 999)
		if err != nil {
			log.Printf("❌ FindExpiredByDepartment error: %v", err)
			return uc.lineRepo.ReplyMessage(msg.ReplyToken, "❌ ไม่สามารถดึงข้อมูลได้ กรุณาลองใหม่ค่ะ")
		}
		nearExpiry, err := uc.equipmentRepo.FindNearExpiryByDepartment(ctx, dept.ID, 999)
		if err != nil {
			log.Printf("❌ FindNearExpiryByDepartment error: %v", err)
			return uc.lineRepo.ReplyMessage(msg.ReplyToken, "❌ ไม่สามารถดึงข้อมูลได้ กรุณาลองใหม่ค่ะ")
		}

		if len(expired) == 0 && len(nearExpiry) == 0 {
			return uc.lineRepo.ReplyMessage(msg.ReplyToken, fmt.Sprintf("✅ ไม่มีเครื่องมือที่หมดอายุหรือใกล้หมดอายุในแผนก %s ค่ะ", dept.Name))
		}
		return uc.lineRepo.ReplyFlexMessage(msg.ReplyToken, fmt.Sprintf("เครื่องมือใกล้หมดอายุ - %s", dept.Name), templates.GetEquipmentExpiryByDeptFlex(expired, nearExpiry, dept.Name, dept.ID, uc.baseURL))
	}

	// เจอหลายแผนก → แสดง Flex ให้เลือก
	// ยังไม่ลบ session เผื่อ user อยากพิมพ์ใหม่ (session จะหมดอายุเองหรือถูกลบเมื่อเลือกแผนกผ่าน postback)
	return uc.lineRepo.ReplyFlexMessage(msg.ReplyToken, "ผลค้นหาแผนก", templates.GetDepartmentSelectionWithInputFlex(departments))
}

// Helper functions for equipment data formatting
func getReplacementYear(year *int) string {
	if year == nil {
		return "N/A"
	}
	return fmt.Sprintf("%d", *year)
}

func getModelName(e *entity.Equipment) string {
	if e.Model.ModelName != "" {
		return e.Model.ModelName
	}
	return "N/A"
}

func getBrandName(e *entity.Equipment) string {
	if e.Model.Brand.Name != "" {
		return e.Model.Brand.Name
	}
	return "N/A"
}

func getDepartmentName(e *entity.Equipment) string {
	if e.Department.Name != "" {
		return e.Department.Name
	}
	return "N/A"
}

func getReceiveDate(e *entity.Equipment) string {
	if e.ReceiveDate != nil {
		return e.ReceiveDate.Format("2006-01-02")
	}
	return "N/A"
}
