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
		message += "✅ ไม่มีอุปกรณ์ที่ต้องแจ้งเตือนในรอบนี้"
		return message
	}

	urgent := s.filterByUrgency(alerts, 0, 3)
	warning := s.filterByUrgency(alerts, 3, 6)
	info := s.filterByUrgency(alerts, 6, 999)

	if len(urgent) > 0 {
		message += "🔴 เร่งด่วน! (เหลือ ≤ 3 เดือน)\n"
		message += "━━━━━━━━━━━━━━━━━━━━━━\n"
		message += s.formatAlertList(urgent)
		message += "\n"
	}

	if len(warning) > 0 {
		message += "🟡 ควรเตรียมการ (3-6 เดือน)\n"
		message += "━━━━━━━━━━━━━━━━━━━━━━\n"
		message += s.formatAlertList(warning)
		message += "\n"
	}

	if len(info) > 0 {
		message += "ℹ️ แจ้งให้ทราบ (> 6 เดือน)\n"
		message += "━━━━━━━━━━━━━━━━━━━━━━\n"
		message += s.formatAlertList(info)
		message += "\n"
	}

	message += "━━━━━━━━━━━━━━━━━━━━━━\n"
	message += fmt.Sprintf("📊 รวมทั้งหมด: %d รายการ", len(alerts))

	return message
}

func (s *NotificationService) formatAlertList(alerts []dto.EquipmentReplacementAlertDTO) string {
	var message string
	for i, alert := range alerts {
		message += fmt.Sprintf("%d. 📦 %s\n", i+1, alert.IDCode)
		message += fmt.Sprintf("   %s - %s\n", alert.BrandName, alert.ModelName)
		message += fmt.Sprintf("   แผนก: %s\n", alert.DepartmentName)
		message += fmt.Sprintf("   ⏰ เหลืออีก %d เดือน\n", alert.MonthsRemaining)

		if i < len(alerts)-1 {
			message += "\n"
		}
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
