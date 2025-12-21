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
	// ข้อมูลพื้นฐาน
	ID           uint    `json:"id"`
	IDCode       string  `json:"id_code"`
	SerialNo     string  `json:"serial_no"`
	ECRIRisk     string  `json:"ecri_risk"`
	AssessmentID *string `json:"assessment_id"`

	// ข้อมูล Brand & Model
	BrandName      string `json:"brand_name"`
	ModelName      string `json:"model_name"`
	CategoryName   string `json:"category_name"`
	Classification string `json:"classification"`

	// ข้อมูล Department
	DepartmentName string `json:"department_name"`

	// ข้อมูลทางการเงิน
	ReceiveDate   *string `json:"receive_date"`
	PurchasePrice float64 `json:"purchase_price"`

	// ข้อมูลวงจรชีวิต
	EquipmentAge    float64 `json:"equipment_age"`
	ComputeDate     *string `json:"compute_date"`
	LifeExpectancy  float64 `json:"life_expectancy"`
	RemainLife      float64 `json:"remain_life"`
	UsefulLifetime  float64 `json:"useful_lifetime"` // % ของอายุการใช้งาน
	ReplacementYear *int    `json:"replacement_year"`

	// คะแนนประเมิน
	Technology      *float64 `json:"technology"`
	UsageStatistics *float64 `json:"usage_statistics"`
	Efficiency      *float64 `json:"efficiency"`
	Others          *string  `json:"others"`

	// Summary (คำนวณจาก maintenance records)
	TotalCM      int64   `json:"total_cm"`       // จำนวนครั้งที่ซ่อม CM
	TotalCost    float64 `json:"total_cost"`     // ค่าใช้จ่ายรวม
	PerCostPrice float64 `json:"per_cost_price"` // % ของราคาซื้อ
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
	AssessmentID *string `json:"assessment_id"`

	ReceiveDate    *string `json:"receive_date"`
	PurchasePrice  float64 `json:"purchase_price"`
	LifeExpectancy float64 `json:"life_expectancy"`

	Technology      *float64 `json:"technology"`
	UsageStatistics *float64 `json:"usage_statistics"`
	Efficiency      *float64 `json:"efficiency"`
	Others          *string  `json:"others"`
}

// EquipmentUpdateDTO - สำหรับแก้ไขอุปกรณ์
type EquipmentUpdateDTO struct {
	ModelID      *uint   `json:"model_id"`
	DepartmentID *uint   `json:"department_id"`
	AssessmentID *string `json:"assessment_id"`

	PurchasePrice   *float64 `json:"purchase_price"`
	LifeExpectancy  *float64 `json:"life_expectancy"`
	ReplacementYear *int     `json:"replacement_year"`

	Technology      *float64 `json:"technology"`
	UsageStatistics *float64 `json:"usage_statistics"`
	Efficiency      *float64 `json:"efficiency"`
	Others          *string  `json:"others"`
}
