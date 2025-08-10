package templates

const ServiceTemplate = `package services

import (
	"context"
	"errors"

	"{{.Module}}/internal/models"
	"{{.Module}}/internal/repositories"
)

// {{.NameVariations.Pascal}}Service handles business logic for {{.NameVariations.PluralLower}}
type {{.NameVariations.Pascal}}Service struct {
	repo *repositories.{{.NameVariations.Pascal}}Repository
}

// New{{.NameVariations.Pascal}}Service creates a new {{.NameVariations.Pascal}}Service
func New{{.NameVariations.Pascal}}Service(repo *repositories.{{.NameVariations.Pascal}}Repository) *{{.NameVariations.Pascal}}Service {
	return &{{.NameVariations.Pascal}}Service{
		repo: repo,
	}
}

// Create{{.NameVariations.Pascal}} creates a new {{.NameVariations.Lower}}
func (s *{{.NameVariations.Pascal}}Service) Create{{.NameVariations.Pascal}}(ctx context.Context, {{.NameVariations.Camel}} *models.{{.NameVariations.Pascal}}) error {
	if err := {{.NameVariations.Camel}}.Validate(); err != nil {
		return err
	}

	return s.repo.Create(ctx, {{.NameVariations.Camel}})
}

// Get{{.NameVariations.Pascal}}ByID retrieves a {{.NameVariations.Lower}} by ID
func (s *{{.NameVariations.Pascal}}Service) Get{{.NameVariations.Pascal}}ByID(ctx context.Context, id string) (*models.{{.NameVariations.Pascal}}, error) {
	{{.NameVariations.Camel}}, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if {{.NameVariations.Camel}} == nil {
		return nil, errors.New("{{.NameVariations.Lower}} not found")
	}

	return {{.NameVariations.Camel}}, nil
}

// Get{{.NameVariations.PluralPascal}} retrieves {{.NameVariations.PluralLower}} with pagination
func (s *{{.NameVariations.Pascal}}Service) Get{{.NameVariations.PluralPascal}}(ctx context.Context, page, limit int) ([]models.{{.NameVariations.Pascal}}, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit
	
	{{.NameVariations.PluralCamel}}, err := s.repo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return {{.NameVariations.PluralCamel}}, total, nil
}

// Update{{.NameVariations.Pascal}} updates a {{.NameVariations.Lower}}
func (s *{{.NameVariations.Pascal}}Service) Update{{.NameVariations.Pascal}}(ctx context.Context, id string, req *models.{{.NameVariations.Pascal}}UpdateRequest) (*models.{{.NameVariations.Pascal}}, error) {
	{{.NameVariations.Camel}}, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if {{.NameVariations.Camel}} == nil {
		return nil, errors.New("{{.NameVariations.Lower}} not found")
	}

	// Update fields if provided
{{- range .Fields}}
	{{- if ne .Type "relation"}}
	if req.{{.Name | ToCamel}} != nil {
		{{$.NameVariations.Camel}}.{{.Name | ToCamel}} = *req.{{.Name | ToCamel}}
	}
	{{- end}}
{{- end}}

	if err := {{.NameVariations.Camel}}.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, {{.NameVariations.Camel}}); err != nil {
		return nil, err
	}

	return {{.NameVariations.Camel}}, nil
}

// Delete{{.NameVariations.Pascal}} deletes a {{.NameVariations.Lower}}
func (s *{{.NameVariations.Pascal}}Service) Delete{{.NameVariations.Pascal}}(ctx context.Context, id string) error {
	{{.NameVariations.Camel}}, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if {{.NameVariations.Camel}} == nil {
		return errors.New("{{.NameVariations.Lower}} not found")
	}

	return s.repo.Delete(ctx, id)
}

// Search{{.NameVariations.PluralPascal}} searches {{.NameVariations.PluralLower}} by criteria
func (s *{{.NameVariations.Pascal}}Service) Search{{.NameVariations.PluralPascal}}(ctx context.Context, query string, page, limit int) ([]models.{{.NameVariations.Pascal}}, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit
	
	{{.NameVariations.PluralCamel}}, err := s.repo.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountBySearch(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	return {{.NameVariations.PluralCamel}}, total, nil
}
`