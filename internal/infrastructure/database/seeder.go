package database

import (
	"log"
	"medical-webhook/internal/domain/line/entity"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedTicketCategories seeds default ticket categories if they don't exist
func SeedTicketCategories(db *gorm.DB) {
	var count int64
	if err := db.Model(&entity.TicketCategory{}).Count(&count).Error; err != nil {
		// Bail out on a count error instead of proceeding to seed, which could
		// insert duplicate categories on a transient failure.
		log.Printf("Failed to count ticket categories, skipping seed: %v", err)
		return
	}
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

// SeedBootstrapAdmin creates the first super-admin from the provided credentials
// when no admin exists yet. It is idempotent (no-op once any admin exists) and
// safe to call on every startup. If no admin exists and no credentials are
// provided, it logs guidance since /api/admin/register is restricted to
// super-admins.
func SeedBootstrapAdmin(db *gorm.DB, username, password, email, fullName string) {
	var count int64
	if err := db.Model(&entity.Admin{}).Count(&count).Error; err != nil {
		log.Printf("bootstrap admin: could not count admins: %v", err)
		return
	}
	if count > 0 {
		return // admins already exist
	}

	if username == "" || password == "" {
		log.Println("⚠️ No admin exists. Set ADMIN_BOOTSTRAP_USERNAME and " +
			"ADMIN_BOOTSTRAP_PASSWORD to create the first super-admin " +
			"(/api/admin/register is restricted to super-admins).")
		return
	}

	if email == "" {
		email = username + "@local"
	}
	if fullName == "" {
		fullName = "Super Admin"
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("bootstrap admin: failed to hash password: %v", err)
		return
	}

	admin := &entity.Admin{
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		FullName:     fullName,
		Role:         string(entity.RoleSuperAdmin),
	}
	if err := db.Create(admin).Error; err != nil {
		log.Printf("bootstrap admin: failed to create super-admin: %v", err)
		return
	}
	log.Printf("✅ Bootstrapped super-admin '%s'", username)
}
