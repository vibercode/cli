package storage

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/vibercode/cli/internal/models"
)

// FileSchemaStorage implements SchemaStorage using file system
type FileSchemaStorage struct {
	basePath string
}

// NewFileSchemaStorage creates a new file-based schema storage
func NewFileSchemaStorage(basePath string) *FileSchemaStorage {
	return &FileSchemaStorage{
		basePath: basePath,
	}
}

// ensureDir ensures the directory exists
func (fs *FileSchemaStorage) ensureDir() error {
	return os.MkdirAll(fs.basePath, 0755)
}

// getSchemaPath returns the path for a schema file
func (fs *FileSchemaStorage) getSchemaPath(id string) string {
	return filepath.Join(fs.basePath, id+".json")
}

// Save saves a schema to file system
func (fs *FileSchemaStorage) Save(schema *models.ResourceSchema) error {
	if err := fs.ensureDir(); err != nil {
		return fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Update timestamps
	if schema.CreatedAt.IsZero() {
		schema.CreatedAt = time.Now()
	}
	schema.UpdatedAt = time.Now()

	// Generate ID if not provided
	if schema.ID == "" {
		schema.ID = generateSchemaID(schema.Name)
	}

	data, err := schema.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal schema: %w", err)
	}

	filePath := fs.getSchemaPath(schema.ID)
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write schema file: %w", err)
	}

	return nil
}

// Load loads a schema by ID
func (fs *FileSchemaStorage) Load(id string) (*models.ResourceSchema, error) {
	filePath := fs.getSchemaPath(id)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("schema not found: %s", id)
		}
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}

	schema, err := models.FromJSON(data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal schema: %w", err)
	}

	return schema, nil
}

// LoadByName loads a schema by name
func (fs *FileSchemaStorage) LoadByName(name string) (*models.ResourceSchema, error) {
	schemas, err := fs.List()
	if err != nil {
		return nil, err
	}

	for _, schema := range schemas {
		if schema.Name == name {
			return schema, nil
		}
	}

	return nil, fmt.Errorf("schema not found: %s", name)
}

// List lists all schemas
func (fs *FileSchemaStorage) List() ([]*models.ResourceSchema, error) {
	if err := fs.ensureDir(); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	files, err := ioutil.ReadDir(fs.basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read storage directory: %w", err)
	}

	var schemas []*models.ResourceSchema
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			id := strings.TrimSuffix(file.Name(), ".json")
			schema, err := fs.Load(id)
			if err != nil {
				continue // Skip corrupted files
			}
			schemas = append(schemas, schema)
		}
	}

	// Sort by creation date (newest first)
	sort.Slice(schemas, func(i, j int) bool {
		return schemas[i].CreatedAt.After(schemas[j].CreatedAt)
	})

	return schemas, nil
}

// Delete deletes a schema by ID
func (fs *FileSchemaStorage) Delete(id string) error {
	filePath := fs.getSchemaPath(id)
	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("schema not found: %s", id)
		}
		return fmt.Errorf("failed to delete schema file: %w", err)
	}
	return nil
}

// Search searches schemas by query
func (fs *FileSchemaStorage) Search(query string) ([]*models.ResourceSchema, error) {
	allSchemas, err := fs.List()
	if err != nil {
		return nil, err
	}

	query = strings.ToLower(query)
	var results []*models.ResourceSchema

	for _, schema := range allSchemas {
		if fs.matchesQuery(schema, query) {
			results = append(results, schema)
		}
	}

	return results, nil
}

// matchesQuery checks if a schema matches the search query
func (fs *FileSchemaStorage) matchesQuery(schema *models.ResourceSchema, query string) bool {
	searchText := strings.ToLower(fmt.Sprintf("%s %s %s %s",
		schema.Name,
		schema.DisplayName,
		schema.Description,
		strings.Join(schema.Tags, " "),
	))

	return strings.Contains(searchText, query)
}

// GetVersions returns all versions of a schema (for now, just returns the current one)
func (fs *FileSchemaStorage) GetVersions(id string) ([]*models.ResourceSchema, error) {
	schema, err := fs.Load(id)
	if err != nil {
		return nil, err
	}

	return []*models.ResourceSchema{schema}, nil
}

// generateSchemaID generates a unique ID for a schema
func generateSchemaID(name string) string {
	// Convert to lowercase and replace spaces/special chars with hyphens
	id := strings.ToLower(name)
	id = strings.ReplaceAll(id, " ", "-")
	id = strings.ReplaceAll(id, "_", "-")
	
	// Add timestamp to ensure uniqueness
	timestamp := time.Now().Format("20060102-150405")
	return fmt.Sprintf("%s-%s", id, timestamp)
}

// SchemaRepository provides higher-level operations on schemas
type SchemaRepository struct {
	storage models.SchemaStorage
}

// Storage returns the underlying storage interface
func (r *SchemaRepository) Storage() models.SchemaStorage {
	return r.storage
}

// Load loads a schema by ID
func (r *SchemaRepository) Load(id string) (*models.ResourceSchema, error) {
	return r.storage.Load(id)
}

// List lists all schemas
func (r *SchemaRepository) List() ([]*models.ResourceSchema, error) {
	return r.storage.List()
}

// Search searches schemas by query
func (r *SchemaRepository) Search(query string) ([]*models.ResourceSchema, error) {
	return r.storage.Search(query)
}

// Delete deletes a schema by ID
func (r *SchemaRepository) Delete(id string) error {
	return r.storage.Delete(id)
}

// NewSchemaRepository creates a new schema repository
func NewSchemaRepository(storage models.SchemaStorage) *SchemaRepository {
	return &SchemaRepository{
		storage: storage,
	}
}

// CreateSchema creates a new schema with validation
func (r *SchemaRepository) CreateSchema(schema *models.ResourceSchema) error {
	// Validate required fields
	if schema.Name == "" {
		return fmt.Errorf("schema name is required")
	}
	
	if schema.DisplayName == "" {
		schema.DisplayName = schema.Name
	}
	
	if len(schema.Fields) == 0 {
		return fmt.Errorf("schema must have at least one field")
	}

	// Generate naming conventions
	if schema.Names == nil {
		schema.Names = generateNamingConventions(schema.Name)
	}

	// Set default version
	if schema.Version == "" {
		schema.Version = "1.0.0"
	}

	// Validate and set field defaults
	for i := range schema.Fields {
		if err := r.validateField(&schema.Fields[i]); err != nil {
			return fmt.Errorf("field %s: %w", schema.Fields[i].Name, err)
		}
	}

	return r.storage.Save(schema)
}

// validateField validates a schema field
func (r *SchemaRepository) validateField(field *models.SchemaField) error {
	if field.Name == "" {
		return fmt.Errorf("field name is required")
	}
	
	if field.Type == "" {
		return fmt.Errorf("field type is required")
	}
	
	if field.DisplayName == "" {
		field.DisplayName = field.Name
	}

	// Set default UI configuration
	if field.UI == nil {
		field.UI = &models.FieldUI{
			Component: getDefaultUIComponent(field.Type),
			Label:     field.DisplayName,
			Order:     0,
		}
	}

	// Set default database configuration
	if field.Database == nil {
		field.Database = &models.DatabaseFieldConfig{
			Nullable: !field.Required,
		}
	}

	// Validate relations
	if field.Type == "relation" || field.Type == "relation_array" {
		if field.Relation == nil {
			return fmt.Errorf("relation configuration is required for relation fields")
		}
		if field.Relation.Target == "" {
			return fmt.Errorf("relation target is required")
		}
	}

	return nil
}

// getDefaultUIComponent returns the default UI component for a field type
func getDefaultUIComponent(fieldType string) string {
	switch fieldType {
	case "string":
		return "input"
	case "text":
		return "textarea"
	case "number", "integer", "float":
		return "number-input"
	case "boolean":
		return "checkbox"
	case "date":
		return "date-picker"
	case "datetime", "timestamp":
		return "datetime-picker"
	case "email":
		return "email-input"
	case "url":
		return "url-input"
	case "file", "image":
		return "file-upload"
	case "relation":
		return "select"
	case "relation_array":
		return "multi-select"
	case "location", "coordinates":
		return "map-picker"
	case "color":
		return "color-picker"
	case "enum":
		return "select"
	case "json", "mixed":
		return "json-editor"
	default:
		return "input"
	}
}

// generateNamingConventions generates naming conventions for a schema
func generateNamingConventions(name string) *models.NamingConventions {
	return &models.NamingConventions{
		Singular:       strings.ToLower(name),
		Plural:         strings.ToLower(name) + "s", // Simple pluralization
		PascalCase:     toPascalCase(name),
		PascalPlural:   toPascalCase(name) + "s",
		CamelCase:      toCamelCase(name),
		CamelPlural:    toCamelCase(name) + "s",
		SnakeCase:      toSnakeCase(name),
		SnakePlural:    toSnakeCase(name) + "s",
		KebabCase:      toKebabCase(name),
		KebabPlural:    toKebabCase(name) + "s",
		TableName:      toSnakeCase(name) + "s",
		CollectionName: toSnakeCase(name) + "s",
	}
}

// String case conversion helpers
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

// GetDefaultSchemaPath returns the default path for schema storage
func GetDefaultSchemaPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".vibercode/schemas"
	}
	return filepath.Join(homeDir, ".vibercode", "schemas")
}

// SchemaTemplate represents a schema template
type SchemaTemplate struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Category    string                 `json:"category"`
	Schema      *models.ResourceSchema `json:"schema"`
	Preview     string                 `json:"preview,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
}

// LoadSchemaTemplates loads built-in schema templates
func LoadSchemaTemplates() []*SchemaTemplate {
	return []*SchemaTemplate{
		{
			Name:        "User Management",
			Description: "Complete user management schema with authentication fields",
			Category:    "Authentication",
			Schema: &models.ResourceSchema{
				Name:        "User",
				DisplayName: "User",
				Description: "User management with authentication",
				Fields: []models.SchemaField{
					{
						Name:        "email",
						Type:        "email",
						DisplayName: "Email",
						Required:    true,
						Database:    &models.DatabaseFieldConfig{Unique: true},
					},
					{
						Name:        "password",
						Type:        "string",
						DisplayName: "Password",
						Required:    true,
						UI:          &models.FieldUI{Component: "password-input"},
					},
					{
						Name:        "first_name",
						Type:        "string",
						DisplayName: "First Name",
						Required:    true,
					},
					{
						Name:        "last_name",
						Type:        "string",
						DisplayName: "Last Name",
						Required:    true,
					},
					{
						Name:        "avatar",
						Type:        "image",
						DisplayName: "Avatar",
						Required:    false,
					},
					{
						Name:        "is_active",
						Type:        "boolean",
						DisplayName: "Active",
						Required:    true,
						DefaultValue: true,
					},
				},
			},
			Tags: []string{"user", "auth", "management"},
		},
		{
			Name:        "Blog Post",
			Description: "Blog post schema with rich content and metadata",
			Category:    "Content",
			Schema: &models.ResourceSchema{
				Name:        "Post",
				DisplayName: "Blog Post",
				Description: "Blog post with rich content",
				Fields: []models.SchemaField{
					{
						Name:        "title",
						Type:        "string",
						DisplayName: "Title",
						Required:    true,
					},
					{
						Name:        "slug",
						Type:        "slug",
						DisplayName: "URL Slug",
						Required:    true,
						Database:    &models.DatabaseFieldConfig{Unique: true},
					},
					{
						Name:        "content",
						Type:        "text",
						DisplayName: "Content",
						Required:    true,
						UI:          &models.FieldUI{Component: "rich-editor"},
					},
					{
						Name:        "excerpt",
						Type:        "text",
						DisplayName: "Excerpt",
						Required:    false,
					},
					{
						Name:        "featured_image",
						Type:        "image",
						DisplayName: "Featured Image",
						Required:    false,
					},
					{
						Name:        "published_at",
						Type:        "datetime",
						DisplayName: "Published At",
						Required:    false,
					},
					{
						Name:        "is_published",
						Type:        "boolean",
						DisplayName: "Published",
						Required:    true,
						DefaultValue: false,
					},
				},
			},
			Tags: []string{"blog", "content", "cms"},
		},
		{
			Name:        "E-commerce Product",
			Description: "Product schema for e-commerce applications",
			Category:    "E-commerce",
			Schema: &models.ResourceSchema{
				Name:        "Product",
				DisplayName: "Product",
				Description: "E-commerce product with pricing and inventory",
				Fields: []models.SchemaField{
					{
						Name:        "name",
						Type:        "string",
						DisplayName: "Product Name",
						Required:    true,
					},
					{
						Name:        "sku",
						Type:        "string",
						DisplayName: "SKU",
						Required:    true,
						Database:    &models.DatabaseFieldConfig{Unique: true},
					},
					{
						Name:        "description",
						Type:        "text",
						DisplayName: "Description",
						Required:    false,
					},
					{
						Name:        "price",
						Type:        "decimal",
						DisplayName: "Price",
						Required:    true,
						Database:    &models.DatabaseFieldConfig{Precision: 10, Scale: 2},
					},
					{
						Name:        "cost",
						Type:        "decimal",
						DisplayName: "Cost",
						Required:    false,
						Database:    &models.DatabaseFieldConfig{Precision: 10, Scale: 2},
					},
					{
						Name:        "stock_quantity",
						Type:        "integer",
						DisplayName: "Stock Quantity",
						Required:    true,
						DefaultValue: 0,
					},
					{
						Name:        "images",
						Type:        "relation_array",
						DisplayName: "Images",
						Required:    false,
						Relation: &models.RelationConfig{
							Target: "File",
							Type:   "one_to_many",
						},
					},
					{
						Name:        "is_active",
						Type:        "boolean",
						DisplayName: "Active",
						Required:    true,
						DefaultValue: true,
					},
				},
			},
			Tags: []string{"product", "ecommerce", "inventory"},
		},
	}
}