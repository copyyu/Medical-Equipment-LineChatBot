// database/database.go
package database

import (
	"database/sql"
	"fmt"
	"log"
	"medical-webhook/internal/config"
	"medical-webhook/internal/domain/line/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB    *gorm.DB
	SqlDB *sql.DB
)

func Connect(cfg *config.Config) error {

	// if cfg.DB.Host == "" || cfg.DB.User == "" || cfg.DB.Name == "" || cfg.DB.Port == "" {
	// 	return fmt.Errorf("database configuration is incomplete: host=%s, user=%s, dbname=%s, port=%s",
	// 		cfg.DB.Host, cfg.DB.User, cfg.DB.Name, cfg.DB.Port)
	// }
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to Postgres: %w", err)
	}

	// Get underlying *sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Store in global variables
	DB = db
	SqlDB = sqlDB

	log.Println("Connected to PostgreSQL successfully!")

	err = db.AutoMigrate(
		&entity.Brand{},
		&entity.Department{},
		&entity.EquipmentCategory{},
		&entity.EquipmentModel{},
		&entity.Equipment{},
		&entity.MaintenanceRecord{},
		&entity.NotificationLog{},
		&entity.NotificationSetting{},
		&entity.Admin{},
		&entity.AdminSession{},
		&entity.Ticket{},
		&entity.TicketCategory{},
		&entity.TicketHistory{},
	)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	SeedTicketCategories(db)

	return nil
}

func Close() error {
	if SqlDB != nil {
		log.Println("Closing database connection...")
		return SqlDB.Close()
	}
	return nil
}

func GetDB() *gorm.DB {
	return DB
}

func GetSqlDB() *sql.DB {
	return SqlDB
}

func HealthCheck() error {
	if SqlDB == nil {
		return fmt.Errorf("database connection is nil")
	}
	return SqlDB.Ping()
}
