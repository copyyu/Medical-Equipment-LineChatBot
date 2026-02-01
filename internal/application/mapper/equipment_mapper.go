// domain/line/dto/equipment_mapper.go
package mapper

import (
	"medical-webhook/internal/application/dto"
	"medical-webhook/internal/domain/line/entity"
	"strconv"
)

// EquipmentMapper - Mapper สำหรับแปลงระหว่าง Entity และ DTO
type EquipmentMapper struct{}

func NewEquipmentMapper() *EquipmentMapper {
	return &EquipmentMapper{}
}

// ToEquipmentEntity - แปลง CreateEquipmentDTO เป็น Equipment Entity
func (m *EquipmentMapper) ToEquipmentEntity(dto *dto.CreateEquipmentDTO) *entity.Equipment {
	return &entity.Equipment{
		IDCode:                dto.IDCode,
		SerialNo:              dto.SerialNo,
		ModelID:               dto.ModelID,
		DepartmentID:          dto.DepartmentID,
		AssessmentID:          dto.AssessmentID,
		ReceiveDate:           dto.ReceiveDate,
		PurchasePrice:         dto.PurchasePrice,
		EquipmentAge:          dto.EquipmentAge,
		ComputeDate:           dto.ComputeDate,
		LifeExpectancy:        dto.LifeExpectancy,
		RemainLife:            dto.RemainLife,
		UsefulLifetimePercent: dto.UsefulLifetimePercent,
		ReplacementYear:       dto.ReplacementYear,
		Technology:            dto.Technology,
		UsageStatistics:       dto.UsageStatistics,
		Efficiency:            dto.Efficiency,
		Others:                dto.Others,
	}
}

// ToCreateEquipmentDTO - แปลง ExcelRowDTO เป็น CreateEquipmentDTO
func (m *EquipmentMapper) ToCreateEquipmentDTO(
	excelRow *dto.ExcelRowDTO,
	modelID uint,
	departmentID uint,
) *dto.CreateEquipmentDTO {
	var assessmentID *string
	if excelRow.AssessmentID != "" {
		assessmentID = &excelRow.AssessmentID
	}

	return &dto.CreateEquipmentDTO{
		IDCode:                excelRow.IDCode,
		SerialNo:              excelRow.SerialNo,
		ModelID:               modelID,
		DepartmentID:          departmentID,
		AssessmentID:          assessmentID,
		ReceiveDate:           excelRow.ReceiveDate,
		PurchasePrice:         excelRow.PurchasePrice,
		EquipmentAge:          excelRow.EquipmentAge,
		ComputeDate:           excelRow.ComputeDate,
		LifeExpectancy:        excelRow.LifeExpectancy,
		RemainLife:            excelRow.RemainLife,
		UsefulLifetimePercent: excelRow.UsefulLifePercent,
		ReplacementYear:       excelRow.ReplacementYear,
		Technology:            excelRow.Technology,
		UsageStatistics:       excelRow.UsageStatistics,
		Efficiency:            excelRow.Efficiency,
		Others:                excelRow.Others,
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
	// Get name from model
	name := ""
	if entity.Model.ModelName != "" {
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

	// Calculate expiry date based on replacement year
	expiry := ""
	isExpiring := false
	if entity.ReplacementYear != nil {
		expiry = formatYear(*entity.ReplacementYear)
		// Check if expiring within 1 year
		if entity.RemainLife <= 1 {
			isExpiring = true
		}
	} else if entity.RemainLife > 0 {
		// Calculate based on remain life
		currentYear := 2026 // Current year
		expiryYear := currentYear + int(entity.RemainLife)
		expiry = formatYear(expiryYear)
		if entity.RemainLife <= 1 {
			isExpiring = true
		}
	}

	// Get last check date from latest maintenance record
	lastCheck := ""
	if len(entity.MaintenanceRecords) > 0 {
		lastCheck = entity.MaintenanceRecords[0].MaintenanceDate.Format("2006-01-02")
	} else if entity.ComputeDate != nil {
		lastCheck = entity.ComputeDate.Format("2006-01-02")
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
	}
}

// mapAssetStatusToFrontend maps backend AssetStatus to frontend EquipmentStatus
func mapAssetStatusToFrontend(status entity.AssetStatus) string {
	switch status {
	case entity.AssetStatusActive:
		return "ready"
	case entity.AssetStatusDefective:
		return "broken"
	case entity.AssetStatusWaitDecom, entity.AssetStatusDecommission:
		return "maintenance"
	case entity.AssetStatusMissing:
		return "broken"
	case entity.AssetStatusPlanToReplace:
		return "expired"
	case entity.AssetStatusActiveReadyToSell:
		return "in_use"
	default:
		return "ready"
	}
}

func formatYear(year int) string {
	return strconv.Itoa(year)
}

// func formatYear(year int) string {
//     return string(rune('0'+year/1000)) + string(rune('0'+(year/100)%10)) + string(rune('0'+(year/10)%10)) + string(rune('0'+year%10))
// }
//

func (m *EquipmentMapper) MapEquipmentToResponse(entity *entity.Equipment) *dto.EquipmentResponse {
	resp := &dto.EquipmentResponse{
		ID:                    entity.ID,
		IDCode:                entity.IDCode,
		SerialNo:              entity.SerialNo,
		AssessmentID:          entity.AssessmentID,
		Status:                string(entity.Status),
		ReceiveDate:           entity.ReceiveDate,
		PurchasePrice:         entity.PurchasePrice,
		EquipmentAge:          entity.EquipmentAge,
		ComputeDate:           entity.ComputeDate,
		LifeExpectancy:        entity.LifeExpectancy,
		RemainLife:            entity.RemainLife,
		UsefulLifetimePercent: entity.UsefulLifetimePercent,
		ReplacementYear:       entity.ReplacementYear,
		Technology:            entity.Technology,
		UsageStatistics:       entity.UsageStatistics,
		Efficiency:            entity.Efficiency,
		Others:                entity.Others,
		CreatedAt:             entity.CreatedAt,
		UpdatedAt:             entity.UpdatedAt,
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
