package dto

import "time"

// ExcelRowDTO - DTO สำหรับข้อมูลจาก Excel แต่ละ row
type ExcelRowDTO struct {
	Department        string
	ECRIRisk          string
	AssessmentID      string
	IDCode            string
	Category          string
	Brand             string
	Model             string
	SerialNo          *string
	Classification    string
	ReceiveDate       *time.Time
	PurchasePrice     float64
	EquipmentAge      float64
	ComputeDate       *time.Time
	LifeExpectancy    float64
	RemainLife        float64
	TotalOfCM         int
	TotalOfCost       float64
	PerCostPrice      float64
	UsefulLifePercent float64
	ReplacementYear   *int
	Technology        *float64
	UsageStatistics   *float64
	Efficiency        *float64
	Others            *string
}

// EquipmentImportResultDTO - DTO สำหรับผลลัพธ์การ import
type EquipmentImportResultDTO struct {
	TotalRows      int      `json:"total_rows"`
	SuccessCount   int      `json:"success_count"`
	FailedCount    int      `json:"failed_count"`
	SkippedCount   int      `json:"skipped_count"`
	NewBrands      int      `json:"new_brands"`
	NewCategories  int      `json:"new_categories"`
	NewDepartments int      `json:"new_departments"`
	NewModels      int      `json:"new_models"`
	FailedRows     []int    `json:"failed_rows"`
	ErrorMessages  []string `json:"error_messages"`
}

// CreateEquipmentDTO - DTO สำหรับสร้าง Equipment
type CreateEquipmentDTO struct {
	IDCode                string
	SerialNo              *string
	ModelID               uint
	DepartmentID          uint
	AssessmentID          *string
	ReceiveDate           *time.Time
	PurchasePrice         float64
	EquipmentAge          float64
	ComputeDate           *time.Time
	LifeExpectancy        float64
	RemainLife            float64
	UsefulLifetimePercent float64
	ReplacementYear       *int
	Technology            *float64
	UsageStatistics       *float64
	Efficiency            *float64
	Others                *string
}

// BrandDTO - DTO สำหรับ Brand
type BrandDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// CategoryDTO - DTO สำหรับ Category
type CategoryDTO struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	ECRIRisk       string `json:"ecri_risk"`
	Classification string `json:"classification"`
}

// DepartmentDTO - DTO สำหรับ Department
type DepartmentDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ModelDTO - DTO สำหรับ Model
type ModelDTO struct {
	ID                    uint    `json:"id"`
	BrandID               uint    `json:"brand_id"`
	CategoryID            uint    `json:"category_id"`
	ModelName             string  `json:"model_name"`
	DefaultLifeExpectancy float64 `json:"default_life_expectancy"`
}
