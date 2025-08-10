package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/vibercode/cli/internal/models"
	"github.com/vibercode/cli/internal/templates"
	"github.com/vibercode/cli/pkg/ui"
)

// APIDocsGenerator handles API documentation generation
type APIDocsGenerator struct{}

// NewAPIDocsGenerator creates a new APIDocsGenerator
func NewAPIDocsGenerator() *APIDocsGenerator {
	return &APIDocsGenerator{}
}

// DocsProject represents a documentation project configuration
type DocsProject struct {
	ProjectName string
	Port        string
	Resources   []*models.Resource
}

// GenerateAPIDocs generates complete API documentation
func (g *APIDocsGenerator) GenerateAPIDocs(projectName, port string, resources []*models.Resource, outputDir string) error {
	project := &DocsProject{
		ProjectName: projectName,
		Port:        port,
		Resources:   resources,
	}

	ui.PrintSuccess("📚 Generating API documentation...")

	// Create docs directory
	docsDir := filepath.Join(outputDir, "docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		return fmt.Errorf("failed to create docs directory: %w", err)
	}

	// Generate OpenAPI specification
	if err := g.generateOpenAPISpec(project, docsDir); err != nil {
		return fmt.Errorf("failed to generate OpenAPI spec: %w", err)
	}

	// Generate Swagger UI handler
	if err := g.generateDocsHandler(project, outputDir); err != nil {
		return fmt.Errorf("failed to generate docs handler: %w", err)
	}

	ui.PrintSuccess("✅ API documentation generated successfully!")
	ui.PrintInfo(fmt.Sprintf("📄 OpenAPI spec: %s/openapi.yaml", docsDir))
	ui.PrintInfo(fmt.Sprintf("🌐 Swagger UI: http://localhost:%s/api/v1/docs", port))

	return nil
}

// generateOpenAPISpec generates the OpenAPI specification file
func (g *APIDocsGenerator) generateOpenAPISpec(project *DocsProject, outputDir string) error {
	tmpl, err := template.New("openapi").Parse(templates.OpenAPISpecTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse OpenAPI template: %w", err)
	}

	outputPath := filepath.Join(outputDir, "openapi.yaml")
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create OpenAPI spec file: %w", err)
	}
	defer file.Close()

	return tmpl.Execute(file, project)
}

// generateDocsHandler generates the documentation handler
func (g *APIDocsGenerator) generateDocsHandler(project *DocsProject, outputDir string) error {
	handlerDir := filepath.Join(outputDir, "internal", "handlers")
	if err := os.MkdirAll(handlerDir, 0755); err != nil {
		return fmt.Errorf("failed to create handlers directory: %w", err)
	}

	outputPath := filepath.Join(handlerDir, "docs_handler.go")
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create docs handler file: %w", err)
	}
	defer file.Close()

	tmpl, err := template.New("docs-handler").Parse(templates.DocsHandlerTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse docs handler template: %w", err)
	}

	return tmpl.Execute(file, project)
}