package dto

// EquipmentSummaryDTO - ข้อมูลสรุปอุปกรณ์ (คำนวณจาก maintenance records)
type EquipmentSummaryDTO struct {
	EquipmentID  uint    `json:"equipment_id"`
	IDCode       string  `json:"id_code"`
	TotalCM      int64   `json:"total_cm"`       // จำนวนครั้งที่ซ่อม CM
	TotalCost    float64 `json:"total_cost"`     // ค่าใช้จ่ายรวมทั้งหมด
	PerCostPrice float64 `json:"per_cost_price"` // % ของราคาซื้อ
}

// EquipmentDetailDTO - ข้อมูลอุปกรณ์แบบเต็ม (สำหรับ API response)
type EquipmentDetailDTO struct {
	ID              uint    `json:"id"`
	IDCode          string  `json:"id_code"`
	SerialNo        string  `json:"serial_no"`
	ECRIRisk        string  `json:"ecri_risk"`
	ECRICode        *string `json:"ecri_code"`
	BrandName       string  `json:"brand_name"`
	ModelName       string  `json:"model_name"`
	CategoryName    string  `json:"category_name"`
	Classification  string  `json:"classification"`
	DepartmentName  string  `json:"department_name"`
	ReceiveDate     *string `json:"receive_date"`
	PurchasePrice   float64 `json:"purchase_price"`
	EquipmentAge    float64 `json:"equipment_age"`
	LifeExpectancy  float64 `json:"life_expectancy"`
	RemainLife      float64 `json:"remain_life"`
	ReplacementYear *int    `json:"replacement_year"`
	AssetTypeName   *string `json:"asset_type_name"`
	AssetName       *string `json:"asset_name"`
	Remark          *string `json:"remark"`
	TotalCM         int64   `json:"total_cm"`
	TotalCost       float64 `json:"total_cost"`
	PerCostPrice    float64 `json:"per_cost_price"`
}

// EquipmentListDTO - สำหรับแสดงรายการอุปกรณ์แบบย่อ
type EquipmentListDTO struct {
	ID             uint    `json:"id"`
	IDCode         string  `json:"id_code"`
	SerialNo       string  `json:"serial_no"`
	BrandName      string  `json:"brand_name"`
	ModelName      string  `json:"model_name"`
	DepartmentName string  `json:"department_name"`
	PurchasePrice  float64 `json:"purchase_price"`
	EquipmentAge   float64 `json:"equipment_age"`
	RemainLife     float64 `json:"remain_life"`
}

// EquipmentCreateDTO - สำหรับสร้างอุปกรณ์ใหม่
type EquipmentCreateDTO struct {
	IDCode       string  `json:"id_code" binding:"required"`
	SerialNo     string  `json:"serial_no" binding:"required"`
	ModelID      uint    `json:"model_id" binding:"required"`
	DepartmentID uint    `json:"department_id" binding:"required"`
	ECRICode     *string `json:"ecri_code"`

	ReceiveDate    *string `json:"receive_date"`
	PurchasePrice  float64 `json:"purchase_price"`
	LifeExpectancy float64 `json:"life_expectancy"`

	Remark *string `json:"remark"`
}

// EquipmentUpdateDTO - สำหรับแก้ไขอุปกรณ์
type EquipmentUpdateDTO struct {
	ModelID      *uint   `json:"model_id"`
	DepartmentID *uint   `json:"department_id"`
	ECRICode     *string `json:"ecri_code"`

	PurchasePrice   *float64 `json:"purchase_price"`
	LifeExpectancy  *float64 `json:"life_expectancy"`
	ReplacementYear *int     `json:"replacement_year"`

	Remark *string `json:"remark"`
}
