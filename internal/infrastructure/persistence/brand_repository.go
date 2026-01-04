package persistence

import (
	"context"
	"log"
	"medical-webhook/internal/domain/line/entity"
	"medical-webhook/internal/infrastructure/database"

	"gorm.io/gorm"
)

// BrandRepository implements repository.BrandRepository using GORM
type BrandRepository struct {
	db *gorm.DB
}

// NewBrandRepository creates a new brand repository
func NewBrandRepository() *BrandRepository {
	return &BrandRepository{
		db: database.DB,
	}
}

// Create creates a new brand
func (r *BrandRepository) Create(ctx context.Context, brand *entity.Brand) error {
	err := r.db.WithContext(ctx).Create(brand).Error
	if err != nil {
		log.Printf("Error creating brand: %v", err)
		return err
	}
	log.Printf("Created brand: %s (ID: %d)", brand.Name, brand.ID)
	return nil
}

// FindByID finds brand by ID
func (r *BrandRepository) FindByID(ctx context.Context, id uint) (*entity.Brand, error) {
	var brand entity.Brand
	err := r.db.WithContext(ctx).First(&brand, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Brand not found with ID: %d", id)
			return nil, nil
		}
		log.Printf("Error finding brand by ID: %v", err)
		return nil, err
	}
	log.Printf("Found brand: %s (ID: %d)", brand.Name, brand.ID)
	return &brand, nil
}

// FindByName finds brand by name
func (r *BrandRepository) FindByName(ctx context.Context, name string) (*entity.Brand, error) {
	var brand entity.Brand
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&brand).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Brand not found with name: %s", name)
			return nil, nil
		}
		log.Printf("Error finding brand by name: %v", err)
		return nil, err
	}
	log.Printf("Found brand: %s (ID: %d)", brand.Name, brand.ID)
	return &brand, nil
}

// FindAll finds all brands
func (r *BrandRepository) FindAll(ctx context.Context) ([]entity.Brand, error) {
	var brands []entity.Brand
	err := r.db.WithContext(ctx).Order("name ASC").Find(&brands).Error
	if err != nil {
		log.Printf("Error finding all brands: %v", err)
		return nil, err
	}
	log.Printf("Found %d brands", len(brands))
	return brands, nil
}

// Update updates brand
func (r *BrandRepository) Update(ctx context.Context, brand *entity.Brand) error {
	err := r.db.WithContext(ctx).Save(brand).Error
	if err != nil {
		log.Printf("Error updating brand: %v", err)
		return err
	}
	log.Printf("Updated brand: %s (ID: %d)", brand.Name, brand.ID)
	return nil
}

// Delete soft deletes brand
func (r *BrandRepository) Delete(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).Delete(&entity.Brand{}, id).Error
	if err != nil {
		log.Printf("Error deleting brand: %v", err)
		return err
	}
	log.Printf("Deleted brand ID: %d", id)
	return nil
}
