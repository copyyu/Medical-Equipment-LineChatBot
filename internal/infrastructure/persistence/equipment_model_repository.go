package persistence

import (
	"context"
	"log"
	"medical-webhook/internal/domain/line/entity"
	"medical-webhook/internal/infrastructure/database"

	"gorm.io/gorm"
)

// EquipmentModelRepository implements repository.EquipmentModelRepository using GORM
type EquipmentModelRepository struct {
	db *gorm.DB
}

// NewEquipmentModelRepository creates a new equipment model repository
func NewEquipmentModelRepository() *EquipmentModelRepository {
	return &EquipmentModelRepository{
		db: database.DB,
	}
}

// Create creates a new equipment model
func (r *EquipmentModelRepository) Create(ctx context.Context, model *entity.EquipmentModel) error {
	err := r.db.WithContext(ctx).Create(model).Error
	if err != nil {
		log.Printf("Error creating equipment model: %v", err)
		return err
	}
	log.Printf("Created equipment model: %s (ID: %d)", model.ModelName, model.ID)
	return nil
}

// FindByID finds equipment model by ID
func (r *EquipmentModelRepository) FindByID(ctx context.Context, id uint) (*entity.EquipmentModel, error) {
	var model entity.EquipmentModel
	err := r.db.WithContext(ctx).
		Preload("Brand").
		Preload("Category").
		First(&model, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Equipment model not found with ID: %d", id)
			return nil, nil
		}
		log.Printf("Error finding equipment model by ID: %v", err)
		return nil, err
	}
	log.Printf("Found equipment model: %s (ID: %d)", model.ModelName, model.ID)
	return &model, nil
}

// FindByBrandCategoryModel finds model by brand_id, category_id, and model_name
func (r *EquipmentModelRepository) FindByBrandCategoryModel(ctx context.Context, brandID, categoryID uint, modelName string) (*entity.EquipmentModel, error) {
	var model entity.EquipmentModel
	err := r.db.WithContext(ctx).
		Preload("Brand").
		Preload("Category").
		Where("brand_id = ? AND category_id = ? AND model_name = ?", brandID, categoryID, modelName).
		First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Equipment model not found: Brand ID %d, Category ID %d, Model %s", brandID, categoryID, modelName)
			return nil, nil
		}
		log.Printf("Error finding equipment model: %v", err)
		return nil, err
	}
	log.Printf("Found equipment model: %s (ID: %d)", model.ModelName, model.ID)
	return &model, nil
}

// FindByBrandID finds all models by brand ID
func (r *EquipmentModelRepository) FindByBrandID(ctx context.Context, brandID uint) ([]entity.EquipmentModel, error) {
	var models []entity.EquipmentModel
	err := r.db.WithContext(ctx).
		Preload("Brand").
		Preload("Category").
		Where("brand_id = ?", brandID).
		Order("model_name ASC").
		Find(&models).Error
	if err != nil {
		log.Printf("❌ Error finding models by brand ID: %v", err)
		return nil, err
	}
	log.Printf("Found %d models for brand ID: %d", len(models), brandID)
	return models, nil
}

// FindByCategoryID finds all models by category ID
func (r *EquipmentModelRepository) FindByCategoryID(ctx context.Context, categoryID uint) ([]entity.EquipmentModel, error) {
	var models []entity.EquipmentModel
	err := r.db.WithContext(ctx).
		Preload("Brand").
		Preload("Category").
		Where("category_id = ?", categoryID).
		Order("model_name ASC").
		Find(&models).Error
	if err != nil {
		log.Printf("❌ Error finding models by category ID: %v", err)
		return nil, err
	}
	log.Printf("Found %d models for category ID: %d", len(models), categoryID)
	return models, nil
}

// FindAll finds all equipment models with preloaded relations
func (r *EquipmentModelRepository) FindAll(ctx context.Context) ([]entity.EquipmentModel, error) {
	var models []entity.EquipmentModel
	err := r.db.WithContext(ctx).
		Preload("Brand").
		Preload("Category").
		Order("model_name ASC").
		Find(&models).Error
	if err != nil {
		log.Printf("Error finding all equipment models: %v", err)
		return nil, err
	}
	log.Printf("Found %d equipment models", len(models))
	return models, nil
}

// Update updates equipment model
func (r *EquipmentModelRepository) Update(ctx context.Context, model *entity.EquipmentModel) error {
	err := r.db.WithContext(ctx).Save(model).Error
	if err != nil {
		log.Printf("Error updating equipment model: %v", err)
		return err
	}
	log.Printf("Updated equipment model: %s (ID: %d)", model.ModelName, model.ID)
	return nil
}

// Delete soft deletes equipment model
func (r *EquipmentModelRepository) Delete(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).Delete(&entity.EquipmentModel{}, id).Error
	if err != nil {
		log.Printf("Error deleting equipment model: %v", err)
		return err
	}
	log.Printf("Deleted equipment model ID: %d", id)
	return nil
}
