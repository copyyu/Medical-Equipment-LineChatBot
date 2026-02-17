package usecase

import (
	"context"
	"log"
	"net/url"
	"strconv"

	"medical-webhook/internal/domain/constants"
	"medical-webhook/internal/domain/line/model"
	"medical-webhook/internal/infrastructure/line/templates"
	"medical-webhook/internal/infrastructure/session"

	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

// HandlePostbackEvent handles postback events from Flex Message buttons
func (uc *MessageUseCase) HandlePostbackEvent(event webhook.PostbackEvent) error {
	data := event.Postback.Data
	replyToken := event.ReplyToken
	log.Printf("📤 Processing postback: %s", data)

	// Parse postback data
	params, _ := url.ParseQuery(data)
	action := params.Get("action")
	serial := params.Get("serial")

	// Get user ID
	var userID string
	if source, ok := event.Source.(webhook.UserSource); ok {
		userID = source.UserId
	}

	switch action {
	case constants.ActionMainMenu:
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgSelectMenu)

	case constants.ActionRequestChange:
		return uc.lineRepo.ReplyFlexMessage(replyToken, "แจ้งเปลี่ยนเครื่อง", uc.messageService.GetEquipmentChangeFlex("https://www.google.com/"))

	case constants.ActionReportProblem:
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgReportProblem)

	case constants.ActionTrackStatus:
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgTrackStatus)

	case constants.ActionContactStaff:
		return uc.lineRepo.ReplyFlexMessage(replyToken, "ติดต่อเจ้าหน้าที่", uc.messageService.GetContactStaffFlex())

	case constants.ActionOCRConfirmYes:
		// User confirmed OCR result - show action menu (ดูข้อมูล/แจ้งปัญหา)
		if serial != "" {
			return uc.lineRepo.ReplyFlexMessage(replyToken, "เลือกการดำเนินการ", templates.GetActionMenuFlex(serial))
		}
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgSelectMenu)

	case constants.ActionOCRConfirmNo:
		if serial == "" {
			return uc.lineRepo.ReplyFlexMessage(
				replyToken,
				"ส่งรูปใหม่",
				templates.GetRetryPhotoFlex(),
			)
		}

		equipments, err := uc.equipmentRepo.FindSimilarSorted(serial, 5)
		if err != nil {
			log.Printf("❌ FindSimilarSorted error: %v", err)
			return uc.lineRepo.ReplyFlexMessage(
				replyToken,
				"ส่งรูปใหม่",
				templates.GetRetryPhotoFlex(),
			)
		}

		if len(equipments) == 0 {
			log.Printf("⚠️ No similar equipment for: %s", serial)
			return uc.lineRepo.ReplyFlexMessage(
				replyToken,
				"ไม่พบในฐานระบบ",
				templates.GetOCRNotFoundFlex(serial),
			)
		}

		log.Printf("✅ Found %d similar equipments (sorted) for: %s", len(equipments), serial)

		return uc.lineRepo.ReplyFlexMessage(
			replyToken,
			"พบข้อมูลใกล้เคียง",
			templates.GetSimilarEquipmentListFlex(serial, equipments),
		)

	case constants.ActionOCRSimilarSelect:
		// ผู้ใช้เลือกจากรายการใกล้เคียง → ถามยืนยันก่อน
		original := params.Get("original")
		if serial == "" {
			return uc.lineRepo.ReplyMessage(replyToken, "ไม่พบหมายเลขที่เลือก กรุณาลองใหม่")
		}

		log.Printf("📋 User selected similar equipment: %s (original OCR: %s)", serial, original)

		return uc.lineRepo.ReplyFlexMessage(
			replyToken,
			"ยืนยันเปลี่ยนหมายเลข",
			templates.GetSimilarConfirmFlex(serial, original),
		)

	case constants.ActionOCRRetake:
		// รีเซ็ตสถานะและให้ผู้ใช้ถ่ายรูปใหม่
		if userID != "" {
			uc.sessionStore.Set(userID, &session.OCRSession{Mode: session.ModeReportProblem})
			log.Printf("📸 User %s requested to retake photo", userID)
		}
		return uc.lineRepo.ReplyMessage(replyToken, "กรุณาถ่ายรูปบาร์โค้ดหรือหมายเลขอุปกรณ์ใหม่อีกครั้ง 📸")

	case constants.ActionViewRepairHist:
		return uc.handleViewRepairHistory(replyToken, serial)

	case constants.ActionViewLifecycle:
		return uc.handleViewLifecycle(replyToken, serial)

	case constants.ActionViewSpecs:
		return uc.handleViewSpecs(replyToken, serial)

	// New handlers for report issue flow
	case constants.ActionShowActionMenu:
		// Show action menu (ดูข้อมูล/แจ้งปัญหา)
		if serial != "" {
			return uc.lineRepo.ReplyFlexMessage(replyToken, "เลือกการดำเนินการ", templates.GetActionMenuFlex(serial))
		}
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgSelectMenu)

	case constants.ActionViewEquipInfo:
		// Go to equipment info menu (existing)
		if serial != "" {
			return uc.lineRepo.ReplyFlexMessage(replyToken, "ข้อมูลเครื่องมือ", templates.GetEquipmentOptionsFlex(serial))
		}
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgSelectMenu)

	case constants.ActionStartReportIssue:
		// Show category selection menu first
		if serial != "" {
			categories, err := uc.ticketUseCase.GetTicketCategories(context.Background())
			if err != nil {
				log.Printf("❌ Failed to get categories: %v", err)
				// Fallback: skip category selection and go to issue input with default category
				return uc.lineRepo.ReplyFlexMessage(replyToken, "แจ้งปัญหา", templates.GetIssueInputFlex(serial, 0))
			}
			return uc.lineRepo.ReplyFlexMessage(replyToken, "เลือกหมวดหมู่", templates.GetCategorySelectionFlex(serial, categories))
		}
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgSelectMenu)

	case constants.ActionConfirmCategory:
		// User selected a category, show issue input
		if serial != "" {
			categoryIDStr := params.Get("category_id")
			categoryID, _ := strconv.ParseUint(categoryIDStr, 10, 32)
			return uc.lineRepo.ReplyFlexMessage(replyToken, "แจ้งปัญหา", templates.GetIssueInputFlex(serial, uint(categoryID)))
		}
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgSelectMenu)

	case constants.ActionInputIssueDesc:
		// Set session mode to wait for issue description
		if serial != "" {
			categoryIDStr := params.Get("category_id")
			categoryID, _ := strconv.ParseUint(categoryIDStr, 10, 32)

			var userID string
			switch source := event.Source.(type) {
			case webhook.UserSource:
				userID = source.UserId
			case webhook.GroupSource:
				userID = source.UserId
			case webhook.RoomSource:
				userID = source.UserId
			}

			if userID != "" {
				uc.sessionStore.Set(userID, &session.OCRSession{
					Mode:         session.ModeInputIssueDesc,
					SerialNumber: serial,
					CategoryID:   uint(categoryID),
				})
			}
			return uc.lineRepo.ReplyMessage(replyToken, constants.MsgInputIssueDesc)
		}
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgSelectMenu)

	case constants.ActionSubmitIssue:
		return uc.handleSubmitIssue(event, replyToken, serial, params)

	case constants.ActionMyTickets:
		// Show user's tickets
		if userID != "" {
			tickets, err := uc.ticketUseCase.GetUserTickets(userID)
			if err != nil {
				log.Printf("❌ GetUserTickets error: %v", err)
				return uc.lineRepo.ReplyMessage(replyToken, "❌ ไม่สามารถดึงข้อมูลได้ กรุณาลองใหม่ค่ะ")
			}
			if len(tickets) == 0 {
				return uc.lineRepo.ReplyMessage(replyToken, "📋 คุณยังไม่มีรายการแจ้งปัญหาค่ะ")
			}
			return uc.lineRepo.ReplyFlexMessage(replyToken, "รายการแจ้งปัญหาของคุณ", templates.GetMyTicketsFlex(tickets))
		}
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgSelectMenu)

	case constants.ActionStartReportMode:
		if userID != "" {
			uc.sessionStore.Set(userID, &session.OCRSession{Mode: session.ModeReportProblem})
		}
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgReportProblem)

	case constants.ActionViewEquipExpiry:
		ctx := context.Background()
		expired, err := uc.equipmentRepo.FindExpired(ctx, 10)
		if err != nil {
			log.Printf("❌ FindExpired error: %v", err)
			return uc.lineRepo.ReplyMessage(replyToken, "❌ ไม่สามารถดึงข้อมูลได้ กรุณาลองใหม่ค่ะ")
		}
		nearExpiry, err := uc.equipmentRepo.FindNearExpiry(ctx, 10)
		if err != nil {
			log.Printf("❌ FindNearExpiry error: %v", err)
			return uc.lineRepo.ReplyMessage(replyToken, "❌ ไม่สามารถดึงข้อมูลได้ กรุณาลองใหม่ค่ะ")
		}
		if len(expired) == 0 && len(nearExpiry) == 0 {
			return uc.lineRepo.ReplyMessage(replyToken, "✅ ไม่มีเครื่องมือที่หมดอายุหรือใกล้หมดอายุในขณะนี้ค่ะ")
		}
		return uc.lineRepo.ReplyFlexMessage(replyToken, "เครื่องมือหมดอายุ/ใกล้หมดอายุ", templates.GetEquipmentExpiryFlex(expired, nearExpiry))

	default:
		log.Printf("⚠️ Unhandled postback action: %s", action)
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgSelectMenu)
	}
}

// handleSubmitIssue handles the submit issue action from postback
func (uc *MessageUseCase) handleSubmitIssue(event webhook.PostbackEvent, replyToken, serial string, params url.Values) error {
	if serial == "" {
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgSelectMenu)
	}

	desc := params.Get("desc") // empty for skip
	categoryIDStr := params.Get("category_id")
	categoryID, _ := strconv.ParseUint(categoryIDStr, 10, 32)
	userID := ""
	var groupID, sourceType string

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

	displayName := ""
	photoURL := ""

	var profile *model.UserProfile
	var err error

	switch sourceType {
	case "group":
		profile, err = uc.lineRepo.GetGroupMemberProfile(groupID, userID)
	case "room":
		profile, err = uc.lineRepo.GetRoomMemberProfile(groupID, userID)
	default:
		profile, err = uc.lineRepo.GetProfile(userID)
	}

	if err != nil {
		log.Printf("❌ Failed to get user profile: %v", err)
		displayName = "LINE User"
	} else if profile != nil {
		displayName = profile.DisplayName
		photoURL = profile.PictureURL
	} else {
		displayName = "LINE User"
	}

	ticket, err := uc.ticketUseCase.CreateTicketFromLINE(
		serial,
		desc,
		userID,
		displayName,
		photoURL,
		uint(categoryID),
	)
	if err != nil {
		// Check if it's a duplicate ticket error
		if err == ErrDuplicateTicket && ticket != nil {
			log.Printf("⚠️ Duplicate ticket found: %s", ticket.TicketNo)
			return uc.lineRepo.ReplyFlexMessage(replyToken, "พบรายการซ้ำ", templates.GetDuplicateTicketFlex(ticket.TicketNo, serial, ticket.Status.GetStatusText()))
		}
		log.Printf("❌ Failed to create ticket: %v", err)
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgIssueReportFailed)
	}
	return uc.lineRepo.ReplyFlexMessage(replyToken, "สร้าง Ticket สำเร็จ", templates.GetTicketCreatedFlex(ticket))
}

// handleViewRepairHistory sends repair history for equipment
func (uc *MessageUseCase) handleViewRepairHistory(replyToken, serial string) error {
	equipment, err := uc.equipmentRepo.FindBySerialOrCode(serial)
	if err != nil || equipment == nil {
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgEquipmentNotFound)
	}

	records, err := uc.equipmentRepo.GetMaintenanceRecords(equipment.ID)
	if err != nil {
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgRepairHistoryFail)
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
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgEquipmentNotFound)
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
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgEquipmentNotFound)
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
