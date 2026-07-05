// database/database.go
package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"medical-webhook/internal/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB    *gorm.DB
	SqlDB *sql.DB
)

func Connect(cfg *config.Config) error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Translate driver errors into gorm sentinels (e.g. ErrDuplicatedKey) so
		// repositories can detect unique-constraint violations portably.
		TranslateError: true,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to Postgres: %w", err)
	}

	// Get underlying *sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Connection pool settings (production-grade)
	sqlDB.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.DB.ConnMaxLifetime)

	// Store in global variables
	DB = db
	SqlDB = sqlDB

	log.Println("Connected to PostgreSQL successfully!")

	// NOTE: Schema management is now handled by golang-migrate.
	// Run migrations with: go run ./cmd/migrate up
	// DO NOT use AutoMigrate in production.

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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return SqlDB.PingContext(ctx)
}
