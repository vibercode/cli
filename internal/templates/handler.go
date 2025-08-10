package templates

const HandlerTemplate = `package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"{{.Module}}/internal/models"
	"{{.Module}}/internal/services"
)

// {{.NameVariations.Pascal}}Handler handles HTTP requests for {{.NameVariations.PluralLower}}
type {{.NameVariations.Pascal}}Handler struct {
	service *services.{{.NameVariations.Pascal}}Service
}

// New{{.NameVariations.Pascal}}Handler creates a new {{.NameVariations.Pascal}}Handler
func New{{.NameVariations.Pascal}}Handler(service *services.{{.NameVariations.Pascal}}Service) *{{.NameVariations.Pascal}}Handler {
	return &{{.NameVariations.Pascal}}Handler{
		service: service,
	}
}

// Create{{.NameVariations.Pascal}} handles POST /{{.NameVariations.PluralKebab}}
func (h *{{.NameVariations.Pascal}}Handler) Create{{.NameVariations.Pascal}}(c *gin.Context) {
	var req models.{{.NameVariations.Pascal}}CreateRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	{{.NameVariations.Camel}} := models.{{.NameVariations.Pascal}}{
{{- range .Fields}}
	{{- if ne .Type "relation"}}
		{{.Name | ToCamel}}: req.{{.Name | ToCamel}},
	{{- end}}
{{- end}}
	}

	if err := h.service.Create{{.NameVariations.Pascal}}(c.Request.Context(), &{{.NameVariations.Camel}}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, {{.NameVariations.Camel}}.ToResponse())
}

// Get{{.NameVariations.Pascal}} handles GET /{{.NameVariations.PluralKebab}}/:id
func (h *{{.NameVariations.Pascal}}Handler) Get{{.NameVariations.Pascal}}(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	{{.NameVariations.Camel}}, err := h.service.Get{{.NameVariations.Pascal}}ByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "{{.NameVariations.Pascal}} not found"})
		return
	}

	c.JSON(http.StatusOK, {{.NameVariations.Camel}}.ToResponse())
}

// Get{{.NameVariations.PluralPascal}} handles GET /{{.NameVariations.PluralKebab}}
func (h *{{.NameVariations.Pascal}}Handler) Get{{.NameVariations.PluralPascal}}(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	{{.NameVariations.PluralCamel}}, total, err := h.service.Get{{.NameVariations.PluralPascal}}(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]models.{{.NameVariations.Pascal}}Response, len({{.NameVariations.PluralCamel}}))
	for i, {{.NameVariations.Camel}} := range {{.NameVariations.PluralCamel}} {
		responses[i] = {{.NameVariations.Camel}}.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  responses,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// Update{{.NameVariations.Pascal}} handles PUT /{{.NameVariations.PluralKebab}}/:id
func (h *{{.NameVariations.Pascal}}Handler) Update{{.NameVariations.Pascal}}(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req models.{{.NameVariations.Pascal}}UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	{{.NameVariations.Camel}}, err := h.service.Update{{.NameVariations.Pascal}}(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, {{.NameVariations.Camel}}.ToResponse())
}

// Delete{{.NameVariations.Pascal}} handles DELETE /{{.NameVariations.PluralKebab}}/:id
func (h *{{.NameVariations.Pascal}}Handler) Delete{{.NameVariations.Pascal}}(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.Delete{{.NameVariations.Pascal}}(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// Setup{{.NameVariations.Pascal}}Routes sets up the routes for {{.NameVariations.Pascal}}Handler
func (h *{{.NameVariations.Pascal}}Handler) Setup{{.NameVariations.Pascal}}Routes(r *gin.RouterGroup) {
	{{.NameVariations.PluralCamel}} := r.Group("/{{.NameVariations.PluralKebab}}")
	{
		{{.NameVariations.PluralCamel}}.POST("", h.Create{{.NameVariations.Pascal}})
		{{.NameVariations.PluralCamel}}.GET("/:id", h.Get{{.NameVariations.Pascal}})
		{{.NameVariations.PluralCamel}}.GET("", h.Get{{.NameVariations.PluralPascal}})
		{{.NameVariations.PluralCamel}}.PUT("/:id", h.Update{{.NameVariations.Pascal}})
		{{.NameVariations.PluralCamel}}.DELETE("/:id", h.Delete{{.NameVariations.Pascal}})
	}
}
`