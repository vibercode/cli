package templates

import (
	"fmt"
	"strings"

	"github.com/vibercode/cli/internal/models"
)

// EnhancedField is a temporary structure for template rendering
type EnhancedField struct {
	*models.SchemaField
	TypeScriptType string
	Filterable     bool
}

type FrontendGenerator struct {
	registry   *TemplateRegistry
	schema     *models.ResourceSchema
	outputDir  string
	framework  string
	apiBaseURL string
}

func NewFrontendGenerator(registry *TemplateRegistry, schema *models.ResourceSchema, outputDir, framework, apiBaseURL string) *FrontendGenerator {
	return &FrontendGenerator{
		registry:   registry,
		schema:     schema,
		outputDir:  outputDir,
		framework:  framework,
		apiBaseURL: apiBaseURL,
	}
}

func (g *FrontendGenerator) GenerateReactComponents() error {
	// Prepare variables for template
	variables := map[string]interface{}{
		"schema":      g.schema,
		"Names":       g.schema.Names,
		"Fields":      g.enhanceFieldsForReact(g.schema.Fields),
		"DisplayName": g.schema.DisplayName,
		"ApiBaseUrl":  g.apiBaseURL,
	}

	return g.registry.GenerateFromTemplate("react-components", g.outputDir, variables)
}

func (g *FrontendGenerator) GenerateVueComponents() error {
	// TODO: Implement Vue component generation
	return fmt.Errorf("Vue components generation not yet implemented")
}

func (g *FrontendGenerator) GenerateAngularComponents() error {
	// TODO: Implement Angular component generation
	return fmt.Errorf("Angular components generation not yet implemented")
}

func (g *FrontendGenerator) GenerateFullstackAI() error {
	templateID := "fullstack-ai"

	// Enhanced fields for React with TypeScript types
	enhancedFields := g.enhanceFieldsForReact(g.schema.Fields)

	// Create interface for AI generation
	interfaceFields := make([]string, 0, len(enhancedFields))
	for _, field := range enhancedFields {
		optional := ""
		if !field.Required {
			optional = "?"
		}

		interfaceFields = append(interfaceFields, fmt.Sprintf("  %s%s: %s;", field.Name, optional, field.TypeScriptType))
	}

	// Create form fields for AI generation
	formFields := make([]string, 0, len(enhancedFields))
	for _, field := range enhancedFields {
		validation := ""
		if field.Required {
			validation = ".required('This field is required')"
		}

		formFields = append(formFields, fmt.Sprintf("  %s: Yup.%s()%s,", field.Name, g.getYupType(field.TypeScriptType), validation))
	}

	// Create filter fields for AI generation
	filterFields := make([]string, 0)
	for _, field := range enhancedFields {
		if field.Filterable {
			filterFields = append(filterFields, fmt.Sprintf("  %s?: %s;", field.Name, field.TypeScriptType))
		}
	}

	// Prepare variables for template
	variables := map[string]interface{}{
		"schema":          g.schema,
		"Names":           g.schema.Names,
		"Fields":          enhancedFields,
		"DisplayName":     g.schema.DisplayName,
		"ApiBaseUrl":      g.apiBaseURL,
		"InterfaceFields": strings.Join(interfaceFields, "\n"),
		"FormFields":      strings.Join(formFields, "\n"),
		"FilterFields":    strings.Join(filterFields, "\n"),
	}

	return g.registry.GenerateFromTemplate(templateID, g.outputDir, variables)
}

// enhanceFieldsForReact enhances SchemaField with additional properties for React templates
func (g *FrontendGenerator) enhanceFieldsForReact(fields []models.SchemaField) []EnhancedField {
	enhancedFields := make([]EnhancedField, len(fields))

	for i, field := range fields {
		enhancedFields[i] = EnhancedField{
			SchemaField:    &field,
			TypeScriptType: g.getTypeScriptType(&field),
			Filterable:     g.isFilterable(&field),
		}
	}

	return enhancedFields
}

// getTypeScriptType returns the TypeScript type for a SchemaField
func (g *FrontendGenerator) getTypeScriptType(field *models.SchemaField) string {
	switch field.Type {
	case "string", "text", "email", "url", "password":
		return "string"
	case "number", "integer", "float", "decimal":
		return "number"
	case "boolean":
		return "boolean"
	case "date", "datetime", "timestamp":
		return "Date"
	case "json", "object":
		return "any"
	case "array":
		return "any[]"
	default:
		return "string"
	}
}

// isFilterable determines if a field should be filterable
func (g *FrontendGenerator) isFilterable(field *models.SchemaField) bool {
	switch field.Type {
	case "string", "text", "email", "url", "number", "integer", "float", "decimal", "boolean", "date", "datetime":
		return true
	default:
		return false
	}
}

// getYupType returns the Yup validation type for a TypeScript type
func (g *FrontendGenerator) getYupType(tsType string) string {
	switch tsType {
	case "string":
		return "string"
	case "number":
		return "number"
	case "boolean":
		return "boolean"
	case "Date":
		return "date"
	default:
		return "string"
	}
}

func (g *FrontendGenerator) generateFullstackConfig() error {
	// Generate package.json
	packageJSON := g.generatePackageJSON()

	// Generate TypeScript config
	tsConfig := g.generateTSConfig()

	// Generate Tailwind config
	tailwindConfig := g.generateTailwindConfig()

	// Generate Vite config
	viteConfig := g.generateViteConfig()

	// Write files
	// Note: In a real implementation, you would write these to actual files
	_ = packageJSON
	_ = tsConfig
	_ = tailwindConfig
	_ = viteConfig

	return nil
}

func (g *FrontendGenerator) generateDockerCompose() string {
	return fmt.Sprintf(`version: '3.8'

services:
  %s-frontend:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "3000:3000"
    environment:
      - VITE_API_BASE_URL=%s
    depends_on:
      - %s-backend
    volumes:
      - .:/app
      - /app/node_modules

  %s-backend:
    build:
      context: ../backend
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - DB_CONNECTION_STRING=postgresql://user:password@db:5432/%s
    depends_on:
      - db

  db:
    image: postgres:15
    environment:
      - POSTGRES_DB=%s
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

volumes:
  postgres_data:
`, g.schema.Names.KebabCase, g.apiBaseURL, g.schema.Names.KebabCase, g.schema.Names.KebabCase, g.schema.Names.SnakeCase, g.schema.Names.SnakeCase)
}

func (g *FrontendGenerator) generateReadme() string {
	return fmt.Sprintf(`# %s Frontend

This is a React + TypeScript frontend application for managing %s resources.

## Features

- React 18 with TypeScript
- Tailwind CSS for styling
- Responsive Design
- Advanced Filtering and Search
- Data Tables with sorting and pagination
- Form Validation with Yup
- Real-time Updates with WebSocket
- Dark Mode support
- Component Library with Storybook
- Testing with Jest and React Testing Library
- CI/CD ready

## Getting Started

### Prerequisites

- Node.js 18+ 
- npm or yarn

### Installation

1. Clone the repository
2. Install dependencies: npm install
3. Set up environment variables: cp .env.example .env
4. Start the development server: npm run dev

### Available Scripts

- npm run dev - Start development server
- npm run build - Build for production
- npm run preview - Preview production build
- npm run test - Run tests
- npm run test:coverage - Run tests with coverage
- npm run lint - Run ESLint
- npm run type-check - Run TypeScript compiler

## Architecture

### Components

#### Atoms
- Input components (text, number, date, etc.)
- Button components
- Typography components

#### Molecules
- Form field components
- Filter components
- Card components

#### Organisms
- Data table
- Form layouts
- Navigation components

#### Templates
- Page layouts
- Modal layouts

#### Pages
- List view
- Detail view
- Create/Edit forms

### State Management

- React Query for server state
- Zustand for client state
- Form state with React Hook Form

### Styling

- Tailwind CSS for utility-first styling
- Headless UI for accessible components
- Heroicons for icons

### API Integration

- Axios for HTTP requests
- React Query for caching and synchronization
- WebSocket for real-time updates

## Field Types

The following field types are supported:

%s

## Development

### Project Structure

src/
|-- components/          # Reusable UI components
|   |-- atoms/          # Basic building blocks
|   |-- molecules/      # Composed components
|   |-- organisms/      # Complex UI components
|   |-- templates/      # Page layouts
|-- pages/              # Route components
|-- hooks/              # Custom React hooks
|-- services/           # API services
|-- types/              # TypeScript type definitions
|-- utils/              # Utility functions
|-- styles/             # Global styles
|-- tests/              # Test files

### Adding New Fields

1. Update the schema definition
2. Regenerate components with: vibercode generate frontend
3. Customize generated components as needed

### Customization

Generated components can be customized by:

1. Modifying the generated files directly
2. Overriding styles with custom CSS
3. Extending component props
4. Adding custom validation rules

## Deployment

### Production Build

npm run build

### Docker

docker build -t %s-frontend .
docker run -p 3000:3000 %s-frontend

### Environment Variables

- VITE_API_BASE_URL - Backend API URL
- VITE_WS_URL - WebSocket URL
- VITE_APP_NAME - Application name
- VITE_ENABLE_ANALYTICS - Enable analytics (true/false)

## License

This project is licensed under the MIT License.
`,
		g.schema.DisplayName,
		g.schema.Names.Plural,
		g.generateFieldDocumentation(),
		g.schema.Names.KebabCase,
		g.schema.Names.KebabCase,
	)
}

func (g *FrontendGenerator) generateFieldDocumentation() string {
	var docs []string

	for _, field := range g.schema.Fields {
		fieldDoc := fmt.Sprintf("- **%s** (%s)", field.Name, field.Type)
		if field.Description != "" {
			fieldDoc += fmt.Sprintf(" - %s", field.Description)
		}
		docs = append(docs, fieldDoc)
	}

	return strings.Join(docs, "\n")
}

func (g *FrontendGenerator) generatePackageJSON() string {
	return `{
  "name": "frontend-app",
  "version": "1.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "tsc && vite build",
    "preview": "vite preview"
  },
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0"
  },
  "devDependencies": {
    "@types/react": "^18.2.43",
    "@types/react-dom": "^18.2.17",
    "@vitejs/plugin-react": "^4.2.1",
    "typescript": "^5.2.2",
    "vite": "^5.0.8"
  }
}`
}

func (g *FrontendGenerator) generateTSConfig() string {
	return `{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "react-jsx",
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true
  },
  "include": ["src"],
  "references": [{ "path": "./tsconfig.node.json" }]
}`
}

func (g *FrontendGenerator) generateTailwindConfig() string {
	return `/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}`
}

func (g *FrontendGenerator) generateViteConfig() string {
	return `import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
})`
}

func (g *FrontendGenerator) GenerateComponents() error {
	switch g.framework {
	case "react":
		return g.GenerateReactComponents()
	case "vue":
		return g.GenerateVueComponents()
	case "angular":
		return g.GenerateAngularComponents()
	default:
		return fmt.Errorf("unsupported framework: %s", g.framework)
	}
}
