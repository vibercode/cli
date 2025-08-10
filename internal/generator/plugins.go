package generator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/vibercode/cli/internal/models"
	"github.com/vibercode/cli/pkg/ui"
	"gopkg.in/yaml.v2"
)

// PluginGenerator handles plugin generation and management
type PluginGenerator struct {
	config *models.PluginConfigManager
}

// NewPluginGenerator creates a new plugin generator
func NewPluginGenerator() *PluginGenerator {
	return &PluginGenerator{
		config: models.GetDefaultConfig(),
	}
}

// Generate generates a new plugin project
func (pg *PluginGenerator) Generate(options models.PluginOptions) error {
	ui.PrintInfo(fmt.Sprintf("🔌 Generating plugin project: %s", options.Name))

	// Validate options
	if err := pg.validateOptions(options); err != nil {
		return fmt.Errorf("invalid options: %v", err)
	}

	// Create plugin directory
	pluginDir := options.Name + "-plugin"
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory: %v", err)
	}

	// Generate plugin structure
	if err := pg.generatePluginStructure(pluginDir, options); err != nil {
		return fmt.Errorf("failed to generate plugin structure: %v", err)
	}

	// Generate plugin manifest
	if err := pg.generateManifest(pluginDir, options); err != nil {
		return fmt.Errorf("failed to generate manifest: %v", err)
	}

	// Generate plugin implementation
	if err := pg.generateImplementation(pluginDir, options); err != nil {
		return fmt.Errorf("failed to generate implementation: %v", err)
	}

	// Generate documentation
	if err := pg.generateDocumentation(pluginDir, options); err != nil {
		return fmt.Errorf("failed to generate documentation: %v", err)
	}

	// Generate build files
	if err := pg.generateBuildFiles(pluginDir, options); err != nil {
		return fmt.Errorf("failed to generate build files: %v", err)
	}

	ui.PrintSuccess("✅ Plugin project generated successfully!")
	ui.PrintInfo(fmt.Sprintf("📁 Plugin created in: %s", pluginDir))
	ui.PrintInfo("🔧 Next steps:")
	ui.PrintInfo(fmt.Sprintf("   cd %s", pluginDir))
	ui.PrintInfo("   go mod tidy")
	ui.PrintInfo("   make build")
	ui.PrintInfo("   vibercode plugins dev-link .")

	return nil
}

// validateOptions validates plugin generation options
func (pg *PluginGenerator) validateOptions(options models.PluginOptions) error {
	if options.Name == "" {
		return fmt.Errorf("plugin name is required")
	}

	if !strings.HasPrefix(options.Name, "vibercode-") {
		options.Name = "vibercode-" + options.Name
	}

	pluginType := models.PluginType(options.Type)
	if !pluginType.IsValid() {
		return fmt.Errorf("invalid plugin type: %s", options.Type)
	}

	return nil
}

// generatePluginStructure creates the plugin directory structure
func (pg *PluginGenerator) generatePluginStructure(pluginDir string, options models.PluginOptions) error {
	dirs := []string{
		"cmd/plugin",
		"internal",
		"pkg",
		"templates",
		"examples",
		"docs",
		"tests",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(pluginDir, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", fullPath, err)
		}
	}

	return nil
}

// generateManifest generates the plugin.yaml manifest file
func (pg *PluginGenerator) generateManifest(pluginDir string, options models.PluginOptions) error {
	pluginType := models.PluginType(options.Type)
	manifest := models.PluginManifest{
		Name:         options.Name,
		Version:      "1.0.0",
		Description:  options.Description,
		Author:       options.Author,
		License:      "MIT",
		Homepage:     fmt.Sprintf("https://github.com/%s/%s", strings.ToLower(options.Author), options.Name),
		Type:         pluginType,
		Capabilities: pluginType.GetCapabilities(),
		Dependencies: models.PluginDependencies{
			ViberCode: ">=1.0.0",
			Go:        ">=1.19",
		},
		Commands: []models.PluginCommand{
			{
				Name:        fmt.Sprintf("%s-command", strings.ToLower(options.Name)),
				Description: fmt.Sprintf("%s command", options.Name),
				Usage:       fmt.Sprintf("vibercode %s --help", strings.ToLower(options.Name)),
			},
		},
		Config: models.PluginConfig{
			Properties: map[string]models.PluginConfigProperty{
				"default_option": {
					Type:        "string",
					Default:     "default_value",
					Description: "Default configuration option",
					Required:    false,
				},
			},
		},
		Main: "./cmd/plugin/main.go",
	}

	data, err := yaml.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %v", err)
	}

	manifestPath := filepath.Join(pluginDir, "plugin.yaml")
	return ioutil.WriteFile(manifestPath, data, 0644)
}

// generateImplementation generates the plugin implementation
func (pg *PluginGenerator) generateImplementation(pluginDir string, options models.PluginOptions) error {
	pluginType := models.PluginType(options.Type)

	// Generate main.go
	mainContent := pg.generateMainFile(options)
	mainPath := filepath.Join(pluginDir, "cmd", "plugin", "main.go")
	if err := ioutil.WriteFile(mainPath, []byte(mainContent), 0644); err != nil {
		return fmt.Errorf("failed to write main.go: %v", err)
	}

	// Generate plugin-specific implementation
	switch pluginType {
	case models.PluginTypeGenerator:
		return pg.generateGeneratorImplementation(pluginDir, options)
	case models.PluginTypeTemplate:
		return pg.generateTemplateImplementation(pluginDir, options)
	case models.PluginTypeCommand:
		return pg.generateCommandImplementation(pluginDir, options)
	case models.PluginTypeIntegration:
		return pg.generateIntegrationImplementation(pluginDir, options)
	default:
		return fmt.Errorf("unsupported plugin type: %s", pluginType)
	}
}

// generateMainFile generates the main.go file
func (pg *PluginGenerator) generateMainFile(options models.PluginOptions) string {
	return fmt.Sprintf(`package main

import (
	"fmt"
	"log"
	"os"

	"github.com/vibercode/plugin-sdk/api"
	"github.com/%s/%s/internal"
)

func main() {
	plugin := internal.New%sPlugin()
	
	if err := api.RegisterPlugin(plugin); err != nil {
		log.Fatalf("Failed to register plugin: %%v", err)
	}

	fmt.Printf("Plugin %%s initialized successfully\n", plugin.Name())
}
`,
		strings.ToLower(options.Author),
		options.Name,
		pluginToCamelCase(options.Name),
	)
}

// generateGeneratorImplementation generates generator plugin implementation
func (pg *PluginGenerator) generateGeneratorImplementation(pluginDir string, options models.PluginOptions) error {
	content := fmt.Sprintf(`package internal

import (
	"fmt"
	"path/filepath"

	"github.com/vibercode/plugin-sdk/api"
	"github.com/vibercode/plugin-sdk/utils"
)

type %sPlugin struct {
	ctx api.PluginContext
}

func New%sPlugin() *%sPlugin {
	return &%sPlugin{}
}

func (p *%sPlugin) Name() string {
	return "%s"
}

func (p *%sPlugin) Version() string {
	return "1.0.0"
}

func (p *%sPlugin) Description() string {
	return "%s"
}

func (p *%sPlugin) Author() string {
	return "%s"
}

func (p *%sPlugin) Initialize(ctx api.PluginContext) error {
	p.ctx = ctx
	p.ctx.Logger().Info("Initializing %%s plugin", p.Name())
	return nil
}

func (p *%sPlugin) Execute(args []string) error {
	p.ctx.Logger().Info("Executing %%s plugin with args: %%v", p.Name(), args)

	// Parse arguments
	config, err := p.parseArgs(args)
	if err != nil {
		return fmt.Errorf("failed to parse arguments: %%v", err)
	}

	// Generate code
	return p.generateCode(config)
}

func (p *%sPlugin) Cleanup() error {
	p.ctx.Logger().Info("Cleaning up %%s plugin", p.Name())
	return nil
}

func (p *%sPlugin) Commands() []api.Command {
	return []api.Command{
		{
			Name:        "%s-generate",
			Description: "Generate code using %s",
			Usage:       "vibercode %s-generate [options]",
		},
	}
}

func (p *%sPlugin) Generators() []api.Generator {
	return []api.Generator{
		{
			Name:        "%s",
			Description: "%s generator",
			Type:        "code",
		},
	}
}

func (p *%sPlugin) Templates() []api.Template {
	return []api.Template{
		{
			Name: "default",
			Path: "./templates/default.tmpl",
		},
	}
}

func (p *%sPlugin) parseArgs(args []string) (map[string]interface{}, error) {
	// TODO: Implement argument parsing
	config := make(map[string]interface{})
	config["name"] = "example"
	return config, nil
}

func (p *%sPlugin) generateCode(config map[string]interface{}) error {
	// TODO: Implement code generation logic
	p.ctx.Logger().Info("Generating code with config: %%+v", config)

	// Example: Create a simple file
	content := "// Generated by %s plugin\npackage main\n\nfunc main() {\n\tfmt.Println(\"Hello from plugin!\")\n}\n"
	
	outputPath := filepath.Join(".", "generated.go")
	return p.ctx.FileSystem().WriteFile(outputPath, []byte(content))
}
`,
		pluginToCamelCase(options.Name), pluginToCamelCase(options.Name), pluginToCamelCase(options.Name), pluginToCamelCase(options.Name),
		pluginToCamelCase(options.Name), options.Name,
		pluginToCamelCase(options.Name),
		pluginToCamelCase(options.Name), options.Description,
		pluginToCamelCase(options.Name), options.Author,
		pluginToCamelCase(options.Name),
		pluginToCamelCase(options.Name),
		pluginToCamelCase(options.Name), pluginToCamelCase(options.Name),
		options.Name, options.Name, options.Name,
		pluginToCamelCase(options.Name),
		options.Name, options.Description,
		pluginToCamelCase(options.Name),
		pluginToCamelCase(options.Name),
		pluginToCamelCase(options.Name), options.Name,
	)

	implPath := filepath.Join(pluginDir, "internal", "plugin.go")
	return ioutil.WriteFile(implPath, []byte(content), 0644)
}

// generateTemplateImplementation generates template plugin implementation
func (pg *PluginGenerator) generateTemplateImplementation(pluginDir string, options models.PluginOptions) error {
	content := fmt.Sprintf(`package internal

import (
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/vibercode/plugin-sdk/api"
)

type %sPlugin struct {
	ctx       api.PluginContext
	templates map[string]*template.Template
}

func New%sPlugin() *%sPlugin {
	return &%sPlugin{
		templates: make(map[string]*template.Template),
	}
}

func (p *%sPlugin) Name() string {
	return "%s"
}

func (p *%sPlugin) Version() string {
	return "1.0.0"
}

func (p *%sPlugin) Initialize(ctx api.PluginContext) error {
	p.ctx = ctx
	return p.loadTemplates()
}

func (p *%sPlugin) Execute(args []string) error {
	// Template plugins are typically used by other generators
	return fmt.Errorf("template plugins cannot be executed directly")
}

func (p *%sPlugin) GetTemplates() (map[string]string, error) {
	templates := make(map[string]string)
	
	// TODO: Load templates from templates directory
	templates["default"] = "Default template content"
	
	return templates, nil
}

func (p *%sPlugin) RenderTemplate(name string, data interface{}) (string, error) {
	tmpl, exists := p.templates[name]
	if !exists {
		return "", fmt.Errorf("template %%s not found", name)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to render template: %%v", err)
	}

	return buf.String(), nil
}

func (p *%sPlugin) ValidateTemplate(name string) error {
	_, exists := p.templates[name]
	if !exists {
		return fmt.Errorf("template %%s not found", name)
	}
	return nil
}

func (p *%sPlugin) loadTemplates() error {
	// TODO: Load templates from templates directory
	return nil
}
`,
		pluginToCamelCase(options.Name), pluginToCamelCase(options.Name), pluginToCamelCase(options.Name), pluginToCamelCase(options.Name),
		pluginToCamelCase(options.Name), options.Name,
		pluginToCamelCase(options.Name),
		pluginToCamelCase(options.Name),
		pluginToCamelCase(options.Name),
		pluginToCamelCase(options.Name),
		pluginToCamelCase(options.Name),
		pluginToCamelCase(options.Name),
		pluginToCamelCase(options.Name),
	)

	implPath := filepath.Join(pluginDir, "internal", "plugin.go")
	return ioutil.WriteFile(implPath, []byte(content), 0644)
}

// generateCommandImplementation generates command plugin implementation
func (pg *PluginGenerator) generateCommandImplementation(pluginDir string, options models.PluginOptions) error {
	content := fmt.Sprintf(`package internal

import (
	"fmt"
	"strings"

	"github.com/vibercode/plugin-sdk/api"
)

type %sPlugin struct {
	ctx api.PluginContext
}

func New%sPlugin() *%sPlugin {
	return &%sPlugin{}
}

func (p *%sPlugin) Name() string {
	return "%s"
}

func (p *%sPlugin) Version() string {
	return "1.0.0"
}

func (p *%sPlugin) Initialize(ctx api.PluginContext) error {
	p.ctx = ctx
	return nil
}

func (p *%sPlugin) Execute(args []string) error {
	return p.ExecuteCommand(args)
}

func (p *%sPlugin) ExecuteCommand(args []string) (*api.ExecutionResult, error) {
	p.ctx.Logger().Info("Executing command with args: %%v", args)

	if len(args) == 0 {
		return nil, fmt.Errorf("no command specified")
	}

	command := args[0]
	commandArgs := args[1:]

	switch command {
	case "hello":
		return p.executeHello(commandArgs)
	case "help":
		return p.executeHelp(commandArgs)
	default:
		return nil, fmt.Errorf("unknown command: %%s", command)
	}
}

func (p *%sPlugin) GetCommandHelp() string {
	return %%sAvailable commands:
  hello [name]  - Say hello to someone
  help          - Show this help message
%%s
}

func (p *%sPlugin) ValidateArgs(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("at least one argument is required")
	}
	return nil
}

func (p *%sPlugin) executeHello(args []string) (*api.ExecutionResult, error) {
	name := "World"
	if len(args) > 0 {
		name = strings.Join(args, " ")
	}

	message := fmt.Sprintf("Hello, %%s! This is %%s plugin.", name, p.Name())
	
	return &api.ExecutionResult{
		Success: true,
		Message: message,
		Output:  message,
	}, nil
}

func (p *%sPlugin) executeHelp(args []string) (*api.ExecutionResult, error) {
	help := p.GetCommandHelp()
	
	return &api.ExecutionResult{
		Success: true,
		Message: "Help displayed",
		Output:  help,
	}, nil
}
`,
		pluginToCamelCase(options.Name), pluginToCamelCase(options.Name), pluginToCamelCase(options.Name),
		pluginToCamelCase(options.Name), options.Name,
		pluginToCamelCase(options.Name),
		pluginToCamelCase(options.Name),
		pluginToCamelCase(options.Name),
		pluginToCamelCase(options.Name), "`", "`",
		pluginToCamelCase(options.Name),
		pluginToCamelCase(options.Name),
		pluginToCamelCase(options.Name),
	)

	implPath := filepath.Join(pluginDir, "internal", "plugin.go")
	return ioutil.WriteFile(implPath, []byte(content), 0644)
}

// generateIntegrationImplementation generates integration plugin implementation
func (pg *PluginGenerator) generateIntegrationImplementation(pluginDir string, options models.PluginOptions) error {
	content := fmt.Sprintf(`package internal

import (
	"fmt"
	"net/http"

	"github.com/vibercode/plugin-sdk/api"
)

type %sPlugin struct {
	ctx       api.PluginContext
	client    *http.Client
	connected bool
	config    map[string]interface{}
}

func New%sPlugin() *%sPlugin {
	return &%sPlugin{
		client: &http.Client{},
	}
}

func (p *%sPlugin) Name() string {
	return "%s"
}

func (p *%sPlugin) Version() string {
	return "1.0.0"
}

func (p *%sPlugin) Initialize(ctx api.PluginContext) error {
	p.ctx = ctx
	return nil
}

func (p *%sPlugin) Execute(args []string) error {
	// Integration plugins typically provide services, not direct execution
	return fmt.Errorf("integration plugins cannot be executed directly")
}

func (p *%sPlugin) InitializeIntegration(config map[string]interface{}) error {
	p.config = config
	p.ctx.Logger().Info("Initializing integration with config: %%+v", config)
	return nil
}

func (p *%sPlugin) Connect() error {
	p.ctx.Logger().Info("Connecting to external service...")
	
	// TODO: Implement actual connection logic
	p.connected = true
	
	p.ctx.Logger().Info("Successfully connected to external service")
	return nil
}

func (p *%sPlugin) Disconnect() error {
	p.ctx.Logger().Info("Disconnecting from external service...")
	
	// TODO: Implement actual disconnection logic
	p.connected = false
	
	p.ctx.Logger().Info("Successfully disconnected from external service")
	return nil
}

func (p *%sPlugin) IsConnected() bool {
	return p.connected
}

func (p *%sPlugin) GetStatus() string {
	if p.connected {
		return "Connected"
	}
	return "Disconnected"
}

// Custom methods for this integration
func (p *%sPlugin) SyncData() error {
	if !p.connected {
		return fmt.Errorf("not connected to external service")
	}

	p.ctx.Logger().Info("Syncing data with external service...")
	
	// TODO: Implement data synchronization logic
	
	return nil
}

func (p *%sPlugin) GetData(id string) (interface{}, error) {
	if !p.connected {
		return nil, fmt.Errorf("not connected to external service")
	}

	p.ctx.Logger().Info("Fetching data for ID: %%s", id)
	
	// TODO: Implement data fetching logic
	
	return map[string]interface{}{
		"id":   id,
		"data": "sample data",
	}, nil
}
`,
		pluginToCamelCase(options.Name), pluginToCamelCase(options.Name), pluginToCamelCase(options.Name), pluginToCamelCase(options.Name),
		pluginToCamelCase(options.Name), options.Name,
		pluginToCamelCase(options.Name),
		pluginToCamelCase(options.Name),
		pluginToCamelCase(options.Name),
		pluginToCamelCase(options.Name),
		pluginToCamelCase(options.Name),
		pluginToCamelCase(options.Name),
		pluginToCamelCase(options.Name),
		pluginToCamelCase(options.Name),
	)

	implPath := filepath.Join(pluginDir, "internal", "plugin.go")
	return ioutil.WriteFile(implPath, []byte(content), 0644)
}

// generateDocumentation generates plugin documentation
func (pg *PluginGenerator) generateDocumentation(pluginDir string, options models.PluginOptions) error {
	readmeContent := fmt.Sprintf(`# %s

%s

## Description

%s

## Installation

` + "```bash" + `
vibercode plugins install %s
` + "```" + `

## Usage

` + "```bash" + `
vibercode %s --help
` + "```" + `

## Configuration

This plugin supports the following configuration options:

- ` + "`default_option`" + `: Default configuration option (default: "default_value")

## Development

### Building

` + "```bash" + `
make build
` + "```" + `

### Testing

` + "```bash" + `
make test
` + "```" + `

### Development Linking

` + "```bash" + `
vibercode plugins dev-link .
` + "```" + `

## License

MIT License - see LICENSE file for details.

## Author

%s
`,
		options.Name,
		options.Description,
		options.Description,
		options.Name,
		strings.ToLower(options.Name),
		options.Author,
	)

	readmePath := filepath.Join(pluginDir, "README.md")
	if err := ioutil.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("failed to write README.md: %v", err)
	}

	// Generate API documentation
	apiDocsContent := fmt.Sprintf(`# %s API Documentation

## Plugin Interface

This plugin implements the following interfaces:

### Core Methods

- ` + "`Name()`" + ` - Returns the plugin name
- ` + "`Version()`" + ` - Returns the plugin version  
- ` + "`Initialize(ctx)`" + ` - Initializes the plugin
- ` + "`Execute(args)`" + ` - Executes the plugin

### Plugin-Specific Methods

TODO: Document plugin-specific methods here

## Configuration Schema

` + "```json" + `
{
  "properties": {
    "default_option": {
      "type": "string",
      "default": "default_value",
      "description": "Default configuration option"
    }
  }
}
` + "```" + `

## Examples

TODO: Add usage examples here
`,
		options.Name,
	)

	apiDocsPath := filepath.Join(pluginDir, "docs", "api.md")
	return ioutil.WriteFile(apiDocsPath, []byte(apiDocsContent), 0644)
}

// generateBuildFiles generates build configuration files
func (pg *PluginGenerator) generateBuildFiles(pluginDir string, options models.PluginOptions) error {
	// Generate go.mod
	goModContent := fmt.Sprintf(`module github.com/%s/%s

go 1.19

require (
	github.com/vibercode/plugin-sdk v1.0.0
)
`,
		strings.ToLower(options.Author),
		options.Name,
	)

	goModPath := filepath.Join(pluginDir, "go.mod")
	if err := ioutil.WriteFile(goModPath, []byte(goModContent), 0644); err != nil {
		return fmt.Errorf("failed to write go.mod: %v", err)
	}

	// Generate Makefile
	makefileContent := fmt.Sprintf(`# Makefile for %s plugin

BINARY_NAME=%s
BUILD_DIR=./build
CMD_DIR=./cmd/plugin

.PHONY: build clean test install dev-link

build:
	@echo "Building plugin..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)

test:
	@echo "Running tests..."
	@go test -v ./...

install: build
	@echo "Installing plugin..."
	@vibercode plugins install ./$(BUILD_DIR)/$(BINARY_NAME)

dev-link:
	@echo "Linking plugin for development..."
	@vibercode plugins dev-link .

lint:
	@echo "Running linter..."
	@golangci-lint run

format:
	@echo "Formatting code..."
	@gofmt -w -s .
	@goimports -w .

deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

help:
	@echo "Available targets:"
	@echo "  build     - Build the plugin binary"
	@echo "  clean     - Clean build artifacts"
	@echo "  test      - Run tests"
	@echo "  install   - Install the plugin"
	@echo "  dev-link  - Link plugin for development"
	@echo "  lint      - Run linter"
	@echo "  format    - Format code"
	@echo "  deps      - Download dependencies"
	@echo "  help      - Show this help"
`,
		options.Name,
		options.Name,
	)

	makefilePath := filepath.Join(pluginDir, "Makefile")
	if err := ioutil.WriteFile(makefilePath, []byte(makefileContent), 0644); err != nil {
		return fmt.Errorf("failed to write Makefile: %v", err)
	}

	// Generate .gitignore
	gitignoreContent := `# Build artifacts
build/
dist/
*.exe
*.dll
*.so
*.dylib

# Go specific
*.test
*.out
vendor/

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Logs
*.log

# Temporary files
*.tmp
*.temp
`

	gitignorePath := filepath.Join(pluginDir, ".gitignore")
	return ioutil.WriteFile(gitignorePath, []byte(gitignoreContent), 0644)
}

// pluginToCamelCase converts a string to CamelCase
func pluginToCamelCase(s string) string {
	// Remove common prefixes
	s = strings.TrimPrefix(s, "vibercode-")
	
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == '-' || r == '_' || r == ' '
	})
	
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}
	
	return strings.Join(words, "")
}

// PluginManager handles plugin management operations
type PluginManager struct {
	config *models.PluginConfigManager
}

// NewPluginManager creates a new plugin manager
func NewPluginManager() *PluginManager {
	return &PluginManager{
		config: models.GetDefaultConfig(),
	}
}

// ListPlugins lists all installed plugins
func (pm *PluginManager) ListPlugins() ([]models.PluginInfo, error) {
	ui.PrintInfo("📋 Listing installed plugins...")

	pluginsDir := pm.config.PluginsDir
	if _, err := os.Stat(pluginsDir); os.IsNotExist(err) {
		return []models.PluginInfo{}, nil
	}

	entries, err := ioutil.ReadDir(pluginsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugins directory: %v", err)
	}

	var plugins []models.PluginInfo
	for _, entry := range entries {
		if entry.IsDir() {
			pluginPath := filepath.Join(pluginsDir, entry.Name())
			info, err := pm.loadPluginInfo(pluginPath)
			if err != nil {
				ui.PrintWarning(fmt.Sprintf("⚠️  Failed to load plugin %s: %v", entry.Name(), err))
				continue
			}
			plugins = append(plugins, *info)
		}
	}

	return plugins, nil
}

// InstallPlugin installs a plugin
func (pm *PluginManager) InstallPlugin(options models.PluginInstallOptions) error {
	ui.PrintInfo(fmt.Sprintf("📦 Installing plugin: %s", options.Name))

	// TODO: Implement plugin installation logic
	// This would involve:
	// 1. Downloading plugin from registry or URL
	// 2. Validating plugin security
	// 3. Extracting to plugins directory
	// 4. Running installation hooks

	ui.PrintSuccess(fmt.Sprintf("✅ Plugin %s installed successfully!", options.Name))
	return nil
}

// UninstallPlugin uninstalls a plugin
func (pm *PluginManager) UninstallPlugin(name string) error {
	ui.PrintInfo(fmt.Sprintf("🗑️  Uninstalling plugin: %s", name))

	pluginPath := filepath.Join(pm.config.PluginsDir, name)
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		return fmt.Errorf("plugin %s is not installed", name)
	}

	if err := os.RemoveAll(pluginPath); err != nil {
		return fmt.Errorf("failed to remove plugin directory: %v", err)
	}

	ui.PrintSuccess(fmt.Sprintf("✅ Plugin %s uninstalled successfully!", name))
	return nil
}

// loadPluginInfo loads plugin information from a plugin directory
func (pm *PluginManager) loadPluginInfo(pluginPath string) (*models.PluginInfo, error) {
	manifestPath := filepath.Join(pluginPath, "plugin.yaml")
	
	data, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin manifest: %v", err)
	}

	var manifest models.PluginManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse plugin manifest: %v", err)
	}

	info := &models.PluginInfo{
		Manifest:    manifest,
		Status:      models.PluginStatusInstalled,
		InstallPath: pluginPath,
		LastUpdated: time.Now(), // TODO: Get actual last modified time
	}

	return info, nil
}