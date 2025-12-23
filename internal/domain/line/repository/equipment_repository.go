package repository

import "medical-webhook/internal/domain/line/entity"

// EquipmentRepository defines interface for equipment database operations
type EquipmentRepository interface {
	// FindByIDCode finds equipment by id_code
	FindByIDCode(idCode string) (*entity.Equipment, error)

	// FindBySerialNo finds equipment by serial_no
	FindBySerialNo(serialNo string) (*entity.Equipment, error)

	// FindBySerialOrCode finds equipment by either serial_no or id_code
	FindBySerialOrCode(query string) (*entity.Equipment, error)

	// GetMaintenanceRecords gets maintenance records for equipment
	GetMaintenanceRecords(equipmentID uint) ([]entity.MaintenanceRecord, error)
}
