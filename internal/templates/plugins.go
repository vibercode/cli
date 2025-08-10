package templates

// PluginManifestTemplate generates plugin.yaml file
const PluginManifestTemplate = `name: "{{.Name}}"
version: "{{.Version}}"
description: "{{.Description}}"
author: "{{.Author}}"
license: "{{.License}}"
homepage: "{{.Homepage}}"

# Plugin type and capabilities
type: "{{.Type}}"
capabilities:
{{- range .Capabilities}}
  - "{{.}}"
{{- end}}

# Dependencies
dependencies:
  vibercode: "{{.Dependencies.ViberCode}}"
  go: "{{.Dependencies.Go}}"

# Commands provided by this plugin
commands:
{{- range .Commands}}
  - name: "{{.Name}}"
    description: "{{.Description}}"
    usage: "{{.Usage}}"
{{- end}}

# Configuration schema
config:
  properties:
{{- range $key, $prop := .Config.Properties}}
    {{$key}}:
      type: "{{$prop.Type}}"
      default: {{$prop.Default}}
      description: "{{$prop.Description}}"
      required: {{$prop.Required}}
{{- end}}

# Entry point
main: "{{.Main}}"
`

// PluginMainTemplate generates main.go file for plugins
const PluginMainTemplate = `package main

import (
	"fmt"
	"log"
	"os"

	"github.com/vibercode/plugin-sdk/api"
	"github.com/{{.Author}}/{{.Name}}/internal"
)

func main() {
	plugin := internal.New{{.CamelName}}Plugin()
	
	if err := api.RegisterPlugin(plugin); err != nil {
		log.Fatalf("Failed to register plugin: %v", err)
	}

	fmt.Printf("Plugin %s initialized successfully\n", plugin.Name())
}
`

// GeneratorPluginTemplate generates implementation for generator plugins
const GeneratorPluginTemplate = `package internal

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/vibercode/plugin-sdk/api"
	"github.com/vibercode/plugin-sdk/utils"
)

type {{.CamelName}}Plugin struct {
	ctx api.PluginContext
}

func New{{.CamelName}}Plugin() *{{.CamelName}}Plugin {
	return &{{.CamelName}}Plugin{}
}

func (p *{{.CamelName}}Plugin) Name() string {
	return "{{.Name}}"
}

func (p *{{.CamelName}}Plugin) Version() string {
	return "{{.Version}}"
}

func (p *{{.CamelName}}Plugin) Description() string {
	return "{{.Description}}"
}

func (p *{{.CamelName}}Plugin) Author() string {
	return "{{.Author}}"
}

func (p *{{.CamelName}}Plugin) Initialize(ctx api.PluginContext) error {
	p.ctx = ctx
	p.ctx.Logger().Info("Initializing %s plugin", p.Name())
	return nil
}

func (p *{{.CamelName}}Plugin) Execute(args []string) error {
	p.ctx.Logger().Info("Executing %s plugin with args: %v", p.Name(), args)

	// Parse arguments
	config, err := p.parseArgs(args)
	if err != nil {
		return fmt.Errorf("failed to parse arguments: %v", err)
	}

	// Generate code
	return p.generateCode(config)
}

func (p *{{.CamelName}}Plugin) Cleanup() error {
	p.ctx.Logger().Info("Cleaning up %s plugin", p.Name())
	return nil
}

func (p *{{.CamelName}}Plugin) Commands() []api.Command {
	return []api.Command{
		{
			Name:        "{{.KebabName}}-generate",
			Description: "Generate code using {{.Name}}",
			Usage:       "vibercode {{.KebabName}}-generate [options]",
			Aliases:     []string{"{{.ShortName}}"},
			Flags: []api.Flag{
				{
					Name:        "name",
					ShortName:   "n",
					Description: "Name of the component to generate",
					Type:        "string",
					Required:    true,
				},
				{
					Name:        "output",
					ShortName:   "o",
					Description: "Output directory",
					Type:        "string",
					Required:    false,
					Default:     ".",
				},
			},
		},
	}
}

func (p *{{.CamelName}}Plugin) Generators() []api.Generator {
	return []api.Generator{
		{
			Name:        "{{.KebabName}}",
			Description: "{{.Description}} generator",
			Type:        "code",
			Extensions:  []string{".go", ".yaml", ".md"},
			Templates:   []string{"default", "advanced", "minimal"},
		},
	}
}

func (p *{{.CamelName}}Plugin) Templates() []api.Template {
	return []api.Template{
		{
			Name:        "default",
			Path:        "./templates/default.tmpl",
			Description: "Default template for {{.Name}}",
			Variables: map[string]string{
				"name":        "Component name",
				"description": "Component description",
				"author":      "Component author",
			},
			Helpers: []string{"toCamel", "toSnake", "toKebab"},
		},
		{
			Name:        "advanced",
			Path:        "./templates/advanced.tmpl",
			Description: "Advanced template with additional features",
			Variables: map[string]string{
				"name":     "Component name",
				"features": "Comma-separated list of features",
			},
		},
	}
}

// Generate implements the GeneratorPlugin interface
func (p *{{.CamelName}}Plugin) Generate(options map[string]interface{}) (*api.ExecutionResult, error) {
	p.ctx.Logger().Info("Generating code with options: %+v", options)

	result := &api.ExecutionResult{
		Success: true,
		Message: "Generation completed successfully",
	}

	// Extract options
	name, ok := options["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("name is required")
	}

	outputDir, ok := options["output"].(string)
	if !ok {
		outputDir = "."
	}

	template, ok := options["template"].(string)
	if !ok {
		template = "default"
	}

	// Generate files
	files, err := p.generateFiles(name, outputDir, template, options)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("Generation failed: %v", err)
		return result, err
	}

	result.FilesCreated = files
	result.Output = fmt.Sprintf("Generated %d files for %s", len(files), name)

	return result, nil
}

// GetTemplates implements the GeneratorPlugin interface
func (p *{{.CamelName}}Plugin) GetTemplates() ([]string, error) {
	return []string{"default", "advanced", "minimal"}, nil
}

// ValidateOptions implements the GeneratorPlugin interface
func (p *{{.CamelName}}Plugin) ValidateOptions(options map[string]interface{}) error {
	// Validate required options
	if _, ok := options["name"]; !ok {
		return fmt.Errorf("name option is required")
	}

	// Validate name format
	name := options["name"].(string)
	if !utils.IsValidIdentifier(name) {
		return fmt.Errorf("invalid name format: %s", name)
	}

	return nil
}

// GetSchema implements the GeneratorPlugin interface
func (p *{{.CamelName}}Plugin) GetSchema() (map[string]interface{}, error) {
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name": map[string]interface{}{
				"type":        "string",
				"description": "Name of the component to generate",
				"pattern":     "^[a-zA-Z][a-zA-Z0-9_]*$",
			},
			"output": map[string]interface{}{
				"type":        "string",
				"description": "Output directory",
				"default":     ".",
			},
			"template": map[string]interface{}{
				"type":        "string",
				"description": "Template to use",
				"enum":        []string{"default", "advanced", "minimal"},
				"default":     "default",
			},
		},
		"required": []string{"name"},
	}

	return schema, nil
}

func (p *{{.CamelName}}Plugin) parseArgs(args []string) (map[string]interface{}, error) {
	config := make(map[string]interface{})
	
	// Simple argument parsing (in real implementation, use a proper flag parser)
	for i, arg := range args {
		switch arg {
		case "--name", "-n":
			if i+1 < len(args) {
				config["name"] = args[i+1]
			}
		case "--output", "-o":
			if i+1 < len(args) {
				config["output"] = args[i+1]
			}
		case "--template", "-t":
			if i+1 < len(args) {
				config["template"] = args[i+1]
			}
		}
	}

	return config, nil
}

func (p *{{.CamelName}}Plugin) generateFiles(name, outputDir, template string, options map[string]interface{}) ([]string, error) {
	var generatedFiles []string

	// Prepare template data
	data := map[string]interface{}{
		"Name":        name,
		"CamelName":   utils.ToCamel(name),
		"SnakeName":   utils.ToSnake(name),
		"KebabName":   utils.ToKebab(name),
		"Author":      p.ctx.Config()["author"],
		"Timestamp":   utils.Now(),
		"Options":     options,
	}

	// Load and render template
	templatePath := filepath.Join("templates", template+".tmpl")
	content, err := p.ctx.TemplateEngine().RenderFile(templatePath, data)
	if err != nil {
		return nil, fmt.Errorf("failed to render template: %v", err)
	}

	// Write output file
	outputFile := filepath.Join(outputDir, strings.ToLower(name)+".go")
	if err := p.ctx.FileSystem().WriteFile(outputFile, []byte(content)); err != nil {
		return nil, fmt.Errorf("failed to write file: %v", err)
	}

	generatedFiles = append(generatedFiles, outputFile)

	// Generate additional files based on template
	switch template {
	case "advanced":
		testFile := filepath.Join(outputDir, strings.ToLower(name)+"_test.go")
		testContent := p.generateTestFile(data)
		if err := p.ctx.FileSystem().WriteFile(testFile, []byte(testContent)); err != nil {
			return nil, fmt.Errorf("failed to write test file: %v", err)
		}
		generatedFiles = append(generatedFiles, testFile)

	case "minimal":
		// Minimal template generates only the basic file
		// No additional files needed
	}

	return generatedFiles, nil
}

func (p *{{.CamelName}}Plugin) generateTestFile(data map[string]interface{}) string {
	// Simple test file template
	return fmt.Sprintf(` + "`" + `package main

import (
	"testing"
)

func Test%s(t *testing.T) {
	// TODO: Implement test for %s
	t.Skip("Test not implemented yet")
}
` + "`" + `, data["CamelName"], data["Name"])
}
`

// TemplatePluginTemplate generates implementation for template plugins
const TemplatePluginTemplate = `package internal

import (
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/vibercode/plugin-sdk/api"
)

type {{.CamelName}}Plugin struct {
	ctx       api.PluginContext
	templates map[string]*template.Template
}

func New{{.CamelName}}Plugin() *{{.CamelName}}Plugin {
	return &{{.CamelName}}Plugin{
		templates: make(map[string]*template.Template),
	}
}

func (p *{{.CamelName}}Plugin) Name() string {
	return "{{.Name}}"
}

func (p *{{.CamelName}}Plugin) Version() string {
	return "{{.Version}}"
}

func (p *{{.CamelName}}Plugin) Description() string {
	return "{{.Description}}"
}

func (p *{{.CamelName}}Plugin) Author() string {
	return "{{.Author}}"
}

func (p *{{.CamelName}}Plugin) Initialize(ctx api.PluginContext) error {
	p.ctx = ctx
	return p.loadTemplates()
}

func (p *{{.CamelName}}Plugin) Execute(args []string) error {
	// Template plugins are typically used by other generators
	return fmt.Errorf("template plugins cannot be executed directly")
}

func (p *{{.CamelName}}Plugin) Cleanup() error {
	return nil
}

func (p *{{.CamelName}}Plugin) Commands() []api.Command {
	return []api.Command{
		{
			Name:        "{{.KebabName}}-list",
			Description: "List available templates",
			Usage:       "vibercode {{.KebabName}}-list",
		},
	}
}

func (p *{{.CamelName}}Plugin) Generators() []api.Generator {
	return []api.Generator{}
}

func (p *{{.CamelName}}Plugin) Templates() []api.Template {
	return []api.Template{
		{
			Name:        "main",
			Path:        "./templates/main.tmpl",
			Description: "Main template for {{.Description}}",
		},
		{
			Name:        "config",
			Path:        "./templates/config.tmpl",
			Description: "Configuration template",
		},
		{
			Name:        "readme",
			Path:        "./templates/readme.tmpl",
			Description: "README template",
		},
	}
}

// GetTemplates implements the TemplatePlugin interface
func (p *{{.CamelName}}Plugin) GetTemplates() (map[string]string, error) {
	templates := make(map[string]string)
	
	templateFiles := []string{"main", "config", "readme"}
	
	for _, name := range templateFiles {
		templatePath := filepath.Join("templates", name+".tmpl")
		content, err := p.ctx.FileSystem().ReadFile(templatePath)
		if err != nil {
			p.ctx.Logger().Warning("Failed to load template %s: %v", name, err)
			continue
		}
		templates[name] = string(content)
	}
	
	return templates, nil
}

// RenderTemplate implements the TemplatePlugin interface
func (p *{{.CamelName}}Plugin) RenderTemplate(name string, data interface{}) (string, error) {
	tmpl, exists := p.templates[name]
	if !exists {
		return "", fmt.Errorf("template %s not found", name)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to render template: %v", err)
	}

	return buf.String(), nil
}

// ValidateTemplate implements the TemplatePlugin interface
func (p *{{.CamelName}}Plugin) ValidateTemplate(name string) error {
	_, exists := p.templates[name]
	if !exists {
		return fmt.Errorf("template %s not found", name)
	}
	return nil
}

// GetTemplateVars implements the TemplatePlugin interface
func (p *{{.CamelName}}Plugin) GetTemplateVars(name string) ([]string, error) {
	tmpl, exists := p.templates[name]
	if !exists {
		return nil, fmt.Errorf("template %s not found", name)
	}

	// Extract template variables (simplified implementation)
	// In a real implementation, you would parse the template to find variables
	vars := []string{"Name", "Description", "Author", "Version"}
	
	return vars, nil
}

func (p *{{.CamelName}}Plugin) loadTemplates() error {
	templateDir := "./templates"
	templateFiles := []string{"main.tmpl", "config.tmpl", "readme.tmpl"}

	for _, file := range templateFiles {
		templatePath := filepath.Join(templateDir, file)
		content, err := p.ctx.FileSystem().ReadFile(templatePath)
		if err != nil {
			p.ctx.Logger().Warning("Failed to load template %s: %v", file, err)
			continue
		}

		name := strings.TrimSuffix(file, ".tmpl")
		tmpl, err := template.New(name).Parse(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %v", name, err)
		}

		p.templates[name] = tmpl
	}

	return nil
}
`

// CommandPluginTemplate generates implementation for command plugins
const CommandPluginTemplate = `package internal

import (
	"fmt"
	"strings"

	"github.com/vibercode/plugin-sdk/api"
)

type {{.CamelName}}Plugin struct {
	ctx api.PluginContext
}

func New{{.CamelName}}Plugin() *{{.CamelName}}Plugin {
	return &{{.CamelName}}Plugin{}
}

func (p *{{.CamelName}}Plugin) Name() string {
	return "{{.Name}}"
}

func (p *{{.CamelName}}Plugin) Version() string {
	return "{{.Version}}"
}

func (p *{{.CamelName}}Plugin) Description() string {
	return "{{.Description}}"
}

func (p *{{.CamelName}}Plugin) Author() string {
	return "{{.Author}}"
}

func (p *{{.CamelName}}Plugin) Initialize(ctx api.PluginContext) error {
	p.ctx = ctx
	return nil
}

func (p *{{.CamelName}}Plugin) Execute(args []string) error {
	return p.ExecuteCommand("default", args)
}

func (p *{{.CamelName}}Plugin) Cleanup() error {
	return nil
}

func (p *{{.CamelName}}Plugin) Commands() []api.Command {
	return []api.Command{
		{
			Name:        "{{.KebabName}}-hello",
			Description: "Say hello",
			Usage:       "vibercode {{.KebabName}}-hello [name]",
			Aliases:     []string{"{{.ShortName}}-hi"},
			Flags: []api.Flag{
				{
					Name:        "name",
					ShortName:   "n",
					Description: "Name to greet",
					Type:        "string",
					Required:    false,
					Default:     "World",
				},
			},
		},
		{
			Name:        "{{.KebabName}}-info",
			Description: "Show plugin information",
			Usage:       "vibercode {{.KebabName}}-info",
		},
	}
}

func (p *{{.CamelName}}Plugin) Generators() []api.Generator {
	return []api.Generator{}
}

func (p *{{.CamelName}}Plugin) Templates() []api.Template {
	return []api.Template{}
}

// ExecuteCommand implements the CommandPlugin interface
func (p *{{.CamelName}}Plugin) ExecuteCommand(command string, args []string) (*api.ExecutionResult, error) {
	p.ctx.Logger().Info("Executing command %s with args: %v", command, args)

	switch command {
	case "hello", "{{.KebabName}}-hello":
		return p.executeHello(args)
	case "info", "{{.KebabName}}-info":
		return p.executeInfo(args)
	case "help":
		return p.executeHelp(args)
	default:
		return nil, fmt.Errorf("unknown command: %s", command)
	}
}

// GetCommandHelp implements the CommandPlugin interface
func (p *{{.CamelName}}Plugin) GetCommandHelp(command string) string {
	switch command {
	case "hello":
		return ` + "`" + `Usage: vibercode {{.KebabName}}-hello [options] [name]

Say hello to someone.

Options:
  -n, --name    Name to greet (default: "World")

Examples:
  vibercode {{.KebabName}}-hello
  vibercode {{.KebabName}}-hello John
  vibercode {{.KebabName}}-hello --name Alice` + "`" + `
	case "info":
		return ` + "`" + `Usage: vibercode {{.KebabName}}-info

Show information about the {{.Name}} plugin.` + "`" + `
	default:
		return p.GetGeneralHelp()
	}
}

// ValidateArgs implements the CommandPlugin interface
func (p *{{.CamelName}}Plugin) ValidateArgs(command string, args []string) error {
	switch command {
	case "hello":
		// Hello command accepts optional arguments
		return nil
	case "info":
		// Info command doesn't accept arguments
		if len(args) > 0 {
			return fmt.Errorf("info command doesn't accept arguments")
		}
		return nil
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

// GetCommands implements the CommandPlugin interface
func (p *{{.CamelName}}Plugin) GetCommands() []api.Command {
	return p.Commands()
}

func (p *{{.CamelName}}Plugin) executeHello(args []string) (*api.ExecutionResult, error) {
	name := "World"
	
	// Parse arguments for name
	for i, arg := range args {
		if arg == "--name" || arg == "-n" {
			if i+1 < len(args) {
				name = args[i+1]
			}
		} else if !strings.HasPrefix(arg, "-") && i == 0 {
			// First non-flag argument is the name
			name = arg
		}
	}

	message := fmt.Sprintf("Hello, %s! This is the %s plugin.", name, p.Name())
	
	return &api.ExecutionResult{
		Success: true,
		Message: "Hello command completed",
		Output:  message,
	}, nil
}

func (p *{{.CamelName}}Plugin) executeInfo(args []string) (*api.ExecutionResult, error) {
	info := fmt.Sprintf(` + "`" + `Plugin Information:
  Name: %s
  Version: %s
  Description: %s
  Author: %s
  
Available Commands:
  hello - Say hello to someone
  info  - Show this information
  help  - Show help for commands` + "`" + `,
		p.Name(), p.Version(), p.Description(), p.Author())
	
	return &api.ExecutionResult{
		Success: true,
		Message: "Plugin information displayed",
		Output:  info,
	}, nil
}

func (p *{{.CamelName}}Plugin) executeHelp(args []string) (*api.ExecutionResult, error) {
	if len(args) > 0 {
		// Help for specific command
		commandHelp := p.GetCommandHelp(args[0])
		return &api.ExecutionResult{
			Success: true,
			Message: "Command help displayed",
			Output:  commandHelp,
		}, nil
	}

	// General help
	help := p.GetGeneralHelp()
	return &api.ExecutionResult{
		Success: true,
		Message: "General help displayed",
		Output:  help,
	}, nil
}

func (p *{{.CamelName}}Plugin) GetGeneralHelp() string {
	return fmt.Sprintf(` + "`" + `%s Plugin

%s

Available commands:
  hello [name]  - Say hello to someone
  info          - Show plugin information
  help [cmd]    - Show help for commands

Use 'vibercode %s help <command>' for more information about a command.` + "`" + `,
		p.Name(), p.Description(), strings.ToLower(p.Name()))
}
`

// IntegrationPluginTemplate generates implementation for integration plugins
const IntegrationPluginTemplate = `package internal

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vibercode/plugin-sdk/api"
)

type {{.CamelName}}Plugin struct {
	ctx       api.PluginContext
	client    *http.Client
	connected bool
	config    map[string]interface{}
	baseURL   string
	apiKey    string
}

func New{{.CamelName}}Plugin() *{{.CamelName}}Plugin {
	return &{{.CamelName}}Plugin{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (p *{{.CamelName}}Plugin) Name() string {
	return "{{.Name}}"
}

func (p *{{.CamelName}}Plugin) Version() string {
	return "{{.Version}}"
}

func (p *{{.CamelName}}Plugin) Description() string {
	return "{{.Description}}"
}

func (p *{{.CamelName}}Plugin) Author() string {
	return "{{.Author}}"
}

func (p *{{.CamelName}}Plugin) Initialize(ctx api.PluginContext) error {
	p.ctx = ctx
	return nil
}

func (p *{{.CamelName}}Plugin) Execute(args []string) error {
	// Integration plugins typically provide services, not direct execution
	return fmt.Errorf("integration plugins cannot be executed directly")
}

func (p *{{.CamelName}}Plugin) Cleanup() error {
	if p.connected {
		return p.Disconnect()
	}
	return nil
}

func (p *{{.CamelName}}Plugin) Commands() []api.Command {
	return []api.Command{
		{
			Name:        "{{.KebabName}}-connect",
			Description: "Connect to {{.Name}} service",
			Usage:       "vibercode {{.KebabName}}-connect",
		},
		{
			Name:        "{{.KebabName}}-status",
			Description: "Show connection status",
			Usage:       "vibercode {{.KebabName}}-status",
		},
		{
			Name:        "{{.KebabName}}-sync",
			Description: "Sync data with {{.Name}}",
			Usage:       "vibercode {{.KebabName}}-sync",
		},
	}
}

func (p *{{.CamelName}}Plugin) Generators() []api.Generator {
	return []api.Generator{}
}

func (p *{{.CamelName}}Plugin) Templates() []api.Template {
	return []api.Template{}
}

// Initialize implements the IntegrationPlugin interface
func (p *{{.CamelName}}Plugin) InitializeIntegration(config map[string]interface{}) error {
	p.config = config
	p.ctx.Logger().Info("Initializing {{.Name}} integration with config")

	// Extract configuration
	if baseURL, ok := config["base_url"].(string); ok {
		p.baseURL = baseURL
	} else {
		return fmt.Errorf("base_url is required in configuration")
	}

	if apiKey, ok := config["api_key"].(string); ok {
		p.apiKey = apiKey
	} else {
		return fmt.Errorf("api_key is required in configuration")
	}

	return nil
}

// Connect implements the IntegrationPlugin interface
func (p *{{.CamelName}}Plugin) Connect() error {
	p.ctx.Logger().Info("Connecting to {{.Name}} service at %s", p.baseURL)
	
	// Test connection
	req, err := http.NewRequest("GET", p.baseURL+"/health", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("User-Agent", "ViberCode-Plugin-{{.Name}}/{{.Version}}")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("connection failed with status: %d", resp.StatusCode)
	}

	p.connected = true
	p.ctx.Logger().Info("Successfully connected to {{.Name}} service")
	return nil
}

// Disconnect implements the IntegrationPlugin interface
func (p *{{.CamelName}}Plugin) Disconnect() error {
	p.ctx.Logger().Info("Disconnecting from {{.Name}} service")
	
	// Perform cleanup if needed
	p.connected = false
	
	p.ctx.Logger().Info("Successfully disconnected from {{.Name}} service")
	return nil
}

// IsConnected implements the IntegrationPlugin interface
func (p *{{.CamelName}}Plugin) IsConnected() bool {
	return p.connected
}

// GetStatus implements the IntegrationPlugin interface
func (p *{{.CamelName}}Plugin) GetStatus() string {
	if p.connected {
		return fmt.Sprintf("Connected to %s", p.baseURL)
	}
	return "Disconnected"
}

// Sync implements the IntegrationPlugin interface
func (p *{{.CamelName}}Plugin) Sync() error {
	if !p.connected {
		return fmt.Errorf("not connected to {{.Name}} service")
	}

	p.ctx.Logger().Info("Syncing data with {{.Name}} service")
	
	// Implement synchronization logic here
	// This is a placeholder implementation
	
	return nil
}

// GetData implements the IntegrationPlugin interface
func (p *{{.CamelName}}Plugin) GetData(query map[string]interface{}) (interface{}, error) {
	if !p.connected {
		return nil, fmt.Errorf("not connected to {{.Name}} service")
	}

	p.ctx.Logger().Info("Fetching data from {{.Name}} with query: %+v", query)
	
	// Implement data fetching logic here
	// This is a placeholder implementation
	
	data := map[string]interface{}{
		"status": "success",
		"data":   []interface{}{},
		"query":  query,
	}

	return data, nil
}

// SendData implements the IntegrationPlugin interface
func (p *{{.CamelName}}Plugin) SendData(data interface{}) error {
	if !p.connected {
		return fmt.Errorf("not connected to {{.Name}} service")
	}

	p.ctx.Logger().Info("Sending data to {{.Name}} service")
	
	// Implement data sending logic here
	// This is a placeholder implementation
	
	return nil
}

// Custom methods for this specific integration

// GetProjects fetches projects from {{.Name}}
func (p *{{.CamelName}}Plugin) GetProjects() ([]interface{}, error) {
	query := map[string]interface{}{
		"type": "projects",
	}
	
	result, err := p.GetData(query)
	if err != nil {
		return nil, err
	}

	if data, ok := result.(map[string]interface{}); ok {
		if projects, ok := data["data"].([]interface{}); ok {
			return projects, nil
		}
	}

	return []interface{}{}, nil
}

// CreateProject creates a new project in {{.Name}}
func (p *{{.CamelName}}Plugin) CreateProject(name, description string) error {
	projectData := map[string]interface{}{
		"name":        name,
		"description": description,
		"created_by":  "vibercode-plugin",
	}

	return p.SendData(projectData)
}

// UpdateProject updates an existing project
func (p *{{.CamelName}}Plugin) UpdateProject(id string, updates map[string]interface{}) error {
	updateData := map[string]interface{}{
		"id":      id,
		"updates": updates,
		"action":  "update_project",
	}

	return p.SendData(updateData)
}
`

// PluginGoModTemplate generates go.mod file for plugins
const PluginGoModTemplate = `module github.com/{{.Author}}/{{.Name}}

go 1.19

require (
	github.com/vibercode/plugin-sdk v1.0.0
)

replace github.com/vibercode/plugin-sdk => ../../plugin-sdk
`

// PluginMakefileTemplate generates Makefile for plugins
const PluginMakefileTemplate = `# Makefile for {{.Name}} plugin

BINARY_NAME={{.Name}}
BUILD_DIR=./build
CMD_DIR=./cmd/plugin

.PHONY: build clean test install dev-link lint format deps help

build:
	@echo "Building {{.Name}} plugin..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)

test:
	@echo "Running tests..."
	@go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

install: build
	@echo "Installing {{.Name}} plugin..."
	@vibercode plugins install ./$(BUILD_DIR)/$(BINARY_NAME)

dev-link:
	@echo "Linking {{.Name}} plugin for development..."
	@vibercode plugins dev-link .

dev-unlink:
	@echo "Unlinking {{.Name}} plugin..."
	@vibercode plugins dev-unlink {{.Name}}

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

package:
	@echo "Packaging {{.Name}} plugin..."
	@vibercode plugins package .

publish: package
	@echo "Publishing {{.Name}} plugin..."
	@vibercode plugins publish ./{{.Name}}-{{.Version}}.tar.gz

validate:
	@echo "Validating {{.Name}} plugin..."
	@vibercode plugins validate .

benchmark:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

help:
	@echo "Available targets:"
	@echo "  build         - Build the plugin binary"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  install       - Install the plugin"
	@echo "  dev-link      - Link plugin for development"
	@echo "  dev-unlink    - Unlink plugin from development"
	@echo "  lint          - Run linter"
	@echo "  format        - Format code"
	@echo "  deps          - Download dependencies"
	@echo "  package       - Package plugin for distribution"
	@echo "  publish       - Publish plugin to registry"
	@echo "  validate      - Validate plugin"
	@echo "  benchmark     - Run benchmarks"
	@echo "  help          - Show this help"
`

// PluginReadmeTemplate generates README.md file for plugins
const PluginReadmeTemplate = `# {{.Name}}

{{.Description}}

## Installation

` + "```bash" + `
vibercode plugins install {{.Name}}
` + "```" + `

## Usage

` + "```bash" + `
vibercode {{.KebabName}} --help
` + "```" + `

### Commands

{{- range .Commands}}
- ` + "`{{.Name}}`" + ` - {{.Description}}
{{- end}}

## Configuration

This plugin supports the following configuration options:

{{- range $key, $prop := .Config.Properties}}
- ` + "`{{$key}}`" + ` ({{$prop.Type}}): {{$prop.Description}}{{if $prop.Default}} (default: {{$prop.Default}}){{end}}
{{- end}}

## Examples

` + "```bash" + `
# Basic usage
vibercode {{.KebabName}} --name myproject

# Advanced usage with custom options
vibercode {{.KebabName}} --name myproject --template advanced --output ./output
` + "```" + `

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
make dev-link
` + "```" + `

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run ` + "`make lint`" + ` and ` + "`make test`" + `
6. Submit a pull request

## License

{{.License}} License - see LICENSE file for details.

## Author

{{.Author}}

## Version

{{.Version}}
`