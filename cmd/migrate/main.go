package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env if present
	_ = godotenv.Load()

	// Parse flags
	dir := flag.String("dir", "migrations", "Directory containing migration files")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		printUsage()
		os.Exit(1)
	}

	command := args[0]

	// Build database URL from environment
	dbURL := buildDatabaseURL()

	// Create migration instance
	sourceURL := fmt.Sprintf("file://%s", *dir)
	m, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		log.Fatalf("❌ Failed to create migrate instance: %v", err)
	}
	defer m.Close()

	switch command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("❌ Migration up failed: %v", err)
		}
		version, dirty, _ := m.Version()
		log.Printf("✅ Migration up complete (version: %d, dirty: %v)", version, dirty)

	case "down":
		steps := 1
		if len(args) > 1 && args[1] == "--all" {
			if err := m.Down(); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("❌ Migration down failed: %v", err)
			}
			log.Println("✅ All migrations rolled back")
		} else {
			if err := m.Steps(-steps); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("❌ Migration down failed: %v", err)
			}
			version, dirty, _ := m.Version()
			log.Printf("✅ Rolled back %d step(s) (version: %d, dirty: %v)", steps, version, dirty)
		}

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("❌ Failed to get version: %v", err)
		}
		fmt.Printf("Version: %d, Dirty: %v\n", version, dirty)

	case "force":
		if len(args) < 2 {
			log.Fatal("❌ Usage: migrate force <version>")
		}
		var version int
		_, err := fmt.Sscanf(args[1], "%d", &version)
		if err != nil {
			log.Fatalf("❌ Invalid version: %v", err)
		}
		if err := m.Force(version); err != nil {
			log.Fatalf("❌ Force failed: %v", err)
		}
		log.Printf("✅ Forced to version %d", version)

	case "drop":
		fmt.Print("⚠️  This will drop all tables. Are you sure? (yes/no): ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "yes" {
			fmt.Println("Aborted.")
			return
		}
		if err := m.Drop(); err != nil {
			log.Fatalf("❌ Drop failed: %v", err)
		}
		log.Println("✅ All tables dropped")

	default:
		fmt.Printf("❌ Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func buildDatabaseURL() string {
	host := getEnvOrDefault("DB_HOST", "localhost")
	port := getEnvOrDefault("DB_PORT", "5432")
	user := getEnvOrDefault("DB_USER", "postgres")
	password := getEnvOrDefault("DB_PASSWORD", "postgres")
	dbname := getEnvOrDefault("DB_NAME", "medical_equipment")
	sslmode := getEnvOrDefault("DB_SSLMODE", "disable")

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode,
	)
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func printUsage() {
	fmt.Println(`Usage: go run ./cmd/migrate <command>

Commands:
  up        Apply all pending migrations
  down      Roll back the last migration (use --all to roll back all)
  version   Show current migration version
  force N   Force migration version to N (fix dirty state)
  drop      Drop all tables (interactive confirmation)

Flags:
  -dir      Directory containing migration files (default: "migrations")

Environment variables:
  DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE`)
}
