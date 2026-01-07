package usecase

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"medical-webhook/internal/domain/line/entity"
	"medical-webhook/internal/domain/line/model"
	"medical-webhook/internal/domain/line/repository"
	"medical-webhook/internal/domain/line/service"
	"medical-webhook/internal/domain/line/templates"
	"medical-webhook/internal/infrastructure/client"

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

// HandleTextMessage handles incoming text message from Rich Menu or direct input
func (uc *MessageUseCase) HandleTextMessage(msg *model.IncomingMessage) error {
	log.Printf("📝 Processing text: %s", msg.Text)
	text := strings.TrimSpace(msg.Text)

	switch {
	// Rich Menu: แจ้งปัญหา / เช็กสถานะ
	case strings.Contains(text, "แจ้งปัญหา") || strings.Contains(text, "เช็กสถานะ"):
		// Set session mode
		uc.sessionStore.Set(msg.UserID, &OCRSession{Mode: ModeReportProblem})
		return uc.lineRepo.ReplyMessage(msg.ReplyToken, "🔧 แจ้งปัญหา / เช็กสถานะเครื่อง\n━━━━━━━━━━━━━━━\nกรุณาถ่ายรูปป้าย Serial Number\nหรือพิมพ์รหัสเครื่อง (ID Code)\n\n📸 ส่งรูปมาได้เลยค่ะ\n✏️ หรือพิมพ์รหัส เช่น SSH12345")

	// Rich Menu: ติดตามสถานะ
	case strings.Contains(text, "ติดตามสถานะ"):
		// Set session mode
		uc.sessionStore.Set(msg.UserID, &OCRSession{Mode: ModeTrackStatus})
		return uc.lineRepo.ReplyMessage(msg.ReplyToken, "📋 ติดตามสถานะ\n━━━━━━━━━━━━━━━\nกรุณาระบุหมายเลข Ticket\nหรือ Serial Number ของเครื่อง\n\nตัวอย่าง: TK-2024001")

	// Rich Menu: แจ้งเปลี่ยนเครื่อง
	case strings.Contains(text, "เปลี่ยนเครื่อง"):
		return uc.lineRepo.ReplyFlexMessage(msg.ReplyToken, "แจ้งเปลี่ยนเครื่อง", uc.messageService.GetEquipmentChangeFlex("https://www.google.com/"))

	// Rich Menu: ติดต่อเจ้าหน้าที่
	case strings.Contains(text, "ติดต่อ"):
		return uc.lineRepo.ReplyFlexMessage(msg.ReplyToken, "ติดต่อเจ้าหน้าที่", uc.messageService.GetContactStaffFlex())

	// เรียกเมนู (ไม่ต้องแสดง Flex แล้ว ใช้ Rich Menu แทน)
	case text == "เมนู" || strings.ToLower(text) == "menu":
		return uc.lineRepo.ReplyMessage(msg.ReplyToken, "กรุณาเลือกบริการจากเมนูด้านล่างค่ะ 👇")

	default:
		// Check if user has selected a menu first
		session := uc.sessionStore.Get(msg.UserID)
		if session == nil || session.Mode == ModeNone {
			// No menu selected - ask to select menu first
			return uc.lineRepo.ReplyMessage(msg.ReplyToken, "ขออภัยค่ะ เพื่อให้สามารถตอบคำถามได้ถูกต้อง กรุณาเลือกบริการที่ต้องการจากเมนูด้านล่างค่ะ 🙇🏻‍♀️")
		}

		// User has selected menu - process based on mode
		if session.Mode == ModeReportProblem {
			// Check if text looks like an equipment id_code or serial
			if len(text) >= 3 && isAlphanumeric(text) {
				// Try to find equipment by id_code or serial_no
				equipment, err := uc.equipmentRepo.FindBySerialOrCode(text)
				if err == nil && equipment != nil {
					log.Printf("✅ Found equipment by text query: %s", text)
					// Clear session after success
					uc.sessionStore.Delete(msg.UserID)
					return uc.lineRepo.ReplyFlexMessage(msg.ReplyToken, "ข้อมูลเครื่องมือ", templates.GetEquipmentOptionsFlex(text))
				}
				if err == nil && equipment == nil {
					log.Printf("⚠️ Equipment not found for text: %s", text)
					return uc.lineRepo.ReplyFlexMessage(msg.ReplyToken, "ไม่พบข้อมูล", templates.GetOCRNotFoundFlex(text))
				}
			}
			return uc.lineRepo.ReplyMessage(msg.ReplyToken, "กรุณาพิมพ์รหัสเครื่อง หรือส่งรูปป้าย Serial Number ค่ะ")
		}

		if session.Mode == ModeTrackStatus {
			// For track status - just acknowledge for now
			return uc.lineRepo.ReplyMessage(msg.ReplyToken, "📋 ระบบได้รับข้อมูล: "+text+"\nกำลังตรวจสอบสถานะให้ค่ะ")
		}

		return uc.lineRepo.ReplyMessage(msg.ReplyToken, "ขออภัยค่ะ กรุณาเลือกบริการจากเมนูด้านล่างค่ะ 🙇🏻‍♀️")
	}
}

// isAlphanumeric checks if string contains only letters, numbers, and common separators
func isAlphanumeric(s string) bool {
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_') {
			return false
		}
	}
	return true
}

// HandleImageMessage handles incoming image message - processes OCR
func (uc *MessageUseCase) HandleImageMessage(msg *model.IncomingMessage) error {
	log.Printf("🖼️ Processing image from user: %s, imageID: %s", msg.UserID, msg.ImageID)

	// Check if user has selected menu first
	session := uc.sessionStore.Get(msg.UserID)
	if session == nil || session.Mode != ModeReportProblem {
		// User hasn't pressed "แจ้งปัญหา" menu first
		return uc.lineRepo.ReplyMessage(msg.ReplyToken, "ขออภัยค่ะ กรุณากดเมนู \"แจ้งปัญหา / เช็กสถานะ\" ก่อนส่งรูปค่ะ 🙇🏻‍♀️")
	}

	// Check if OCR client is configured
	if uc.ocrClient == nil {
		log.Println("⚠️ OCR client not configured, using default response")
		return uc.lineRepo.ReplyMessage(msg.ReplyToken, "ได้รับรูปภาพเรียบร้อยแล้ว กรุณารอเจ้าหน้าที่ตรวจสอบ")
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
	return uc.lineRepo.ReplyMessage(msg.ReplyToken, "ได้รับตำแหน่งของคุณแล้ว")
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
	case "main_menu":
		return uc.lineRepo.ReplyMessage(replyToken, "กรุณาเลือกบริการจากเมนูด้านล่างค่ะ 👇")

	case "request_change":
		return uc.lineRepo.ReplyFlexMessage(replyToken, "แจ้งเปลี่ยนเครื่อง", uc.messageService.GetEquipmentChangeFlex("https://www.google.com/"))

	case "report_problem":
		return uc.lineRepo.ReplyMessage(replyToken, " แจ้งปัญหา / เช็กสถานะเครื่อง\n━━━━━━━━━━━━━━━\nกรุณาถ่ายรูปป้าย Serial Number\nหรือสแกน QR Code บนเครื่อง\n\n📸 ส่งรูปมาได้เลยครับ")

	case "track_status":
		return uc.lineRepo.ReplyMessage(replyToken, " ติดตามสถานะ\n━━━━━━━━━━━━━━━\nกรุณาระบุหมายเลข Ticket\nหรือ Serial Number ของเครื่อง\n\nตัวอย่าง: TK-2024001")

	case "contact_staff":
		return uc.lineRepo.ReplyFlexMessage(replyToken, "ติดต่อเจ้าหน้าที่", uc.messageService.GetContactStaffFlex())

	case "ocr_confirm_yes":
		// User confirmed OCR result - show equipment options
		if serial != "" {
			return uc.lineRepo.ReplyFlexMessage(replyToken, "ข้อมูลเครื่องมือ", templates.GetEquipmentOptionsFlex(serial))
		}
		return uc.lineRepo.ReplyMessage(replyToken, "กรุณาเลือกบริการจากเมนูด้านล่างค่ะ 👇")

	case "ocr_confirm_no":
		// User denied OCR result - ask for new photo
		return uc.lineRepo.ReplyFlexMessage(replyToken, "ส่งรูปใหม่", templates.GetRetryPhotoFlex())

	case "view_repair_history":
		return uc.handleViewRepairHistory(replyToken, serial)

	case "view_lifecycle":
		return uc.handleViewLifecycle(replyToken, serial)

	case "view_specs":
		return uc.handleViewSpecs(replyToken, serial)

	default:
		return uc.lineRepo.ReplyMessage(replyToken, "กรุณาเลือกบริการจากเมนูด้านล่างค่ะ 👇")
	}
}

// handleViewRepairHistory sends repair history for equipment
func (uc *MessageUseCase) handleViewRepairHistory(replyToken, serial string) error {
	equipment, err := uc.equipmentRepo.FindBySerialOrCode(serial)
	if err != nil || equipment == nil {
		return uc.lineRepo.ReplyMessage(replyToken, "❌ ไม่พบข้อมูลเครื่องมือ")
	}

	records, err := uc.equipmentRepo.GetMaintenanceRecords(equipment.ID)
	if err != nil {
		return uc.lineRepo.ReplyMessage(replyToken, "❌ ไม่สามารถดึงประวัติการซ่อมได้")
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
		return uc.lineRepo.ReplyMessage(replyToken, "❌ ไม่พบข้อมูลเครื่องมือ")
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
		return uc.lineRepo.ReplyMessage(replyToken, "❌ ไม่พบข้อมูลเครื่องมือ")
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
		Text: "ยินดีต้อนรับสู่ระบบเครื่องมือแพทย์ค่ะ 🏥\n━━━━━━━━━━━━━━━\nกรุณาเลือกบริการจากเมนูด้านล่างได้เลยค่ะ 👇",
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
