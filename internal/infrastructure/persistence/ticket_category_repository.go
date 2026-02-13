package persistence

import (
	"medical-webhook/internal/domain/line/entity"

	"gorm.io/gorm"
)

// TicketCategoryRepository implements repository.TicketCategoryRepository
type TicketCategoryRepository struct {
	db *gorm.DB
}

// NewTicketCategoryRepository creates a new ticket category repository
func NewTicketCategoryRepository(db *gorm.DB) *TicketCategoryRepository {
	return &TicketCategoryRepository{db: db}
}

func (r *TicketCategoryRepository) GetTicketCategories() ([]entity.TicketCategory, error) {
	var categories []entity.TicketCategory
	err := r.db.Where("is_active = ?", true).
		Order("sort_order ASC").
		Find(&categories).Error
	return categories, err
}

func (r *TicketCategoryRepository) FindCategoryByName(name string) (*entity.TicketCategory, error) {
	var category entity.TicketCategory
	err := r.db.Where("name = ?", name).First(&category).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &category, err
}

func (r *TicketCategoryRepository) CreateCategory(category *entity.TicketCategory) error {
	return r.db.Create(category).Error
}
