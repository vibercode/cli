package models

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
)

// FieldType represents the type of a field
type FieldType string

const (
	// Basic types
	FieldTypeString      FieldType = "string"
	FieldTypeText        FieldType = "text"
	FieldTypeNumber      FieldType = "number"
	FieldTypeFloat       FieldType = "float"
	FieldTypeBoolean     FieldType = "boolean"
	FieldTypeDate        FieldType = "date"
	FieldTypeUUID        FieldType = "uuid"
	FieldTypeJSON        FieldType = "json"
	
	// Enhanced types
	FieldTypeEmail       FieldType = "email"
	FieldTypeURL         FieldType = "url"
	FieldTypeSlug        FieldType = "slug"
	FieldTypeColor       FieldType = "color"
	FieldTypeFile        FieldType = "file"
	FieldTypeImage       FieldType = "image"
	FieldTypeCoordinates FieldType = "coordinates"
	FieldTypeCurrency    FieldType = "currency"
	FieldTypeEnum        FieldType = "enum"
	FieldTypePassword    FieldType = "password"
	FieldTypePhone       FieldType = "phone"
	
	// Relations
	FieldTypeRelation    FieldType = "relation"
	FieldTypeRelationArray FieldType = "relation-array"
)

// ValidationRule represents a validation rule for a field
type ValidationRule struct {
	Type     string      `json:"type"`
	Value    interface{} `json:"value,omitempty"`
	Message  string      `json:"message,omitempty"`
	Depends  []string    `json:"depends,omitempty"` // For cross-field validation
}

// Field represents a field in a resource
type Field struct {
	Name            string            `json:"name"`
	Type            FieldType         `json:"type"`
	DisplayName     string            `json:"display_name"`
	Description     string            `json:"description"`
	Required        bool              `json:"required"`
	Reference       string            `json:"reference,omitempty"`       // For relations
	Package         string            `json:"package,omitempty"`         // For relations
	EnumValues      []string          `json:"enum_values,omitempty"`     // For enum type
	DefaultValue    interface{}       `json:"default_value,omitempty"`   // Default value
	MinLength       *int              `json:"min_length,omitempty"`      // For strings
	MaxLength       *int              `json:"max_length,omitempty"`      // For strings
	MinValue        *float64          `json:"min_value,omitempty"`       // For numbers
	MaxValue        *float64          `json:"max_value,omitempty"`       // For numbers
	Pattern         string            `json:"pattern,omitempty"`         // Regex pattern
	ValidationRules []ValidationRule  `json:"validation_rules,omitempty"`
	CustomValidator string            `json:"custom_validator,omitempty"` // Custom validation function
	Index           bool              `json:"index,omitempty"`           // Database index
	Unique          bool              `json:"unique,omitempty"`          // Unique constraint
	Nullable        bool              `json:"nullable,omitempty"`        // Nullable in database
}

// GoType returns the Go type for the field
func (f *Field) GoType() string {
	switch f.Type {
	case FieldTypeString, FieldTypeText, FieldTypeEmail, FieldTypeURL, 
		 FieldTypeSlug, FieldTypeColor, FieldTypePassword, FieldTypePhone:
		return "string"
	case FieldTypeNumber:
		return "int"
	case FieldTypeFloat, FieldTypeCurrency:
		return "float64"
	case FieldTypeBoolean:
		return "bool"
	case FieldTypeDate:
		return "time.Time"
	case FieldTypeUUID:
		return "uuid.UUID"
	case FieldTypeJSON:
		return "json.RawMessage"
	case FieldTypeFile, FieldTypeImage:
		return "string" // File path/URL
	case FieldTypeCoordinates:
		return "Coordinates" // Custom struct
	case FieldTypeEnum:
		return f.getEnumTypeName()
	case FieldTypeRelation:
		return "uint" // Foreign key
	case FieldTypeRelationArray:
		return "[]" + f.Reference
	default:
		return "string"
	}
}

// getEnumTypeName returns the enum type name for the field
func (f *Field) getEnumTypeName() string {
	return strcase.ToCamel(f.Name) + "Type"
}

// GoStructField returns the Go struct field definition
func (f *Field) GoStructField() string {
	fieldName := strcase.ToCamel(f.Name)
	fieldType := f.GoType()
	
	// Build tags
	var tags []string
	
	// JSON tag
	jsonTag := `json:"` + strcase.ToSnake(f.Name) + `"`
	if !f.Required || f.Nullable {
		jsonTag = `json:"` + strcase.ToSnake(f.Name) + `,omitempty"`
	}
	tags = append(tags, jsonTag)
	
	// GORM tags
	gormTags := f.buildGORMTags()
	if gormTags != "" {
		tags = append(tags, gormTags)
	}
	
	// Binding tags
	bindingTags := f.buildBindingTags()
	if bindingTags != "" {
		tags = append(tags, bindingTags)
	}
	
	tagString := "`" + strings.Join(tags, " ") + "`"
	
	return fieldName + " " + fieldType + " " + tagString
}

// buildGORMTags builds GORM tags for the field
func (f *Field) buildGORMTags() string {
	var gormParts []string
	
	// Relations
	if f.Type == FieldTypeRelation {
		gormParts = append(gormParts, "foreignKey:"+f.Reference+"ID")
	}
	
	// Unique constraint
	if f.Unique {
		gormParts = append(gormParts, "unique")
	}
	
	// Index
	if f.Index {
		gormParts = append(gormParts, "index")
	}
	
	// Not null
	if f.Required && !f.Nullable {
		gormParts = append(gormParts, "not null")
	}
	
	// Default value
	if f.DefaultValue != nil {
		gormParts = append(gormParts, "default:"+fmt.Sprintf("%v", f.DefaultValue))
	}
	
	// Size constraints for strings
	if f.MaxLength != nil && (f.Type == FieldTypeString || f.Type == FieldTypeText || 
		f.Type == FieldTypeEmail || f.Type == FieldTypeURL || f.Type == FieldTypeSlug) {
		gormParts = append(gormParts, fmt.Sprintf("size:%d", *f.MaxLength))
	}
	
	if len(gormParts) > 0 {
		return `gorm:"` + strings.Join(gormParts, ";") + `"`
	}
	
	return ""
}

// buildBindingTags builds validation binding tags for the field
func (f *Field) buildBindingTags() string {
	var bindingParts []string
	
	// Required validation
	if f.Required {
		bindingParts = append(bindingParts, "required")
	}
	
	// Email validation
	if f.Type == FieldTypeEmail {
		bindingParts = append(bindingParts, "email")
	}
	
	// URL validation
	if f.Type == FieldTypeURL {
		bindingParts = append(bindingParts, "url")
	}
	
	// Length constraints
	if f.MinLength != nil && f.MaxLength != nil {
		bindingParts = append(bindingParts, fmt.Sprintf("min=%d,max=%d", *f.MinLength, *f.MaxLength))
	} else if f.MinLength != nil {
		bindingParts = append(bindingParts, fmt.Sprintf("min=%d", *f.MinLength))
	} else if f.MaxLength != nil {
		bindingParts = append(bindingParts, fmt.Sprintf("max=%d", *f.MaxLength))
	}
	
	// Numeric constraints
	if f.MinValue != nil && f.MaxValue != nil {
		bindingParts = append(bindingParts, fmt.Sprintf("min=%v,max=%v", *f.MinValue, *f.MaxValue))
	} else if f.MinValue != nil {
		bindingParts = append(bindingParts, fmt.Sprintf("min=%v", *f.MinValue))
	} else if f.MaxValue != nil {
		bindingParts = append(bindingParts, fmt.Sprintf("max=%v", *f.MaxValue))
	}
	
	// Pattern validation
	if f.Pattern != "" {
		bindingParts = append(bindingParts, fmt.Sprintf("regexp=%s", f.Pattern))
	}
	
	// Enum validation
	if f.Type == FieldTypeEnum && len(f.EnumValues) > 0 {
		bindingParts = append(bindingParts, "oneof="+strings.Join(f.EnumValues, " "))
	}
	
	if len(bindingParts) > 0 {
		return `binding:"` + strings.Join(bindingParts, ",") + `"`
	}
	
	return ""
}

// GoValidation returns validation code for the field
func (f *Field) GoValidation() string {
	fieldName := strcase.ToCamel(f.Name)
	var validations []string
	
	// Required validation
	if f.Required {
		switch f.Type {
		case FieldTypeString, FieldTypeText, FieldTypeEmail, FieldTypeURL, 
			 FieldTypeSlug, FieldTypeColor, FieldTypePassword, FieldTypePhone:
			validations = append(validations, `if `+fieldName+` == "" {
		return errors.New("`+f.DisplayName+` is required")
	}`)
		case FieldTypeNumber:
			validations = append(validations, `if `+fieldName+` <= 0 {
		return errors.New("`+f.DisplayName+` must be greater than 0")
	}`)
		}
	}
	
	// Type-specific validations
	switch f.Type {
	case FieldTypeEmail:
		validations = append(validations, `if !regexp.MustCompile(`+"`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$`"+`).MatchString(`+fieldName+`) {
		return errors.New("`+f.DisplayName+` must be a valid email address")
	}`)
	case FieldTypeURL:
		validations = append(validations, `if _, err := url.Parse(`+fieldName+`); err != nil {
		return errors.New("`+f.DisplayName+` must be a valid URL")
	}`)
	case FieldTypeSlug:
		validations = append(validations, `if !regexp.MustCompile(`+"`^[a-z0-9-]+$`"+`).MatchString(`+fieldName+`) {
		return errors.New("`+f.DisplayName+` must be a valid slug (lowercase letters, numbers, and hyphens only)")
	}`)
	case FieldTypeColor:
		validations = append(validations, `if !regexp.MustCompile(`+"`^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`"+`).MatchString(`+fieldName+`) {
		return errors.New("`+f.DisplayName+` must be a valid hex color")
	}`)
	case FieldTypePhone:
		validations = append(validations, `if !regexp.MustCompile(`+"`^\\+?[1-9]\\d{1,14}$`"+`).MatchString(`+fieldName+`) {
		return errors.New("`+f.DisplayName+` must be a valid phone number")
	}`)
	}
	
	// Length validations
	if f.MinLength != nil || f.MaxLength != nil {
		if f.MinLength != nil && f.MaxLength != nil {
			validations = append(validations, `if len(`+fieldName+`) < `+fmt.Sprintf("%d", *f.MinLength)+` || len(`+fieldName+`) > `+fmt.Sprintf("%d", *f.MaxLength)+` {
		return errors.New("`+f.DisplayName+` must be between `+fmt.Sprintf("%d", *f.MinLength)+` and `+fmt.Sprintf("%d", *f.MaxLength)+` characters")
	}`)
		} else if f.MinLength != nil {
			validations = append(validations, `if len(`+fieldName+`) < `+fmt.Sprintf("%d", *f.MinLength)+` {
		return errors.New("`+f.DisplayName+` must be at least `+fmt.Sprintf("%d", *f.MinLength)+` characters")
	}`)
		} else if f.MaxLength != nil {
			validations = append(validations, `if len(`+fieldName+`) > `+fmt.Sprintf("%d", *f.MaxLength)+` {
		return errors.New("`+f.DisplayName+` cannot exceed `+fmt.Sprintf("%d", *f.MaxLength)+` characters")
	}`)
		}
	}
	
	// Numeric range validations
	if f.MinValue != nil || f.MaxValue != nil {
		if f.MinValue != nil && f.MaxValue != nil {
			validations = append(validations, `if `+fieldName+` < `+fmt.Sprintf("%v", *f.MinValue)+` || `+fieldName+` > `+fmt.Sprintf("%v", *f.MaxValue)+` {
		return errors.New("`+f.DisplayName+` must be between `+fmt.Sprintf("%v", *f.MinValue)+` and `+fmt.Sprintf("%v", *f.MaxValue)+`")
	}`)
		} else if f.MinValue != nil {
			validations = append(validations, `if `+fieldName+` < `+fmt.Sprintf("%v", *f.MinValue)+` {
		return errors.New("`+f.DisplayName+` must be at least `+fmt.Sprintf("%v", *f.MinValue)+`")
	}`)
		} else if f.MaxValue != nil {
			validations = append(validations, `if `+fieldName+` > `+fmt.Sprintf("%v", *f.MaxValue)+` {
		return errors.New("`+f.DisplayName+` cannot exceed `+fmt.Sprintf("%v", *f.MaxValue)+`")
	}`)
		}
	}
	
	// Pattern validation
	if f.Pattern != "" {
		validations = append(validations, `if !regexp.MustCompile(`+"`"+f.Pattern+"`"+`).MatchString(`+fieldName+`) {
		return errors.New("`+f.DisplayName+` does not match required pattern")
	}`)
	}
	
	// Enum validation
	if f.Type == FieldTypeEnum && len(f.EnumValues) > 0 {
		validValues := strings.Join(f.EnumValues, `", "`)
		validations = append(validations, `validValues := []string{"`+validValues+`"}
	valid := false
	for _, v := range validValues {
		if `+fieldName+` == v {
			valid = true
			break
		}
	}
	if !valid {
		return errors.New("`+f.DisplayName+` must be one of: `+validValues+`")
	}`)
	}
	
	// Custom validations
	for _, rule := range f.ValidationRules {
		if rule.Message != "" {
			validations = append(validations, "// Custom validation: "+rule.Message)
		}
	}
	
	if len(validations) == 0 {
		return "// " + f.DisplayName + " validation can be added here if needed"
	}
	
	return strings.Join(validations, "\n\n\t")
}

// Resource represents a resource to be generated
type Resource struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	TableName   string  `json:"table_name"`
	Module      string  `json:"module"`
	Fields      []Field `json:"fields"`
}

// NameVariations returns various name formats
func (r *Resource) NameVariations() map[string]string {
	return map[string]string{
		"Pascal":       strcase.ToCamel(r.Name),
		"Camel":        strcase.ToLowerCamel(r.Name),
		"Snake":        strcase.ToSnake(r.Name),
		"Kebab":        strcase.ToKebab(r.Name),
		"Lower":        strings.ToLower(r.Name),
		"PluralPascal": strcase.ToCamel(r.Name) + "s", // Simple pluralization
		"PluralCamel":  strcase.ToLowerCamel(r.Name) + "s",
		"PluralSnake":  strcase.ToSnake(r.Name) + "s",
		"PluralKebab":  strcase.ToKebab(r.Name) + "s",
		"PluralLower":  strings.ToLower(r.Name) + "s",
	}
}

// RequiredImports returns the imports needed for the resource
func (r *Resource) RequiredImports() []string {
	imports := []string{
		"errors",
		"gorm.io/gorm",
	}
	
	needsRegexp := false
	needsURL := false
	
	for _, field := range r.Fields {
		switch field.Type {
		case FieldTypeDate:
			imports = append(imports, "time")
		case FieldTypeUUID:
			imports = append(imports, "github.com/google/uuid")
		case FieldTypeJSON:
			imports = append(imports, "encoding/json")
		case FieldTypeEmail, FieldTypeSlug, FieldTypeColor, FieldTypePhone:
			needsRegexp = true
		case FieldTypeURL:
			needsURL = true
		case FieldTypeCoordinates:
			// Custom struct will be defined in the same package
		case FieldTypeFile, FieldTypeImage:
			imports = append(imports, "mime/multipart")
		case FieldTypeCurrency:
			imports = append(imports, "math")
		}
		
		// Check for pattern validation
		if field.Pattern != "" {
			needsRegexp = true
		}
	}
	
	if needsRegexp {
		imports = append(imports, "regexp")
	}
	
	if needsURL {
		imports = append(imports, "net/url")
	}
	
	return removeDuplicates(imports)
}

// removeDuplicates removes duplicate strings from slice
func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	result := []string{}
	
	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}
	
	return result
}

// GetSupportedFieldTypes returns all supported field types
func GetSupportedFieldTypes() []FieldType {
	return []FieldType{
		FieldTypeString, FieldTypeText, FieldTypeNumber, FieldTypeFloat, FieldTypeBoolean,
		FieldTypeDate, FieldTypeUUID, FieldTypeJSON, FieldTypeEmail, FieldTypeURL,
		FieldTypeSlug, FieldTypeColor, FieldTypeFile, FieldTypeImage, FieldTypeCoordinates,
		FieldTypeCurrency, FieldTypeEnum, FieldTypePassword, FieldTypePhone,
		FieldTypeRelation, FieldTypeRelationArray,
	}
}

// GetFieldTypeDescription returns a human-readable description of the field type
func GetFieldTypeDescription(fieldType FieldType) string {
	switch fieldType {
	case FieldTypeString:
		return "Short text field"
	case FieldTypeText:
		return "Long text field"
	case FieldTypeNumber:
		return "Integer number"
	case FieldTypeFloat:
		return "Decimal number"
	case FieldTypeBoolean:
		return "True/false value"
	case FieldTypeDate:
		return "Date and time"
	case FieldTypeUUID:
		return "Unique identifier"
	case FieldTypeJSON:
		return "JSON data"
	case FieldTypeEmail:
		return "Email address with validation"
	case FieldTypeURL:
		return "URL with validation"
	case FieldTypeSlug:
		return "URL-friendly text"
	case FieldTypeColor:
		return "Hex color code"
	case FieldTypeFile:
		return "File upload"
	case FieldTypeImage:
		return "Image upload"
	case FieldTypeCoordinates:
		return "GPS coordinates"
	case FieldTypeCurrency:
		return "Currency amount"
	case FieldTypeEnum:
		return "Predefined options"
	case FieldTypePassword:
		return "Password field"
	case FieldTypePhone:
		return "Phone number"
	case FieldTypeRelation:
		return "Reference to another resource"
	case FieldTypeRelationArray:
		return "Multiple references"
	default:
		return "Unknown field type"
	}
}

// GenerateCoordinatesStruct generates the Coordinates struct definition
func GenerateCoordinatesStruct() string {
	return `// Coordinates represents GPS coordinates
type Coordinates struct {
	Latitude  float64 ` + "`" + `json:"latitude" gorm:"not null"` + "`" + `
	Longitude float64 ` + "`" + `json:"longitude" gorm:"not null"` + "`" + `
}

// String returns a string representation of coordinates
func (c Coordinates) String() string {
	return fmt.Sprintf("%.6f,%.6f", c.Latitude, c.Longitude)
}

// IsValid checks if coordinates are valid
func (c Coordinates) IsValid() bool {
	return c.Latitude >= -90 && c.Latitude <= 90 && 
		   c.Longitude >= -180 && c.Longitude <= 180
}`
}

// GenerateEnumType generates an enum type definition for a field
func (f *Field) GenerateEnumType() string {
	if f.Type != FieldTypeEnum || len(f.EnumValues) == 0 {
		return ""
	}
	
	enumName := f.getEnumTypeName()
	var enumConsts []string
	
	for i, value := range f.EnumValues {
		constName := enumName + strcase.ToCamel(value)
		if i == 0 {
			enumConsts = append(enumConsts, fmt.Sprintf(`	%s %s = "%s"`, constName, enumName, value))
		} else {
			enumConsts = append(enumConsts, fmt.Sprintf(`	%s = "%s"`, constName, value))
		}
	}
	
	return fmt.Sprintf(`// %s represents the enum type for %s
type %s string

const (
%s
)

// String returns the string representation of %s
func (e %s) String() string {
	return string(e)
}

// IsValid checks if the enum value is valid
func (e %s) IsValid() bool {
	switch e {
%s
		return true
	default:
		return false
	}
}`, 
		enumName, f.DisplayName, enumName,
		strings.Join(enumConsts, "\n"),
		enumName, enumName, enumName,
		strings.Join(func() []string {
			var cases []string
			for _, value := range f.EnumValues {
				constName := enumName + strcase.ToCamel(value)
				cases = append(cases, fmt.Sprintf("\tcase %s:", constName))
			}
			return cases
		}(), "\n"))
}