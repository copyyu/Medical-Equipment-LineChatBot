// internal/domain/line/service/excel_parser_service.go
package service

import (
	"fmt"
	"io"
	"medical-webhook/internal/application/dto"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// ExcelParserService - Service สำหรับ parse Excel file
type ExcelParserService interface {
	ParseExcelFile(file io.Reader) ([]*dto.ExcelRowDTO, error)
	ParseExcelRow(row []string, rowNum int) (*dto.ExcelRowDTO, error)
}

type excelParserService struct{}

func NewExcelParserService() ExcelParserService {
	return &excelParserService{}
}

// ParseExcelFile - อ่านและ parse ไฟล์ Excel ทั้งหมด
func (s *excelParserService) ParseExcelFile(file io.Reader) ([]*dto.ExcelRowDTO, error) {
	// Open Excel file
	xlFile, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open excel file: %w", err)
	}
	defer xlFile.Close()

	// Get first sheet
	sheets := xlFile.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("no sheets found in excel file")
	}

	// Read all rows
	rows, err := xlFile.GetRows(sheets[0])
	if err != nil {
		return nil, fmt.Errorf("failed to read rows: %w", err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("excel file is empty or has no data rows")
	}

	// Parse each row (skip header)
	result := make([]*dto.ExcelRowDTO, 0)
	for i, row := range rows[1:] {
		rowNum := i + 2 // +2 because skip header and 0-indexed

		// ⭐ Skip empty rows
		if isEmptyRow(row) {
			continue
		}

		rowData, err := s.ParseExcelRow(row, rowNum)
		if err != nil {
			// Skip invalid rows
			continue
		}

		result = append(result, rowData)
	}

	return result, nil
}

// isEmptyRow - ตรวจสอบว่า row นี้ว่างเปล่าหรือไม่
func isEmptyRow(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}

// ParseExcelRow - parse แต่ละ row จาก Excel
func (s *excelParserService) ParseExcelRow(row []string, rowNum int) (*dto.ExcelRowDTO, error) {
	// Helper function to safely get column value
	getCol := func(idx int) string {
		if idx < len(row) {
			return strings.TrimSpace(row[idx])
		}
		return ""
	}

	data := &dto.ExcelRowDTO{
		Department:     getCol(0),         // Department
		ECRIRisk:       getCol(1),         // ECRI Risk
		AssessmentID:   getCol(2),         // Assessment ID
		IDCode:         getCol(3),         // ID CODE
		Category:       getCol(4),         // Equipment Category
		Brand:          getCol(5),         // Brand
		Model:          getCol(6),         // Model
		SerialNo:       strPtr(getCol(7)), // Serial No
		Classification: getCol(8),         // Classification
	}

	if data.IDCode == "" {
		return nil, fmt.Errorf("ID CODE is empty")
	}

	// Receive Date (column 9)
	if dateStr := getCol(9); dateStr != "" {
		if date, err := s.parseDate(dateStr); err == nil {
			data.ReceiveDate = &date
		}
	}

	// Purchase Price (column 10)
	if priceStr := getCol(10); priceStr != "" {
		if price, err := strconv.ParseFloat(priceStr, 64); err == nil {
			data.PurchasePrice = price
		}
	}

	// Equipment Age (column 11)
	if ageStr := getCol(11); ageStr != "" {
		if age, err := strconv.ParseFloat(ageStr, 64); err == nil {
			data.EquipmentAge = age
		}
	}

	// Compute Date (column 12)
	if dateStr := getCol(12); dateStr != "" {
		if date, err := s.parseDate(dateStr); err == nil {
			data.ComputeDate = &date
		}
	}

	// Life Expectancy (column 13)
	if lifeStr := getCol(13); lifeStr != "" {
		if life, err := strconv.ParseFloat(lifeStr, 64); err == nil {
			data.LifeExpectancy = life
		}
	}

	// Remain Life (column 14)
	if remainStr := getCol(14); remainStr != "" {
		if remain, err := strconv.ParseFloat(remainStr, 64); err == nil {
			data.RemainLife = remain
		}
	}

	// Total of CM (column 15)
	if cmStr := getCol(15); cmStr != "" {
		if cm, err := strconv.Atoi(cmStr); err == nil {
			data.TotalOfCM = cm
		}
	}

	// Total of Cost (column 16)
	if costStr := getCol(16); costStr != "" {
		if cost, err := strconv.ParseFloat(costStr, 64); err == nil {
			data.TotalOfCost = cost
		}
	}

	// Per Cost Price (column 17)
	if perCostStr := getCol(17); perCostStr != "" {
		if perCost, err := strconv.ParseFloat(perCostStr, 64); err == nil {
			data.PerCostPrice = perCost
		}
	}

	// % of useful lifetime (column 18)
	if percentStr := getCol(18); percentStr != "" {
		// Remove % sign if present
		percentStr = strings.TrimSuffix(percentStr, "%")
		if percent, err := strconv.ParseFloat(percentStr, 64); err == nil {
			data.UsefulLifePercent = percent
		}
	}

	// Replacement Year (column 19)
	if yearStr := getCol(19); yearStr != "" {
		if year, err := strconv.Atoi(yearStr); err == nil {
			data.ReplacementYear = &year
		}
	}

	// Technology score (column 20)
	if techStr := getCol(20); techStr != "" {
		if tech, err := strconv.ParseFloat(techStr, 64); err == nil {
			data.Technology = &tech
		}
	}

	// Usage Statistics score (column 21)
	if usageStr := getCol(21); usageStr != "" {
		if usage, err := strconv.ParseFloat(usageStr, 64); err == nil {
			data.UsageStatistics = &usage
		}
	}

	// Efficiency score (column 22)
	if effStr := getCol(22); effStr != "" {
		if eff, err := strconv.ParseFloat(effStr, 64); err == nil {
			data.Efficiency = &eff
		}
	}

	// Others (column 23)
	if others := getCol(23); others != "" {
		data.Others = &others
	}

	return data, nil
}

// parseDate - parse วันที่จากหลายรูปแบบ
func (s *excelParserService) parseDate(dateStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02",
		"02/01/2006",
		"01/02/2006",
		"2006/01/02",
		"02-01-2006",
		"01-02-2006",
		"2-Jan-2006",
		"02-Jan-06",
		"2006-01-02 15:04:05",
		time.RFC3339,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
