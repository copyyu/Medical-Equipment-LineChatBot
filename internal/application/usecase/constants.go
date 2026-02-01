package usecase

// Postback action constants - action ที่ใช้ใน postback data
const (
	ActionMainMenu         = "main_menu"
	ActionRequestChange    = "request_change"
	ActionReportProblem    = "report_problem"
	ActionTrackStatus      = "track_status"
	ActionContactStaff     = "contact_staff"
	ActionOCRConfirmYes    = "ocr_confirm_yes"
	ActionOCRConfirmNo     = "ocr_confirm_no"
	ActionViewRepairHist   = "view_repair_history"
	ActionViewLifecycle    = "view_lifecycle"
	ActionViewSpecs        = "view_specs"
	ActionShowActionMenu   = "show_action_menu"    // แสดงเมนูเลือก (ดูข้อมูล/แจ้งปัญหา)
	ActionViewEquipInfo    = "view_equipment_info" // ไปหน้าข้อมูลเครื่อง
	ActionStartReportIssue = "start_report_issue"  // เริ่มแจ้งปัญหา
	ActionInputIssueDesc   = "input_issue_desc"    // รอพิมพ์รายละเอียด
	ActionSubmitIssue      = "submit_issue"        // บันทึกปัญหา
)

// Validation constants
const (
	MinSerialLength = 3
	MaxInputLength  = 100
)
