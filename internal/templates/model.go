package templates

const ModelTemplate = `package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
{{- range .RequiredImports}}
	"{{.}}"
{{- end}}
)

// {{.NameVariations.Pascal}} represents a {{.NameVariations.Lower}} entity
type {{.NameVariations.Pascal}} struct {
	ID        primitive.ObjectID ` + "`" + `json:"id" bson:"_id,omitempty"` + "`" + `
	CreatedAt time.Time          ` + "`" + `json:"created_at" bson:"created_at"` + "`" + `
	UpdatedAt time.Time          ` + "`" + `json:"updated_at" bson:"updated_at"` + "`" + `
{{- range .Fields}}
	{{.GoStructField}}
{{- end}}
}

// {{.NameVariations.Pascal}}CreateRequest represents the request payload for creating a {{.NameVariations.Lower}}
type {{.NameVariations.Pascal}}CreateRequest struct {
{{- range .Fields}}
	{{- if ne .Type "relation"}}
	{{.GoStructField}}
	{{- end}}
{{- end}}
}

// {{.NameVariations.Pascal}}UpdateRequest represents the request payload for updating a {{.NameVariations.Lower}}
type {{.NameVariations.Pascal}}UpdateRequest struct {
{{- range .Fields}}
	{{- if ne .Type "relation"}}
	{{.Name | ToCamel}} *{{.GoType}} ` + "`" + `json:"{{.Name | ToSnake}},omitempty"` + "`" + `
	{{- end}}
{{- end}}
}

// {{.NameVariations.Pascal}}Response represents the response payload for {{.NameVariations.Lower}} operations
type {{.NameVariations.Pascal}}Response struct {
	ID        primitive.ObjectID ` + "`" + `json:"id"` + "`" + `
	CreatedAt time.Time          ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time          ` + "`" + `json:"updated_at"` + "`" + `
{{- range .Fields}}
	{{.GoStructField}}
{{- end}}
}

// CollectionName returns the MongoDB collection name for {{.NameVariations.Pascal}}
func ({{.NameVariations.Pascal}}) CollectionName() string {
	return "{{.TableName}}"
}

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

// Validate validates the {{.NameVariations.Pascal}} fields
func ({{.NameVariations.Camel}} *{{.NameVariations.Pascal}}) Validate() error {
{{- range .Fields}}
	{{- if .Required}}
	{{.GoValidation}}
	{{- end}}
{{- end}}
	return nil
}
`