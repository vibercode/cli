package templates

// OpenAPISpecTemplate generates OpenAPI 3.0 specification
const OpenAPISpecTemplate = `openapi: 3.0.3
info:
  title: {{.ProjectName}} API
  description: Auto-generated API documentation for {{.ProjectName}}
  version: 1.0.0
  contact:
    name: API Support
    email: support@{{.ProjectName}}.com

servers:
  - url: http://localhost:{{.Port}}/api/v1
    description: Development server

paths:
  /health:
    get:
      summary: Health check
      operationId: healthCheck
      tags: [System]
      responses:
        '200':
          description: API is healthy
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string
                    example: ok

{{range .Resources}}
  /{{.Names.KebabPlural}}:
    get:
      summary: List {{.Names.PluralTitle}}
      operationId: list{{.Names.PluralTitle}}
      tags: [{{.Names.Title}}]
      parameters:
        - name: page
          in: query
          schema:
            type: integer
            minimum: 1
            default: 1
        - name: limit
          in: query
          schema:
            type: integer
            minimum: 1
            maximum: 100
            default: 10
      responses:
        '200':
          description: List of {{.Names.LowerPlural}}
          content:
            application/json:
              schema:
                type: object
                properties:
                  data:
                    type: array
                    items:
                      $ref: '#/components/schemas/{{.Names.Title}}'
                  pagination:
                    $ref: '#/components/schemas/Pagination'

    post:
      summary: Create {{.Names.Title}}
      operationId: create{{.Names.Title}}
      tags: [{{.Names.Title}}]
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/{{.Names.Title}}Request'
      responses:
        '201':
          description: {{.Names.Title}} created successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/{{.Names.Title}}'

  /{{.Names.KebabPlural}}/{id}:
    get:
      summary: Get {{.Names.Title}}
      operationId: get{{.Names.Title}}
      tags: [{{.Names.Title}}]
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: {{.Names.Title}} details
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/{{.Names.Title}}'

    put:
      summary: Update {{.Names.Title}}
      operationId: update{{.Names.Title}}
      tags: [{{.Names.Title}}]
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/{{.Names.Title}}Request'
      responses:
        '200':
          description: {{.Names.Title}} updated successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/{{.Names.Title}}'

    delete:
      summary: Delete {{.Names.Title}}
      operationId: delete{{.Names.Title}}
      tags: [{{.Names.Title}}]
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '204':
          description: {{.Names.Title}} deleted successfully
{{end}}

components:
  schemas:
{{range .Resources}}
    {{.Names.Title}}:
      type: object
      properties:
        id:
          type: string
          format: uuid
          readOnly: true
{{range .Fields}}
        {{.Names.Snake}}:
          {{- if eq .Type "string" "text"}}
          type: string
          {{- else if eq .Type "number"}}
          type: integer
          {{- else if eq .Type "float"}}
          type: number
          {{- else if eq .Type "boolean"}}
          type: boolean
          {{- else if eq .Type "date"}}
          type: string
          format: date-time
          {{- else if eq .Type "uuid"}}
          type: string
          format: uuid
          {{- else if eq .Type "json"}}
          type: object
          {{- else if eq .Type "relation"}}
          type: string
          format: uuid
          {{- else if eq .Type "relation-array"}}
          type: array
          items:
            type: string
            format: uuid
          {{- else}}
          type: string
          {{- end}}
{{end}}
        created_at:
          type: string
          format: date-time
          readOnly: true
        updated_at:
          type: string
          format: date-time
          readOnly: true

    {{.Names.Title}}Request:
      type: object
      properties:
{{range .Fields}}
        {{.Names.Snake}}:
          {{- if eq .Type "string" "text"}}
          type: string
          {{- else if eq .Type "number"}}
          type: integer
          {{- else if eq .Type "float"}}
          type: number
          {{- else if eq .Type "boolean"}}
          type: boolean
          {{- else if eq .Type "date"}}
          type: string
          format: date-time
          {{- else if eq .Type "uuid"}}
          type: string
          format: uuid
          {{- else if eq .Type "json"}}
          type: object
          {{- else if eq .Type "relation"}}
          type: string
          format: uuid
          {{- else if eq .Type "relation-array"}}
          type: array
          items:
            type: string
            format: uuid
          {{- else}}
          type: string
          {{- end}}
{{end}}
{{end}}

    Pagination:
      type: object
      properties:
        current_page:
          type: integer
        per_page:
          type: integer
        total_pages:
          type: integer
        total_items:
          type: integer

tags:
  - name: System
    description: System endpoints
{{range .Resources}}
  - name: {{.Names.Title}}
    description: {{.Names.Title}} management
{{end}}
`

// SwaggerUITemplate generates the Swagger UI HTML page
const SwaggerUITemplate = `<!DOCTYPE html>
<html>
<head>
    <title>{{.ProjectName}} API Documentation</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui.css" />
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui-bundle.js"></script>
    <script>
        SwaggerUIBundle({
            url: '/api/v1/docs/openapi.yaml',
            dom_id: '#swagger-ui',
            presets: [SwaggerUIBundle.presets.apis],
            layout: "StandaloneLayout"
        });
    </script>
</body>
</html>`

// DocsHandlerTemplate generates the documentation handler
const DocsHandlerTemplate = `package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type DocsHandler struct{}

func NewDocsHandler() *DocsHandler {
	return &DocsHandler{}
}

func (h *DocsHandler) ServeSwaggerUI(c *gin.Context) {
	html := ` + "`{{.SwaggerUI}}`" + `
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

func (h *DocsHandler) ServeOpenAPISpec(c *gin.Context) {
	c.File("./docs/openapi.yaml")
}
`