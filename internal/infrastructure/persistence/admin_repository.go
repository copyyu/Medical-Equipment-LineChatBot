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

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository() *AdminRepository {
	return &AdminRepository{
		db: database.DB,
	}
}

func (r *AdminRepository) Create(ctx context.Context, admin *entity.Admin) error {
	admin.ID = uuid.New()
	admin.CreatedAt = time.Now()
	admin.UpdatedAt = time.Now()

	return r.db.WithContext(ctx).Create(admin).Error
}

func (r *AdminRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Admin, error) {
	var admin entity.Admin
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &admin, nil
}

func (r *AdminRepository) GetByUsername(ctx context.Context, username string) (*entity.Admin, error) {
	var admin entity.Admin
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &admin, nil
}

func (r *AdminRepository) GetByEmail(ctx context.Context, email string) (*entity.Admin, error) {
	var admin entity.Admin
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&admin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &admin, nil
}

func (r *AdminRepository) Update(ctx context.Context, admin *entity.Admin) error {
	admin.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Model(&entity.Admin{}).Where("id = ?", admin.ID).Updates(admin).Error
}

func (r *AdminRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&entity.Admin{}).Where("id = ?", id).Update("last_login_at", now).Error
}

func (r *AdminRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Admin{}, "id = ?", id).Error
}

func (r *AdminRepository) List(ctx context.Context, limit, offset int) ([]*entity.Admin, error) {
	var admins []*entity.Admin
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&admins).Error; err != nil {
		return nil, err
	}

	return admins, nil
}
