package database

import (
	"log"
	"medical-webhook/internal/domain/line/entity"

	"gorm.io/gorm"
)

// EquipmentRepository implements repository.EquipmentRepository using GORM
type EquipmentRepository struct {
	db *gorm.DB
}

// NewEquipmentRepository creates a new equipment repository
func NewEquipmentRepository() *EquipmentRepository {
	return &EquipmentRepository{
		db: DB, // Use global DB from db_connect.go
	}
}

// FindByIDCode finds equipment by id_code
func (r *EquipmentRepository) FindByIDCode(idCode string) (*entity.Equipment, error) {
	var equipment entity.Equipment
	err := r.db.Preload("Model").Preload("Model.Brand").Preload("Department").
		Where("id_code = ?", idCode).First(&equipment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("⚠️ Equipment not found with id_code: %s", idCode)
			return nil, nil
		}
		log.Printf("❌ Error finding equipment by id_code: %v", err)
		return nil, err
	}
	log.Printf("✅ Found equipment: %s (ID: %d)", equipment.IDCode, equipment.ID)
	return &equipment, nil
}

// FindBySerialNo finds equipment by serial_no
func (r *EquipmentRepository) FindBySerialNo(serialNo string) (*entity.Equipment, error) {
	var equipment entity.Equipment
	err := r.db.Preload("Model").Preload("Model.Brand").Preload("Department").
		Where("serial_no = ?", serialNo).First(&equipment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("⚠️ Equipment not found with serial_no: %s", serialNo)
			return nil, nil
		}
		log.Printf("❌ Error finding equipment by serial_no: %v", err)
		return nil, err
	}
	log.Printf("✅ Found equipment: %s (ID: %d)", equipment.SerialNo, equipment.ID)
	return &equipment, nil
}

// FindBySerialOrCode finds equipment by either serial_no or id_code
func (r *EquipmentRepository) FindBySerialOrCode(query string) (*entity.Equipment, error) {
	var equipment entity.Equipment
	err := r.db.Preload("Model").Preload("Model.Brand").Preload("Department").
		Where("serial_no = ? OR id_code = ?", query, query).First(&equipment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("⚠️ Equipment not found: %s", query)
			return nil, nil
		}
		log.Printf("❌ Error finding equipment: %v", err)
		return nil, err
	}
	log.Printf("✅ Found equipment: %s/%s (ID: %d)", equipment.SerialNo, equipment.IDCode, equipment.ID)
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
		log.Printf("❌ Error getting maintenance records: %v", err)
		return nil, err
	}
	log.Printf("✅ Found %d maintenance records for equipment ID: %d", len(records), equipmentID)
	return records, nil
}
