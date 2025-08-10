package templates

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/vibercode/cli/pkg/ui"
)

// TemplateCategory represents the category of a template
type TemplateCategory string

const (
	CategoryFullstackAI TemplateCategory = "fullstack-ai"
	CategoryFrontend    TemplateCategory = "frontend"
	CategoryBackend     TemplateCategory = "backend"
	CategoryAPI         TemplateCategory = "api"
	CategoryResource    TemplateCategory = "resource"
)

// TemplateType represents the type/framework of a template
type TemplateType string

const (
	TypeGo      TemplateType = "go"
	TypeReact   TemplateType = "react"
	TypeVue     TemplateType = "vue"
	TypeAngular TemplateType = "angular"
	TypeNextJS  TemplateType = "nextjs"
	TypeSvelte  TemplateType = "svelte"
)

// TemplateMetadata contains information about a template
type TemplateMetadata struct {
	ID            string           `json:"id" yaml:"id"`
	Name          string           `json:"name" yaml:"name"`
	DisplayName   string           `json:"display_name" yaml:"display_name"`
	Description   string           `json:"description" yaml:"description"`
	Category      TemplateCategory `json:"category" yaml:"category"`
	Type          TemplateType     `json:"type" yaml:"type"`
	Version       string           `json:"version" yaml:"version"`
	Author        string           `json:"author" yaml:"author"`
	Tags          []string         `json:"tags" yaml:"tags"`
	Dependencies  []string         `json:"dependencies" yaml:"dependencies"`
	Requirements  []string         `json:"requirements" yaml:"requirements"`
	Framework     string           `json:"framework" yaml:"framework"`
	Language      string           `json:"language" yaml:"language"`
	CreatedAt     time.Time        `json:"created_at" yaml:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at" yaml:"updated_at"`
	Deprecated    bool             `json:"deprecated" yaml:"deprecated"`
	Experimental  bool             `json:"experimental" yaml:"experimental"`
	Documentation string           `json:"documentation" yaml:"documentation"`
	Examples      []string         `json:"examples" yaml:"examples"`
	Variables     []TemplateVar    `json:"variables" yaml:"variables"`
	Files         []TemplateFile   `json:"files" yaml:"files"`
}

// TemplateVar represents a configurable variable in a template
type TemplateVar struct {
	Name         string      `json:"name" yaml:"name"`
	DisplayName  string      `json:"display_name" yaml:"display_name"`
	Description  string      `json:"description" yaml:"description"`
	Type         string      `json:"type" yaml:"type"` // string, number, boolean, array, object
	Default      interface{} `json:"default" yaml:"default"`
	Required     bool        `json:"required" yaml:"required"`
	Options      []string    `json:"options,omitempty" yaml:"options,omitempty"`
	Pattern      string      `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	Min          *float64    `json:"min,omitempty" yaml:"min,omitempty"`
	Max          *float64    `json:"max,omitempty" yaml:"max,omitempty"`
	Placeholder  string      `json:"placeholder,omitempty" yaml:"placeholder,omitempty"`
}

// TemplateFile represents a file to be generated from a template
type TemplateFile struct {
	Path        string            `json:"path" yaml:"path"`
	Template    string            `json:"template" yaml:"template"`
	Condition   string            `json:"condition,omitempty" yaml:"condition,omitempty"`
	Permissions string            `json:"permissions,omitempty" yaml:"permissions,omitempty"`
	Variables   map[string]string `json:"variables,omitempty" yaml:"variables,omitempty"`
}

// TemplateRegistry manages all available templates
type TemplateRegistry struct {
	templates   map[string]*TemplateMetadata
	templateDir string
	funcMap     template.FuncMap
}

// NewTemplateRegistry creates a new template registry
func NewTemplateRegistry(templateDir string) *TemplateRegistry {
	return &TemplateRegistry{
		templates:   make(map[string]*TemplateMetadata),
		templateDir: templateDir,
		funcMap:     CreateTemplateFuncMap(),
	}
}

// LoadTemplates loads all templates from the template directory
func (r *TemplateRegistry) LoadTemplates() error {
	if _, err := os.Stat(r.templateDir); os.IsNotExist(err) {
		ui.PrintWarning(fmt.Sprintf("Template directory does not exist: %s", r.templateDir))
		return r.loadBuiltinTemplates()
	}

	ui.PrintInfo(fmt.Sprintf("Loading templates from: %s", r.templateDir))

	err := filepath.WalkDir(r.templateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, "template.json") || strings.HasSuffix(path, "template.yaml") {
			return r.loadTemplateMetadata(path)
		}

		return nil
	})

	if err != nil {
		ui.PrintError(fmt.Sprintf("Failed to load templates: %v", err))
		return err
	}

	// Also load builtin templates
	if err := r.loadBuiltinTemplates(); err != nil {
		ui.PrintWarning(fmt.Sprintf("Failed to load builtin templates: %v", err))
	}

	ui.PrintSuccess(fmt.Sprintf("Loaded %d templates", len(r.templates)))
	return nil
}

// loadTemplateMetadata loads a single template metadata file
func (r *TemplateRegistry) loadTemplateMetadata(metadataPath string) error {
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to read template metadata %s: %w", metadataPath, err)
	}

	var metadata TemplateMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return fmt.Errorf("failed to parse template metadata %s: %w", metadataPath, err)
	}

	// Validate required fields
	if metadata.ID == "" || metadata.Name == "" || metadata.Category == "" || metadata.Type == "" {
		return fmt.Errorf("template metadata %s missing required fields", metadataPath)
	}

	// Set template directory relative to metadata file
	templateDir := filepath.Dir(metadataPath)
	
	// Load and validate template files
	for i, file := range metadata.Files {
		templatePath := filepath.Join(templateDir, file.Template)
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			return fmt.Errorf("template file not found: %s", templatePath)
		}
		metadata.Files[i].Template = templatePath
	}

	r.templates[metadata.ID] = &metadata
	ui.PrintInfo(fmt.Sprintf("Loaded template: %s (%s)", metadata.DisplayName, metadata.ID))
	return nil
}

// loadBuiltinTemplates loads the existing Go CLI templates as builtin templates
func (r *TemplateRegistry) loadBuiltinTemplates() error {
	// Load existing Go CLI templates
	builtinTemplates := []*TemplateMetadata{
		{
			ID:          "go-api-resource",
			Name:        "go-api-resource",
			DisplayName: "Go API Resource",
			Description: "Generate Go API resource with Clean Architecture (Model, Repository, Service, Handler)",
			Category:    CategoryAPI,
			Type:        TypeGo,
			Version:     "1.0.0",
			Author:      "ViberCode CLI",
			Tags:        []string{"go", "api", "crud", "clean-architecture"},
			Framework:   "Gin",
			Language:    "Go",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Variables: []TemplateVar{
				{
					Name:        "schema",
					DisplayName: "Resource Schema",
					Description: "The resource schema definition",
					Type:        "object",
					Required:    true,
				},
				{
					Name:        "module",
					DisplayName: "Go Module",
					Description: "The Go module name",
					Type:        "string",
					Required:    true,
					Default:     "github.com/example/project",
				},
			},
			Files: []TemplateFile{
				{Path: "internal/models/{{.Names.SnakeCase}}.go", Template: "builtin:model"},
				{Path: "internal/repositories/{{.Names.SnakeCase}}.go", Template: "builtin:repository"},
				{Path: "internal/services/{{.Names.SnakeCase}}.go", Template: "builtin:service"},
				{Path: "internal/handlers/{{.Names.SnakeCase}}.go", Template: "builtin:handler"},
			},
		},
		{
			ID:          "react-crud-component",
			Name:        "react-crud-component",
			DisplayName: "React CRUD Component",
			Description: "Generate React CRUD components with TypeScript and hooks",
			Category:    CategoryFrontend,
			Type:        TypeReact,
			Version:     "1.0.0",
			Author:      "ViberCode CLI",
			Tags:        []string{"react", "typescript", "crud", "hooks"},
			Framework:   "React",
			Language:    "TypeScript",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Experimental: true,
			Variables: []TemplateVar{
				{
					Name:        "schema",
					DisplayName: "Resource Schema",
					Description: "The resource schema definition",
					Type:        "object",
					Required:    true,
				},
				{
					Name:        "apiBaseUrl",
					DisplayName: "API Base URL",
					Description: "The base URL for the API",
					Type:        "string",
					Required:    true,
					Default:     "/api/v1",
				},
			},
			Files: []TemplateFile{
				{Path: "src/components/{{.Names.PascalCase}}/{{.Names.PascalCase}}List.tsx", Template: "builtin:react-list"},
				{Path: "src/components/{{.Names.PascalCase}}/{{.Names.PascalCase}}Form.tsx", Template: "builtin:react-form"},
				{Path: "src/components/{{.Names.PascalCase}}/{{.Names.PascalCase}}Detail.tsx", Template: "builtin:react-detail"},
				{Path: "src/types/{{.Names.CamelCase}}.ts", Template: "builtin:react-types"},
				{Path: "src/hooks/use{{.Names.PascalCase}}.ts", Template: "builtin:react-hooks"},
			},
		},
	}

	for _, template := range builtinTemplates {
		r.templates[template.ID] = template
		ui.PrintInfo(fmt.Sprintf("Loaded builtin template: %s (%s)", template.DisplayName, template.ID))
	}

	return nil
}

// GetTemplate retrieves a template by ID
func (r *TemplateRegistry) GetTemplate(id string) (*TemplateMetadata, error) {
	template, exists := r.templates[id]
	if !exists {
		return nil, fmt.Errorf("template not found: %s", id)
	}
	return template, nil
}

// ListTemplates returns all available templates
func (r *TemplateRegistry) ListTemplates() []*TemplateMetadata {
	templates := make([]*TemplateMetadata, 0, len(r.templates))
	for _, template := range r.templates {
		templates = append(templates, template)
	}
	return templates
}

// FilterTemplates returns templates filtered by category and/or type
func (r *TemplateRegistry) FilterTemplates(category TemplateCategory, templateType TemplateType) []*TemplateMetadata {
	templates := make([]*TemplateMetadata, 0)
	for _, template := range r.templates {
		if category != "" && template.Category != category {
			continue
		}
		if templateType != "" && template.Type != templateType {
			continue
		}
		templates = append(templates, template)
	}
	return templates
}

// ValidateTemplate validates a template's structure and dependencies
func (r *TemplateRegistry) ValidateTemplate(id string) error {
	template, err := r.GetTemplate(id)
	if err != nil {
		return err
	}

	// Validate template files exist
	for _, file := range template.Files {
		if strings.HasPrefix(file.Template, "builtin:") {
			// Skip validation for builtin templates
			continue
		}
		if _, err := os.Stat(file.Template); os.IsNotExist(err) {
			return fmt.Errorf("template file not found: %s", file.Template)
		}
	}

	// Validate dependencies
	for _, dep := range template.Dependencies {
		if _, err := r.GetTemplate(dep); err != nil {
			return fmt.Errorf("dependency template not found: %s", dep)
		}
	}

	return nil
}

// GenerateFromTemplate generates files from a template
func (r *TemplateRegistry) GenerateFromTemplate(templateID string, outputDir string, variables map[string]interface{}) error {
	template, err := r.GetTemplate(templateID)
	if err != nil {
		return err
	}

	if err := r.ValidateTemplate(templateID); err != nil {
		return fmt.Errorf("template validation failed: %w", err)
	}

	ui.PrintInfo(fmt.Sprintf("Generating from template: %s", template.DisplayName))

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate each file
	for _, file := range template.Files {
		if err := r.generateFile(file, outputDir, variables); err != nil {
			return fmt.Errorf("failed to generate file %s: %w", file.Path, err)
		}
	}

	ui.PrintSuccess(fmt.Sprintf("Generated %d files from template %s", len(template.Files), template.DisplayName))
	return nil
}

// generateFile generates a single file from a template file definition
func (r *TemplateRegistry) generateFile(file TemplateFile, outputDir string, variables map[string]interface{}) error {
	// Evaluate condition if present
	if file.Condition != "" {
		// TODO: Implement condition evaluation
		ui.PrintWarning(fmt.Sprintf("Condition evaluation not yet implemented for file: %s", file.Path))
	}

	// Parse file path template
	pathTemplate, err := template.New("path").Funcs(r.funcMap).Parse(file.Path)
	if err != nil {
		return fmt.Errorf("failed to parse path template: %w", err)
	}

	var pathBuf strings.Builder
	if err := pathTemplate.Execute(&pathBuf, variables); err != nil {
		return fmt.Errorf("failed to execute path template: %w", err)
	}

	outputPath := filepath.Join(outputDir, pathBuf.String())

	// Create directory if needed
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Generate file content
	var content string
	if strings.HasPrefix(file.Template, "builtin:") {
		// Handle builtin templates
		content, err = r.getBuiltinTemplate(file.Template, variables)
		if err != nil {
			return fmt.Errorf("failed to get builtin template: %w", err)
		}
	} else {
		// Handle file-based templates
		content, err = r.executeFileTemplate(file.Template, variables)
		if err != nil {
			return fmt.Errorf("failed to execute template: %w", err)
		}
	}

	// Write file
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	ui.PrintFileCreated(outputPath)
	return nil
}

// getBuiltinTemplate returns content for builtin templates
func (r *TemplateRegistry) getBuiltinTemplate(templateName string, variables map[string]interface{}) (string, error) {
	switch templateName {
	case "builtin:model":
		return r.executeTemplate(SchemaModelTemplate, variables)
	case "builtin:repository":
		return r.executeTemplate(SchemaRepositoryTemplate, variables)
	case "builtin:service":
		return r.executeTemplate(SchemaServiceTemplate, variables)
	case "builtin:handler":
		return r.executeTemplate(SchemaHandlerTemplate, variables)
	case "builtin:react-list":
		return r.executeTemplate(ReactListTemplate, variables)
	case "builtin:react-form":
		return r.executeTemplate(ReactFormTemplate, variables)
	case "builtin:react-detail":
		return r.executeTemplate(ReactDetailTemplate, variables)
	case "builtin:react-types":
		return r.executeTemplate(ReactTypesTemplate, variables)
	case "builtin:react-hooks":
		return r.executeTemplate(ReactHooksTemplate, variables)
	default:
		return "", fmt.Errorf("unknown builtin template: %s", templateName)
	}
}

// executeFileTemplate reads and executes a template file
func (r *TemplateRegistry) executeFileTemplate(templatePath string, variables map[string]interface{}) (string, error) {
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %w", err)
	}

	return r.executeTemplate(string(templateContent), variables)
}

// executeTemplate executes a template string with variables
func (r *TemplateRegistry) executeTemplate(templateContent string, variables map[string]interface{}) (string, error) {
	tmpl, err := template.New("template").Funcs(r.funcMap).Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, variables); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// CreateTemplateFuncMap creates the function map for template execution
func CreateTemplateFuncMap() template.FuncMap {
	funcMap := make(template.FuncMap)
	
	// Copy existing schema helper functions
	for k, v := range SchemaHelperFunctions {
		funcMap[k] = v
	}

	// Add additional helper functions
	funcMap["toJSON"] = func(v interface{}) string {
		data, _ := json.Marshal(v)
		return string(data)
	}

	funcMap["indent"] = func(spaces int, text string) string {
		indentation := strings.Repeat(" ", spaces)
		lines := strings.Split(text, "\n")
		for i, line := range lines {
			if line != "" {
				lines[i] = indentation + line
			}
		}
		return strings.Join(lines, "\n")
	}

	funcMap["camelCase"] = func(s string) string {
		// Convert to camelCase
		words := strings.FieldsFunc(s, func(c rune) bool {
			return !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9'))
		})
		if len(words) == 0 {
			return s
		}
		result := strings.ToLower(words[0])
		for _, word := range words[1:] {
			result += strings.Title(strings.ToLower(word))
		}
		return result
	}

	funcMap["pascalCase"] = func(s string) string {
		words := strings.FieldsFunc(s, func(c rune) bool {
			return !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9'))
		})
		var result string
		for _, word := range words {
			result += strings.Title(strings.ToLower(word))
		}
		return result
	}

	funcMap["kebabCase"] = func(s string) string {
		words := strings.FieldsFunc(s, func(c rune) bool {
			return !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9'))
		})
		var result []string
		for _, word := range words {
			result = append(result, strings.ToLower(word))
		}
		return strings.Join(result, "-")
	}

	funcMap["snakeCase"] = func(s string) string {
		words := strings.FieldsFunc(s, func(c rune) bool {
			return !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9'))
		})
		var result []string
		for _, word := range words {
			result = append(result, strings.ToLower(word))
		}
		return strings.Join(result, "_")
	}

	return funcMap
}