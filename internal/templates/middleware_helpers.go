package templates

import (
	"fmt"
	"strings"

	"github.com/vibercode/cli/internal/models"
)

// GetMiddlewareRegistryTemplate generates middleware registry template
func GetMiddlewareRegistryTemplate(middlewares []models.MiddlewareConfig) string {
	var imports []string
	var structs []string
	var initFunctions []string
	var registrationCalls []string

	// Common imports
	imports = append(imports,
		"github.com/gin-gonic/gin",
		"github.com/sirupsen/logrus",
	)

	for _, middleware := range middlewares {
		structName := middleware.GetStructName()
		fieldName := getValidFieldName(string(middleware.Type))

		structs = append(structs, fmt.Sprintf("\t%s *%s", fieldName, structName))

		switch middleware.Type {
		case models.AuthMiddleware:
			switch middleware.Options.AuthStrategy {
			case models.JWTAuth:
				initFunctions = append(initFunctions, fmt.Sprintf(`
	// Initialize %s middleware
	m.%s = New%s(
		config.JWT.Secret,
		config.JWT.Issuer,
	)`, middleware.Name, fieldName, structName))
			case models.APIKeyAuth:
				initFunctions = append(initFunctions, fmt.Sprintf(`
	// Initialize %s middleware
	validKeys := []string{config.Auth.APIKeys...}
	m.%s = New%s(validKeys, "%s")`, 
					middleware.Name, fieldName, structName, middleware.Options.APIKeyHeader))
			}
			registrationCalls = append(registrationCalls, fmt.Sprintf("\tr.Use(m.%s.JWT())", fieldName))

		case models.LoggingMiddleware:
			initFunctions = append(initFunctions, fmt.Sprintf(`
	// Initialize %s middleware
	m.%s = New%s(
		logger,
		config.Logging.ExcludePaths,
	)`, middleware.Name, fieldName, structName))
			registrationCalls = append(registrationCalls, fmt.Sprintf("\tr.Use(m.%s.RequestLogger())", fieldName))

		case models.CORSMiddleware:
			initFunctions = append(initFunctions, fmt.Sprintf(`
	// Initialize %s middleware
	m.%s = New%s()`, middleware.Name, fieldName, structName))
			registrationCalls = append(registrationCalls, fmt.Sprintf("\tr.Use(m.%s.CORS())", fieldName))

		case models.RateLimitMiddleware:
			initFunctions = append(initFunctions, fmt.Sprintf(`
	// Initialize %s middleware
	m.%s = New%s(
		config.RateLimit.RequestsPerSecond,
		config.RateLimit.BurstSize,
	)`, middleware.Name, fieldName, structName))
			registrationCalls = append(registrationCalls, fmt.Sprintf("\tr.Use(m.%s.RateLimit())", fieldName))

		case models.CustomMiddleware:
			initFunctions = append(initFunctions, fmt.Sprintf(`
	// Initialize %s middleware
	m.%s = New%s()`, middleware.Name, fieldName, structName))
			registrationCalls = append(registrationCalls, fmt.Sprintf("\tr.Use(m.%s.Handle())", fieldName))
		}
	}

	return fmt.Sprintf(`package middleware

import (
%s
)

// Manager manages all middleware components
type Manager struct {
%s
}

// NewManager creates a new middleware manager
func NewManager(config *MiddlewareConfig, logger *logrus.Logger) *Manager {
	m := &Manager{}
%s
	return m
}

// RegisterMiddleware registers all middleware with the router
func (m *Manager) RegisterMiddleware(r *gin.Engine) {
%s
}

// RegisterAuthMiddleware registers authentication middleware for protected routes
func (m *Manager) RegisterAuthMiddleware(r gin.IRouter) {
	if m.auth != nil {
		r.Use(m.auth.JWT())
	}
}

// RegisterCORSMiddleware registers CORS middleware
func (m *Manager) RegisterCORSMiddleware(r *gin.Engine) {
	if m.cors != nil {
		r.Use(m.cors.CORS())
	}
}
`, generateMiddlewareImports(imports), strings.Join(structs, "\n"), 
		strings.Join(initFunctions, "\n"), strings.Join(registrationCalls, "\n"))
}

// GetMiddlewareConfigTemplate generates middleware configuration template
func GetMiddlewareConfigTemplate(middlewares []models.MiddlewareConfig) string {
	var configStructs []string
	var defaultConfigs []string

	// Base middleware config
	configStructs = append(configStructs, `// MiddlewareConfig holds all middleware configuration
type MiddlewareConfig struct {`)

	for _, middleware := range middlewares {
		switch middleware.Type {
		case models.AuthMiddleware:
			if !containsConfigStruct(configStructs, "Auth") {
				configStructs = append(configStructs, `	Auth   AuthConfig   `+"`yaml:\"auth\" json:\"auth\"`")
			}
		case models.LoggingMiddleware:
			if !containsConfigStruct(configStructs, "Logging") {
				configStructs = append(configStructs, `	Logging LoggingConfig `+"`yaml:\"logging\" json:\"logging\"`")
			}
		case models.CORSMiddleware:
			if !containsConfigStruct(configStructs, "CORS") {
				configStructs = append(configStructs, `	CORS   CORSConfig   `+"`yaml:\"cors\" json:\"cors\"`")
			}
		case models.RateLimitMiddleware:
			if !containsConfigStruct(configStructs, "RateLimit") {
				configStructs = append(configStructs, `	RateLimit RateLimitConfig `+"`yaml:\"rate_limit\" json:\"rate_limit\"`")
			}
		}
	}

	configStructs = append(configStructs, "}")

	// Individual config structs
	for _, middleware := range middlewares {
		switch middleware.Type {
		case models.AuthMiddleware:
			configStructs = append(configStructs, `
// AuthConfig holds authentication middleware configuration
type AuthConfig struct {
	Enabled   bool     `+"`yaml:\"enabled\" json:\"enabled\"`"+`
	Strategy  string   `+"`yaml:\"strategy\" json:\"strategy\"`"+`
	JWTSecret string   `+"`yaml:\"jwt_secret\" json:\"jwt_secret\" env:\"JWT_SECRET\"`"+`
	JWTIssuer string   `+"`yaml:\"jwt_issuer\" json:\"jwt_issuer\"`"+`
	APIKeys   []string `+"`yaml:\"api_keys\" json:\"api_keys\"`"+`
}`)

		case models.LoggingMiddleware:
			configStructs = append(configStructs, `
// LoggingConfig holds logging middleware configuration
type LoggingConfig struct {
	Enabled      bool     `+"`yaml:\"enabled\" json:\"enabled\"`"+`
	Level        string   `+"`yaml:\"level\" json:\"level\"`"+`
	Format       string   `+"`yaml:\"format\" json:\"format\"`"+`
	LogRequests  bool     `+"`yaml:\"log_requests\" json:\"log_requests\"`"+`
	LogResponses bool     `+"`yaml:\"log_responses\" json:\"log_responses\"`"+`
	ExcludePaths []string `+"`yaml:\"exclude_paths\" json:\"exclude_paths\"`"+`
}`)

		case models.CORSMiddleware:
			configStructs = append(configStructs, `
// CORSConfig holds CORS middleware configuration
type CORSConfig struct {
	Enabled          bool     `+"`yaml:\"enabled\" json:\"enabled\"`"+`
	AllowedOrigins   []string `+"`yaml:\"allowed_origins\" json:\"allowed_origins\"`"+`
	AllowedMethods   []string `+"`yaml:\"allowed_methods\" json:\"allowed_methods\"`"+`
	AllowedHeaders   []string `+"`yaml:\"allowed_headers\" json:\"allowed_headers\"`"+`
	ExposeHeaders    []string `+"`yaml:\"expose_headers\" json:\"expose_headers\"`"+`
	AllowCredentials bool     `+"`yaml:\"allow_credentials\" json:\"allow_credentials\"`"+`
	MaxAge           int      `+"`yaml:\"max_age\" json:\"max_age\"`"+`
}`)

		case models.RateLimitMiddleware:
			configStructs = append(configStructs, `
// RateLimitConfig holds rate limiting middleware configuration  
type RateLimitConfig struct {
	Enabled           bool   `+"`yaml:\"enabled\" json:\"enabled\"`"+`
	RequestsPerSecond int    `+"`yaml:\"requests_per_second\" json:\"requests_per_second\"`"+`
	BurstSize         int    `+"`yaml:\"burst_size\" json:\"burst_size\"`"+`
	Strategy          string `+"`yaml:\"strategy\" json:\"strategy\"`"+`
	UseRedis          bool   `+"`yaml:\"use_redis\" json:\"use_redis\"`"+`
	RedisURL          string `+"`yaml:\"redis_url\" json:\"redis_url\" env:\"REDIS_URL\"`"+`
}`)
		}
	}

	// Default configuration function
	defaultConfigs = append(defaultConfigs, `
// DefaultMiddlewareConfig returns default middleware configuration
func DefaultMiddlewareConfig() *MiddlewareConfig {
	return &MiddlewareConfig{`)

	for _, middleware := range middlewares {
		switch middleware.Type {
		case models.AuthMiddleware:
			strategy := string(middleware.Options.AuthStrategy)
			defaultConfigs = append(defaultConfigs, fmt.Sprintf(`		Auth: AuthConfig{
			Enabled:   true,
			Strategy:  "%s",
			JWTSecret: "your-secret-key",
			JWTIssuer: "%s",
			APIKeys:   []string{},
		},`, strategy, middleware.Options.JWTIssuer))

		case models.LoggingMiddleware:
			defaultConfigs = append(defaultConfigs, fmt.Sprintf(`		Logging: LoggingConfig{
			Enabled:      true,
			Level:        "%s",
			Format:       "%s",
			LogRequests:  %t,
			LogResponses: %t,
			ExcludePaths: %s,
		},`, middleware.Options.LogLevel, middleware.Options.LogFormat,
				middleware.Options.LogRequests, middleware.Options.LogResponses,
				formatStringSliceForConfig(middleware.Options.ExcludePaths)))

		case models.CORSMiddleware:
			defaultConfigs = append(defaultConfigs, fmt.Sprintf(`		CORS: CORSConfig{
			Enabled:          true,
			AllowedOrigins:   %s,
			AllowedMethods:   %s,
			AllowedHeaders:   %s,
			ExposeHeaders:    []string{},
			AllowCredentials: %t,
			MaxAge:           %d,
		},`, formatStringSliceForConfig(middleware.Options.AllowedOrigins),
				formatStringSliceForConfig(middleware.Options.AllowedMethods),
				formatStringSliceForConfig(middleware.Options.AllowedHeaders),
				middleware.Options.AllowCredentials, middleware.Options.MaxAge))

		case models.RateLimitMiddleware:
			defaultConfigs = append(defaultConfigs, fmt.Sprintf(`		RateLimit: RateLimitConfig{
			Enabled:           true,
			RequestsPerSecond: %d,
			BurstSize:         %d,
			Strategy:          "%s",
			UseRedis:          %t,
			RedisURL:          "redis://localhost:6379",
		},`, middleware.Options.RequestsPerSecond, middleware.Options.BurstSize,
				middleware.Options.Strategy, middleware.Options.UseRedis))
		}
	}

	defaultConfigs = append(defaultConfigs, `	}
}`)

	return fmt.Sprintf(`package config

%s

%s
`, strings.Join(configStructs, "\n"), strings.Join(defaultConfigs, "\n"))
}

// Helper functions

func generateMiddlewareImports(imports []string) string {
	var importLines []string
	for _, imp := range imports {
		importLines = append(importLines, fmt.Sprintf("\t\"%s\"", imp))
	}
	return strings.Join(importLines, "\n")
}

func containsConfigStruct(structs []string, configName string) bool {
	for _, s := range structs {
		if strings.Contains(s, configName) {
			return true
		}
	}
	return false
}

func formatStringSliceForConfig(slice []string) string {
	if len(slice) == 0 {
		return "[]string{}"
	}

	var quoted []string
	for _, s := range slice {
		quoted = append(quoted, fmt.Sprintf("\"%s\"", s))
	}

	return fmt.Sprintf("[]string{%s}", strings.Join(quoted, ", "))
}

// getValidFieldName converts middleware type to valid Go field name
func getValidFieldName(middlewareType string) string {
	// Convert kebab-case to camelCase for valid Go identifiers
	switch middlewareType {
	case "rate-limit":
		return "rateLimit"
	case "auth":
		return "auth"
	case "logging":
		return "logging"
	case "cors":
		return "cors"
	default:
		// Remove hyphens and convert to lowercase
		return strings.ReplaceAll(strings.ToLower(middlewareType), "-", "")
	}
}