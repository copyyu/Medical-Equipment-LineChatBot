package persistence

import (
	"context"
	"log"
	"medical-webhook/internal/domain/line/entity"
	"medical-webhook/internal/infrastructure/database"

	"gorm.io/gorm"
)

// MaintenanceRecordRepository implements repository.MaintenanceRecordRepository using GORM
type MaintenanceRecordRepository struct {
	db *gorm.DB
}

// NewMaintenanceRecordRepository creates a new maintenance record repository
func NewMaintenanceRecordRepository() *MaintenanceRecordRepository {
	return &MaintenanceRecordRepository{
		db: database.DB,
	}
}

// Create creates a new maintenance record
func (r *MaintenanceRecordRepository) Create(ctx context.Context, record *entity.MaintenanceRecord) error {
	err := r.db.WithContext(ctx).Create(record).Error
	if err != nil {
		log.Printf("Error creating maintenance record: %v", err)
		return err
	}
	log.Printf("Created maintenance record (ID: %d) for equipment ID: %d", record.ID, record.EquipmentID)
	return nil
}

// FindByID finds maintenance record by ID
func (r *MaintenanceRecordRepository) FindByID(ctx context.Context, id uint) (*entity.MaintenanceRecord, error) {
	var record entity.MaintenanceRecord
	err := r.db.WithContext(ctx).
		Preload("Equipment").
		First(&record, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Maintenance record not found with ID: %d", id)
			return nil, nil
		}
		log.Printf("Error finding maintenance record by ID: %v", err)
		return nil, err
	}
	log.Printf("Found maintenance record ID: %d", record.ID)
	return &record, nil
}

// FindByEquipmentID finds all maintenance records for specific equipment
func (r *MaintenanceRecordRepository) FindByEquipmentID(ctx context.Context, equipmentID uint) ([]entity.MaintenanceRecord, error) {
	var records []entity.MaintenanceRecord
	err := r.db.WithContext(ctx).
		Where("equipment_id = ?", equipmentID).
		Order("maintenance_date DESC").
		Find(&records).Error
	if err != nil {
		log.Printf("Error finding maintenance records for equipment ID %d: %v", equipmentID, err)
		return nil, err
	}
	log.Printf("Found %d maintenance records for equipment ID: %d", len(records), equipmentID)
	return records, nil
}

// FindByEquipmentIDAndType finds maintenance records by equipment ID and type
func (r *MaintenanceRecordRepository) FindByEquipmentIDAndType(ctx context.Context, equipmentID uint, maintenanceType entity.MaintenanceType) ([]entity.MaintenanceRecord, error) {
	var records []entity.MaintenanceRecord
	err := r.db.WithContext(ctx).
		Where("equipment_id = ? AND maintenance_type = ?", equipmentID, maintenanceType).
		Order("maintenance_date DESC").
		Find(&records).Error
	if err != nil {
		log.Printf("Error finding %s records for equipment ID %d: %v", maintenanceType, equipmentID, err)
		return nil, err
	}
	log.Printf("Found %d %s records for equipment ID: %d", len(records), maintenanceType, equipmentID)
	return records, nil
}

// GetTotalCostByEquipmentID calculates total maintenance cost for equipment
func (r *MaintenanceRecordRepository) GetTotalCostByEquipmentID(ctx context.Context, equipmentID uint) (float64, error) {
	var totalCost float64
	err := r.db.WithContext(ctx).
		Model(&entity.MaintenanceRecord{}).
		Where("equipment_id = ?", equipmentID).
		Select("COALESCE(SUM(cost), 0)").
		Scan(&totalCost).Error
	if err != nil {
		log.Printf("Error calculating total cost for equipment ID %d: %v", equipmentID, err)
		return 0, err
	}
	log.Printf("Total maintenance cost for equipment ID %d: %.2f", equipmentID, totalCost)
	return totalCost, nil
}

// GetCMCountByEquipmentID counts corrective maintenance records for equipment
func (r *MaintenanceRecordRepository) GetCMCountByEquipmentID(ctx context.Context, equipmentID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.MaintenanceRecord{}).
		Where("equipment_id = ? AND maintenance_type = ?", equipmentID, entity.MaintenanceCM).
		Count(&count).Error
	if err != nil {
		log.Printf("Error counting CM records for equipment ID %d: %v", equipmentID, err)
		return 0, err
	}
	log.Printf("CM count for equipment ID %d: %d", equipmentID, count)
	return count, nil
}

// Update updates maintenance record
func (r *MaintenanceRecordRepository) Update(ctx context.Context, record *entity.MaintenanceRecord) error {
	err := r.db.WithContext(ctx).Save(record).Error
	if err != nil {
		log.Printf("Error updating maintenance record: %v", err)
		return err
	}
	log.Printf("Updated maintenance record ID: %d", record.ID)
	return nil
}

// Delete soft deletes maintenance record
func (r *MaintenanceRecordRepository) Delete(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).Delete(&entity.MaintenanceRecord{}, id).Error
	if err != nil {
		log.Printf("Error deleting maintenance record: %v", err)
		return err
	}
	log.Printf("Deleted maintenance record ID: %d", id)
	return nil
}
