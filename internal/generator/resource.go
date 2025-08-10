package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/manifoldco/promptui"
	"github.com/vibercode/cli/internal/models"
	"github.com/vibercode/cli/internal/templates"
	"github.com/vibercode/cli/pkg/ui"
)

// ResourceGenerator handles resource generation
type ResourceGenerator struct{}

// NewResourceGenerator creates a new ResourceGenerator
func NewResourceGenerator() *ResourceGenerator {
	return &ResourceGenerator{}
}

// Generate generates a complete CRUD resource
func (g *ResourceGenerator) Generate() error {
	ui.PrintHeader("CRUD Resource Generator")
	
	resource, err := g.collectResourceInfo()
	if err != nil {
		return fmt.Errorf("failed to collect resource info: %w", err)
	}

	// Set table name if not provided
	if resource.TableName == "" {
		resource.TableName = strcase.ToSnake(resource.Name) + "s"
	}

	// Show resource summary
	fieldNames := make([]string, len(resource.Fields))
	for i, field := range resource.Fields {
		fieldNames[i] = fmt.Sprintf("%s (%s)", field.Name, field.Type)
	}
	ui.PrintResourceSummary(resource.Name, fieldNames)
	fmt.Println()

	if !ui.ConfirmAction("Generate CRUD resource with this configuration?") {
		ui.PrintInfo("Resource generation cancelled")
		return nil
	}

	// Generate files with progress
	files := []struct {
		name string
		fn   func(*models.Resource) error
	}{
		{"Model", g.generateModel},
		{"Handler", g.generateHandler},
		{"Service", g.generateService},
		{"Repository", g.generateRepository},
	}

	ui.PrintStep(1, 1, "Generating CRUD resource files")
	for _, file := range files {
		spinner := ui.ShowSpinner(fmt.Sprintf("Generating %s...", file.name))
		time.Sleep(300 * time.Millisecond) // Brief pause for UX
		
		if err := file.fn(resource); err != nil {
			spinner.Stop()
			return fmt.Errorf("failed to generate %s: %w", file.name, err)
		}
		
		spinner.Stop()
		ui.PrintFileCreated(fmt.Sprintf("internal/%ss/%s_%s.go", 
			strings.ToLower(file.name), 
			strcase.ToSnake(resource.Name), 
			strings.ToLower(file.name)))
	}

	// Show success message
	fmt.Println()
	ui.PrintSuccess(fmt.Sprintf("CRUD resource '%s' generated successfully!", resource.Name))
	
	generatedFiles := []string{
		fmt.Sprintf("internal/models/%s.go", strcase.ToSnake(resource.Name)),
		fmt.Sprintf("internal/handlers/%s_handler.go", strcase.ToSnake(resource.Name)),
		fmt.Sprintf("internal/services/%s_service.go", strcase.ToSnake(resource.Name)),
		fmt.Sprintf("internal/repositories/%s_repository.go", strcase.ToSnake(resource.Name)),
	}
	
	ui.PrintGeneratedFiles(generatedFiles)

	return nil
}

// collectResourceInfo collects information about the resource from user input
func (g *ResourceGenerator) collectResourceInfo() (*models.Resource, error) {
	resource := &models.Resource{}

	ui.PrintInfo("Let's configure your new CRUD resource")
	fmt.Println()

	// Resource name
	name, err := ui.TextInput(ui.IconCode + " Resource name:")
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(name) == "" {
		ui.ExitWithError("Resource name cannot be empty")
	}
	resource.Name = strings.TrimSpace(name)

	// Description
	desc, err := ui.TextInput(ui.IconInfo + " Resource description:")
	if err != nil {
		return nil, err
	}
	resource.Description = strings.TrimSpace(desc)

	// Module name
	module, err := ui.TextInput(ui.IconPackage + " Go module name:", "github.com/user/project")
	if err != nil {
		return nil, err
	}
	resource.Module = strings.TrimSpace(module)

	// Table name (optional)
	defaultTable := strcase.ToSnake(resource.Name) + "s"
	table, err := ui.TextInput(ui.IconDatabase + " Database table name:", defaultTable)
	if err != nil {
		return nil, err
	}
	resource.TableName = strings.TrimSpace(table)

	// Collect fields
	fields, err := g.collectFields()
	if err != nil {
		return nil, err
	}
	resource.Fields = fields

	return resource, nil
}

// collectFields collects field information from user input
func (g *ResourceGenerator) collectFields() ([]models.Field, error) {
	var fields []models.Field

	for {
		field, err := g.collectField()
		if err != nil {
			return nil, err
		}

		fields = append(fields, *field)

		// Ask if user wants to add another field
		continuePrompt := promptui.Prompt{
			Label:     "Add another field",
			IsConfirm: true,
		}
		
		result, err := continuePrompt.Run()
		if err != nil {
			if err == promptui.ErrAbort {
				break
			}
			return nil, err
		}
		
		if result != "y" {
			break
		}
	}

	return fields, nil
}

// collectField collects information about a single field
func (g *ResourceGenerator) collectField() (*models.Field, error) {
	field := &models.Field{}

	// Field name
	namePrompt := promptui.Prompt{
		Label: "Field name (camelCase)",
		Validate: func(input string) error {
			if len(strings.TrimSpace(input)) == 0 {
				return fmt.Errorf("field name cannot be empty")
			}
			return nil
		},
	}
	name, err := namePrompt.Run()
	if err != nil {
		return nil, err
	}
	field.Name = strings.TrimSpace(name)

	// Field type
	typePrompt := promptui.Select{
		Label: "Field type",
		Items: []string{
			"string",
			"text",
			"number",
			"float",
			"boolean",
			"date",
			"uuid",
			"json",
			"relation",
			"relation-array",
		},
	}
	_, fieldType, err := typePrompt.Run()
	if err != nil {
		return nil, err
	}
	field.Type = models.FieldType(fieldType)

	// Display name
	displayPrompt := promptui.Prompt{
		Label:   "Field display name",
		Default: field.Name,
	}
	display, err := displayPrompt.Run()
	if err != nil {
		return nil, err
	}
	field.DisplayName = strings.TrimSpace(display)

	// Description
	descPrompt := promptui.Prompt{
		Label: "Field description",
	}
	desc, err := descPrompt.Run()
	if err != nil {
		return nil, err
	}
	field.Description = strings.TrimSpace(desc)

	// Required
	requiredPrompt := promptui.Prompt{
		Label:     "Is this field required",
		IsConfirm: true,
	}
	result, err := requiredPrompt.Run()
	if err != nil && err != promptui.ErrAbort {
		return nil, err
	}
	field.Required = result == "y"

	// Handle relation fields
	if field.Type == models.FieldTypeRelation || field.Type == models.FieldTypeRelationArray {
		refPrompt := promptui.Prompt{
			Label: "Related model name",
		}
		ref, err := refPrompt.Run()
		if err != nil {
			return nil, err
		}
		field.Reference = strings.TrimSpace(ref)

		pkgPrompt := promptui.Prompt{
			Label: "Related model package",
		}
		pkg, err := pkgPrompt.Run()
		if err != nil {
			return nil, err
		}
		field.Package = strings.TrimSpace(pkg)
	}

	return field, nil
}

// generateModel generates the model file
func (g *ResourceGenerator) generateModel(resource *models.Resource) error {
	return g.generateFromTemplate(
		resource,
		templates.ModelTemplate,
		"internal/models",
		strcase.ToSnake(resource.Name)+".go",
	)
}

// generateHandler generates the handler file
func (g *ResourceGenerator) generateHandler(resource *models.Resource) error {
	return g.generateFromTemplate(
		resource,
		templates.HandlerTemplate,
		"internal/handlers",
		strcase.ToSnake(resource.Name)+"_handler.go",
	)
}

// generateService generates the service file
func (g *ResourceGenerator) generateService(resource *models.Resource) error {
	return g.generateFromTemplate(
		resource,
		templates.ServiceTemplate,
		"internal/services",
		strcase.ToSnake(resource.Name)+"_service.go",
	)
}

// generateRepository generates the repository file
func (g *ResourceGenerator) generateRepository(resource *models.Resource) error {
	return g.generateFromTemplate(
		resource,
		templates.RepositoryTemplate,
		"internal/repositories",
		strcase.ToSnake(resource.Name)+"_repository.go",
	)
}

// generateFromTemplate generates a file from a template
func (g *ResourceGenerator) generateFromTemplate(resource *models.Resource, templateStr, dir, filename string) error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Create template with helper functions
	tmpl, err := template.New("generator").Funcs(template.FuncMap{
		"ToCamel":      strcase.ToCamel,
		"ToLowerCamel": strcase.ToLowerCamel,
		"ToSnake":      strcase.ToSnake,
		"ToKebab":      strcase.ToKebab,
	}).Parse(templateStr)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Create output file
	outputPath := filepath.Join(dir, filename)
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outputPath, err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, resource); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}