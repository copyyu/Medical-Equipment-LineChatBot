package service

// MessageService handles message business logic
type MessageService struct{}

// NewMessageService creates a new message service
func NewMessageService() *MessageService {
	return &MessageService{}
}

// GetMenuMessage returns menu message text
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

// GetRepairFormMessage returns repair form message
func (s *MessageService) GetRepairFormMessage() string {
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
}

// GetTrackingFormMessage returns tracking form message
func (s *MessageService) GetTrackingFormMessage() string {
	return `🔍 ติดตามสถานะการซ่อม
━━━━━━━━━━━━━━━
กรุณาระบุหมายเลข Ticket
หรือรหัสเครื่องมือที่ต้องการติดตาม

ตัวอย่าง:
ติดตาม TK-2024001
หรือ
ติดตาม ECG-001`
}

// GetInquiryFormMessage returns inquiry form message
func (s *MessageService) GetInquiryFormMessage() string {
	return `ℹ️ สอบถามข้อมูลเครื่องมือ
━━━━━━━━━━━━━━━
กรุณาพิมพ์ชื่อหรือรหัสเครื่องมือ
ที่ต้องการสอบถาม

ตัวอย่าง:
สอบถาม Defibrillator
หรือ
สอบถาม DEF-001`
}

// GetContactMessage returns contact information message
func (s *MessageService) GetContactMessage() string {
	return `📞 ติดต่อเจ้าหน้าที่
━━━━━━━━━━━━━━━
🏥 ศูนย์เครื่องมือแพทย์

📱 โทร: 02-XXX-XXXX
📧 Email: medical-equipment@hospital.com
⏰ เวลาทำการ: จ-ศ 08:00-17:00

🚨 กรณีฉุกเฉิน: 02-XXX-YYYY (24 ชม.)`
}

// GetWelcomeMessage returns welcome message for new follower
func (s *MessageService) GetWelcomeMessage() string {
	return `👋 สวัสดีครับ ยินดีต้อนรับสู่
🏥 ระบบเครื่องมือแพทย์

พิมพ์ "เมนู" เพื่อดูบริการของเรา`
}

// GetDefaultMessage returns default message for unknown input
func (s *MessageService) GetDefaultMessage() string {
	return s.GetWelcomeMessage()
}

// GetFollowerWelcomeMessage returns welcome message when user follows
func (s *MessageService) GetFollowerWelcomeMessage() string {
	return `🏥 ยินดีต้อนรับสู่ระบบเครื่องมือแพทย์!
━━━━━━━━━━━━━━━
ขอบคุณที่เพิ่มเราเป็นเพื่อน

พิมพ์ "เมนู" เพื่อเริ่มใช้งาน`
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
