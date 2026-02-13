package database

import (
	"log"
	"medical-webhook/internal/domain/line/entity"

	"gorm.io/gorm"
)

// SeedTicketCategories seeds default ticket categories if they don't exist
func SeedTicketCategories(db *gorm.DB) {
	var count int64
	db.Model(&entity.TicketCategory{}).Count(&count)
	if count > 0 {
		return
	}

	categories := []entity.TicketCategory{
		{Name: "แจ้งซ่อม", Color: "#EF5350", Icon: "🔧", SortOrder: 1},
		{Name: "บำรุงรักษา", Color: "#FFA726", Icon: "🛠️", SortOrder: 2},
		{Name: "สอบถามการใช้งาน", Color: "#42A5F5", Icon: "❓", SortOrder: 3},
		{Name: "อื่นๆ", Color: "#78909C", Icon: "📝", SortOrder: 4},
	}

	if err := db.Create(&categories).Error; err != nil {
		log.Printf("Failed to seed categories: %v", err)
	} else {
		log.Println("Default ticket categories seeded successfully.")
	}
}
