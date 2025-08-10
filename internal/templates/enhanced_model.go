package templates

const EnhancedModelTemplate = `package models

import (
{{- if .DatabaseProvider.Type == "mongodb"}}
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
{{- else}}
	"time"
	"gorm.io/gorm"
{{- end}}
{{- range .RequiredImports}}
	"{{.}}"
{{- end}}
)

{{- range .Fields}}
{{- if eq .Type "enum"}}
{{.GenerateEnumType}}

{{- end}}
{{- end}}

{{- $hasCoordinates := false}}
{{- range .Fields}}
{{- if eq .Type "coordinates"}}
{{- $hasCoordinates = true}}
{{- end}}
{{- end}}

{{- if $hasCoordinates}}
{{GenerateCoordinatesStruct}}

{{- end}}

// {{.NameVariations.Pascal}} represents a {{.NameVariations.Lower}} entity
type {{.NameVariations.Pascal}} struct {
{{- if .DatabaseProvider.Type == "mongodb"}}
	ID        primitive.ObjectID ` + "`" + `json:"id" bson:"_id,omitempty"` + "`" + `
	CreatedAt time.Time          ` + "`" + `json:"created_at" bson:"created_at"` + "`" + `
	UpdatedAt time.Time          ` + "`" + `json:"updated_at" bson:"updated_at"` + "`" + `
{{- else}}
	ID        uint           ` + "`" + `json:"id" gorm:"primaryKey"` + "`" + `
	CreatedAt time.Time      ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time      ` + "`" + `json:"updated_at"` + "`" + `
	DeletedAt gorm.DeletedAt ` + "`" + `json:"-" gorm:"index"` + "`" + `
{{- end}}
{{- range .Fields}}
	{{.GoStructField}}
{{- end}}
}

// {{.NameVariations.Pascal}}CreateRequest represents the request payload for creating a {{.NameVariations.Lower}}
type {{.NameVariations.Pascal}}CreateRequest struct {
{{- range .Fields}}
	{{- if and (ne .Type "relation") (ne .Type "relation-array")}}
	{{.GoStructField}}
	{{- end}}
{{- end}}
}

// {{.NameVariations.Pascal}}UpdateRequest represents the request payload for updating a {{.NameVariations.Lower}}
type {{.NameVariations.Pascal}}UpdateRequest struct {
{{- range .Fields}}
	{{- if and (ne .Type "relation") (ne .Type "relation-array")}}
	{{.Name | ToCamel}} *{{.GoType}} ` + "`" + `json:"{{.Name | ToSnake}},omitempty"{{- if .Required}} binding:"required"{{- end}}` + "`" + `
	{{- end}}
{{- end}}
}

// {{.NameVariations.Pascal}}Response represents the response payload for {{.NameVariations.Lower}} operations
type {{.NameVariations.Pascal}}Response struct {
{{- if .DatabaseProvider.Type == "mongodb"}}
	ID        primitive.ObjectID ` + "`" + `json:"id"` + "`" + `
{{- else}}
	ID        uint               ` + "`" + `json:"id"` + "`" + `
{{- end}}
	CreatedAt time.Time          ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time          ` + "`" + `json:"updated_at"` + "`" + `
{{- range .Fields}}
	{{.GoStructField}}
{{- end}}
}

{{- if .DatabaseProvider.Type == "mongodb"}}
// CollectionName returns the MongoDB collection name for {{.NameVariations.Pascal}}
func ({{.NameVariations.Pascal}}) CollectionName() string {
	return "{{.TableName}}"
}
{{- else}}
// TableName returns the database table name for {{.NameVariations.Pascal}}
func ({{.NameVariations.Pascal}}) TableName() string {
	return "{{.TableName}}"
}
{{- end}}

// ToResponse converts {{.NameVariations.Pascal}} to {{.NameVariations.Pascal}}Response
func ({{.NameVariations.Camel}} *{{.NameVariations.Pascal}}) ToResponse() {{.NameVariations.Pascal}}Response {
	return {{.NameVariations.Pascal}}Response{
		ID:        {{.NameVariations.Camel}}.ID,
		CreatedAt: {{.NameVariations.Camel}}.CreatedAt,
		UpdatedAt: {{.NameVariations.Camel}}.UpdatedAt,
{{- range .Fields}}
		{{.Name | ToCamel}}: {{$.NameVariations.Camel}}.{{.Name | ToCamel}},
{{- end}}
	}
}

// FromCreateRequest populates {{.NameVariations.Pascal}} from {{.NameVariations.Pascal}}CreateRequest
func ({{.NameVariations.Camel}} *{{.NameVariations.Pascal}}) FromCreateRequest(req {{.NameVariations.Pascal}}CreateRequest) {
{{- range .Fields}}
	{{- if and (ne .Type "relation") (ne .Type "relation-array")}}
	{{$.NameVariations.Camel}}.{{.Name | ToCamel}} = req.{{.Name | ToCamel}}
	{{- end}}
{{- end}}
}

// UpdateFromRequest updates {{.NameVariations.Pascal}} from {{.NameVariations.Pascal}}UpdateRequest
func ({{.NameVariations.Camel}} *{{.NameVariations.Pascal}}) UpdateFromRequest(req {{.NameVariations.Pascal}}UpdateRequest) {
{{- range .Fields}}
	{{- if and (ne .Type "relation") (ne .Type "relation-array")}}
	if req.{{.Name | ToCamel}} != nil {
		{{$.NameVariations.Camel}}.{{.Name | ToCamel}} = *req.{{.Name | ToCamel}}
	}
	{{- end}}
{{- end}}
}

// Validate validates the {{.NameVariations.Pascal}} fields
func ({{.NameVariations.Camel}} *{{.NameVariations.Pascal}}) Validate() error {
{{- range .Fields}}
	{{.GoValidation}}

{{- end}}
	return nil
}

// ValidateCreate validates the {{.NameVariations.Pascal}}CreateRequest
func (req *{{.NameVariations.Pascal}}CreateRequest) Validate() error {
{{- range .Fields}}
	{{- if and (ne .Type "relation") (ne .Type "relation-array")}}
	{{.GoValidation}}

	{{- end}}
{{- end}}
	return nil
}

// ValidateUpdate validates the {{.NameVariations.Pascal}}UpdateRequest
func (req *{{.NameVariations.Pascal}}UpdateRequest) Validate() error {
{{- range .Fields}}
	{{- if and (ne .Type "relation") (ne .Type "relation-array")}}
	if req.{{.Name | ToCamel}} != nil {
		{{.Name | ToCamel}} := *req.{{.Name | ToCamel}}
		{{.GoValidation}}
	}

	{{- end}}
{{- end}}
	return nil
}

{{- range .Fields}}
{{- if eq .Type "slug"}}

// Generate{{.Name | ToCamel}} generates a slug from the {{.NameVariations.Lower}} name or title
func ({{$.NameVariations.Camel}} *{{$.NameVariations.Pascal}}) Generate{{.Name | ToCamel}}() {
	// This would typically generate from a name or title field
	// Implementation depends on source field
	{{$.NameVariations.Camel}}.{{.Name | ToCamel}} = "generated-slug"
}
{{- end}}
{{- if eq .Type "password"}}

// Set{{.Name | ToCamel}} hashes and sets the password
func ({{$.NameVariations.Camel}} *{{$.NameVariations.Pascal}}) Set{{.Name | ToCamel}}(password string) error {
	// Hash password using bcrypt or similar
	// This is a placeholder - implement proper password hashing
	{{$.NameVariations.Camel}}.{{.Name | ToCamel}} = password
	return nil
}

// Check{{.Name | ToCamel}} checks if the provided password matches
func ({{$.NameVariations.Camel}} *{{$.NameVariations.Pascal}}) Check{{.Name | ToCamel}}(password string) bool {
	// Compare hashed password
	// This is a placeholder - implement proper password checking
	return {{$.NameVariations.Camel}}.{{.Name | ToCamel}} == password
}
{{- end}}
{{- end}}
`

const EnhancedHandlerTemplate = `package handlers

import (
	"net/http"
	"strconv"
{{- if .DatabaseProvider.Type == "mongodb"}}
	"go.mongodb.org/mongo-driver/bson/primitive"
{{- end}}

	"github.com/gin-gonic/gin"
	"{{.Module}}/internal/models"
	"{{.Module}}/internal/services"
)

type {{.NameVariations.Pascal}}Handler struct {
	service *services.{{.NameVariations.Pascal}}Service
}

func New{{.NameVariations.Pascal}}Handler(service *services.{{.NameVariations.Pascal}}Service) *{{.NameVariations.Pascal}}Handler {
	return &{{.NameVariations.Pascal}}Handler{
		service: service,
	}
}

// Create{{.NameVariations.Pascal}} creates a new {{.NameVariations.Lower}}
// @Summary Create {{.NameVariations.Lower}}
// @Description Create a new {{.NameVariations.Lower}}
// @Tags {{.NameVariations.PluralLower}}
// @Accept json
// @Produce json
// @Param {{.NameVariations.Lower}} body models.{{.NameVariations.Pascal}}CreateRequest true "{{.NameVariations.Pascal}} data"
// @Success 201 {object} models.{{.NameVariations.Pascal}}Response
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /{{.NameVariations.PluralKebab}} [post]
func (h *{{.NameVariations.Pascal}}Handler) Create{{.NameVariations.Pascal}}(c *gin.Context) {
	var req models.{{.NameVariations.Pascal}}CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	{{.NameVariations.Lower}}, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, {{.NameVariations.Lower}}.ToResponse())
}

// Get{{.NameVariations.Pascal}} retrieves a {{.NameVariations.Lower}} by ID
// @Summary Get {{.NameVariations.Lower}}
// @Description Get a {{.NameVariations.Lower}} by ID
// @Tags {{.NameVariations.PluralLower}}
// @Produce json
// @Param id path string true "{{.NameVariations.Pascal}} ID"
// @Success 200 {object} models.{{.NameVariations.Pascal}}Response
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /{{.NameVariations.PluralKebab}}/{id} [get]
func (h *{{.NameVariations.Pascal}}Handler) Get{{.NameVariations.Pascal}}(c *gin.Context) {
	idStr := c.Param("id")
{{- if .DatabaseProvider.Type == "mongodb"}}
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
{{- else}}
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
{{- end}}

	{{.NameVariations.Lower}}, err := h.service.GetByID(c.Request.Context(), {{- if ne .DatabaseProvider.Type "mongodb"}}uint(id){{- else}}id{{- end}})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "{{.NameVariations.Pascal}} not found"})
		return
	}

	c.JSON(http.StatusOK, {{.NameVariations.Lower}}.ToResponse())
}

// List{{.NameVariations.PluralPascal}} retrieves all {{.NameVariations.PluralLower}}
// @Summary List {{.NameVariations.PluralLower}}
// @Description Get all {{.NameVariations.PluralLower}}
// @Tags {{.NameVariations.PluralLower}}
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
{{- range .Fields}}
{{- if or (eq .Type "string") (eq .Type "email") (eq .Type "slug")}}
// @Param {{.Name}} query string false "Filter by {{.DisplayName}}"
{{- end}}
{{- end}}
// @Success 200 {array} models.{{.NameVariations.Pascal}}Response
// @Failure 500 {object} map[string]interface{}
// @Router /{{.NameVariations.PluralKebab}} [get]
func (h *{{.NameVariations.Pascal}}Handler) List{{.NameVariations.PluralPascal}}(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Parse filter parameters
	filters := make(map[string]interface{})
{{- range .Fields}}
{{- if or (eq .Type "string") (eq .Type "email") (eq .Type "slug")}}
	if {{.Name}} := c.Query("{{.Name}}"); {{.Name}} != "" {
		filters["{{.Name}}"] = {{.Name}}
	}
{{- end}}
{{- end}}

	{{.NameVariations.PluralLower}}, total, err := h.service.List(c.Request.Context(), filters, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []models.{{.NameVariations.Pascal}}Response
	for _, {{.NameVariations.Lower}} := range {{.NameVariations.PluralLower}} {
		responses = append(responses, {{.NameVariations.Lower}}.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  responses,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// Update{{.NameVariations.Pascal}} updates a {{.NameVariations.Lower}}
// @Summary Update {{.NameVariations.Lower}}
// @Description Update a {{.NameVariations.Lower}} by ID
// @Tags {{.NameVariations.PluralLower}}
// @Accept json
// @Produce json
// @Param id path string true "{{.NameVariations.Pascal}} ID"
// @Param {{.NameVariations.Lower}} body models.{{.NameVariations.Pascal}}UpdateRequest true "{{.NameVariations.Pascal}} data"
// @Success 200 {object} models.{{.NameVariations.Pascal}}Response
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /{{.NameVariations.PluralKebab}}/{id} [put]
func (h *{{.NameVariations.Pascal}}Handler) Update{{.NameVariations.Pascal}}(c *gin.Context) {
	idStr := c.Param("id")
{{- if .DatabaseProvider.Type == "mongodb"}}
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
{{- else}}
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
{{- end}}

	var req models.{{.NameVariations.Pascal}}UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	{{.NameVariations.Lower}}, err := h.service.Update(c.Request.Context(), {{- if ne .DatabaseProvider.Type "mongodb"}}uint(id){{- else}}id{{- end}}, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, {{.NameVariations.Lower}}.ToResponse())
}

// Delete{{.NameVariations.Pascal}} deletes a {{.NameVariations.Lower}}
// @Summary Delete {{.NameVariations.Lower}}
// @Description Delete a {{.NameVariations.Lower}} by ID
// @Tags {{.NameVariations.PluralLower}}
// @Param id path string true "{{.NameVariations.Pascal}} ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /{{.NameVariations.PluralKebab}}/{id} [delete]
func (h *{{.NameVariations.Pascal}}Handler) Delete{{.NameVariations.Pascal}}(c *gin.Context) {
	idStr := c.Param("id")
{{- if .DatabaseProvider.Type == "mongodb"}}
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
{{- else}}
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
{{- end}}

	err = h.service.Delete(c.Request.Context(), {{- if ne .DatabaseProvider.Type "mongodb"}}uint(id){{- else}}id{{- end}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Setup{{.NameVariations.Pascal}}Routes sets up the routes for {{.NameVariations.Lower}} handlers
func (h *{{.NameVariations.Pascal}}Handler) Setup{{.NameVariations.Pascal}}Routes(r *gin.RouterGroup) {
	{{.NameVariations.PluralLower}} := r.Group("/{{.NameVariations.PluralKebab}}")
	{
		{{.NameVariations.PluralLower}}.POST("", h.Create{{.NameVariations.Pascal}})
		{{.NameVariations.PluralLower}}.GET("", h.List{{.NameVariations.PluralPascal}})
		{{.NameVariations.PluralLower}}.GET("/:id", h.Get{{.NameVariations.Pascal}})
		{{.NameVariations.PluralLower}}.PUT("/:id", h.Update{{.NameVariations.Pascal}})
		{{.NameVariations.PluralLower}}.DELETE("/:id", h.Delete{{.NameVariations.Pascal}})
	}
}
`