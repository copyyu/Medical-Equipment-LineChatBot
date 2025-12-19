package usecase

import (
	"log"
	"medical-webhook/internal/domain/line/model"
	"medical-webhook/internal/domain/line/repository"
	"medical-webhook/internal/domain/line/service"
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

// HandleTextMessage handles incoming text message
func (uc *MessageUseCase) HandleTextMessage(msg *model.IncomingMessage) error {
	log.Printf("📝 Processing text: %s", msg.Text)

	// Use service to get response
	responseText := uc.messageService.ProcessTextCommand(msg.Text)

	// Use repository to send reply
	return uc.lineRepo.ReplyMessage(msg.ReplyToken, responseText)
}

// HandleImageMessage handles incoming image message
func (uc *MessageUseCase) HandleImageMessage(msg *model.IncomingMessage) error {
	log.Printf("🖼️ Processing image from user: %s", msg.UserID)

	responseText := "ได้รับรูปภาพเรียบร้อยแล้ว กรุณารอเจ้าหน้าที่ตรวจสอบ"
	return uc.lineRepo.ReplyMessage(msg.ReplyToken, responseText)
}

// HandleLocationMessage handles incoming location message
func (uc *MessageUseCase) HandleLocationMessage(msg *model.IncomingMessage) error {
	log.Printf("📍 Processing location from user: %s", msg.UserID)

	responseText := "ได้รับตำแหน่งของคุณแล้ว"
	return uc.lineRepo.ReplyMessage(msg.ReplyToken, responseText)
}

// SendWelcomeMessage sends welcome message to new follower
func (uc *MessageUseCase) SendWelcomeMessage(userID string) error {
	log.Printf("👋 Sending welcome to user: %s", userID)

	welcomeText := uc.messageService.GetFollowerWelcomeMessage()

	msg := &model.OutgoingMessage{
		To:   userID,
		Text: welcomeText,
	}

	return uc.lineRepo.PushMessage(msg)
}
