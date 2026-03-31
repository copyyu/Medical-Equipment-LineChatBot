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

// ParseExcelRow - parse แต่ละ row จาก Excel (54 columns)
func (s *excelParserService) ParseExcelRow(row []string, rowNum int) (*dto.ExcelRowDTO, error) {
	// Helper function to safely get column value
	getCol := func(idx int) string {
		if idx < len(row) {
			return strings.TrimSpace(row[idx])
		}
		return ""
	}

	// ID Code (col 11) is required
	idCode := getCol(11)
	if idCode == "" {
		return nil, fmt.Errorf("row %d: ID Code is empty", rowNum)
	}

	data := &dto.ExcelRowDTO{
		AssetTypeName:  getCol(0),         // Asset Type Name
		Category:       getCol(1),         // Category
		ECRICode:       getCol(2),         // ECRI Code
		Brand:          getCol(3),         // Brand Name
		Model:          getCol(4),         // Model Name
		SerialNo:       strPtr(getCol(5)), // Serial Number
		Building:       strPtr(getCol(6)), // Building
		AssetStatus:    getCol(7),         // Asset Status
		IDCode:         idCode,            // ID Code
		ECRIRisk:       getCol(26),        // Risk Level
		Classification: getCol(27),        // Classification
		Department:     getCol(16),        // Department Name
	}

	// col 8: Asset Status Internal
	if v := getCol(8); v != "" {
		data.AssetStatusInternal = &v
	}

	// col 9: Rental Status
	if v := getCol(9); v != "" {
		data.RentalStatus = &v
	}

	// col 10: Business Name
	if v := getCol(10); v != "" {
		data.BusinessName = &v
	}

	// col 12: Item No
	if v := getCol(12); v != "" {
		data.ItemNo = &v
	}

	// col 13: SKU No
	if v := getCol(13); v != "" {
		data.SKUNo = &v
	}

	// col 14: Updated Date
	if dateStr := getCol(14); dateStr != "" {
		if date, err := s.parseDate(dateStr); err == nil {
			data.UpdatedDate = &date
		}
	}

	// col 15: Updated By
	if v := getCol(15); v != "" {
		data.UpdatedBy = &v
	}

	// col 17: Warranty Period
	if v := getCol(17); v != "" {
		data.WarrantyPeriod = &v
	}

	// col 18: Warranty Start Date
	if dateStr := getCol(18); dateStr != "" {
		if date, err := s.parseDate(dateStr); err == nil {
			data.WarrantyStartDate = &date
		}
	}

	// col 19: Warranty End Date
	if dateStr := getCol(19); dateStr != "" {
		if date, err := s.parseDate(dateStr); err == nil {
			data.WarrantyEndDate = &date
		}
	}

	// col 20: Warranty PM
	if v := getCol(20); v != "" {
		data.WarrantyPM = &v
	}

	// col 21: Warranty Cal
	if v := getCol(21); v != "" {
		data.WarrantyCal = &v
	}

	// col 22: Floor
	if v := getCol(22); v != "" {
		data.Floor = &v
	}

	// col 23: Room
	if v := getCol(23); v != "" {
		data.Room = &v
	}

	// col 24: Phone No
	if v := getCol(24); v != "" {
		data.PhoneNo = &v
	}

	// col 25: Power Consumption
	if v := getCol(25); v != "" {
		data.PowerConsumption = &v
	}

	// col 28: Estimated Use Life → LifeExpectancy
	if lifeStr := getCol(28); lifeStr != "" {
		if life, err := strconv.ParseFloat(lifeStr, 64); err == nil {
			data.LifeExpectancy = life
		}
	}

	// col 29: Cal Period
	if v := getCol(29); v != "" {
		data.CalPeriod = &v
	}

	// col 30: Vendor PM
	if v := getCol(30); v != "" {
		data.VendorPM = &v
	}

	// col 31: Vendor Cal
	if v := getCol(31); v != "" {
		data.VendorCal = &v
	}

	// col 32: Tor No
	if v := getCol(32); v != "" {
		data.TorNo = &v
	}

	// col 33: Purchase Date
	if dateStr := getCol(33); dateStr != "" {
		if date, err := s.parseDate(dateStr); err == nil {
			data.PurchaseDate = &date
		}
	}

	// col 34: Price → PurchasePrice
	if priceStr := getCol(34); priceStr != "" {
		if price, err := strconv.ParseFloat(priceStr, 64); err == nil {
			data.PurchasePrice = price
		}
	}

	// col 35: Receive Date
	if dateStr := getCol(35); dateStr != "" {
		if date, err := s.parseDate(dateStr); err == nil {
			data.ReceiveDate = &date
		}
	}

	// col 36: Registeration Date
	if dateStr := getCol(36); dateStr != "" {
		if date, err := s.parseDate(dateStr); err == nil {
			data.RegistrationDate = &date
		}
	}

	// col 37: Supplier
	if v := getCol(37); v != "" {
		data.Supplier = &v
	}

	// col 38: Ownership
	if v := getCol(38); v != "" {
		data.Ownership = &v
	}

	// col 39: Po No
	if v := getCol(39); v != "" {
		data.PoNo = &v
	}

	// col 40: Contract No
	if v := getCol(40); v != "" {
		data.ContractNo = &v
	}

	// col 41: Invoice No
	if v := getCol(41); v != "" {
		data.InvoiceNo = &v
	}

	// col 42: Document No
	if v := getCol(42); v != "" {
		data.DocumentNo = &v
	}

	// col 43: Manufacturing Country
	if v := getCol(43); v != "" {
		data.ManufacturingCountry = &v
	}

	// col 44: Revenue Per Month
	if revStr := getCol(44); revStr != "" {
		if rev, err := strconv.ParseFloat(revStr, 64); err == nil {
			data.RevenuePerMonth = &rev
		}
	}

	// col 45: Remark
	if v := getCol(45); v != "" {
		data.Remark = &v
	}

	// col 46: Approved By
	if v := getCol(46); v != "" {
		data.ApprovedBy = &v
	}

	// col 47: Nsmart Item Code
	if v := getCol(47); v != "" {
		data.NsmartItemCode = &v
	}

	// col 48: Asset Name
	if v := getCol(48); v != "" {
		data.AssetName = &v
	}

	// col 49: Asset ID
	if v := getCol(49); v != "" {
		data.AssetID = &v
	}

	// col 50: Last PM Date
	if dateStr := getCol(50); dateStr != "" {
		if date, err := s.parseDate(dateStr); err == nil {
			data.LastPMDate = &date
		}
	}

	// col 51: Last Cal Date
	if dateStr := getCol(51); dateStr != "" {
		if date, err := s.parseDate(dateStr); err == nil {
			data.LastCalDate = &date
		}
	}

	// col 52: PM Period
	if v := getCol(52); v != "" {
		data.PMPeriod = &v
	}

	// col 53: Borrow status
	if v := getCol(53); v != "" {
		data.BorrowStatus = &v
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
