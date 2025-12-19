package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

// WebhookHandler handles LINE webhook events
type WebhookHandler struct {
	bot    *messaging_api.MessagingApiAPI
	secret string
	token  string
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(channelToken, channelSecret string) (*WebhookHandler, error) {
	bot, err := messaging_api.NewMessagingApiAPI(channelToken)
	if err != nil {
		return nil, err
	}

	return &WebhookHandler{
		bot:    bot,
		secret: channelSecret,
		token:  channelToken,
	}, nil
}

// HandleCallback handles webhook callback from LINE
func (h *WebhookHandler) HandleCallback(c *fiber.Ctx) error {
	log.Println("📩 Received webhook callback")

	req, err := http.NewRequest(c.Method(), c.OriginalURL(), nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	c.Request().Header.VisitAll(func(key, value []byte) {
		req.Header.Set(string(key), string(value))
	})

	bodyData := c.Body()
	req.Body = &readCloser{Reader: bytes.NewReader(bodyData)}

	cb, err := webhook.ParseRequest(h.secret, req)
	if err != nil {
		log.Printf("❌ Error parsing webhook: %v", err)
		if err == webhook.ErrInvalidSignature {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid signature"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	for _, event := range cb.Events {
		h.handleEvent(event)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}

type readCloser struct {
	*bytes.Reader
}

func (rc *readCloser) Close() error { return nil }

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

func (h *WebhookHandler) handleMessageEvent(event webhook.MessageEvent) {
	replyToken := event.ReplyToken

	switch msg := event.Message.(type) {
	case webhook.TextMessageContent:
		log.Printf("📝 Text message: %s", msg.Text)
		h.processTextMessage(replyToken, msg.Text)
	case webhook.ImageMessageContent:
		log.Println("🖼️ Received image message")
		h.replyText(replyToken, "ได้รับรูปภาพเรียบร้อยแล้ว กรุณารอเจ้าหน้าที่ตรวจสอบ")
	case webhook.LocationMessageContent:
		log.Printf("📍 Location: %s", msg.Address)
		h.replyText(replyToken, "ได้รับตำแหน่งของคุณแล้ว")
	default:
		log.Printf("Unhandled message type: %T", msg)
	}
}

func (h *WebhookHandler) processTextMessage(replyToken, text string) {
	textLower := strings.ToLower(strings.TrimSpace(text))

	switch {
	case textLower == "แจ้งเปลี่ยนเครื่อง" || strings.Contains(textLower, "เปลี่ยนเครื่อง"):
		h.replyFlexMessage(replyToken, "แจ้งเปลี่ยนเครื่อง", GetEquipmentChangeFlex("https://www.google.com/"))
	case textLower == "ติดต่อ" || textLower == "ติดต่อเจ้าหน้าที่":
		h.replyFlexMessage(replyToken, "ติดต่อเจ้าหน้าที่", GetContactStaffFlex())
	default:
		h.replyFlexMessage(replyToken, "เมนูหลัก", GetMainMenuFlex())
	}
}

func (h *WebhookHandler) handleFollowEvent(event webhook.FollowEvent) {
	log.Println("👋 New follower!")
	var userID string
	switch source := event.Source.(type) {
	case webhook.UserSource:
		userID = source.UserId
	}
	if userID != "" {
		h.pushFlexMessage(userID, "ยินดีต้อนรับ", GetMainMenuFlex())
	}
}

func (h *WebhookHandler) handleUnfollowEvent(event webhook.UnfollowEvent) {
	log.Println("😢 User unfollowed")
}

func (h *WebhookHandler) handlePostbackEvent(event webhook.PostbackEvent) {
	data := event.Postback.Data
	replyToken := event.ReplyToken
	log.Printf("📤 Postback data: %s", data)

	switch data {
	case "action=main_menu":
		h.replyFlexMessage(replyToken, "เมนูหลัก", GetMainMenuFlex())
	case "action=request_change":
		h.replyFlexMessage(replyToken, "แจ้งเปลี่ยนเครื่อง", GetEquipmentChangeFlex("https://www.google.com/"))
	case "action=report_problem":
		h.replyText(replyToken, "🔧 แจ้งปัญหา / เช็กสถานะเครื่อง\n━━━━━━━━━━━━━━━\nกรุณาถ่ายรูปป้าย Serial Number\nหรือสแกน QR Code บนเครื่อง\n\n📸 ส่งรูปมาได้เลยครับ")
	case "action=track_status":
		h.replyText(replyToken, "📋 ติดตามสถานะ\n━━━━━━━━━━━━━━━\nกรุณาระบุหมายเลข Ticket\nหรือ Serial Number ของเครื่อง\n\nตัวอย่าง: TK-2024001")
	case "action=contact_staff":
		h.replyFlexMessage(replyToken, "ติดต่อเจ้าหน้าที่", GetContactStaffFlex())
	default:
		h.replyFlexMessage(replyToken, "เมนูหลัก", GetMainMenuFlex())
	}
}

func (h *WebhookHandler) replyText(replyToken, text string) {
	if replyToken == "" {
		return
	}
	_, err := h.bot.ReplyMessage(&messaging_api.ReplyMessageRequest{
		ReplyToken: replyToken,
		Messages: []messaging_api.MessageInterface{
			&messaging_api.TextMessage{Text: text},
		},
	})
	if err != nil {
		log.Printf("❌ Error replying text: %v", err)
	}
}

func (h *WebhookHandler) replyFlexMessage(replyToken, altText string, flexContent map[string]interface{}) {
	if replyToken == "" {
		return
	}

	requestBody := map[string]interface{}{
		"replyToken": replyToken,
		"messages": []map[string]interface{}{
			{"type": "flex", "altText": altText, "contents": flexContent},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("❌ Error marshaling request: %v", err)
		return
	}

	req, err := http.NewRequest("POST", "https://api.line.me/v2/bot/message/reply", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("❌ Error creating request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Error sending flex message: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("❌ LINE API error %d: %s", resp.StatusCode, string(bodyBytes))
	} else {
		log.Println("✅ Flex message sent successfully")
	}
}

func (h *WebhookHandler) pushFlexMessage(userID, altText string, flexContent map[string]interface{}) {
	requestBody := map[string]interface{}{
		"to": userID,
		"messages": []map[string]interface{}{
			{"type": "flex", "altText": altText, "contents": flexContent},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("❌ Error marshaling request: %v", err)
		return
	}

	req, err := http.NewRequest("POST", "https://api.line.me/v2/bot/message/push", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("❌ Error creating request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Error sending push flex message: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("❌ LINE API error %d: %s", resp.StatusCode, string(bodyBytes))
	} else {
		log.Println("✅ Push flex message sent successfully")
	}
}
