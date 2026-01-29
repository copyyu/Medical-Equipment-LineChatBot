package usecase

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"medical-webhook/internal/application/service"
	"medical-webhook/internal/domain/line/constants"
	"medical-webhook/internal/domain/line/entity"
	"medical-webhook/internal/domain/line/model"
	"medical-webhook/internal/domain/line/repository"
	"medical-webhook/internal/infrastructure/client"
	"medical-webhook/internal/infrastructure/line/templates"

	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

// MessageUseCase orchestrates message processing flow
type MessageUseCase struct {
	lineRepo       repository.LineRepository
	equipmentRepo  repository.EquipmentRepository
	ocrClient      *client.OCRClient
	sessionStore   *SessionStore
	messageService *service.MessageService
}

// NewMessageUseCase creates a new message use case
func NewMessageUseCase(
	lineRepo repository.LineRepository,
	equipmentRepo repository.EquipmentRepository,
	ocrClient *client.OCRClient,
	sessionStore *SessionStore,
	messageService *service.MessageService,
) *MessageUseCase {
	return &MessageUseCase{
		lineRepo:       lineRepo,
		equipmentRepo:  equipmentRepo,
		ocrClient:      ocrClient,
		sessionStore:   sessionStore,
		messageService: messageService,
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
	case strings.Contains(text, "แจ้งปัญหา") || strings.Contains(text, "เช็กสถานะ"):
		uc.sessionStore.Set(msg.UserID, &OCRSession{Mode: ModeReportProblem})
		return true, uc.lineRepo.ReplyMessage(msg.ReplyToken, constants.MsgReportProblem)

	case strings.Contains(text, "ติดตามสถานะ"):
		uc.sessionStore.Set(msg.UserID, &OCRSession{Mode: ModeTrackStatus})
		return true, uc.lineRepo.ReplyMessage(msg.ReplyToken, constants.MsgTrackStatus)

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
	session := uc.sessionStore.Get(msg.UserID)
	if session == nil || session.Mode == ModeNone {
		return uc.lineRepo.ReplyMessage(msg.ReplyToken, constants.MsgSelectMenuFirst)
	}

	switch session.Mode {
	case ModeReportProblem:
		return uc.handleReportProblemInput(msg, text)
	case ModeTrackStatus:
		return uc.handleTrackStatusInput(msg, text)
	case ModeInputIssueDesc:
		return uc.handleInputIssueDescInput(msg, text)
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
	// Sanitize input before using
	sanitizedText := SanitizeInput(text)
	return uc.lineRepo.ReplyMessage(msg.ReplyToken, "📋 ระบบได้รับข้อมูล: "+sanitizedText+"\nกำลังตรวจสอบสถานะให้ค่ะ")
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

	// Create maintenance record
	err := uc.createMaintenanceRecord(session.SerialNumber, description)
	if err != nil {
		log.Printf("❌ Failed to create maintenance record: %v", err)
		uc.sessionStore.Delete(msg.UserID)
		return uc.lineRepo.ReplyMessage(msg.ReplyToken, constants.MsgIssueReportFailed)
	}

	// Clear session and show success
	uc.sessionStore.Delete(msg.UserID)
	return uc.lineRepo.ReplyFlexMessage(msg.ReplyToken, "บันทึกสำเร็จ", templates.GetIssueSuccessFlex(session.SerialNumber))
}

// createMaintenanceRecord creates a new maintenance record for equipment
func (uc *MessageUseCase) createMaintenanceRecord(serialOrCode, description string) error {
	// Find equipment first
	equipment, err := uc.equipmentRepo.FindBySerialOrCode(serialOrCode)
	if err != nil || equipment == nil {
		return fmt.Errorf("equipment not found: %s", serialOrCode)
	}

	// Create maintenance record
	record := &entity.MaintenanceRecord{
		EquipmentID:     equipment.ID,
		MaintenanceType: entity.MaintenanceCM, // CM = Corrective Maintenance (ซ่อมแซม)
		MaintenanceDate: time.Now(),
		Description:     description,
		Status:          entity.JobStatusInProcess,
	}

	return uc.equipmentRepo.CreateMaintenanceRecord(record)
}

// HandleImageMessage handles incoming image message - processes OCR
func (uc *MessageUseCase) HandleImageMessage(msg *model.IncomingMessage) error {
	log.Printf("🖼️ Processing image from user: %s, imageID: %s", msg.UserID, msg.ImageID)

	// Check if user has selected menu first
	session := uc.sessionStore.Get(msg.UserID)
	if session == nil || session.Mode != ModeReportProblem {
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
		// Equipment not found in DB
		log.Printf("⚠️ Equipment not found: %s", detectedText)
		return uc.lineRepo.ReplyFlexMessage(msg.ReplyToken, "ไม่พบข้อมูล", templates.GetOCRNotFoundFlex(detectedText))
	}

	// Step 5: Store session for confirmation
	uc.sessionStore.Set(msg.UserID, &OCRSession{
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

// HandlePostbackEvent handles postback events from Flex Message buttons
func (uc *MessageUseCase) HandlePostbackEvent(event webhook.PostbackEvent) error {
	data := event.Postback.Data
	replyToken := event.ReplyToken
	log.Printf("📤 Processing postback: %s", data)

	// Parse postback data
	params, _ := url.ParseQuery(data)
	action := params.Get("action")
	serial := params.Get("serial")

	switch action {
	case ActionMainMenu:
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgSelectMenu)

	case ActionRequestChange:
		return uc.lineRepo.ReplyFlexMessage(replyToken, "แจ้งเปลี่ยนเครื่อง", uc.messageService.GetEquipmentChangeFlex("https://www.google.com/"))

	case ActionReportProblem:
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgReportProblem)

	case ActionTrackStatus:
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgTrackStatus)

	case ActionContactStaff:
		return uc.lineRepo.ReplyFlexMessage(replyToken, "ติดต่อเจ้าหน้าที่", uc.messageService.GetContactStaffFlex())

	case ActionOCRConfirmYes:
		// User confirmed OCR result - show action menu (ดูข้อมูล/แจ้งปัญหา)
		if serial != "" {
			return uc.lineRepo.ReplyFlexMessage(replyToken, "เลือกการดำเนินการ", templates.GetActionMenuFlex(serial))
		}
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgSelectMenu)

	case ActionOCRConfirmNo:
		// User denied OCR result - ask for new photo
		return uc.lineRepo.ReplyFlexMessage(replyToken, "ส่งรูปใหม่", templates.GetRetryPhotoFlex())

	case ActionViewRepairHist:
		return uc.handleViewRepairHistory(replyToken, serial)

	case ActionViewLifecycle:
		return uc.handleViewLifecycle(replyToken, serial)

	case ActionViewSpecs:
		return uc.handleViewSpecs(replyToken, serial)

	// New handlers for report issue flow
	case ActionShowActionMenu:
		// Show action menu (ดูข้อมูล/แจ้งปัญหา)
		if serial != "" {
			return uc.lineRepo.ReplyFlexMessage(replyToken, "เลือกการดำเนินการ", templates.GetActionMenuFlex(serial))
		}
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgSelectMenu)

	case ActionViewEquipInfo:
		// Go to equipment info menu (existing)
		if serial != "" {
			return uc.lineRepo.ReplyFlexMessage(replyToken, "ข้อมูลเครื่องมือ", templates.GetEquipmentOptionsFlex(serial))
		}
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgSelectMenu)

	case ActionStartReportIssue:
		// Show issue input menu (พิมพ์รายละเอียด/ข้าม)
		if serial != "" {
			return uc.lineRepo.ReplyFlexMessage(replyToken, "แจ้งปัญหา", templates.GetIssueInputFlex(serial))
		}
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgSelectMenu)

	case ActionInputIssueDesc:
		// Set session mode to wait for issue description
		if serial != "" {
			uc.sessionStore.Set(event.Source.(webhook.UserSource).UserId, &OCRSession{
				Mode:         ModeInputIssueDesc,
				SerialNumber: serial,
			})
			return uc.lineRepo.ReplyMessage(replyToken, constants.MsgInputIssueDesc)
		}
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgSelectMenu)

	case ActionSubmitIssue:
		// Submit issue without description (skip)
		if serial != "" {
			desc := params.Get("desc") // empty for skip
			err := uc.createMaintenanceRecord(serial, desc)
			if err != nil {
				log.Printf("❌ Failed to create maintenance record: %v", err)
				return uc.lineRepo.ReplyMessage(replyToken, constants.MsgIssueReportFailed)
			}
			return uc.lineRepo.ReplyFlexMessage(replyToken, "บันทึกสำเร็จ", templates.GetIssueSuccessFlex(serial))
		}
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgSelectMenu)

	default:
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgSelectMenu)
	}
}

// handleViewRepairHistory sends repair history for equipment
func (uc *MessageUseCase) handleViewRepairHistory(replyToken, serial string) error {
	equipment, err := uc.equipmentRepo.FindBySerialOrCode(serial)
	if err != nil || equipment == nil {
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgEquipmentNotFound)
	}

	records, err := uc.equipmentRepo.GetMaintenanceRecords(equipment.ID)
	if err != nil {
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgRepairHistoryFail)
	}

	// Convert to map format for template
	recordMaps := make([]map[string]interface{}, len(records))
	for i, r := range records {
		recordMaps[i] = map[string]interface{}{
			"date":        r.MaintenanceDate.Format("2006-01-02"),
			"type":        string(r.MaintenanceType),
			"description": r.Description,
			"cost":        r.Cost,
		}
	}

	return uc.lineRepo.ReplyFlexMessage(replyToken, "ประวัติการซ่อม", templates.GetRepairHistoryFlex(serial, recordMaps))
}

// handleViewLifecycle sends lifecycle info for equipment
func (uc *MessageUseCase) handleViewLifecycle(replyToken, serial string) error {
	equipment, err := uc.equipmentRepo.FindBySerialOrCode(serial)
	if err != nil || equipment == nil {
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgEquipmentNotFound)
	}

	data := map[string]interface{}{
		"equipment_age":    equipment.EquipmentAge,
		"life_expectancy":  equipment.LifeExpectancy,
		"remain_life":      equipment.RemainLife,
		"useful_percent":   equipment.UsefulLifetimePercent,
		"replacement_year": getReplacementYear(equipment.ReplacementYear),
	}

	return uc.lineRepo.ReplyFlexMessage(replyToken, "อายุ/วงจรชีวิต", templates.GetLifecycleFlex(serial, data))
}

// handleViewSpecs sends specs info for equipment
func (uc *MessageUseCase) handleViewSpecs(replyToken, serial string) error {
	equipment, err := uc.equipmentRepo.FindBySerialOrCode(serial)
	if err != nil || equipment == nil {
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgEquipmentNotFound)
	}

	data := map[string]interface{}{
		"model_name":   getModelName(equipment),
		"brand":        getBrandName(equipment),
		"department":   getDepartmentName(equipment),
		"receive_date": getReceiveDate(equipment),
		"price":        equipment.PurchasePrice,
	}

	return uc.lineRepo.ReplyFlexMessage(replyToken, "สเปกเครื่อง", templates.GetSpecsFlex(serial, data))
}

// SendWelcomeMessage sends welcome message to new follower
func (uc *MessageUseCase) SendWelcomeMessage(userID string) error {
	log.Printf("👋 Sending welcome to user: %s", userID)
	return uc.lineRepo.PushMessage(&model.OutgoingMessage{
		To:   userID,
		Text: constants.MsgWelcome,
	})
}

// Helper functions
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
