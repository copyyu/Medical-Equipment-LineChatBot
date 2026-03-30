package dto

import "time"

// ExcelRowDTO - DTO สำหรับข้อมูลจาก Excel แต่ละ row
type ExcelRowDTO struct {
	AssetTypeName        string
	Category             string
	ECRICode             string
	Brand                string
	Model                string
	SerialNo             *string
	Building             *string
	AssetStatus          string
	AssetStatusInternal  *string
	RentalStatus         *string
	BusinessName         *string
	IDCode               string
	ItemNo               *string
	SKUNo                *string
	UpdatedDate          *time.Time
	UpdatedBy            *string
	Department           string
	WarrantyPeriod       *string
	WarrantyStartDate    *time.Time
	WarrantyEndDate      *time.Time
	WarrantyPM           *string
	WarrantyCal          *string
	Floor                *string
	Room                 *string
	PhoneNo              *string
	PowerConsumption     *string
	ECRIRisk             string
	Classification       string
	LifeExpectancy       float64
	CalPeriod            *string
	VendorPM             *string
	VendorCal            *string
	TorNo                *string
	PurchaseDate         *time.Time
	PurchasePrice        float64
	ReceiveDate          *time.Time
	RegistrationDate     *time.Time
	Supplier             *string
	Ownership            *string
	PoNo                 *string
	ContractNo           *string
	InvoiceNo            *string
	DocumentNo           *string
	ManufacturingCountry *string
	RevenuePerMonth      *float64
	Remark               *string
	ApprovedBy           *string
	NsmartItemCode       *string
	AssetName            *string
	AssetID              *string
	LastPMDate           *time.Time
	LastCalDate          *time.Time
	PMPeriod             *string
	BorrowStatus         *string
}

// EquipmentImportResultDTO - DTO สำหรับผลลัพธ์การ import
type EquipmentImportResultDTO struct {
	TotalRows      int      `json:"total_rows"`
	SuccessCount   int      `json:"success_count"`
	UpdatedCount   int      `json:"updated_count"`
	FailedCount    int      `json:"failed_count"`
	SkippedCount   int      `json:"skipped_count"`
	NewBrands      int      `json:"new_brands"`
	NewCategories  int      `json:"new_categories"`
	NewDepartments int      `json:"new_departments"`
	NewModels      int      `json:"new_models"`
	FailedRows     []int    `json:"failed_rows"`
	ErrorMessages  []string `json:"error_messages"`
}

// CreateEquipmentDTO - DTO สำหรับสร้าง Equipment (from Excel import)
type CreateEquipmentDTO struct {
	IDCode               string
	SerialNo             *string
	ModelID              uint
	DepartmentID         uint
	AssetTypeName        *string
	ECRICode             *string
	AssetName            *string
	AssetID              *string
	AssetStatus          string
	AssetStatusInternal  *string
	RentalStatus         *string
	BorrowStatus         *string
	Building             *string
	Floor                *string
	Room                 *string
	PhoneNo              *string
	BusinessName         *string
	ItemNo               *string
	SKUNo                *string
	ReceiveDate          *time.Time
	PurchaseDate         *time.Time
	RegistrationDate     *time.Time
	PurchasePrice        float64
	LifeExpectancy       float64
	WarrantyPeriod       *string
	WarrantyStartDate    *time.Time
	WarrantyEndDate      *time.Time
	WarrantyPM           *string
	WarrantyCal          *string
	LastPMDate           *time.Time
	LastCalDate          *time.Time
	PMPeriod             *string
	CalPeriod            *string
	VendorPM             *string
	VendorCal            *string
	PowerConsumption     *string
	Supplier             *string
	Ownership            *string
	PoNo                 *string
	ContractNo           *string
	InvoiceNo            *string
	DocumentNo           *string
	TorNo                *string
	ManufacturingCountry *string
	RevenuePerMonth      *float64
	Remark               *string
	ApprovedBy           *string
	NsmartItemCode       *string
	UpdatedBy            *string
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
