package handlers

import (
	"bytes"
	"log"
	"net/http"

	"medical-webhook/internal/application/usecase"
	"medical-webhook/internal/domain/line/model"

	"github.com/gofiber/fiber/v2"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

// WebhookHandler handles LINE webhook events
type WebhookHandler struct {
	secret         string
	messageUseCase *usecase.MessageUseCase
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(secret string, messageUseCase *usecase.MessageUseCase) *WebhookHandler {
	return &WebhookHandler{
		secret:         secret,
		messageUseCase: messageUseCase,
	}
}

// HandleCallback handles webhook callback from LINE
func (h *WebhookHandler) HandleCallback(c *fiber.Ctx) error {
	log.Println("📩 Received webhook callback")

	// Create a standard http.Request from Fiber context
	req, err := http.NewRequest(c.Method(), c.OriginalURL(), nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Copy headers
	c.Request().Header.VisitAll(func(key, value []byte) {
		req.Header.Set(string(key), string(value))
	})

	// Set body
	bodyData := c.Body()
	req.Body = &readCloser{Reader: bytes.NewReader(bodyData)}

	// Parse webhook request
	cb, err := webhook.ParseRequest(h.secret, req)
	if err != nil {
		log.Printf("❌ Error parsing webhook: %v", err)
		if err == webhook.ErrInvalidSignature {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid signature"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Process events
	for _, event := range cb.Events {
		h.handleEvent(event)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}

// readCloser implements io.ReadCloser for bytes
type readCloser struct {
	*bytes.Reader
}

func (rc *readCloser) Close() error { return nil }

// handleEvent processes individual LINE events
func (h *WebhookHandler) handleEvent(event webhook.EventInterface) {
	switch e := event.(type) {
	case webhook.MessageEvent:
		h.handleMessageEvent(e)
	case webhook.FollowEvent:
		h.handleFollowEvent(e)
	case webhook.UnfollowEvent:
		h.handleUnfollowEvent(e)
	case webhook.PostbackEvent:
		h.handlePostbackEvent(e)
	default:
		log.Printf("Unhandled event type: %T", e)
	}
}

// handleMessageEvent handles message events
func (h *WebhookHandler) handleMessageEvent(event webhook.MessageEvent) {
	var replyToken string
	if event.ReplyToken != "" {
		replyToken = event.ReplyToken
	}

	var userID, groupID, sourceType string

	switch source := event.Source.(type) {
	case webhook.UserSource:
		userID = source.UserId
		sourceType = "user"
	case webhook.GroupSource:
		userID = source.UserId
		groupID = source.GroupId
		sourceType = "group"
	case webhook.RoomSource:
		userID = source.UserId
		groupID = source.RoomId
		sourceType = "room"
	}

	switch msg := event.Message.(type) {
	case webhook.TextMessageContent:
		incomingMsg := &model.IncomingMessage{
			UserID:      userID,
			GroupID:     groupID,
			SourceType:  sourceType,
			MessageType: model.MessageTypeText,
			Text:        msg.Text,
			ReplyToken:  replyToken,
		}
		h.messageUseCase.HandleTextMessage(incomingMsg)

	case webhook.ImageMessageContent:
		incomingMsg := &model.IncomingMessage{
			UserID:      userID,
			GroupID:     groupID,
			SourceType:  sourceType,
			MessageType: model.MessageTypeImage,
			ImageID:     msg.Id,
			ReplyToken:  replyToken,
		}
		h.messageUseCase.HandleImageMessage(incomingMsg)

	case webhook.LocationMessageContent:
		incomingMsg := &model.IncomingMessage{
			UserID:      userID,
			GroupID:     groupID,
			SourceType:  sourceType,
			MessageType: model.MessageTypeLocation,
			Location: &model.Location{
				Latitude:  msg.Latitude,
				Longitude: msg.Longitude,
				Address:   msg.Address,
			},
			ReplyToken: replyToken,
		}
		h.messageUseCase.HandleLocationMessage(incomingMsg)

	default:
		log.Printf("Unhandled message type: %T", msg)
	}
}

// handleFollowEvent handles new follower events
func (h *WebhookHandler) handleFollowEvent(event webhook.FollowEvent) {
	log.Println("👋 New follower!")

	var userID string
	if source, ok := event.Source.(webhook.UserSource); ok {
		userID = source.UserId
	}

	if userID != "" {
		h.messageUseCase.SendWelcomeMessage(userID)
	}
}

// handleUnfollowEvent handles unfollow events
func (h *WebhookHandler) handleUnfollowEvent(event webhook.UnfollowEvent) {
	log.Println("😢 User unfollowed")
}

// handlePostbackEvent handles postback events
func (h *WebhookHandler) handlePostbackEvent(event webhook.PostbackEvent) {
	h.messageUseCase.HandlePostbackEvent(event)
}
