package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vibercode/cli/internal/models"
)

// TestAPIDocsGenerator tests the API documentation generator
func TestAPIDocsGenerator(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()

	// Create test options
	opts := DefaultDocumentationGeneratorOptions("test-api", tempDir)
	
	// Add test models
	opts.Models = []models.ModelDefinition{
		{
			Name:        "User",
			Description: "User model for testing",
			Fields: []models.Field{
				{
					Name:        "id",
					Type:        models.FieldTypeNumber,
					DisplayName: "ID",
					Description: "User ID",
					Required:    true,
				},
				{
					Name:        "email",
					Type:        models.FieldTypeEmail,
					DisplayName: "Email",
					Description: "User email address",
					Required:    true,
				},
				{
					Name:        "name",
					Type:        models.FieldTypeString,
					DisplayName: "Name",
					Description: "User full name",
					Required:    false,
				},
			},
		},
	}
	
	// Add test endpoints
	opts.Endpoints = []models.EndpointDefinition{
		{
			Method:      "GET",
			Path:        "/users",
			Summary:     "List users",
			Description: "Get a list of all users",
			Tags:        []string{"Users"},
		},
		{
			Method:      "POST",
			Path:        "/users",
			Summary:     "Create user",
			Description: "Create a new user",
			Tags:        []string{"Users"},
		},
	}

	// Create generator
	generator := NewAPIDocsGenerator(opts)

	t.Run("GenerateAPIDocs", func(t *testing.T) {
		err := generator.GenerateAPIDocs()
		if err != nil {
			t.Errorf("Failed to generate API docs: %v", err)
		}

		// Check if OpenAPI spec was generated
		specPath := filepath.Join(tempDir, "docs", "openapi.json")
		if _, err := os.Stat(specPath); os.IsNotExist(err) {
			t.Error("OpenAPI spec file was not created")
		}

		// Check if YAML spec was generated
		yamlPath := filepath.Join(tempDir, "docs", "openapi.yaml")
		if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
			t.Error("OpenAPI YAML spec file was not created")
		}

		// Check if Swagger UI was generated
		swaggerPath := filepath.Join(tempDir, "docs", "swagger-ui", "index.html")
		if _, err := os.Stat(swaggerPath); os.IsNotExist(err) {
			t.Error("Swagger UI file was not created")
		}

		// Check if middleware was generated
		middlewarePath := filepath.Join(tempDir, "internal", "middleware", "docs.go")
		if _, err := os.Stat(middlewarePath); os.IsNotExist(err) {
			t.Error("Docs middleware file was not created")
		}

		// Check if handlers were generated
		handlerPath := filepath.Join(tempDir, "internal", "handlers", "docs.go")
		if _, err := os.Stat(handlerPath); os.IsNotExist(err) {
			t.Error("Docs handler file was not created")
		}
	})

	t.Run("ValidateOpenAPISpec", func(t *testing.T) {
		specPath := filepath.Join(tempDir, "docs", "openapi.json")
		specData, err := os.ReadFile(specPath)
		if err != nil {
			t.Fatalf("Failed to read OpenAPI spec: %v", err)
		}

		var spec models.OpenAPISpec
		if err := json.Unmarshal(specData, &spec); err != nil {
			t.Errorf("Failed to parse OpenAPI spec: %v", err)
		}

		// Validate basic spec structure
		if spec.OpenAPI != models.OpenAPIVersion {
			t.Errorf("Expected OpenAPI version %s, got %s", models.OpenAPIVersion, spec.OpenAPI)
		}

		if spec.Info.Title != "test-api API" {
			t.Errorf("Expected title 'test-api API', got '%s'", spec.Info.Title)
		}

		// Check if schemas were added
		if spec.Components == nil || spec.Components.Schemas == nil {
			t.Error("Components.Schemas is missing")
		} else {
			if _, exists := spec.Components.Schemas["User"]; !exists {
				t.Error("User schema is missing from components")
			}
		}

		// Check if paths were added
		if spec.Paths == nil {
			t.Error("Paths are missing")
		} else {
			if _, exists := spec.Paths["/users"]; !exists {
				t.Error("Users path is missing")
			}
		}
	})
}

// TestGenerateOpenAPISpec tests OpenAPI spec generation
func TestGenerateOpenAPISpec(t *testing.T) {
	tempDir := t.TempDir()
	opts := DefaultDocumentationGeneratorOptions("spec-test", tempDir)
	
	generator := NewAPIDocsGenerator(opts)
	
	t.Run("GenerateBasicSpec", func(t *testing.T) {
		err := generator.generateOpenAPISpec()
		if err != nil {
			t.Errorf("Failed to generate OpenAPI spec: %v", err)
		}

		// Check if spec was created
		if generator.spec == nil {
			t.Error("OpenAPI spec was not created")
		}

		// Validate spec structure
		if generator.spec.OpenAPI != models.OpenAPIVersion {
			t.Errorf("Expected OpenAPI version %s, got %s", models.OpenAPIVersion, generator.spec.OpenAPI)
		}

		if generator.spec.Info.Title != "spec-test API" {
			t.Errorf("Expected title 'spec-test API', got '%s'", generator.spec.Info.Title)
		}
	})
}

// TestAddSecuritySchemes tests security scheme addition
func TestAddSecuritySchemes(t *testing.T) {
	tempDir := t.TempDir()
	opts := DefaultDocumentationGeneratorOptions("security-test", tempDir)
	
	// Configure JWT security
	opts.Config.Security.JWT = &models.JWTSecurityScheme{
		Type:         "http",
		Scheme:       "bearer",
		BearerFormat: "JWT",
		Description:  "JWT Authorization header",
	}

	generator := NewAPIDocsGenerator(opts)
	generator.spec = &models.OpenAPISpec{
		Components: &models.ComponentsObject{
			SecuritySchemes: make(map[string]*models.SecuritySchemeObject),
		},
	}

	t.Run("AddJWTSecurityScheme", func(t *testing.T) {
		generator.addSecuritySchemes()

		if generator.spec.Components.SecuritySchemes["bearerAuth"] == nil {
			t.Error("Bearer auth security scheme was not added")
		}

		bearerAuth := generator.spec.Components.SecuritySchemes["bearerAuth"]
		if bearerAuth.Type != "http" {
			t.Errorf("Expected type 'http', got '%s'", bearerAuth.Type)
		}

		if bearerAuth.Scheme != "bearer" {
			t.Errorf("Expected scheme 'bearer', got '%s'", bearerAuth.Scheme)
		}
	})
}

// TestAddModelSchemas tests model schema addition
func TestAddModelSchemas(t *testing.T) {
	tempDir := t.TempDir()
	opts := DefaultDocumentationGeneratorOptions("model-test", tempDir)
	
	// Add test model
	opts.Models = []models.ModelDefinition{
		{
			Name:        "Product",
			Description: "Product model",
			Fields: []models.Field{
				{
					Name:        "id",
					Type:        models.FieldTypeNumber,
					DisplayName: "ID",
					Required:    true,
				},
				{
					Name:        "name",
					Type:        models.FieldTypeString,
					DisplayName: "Product Name",
					Required:    true,
					MinLength:   &[]int{1}[0],
					MaxLength:   &[]int{100}[0],
				},
				{
					Name:        "price",
					Type:        models.FieldTypeCurrency,
					DisplayName: "Price",
					Required:    true,
					MinValue:    &[]float64{0}[0],
				},
				{
					Name:        "category",
					Type:        models.FieldTypeEnum,
					DisplayName: "Category",
					EnumValues:  []string{"electronics", "clothing", "books"},
					Required:    true,
				},
			},
		},
	}

	generator := NewAPIDocsGenerator(opts)
	generator.spec = &models.OpenAPISpec{
		Components: &models.ComponentsObject{
			Schemas: make(map[string]*models.SchemaObject),
		},
	}

	t.Run("AddProductSchema", func(t *testing.T) {
		generator.addModelSchemas()

		productSchema := generator.spec.Components.Schemas["Product"]
		if productSchema == nil {
			t.Error("Product schema was not added")
		}

		// Check properties
		if productSchema.Properties == nil {
			t.Error("Product schema properties are missing")
		}

		idProp := productSchema.Properties["id"]
		if idProp == nil || idProp.Type != "integer" {
			t.Error("Product ID property is invalid")
		}

		nameProp := productSchema.Properties["name"]
		if nameProp == nil || nameProp.Type != "string" {
			t.Error("Product name property is invalid")
		}

		// Check required fields
		if len(productSchema.Required) != 4 {
			t.Errorf("Expected 4 required fields, got %d", len(productSchema.Required))
		}
	})
}

// TestAddEndpointPaths tests endpoint path addition
func TestAddEndpointPaths(t *testing.T) {
	tempDir := t.TempDir()
	opts := DefaultDocumentationGeneratorOptions("endpoint-test", tempDir)
	
	// Add test endpoints
	opts.Endpoints = []models.EndpointDefinition{
		{
			Method:      "GET",
			Path:        "/products",
			Summary:     "List products",
			Description: "Get all products",
			Tags:        []string{"Products"},
			Parameters: []models.Parameter{
				{
					Name:        "page",
					In:          "query",
					Description: "Page number",
					Required:    false,
					Schema: &models.SchemaObject{
						Type: "integer",
					},
				},
			},
		},
		{
			Method:      "POST",
			Path:        "/products",
			Summary:     "Create product",
			Description: "Create a new product",
			Tags:        []string{"Products"},
		},
	}

	generator := NewAPIDocsGenerator(opts)
	generator.spec = &models.OpenAPISpec{
		Paths: make(map[string]*models.PathItem),
	}

	t.Run("AddProductEndpoints", func(t *testing.T) {
		generator.addEndpointPaths()

		productPath := generator.spec.Paths["/products"]
		if productPath == nil {
			t.Error("Products path was not added")
		}

		// Check GET operation
		if productPath.Get == nil {
			t.Error("GET operation is missing")
		}

		if productPath.Get.Summary != "List products" {
			t.Errorf("Expected summary 'List products', got '%s'", productPath.Get.Summary)
		}

		// Check POST operation
		if productPath.Post == nil {
			t.Error("POST operation is missing")
		}

		if productPath.Post.Summary != "Create product" {
			t.Errorf("Expected summary 'Create product', got '%s'", productPath.Post.Summary)
		}
	})
}

// TestGenerateOperationID tests operation ID generation
func TestGenerateOperationID(t *testing.T) {
	tempDir := t.TempDir()
	opts := DefaultDocumentationGeneratorOptions("operation-test", tempDir)
	generator := NewAPIDocsGenerator(opts)

	tests := []struct {
		method   string
		path     string
		expected string
	}{
		{"GET", "/users", "getUsers"},
		{"POST", "/users/{id}/posts", "postUsers_id_posts"},
		{"PUT", "/", "putRoot"},
		{"DELETE", "/products/{productId}", "deleteProducts_productId"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s_%s", test.method, test.path), func(t *testing.T) {
			operationID := generator.generateOperationID(test.method, test.path)
			if operationID != test.expected {
				t.Errorf("Expected operation ID '%s', got '%s'", test.expected, operationID)
			}
		})
	}
}

// TestSwaggerUIGeneration tests Swagger UI generation
func TestSwaggerUIGeneration(t *testing.T) {
	tempDir := t.TempDir()
	opts := DefaultDocumentationGeneratorOptions("swagger-test", tempDir)
	
	generator := NewAPIDocsGenerator(opts)
	
	// Create directories
	err := generator.createOutputDirectory()
	if err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	t.Run("GenerateSwaggerUI", func(t *testing.T) {
		err := generator.generateSwaggerUI()
		if err != nil {
			t.Errorf("Failed to generate Swagger UI: %v", err)
		}

		// Check if Swagger UI HTML was generated
		swaggerPath := filepath.Join(tempDir, "docs", "swagger-ui", "index.html")
		if _, err := os.Stat(swaggerPath); os.IsNotExist(err) {
			t.Error("Swagger UI HTML file was not created")
		}

		// Read and validate HTML content
		htmlContent, err := os.ReadFile(swaggerPath)
		if err != nil {
			t.Fatalf("Failed to read Swagger UI HTML: %v", err)
		}

		htmlStr := string(htmlContent)
		if !strings.Contains(htmlStr, "swagger-ui") {
			t.Error("Swagger UI HTML does not contain 'swagger-ui'")
		}

		if !strings.Contains(htmlStr, "SwaggerUIBundle") {
			t.Error("Swagger UI HTML does not contain 'SwaggerUIBundle'")
		}
	})
}

// TestDefaultDocumentationGeneratorOptions tests default options
func TestDefaultDocumentationGeneratorOptions(t *testing.T) {
	opts := DefaultDocumentationGeneratorOptions("default-test", "/tmp/test")

	// Validate basic options
	if opts.ProjectName != "default-test" {
		t.Errorf("Expected project name 'default-test', got '%s'", opts.ProjectName)
	}

	if opts.OutputPath != "/tmp/test" {
		t.Errorf("Expected output path '/tmp/test', got '%s'", opts.OutputPath)
	}

	// Validate config
	if opts.Config.ProjectName != "default-test" {
		t.Errorf("Expected config project name 'default-test', got '%s'", opts.Config.ProjectName)
	}

	// Validate boolean flags
	if !opts.IncludeAuth {
		t.Error("IncludeAuth should be true by default")
	}

	if !opts.IncludeModels {
		t.Error("IncludeModels should be true by default")
	}

	if !opts.IncludeEndpoints {
		t.Error("IncludeEndpoints should be true by default")
	}

	if !opts.GenerateExamples {
		t.Error("GenerateExamples should be true by default")
	}

	if !opts.GenerateSwaggerUI {
		t.Error("GenerateSwaggerUI should be true by default")
	}

	// Validate Swagger UI config
	expectedTitle := "default-test API Documentation"
	if opts.SwaggerUIConfig.Title != expectedTitle {
		t.Errorf("Expected Swagger UI title '%s', got '%s'", expectedTitle, opts.SwaggerUIConfig.Title)
	}
}

// TestAPIDocsGetModuleName tests module name extraction for API docs
func TestAPIDocsGetModuleName(t *testing.T) {
	tempDir := t.TempDir()

	// Create a fake go.mod file
	goModContent := `module github.com/test/api-docs-project

go 1.21

require (
    github.com/gin-gonic/gin v1.9.0
)`

	goModPath := filepath.Join(tempDir, "go.mod")
	err := os.WriteFile(goModPath, []byte(goModContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create go.mod file: %v", err)
	}

	opts := DefaultDocumentationGeneratorOptions("module-test", tempDir)
	generator := NewAPIDocsGenerator(opts)

	moduleName := generator.getModuleName()
	expectedModule := "github.com/test/api-docs-project"

	if moduleName != expectedModule {
		t.Errorf("Expected module name '%s', got '%s'", expectedModule, moduleName)
	}
}

// Benchmark tests
func BenchmarkAPIDocsGenerator(b *testing.B) {
	tempDir := b.TempDir()
	opts := DefaultDocumentationGeneratorOptions("bench-test", tempDir)
	
	// Add some test data
	opts.Models = []models.ModelDefinition{
		{
			Name: "BenchmarkModel",
			Fields: []models.Field{
				{Name: "id", Type: models.FieldTypeNumber, Required: true},
				{Name: "name", Type: models.FieldTypeString, Required: true},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generator := NewAPIDocsGenerator(opts)
		_ = generator.generateOpenAPISpec()
	}
}

// TestAPIDocsWithAuthentication tests docs generation with auth
func TestAPIDocsWithAuthentication(t *testing.T) {
	tempDir := t.TempDir()
	opts := DefaultDocumentationGeneratorOptions("auth-docs-test", tempDir)
	
	// Configure authentication
	opts.IncludeAuth = true
	opts.AuthConfig = &models.AuthConfig{
		Provider: models.AuthProviderJWT,
		OAuth2Providers: []models.OAuth2Provider{
			{
				Name:        "google",
				ClientID:    "test-client-id",
				AuthURL:     "https://accounts.google.com/o/oauth2/auth",
				TokenURL:    "https://oauth2.googleapis.com/token",
				UserInfoURL: "https://www.googleapis.com/oauth2/v2/userinfo",
			},
		},
	}

	generator := NewAPIDocsGenerator(opts)

	t.Run("GenerateWithAuth", func(t *testing.T) {
		err := generator.GenerateAPIDocs()
		if err != nil {
			t.Errorf("Failed to generate API docs with auth: %v", err)
		}

		// Validate that auth paths were added
		if generator.spec.Paths["/auth/login"] == nil {
			t.Error("Auth login path was not added")
		}

		if generator.spec.Paths["/auth/register"] == nil {
			t.Error("Auth register path was not added")
		}

		// Check OAuth2 paths
		if generator.spec.Paths["/auth/google"] == nil {
			t.Error("OAuth2 Google path was not added")
		}
	})
}