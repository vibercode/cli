package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vibercode/cli/internal/models"
)

func TestMiddlewareGenerator_Generate(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "middleware-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Change to temp directory
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	err = os.Chdir(tempDir)
	require.NoError(t, err)
	defer os.Chdir(originalDir)

	tests := []struct {
		name    string
		options MiddlewareOptions
		wantErr bool
	}{
		{
			name: "Generate JWT auth middleware",
			options: MiddlewareOptions{
				Type: "auth",
			},
			wantErr: false,
		},
		{
			name: "Generate logging middleware",
			options: MiddlewareOptions{
				Type: "logging",
			},
			wantErr: false,
		},
		{
			name: "Generate CORS middleware",
			options: MiddlewareOptions{
				Type: "cors",
			},
			wantErr: false,
		},
		{
			name: "Generate rate limit middleware",
			options: MiddlewareOptions{
				Type: "rate-limit",
			},
			wantErr: false,
		},
		{
			name: "Generate custom middleware",
			options: MiddlewareOptions{
				Name:   "CustomValidator",
				Custom: true,
			},
			wantErr: false,
		},
		{
			name: "Generate API security preset",
			options: MiddlewareOptions{
				Preset: "api-security",
			},
			wantErr: false,
		},
		{
			name: "Invalid middleware type",
			options: MiddlewareOptions{
				Type: "invalid",
			},
			wantErr: true,
		},
		{
			name: "Invalid preset",
			options: MiddlewareOptions{
				Preset: "invalid-preset",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := NewMiddlewareGenerator()
			err := gen.Generate(tt.options)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Verify files were created
			if tt.options.Preset != "" {
				// For presets, check that middleware directory exists
				assert.DirExists(t, "internal/middleware")
				assert.DirExists(t, "internal/config")
			} else if tt.options.Custom {
				// For custom middleware, check specific file
				expectedFile := filepath.Join("internal/middleware", "customvalidator.go")
				assert.FileExists(t, expectedFile)
			} else if tt.options.Type != "" {
				// For single middleware, check specific file
				expectedFile := filepath.Join("internal/middleware", tt.options.Type+".go")
				assert.FileExists(t, expectedFile)
			}

			// Clean up for next test
			os.RemoveAll("internal")
		})
	}
}

func TestMiddlewareGenerator_generatePreset(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "middleware-preset-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	originalDir, err := os.Getwd()
	require.NoError(t, err)
	err = os.Chdir(tempDir)
	require.NoError(t, err)
	defer os.Chdir(originalDir)

	gen := NewMiddlewareGenerator()
	gen.options = MiddlewareOptions{
		Preset: "api-security",
	}

	err = gen.generatePreset()
	require.NoError(t, err)

	// Verify all preset files were created
	expectedFiles := []string{
		"internal/middleware/auth.go",
		"internal/middleware/cors.go", 
		"internal/middleware/rate-limit.go",
		"internal/middleware/logging.go",
		"internal/middleware/middleware.go",
		"internal/config/middleware.go",
	}

	for _, file := range expectedFiles {
		assert.FileExists(t, file, "Expected file %s to be created", file)
	}

	// Verify file contents contain expected structures
	authContent, err := os.ReadFile("internal/middleware/auth.go")
	require.NoError(t, err)
	assert.Contains(t, string(authContent), "type AuthMiddleware struct")
	assert.Contains(t, string(authContent), "func NewAuthMiddleware")
	assert.Contains(t, string(authContent), "func (m *AuthMiddleware) JWT()")
}

func TestMiddlewareGenerator_getStandardMiddlewareConfig(t *testing.T) {
	gen := NewMiddlewareGenerator()

	tests := []struct {
		name           string
		middlewareType string
		wantType       models.MiddlewareType
		wantName       string
	}{
		{
			name:           "Auth middleware config",
			middlewareType: "auth",
			wantType:       models.AuthMiddleware,
			wantName:       "Auth",
		},
		{
			name:           "Logging middleware config",
			middlewareType: "logging",
			wantType:       models.LoggingMiddleware,
			wantName:       "Logging",
		},
		{
			name:           "CORS middleware config",
			middlewareType: "cors",
			wantType:       models.CORSMiddleware,
			wantName:       "Cors",
		},
		{
			name:           "Rate limit middleware config",
			middlewareType: "rate-limit",
			wantType:       models.RateLimitMiddleware,
			wantName:       "Rate-limit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen.options.Type = tt.middlewareType
			
			config, err := gen.getStandardMiddlewareConfig()
			require.NoError(t, err)

			assert.Equal(t, tt.wantType, config.Type)
			assert.Equal(t, tt.wantName, config.Name)
			assert.NotEmpty(t, config.Description)
		})
	}
}

func TestMiddlewareGenerator_createMiddlewareDirectory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "middleware-dir-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	originalDir, err := os.Getwd()
	require.NoError(t, err)
	err = os.Chdir(tempDir)
	require.NoError(t, err)
	defer os.Chdir(originalDir)

	gen := NewMiddlewareGenerator()
	err = gen.createMiddlewareDirectory()
	require.NoError(t, err)

	// Verify directories were created
	assert.DirExists(t, "internal/middleware")
	assert.DirExists(t, "internal/config")
}

func TestMiddlewareModels_IsValidMiddlewareType(t *testing.T) {
	tests := []struct {
		name           string
		middlewareType string
		want           bool
	}{
		{"Valid auth type", "auth", true},
		{"Valid logging type", "logging", true},
		{"Valid cors type", "cors", true},
		{"Valid rate-limit type", "rate-limit", true},
		{"Valid custom type", "custom", true},
		{"Invalid type", "invalid", false},
		{"Empty type", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := models.IsValidMiddlewareType(tt.middlewareType)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestMiddlewareModels_IsValidPreset(t *testing.T) {
	tests := []struct {
		name   string
		preset string
		want   bool
	}{
		{"Valid api-security preset", "api-security", true},
		{"Valid web-app preset", "web-app", true},
		{"Valid microservice preset", "microservice", true},
		{"Valid public-api preset", "public-api", true},
		{"Invalid preset", "invalid", false},
		{"Empty preset", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := models.IsValidPreset(tt.preset)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestMiddlewareModels_GetDefaultPresetMiddlewares(t *testing.T) {
	tests := []struct {
		name        string
		preset      models.MiddlewarePreset
		wantCount   int
		wantTypes   []models.MiddlewareType
	}{
		{
			name:      "API Security preset",
			preset:    models.APISecurityPreset,
			wantCount: 4,
			wantTypes: []models.MiddlewareType{
				models.AuthMiddleware,
				models.CORSMiddleware,
				models.RateLimitMiddleware,
				models.LoggingMiddleware,
			},
		},
		{
			name:      "Web App preset",
			preset:    models.WebAppPreset,
			wantCount: 3,
			wantTypes: []models.MiddlewareType{
				models.CORSMiddleware,
				models.AuthMiddleware,
				models.LoggingMiddleware,
			},
		},
		{
			name:      "Microservice preset",
			preset:    models.MicroservicePreset,
			wantCount: 3,
			wantTypes: []models.MiddlewareType{
				models.AuthMiddleware,
				models.RateLimitMiddleware,
				models.LoggingMiddleware,
			},
		},
		{
			name:      "Public API preset",
			preset:    models.PublicAPIPreset,
			wantCount: 3,
			wantTypes: []models.MiddlewareType{
				models.CORSMiddleware,
				models.RateLimitMiddleware,
				models.LoggingMiddleware,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middlewares := models.GetDefaultPresetMiddlewares(tt.preset)
			assert.Len(t, middlewares, tt.wantCount)

			// Check that all expected types are present
			foundTypes := make(map[models.MiddlewareType]bool)
			for _, middleware := range middlewares {
				foundTypes[middleware.Type] = true
			}

			for _, expectedType := range tt.wantTypes {
				assert.True(t, foundTypes[expectedType], 
					"Expected middleware type %s not found in preset %s", expectedType, tt.preset)
			}
		})
	}
}