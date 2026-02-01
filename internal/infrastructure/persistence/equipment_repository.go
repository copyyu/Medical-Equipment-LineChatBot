package persistence

import (
	"context"
	"log"
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

// CountNearExpiry returns count of equipments with remain_life <= 1 year
func (r *EquipmentRepository) CountNearExpiry(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Equipment{}).
		Where("remain_life <= ?", 1.0).
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
