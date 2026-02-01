package service

import (
	"context"
	"fmt"
	"log"
	"medical-webhook/internal/application/mapper"
	"medical-webhook/internal/domain/line/entity"
	"medical-webhook/internal/domain/line/repository"
	"strings"
)

// MasterDataService - Service สำหรับจัดการ Master Data (Brand, Category, Department, Model)
// ไม่รวม Equipment เพราะ Equipment เป็น Transaction Data ไม่ใช่ Master Data
type MasterDataService interface {
	GetOrCreateBrand(ctx context.Context, name string) (*entity.Brand, bool, error)
	GetOrCreateCategory(ctx context.Context, name, ecriRisk, classification string) (*entity.EquipmentCategory, bool, error)
	GetOrCreateDepartment(ctx context.Context, name string) (*entity.Department, bool, error)
	GetOrCreateModel(ctx context.Context, brandID, categoryID uint, modelName string, lifeExpectancy float64) (*entity.EquipmentModel, bool, error)
	ClearCache() // Clear cache after import session
}

type masterDataService struct {
	brandRepo      repository.BrandRepository
	categoryRepo   repository.EquipmentCategoryRepository
	departmentRepo repository.DepartmentRepository
	modelRepo      repository.EquipmentModelRepository
	mapper         *mapper.EquipmentMapper

	// Cache for reducing duplicate queries during import
	brandCache      map[string]*entity.Brand
	categoryCache   map[string]*entity.EquipmentCategory
	departmentCache map[string]*entity.Department
	modelCache      map[string]*entity.EquipmentModel
}

func NewMasterDataService(
	brandRepo repository.BrandRepository,
	categoryRepo repository.EquipmentCategoryRepository,
	departmentRepo repository.DepartmentRepository,
	modelRepo repository.EquipmentModelRepository,
	mapper *mapper.EquipmentMapper,
) MasterDataService {
	return &masterDataService{
		brandRepo:       brandRepo,
		categoryRepo:    categoryRepo,
		departmentRepo:  departmentRepo,
		modelRepo:       modelRepo,
		mapper:          mapper,
		brandCache:      make(map[string]*entity.Brand),
		categoryCache:   make(map[string]*entity.EquipmentCategory),
		departmentCache: make(map[string]*entity.Department),
		modelCache:      make(map[string]*entity.EquipmentModel),
	}
}

// GetOrCreateBrand - หา Brand หรือสร้างใหม่ถ้ายังไม่มี
func (s *masterDataService) GetOrCreateBrand(ctx context.Context, name string) (*entity.Brand, bool, error) {
	if name == "" {
		return nil, false, fmt.Errorf("brand name is required")
	}

	// Normalize name
	name = strings.TrimSpace(name)

	// Check cache first
	if brand, exists := s.brandCache[name]; exists {
		log.Printf("🔄 Using cached brand: %s", name)
		return brand, false, nil
	}

	// Query database
	brand, err := s.brandRepo.FindByName(ctx, name)
	if err != nil {
		return nil, false, fmt.Errorf("failed to find brand: %w", err)
	}

	// Found in database
	if brand != nil {
		log.Printf("✅ Found existing brand: %s (ID: %d)", brand.Name, brand.ID)
		s.brandCache[name] = brand
		return brand, false, nil
	}

	// Create new brand
	newBrand := s.mapper.ToBrandEntity(name)
	if err := s.brandRepo.Create(ctx, newBrand); err != nil {
		return nil, false, fmt.Errorf("failed to create brand: %w", err)
	}

	log.Printf("🆕 Created new brand: %s (ID: %d)", newBrand.Name, newBrand.ID)
	s.brandCache[name] = newBrand
	return newBrand, true, nil
}

// GetOrCreateCategory - หา Category หรือสร้างใหม่ถ้ายังไม่มี
func (s *masterDataService) GetOrCreateCategory(ctx context.Context, name, ecriRisk, classification string) (*entity.EquipmentCategory, bool, error) {
	if name == "" {
		return nil, false, fmt.Errorf("category name is required")
	}

	// Normalize name
	name = strings.TrimSpace(name)

	// Check cache first
	if category, exists := s.categoryCache[name]; exists {
		log.Printf("🔄 Using cached category: %s", name)
		return category, false, nil
	}

	// Query database
	category, err := s.categoryRepo.FindByName(ctx, name)
	if err != nil {
		return nil, false, fmt.Errorf("failed to find category: %w", err)
	}

	// Found in database
	if category != nil {
		log.Printf("✅ Found existing category: %s (ID: %d)", category.Name, category.ID)
		s.categoryCache[name] = category
		return category, false, nil
	}

	// Normalize ECRI Risk
	normalizedRisk := s.normalizeECRIRisk(ecriRisk)

	// Create new category
	newCategory := s.mapper.ToCategoryEntity(name, normalizedRisk, classification)
	if err := s.categoryRepo.Create(ctx, newCategory); err != nil {
		return nil, false, fmt.Errorf("failed to create category: %w", err)
	}

	log.Printf("🆕 Created new category: %s (ID: %d, Risk: %s)", newCategory.Name, newCategory.ID, newCategory.ECRIRisk)
	s.categoryCache[name] = newCategory
	return newCategory, true, nil
}

// GetOrCreateDepartment - หา Department หรือสร้างใหม่ถ้ายังไม่มี
func (s *masterDataService) GetOrCreateDepartment(ctx context.Context, name string) (*entity.Department, bool, error) {
	if name == "" {
		return nil, false, fmt.Errorf("department name is required")
	}

	// Normalize name
	name = strings.TrimSpace(name)

	// Check cache first
	if dept, exists := s.departmentCache[name]; exists {
		log.Printf("🔄 Using cached department: %s", name)
		return dept, false, nil
	}

	// Query database
	dept, err := s.departmentRepo.FindByName(ctx, name)
	if err != nil {
		return nil, false, fmt.Errorf("failed to find department: %w", err)
	}

	// Found in database
	if dept != nil {
		log.Printf("✅ Found existing department: %s (ID: %d)", dept.Name, dept.ID)
		s.departmentCache[name] = dept
		return dept, false, nil
	}

	// Create new department
	newDept := s.mapper.ToDepartmentEntity(name)
	if err := s.departmentRepo.Create(ctx, newDept); err != nil {
		return nil, false, fmt.Errorf("failed to create department: %w", err)
	}

	log.Printf("🆕 Created new department: %s (ID: %d)", newDept.Name, newDept.ID)
	s.departmentCache[name] = newDept
	return newDept, true, nil
}

// GetOrCreateModel - หา Model หรือสร้างใหม่ถ้ายังไม่มี
func (s *masterDataService) GetOrCreateModel(ctx context.Context, brandID, categoryID uint, modelName string, lifeExpectancy float64) (*entity.EquipmentModel, bool, error) {
	if modelName == "" {
		return nil, false, fmt.Errorf("model name is required")
	}

	// Normalize model name
	modelName = strings.TrimSpace(modelName)

	// Create cache key
	cacheKey := fmt.Sprintf("%d-%d-%s", brandID, categoryID, modelName)

	// Check cache first
	if model, exists := s.modelCache[cacheKey]; exists {
		log.Printf("🔄 Using cached model: %s", modelName)
		return model, false, nil
	}

	// Query database
	model, err := s.modelRepo.FindByBrandCategoryModel(ctx, brandID, categoryID, modelName)
	if err != nil {
		return nil, false, fmt.Errorf("failed to find model: %w", err)
	}

	// Found in database
	if model != nil {
		log.Printf("✅ Found existing model: %s (ID: %d)", model.ModelName, model.ID)
		s.modelCache[cacheKey] = model
		return model, false, nil
	}

	// Create new model
	newModel := s.mapper.ToModelEntity(brandID, categoryID, modelName, lifeExpectancy)
	if err := s.modelRepo.Create(ctx, newModel); err != nil {
		return nil, false, fmt.Errorf("failed to create model: %w", err)
	}

	log.Printf("🆕 Created new model: %s (ID: %d, Life: %.0f years)", newModel.ModelName, newModel.ID, newModel.DefaultLifeExpectancy)
	s.modelCache[cacheKey] = newModel
	return newModel, true, nil
}

// normalizeECRIRisk - normalize ECRI Risk value
func (s *masterDataService) normalizeECRIRisk(risk string) string {
	normalized := strings.ToUpper(strings.TrimSpace(risk))
	switch normalized {
	case "HIGH", "H", "สูง":
		return "HIGH"
	case "LOW", "L", "ต่ำ":
		return "LOW"
	case "MEDIUM", "MED", "MODERATE", "M", "ปานกลาง", "":
		return "MEDIUM"
	default:
		log.Printf("⚠️ Unknown ECRI risk value: %s, defaulting to MEDIUM", risk)
		return "MEDIUM"
	}
}

// ClearCache - ล้าง cache (เรียกหลังจบ import แต่ละครั้ง)
func (s *masterDataService) ClearCache() {
	log.Println("🧹 Clearing master data cache...")
	s.brandCache = make(map[string]*entity.Brand)
	s.categoryCache = make(map[string]*entity.EquipmentCategory)
	s.departmentCache = make(map[string]*entity.Department)
	s.modelCache = make(map[string]*entity.EquipmentModel)
	log.Println("✅ Cache cleared")
}
