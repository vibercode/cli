package models

import (
	"fmt"
	"strings"
)

// MiddlewareType represents different types of middleware
type MiddlewareType string

const (
	AuthMiddleware      MiddlewareType = "auth"
	LoggingMiddleware   MiddlewareType = "logging"
	CORSMiddleware      MiddlewareType = "cors"
	RateLimitMiddleware MiddlewareType = "rate-limit"
	CustomMiddleware    MiddlewareType = "custom"
)

// AuthStrategy represents different authentication strategies
type AuthStrategy string

const (
	JWTAuth    AuthStrategy = "jwt"
	APIKeyAuth AuthStrategy = "apikey"
	SessionAuth AuthStrategy = "session"
	BasicAuth  AuthStrategy = "basic"
)

// MiddlewareConfig represents middleware configuration
type MiddlewareConfig struct {
	Name        string           `json:"name"`
	Type        MiddlewareType   `json:"type"`
	Description string           `json:"description"`
	Options     MiddlewareOptions `json:"options"`
}

// MiddlewareOptions contains configuration options for different middleware types
type MiddlewareOptions struct {
	// Auth options
	AuthStrategy   AuthStrategy `json:"auth_strategy,omitempty"`
	JWTSecret      string       `json:"jwt_secret,omitempty"`
	JWTIssuer      string       `json:"jwt_issuer,omitempty"`
	APIKeyHeader   string       `json:"api_key_header,omitempty"`
	RequireRole    bool         `json:"require_role,omitempty"`
	
	// Logging options
	LogLevel       string   `json:"log_level,omitempty"`
	LogFormat      string   `json:"log_format,omitempty"`
	LogRequests    bool     `json:"log_requests,omitempty"`
	LogResponses   bool     `json:"log_responses,omitempty"`
	ExcludePaths   []string `json:"exclude_paths,omitempty"`
	
	// CORS options
	AllowedOrigins []string `json:"allowed_origins,omitempty"`
	AllowedMethods []string `json:"allowed_methods,omitempty"`
	AllowedHeaders []string `json:"allowed_headers,omitempty"`
	ExposeHeaders  []string `json:"expose_headers,omitempty"`
	AllowCredentials bool   `json:"allow_credentials,omitempty"`
	MaxAge         int      `json:"max_age,omitempty"`
	
	// Rate limiting options
	RequestsPerSecond int    `json:"requests_per_second,omitempty"`
	BurstSize        int    `json:"burst_size,omitempty"`
	WindowSize       string `json:"window_size,omitempty"`
	Strategy         string `json:"strategy,omitempty"`
	UseRedis         bool   `json:"use_redis,omitempty"`
	
	// Custom middleware options
	CustomLogic    string            `json:"custom_logic,omitempty"`
	Dependencies   []string          `json:"dependencies,omitempty"`
	ConfigFields   map[string]string `json:"config_fields,omitempty"`
}

// MiddlewarePreset represents predefined middleware configurations
type MiddlewarePreset string

const (
	APISecurityPreset    MiddlewarePreset = "api-security"
	WebAppPreset        MiddlewarePreset = "web-app" 
	MicroservicePreset  MiddlewarePreset = "microservice"
	PublicAPIPreset     MiddlewarePreset = "public-api"
)

// GetFileName returns the middleware file name
func (m *MiddlewareConfig) GetFileName() string {
	if m.Type == CustomMiddleware {
		return fmt.Sprintf("%s.go", strings.ToLower(m.Name))
	}
	return fmt.Sprintf("%s.go", string(m.Type))
}

// GetStructName returns the middleware struct name
func (m *MiddlewareConfig) GetStructName() string {
	if m.Type == CustomMiddleware {
		return fmt.Sprintf("%sMiddleware", m.Name)
	}
	
	switch m.Type {
	case AuthMiddleware:
		return "AuthMiddleware"
	case LoggingMiddleware:
		return "LoggingMiddleware"
	case CORSMiddleware:
		return "CORSMiddleware"
	case RateLimitMiddleware:
		return "RateLimitMiddleware"
	default:
		return "Middleware"
	}
}

// GetPackageName returns the package name for the middleware
func (m *MiddlewareConfig) GetPackageName() string {
	return "middleware"
}

// GetImports returns required imports for the middleware
func (m *MiddlewareConfig) GetImports() []string {
	var imports []string
	
	// Common imports
	imports = append(imports, 
		"context",
		"net/http",
		"time",
		"github.com/gin-gonic/gin",
	)
	
	switch m.Type {
	case AuthMiddleware:
		switch m.Options.AuthStrategy {
		case JWTAuth:
			imports = append(imports, 
				"strings",
				"github.com/golang-jwt/jwt/v5",
				"errors",
			)
		case APIKeyAuth:
			imports = append(imports, "crypto/subtle")
		}
		
	case LoggingMiddleware:
		imports = append(imports, 
			"github.com/sirupsen/logrus",
			"bytes",
			"io",
		)
		
	case CORSMiddleware:
		imports = append(imports, "strings")
		
	case RateLimitMiddleware:
		imports = append(imports, 
			"sync",
			"golang.org/x/time/rate",
		)
		if m.Options.UseRedis {
			imports = append(imports, 
				"github.com/go-redis/redis/v8",
				"strconv",
			)
		}
	}
	
	return imports
}

// GetDefaultPresetMiddlewares returns middleware configs for a preset
func GetDefaultPresetMiddlewares(preset MiddlewarePreset) []MiddlewareConfig {
	switch preset {
	case APISecurityPreset:
		return []MiddlewareConfig{
			{
				Name:        "Auth",
				Type:        AuthMiddleware,
				Description: "JWT authentication middleware",
				Options: MiddlewareOptions{
					AuthStrategy: JWTAuth,
					JWTSecret:    "${JWT_SECRET}",
					JWTIssuer:    "vibercode-api",
					RequireRole:  false,
				},
			},
			{
				Name:        "CORS",
				Type:        CORSMiddleware,
				Description: "CORS configuration middleware",
				Options: MiddlewareOptions{
					AllowedOrigins: []string{"*"},
					AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
					AllowedHeaders: []string{"Authorization", "Content-Type"},
					MaxAge:         3600,
				},
			},
			{
				Name:        "RateLimit",
				Type:        RateLimitMiddleware,
				Description: "Rate limiting middleware",
				Options: MiddlewareOptions{
					RequestsPerSecond: 100,
					BurstSize:        10,
					Strategy:         "token_bucket",
				},
			},
			{
				Name:        "Logging",
				Type:        LoggingMiddleware,
				Description: "Request/response logging middleware",
				Options: MiddlewareOptions{
					LogLevel:     "info",
					LogFormat:    "json",
					LogRequests:  true,
					LogResponses: false,
					ExcludePaths: []string{"/health", "/metrics"},
				},
			},
		}
		
	case WebAppPreset:
		return []MiddlewareConfig{
			{
				Name:        "CORS",
				Type:        CORSMiddleware,
				Description: "CORS for web application",
				Options: MiddlewareOptions{
					AllowedOrigins:   []string{"http://localhost:3000", "https://app.example.com"},
					AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
					AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Requested-With"},
					AllowCredentials: true,
					MaxAge:          86400,
				},
			},
			{
				Name:        "Auth",
				Type:        AuthMiddleware,
				Description: "Session-based authentication",
				Options: MiddlewareOptions{
					AuthStrategy: SessionAuth,
				},
			},
			{
				Name:        "Logging",
				Type:        LoggingMiddleware,
				Description: "Web application logging",
				Options: MiddlewareOptions{
					LogLevel:     "info",
					LogFormat:    "text",
					LogRequests:  true,
					LogResponses: false,
					ExcludePaths: []string{"/static", "/assets"},
				},
			},
		}
		
	case MicroservicePreset:
		return []MiddlewareConfig{
			{
				Name:        "Auth",
				Type:        AuthMiddleware,
				Description: "Service-to-service authentication",
				Options: MiddlewareOptions{
					AuthStrategy: APIKeyAuth,
					APIKeyHeader: "X-API-Key",
				},
			},
			{
				Name:        "RateLimit",
				Type:        RateLimitMiddleware,
				Description: "Service rate limiting",
				Options: MiddlewareOptions{
					RequestsPerSecond: 1000,
					BurstSize:        100,
					Strategy:         "sliding_window",
					UseRedis:         true,
				},
			},
			{
				Name:        "Logging",
				Type:        LoggingMiddleware,
				Description: "Microservice logging",
				Options: MiddlewareOptions{
					LogLevel:     "debug",
					LogFormat:    "json",
					LogRequests:  true,
					LogResponses: true,
				},
			},
		}
		
	case PublicAPIPreset:
		return []MiddlewareConfig{
			{
				Name:        "CORS",
				Type:        CORSMiddleware,
				Description: "Public API CORS",
				Options: MiddlewareOptions{
					AllowedOrigins: []string{"*"},
					AllowedMethods: []string{"GET", "POST", "OPTIONS"},
					AllowedHeaders: []string{"Authorization", "Content-Type"},
					MaxAge:         3600,
				},
			},
			{
				Name:        "RateLimit",
				Type:        RateLimitMiddleware,
				Description: "Public API rate limiting",
				Options: MiddlewareOptions{
					RequestsPerSecond: 10,
					BurstSize:        5,
					Strategy:         "token_bucket",
				},
			},
			{
				Name:        "Logging",
				Type:        LoggingMiddleware,
				Description: "Public API logging",
				Options: MiddlewareOptions{
					LogLevel:     "info",
					LogFormat:    "json",
					LogRequests:  true,
					LogResponses: false,
				},
			},
		}
		
	default:
		return []MiddlewareConfig{}
	}
}

// IsValidMiddlewareType checks if a middleware type is valid
func IsValidMiddlewareType(middlewareType string) bool {
	validTypes := []string{
		string(AuthMiddleware),
		string(LoggingMiddleware),
		string(CORSMiddleware),
		string(RateLimitMiddleware),
		string(CustomMiddleware),
	}
	
	for _, validType := range validTypes {
		if middlewareType == validType {
			return true
		}
	}
	return false
}

// IsValidPreset checks if a preset is valid
func IsValidPreset(preset string) bool {
	validPresets := []string{
		string(APISecurityPreset),
		string(WebAppPreset),
		string(MicroservicePreset),
		string(PublicAPIPreset),
	}
	
	for _, validPreset := range validPresets {
		if preset == validPreset {
			return true
		}
	}
	return false
}