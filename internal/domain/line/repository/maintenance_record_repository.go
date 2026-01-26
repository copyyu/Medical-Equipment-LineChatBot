package repository

import (
	"context"
	"medical-webhook/internal/domain/line/entity"
)

// MaintenanceRecordRepository defines interface for maintenance record database operations
// This interface combines CRUD operations with aggregate query operations for dashboard
type MaintenanceRecordRepository interface {
	// CRUD Operations
	Create(ctx context.Context, record *entity.MaintenanceRecord) error
	FindByID(ctx context.Context, id uint) (*entity.MaintenanceRecord, error)
	FindByEquipmentID(ctx context.Context, equipmentID uint) ([]entity.MaintenanceRecord, error)
	FindByEquipmentIDAndType(ctx context.Context, equipmentID uint, maintenanceType entity.MaintenanceType) ([]entity.MaintenanceRecord, error)
	Update(ctx context.Context, record *entity.MaintenanceRecord) error
	Delete(ctx context.Context, id uint) error

	// Aggregate Query Operations (for Dashboard)
	Count(ctx context.Context) (int64, error)
	CountByType(ctx context.Context) (map[string]int64, error)
	GetRecent(ctx context.Context, limit int) ([]entity.MaintenanceRecord, error)
	GetTotalCostByEquipmentID(ctx context.Context, equipmentID uint) (float64, error)
	GetCMCountByEquipmentID(ctx context.Context, equipmentID uint) (int64, error)
}
