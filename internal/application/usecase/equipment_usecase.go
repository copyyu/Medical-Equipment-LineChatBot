package usecase

import (
	"context"
	"errors"
	"log"
	"medical-webhook/internal/application/dto"
	"medical-webhook/internal/application/mapper"
	"medical-webhook/internal/application/service"
	"medical-webhook/internal/domain/line/entity"
	"time"
)

type EquipmentUsecase interface {
	GetEquipmentList(ctx context.Context, req dto.EquipmentListRequest) (*dto.EquipmentListResponse, error)
	GetByIDCode(ctx context.Context, idCode string) (*dto.EquipmentDetailResponse, error)
	UpdateEquipment(ctx context.Context, idCode string, req dto.EquipmentUpdateRequest) error
	DeleteEquipment(ctx context.Context, idCode string) error
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

// GetEquipmentList - ดึงรายการ Equipment แบบ pagination พร้อม filter
func (u *equipmentUsecase) GetEquipmentList(ctx context.Context, req dto.EquipmentListRequest) (*dto.EquipmentListResponse, error) {
	log.Printf("Usecase: GetEquipmentList - Page: %d, Limit: %d, Status: %s, Search: %s", req.Page, req.Limit, req.Status, req.Search)

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

	var total int64
	var equipments []entity.Equipment
	var err error

	// Check if we need to apply filters
	hasFilter := req.Status != "" || req.Search != ""

	if hasFilter {
		// Get total count with filters
		total, err = u.equipmentService.CountEquipmentsWithFilter(ctx, req.Status, req.Search)
		if err != nil {
			log.Printf("Usecase: Error counting equipments with filter: %v", err)
			return nil, err
		}

		// Get equipment list with pagination and filters
		equipments, err = u.equipmentService.FindAllEquipmentsWithFilter(ctx, req.Limit, offset, req.Status, req.Search)
		if err != nil {
			log.Printf("Usecase: Error getting equipments with filter: %v", err)
			return nil, err
		}
	} else {
		// Get total count without filters
		total, err = u.equipmentService.CountEquipments(ctx)
		if err != nil {
			log.Printf("Usecase: Error counting equipments: %v", err)
			return nil, err
		}

		// Get equipment list with pagination without filters
		equipments, err = u.equipmentService.FindAllEquipments(ctx, req.Limit, offset)
		if err != nil {
			log.Printf("Usecase: Error getting equipments: %v", err)
			return nil, err
		}
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

	log.Printf("Usecase: Successfully retrieved %d equipments (total: %d, filtered: %v)", len(items), total, hasFilter)

	return &dto.EquipmentListResponse{
		Data:       items,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetByIDCode returns equipment detail by ID code
func (u *equipmentUsecase) GetByIDCode(ctx context.Context, idCode string) (*dto.EquipmentDetailResponse, error) {
	equipment, err := u.equipmentService.FindEquipmentByIDCode(ctx, idCode)
	if err != nil {
		return nil, err
	}
	if equipment == nil {
		return nil, errors.New("equipment not found")
	}

	result := u.mapEquipmentToDetailResponse(equipment)
	return result, nil
}

// UpdateEquipment updates equipment by ID code
func (u *equipmentUsecase) UpdateEquipment(ctx context.Context, idCode string, req dto.EquipmentUpdateRequest) error {
	// Find existing equipment
	equipment, err := u.equipmentService.FindEquipmentByIDCode(ctx, idCode)
	if err != nil {
		return err
	}
	if equipment == nil {
		return errors.New("equipment not found")
	}

	// Update status if provided
	if req.Status != "" {
		equipment.Status = entity.AssetStatus(req.Status)
	}

	// Update department if location provided
	if req.Location != "" {
		dept, err := u.equipmentService.FindOrCreateDepartment(ctx, req.Location)
		if err != nil {
			return err
		}
		equipment.DepartmentID = dept.ID
	}

	// Update expiry date and calculate RemainLife if provided
	if req.ExpiryDate != "" {
		expiryDate, err := time.Parse("2006-01-02", req.ExpiryDate)
		if err == nil {
			// Calculate RemainLife from expiry date
			now := time.Now()
			daysInYear := 365.25
			daysRemaining := expiryDate.Sub(now).Hours() / 24
			remainLife := daysRemaining / daysInYear
			equipment.RemainLife = remainLife

			// Update ReplacementYear based on expiry date
			replacementYear := expiryDate.Year()
			equipment.ReplacementYear = &replacementYear

			// Update LifeExpectancy based on new expiry date and receive date
			if equipment.ReceiveDate != nil {
				newLifeExpectancy := expiryDate.Sub(*equipment.ReceiveDate).Hours() / (24 * daysInYear)
				equipment.LifeExpectancy = newLifeExpectancy
			}

			// Auto-update status based on new RemainLife
			equipment.Status = u.calculateStatus(remainLife)
		}
	}

	// Save updated equipment via service
	if err := u.equipmentService.UpdateEquipment(ctx, equipment); err != nil {
		log.Printf("Usecase: UpdateEquipment - Error: %v", err)
		return err
	}

	log.Printf("Usecase: UpdateEquipment - Equipment ID: %s updated successfully", idCode)
	return nil
}

// DeleteEquipment soft deletes equipment by ID code
func (u *equipmentUsecase) DeleteEquipment(ctx context.Context, idCode string) error {
	// Find existing equipment
	equipment, err := u.equipmentService.FindEquipmentByIDCode(ctx, idCode)
	if err != nil {
		return err
	}
	if equipment == nil {
		return errors.New("equipment not found")
	}

	// Delete equipment via service
	if err := u.equipmentService.DeleteEquipment(ctx, equipment.ID); err != nil {
		log.Printf("Usecase: DeleteEquipment - Error: %v", err)
		return err
	}

	log.Printf("Usecase: DeleteEquipment - Equipment ID: %s deleted successfully", idCode)
	return nil
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
		IDCode:         req.IDCode,
		SerialNo:       &req.SerialNo,
		ModelID:        model.ID,
		DepartmentID:   department.ID,
		ReceiveDate:    &receiveDate,
		PurchasePrice:  req.PurchasePrice,
		LifeExpectancy: req.LifeExpectancy,
	}

	// Set optional fields
	if req.ECRICode != "" {
		equipment.ECRICode = &req.ECRICode
	}
	if req.AssetTypeName != "" {
		equipment.AssetTypeName = &req.AssetTypeName
	}
	if req.AssetName != "" {
		equipment.AssetName = &req.AssetName
	}
	if req.Building != "" {
		equipment.Building = &req.Building
	}
	if req.Floor != "" {
		equipment.Floor = &req.Floor
	}
	if req.Room != "" {
		equipment.Room = &req.Room
	}
	if req.Remark != "" {
		equipment.Remark = &req.Remark
	}
	if req.PurchaseDate != "" {
		purchaseDate, err := time.Parse("2006-01-02", req.PurchaseDate)
		if err == nil {
			equipment.PurchaseDate = &purchaseDate
		}
	}
	if req.WarrantyStartDate != "" {
		wd, err := time.Parse("2006-01-02", req.WarrantyStartDate)
		if err == nil {
			equipment.WarrantyStartDate = &wd
		}
	}
	if req.WarrantyEndDate != "" {
		wd, err := time.Parse("2006-01-02", req.WarrantyEndDate)
		if err == nil {
			equipment.WarrantyEndDate = &wd
		}
	}
	if req.WarrantyPeriod != "" {
		equipment.WarrantyPeriod = &req.WarrantyPeriod
	}

	// 9. ✅ Compute lifecycle fields (EquipmentAge, RemainLife, ReplacementYear)
	u.mapper.ComputeLifecycleFieldsPublic(equipment)

	// 10. Calculate and set Status based on RemainLife
	equipment.Status = u.calculateStatus(equipment.RemainLife)

	// 11. Save to database via service
	err = u.equipmentService.CreateEquipment(ctx, equipment)
	if err != nil {
		log.Printf("Usecase: Error creating equipment: %v", err)
		return nil, err
	}

	// 12. Load relations and return
	createdEquipment, err := u.equipmentService.FindEquipmentByID(ctx, equipment.ID)
	if err != nil {
		log.Printf("Usecase: Error loading created equipment: %v", err)
		return nil, err
	}

	log.Printf("Usecase: Successfully created equipment: %s (ID: %d)", equipment.IDCode, equipment.ID)

	// 13. Map to response DTO using mapper
	return u.mapper.MapEquipmentToResponse(createdEquipment), nil
}

// ===== HELPER FUNCTIONS =====

func (u *equipmentUsecase) mapEquipmentToDetailResponse(e *entity.Equipment) *dto.EquipmentDetailResponse {
	item := u.mapper.MapEquipmentToListItem(e)

	serialNo := ""
	if e.SerialNo != nil {
		serialNo = *e.SerialNo
	}

	brand := ""
	if e.Model.Brand.Name != "" {
		brand = e.Model.Brand.Name
	}

	return &dto.EquipmentDetailResponse{
		ID:           item.ID,
		Name:         item.Name,
		Category:     item.Category,
		Status:       item.Status,
		Location:     item.Location,
		LastCheck:    item.LastCheck,
		Expiry:       item.Expiry,
		IsExpiring:   item.IsExpiring,
		SerialNo:     serialNo,
		Brand:        brand,
		DepartmentID: e.DepartmentID,
	}
}

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

	return nil
}

// calculateStatus determines the asset status based on remain_life
func (u *equipmentUsecase) calculateStatus(remainLife float64) entity.AssetStatus {
	if remainLife <= 0 {
		return entity.AssetStatusPlanToReplace // หมดอายุแล้ว → รอเปลี่ยนใหม่
	}
	// ใกล้หมดอายุ (0 < remainLife <= 1) หรือปกติ → ยังใช้งานอยู่
	return entity.AssetStatusActive
}
