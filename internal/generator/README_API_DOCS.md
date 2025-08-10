# API Documentation Generator

## Overview

The Vibercode CLI API Documentation Generator provides comprehensive OpenAPI/Swagger specification generation with interactive documentation, code examples, and automatic endpoint documentation for Go APIs.

## Features

- ✅ **OpenAPI 3.0 Specification**: Complete OpenAPI 3.0 compliant specification generation
- ✅ **Swagger UI Integration**: Interactive API documentation with try-it-now functionality
- ✅ **Model Schema Generation**: Automatic schema generation from Go structs and field definitions
- ✅ **Endpoint Documentation**: Automatic CRUD endpoint documentation with examples
- ✅ **Authentication Documentation**: JWT, OAuth2, and API Key authentication documentation
- ✅ **Multiple Export Formats**: JSON, YAML, HTML, and PDF export capabilities
- ✅ **Request/Response Examples**: Automatic example generation based on field types
- ✅ **Security Schemes**: Multiple authentication method documentation
- ✅ **Custom Validation**: Integration with field validation rules and constraints

## Quick Start

### Basic API Documentation

```go
package main

import (
    "github.com/vibercode/cli/internal/generator"
    "github.com/vibercode/cli/internal/models"
)

func main() {
    // Create default documentation options
    opts := generator.DefaultDocumentationGeneratorOptions("my-api", "./docs-output")
    
    // Generate API documentation
    docGen := generator.NewAPIDocsGenerator(opts)
    if err := docGen.GenerateAPIDocs(); err != nil {
        panic(err)
    }
}
```

### With Custom Models

```go
// Add custom models to documentation
opts := generator.DefaultDocumentationGeneratorOptions("my-api", "./docs-output")

opts.Models = []models.ModelDefinition{
    {
        Name:        "User",
        Description: "User account information",
        Fields: []models.Field{
            {
                Name:        "id",
                Type:        models.FieldTypeNumber,
                DisplayName: "User ID",
                Description: "Unique identifier for the user",
                Required:    true,
            },
            {
                Name:        "email",
                Type:        models.FieldTypeEmail,
                DisplayName: "Email Address",
                Description: "User's email address for login and notifications",
                Required:    true,
            },
            {
                Name:        "profile",
                Type:        models.FieldTypeJSON,
                DisplayName: "User Profile",
                Description: "Additional user profile information",
                Required:    false,
            },
        },
        Examples: map[string]interface{}{
            "basic_user": map[string]interface{}{
                "id":      1,
                "email":   "user@example.com",
                "profile": map[string]interface{}{"name": "John Doe"},
            },
        },
    },
}

docGen := generator.NewAPIDocsGenerator(opts)
err := docGen.GenerateAPIDocs()
```

### With Custom Endpoints

```go
// Add custom API endpoints
opts.Endpoints = []models.EndpointDefinition{
    {
        Method:      "GET",
        Path:        "/api/users",
        Summary:     "List Users",
        Description: "Retrieve a paginated list of users with optional filtering",
        Tags:        []string{"Users", "Management"},
        Parameters: []models.Parameter{
            {
                Name:        "page",
                In:          "query",
                Description: "Page number for pagination",
                Required:    false,
                Schema: &models.SchemaObject{
                    Type:    "integer",
                    Default: 1,
                    Minimum: &[]float64{1}[0],
                },
            },
            {
                Name:        "search",
                In:          "query",
                Description: "Search term to filter users",
                Required:    false,
                Schema: &models.SchemaObject{
                    Type: "string",
                },
            },
        },
        Responses: map[string]*models.Response{
            "200": {
                Description: "Successful response with user list",
                Content: map[string]*models.MediaTypeObject{
                    "application/json": {
                        Schema: &models.SchemaObject{
                            Type: "object",
                            Properties: map[string]*models.SchemaObject{
                                "data": {
                                    Type: "array",
                                    Items: &models.SchemaObject{
                                        Type: "object", // Reference to User schema
                                    },
                                },
                                "total":     {Type: "integer"},
                                "page":      {Type: "integer"},
                                "per_page":  {Type: "integer"},
                            },
                        },
                    },
                },
            },
        },
        Examples: models.EndpointExamples{
            Request: map[string]interface{}{
                "query_params": map[string]interface{}{
                    "page":   1,
                    "search": "john",
                },
            },
            Response: map[string]interface{}{
                "data": []interface{}{
                    map[string]interface{}{
                        "id":    1,
                        "email": "john@example.com",
                    },
                },
                "total":    1,
                "page":     1,
                "per_page": 20,
            },
            Curl: `curl -X GET "http://localhost:8080/api/users?page=1&search=john" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"`,
        },
    },
}
```

### With Authentication Documentation

```go
// Configure authentication documentation
opts.IncludeAuth = true
opts.AuthConfig = &models.AuthConfig{
    Provider: models.AuthProviderJWT,
    Methods:  []models.AuthMethod{models.AuthMethodEmail, models.AuthMethodUsername},
    OAuth2Providers: []models.OAuth2Provider{
        {
            Name:         "google",
            ClientID:     "your-google-client-id",
            ClientSecret: "your-google-client-secret",
            RedirectURL:  "http://localhost:8080/auth/google/callback",
            Scopes:       []string{"openid", "email", "profile"},
            AuthURL:      "https://accounts.google.com/o/oauth2/auth",
            TokenURL:     "https://oauth2.googleapis.com/token",
            UserInfoURL:  "https://www.googleapis.com/oauth2/v2/userinfo",
        },
        {
            Name:         "github",
            ClientID:     "your-github-client-id",
            ClientSecret: "your-github-client-secret",
            RedirectURL:  "http://localhost:8080/auth/github/callback",
            Scopes:       []string{"user:email"},
            AuthURL:      "https://github.com/login/oauth/authorize",
            TokenURL:     "https://github.com/login/oauth/access_token",
            UserInfoURL:  "https://api.github.com/user",
        },
    },
}
```

## Generated Structure

The API documentation generator creates a comprehensive documentation system:

```
docs-output/
├── docs/
│   ├── openapi.json           # OpenAPI 3.0 specification (JSON)
│   ├── openapi.yaml           # OpenAPI 3.0 specification (YAML)
│   └── swagger-ui/
│       └── index.html         # Interactive Swagger UI
├── internal/
│   ├── middleware/
│   │   └── docs.go            # Documentation middleware
│   ├── handlers/
│   │   ├── docs.go            # Documentation endpoints
│   │   └── endpoint_docs.go   # Generated endpoint documentation
│   └── models/
│       └── schemas.go         # Model schema definitions
└── go.mod
```

## Configuration Options

### DocumentationGeneratorOptions

```go
type DocumentationGeneratorOptions struct {
    ProjectName        string                  // API project name
    OutputPath         string                  // Output directory path
    Config             APIDocumentationConfig  // Documentation configuration
    IncludeAuth        bool                    // Include authentication endpoints
    AuthConfig         *AuthConfig             // Authentication configuration
    IncludeModels      bool                    // Include model schemas
    Models             []ModelDefinition       // Model definitions
    IncludeEndpoints   bool                    // Include API endpoints
    Endpoints          []EndpointDefinition    // Endpoint definitions
    CustomSchemas      map[string]*SchemaObject // Custom schema definitions
    GenerateExamples   bool                    // Generate request/response examples
    GenerateSwaggerUI  bool                    // Generate Swagger UI interface
    SwaggerUIConfig    SwaggerUIConfig         // Swagger UI configuration
}
```

### APIDocumentationConfig

```go
type APIDocumentationConfig struct {
    ProjectName     string             // Project name
    Version         string             // API version
    Description     string             // API description
    Contact         ContactInfo        // Contact information
    License         LicenseInfo        // License information
    Servers         []ServerInfo       // Server configurations
    BasePath        string             // API base path
    EnableSwaggerUI bool               // Enable Swagger UI
    SwaggerUIPath   string             // Swagger UI endpoint path
    OutputFormat    []DocumentFormat   // Output formats (JSON, YAML, HTML, PDF)
    Security        SecuritySchemes    // Security configurations
    Tags            []Tag              // API tags for grouping
    ExternalDocs    *ExternalDocsInfo  // External documentation links
}
```

### SwaggerUIConfig

```go
type SwaggerUIConfig struct {
    Title                    string            // UI title
    Theme                    string            // UI theme
    DeepLinking              bool              // Enable deep linking
    DisplayOperationId       bool              // Show operation IDs
    DefaultModelsExpandDepth int               // Default model expansion depth
    DefaultModelExpandDepth  int               // Default model property depth
    DocExpansion             string            // Documentation expansion mode
    Filter                   bool              // Enable filtering
    ShowExtensions           bool              // Show spec extensions
    ShowCommonExtensions     bool              // Show common extensions
    TryItOutEnabled          bool              // Enable try-it-out functionality
    CustomCSS                string            // Custom CSS styles
    CustomJS                 string            // Custom JavaScript
    OAuth2Config             *OAuth2UIConfig   // OAuth2 UI configuration
}
```

## Field Type Integration

The documentation generator integrates seamlessly with the enhanced field types system:

### Supported Field Types in Documentation

| Field Type | OpenAPI Type | Format | Validation Support |
|------------|--------------|--------|-------------------|
| `string` | string | - | min/max length, pattern |
| `text` | string | - | min/max length |
| `email` | string | email | email validation |
| `url` | string | uri | URL validation |
| `password` | string | password | min length |
| `phone` | string | - | phone format |
| `slug` | string | - | slug pattern |
| `color` | string | - | hex color pattern |
| `number` | integer | int64 | min/max value |
| `float` | number | float | min/max value |
| `currency` | number | double | min/max value |
| `boolean` | boolean | - | - |
| `date` | string | date-time | - |
| `uuid` | string | uuid | - |
| `json` | object | - | - |
| `coordinates` | object | - | lat/lng validation |
| `file` | string | binary | - |
| `image` | string | binary | - |
| `enum` | string | - | enum values |
| `relation` | integer | int64 | - |
| `relation-array` | array | - | array of integers |

### Example Field Documentation Generation

```go
// Field definition
field := models.Field{
    Name:        "email",
    Type:        models.FieldTypeEmail,
    DisplayName: "Email Address",
    Description: "User's primary email address",
    Required:    true,
    MinLength:   &[]int{5}[0],
    MaxLength:   &[]int{100}[0],
    Pattern:     `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
}

// Generated OpenAPI schema
{
  "type": "string",
  "format": "email",
  "title": "Email Address",
  "description": "User's primary email address",
  "minLength": 5,
  "maxLength": 100,
  "pattern": "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
  "example": "user@example.com"
}
```

## Security Documentation

### JWT Authentication

```go
// JWT configuration generates this security scheme:
{
  "bearerAuth": {
    "type": "http",
    "scheme": "bearer",
    "bearerFormat": "JWT",
    "description": "JWT Authorization header using the Bearer scheme"
  }
}
```

### OAuth2 Authentication

```go
// OAuth2 configuration generates this security scheme:
{
  "oauth2": {
    "type": "oauth2",
    "description": "OAuth2 authentication",
    "flows": {
      "authorizationCode": {
        "authorizationUrl": "https://accounts.google.com/o/oauth2/auth",
        "tokenUrl": "https://oauth2.googleapis.com/token",
        "refreshUrl": "https://oauth2.googleapis.com/token",
        "scopes": {
          "openid": "OpenID Connect",
          "email": "Access email address",
          "profile": "Access profile information"
        }
      }
    }
  }
}
```

### API Key Authentication

```go
// API Key configuration
{
  "apiKeyAuth": {
    "type": "apiKey",
    "in": "header",
    "name": "X-API-Key",
    "description": "API key authentication"
  }
}
```

## Middleware Integration

### Documentation Middleware

The generated middleware provides endpoints for serving documentation:

```go
// Setup documentation routes
func SetupDocsRoutes(mux *http.ServeMux, config *models.APIDocumentationConfig) {
    docsHandler := handlers.NewDocsHandler(config)
    
    // OpenAPI spec endpoints
    mux.HandleFunc("GET /api/docs/openapi.json", docsHandler.GetOpenAPISpec)
    mux.HandleFunc("GET /api/docs/openapi.yaml", docsHandler.GetOpenAPISpec)
    
    // Swagger UI endpoint
    mux.HandleFunc("GET /api/docs", docsHandler.GetSwaggerUI)
    
    // Health check
    mux.HandleFunc("GET /api/docs/health", docsHandler.HealthCheck)
}
```

### CORS Support

The documentation middleware includes CORS headers for cross-origin access:

```go
// CORS headers for API specification
w.Header().Set("Access-Control-Allow-Origin", "*")
w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
```

## Advanced Features

### Custom Schema Definitions

```go
// Add custom schemas to the documentation
opts.CustomSchemas = map[string]*models.SchemaObject{
    "PaginationMetadata": {
        Type: "object",
        Properties: map[string]*models.SchemaObject{
            "total":       {Type: "integer", Description: "Total number of items"},
            "page":        {Type: "integer", Description: "Current page number"},
            "per_page":    {Type: "integer", Description: "Items per page"},
            "total_pages": {Type: "integer", Description: "Total number of pages"},
        },
    },
    "ErrorResponse": {
        Type: "object",
        Properties: map[string]*models.SchemaObject{
            "error":   {Type: "string", Description: "Error message"},
            "code":    {Type: "string", Description: "Error code"},
            "details": {Type: "object", Description: "Additional error details"},
        },
    },
}
```

### Example Generation

The system automatically generates examples based on field types:

```go
// Automatic example generation for different field types
examples := map[string]interface{}{
    "string":      "example string",
    "text":        "This is an example text field with more content",
    "email":       "user@example.com",
    "url":         "https://example.com",
    "slug":        "example-slug",
    "color":       "#FF5733",
    "number":      42,
    "float":       3.14,
    "currency":    99.99,
    "boolean":     true,
    "date":        "2024-01-15T10:30:00Z",
    "uuid":        "550e8400-e29b-41d4-a716-446655440000",
    "coordinates": map[string]interface{}{"latitude": 40.7128, "longitude": -74.0060},
}
```

### Multi-Format Export

```go
// Configure multiple export formats
opts.Config.OutputFormat = []models.DocumentFormat{
    models.FormatJSON,  // OpenAPI JSON specification
    models.FormatYAML,  // OpenAPI YAML specification
    models.FormatHTML,  // Swagger UI HTML
    models.FormatPDF,   // PDF documentation (if enabled)
}
```

## Testing

The documentation generator includes comprehensive test coverage:

### Unit Tests

```bash
# Run documentation generator tests
go test ./internal/generator -run TestAPIDocsGenerator -v

# Run specific test suites
go test ./internal/generator -run TestGenerateOpenAPISpec -v
go test ./internal/generator -run TestAddModelSchemas -v
go test ./internal/generator -run TestSwaggerUIGeneration -v
```

### Integration Tests

```bash
# Test complete documentation generation
go test ./internal/generator -run TestAPIDocsGenerator -v

# Test with authentication
go test ./internal/generator -run TestAPIDocsWithAuthentication -v

# Benchmark documentation generation
go test ./internal/generator -bench=BenchmarkAPIDocsGenerator -v
```

### Manual Testing

```bash
# Generate test API documentation
./vibercode generate docs \
  --project-name test-api \
  --output ./test-docs \
  --include-auth \
  --include-models \
  --generate-swagger-ui

# Serve documentation locally
cd test-docs
go run cmd/server/main.go

# Access documentation at http://localhost:8080/docs
```

## Best Practices

### Documentation Quality

1. **Clear Descriptions**: Provide detailed descriptions for all models, fields, and endpoints
2. **Meaningful Examples**: Include realistic examples that demonstrate actual usage
3. **Proper Tagging**: Use consistent tags to organize endpoints logically
4. **Version Management**: Keep API versions synchronized with documentation
5. **Security Documentation**: Clearly document authentication requirements

### Performance Optimization

1. **Spec Caching**: Cache generated OpenAPI specifications for better performance
2. **Lazy Loading**: Generate documentation on-demand rather than at startup
3. **Compression**: Enable gzip compression for documentation endpoints
4. **CDN Usage**: Serve static assets (Swagger UI) from CDN in production

### Maintenance

1. **Automated Generation**: Integrate documentation generation into CI/CD pipeline
2. **Validation**: Validate OpenAPI specifications against the official schema
3. **Testing**: Include documentation generation in automated tests
4. **Monitoring**: Monitor documentation endpoint usage and performance

### Development Workflow

1. **Code-First Approach**: Generate documentation from code annotations and models
2. **Continuous Updates**: Regenerate documentation on model or endpoint changes
3. **Review Process**: Include documentation reviews in code review process
4. **Feedback Integration**: Collect and integrate feedback from API consumers

## Troubleshooting

### Common Issues

1. **Missing Schemas**
   - Ensure models are properly defined in the configuration
   - Check field type mappings for custom types
   - Verify required imports in generated code

2. **Swagger UI Not Loading**
   - Check CORS configuration for cross-origin requests
   - Verify Swagger UI assets are accessible
   - Validate OpenAPI specification JSON syntax

3. **Authentication Documentation Issues**
   - Ensure auth configuration is properly set
   - Check security scheme definitions
   - Verify OAuth2 flow configurations

4. **Example Generation Problems**
   - Check field default values and example values
   - Verify enum value definitions
   - Ensure field validation rules are properly set

### Debug Mode

Enable detailed logging for troubleshooting:

```go
opts := generator.DefaultDocumentationGeneratorOptions("debug-api", "./debug-docs")
opts.Config.Debug = true // Enable debug logging

// Generated code will include detailed logging
```

## Contributing

To contribute to the API documentation generator:

1. **Add New Features**: Extend the generator for additional OpenAPI 3.0 features
2. **Improve Templates**: Enhance Swagger UI templates and customization options
3. **Add Export Formats**: Implement additional export formats (PDF, Markdown)
4. **Enhance Examples**: Improve automatic example generation algorithms
5. **Testing**: Add more comprehensive test coverage for edge cases

## API Reference

### Core Types

- `APIDocumentationConfig` - Main documentation configuration
- `DocumentationGeneratorOptions` - Generator options and settings
- `OpenAPISpec` - Complete OpenAPI 3.0 specification structure
- `SwaggerUIConfig` - Swagger UI customization options
- `ModelDefinition` - Model definition for schema generation
- `EndpointDefinition` - API endpoint definition for path generation

### Key Functions

```go
// Create new documentation generator
func NewAPIDocsGenerator(options models.DocumentationGeneratorOptions) *APIDocsGenerator

// Generate complete API documentation
func (g *APIDocsGenerator) GenerateAPIDocs() error

// Configuration helpers
func DefaultDocumentationGeneratorOptions(projectName, outputPath string) models.DocumentationGeneratorOptions
func DefaultAPIDocumentationConfig(projectName string) models.APIDocumentationConfig
func DefaultSwaggerUIConfig(projectName string) models.SwaggerUIConfig
```

This API documentation generator provides a complete solution for creating professional, interactive API documentation that integrates seamlessly with the Vibercode CLI's field type system and authentication framework.