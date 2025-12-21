package line

import (
	"context"
	"fmt"
	"log"
	"medical-webhook/internal/application/dto"
	"medical-webhook/internal/domain/line/entity"
	"time"

	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// GetSettings - ดึงการตั้งค่า
func (r *NotificationRepository) GetSettings(ctx context.Context) (*entity.NotificationSetting, error) {
	var settings entity.NotificationSetting
	err := r.db.WithContext(ctx).First(&settings).Error
	if err == gorm.ErrRecordNotFound {
		// สร้างค่าเริ่มต้นถ้าไม่มี
		settings = entity.NotificationSetting{
			IsEnabled: false,
		}
		r.db.WithContext(ctx).Create(&settings)
		return &settings, nil
	}
	return &settings, err
}

// UpdateSettings - อัพเดทการตั้งค่า
func (r *NotificationRepository) UpdateSettings(ctx context.Context, settings *entity.NotificationSetting) error {
	return r.db.WithContext(ctx).Save(settings).Error
}

func (r *NotificationRepository) CountAllEquipments(ctx context.Context) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Equipment{}). // เปลี่ยนเป็น entity ของคุณ
		Where("deleted_at IS NULL").
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count equipments: %w", err)
	}

	return int(count), nil
}

// CreateLog - บันทึก log
func (r *NotificationRepository) CreateLog(ctx context.Context, log *entity.NotificationLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// GetLastNotification - ดึงการแจ้งเตือนล่าสุด
func (r *NotificationRepository) GetLastNotification(ctx context.Context) (*entity.NotificationLog, error) {
	var log entity.NotificationLog
	err := r.db.WithContext(ctx).
		Order("sent_at DESC").
		First(&log).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &log, err
}

// GetLogsByMonth - ดึง logs ตามเดือน
func (r *NotificationRepository) GetLogsByMonth(ctx context.Context, year int, month int) ([]entity.NotificationLog, error) {
	var logs []entity.NotificationLog
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, 0)

	err := r.db.WithContext(ctx).
		Where("sent_at >= ? AND sent_at < ?", startDate, endDate).
		Order("sent_at DESC").
		Find(&logs).Error

	return logs, err
}

// August Alert - แจ้งทุกเครื่องที่เหลือ 11-13 เดือน
func (r *NotificationRepository) GetEquipmentsForAugustAlert(ctx context.Context) ([]dto.EquipmentReplacementAlertDTO, error) {
	var results []dto.EquipmentReplacementAlertDTO

	currentYear := time.Now().Year()
	nextYear := currentYear + 1

	// TEST: แจ้งเครื่องปีหน้า (2026) เหลือ 7-9 เดือน เพื่อทดสอบ
	// targetYear := nextYear
	// monthsRangeMin := 7
	// monthsRangeMax := 9

	// PRODUCTION: แจ้งเครื่องที่ครบปีหน้า (เหลือ 12 เดือน)
	targetYear := nextYear
	monthsRangeMin := 11
	monthsRangeMax := 13

	query := `
		SELECT 
			e.id as equipment_id,
			e.id_code,
			COALESCE(e.serial_no, '') as serial_no,
			b.name as brand_name,
			m.model_name,
			d.name as department_name,
			e.replacement_year,
			(e.replacement_year - EXTRACT(YEAR FROM NOW())::INTEGER) * 12 - EXTRACT(MONTH FROM NOW())::INTEGER + 8 as months_remaining,
			COALESCE(e.purchase_price, 0) as purchase_price,
			'AUGUST' as notify_round
		FROM equipments e
		JOIN equipment_models m ON e.model_id = m.id
		JOIN brands b ON m.brand_id = b.id
		JOIN departments d ON e.department_id = d.id
		WHERE e.deleted_at IS NULL
		AND e.replacement_year = ?
		AND (e.replacement_year - EXTRACT(YEAR FROM NOW())::INTEGER) * 12 - EXTRACT(MONTH FROM NOW())::INTEGER + 8 BETWEEN ? AND ?
		ORDER BY e.id_code ASC
	`

	err := r.db.WithContext(ctx).Raw(query, targetYear, monthsRangeMin, monthsRangeMax).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get August alerts: %w", err)
	}

	log.Printf("📊 August Alert: Found %d equipments for year %d (range: %d-%d months)",
		len(results), targetYear, monthsRangeMin, monthsRangeMax)
	for _, r := range results {
		log.Printf("  ✓ %s: year=%d, months=%d", r.IDCode, r.ReplacementYear, r.MonthsRemaining)
	}

	return results, nil
}

// June Alert - แจ้งทุกเครื่องที่เหลือ 1-3 เดือน
func (r *NotificationRepository) GetEquipmentsForJuneAlert(ctx context.Context) ([]dto.EquipmentReplacementAlertDTO, error) {
	var results []dto.EquipmentReplacementAlertDTO

	currentYear := time.Now().Year()

	// TEST: แจ้งเครื่องปีหน้า (2026) เพื่อทดสอบ
	// targetYear := currentYear + 1
	// monthsRangeMin := 7
	// monthsRangeMax := 9

	// PRODUCTION: แจ้งเครื่องที่ครบปีนี้ (เหลือ 2 เดือน)
	targetYear := currentYear
	monthsRangeMin := 1
	monthsRangeMax := 3

	query := `
		SELECT 
			e.id as equipment_id,
			e.id_code,
			COALESCE(e.serial_no, '') as serial_no,
			b.name as brand_name,
			m.model_name,
			d.name as department_name,
			e.replacement_year,
			(e.replacement_year - EXTRACT(YEAR FROM NOW())::INTEGER) * 12 - EXTRACT(MONTH FROM NOW())::INTEGER + 8 as months_remaining,
			COALESCE(e.purchase_price, 0) as purchase_price,
			'JUNE' as notify_round
		FROM equipments e
		JOIN equipment_models m ON e.model_id = m.id
		JOIN brands b ON m.brand_id = b.id
		JOIN departments d ON e.department_id = d.id
		WHERE e.deleted_at IS NULL
		AND e.replacement_year = ?
		AND (e.replacement_year - EXTRACT(YEAR FROM NOW())::INTEGER) * 12 - EXTRACT(MONTH FROM NOW())::INTEGER + 8 BETWEEN ? AND ?
		ORDER BY e.id_code ASC
	`

	err := r.db.WithContext(ctx).Raw(query, targetYear, monthsRangeMin, monthsRangeMax).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get June alerts: %w", err)
	}

	log.Printf("📊 June Alert: Found %d equipments for year %d (range: %d-%d months)",
		len(results), targetYear, monthsRangeMin, monthsRangeMax)
	for _, r := range results {
		log.Printf("  ✓ %s: year=%d, months=%d", r.IDCode, r.ReplacementYear, r.MonthsRemaining)
	}

	return results, nil
}
