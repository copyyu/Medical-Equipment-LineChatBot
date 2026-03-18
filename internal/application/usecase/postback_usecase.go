package usecase

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

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

	case constants.ActionShowActionMenu:
		if serial != "" {
			return uc.lineRepo.ReplyFlexMessage(replyToken, "เลือกการดำเนินการ", templates.GetActionMenuFlex(serial))
		}
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgSelectMenu)

	case constants.ActionViewEquipInfo:
		if serial != "" {
			return uc.lineRepo.ReplyFlexMessage(replyToken, "ข้อมูลเครื่องมือ", templates.GetEquipmentOptionsFlex(serial))
		}
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgSelectMenu)

	case constants.ActionStartReportIssue:
		if serial != "" {
			categories, err := uc.ticketUseCase.GetTicketCategories(context.Background())
			if err != nil {
				log.Printf("❌ Failed to get categories: %v", err)
				return uc.lineRepo.ReplyFlexMessage(replyToken, "แจ้งปัญหา", templates.GetIssueInputFlex(serial, 0))
			}
			return uc.lineRepo.ReplyFlexMessage(replyToken, "เลือกหมวดหมู่", templates.GetCategorySelectionFlex(serial, categories))
		}
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgSelectMenu)

	case constants.ActionConfirmCategory:
		if serial != "" {
			categoryIDStr := params.Get("category_id")
			categoryID, _ := strconv.ParseUint(categoryIDStr, 10, 32)
			return uc.lineRepo.ReplyFlexMessage(replyToken, "แจ้งปัญหา", templates.GetIssueInputFlex(serial, uint(categoryID)))
		}
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgSelectMenu)

	case constants.ActionInputIssueDesc:
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
		departments, err := uc.departmentRepo.FindAll(ctx)
		if err != nil {
			log.Printf("❌ FindAll departments error: %v", err)
			return uc.lineRepo.ReplyMessage(replyToken, "❌ ไม่สามารถดึงข้อมูลแผนกได้ กรุณาลองใหม่ค่ะ")
		}
		if len(departments) == 0 {
			return uc.lineRepo.ReplyMessage(replyToken, "⚠️ ไม่พบข้อมูลแผนกในระบบค่ะ")
		}
		if userID != "" {
			uc.sessionStore.Set(userID, &session.OCRSession{Mode: session.ModeSelectDeptForExpiry})
		}
		return uc.lineRepo.ReplyFlexMessage(replyToken, "เลือกแผนก", templates.GetDepartmentSelectionWithInputFlex(departments))

	case constants.ActionViewEquipExpiryByDept:
		if userID != "" {
			uc.sessionStore.Delete(userID)
		}
		ctx := context.Background()
		departmentIDStr := params.Get("department_id")
		departmentID, err := strconv.ParseUint(departmentIDStr, 10, 32)
		if err != nil {
			log.Printf("❌ Invalid department_id: %v", err)
			return uc.lineRepo.ReplyMessage(replyToken, "❌ รหัสแผนกไม่ถูกต้องค่ะ")
		}

		dept, err := uc.departmentRepo.FindByID(ctx, uint(departmentID))
		if err != nil || dept == nil {
			log.Printf("❌ Department not found: %v", err)
			return uc.lineRepo.ReplyMessage(replyToken, "❌ ไม่พบแผนกที่เลือกค่ะ")
		}

		thisYear := time.Now().Year()
		nextYear := thisYear + 1
		deptIDPtr := func(id uint) *uint { return &id }(uint(departmentID))

		thisYearItems, err := uc.equipmentRepo.FindByReplacementYear(ctx, thisYear, deptIDPtr)
		if err != nil {
			return uc.lineRepo.ReplyMessage(replyToken, "❌ ไม่สามารถดึงข้อมูลได้ กรุณาลองใหม่ค่ะ")
		}
		nextYearItems, err := uc.equipmentRepo.FindByReplacementYear(ctx, nextYear, deptIDPtr)
		if err != nil {
			return uc.lineRepo.ReplyMessage(replyToken, "❌ ไม่สามารถดึงข้อมูลได้ กรุณาลองใหม่ค่ะ")
		}
		if len(thisYearItems) == 0 && len(nextYearItems) == 0 {
			return uc.lineRepo.ReplyMessage(replyToken, fmt.Sprintf("✅ ไม่มีเครื่องมือที่หมดอายุหรือใกล้หมดอายุในแผนก %s ค่ะ", dept.Name))
		}
		return uc.lineRepo.ReplyFlexMessage(replyToken,
			fmt.Sprintf("เครื่องใกล้หมดอายุ - %s", dept.Name),
			templates.GetExpiryYearFilterFlex(dept.Name, uint(departmentID), len(thisYearItems), len(nextYearItems)))

	case constants.ActionViewExpiryFiltered:
		ctx := context.Background()
		departmentIDStr := params.Get("department_id")
		departmentID, err := strconv.ParseUint(departmentIDStr, 10, 32)
		if err != nil {
			log.Printf("❌ Invalid department_id: %v", err)
			return uc.lineRepo.ReplyMessage(replyToken, "❌ รหัสแผนกไม่ถูกต้องค่ะ")
		}
		filter := params.Get("filter")

		dept, err := uc.departmentRepo.FindByID(ctx, uint(departmentID))
		if err != nil || dept == nil {
			return uc.lineRepo.ReplyMessage(replyToken, "❌ ไม่พบแผนกที่เลือกค่ะ")
		}

		thisYear := time.Now().Year()
		nextYear := thisYear + 1
		deptIDPtr := func(id uint) *uint { return &id }(uint(departmentID))

		thisYearItems, err := uc.equipmentRepo.FindByReplacementYear(ctx, thisYear, deptIDPtr)
		if err != nil {
			return uc.lineRepo.ReplyMessage(replyToken, "❌ ไม่สามารถดึงข้อมูลได้ กรุณาลองใหม่ค่ะ")
		}
		nextYearItems, err := uc.equipmentRepo.FindByReplacementYear(ctx, nextYear, deptIDPtr)
		if err != nil {
			return uc.lineRepo.ReplyMessage(replyToken, "❌ ไม่สามารถดึงข้อมูลได้ กรุณาลองใหม่ค่ะ")
		}

		if len(thisYearItems) == 0 && len(nextYearItems) == 0 {
			return uc.lineRepo.ReplyMessage(replyToken, fmt.Sprintf("✅ ไม่มีเครื่องมือที่หมดอายุหรือใกล้หมดอายุในแผนก %s ค่ะ", dept.Name))
		}

		switch filter {
		case "this_year":
			return uc.lineRepo.ReplyFlexMessage(replyToken,
				fmt.Sprintf("หมดอายุปี %d - %s", thisYear, dept.Name),
				templates.GetEquipmentExpiryFilteredFlex(thisYearItems,
					fmt.Sprintf("🔴 หมดอายุภายในปี %d", thisYear), templates.ColorDanger,
					dept.Name, uint(departmentID), uc.baseURL, "this_year"))
		case "next_year":
			return uc.lineRepo.ReplyFlexMessage(replyToken,
				fmt.Sprintf("หมดอายุปี %d - %s", nextYear, dept.Name),
				templates.GetEquipmentExpiryFilteredFlex(nextYearItems,
					fmt.Sprintf("🟡 หมดอายุปี %d", nextYear), templates.ColorWarning,
					dept.Name, uint(departmentID), uc.baseURL, "next_year"))
		default:
			return uc.lineRepo.ReplyFlexMessage(replyToken,
				fmt.Sprintf("เครื่องมือใกล้หมดอายุ - %s", dept.Name),
				templates.GetEquipmentExpiryByDeptFlex(thisYearItems, nextYearItems, dept.Name, uint(departmentID), uc.baseURL))
		}

	case constants.ActionNavDeptPage:
		ctx := context.Background()
		pageStr := params.Get("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}

		departments, err := uc.departmentRepo.FindAll(ctx)
		if err != nil {
			log.Printf("❌ FindAll departments error: %v", err)
			return uc.lineRepo.ReplyMessage(replyToken, "❌ ไม่สามารถดึงข้อมูลแผนกได้ กรุณาลองใหม่ค่ะ")
		}
		if len(departments) == 0 {
			return uc.lineRepo.ReplyMessage(replyToken, "⚠️ ไม่พบข้อมูลแผนกในระบบค่ะ")
		}

		return uc.lineRepo.ReplyFlexMessage(replyToken, "เลือกแผนก", templates.GetDepartmentSelectionPageFlex(departments, page))

	default:
		log.Printf("⚠️ Unhandled postback action: %s", action)
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgSelectMenu)
	}
}

func (uc *MessageUseCase) handleSubmitIssue(event webhook.PostbackEvent, replyToken, serial string, params url.Values) error {
	if serial == "" {
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgSelectMenu)
	}

	desc := params.Get("desc")
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

	tickets, err := uc.ticketUseCase.GetTicketsByEquipmentID(equipment.ID)
	if err != nil {
		return uc.lineRepo.ReplyMessage(replyToken, constants.MsgRepairHistoryFail)
	}

	// Convert to map format for template
	recordMaps := make([]map[string]interface{}, len(tickets))
	for i, t := range tickets {
		description := ""
		if t.Description != nil {
			description = *t.Description
		}
		if description == "" {
			description = "แจ้งซ่อม"
		}

		typeStr := "อื่น ๆ"
		if t.Category.ID != 0 {
			typeStr = t.Category.Name
		}

		recordMaps[i] = map[string]interface{}{
			"date":        t.ReportedAt.Format("2006-01-02"),
			"type":        typeStr,
			"description": fmt.Sprintf("[%s] %s", t.TicketNo, description),
			"cost":        "-",
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

	// Compute useful percent dynamically
	usefulPercent := 0.0
	if equipment.LifeExpectancy > 0 {
		usefulPercent = (equipment.EquipmentAge / equipment.LifeExpectancy) * 100
		if usefulPercent > 100 {
			usefulPercent = 100
		}
	}

	data := map[string]interface{}{
		"equipment_age":    equipment.EquipmentAge,
		"life_expectancy":  equipment.LifeExpectancy,
		"remain_life":      equipment.RemainLife,
		"useful_percent":   usefulPercent,
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
