package persistence

import (
	"context"
	"log"
	"medical-webhook/internal/domain/line/entity"
	"medical-webhook/internal/infrastructure/database"

	"gorm.io/gorm"
)

// DepartmentRepository implements repository.DepartmentRepository using GORM
type DepartmentRepository struct {
	db *gorm.DB
}

// NewDepartmentRepository creates a new department repository
func NewDepartmentRepository() *DepartmentRepository {
	return &DepartmentRepository{
		db: database.DB,
	}
}

// Create creates a new department
func (r *DepartmentRepository) Create(ctx context.Context, department *entity.Department) error {
	err := r.db.WithContext(ctx).Create(department).Error
	if err != nil {
		log.Printf("Error creating department: %v", err)
		return err
	}
	log.Printf("Created department: %s (ID: %d)", department.Name, department.ID)
	return nil
}

// FindByID finds department by ID
func (r *DepartmentRepository) FindByID(ctx context.Context, id uint) (*entity.Department, error) {
	var department entity.Department
	err := r.db.WithContext(ctx).First(&department, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Department not found with ID: %d", id)
			return nil, nil
		}
		log.Printf("Error finding department by ID: %v", err)
		return nil, err
	}
	log.Printf("Found department: %s (ID: %d)", department.Name, department.ID)
	return &department, nil
}

// FindByName finds department by name
func (r *DepartmentRepository) FindByName(ctx context.Context, name string) (*entity.Department, error) {
	var department entity.Department
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&department).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Department not found with name: %s", name)
			return nil, nil
		}
		log.Printf("Error finding department by name: %v", err)
		return nil, err
	}
	log.Printf("Found department: %s (ID: %d)", department.Name, department.ID)
	return &department, nil
}

// FindAll finds all departments
func (r *DepartmentRepository) FindAll(ctx context.Context) ([]entity.Department, error) {
	var departments []entity.Department
	err := r.db.WithContext(ctx).Order("name ASC").Find(&departments).Error
	if err != nil {
		log.Printf("Error finding all departments: %v", err)
		return nil, err
	}
	log.Printf("Found %d departments", len(departments))
	return departments, nil
}

// Update updates department
func (r *DepartmentRepository) Update(ctx context.Context, department *entity.Department) error {
	err := r.db.WithContext(ctx).Save(department).Error
	if err != nil {
		log.Printf("Error updating department: %v", err)
		return err
	}
	log.Printf("Updated department: %s (ID: %d)", department.Name, department.ID)
	return nil
}

// Delete soft deletes department
func (r *DepartmentRepository) Delete(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).Delete(&entity.Department{}, id).Error
	if err != nil {
		log.Printf("Error deleting department: %v", err)
		return err
	}
	log.Printf("Deleted department ID: %d", id)
	return nil
}

// FindOrCreate finds department by name or creates a new one
func (r *DepartmentRepository) FindOrCreate(ctx context.Context, name string) (*entity.Department, error) {
	// Try to find existing department
	dept, err := r.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if dept != nil {
		return dept, nil
	}

	// Create new department
	newDept := &entity.Department{Name: name}
	if err := r.Create(ctx, newDept); err != nil {
		return nil, err
	}
	return newDept, nil
}

// SearchByNameLike searches departments by keyword using LIKE
func (r *DepartmentRepository) SearchByNameLike(ctx context.Context, keyword string, limit int) ([]entity.Department, error) {
	var departments []entity.Department
	err := r.db.WithContext(ctx).
		Where("name LIKE ?", "%"+keyword+"%").
		Order("name ASC").
		Limit(limit).
		Find(&departments).Error
	if err != nil {
		log.Printf("Error searching departments by keyword '%s': %v", keyword, err)
		return nil, err
	}
	log.Printf("Found %d departments matching '%s'", len(departments), keyword)
	return departments, nil
}
