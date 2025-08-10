package config

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/vibercode/cli/internal/templates"
)

// GeneratorOptions contains options for configuration generation
type GeneratorOptions struct {
	ProjectName      string
	OutputPath       string
	Environment      string
	DatabaseProvider string
	DatabasePort     int
	Domain           string
	GenerateEnvFiles bool
	GenerateDocker   bool
	GenerateMakefile bool
}

// Generator handles configuration file generation
type Generator struct {
	options GeneratorOptions
	config  *Config
}

// NewGenerator creates a new configuration generator
func NewGenerator(opts GeneratorOptions) *Generator {
	return &Generator{
		options: opts,
		config:  DefaultConfig(),
	}
}

// GenerateProjectConfig generates all configuration files for a project
func (g *Generator) GenerateProjectConfig() error {
	// Ensure output directory exists
	if err := os.MkdirAll(g.options.OutputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Customize default config based on options
	g.customizeConfig()

	// Generate main configuration file
	if err := g.generateMainConfig(); err != nil {
		return fmt.Errorf("failed to generate main config: %w", err)
	}

	// Generate environment-specific configurations
	if err := g.generateEnvironmentConfigs(); err != nil {
		return fmt.Errorf("failed to generate environment configs: %w", err)
	}

	// Generate environment files if requested
	if g.options.GenerateEnvFiles {
		if err := g.generateEnvFiles(); err != nil {
			return fmt.Errorf("failed to generate env files: %w", err)
		}
	}

	// Generate Docker configuration if requested
	if g.options.GenerateDocker {
		if err := g.generateDockerConfig(); err != nil {
			return fmt.Errorf("failed to generate Docker config: %w", err)
		}
	}

	// Generate Makefile if requested
	if g.options.GenerateMakefile {
		if err := g.generateMakefile(); err != nil {
			return fmt.Errorf("failed to generate Makefile: %w", err)
		}
	}

	return nil
}

// customizeConfig customizes the default configuration based on options
func (g *Generator) customizeConfig() {
	// Set environment
	g.config.Environment = g.options.Environment

	// Customize database configuration
	g.config.Database.Provider = g.options.DatabaseProvider
	switch g.options.DatabaseProvider {
	case "postgres":
		g.config.Database.Port = 5432
		g.config.Database.Username = g.options.ProjectName
		g.config.Database.Database = g.options.ProjectName
		g.config.Database.SSLMode = "require"
	case "mysql":
		g.config.Database.Port = 3306
		g.config.Database.Username = g.options.ProjectName
		g.config.Database.Database = g.options.ProjectName
	case "mongodb":
		g.config.Database.Port = 27017
		g.config.Database.Username = g.options.ProjectName
		g.config.Database.Database = g.options.ProjectName
	case "sqlite":
		g.config.Database.Database = fmt.Sprintf("%s.db", g.options.ProjectName)
	case "supabase":
		g.config.Database.Supabase.EnableAuth = true
		g.config.Database.Supabase.EnableStorage = true
		g.config.Database.Supabase.EnableRealtime = false
	}

	// Override with custom port if provided
	if g.options.DatabasePort > 0 {
		g.config.Database.Port = g.options.DatabasePort
	}

	// Environment-specific customizations
	switch g.options.Environment {
	case "development":
		g.config.Logging.Level = "debug"
		g.config.Logging.Format = "text"
		g.config.Security.EnableRateLimiting = false
		g.config.Security.EnableCSRF = false
		g.config.Server.CORS.AllowedOrigins = []string{"*"}
	case "staging":
		g.config.Logging.Level = "info"
		g.config.Logging.Format = "json"
		g.config.Security.EnableRateLimiting = true
		g.config.Security.EnableCSRF = true
		if g.options.Domain != "" {
			g.config.Server.CORS.AllowedOrigins = []string{
				fmt.Sprintf("https://staging.%s", g.options.Domain),
			}
		}
	case "production":
		g.config.Logging.Level = "warn"
		g.config.Logging.Format = "json"
		g.config.Logging.Output = "file"
		g.config.Logging.Filename = fmt.Sprintf("/var/log/%s/app.log", g.options.ProjectName)
		g.config.Security.EnableRateLimiting = true
		g.config.Security.EnableCSRF = true
		g.config.Server.EnableHTTPS = true
		g.config.Database.Migrations.AutoMigrate = false
		g.config.Database.Migrations.BackupBefore = true
		if g.options.Domain != "" {
			g.config.Server.CORS.AllowedOrigins = []string{
				fmt.Sprintf("https://%s", g.options.Domain),
			}
		}
	}
}

// generateMainConfig generates the main configuration file
func (g *Generator) generateMainConfig() error {
	configPath := filepath.Join(g.options.OutputPath, "config.json")
	return g.generateConfigFile(templates.ConfigTemplate, configPath, g.config)
}

// generateEnvironmentConfigs generates environment-specific configuration files
func (g *Generator) generateEnvironmentConfigs() error {
	environments := []struct {
		name     string
		template string
	}{
		{"development", templates.DevelopmentConfigTemplate},
		{"staging", templates.StagingConfigTemplate},
		{"production", templates.ProductionConfigTemplate},
	}

	for _, env := range environments {
		templateData := struct {
			ProjectName      string
			DatabaseProvider string
			DatabasePort     int
			Domain          string
		}{
			ProjectName:      g.options.ProjectName,
			DatabaseProvider: g.options.DatabaseProvider,
			DatabasePort:     g.getDatabasePort(),
			Domain:          g.options.Domain,
		}

		configPath := filepath.Join(g.options.OutputPath, fmt.Sprintf("config.%s.json", env.name))
		if err := g.generateConfigFile(env.template, configPath, templateData); err != nil {
			return err
		}
	}

	return nil
}

// generateEnvFiles generates environment variable files
func (g *Generator) generateEnvFiles() error {
	environments := []string{"development", "staging", "production"}

	for _, env := range environments {
		// Create environment-specific config
		envConfig := *g.config
		envConfig.Environment = env

		// Customize for environment
		switch env {
		case "development":
			envConfig.Database.Database = fmt.Sprintf("%s_dev", g.options.ProjectName)
			envConfig.Database.SSLMode = "disable"
		case "staging":
			envConfig.Database.Database = fmt.Sprintf("%s_staging", g.options.ProjectName)
		case "production":
			envConfig.Database.Database = fmt.Sprintf("%s_prod", g.options.ProjectName)
		}

		envPath := filepath.Join(g.options.OutputPath, fmt.Sprintf(".env.%s", env))
		if err := g.generateConfigFile(templates.EnvTemplate, envPath, &envConfig); err != nil {
			return err
		}
	}

	// Generate example env file
	examplePath := filepath.Join(g.options.OutputPath, ".env.example")
	return g.generateConfigFile(templates.EnvTemplate, examplePath, g.config)
}

// generateDockerConfig generates Docker-related configuration files
func (g *Generator) generateDockerConfig() error {
	// Generate docker-compose.yml
	dockerComposePath := filepath.Join(g.options.OutputPath, "docker-compose.yml")
	if err := g.generateConfigFile(templates.DockerComposeTemplate, dockerComposePath, g.config); err != nil {
		return err
	}

	// Generate Dockerfile
	dockerfilePath := filepath.Join(g.options.OutputPath, "Dockerfile")
	return g.generateConfigFile(templates.DockerfileTemplate, dockerfilePath, g.config)
}

// generateMakefile generates Makefile for the project
func (g *Generator) generateMakefile() error {
	makefilePath := filepath.Join(g.options.OutputPath, "Makefile")
	return g.generateConfigFile(templates.MakefileTemplate, makefilePath, g.config)
}

// generateConfigFile generates a configuration file from a template
func (g *Generator) generateConfigFile(templateStr, outputPath string, data interface{}) error {
	// Create template with helper functions
	tmpl := template.New("config").Funcs(template.FuncMap{
		"toJSON": func(v interface{}) string {
			// Simple JSON serialization for arrays
			switch val := v.(type) {
			case []string:
				result := "["
				for i, item := range val {
					if i > 0 {
						result += ", "
					}
					result += fmt.Sprintf("\"%s\"", item)
				}
				result += "]"
				return result
			default:
				return fmt.Sprintf("%v", v)
			}
		},
	})

	// Parse template
	tmpl, err := tmpl.Parse(templateStr)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// getDatabasePort returns the appropriate database port
func (g *Generator) getDatabasePort() int {
	if g.options.DatabasePort > 0 {
		return g.options.DatabasePort
	}

	switch g.options.DatabaseProvider {
	case "postgres":
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

// GenerateEnvironmentTemplate generates environment-specific configuration template
func GenerateEnvironmentTemplate(environment, databaseProvider string) (*Config, error) {
	config := DefaultConfig()
	config.Environment = environment

	// Apply environment-specific defaults
	switch environment {
	case "development":
		config.Logging.Level = "debug"
		config.Logging.Format = "text"
		config.Security.EnableRateLimiting = false
		config.Database.SSLMode = "disable"
		config.Database.Migrations.AutoMigrate = true
	case "staging":
		config.Logging.Level = "info"
		config.Logging.Format = "json"
		config.Security.EnableRateLimiting = true
		config.Database.SSLMode = "require"
		config.Database.Migrations.AutoMigrate = false
		config.Database.Migrations.BackupBefore = true
	case "production":
		config.Logging.Level = "error"
		config.Logging.Format = "json"
		config.Logging.Output = "file"
		config.Security.EnableRateLimiting = true
		config.Security.EnableCSRF = true
		config.Server.EnableHTTPS = true
		config.Database.SSLMode = "require"
		config.Database.Migrations.AutoMigrate = false
		config.Database.Migrations.BackupBefore = true
		config.Database.MaxOpenConns = 50
		config.Database.MaxIdleConns = 25
	}

	// Apply database-specific defaults
	config.Database.Provider = databaseProvider
	switch databaseProvider {
	case "postgres":
		config.Database.Port = 5432
	case "mysql":
		config.Database.Port = 3306
	case "mongodb":
		config.Database.Port = 27017
	case "supabase":
		config.Database.Supabase.EnableAuth = true
		config.Database.Supabase.EnableStorage = true
	}

	return config, nil
}

// ValidateGeneratorOptions validates the generator options
func ValidateGeneratorOptions(opts GeneratorOptions) error {
	if opts.ProjectName == "" {
		return fmt.Errorf("project name is required")
	}

	if opts.OutputPath == "" {
		return fmt.Errorf("output path is required")
	}

	if opts.Environment == "" {
		opts.Environment = "development"
	}

	validEnvironments := []string{"development", "staging", "production"}
	validEnv := false
	for _, env := range validEnvironments {
		if opts.Environment == env {
			validEnv = true
			break
		}
	}
	if !validEnv {
		return fmt.Errorf("environment must be one of: %v", validEnvironments)
	}

	if opts.DatabaseProvider == "" {
		opts.DatabaseProvider = "postgres"
	}

	validProviders := []string{"postgres", "mysql", "sqlite", "mongodb", "supabase", "redis"}
	validProvider := false
	for _, provider := range validProviders {
		if opts.DatabaseProvider == provider {
			validProvider = true
			break
		}
	}
	if !validProvider {
		return fmt.Errorf("database provider must be one of: %v", validProviders)
	}

	return nil
}

// GetConfigurationSummary returns a summary of the generated configuration
func (g *Generator) GetConfigurationSummary() map[string]interface{} {
	return map[string]interface{}{
		"project_name":       g.options.ProjectName,
		"environment":        g.config.Environment,
		"server_port":        g.config.Server.Port,
		"database_provider":  g.config.Database.Provider,
		"database_port":      g.config.Database.Port,
		"auth_provider":      g.config.Auth.Provider,
		"storage_provider":   g.config.Storage.Provider,
		"cache_provider":     g.config.Cache.Provider,
		"logging_level":      g.config.Logging.Level,
		"monitoring_enabled": g.config.Monitoring.EnableMetrics,
		"security_features": map[string]bool{
			"rate_limiting": g.config.Security.EnableRateLimiting,
			"csrf":          g.config.Security.EnableCSRF,
			"https":         g.config.Server.EnableHTTPS,
		},
		"features": map[string]bool{
			"graphql":       g.config.Features.EnableGraphQL,
			"websocket":     g.config.Features.EnableWebSocket,
			"file_upload":   g.config.Features.EnableFileUpload,
			"notifications": g.config.Features.EnableNotifications,
			"search":        g.config.Features.EnableSearch,
			"caching":       g.config.Features.EnableCaching,
		},
		"generated_at": time.Now().Format("2006-01-02 15:04:05"),
	}
}