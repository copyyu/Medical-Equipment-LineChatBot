package persistence

import (
	"context"
	"log"
	"math/rand"
	"medical-webhook/internal/domain/line/entity"
	"medical-webhook/internal/infrastructure/database"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// EquipmentRepository implements repository.EquipmentRepository using GORM
type EquipmentRepository struct {
	db *gorm.DB
}

// NewEquipmentRepository creates a new equipment repository
func NewEquipmentRepository() *EquipmentRepository {
	return &EquipmentRepository{
		db: database.DB,
	}
}

// FindByIDCode finds equipment by id_code
func (r *EquipmentRepository) FindByIDCode(idCode string) (*entity.Equipment, error) {
	var equipment entity.Equipment
	err := r.db.Preload("Model").Preload("Model.Brand").Preload("Department").
		Where("id_code = ?", idCode).First(&equipment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Equipment not found with id_code: %s", idCode)
			return nil, nil
		}
		log.Printf("Error finding equipment by id_code: %v", err)
		return nil, err
	}
	log.Printf("Found equipment: %s (ID: %d)", equipment.IDCode, equipment.ID)
	return &equipment, nil
}

// FindBySerialNo finds equipment by serial_no
func (r *EquipmentRepository) FindBySerialNo(serialNo string) (*entity.Equipment, error) {
	var equipment entity.Equipment
	err := r.db.Preload("Model").Preload("Model.Brand").Preload("Department").
		Where("serial_no = ?", serialNo).First(&equipment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Equipment not found with serial_no: %s", serialNo)
			return nil, nil
		}
		log.Printf("Error finding equipment by serial_no: %v", err)
		return nil, err
	}
	serial := "N/A"
	if equipment.SerialNo != nil {
		serial = *equipment.SerialNo
	}
	log.Printf("Found equipment: %s (ID: %d)", serial, equipment.ID)
	return &equipment, nil
}

// FindBySerialOrCode finds equipment by either serial_no or id_code
func (r *EquipmentRepository) FindBySerialOrCode(query string) (*entity.Equipment, error) {
	var equipment entity.Equipment
	err := r.db.Preload("Model").Preload("Model.Brand").Preload("Department").
		Where("serial_no = ? OR id_code = ?", query, query).First(&equipment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Equipment not found: %s", query)
			return nil, nil
		}
		log.Printf("Error finding equipment: %v", err)
		return nil, err
	}
	serial := "N/A"
	if equipment.SerialNo != nil {
		serial = *equipment.SerialNo
	}
	log.Printf("Found equipment: %s/%s (ID: %d)", serial, equipment.IDCode, equipment.ID)
	return &equipment, nil
}

// GetMaintenanceRecords gets maintenance records for equipment
func (r *EquipmentRepository) GetMaintenanceRecords(equipmentID uint) ([]entity.MaintenanceRecord, error) {
	var records []entity.MaintenanceRecord
	err := r.db.Where("equipment_id = ?", equipmentID).
		Order("maintenance_date DESC").
		Limit(10).
		Find(&records).Error
	if err != nil {
		log.Printf("Error getting maintenance records: %v", err)
		return nil, err
	}
	log.Printf("Found %d maintenance records for equipment ID: %d", len(records), equipmentID)
	return records, nil
}

// CreateMaintenanceRecord creates a new maintenance record
func (r *EquipmentRepository) CreateMaintenanceRecord(record *entity.MaintenanceRecord) error {
	err := r.db.Create(record).Error
	if err != nil {
		log.Printf("Error creating maintenance record: %v", err)
		return err
	}
	log.Printf("Created maintenance record ID: %d for equipment ID: %d", record.ID, record.EquipmentID)
	return nil
}

// Create creates a new equipment
func (r *EquipmentRepository) Create(ctx context.Context, equipment *entity.Equipment) error {
	err := r.db.WithContext(ctx).Create(equipment).Error
	if err != nil {
		log.Printf("Error creating equipment: %v", err)
		return err
	}
	log.Printf("Created equipment: %s (ID: %d)", equipment.IDCode, equipment.ID)
	return nil
}

// CreateOrUpdate creates or updates equipment based on id_code
// This is used for Excel import to handle duplicate entries
func (r *EquipmentRepository) CreateOrUpdate(ctx context.Context, equipment *entity.Equipment) error {
	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id_code"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"serial_no", "model_id", "department_id", "assessment_id",
			"receive_date", "purchase_price", "equipment_age", "compute_date",
			"life_expectancy", "remain_life", "useful_lifetime_percent",
			"replacement_year", "technology", "usage_statistics", "efficiency", "others",
		}),
	}).Create(equipment).Error

	if err != nil {
		log.Printf("Error creating/updating equipment: %v", err)
		return err
	}
	log.Printf("Created/Updated equipment: %s (ID: %d)", equipment.IDCode, equipment.ID)
	return nil
}

// Update updates equipment
func (r *EquipmentRepository) Update(ctx context.Context, equipment *entity.Equipment) error {
	err := r.db.WithContext(ctx).Save(equipment).Error
	if err != nil {
		log.Printf("Error updating equipment: %v", err)
		return err
	}
	log.Printf("Updated equipment: %s (ID: %d)", equipment.IDCode, equipment.ID)
	return nil
}

// Delete soft deletes equipment
func (r *EquipmentRepository) Delete(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).Delete(&entity.Equipment{}, id).Error
	if err != nil {
		log.Printf("Error deleting equipment: %v", err)
		return err
	}
	log.Printf("Deleted equipment ID: %d", id)
	return nil
}

// FindByID finds equipment by ID with all relations preloaded
func (r *EquipmentRepository) FindByID(ctx context.Context, id uint) (*entity.Equipment, error) {
	var equipment entity.Equipment
	err := r.db.WithContext(ctx).
		Preload("Model").
		Preload("Model.Brand").
		Preload("Model.Category").
		Preload("Department").
		First(&equipment, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Equipment not found with ID: %d", id)
			return nil, nil
		}
		log.Printf("Error finding equipment by ID: %v", err)
		return nil, err
	}
	log.Printf("Found equipment: %s (ID: %d)", equipment.IDCode, equipment.ID)
	return &equipment, nil
}

// FindAll finds all equipments with pagination
func (r *EquipmentRepository) FindAll(ctx context.Context, limit, offset int) ([]entity.Equipment, error) {
	var equipments []entity.Equipment
	query := r.db.WithContext(ctx).
		Preload("Model").
		Preload("Model.Brand").
		Preload("Model.Category").
		Preload("Department").
		Order("id DESC")

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	err := query.Find(&equipments).Error
	if err != nil {
		log.Printf("Error finding all equipments: %v", err)
		return nil, err
	}
	log.Printf("Found %d equipments", len(equipments))
	return equipments, nil
}

// Count returns total count of equipments
func (r *EquipmentRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Equipment{}).Count(&count).Error
	if err != nil {
		log.Printf("Error counting equipments: %v", err)
		return 0, err
	}
	return count, nil
}

// CountWithFilter returns total count of equipments with filters
func (r *EquipmentRepository) CountWithFilter(ctx context.Context, status, search string) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entity.Equipment{})

	// Apply status filter
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Apply search filter (search in id_code, serial_no, and model name via join)
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Joins("LEFT JOIN equipment_models ON equipment_models.id = equipments.model_id").
			Where("equipments.id_code LIKE ? OR equipments.serial_no LIKE ? OR equipment_models.model_name LIKE ?",
				searchPattern, searchPattern, searchPattern)
	}

	err := query.Count(&count).Error
	if err != nil {
		log.Printf("Error counting equipments with filter: %v", err)
		return 0, err
	}
	return count, nil
}

// FindAllWithFilter finds all equipments with pagination and filters
func (r *EquipmentRepository) FindAllWithFilter(ctx context.Context, limit, offset int, status, search string) ([]entity.Equipment, error) {
	var equipments []entity.Equipment
	query := r.db.WithContext(ctx).
		Preload("Model").
		Preload("Model.Brand").
		Preload("Model.Category").
		Preload("Department")

	// Apply status filter
	if status != "" {
		query = query.Where("equipments.status = ?", status)
	}

	// Apply search filter (search in id_code, serial_no, and model name via join)
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Joins("LEFT JOIN equipment_models ON equipment_models.id = equipments.model_id").
			Where("equipments.id_code LIKE ? OR equipments.serial_no LIKE ? OR equipment_models.model_name LIKE ?",
				searchPattern, searchPattern, searchPattern)
	}

	query = query.Order("equipments.id DESC")

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	err := query.Find(&equipments).Error
	if err != nil {
		log.Printf("Error finding equipments with filter: %v", err)
		return nil, err
	}
	log.Printf("Found %d equipments with filter (status=%s, search=%s)", len(equipments), status, search)
	return equipments, nil
}

// CountNearExpiry returns count of equipments with dynamic remain_life between 0 and 1 year
// Uses PostgreSQL: life_expectancy - (NOW()::date - receive_date::date) / 365.25
func (r *EquipmentRepository) CountNearExpiry(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Equipment{}).
		Where("receive_date IS NOT NULL AND life_expectancy > 0").
		Where("(life_expectancy - (NOW()::date - receive_date::date) / 365.25) > 0").
		Where("(life_expectancy - (NOW()::date - receive_date::date) / 365.25) <= 1").
		Count(&count).Error
	if err != nil {
		log.Printf("Error counting near expiry equipments: %v", err)
		return 0, err
	}
	return count, nil
}

// CountByStatus returns count of equipments grouped by status
func (r *EquipmentRepository) CountByStatus(ctx context.Context) (map[entity.AssetStatus]int64, error) {
	var results []struct {
		Status string
		Count  int64
	}

	err := r.db.WithContext(ctx).
		Model(&entity.Equipment{}).
		Select("status, count(*) as count").
		Group("status").
		Find(&results).Error

	if err != nil {
		log.Printf("Error counting equipments by status: %v", err)
		return nil, err
	}

	counts := make(map[entity.AssetStatus]int64)
	for _, r := range results {
		counts[entity.AssetStatus(r.Status)] = r.Count
	}
	return counts, nil
}

// CountExpired returns count of equipments with dynamic remain_life <= 0
// Uses PostgreSQL: life_expectancy - (NOW()::date - receive_date::date) / 365.25
func (r *EquipmentRepository) CountExpired(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Equipment{}).
		Where("receive_date IS NOT NULL AND life_expectancy > 0").
		Where("(life_expectancy - (NOW()::date - receive_date::date) / 365.25) <= 0").
		Count(&count).Error

	if err != nil {
		log.Printf("Error counting expired equipments: %v", err)
		return 0, err
	}
	return count, nil
}

// FindExpired returns equipments with dynamic remain_life <= 0 (expired)
func (r *EquipmentRepository) FindExpired(ctx context.Context, limit int) ([]entity.Equipment, error) {
	var equipments []entity.Equipment
	query := r.db.WithContext(ctx).
		Preload("Model").
		Preload("Model.Brand").
		Preload("Department").
		Where("receive_date IS NOT NULL AND life_expectancy > 0").
		Where("(life_expectancy - (NOW()::date - receive_date::date) / 365.25) <= 0").
		Order("(life_expectancy - (NOW()::date - receive_date::date) / 365.25) ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&equipments).Error
	if err != nil {
		log.Printf("Error finding expired equipments: %v", err)
		return nil, err
	}
	log.Printf("Found %d expired equipments", len(equipments))
	return equipments, nil
}

// FindNearExpiry returns equipments with dynamic remain_life between 0 and 1 year
func (r *EquipmentRepository) FindNearExpiry(ctx context.Context, limit int) ([]entity.Equipment, error) {
	var equipments []entity.Equipment
	query := r.db.WithContext(ctx).
		Preload("Model").
		Preload("Model.Brand").
		Preload("Department").
		Where("receive_date IS NOT NULL AND life_expectancy > 0").
		Where("(life_expectancy - (NOW()::date - receive_date::date) / 365.25) > 0").
		Where("(life_expectancy - (NOW()::date - receive_date::date) / 365.25) <= 1").
		Order("(life_expectancy - (NOW()::date - receive_date::date) / 365.25) ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&equipments).Error
	if err != nil {
		log.Printf("Error finding near expiry equipments: %v", err)
		return nil, err
	}
	log.Printf("Found %d near expiry equipments", len(equipments))
	return equipments, nil
}

// FindExpiredByDepartment returns expired equipments filtered by department ID
func (r *EquipmentRepository) FindExpiredByDepartment(ctx context.Context, departmentID uint, limit int) ([]entity.Equipment, error) {
	var equipments []entity.Equipment
	query := r.db.WithContext(ctx).
		Preload("Model").
		Preload("Model.Brand").
		Preload("Department").
		Where("department_id = ?", departmentID).
		Where("receive_date IS NOT NULL AND life_expectancy > 0").
		Where("(life_expectancy - (NOW()::date - receive_date::date) / 365.25) <= 0").
		Order("(life_expectancy - (NOW()::date - receive_date::date) / 365.25) ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&equipments).Error
	if err != nil {
		log.Printf("Error finding expired equipments by department: %v", err)
		return nil, err
	}
	log.Printf("Found %d expired equipments for department %d", len(equipments), departmentID)
	return equipments, nil
}

// FindNearExpiryByDepartment returns near-expiry equipments filtered by department ID
func (r *EquipmentRepository) FindNearExpiryByDepartment(ctx context.Context, departmentID uint, limit int) ([]entity.Equipment, error) {
	var equipments []entity.Equipment
	query := r.db.WithContext(ctx).
		Preload("Model").
		Preload("Model.Brand").
		Preload("Department").
		Where("department_id = ?", departmentID).
		Where("receive_date IS NOT NULL AND life_expectancy > 0").
		Where("(life_expectancy - (NOW()::date - receive_date::date) / 365.25) > 0").
		Where("(life_expectancy - (NOW()::date - receive_date::date) / 365.25) <= 1").
		Order("(life_expectancy - (NOW()::date - receive_date::date) / 365.25) ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&equipments).Error
	if err != nil {
		log.Printf("Error finding near expiry equipments by department: %v", err)
		return nil, err
	}
	log.Printf("Found %d near expiry equipments for department %d", len(equipments), departmentID)
	return equipments, nil
}

func (r *EquipmentRepository) FindSimilarByIDCodePrefix(prefix string, limit int) ([]*entity.Equipment, error) {
	var equipments []*entity.Equipment

	// ดึงมาทั้งหมดที่ match prefix ก่อน
	err := r.db.
		Where("id_code LIKE ?", prefix+"%").
		Find(&equipments).Error
	if err != nil {
		return nil, err
	}

	// สุ่ม shuffle
	rand.Shuffle(len(equipments), func(i, j int) {
		equipments[i], equipments[j] = equipments[j], equipments[i]
	})

	// ตัดให้เหลือแค่ limit
	if limit > 0 && len(equipments) > limit {
		equipments = equipments[:limit]
	}

	return equipments, nil
}

// FindBestMatch finds the single most similar equipment by id_code using pg_trgm trigram similarity.
// Returns the best matching equipment, similarity percentage (0-100), and error.
// This handles OCR misread characters (e.g., SSH12345 → SSH123S5) robustly.
func (r *EquipmentRepository) FindBestMatch(query string) (*entity.Equipment, int, error) {
	var result struct {
		entity.Equipment
		Sim float64 `gorm:"column:sim"`
	}

	err := r.db.
		Table("equipments").
		Select("equipments.*, similarity(id_code, ?) AS sim", query).
		Where("similarity(id_code, ?) > 0.3", query).
		Where("deleted_at IS NULL").
		Order("sim DESC").
		Limit(1).
		Scan(&result).Error

	if err != nil {
		log.Printf("❌ FindBestMatch error: %v", err)
		return nil, 0, err
	}

	if result.ID == 0 {
		log.Printf("⚠️ No match found for query: %s", query)
		return nil, 0, nil
	}

	pct := int(result.Sim * 100)
	log.Printf("✅ Best match: %s (similarity: %d%%) for query: %s", result.IDCode, pct, query)
	return &result.Equipment, pct, nil
}

// FindSimilarSorted finds similar equipment by prefix, sorted by pg_trgm similarity (most similar first).
func (r *EquipmentRepository) FindSimilarSorted(query string, limit int) ([]*entity.Equipment, error) {
	prefix := query
	if len(query) >= 6 {
		prefix = query[:6]
	}

	var equipments []*entity.Equipment
	err := r.db.
		Where("id_code LIKE ? AND deleted_at IS NULL", prefix+"%").
		Order("similarity(id_code, '" + query + "') DESC").
		Limit(limit).
		Find(&equipments).Error

	if err != nil {
		log.Printf("❌ FindSimilarSorted error: %v", err)
		return nil, err
	}

	return equipments, nil
}

// FindByReplacementYear finds equipment where replacement_year matches the given year
// If departmentID is not nil, filter by department
func (r *EquipmentRepository) FindByReplacementYear(ctx context.Context, year int, departmentID *uint) ([]entity.Equipment, error) {
	var equipments []entity.Equipment
	query := r.db.WithContext(ctx).
		Preload("Model").
		Preload("Model.Brand").
		Preload("Model.Category").
		Preload("Department").
		Where("replacement_year = ?", year)

	if departmentID != nil {
		query = query.Where("department_id = ?", *departmentID)
	}

	query = query.Order("id_code ASC")

	err := query.Find(&equipments).Error
	if err != nil {
		log.Printf("Error finding equipment by replacement year %d: %v", year, err)
		return nil, err
	}
	log.Printf("Found %d equipments for replacement year %d", len(equipments), year)
	return equipments, nil
}
