package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vibercode/cli/internal/models"
)

// TestAuthGenerator tests the authentication generator
func TestAuthGenerator(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()

	// Create test options
	opts := models.AuthGeneratorOptions{
		ProjectName:      "test-auth",
		OutputPath:       tempDir,
		AuthConfig:       models.DefaultAuthConfig(),
		UserModel:        models.DefaultUserModel(),
		DatabaseProvider: "postgres",
		Endpoints:        models.GetDefaultAuthEndpoints(),
	}

	// Create generator
	generator := NewAuthGenerator(opts)

	t.Run("GenerateMiddleware", func(t *testing.T) {
		err := generator.generateMiddleware()
		if err != nil {
			t.Errorf("Failed to generate middleware: %v", err)
		}

		// Check if JWT middleware file was created
		middlewarePath := filepath.Join(tempDir, "internal", "middleware", "jwt.go")
		if _, err := os.Stat(middlewarePath); os.IsNotExist(err) {
			t.Error("JWT middleware file was not created")
		}
	})

	t.Run("GenerateHandlers", func(t *testing.T) {
		err := generator.generateHandlers()
		if err != nil {
			t.Errorf("Failed to generate handlers: %v", err)
		}

		// Check if auth handler file was created
		handlerPath := filepath.Join(tempDir, "internal", "handlers", "auth.go")
		if _, err := os.Stat(handlerPath); os.IsNotExist(err) {
			t.Error("Auth handler file was not created")
		}
	})

	t.Run("GenerateServices", func(t *testing.T) {
		err := generator.generateServices()
		if err != nil {
			t.Errorf("Failed to generate services: %v", err)
		}

		// Check if auth service file was created
		servicePath := filepath.Join(tempDir, "internal", "services", "auth.go")
		if _, err := os.Stat(servicePath); os.IsNotExist(err) {
			t.Error("Auth service file was not created")
		}
	})

	t.Run("GenerateModels", func(t *testing.T) {
		err := generator.generateModels()
		if err != nil {
			t.Errorf("Failed to generate models: %v", err)
		}

		// Check if user model file was created
		modelPath := filepath.Join(tempDir, "internal", "models", "user.go")
		if _, err := os.Stat(modelPath); os.IsNotExist(err) {
			t.Error("User model file was not created")
		}
	})

	t.Run("GenerateRepositories", func(t *testing.T) {
		err := generator.generateRepositories()
		if err != nil {
			t.Errorf("Failed to generate repositories: %v", err)
		}

		// Check if user repository file was created
		repoPath := filepath.Join(tempDir, "internal", "repositories", "user.go")
		if _, err := os.Stat(repoPath); os.IsNotExist(err) {
			t.Error("User repository file was not created")
		}
	})
}

// TestAuthGeneratorWithOAuth2 tests OAuth2 generation
func TestAuthGeneratorWithOAuth2(t *testing.T) {
	tempDir := t.TempDir()

	// Create OAuth2 providers
	oauth2Providers := []models.OAuth2Provider{
		{
			Name:         "google",
			ClientID:     "google-client-id",
			ClientSecret: "google-client-secret",
			RedirectURL:  "http://localhost:8080/auth/google/callback",
			Scopes:       []string{"openid", "email", "profile"},
			AuthURL:      "https://accounts.google.com/o/oauth2/auth",
			TokenURL:     "https://oauth2.googleapis.com/token",
			UserInfoURL:  "https://www.googleapis.com/oauth2/v2/userinfo",
		},
		{
			Name:         "github",
			ClientID:     "github-client-id",
			ClientSecret: "github-client-secret",
			RedirectURL:  "http://localhost:8080/auth/github/callback",
			Scopes:       []string{"user:email"},
			AuthURL:      "https://github.com/login/oauth/authorize",
			TokenURL:     "https://github.com/login/oauth/access_token",
			UserInfoURL:  "https://api.github.com/user",
		},
	}

	opts := DefaultAuthGeneratorOptions("test-oauth2", tempDir)
	opts = WithOAuth2(opts, oauth2Providers)

	generator := NewAuthGenerator(opts)

	t.Run("GenerateOAuth2Components", func(t *testing.T) {
		err := generator.generateOAuth2Components()
		if err != nil {
			t.Errorf("Failed to generate OAuth2 components: %v", err)
		}

		// Check if OAuth2 middleware was created
		oauth2MiddlewarePath := filepath.Join(tempDir, "internal", "middleware", "oauth2.go")
		if _, err := os.Stat(oauth2MiddlewarePath); os.IsNotExist(err) {
			t.Error("OAuth2 middleware file was not created")
		}

		// Check if OAuth2 handlers were created
		oauth2HandlerPath := filepath.Join(tempDir, "internal", "handlers", "oauth2.go")
		if _, err := os.Stat(oauth2HandlerPath); os.IsNotExist(err) {
			t.Error("OAuth2 handler file was not created")
		}
	})
}

// TestAuthGeneratorWithRBAC tests RBAC generation
func TestAuthGeneratorWithRBAC(t *testing.T) {
	tempDir := t.TempDir()

	opts := DefaultAuthGeneratorOptions("test-rbac", tempDir)
	opts = WithRBAC(opts)

	generator := NewAuthGenerator(opts)

	t.Run("GenerateRBACModels", func(t *testing.T) {
		err := generator.generateModels()
		if err != nil {
			t.Errorf("Failed to generate RBAC models: %v", err)
		}

		// Check if user model with RBAC relationships was created
		modelPath := filepath.Join(tempDir, "internal", "models", "user.go")
		if _, err := os.Stat(modelPath); os.IsNotExist(err) {
			t.Error("User model with RBAC file was not created")
		}

		// Read file content to verify RBAC fields are included
		content, err := os.ReadFile(modelPath)
		if err != nil {
			t.Errorf("Failed to read generated model file: %v", err)
		}

		contentStr := string(content)
		if !containsString(contentStr, "Roles []Role") {
			t.Error("Generated user model does not contain RBAC relationship")
		}
	})
}

// TestAuthGeneratorWithSupabase tests Supabase integration
func TestAuthGeneratorWithSupabase(t *testing.T) {
	tempDir := t.TempDir()

	supabaseConfig := models.SupabaseAuthConfig{
		ProjectURL:      "https://test-project.supabase.co",
		APIKey:          "test-api-key",
		ServiceKey:      "test-service-key",
		JWTSecret:       "test-jwt-secret",
		EnableAuth:      true,
		EnableSocial:    true,
		SocialProviders: []string{"google", "github"},
	}

	opts := DefaultAuthGeneratorOptions("test-supabase", tempDir)
	opts = WithSupabase(opts, supabaseConfig)

	generator := NewAuthGenerator(opts)

	t.Run("GenerateSupabaseIntegration", func(t *testing.T) {
		err := generator.generateSupabaseIntegration()
		if err != nil {
			t.Errorf("Failed to generate Supabase integration: %v", err)
		}

		// Check if Supabase auth service was created
		supabaseServicePath := filepath.Join(tempDir, "internal", "services", "supabase_auth.go")
		if _, err := os.Stat(supabaseServicePath); os.IsNotExist(err) {
			t.Error("Supabase auth service file was not created")
		}
	})
}

// TestDefaultAuthGeneratorOptions tests default options generation
func TestDefaultAuthGeneratorOptions(t *testing.T) {
	opts := DefaultAuthGeneratorOptions("test-project", "/tmp/test")

	// Validate basic options
	if opts.ProjectName != "test-project" {
		t.Errorf("Expected project name 'test-project', got '%s'", opts.ProjectName)
	}

	if opts.OutputPath != "/tmp/test" {
		t.Errorf("Expected output path '/tmp/test', got '%s'", opts.OutputPath)
	}

	if opts.AuthConfig.Provider != models.AuthProviderJWT {
		t.Errorf("Expected JWT provider, got '%s'", opts.AuthConfig.Provider)
	}

	if opts.DatabaseProvider != "postgres" {
		t.Errorf("Expected postgres provider, got '%s'", opts.DatabaseProvider)
	}

	// Validate user model defaults
	if opts.UserModel.TableName != "users" {
		t.Errorf("Expected table name 'users', got '%s'", opts.UserModel.TableName)
	}

	if opts.UserModel.StructName != "User" {
		t.Errorf("Expected struct name 'User', got '%s'", opts.UserModel.StructName)
	}

	if opts.UserModel.PrimaryKeyType != "uint" {
		t.Errorf("Expected primary key type 'uint', got '%s'", opts.UserModel.PrimaryKeyType)
	}
}

// TestWithRBAC tests RBAC configuration addition
func TestWithRBAC(t *testing.T) {
	opts := DefaultAuthGeneratorOptions("test-project", "/tmp/test")
	rbacOpts := WithRBAC(opts)

	// Check if RBAC is enabled
	if !rbacOpts.AuthConfig.EnableRBAC {
		t.Error("RBAC should be enabled")
	}

	if !rbacOpts.AuthConfig.EnablePermissions {
		t.Error("Permissions should be enabled")
	}

	// Check if models are set
	if rbacOpts.RoleModel == nil {
		t.Error("Role model should be set")
	}

	if rbacOpts.PermissionModel == nil {
		t.Error("Permission model should be set")
	}

	if rbacOpts.UserRoleModel == nil {
		t.Error("User role model should be set")
	}

	if rbacOpts.RolePermModel == nil {
		t.Error("Role permission model should be set")
	}

	// Validate role model defaults
	if rbacOpts.RoleModel.StructName != "Role" {
		t.Errorf("Expected role struct name 'Role', got '%s'", rbacOpts.RoleModel.StructName)
	}

	if rbacOpts.PermissionModel.StructName != "Permission" {
		t.Errorf("Expected permission struct name 'Permission', got '%s'", rbacOpts.PermissionModel.StructName)
	}
}

// TestWithOAuth2 tests OAuth2 configuration addition
func TestWithOAuth2(t *testing.T) {
	providers := []models.OAuth2Provider{
		{
			Name:        "google",
			ClientID:    "google-id",
			Scopes:      []string{"openid", "email"},
			AuthURL:     "https://accounts.google.com/o/oauth2/auth",
			TokenURL:    "https://oauth2.googleapis.com/token",
			UserInfoURL: "https://www.googleapis.com/oauth2/v2/userinfo",
		},
	}

	opts := DefaultAuthGeneratorOptions("test-project", "/tmp/test")
	oauth2Opts := WithOAuth2(opts, providers)

	// Check if OAuth2 providers are set
	if len(oauth2Opts.AuthConfig.OAuth2Providers) != 1 {
		t.Errorf("Expected 1 OAuth2 provider, got %d", len(oauth2Opts.AuthConfig.OAuth2Providers))
	}

	provider := oauth2Opts.AuthConfig.OAuth2Providers[0]
	if provider.Name != "google" {
		t.Errorf("Expected provider name 'google', got '%s'", provider.Name)
	}

	if provider.ClientID != "google-id" {
		t.Errorf("Expected client ID 'google-id', got '%s'", provider.ClientID)
	}
}

// TestWithSupabase tests Supabase configuration
func TestWithSupabase(t *testing.T) {
	supabaseConfig := models.SupabaseAuthConfig{
		ProjectURL:   "https://test.supabase.co",
		APIKey:       "test-key",
		EnableAuth:   true,
		EnableSocial: true,
	}

	opts := DefaultAuthGeneratorOptions("test-project", "/tmp/test")
	supabaseOpts := WithSupabase(opts, supabaseConfig)

	// Check if provider is set to Supabase
	if supabaseOpts.AuthConfig.Provider != models.AuthProviderSupabase {
		t.Errorf("Expected Supabase provider, got '%s'", supabaseOpts.AuthConfig.Provider)
	}

	// Check if database provider is set to Supabase
	if supabaseOpts.DatabaseProvider != "supabase" {
		t.Errorf("Expected supabase database provider, got '%s'", supabaseOpts.DatabaseProvider)
	}

	// Check if Supabase config is set
	if supabaseOpts.AuthConfig.Supabase.ProjectURL != "https://test.supabase.co" {
		t.Errorf("Expected project URL 'https://test.supabase.co', got '%s'", supabaseOpts.AuthConfig.Supabase.ProjectURL)
	}

	if !supabaseOpts.AuthConfig.Supabase.EnableAuth {
		t.Error("Supabase auth should be enabled")
	}
}

// TestTemplateData tests template data generation
func TestTemplateData(t *testing.T) {
	opts := DefaultAuthGeneratorOptions("test-project", "/tmp/test")
	opts = WithRBAC(opts)

	generator := NewAuthGenerator(opts)
	data := generator.getTemplateData()

	// Check basic data
	if projectName, ok := data["ProjectName"]; !ok || projectName != "test-project" {
		t.Errorf("Expected project name 'test-project', got '%v'", projectName)
	}

	if databaseProvider, ok := data["DatabaseProvider"]; !ok || databaseProvider != "postgres" {
		t.Errorf("Expected database provider 'postgres', got '%v'", databaseProvider)
	}

	// Check user model
	if userModelData, ok := data["UserModel"]; ok {
		userModel := userModelData.(models.UserModel)
		if userModel.StructName != "User" {
			t.Errorf("Expected user struct name 'User', got '%s'", userModel.StructName)
		}
		if userModel.PrimaryKeyType != "uint" {
			t.Errorf("Expected primary key type 'uint', got '%s'", userModel.PrimaryKeyType)
		}
	} else {
		t.Error("UserModel data should be present")
	}

	// Check role model (should be present because of WithRBAC)
	if roleModelData, ok := data["RoleModel"]; ok {
		roleModel := roleModelData.(models.RoleModel)
		if roleModel.StructName != "Role" {
			t.Errorf("Expected role struct name 'Role', got '%s'", roleModel.StructName)
		}
	} else {
		t.Error("RoleModel data should be present")
	}
}

// TestGetModuleName tests module name extraction
func TestGetModuleName(t *testing.T) {
	tempDir := t.TempDir()

	// Create a fake go.mod file
	goModContent := `module github.com/test/project

go 1.21

require (
    github.com/some/dependency v1.0.0
)`

	goModPath := filepath.Join(tempDir, "go.mod")
	err := os.WriteFile(goModPath, []byte(goModContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create go.mod file: %v", err)
	}

	opts := DefaultAuthGeneratorOptions("test-project", tempDir)
	generator := NewAuthGenerator(opts)

	moduleName := generator.getModuleName()
	expectedModule := "github.com/test/project"

	if moduleName != expectedModule {
		t.Errorf("Expected module name '%s', got '%s'", expectedModule, moduleName)
	}
}

// Helper function for string contains check
func containsString(s, substr string) bool {
	return len(substr) == 0 || len(s) >= len(substr) && (s == substr || containsString(s[1:], substr) || (len(s) > 0 && containsString(s[1:], substr)))
}

// Benchmark tests
func BenchmarkAuthGenerator(b *testing.B) {
	tempDir := b.TempDir()
	opts := DefaultAuthGeneratorOptions("bench-test", tempDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generator := NewAuthGenerator(opts)
		_ = generator.generateMiddleware()
	}
}

// Test validation
func TestAuthConfigValidation(t *testing.T) {
	t.Run("ValidConfig", func(t *testing.T) {
		config := models.DefaultAuthConfig()
		err := models.ValidateAuthConfig(config)
		if err != nil {
			t.Errorf("Expected valid config, got error: %v", err)
		}
	})

	t.Run("InvalidProvider", func(t *testing.T) {
		config := models.DefaultAuthConfig()
		config.Provider = "" // Invalid empty provider
		err := models.ValidateAuthConfig(config)
		if err == nil {
			t.Error("Expected validation error for empty provider")
		}
	})

	t.Run("MissingJWTSecret", func(t *testing.T) {
		config := models.DefaultAuthConfig()
		config.Provider = models.AuthProviderJWT
		config.JWTSecret = "" // Missing JWT secret
		err := models.ValidateAuthConfig(config)
		if err == nil {
			t.Error("Expected validation error for missing JWT secret")
		}
	})

	t.Run("WeakPassword", func(t *testing.T) {
		config := models.DefaultAuthConfig()
		config.PasswordMinLength = 3 // Too weak
		err := models.ValidateAuthConfig(config)
		if err == nil {
			t.Error("Expected validation error for weak password requirements")
		}
	})
}