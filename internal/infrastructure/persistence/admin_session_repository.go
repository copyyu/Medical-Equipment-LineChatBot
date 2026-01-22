package persistence

import (
	"context"
	"errors"
	"medical-webhook/internal/domain/line/entity"
	"medical-webhook/internal/infrastructure/database"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminSessionRepository struct {
	db *gorm.DB
}

func NewAdminSessionRepository() *AdminSessionRepository {
	return &AdminSessionRepository{
		db: database.DB,
	}
}

func (r *AdminSessionRepository) Create(ctx context.Context, session *entity.AdminSession) error {
	session.ID = uuid.New()
	session.CreatedAt = time.Now()

	return r.db.WithContext(ctx).Create(session).Error
}

func (r *AdminSessionRepository) GetByToken(ctx context.Context, token string) (*entity.AdminSession, error) {
	var session entity.AdminSession
	if err := r.db.WithContext(ctx).Where("token = ? AND expires_at > ?", token, time.Now()).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &session, nil
}

func (r *AdminSessionRepository) GetByAdminID(ctx context.Context, adminID uuid.UUID) ([]*entity.AdminSession, error) {
	var sessions []*entity.AdminSession
	if err := r.db.WithContext(ctx).Where("admin_id = ?", adminID).Find(&sessions).Error; err != nil {
		return nil, err
	}

	return sessions, nil
}

func (r *AdminSessionRepository) DeleteByToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Where("token = ?", token).Delete(&entity.AdminSession{}).Error
}

func (r *AdminSessionRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&entity.AdminSession{}).Error
}

func (r *AdminSessionRepository) DeleteByAdminID(ctx context.Context, adminID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("admin_id = ?", adminID).Delete(&entity.AdminSession{}).Error
}
