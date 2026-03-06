package service

import (
	"fmt"
	"medical-webhook/internal/application/dto"
	"time"
)

type NotificationService struct{}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (s *NotificationService) FormatJuneAlert(alerts []dto.EquipmentReplacementAlertDTO) string {
	return s.formatAlert(alerts, "มิถุนายน (6 เดือนก่อนหมดอายุ)")
}

func (s *NotificationService) FormatAugustAlert(alerts []dto.EquipmentReplacementAlertDTO) string {
	return s.formatAlert(alerts, "สิงหาคม (1 ปีก่อนหมดอายุ)")
}

func (s *NotificationService) formatAlert(alerts []dto.EquipmentReplacementAlertDTO, roundName string) string {
	message := "\n🔔 แจ้งเตือนอุปกรณ์ที่ใกล้ครบกำหนดเปลี่ยน\n"
	message += fmt.Sprintf("📅 รอบแจ้งเตือน: %s\n", roundName)
	message += fmt.Sprintf("📆 วันที่: %s\n", time.Now().Format("02/01/2006 15:04"))
	message += "━━━━━━━━━━━━━━━━━━━━━━\n\n"

	if len(alerts) == 0 {
		message += "ไม่มีอุปกรณ์ที่ต้องแจ้งเตือนในรอบนี้"
		return message
	}

	urgent := s.filterByUrgency(alerts, 0, 3)
	warning := s.filterByUrgency(alerts, 3, 6)
	info := s.filterByUrgency(alerts, 6, 999)

	var displayAlerts []dto.EquipmentReplacementAlertDTO
	displayAlerts = append(displayAlerts, urgent...)
	displayAlerts = append(displayAlerts, warning...)
	displayAlerts = append(displayAlerts, info...)

	limit := 5
	if len(displayAlerts) < limit {
		limit = len(displayAlerts)
	}

	message += fmt.Sprintf("⚠ พบอุปกรณ์ใกล้หมดอายุจำนวน %d รายการ\n", len(alerts))
	message += "ตัวอย่างคร่าวๆ:\n\n"

	for i := 0; i < limit; i++ {
		alert := displayAlerts[i]
		message += fmt.Sprintf("%d. %s\n", i+1, alert.IDCode)
		message += fmt.Sprintf("   📦 %s - %s\n", alert.BrandName, alert.ModelName)
		message += fmt.Sprintf("   🏥 แผนก: %s\n", alert.DepartmentName)
		message += fmt.Sprintf("   ⏰ เหลืออีก %d เดือน\n\n", alert.MonthsRemaining)
	}

	message += "━━━━━━━━━━━━━━━━━━━━━━\n"
	if len(alerts) > limit {
		message += fmt.Sprintf("... และอื่นๆ อีก %d รายการ\n", len(alerts)-limit)
	}

	return message
}

func (s *NotificationService) filterByUrgency(alerts []dto.EquipmentReplacementAlertDTO, minMonths, maxMonths int) []dto.EquipmentReplacementAlertDTO {
	var filtered []dto.EquipmentReplacementAlertDTO
	for _, alert := range alerts {
		if alert.MonthsRemaining >= minMonths && alert.MonthsRemaining < maxMonths {
			filtered = append(filtered, alert)
		}
	}
	return filtered
}
