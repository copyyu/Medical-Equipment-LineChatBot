// internal/interfaces/http/handlers/equipment_import_handler.go
package handlers

import (
	"log"
	"medical-webhook/internal/application/dto"
	"medical-webhook/internal/application/usecase"

	"github.com/gofiber/fiber/v2"
)

type EquipmentImportHandler struct {
	importUseCase usecase.EquipmentImportUseCase
}

func NewEquipmentImportHandler(importUseCase usecase.EquipmentImportUseCase) *EquipmentImportHandler {
	return &EquipmentImportHandler{
		importUseCase: importUseCase,
	}
}

// ImportExcel handles the Excel file upload and import
func (h *EquipmentImportHandler) ImportExcel(c *fiber.Ctx) error {
	log.Println("Received equipment import request")

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("Failed to get file: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to get file. Please upload a file with key 'file'",
		})
	}

	log.Printf("Received file: %s (%.2f KB)", file.Filename, float64(file.Size)/1024)

	// Validate file extension
	if !isValidExcelFile(file.Filename) {
		log.Printf("Invalid file type: %s", file.Filename)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid file type. Only .xlsx and .xls files are allowed",
		})
	}

	// Open the file
	fileReader, err := file.Open()
	if err != nil {
		log.Printf("Failed to open file: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to open file: " + err.Error(),
		})
	}
	defer fileReader.Close()

	// Execute UseCase
	log.Println("Starting Excel import process...")
	result, err := h.importUseCase.Execute(c.Context(), fileReader)
	if err != nil {
		log.Printf("Import failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to import Excel: " + err.Error(),
		})
	}

	// Build success message
	message := buildSuccessMessage(result)
	log.Printf("Import completed: %s", message)
	log.Printf("Stats: Total=%d, Success=%d, Failed=%d, Skipped=%d",
		result.TotalRows, result.SuccessCount, result.FailedCount, result.SkippedCount)
	log.Printf("New records: Brands=%d, Categories=%d, Departments=%d, Models=%d",
		result.NewBrands, result.NewCategories, result.NewDepartments, result.NewModels)

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    result,
	})
}

// ImportExcelBatch handles multiple file uploads
func (h *EquipmentImportHandler) ImportExcelBatch(c *fiber.Ctx) error {
	log.Println("📥 Received batch equipment import request")

	// Get multipart form
	form, err := c.MultipartForm()
	if err != nil {
		log.Printf("Failed to parse multipart form: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to parse form: " + err.Error(),
		})
	}

	// Get all files
	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "No files provided. Please upload files with key 'files'",
		})
	}

	log.Printf("📄 Received %d files for batch import", len(files))

	var results []fiber.Map
	totalSuccess := 0
	totalFailed := 0

	// Process each file
	for i, fileHeader := range files {
		log.Printf("📄 Processing file %d/%d: %s", i+1, len(files), fileHeader.Filename)

		// Validate file extension
		if !isValidExcelFile(fileHeader.Filename) {
			log.Printf("Invalid file type: %s", fileHeader.Filename)
			results = append(results, fiber.Map{
				"filename": fileHeader.Filename,
				"success":  false,
				"error":    "Invalid file type",
			})
			continue
		}

		// Open file
		file, err := fileHeader.Open()
		if err != nil {
			log.Printf("Failed to open file %s: %v", fileHeader.Filename, err)
			results = append(results, fiber.Map{
				"filename": fileHeader.Filename,
				"success":  false,
				"error":    "Failed to open file",
			})
			continue
		}

		// Execute import
		result, err := h.importUseCase.Execute(c.Context(), file)
		file.Close()

		if err != nil {
			log.Printf("Import failed for %s: %v", fileHeader.Filename, err)
			results = append(results, fiber.Map{
				"filename": fileHeader.Filename,
				"success":  false,
				"error":    err.Error(),
			})
			continue
		}

		totalSuccess += result.SuccessCount
		totalFailed += result.FailedCount

		results = append(results, fiber.Map{
			"filename": fileHeader.Filename,
			"success":  true,
			"data":     result,
		})

		log.Printf("Completed %s: Success=%d, Failed=%d",
			fileHeader.Filename, result.SuccessCount, result.FailedCount)
	}

	log.Printf("Batch import completed: Total Success=%d, Total Failed=%d", totalSuccess, totalFailed)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":       true,
		"message":       "Batch import completed",
		"files_count":   len(files),
		"total_success": totalSuccess,
		"total_failed":  totalFailed,
		"results":       results,
	})
}

// GetImportHistory returns import history (optional - if you want to track imports)
func (h *EquipmentImportHandler) GetImportHistory(c *fiber.Ctx) error {
	// TODO: Implement if you want to track import history
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Import history endpoint - to be implemented",
		"data":    []interface{}{},
	})
}

// buildSuccessMessage builds appropriate message based on import result
func buildSuccessMessage(result *dto.EquipmentImportResultDTO) string {
	if result.FailedCount == 0 && result.SkippedCount == 0 {
		return "Excel file imported successfully! All rows processed."
	}
	if result.FailedCount > 0 && result.SkippedCount > 0 {
		return "Excel file imported with some errors and skipped rows. Please check the details."
	}
	if result.FailedCount > 0 {
		return "Excel file imported with some errors. Please check the error messages."
	}
	if result.SkippedCount > 0 {
		return "Excel file imported successfully. Some rows were skipped."
	}
	return "Excel file imported successfully!"
}

// isValidExcelFile validates the file extension
func isValidExcelFile(filename string) bool {
	if len(filename) < 4 {
		return false
	}
	// Check for .xlsx (5 chars)
	if len(filename) >= 5 && filename[len(filename)-5:] == ".xlsx" {
		return true
	}
	// Check for .xls (4 chars)
	if filename[len(filename)-4:] == ".xls" {
		return true
	}
	return false
}
