# Task 07: Middleware Generator

## Overview
Implement a middleware generator that creates custom middleware components for Go web APIs. This includes authentication, logging, CORS, rate limiting, and custom business logic middleware.

## Objectives
- Generate standard middleware (auth, logging, CORS, rate limiting)
- Create custom middleware from templates
- Integrate middleware with existing project structure
- Provide configuration options for middleware behavior

## Implementation Details

### Command Structure
```bash
# Generate standard middleware
vibercode generate middleware --type auth
vibercode generate middleware --type logging
vibercode generate middleware --type cors
vibercode generate middleware --type rate-limit

# Generate custom middleware
vibercode generate middleware --name CustomValidator --custom

# Generate multiple middleware at once
vibercode generate middleware --preset api-security
```

### Middleware Types

#### 1. Authentication Middleware
- JWT validation
- API key validation
- Session-based authentication
- Role-based access control (RBAC)

#### 2. Logging Middleware
- Request/response logging
- Performance metrics
- Error tracking
- Structured logging with levels

#### 3. CORS Middleware
- Cross-origin resource sharing
- Configurable origins, methods, headers
- Preflight request handling

#### 4. Rate Limiting Middleware
- Request rate limiting
- IP-based or user-based limits
- Different strategies (sliding window, token bucket)
- Redis integration for distributed systems

#### 5. Custom Middleware
- Template-based generation
- Business logic integration
- Request/response transformation
- Custom validation logic

### File Structure
```
internal/
├── middleware/
│   ├── auth.go                 # Authentication middleware
│   ├── logging.go             # Logging middleware
│   ├── cors.go                # CORS middleware
│   ├── rate_limit.go          # Rate limiting middleware
│   ├── custom_validator.go    # Custom middleware example
│   └── middleware.go          # Middleware registry
├── config/
│   └── middleware.go          # Middleware configuration
```

### Templates Required

#### Authentication Middleware Template
```go
package middleware

import (
    "context"
    "net/http"
    "strings"
    
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
    secretKey []byte
    issuer    string
}

func NewAuthMiddleware(secretKey string, issuer string) *AuthMiddleware {
    return &AuthMiddleware{
        secretKey: []byte(secretKey),
        issuer:    issuer,
    }
}

func (m *AuthMiddleware) JWT() gin.HandlerFunc {
    return func(c *gin.Context) {
        // JWT validation logic
    }
}

func (m *AuthMiddleware) APIKey() gin.HandlerFunc {
    return func(c *gin.Context) {
        // API key validation logic
    }
}
```

#### Logging Middleware Template
```go
package middleware

import (
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
)

type LoggingMiddleware struct {
    logger *logrus.Logger
}

func NewLoggingMiddleware(logger *logrus.Logger) *LoggingMiddleware {
    return &LoggingMiddleware{logger: logger}
}

func (m *LoggingMiddleware) RequestLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Request logging logic
    }
}
```

### Configuration Integration

#### Middleware Configuration Structure
```go
type MiddlewareConfig struct {
    Auth struct {
        Enabled   bool   `yaml:"enabled"`
        JWTSecret string `yaml:"jwt_secret"`
        Issuer    string `yaml:"issuer"`
    } `yaml:"auth"`
    
    Logging struct {
        Enabled bool   `yaml:"enabled"`
        Level   string `yaml:"level"`
        Format  string `yaml:"format"`
    } `yaml:"logging"`
    
    CORS struct {
        Enabled        bool     `yaml:"enabled"`
        AllowedOrigins []string `yaml:"allowed_origins"`
        AllowedMethods []string `yaml:"allowed_methods"`
        AllowedHeaders []string `yaml:"allowed_headers"`
    } `yaml:"cors"`
    
    RateLimit struct {
        Enabled bool `yaml:"enabled"`
        RPS     int  `yaml:"rps"`
        Burst   int  `yaml:"burst"`
    } `yaml:"rate_limit"`
}
```

### Generator Implementation

#### Command Integration
```go
// Add to cmd/generate.go
var generateMiddlewareCmd = &cobra.Command{
    Use:   "middleware",
    Short: "🔧 Generate middleware components",
    Long: ui.Bold.Sprint("Generate middleware components") + "\n\n" +
        "This command creates middleware with:\n" +
        "  " + ui.IconGear + " Authentication (JWT, API Key)\n" +
        "  " + ui.IconDoc + " Logging and monitoring\n" +
        "  " + ui.IconCORS + " CORS configuration\n" +
        "  " + ui.IconSpeed + " Rate limiting\n" +
        "  " + ui.IconCode + " Custom middleware templates\n",
    RunE: func(cmd *cobra.Command, args []string) error {
        middlewareType, _ := cmd.Flags().GetString("type")
        customName, _ := cmd.Flags().GetString("name")
        isCustom, _ := cmd.Flags().GetBool("custom")
        preset, _ := cmd.Flags().GetString("preset")

        gen := generator.NewMiddlewareGenerator()
        return gen.Generate(generator.MiddlewareOptions{
            Type:    middlewareType,
            Name:    customName,
            Custom:  isCustom,
            Preset:  preset,
        })
    },
}
```

### Middleware Integration Points

#### 1. Router Integration
- Automatic middleware registration
- Order of middleware execution
- Conditional middleware application

#### 2. Configuration Management
- Environment-based configuration
- Runtime configuration updates
- Validation of middleware settings

#### 3. Dependency Injection
- Service integration (database, cache, external APIs)
- Logger integration
- Metrics collection integration

### Testing Strategy

#### Unit Tests
- Middleware logic testing
- Configuration validation
- Error handling scenarios

#### Integration Tests
- Middleware chain testing
- Performance impact testing
- Security validation

### Documentation Generation

#### Middleware Documentation
- Purpose and functionality
- Configuration options
- Usage examples
- Performance considerations

## Dependencies
- Task 02: Template System Enhancement (for middleware templates)
- Task 03: Configuration Management (for middleware config)

## Deliverables
1. Middleware generator implementation
2. Standard middleware templates (auth, logging, CORS, rate limiting)
3. Custom middleware generation capability
4. Configuration integration
5. Documentation and examples
6. Unit and integration tests

## Acceptance Criteria
- [x] Generate standard middleware types
- [x] Create custom middleware from templates
- [x] Integrate with existing project configuration
- [x] Support middleware presets for common use cases
- [x] Include comprehensive documentation
- [x] Pass all unit and integration tests
- [x] Support multiple authentication strategies
- [x] Configurable middleware behavior

## Implementation Priority
High - Essential for API security and monitoring capabilities

## Estimated Effort
3-4 days

## Status
✅ **COMPLETED** - December 2024

## Implementation Summary

### ✅ Completed Features
1. **CLI Command Integration**
   - Added `vibercode generate middleware` command with full flag support
   - Interactive and non-interactive modes supported
   - Help documentation and examples

2. **Middleware Generator**
   - Complete generator logic in `internal/generator/middleware.go`
   - Support for single middleware, custom middleware, and presets
   - Error handling and validation

3. **Data Models**
   - Comprehensive models in `internal/models/middleware.go`
   - Support for all middleware types and configurations
   - Preset definitions with validation

4. **Template System**
   - Complete templates for all middleware types in `internal/templates/middleware.go`
   - Helper templates in `internal/templates/middleware_helpers.go`
   - Registry and configuration templates

5. **Middleware Types Implemented**
   - **Authentication**: JWT, API Key, Session, Basic Auth
   - **Logging**: Request/response logging with structured logs
   - **CORS**: Full CORS configuration with origin validation
   - **Rate Limiting**: Token bucket and sliding window with Redis support
   - **Custom**: Template-based custom middleware generation

6. **Presets Available**
   - `api-security`: Auth + CORS + RateLimit + Logging
   - `web-app`: CORS + Auth + Logging
   - `microservice`: Auth + RateLimit + Logging
   - `public-api`: CORS + RateLimit + Logging

7. **Configuration System**
   - YAML/JSON configuration support
   - Environment variable integration
   - Default configurations for each middleware type

8. **Registry System**
   - Middleware manager for centralized registration
   - Automatic middleware chaining
   - Conditional middleware registration

9. **Testing**
   - Comprehensive test suite in `internal/generator/middleware_test.go`
   - Unit tests for generator, models, and validation functions
   - Integration tests for file generation

### 🔧 Technical Implementation Details
- **Generated Files**: Creates clean Go middleware files with proper imports
- **Architecture**: Follows clean architecture principles
- **Security**: Implements security best practices (constant-time comparison, JWT validation)
- **Performance**: Efficient middleware with minimal overhead
- **Extensibility**: Easy to add new middleware types and strategies

### 🚀 Usage Examples
```bash
# Single middleware
vibercode generate middleware --type auth
vibercode generate middleware --type logging
vibercode generate middleware --type cors
vibercode generate middleware --type rate-limit

# Custom middleware
vibercode generate middleware --name CustomValidator --custom

# Complete presets
vibercode generate middleware --preset api-security
vibercode generate middleware --preset web-app
```

### 📁 Generated Structure
```
internal/
├── middleware/
│   ├── auth.go                 # Authentication middleware
│   ├── logging.go             # Logging middleware
│   ├── cors.go                # CORS middleware
│   ├── rate_limit.go          # Rate limiting middleware
│   └── middleware.go          # Middleware registry
├── config/
│   └── middleware.go          # Middleware configuration
```

## Notes
- ✅ Middleware performance impact minimized through efficient implementations
- ✅ Security best practices implemented (JWT validation, constant-time comparison)
- ✅ Distributed systems support via Redis integration for rate limiting
- ✅ Extensible architecture allows easy addition of new middleware types
- ✅ Fixed field naming issues in generated registry templates
- ✅ Comprehensive error handling and validation throughout