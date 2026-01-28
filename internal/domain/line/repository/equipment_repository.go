package repository

import (
	"context"
	"medical-webhook/internal/domain/line/entity"
)

type EquipmentRepository interface {
	// FindByIDCode finds equipment by id_code
	FindByIDCode(idCode string) (*entity.Equipment, error)
	// FindBySerialNo finds equipment by serial_no
	FindBySerialNo(serialNo string) (*entity.Equipment, error)
	// FindBySerialOrCode finds equipment by either serial_no or id_code
	FindBySerialOrCode(query string) (*entity.Equipment, error)
	// GetMaintenanceRecords gets maintenance records for equipment
	GetMaintenanceRecords(equipmentID uint) ([]entity.MaintenanceRecord, error)

	// CRUD methods for Excel import
	Create(ctx context.Context, equipment *entity.Equipment) error
	CreateOrUpdate(ctx context.Context, equipment *entity.Equipment) error
	Update(ctx context.Context, equipment *entity.Equipment) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*entity.Equipment, error)
	FindAll(ctx context.Context, limit, offset int) ([]entity.Equipment, error)

	// Aggregate Query Operations (for Dashboard)
	Count(ctx context.Context) (int64, error)
	CountNearExpiry(ctx context.Context) (int64, error)
	CountByStatus(ctx context.Context) (map[entity.AssetStatus]int64, error) // นับจำนวนอุปกรณ์แยกตาม Status
}
