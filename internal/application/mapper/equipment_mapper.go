// domain/line/dto/equipment_mapper.go
package mapper

import (
	"medical-webhook/internal/application/dto"
	"medical-webhook/internal/domain/line/entity"
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
