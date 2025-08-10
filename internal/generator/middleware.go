package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/vibercode/cli/internal/models"
	"github.com/vibercode/cli/internal/templates"
	"github.com/vibercode/cli/pkg/ui"
)

// MiddlewareOptions contains configuration for middleware generation
type MiddlewareOptions struct {
	Type   string
	Name   string
	Custom bool
	Preset string
}

// MiddlewareGenerator handles middleware generation
type MiddlewareGenerator struct {
	options MiddlewareOptions
}

// NewMiddlewareGenerator creates a new middleware generator
func NewMiddlewareGenerator() *MiddlewareGenerator {
	return &MiddlewareGenerator{}
}

// Generate generates middleware based on options
func (g *MiddlewareGenerator) Generate(options MiddlewareOptions) error {
	g.options = options

	ui.PrintStep(1, 1, "Starting middleware generation...")

	// Handle preset generation
	if options.Preset != "" {
		return g.generatePreset()
	}

	// Handle custom middleware generation
	if options.Custom {
		return g.generateCustomMiddleware()
	}

	// Handle standard middleware generation
	if options.Type != "" {
		return g.generateStandardMiddleware()
	}

	// Interactive mode - ask user what they want to generate
	return g.generateInteractiveMiddleware()
}

// generatePreset generates a preset of middleware
func (g *MiddlewareGenerator) generatePreset() error {
	if !models.IsValidPreset(g.options.Preset) {
		return fmt.Errorf("invalid preset: %s", g.options.Preset)
	}

	ui.PrintStep(1, 3, fmt.Sprintf("Generating %s preset...", g.options.Preset))

	middlewares := models.GetDefaultPresetMiddlewares(models.MiddlewarePreset(g.options.Preset))
	
	// Create middleware directory
	if err := g.createMiddlewareDirectory(); err != nil {
		return fmt.Errorf("failed to create middleware directory: %w", err)
	}

	ui.PrintStep(2, 3, "Generating middleware files...")
	
	// Generate each middleware in the preset
	for _, middleware := range middlewares {
		if err := g.generateMiddlewareFile(middleware); err != nil {
			return fmt.Errorf("failed to generate %s middleware: %w", middleware.Name, err)
		}
		ui.PrintFileCreated(fmt.Sprintf("internal/middleware/%s", middleware.GetFileName()))
	}

	ui.PrintStep(3, 3, "Generating middleware registry...")
	
	// Generate middleware registry
	if err := g.generateMiddlewareRegistry(middlewares); err != nil {
		return fmt.Errorf("failed to generate middleware registry: %w", err)
	}
	ui.PrintFileCreated("internal/middleware/middleware.go")

	// Generate configuration
	if err := g.generateMiddlewareConfig(middlewares); err != nil {
		return fmt.Errorf("failed to generate middleware config: %w", err)
	}
	ui.PrintFileCreated("internal/config/middleware.go")

	ui.PrintSuccess(fmt.Sprintf("Preset %s generated successfully!", g.options.Preset))
	g.showPresetSummary(middlewares)

	return nil
}

// generateCustomMiddleware generates custom middleware
func (g *MiddlewareGenerator) generateCustomMiddleware() error {
	ui.PrintStep(1, 2, "Generating custom middleware...")

	// Get custom middleware configuration
	config, err := g.getCustomMiddlewareConfig()
	if err != nil {
		return fmt.Errorf("failed to get custom middleware config: %w", err)
	}

	// Create middleware directory
	if err := g.createMiddlewareDirectory(); err != nil {
		return fmt.Errorf("failed to create middleware directory: %w", err)
	}

	ui.PrintStep(2, 2, "Creating middleware file...")

	// Generate middleware file
	if err := g.generateMiddlewareFile(config); err != nil {
		return fmt.Errorf("failed to generate middleware file: %w", err)
	}

	ui.PrintSuccess(fmt.Sprintf("Custom middleware %s generated successfully!", config.Name))
	ui.PrintFileCreated(fmt.Sprintf("internal/middleware/%s", config.GetFileName()))

	return nil
}

// generateStandardMiddleware generates a single standard middleware
func (g *MiddlewareGenerator) generateStandardMiddleware() error {
	if !models.IsValidMiddlewareType(g.options.Type) {
		return fmt.Errorf("invalid middleware type: %s", g.options.Type)
	}

	ui.PrintStep(1, 2, fmt.Sprintf("Generating %s middleware...", g.options.Type))

	// Get standard middleware configuration
	config, err := g.getStandardMiddlewareConfig()
	if err != nil {
		return fmt.Errorf("failed to get middleware config: %w", err)
	}

	// Create middleware directory
	if err := g.createMiddlewareDirectory(); err != nil {
		return fmt.Errorf("failed to create middleware directory: %w", err)
	}

	ui.PrintStep(2, 2, "Creating middleware file...")

	// Generate middleware file
	if err := g.generateMiddlewareFile(config); err != nil {
		return fmt.Errorf("failed to generate middleware file: %w", err)
	}

	ui.PrintSuccess(fmt.Sprintf("Middleware %s generated successfully!", config.Name))
	ui.PrintFileCreated(fmt.Sprintf("internal/middleware/%s", config.GetFileName()))

	return nil
}

// generateInteractiveMiddleware handles interactive middleware generation
func (g *MiddlewareGenerator) generateInteractiveMiddleware() error {
	ui.PrintStep(1, 1, "Interactive middleware generation...")

	// Ask what type of middleware to generate
	typePrompt := promptui.Select{
		Label: ui.IconGear + " What would you like to generate?",
		Items: []string{
			"Single middleware",
			"Multiple middleware (preset)",
			"Custom middleware",
		},
	}

	_, choice, err := typePrompt.Run()
	if err != nil {
		return err
	}

	switch choice {
	case "Single middleware":
		return g.generateInteractiveSingleMiddleware()
	case "Multiple middleware (preset)":
		return g.generateInteractivePreset()
	case "Custom middleware":
		g.options.Custom = true
		return g.generateCustomMiddleware()
	}

	return nil
}

// generateInteractiveSingleMiddleware handles interactive single middleware generation
func (g *MiddlewareGenerator) generateInteractiveSingleMiddleware() error {
	typePrompt := promptui.Select{
		Label: ui.IconCode + " Select middleware type",
		Items: []string{"auth", "logging", "cors", "rate-limit"},
	}

	_, middlewareType, err := typePrompt.Run()
	if err != nil {
		return err
	}

	g.options.Type = middlewareType
	return g.generateStandardMiddleware()
}

// generateInteractivePreset handles interactive preset generation
func (g *MiddlewareGenerator) generateInteractivePreset() error {
	presetPrompt := promptui.Select{
		Label: ui.IconPackage + " Select middleware preset",
		Items: []string{"api-security", "web-app", "microservice", "public-api"},
	}

	_, preset, err := presetPrompt.Run()
	if err != nil {
		return err
	}

	g.options.Preset = preset
	return g.generatePreset()
}

// getCustomMiddlewareConfig gets configuration for custom middleware
func (g *MiddlewareGenerator) getCustomMiddlewareConfig() (models.MiddlewareConfig, error) {
	var config models.MiddlewareConfig

	// Middleware name
	namePrompt := promptui.Prompt{
		Label: ui.IconCode + " Middleware name",
	}
	name, err := namePrompt.Run()
	if err != nil {
		return config, err
	}

	// Description
	descPrompt := promptui.Prompt{
		Label:   ui.IconDoc + " Description",
		Default: fmt.Sprintf("Custom %s middleware", name),
	}
	desc, err := descPrompt.Run()
	if err != nil {
		return config, err
	}

	config = models.MiddlewareConfig{
		Name:        strings.Title(name),
		Type:        models.CustomMiddleware,
		Description: desc,
		Options:     models.MiddlewareOptions{},
	}

	return config, nil
}

// getStandardMiddlewareConfig gets configuration for standard middleware
func (g *MiddlewareGenerator) getStandardMiddlewareConfig() (models.MiddlewareConfig, error) {
	var config models.MiddlewareConfig
	middlewareType := models.MiddlewareType(g.options.Type)

	config.Type = middlewareType
	config.Name = strings.Title(string(middlewareType))

	switch middlewareType {
	case models.AuthMiddleware:
		return g.getAuthMiddlewareConfig(config)
	case models.LoggingMiddleware:
		return g.getLoggingMiddlewareConfig(config)
	case models.CORSMiddleware:
		return g.getCORSMiddlewareConfig(config)
	case models.RateLimitMiddleware:
		return g.getRateLimitMiddlewareConfig(config)
	}

	return config, nil
}

// getAuthMiddlewareConfig gets auth middleware specific configuration
func (g *MiddlewareGenerator) getAuthMiddlewareConfig(config models.MiddlewareConfig) (models.MiddlewareConfig, error) {
	config.Description = "Authentication middleware"

	// Auth strategy
	strategyPrompt := promptui.Select{
		Label: ui.IconGear + " Authentication strategy",
		Items: []string{"jwt", "apikey", "session", "basic"},
	}

	_, strategy, err := strategyPrompt.Run()
	if err != nil {
		return config, err
	}

	config.Options.AuthStrategy = models.AuthStrategy(strategy)

	// Strategy-specific configuration
	switch models.AuthStrategy(strategy) {
	case models.JWTAuth:
		config.Options.JWTSecret = "${JWT_SECRET}"
		config.Options.JWTIssuer = "vibercode-api"
	case models.APIKeyAuth:
		config.Options.APIKeyHeader = "X-API-Key"
	}

	return config, nil
}

// getLoggingMiddlewareConfig gets logging middleware configuration
func (g *MiddlewareGenerator) getLoggingMiddlewareConfig(config models.MiddlewareConfig) (models.MiddlewareConfig, error) {
	config.Description = "Request/response logging middleware"

	// Log level
	levelPrompt := promptui.Select{
		Label: ui.IconDoc + " Log level",
		Items: []string{"debug", "info", "warn", "error"},
	}

	_, level, err := levelPrompt.Run()
	if err != nil {
		return config, err
	}

	config.Options.LogLevel = level
	config.Options.LogFormat = "json"
	config.Options.LogRequests = true
	config.Options.LogResponses = false
	config.Options.ExcludePaths = []string{"/health", "/metrics"}

	return config, nil
}

// getCORSMiddlewareConfig gets CORS middleware configuration
func (g *MiddlewareGenerator) getCORSMiddlewareConfig(config models.MiddlewareConfig) (models.MiddlewareConfig, error) {
	config.Description = "CORS configuration middleware"

	config.Options.AllowedOrigins = []string{"*"}
	config.Options.AllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.Options.AllowedHeaders = []string{"Authorization", "Content-Type"}
	config.Options.MaxAge = 3600

	return config, nil
}

// getRateLimitMiddlewareConfig gets rate limit middleware configuration
func (g *MiddlewareGenerator) getRateLimitMiddlewareConfig(config models.MiddlewareConfig) (models.MiddlewareConfig, error) {
	config.Description = "Rate limiting middleware"

	config.Options.RequestsPerSecond = 100
	config.Options.BurstSize = 10
	config.Options.Strategy = "token_bucket"

	return config, nil
}

// createMiddlewareDirectory creates the middleware directory structure
func (g *MiddlewareGenerator) createMiddlewareDirectory() error {
	dirs := []string{
		"internal/middleware",
		"internal/config",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// generateMiddlewareFile generates a single middleware file
func (g *MiddlewareGenerator) generateMiddlewareFile(config models.MiddlewareConfig) error {
	content := templates.GetMiddlewareTemplate(config)
	filename := filepath.Join("internal/middleware", config.GetFileName())

	return g.writeFile(filename, content)
}

// generateMiddlewareRegistry generates the middleware registry file
func (g *MiddlewareGenerator) generateMiddlewareRegistry(middlewares []models.MiddlewareConfig) error {
	content := templates.GetMiddlewareRegistryTemplate(middlewares)
	filename := "internal/middleware/middleware.go"

	return g.writeFile(filename, content)
}

// generateMiddlewareConfig generates middleware configuration file
func (g *MiddlewareGenerator) generateMiddlewareConfig(middlewares []models.MiddlewareConfig) error {
	content := templates.GetMiddlewareConfigTemplate(middlewares)
	filename := "internal/config/middleware.go"

	return g.writeFile(filename, content)
}

// writeFile writes content to a file
func (g *MiddlewareGenerator) writeFile(filePath, content string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Write file
	return os.WriteFile(filePath, []byte(content), 0644)
}

// showPresetSummary shows a summary of the generated preset
func (g *MiddlewareGenerator) showPresetSummary(middlewares []models.MiddlewareConfig) {
	ui.PrintInfo(fmt.Sprintf("Generated %d middleware components:", len(middlewares)))
	
	for _, middleware := range middlewares {
		fmt.Printf("  %s %s: %s\n", 
			ui.IconGear, 
			ui.Bold.Sprint(middleware.Name), 
			ui.Muted.Sprint(middleware.Description))
	}
	
	ui.PrintSeparator()
	ui.PrintInfo("Next steps:")
	fmt.Println("  1. Update your main.go to register middleware")
	fmt.Println("  2. Configure environment variables for middleware settings")
	fmt.Println("  3. Test middleware functionality")
}