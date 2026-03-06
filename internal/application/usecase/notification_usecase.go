package usecase

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"medical-webhook/internal/application/dto"
	"medical-webhook/internal/application/service"
	"medical-webhook/internal/domain/line/entity"
	"medical-webhook/internal/domain/line/repository"
	notificationRepo "medical-webhook/internal/domain/line/repository"

	excelize "github.com/xuri/excelize/v2"
)

type NotificationUseCase struct {
	notificationRepo    notificationRepo.NotificationRepository
	notificationService *service.NotificationService
	lineRepo            repository.LineRepository
	equipmentRepo       repository.EquipmentRepository
}

func NewNotificationUseCase(
	notificationRepo notificationRepo.NotificationRepository,
	notificationService *service.NotificationService,
	lineRepo repository.LineRepository,
	equipmentRepo repository.EquipmentRepository,
) *NotificationUseCase {
	return &NotificationUseCase{
		notificationRepo:    notificationRepo,
		notificationService: notificationService,
		lineRepo:            lineRepo,
		equipmentRepo:       equipmentRepo,
	}
}

func (uc *NotificationUseCase) SendJuneAlerts(ctx context.Context) error {
	return uc.sendAlerts(ctx, "JUNE")
}

func (uc *NotificationUseCase) SendAugustAlerts(ctx context.Context) error {
	return uc.sendAlerts(ctx, "AUGUST")
}

func (uc *NotificationUseCase) sendAlerts(ctx context.Context, notifyRound string) error {
	log.Printf("🔔 Starting %s notification round...", notifyRound)

	settings, err := uc.notificationRepo.GetSettings(ctx)
	if err != nil {
		return fmt.Errorf("failed to get settings: %w", err)
	}

	if !settings.IsEnabled {
		log.Println("⚠️ Notification is disabled")
		return nil
	}

	var alerts []dto.EquipmentReplacementAlertDTO
	if notifyRound == "JUNE" {
		alerts, err = uc.notificationRepo.GetEquipmentsForJuneAlert(ctx)
	} else {
		alerts, err = uc.notificationRepo.GetEquipmentsForAugustAlert(ctx)
	}

	if err != nil {
		return fmt.Errorf("failed to get equipment alerts: %w", err)
	}

	if len(alerts) == 0 {
		log.Printf("ℹ️ No equipment needs notification for %s round", notifyRound)
		return nil
	}

	var message string
	if notifyRound == "JUNE" {
		message = uc.notificationService.FormatJuneAlert(alerts)
	} else {
		message = uc.notificationService.FormatAugustAlert(alerts)
	}

	// ✅ Broadcast ไปยังทุกคนที่เพิ่มเพื่อน
	err = uc.lineRepo.BroadcastMessage(message)

	// บันทึก log
	now := time.Now()
	for _, alert := range alerts {
		status := entity.NotificationStatusSent
		var errorMsg *string
		if err != nil {
			status = entity.NotificationStatusFailed
			msg := err.Error()
			errorMsg = &msg
		}

		notifLog := &entity.NotificationLog{
			EquipmentID: alert.EquipmentID,
			NotifyRound: notifyRound,
			Message:     message,
			Status:      status,
			SentAt:      now,
			ErrorMsg:    errorMsg,
		}
		uc.notificationRepo.CreateLog(ctx, notifLog)
	}

	if err != nil {
		return fmt.Errorf("failed to send broadcast: %w", err)
	}

	log.Printf("✅ Broadcast sent to all Bot friends for %s round", notifyRound)
	return nil
}

func (uc *NotificationUseCase) GetNotificationSummary(ctx context.Context) (*dto.NotificationSummaryDTO, error) {
	// ✅ นับเครื่องมือทั้งหมด
	totalEquipments, err := uc.notificationRepo.CountAllEquipments(ctx)
	if err != nil {
		log.Printf("Error counting total equipments: %v", err)
		totalEquipments = 0 // fallback
	}

	// นับเครื่องที่ต้อง alert แยกตามรอบ
	juneAlerts, _ := uc.notificationRepo.GetEquipmentsForJuneAlert(ctx)
	augustAlerts, _ := uc.notificationRepo.GetEquipmentsForAugustAlert(ctx)

	summary := &dto.NotificationSummaryDTO{
		TotalEquipments: totalEquipments,   // ✅ ใช้จำนวนทั้งหมด
		JuneAlerts:      len(juneAlerts),   // จำนวนที่ต้อง alert เดือน 6
		AugustAlerts:    len(augustAlerts), // จำนวนที่ต้อง alert เดือน 8
	}

	// ดึง notification ล่าสุด
	lastLog, err := uc.notificationRepo.GetLastNotification(ctx)
	if err == nil && lastLog != nil {
		summary.LastNotification = &lastLog.SentAt
	}

	return summary, nil
}

func (uc *NotificationUseCase) UpdateSettings(ctx context.Context, settingsDTO *dto.NotificationSettingDTO) error {
	settings, err := uc.notificationRepo.GetSettings(ctx)
	if err != nil {
		return err
	}

	settings.IsEnabled = settingsDTO.IsEnabled

	return uc.notificationRepo.UpdateSettings(ctx, settings)
}

// BuildExpiryExcel สร้างไฟล์ Excel เครื่องใกล้หมดอายุ กรองตาม filter (this_year / next_year / all)
// ถ้า departmentID == nil จะดึงทุกแผนก
func (uc *NotificationUseCase) BuildExpiryExcel(ctx context.Context, departmentID *uint, filter string) ([]byte, string, error) {

	now := time.Now()
	thisYear := now.Year()
	nextYear := thisYear + 1

	// ดึงเครื่องตาม ReplacementYear ตรงๆ — ไม่ใช้ expired/near-expiry
	thisYearItems, err := uc.equipmentRepo.FindByReplacementYear(ctx, thisYear, departmentID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get this year equipments: %w", err)
	}

	nextYearItems, err := uc.equipmentRepo.FindByReplacementYear(ctx, nextYear, departmentID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get next year equipments: %w", err)
	}

	// สร้างไฟล์ Excel
	f := excelize.NewFile()
	defer f.Close()

	sheet := "เครื่องใกล้หมดอายุ"
	f.SetSheetName("Sheet1", sheet)

	// Header Style (dark green)
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF", Size: 11},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"1B5E20"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "FFFFFF", Style: 1},
			{Type: "right", Color: "FFFFFF", Style: 1},
			{Type: "top", Color: "FFFFFF", Style: 1},
			{Type: "bottom", Color: "FFFFFF", Style: 1},
		},
	})

	// Year section header style
	thisYearStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF", Size: 11},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"E53935"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
	})
	nextYearStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF", Size: 11},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"FF9800"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
	})

	// Normal alternating row styles
	normalStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "CCCCCC", Style: 1},
			{Type: "right", Color: "CCCCCC", Style: 1},
			{Type: "top", Color: "CCCCCC", Style: 1},
			{Type: "bottom", Color: "CCCCCC", Style: 1},
		},
	})
	altStyle, _ := f.NewStyle(&excelize.Style{
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"F5F5F5"}, Pattern: 1},
		Alignment: &excelize.Alignment{Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "CCCCCC", Style: 1},
			{Type: "right", Color: "CCCCCC", Style: 1},
			{Type: "top", Color: "CCCCCC", Style: 1},
			{Type: "bottom", Color: "CCCCCC", Style: 1},
		},
	})
	runningNoStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "CCCCCC", Style: 1},
			{Type: "right", Color: "CCCCCC", Style: 1},
			{Type: "top", Color: "CCCCCC", Style: 1},
			{Type: "bottom", Color: "CCCCCC", Style: 1},
		},
	})

	// Columns: A=ลำดับ B=รหัสเครื่อง C=Serial D=ยี่ห้อ E=รุ่น F=แผนก G=ปีที่เปลี่ยน
	f.SetColWidth(sheet, "A", "A", 7)
	f.SetColWidth(sheet, "B", "B", 16)
	f.SetColWidth(sheet, "C", "C", 18)
	f.SetColWidth(sheet, "D", "D", 24)
	f.SetColWidth(sheet, "E", "E", 22)
	f.SetColWidth(sheet, "F", "F", 14)
	f.SetColWidth(sheet, "G", "G", 14)

	// Title row
	f.MergeCell(sheet, "A1", "G1")
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14, Color: "1B5E20"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	f.SetCellValue(sheet, "A1", fmt.Sprintf("รายงานเครื่องมือใกล้หมดอายุ ปี %d-%d  —  ข้อมูล ณ %s",
		thisYear, nextYear, now.Format("02/01/2006 15:04")))
	f.SetCellStyle(sheet, "A1", "G1", titleStyle)
	f.SetRowHeight(sheet, 1, 28)

	// Column headers (row 2)
	colHeaders := []string{"ลำดับ", "รหัสเครื่อง", "Serial No", "ยี่ห้อ", "รุ่น", "แผนก", "ปีที่เปลี่ยน"}
	for i, h := range colHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 2)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}
	f.SetRowHeight(sheet, 2, 22)

	row := 3
	totalCount := 0

	// Helper: write one equipment row with running number
	writeEquipRow := func(e entity.Equipment, runNo int) {
		serialNo := ""
		if e.SerialNo != nil {
			serialNo = *e.SerialNo
		}
		modelName := ""
		if e.Model.ModelName != "" {
			modelName = e.Model.ModelName
		}
		brandName := ""
		if e.Model.Brand.Name != "" {
			brandName = e.Model.Brand.Name
		}
		deptName := ""
		if e.Department.Name != "" {
			deptName = e.Department.Name
		}
		replYear := 0
		if e.ReplacementYear != nil {
			replYear = *e.ReplacementYear
		}

		rowStyle := normalStyle
		if runNo%2 == 0 {
			rowStyle = altStyle
		}

		// Running number (col A)
		cellA, _ := excelize.CoordinatesToCellName(1, row)
		f.SetCellValue(sheet, cellA, runNo)
		f.SetCellStyle(sheet, cellA, cellA, runningNoStyle)

		// Data columns (B-G)
		values := []interface{}{e.IDCode, serialNo, brandName, modelName, deptName, replYear}
		for col, val := range values {
			cell, _ := excelize.CoordinatesToCellName(col+2, row)
			f.SetCellValue(sheet, cell, val)
			f.SetCellStyle(sheet, cell, cell, rowStyle)
		}
		f.SetRowHeight(sheet, row, 18)
		row++
		totalCount++
	}

	// Helper: write year section header
	writeYearHeader := func(label string, count int, style int) {
		sc := fmt.Sprintf("A%d", row)
		ec := fmt.Sprintf("G%d", row)
		f.MergeCell(sheet, sc, ec)
		f.SetCellValue(sheet, sc, fmt.Sprintf("  %s  (%d รายการ)", label, count))
		f.SetCellStyle(sheet, sc, ec, style)
		f.SetRowHeight(sheet, row, 20)
		row++
	}

	showThisYear := filter == "this_year" || filter == "all" || filter == ""
	showNextYear := filter == "next_year" || filter == "all" || filter == ""

	// Section 1: ปีนี้
	if showThisYear {
		writeYearHeader(fmt.Sprintf("🔴 ปีนี้ — ต้องเปลี่ยนปี %d", thisYear), len(thisYearItems), thisYearStyle)
		if len(thisYearItems) == 0 {
			mc := fmt.Sprintf("A%d", row)
			me := fmt.Sprintf("G%d", row)
			f.MergeCell(sheet, mc, me)
			f.SetCellValue(sheet, mc, "✅ ไม่มีเครื่องที่ต้องเปลี่ยนในปีนี้")
			row++
		}
		for i, e := range thisYearItems {
			writeEquipRow(e, i+1)
		}
	}

	// Section 2: ปีหน้า
	if showNextYear {
		writeYearHeader(fmt.Sprintf("🟡 ปีหน้า — ต้องเปลี่ยนปี %d", nextYear), len(nextYearItems), nextYearStyle)
		if len(nextYearItems) == 0 {
			mc := fmt.Sprintf("A%d", row)
			me := fmt.Sprintf("G%d", row)
			f.MergeCell(sheet, mc, me)
			f.SetCellValue(sheet, mc, "✅ ไม่มีเครื่องที่ต้องเปลี่ยนในปีหน้า")
			row++
		}
		for i, e := range nextYearItems {
			writeEquipRow(e, i+1)
		}
	}

	// Summary row
	sc := fmt.Sprintf("A%d", row)
	se := fmt.Sprintf("G%d", row)
	f.MergeCell(sheet, sc, se)
	summaryStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "1B5E20"},
		Alignment: &excelize.Alignment{Horizontal: "right", Vertical: "center"},
	})
	f.SetCellValue(sheet, sc, fmt.Sprintf("รวมทั้งหมด: %d รายการ  (ปี %d: %d | ปี %d: %d)",
		totalCount, thisYear, len(thisYearItems), nextYear, len(nextYearItems)))
	f.SetCellStyle(sheet, sc, se, summaryStyle)

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, "", fmt.Errorf("failed to write excel: %w", err)
	}

	// Filename reflects filter
	filterLabel := fmt.Sprintf("%d-%d", thisYear, nextYear)
	if filter == "this_year" {
		filterLabel = fmt.Sprintf("%d", thisYear)
	} else if filter == "next_year" {
		filterLabel = fmt.Sprintf("%d", nextYear)
	}
	filename := fmt.Sprintf("expiry_report_%s_%s.xlsx", filterLabel, now.Format("20060102"))
	return buf.Bytes(), filename, nil
}
