package repository

import (
	"context"
	"medical-webhook/internal/application/dto"
	"medical-webhook/internal/domain/line/entity"
)

// NotificationRepository defines interface for notification operations
type NotificationRepository interface {
	// Settings
	GetSettings(ctx context.Context) (*entity.NotificationSetting, error)
	UpdateSettings(ctx context.Context, settings *entity.NotificationSetting) error
	CountAllEquipments(ctx context.Context) (int, error)

	// Logs
	CreateLog(ctx context.Context, log *entity.NotificationLog) error
	GetLastNotification(ctx context.Context) (*entity.NotificationLog, error)
	GetLogsByMonth(ctx context.Context, year int, month int) ([]entity.NotificationLog, error)

	// Equipment Alerts
	GetEquipmentsForJuneAlert(ctx context.Context) ([]dto.EquipmentReplacementAlertDTO, error)
	GetEquipmentsForAugustAlert(ctx context.Context) ([]dto.EquipmentReplacementAlertDTO, error)
	GetEquipmentsForTestAlert(ctx context.Context, targetYear int, notifyRound string) ([]dto.EquipmentReplacementAlertDTO, error)
}
