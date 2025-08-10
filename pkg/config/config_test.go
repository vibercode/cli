package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestDefaultConfig tests the default configuration
func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	
	// Test basic values
	if config.Environment != "development" {
		t.Errorf("Expected environment 'development', got %s", config.Environment)
	}
	
	if config.Server.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", config.Server.Port)
	}
	
	if config.Database.Provider != "postgres" {
		t.Errorf("Expected database provider 'postgres', got %s", config.Database.Provider)
	}
	
	if config.Logging.Level != "info" {
		t.Errorf("Expected log level 'info', got %s", config.Logging.Level)
	}
}

// TestConfigValidation tests configuration validation
func TestConfigValidation(t *testing.T) {
	validator := NewValidator()
	
	t.Run("ValidConfig", func(t *testing.T) {
		config := DefaultConfig()
		err := validator.ValidateConfig(config)
		if err != nil {
			t.Errorf("Expected valid config, got error: %v", err)
		}
	})
	
	t.Run("InvalidPort", func(t *testing.T) {
		config := DefaultConfig()
		config.Server.Port = 70000 // Invalid port
		err := validator.ValidateConfig(config)
		if err == nil {
			t.Error("Expected validation error for invalid port")
		}
	})
	
	t.Run("InvalidDatabaseProvider", func(t *testing.T) {
		config := DefaultConfig()
		config.Database.Provider = "invalid"
		err := validator.ValidateConfig(config)
		if err == nil {
			t.Error("Expected validation error for invalid database provider")
		}
	})
	
	t.Run("MissingJWTSecret", func(t *testing.T) {
		config := DefaultConfig()
		config.Auth.Provider = "jwt"
		config.Auth.JWTSecret = ""
		err := validator.ValidateConfig(config)
		if err == nil {
			t.Error("Expected validation error for missing JWT secret")
		}
	})
}

// TestConfigManager tests the configuration manager
func TestConfigManager(t *testing.T) {
	manager := NewManager()
	
	t.Run("LoadDefaultConfig", func(t *testing.T) {
		opts := LoadOptions{
			Environment: "development",
			Validate:    true,
		}
		
		err := manager.Load(opts)
		if err != nil {
			t.Errorf("Failed to load default config: %v", err)
		}
		
		config := manager.GetConfig()
		if config.Environment != "development" {
			t.Errorf("Expected environment 'development', got %s", config.Environment)
		}
	})
	
	t.Run("LoadFromJSON", func(t *testing.T) {
		// Create temporary config file
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "config.json")
		
		testConfig := map[string]interface{}{
			"environment": "testing",
			"server": map[string]interface{}{
				"host": "0.0.0.0",
				"port": 9000,
			},
			"database": map[string]interface{}{
				"provider": "sqlite",
				"database": "test.db",
			},
		}
		
		configData, _ := json.Marshal(testConfig)
		err := os.WriteFile(configPath, configData, 0644)
		if err != nil {
			t.Fatalf("Failed to write test config: %v", err)
		}
		
		opts := LoadOptions{
			ConfigPath:  configPath,
			Environment: "testing",
			Validate:    true,
		}
		
		err = manager.Load(opts)
		if err != nil {
			t.Errorf("Failed to load config from file: %v", err)
		}
		
		config := manager.GetConfig()
		if config.Environment != "testing" {
			t.Errorf("Expected environment 'testing', got %s", config.Environment)
		}
		if config.Server.Port != 9000 {
			t.Errorf("Expected port 9000, got %d", config.Server.Port)
		}
	})
	
	t.Run("LoadFromEnvFile", func(t *testing.T) {
		tempDir := t.TempDir()
		envPath := filepath.Join(tempDir, ".env.test")
		
		envContent := `SERVER_HOST=testhost
SERVER_PORT=8888
DATABASE_PROVIDER=mysql
LOGGING_LEVEL=debug`
		
		err := os.WriteFile(envPath, []byte(envContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write test env file: %v", err)
		}
		
		// Change to temp directory
		oldWd, _ := os.Getwd()
		os.Chdir(tempDir)
		defer os.Chdir(oldWd)
		
		opts := LoadOptions{
			Environment: "test",
			EnvFiles:    []string{".env.test"},
			Validate:    false, // Skip validation for this test
		}
		
		err = manager.Load(opts)
		if err != nil {
			t.Errorf("Failed to load config with env file: %v", err)
		}
		
		config := manager.GetConfig()
		if config.Server.Host != "testhost" {
			t.Errorf("Expected host 'testhost', got %s", config.Server.Host)
		}
		if config.Server.Port != 8888 {
			t.Errorf("Expected port 8888, got %d", config.Server.Port)
		}
	})
	
	t.Run("UpdateConfig", func(t *testing.T) {
		opts := LoadOptions{
			Environment: "development",
			Validate:    true,
		}
		
		err := manager.Load(opts)
		if err != nil {
			t.Fatalf("Failed to load initial config: %v", err)
		}
		
		err = manager.UpdateConfig(func(config *Config) error {
			config.Server.Port = 9999
			return nil
		})
		
		if err != nil {
			t.Errorf("Failed to update config: %v", err)
		}
		
		config := manager.GetConfig()
		if config.Server.Port != 9999 {
			t.Errorf("Expected port 9999 after update, got %d", config.Server.Port)
		}
	})
	
	t.Run("UpdateConfigValidationFailure", func(t *testing.T) {
		opts := LoadOptions{
			Environment: "development",
			Validate:    true,
		}
		
		err := manager.Load(opts)
		if err != nil {
			t.Fatalf("Failed to load initial config: %v", err)
		}
		
		originalPort := manager.GetConfig().Server.Port
		
		err = manager.UpdateConfig(func(config *Config) error {
			config.Server.Port = 70000 // Invalid port
			return nil
		})
		
		if err == nil {
			t.Error("Expected validation error when updating with invalid port")
		}
		
		// Config should remain unchanged
		config := manager.GetConfig()
		if config.Server.Port != originalPort {
			t.Errorf("Config should not have changed after validation failure")
		}
	})
}

// TestConfigGenerator tests the configuration generator
func TestConfigGenerator(t *testing.T) {
	t.Run("ValidateOptions", func(t *testing.T) {
		// Valid options
		opts := GeneratorOptions{
			ProjectName:      "test-project",
			OutputPath:       "/tmp/test",
			Environment:      "development",
			DatabaseProvider: "postgres",
		}
		
		err := ValidateGeneratorOptions(opts)
		if err != nil {
			t.Errorf("Expected valid options, got error: %v", err)
		}
		
		// Invalid project name
		opts.ProjectName = ""
		err = ValidateGeneratorOptions(opts)
		if err == nil {
			t.Error("Expected error for empty project name")
		}
		
		// Invalid environment
		opts.ProjectName = "test"
		opts.Environment = "invalid"
		err = ValidateGeneratorOptions(opts)
		if err == nil {
			t.Error("Expected error for invalid environment")
		}
		
		// Invalid database provider
		opts.Environment = "development"
		opts.DatabaseProvider = "invalid"
		err = ValidateGeneratorOptions(opts)
		if err == nil {
			t.Error("Expected error for invalid database provider")
		}
	})
	
	t.Run("GenerateConfig", func(t *testing.T) {
		tempDir := t.TempDir()
		
		opts := GeneratorOptions{
			ProjectName:      "test-project",
			OutputPath:       tempDir,
			Environment:      "development",
			DatabaseProvider: "postgres",
			GenerateEnvFiles: true,
			GenerateDocker:   true,
			GenerateMakefile: true,
		}
		
		generator := NewGenerator(opts)
		err := generator.GenerateProjectConfig()
		if err != nil {
			t.Errorf("Failed to generate project config: %v", err)
		}
		
		// Check that files were created
		expectedFiles := []string{
			"config.json",
			"config.development.json",
			"config.staging.json",
			"config.production.json",
			".env.development",
			".env.staging",
			".env.production",
			".env.example",
			"docker-compose.yml",
			"Dockerfile",
			"Makefile",
		}
		
		for _, file := range expectedFiles {
			filePath := filepath.Join(tempDir, file)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("Expected file %s was not created", file)
			}
		}
		
		// Check content of main config file
		configPath := filepath.Join(tempDir, "config.json")
		configData, err := os.ReadFile(configPath)
		if err != nil {
			t.Errorf("Failed to read generated config: %v", err)
		}
		
		if !strings.Contains(string(configData), "development") {
			t.Error("Generated config should contain environment setting")
		}
		if !strings.Contains(string(configData), "postgres") {
			t.Error("Generated config should contain database provider")
		}
	})
}

// TestEnvironmentTemplates tests environment-specific configuration templates
func TestEnvironmentTemplates(t *testing.T) {
	environments := []string{"development", "staging", "production"}
	
	for _, env := range environments {
		t.Run(env, func(t *testing.T) {
			config, err := GenerateEnvironmentTemplate(env, "postgres")
			if err != nil {
				t.Errorf("Failed to generate %s template: %v", env, err)
			}
			
			if config.Environment != env {
				t.Errorf("Expected environment %s, got %s", env, config.Environment)
			}
			
			// Environment-specific assertions
			switch env {
			case "development":
				if config.Logging.Level != "debug" {
					t.Errorf("Development should have debug logging, got %s", config.Logging.Level)
				}
				if config.Security.EnableRateLimiting {
					t.Error("Development should not have rate limiting enabled")
				}
			case "staging":
				if config.Logging.Level != "info" {
					t.Errorf("Staging should have info logging, got %s", config.Logging.Level)
				}
				if !config.Security.EnableRateLimiting {
					t.Error("Staging should have rate limiting enabled")
				}
			case "production":
				if config.Logging.Level != "error" {
					t.Errorf("Production should have error logging, got %s", config.Logging.Level)
				}
				if !config.Server.EnableHTTPS {
					t.Error("Production should have HTTPS enabled")
				}
				if !config.Security.EnableCSRF {
					t.Error("Production should have CSRF enabled")
				}
			}
		})
	}
}

// TestCustomValidators tests custom validation functions
func TestCustomValidators(t *testing.T) {
	validator := NewValidator()
	
	t.Run("PortValidation", func(t *testing.T) {
		// Valid ports
		validPorts := []int{1, 80, 443, 8080, 65535}
		for _, port := range validPorts {
			err := validator.ValidateField(port, "port")
			if err != nil {
				t.Errorf("Port %d should be valid, got error: %v", port, err)
			}
		}
		
		// Invalid ports
		invalidPorts := []int{0, -1, 65536, 100000}
		for _, port := range invalidPorts {
			err := validator.ValidateField(port, "port")
			if err == nil {
				t.Errorf("Port %d should be invalid", port)
			}
		}
	})
	
	t.Run("HostValidation", func(t *testing.T) {
		// Valid hosts
		validHosts := []string{"localhost", "127.0.0.1", "example.com", "192.168.1.1"}
		for _, host := range validHosts {
			err := validator.ValidateField(host, "host")
			if err != nil {
				t.Errorf("Host %s should be valid, got error: %v", host, err)
			}
		}
		
		// Invalid hosts
		invalidHosts := []string{"", "invalid..host", "host:with:colons"}
		for _, host := range invalidHosts {
			err := validator.ValidateField(host, "host")
			if err == nil {
				t.Errorf("Host %s should be invalid", host)
			}
		}
	})
	
	t.Run("DatabaseProviderValidation", func(t *testing.T) {
		// Valid providers
		validProviders := []string{"postgres", "mysql", "sqlite", "mongodb", "supabase", "redis"}
		for _, provider := range validProviders {
			err := validator.ValidateField(provider, "database_provider")
			if err != nil {
				t.Errorf("Provider %s should be valid, got error: %v", provider, err)
			}
		}
		
		// Invalid providers
		invalidProviders := []string{"", "oracle", "invalid"}
		for _, provider := range invalidProviders {
			err := validator.ValidateField(provider, "database_provider")
			if err == nil {
				t.Errorf("Provider %s should be invalid", provider)
			}
		}
	})
}

// TestConfigSerialization tests configuration serialization and deserialization
func TestConfigSerialization(t *testing.T) {
	manager := NewManager()
	
	t.Run("JSONSerialization", func(t *testing.T) {
		config := DefaultConfig()
		config.Environment = "test"
		config.Server.Port = 9999
		
		// Serialize to JSON
		data, err := json.Marshal(config)
		if err != nil {
			t.Errorf("Failed to marshal config to JSON: %v", err)
		}
		
		// Deserialize from JSON
		var newConfig Config
		err = json.Unmarshal(data, &newConfig)
		if err != nil {
			t.Errorf("Failed to unmarshal config from JSON: %v", err)
		}
		
		// Verify values
		if newConfig.Environment != "test" {
			t.Errorf("Expected environment 'test', got %s", newConfig.Environment)
		}
		if newConfig.Server.Port != 9999 {
			t.Errorf("Expected port 9999, got %d", newConfig.Server.Port)
		}
	})
	
	t.Run("SaveAndLoadConfig", func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "test_config.json")
		
		// Load default config
		opts := LoadOptions{
			Environment: "development",
			Validate:    true,
		}
		err := manager.Load(opts)
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}
		
		// Modify config
		err = manager.UpdateConfig(func(config *Config) error {
			config.Server.Port = 7777
			config.Database.Provider = "mysql"
			return nil
		})
		if err != nil {
			t.Fatalf("Failed to update config: %v", err)
		}
		
		// Save config
		err = manager.SaveConfig(configPath)
		if err != nil {
			t.Errorf("Failed to save config: %v", err)
		}
		
		// Load saved config
		newManager := NewManager()
		newOpts := LoadOptions{
			ConfigPath: configPath,
			Validate:   true,
		}
		err = newManager.Load(newOpts)
		if err != nil {
			t.Errorf("Failed to load saved config: %v", err)
		}
		
		// Verify loaded config
		loadedConfig := newManager.GetConfig()
		if loadedConfig.Server.Port != 7777 {
			t.Errorf("Expected port 7777, got %d", loadedConfig.Server.Port)
		}
		if loadedConfig.Database.Provider != "mysql" {
			t.Errorf("Expected provider 'mysql', got %s", loadedConfig.Database.Provider)
		}
	})
}

// BenchmarkConfigValidation benchmarks configuration validation
func BenchmarkConfigValidation(b *testing.B) {
	validator := NewValidator()
	config := DefaultConfig()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.ValidateConfig(config)
	}
}

// BenchmarkConfigLoad benchmarks configuration loading
func BenchmarkConfigLoad(b *testing.B) {
	opts := LoadOptions{
		Environment: "development",
		Validate:    true,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager := NewManager()
		manager.Load(opts)
	}
}