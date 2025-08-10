package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// SchemaField represents a field in a resource schema
type SchemaField struct {
	Name         string                 `json:"name"`
	Type         string                 `json:"type"`
	DisplayName  string                 `json:"display_name"`
	Description  string                 `json:"description"`
	Required     bool                   `json:"required"`
	DefaultValue interface{}            `json:"default_value,omitempty"`
	Validation   *FieldValidation       `json:"validation,omitempty"`
	UI           *FieldUI               `json:"ui,omitempty"`
	Relation     *RelationConfig        `json:"relation,omitempty"`
	Database     *DatabaseFieldConfig   `json:"database,omitempty"`
	Frontend     *FrontendFieldConfig   `json:"frontend,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// FieldValidation contains validation rules for a field
type FieldValidation struct {
	MinLength    *int     `json:"min_length,omitempty"`
	MaxLength    *int     `json:"max_length,omitempty"`
	Pattern      string   `json:"pattern,omitempty"`
	Min          *float64 `json:"min,omitempty"`
	Max          *float64 `json:"max,omitempty"`
	AllowedValues []string `json:"allowed_values,omitempty"`
	CustomRules  []string `json:"custom_rules,omitempty"`
}

// FieldUI contains UI configuration for the field
type FieldUI struct {
	Component    string                 `json:"component"`     // "input", "textarea", "select", "checkbox", "datepicker", etc.
	Label        string                 `json:"label"`
	Placeholder  string                 `json:"placeholder,omitempty"`
	HelpText     string                 `json:"help_text,omitempty"`
	Order        int                    `json:"order"`
	Hidden       bool                   `json:"hidden,omitempty"`
	ReadOnly     bool                   `json:"readonly,omitempty"`
	Width        string                 `json:"width,omitempty"`        // "full", "half", "third", etc.
	Group        string                 `json:"group,omitempty"`        // For grouping fields
	Conditional  *ConditionalLogic      `json:"conditional,omitempty"`
	Props        map[string]interface{} `json:"props,omitempty"`        // Component-specific props
}

// ConditionalLogic defines when a field should be shown/hidden
type ConditionalLogic struct {
	Field     string      `json:"field"`
	Operator  string      `json:"operator"` // "equals", "not_equals", "contains", "greater_than", etc.
	Value     interface{} `json:"value"`
	Action    string      `json:"action"`   // "show", "hide", "enable", "disable"
}

// RelationConfig contains relationship configuration
type RelationConfig struct {
	Target       string   `json:"target"`        // Target model name
	Type         string   `json:"type"`          // "one_to_one", "one_to_many", "many_to_many"
	ForeignKey   string   `json:"foreign_key"`   // Foreign key field name
	LocalKey     string   `json:"local_key"`     // Local key field name
	PivotTable   string   `json:"pivot_table,omitempty"`
	Cascade      bool     `json:"cascade,omitempty"`
	Populate     bool     `json:"populate,omitempty"`
	PopulateDepth int     `json:"populate_depth,omitempty"`
}

// DatabaseFieldConfig contains database-specific configuration
type DatabaseFieldConfig struct {
	ColumnName   string                 `json:"column_name,omitempty"`
	Type         string                 `json:"type,omitempty"`      // Database-specific type
	Size         int                    `json:"size,omitempty"`      // For varchar, etc.
	Precision    int                    `json:"precision,omitempty"` // For decimal
	Scale        int                    `json:"scale,omitempty"`     // For decimal
	Nullable     bool                   `json:"nullable"`
	Index        bool                   `json:"index,omitempty"`
	Unique       bool                   `json:"unique,omitempty"`
	Primary      bool                   `json:"primary,omitempty"`
	AutoIncrement bool                  `json:"auto_increment,omitempty"`
	Default      interface{}            `json:"default,omitempty"`
	Comment      string                 `json:"comment,omitempty"`
	Extras       map[string]interface{} `json:"extras,omitempty"` // Database-specific extras
}

// FrontendFieldConfig contains frontend-specific configuration
type FrontendFieldConfig struct {
	AtomicComponent string                 `json:"atomic_component"` // "atom", "molecule", "organism"
	ComponentName   string                 `json:"component_name"`
	Props           map[string]interface{} `json:"props,omitempty"`
	Styling         *StylingConfig         `json:"styling,omitempty"`
	Validation      *FrontendValidation    `json:"validation,omitempty"`
}

// StylingConfig contains styling configuration
type StylingConfig struct {
	Classes     []string               `json:"classes,omitempty"`
	Styles      map[string]string      `json:"styles,omitempty"`
	Responsive  map[string]interface{} `json:"responsive,omitempty"`
	Theme       string                 `json:"theme,omitempty"`
}

// FrontendValidation contains frontend validation rules
type FrontendValidation struct {
	Required     bool     `json:"required"`
	Pattern      string   `json:"pattern,omitempty"`
	MinLength    *int     `json:"min_length,omitempty"`
	MaxLength    *int     `json:"max_length,omitempty"`
	Min          *float64 `json:"min,omitempty"`
	Max          *float64 `json:"max,omitempty"`
	ErrorMessage string   `json:"error_message,omitempty"`
}

// ResourceSchema represents a complete resource schema
type ResourceSchema struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	DisplayName  string                 `json:"display_name"`
	Description  string                 `json:"description"`
	Version      string                 `json:"version"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	CreatedBy    string                 `json:"created_by,omitempty"`
	Tags         []string               `json:"tags,omitempty"`
	
	// Naming conventions
	Names        *NamingConventions     `json:"names"`
	
	// Schema configuration
	Fields       []SchemaField          `json:"fields"`
	Indexes      []IndexConfig          `json:"indexes,omitempty"`
	Constraints  []ConstraintConfig     `json:"constraints,omitempty"`
	
	// Generation options
	Options      *GenerationOptions     `json:"options,omitempty"`
	
	// Database configuration
	Database     *DatabaseConfig        `json:"database,omitempty"`
	
	// Frontend configuration
	Frontend     *FrontendConfig        `json:"frontend,omitempty"`
	
	// Metadata
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// NamingConventions holds various naming formats
type NamingConventions struct {
	Singular       string `json:"singular"`        // user
	Plural         string `json:"plural"`          // users
	PascalCase     string `json:"pascal_case"`     // User
	PascalPlural   string `json:"pascal_plural"`   // Users
	CamelCase      string `json:"camel_case"`      // user
	CamelPlural    string `json:"camel_plural"`    // users
	SnakeCase      string `json:"snake_case"`      // user
	SnakePlural    string `json:"snake_plural"`    // users
	KebabCase      string `json:"kebab_case"`      // user
	KebabPlural    string `json:"kebab_plural"`    // users
	TableName      string `json:"table_name"`      // users
	CollectionName string `json:"collection_name"` // users
}

// IndexConfig represents database index configuration
type IndexConfig struct {
	Name    string   `json:"name"`
	Fields  []string `json:"fields"`
	Type    string   `json:"type"`    // "btree", "hash", "gist", "gin", etc.
	Unique  bool     `json:"unique,omitempty"`
	Partial string   `json:"partial,omitempty"` // Partial index condition
}

// ConstraintConfig represents database constraint configuration
type ConstraintConfig struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`       // "check", "foreign_key", "unique", etc.
	Fields     []string `json:"fields"`
	Reference  string   `json:"reference,omitempty"`
	OnDelete   string   `json:"on_delete,omitempty"`
	OnUpdate   string   `json:"on_update,omitempty"`
	Condition  string   `json:"condition,omitempty"`
}

// GenerationOptions contains code generation options
type GenerationOptions struct {
	GenerateAPI      bool     `json:"generate_api"`
	GenerateModel    bool     `json:"generate_model"`
	GenerateService  bool     `json:"generate_service"`
	GenerateRepo     bool     `json:"generate_repo"`
	GenerateHandler  bool     `json:"generate_handler"`
	GenerateTests    bool     `json:"generate_tests"`
	GenerateMocks    bool     `json:"generate_mocks"`
	GenerateDocs     bool     `json:"generate_docs"`
	GenerateFrontend bool     `json:"generate_frontend"`
	Features         []string `json:"features,omitempty"` // "auth", "validation", "caching", etc.
}

// DatabaseConfig contains database-specific configuration
type DatabaseConfig struct {
	Provider     string                 `json:"provider"`      // "postgres", "mysql", "sqlite", "mongodb", "supabase"
	TableName    string                 `json:"table_name"`
	Schema       string                 `json:"schema,omitempty"`
	Engine       string                 `json:"engine,omitempty"`
	Charset      string                 `json:"charset,omitempty"`
	Collation    string                 `json:"collation,omitempty"`
	Options      map[string]interface{} `json:"options,omitempty"`
	Migrations   *MigrationConfig       `json:"migrations,omitempty"`
}

// MigrationConfig contains migration configuration
type MigrationConfig struct {
	AutoMigrate  bool     `json:"auto_migrate"`
	BackupBefore bool     `json:"backup_before"`
	Versioned    bool     `json:"versioned"`
	Path         string   `json:"path,omitempty"`
	Seeds        []string `json:"seeds,omitempty"`
}

// FrontendConfig contains frontend generation configuration
type FrontendConfig struct {
	Framework       string                 `json:"framework"`        // "react", "vue", "angular", "svelte"
	ComponentStyle  string                 `json:"component_style"`  // "atomic", "feature", "page"
	StateManagement string                 `json:"state_management"` // "redux", "zustand", "pinia", etc.
	Styling         string                 `json:"styling"`          // "tailwind", "styled-components", "css-modules"
	TypeScript      bool                   `json:"typescript"`
	Forms           *FormConfig            `json:"forms,omitempty"`
	Tables          *TableConfig           `json:"tables,omitempty"`
	Options         map[string]interface{} `json:"options,omitempty"`
}

// FormConfig contains form generation configuration
type FormConfig struct {
	Library     string   `json:"library"`      // "formik", "react-hook-form", "vue-form", etc.
	Validation  string   `json:"validation"`   // "yup", "zod", "joi", etc.
	Components  []string `json:"components"`   // Custom component mappings
	Layout      string   `json:"layout"`       // "vertical", "horizontal", "inline"
	Grouping    bool     `json:"grouping"`     // Group fields by sections
}

// TableConfig contains table generation configuration
type TableConfig struct {
	Library     string   `json:"library"`      // "react-table", "ant-table", etc.
	Features    []string `json:"features"`     // "pagination", "sorting", "filtering", "export"
	Columns     []string `json:"columns"`      // Which fields to show as columns
	Actions     []string `json:"actions"`      // "view", "edit", "delete", "bulk"
	Responsive  bool     `json:"responsive"`
}

// SchemaStorage represents the storage interface for schemas
type SchemaStorage interface {
	Save(schema *ResourceSchema) error
	Load(id string) (*ResourceSchema, error)
	LoadByName(name string) (*ResourceSchema, error)
	List() ([]*ResourceSchema, error)
	Delete(id string) error
	Search(query string) ([]*ResourceSchema, error)
	GetVersions(id string) ([]*ResourceSchema, error)
}

// ToJSON converts the schema to JSON
func (s *ResourceSchema) ToJSON() ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}

// FromJSON creates a ResourceSchema from JSON
func FromJSON(data []byte) (*ResourceSchema, error) {
	var schema ResourceSchema
	err := json.Unmarshal(data, &schema)
	return &schema, err
}

// GetGoType returns the Go type for a field
func (f *SchemaField) GetGoType() string {
	switch f.Type {
	case "string", "text", "email", "url", "slug", "color", "file", "image":
		return "string"
	case "number", "integer":
		return "int64"
	case "float", "decimal":
		return "float64"
	case "boolean":
		return "bool"
	case "date", "datetime", "timestamp":
		return "time.Time"
	case "uuid":
		return "uuid.UUID"
	case "json", "mixed":
		return "json.RawMessage"
	case "relation":
		if f.Relation != nil {
			return "*" + f.Relation.Target
		}
		return "interface{}"
	case "relation_array":
		if f.Relation != nil {
			return "[]*" + f.Relation.Target
		}
		return "[]interface{}"
	case "location", "coordinates":
		return "*Location"
	case "currency":
		return "decimal.Decimal"
	case "enum":
		return "string" // Could be custom enum type
	default:
		return "interface{}"
	}
}

// GetGORMTag returns the GORM tag for a field
func (f *SchemaField) GetGORMTag(dbProvider string) string {
	if f.Database == nil {
		f.Database = &DatabaseFieldConfig{}
	}
	
	var tags []string
	
	// Column name
	if f.Database.ColumnName != "" {
		tags = append(tags, "column:"+f.Database.ColumnName)
	}
	
	// Type
	if f.Database.Type != "" {
		tags = append(tags, "type:"+f.Database.Type)
	} else {
		// Auto-generate type based on field type and provider
		dbType := f.getDBType(dbProvider)
		if dbType != "" {
			tags = append(tags, "type:"+dbType)
		}
	}
	
	// Size
	if f.Database.Size > 0 {
		tags = append(tags, fmt.Sprintf("size:%d", f.Database.Size))
	}
	
	// Precision and Scale
	if f.Database.Precision > 0 {
		tags = append(tags, fmt.Sprintf("precision:%d", f.Database.Precision))
	}
	if f.Database.Scale > 0 {
		tags = append(tags, fmt.Sprintf("scale:%d", f.Database.Scale))
	}
	
	// Nullable
	if !f.Database.Nullable {
		tags = append(tags, "not null")
	}
	
	// Index
	if f.Database.Index {
		tags = append(tags, "index")
	}
	
	// Unique
	if f.Database.Unique {
		tags = append(tags, "unique")
	}
	
	// Primary key
	if f.Database.Primary {
		tags = append(tags, "primaryKey")
	}
	
	// Auto increment
	if f.Database.AutoIncrement {
		tags = append(tags, "autoIncrement")
	}
	
	// Default value
	if f.Database.Default != nil {
		tags = append(tags, fmt.Sprintf("default:%v", f.Database.Default))
	}
	
	// Relations
	if f.Type == "relation" && f.Relation != nil {
		tags = append(tags, "foreignKey:"+f.Relation.ForeignKey)
	}
	
	// Many-to-many
	if f.Type == "relation_array" && f.Relation != nil && f.Relation.Type == "many_to_many" {
		tags = append(tags, "many2many:"+f.Relation.PivotTable)
	}
	
	return strings.Join(tags, ";")
}

// getDBType returns the database-specific type
func (f *SchemaField) getDBType(provider string) string {
	switch provider {
	case "postgres", "supabase":
		return f.getPostgresType()
	case "mysql":
		return f.getMySQLType()
	case "sqlite":
		return f.getSQLiteType()
	case "mongodb":
		return f.getMongoType()
	default:
		return ""
	}
}

// getPostgresType returns PostgreSQL-specific type
func (f *SchemaField) getPostgresType() string {
	switch f.Type {
	case "string":
		if f.Database != nil && f.Database.Size > 0 {
			return fmt.Sprintf("varchar(%d)", f.Database.Size)
		}
		return "varchar(255)"
	case "text":
		return "text"
	case "number", "integer":
		return "bigint"
	case "float", "decimal":
		return "decimal"
	case "boolean":
		return "boolean"
	case "date":
		return "date"
	case "datetime", "timestamp":
		return "timestamp"
	case "uuid":
		return "uuid"
	case "json", "mixed":
		return "jsonb"
	case "location", "coordinates":
		return "geometry(Point,4326)"
	default:
		return "text"
	}
}

// getMySQLType returns MySQL-specific type
func (f *SchemaField) getMySQLType() string {
	switch f.Type {
	case "string":
		if f.Database != nil && f.Database.Size > 0 {
			return fmt.Sprintf("varchar(%d)", f.Database.Size)
		}
		return "varchar(255)"
	case "text":
		return "text"
	case "number", "integer":
		return "bigint"
	case "float", "decimal":
		return "decimal"
	case "boolean":
		return "boolean"
	case "date":
		return "date"
	case "datetime", "timestamp":
		return "datetime"
	case "uuid":
		return "char(36)"
	case "json", "mixed":
		return "json"
	case "location", "coordinates":
		return "point"
	default:
		return "text"
	}
}

// getSQLiteType returns SQLite-specific type
func (f *SchemaField) getSQLiteType() string {
	switch f.Type {
	case "string", "text", "uuid":
		return "text"
	case "number", "integer":
		return "integer"
	case "float", "decimal":
		return "real"
	case "boolean":
		return "boolean"
	case "date", "datetime", "timestamp":
		return "datetime"
	case "json", "mixed":
		return "text"
	default:
		return "text"
	}
}

// getMongoType returns MongoDB-specific type
func (f *SchemaField) getMongoType() string {
	switch f.Type {
	case "string", "text":
		return "string"
	case "number", "integer":
		return "number"
	case "float", "decimal":
		return "number"
	case "boolean":
		return "boolean"
	case "date", "datetime", "timestamp":
		return "date"
	case "uuid":
		return "string"
	case "json", "mixed":
		return "object"
	case "location", "coordinates":
		return "2dsphere"
	default:
		return "string"
	}
}