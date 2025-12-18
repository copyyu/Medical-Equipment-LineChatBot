package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

// WebhookHandler handles LINE webhook events
type WebhookHandler struct {
	bot    *messaging_api.MessagingApiAPI
	secret string
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
	}, nil
}

// HandleCallback handles webhook callback from LINE
func (h *WebhookHandler) HandleCallback(c *gin.Context) {
	log.Println("📩 Received webhook callback")

	cb, err := webhook.ParseRequest(h.secret, c.Request)
	if err != nil {
		log.Printf("❌ Error parsing webhook: %v", err)
		if err == webhook.ErrInvalidSignature {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Process events
	for _, event := range cb.Events {
		h.handleEvent(event)
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

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

// handleMessageEvent handles text and other message events
func (h *WebhookHandler) handleMessageEvent(event webhook.MessageEvent) {
	var replyToken string
	if event.ReplyToken != "" {
		replyToken = event.ReplyToken
	}

	switch msg := event.Message.(type) {
	case webhook.TextMessageContent:
		log.Printf("📝 Text message: %s", msg.Text)
		h.replyText(replyToken, h.processTextMessage(msg.Text))
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

// processTextMessage processes text messages and returns appropriate response
func (h *WebhookHandler) processTextMessage(text string) string {
	// Medical equipment related keywords
	switch text {
	case "เมนู", "menu", "Menu":
		return `🏥 ระบบเครื่องมือแพทย์
━━━━━━━━━━━━━━━
📋 บริการของเรา:

1️⃣ แจ้งซ่อมเครื่องมือแพทย์
   พิมพ์: แจ้งซ่อม

2️⃣ ติดตามสถานะการซ่อม
   พิมพ์: ติดตาม

3️⃣ สอบถามข้อมูลเครื่องมือ
   พิมพ์: สอบถาม

4️⃣ ติดต่อเจ้าหน้าที่
   พิมพ์: ติดต่อ

━━━━━━━━━━━━━━━
พิมพ์ "เมนู" เพื่อดูเมนูอีกครั้ง`

	case "แจ้งซ่อม":
		return `🔧 แจ้งซ่อมเครื่องมือแพทย์
━━━━━━━━━━━━━━━
กรุณาระบุข้อมูลดังนี้:

📍 ชื่อเครื่องมือ:
📍 รหัสเครื่อง:
📍 แผนก/หน่วยงาน:
📍 อาการเสีย:
📍 ชื่อผู้แจ้ง:
📍 เบอร์ติดต่อ:

ตัวอย่าง:
เครื่อง: Monitor ECG
รหัส: ECG-001
แผนก: ICU
อาการ: หน้าจอไม่ติด
ผู้แจ้ง: พยาบาล สมหญิง
เบอร์: 1234`

	case "ติดตาม":
		return `🔍 ติดตามสถานะการซ่อม
━━━━━━━━━━━━━━━
กรุณาระบุหมายเลข Ticket
หรือรหัสเครื่องมือที่ต้องการติดตาม

ตัวอย่าง:
ติดตาม TK-2024001
หรือ
ติดตาม ECG-001`

	case "สอบถาม":
		return `ℹ️ สอบถามข้อมูลเครื่องมือ
━━━━━━━━━━━━━━━
กรุณาพิมพ์ชื่อหรือรหัสเครื่องมือ
ที่ต้องการสอบถาม

ตัวอย่าง:
สอบถาม Defibrillator
หรือ
สอบถาม DEF-001`

	case "ติดต่อ":
		return `📞 ติดต่อเจ้าหน้าที่
━━━━━━━━━━━━━━━
🏥 ศูนย์เครื่องมือแพทย์

📱 โทร: 02-XXX-XXXX
📧 Email: medical-equipment@hospital.com
⏰ เวลาทำการ: จ-ศ 08:00-17:00

🚨 กรณีฉุกเฉิน: 02-XXX-YYYY (24 ชม.)`

	default:
		return `👋 สวัสดีครับ ยินดีต้อนรับสู่
🏥 ระบบเครื่องมือแพทย์

พิมพ์ "เมนู" เพื่อดูบริการของเรา`
	}
}

// handleFollowEvent handles new follower events
func (h *WebhookHandler) handleFollowEvent(event webhook.FollowEvent) {
	log.Println("👋 New follower!")

	var userID string
	switch source := event.Source.(type) {
	case webhook.UserSource:
		userID = source.UserId
	}

	if userID != "" {
		h.pushMessage(userID, `🏥 ยินดีต้อนรับสู่ระบบเครื่องมือแพทย์!
━━━━━━━━━━━━━━━
ขอบคุณที่เพิ่มเราเป็นเพื่อน

พิมพ์ "เมนู" เพื่อเริ่มใช้งาน`)
	}
}

// handleUnfollowEvent handles unfollow events
func (h *WebhookHandler) handleUnfollowEvent(event webhook.UnfollowEvent) {
	log.Println("😢 User unfollowed")
}

// handlePostbackEvent handles postback events
func (h *WebhookHandler) handlePostbackEvent(event webhook.PostbackEvent) {
	log.Printf("📤 Postback data: %s", event.Postback.Data)

	if event.ReplyToken != "" {
		h.replyText(event.ReplyToken, "ได้รับข้อมูลเรียบร้อยแล้ว")
	}
}

// replyText sends a text reply message
func (h *WebhookHandler) replyText(replyToken, text string) {
	if replyToken == "" {
		return
	}

	_, err := h.bot.ReplyMessage(&messaging_api.ReplyMessageRequest{
		ReplyToken: replyToken,
		Messages: []messaging_api.MessageInterface{
			&messaging_api.TextMessage{
				Text: text,
			},
		},
	})

	if err != nil {
		log.Printf("❌ Error replying: %v", err)
	}
}

// pushMessage sends a push message to a user
func (h *WebhookHandler) pushMessage(userID, text string) {
	_, err := h.bot.PushMessage(&messaging_api.PushMessageRequest{
		To: userID,
		Messages: []messaging_api.MessageInterface{
			&messaging_api.TextMessage{
				Text: text,
			},
		},
	}, "")

	if err != nil {
		log.Printf("❌ Error pushing message: %v", err)
	}
}

// SendMessage sends a message to a specific user (for external use)
func (h *WebhookHandler) SendMessage(userID, text string) error {
	_, err := h.bot.PushMessage(&messaging_api.PushMessageRequest{
		To: userID,
		Messages: []messaging_api.MessageInterface{
			&messaging_api.TextMessage{
				Text: text,
			},
		},
	}, "")

	return err
}
