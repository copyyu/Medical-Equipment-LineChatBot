package usecase

// Message constants - ข้อความที่ใช้ตอบกลับผู้ใช้
const (
	// เมนูและการนำทาง
	MsgSelectMenu      = "กรุณาเลือกบริการจากเมนูด้านล่างค่ะ 👇"
	MsgSelectMenuFirst = "ขออภัยค่ะ เพื่อให้สามารถตอบคำถามได้ถูกต้อง กรุณาเลือกบริการที่ต้องการจากเมนูด้านล่างค่ะ 🙇🏻‍♀️"

	// ข้อผิดพลาด
	MsgEquipmentNotFound = "❌ ไม่พบข้อมูลเครื่องมือ"
	MsgDBLookupFailed    = "❌ ไม่สามารถดึงข้อมูลได้"
	MsgRepairHistoryFail = "❌ ไม่สามารถดึงประวัติการซ่อมได้"

	// แจ้งปัญหา / เช็กสถานะ
	MsgReportProblem = "🔧 แจ้งปัญหา / เช็กสถานะเครื่อง\n━━━━━━━━━━━━━━━\nกรุณาถ่ายรูปป้าย Serial Number\nหรือพิมพ์รหัสเครื่อง (ID Code)\n\n📸 ส่งรูปมาได้เลยค่ะ\n✏️ หรือพิมพ์รหัส เช่น SSH12345"

	// ติดตามสถานะ
	MsgTrackStatus = "📋 ติดตามสถานะ\n━━━━━━━━━━━━━━━\nกรุณาระบุหมายเลข Ticket\nหรือ Serial Number ของเครื่อง\n\nตัวอย่าง: TK-2024001"

	// การส่งรูป
	MsgRequestPhoto       = "กรุณาพิมพ์รหัสเครื่อง หรือส่งรูปป้าย Serial Number ค่ะ"
	MsgPleaseSelectReport = "ขออภัยค่ะ กรุณากดเมนู \"แจ้งปัญหา / เช็กสถานะ\" ก่อนส่งรูปค่ะ 🙇🏻‍♀️"
	MsgImageReceived      = "ได้รับรูปภาพเรียบร้อยแล้ว กรุณารอเจ้าหน้าที่ตรวจสอบ"
	MsgLocationReceived   = "ได้รับตำแหน่งของคุณแล้ว"

	// ข้อความต้อนรับ
	MsgWelcome = "ยินดีต้อนรับสู่ระบบเครื่องมือแพทย์ค่ะ 🏥\n━━━━━━━━━━━━━━━\nกรุณาเลือกบริการจากเมนูด้านล่างได้เลยค่ะ 👇"
)

// Postback action constants - action ที่ใช้ใน postback data
const (
	ActionMainMenu       = "main_menu"
	ActionRequestChange  = "request_change"
	ActionReportProblem  = "report_problem"
	ActionTrackStatus    = "track_status"
	ActionContactStaff   = "contact_staff"
	ActionOCRConfirmYes  = "ocr_confirm_yes"
	ActionOCRConfirmNo   = "ocr_confirm_no"
	ActionViewRepairHist = "view_repair_history"
	ActionViewLifecycle  = "view_lifecycle"
	ActionViewSpecs      = "view_specs"
)

// Validation constants
const (
	MinSerialLength = 3
	MaxInputLength  = 100
)
