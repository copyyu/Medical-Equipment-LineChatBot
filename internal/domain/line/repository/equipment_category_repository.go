package repository

import (
	"context"
	"medical-webhook/internal/domain/line/entity"
)

// EquipmentCategoryRepository defines interface for equipment category database operations
type EquipmentCategoryRepository interface {
	Create(ctx context.Context, category *entity.EquipmentCategory) error
	FindByID(ctx context.Context, id uint) (*entity.EquipmentCategory, error)
	FindByName(ctx context.Context, name string) (*entity.EquipmentCategory, error)
	FindAll(ctx context.Context) ([]entity.EquipmentCategory, error)
	FindByECRIRisk(ctx context.Context, risk entity.ECRIRiskLevel) ([]entity.EquipmentCategory, error)
	Update(ctx context.Context, category *entity.EquipmentCategory) error
	Delete(ctx context.Context, id uint) error
}
