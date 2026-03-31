package mapper

import (
	"medical-webhook/internal/application/dto"
	"medical-webhook/internal/domain/line/entity"
	"strconv"
	"time"
)

// EquipmentMapper - Mapper สำหรับแปลงระหว่าง Entity และ DTO
type EquipmentMapper struct{}

func NewEquipmentMapper() *EquipmentMapper {
	return &EquipmentMapper{}
}

// ToEquipmentEntity - แปลง CreateEquipmentDTO เป็น Equipment Entity
func (m *EquipmentMapper) ToEquipmentEntity(dto *dto.CreateEquipmentDTO) *entity.Equipment {
	eq := &entity.Equipment{
		IDCode:       dto.IDCode,
		SerialNo:     dto.SerialNo,
		ModelID:      dto.ModelID,
		DepartmentID: dto.DepartmentID,

		// Basic Info
		AssetTypeName: dto.AssetTypeName,
		ECRICode:      dto.ECRICode,
		AssetName:     dto.AssetName,
		AssetID:       dto.AssetID,

		// Status
		AssetStatusInternal: dto.AssetStatusInternal,
		RentalStatus:        dto.RentalStatus,
		BorrowStatus:        dto.BorrowStatus,

		// Location
		Building: dto.Building,
		Floor:    dto.Floor,
		Room:     dto.Room,
		PhoneNo:  dto.PhoneNo,

		// Business
		BusinessName: dto.BusinessName,
		ItemNo:       dto.ItemNo,
		SKUNo:        dto.SKUNo,

		// Dates
		ReceiveDate:      dto.ReceiveDate,
		PurchaseDate:     dto.PurchaseDate,
		RegistrationDate: dto.RegistrationDate,
		PurchasePrice:    dto.PurchasePrice,

		// Lifecycle
		LifeExpectancy: dto.LifeExpectancy,

		// Warranty
		WarrantyPeriod:    dto.WarrantyPeriod,
		WarrantyStartDate: dto.WarrantyStartDate,
		WarrantyEndDate:   dto.WarrantyEndDate,
		WarrantyPM:        dto.WarrantyPM,
		WarrantyCal:       dto.WarrantyCal,

		// PM & Calibration
		LastPMDate:  dto.LastPMDate,
		LastCalDate: dto.LastCalDate,
		PMPeriod:    dto.PMPeriod,
		CalPeriod:   dto.CalPeriod,
		VendorPM:    dto.VendorPM,
		VendorCal:   dto.VendorCal,

		// Power
		PowerConsumption: dto.PowerConsumption,

		// Procurement
		Supplier:             dto.Supplier,
		Ownership:            dto.Ownership,
		PoNo:                 dto.PoNo,
		ContractNo:           dto.ContractNo,
		InvoiceNo:            dto.InvoiceNo,
		DocumentNo:           dto.DocumentNo,
		TorNo:                dto.TorNo,
		ManufacturingCountry: dto.ManufacturingCountry,

		// Financial
		RevenuePerMonth: dto.RevenuePerMonth,

		// Misc
		Remark:         dto.Remark,
		ApprovedBy:     dto.ApprovedBy,
		NsmartItemCode: dto.NsmartItemCode,
		UpdatedBy:      dto.UpdatedBy,
	}

	// Compute lifecycle fields
	m.computeLifecycleFields(eq)

	return eq
}

// computeLifecycleFields - คำนวณ EquipmentAge, RemainLife, ReplacementYear จาก ReceiveDate + LifeExpectancy
func (m *EquipmentMapper) computeLifecycleFields(eq *entity.Equipment) {
	if eq.ReceiveDate != nil && eq.LifeExpectancy > 0 {
		now := time.Now()
		// EquipmentAge = (today - receive_date) in fractional years
		equipmentAge := now.Sub(*eq.ReceiveDate).Hours() / (24 * 365.25)
		eq.EquipmentAge = equipmentAge

		// RemainLife = LifeExpectancy - EquipmentAge
		eq.RemainLife = eq.LifeExpectancy - equipmentAge

		// ReplacementYear = ReceiveDate.Year + LifeExpectancy
		replacementYear := eq.ReceiveDate.Year() + int(eq.LifeExpectancy)
		eq.ReplacementYear = &replacementYear
	}
}

// ComputeLifecycleFieldsPublic - public wrapper for computeLifecycleFields
func (m *EquipmentMapper) ComputeLifecycleFieldsPublic(eq *entity.Equipment) {
	m.computeLifecycleFields(eq)
}

// ToCreateEquipmentDTO - แปลง ExcelRowDTO เป็น CreateEquipmentDTO
func (m *EquipmentMapper) ToCreateEquipmentDTO(
	excelRow *dto.ExcelRowDTO,
	modelID uint,
	departmentID uint,
) *dto.CreateEquipmentDTO {
	return &dto.CreateEquipmentDTO{
		IDCode:       excelRow.IDCode,
		SerialNo:     excelRow.SerialNo,
		ModelID:      modelID,
		DepartmentID: departmentID,

		// Basic Info
		AssetTypeName: strPtr(excelRow.AssetTypeName),
		ECRICode:      strPtr(excelRow.ECRICode),
		AssetName:     excelRow.AssetName,
		AssetID:       excelRow.AssetID,

		// Status
		AssetStatus:         excelRow.AssetStatus,
		AssetStatusInternal: excelRow.AssetStatusInternal,
		RentalStatus:        excelRow.RentalStatus,
		BorrowStatus:        excelRow.BorrowStatus,

		// Location
		Building: excelRow.Building,
		Floor:    excelRow.Floor,
		Room:     excelRow.Room,
		PhoneNo:  excelRow.PhoneNo,

		// Business
		BusinessName: excelRow.BusinessName,
		ItemNo:       excelRow.ItemNo,
		SKUNo:        excelRow.SKUNo,

		// Dates
		ReceiveDate:      excelRow.ReceiveDate,
		PurchaseDate:     excelRow.PurchaseDate,
		RegistrationDate: excelRow.RegistrationDate,
		PurchasePrice:    excelRow.PurchasePrice,

		// Lifecycle
		LifeExpectancy: excelRow.LifeExpectancy,

		// Warranty
		WarrantyPeriod:    excelRow.WarrantyPeriod,
		WarrantyStartDate: excelRow.WarrantyStartDate,
		WarrantyEndDate:   excelRow.WarrantyEndDate,
		WarrantyPM:        excelRow.WarrantyPM,
		WarrantyCal:       excelRow.WarrantyCal,

		// PM & Calibration
		LastPMDate:  excelRow.LastPMDate,
		LastCalDate: excelRow.LastCalDate,
		PMPeriod:    excelRow.PMPeriod,
		CalPeriod:   excelRow.CalPeriod,
		VendorPM:    excelRow.VendorPM,
		VendorCal:   excelRow.VendorCal,

		// Power
		PowerConsumption: excelRow.PowerConsumption,

		// Procurement
		Supplier:             excelRow.Supplier,
		Ownership:            excelRow.Ownership,
		PoNo:                 excelRow.PoNo,
		ContractNo:           excelRow.ContractNo,
		InvoiceNo:            excelRow.InvoiceNo,
		DocumentNo:           excelRow.DocumentNo,
		TorNo:                excelRow.TorNo,
		ManufacturingCountry: excelRow.ManufacturingCountry,

		// Financial
		RevenuePerMonth: excelRow.RevenuePerMonth,

		// Misc
		Remark:         excelRow.Remark,
		ApprovedBy:     excelRow.ApprovedBy,
		NsmartItemCode: excelRow.NsmartItemCode,
		UpdatedBy:      excelRow.UpdatedBy,
	}
}

// ToBrandEntity - แปลง name เป็น Brand Entity
func (m *EquipmentMapper) ToBrandEntity(name string) *entity.Brand {
	return &entity.Brand{
		Name: name,
	}
}

// ToCategoryEntity - แปลง data เป็น EquipmentCategory Entity
func (m *EquipmentMapper) ToCategoryEntity(name, ecriRisk, classification string) *entity.EquipmentCategory {
	var riskLevel entity.ECRIRiskLevel
	switch ecriRisk {
	case "HIGH":
		riskLevel = entity.RiskHigh
	case "LOW":
		riskLevel = entity.RiskLow
	default:
		riskLevel = entity.RiskMedium
	}

	return &entity.EquipmentCategory{
		Name:           name,
		ECRIRisk:       riskLevel,
		Classification: classification,
	}
}

// ToDepartmentEntity - แปลง name เป็น Department Entity
func (m *EquipmentMapper) ToDepartmentEntity(name string) *entity.Department {
	return &entity.Department{
		Name: name,
	}
}

// ToModelEntity - แปลง data เป็น EquipmentModel Entity
func (m *EquipmentMapper) ToModelEntity(
	brandID uint,
	categoryID uint,
	modelName string,
	lifeExpectancy float64,
) *entity.EquipmentModel {
	if lifeExpectancy == 0 {
		lifeExpectancy = 10 // default
	}

	return &entity.EquipmentModel{
		BrandID:               brandID,
		CategoryID:            categoryID,
		ModelName:             modelName,
		DefaultLifeExpectancy: lifeExpectancy,
	}
}

// ToBrandDTO - แปลง Brand Entity เป็น DTO
func (m *EquipmentMapper) ToBrandDTO(entity *entity.Brand) *dto.BrandDTO {
	return &dto.BrandDTO{
		ID:   entity.ID,
		Name: entity.Name,
	}
}

// ToCategoryDTO - แปลง EquipmentCategory Entity เป็น DTO
func (m *EquipmentMapper) ToCategoryDTO(entity *entity.EquipmentCategory) *dto.CategoryDTO {
	return &dto.CategoryDTO{
		ID:             entity.ID,
		Name:           entity.Name,
		ECRIRisk:       string(entity.ECRIRisk),
		Classification: entity.Classification,
	}
}

// ToDepartmentDTO - แปลง Department Entity เป็น DTO
func (m *EquipmentMapper) ToDepartmentDTO(entity *entity.Department) *dto.DepartmentDTO {
	return &dto.DepartmentDTO{
		ID:   entity.ID,
		Name: entity.Name,
	}
}

// ToModelDTO - แปลง EquipmentModel Entity เป็น DTO
func (m *EquipmentMapper) ToModelDTO(entity *entity.EquipmentModel) *dto.ModelDTO {
	return &dto.ModelDTO{
		ID:                    entity.ID,
		BrandID:               entity.BrandID,
		CategoryID:            entity.CategoryID,
		ModelName:             entity.ModelName,
		DefaultLifeExpectancy: entity.DefaultLifeExpectancy,
	}
}

func (m *EquipmentMapper) MapEquipmentToListItem(entity *entity.Equipment) *dto.EquipmentListItem {
	// Get name: prefer AssetName, then Model name, then IDCode
	name := ""
	if entity.AssetName != nil && *entity.AssetName != "" {
		name = *entity.AssetName
	} else if entity.Model.ModelName != "" {
		name = entity.Model.ModelName
	} else {
		name = entity.IDCode
	}

	// Get category from model
	category := ""
	if entity.Model.Category.Name != "" {
		category = entity.Model.Category.Name
	}

	// Get location from department
	location := ""
	if entity.Department.Name != "" {
		location = entity.Department.Name
	}

	// Calculate expiry and remain_life dynamically based on current date
	expiry := ""
	isExpiring := false
	now := time.Now()
	currentYear := now.Year()

	// Dynamic remain_life calculation: use today's date
	dynamicRemainLife := entity.RemainLife // fallback to stored value

	if entity.ReceiveDate != nil && entity.LifeExpectancy > 0 {
		// Primary: equipment_age = (today - receive_date) in fractional years
		equipmentAge := now.Sub(*entity.ReceiveDate).Hours() / (24 * 365.25)
		dynamicRemainLife = entity.LifeExpectancy - equipmentAge
	}

	// Calculate expiry year for display
	if entity.ReplacementYear != nil {
		expiry = formatYear(*entity.ReplacementYear)
	} else if entity.ReceiveDate != nil && entity.LifeExpectancy > 0 {
		expiryYear := entity.ReceiveDate.Year() + int(entity.LifeExpectancy)
		expiry = formatYear(expiryYear)
	} else if entity.RemainLife > 0 {
		expiryYear := currentYear + int(entity.RemainLife)
		expiry = formatYear(expiryYear)
	}

	// Set isExpiring based on dynamic remain_life
	if dynamicRemainLife <= 1 {
		isExpiring = true
	}

	// Get last check date from latest maintenance record or LastPMDate
	lastCheck := ""
	if len(entity.MaintenanceRecords) > 0 {
		lastCheck = entity.MaintenanceRecords[0].MaintenanceDate.Format("2006-01-02")
	} else if entity.LastPMDate != nil {
		lastCheck = entity.LastPMDate.Format("2006-01-02")
	}

	// Map asset status to frontend status
	status := mapAssetStatusToFrontend(entity.Status)

	return &dto.EquipmentListItem{
		ID:         entity.IDCode,
		Name:       name,
		Category:   category,
		Status:     status,
		Location:   location,
		LastCheck:  lastCheck,
		Expiry:     expiry,
		IsExpiring: isExpiring,
		RemainLife: dynamicRemainLife,
	}
}

// mapAssetStatusToFrontend maps backend AssetStatus to frontend EquipmentStatus
func mapAssetStatusToFrontend(status entity.AssetStatus) string {
	return string(status)
}

func formatYear(year int) string {
	return strconv.Itoa(year)
}

func (m *EquipmentMapper) MapEquipmentToResponse(entity *entity.Equipment) *dto.EquipmentResponse {
	resp := &dto.EquipmentResponse{
		ID:                   entity.ID,
		IDCode:               entity.IDCode,
		SerialNo:             entity.SerialNo,
		ECRICode:             entity.ECRICode,
		Status:               string(entity.Status),
		ReceiveDate:          entity.ReceiveDate,
		PurchasePrice:        entity.PurchasePrice,
		EquipmentAge:         entity.EquipmentAge,
		LifeExpectancy:       entity.LifeExpectancy,
		RemainLife:           entity.RemainLife,
		ReplacementYear:      entity.ReplacementYear,
		AssetTypeName:        entity.AssetTypeName,
		AssetName:            entity.AssetName,
		AssetID:              entity.AssetID,
		AssetStatusInternal:  entity.AssetStatusInternal,
		RentalStatus:         entity.RentalStatus,
		BorrowStatus:         entity.BorrowStatus,
		Building:             entity.Building,
		Floor:                entity.Floor,
		Room:                 entity.Room,
		PhoneNo:              entity.PhoneNo,
		BusinessName:         entity.BusinessName,
		ItemNo:               entity.ItemNo,
		SKUNo:                entity.SKUNo,
		PurchaseDate:         entity.PurchaseDate,
		RegistrationDate:     entity.RegistrationDate,
		WarrantyPeriod:       entity.WarrantyPeriod,
		WarrantyStartDate:    entity.WarrantyStartDate,
		WarrantyEndDate:      entity.WarrantyEndDate,
		WarrantyPM:           entity.WarrantyPM,
		WarrantyCal:          entity.WarrantyCal,
		LastPMDate:           entity.LastPMDate,
		LastCalDate:          entity.LastCalDate,
		PMPeriod:             entity.PMPeriod,
		CalPeriod:            entity.CalPeriod,
		VendorPM:             entity.VendorPM,
		VendorCal:            entity.VendorCal,
		PowerConsumption:     entity.PowerConsumption,
		Supplier:             entity.Supplier,
		Ownership:            entity.Ownership,
		PoNo:                 entity.PoNo,
		ContractNo:           entity.ContractNo,
		InvoiceNo:            entity.InvoiceNo,
		DocumentNo:           entity.DocumentNo,
		TorNo:                entity.TorNo,
		ManufacturingCountry: entity.ManufacturingCountry,
		RevenuePerMonth:      entity.RevenuePerMonth,
		Remark:               entity.Remark,
		ApprovedBy:           entity.ApprovedBy,
		NsmartItemCode:       entity.NsmartItemCode,
		UpdatedBy:            entity.UpdatedBy,

		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}

	// Map Model if exists
	if entity.Model.ID != 0 {
		resp.Model = &dto.EquipmentModelDTO{
			ID:                    entity.Model.ID,
			ModelName:             entity.Model.ModelName,
			DefaultLifeExpectancy: entity.Model.DefaultLifeExpectancy,
		}

		// Map Brand if exists
		if entity.Model.Brand.ID != 0 {
			resp.Model.Brand = &dto.BrandDTO{
				ID:   entity.Model.Brand.ID,
				Name: entity.Model.Brand.Name,
			}
		}

		// Map Category if exists
		if entity.Model.Category.ID != 0 {
			resp.Model.Category = &dto.CategoryDTO{
				ID:             entity.Model.Category.ID,
				Name:           entity.Model.Category.Name,
				ECRIRisk:       string(entity.Model.Category.ECRIRisk),
				Classification: entity.Model.Category.Classification,
			}
		}
	}

	// Map Department if exists
	if entity.Department.ID != 0 {
		resp.Department = &dto.DepartmentDTO{
			ID:   entity.Department.ID,
			Name: entity.Department.Name,
		}
	}

	return resp
}

// strPtr - helper to convert non-empty string to *string
func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
