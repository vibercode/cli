package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vibercode/cli/internal/generator"
	"github.com/vibercode/cli/internal/models"
	"github.com/vibercode/cli/internal/storage"
	"github.com/vibercode/cli/pkg/ui"
)

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "📋 Manage resource schemas",
	Long: ui.Bold.Sprint("Schema Management") + "\n\n" +
		"Create, manage and generate code from resource schemas.\n\n" +
		ui.Bold.Sprint("Available commands:") + "\n" +
		"  " + ui.IconCode + " create    - Create a new resource schema\n" +
		"  " + ui.IconDatabase + " list      - List all schemas\n" +
		"  " + ui.IconGear + " generate  - Generate code from schema\n" +
		"  " + ui.IconDoc + " show      - Show schema details\n" +
		"  " + ui.IconBuild + " delete    - Delete a schema\n",
}

var schemaCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "📝 Create a new resource schema",
	Long: ui.Bold.Sprint("Create a new resource schema") + "\n\n" +
		"Interactive schema creation with field definitions, validation rules,\n" +
		"and UI configuration for frontend generation.\n",
	RunE: func(cmd *cobra.Command, args []string) error {
		return createSchema()
	},
}

var schemaListCmd = &cobra.Command{
	Use:   "list",
	Short: "📋 List all schemas",
	Long: ui.Bold.Sprint("List all available schemas") + "\n\n" +
		"Shows all stored schemas with basic information.\n",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listSchemas()
	},
}

var schemaGenerateCmd = &cobra.Command{
	Use:   "generate [schema-name]",
	Short: "🔨 Generate code from schema",
	Long: ui.Bold.Sprint("Generate code from a stored schema") + "\n\n" +
		"Generates complete CRUD code including models, repositories,\n" +
		"services, handlers, and migrations.\n",
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var schemaName string
		if len(args) > 0 {
			schemaName = args[0]
		}
		return generateFromSchema(schemaName)
	},
}

var schemaShowCmd = &cobra.Command{
	Use:   "show [schema-name]",
	Short: "👁️  Show schema details",
	Long: ui.Bold.Sprint("Show detailed schema information") + "\n\n" +
		"Displays complete schema definition including fields,\n" +
		"validation rules, and configuration.\n",
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var schemaName string
		if len(args) > 0 {
			schemaName = args[0]
		}
		return showSchema(schemaName)
	},
}

var schemaDeleteCmd = &cobra.Command{
	Use:   "delete [schema-name]",
	Short: "🗑️  Delete a schema",
	Long: ui.Bold.Sprint("Delete a stored schema") + "\n\n" +
		"Permanently removes a schema from storage.\n",
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var schemaName string
		if len(args) > 0 {
			schemaName = args[0]
		}
		return deleteSchema(schemaName)
	},
}

// Command flags
var (
	outputDir    string
	module       string
	dbProvider   string
	templateName string
)

func init() {
	// Add schema commands
	schemaCmd.AddCommand(schemaCreateCmd)
	schemaCmd.AddCommand(schemaListCmd)
	schemaCmd.AddCommand(schemaGenerateCmd)
	schemaCmd.AddCommand(schemaShowCmd)
	schemaCmd.AddCommand(schemaDeleteCmd)

	// Add flags
	schemaGenerateCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory for generated code")
	schemaGenerateCmd.Flags().StringVarP(&module, "module", "m", "", "Go module name")
	schemaGenerateCmd.Flags().StringVarP(&dbProvider, "database", "d", "postgres", "Database provider (postgres, mysql, sqlite, supabase, mongodb)")

	schemaCreateCmd.Flags().StringVarP(&templateName, "template", "t", "", "Use a predefined template")
}

// createSchema creates a new resource schema interactively
func createSchema() error {
	ui.PrintHeader("Create Resource Schema")

	// Get storage
	storagePath := storage.GetDefaultSchemaPath()
	schemaStorage := storage.NewFileSchemaStorage(storagePath)
	repo := storage.NewSchemaRepository(schemaStorage)

	// Check if user wants to use a template
	useTemplate := ui.ConfirmAction("Would you like to start from a template?")

	var schema *models.ResourceSchema

	if useTemplate {
		var err error
		schema, err = selectTemplate()
		if err != nil {
			return err
		}
	} else {
		schema = &models.ResourceSchema{}
	}

	// Collect basic information
	if err := collectSchemaBasicInfo(schema); err != nil {
		return err
	}

	// Collect fields
	if err := collectSchemaFields(schema); err != nil {
		return err
	}

	// Collect database configuration
	if err := collectDatabaseConfig(schema); err != nil {
		return err
	}

	// Collect frontend configuration
	if err := collectFrontendConfig(schema); err != nil {
		return err
	}

	// Show summary and confirm
	if err := showSchemaSummary(schema); err != nil {
		return err
	}

	if !ui.ConfirmAction("Create schema with this configuration?") {
		ui.PrintInfo("Schema creation cancelled")
		return nil
	}

	// Save schema
	if err := repo.CreateSchema(schema); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	ui.PrintSuccess(fmt.Sprintf("Schema '%s' created successfully!", schema.Name))
	ui.PrintInfo(fmt.Sprintf("Schema ID: %s", schema.ID))

	return nil
}

// selectTemplate allows user to select from predefined templates
func selectTemplate() (*models.ResourceSchema, error) {
	templates := storage.LoadSchemaTemplates()
	
	var options []string
	for _, tmpl := range templates {
		options = append(options, fmt.Sprintf("%s - %s", tmpl.Name, tmpl.Description))
	}

	selected, err := ui.SelectOption("Select a template:", options)
	if err != nil {
		return nil, err
	}

	// Find selected template
	for i, option := range options {
		if option == selected {
			return templates[i].Schema, nil
		}
	}

	return nil, fmt.Errorf("template not found")
}

// collectSchemaBasicInfo collects basic schema information
func collectSchemaBasicInfo(schema *models.ResourceSchema) error {
	ui.PrintSubHeader("Basic Information")

	// Name (if not from template)
	if schema.Name == "" {
		name, err := ui.TextInput(ui.IconAPI + " Resource name:")
		if err != nil {
			return err
		}
		schema.Name = strings.TrimSpace(name)
	} else {
		// Allow editing template name
		name, err := ui.TextInput(ui.IconAPI + " Resource name:", schema.Name)
		if err != nil {
			return err
		}
		schema.Name = strings.TrimSpace(name)
	}

	// Display name
	displayName, err := ui.TextInput(ui.IconDoc + " Display name:", schema.Name)
	if err != nil {
		return err
	}
	schema.DisplayName = strings.TrimSpace(displayName)

	// Description
	description, err := ui.TextInput(ui.IconDoc + " Description:", schema.Description)
	if err != nil {
		return err
	}
	schema.Description = strings.TrimSpace(description)

	return nil
}

// collectSchemaFields collects field definitions
func collectSchemaFields(schema *models.ResourceSchema) error {
	ui.PrintSubHeader("Field Definitions")

	// If schema has fields (from template), ask if user wants to modify
	if len(schema.Fields) > 0 {
		modify := ui.ConfirmAction("Template has fields defined. Would you like to modify them?")
		if !modify {
			return nil
		}
	}

	for {
		field := models.SchemaField{}

		// Field name
		name, err := ui.TextInput(ui.IconCode + " Field name (empty to finish):")
		if err != nil {
			return err
		}
		name = strings.TrimSpace(name)
		if name == "" {
			break
		}
		field.Name = name

		// Field type
		fieldTypes := []string{
			"string", "text", "number", "float", "boolean", "date", "datetime",
			"email", "url", "slug", "color", "file", "image", "location",
			"uuid", "json", "currency", "enum", "relation", "relation_array",
		}
		fieldType, err := ui.SelectOption(ui.IconGear + " Field type:", fieldTypes)
		if err != nil {
			return err
		}
		field.Type = fieldType

		// Display name
		displayName, err := ui.TextInput(ui.IconDoc + " Display name:", field.Name)
		if err != nil {
			return err
		}
		field.DisplayName = strings.TrimSpace(displayName)

		// Required
		required := ui.ConfirmAction("Is this field required?")
		field.Required = required

		// Handle relations
		if field.Type == "relation" || field.Type == "relation_array" {
			if err := collectRelationConfig(&field); err != nil {
				return err
			}
		}

		// Collect validation rules
		if err := collectValidationRules(&field); err != nil {
			return err
		}

		schema.Fields = append(schema.Fields, field)
		ui.PrintSuccess(fmt.Sprintf("Field '%s' added", field.Name))
	}

	return nil
}

// collectRelationConfig collects relation configuration
func collectRelationConfig(field *models.SchemaField) error {
	field.Relation = &models.RelationConfig{}

	// Target model
	target, err := ui.TextInput(ui.IconDatabase + " Target model:")
	if err != nil {
		return err
	}
	field.Relation.Target = strings.TrimSpace(target)

	// Relation type
	if field.Type == "relation" {
		field.Relation.Type = "one_to_one"
	} else {
		field.Relation.Type = "one_to_many"
	}

	// Foreign key
	foreignKey, err := ui.TextInput(ui.IconGear + " Foreign key field:", field.Name+"_id")
	if err != nil {
		return err
	}
	field.Relation.ForeignKey = strings.TrimSpace(foreignKey)

	// Populate
	populate := ui.ConfirmAction("Auto-populate this relation?")
	field.Relation.Populate = populate

	return nil
}

// collectValidationRules collects validation rules
func collectValidationRules(field *models.SchemaField) error {
	addValidation := ui.ConfirmAction("Add validation rules?")

	if !addValidation {
		return nil
	}

	field.Validation = &models.FieldValidation{}

	switch field.Type {
	case "string", "text", "email", "url":
		// Min/Max length
		if minLen, err := ui.TextInput("Minimum length (optional):"); err == nil && minLen != "" {
			if val, err := parseIntPointer(minLen); err == nil {
				field.Validation.MinLength = val
			}
		}
		if maxLen, err := ui.TextInput("Maximum length (optional):"); err == nil && maxLen != "" {
			if val, err := parseIntPointer(maxLen); err == nil {
				field.Validation.MaxLength = val
			}
		}

		// Pattern
		if pattern, err := ui.TextInput("Regex pattern (optional):"); err == nil && pattern != "" {
			field.Validation.Pattern = pattern
		}

	case "number", "float", "currency":
		// Min/Max value
		if min, err := ui.TextInput("Minimum value (optional):"); err == nil && min != "" {
			if val, err := parseFloatPointer(min); err == nil {
				field.Validation.Min = val
			}
		}
		if max, err := ui.TextInput("Maximum value (optional):"); err == nil && max != "" {
			if val, err := parseFloatPointer(max); err == nil {
				field.Validation.Max = val
			}
		}
	}

	return nil
}

// collectDatabaseConfig collects database configuration
func collectDatabaseConfig(schema *models.ResourceSchema) error {
	ui.PrintSubHeader("Database Configuration")

	// Database provider
	providers := []string{"postgres", "mysql", "sqlite", "supabase", "mongodb"}
	provider, err := ui.SelectOption(ui.IconDatabase + " Database provider:", providers)
	if err != nil {
		return err
	}

	schema.Database = &models.DatabaseConfig{
		Provider: provider,
	}

	// Table name
	tableName, err := ui.TextInput(ui.IconGear + " Table name:", strings.ToLower(schema.Name)+"s")
	if err != nil {
		return err
	}
	schema.Database.TableName = strings.TrimSpace(tableName)

	return nil
}

// collectFrontendConfig collects frontend configuration
func collectFrontendConfig(schema *models.ResourceSchema) error {
	ui.PrintSubHeader("Frontend Configuration")

	generateFrontend := ui.ConfirmAction("Generate frontend components?")

	if !generateFrontend {
		return nil
	}

	schema.Frontend = &models.FrontendConfig{}

	// Framework
	frameworks := []string{"react", "vue", "angular", "svelte"}
	framework, err := ui.SelectOption(ui.IconCode + " Frontend framework:", frameworks)
	if err != nil {
		return err
	}
	schema.Frontend.Framework = framework

	// Component style
	componentStyles := []string{"atomic", "feature", "page"}
	componentStyle, err := ui.SelectOption(ui.IconGear + " Component style:", componentStyles)
	if err != nil {
		return err
	}
	schema.Frontend.ComponentStyle = componentStyle

	// TypeScript
	typescript := ui.ConfirmAction("Use TypeScript?")
	schema.Frontend.TypeScript = typescript

	return nil
}

// showSchemaSummary shows schema summary before creation
func showSchemaSummary(schema *models.ResourceSchema) error {
	ui.PrintSubHeader("Schema Summary")

	ui.PrintFeature(ui.IconAPI, "Name", schema.Name)
	ui.PrintFeature(ui.IconDoc, "Display Name", schema.DisplayName)
	ui.PrintFeature(ui.IconDoc, "Description", schema.Description)

	if schema.Database != nil {
		ui.PrintFeature(ui.IconDatabase, "Database", schema.Database.Provider)
		ui.PrintFeature(ui.IconGear, "Table Name", schema.Database.TableName)
	}

	if schema.Frontend != nil {
		ui.PrintFeature(ui.IconCode, "Frontend", schema.Frontend.Framework)
		ui.PrintFeature(ui.IconGear, "Component Style", schema.Frontend.ComponentStyle)
		ui.PrintFeature(ui.IconBuild, "TypeScript", fmt.Sprintf("%t", schema.Frontend.TypeScript))
	}

	ui.PrintSubHeader("Fields")
	for _, field := range schema.Fields {
		requiredStr := ""
		if field.Required {
			requiredStr = " (required)"
		}
		ui.PrintFeature(ui.IconCode, field.Name, fmt.Sprintf("%s%s - %s", field.Type, requiredStr, field.DisplayName))
	}

	return nil
}

// listSchemas lists all available schemas
func listSchemas() error {
	storagePath := storage.GetDefaultSchemaPath()
	schemaStorage := storage.NewFileSchemaStorage(storagePath)

	schemas, err := schemaStorage.List()
	if err != nil {
		return fmt.Errorf("failed to list schemas: %w", err)
	}

	if len(schemas) == 0 {
		ui.PrintInfo("No schemas found. Create one with 'vibercode schema create'")
		return nil
	}

	ui.PrintHeader("Available Schemas")

	for _, schema := range schemas {
		ui.PrintFeature(ui.IconAPI, schema.Name, schema.Description)
		ui.PrintInfo(fmt.Sprintf("  ID: %s, Fields: %d, Created: %s",
			schema.ID,
			len(schema.Fields),
			schema.CreatedAt.Format("2006-01-02"),
		))
		fmt.Println()
	}

	return nil
}

// generateFromSchema generates code from a stored schema
func generateFromSchema(schemaName string) error {
	storagePath := storage.GetDefaultSchemaPath()
	schemaStorage := storage.NewFileSchemaStorage(storagePath)

	// If no schema name provided, list and select
	if schemaName == "" {
		schemas, err := schemaStorage.List()
		if err != nil {
			return fmt.Errorf("failed to list schemas: %w", err)
		}

		if len(schemas) == 0 {
			ui.PrintInfo("No schemas found. Create one with 'vibercode schema create'")
			return nil
		}

		var options []string
		for _, schema := range schemas {
			options = append(options, fmt.Sprintf("%s - %s", schema.Name, schema.Description))
		}

		selected, err := ui.SelectOption("Select schema to generate:", options)
		if err != nil {
			return err
		}

		// Extract schema name
		schemaName = strings.Split(selected, " - ")[0]
	}

	// Get schema
	schema, err := schemaStorage.LoadByName(schemaName)
	if err != nil {
		return fmt.Errorf("schema not found: %w", err)
	}

	// Get generation options
	if module == "" {
		moduleInput, err := ui.TextInput(ui.IconPackage + " Go module name:")
		if err != nil {
			return err
		}
		module = strings.TrimSpace(moduleInput)
	}

	if outputDir == "" {
		outputInput, err := ui.TextInput(ui.IconGear + " Output directory:", ".")
		if err != nil {
			return err
		}
		outputDir = strings.TrimSpace(outputInput)
	}

	// Convert to absolute path
	outputDir, err = filepath.Abs(outputDir)
	if err != nil {
		return fmt.Errorf("invalid output directory: %w", err)
	}

	// Use schema's database provider if not specified
	if dbProvider == "" {
		if schema.Database != nil && schema.Database.Provider != "" {
			dbProvider = schema.Database.Provider
		} else {
			dbProvider = "postgres"
		}
	}

	ui.PrintHeader("Generating Code")
	ui.PrintFeature(ui.IconAPI, "Schema", schema.Name)
	ui.PrintFeature(ui.IconPackage, "Module", module)
	ui.PrintFeature(ui.IconDatabase, "Database", dbProvider)
	ui.PrintFeature(ui.IconGear, "Output", outputDir)

	if !ui.ConfirmAction("Generate code with these settings?") {
		ui.PrintInfo("Code generation cancelled")
		return nil
	}

	// Generate code
	generator := generator.NewSchemaGenerator(schemaStorage)
	if err := generator.GenerateFromSchema(schema.ID, outputDir, module, dbProvider); err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}

	ui.PrintSuccess("Code generated successfully!")
	ui.PrintInfo(fmt.Sprintf("Generated files in: %s", outputDir))

	return nil
}

// showSchema shows detailed schema information
func showSchema(schemaName string) error {
	storagePath := storage.GetDefaultSchemaPath()
	schemaStorage := storage.NewFileSchemaStorage(storagePath)

	// If no schema name provided, list and select
	if schemaName == "" {
		schemas, err := schemaStorage.List()
		if err != nil {
			return fmt.Errorf("failed to list schemas: %w", err)
		}

		if len(schemas) == 0 {
			ui.PrintInfo("No schemas found. Create one with 'vibercode schema create'")
			return nil
		}

		var options []string
		for _, schema := range schemas {
			options = append(options, schema.Name)
		}

		selected, err := ui.SelectOption("Select schema to show:", options)
		if err != nil {
			return err
		}
		schemaName = selected
	}

	// Get schema
	schema, err := schemaStorage.LoadByName(schemaName)
	if err != nil {
		return fmt.Errorf("schema not found: %w", err)
	}

	// Display schema as JSON
	data, err := schema.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize schema: %w", err)
	}

	ui.PrintHeader(fmt.Sprintf("Schema: %s", schema.Name))
	fmt.Println(string(data))

	return nil
}

// deleteSchema deletes a schema
func deleteSchema(schemaName string) error {
	storagePath := storage.GetDefaultSchemaPath()
	schemaStorage := storage.NewFileSchemaStorage(storagePath)

	// If no schema name provided, list and select
	if schemaName == "" {
		schemas, err := schemaStorage.List()
		if err != nil {
			return fmt.Errorf("failed to list schemas: %w", err)
		}

		if len(schemas) == 0 {
			ui.PrintInfo("No schemas found")
			return nil
		}

		var options []string
		for _, schema := range schemas {
			options = append(options, schema.Name)
		}

		selected, err := ui.SelectOption("Select schema to delete:", options)
		if err != nil {
			return err
		}
		schemaName = selected
	}

	// Get schema
	schema, err := schemaStorage.LoadByName(schemaName)
	if err != nil {
		return fmt.Errorf("schema not found: %w", err)
	}

	// Confirm deletion
	if !ui.ConfirmAction(fmt.Sprintf("Are you sure you want to delete schema '%s'?", schema.Name)) {
		ui.PrintInfo("Deletion cancelled")
		return nil
	}

	// Delete schema
	if err := schemaStorage.Delete(schema.ID); err != nil {
		return fmt.Errorf("failed to delete schema: %w", err)
	}

	ui.PrintSuccess(fmt.Sprintf("Schema '%s' deleted successfully", schema.Name))
	return nil
}

// Helper functions
func parseIntPointer(s string) (*int, error) {
	if s == "" {
		return nil, nil
	}
	var val int
	if _, err := fmt.Sscanf(s, "%d", &val); err != nil {
		return nil, err
	}
	return &val, nil
}

func parseFloatPointer(s string) (*float64, error) {
	if s == "" {
		return nil, nil
	}
	var val float64
	if _, err := fmt.Sscanf(s, "%f", &val); err != nil {
		return nil, err
	}
	return &val, nil
}