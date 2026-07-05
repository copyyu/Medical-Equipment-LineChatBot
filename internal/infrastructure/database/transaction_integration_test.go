//go:build integration

// Integration tests for the transaction manager against a real PostgreSQL
// instance. They are excluded from the normal build/test run and only compile
// and run under the "integration" tag:
//
//	TEST_DATABASE_URL=postgres://user:pass@localhost:5432/db?sslmode=disable \
//	    go test -tags=integration ./internal/infrastructure/database/...
//
// If TEST_DATABASE_URL is unset the tests skip, so `make test-integration`
// is safe to run without a database.
package database

import (
	"context"
	"errors"
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type txTestRow struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set; skipping integration test")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	if err := db.AutoMigrate(&txTestRow{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	t.Cleanup(func() { _ = db.Migrator().DropTable(&txTestRow{}) })
	return db
}

func TestWithTransaction_CommitsOnSuccess(t *testing.T) {
	db := openTestDB(t)
	m := NewTxManager(db)

	err := m.WithTransaction(context.Background(), func(ctx context.Context) error {
		return DBFromContext(ctx, db).Create(&txTestRow{Name: "keep"}).Error
	})
	if err != nil {
		t.Fatalf("WithTransaction returned error: %v", err)
	}

	var count int64
	db.Model(&txTestRow{}).Where("name = ?", "keep").Count(&count)
	if count != 1 {
		t.Fatalf("expected committed row, got count=%d", count)
	}
}

func TestWithTransaction_RollsBackOnError(t *testing.T) {
	db := openTestDB(t)
	m := NewTxManager(db)
	sentinel := errors.New("boom")

	err := m.WithTransaction(context.Background(), func(ctx context.Context) error {
		if e := DBFromContext(ctx, db).Create(&txTestRow{Name: "drop"}).Error; e != nil {
			return e
		}
		return sentinel // triggers rollback
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error to propagate, got: %v", err)
	}

	var count int64
	db.Model(&txTestRow{}).Where("name = ?", "drop").Count(&count)
	if count != 0 {
		t.Fatalf("expected rolled-back row to be absent, got count=%d", count)
	}
}
