// File: backend/scripts/run_migrations.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	// Database connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://admindb:gonsgarage123@localhost:5432/gonsgarage?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Create migrations table if it doesn't exist
	createMigrationsTable(db)

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	fmt.Println("✅ All migrations completed successfully!")
}

func createMigrationsTable(db *sql.DB) {
	query := `
    CREATE TABLE IF NOT EXISTS schema_migrations (
        version VARCHAR(255) PRIMARY KEY,
        applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

	if _, err := db.Exec(query); err != nil {
		log.Fatal("Failed to create schema_migrations table:", err)
	}
}

func runMigrations(db *sql.DB) error {
	// Get applied migrations
	appliedMigrations := getAppliedMigrations(db)

	// Get migration files
	migrationFiles, err := getMigrationFiles()
	if err != nil {
		return err
	}

	// Run pending migrations
	for _, file := range migrationFiles {
		if !strings.HasSuffix(file, ".up.sql") {
			continue
		}

		version := strings.TrimSuffix(file, ".up.sql")

		if _, exists := appliedMigrations[version]; exists {
			fmt.Printf("⏭️  Skipping already applied migration: %s\n", version)
			continue
		}

		fmt.Printf("🔄 Running migration: %s\n", version)

		if err := runMigration(db, file, version); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", version, err)
		}

		fmt.Printf("✅ Completed migration: %s\n", version)
	}

	return nil
}

func getAppliedMigrations(db *sql.DB) map[string]bool {
	applied := make(map[string]bool)

	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return applied
	}
	defer rows.Close()

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err == nil {
			applied[version] = true
		}
	}

	return applied
}

func getMigrationFiles() ([]string, error) {
	migrationsDir := "migrations"

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}

	sort.Strings(migrationFiles)
	return migrationFiles, nil
}

func runMigration(db *sql.DB, filename, version string) error {
	// Read migration file
	content, err := os.ReadFile(filepath.Join("migrations", filename))
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Execute migration
	if _, err := tx.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	// Record migration as applied
	if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration: %w", err)
	}

	return nil
}
