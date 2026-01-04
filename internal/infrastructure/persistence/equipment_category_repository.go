package persistence

import (
	"context"
	"log"
	"medical-webhook/internal/domain/line/entity"
	"medical-webhook/internal/infrastructure/database"

	"gorm.io/gorm"
)

// EquipmentCategoryRepository implements repository.EquipmentCategoryRepository using GORM
type EquipmentCategoryRepository struct {
	db *gorm.DB
}

// NewEquipmentCategoryRepository creates a new equipment category repository
func NewEquipmentCategoryRepository() *EquipmentCategoryRepository {
	return &EquipmentCategoryRepository{
		db: database.DB,
	}
}

// Create creates a new equipment category
func (r *EquipmentCategoryRepository) Create(ctx context.Context, category *entity.EquipmentCategory) error {
	err := r.db.WithContext(ctx).Create(category).Error
	if err != nil {
		log.Printf("Error creating equipment category: %v", err)
		return err
	}
	log.Printf("Created equipment category: %s (ID: %d)", category.Name, category.ID)
	return nil
}

// FindByID finds equipment category by ID
func (r *EquipmentCategoryRepository) FindByID(ctx context.Context, id uint) (*entity.EquipmentCategory, error) {
	var category entity.EquipmentCategory
	err := r.db.WithContext(ctx).First(&category, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Equipment category not found with ID: %d", id)
			return nil, nil
		}
		log.Printf("Error finding equipment category by ID: %v", err)
		return nil, err
	}
	log.Printf("Found equipment category: %s (ID: %d)", category.Name, category.ID)
	return &category, nil
}

// FindByName finds equipment category by name
func (r *EquipmentCategoryRepository) FindByName(ctx context.Context, name string) (*entity.EquipmentCategory, error) {
	var category entity.EquipmentCategory
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Equipment category not found with name: %s", name)
			return nil, nil
		}
		log.Printf("Error finding equipment category by name: %v", err)
		return nil, err
	}
	log.Printf("Found equipment category: %s (ID: %d)", category.Name, category.ID)
	return &category, nil
}

// FindAll finds all equipment categories
func (r *EquipmentCategoryRepository) FindAll(ctx context.Context) ([]entity.EquipmentCategory, error) {
	var categories []entity.EquipmentCategory
	err := r.db.WithContext(ctx).Order("name ASC").Find(&categories).Error
	if err != nil {
		log.Printf("Error finding all equipment categories: %v", err)
		return nil, err
	}
	log.Printf("Found %d equipment categories", len(categories))
	return categories, nil
}

// FindByECRIRisk finds categories by ECRI risk level
func (r *EquipmentCategoryRepository) FindByECRIRisk(ctx context.Context, risk entity.ECRIRiskLevel) ([]entity.EquipmentCategory, error) {
	var categories []entity.EquipmentCategory
	err := r.db.WithContext(ctx).Where("ecri_risk = ?", risk).Order("name ASC").Find(&categories).Error
	if err != nil {
		log.Printf("Error finding equipment categories by ECRI risk: %v", err)
		return nil, err
	}
	log.Printf("Found %d equipment categories with ECRI risk: %s", len(categories), risk)
	return categories, nil
}

// Update updates equipment category
func (r *EquipmentCategoryRepository) Update(ctx context.Context, category *entity.EquipmentCategory) error {
	err := r.db.WithContext(ctx).Save(category).Error
	if err != nil {
		log.Printf("Error updating equipment category: %v", err)
		return err
	}
	log.Printf("Updated equipment category: %s (ID: %d)", category.Name, category.ID)
	return nil
}

// Delete soft deletes equipment category
func (r *EquipmentCategoryRepository) Delete(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).Delete(&entity.EquipmentCategory{}, id).Error
	if err != nil {
		log.Printf("Error deleting equipment category: %v", err)
		return err
	}
	log.Printf("Deleted equipment category ID: %d", id)
	return nil
}
