package templates

// MigrationTemplate generates database migration files
const MigrationTemplate = `-- Migration: {{.Version}}_{{.Name}}
-- Generated at: {{.CreatedAt}}
-- Description: {{.Description}}

-- +migrate Up
{{.UpSQL}}

-- +migrate Down
{{.DownSQL}}
`

// MigrationRunnerTemplate generates the migration runner
const MigrationRunnerTemplate = `package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/joho/godotenv"
	{{- if eq .Database.Type "postgres"}}
	_ "github.com/lib/pq"
	{{- else if eq .Database.Type "mysql"}}
	_ "github.com/go-sql-driver/mysql"
	{{- else if eq .Database.Type "sqlite"}}
	_ "github.com/mattn/go-sqlite3"
	{{- end}}
)

// Migration represents a database migration
type Migration struct {
	Version   string
	Name      string
	FilePath  string
	UpSQL     string
	DownSQL   string
	AppliedAt *time.Time
}

// MigrationRunner handles database migrations
type MigrationRunner struct {
	db            *sql.DB
	migrationsDir string
}

// NewMigrationRunner creates a new migration runner
func NewMigrationRunner(db *sql.DB, migrationsDir string) *MigrationRunner {
	return &MigrationRunner{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

// Run executes migrations
func (mr *MigrationRunner) Run() error {
	// Create migrations table if not exists
	if err := mr.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get pending migrations
	migrations, err := mr.getPendingMigrations()
	if err != nil {
		return fmt.Errorf("failed to get pending migrations: %w", err)
	}

	if len(migrations) == 0 {
		fmt.Println("✅ No pending migrations")
		return nil
	}

	// Execute migrations
	for _, migration := range migrations {
		if err := mr.executeMigration(migration); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", migration.Version, err)
		}
		fmt.Printf("✅ Applied migration: %s_%s\n", migration.Version, migration.Name)
	}

	fmt.Printf("🎉 Applied %d migrations successfully\n", len(migrations))
	return nil
}

// Rollback rolls back the last migration
func (mr *MigrationRunner) Rollback() error {
	// Get last applied migration
	migration, err := mr.getLastAppliedMigration()
	if err != nil {
		return fmt.Errorf("failed to get last migration: %w", err)
	}

	if migration == nil {
		fmt.Println("ℹ️  No migrations to rollback")
		return nil
	}

	// Execute down migration
	if err := mr.executeDownMigration(migration); err != nil {
		return fmt.Errorf("failed to rollback migration %s: %w", migration.Version, err)
	}

	fmt.Printf("↩️  Rolled back migration: %s_%s\n", migration.Version, migration.Name)
	return nil
}

// createMigrationsTable creates the migrations tracking table
func (mr *MigrationRunner) createMigrationsTable() error {
	query := ` + "`" + `
		CREATE TABLE IF NOT EXISTS migrations (
			version VARCHAR(14) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	` + "`" + `

	_, err := mr.db.Exec(query)
	return err
}

// getPendingMigrations returns migrations that haven't been applied
func (mr *MigrationRunner) getPendingMigrations() ([]*Migration, error) {
	// Get all migration files
	files, err := filepath.Glob(filepath.Join(mr.migrationsDir, "*.sql"))
	if err != nil {
		return nil, err
	}

	// Get applied migrations
	appliedVersions, err := mr.getAppliedVersions()
	if err != nil {
		return nil, err
	}

	var pendingMigrations []*Migration
	for _, file := range files {
		migration, err := mr.parseMigrationFile(file)
		if err != nil {
			continue // Skip invalid files
		}

		// Check if already applied
		if _, applied := appliedVersions[migration.Version]; !applied {
			pendingMigrations = append(pendingMigrations, migration)
		}
	}

	// Sort by version
	sort.Slice(pendingMigrations, func(i, j int) bool {
		return pendingMigrations[i].Version < pendingMigrations[j].Version
	})

	return pendingMigrations, nil
}

// getAppliedVersions returns a map of applied migration versions
func (mr *MigrationRunner) getAppliedVersions() (map[string]bool, error) {
	rows, err := mr.db.Query("SELECT version FROM migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	versions := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		versions[version] = true
	}

	return versions, nil
}

// parseMigrationFile parses a migration file
func (mr *MigrationRunner) parseMigrationFile(filePath string) (*Migration, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	filename := filepath.Base(filePath)
	parts := strings.SplitN(filename, "_", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid migration filename: %s", filename)
	}

	version := parts[0]
	name := strings.TrimSuffix(parts[1], ".sql")

	// Parse up and down SQL
	contentStr := string(content)
	upSQL, downSQL := mr.parseUpDown(contentStr)

	return &Migration{
		Version:  version,
		Name:     name,
		FilePath: filePath,
		UpSQL:    upSQL,
		DownSQL:  downSQL,
	}, nil
}

// parseUpDown parses up and down SQL from migration content
func (mr *MigrationRunner) parseUpDown(content string) (string, string) {
	lines := strings.Split(content, "\n")
	var upSQL, downSQL strings.Builder
	var currentSection string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		if strings.Contains(line, "+migrate Up") {
			currentSection = "up"
			continue
		} else if strings.Contains(line, "+migrate Down") {
			currentSection = "down"
			continue
		}

		if strings.HasPrefix(line, "--") {
			continue // Skip comments
		}

		if currentSection == "up" {
			upSQL.WriteString(line)
			upSQL.WriteString("\n")
		} else if currentSection == "down" {
			downSQL.WriteString(line)
			downSQL.WriteString("\n")
		}
	}

	return strings.TrimSpace(upSQL.String()), strings.TrimSpace(downSQL.String())
}

// executeMigration executes an up migration
func (mr *MigrationRunner) executeMigration(migration *Migration) error {
	// Execute up SQL
	if _, err := mr.db.Exec(migration.UpSQL); err != nil {
		return err
	}

	// Record migration
	_, err := mr.db.Exec(
		"INSERT INTO migrations (version, name) VALUES ($1, $2)",
		migration.Version, migration.Name,
	)
	return err
}

// getLastAppliedMigration returns the most recently applied migration
func (mr *MigrationRunner) getLastAppliedMigration() (*Migration, error) {
	row := mr.db.QueryRow(` + "`" + `
		SELECT version, name, applied_at 
		FROM migrations 
		ORDER BY version DESC 
		LIMIT 1
	` + "`" + `)

	var version, name string
	var appliedAt time.Time
	err := row.Scan(&version, &name, &appliedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Find the migration file to get down SQL
	filePath := filepath.Join(mr.migrationsDir, fmt.Sprintf("%s_%s.sql", version, name))
	migration, err := mr.parseMigrationFile(filePath)
	if err != nil {
		return nil, err
	}

	migration.AppliedAt = &appliedAt
	return migration, nil
}

// executeDownMigration executes a down migration
func (mr *MigrationRunner) executeDownMigration(migration *Migration) error {
	// Execute down SQL
	if _, err := mr.db.Exec(migration.DownSQL); err != nil {
		return err
	}

	// Remove migration record
	_, err := mr.db.Exec("DELETE FROM migrations WHERE version = $1", migration.Version)
	return err
}

func main() {
	// Load environment
	godotenv.Load()

	// Get database URL
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Connect to database
	db, err := sql.Open("{{.Database.GetDriverName}}", databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Get migrations directory
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		migrationsDir = "./migrations"
	}

	// Create migration runner
	runner := NewMigrationRunner(db, migrationsDir)

	// Get command
	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "up", "migrate":
		if err := runner.Run(); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
	case "down", "rollback":
		if err := runner.Rollback(); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
	default:
		fmt.Println("Usage: migrate [up|down]")
		fmt.Println("  up:   Run pending migrations (default)")
		fmt.Println("  down: Rollback last migration")
	}
}
`

// MigrationMakefileTemplate generates migration-related Makefile targets
const MigrationMakefileTemplate = `
# Migration commands
.PHONY: migrate migrate-up migrate-down migrate-create migrate-status

migrate: migrate-up

migrate-up:
	@echo "Running migrations..."
	@go run cmd/migrate/main.go up

migrate-down:
	@echo "Rolling back last migration..."
	@go run cmd/migrate/main.go down

migrate-create:
	@read -p "Migration name: " name; \
	timestamp=$$(date +%Y%m%d%H%M%S); \
	mkdir -p migrations; \
	touch migrations/$${timestamp}_$${name}.sql; \
	echo "Created migration: migrations/$${timestamp}_$${name}.sql"

migrate-status:
	@echo "Migration status:"
	@go run cmd/migrate/main.go status 2>/dev/null || echo "No migrations table found"
`