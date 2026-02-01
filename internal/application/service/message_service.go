package service

import (
	"fmt"

	"medical-webhook/internal/config"
	"medical-webhook/internal/infrastructure/line/templates"
)

// MessageService handles message business logic
type MessageService struct {
	contact config.ContactConfig
}

// NewMessageService creates a new message service with contact configuration
func NewMessageService(contact config.ContactConfig) *MessageService {
	return &MessageService{
		contact: contact,
	}
}

// ProcessTextCommand processes text command and returns appropriate response
func (s *MessageService) ProcessTextCommand(text string) string {
	switch text {
	case "เมนู", "menu", "Menu":
		return s.GetMenuMessage()
	case "แจ้งซ่อม":
		return s.GetRepairFormMessage()
	case "ติดตาม":
		return s.GetTrackingFormMessage()
	case "สอบถาม":
		return s.GetInquiryFormMessage()
	case "ติดต่อ":
		return s.GetContactMessage()
	default:
		return s.GetDefaultMessage()
	}
}

// Text Message Methods
func (s *MessageService) GetMenuMessage() string {
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
}

func (s *MessageService) GetRepairFormMessage() string {
	return `🔧 แจ้งซ่อมเครื่องมือแพทย์
━━━━━━━━━━━━━━━
กรุณาระบุข้อมูลดังนี้:

📍 ชื่อเครื่องมือ:
📍 รหัสเครื่อง:
📍 แผนก/หน่วยงาน:
📍 อาการเสีย:
📍 ชื่อผู้แจ้ง:
📍 เบอร์ติดต่อ:`
}

func (s *MessageService) GetTrackingFormMessage() string {
	return `🔍 ติดตามสถานะการซ่อม
━━━━━━━━━━━━━━━
กรุณาระบุหมายเลข Ticket
หรือรหัสเครื่องมือที่ต้องการติดตาม

ตัวอย่าง: TK-2024001`
}

func (s *MessageService) GetInquiryFormMessage() string {
	return `ℹ️ สอบถามข้อมูลเครื่องมือ
━━━━━━━━━━━━━━━
กรุณาพิมพ์ชื่อหรือรหัสเครื่องมือ
ที่ต้องการสอบถาม`
}

// GetContactMessage returns contact information from config
func (s *MessageService) GetContactMessage() string {
	msg := fmt.Sprintf(`📞 ติดต่อเจ้าหน้าที่
━━━━━━━━━━━━━━━
🏥 %s`, s.contact.CenterName)

	if s.contact.Phone != "" {
		msg += fmt.Sprintf("\n\n📱 โทร: %s", s.contact.Phone)
	}
	if s.contact.Email != "" {
		msg += fmt.Sprintf("\n📧 Email: %s", s.contact.Email)
	}
	if s.contact.WorkingHours != "" {
		msg += fmt.Sprintf("\n⏰ เวลาทำการ: %s", s.contact.WorkingHours)
	}
	if s.contact.EmergencyPhone != "" {
		msg += fmt.Sprintf("\n\n🚨 กรณีฉุกเฉิน: %s (24 ชม.)", s.contact.EmergencyPhone)
	}

	return msg
}

func (s *MessageService) GetDefaultMessage() string {
	return s.GetMenuMessage()
}

func (s *MessageService) GetFollowerWelcomeMessage() string {
	return `🏥 ยินดีต้อนรับสู่ระบบเครื่องมือแพทย์!
━━━━━━━━━━━━━━━
ขอบคุณที่เพิ่มเราเป็นเพื่อน

พิมพ์ "เมนู" เพื่อเริ่มใช้งาน`
}

// Flex Message Methods
// GetEquipmentChangeFlex returns a Flex Message for equipment change request
func (s *MessageService) GetEquipmentChangeFlex(linkURL string) map[string]interface{} {
	return templates.GetEquipmentChangeFlex(linkURL)
}

// GetContactStaffFlex returns a Flex Message for contact information
func (s *MessageService) GetContactStaffFlex() map[string]interface{} {
	return templates.GetContactStaffFlex()
}
