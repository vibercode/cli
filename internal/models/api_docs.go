package models

import (
	"fmt"
	"time"
)

// OpenAPIVersion represents the OpenAPI specification version
const OpenAPIVersion = "3.0.3"

// APIDocumentationConfig represents configuration for API documentation generation
type APIDocumentationConfig struct {
	ProjectName     string             `json:"project_name" yaml:"project_name"`
	Version         string             `json:"version" yaml:"version"`
	Description     string             `json:"description" yaml:"description"`
	Contact         ContactInfo        `json:"contact" yaml:"contact"`
	License         LicenseInfo        `json:"license" yaml:"license"`
	Servers         []ServerInfo       `json:"servers" yaml:"servers"`
	BasePath        string             `json:"base_path" yaml:"base_path"`
	EnableSwaggerUI bool               `json:"enable_swagger_ui" yaml:"enable_swagger_ui"`
	SwaggerUIPath   string             `json:"swagger_ui_path" yaml:"swagger_ui_path"`
	OutputFormat    []DocumentFormat   `json:"output_formats" yaml:"output_formats"`
	Security        SecuritySchemes    `json:"security" yaml:"security"`
	Tags            []Tag              `json:"tags" yaml:"tags"`
	ExternalDocs    *ExternalDocsInfo  `json:"external_docs" yaml:"external_docs"`
	Extensions      map[string]interface{} `json:"extensions" yaml:"extensions"`
}

// ContactInfo represents API contact information
type ContactInfo struct {
	Name  string `json:"name" yaml:"name"`
	Email string `json:"email" yaml:"email"`
	URL   string `json:"url" yaml:"url"`
}

// LicenseInfo represents API license information
type LicenseInfo struct {
	Name string `json:"name" yaml:"name"`
	URL  string `json:"url" yaml:"url"`
}

// ServerInfo represents API server information
type ServerInfo struct {
	URL         string                 `json:"url" yaml:"url"`
	Description string                 `json:"description" yaml:"description"`
	Variables   map[string]interface{} `json:"variables" yaml:"variables"`
}

// DocumentFormat represents available documentation formats
type DocumentFormat string

const (
	FormatJSON DocumentFormat = "json"
	FormatYAML DocumentFormat = "yaml"
	FormatHTML DocumentFormat = "html"
	FormatPDF  DocumentFormat = "pdf"
)

// SecuritySchemes represents API security configurations
type SecuritySchemes struct {
	JWT    *JWTSecurityScheme    `json:"jwt" yaml:"jwt"`
	APIKey *APIKeySecurityScheme `json:"api_key" yaml:"api_key"`
	OAuth2 *OAuth2SecurityScheme `json:"oauth2" yaml:"oauth2"`
}

// JWTSecurityScheme represents JWT authentication scheme
type JWTSecurityScheme struct {
	Type         string `json:"type" yaml:"type"`
	Scheme       string `json:"scheme" yaml:"scheme"`
	BearerFormat string `json:"bearerFormat" yaml:"bearerFormat"`
	Description  string `json:"description" yaml:"description"`
}

// APIKeySecurityScheme represents API key authentication scheme
type APIKeySecurityScheme struct {
	Type        string `json:"type" yaml:"type"`
	In          string `json:"in" yaml:"in"`
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
}

// OAuth2SecurityScheme represents OAuth2 authentication scheme
type OAuth2SecurityScheme struct {
	Type             string                    `json:"type" yaml:"type"`
	Flows            OAuth2Flows               `json:"flows" yaml:"flows"`
	Description      string                    `json:"description" yaml:"description"`
	OpenIdConnectUrl string                    `json:"openIdConnectUrl" yaml:"openIdConnectUrl"`
}

// OAuth2Flows represents OAuth2 flow configurations
type OAuth2Flows struct {
	Implicit          *OAuth2Flow `json:"implicit" yaml:"implicit"`
	Password          *OAuth2Flow `json:"password" yaml:"password"`
	ClientCredentials *OAuth2Flow `json:"clientCredentials" yaml:"clientCredentials"`
	AuthorizationCode *OAuth2Flow `json:"authorizationCode" yaml:"authorizationCode"`
}

// OAuth2Flow represents a specific OAuth2 flow
type OAuth2Flow struct {
	AuthorizationUrl string            `json:"authorizationUrl" yaml:"authorizationUrl"`
	TokenUrl         string            `json:"tokenUrl" yaml:"tokenUrl"`
	RefreshUrl       string            `json:"refreshUrl" yaml:"refreshUrl"`
	Scopes           map[string]string `json:"scopes" yaml:"scopes"`
}

// Tag represents an API tag for grouping operations
type Tag struct {
	Name         string            `json:"name" yaml:"name"`
	Description  string            `json:"description" yaml:"description"`
	ExternalDocs *ExternalDocsInfo `json:"externalDocs" yaml:"externalDocs"`
}

// ExternalDocsInfo represents external documentation
type ExternalDocsInfo struct {
	Description string `json:"description" yaml:"description"`
	URL         string `json:"url" yaml:"url"`
}

// OpenAPISpec represents the complete OpenAPI specification
type OpenAPISpec struct {
	OpenAPI      string                    `json:"openapi" yaml:"openapi"`
	Info         InfoObject                `json:"info" yaml:"info"`
	Servers      []ServerInfo              `json:"servers" yaml:"servers"`
	Paths        map[string]*PathItem      `json:"paths" yaml:"paths"`
	Components   *ComponentsObject         `json:"components" yaml:"components"`
	Security     []map[string][]string     `json:"security" yaml:"security"`
	Tags         []Tag                     `json:"tags" yaml:"tags"`
	ExternalDocs *ExternalDocsInfo         `json:"externalDocs" yaml:"externalDocs"`
}

// InfoObject represents API information
type InfoObject struct {
	Title          string       `json:"title" yaml:"title"`
	Description    string       `json:"description" yaml:"description"`
	TermsOfService string       `json:"termsOfService" yaml:"termsOfService"`
	Contact        *ContactInfo `json:"contact" yaml:"contact"`
	License        *LicenseInfo `json:"license" yaml:"license"`
	Version        string       `json:"version" yaml:"version"`
}

// PathItem represents a single API path
type PathItem struct {
	Summary     string     `json:"summary" yaml:"summary"`
	Description string     `json:"description" yaml:"description"`
	Get         *Operation `json:"get" yaml:"get"`
	Put         *Operation `json:"put" yaml:"put"`
	Post        *Operation `json:"post" yaml:"post"`
	Delete      *Operation `json:"delete" yaml:"delete"`
	Options     *Operation `json:"options" yaml:"options"`
	Head        *Operation `json:"head" yaml:"head"`
	Patch       *Operation `json:"patch" yaml:"patch"`
	Trace       *Operation `json:"trace" yaml:"trace"`
	Parameters  []Parameter `json:"parameters" yaml:"parameters"`
}

// Operation represents a single API operation
type Operation struct {
	Tags         []string              `json:"tags" yaml:"tags"`
	Summary      string                `json:"summary" yaml:"summary"`
	Description  string                `json:"description" yaml:"description"`
	ExternalDocs *ExternalDocsInfo     `json:"externalDocs" yaml:"externalDocs"`
	OperationId  string                `json:"operationId" yaml:"operationId"`
	Parameters   []Parameter           `json:"parameters" yaml:"parameters"`
	RequestBody  *RequestBody          `json:"requestBody" yaml:"requestBody"`
	Responses    map[string]*Response  `json:"responses" yaml:"responses"`
	Callbacks    map[string]*Callback  `json:"callbacks" yaml:"callbacks"`
	Deprecated   bool                  `json:"deprecated" yaml:"deprecated"`
	Security     []map[string][]string `json:"security" yaml:"security"`
	Servers      []ServerInfo          `json:"servers" yaml:"servers"`
}

// Parameter represents an API parameter
type Parameter struct {
	Name            string       `json:"name" yaml:"name"`
	In              string       `json:"in" yaml:"in"` // query, header, path, cookie
	Description     string       `json:"description" yaml:"description"`
	Required        bool         `json:"required" yaml:"required"`
	Deprecated      bool         `json:"deprecated" yaml:"deprecated"`
	AllowEmptyValue bool         `json:"allowEmptyValue" yaml:"allowEmptyValue"`
	Style           string       `json:"style" yaml:"style"`
	Explode         bool         `json:"explode" yaml:"explode"`
	AllowReserved   bool         `json:"allowReserved" yaml:"allowReserved"`
	Schema          *SchemaObject `json:"schema" yaml:"schema"`
	Example         interface{}  `json:"example" yaml:"example"`
	Examples        map[string]*ExampleObject `json:"examples" yaml:"examples"`
}

// RequestBody represents a request body
type RequestBody struct {
	Description string                       `json:"description" yaml:"description"`
	Content     map[string]*MediaTypeObject  `json:"content" yaml:"content"`
	Required    bool                         `json:"required" yaml:"required"`
}

// Response represents an API response
type Response struct {
	Description string                      `json:"description" yaml:"description"`
	Headers     map[string]*HeaderObject    `json:"headers" yaml:"headers"`
	Content     map[string]*MediaTypeObject `json:"content" yaml:"content"`
	Links       map[string]*LinkObject      `json:"links" yaml:"links"`
}

// MediaTypeObject represents a media type
type MediaTypeObject struct {
	Schema   *SchemaObject             `json:"schema" yaml:"schema"`
	Example  interface{}               `json:"example" yaml:"example"`
	Examples map[string]*ExampleObject `json:"examples" yaml:"examples"`
	Encoding map[string]*EncodingObject `json:"encoding" yaml:"encoding"`
}

// SchemaObject represents a schema
type SchemaObject struct {
	Type                 string                    `json:"type" yaml:"type"`
	Format              string                    `json:"format" yaml:"format"`
	Title               string                    `json:"title" yaml:"title"`
	Description         string                    `json:"description" yaml:"description"`
	Default             interface{}               `json:"default" yaml:"default"`
	Example             interface{}               `json:"example" yaml:"example"`
	Enum                []interface{}             `json:"enum" yaml:"enum"`
	Required            []string                  `json:"required" yaml:"required"`
	Properties          map[string]*SchemaObject  `json:"properties" yaml:"properties"`
	Items               *SchemaObject             `json:"items" yaml:"items"`
	AdditionalProperties interface{}               `json:"additionalProperties" yaml:"additionalProperties"`
	AllOf               []*SchemaObject           `json:"allOf" yaml:"allOf"`
	OneOf               []*SchemaObject           `json:"oneOf" yaml:"oneOf"`
	AnyOf               []*SchemaObject           `json:"anyOf" yaml:"anyOf"`
	Not                 *SchemaObject             `json:"not" yaml:"not"`
	Discriminator       *DiscriminatorObject      `json:"discriminator" yaml:"discriminator"`
	ReadOnly            bool                      `json:"readOnly" yaml:"readOnly"`
	WriteOnly           bool                      `json:"writeOnly" yaml:"writeOnly"`
	Xml                 *XmlObject                `json:"xml" yaml:"xml"`
	ExternalDocs        *ExternalDocsInfo         `json:"externalDocs" yaml:"externalDocs"`
	Deprecated          bool                      `json:"deprecated" yaml:"deprecated"`
	Minimum             *float64                  `json:"minimum" yaml:"minimum"`
	Maximum             *float64                  `json:"maximum" yaml:"maximum"`
	ExclusiveMinimum    bool                      `json:"exclusiveMinimum" yaml:"exclusiveMinimum"`
	ExclusiveMaximum    bool                      `json:"exclusiveMaximum" yaml:"exclusiveMaximum"`
	MinLength           *int                      `json:"minLength" yaml:"minLength"`
	MaxLength           *int                      `json:"maxLength" yaml:"maxLength"`
	Pattern             string                    `json:"pattern" yaml:"pattern"`
	MinItems            *int                      `json:"minItems" yaml:"minItems"`
	MaxItems            *int                      `json:"maxItems" yaml:"maxItems"`
	UniqueItems         bool                      `json:"uniqueItems" yaml:"uniqueItems"`
	MinProperties       *int                      `json:"minProperties" yaml:"minProperties"`
	MaxProperties       *int                      `json:"maxProperties" yaml:"maxProperties"`
	MultipleOf          *float64                  `json:"multipleOf" yaml:"multipleOf"`
}

// ComponentsObject represents reusable components
type ComponentsObject struct {
	Schemas         map[string]*SchemaObject         `json:"schemas" yaml:"schemas"`
	Responses       map[string]*Response             `json:"responses" yaml:"responses"`
	Parameters      map[string]*Parameter            `json:"parameters" yaml:"parameters"`
	Examples        map[string]*ExampleObject        `json:"examples" yaml:"examples"`
	RequestBodies   map[string]*RequestBody          `json:"requestBodies" yaml:"requestBodies"`
	Headers         map[string]*HeaderObject         `json:"headers" yaml:"headers"`
	SecuritySchemes map[string]*SecuritySchemeObject `json:"securitySchemes" yaml:"securitySchemes"`
	Links           map[string]*LinkObject           `json:"links" yaml:"links"`
	Callbacks       map[string]*Callback             `json:"callbacks" yaml:"callbacks"`
}

// ExampleObject represents an example
type ExampleObject struct {
	Summary       string      `json:"summary" yaml:"summary"`
	Description   string      `json:"description" yaml:"description"`
	Value         interface{} `json:"value" yaml:"value"`
	ExternalValue string      `json:"externalValue" yaml:"externalValue"`
}

// HeaderObject represents a header
type HeaderObject struct {
	Description     string       `json:"description" yaml:"description"`
	Required        bool         `json:"required" yaml:"required"`
	Deprecated      bool         `json:"deprecated" yaml:"deprecated"`
	AllowEmptyValue bool         `json:"allowEmptyValue" yaml:"allowEmptyValue"`
	Schema          *SchemaObject `json:"schema" yaml:"schema"`
	Example         interface{}  `json:"example" yaml:"example"`
	Examples        map[string]*ExampleObject `json:"examples" yaml:"examples"`
}

// SecuritySchemeObject represents a security scheme
type SecuritySchemeObject struct {
	Type             string      `json:"type" yaml:"type"`
	Description      string      `json:"description" yaml:"description"`
	Name             string      `json:"name" yaml:"name"`
	In               string      `json:"in" yaml:"in"`
	Scheme           string      `json:"scheme" yaml:"scheme"`
	BearerFormat     string      `json:"bearerFormat" yaml:"bearerFormat"`
	Flows            OAuth2Flows `json:"flows" yaml:"flows"`
	OpenIdConnectUrl string      `json:"openIdConnectUrl" yaml:"openIdConnectUrl"`
}

// LinkObject represents a link
type LinkObject struct {
	OperationRef string                 `json:"operationRef" yaml:"operationRef"`
	OperationId  string                 `json:"operationId" yaml:"operationId"`
	Parameters   map[string]interface{} `json:"parameters" yaml:"parameters"`
	RequestBody  interface{}            `json:"requestBody" yaml:"requestBody"`
	Description  string                 `json:"description" yaml:"description"`
	Server       *ServerInfo            `json:"server" yaml:"server"`
}

// Callback represents a callback
type Callback map[string]*PathItem

// EncodingObject represents encoding
type EncodingObject struct {
	ContentType   string              `json:"contentType" yaml:"contentType"`
	Headers       map[string]*HeaderObject `json:"headers" yaml:"headers"`
	Style         string              `json:"style" yaml:"style"`
	Explode       bool                `json:"explode" yaml:"explode"`
	AllowReserved bool                `json:"allowReserved" yaml:"allowReserved"`
}

// DiscriminatorObject represents a discriminator
type DiscriminatorObject struct {
	PropertyName string            `json:"propertyName" yaml:"propertyName"`
	Mapping      map[string]string `json:"mapping" yaml:"mapping"`
}

// XmlObject represents XML metadata
type XmlObject struct {
	Name      string `json:"name" yaml:"name"`
	Namespace string `json:"namespace" yaml:"namespace"`
	Prefix    string `json:"prefix" yaml:"prefix"`
	Attribute bool   `json:"attribute" yaml:"attribute"`
	Wrapped   bool   `json:"wrapped" yaml:"wrapped"`
}

// DocumentationGeneratorOptions represents options for documentation generation
type DocumentationGeneratorOptions struct {
	ProjectName        string                  `json:"project_name" yaml:"project_name"`
	OutputPath         string                  `json:"output_path" yaml:"output_path"`
	Config             APIDocumentationConfig  `json:"config" yaml:"config"`
	IncludeAuth        bool                    `json:"include_auth" yaml:"include_auth"`
	AuthConfig         *AuthConfig             `json:"auth_config" yaml:"auth_config"`
	IncludeModels      bool                    `json:"include_models" yaml:"include_models"`
	Models             []ModelDefinition       `json:"models" yaml:"models"`
	IncludeEndpoints   bool                    `json:"include_endpoints" yaml:"include_endpoints"`
	Endpoints          []EndpointDefinition    `json:"endpoints" yaml:"endpoints"`
	CustomSchemas      map[string]*SchemaObject `json:"custom_schemas" yaml:"custom_schemas"`
	GenerateExamples   bool                    `json:"generate_examples" yaml:"generate_examples"`
	GenerateSwaggerUI  bool                    `json:"generate_swagger_ui" yaml:"generate_swagger_ui"`
	SwaggerUIConfig    SwaggerUIConfig         `json:"swagger_ui_config" yaml:"swagger_ui_config"`
}

// ModelDefinition represents a model for documentation
type ModelDefinition struct {
	Name         string                   `json:"name" yaml:"name"`
	Description  string                   `json:"description" yaml:"description"`
	Fields       []Field                  `json:"fields" yaml:"fields"`
	Examples     map[string]interface{}   `json:"examples" yaml:"examples"`
	Validations  []ValidationRule         `json:"validations" yaml:"validations"`
	Relationships []RelationshipDefinition `json:"relationships" yaml:"relationships"`
}

// EndpointDefinition represents an endpoint for documentation
type EndpointDefinition struct {
	Method      string            `json:"method" yaml:"method"`
	Path        string            `json:"path" yaml:"path"`
	Summary     string            `json:"summary" yaml:"summary"`
	Description string            `json:"description" yaml:"description"`
	Tags        []string          `json:"tags" yaml:"tags"`
	Parameters  []Parameter       `json:"parameters" yaml:"parameters"`
	RequestBody *RequestBody      `json:"request_body" yaml:"request_body"`
	Responses   map[string]*Response `json:"responses" yaml:"responses"`
	Security    []string          `json:"security" yaml:"security"`
	Deprecated  bool              `json:"deprecated" yaml:"deprecated"`
	Examples    EndpointExamples  `json:"examples" yaml:"examples"`
}

// EndpointExamples represents examples for an endpoint
type EndpointExamples struct {
	Request  map[string]interface{} `json:"request" yaml:"request"`
	Response map[string]interface{} `json:"response" yaml:"response"`
	Headers  map[string]string      `json:"headers" yaml:"headers"`
	Curl     string                 `json:"curl" yaml:"curl"`
}

// RelationshipDefinition represents model relationships
type RelationshipDefinition struct {
	Name         string `json:"name" yaml:"name"`
	Type         string `json:"type" yaml:"type"` // belongs_to, has_one, has_many
	TargetModel  string `json:"target_model" yaml:"target_model"`
	ForeignKey   string `json:"foreign_key" yaml:"foreign_key"`
	Description  string `json:"description" yaml:"description"`
}

// SwaggerUIConfig represents Swagger UI configuration
type SwaggerUIConfig struct {
	Title                string            `json:"title" yaml:"title"`
	Theme                string            `json:"theme" yaml:"theme"`
	DeepLinking          bool              `json:"deep_linking" yaml:"deep_linking"`
	DisplayOperationId   bool              `json:"display_operation_id" yaml:"display_operation_id"`
	DefaultModelsExpandDepth int           `json:"default_models_expand_depth" yaml:"default_models_expand_depth"`
	DefaultModelExpandDepth  int           `json:"default_model_expand_depth" yaml:"default_model_expand_depth"`
	DocExpansion         string            `json:"doc_expansion" yaml:"doc_expansion"`
	Filter               bool              `json:"filter" yaml:"filter"`
	ShowExtensions       bool              `json:"show_extensions" yaml:"show_extensions"`
	ShowCommonExtensions bool              `json:"show_common_extensions" yaml:"show_common_extensions"`
	TryItOutEnabled      bool              `json:"try_it_out_enabled" yaml:"try_it_out_enabled"`
	CustomCSS            string            `json:"custom_css" yaml:"custom_css"`
	CustomJS             string            `json:"custom_js" yaml:"custom_js"`
	OAuth2Config         *OAuth2UIConfig   `json:"oauth2_config" yaml:"oauth2_config"`
}

// OAuth2UIConfig represents OAuth2 configuration for Swagger UI
type OAuth2UIConfig struct {
	ClientId     string            `json:"clientId" yaml:"clientId"`
	ClientSecret string            `json:"clientSecret" yaml:"clientSecret"`
	Realm        string            `json:"realm" yaml:"realm"`
	AppName      string            `json:"appName" yaml:"appName"`
	ScopeSeparator string          `json:"scopeSeparator" yaml:"scopeSeparator"`
	Scopes       []string          `json:"scopes" yaml:"scopes"`
	AdditionalQueryStringParams map[string]string `json:"additionalQueryStringParams" yaml:"additionalQueryStringParams"`
	UseBasicAuthenticationWithAccessCodeGrant bool `json:"useBasicAuthenticationWithAccessCodeGrant" yaml:"useBasicAuthenticationWithAccessCodeGrant"`
	UsePkceWithAuthorizationCodeGrant bool `json:"usePkceWithAuthorizationCodeGrant" yaml:"usePkceWithAuthorizationCodeGrant"`
}

// DefaultAPIDocumentationConfig returns default configuration for API documentation
func DefaultAPIDocumentationConfig(projectName string) APIDocumentationConfig {
	return APIDocumentationConfig{
		ProjectName:     projectName,
		Version:         "1.0.0",
		Description:     fmt.Sprintf("API documentation for %s", projectName),
		Contact: ContactInfo{
			Name:  "API Team",
			Email: "api@example.com",
		},
		License: LicenseInfo{
			Name: "MIT",
			URL:  "https://opensource.org/licenses/MIT",
		},
		Servers: []ServerInfo{
			{
				URL:         "http://localhost:8080",
				Description: "Development server",
			},
		},
		BasePath:        "/api/v1",
		EnableSwaggerUI: true,
		SwaggerUIPath:   "/docs",
		OutputFormat:    []DocumentFormat{FormatJSON, FormatYAML, FormatHTML},
		Security: SecuritySchemes{
			JWT: &JWTSecurityScheme{
				Type:         "http",
				Scheme:       "bearer",
				BearerFormat: "JWT",
				Description:  "JWT Authorization header using the Bearer scheme",
			},
		},
		Tags: []Tag{
			{
				Name:        "auth",
				Description: "Authentication operations",
			},
			{
				Name:        "users",
				Description: "User management operations",
			},
		},
	}
}

// DefaultSwaggerUIConfig returns default Swagger UI configuration
func DefaultSwaggerUIConfig(projectName string) SwaggerUIConfig {
	return SwaggerUIConfig{
		Title:                    fmt.Sprintf("%s API Documentation", projectName),
		Theme:                    "default",
		DeepLinking:              true,
		DisplayOperationId:       false,
		DefaultModelsExpandDepth: 1,
		DefaultModelExpandDepth:  1,
		DocExpansion:            "list",
		Filter:                  true,
		ShowExtensions:          false,
		ShowCommonExtensions:    false,
		TryItOutEnabled:         true,
	}
}

// GetGoTypeForOpenAPI returns OpenAPI type for Go field type
func GetGoTypeForOpenAPI(fieldType FieldType) (string, string) {
	switch fieldType {
	case FieldTypeString, FieldTypeText, FieldTypeSlug:
		return "string", ""
	case FieldTypeEmail:
		return "string", "email"
	case FieldTypeURL:
		return "string", "uri"
	case FieldTypePassword:
		return "string", "password"
	case FieldTypePhone:
		return "string", ""
	case FieldTypeColor:
		return "string", ""
	case FieldTypeNumber:
		return "integer", "int64"
	case FieldTypeFloat:
		return "number", "float"
	case FieldTypeCurrency:
		return "number", "double"
	case FieldTypeBoolean:
		return "boolean", ""
	case FieldTypeDate:
		return "string", "date-time"
	case FieldTypeUUID:
		return "string", "uuid"
	case FieldTypeJSON:
		return "object", ""
	case FieldTypeFile, FieldTypeImage:
		return "string", "binary"
	case FieldTypeCoordinates:
		return "object", ""
	case FieldTypeEnum:
		return "string", ""
	case FieldTypeRelation, FieldTypeRelationArray:
		return "object", ""
	default:
		return "string", ""
	}
}

// GenerateSchemaFromField generates OpenAPI schema from field definition
func GenerateSchemaFromField(field Field) *SchemaObject {
	openAPIType, format := GetGoTypeForOpenAPI(field.Type)
	
	schema := &SchemaObject{
		Type:        openAPIType,
		Format:      format,
		Title:       field.DisplayName,
		Description: field.Description,
	}
	
	// Add validation constraints
	if field.MinLength != nil {
		schema.MinLength = field.MinLength
	}
	if field.MaxLength != nil {
		schema.MaxLength = field.MaxLength
	}
	if field.Pattern != "" {
		schema.Pattern = field.Pattern
	}
	if field.MinValue != nil {
		schema.Minimum = field.MinValue
	}
	if field.MaxValue != nil {
		schema.Maximum = field.MaxValue
	}
	
	// Handle special field types
	switch field.Type {
	case FieldTypeEnum:
		if len(field.EnumValues) > 0 {
			enumValues := make([]interface{}, len(field.EnumValues))
			for i, v := range field.EnumValues {
				enumValues[i] = v
			}
			schema.Enum = enumValues
		}
	case FieldTypeCoordinates:
		schema = &SchemaObject{
			Type: "object",
			Properties: map[string]*SchemaObject{
				"latitude": {
					Type:   "number",
					Format: "double",
				},
				"longitude": {
					Type:   "number",
					Format: "double",
				},
			},
			Required: []string{"latitude", "longitude"},
		}
	}
	
	// Set default value
	if field.DefaultValue != nil {
		schema.Default = field.DefaultValue
	}
	
	return schema
}

// GenerateExampleFromField generates example value from field definition
func GenerateExampleFromField(field Field) interface{} {
	if field.DefaultValue != nil {
		return field.DefaultValue
	}
	
	// Generate example based on field type
	switch field.Type {
	case FieldTypeString:
		return "example string"
	case FieldTypeText:
		return "This is an example text field with more content"
	case FieldTypeEmail:
		return "user@example.com"
	case FieldTypeURL:
		return "https://example.com"
	case FieldTypePassword:
		return "********"
	case FieldTypePhone:
		return "+1234567890"
	case FieldTypeSlug:
		return "example-slug"
	case FieldTypeColor:
		return "#FF5733"
	case FieldTypeNumber:
		return 42
	case FieldTypeFloat:
		return 3.14
	case FieldTypeCurrency:
		return 99.99
	case FieldTypeBoolean:
		return true
	case FieldTypeDate:
		return time.Now().Format(time.RFC3339)
	case FieldTypeUUID:
		return "550e8400-e29b-41d4-a716-446655440000"
	case FieldTypeJSON:
		return map[string]interface{}{"key": "value"}
	case FieldTypeCoordinates:
		return map[string]interface{}{
			"latitude":  40.7128,
			"longitude": -74.0060,
		}
	case FieldTypeEnum:
		if len(field.EnumValues) > 0 {
			return field.EnumValues[0]
		}
		return "option1"
	default:
		return "example value"
	}
}