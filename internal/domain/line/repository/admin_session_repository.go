package repository

import (
	"context"
	"medical-webhook/internal/domain/line/entity"

	"github.com/google/uuid"
)

type AdminSessionRepository interface {
	Create(ctx context.Context, session *entity.AdminSession) error
	GetByToken(ctx context.Context, token string) (*entity.AdminSession, error)
	GetByAdminID(ctx context.Context, adminID uuid.UUID) ([]*entity.AdminSession, error)
	DeleteByToken(ctx context.Context, token string) error
	DeleteExpired(ctx context.Context) error
	DeleteByAdminID(ctx context.Context, adminID uuid.UUID) error
}
