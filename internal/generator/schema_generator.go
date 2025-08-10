package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/vibercode/cli/internal/models"
	"github.com/vibercode/cli/internal/templates"
)

// SchemaGenerator generates code from resource schemas
type SchemaGenerator struct {
	storage models.SchemaStorage
}

// NewSchemaGenerator creates a new schema generator
func NewSchemaGenerator(storage models.SchemaStorage) *SchemaGenerator {
	return &SchemaGenerator{
		storage: storage,
	}
}

// GenerateFromSchema generates complete CRUD code from a schema
func (g *SchemaGenerator) GenerateFromSchema(schemaID, outputPath, module string, dbProvider string) error {
	// Load schema
	schema, err := g.storage.Load(schemaID)
	if err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}

	// Prepare template data
	data := g.prepareTemplateData(schema, module, dbProvider)

	// Generate files
	generators := map[string]string{
		"model":      filepath.Join("internal", "models", schema.Names.SnakeCase+".go"),
		"repository": filepath.Join("internal", "repositories", schema.Names.SnakeCase+"_repository.go"),
		"service":    filepath.Join("internal", "services", schema.Names.SnakeCase+"_service.go"),
		"handler":    filepath.Join("internal", "handlers", schema.Names.SnakeCase+"_handler.go"),
	}

	schemaTemplates := templates.GetSchemaTemplates()

	for templateName, relativePath := range generators {
		template := schemaTemplates[templateName]
		fullPath := filepath.Join(outputPath, relativePath)

		if err := g.generateFile(template, data, fullPath); err != nil {
			return fmt.Errorf("failed to generate %s: %w", templateName, err)
		}
	}

	// Generate migration file
	if err := g.generateMigration(schema, outputPath, dbProvider); err != nil {
		return fmt.Errorf("failed to generate migration: %w", err)
	}

	return nil
}

// GenerateFromSchemaName generates code from schema name
func (g *SchemaGenerator) GenerateFromSchemaName(schemaName, outputPath, module string, dbProvider string) error {
	schema, err := g.storage.LoadByName(schemaName)
	if err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}

	return g.GenerateFromSchema(schema.ID, outputPath, module, dbProvider)
}

// prepareTemplateData prepares data for template execution
func (g *SchemaGenerator) prepareTemplateData(schema *models.ResourceSchema, module string, dbProvider string) *EnhancedSchema {
	return g.prepareEnhancedTemplateData(schema, module, dbProvider)
}

// TemplateData represents data passed to templates (deprecated, use EnhancedSchema)
type TemplateData struct {
	*models.ResourceSchema
	Module          string
	DBProvider      string
	Relations       []RelationInfo
	RequiredImports []string
	GetPreloads     string
	GetSearchFields string
	GetSearchValues string
}

// RelationInfo represents relation information for templates
type RelationInfo struct {
	Name        string
	DisplayName string
	Type        string
	Target      string
}

// FieldNamingConventions represents naming conventions for a field
type FieldNamingConventions struct {
	PascalCase string
	CamelCase  string
	SnakeCase  string
	KebabCase  string
}

// extractRelations extracts relation information from schema
func (g *SchemaGenerator) extractRelations(schema *models.ResourceSchema) []RelationInfo {
	var relations []RelationInfo
	seen := make(map[string]bool)

	for _, field := range schema.Fields {
		if (field.Type == "relation" || field.Type == "relation_array") && field.Relation != nil {
			target := field.Relation.Target
			if !seen[target] {
				relations = append(relations, RelationInfo{
					Name:        target,
					DisplayName: target,
					Type:        field.Relation.Type,
					Target:      target,
				})
				seen[target] = true
			}
		}
	}

	return relations
}

// getRequiredImports returns required imports for the schema
func (g *SchemaGenerator) getRequiredImports(schema *models.ResourceSchema, dbProvider string) []string {
	imports := make(map[string]bool)

	for _, field := range schema.Fields {
		switch field.Type {
		case "uuid":
			imports["github.com/google/uuid"] = true
		case "decimal", "currency":
			imports["github.com/shopspring/decimal"] = true
		case "json", "mixed":
			imports["encoding/json"] = true
		case "date", "datetime", "timestamp":
			imports["time"] = true
		}

		// Database-specific imports
		if dbProvider == "supabase" {
			imports["github.com/supabase-community/gotrue-go"] = true
		}
	}

	var result []string
	for imp := range imports {
		result = append(result, imp)
	}

	return result
}

// generateGoStructField generates Go struct field definition
func (g *SchemaGenerator) generateGoStructField(field *models.SchemaField, dbProvider string) string {
	fieldName := toPascalCase(field.Name)
	fieldType := field.GetGoType()
	gormTag := field.GetGORMTag(dbProvider)
	jsonTag := fmt.Sprintf(`json:"%s"`, toSnakeCase(field.Name))

	var tags []string
	tags = append(tags, jsonTag)
	if gormTag != "" {
		tags = append(tags, fmt.Sprintf(`gorm:"%s"`, gormTag))
	}

	// Add validation tags
	if field.Required {
		tags = append(tags, `binding:"required"`)
	}

	tagString := "`" + strings.Join(tags, " ") + "`"
	return fmt.Sprintf("%s %s %s", fieldName, fieldType, tagString)
}

// generateGoRequestField generates Go request field definition
func (g *SchemaGenerator) generateGoRequestField(field *models.SchemaField) string {
	fieldName := toPascalCase(field.Name)
	fieldType := field.GetGoType()
	jsonTag := fmt.Sprintf(`json:"%s"`, toSnakeCase(field.Name))

	var tags []string
	tags = append(tags, jsonTag)

	if field.Required {
		tags = append(tags, `binding:"required"`)
	}

	// Add validation tags based on field type
	switch field.Type {
	case "email":
		tags = append(tags, `binding:"email"`)
	case "url":
		tags = append(tags, `binding:"url"`)
	case "number", "integer":
		if field.Validation != nil {
			if field.Validation.Min != nil {
				tags = append(tags, fmt.Sprintf(`binding:"min=%g"`, *field.Validation.Min))
			}
			if field.Validation.Max != nil {
				tags = append(tags, fmt.Sprintf(`binding:"max=%g"`, *field.Validation.Max))
			}
		}
	case "string", "text":
		if field.Validation != nil {
			if field.Validation.MinLength != nil {
				tags = append(tags, fmt.Sprintf(`binding:"min=%d"`, *field.Validation.MinLength))
			}
			if field.Validation.MaxLength != nil {
				tags = append(tags, fmt.Sprintf(`binding:"max=%d"`, *field.Validation.MaxLength))
			}
		}
	}

	tagString := "`" + strings.Join(tags, " ") + "`"
	return fmt.Sprintf("%s %s %s", fieldName, fieldType, tagString)
}

// generateGoResponseField generates Go response field definition
func (g *SchemaGenerator) generateGoResponseField(field *models.SchemaField) string {
	fieldName := toPascalCase(field.Name)
	fieldType := field.GetGoType()
	jsonTag := fmt.Sprintf(`json:"%s"`, toSnakeCase(field.Name))

	tagString := "`" + jsonTag + "`"
	return fmt.Sprintf("%s %s %s", fieldName, fieldType, tagString)
}

// generateGoFilterField generates Go filter field definition
func (g *SchemaGenerator) generateGoFilterField(field *models.SchemaField) string {
	fieldName := toPascalCase(field.Name)
	fieldType := field.GetGoType()
	
	// Make filter fields optional (pointers for primitive types)
	switch fieldType {
	case "string", "int64", "float64", "bool":
		fieldType = "*" + fieldType
	}

	jsonTag := fmt.Sprintf(`json:"%s,omitempty" form:"%s"`, toSnakeCase(field.Name), toSnakeCase(field.Name))
	tagString := "`" + jsonTag + "`"
	
	return fmt.Sprintf("%s %s %s", fieldName, fieldType, tagString)
}

// generateGoValidation generates Go validation code
func (g *SchemaGenerator) generateGoValidation(field *models.SchemaField) string {
	fieldName := toPascalCase(field.Name)
	
	switch field.Type {
	case "string", "text", "email", "url":
		if field.Required {
			return fmt.Sprintf(`if r.%s == "" {
		return fmt.Errorf("%s is required")
	}`, fieldName, field.DisplayName)
		}
	case "number", "integer":
		if field.Required {
			return fmt.Sprintf(`if r.%s <= 0 {
		return fmt.Errorf("%s must be greater than 0")
	}`, fieldName, field.DisplayName)
		}
	}
	
	return fmt.Sprintf("// %s validation can be added here if needed", field.DisplayName)
}

// generateGoFilterQuery generates Go filter query code
func (g *SchemaGenerator) generateGoFilterQuery(field *models.SchemaField) string {
	fieldName := toPascalCase(field.Name)
	columnName := toSnakeCase(field.Name)
	
	switch field.Type {
	case "string", "text", "email", "url":
		return fmt.Sprintf(`if filter.%s != nil {
		query = query.Where("%s ILIKE ?", "%%"+*filter.%s+"%%")
	}`, fieldName, columnName, fieldName)
	case "number", "integer", "float", "decimal":
		return fmt.Sprintf(`if filter.%s != nil {
		query = query.Where("%s = ?", *filter.%s)
	}`, fieldName, columnName, fieldName)
	case "boolean":
		return fmt.Sprintf(`if filter.%s != nil {
		query = query.Where("%s = ?", *filter.%s)
	}`, fieldName, columnName, fieldName)
	case "date", "datetime", "timestamp":
		return fmt.Sprintf(`if filter.%s != nil {
		query = query.Where("DATE(%s) = DATE(?)", *filter.%s)
	}`, fieldName, columnName, fieldName)
	}
	
	return ""
}

// isFieldFilterable determines if a field should be filterable
func (g *SchemaGenerator) isFieldFilterable(field *models.SchemaField) bool {
	switch field.Type {
	case "string", "text", "email", "url", "number", "integer", "float", "decimal", "boolean", "date", "datetime", "timestamp":
		return true
	case "relation":
		return true
	default:
		return false
	}
}

// isFieldReadOnly determines if a field is read-only
func (g *SchemaGenerator) isFieldReadOnly(field *models.SchemaField) bool {
	// Typically timestamps and IDs are read-only
	readOnlyFields := []string{"id", "created_at", "updated_at", "deleted_at"}
	for _, ro := range readOnlyFields {
		if field.Name == ro {
			return true
		}
	}
	return false
}

// generatePreloads generates GORM preload statements
func (g *SchemaGenerator) generatePreloads(schema *models.ResourceSchema) string {
	var preloads []string
	
	for _, field := range schema.Fields {
		if field.Type == "relation" || field.Type == "relation_array" {
			if field.Relation != nil && field.Relation.Populate {
				preloads = append(preloads, fmt.Sprintf(`.Preload("%s")`, toPascalCase(field.Name)))
			}
		}
	}
	
	return strings.Join(preloads, "")
}

// generateSearchFields generates search field conditions
func (g *SchemaGenerator) generateSearchFields(schema *models.ResourceSchema) string {
	var conditions []string
	
	for _, field := range schema.Fields {
		if field.Type == "string" || field.Type == "text" {
			conditions = append(conditions, fmt.Sprintf("%s ILIKE ?", toSnakeCase(field.Name)))
		}
	}
	
	if len(conditions) == 0 {
		return "name ILIKE ?"
	}
	
	return strings.Join(conditions, " OR ")
}

// generateSearchValues generates search value parameters
func (g *SchemaGenerator) generateSearchValues(schema *models.ResourceSchema) string {
	var values []string
	
	for _, field := range schema.Fields {
		if field.Type == "string" || field.Type == "text" {
			values = append(values, "searchQuery")
		}
	}
	
	if len(values) == 0 {
		return "searchQuery"
	}
	
	return strings.Join(values, ", ")
}

// generateFieldNamingConventions generates naming conventions for a field
func generateFieldNamingConventions(name string) *FieldNamingConventions {
	return &FieldNamingConventions{
		PascalCase: toPascalCase(name),
		CamelCase:  toCamelCase(name),
		SnakeCase:  toSnakeCase(name),
		KebabCase:  toKebabCase(name),
	}
}

// generateFile generates a file from template
func (g *SchemaGenerator) generateFile(templateStr string, data interface{}, outputPath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Parse template
	tmpl, err := template.New("generator").Funcs(templates.SchemaHelperFunctions).Parse(templateStr)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Create file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outputPath, err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// generateMigration generates database migration file
func (g *SchemaGenerator) generateMigration(schema *models.ResourceSchema, outputPath, dbProvider string) error {
	// Skip migration generation for MongoDB as it doesn't require schema migrations
	if dbProvider == "mongodb" {
		return nil
	}
	
	migrationTemplate := `package migrations

import (
	"gorm.io/gorm"
	"{{.Module}}/internal/models"
)

// Migration{{.Names.PascalCase}} migrates {{.DisplayName}} table
func Migration{{.Names.PascalCase}}(db *gorm.DB) error {
	return db.AutoMigrate(&models.{{.Names.PascalCase}}{})
}

// Rollback{{.Names.PascalCase}} rolls back {{.DisplayName}} table
func Rollback{{.Names.PascalCase}}(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.{{.Names.PascalCase}}{})
}
`

	data := g.prepareTemplateData(schema, "", dbProvider)
	migrationPath := filepath.Join(outputPath, "migrations", fmt.Sprintf("%s_migration.go", schema.Names.SnakeCase))
	
	return g.generateFile(migrationTemplate, data, migrationPath)
}

// Helper functions (reuse from previous implementation)
func toPascalCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == ' ' || r == '_' || r == '-'
	})
	for i, word := range words {
		words[i] = strings.Title(strings.ToLower(word))
	}
	return strings.Join(words, "")
}

func toCamelCase(s string) string {
	pascal := toPascalCase(s)
	if len(pascal) == 0 {
		return pascal
	}
	return strings.ToLower(pascal[:1]) + pascal[1:]
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && (r >= 'A' && r <= 'Z') {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

func toKebabCase(s string) string {
	return strings.ReplaceAll(toSnakeCase(s), "_", "-")
}