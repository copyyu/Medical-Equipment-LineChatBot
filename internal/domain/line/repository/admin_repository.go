package repository

import (
	"context"
	"medical-webhook/internal/domain/line/entity"

	"github.com/google/uuid"
)

type AdminRepository interface {
	Create(ctx context.Context, admin *entity.Admin) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Admin, error)
	GetByUsername(ctx context.Context, username string) (*entity.Admin, error)
	GetByEmail(ctx context.Context, email string) (*entity.Admin, error)
	Update(ctx context.Context, admin *entity.Admin) error
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*entity.Admin, error)
}
