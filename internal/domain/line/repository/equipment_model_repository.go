package repository

import (
	"context"
	"medical-webhook/internal/domain/line/entity"
)

// EquipmentModelRepository defines interface for equipment model database operations
type EquipmentModelRepository interface {
	Create(ctx context.Context, model *entity.EquipmentModel) error
	FindByID(ctx context.Context, id uint) (*entity.EquipmentModel, error)
	FindByBrandCategoryModel(ctx context.Context, brandID, categoryID uint, modelName string) (*entity.EquipmentModel, error)
	FindByBrandID(ctx context.Context, brandID uint) ([]entity.EquipmentModel, error)
	FindByCategoryID(ctx context.Context, categoryID uint) ([]entity.EquipmentModel, error)
	FindAll(ctx context.Context) ([]entity.EquipmentModel, error)
	Update(ctx context.Context, model *entity.EquipmentModel) error
	Delete(ctx context.Context, id uint) error
}
