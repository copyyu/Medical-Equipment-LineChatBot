package repository

import (
	"context"
	"medical-webhook/internal/domain/line/entity"
)

// MaintenanceRecordRepository defines interface for maintenance record database operations
type MaintenanceRecordRepository interface {
	Create(ctx context.Context, record *entity.MaintenanceRecord) error
	FindByID(ctx context.Context, id uint) (*entity.MaintenanceRecord, error)
	FindByEquipmentID(ctx context.Context, equipmentID uint) ([]entity.MaintenanceRecord, error)
	FindByEquipmentIDAndType(ctx context.Context, equipmentID uint, maintenanceType entity.MaintenanceType) ([]entity.MaintenanceRecord, error)
	GetTotalCostByEquipmentID(ctx context.Context, equipmentID uint) (float64, error)
	GetCMCountByEquipmentID(ctx context.Context, equipmentID uint) (int64, error)
	Update(ctx context.Context, record *entity.MaintenanceRecord) error
	Delete(ctx context.Context, id uint) error
}
