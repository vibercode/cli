package generator

import (
	"github.com/vibercode/cli/internal/models"
)

// EnhancedField extends SchemaField with template-specific data
type EnhancedField struct {
	*models.SchemaField
	GoStructField  string
	GoRequestField string
	GoResponseField string
	GoFilterField  string
	GoValidation   string
	GoFilterQuery  string
	Names          *FieldNamingConventions
	Filterable     bool
	ReadOnly       bool
}

// EnhancedSchema extends ResourceSchema with enhanced fields
type EnhancedSchema struct {
	*models.ResourceSchema
	Fields          []EnhancedField
	Module          string
	DBProvider      string
	Relations       []RelationInfo
	RequiredImports []string
	GetPreloads     string
	GetSearchFields string
	GetSearchValues string
}

// enhanceField converts a SchemaField to EnhancedField
func (g *SchemaGenerator) enhanceField(field *models.SchemaField, dbProvider string) EnhancedField {
	enhanced := EnhancedField{
		SchemaField: field,
		Names:       generateFieldNamingConventions(field.Name),
		Filterable:  g.isFieldFilterable(field),
		ReadOnly:    g.isFieldReadOnly(field),
	}

	enhanced.GoStructField = g.generateGoStructField(field, dbProvider)
	enhanced.GoRequestField = g.generateGoRequestField(field)
	enhanced.GoResponseField = g.generateGoResponseField(field)
	enhanced.GoFilterField = g.generateGoFilterField(field)
	enhanced.GoValidation = g.generateGoValidation(field)
	enhanced.GoFilterQuery = g.generateGoFilterQuery(field)

	return enhanced
}

// prepareEnhancedTemplateData prepares enhanced data for template execution
func (g *SchemaGenerator) prepareEnhancedTemplateData(schema *models.ResourceSchema, module string, dbProvider string) *EnhancedSchema {
	enhanced := &EnhancedSchema{
		ResourceSchema:  schema,
		Module:         module,
		DBProvider:     dbProvider,
		Relations:      g.extractRelations(schema),
		RequiredImports: g.getRequiredImports(schema, dbProvider),
	}

	// Enhance fields
	enhanced.Fields = make([]EnhancedField, len(schema.Fields))
	for i, field := range schema.Fields {
		enhanced.Fields[i] = g.enhanceField(&field, dbProvider)
	}

	// Add helper methods
	enhanced.GetPreloads = g.generatePreloads(schema)
	enhanced.GetSearchFields = g.generateSearchFields(schema)
	enhanced.GetSearchValues = g.generateSearchValues(schema)

	return enhanced
}