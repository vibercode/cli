package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/vibercode/cli/internal/models"
	"github.com/vibercode/cli/internal/templates"
	"github.com/vibercode/cli/pkg/ui"
)

// APIGenerator handles API project generation
type APIGenerator struct{}

// NewAPIGenerator creates a new APIGenerator
func NewAPIGenerator() *APIGenerator {
	return &APIGenerator{}
}

// APIProject represents an API project configuration
type APIProject struct {
	Name     string
	Port     string
	Database *models.DatabaseProvider
	Module   string
}

// VibercodeManifest represents the project configuration that can be used to regenerate the project
type VibercodeManifest struct {
	Version     string                    `json:"version"`
	ProjectType string                    `json:"project_type"`
	Name        string                    `json:"name"`
	Port        string                    `json:"port"`
	Database    *models.DatabaseProvider  `json:"database"`
	Module      string                    `json:"module"`
	GeneratedAt string                    `json:"generated_at"`
	UpdatedAt   string                    `json:"updated_at,omitempty"`
	CLI         VibercodeManifestCLI      `json:"cli"`
	History     []VibercodeManifestEvent  `json:"history,omitempty"`
	Resources   []VibercodeResource       `json:"resources,omitempty"`
}

// VibercodeManifestCLI represents CLI-specific information
type VibercodeManifestCLI struct {
	Version string `json:"version"`
	Command string `json:"command"`
}

// VibercodeManifestEvent represents a change/generation event
type VibercodeManifestEvent struct {
	Type        string `json:"type"`        // "create", "generate_resource", "update"
	Description string `json:"description"`
	Timestamp   string `json:"timestamp"`
	CLI         VibercodeManifestCLI `json:"cli"`
}

// VibercodeResource represents a generated resource
type VibercodeResource struct {
	Name        string                   `json:"name"`
	Type        string                   `json:"type"`        // "crud", "model", "handler", etc.
	Fields      []VibercodeResourceField `json:"fields,omitempty"`
	GeneratedAt string                   `json:"generated_at"`
}

// VibercodeResourceField represents a field in a resource
type VibercodeResourceField struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

// Generate generates a complete API project
func (g *APIGenerator) Generate() error {
	ui.PrintHeader("API Project Generator")
	
	project, err := g.collectProjectInfo()
	if err != nil {
		return fmt.Errorf("failed to collect project info: %w", err)
	}

	// Show project summary
	ui.PrintSubHeader("Project Configuration")
	ui.PrintFeature(ui.IconAPI, "Project Name", project.Name)
	ui.PrintFeature(ui.IconGear, "Port", project.Port)
	ui.PrintFeature(ui.IconDatabase, "Database", project.Database.GetDisplayName())
	ui.PrintFeature(ui.IconPackage, "Module", project.Module)
	fmt.Println()

	if !ui.ConfirmAction("Generate project with this configuration?") {
		ui.PrintInfo("Project generation cancelled")
		return nil
	}

	// Start generation with spinner
	spinner := ui.ShowSpinner("Creating project structure...")
	time.Sleep(500 * time.Millisecond) // Brief pause for UX

	// Create project directory
	if err := os.MkdirAll(project.Name, 0755); err != nil {
		spinner.Stop()
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Generate project structure
	if err := g.createProjectStructure(project); err != nil {
		spinner.Stop()
		return fmt.Errorf("failed to create project structure: %w", err)
	}
	
	spinner.Stop()
	ui.PrintSuccess("Project structure created")

	// Generate files with progress
	files := []struct {
		name string
		fn   func(*APIProject) error
	}{
		{"go.mod", g.generateGoMod},
		{"main.go", g.generateMain},
		{"database package", g.generateDatabase},
		{"handlers", g.generateHandlers},
		{"Dockerfile", g.generateDockerfile},
		{".env.example", g.generateEnvExample},
		{"docker-compose.yml", g.generateDockerCompose},
		{"Makefile", g.generateMakefile},
		{"README.md", g.generateReadme},
	}

	ui.PrintStep(1, 2, "Generating project files")
	for _, file := range files {
		spinner := ui.ShowSpinner(fmt.Sprintf("Generating %s...", file.name))
		time.Sleep(200 * time.Millisecond) // Brief pause for UX
		
		if err := file.fn(project); err != nil {
			spinner.Stop()
			return fmt.Errorf("failed to generate %s: %w", file.name, err)
		}
		
		spinner.Stop()
		ui.PrintFileCreated(file.name)
	}

	ui.PrintStep(2, 2, "Finalizing project")
	time.Sleep(300 * time.Millisecond)

	// Generate manifest file
	if err := g.generateManifest(project); err != nil {
		ui.PrintWarning(fmt.Sprintf("Failed to generate manifest file: %v", err))
	}

	// Show success message
	fmt.Println()
	ui.PrintSuccess(fmt.Sprintf("API project '%s' generated successfully!", project.Name))
	
	// Show project structure
	ui.PrintProjectStructure(project.Name)
	
	// Show database info
	ui.PrintDatabaseInfo(project.Database.Type, project.Name)
	
	// Show next steps
	ui.PrintNextSteps(project.Name)

	return nil
}

// collectProjectInfo collects project information from user input
func (g *APIGenerator) collectProjectInfo() (*APIProject, error) {
	project := &APIProject{}

	ui.PrintInfo("Let's configure your new Go API project")
	fmt.Println()

	// Project name
	name, err := ui.TextInput(ui.IconAPI + " Project name:")
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(name) == "" {
		ui.ExitWithError("Project name cannot be empty")
	}
	project.Name = strings.TrimSpace(name)

	// Port
	port, err := ui.TextInput(ui.IconGear + " API port:", "8080")
	if err != nil {
		return nil, err
	}
	project.Port = strings.TrimSpace(port)

	// Database
	dbType, err := ui.SelectOption(ui.IconDatabase + " Database type:", models.SupportedDatabaseTypes())
	if err != nil {
		return nil, err
	}
	
	// Create database provider with basic configuration
	project.Database = &models.DatabaseProvider{
		Type: dbType,
		Host: "localhost",
		Port: getDefaultPort(dbType),
	}
	
	// Set database-specific defaults
	switch dbType {
	case "postgres":
		project.Database.Database = project.Name
		project.Database.Username = "postgres"
		project.Database.Password = "postgres"
		project.Database.SSLMode = "disable"
	case "mysql":
		project.Database.Database = project.Name
		project.Database.Username = "root"
		project.Database.Password = "password"
	case "sqlite":
		project.Database.Database = project.Name + ".db"
	case "supabase":
		project.Database.Database = "postgres"
		project.Database.Username = "postgres"
		project.Database.SSLMode = "require"
	case "mongodb":
		project.Database.Database = project.Name
	case "redis":
		project.Database.Database = "0"
	}

	// Module name
	defaultModule := "github.com/user/" + project.Name
	module, err := ui.TextInput(ui.IconPackage + " Go module name:", defaultModule)
	if err != nil {
		return nil, err
	}
	project.Module = strings.TrimSpace(module)

	return project, nil
}

// getDefaultPort returns the default port for a database type
func getDefaultPort(dbType string) int {
	switch dbType {
	case "postgres", "supabase":
		return 5432
	case "mysql":
		return 3306
	case "mongodb":
		return 27017
	case "redis":
		return 6379
	default:
		return 5432
	}
}

// createProjectStructure creates the basic project directory structure
func (g *APIGenerator) createProjectStructure(project *APIProject) error {
	dirs := []string{
		filepath.Join(project.Name, "cmd", "server"),
		filepath.Join(project.Name, "internal", "handlers"),
		filepath.Join(project.Name, "internal", "services"),
		filepath.Join(project.Name, "internal", "repositories"),
		filepath.Join(project.Name, "internal", "models"),
		filepath.Join(project.Name, "internal", "middleware"),
		filepath.Join(project.Name, "pkg", "database"),
		filepath.Join(project.Name, "pkg", "config"),
		filepath.Join(project.Name, "pkg", "utils"),
		filepath.Join(project.Name, "docs"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// generateGoMod generates the go.mod file
func (g *APIGenerator) generateGoMod(project *APIProject) error {
	template := `module {{.Module}}

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/joho/godotenv v1.4.0
{{- if eq .Database.Type "mongodb"}}
	go.mongodb.org/mongo-driver v1.13.1
{{- else if eq .Database.Type "redis"}}
	github.com/go-redis/redis/v8 v8.11.5
{{- else}}
	gorm.io/gorm v1.25.5
{{- if eq .Database.Type "postgres"}}
	gorm.io/driver/postgres v1.5.4
{{- else if eq .Database.Type "mysql"}}
	gorm.io/driver/mysql v1.5.2
{{- else if eq .Database.Type "sqlite"}}
	gorm.io/driver/sqlite v1.5.4
{{- else if eq .Database.Type "supabase"}}
	gorm.io/driver/postgres v1.5.4
	github.com/supabase-community/gotrue-go v1.0.1
	github.com/supabase-community/storage-go v0.7.0
{{- end}}
{{- end}}
	github.com/google/uuid v1.4.0
)
`
	return g.generateFromTemplate(project, template, filepath.Join(project.Name, "go.mod"))
}

// generateMain generates the main.go file
func (g *APIGenerator) generateMain(project *APIProject) error {
	template := `package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"{{.Module}}/internal/handlers"
	"{{.Module}}/pkg/database"{{- if eq .Database.Type "mongodb"}}
	_ "go.mongodb.org/mongo-driver/mongo"{{- else if eq .Database.Type "redis"}}
	_ "github.com/go-redis/redis/v8"{{- end}}
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to database
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Add middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Setup routes
	api := r.Group("/api/v1")
	handlers.SetupRoutes(api, db)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"service": "{{.Name}}",
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "{{.Port}}"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
`
	return g.generateFromTemplate(project, template, filepath.Join(project.Name, "cmd", "server", "main.go"))
}

// generateDatabase generates the database package
func (g *APIGenerator) generateDatabase(project *APIProject) error {
	var driverImport string
	var driverCall string

	switch project.Database.Type {
	case "postgres":
		driverImport = `"gorm.io/driver/postgres"`
		driverCall = `postgres.Open(dsn)`
	case "mysql":
		driverImport = `"gorm.io/driver/mysql"`
		driverCall = `mysql.Open(dsn)`
	case "sqlite":
		driverImport = `"gorm.io/driver/sqlite"`
		driverCall = `sqlite.Open(dsn)`
	case "supabase":
		driverImport = `"gorm.io/driver/postgres"`
		driverCall = `postgres.Open(dsn)`
	case "mongodb":
		driverImport = `"go.mongodb.org/mongo-driver/mongo"`
		driverCall = `mongo.Connect(context.TODO(), options.Client().ApplyURI(dsn))`
	default:
		driverImport = `"gorm.io/driver/postgres"`
		driverCall = `postgres.Open(dsn)`
	}

	var template string
	
	if project.Database.Type == "supabase" {
		return g.generateSupabaseDatabase(project)
	} else if project.Database.Type == "mongodb" {
		return g.generateMongoDBDatabase(project)
	} else if project.Database.Type == "redis" {
		return g.generateRedisDatabase(project)
	} else {
		template = `package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	` + driverImport + `
)

var DB *gorm.DB

// Connect establishes a database connection
func Connect() (*gorm.DB, error) {
	var err error
	
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = getDefaultDSN()
	}

	DB, err = gorm.Open(` + driverCall + `, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connected successfully")
	return DB, nil
}

// getDefaultDSN returns the default database connection string
func getDefaultDSN() string {
	return "{{call .Database.GetDSN}}"
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
	}
	return g.generateFromTemplate(project, template, filepath.Join(project.Name, "pkg", "database", "database.go"))
}

// generateHandlers generates the handlers setup
func (g *APIGenerator) generateHandlers(project *APIProject) error {
	var template string
	
	if project.Database.Type == "mongodb" {
		template = `package handlers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetupRoutes sets up all API routes
func SetupRoutes(r *gin.RouterGroup, db *mongo.Database) {
	// Initialize repositories
	// userRepo := repositories.NewUserRepository(db)

	// Initialize services
	// userService := services.NewUserService(userRepo)

	// Initialize handlers
	// userHandler := NewUserHandler(userService)

	// Setup routes
	// userHandler.SetupUserRoutes(r)

	// Example endpoint
	r.GET("/example", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "{{.Name}} API is running!",
		})
	})
}
`
	} else {
		template = `package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes sets up all API routes
func SetupRoutes(r *gin.RouterGroup, db *gorm.DB) {
	// Initialize repositories
	// userRepo := repositories.NewUserRepository(db)

	// Initialize services
	// userService := services.NewUserService(userRepo)

	// Initialize handlers
	// userHandler := NewUserHandler(userService)

	// Setup routes
	// userHandler.SetupUserRoutes(r)

	// Example endpoint
	r.GET("/example", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "{{.Name}} API is running!",
		})
	})
}
`
	}
	
	return g.generateFromTemplate(project, template, filepath.Join(project.Name, "internal", "handlers", "routes.go"))
}

// generateDockerfile generates a Dockerfile
func (g *APIGenerator) generateDockerfile(project *APIProject) error {
	template := `FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY go.sum* ./
RUN go mod download

COPY . .
RUN go mod tidy
RUN go build -o main cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE {{.Port}}

CMD ["./main"]
`
	return g.generateFromTemplate(project, template, filepath.Join(project.Name, "Dockerfile"))
}

// generateEnvExample generates .env.example file
func (g *APIGenerator) generateEnvExample(project *APIProject) error {
	template := `# Server Configuration
PORT={{.Port}}

# Database Configuration
{{- if eq .Database.Type "postgres"}}
DATABASE_URL=postgres://username:password@localhost:5432/{{.Name}}?sslmode=disable
{{- else if eq .Database.Type "mysql"}}
DATABASE_URL=username:password@tcp(localhost:3306)/{{.Name}}?charset=utf8mb4&parseTime=True&loc=Local
{{- else if eq .Database.Type "sqlite"}}
DATABASE_URL={{.Name}}.db
{{- else if eq .Database.Type "supabase"}}
DATABASE_URL=postgresql://postgres:[YOUR-PASSWORD]@db.[YOUR-PROJECT-REF].supabase.co:5432/postgres?sslmode=require

# Supabase Configuration
SUPABASE_URL=https://[YOUR-PROJECT-REF].supabase.co
SUPABASE_ANON_KEY=[YOUR-ANON-KEY]
SUPABASE_SERVICE_KEY=[YOUR-SERVICE-KEY]
SUPABASE_JWT_SECRET=[YOUR-JWT-SECRET]
{{- else if eq .Database.Type "mongodb"}}
DATABASE_URL=mongodb://localhost:27017/{{.Name}}
{{- else if eq .Database.Type "redis"}}
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_MAX_RETRIES=3
REDIS_POOL_SIZE=10
{{- end}}

# JWT Configuration
JWT_SECRET=your-secret-key-here

# Other Configuration
GIN_MODE=release
LOG_LEVEL=info
`
	return g.generateFromTemplate(project, template, filepath.Join(project.Name, ".env.example"))
}

// generateDockerCompose generates docker-compose.yml file
func (g *APIGenerator) generateDockerCompose(project *APIProject) error {
	template := `# Docker Compose file for {{.Name}}

services:
  app:
    build: .
    ports:
      - "{{.Port}}:{{.Port}}"
    depends_on:
{{- if eq .Database.Type "postgres"}}
      - db
{{- else if eq .Database.Type "mysql"}}
      - db
{{- else if eq .Database.Type "mongodb"}}
      - db
{{- end}}
    environment:
      - PORT={{.Port}}
{{- if eq .Database.Type "postgres"}}
      - DATABASE_URL=postgres://postgres:postgres@db:5432/{{.Name}}?sslmode=disable
{{- else if eq .Database.Type "mysql"}}
      - DATABASE_URL=root:password@tcp(db:3306)/{{.Name}}?charset=utf8mb4&parseTime=True&loc=Local
{{- else if eq .Database.Type "sqlite"}}
      - DATABASE_URL={{.Name}}.db
{{- else if eq .Database.Type "mongodb"}}
      - DATABASE_URL=mongodb://db:27017/{{.Name}}
{{- else if eq .Database.Type "redis"}}
      - REDIS_ADDR=redis:6379
      - REDIS_PASSWORD=
      - REDIS_DB=0
{{- end}}
    volumes:
      - .:/app
      - /app/tmp
    networks:
      - {{.Name}}-network

{{- if eq .Database.Type "postgres"}}
  db:
    image: postgres:15-alpine
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB={{.Name}}
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - {{.Name}}-network

{{- else if eq .Database.Type "mysql"}}
  db:
    image: mysql:8.0
    environment:
      - MYSQL_ROOT_PASSWORD=password
      - MYSQL_DATABASE={{.Name}}
      - MYSQL_USER=user
      - MYSQL_PASSWORD=password
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
    networks:
      - {{.Name}}-network

{{- else if eq .Database.Type "mongodb"}}
  db:
    image: mongo:7.0
    environment:
      - MONGO_INITDB_DATABASE={{.Name}}
    ports:
      - "27017:27017"
    volumes:
      - mongo_data:/data/db
    networks:
      - {{.Name}}-network

{{- else if eq .Database.Type "redis"}}
  redis:
    image: redis:7.2-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - {{.Name}}-network

{{- end}}

{{- if ne .Database.Type "sqlite"}}
volumes:
{{- if eq .Database.Type "postgres"}}
  postgres_data:
{{- else if eq .Database.Type "mysql"}}
  mysql_data:
{{- else if eq .Database.Type "mongodb"}}
  mongo_data:
{{- else if eq .Database.Type "redis"}}
  redis_data:
{{- end}}

{{- end}}
networks:
  {{.Name}}-network:
    driver: bridge
`
	return g.generateFromTemplate(project, template, filepath.Join(project.Name, "docker-compose.yml"))
}

// generateMakefile generates Makefile for common operations
func (g *APIGenerator) generateMakefile(project *APIProject) error {
	template := `.PHONY: build run test clean docker-up docker-down help

# Build the application
build:
	go build -o bin/{{.Name}} cmd/server/main.go

# Run the application locally
run:
	go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Start with Docker Compose
docker-up:
	docker-compose up --build

# Stop Docker Compose
docker-down:
	docker-compose down

# Stop and remove volumes
docker-clean:
	docker-compose down -v
	docker system prune -f

# Install dependencies
deps:
	go mod tidy
	go mod download

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Live reload for development
dev:
	air

# Database migration (if using migrate tool)
migrate-up:
	migrate -path migrations -database "{{- if eq .Database.Type "postgres"}}postgres://postgres:postgres@localhost:5432/{{.Name}}?sslmode=disable{{- else if eq .Database.Type "mysql"}}mysql://root:password@tcp(localhost:3306)/{{.Name}}{{- else}}sqlite3://{{.Name}}.db{{- end}}" up

migrate-down:
	migrate -path migrations -database "{{- if eq .Database.Type "postgres"}}postgres://postgres:postgres@localhost:5432/{{.Name}}?sslmode=disable{{- else if eq .Database.Type "mysql"}}mysql://root:password@tcp(localhost:3306)/{{.Name}}{{- else}}sqlite3://{{.Name}}.db{{- end}}" down

# Help
help:
	@echo "Available commands:"
	@echo "  build       - Build the application"
	@echo "  run         - Run the application locally"
	@echo "  test        - Run tests"
	@echo "  clean       - Clean build artifacts"
	@echo "  docker-up   - Start with Docker Compose"
	@echo "  docker-down - Stop Docker Compose"
	@echo "  docker-clean- Stop and remove volumes"
	@echo "  deps        - Install dependencies"
	@echo "  fmt         - Format code"
	@echo "  lint        - Lint code"
	@echo "  dev         - Live reload for development"
	@echo "  migrate-up  - Run database migrations up"
	@echo "  migrate-down- Run database migrations down"
`
	return g.generateFromTemplate(project, template, filepath.Join(project.Name, "Makefile"))
}

// generateReadme generates README.md file
func (g *APIGenerator) generateReadme(project *APIProject) error {
	template := `# {{.Name}}

A Go web API built with clean architecture principles.

## Features

- ✅ **Clean Architecture**: Handlers, Services, Repositories, Models
- ✅ **{{.Database.Type | ToCamel}} Database**: Ready to use {{.Database.Type}} database
- ✅ **Docker Support**: Complete Docker setup with docker-compose
- ✅ **GORM Integration**: Database operations with GORM
- ✅ **Gin Framework**: Fast HTTP router and middleware
- ✅ **Environment Configuration**: Flexible configuration management
- ✅ **Health Check**: Built-in health check endpoint

## Quick Start

### Option 1: Docker (Recommended)

1. Start the application with Docker:
   ` + "```bash" + `
   docker-compose up --build
   ` + "```" + `

2. The API will be available at: http://localhost:{{.Port}}

3. Health check: http://localhost:{{.Port}}/health

### Option 2: Local Development

1. Install dependencies:
   ` + "```bash" + `
   go mod tidy
   ` + "```" + `

2. {{- if eq .Database.Type "postgres"}}
   Start PostgreSQL database:
   ` + "```bash" + `
   docker run --name {{.Name}}-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB={{.Name}} -p 5432:5432 -d postgres:15-alpine
   ` + "```" + `
   {{- else if eq .Database.Type "mysql"}}
   Start MySQL database:
   ` + "```bash" + `
   docker run --name {{.Name}}-mysql -e MYSQL_ROOT_PASSWORD=password -e MYSQL_DATABASE={{.Name}} -p 3306:3306 -d mysql:8.0
   ` + "```" + `
   {{- end}}

3. Run the application:
   ` + "```bash" + `
   go run cmd/server/main.go
   ` + "```" + `

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | /health  | Health check |
| GET    | /api/v1/example | Example endpoint |

## Project Structure

` + "```" + `
{{.Name}}/
├── cmd/server/         # Application entry point
├── internal/
│   ├── handlers/       # HTTP handlers
│   ├── services/       # Business logic
│   ├── repositories/   # Data access
│   ├── models/         # Data models
│   └── middleware/     # HTTP middleware
├── pkg/
│   ├── database/       # Database connection
│   ├── config/         # Configuration
│   └── utils/          # Utilities
├── docs/               # Documentation
├── docker-compose.yml  # Docker configuration
├── Dockerfile          # Docker image
├── Makefile           # Build commands
└── .env.example       # Environment variables
` + "```" + `

## Available Commands

Run ` + "`make help`" + ` to see all available commands:

` + "```bash" + `
make build        # Build the application
make run          # Run locally
make test         # Run tests
make docker-up    # Start with Docker
make docker-down  # Stop Docker
make deps         # Install dependencies
make fmt          # Format code
make lint         # Lint code
` + "```" + `

## Database

This project uses **{{.Database.Type | ToCamel}}** as the database.

{{- if eq .Database.Type "postgres"}}
### PostgreSQL Configuration
- Host: db (in Docker) / localhost (local)
- Port: 5432
- Database: {{.Name}}
- User: postgres
- Password: postgres
{{- else if eq .Database.Type "mysql"}}
### MySQL Configuration
- Host: db (in Docker) / localhost (local)
- Port: 3306
- Database: {{.Name}}
- User: root
- Password: password
{{- else if eq .Database.Type "sqlite"}}
### SQLite Configuration
- Database file: {{.Name}}.db
{{- else if eq .Database.Type "mongodb"}}
### MongoDB Configuration
- Host: db (in Docker) / localhost (local)
- Port: 27017
- Database: {{.Name}}
{{- end}}

## Environment Variables

Copy ` + "`.env.example`" + ` to ` + "`.env`" + ` and configure:

` + "```bash" + `
cp .env.example .env
` + "```" + `

## Development

1. **Add a new resource:**
   ` + "```bash" + `
   vibercode generate resource
   ` + "```" + `

2. **Run tests:**
   ` + "```bash" + `
   make test
   ` + "```" + `

3. **Format code:**
   ` + "```bash" + `
   make fmt
   ` + "```" + `

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License.
`
	return g.generateFromTemplate(project, template, filepath.Join(project.Name, "README.md"))
}

// generateSupabaseDatabase generates Supabase-specific database configuration
func (g *APIGenerator) generateSupabaseDatabase(project *APIProject) error {
	// Use template from templates package
	return g.generateFromTemplate(project, templates.SupabaseDatabaseTemplate, filepath.Join(project.Name, "pkg", "database", "database.go"))
}

// generateMongoDBDatabase generates MongoDB-specific database configuration
func (g *APIGenerator) generateMongoDBDatabase(project *APIProject) error {
	template := `package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

// Connect establishes a MongoDB connection
func Connect() (*mongo.Database, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = getDefaultDSN()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	DB = client.Database("{{.Name}}")
	log.Println("MongoDB connected successfully")
	return DB, nil
}

// getDefaultDSN returns the default MongoDB connection string
func getDefaultDSN() string {
	return "mongodb://localhost:27017/{{.Name}}"
}

// Close closes the MongoDB connection
func Close() error {
	if DB != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return DB.Client().Disconnect(ctx)
	}
	return nil
}

// GetCollection returns a MongoDB collection
func GetCollection(name string) *mongo.Collection {
	return DB.Collection(name)
}

// HealthCheck verifies the MongoDB connection
func HealthCheck() error {
	if DB == nil {
		return fmt.Errorf("database not connected")
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	return DB.Client().Ping(ctx, nil)
}
`
	return g.generateFromTemplate(project, template, filepath.Join(project.Name, "pkg", "database", "database.go"))
}

// generateRedisDatabase generates Redis-specific database configuration
func (g *APIGenerator) generateRedisDatabase(project *APIProject) error {
	template := `package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

// Connect establishes a Redis connection
func Connect() (*redis.Client, error) {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	password := os.Getenv("REDIS_PASSWORD")
	db := 0
	if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
		if parsedDB, err := strconv.Atoi(dbStr); err == nil {
			db = parsedDB
		}
	}

	maxRetries := 3
	if retriesStr := os.Getenv("REDIS_MAX_RETRIES"); retriesStr != "" {
		if parsedRetries, err := strconv.Atoi(retriesStr); err == nil {
			maxRetries = parsedRetries
		}
	}

	poolSize := 10
	if poolStr := os.Getenv("REDIS_POOL_SIZE"); poolStr != "" {
		if parsedPool, err := strconv.Atoi(poolStr); err == nil {
			poolSize = parsedPool
		}
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:        addr,
		Password:    password,
		DB:          db,
		MaxRetries:  maxRetries,
		PoolSize:    poolSize,
		DialTimeout: 5 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	log.Println("Redis connected successfully")
	return RedisClient, nil
}

// Close closes the Redis connection
func Close() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}

// GetClient returns the Redis client
func GetClient() *redis.Client {
	return RedisClient
}

// HealthCheck verifies the Redis connection
func HealthCheck() error {
	if RedisClient == nil {
		return fmt.Errorf("redis not connected")
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	_, err := RedisClient.Ping(ctx).Result()
	return err
}

// Set stores a key-value pair with expiration
func Set(key string, value interface{}, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value by key
func Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	return RedisClient.Get(ctx, key).Result()
}

// Delete removes a key
func Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	return RedisClient.Del(ctx, key).Err()
}

// Exists checks if a key exists
func Exists(key string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	count, err := RedisClient.Exists(ctx, key).Result()
	return count > 0, err
}
`
	return g.generateFromTemplate(project, template, filepath.Join(project.Name, "pkg", "database", "database.go"))
}

// generateFromTemplate generates a file from a template string
func (g *APIGenerator) generateFromTemplate(project *APIProject, templateStr, outputPath string) error {
	tmpl, err := template.New("generator").Funcs(template.FuncMap{
		"ToCamel":      func(s string) string { return strings.Title(s) },
		"ToLowerCamel": func(s string) string { return strings.ToLower(s[:1]) + s[1:] },
		"ToSnake":      func(s string) string { return strings.ToLower(s) },
		"ToKebab":      func(s string) string { return strings.ToLower(s) },
		"GetEnvVars":   func(db *models.DatabaseProvider) map[string]string { return db.GetEnvironmentVars() },
	}).Parse(templateStr)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outputPath, err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, project); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// generateManifest creates a .vibercode/manifest.vibe file with project configuration
func (g *APIGenerator) generateManifest(project *APIProject) error {
	// Create .vibercode directory
	vibercodeDir := filepath.Join(project.Name, ".vibercode")
	if err := os.MkdirAll(vibercodeDir, 0755); err != nil {
		return fmt.Errorf("failed to create .vibercode directory: %w", err)
	}

	// Create manifest structure
	now := time.Now().Format(time.RFC3339)
	manifest := VibercodeManifest{
		Version:     "1.0.0",
		ProjectType: "api",
		Name:        project.Name,
		Port:        project.Port,
		Database:    project.Database,
		Module:      project.Module,
		GeneratedAt: now,
		UpdatedAt:   now,
		CLI: VibercodeManifestCLI{
			Version: "1.0.0", // TODO: Get this from build info
			Command: "vibercode generate api",
		},
		History: []VibercodeManifestEvent{
			{
				Type:        "create",
				Description: "Initial API project generation",
				Timestamp:   now,
				CLI: VibercodeManifestCLI{
					Version: "1.0.0",
					Command: "vibercode generate api",
				},
			},
		},
		Resources: []VibercodeResource{}, // Initialize empty resources array
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal manifest to JSON: %w", err)
	}

	// Write to file
	manifestPath := filepath.Join(vibercodeDir, "manifest.vibe")
	if err := os.WriteFile(manifestPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write manifest file: %w", err)
	}

	return nil
}