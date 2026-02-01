package usecase

import (
	"context"
	"errors"
	"log"
	"medical-webhook/internal/application/dto"
	"medical-webhook/internal/application/mapper"
	"medical-webhook/internal/domain/line/entity"
	"medical-webhook/internal/domain/line/service"
	"time"
)

type EquipmentUsecase interface {
	GetEquipmentList(ctx context.Context, req dto.EquipmentListRequest) (*dto.EquipmentListResponse, error)
	CreateEquipment(ctx context.Context, req dto.CreateEquipmentRequest) (*dto.EquipmentResponse, error)
}

type equipmentUsecase struct {
	equipmentService service.EquipmentService
	mapper           *mapper.EquipmentMapper
}

func NewEquipmentUsecase(equipmentService service.EquipmentService) EquipmentUsecase {
	return &equipmentUsecase{
		equipmentService: equipmentService,
		mapper:           mapper.NewEquipmentMapper(),
	}
}

// GetEquipmentList - ดึงรายการ Equipment แบบ pagination
func (u *equipmentUsecase) GetEquipmentList(ctx context.Context, req dto.EquipmentListRequest) (*dto.EquipmentListResponse, error) {
	log.Printf("Usecase: GetEquipmentList - Page: %d, Limit: %d", req.Page, req.Limit)

	// Set default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100 // Max limit
	}

	// Calculate offset
	offset := (req.Page - 1) * req.Limit

	// Get total count from service
	total, err := u.equipmentService.CountEquipments(ctx)
	if err != nil {
		log.Printf("Usecase: Error counting equipments: %v", err)
		return nil, err
	}

	// Get equipment list with pagination from service
	equipments, err := u.equipmentService.FindAllEquipments(ctx, req.Limit, offset)
	if err != nil {
		log.Printf("Usecase: Error getting equipments: %v", err)
		return nil, err
	}

	// Map to DTO using mapper
	items := make([]dto.EquipmentListItem, 0, len(equipments))
	for _, e := range equipments {
		item := u.mapper.MapEquipmentToListItem(&e)
		items = append(items, *item)
	}

	// Calculate total pages
	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	log.Printf("Usecase: Successfully retrieved %d equipments", len(items))

	return &dto.EquipmentListResponse{
		Data:       items,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// CreateEquipment - สร้าง Equipment ใหม่จากข้อมูลที่กรอกในฟอร์ม
func (u *equipmentUsecase) CreateEquipment(ctx context.Context, req dto.CreateEquipmentRequest) (*dto.EquipmentResponse, error) {
	log.Printf("Usecase: CreateEquipment - IDCode: %s, SerialNo: %s", req.IDCode, req.SerialNo)

	// 1. Validate required fields
	if err := u.validateCreateRequest(req); err != nil {
		log.Printf("Usecase: Validation failed: %v", err)
		return nil, err
	}

	// 2. Check if equipment with same IDCode already exists
	existingEquip, err := u.equipmentService.FindEquipmentByIDCode(ctx, req.IDCode)
	if err != nil {
		log.Printf("Usecase: Error checking existing equipment: %v", err)
		return nil, err
	}
	if existingEquip != nil {
		return nil, errors.New("equipment with this ID code already exists")
	}

	// 3. Find or Create Brand
	brand, err := u.equipmentService.FindOrCreateBrand(ctx, req.Brand)
	if err != nil {
		log.Printf("Usecase: Error with brand: %v", err)
		return nil, err
	}

	// 4. Find or Create Category
	category, err := u.equipmentService.FindOrCreateCategory(ctx, req.Category, "MEDIUM", "General")
	if err != nil {
		log.Printf("Usecase: Error with category: %v", err)
		return nil, err
	}

	// 5. Find or Create Department
	department, err := u.equipmentService.FindOrCreateDepartment(ctx, req.Department)
	if err != nil {
		log.Printf("Usecase: Error with department: %v", err)
		return nil, err
	}

	// 6. Find or Create Model
	model, err := u.equipmentService.FindOrCreateModel(ctx, req.Model, brand.ID, category.ID, req.LifeExpectancy)
	if err != nil {
		log.Printf("Usecase: Error with model: %v", err)
		return nil, err
	}

	// 7. Parse dates
	receiveDate, err := time.Parse("2006-01-02", req.ReceiveDate)
	if err != nil {
		return nil, errors.New("invalid receive_date format, use YYYY-MM-DD")
	}

	// 8. Create Equipment entity
	equipment := &entity.Equipment{
		IDCode:                req.IDCode,
		SerialNo:              &req.SerialNo,
		ModelID:               model.ID,
		DepartmentID:          department.ID,
		ReceiveDate:           &receiveDate,
		PurchasePrice:         req.PurchasePrice,
		EquipmentAge:          req.EquipmentAge,
		LifeExpectancy:        req.LifeExpectancy,
		RemainLife:            req.RemainLife,
		UsefulLifetimePercent: req.UsefulLifetimePercent,
	}

	// Set optional fields
	if req.AssessmentID != "" {
		equipment.AssessmentID = &req.AssessmentID
	}
	if req.ComputeDate != "" {
		computeDate, err := time.Parse("2006-01-02", req.ComputeDate)
		if err == nil {
			equipment.ComputeDate = &computeDate
		}
	}
	if req.ReplacementYear > 0 {
		equipment.ReplacementYear = &req.ReplacementYear
	}
	if req.Technology != nil {
		equipment.Technology = req.Technology
	}
	if req.UsageStatistics != nil {
		equipment.UsageStatistics = req.UsageStatistics
	}
	if req.Efficiency != nil {
		equipment.Efficiency = req.Efficiency
	}
	if req.Others != "" {
		equipment.Others = &req.Others
	}

	// 9. Calculate and set Status based on RemainLife
	equipment.Status = u.calculateStatus(req.RemainLife)

	// 10. Save to database via service
	err = u.equipmentService.CreateEquipment(ctx, equipment)
	if err != nil {
		log.Printf("Usecase: Error creating equipment: %v", err)
		return nil, err
	}

	// 11. Load relations and return
	createdEquipment, err := u.equipmentService.FindEquipmentByID(ctx, equipment.ID)
	if err != nil {
		log.Printf("Usecase: Error loading created equipment: %v", err)
		return nil, err
	}

	log.Printf("Usecase: Successfully created equipment: %s (ID: %d)", equipment.IDCode, equipment.ID)

	// 12. Map to response DTO using mapper
	return u.mapper.MapEquipmentToResponse(createdEquipment), nil
}

// ===== VALIDATION =====

func (u *equipmentUsecase) validateCreateRequest(req dto.CreateEquipmentRequest) error {
	if req.IDCode == "" {
		return errors.New("id_code is required")
	}
	if req.SerialNo == "" {
		return errors.New("serial_no is required")
	}
	if req.Department == "" {
		return errors.New("department is required")
	}
	if req.Brand == "" {
		return errors.New("brand is required")
	}
	if req.Model == "" {
		return errors.New("model is required")
	}
	if req.Category == "" {
		return errors.New("category is required")
	}
	if req.ReceiveDate == "" {
		return errors.New("receive_date is required")
	}
	if req.PurchasePrice < 0 {
		return errors.New("purchase_price must be greater than or equal to 0")
	}

	// Business logic validations
	if req.EquipmentAge > req.LifeExpectancy {
		return errors.New("equipment_age cannot exceed life_expectancy")
	}
	if req.UsefulLifetimePercent < 0 || req.UsefulLifetimePercent > 100 {
		return errors.New("useful_lifetime_percent must be between 0 and 100")
	}

	return nil
}

// ===== HELPER FUNCTIONS =====

func (u *equipmentUsecase) calculateStatus(remainLife float64) entity.AssetStatus {
	if remainLife <= 0 {
		return entity.AssetStatusPlanToReplace // ถึงเวลาเปลี่ยนแล้ว
	} else if remainLife <= 1 {
		return entity.AssetStatusPlanToReplace // เหลืออายุไม่ถึง 1 ปี
	}
	return entity.AssetStatusActive // ใช้งานปกติ
}
