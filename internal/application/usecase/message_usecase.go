package usecase

import (
	"log"
	"strings"

	"medical-webhook/internal/domain/line/model"
	"medical-webhook/internal/domain/line/repository"
	"medical-webhook/internal/domain/line/service"

	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

// MessageUseCase orchestrates message processing flow
type MessageUseCase struct {
	lineRepo       repository.LineRepository
	messageService *service.MessageService
}

// NewMessageUseCase creates a new message use case
func NewMessageUseCase(
	lineRepo repository.LineRepository,
	messageService *service.MessageService,
) *MessageUseCase {
	return &MessageUseCase{
		lineRepo:       lineRepo,
		messageService: messageService,
	}
}

// HandleTextMessage handles incoming text message - shows Flex Menu for all messages
func (uc *MessageUseCase) HandleTextMessage(msg *model.IncomingMessage) error {
	log.Printf("📝 Processing text: %s", msg.Text)
	textLower := strings.ToLower(strings.TrimSpace(msg.Text))

	switch {
	case textLower == "แจ้งเปลี่ยนเครื่อง" || strings.Contains(textLower, "เปลี่ยนเครื่อง"):
		return uc.lineRepo.ReplyFlexMessage(msg.ReplyToken, "แจ้งเปลี่ยนเครื่อง", uc.messageService.GetEquipmentChangeFlex("https://www.google.com/"))

	case textLower == "ติดต่อ" || textLower == "ติดต่อเจ้าหน้าที่":
		return uc.lineRepo.ReplyFlexMessage(msg.ReplyToken, "ติดต่อเจ้าหน้าที่", uc.messageService.GetContactStaffFlex())

	default:
		// Default: Show main menu Flex Message
		return uc.lineRepo.ReplyFlexMessage(msg.ReplyToken, "เมนูหลัก", uc.messageService.GetMainMenuFlex())
	}
}

// HandleImageMessage handles incoming image message
func (uc *MessageUseCase) HandleImageMessage(msg *model.IncomingMessage) error {
	log.Printf("🖼️ Processing image from user: %s", msg.UserID)
	return uc.lineRepo.ReplyMessage(msg.ReplyToken, "ได้รับรูปภาพเรียบร้อยแล้ว กรุณารอเจ้าหน้าที่ตรวจสอบ")
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

	switch data {
	case "action=main_menu":
		return uc.lineRepo.ReplyFlexMessage(replyToken, "เมนูหลัก", uc.messageService.GetMainMenuFlex())

	case "action=request_change":
		return uc.lineRepo.ReplyFlexMessage(replyToken, "แจ้งเปลี่ยนเครื่อง", uc.messageService.GetEquipmentChangeFlex("https://www.google.com/"))

	case "action=report_problem":
		return uc.lineRepo.ReplyMessage(replyToken, "🔧 แจ้งปัญหา / เช็กสถานะเครื่อง\n━━━━━━━━━━━━━━━\nกรุณาถ่ายรูปป้าย Serial Number\nหรือสแกน QR Code บนเครื่อง\n\n📸 ส่งรูปมาได้เลยครับ")

	case "action=track_status":
		return uc.lineRepo.ReplyMessage(replyToken, "📋 ติดตามสถานะ\n━━━━━━━━━━━━━━━\nกรุณาระบุหมายเลข Ticket\nหรือ Serial Number ของเครื่อง\n\nตัวอย่าง: TK-2024001")

	case "action=contact_staff":
		return uc.lineRepo.ReplyFlexMessage(replyToken, "ติดต่อเจ้าหน้าที่", uc.messageService.GetContactStaffFlex())

	default:
		return uc.lineRepo.ReplyFlexMessage(replyToken, "เมนูหลัก", uc.messageService.GetMainMenuFlex())
	}
}

// SendWelcomeMessage sends welcome Flex Message to new follower
func (uc *MessageUseCase) SendWelcomeMessage(userID string) error {
	log.Printf("👋 Sending welcome to user: %s", userID)
	return uc.lineRepo.PushFlexMessage(userID, "ยินดีต้อนรับ", uc.messageService.GetMainMenuFlex())
}
