# Configuration Management System

## Overview

The Vibercode CLI configuration system provides comprehensive configuration management with support for multiple environments, validation, hot-reloading, and deployment configurations.

## Features

- ✅ **Hierarchical Configuration**: Structured configuration with nested sections
- ✅ **Multi-Environment Support**: Development, staging, and production environments
- ✅ **Validation System**: Comprehensive validation with custom rules
- ✅ **Multiple Formats**: JSON and YAML configuration files
- ✅ **Environment Variables**: Override configuration with environment variables
- ✅ **Template Generation**: Generate configuration files for new projects
- ✅ **Deployment Configs**: Docker, Docker Compose, and Makefile generation
- 🔄 **Hot-Reloading**: Automatic configuration reloading (when fsnotify is available)

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/vibercode/cli/pkg/config"
)

func main() {
    // Create a configuration manager
    manager := config.NewManager()
    
    // Load configuration
    opts := config.LoadOptions{
        Environment: "development",
        Validate:    true,
    }
    
    if err := manager.Load(opts); err != nil {
        panic(err)
    }
    
    // Get current configuration
    cfg := manager.GetConfig()
    fmt.Printf("Server running on %s:%d\n", cfg.Server.Host, cfg.Server.Port)
}
```

### Loading from Files

```go
opts := config.LoadOptions{
    ConfigPath:      "./config.json",
    Environment:     "production",
    EnvFiles:        []string{".env.production"},
    WatchForChanges: true,
    Validate:        true,
}

err := manager.Load(opts)
```

### Generating Project Configuration

```go
opts := config.GeneratorOptions{
    ProjectName:      "my-api",
    OutputPath:       "./",
    Environment:      "development",
    DatabaseProvider: "postgres",
    GenerateEnvFiles: true,
    GenerateDocker:   true,
    GenerateMakefile: true,
}

generator := config.NewGenerator(opts)
err := generator.GenerateProjectConfig()
```

## Configuration Structure

### Main Configuration Sections

```json
{
  "environment": "development",
  "server": {
    "host": "localhost",
    "port": 8080,
    "enable_https": false,
    "cors": {
      "allowed_origins": ["*"],
      "allowed_methods": ["GET", "POST", "PUT", "DELETE"],
      "allow_credentials": false
    }
  },
  "database": {
    "provider": "postgres",
    "host": "localhost",
    "port": 5432,
    "database": "myapp",
    "migrations": {
      "auto_migrate": true,
      "versioned": true
    }
  },
  "auth": {
    "provider": "jwt",
    "token_expiry": "24h",
    "password_min_length": 8
  },
  "storage": {
    "provider": "local",
    "local_path": "./uploads",
    "max_file_size": 10485760
  },
  "cache": {
    "provider": "memory",
    "ttl": "1h"
  },
  "logging": {
    "level": "info",
    "format": "json",
    "output": "stdout"
  },
  "monitoring": {
    "enable_metrics": true,
    "enable_health": true,
    "metrics_port": 9090
  },
  "security": {
    "enable_rate_limiting": true,
    "enable_csrf": true,
    "rate_limit": {
      "requests_per_second": 100,
      "burst_size": 200
    }
  },
  "features": {
    "enable_graphql": false,
    "enable_websocket": false,
    "enable_file_upload": true
  }
}
```

## Environment Management

### Environment Hierarchy

1. **Default Configuration**: Base configuration with sensible defaults
2. **Environment File**: Environment-specific configuration (e.g., `config.production.json`)
3. **Environment Variables**: Override individual settings
4. **Local Overrides**: Local development overrides (`.env.local`)

### Environment Files

```bash
# Load configuration for different environments
config.json                 # Base configuration
config.development.json     # Development overrides
config.staging.json         # Staging overrides  
config.production.json      # Production overrides

.env                        # Base environment variables
.env.development           # Development environment variables
.env.staging              # Staging environment variables
.env.production           # Production environment variables
.env.local                # Local overrides (git-ignored)
```

### Environment Variables

Configuration can be overridden using environment variables with the pattern `SECTION_FIELD`:

```bash
# Server configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_ENABLE_HTTPS=true

# Database configuration
DATABASE_PROVIDER=postgres
DATABASE_HOST=db.example.com
DATABASE_PORT=5432
DATABASE_USERNAME=myuser
DATABASE_PASSWORD=mypassword

# Authentication
AUTH_PROVIDER=jwt
AUTH_JWT_SECRET=your-secret-key

# Supabase configuration
SUPABASE_PROJECT_URL=https://your-project.supabase.co
SUPABASE_API_KEY=your-api-key
```

## Validation System

### Built-in Validators

- **Required Fields**: Ensure critical configuration is present
- **Port Validation**: Valid port numbers (1-65535)
- **Host Validation**: Valid hostnames and IP addresses
- **URL Validation**: Valid HTTP/HTTPS URLs
- **Provider Validation**: Valid database/auth/storage providers
- **File Path Validation**: Valid file system paths

### Custom Validation Rules

```go
// Example: Custom validation for specific providers
func (m *Manager) validateCustomRules(config *Config) error {
    if config.Database.Provider == "supabase" {
        if config.Database.Supabase.ProjectURL == "" {
            return fmt.Errorf("Supabase project URL is required")
        }
    }
    return nil
}
```

### Validation Usage

```go
validator := config.NewValidator()

// Validate entire configuration
err := validator.ValidateConfig(cfg)

// Validate specific section
err := validator.ValidateSection(&cfg.Database)

// Validate single field
err := validator.ValidateField(8080, "port")
```

## Database Providers

### Supported Providers

| Provider   | Description                    | Configuration Required                |
|------------|--------------------------------|---------------------------------------|
| `postgres` | PostgreSQL database            | host, port, username, password, database |
| `mysql`    | MySQL/MariaDB database         | host, port, username, password, database |
| `sqlite`   | SQLite file database           | database (file path)                  |
| `mongodb`  | MongoDB document database      | host, port, username, password, database |
| `supabase` | Supabase backend-as-a-service  | project_url, api_key, service_key     |
| `redis`    | Redis key-value store          | host, port, password, database        |

### Database-Specific Configuration

#### PostgreSQL
```json
{
  "database": {
    "provider": "postgres",
    "host": "localhost",
    "port": 5432,
    "username": "myuser",
    "password": "mypassword", 
    "database": "myapp",
    "ssl_mode": "require",
    "max_open_conns": 25,
    "max_idle_conns": 10
  }
}
```

#### Supabase
```json
{
  "database": {
    "provider": "supabase",
    "supabase": {
      "project_url": "https://your-project.supabase.co",
      "api_key": "your-api-key",
      "service_key": "your-service-key",
      "jwt_secret": "your-jwt-secret",
      "enable_auth": true,
      "enable_storage": true,
      "enable_realtime": false
    }
  }
}
```

## Configuration Templates

### Generating New Projects

```bash
# Generate configuration for a new project
vibercode generate config \
  --project-name my-api \
  --environment development \
  --database postgres \
  --output ./config \
  --env-files \
  --docker \
  --makefile
```

This generates:
- `config.json` - Main configuration
- `config.{env}.json` - Environment-specific configs
- `.env.{env}` - Environment variable files
- `docker-compose.yml` - Docker composition
- `Dockerfile` - Container image
- `Makefile` - Build and development tasks

### Template Customization

Templates support conditional sections based on providers:

```json
{
  "database": {
    "provider": "{{.DatabaseProvider}}",
    {{if eq .DatabaseProvider "supabase"}}
    "supabase": {
      "project_url": "{{.Database.Supabase.ProjectURL}}",
      "enable_auth": {{.Database.Supabase.EnableAuth}}
    }
    {{end}}
  }
}
```

## Security Configuration

### Authentication Providers

- **JWT**: JSON Web Token authentication
- **OAuth2**: OAuth2 provider integration (Google, GitHub)
- **Supabase**: Supabase Auth integration

### Security Features

```json
{
  "security": {
    "enable_rate_limiting": true,
    "rate_limit": {
      "requests_per_second": 100,
      "burst_size": 200,
      "cleanup_interval": "1m"
    },
    "enable_csrf": true,
    "csrf_secret": "your-csrf-secret",
    "trusted_proxies": ["127.0.0.1"],
    "secret_keys": {
      "encryption_key": "your-encryption-key",
      "signing_key": "your-signing-key"
    }
  }
}
```

## Monitoring and Observability

### Metrics and Health Checks

```json
{
  "monitoring": {
    "enable_metrics": true,
    "enable_tracing": true,
    "enable_health": true,
    "metrics_port": 9090,
    "prometheus": {
      "endpoint": "http://prometheus:9090",
      "path": "/metrics"
    },
    "jaeger": {
      "endpoint": "http://jaeger:14268",
      "service_name": "my-api",
      "sample_rate": 0.1
    }
  }
}
```

### Logging Configuration

```json
{
  "logging": {
    "level": "info",
    "format": "json",
    "output": "stdout",
    "filename": "/var/log/myapp/app.log",
    "max_size": 100,
    "max_backups": 3,
    "max_age": 28,
    "compress": true
  }
}
```

## Deployment Configuration

### Docker Support

Generated `docker-compose.yml` includes:
- Application container
- Database service (PostgreSQL, MySQL, MongoDB)
- Redis cache (if enabled)
- Prometheus monitoring (if enabled)

### Environment-Specific Deployments

#### Development
- Debug logging
- Hot-reloading enabled
- Permissive CORS
- Simplified security

#### Staging
- Production-like configuration
- Full monitoring enabled
- Restricted CORS
- Migration safety checks

#### Production
- Error-level logging
- HTTPS enforced
- Full security features
- Performance optimizations

## Advanced Features

### Configuration Updates

```go
// Update configuration at runtime
err := manager.UpdateConfig(func(config *Config) error {
    config.Server.Port = 9000
    config.Logging.Level = "debug"
    return nil
})
```

### Reload Callbacks

```go
// Register callback for configuration changes
manager.AddReloadCallback(func(oldConfig, newConfig *Config) error {
    if oldConfig.Server.Port != newConfig.Server.Port {
        fmt.Println("Server port changed, restart required")
    }
    return nil
})
```

### Configuration Export

```go
// Save current configuration
err := manager.SaveConfig("./current-config.json")

// Generate environment file
err := manager.GenerateEnvFile("./.env.current")
```

## Best Practices

### Development
1. Use environment-specific configuration files
2. Never commit sensitive data (use `.env.local`)
3. Enable validation in development
4. Use hot-reloading for rapid iteration

### Production
1. Use environment variables for secrets
2. Enable HTTPS and security features
3. Configure proper logging and monitoring
4. Use versioned migrations

### Security
1. Rotate JWT secrets regularly
2. Use strong encryption keys
3. Enable CSRF protection
4. Configure trusted proxies correctly
5. Use rate limiting to prevent abuse

## API Reference

### Types

- `Config` - Main configuration structure
- `Manager` - Configuration management
- `Validator` - Configuration validation
- `Generator` - Configuration generation
- `LoadOptions` - Configuration loading options
- `GeneratorOptions` - Configuration generation options

### Key Methods

```go
// Manager methods
func NewManager() *Manager
func (m *Manager) Load(opts LoadOptions) error
func (m *Manager) GetConfig() *Config
func (m *Manager) UpdateConfig(updater func(*Config) error) error
func (m *Manager) SaveConfig(filePath string) error
func (m *Manager) Reload() error
func (m *Manager) Close() error

// Validator methods
func NewValidator() *Validator
func (v *Validator) ValidateConfig(config *Config) error
func (v *Validator) ValidateSection(section interface{}) error

// Generator methods
func NewGenerator(opts GeneratorOptions) *Generator
func (g *Generator) GenerateProjectConfig() error
```

## Examples

See the `config_test.go` file for comprehensive usage examples including:
- Loading configuration from files
- Environment variable overrides
- Validation scenarios
- Configuration generation
- Serialization and deserialization

## Troubleshooting

### Common Issues

1. **Validation Errors**: Check that all required fields are present and valid
2. **File Not Found**: Ensure configuration files exist and paths are correct
3. **Permission Errors**: Check file permissions for configuration files
4. **Invalid JSON/YAML**: Validate configuration file syntax
5. **Environment Variables**: Ensure environment variables use correct naming pattern

### Debug Mode

Enable debug logging to troubleshoot configuration issues:

```go
opts := config.LoadOptions{
    Environment: "development",
    Validate:    true,
}

// Check validation errors
if err := manager.Load(opts); err != nil {
    fmt.Printf("Configuration error: %v\n", err)
}
```