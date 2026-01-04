package repository

import (
	"context"
	"medical-webhook/internal/domain/line/entity"
)

// BrandRepository defines interface for brand database operations
type BrandRepository interface {
	Create(ctx context.Context, brand *entity.Brand) error
	FindByID(ctx context.Context, id uint) (*entity.Brand, error)
	FindByName(ctx context.Context, name string) (*entity.Brand, error)
	FindAll(ctx context.Context) ([]entity.Brand, error)
	Update(ctx context.Context, brand *entity.Brand) error
	Delete(ctx context.Context, id uint) error
}
