package models

import (
	"strings"
	"testing"
)

func TestField_GoType(t *testing.T) {
	tests := []struct {
		name     string
		field    Field
		expected string
	}{
		{
			name:     "String type",
			field:    Field{Type: FieldTypeString},
			expected: "string",
		},
		{
			name:     "Email type",
			field:    Field{Type: FieldTypeEmail},
			expected: "string",
		},
		{
			name:     "URL type",
			field:    Field{Type: FieldTypeURL},
			expected: "string",
		},
		{
			name:     "Number type",
			field:    Field{Type: FieldTypeNumber},
			expected: "int",
		},
		{
			name:     "Float type",
			field:    Field{Type: FieldTypeFloat},
			expected: "float64",
		},
		{
			name:     "Currency type",
			field:    Field{Type: FieldTypeCurrency},
			expected: "float64",
		},
		{
			name:     "Boolean type",
			field:    Field{Type: FieldTypeBoolean},
			expected: "bool",
		},
		{
			name:     "Date type",
			field:    Field{Type: FieldTypeDate},
			expected: "time.Time",
		},
		{
			name:     "UUID type",
			field:    Field{Type: FieldTypeUUID},
			expected: "uuid.UUID",
		},
		{
			name:     "File type",
			field:    Field{Type: FieldTypeFile},
			expected: "string",
		},
		{
			name:     "Image type",
			field:    Field{Type: FieldTypeImage},
			expected: "string",
		},
		{
			name:     "Coordinates type",
			field:    Field{Type: FieldTypeCoordinates},
			expected: "Coordinates",
		},
		{
			name:     "Enum type",
			field:    Field{Type: FieldTypeEnum, Name: "status"},
			expected: "StatusType",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.field.GoType()
			if result != tt.expected {
				t.Errorf("GoType() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestField_GoValidation(t *testing.T) {
	tests := []struct {
		name     string
		field    Field
		contains []string // Strings that should be present in validation
	}{
		{
			name: "Required email field",
			field: Field{
				Type:        FieldTypeEmail,
				Name:        "email",
				DisplayName: "Email",
				Required:    true,
			},
			contains: []string{"required", "valid email"},
		},
		{
			name: "URL field with validation",
			field: Field{
				Type:        FieldTypeURL,
				Name:        "website",
				DisplayName: "Website",
				Required:    false,
			},
			contains: []string{"valid URL"},
		},
		{
			name: "String with length constraints",
			field: Field{
				Type:        FieldTypeString,
				Name:        "name",
				DisplayName: "Name",
				Required:    true,
				MinLength:   &[]int{3}[0],
				MaxLength:   &[]int{50}[0],
			},
			contains: []string{"required", "between", "3", "50"},
		},
		{
			name: "Number with range constraints",
			field: Field{
				Type:        FieldTypeNumber,
				Name:        "age",
				DisplayName: "Age",
				MinValue:    &[]float64{18}[0],
				MaxValue:    &[]float64{100}[0],
			},
			contains: []string{"between", "18", "100"},
		},
		{
			name: "Enum field",
			field: Field{
				Type:        FieldTypeEnum,
				Name:        "status",
				DisplayName: "Status",
				EnumValues:  []string{"active", "inactive", "pending"},
			},
			contains: []string{"one of", "active", "inactive", "pending"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.field.GoValidation()
			
			for _, contain := range tt.contains {
				if !containsString(result, contain) {
					t.Errorf("GoValidation() = %v, should contain %v", result, contain)
				}
			}
		})
	}
}

func TestField_GoStructField(t *testing.T) {
	tests := []struct {
		name     string
		field    Field
		contains []string
	}{
		{
			name: "Required email field",
			field: Field{
				Type:        FieldTypeEmail,
				Name:        "email",
				DisplayName: "Email",
				Required:    true,
			},
			contains: []string{"Email", "string", "json:", "binding:", "required", "email"},
		},
		{
			name: "Optional URL field",
			field: Field{
				Type:        FieldTypeURL,
				Name:        "website_url",
				DisplayName: "Website URL",
				Required:    false,
			},
			contains: []string{"WebsiteUrl", "string", "json:", "omitempty", "url"},
		},
		{
			name: "Unique slug field",
			field: Field{
				Type:        FieldTypeSlug,
				Name:        "slug",
				DisplayName: "Slug",
				Required:    true,
				Unique:      true,
			},
			contains: []string{"Slug", "string", "json:", "gorm:", "unique"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.field.GoStructField()
			
			for _, contain := range tt.contains {
				if !containsString(result, contain) {
					t.Errorf("GoStructField() = %v, should contain %v", result, contain)
				}
			}
		})
	}
}

func TestGenerateEnumType(t *testing.T) {
	field := Field{
		Type:        FieldTypeEnum,
		Name:        "status",
		DisplayName: "Status",
		EnumValues:  []string{"active", "inactive", "pending"},
	}

	result := field.GenerateEnumType()
	
	expectedStrings := []string{
		"StatusType",
		"StatusTypeActive",
		"StatusTypeInactive", 
		"StatusTypePending",
		"const (",
		"IsValid()",
		"String()",
	}

	for _, expected := range expectedStrings {
		if !containsString(result, expected) {
			t.Errorf("GenerateEnumType() should contain %v, got: %v", expected, result)
		}
	}
}

func TestGenerateCoordinatesStruct(t *testing.T) {
	result := GenerateCoordinatesStruct()
	
	expectedStrings := []string{
		"type Coordinates struct",
		"Latitude",
		"Longitude",
		"float64",
		"IsValid()",
		"String()",
	}

	for _, expected := range expectedStrings {
		if !containsString(result, expected) {
			t.Errorf("GenerateCoordinatesStruct() should contain %v, got: %v", expected, result)
		}
	}
}

func TestGetSupportedFieldTypes(t *testing.T) {
	types := GetSupportedFieldTypes()
	
	// Check that we have all expected types
	expectedTypes := []FieldType{
		FieldTypeString, FieldTypeEmail, FieldTypeURL, FieldTypeSlug,
		FieldTypeNumber, FieldTypeFloat, FieldTypeCurrency,
		FieldTypeBoolean, FieldTypeDate, FieldTypeUUID,
		FieldTypeFile, FieldTypeImage, FieldTypeCoordinates,
		FieldTypeEnum, FieldTypePassword, FieldTypePhone,
	}

	for _, expectedType := range expectedTypes {
		found := false
		for _, actualType := range types {
			if actualType == expectedType {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected field type %v not found in supported types", expectedType)
		}
	}
}

func TestGetFieldTypeDescription(t *testing.T) {
	tests := []struct {
		fieldType   FieldType
		description string
	}{
		{FieldTypeEmail, "Email address with validation"},
		{FieldTypeURL, "URL with validation"},
		{FieldTypeSlug, "URL-friendly text"},
		{FieldTypeColor, "Hex color code"},
		{FieldTypeFile, "File upload"},
		{FieldTypeImage, "Image upload"},
		{FieldTypeCoordinates, "GPS coordinates"},
		{FieldTypeCurrency, "Currency amount"},
		{FieldTypeEnum, "Predefined options"},
		{FieldTypePassword, "Password field"},
		{FieldTypePhone, "Phone number"},
	}

	for _, tt := range tests {
		t.Run(string(tt.fieldType), func(t *testing.T) {
			result := GetFieldTypeDescription(tt.fieldType)
			if result != tt.description {
				t.Errorf("GetFieldTypeDescription(%v) = %v, expected %v", tt.fieldType, result, tt.description)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return strings.Contains(s, substr)
}

func TestResource_RequiredImports(t *testing.T) {
	resource := Resource{
		Name: "User",
		Fields: []Field{
			{Type: FieldTypeEmail, Name: "email"},
			{Type: FieldTypeDate, Name: "birth_date"},
			{Type: FieldTypeUUID, Name: "uuid"},
			{Type: FieldTypeCoordinates, Name: "location"},
			{Type: FieldTypeFile, Name: "avatar"},
		},
	}

	imports := resource.RequiredImports()
	
	expectedImports := []string{
		"time",
		"github.com/google/uuid",
		"regexp",
		"mime/multipart",
	}

	for _, expectedImport := range expectedImports {
		found := false
		for _, actualImport := range imports {
			if actualImport == expectedImport {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected import %v not found in required imports: %v", expectedImport, imports)
		}
	}
}