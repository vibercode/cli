package templates

// SupabaseDatabaseTemplate provides database connection setup for Supabase
const SupabaseDatabaseTemplate = `package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	gotrue "github.com/supabase-community/gotrue-go"
	storage "github.com/supabase-community/storage-go"
)

var (
	DB      *gorm.DB
	Auth    *gotrue.Client
	Storage *storage.Client
)

// DatabaseConfig holds the database configuration
type DatabaseConfig struct {
	URL        string
	AnonKey    string
	ServiceKey string
	JWTSecret  string
}

// Connect establishes connection to Supabase
func Connect() error {
	config := DatabaseConfig{
		URL:        os.Getenv("SUPABASE_URL"),
		AnonKey:    os.Getenv("SUPABASE_ANON_KEY"),
		ServiceKey: os.Getenv("SUPABASE_SERVICE_KEY"),
		JWTSecret:  os.Getenv("SUPABASE_JWT_SECRET"),
	}

	if err := validateConfig(config); err != nil {
		return fmt.Errorf("invalid supabase configuration: %w", err)
	}

	// Connect to PostgreSQL database
	if err := connectDatabase(config); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Initialize Auth client
	if err := initAuthClient(config); err != nil {
		return fmt.Errorf("failed to initialize auth client: %w", err)
	}

	// Initialize Storage client
	if err := initStorageClient(config); err != nil {
		return fmt.Errorf("failed to initialize storage client: %w", err)
	}

	log.Println("✅ Connected to Supabase successfully")
	return nil
}

// connectDatabase connects to the PostgreSQL database
func connectDatabase(config DatabaseConfig) error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is required")
	}

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	DB = db
	return nil
}

// initAuthClient initializes the Supabase Auth client
func initAuthClient(config DatabaseConfig) error {
	client := gotrue.New(config.URL, config.ServiceKey)
	Auth = client
	return nil
}

// initStorageClient initializes the Supabase Storage client
func initStorageClient(config DatabaseConfig) error {
	client := storage.NewClient(config.URL, config.ServiceKey, nil)
	Storage = client
	return nil
}

// validateConfig validates the Supabase configuration
func validateConfig(config DatabaseConfig) error {
	if config.URL == "" {
		return fmt.Errorf("SUPABASE_URL is required")
	}
	if config.AnonKey == "" {
		return fmt.Errorf("SUPABASE_ANON_KEY is required")
	}
	if config.ServiceKey == "" {
		return fmt.Errorf("SUPABASE_SERVICE_KEY is required")
	}
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// GetAuth returns the Auth client
func GetAuth() *gotrue.Client {
	return Auth
}

// GetStorage returns the Storage client
func GetStorage() *storage.Client {
	return Storage
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
`

// SupabaseEnvTemplate provides environment variables template for Supabase
const SupabaseEnvTemplate = `# Supabase Configuration
SUPABASE_URL=https://{{.Database.ProjectRef}}.supabase.co
SUPABASE_ANON_KEY={{.Database.AnonKey}}
SUPABASE_SERVICE_KEY={{.Database.ServiceKey}}
SUPABASE_JWT_SECRET={{.Database.JWTSecret}}

# Database Connection (for direct PostgreSQL access)
DATABASE_URL={{.Database.GetDSN}}

# Application Settings
{{range $key, $value := .Database.GetEnvironmentVars}}
{{$key}}={{$value}}
{{end}}

# Server Configuration
PORT={{.Port}}
GIN_MODE=release

# CORS Configuration
CORS_ORIGINS=http://localhost:3000,https://{{.Database.ProjectRef}}.supabase.co
`