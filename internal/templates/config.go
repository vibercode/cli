package templates

// ConfigTemplate generates the main configuration file
const ConfigTemplate = `{
  "environment": "{{.Environment}}",
  "server": {
    "host": "{{.Server.Host}}",
    "port": {{.Server.Port}},
    "read_timeout": "{{.Server.ReadTimeout}}",
    "write_timeout": "{{.Server.WriteTimeout}}",
    "idle_timeout": "{{.Server.IdleTimeout}}",
    "max_header_bytes": {{.Server.MaxHeaderBytes}},
    "enable_https": {{.Server.EnableHTTPS}},
    "cors": {
      "allowed_origins": {{.Server.CORS.AllowedOrigins | toJSON}},
      "allowed_methods": {{.Server.CORS.AllowedMethods | toJSON}},
      "allowed_headers": {{.Server.CORS.AllowedHeaders | toJSON}},
      "allow_credentials": {{.Server.CORS.AllowCredentials}},
      "max_age": {{.Server.CORS.MaxAge}}
    }
  },
  "database": {
    "provider": "{{.Database.Provider}}",
    "host": "{{.Database.Host}}",
    "port": {{.Database.Port}},
    "database": "{{.Database.Database}}",
    "schema": "{{.Database.Schema}}",
    "ssl_mode": "{{.Database.SSLMode}}",
    "max_open_conns": {{.Database.MaxOpenConns}},
    "max_idle_conns": {{.Database.MaxIdleConns}},
    "conn_max_lifetime": "{{.Database.ConnMaxLifetime}}",
    "conn_max_idle_time": "{{.Database.ConnMaxIdleTime}}",
    "migrations": {
      "auto_migrate": {{.Database.Migrations.AutoMigrate}},
      "backup_before": {{.Database.Migrations.BackupBefore}},
      "versioned": {{.Database.Migrations.Versioned}},
      "migrations_path": "{{.Database.Migrations.MigrationsPath}}"
    }{{if eq .Database.Provider "supabase"}},
    "supabase": {
      "project_url": "{{.Database.Supabase.ProjectURL}}",
      "enable_auth": {{.Database.Supabase.EnableAuth}},
      "enable_storage": {{.Database.Supabase.EnableStorage}},
      "enable_realtime": {{.Database.Supabase.EnableRealtime}}
    }{{end}}
  },
  "auth": {
    "provider": "{{.Auth.Provider}}",
    "token_expiry": "{{.Auth.TokenExpiry}}",
    "refresh_expiry": "{{.Auth.RefreshExpiry}}",
    "password_min_length": {{.Auth.PasswordMinLen}},
    "enable_mfa": {{.Auth.EnableMFA}}{{if eq .Auth.Provider "oauth2"}},
    "oauth2": {
      "redirect_url": "{{.Auth.OAuth2.RedirectURL}}"
    }{{end}}
  },
  "storage": {
    "provider": "{{.Storage.Provider}}",
    "max_file_size": {{.Storage.MaxFileSize}},
    "allowed_types": {{.Storage.AllowedTypes | toJSON}}{{if eq .Storage.Provider "local"}},
    "local_path": "{{.Storage.LocalPath}}"{{end}}{{if eq .Storage.Provider "s3"}},
    "s3": {
      "bucket": "{{.Storage.S3.Bucket}}",
      "region": "{{.Storage.S3.Region}}"
    }{{end}}{{if eq .Storage.Provider "supabase"}},
    "supabase": {
      "bucket_name": "{{.Storage.Supabase.BucketName}}"
    }{{end}}
  },
  "cache": {
    "provider": "{{.Cache.Provider}}",
    "ttl": "{{.Cache.TTL}}"{{if eq .Cache.Provider "redis"}},
    "redis": {
      "host": "{{.Cache.Redis.Host}}",
      "port": {{.Cache.Redis.Port}},
      "database": {{.Cache.Redis.Database}},
      "pool_size": {{.Cache.Redis.PoolSize}}
    }{{end}}
  },
  "logging": {
    "level": "{{.Logging.Level}}",
    "format": "{{.Logging.Format}}",
    "output": "{{.Logging.Output}}",
    "max_size": {{.Logging.MaxSize}},
    "max_backups": {{.Logging.MaxBackups}},
    "max_age": {{.Logging.MaxAge}},
    "compress": {{.Logging.Compress}}
  },
  "monitoring": {
    "enable_metrics": {{.Monitoring.EnableMetrics}},
    "enable_tracing": {{.Monitoring.EnableTracing}},
    "enable_health": {{.Monitoring.EnableHealth}},
    "metrics_port": {{.Monitoring.MetricsPort}},
    "prometheus": {
      "path": "{{.Monitoring.Prometheus.Path}}"
    }
  },
  "security": {
    "enable_rate_limiting": {{.Security.EnableRateLimiting}},
    "rate_limit": {
      "requests_per_second": {{.Security.RateLimit.RequestsPerSecond}},
      "burst_size": {{.Security.RateLimit.BurstSize}},
      "cleanup_interval": "{{.Security.RateLimit.CleanupInterval}}"
    },
    "enable_csrf": {{.Security.EnableCSRF}},
    "trusted_proxies": {{.Security.TrustedProxies | toJSON}}
  },
  "features": {
    "enable_graphql": {{.Features.EnableGraphQL}},
    "enable_websocket": {{.Features.EnableWebSocket}},
    "enable_file_upload": {{.Features.EnableFileUpload}},
    "enable_notifications": {{.Features.EnableNotifications}},
    "enable_search": {{.Features.EnableSearch}},
    "enable_caching": {{.Features.EnableCaching}}
  }
}`

// DevelopmentConfigTemplate generates development environment configuration
const DevelopmentConfigTemplate = `{
  "environment": "development",
  "server": {
    "host": "localhost",
    "port": 8080,
    "cors": {
      "allowed_origins": ["*"],
      "allowed_methods": ["GET", "POST", "PUT", "DELETE", "OPTIONS"],
      "allowed_headers": ["*"]
    }
  },
  "database": {
    "provider": "{{.DatabaseProvider}}",
    "host": "localhost",
    "port": {{.DatabasePort}},
    "database": "{{.ProjectName}}_dev",
    "ssl_mode": "disable",
    "migrations": {
      "auto_migrate": true,
      "versioned": true
    }
  },
  "logging": {
    "level": "debug",
    "format": "text",
    "output": "stdout"
  },
  "monitoring": {
    "enable_metrics": true,
    "enable_health": true
  },
  "security": {
    "enable_rate_limiting": false,
    "enable_csrf": false
  }
}`

// StagingConfigTemplate generates staging environment configuration
const StagingConfigTemplate = `{
  "environment": "staging",
  "server": {
    "host": "0.0.0.0",
    "port": 8080,
    "cors": {
      "allowed_origins": ["https://staging.{{.Domain}}"],
      "allowed_methods": ["GET", "POST", "PUT", "DELETE", "OPTIONS"],
      "allow_credentials": true
    }
  },
  "database": {
    "provider": "{{.DatabaseProvider}}",
    "database": "{{.ProjectName}}_staging",
    "ssl_mode": "require",
    "migrations": {
      "auto_migrate": false,
      "backup_before": true,
      "versioned": true
    }
  },
  "logging": {
    "level": "info",
    "format": "json",
    "output": "stdout"
  },
  "monitoring": {
    "enable_metrics": true,
    "enable_tracing": true,
    "enable_health": true
  },
  "security": {
    "enable_rate_limiting": true,
    "rate_limit": {
      "requests_per_second": 100,
      "burst_size": 200
    },
    "enable_csrf": true
  }
}`

// ProductionConfigTemplate generates production environment configuration
const ProductionConfigTemplate = `{
  "environment": "production",
  "server": {
    "host": "0.0.0.0",
    "port": 8080,
    "enable_https": true,
    "cors": {
      "allowed_origins": ["https://{{.Domain}}"],
      "allowed_methods": ["GET", "POST", "PUT", "DELETE"],
      "allow_credentials": true,
      "max_age": 86400
    }
  },
  "database": {
    "provider": "{{.DatabaseProvider}}",
    "database": "{{.ProjectName}}_prod",
    "ssl_mode": "require",
    "max_open_conns": 50,
    "max_idle_conns": 25,
    "migrations": {
      "auto_migrate": false,
      "backup_before": true,
      "versioned": true
    }
  },
  "logging": {
    "level": "warn",
    "format": "json",
    "output": "file",
    "filename": "/var/log/{{.ProjectName}}/app.log",
    "max_size": 100,
    "max_backups": 10,
    "max_age": 30,
    "compress": true
  },
  "monitoring": {
    "enable_metrics": true,
    "enable_tracing": true,
    "enable_health": true,
    "metrics_port": 9090
  },
  "security": {
    "enable_rate_limiting": true,
    "rate_limit": {
      "requests_per_second": 50,
      "burst_size": 100,
      "cleanup_interval": "1m"
    },
    "enable_csrf": true,
    "trusted_proxies": ["10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"]
  },
  "features": {
    "enable_file_upload": true,
    "enable_notifications": true,
    "enable_search": true,
    "enable_caching": true
  }
}`

// EnvTemplate generates environment variables file
const EnvTemplate = `# {{.ProjectName}} Configuration
# Environment: {{.Environment}}

# Server Configuration
SERVER_HOST={{.Server.Host}}
SERVER_PORT={{.Server.Port}}
{{if .Server.EnableHTTPS}}
SERVER_TLS_CERT_FILE={{.Server.TLSCertFile}}
SERVER_TLS_KEY_FILE={{.Server.TLSKeyFile}}
{{end}}

# Database Configuration
DATABASE_PROVIDER={{.Database.Provider}}
{{if ne .Database.Provider "sqlite"}}
DATABASE_HOST={{.Database.Host}}
DATABASE_PORT={{.Database.Port}}
DATABASE_USERNAME={{.Database.Username}}
DATABASE_PASSWORD={{.Database.Password}}
{{end}}
DATABASE_DATABASE={{.Database.Database}}
{{if eq .Database.Provider "postgres" "mysql"}}
DATABASE_SSL_MODE={{.Database.SSLMode}}
{{end}}

{{if eq .Database.Provider "supabase"}}
# Supabase Configuration
SUPABASE_PROJECT_URL={{.Database.Supabase.ProjectURL}}
SUPABASE_API_KEY={{.Database.Supabase.APIKey}}
SUPABASE_SERVICE_KEY={{.Database.Supabase.ServiceKey}}
SUPABASE_JWT_SECRET={{.Database.Supabase.JWTSecret}}
{{end}}

# Authentication Configuration
AUTH_PROVIDER={{.Auth.Provider}}
{{if eq .Auth.Provider "jwt"}}
AUTH_JWT_SECRET={{.Auth.JWTSecret}}
{{end}}
{{if eq .Auth.Provider "oauth2"}}
# OAuth2 Configuration
GOOGLE_CLIENT_ID={{.Auth.OAuth2.GoogleClientID}}
GOOGLE_CLIENT_SECRET={{.Auth.OAuth2.GoogleClientSecret}}
GITHUB_CLIENT_ID={{.Auth.OAuth2.GitHubClientID}}
GITHUB_CLIENT_SECRET={{.Auth.OAuth2.GitHubClientSecret}}
OAUTH2_REDIRECT_URL={{.Auth.OAuth2.RedirectURL}}
{{end}}

# Storage Configuration
STORAGE_PROVIDER={{.Storage.Provider}}
{{if eq .Storage.Provider "local"}}
STORAGE_LOCAL_PATH={{.Storage.LocalPath}}
{{end}}
{{if eq .Storage.Provider "s3"}}
S3_BUCKET={{.Storage.S3.Bucket}}
S3_REGION={{.Storage.S3.Region}}
S3_ACCESS_KEY={{.Storage.S3.AccessKey}}
S3_SECRET_KEY={{.Storage.S3.SecretKey}}
{{end}}

# Cache Configuration
CACHE_PROVIDER={{.Cache.Provider}}
{{if eq .Cache.Provider "redis"}}
REDIS_HOST={{.Cache.Redis.Host}}
REDIS_PORT={{.Cache.Redis.Port}}
REDIS_PASSWORD={{.Cache.Redis.Password}}
REDIS_DATABASE={{.Cache.Redis.Database}}
{{end}}

# Security Configuration
{{if .Security.EnableCSRF}}
CSRF_SECRET={{.Security.CSRFSecret}}
{{end}}
ENCRYPTION_KEY={{.Security.SecretKeys.EncryptionKey}}
SIGNING_KEY={{.Security.SecretKeys.SigningKey}}
`

// DockerComposeTemplate generates docker-compose.yml for development
const DockerComposeTemplate = `version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "{{.Server.Port}}:{{.Server.Port}}"
    environment:
      - APP_ENV={{.Environment}}
    env_file:
      - .env.{{.Environment}}
    depends_on:
      - database
      {{if eq .Cache.Provider "redis"}}- redis{{end}}
    networks:
      - {{.ProjectName}}_network
    volumes:
      - ./logs:/app/logs
      {{if eq .Storage.Provider "local"}}- ./uploads:/app/uploads{{end}}

  database:
    {{if eq .Database.Provider "postgres"}}
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: {{.Database.Database}}
      POSTGRES_USER: {{.Database.Username}}
      POSTGRES_PASSWORD: {{.Database.Password}}
    ports:
      - "{{.Database.Port}}:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    {{else if eq .Database.Provider "mysql"}}
    image: mysql:8.0
    environment:
      MYSQL_DATABASE: {{.Database.Database}}
      MYSQL_USER: {{.Database.Username}}
      MYSQL_PASSWORD: {{.Database.Password}}
      MYSQL_ROOT_PASSWORD: rootpassword
    ports:
      - "{{.Database.Port}}:3306"
    volumes:
      - mysql_data:/var/lib/mysql
    {{else if eq .Database.Provider "mongodb"}}
    image: mongo:6.0
    environment:
      MONGO_INITDB_DATABASE: {{.Database.Database}}
      MONGO_INITDB_ROOT_USERNAME: {{.Database.Username}}
      MONGO_INITDB_ROOT_PASSWORD: {{.Database.Password}}
    ports:
      - "{{.Database.Port}}:27017"
    volumes:
      - mongodb_data:/data/db
    {{end}}
    networks:
      - {{.ProjectName}}_network

  {{if eq .Cache.Provider "redis"}}
  redis:
    image: redis:7-alpine
    ports:
      - "{{.Cache.Redis.Port}}:6379"
    volumes:
      - redis_data:/data
    networks:
      - {{.ProjectName}}_network
    command: redis-server --appendonly yes
  {{end}}

  {{if .Monitoring.EnableMetrics}}
  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./monitoring/prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    networks:
      - {{.ProjectName}}_network
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--web.console.libraries=/etc/prometheus/console_libraries'
      - '--web.console.templates=/etc/prometheus/consoles'
  {{end}}

volumes:
  {{if eq .Database.Provider "postgres"}}postgres_data:{{end}}
  {{if eq .Database.Provider "mysql"}}mysql_data:{{end}}
  {{if eq .Database.Provider "mongodb"}}mongodb_data:{{end}}
  {{if eq .Cache.Provider "redis"}}redis_data:{{end}}
  {{if .Monitoring.EnableMetrics}}prometheus_data:{{end}}

networks:
  {{.ProjectName}}_network:
    driver: bridge
`

// DockerfileTemplate generates Dockerfile for the application
const DockerfileTemplate = `# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main cmd/server/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -S appuser -u 1001 -G appgroup

# Copy binary from builder
COPY --from=builder /app/main .

# Create necessary directories
RUN mkdir -p logs uploads migrations
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE {{.Server.Port}}

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:{{.Server.Port}}/health || exit 1

# Run the application
CMD ["./main"]
`

// MakefileTemplate generates Makefile for the project
const MakefileTemplate = `# {{.ProjectName}} Makefile

.PHONY: help build run test clean docker-build docker-run dev deps lint format

# Variables
APP_NAME={{.ProjectName}}
DOCKER_IMAGE={{.ProjectName}}:latest
PORT={{.Server.Port}}

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Install dependencies
	go mod download
	go mod tidy

build: ## Build the application
	go build -ldflags="-w -s" -o bin/$(APP_NAME) cmd/server/main.go

run: ## Run the application
	go run cmd/server/main.go

dev: ## Run in development mode with hot reload
	air -c .air.toml

test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

bench: ## Run benchmarks
	go test -bench=. -benchmem ./...

lint: ## Run linter
	golangci-lint run

format: ## Format code
	gofmt -s -w .
	goimports -w .

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE) .

docker-run: ## Run Docker container
	docker run -p $(PORT):$(PORT) --env-file .env.development $(DOCKER_IMAGE)

docker-compose-up: ## Start all services with docker-compose
	docker-compose up -d

docker-compose-down: ## Stop all services
	docker-compose down

migrate-up: ## Run database migrations
	go run cmd/migrate/main.go up

migrate-down: ## Rollback database migrations
	go run cmd/migrate/main.go down

migrate-create: ## Create new migration (usage: make migrate-create NAME=migration_name)
	go run cmd/migrate/main.go create $(NAME)

gen-docs: ## Generate API documentation
	swag init -g cmd/server/main.go -o docs

security-scan: ## Run security scan
	gosec ./...

install-tools: ## Install development tools
	go install github.com/cosmtrek/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
`