package service

import (
	"context"
	"errors"
	"log"
	"medical-webhook/internal/domain/line/entity"
	"medical-webhook/internal/domain/line/repository"

	"gorm.io/gorm"
)

type EquipmentService interface {
	// Equipment CRUD
	FindEquipmentByIDCode(ctx context.Context, idCode string) (*entity.Equipment, error)
	FindEquipmentByID(ctx context.Context, id uint) (*entity.Equipment, error)
	FindAllEquipments(ctx context.Context, limit, offset int) ([]entity.Equipment, error)
	FindAllEquipmentsWithFilter(ctx context.Context, limit, offset int, status, search, expiryFilter string) ([]entity.Equipment, error)
	CountEquipments(ctx context.Context) (int64, error)
	CountEquipmentsWithFilter(ctx context.Context, status, search, expiryFilter string) (int64, error)
	CreateEquipment(ctx context.Context, equipment *entity.Equipment) error
	UpdateEquipment(ctx context.Context, equipment *entity.Equipment) error
	DeleteEquipment(ctx context.Context, id uint) error

	// Categories
	GetAllCategories(ctx context.Context) ([]entity.EquipmentCategory, error)

	// Brand
	FindOrCreateBrand(ctx context.Context, brandName string) (*entity.Brand, error)

	// Category
	FindOrCreateCategory(ctx context.Context, categoryName string, riskLevel entity.ECRIRiskLevel, description string) (*entity.EquipmentCategory, error)

	// Department
	FindOrCreateDepartment(ctx context.Context, departmentName string) (*entity.Department, error)

	// Model
	FindOrCreateModel(ctx context.Context, modelName string, brandID, categoryID uint, lifeExpectancy float64) (*entity.EquipmentModel, error)
}

type equipmentService struct {
	equipmentRepo  repository.EquipmentRepository
	brandRepo      repository.BrandRepository
	categoryRepo   repository.EquipmentCategoryRepository
	departmentRepo repository.DepartmentRepository
	modelRepo      repository.EquipmentModelRepository
}

func NewEquipmentService(
	equipmentRepo repository.EquipmentRepository,
	brandRepo repository.BrandRepository,
	categoryRepo repository.EquipmentCategoryRepository,
	departmentRepo repository.DepartmentRepository,
	modelRepo repository.EquipmentModelRepository,
) EquipmentService {
	return &equipmentService{
		equipmentRepo:  equipmentRepo,
		brandRepo:      brandRepo,
		categoryRepo:   categoryRepo,
		departmentRepo: departmentRepo,
		modelRepo:      modelRepo,
	}
}

func (s *equipmentService) FindEquipmentByIDCode(ctx context.Context, idCode string) (*entity.Equipment, error) {
	equipment, err := s.equipmentRepo.FindByIDCode(idCode)
	if err != nil {
		log.Printf("Service: Error finding equipment by ID code: %v", err)
		return nil, err
	}
	return equipment, nil
}

func (s *equipmentService) FindEquipmentByID(ctx context.Context, id uint) (*entity.Equipment, error) {
	equipment, err := s.equipmentRepo.FindByID(ctx, id)
	if err != nil {
		log.Printf("Service: Error finding equipment by ID: %v", err)
		return nil, err
	}
	return equipment, nil
}

func (s *equipmentService) FindAllEquipments(ctx context.Context, limit, offset int) ([]entity.Equipment, error) {
	equipments, err := s.equipmentRepo.FindAll(ctx, limit, offset)
	if err != nil {
		log.Printf("Service: Error finding all equipments: %v", err)
		return nil, err
	}
	return equipments, nil
}

func (s *equipmentService) CountEquipments(ctx context.Context) (int64, error) {
	count, err := s.equipmentRepo.Count(ctx)
	if err != nil {
		log.Printf("Service: Error counting equipments: %v", err)
		return 0, err
	}
	return count, nil
}

func (s *equipmentService) FindAllEquipmentsWithFilter(ctx context.Context, limit, offset int, status, search, expiryFilter string) ([]entity.Equipment, error) {
	equipments, err := s.equipmentRepo.FindAllWithFilter(ctx, limit, offset, status, search, expiryFilter)
	if err != nil {
		log.Printf("Service: Error finding equipments with filter: %v", err)
		return nil, err
	}
	return equipments, nil
}

func (s *equipmentService) CountEquipmentsWithFilter(ctx context.Context, status, search, expiryFilter string) (int64, error) {
	count, err := s.equipmentRepo.CountWithFilter(ctx, status, search, expiryFilter)
	if err != nil {
		log.Printf("Service: Error counting equipments with filter: %v", err)
		return 0, err
	}
	return count, nil
}

func (s *equipmentService) CreateEquipment(ctx context.Context, equipment *entity.Equipment) error {
	// Validate equipment entity
	if err := s.validateEquipment(equipment); err != nil {
		return err
	}

	err := s.equipmentRepo.Create(ctx, equipment)
	if err != nil {
		log.Printf("Service: Error creating equipment: %v", err)
		return err
	}

	log.Printf("Service: Successfully created equipment ID: %d", equipment.ID)
	return nil
}

func (s *equipmentService) UpdateEquipment(ctx context.Context, equipment *entity.Equipment) error {
	// Validate equipment entity
	if err := s.validateEquipment(equipment); err != nil {
		return err
	}

	err := s.equipmentRepo.Update(ctx, equipment)
	if err != nil {
		log.Printf("Service: Error updating equipment: %v", err)
		return err
	}

	log.Printf("Service: Successfully updated equipment ID: %d", equipment.ID)
	return nil
}

func (s *equipmentService) DeleteEquipment(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("equipment ID is required")
	}

	err := s.equipmentRepo.Delete(ctx, id)
	if err != nil {
		log.Printf("Service: Error deleting equipment: %v", err)
		return err
	}

	log.Printf("Service: Successfully deleted equipment ID: %d", id)
	return nil
}

func (s *equipmentService) GetAllCategories(ctx context.Context) ([]entity.EquipmentCategory, error) {
	categories, err := s.categoryRepo.FindAll(ctx)
	if err != nil {
		log.Printf("Service: Error finding all categories: %v", err)
		return nil, err
	}
	log.Printf("Service: Found %d categories", len(categories))
	return categories, nil
}

func (s *equipmentService) FindOrCreateBrand(ctx context.Context, brandName string) (*entity.Brand, error) {
	// Try to find existing brand
	brand, err := s.brandRepo.FindByName(ctx, brandName)
	if err != nil {
		log.Printf("Service: Error finding brand: %v", err)
		return nil, err
	}

	// If brand exists, return it
	if brand != nil {
		log.Printf("Service: Found existing brand: %s (ID: %d)", brandName, brand.ID)
		return brand, nil
	}

	// Create new brand
	newBrand := &entity.Brand{
		Name: brandName,
	}

	err = s.brandRepo.Create(ctx, newBrand)
	if err != nil {
		log.Printf("Service: Error creating brand: %v", err)
		return nil, err
	}

	log.Printf("Service: Created new brand: %s (ID: %d)", brandName, newBrand.ID)
	return newBrand, nil
}

func (s *equipmentService) FindOrCreateCategory(ctx context.Context, categoryName string, riskLevel entity.ECRIRiskLevel, description string) (*entity.EquipmentCategory, error) {
	// Try to find existing category
	category, err := s.categoryRepo.FindByName(ctx, categoryName)
	if err != nil {
		log.Printf("Service: Error finding category: %v", err)
		return nil, err
	}

	// If category exists, return it
	if category != nil {
		log.Printf("Service: Found existing category: %s (ID: %d)", categoryName, category.ID)
		return category, nil
	}

	// Create new category
	newCategory := &entity.EquipmentCategory{
		Name:           categoryName,
		ECRIRisk:       riskLevel,
		Classification: description,
	}

	err = s.categoryRepo.Create(ctx, newCategory)
	if err != nil {
		log.Printf("Service: Error creating category: %v", err)
		return nil, err
	}

	log.Printf("Service: Created new category: %s (ID: %d)", categoryName, newCategory.ID)
	return newCategory, nil
}

func (s *equipmentService) FindOrCreateDepartment(ctx context.Context, departmentName string) (*entity.Department, error) {
	// Try to find existing department
	department, err := s.departmentRepo.FindByName(ctx, departmentName)
	if err != nil {
		log.Printf("Service: Error finding department: %v", err)
		return nil, err
	}

	// If department exists, return it
	if department != nil {
		log.Printf("Service: Found existing department: %s (ID: %d)", departmentName, department.ID)
		return department, nil
	}

	// Create new department
	newDepartment := &entity.Department{
		Name: departmentName,
	}

	err = s.departmentRepo.Create(ctx, newDepartment)
	if err != nil {
		log.Printf("Service: Error creating department: %v", err)
		return nil, err
	}

	log.Printf("Service: Created new department: %s (ID: %d)", departmentName, newDepartment.ID)
	return newDepartment, nil
}

func (s *equipmentService) FindOrCreateModel(ctx context.Context, modelName string, brandID, categoryID uint, lifeExpectancy float64) (*entity.EquipmentModel, error) {
	// Try to find existing model with same name, brand, and category
	models, err := s.modelRepo.FindAll(ctx)
	if err != nil {
		log.Printf("Service: Error finding models: %v", err)
		return nil, err
	}

	for _, m := range models {
		if m.ModelName == modelName && m.BrandID == brandID && m.CategoryID == categoryID {
			log.Printf("Service: Found existing model: %s (ID: %d)", modelName, m.ID)
			return &m, nil
		}
	}

	newModel := &entity.EquipmentModel{
		BrandID:               brandID,
		CategoryID:            categoryID,
		ModelName:             modelName,
		DefaultLifeExpectancy: lifeExpectancy,
	}

	err = s.modelRepo.Create(ctx, newModel)
	if err != nil {
		log.Printf("Service: Error creating model: %v", err)
		return nil, err
	}

	log.Printf("Service: Created new model: %s (ID: %d)", modelName, newModel.ID)
	return newModel, nil
}

func (s *equipmentService) validateEquipment(equipment *entity.Equipment) error {
	if equipment.IDCode == "" {
		return errors.New("equipment ID code is required")
	}
	if equipment.ModelID == 0 {
		return errors.New("equipment model ID is required")
	}
	if equipment.DepartmentID == 0 {
		return errors.New("equipment department ID is required")
	}
	return nil
}
func FindSimilarByIDCode(db *gorm.DB, idCode string) ([]*entity.Equipment, error) {
	prefix := idCode[:6] // SSH017

	var equipments []*entity.Equipment
	err := db.
		Where("id_code LIKE ?", prefix+"%").
		Order("id_code").
		Find(&equipments).Error

	return equipments, err
}
