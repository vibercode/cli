package templates

import (
	"strings"
	"text/template"
)

// SchemaModelTemplate generates Go model from schema
const SchemaModelTemplate = `package models

import (
	"time"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
{{- range .RequiredImports}}
	"{{.}}"
{{- end}}
{{- if .Database.Provider | eq "supabase"}}
	"github.com/google/uuid"
{{- end}}
)

{{- range .Relations}}
// {{.Name}} represents the {{.DisplayName}} relation
type {{.Name}} struct {
	// Relations will be populated based on related schemas
}
{{- end}}

// {{.Names.PascalCase}} represents the {{.DisplayName}} model
type {{.Names.PascalCase}} struct {
	ID        primitive.ObjectID ` + "`" + `json:"id" bson:"_id,omitempty"` + "`" + `
	CreatedAt time.Time          ` + "`" + `json:"created_at" bson:"created_at"` + "`" + `
	UpdatedAt time.Time          ` + "`" + `json:"updated_at" bson:"updated_at"` + "`" + `

{{- range .Fields}}
	{{.GoStructField}}
{{- end}}
}

// CollectionName returns the MongoDB collection name for {{.Names.PascalCase}}
func ({{.Names.PascalCase}}) CollectionName() string {
	return "{{.Names.TableName}}"
}

// {{.Names.PascalCase}}Request represents the request payload for creating/updating {{.DisplayName}}
type {{.Names.PascalCase}}Request struct {
{{- range .Fields}}
{{- if not .ReadOnly}}
	{{.GoRequestField}}
{{- end}}
{{- end}}
}

// {{.Names.PascalCase}}Response represents the response payload for {{.DisplayName}}
type {{.Names.PascalCase}}Response struct {
	ID        primitive.ObjectID ` + "`" + `json:"id"` + "`" + `
	CreatedAt time.Time          ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time          ` + "`" + `json:"updated_at"` + "`" + `

{{- range .Fields}}
	{{.GoResponseField}}
{{- end}}
}

// To{{.Names.PascalCase}}Response converts model to response
func (m *{{.Names.PascalCase}}) To{{.Names.PascalCase}}Response() *{{.Names.PascalCase}}Response {
	return &{{.Names.PascalCase}}Response{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
{{- range .Fields}}
		{{.Names.PascalCase}}: m.{{.Names.PascalCase}},
{{- end}}
	}
}

// {{.Names.PascalCase}}Filter represents filter options for {{.DisplayName}}
type {{.Names.PascalCase}}Filter struct {
	Page     int    ` + "`" + `json:"page" form:"page"` + "`" + `
	PageSize int    ` + "`" + `json:"page_size" form:"page_size"` + "`" + `
	Sort     string ` + "`" + `json:"sort" form:"sort"` + "`" + `
	Order    string ` + "`" + `json:"order" form:"order"` + "`" + `
	Search   string ` + "`" + `json:"search" form:"search"` + "`" + `

{{- range .Fields}}
{{- if .Filterable}}
	{{.GoFilterField}}
{{- end}}
{{- end}}
}

// Validate validates the {{.Names.PascalCase}}Request
func (r *{{.Names.PascalCase}}Request) Validate() error {
{{- range .Fields}}
{{- if .Required}}
	{{.GoValidation}}
{{- end}}
{{- end}}
	return nil
}
`

// SchemaRepositoryTemplate generates repository layer
const SchemaRepositoryTemplate = `package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"
	"{{.Module}}/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// {{.Names.PascalCase}}Repository handles database operations for {{.DisplayName}}
type {{.Names.PascalCase}}Repository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// New{{.Names.PascalCase}}Repository creates a new {{.Names.PascalCase}} repository
func New{{.Names.PascalCase}}Repository(db *mongo.Database) *{{.Names.PascalCase}}Repository {
	return &{{.Names.PascalCase}}Repository{
		db:         db,
		collection: db.Collection("{{.Names.TableName}}"),
	}
}

// Create creates a new {{.Names.Singular}}
func (r *{{.Names.PascalCase}}Repository) Create(ctx context.Context, {{.Names.CamelCase}} *models.{{.Names.PascalCase}}) error {
	{{.Names.CamelCase}}.ID = primitive.NewObjectID()
	{{.Names.CamelCase}}.CreatedAt = time.Now()
	{{.Names.CamelCase}}.UpdatedAt = time.Now()
	
	_, err := r.collection.InsertOne(ctx, {{.Names.CamelCase}})
	return err
}

// GetByID retrieves a {{.Names.Singular}} by ID
func (r *{{.Names.PascalCase}}Repository) GetByID(ctx context.Context, idStr string) (*models.{{.Names.PascalCase}}, error) {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	var {{.Names.CamelCase}} models.{{.Names.PascalCase}}
	err = r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&{{.Names.CamelCase}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &{{.Names.CamelCase}}, nil
}

// GetAll retrieves all {{.Names.Plural}} with filtering
func (r *{{.Names.PascalCase}}Repository) GetAll(ctx context.Context, filter *models.{{.Names.PascalCase}}Filter) ([]*models.{{.Names.PascalCase}}, int64, error) {
	var {{.Names.CamelPlural}} []*models.{{.Names.PascalCase}}

	// Build filter
	mongoFilter := bson.M{}
	if filter.Search != "" {
		searchConditions := bson.A{}
		// Add search conditions for string fields dynamically
		searchConditions = append(searchConditions, bson.M{"name": bson.M{"$regex": filter.Search, "$options": "i"}})
		if len(searchConditions) > 0 {
			mongoFilter["$or"] = searchConditions
		}
	}

	// Count total records
	total, err := r.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination and sorting
	opts := options.Find()
	if filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		opts.SetSkip(int64(offset))
		opts.SetLimit(int64(filter.PageSize))
	}

	// Apply sorting
	sortField := "created_at"
	sortOrder := -1
	if filter.Sort != "" {
		sortField = filter.Sort
		if strings.ToUpper(filter.Order) == "ASC" {
			sortOrder = 1
		}
	}
	opts.SetSort(bson.D{{ "{" }}{{ "{" }}sortField, sortOrder{{ "}" }}{{ "}" }})

	// Execute query
	cursor, err := r.collection.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &{{.Names.CamelPlural}}); err != nil {
		return nil, 0, err
	}

	return {{.Names.CamelPlural}}, total, nil
}

// Update updates a {{.Names.Singular}}
func (r *{{.Names.PascalCase}}Repository) Update(ctx context.Context, {{.Names.CamelCase}} *models.{{.Names.PascalCase}}) error {
	{{.Names.CamelCase}}.UpdatedAt = time.Now()
	
	filter := bson.M{"_id": {{.Names.CamelCase}}.ID}
	update := bson.M{"$set": {{.Names.CamelCase}}}
	
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete deletes a {{.Names.Singular}}
func (r *{{.Names.PascalCase}}Repository) Delete(ctx context.Context, idStr string) error {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}
	
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// HardDelete permanently deletes a {{.Names.Singular}} (same as Delete in MongoDB)
func (r *{{.Names.PascalCase}}Repository) HardDelete(ctx context.Context, idStr string) error {
	return r.Delete(ctx, idStr)
}

// Exists checks if a {{.Names.Singular}} exists
func (r *{{.Names.PascalCase}}Repository) Exists(ctx context.Context, idStr string) (bool, error) {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return false, fmt.Errorf("invalid ID format: %w", err)
	}
	
	count, err := r.collection.CountDocuments(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

{{- range .Fields}}
{{- if .Database.Unique}}
// GetBy{{.Names.PascalCase}} retrieves a {{$.Names.Singular}} by {{.DisplayName}}
func (r *{{$.Names.PascalCase}}Repository) GetBy{{.Names.PascalCase}}(ctx context.Context, {{.Names.CamelCase}} {{.GoType}}) (*models.{{$.Names.PascalCase}}, error) {
	var {{$.Names.CamelCase}} models.{{$.Names.PascalCase}}
	err := r.collection.FindOne(ctx, bson.M{"{{.Database.ColumnName}}": {{.Names.CamelCase}}}).Decode(&{{$.Names.CamelCase}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &{{$.Names.CamelCase}}, nil
}
{{- end}}
{{- end}}

// Repository interface for dependency injection
type {{.Names.PascalCase}}RepositoryInterface interface {
	Create(ctx context.Context, {{.Names.CamelCase}} *models.{{.Names.PascalCase}}) error
	GetByID(ctx context.Context, id string) (*models.{{.Names.PascalCase}}, error)
	GetAll(ctx context.Context, filter *models.{{.Names.PascalCase}}Filter) ([]*models.{{.Names.PascalCase}}, int64, error)
	Update(ctx context.Context, {{.Names.CamelCase}} *models.{{.Names.PascalCase}}) error
	Delete(ctx context.Context, id string) error
	HardDelete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
{{- range .Fields}}
{{- if .Database.Unique}}
	GetBy{{.Names.PascalCase}}(ctx context.Context, {{.Names.CamelCase}} {{.GoType}}) (*models.{{.Names.PascalCase}}, error)
{{- end}}
{{- end}}
}
`

// SchemaServiceTemplate generates service layer
const SchemaServiceTemplate = `package services

import (
	"fmt"
	"{{.Module}}/internal/models"
	"{{.Module}}/internal/repositories"
)

// {{.Names.PascalCase}}Service handles business logic for {{.DisplayName}}
type {{.Names.PascalCase}}Service struct {
	repo repositories.{{.Names.PascalCase}}RepositoryInterface
}

// New{{.Names.PascalCase}}Service creates a new {{.Names.PascalCase}} service
func New{{.Names.PascalCase}}Service(repo repositories.{{.Names.PascalCase}}RepositoryInterface) *{{.Names.PascalCase}}Service {
	return &{{.Names.PascalCase}}Service{repo: repo}
}

// Create creates a new {{.Names.Singular}}
func (s *{{.Names.PascalCase}}Service) Create(ctx context.Context, req *models.{{.Names.PascalCase}}Request) (*models.{{.Names.PascalCase}}, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Business logic validations
	if err := s.validateCreate(ctx, req); err != nil {
		return nil, err
	}

	// Convert request to model
	{{.Names.CamelCase}} := &models.{{.Names.PascalCase}}{
{{- range .Fields}}
{{- if not .ReadOnly}}
		{{.Names.PascalCase}}: req.{{.Names.PascalCase}},
{{- end}}
{{- end}}
	}

	// Create in database
	if err := s.repo.Create(ctx, {{.Names.CamelCase}}); err != nil {
		return nil, fmt.Errorf("failed to create {{.Names.Singular}}: %w", err)
	}

	return {{.Names.CamelCase}}, nil
}

// GetByID retrieves a {{.Names.Singular}} by ID
func (s *{{.Names.PascalCase}}Service) GetByID(ctx context.Context, id string) (*models.{{.Names.PascalCase}}, error) {
	{{.Names.CamelCase}}, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get {{.Names.Singular}}: %w", err)
	}
	return {{.Names.CamelCase}}, nil
}

// GetAll retrieves all {{.Names.Plural}} with filtering
func (s *{{.Names.PascalCase}}Service) GetAll(ctx context.Context, filter *models.{{.Names.PascalCase}}Filter) ([]*models.{{.Names.PascalCase}}, int64, error) {
	// Apply default pagination
	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.PageSize == 0 {
		filter.PageSize = 20
	}

	{{.Names.CamelPlural}}, total, err := s.repo.GetAll(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get {{.Names.Plural}}: %w", err)
	}

	return {{.Names.CamelPlural}}, total, nil
}

// Update updates a {{.Names.Singular}}
func (s *{{.Names.PascalCase}}Service) Update(ctx context.Context, id string, req *models.{{.Names.PascalCase}}Request) (*models.{{.Names.PascalCase}}, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get existing {{.Names.Singular}}
	{{.Names.CamelCase}}, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("{{.Names.Singular}} not found: %w", err)
	}

	// Business logic validations
	if err := s.validateUpdate(ctx, {{.Names.CamelCase}}, req); err != nil {
		return nil, err
	}

	// Update fields
{{- range .Fields}}
{{- if not .ReadOnly}}
	{{$.Names.CamelCase}}.{{.Names.PascalCase}} = req.{{.Names.PascalCase}}
{{- end}}
{{- end}}

	// Update in database
	if err := s.repo.Update(ctx, {{.Names.CamelCase}}); err != nil {
		return nil, fmt.Errorf("failed to update {{.Names.Singular}}: %w", err)
	}

	return {{.Names.CamelCase}}, nil
}

// Delete deletes a {{.Names.Singular}}
func (s *{{.Names.PascalCase}}Service) Delete(ctx context.Context, id string) error {
	// Check if {{.Names.Singular}} exists
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check {{.Names.Singular}} existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("{{.Names.Singular}} not found")
	}

	// Business logic validations
	if err := s.validateDelete(ctx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete {{.Names.Singular}}: %w", err)
	}

	return nil
}

// validateCreate validates business rules for creating {{.Names.Singular}}
func (s *{{.Names.PascalCase}}Service) validateCreate(ctx context.Context, req *models.{{.Names.PascalCase}}Request) error {
{{- range .Fields}}
{{- if .Database.Unique}}
	// Check if {{.DisplayName}} already exists
	if _, err := s.repo.GetBy{{.Names.PascalCase}}(ctx, req.{{.Names.PascalCase}}); err == nil {
		return fmt.Errorf("{{.DisplayName}} already exists")
	}
{{- end}}
{{- end}}
	return nil
}

// validateUpdate validates business rules for updating {{.Names.Singular}}
func (s *{{.Names.PascalCase}}Service) validateUpdate(ctx context.Context, existing *models.{{.Names.PascalCase}}, req *models.{{.Names.PascalCase}}Request) error {
{{- range .Fields}}
{{- if .Database.Unique}}
	// Check if {{.DisplayName}} already exists (excluding current record)
	if existing.{{.Names.PascalCase}} != req.{{.Names.PascalCase}} {
		if _, err := s.repo.GetBy{{.Names.PascalCase}}(ctx, req.{{.Names.PascalCase}}); err == nil {
			return fmt.Errorf("{{.DisplayName}} already exists")
		}
	}
{{- end}}
{{- end}}
	return nil
}

// validateDelete validates business rules for deleting {{.Names.Singular}}
func (s *{{.Names.PascalCase}}Service) validateDelete(ctx context.Context, id string) error {
	// Add custom delete validations here
	return nil
}

// Service interface for dependency injection
type {{.Names.PascalCase}}ServiceInterface interface {
	Create(ctx context.Context, req *models.{{.Names.PascalCase}}Request) (*models.{{.Names.PascalCase}}, error)
	GetByID(ctx context.Context, id string) (*models.{{.Names.PascalCase}}, error)
	GetAll(ctx context.Context, filter *models.{{.Names.PascalCase}}Filter) ([]*models.{{.Names.PascalCase}}, int64, error)
	Update(ctx context.Context, id string, req *models.{{.Names.PascalCase}}Request) (*models.{{.Names.PascalCase}}, error)
	Delete(ctx context.Context, id string) error
}
`

// SchemaHandlerTemplate generates HTTP handler
const SchemaHandlerTemplate = `package handlers

import (
	"net/http"
	"strconv"
	"{{.Module}}/internal/models"
	"{{.Module}}/internal/services"
	"github.com/gin-gonic/gin"
)

// {{.Names.PascalCase}}Handler handles HTTP requests for {{.DisplayName}}
type {{.Names.PascalCase}}Handler struct {
	service services.{{.Names.PascalCase}}ServiceInterface
}

// New{{.Names.PascalCase}}Handler creates a new {{.Names.PascalCase}} handler
func New{{.Names.PascalCase}}Handler(service services.{{.Names.PascalCase}}ServiceInterface) *{{.Names.PascalCase}}Handler {
	return &{{.Names.PascalCase}}Handler{service: service}
}

// Create handles POST /{{.Names.KebabPlural}}
func (h *{{.Names.PascalCase}}Handler) Create(c *gin.Context) {
	var req models.{{.Names.PascalCase}}Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	{{.Names.CamelCase}}, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, {{.Names.CamelCase}}.To{{.Names.PascalCase}}Response())
}

// GetByID handles GET /{{.Names.KebabPlural}}/:id
func (h *{{.Names.PascalCase}}Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	{{.Names.CamelCase}}, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, {{.Names.CamelCase}}.To{{.Names.PascalCase}}Response())
}

// GetAll handles GET /{{.Names.KebabPlural}}
func (h *{{.Names.PascalCase}}Handler) GetAll(c *gin.Context) {
	var filter models.{{.Names.PascalCase}}Filter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	{{.Names.CamelPlural}}, total, err := h.service.GetAll(c.Request.Context(), &filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to response format
	responses := make([]*models.{{.Names.PascalCase}}Response, len({{.Names.CamelPlural}}))
	for i, {{.Names.CamelCase}} := range {{.Names.CamelPlural}} {
		responses[i] = {{.Names.CamelCase}}.To{{.Names.PascalCase}}Response()
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      responses,
		"total":     total,
		"page":      filter.Page,
		"page_size": filter.PageSize,
	})
}

// Update handles PUT /{{.Names.KebabPlural}}/:id
func (h *{{.Names.PascalCase}}Handler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req models.{{.Names.PascalCase}}Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	{{.Names.CamelCase}}, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, {{.Names.CamelCase}}.To{{.Names.PascalCase}}Response())
}

// Delete handles DELETE /{{.Names.KebabPlural}}/:id
func (h *{{.Names.PascalCase}}Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "{{.DisplayName}} deleted successfully"})
}

// Setup{{.Names.PascalCase}}Routes sets up routes for {{.DisplayName}}
func Setup{{.Names.PascalCase}}Routes(r *gin.RouterGroup, handler *{{.Names.PascalCase}}Handler) {
	{{.Names.CamelPlural}} := r.Group("/{{.Names.KebabPlural}}")
	{
		{{.Names.CamelPlural}}.POST("", handler.Create)
		{{.Names.CamelPlural}}.GET("", handler.GetAll)
		{{.Names.CamelPlural}}.GET("/:id", handler.GetByID)
		{{.Names.CamelPlural}}.PUT("/:id", handler.Update)
		{{.Names.CamelPlural}}.DELETE("/:id", handler.Delete)
	}
}
`

// SchemaHelperFunctions contains helper functions for templates
var SchemaHelperFunctions = template.FuncMap{
	"eq": func(a, b interface{}) bool {
		return a == b
	},
	"ne": func(a, b interface{}) bool {
		return a != b
	},
	"contains": func(slice []string, item string) bool {
		for _, s := range slice {
			if s == item {
				return true
			}
		}
		return false
	},
	"join": func(slice []string, separator string) string {
		return strings.Join(slice, separator)
	},
	"lower": func(s string) string {
		return strings.ToLower(s)
	},
	"upper": func(s string) string {
		return strings.ToUpper(s)
	},
	"title": func(s string) string {
		return strings.Title(s)
	},
	"replace": func(s, old, new string) string {
		return strings.ReplaceAll(s, old, new)
	},
	"hasPrefix": func(s, prefix string) bool {
		return strings.HasPrefix(s, prefix)
	},
	"hasSuffix": func(s, suffix string) bool {
		return strings.HasSuffix(s, suffix)
	},
}

// GetSchemaTemplates returns all schema templates
func GetSchemaTemplates() map[string]string {
	return map[string]string{
		"model":      SchemaModelTemplate,
		"repository": SchemaRepositoryTemplate,
		"service":    SchemaServiceTemplate,
		"handler":    SchemaHandlerTemplate,
	}
}