package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"medical-webhook/internal/application/dto"
	"medical-webhook/internal/domain/line/entity"
	"medical-webhook/internal/domain/line/repository"
	notificationRepo "medical-webhook/internal/domain/line/repository"
	"medical-webhook/internal/domain/line/service"
)

type NotificationUseCase struct {
	notificationRepo    notificationRepo.NotificationRepository
	notificationService *service.NotificationService
	lineRepo            repository.LineRepository
}

func NewNotificationUseCase(
	notificationRepo notificationRepo.NotificationRepository,
	notificationService *service.NotificationService,
	lineRepo repository.LineRepository,
) *NotificationUseCase {
	return &NotificationUseCase{
		notificationRepo:    notificationRepo,
		notificationService: notificationService,
		lineRepo:            lineRepo,
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

		log := &entity.NotificationLog{
			EquipmentID: alert.EquipmentID,
			NotifyRound: notifyRound,
			Message:     message,
			Status:      status,
			SentAt:      now,
			ErrorMsg:    errorMsg,
		}
		uc.notificationRepo.CreateLog(ctx, log)
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
