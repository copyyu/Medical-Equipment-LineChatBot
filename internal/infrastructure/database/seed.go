package database

import (
	"log"
	"medical-webhook/internal/domain/line/entity"
	"time"
)

// SeedMockData inserts mock data for testing
func SeedMockData() error {
	log.Println("🌱 Seeding mock data...")

	// Check if data already exists
	var count int64
	DB.Model(&entity.Equipment{}).Count(&count)
	if count > 0 {
		log.Println("✅ Data already exists, skipping seed")
		return nil
	}

	// Create Brands
	brands := []entity.Brand{
		{Name: "D-Link"},
		{Name: "Philips"},
		{Name: "GE Healthcare"},
		{Name: "Siemens"},
	}
	if err := DB.Create(&brands).Error; err != nil {
		log.Printf("❌ Failed to seed brands: %v", err)
		return err
	}
	log.Printf("✅ Created %d brands", len(brands))

	// Create Departments
	departments := []entity.Department{
		{Name: "ห้องฉุกเฉิน"},
		{Name: "ศัลยกรรม"},
		{Name: "อายุรกรรม"},
		{Name: "กุมารเวช"},
	}
	if err := DB.Create(&departments).Error; err != nil {
		log.Printf("❌ Failed to seed departments: %v", err)
		return err
	}
	log.Printf("✅ Created %d departments", len(departments))

	// Create Categories
	categories := []entity.EquipmentCategory{
		{Name: "กล้องวงจรปิด"},
		{Name: "เครื่องตรวจวัด"},
		{Name: "เครื่องมือผ่าตัด"},
	}
	if err := DB.Create(&categories).Error; err != nil {
		log.Printf("❌ Failed to seed categories: %v", err)
		return err
	}
	log.Printf("✅ Created %d categories", len(categories))

	// Create Equipment Models
	models := []entity.EquipmentModel{
		{BrandID: 1, CategoryID: 1, ModelName: "DCS-942L", DefaultLifeExpectancy: 5},
		{BrandID: 2, CategoryID: 2, ModelName: "IntelliVue MX800", DefaultLifeExpectancy: 10},
		{BrandID: 3, CategoryID: 3, ModelName: "OEC 9900", DefaultLifeExpectancy: 15},
	}
	if err := DB.Create(&models).Error; err != nil {
		log.Printf("❌ Failed to seed models: %v", err)
		return err
	}
	log.Printf("✅ Created %d models", len(models))

	// Create Equipment with various test serial numbers
	receiveDate := time.Now().AddDate(-2, 0, 0) // 2 years ago
	computeDate := time.Now()
	replacementYear := 2030

	equipments := []entity.Equipment{
		{
			IDCode:                "ABCD123456789",
			SerialNo:              "SN-ABCD123456789",
			ModelID:               1,
			DepartmentID:          1,
			ReceiveDate:           &receiveDate,
			PurchasePrice:         45000.00,
			EquipmentAge:          2.0,
			ComputeDate:           &computeDate,
			LifeExpectancy:        5.0,
			RemainLife:            3.0,
			UsefulLifetimePercent: 40.0,
			ReplacementYear:       &replacementYear,
		},
		{
			IDCode:                "8001",
			SerialNo:              "SN-8001",
			ModelID:               2,
			DepartmentID:          2,
			ReceiveDate:           &receiveDate,
			PurchasePrice:         250000.00,
			EquipmentAge:          2.0,
			ComputeDate:           &computeDate,
			LifeExpectancy:        10.0,
			RemainLife:            8.0,
			UsefulLifetimePercent: 20.0,
			ReplacementYear:       &replacementYear,
		},
		{
			IDCode:                "TEST001",
			SerialNo:              "SN-TEST001",
			ModelID:               3,
			DepartmentID:          3,
			ReceiveDate:           &receiveDate,
			PurchasePrice:         1500000.00,
			EquipmentAge:          2.0,
			ComputeDate:           &computeDate,
			LifeExpectancy:        15.0,
			RemainLife:            13.0,
			UsefulLifetimePercent: 13.33,
			ReplacementYear:       &replacementYear,
		},
	}
	if err := DB.Create(&equipments).Error; err != nil {
		log.Printf("❌ Failed to seed equipments: %v", err)
		return err
	}
	log.Printf("✅ Created %d equipments", len(equipments))

	// Create Maintenance Records
	maintenanceDate := time.Now().AddDate(0, -6, 0) // 6 months ago
	records := []entity.MaintenanceRecord{
		{
			EquipmentID:     1,
			MaintenanceType: entity.MaintenancePM,
			MaintenanceDate: maintenanceDate,
			Cost:            5000.00,
			Description:     "ตรวจเช็คประจำปี ปรับแต่ง focus",
			Technician:      "สมชาย ใจดี",
		},
		{
			EquipmentID:     1,
			MaintenanceType: entity.MaintenanceCM,
			MaintenanceDate: maintenanceDate.AddDate(0, -3, 0),
			Cost:            12000.00,
			Description:     "เปลี่ยนเลนส์กล้อง",
			Technician:      "สมหญิง รักงาน",
		},
		{
			EquipmentID:     2,
			MaintenanceType: entity.MaintenancePM,
			MaintenanceDate: maintenanceDate,
			Cost:            15000.00,
			Description:     "Calibrate sensors ตามมาตรฐาน",
			Technician:      "สมชาย ใจดี",
		},
	}
	if err := DB.Create(&records).Error; err != nil {
		log.Printf("❌ Failed to seed maintenance records: %v", err)
		return err
	}
	log.Printf("✅ Created %d maintenance records", len(records))

	log.Println("🎉 Mock data seeded successfully!")
	log.Println("📋 Test IDs: ABCD123456789, 8001, TEST001")
	return nil
}
