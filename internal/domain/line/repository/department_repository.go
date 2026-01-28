package repository

import (
	"context"
	"medical-webhook/internal/domain/line/entity"
)

// DepartmentRepository defines interface for department database operations
type DepartmentRepository interface {
	Create(ctx context.Context, department *entity.Department) error
	FindByID(ctx context.Context, id uint) (*entity.Department, error)
	FindByName(ctx context.Context, name string) (*entity.Department, error)
	FindAll(ctx context.Context) ([]entity.Department, error)
	Update(ctx context.Context, department *entity.Department) error
	Delete(ctx context.Context, id uint) error
	FindOrCreate(ctx context.Context, name string) (*entity.Department, error)
}
