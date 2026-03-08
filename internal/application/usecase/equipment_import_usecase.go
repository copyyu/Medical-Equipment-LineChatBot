package usecase

import (
	"context"
	"fmt"
	"io"
	"log"
	"medical-webhook/internal/application/dto"
	"medical-webhook/internal/application/mapper"
	"medical-webhook/internal/application/service"
	"medical-webhook/internal/domain/line/repository"
)

// EquipmentImportUseCase - UseCase สำหรับ import equipment จาก Excel
type EquipmentImportUseCase interface {
	Execute(ctx context.Context, file io.Reader) (*dto.EquipmentImportResultDTO, error)
}

type equipmentImportUseCase struct {
	equipmentRepo     repository.EquipmentRepository
	excelParser       service.ExcelParserService
	masterDataService service.MasterDataService
	mapper            *mapper.EquipmentMapper
}

func NewEquipmentImportUseCase(
	equipmentRepo repository.EquipmentRepository,
	excelParser service.ExcelParserService,
	masterDataService service.MasterDataService,
	mapper *mapper.EquipmentMapper,
) EquipmentImportUseCase {
	return &equipmentImportUseCase{
		equipmentRepo:     equipmentRepo,
		excelParser:       excelParser,
		masterDataService: masterDataService,
		mapper:            mapper,
	}
}

// Execute - ดำเนินการ import equipment จาก Excel file
func (uc *equipmentImportUseCase) Execute(ctx context.Context, file io.Reader) (*dto.EquipmentImportResultDTO, error) {
	result := &dto.EquipmentImportResultDTO{
		FailedRows:    make([]int, 0),
		ErrorMessages: make([]string, 0),
	}

	// ⭐ Clear cache when done
	defer uc.masterDataService.ClearCache()

	// 1. Parse Excel file
	excelRows, err := uc.excelParser.ParseExcelFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse excel file: %w", err)
	}

	// ⭐ Log total rows
	log.Printf("📊 Starting import: %d rows to process", len(excelRows))

	// 2. Process each row
	for i, excelRow := range excelRows {
		rowNum := i + 2

		// ⭐ Progress reporting for large files
		if (i+1)%100 == 0 {
			log.Printf("📊 Progress: %d/%d rows processed", i+1, len(excelRows))
		}

		// Validate required fields
		if err := uc.validateExcelRow(excelRow); err != nil {
			result.FailedCount++
			result.FailedRows = append(result.FailedRows, rowNum)
			result.ErrorMessages = append(result.ErrorMessages, fmt.Sprintf("Row %d: %s", rowNum, err.Error()))
			continue
		}

		// Skip if no ID Code
		if excelRow.IDCode == "" {
			result.SkippedCount++
			continue
		}

		result.TotalRows++

		// 3. Get or Create Master Data
		department, isNewDept, err := uc.masterDataService.GetOrCreateDepartment(ctx, excelRow.Department)
		if err != nil {
			result.FailedCount++
			result.FailedRows = append(result.FailedRows, rowNum)
			result.ErrorMessages = append(result.ErrorMessages, fmt.Sprintf("Row %d: %s", rowNum, err.Error()))
			continue
		}
		if isNewDept {
			result.NewDepartments++
		}

		brand, isNewBrand, err := uc.masterDataService.GetOrCreateBrand(ctx, excelRow.Brand)
		if err != nil {
			result.FailedCount++
			result.FailedRows = append(result.FailedRows, rowNum)
			result.ErrorMessages = append(result.ErrorMessages, fmt.Sprintf("Row %d: %s", rowNum, err.Error()))
			continue
		}
		if isNewBrand {
			result.NewBrands++
		}

		category, isNewCategory, err := uc.masterDataService.GetOrCreateCategory(
			ctx,
			excelRow.Category,
			excelRow.ECRIRisk,
			excelRow.Classification,
		)
		if err != nil {
			result.FailedCount++
			result.FailedRows = append(result.FailedRows, rowNum)
			result.ErrorMessages = append(result.ErrorMessages, fmt.Sprintf("Row %d: %s", rowNum, err.Error()))
			continue
		}
		if isNewCategory {
			result.NewCategories++
		}

		model, isNewModel, err := uc.masterDataService.GetOrCreateModel(
			ctx,
			brand.ID,
			category.ID,
			excelRow.Model,
			excelRow.LifeExpectancy,
		)
		if err != nil {
			result.FailedCount++
			result.FailedRows = append(result.FailedRows, rowNum)
			result.ErrorMessages = append(result.ErrorMessages, fmt.Sprintf("Row %d: %s", rowNum, err.Error()))
			continue
		}
		if isNewModel {
			result.NewModels++
		}

		// 4. Map to CreateEquipmentDTO
		createEquipmentDTO := uc.mapper.ToCreateEquipmentDTO(excelRow, model.ID, department.ID)

		// 5. Map to Entity
		equipment := uc.mapper.ToEquipmentEntity(createEquipmentDTO)

		// 6. Save Equipment (CreateOrUpdate)
		if err := uc.equipmentRepo.CreateOrUpdate(ctx, equipment); err != nil {
			result.FailedCount++
			result.FailedRows = append(result.FailedRows, rowNum)
			result.ErrorMessages = append(result.ErrorMessages, fmt.Sprintf("Row %d: failed to save equipment: %s", rowNum, err.Error()))
			continue
		}

		result.SuccessCount++
	}

	// ⭐ Log summary
	log.Printf("🎉 Import completed: Success=%d, Failed=%d, Skipped=%d",
		result.SuccessCount, result.FailedCount, result.SkippedCount)

	return result, nil
}

// validateExcelRow
func (uc *equipmentImportUseCase) validateExcelRow(row *dto.ExcelRowDTO) error {
	if row.IDCode == "" {
		return fmt.Errorf("ID CODE is required")
	}

	if row.Department == "" {
		row.Department = "-"
	}
	if row.Brand == "" {
		row.Brand = "-"
	}
	if row.Category == "" {
		row.Category = "-"
	}
	if row.Model == "" {
		row.Model = "-"
	}

	return nil
}
