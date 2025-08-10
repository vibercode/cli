package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"strings"

	"github.com/vibercode/cli/internal/models"
	"github.com/vibercode/cli/internal/templates"
	"github.com/iancoleman/strcase"
)

// AuthGenerator handles authentication system generation
type AuthGenerator struct {
	options    models.AuthGeneratorOptions
	templateFuncs template.FuncMap
}

// NewAuthGenerator creates a new authentication generator
func NewAuthGenerator(options models.AuthGeneratorOptions) *AuthGenerator {
	generator := &AuthGenerator{
		options: options,
	}

	// Setup template functions
	generator.templateFuncs = template.FuncMap{
		"title":      strings.Title,
		"lower":      strings.ToLower,
		"upper":      strings.ToUpper,
		"camelCase":  strcase.ToCamel,
		"snakeCase":  strcase.ToSnake,
		"kebabCase":  strcase.ToKebab,
		"contains":   generator.containsMethod,
		"toJSON":     generator.toJSON,
	}

	return generator
}

// GenerateAuthSystem generates complete authentication system
func (g *AuthGenerator) GenerateAuthSystem() error {
	// Validate options
	if err := models.ValidateAuthConfig(g.options.AuthConfig); err != nil {
		return fmt.Errorf("invalid auth config: %w", err)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(g.options.OutputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate core authentication files
	if err := g.generateMiddleware(); err != nil {
		return fmt.Errorf("failed to generate middleware: %w", err)
	}

	if err := g.generateHandlers(); err != nil {
		return fmt.Errorf("failed to generate handlers: %w", err)
	}

	if err := g.generateServices(); err != nil {
		return fmt.Errorf("failed to generate services: %w", err)
	}

	if err := g.generateModels(); err != nil {
		return fmt.Errorf("failed to generate models: %w", err)
	}

	if err := g.generateRepositories(); err != nil {
		return fmt.Errorf("failed to generate repositories: %w", err)
	}

	// Generate OAuth2 components if enabled
	if len(g.options.AuthConfig.OAuth2Providers) > 0 {
		if err := g.generateOAuth2Components(); err != nil {
			return fmt.Errorf("failed to generate OAuth2 components: %w", err)
		}
	}

	// Generate Supabase integration if enabled
	if g.options.AuthConfig.Provider == models.AuthProviderSupabase {
		if err := g.generateSupabaseIntegration(); err != nil {
			return fmt.Errorf("failed to generate Supabase integration: %w", err)
		}
	}

	// Generate migrations
	if err := g.generateMigrations(); err != nil {
		return fmt.Errorf("failed to generate migrations: %w", err)
	}

	// Generate route setup
	if err := g.generateRouteSetup(); err != nil {
		return fmt.Errorf("failed to generate route setup: %w", err)
	}

	return nil
}

// generateMiddleware generates JWT middleware
func (g *AuthGenerator) generateMiddleware() error {
	middlewarePath := filepath.Join(g.options.OutputPath, "internal", "middleware")
	if err := os.MkdirAll(middlewarePath, 0755); err != nil {
		return err
	}

	// Generate JWT middleware
	data := g.getTemplateData()
	return g.executeTemplate(
		templates.JWTMiddlewareTemplate,
		filepath.Join(middlewarePath, "jwt.go"),
		data,
	)
}

// generateHandlers generates authentication handlers
func (g *AuthGenerator) generateHandlers() error {
	handlersPath := filepath.Join(g.options.OutputPath, "internal", "handlers")
	if err := os.MkdirAll(handlersPath, 0755); err != nil {
		return err
	}

	// Generate auth handlers
	data := g.getTemplateData()
	return g.executeTemplate(
		templates.AuthHandlersTemplate,
		filepath.Join(handlersPath, "auth.go"),
		data,
	)
}

// generateServices generates authentication services
func (g *AuthGenerator) generateServices() error {
	servicesPath := filepath.Join(g.options.OutputPath, "internal", "services")
	if err := os.MkdirAll(servicesPath, 0755); err != nil {
		return err
	}

	// Generate auth service
	data := g.getTemplateData()
	return g.executeTemplate(
		templates.AuthServiceTemplate,
		filepath.Join(servicesPath, "auth.go"),
		data,
	)
}

// generateModels generates user and auth-related models
func (g *AuthGenerator) generateModels() error {
	modelsPath := filepath.Join(g.options.OutputPath, "internal", "models")
	if err := os.MkdirAll(modelsPath, 0755); err != nil {
		return err
	}

	// Generate user model
	data := g.getTemplateData()
	return g.executeTemplate(
		templates.UserModelTemplate,
		filepath.Join(modelsPath, "user.go"),
		data,
	)
}

// generateRepositories generates auth repositories
func (g *AuthGenerator) generateRepositories() error {
	repoPath := filepath.Join(g.options.OutputPath, "internal", "repositories")
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		return err
	}

	// Generate user repository
	data := g.getTemplateData()
	
	// Generate based on database provider
	switch g.options.DatabaseProvider {
	case "postgres", "mysql", "sqlite":
		return g.generateGORMRepository(repoPath, data)
	case "mongodb":
		return g.generateMongoRepository(repoPath, data)
	case "supabase":
		return g.generateSupabaseRepository(repoPath, data)
	default:
		return fmt.Errorf("unsupported database provider: %s", g.options.DatabaseProvider)
	}
}

// generateOAuth2Components generates OAuth2 middleware and handlers
func (g *AuthGenerator) generateOAuth2Components() error {
	// Generate OAuth2 middleware
	middlewarePath := filepath.Join(g.options.OutputPath, "internal", "middleware")
	data := g.getTemplateData()
	
	if err := g.executeTemplate(
		templates.OAuth2MiddlewareTemplate,
		filepath.Join(middlewarePath, "oauth2.go"),
		data,
	); err != nil {
		return err
	}

	// Generate OAuth2 handlers
	handlersPath := filepath.Join(g.options.OutputPath, "internal", "handlers")
	return g.executeTemplate(
		templates.OAuth2HandlersTemplate,
		filepath.Join(handlersPath, "oauth2.go"),
		data,
	)
}

// generateSupabaseIntegration generates Supabase authentication integration
func (g *AuthGenerator) generateSupabaseIntegration() error {
	servicesPath := filepath.Join(g.options.OutputPath, "internal", "services")
	data := g.getTemplateData()
	
	return g.executeTemplate(
		templates.SupabaseAuthTemplate,
		filepath.Join(servicesPath, "supabase_auth.go"),
		data,
	)
}

// generateMigrations generates database migrations for auth tables
func (g *AuthGenerator) generateMigrations() error {
	migrationsPath := filepath.Join(g.options.OutputPath, "migrations")
	if err := os.MkdirAll(migrationsPath, 0755); err != nil {
		return err
	}

	// Generate user table migration
	if err := g.generateUserMigration(migrationsPath); err != nil {
		return err
	}

	// Generate RBAC migrations if enabled
	if g.options.AuthConfig.EnableRBAC {
		if err := g.generateRBACMigrations(migrationsPath); err != nil {
			return err
		}
	}

	return nil
}

// generateRouteSetup generates route setup for authentication
func (g *AuthGenerator) generateRouteSetup() error {
	routesPath := filepath.Join(g.options.OutputPath, "internal", "routes")
	if err := os.MkdirAll(routesPath, 0755); err != nil {
		return err
	}

	data := g.getTemplateData()
	routeTemplate := g.generateRouteTemplate()
	
	return g.executeTemplate(
		routeTemplate,
		filepath.Join(routesPath, "auth.go"),
		data,
	)
}

// getTemplateData returns data for template execution
func (g *AuthGenerator) getTemplateData() map[string]interface{} {
	// Determine primary key type
	primaryKeyType := "string"
	if g.options.UserModel.PrimaryKey == "id" {
		primaryKeyType = "uint"
	}

	// Set struct name if not provided
	userModel := g.options.UserModel
	if userModel.StructName == "" {
		userModel.StructName = "User"
	}

	// Set role model struct name
	var roleModel models.RoleModel
	if g.options.RoleModel != nil {
		roleModel = *g.options.RoleModel
		if roleModel.StructName == "" {
			roleModel.StructName = "Role"
		}
	}

	// Set permission model struct name
	var permissionModel models.PermissionModel
	if g.options.PermissionModel != nil {
		permissionModel = *g.options.PermissionModel
		if permissionModel.StructName == "" {
			permissionModel.StructName = "Permission"
		}
	}

	data := map[string]interface{}{
		"Module":            g.getModuleName(),
		"ProjectName":       g.options.ProjectName,
		"AuthConfig":        g.options.AuthConfig,
		"UserModel":         userModel,
		"RoleModel":         roleModel,
		"PermissionModel":   permissionModel,
		"DatabaseProvider":  g.options.DatabaseProvider,
		"OAuth2Providers":   g.options.AuthConfig.OAuth2Providers,
		"SupabaseConfig":    g.options.AuthConfig.Supabase,
		"Imports":           g.getRequiredImports(),
		"DatabaseTags":      g.getDatabaseTags(),
		"Endpoints":         g.getAuthEndpoints(),
	}

	// Add primary key type
	userModel.PrimaryKeyType = primaryKeyType
	data["UserModel"] = userModel

	return data
}

// getModuleName extracts module name from output path
func (g *AuthGenerator) getModuleName() string {
	// Try to read go.mod file
	goModPath := filepath.Join(g.options.OutputPath, "go.mod")
	if data, err := os.ReadFile(goModPath); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "module ") {
				return strings.TrimSpace(strings.TrimPrefix(line, "module"))
			}
		}
	}

	// Fallback to project name
	return fmt.Sprintf("github.com/%s/%s", g.options.ProjectName, g.options.ProjectName)
}

// getRequiredImports returns required imports for templates
func (g *AuthGenerator) getRequiredImports() []string {
	imports := g.options.AuthConfig.GetRequiredImports()
	
	// Add database-specific imports
	switch g.options.DatabaseProvider {
	case "postgres":
		imports = append(imports, "gorm.io/driver/postgres", "gorm.io/gorm")
	case "mysql":
		imports = append(imports, "gorm.io/driver/mysql", "gorm.io/gorm")
	case "sqlite":
		imports = append(imports, "gorm.io/driver/sqlite", "gorm.io/gorm")
	case "mongodb":
		imports = append(imports, "go.mongodb.org/mongo-driver/mongo")
	}

	return imports
}

// getDatabaseTags returns database tags for model fields
func (g *AuthGenerator) getDatabaseTags() map[string]string {
	tags := make(map[string]string)
	
	switch g.options.DatabaseProvider {
	case "postgres", "mysql", "sqlite":
		// GORM tags
		tags["Primary"] = fmt.Sprintf(`gorm:"primaryKey;column:%s" json:"%s"`, 
			g.options.UserModel.PrimaryKey, g.options.UserModel.PrimaryKey)
		tags["Email"] = fmt.Sprintf(`gorm:"unique;not null;column:%s" json:"%s"`, 
			g.options.UserModel.EmailField, g.options.UserModel.EmailField)
		tags["Username"] = fmt.Sprintf(`gorm:"unique;column:%s" json:"%s"`, 
			g.options.UserModel.UsernameField, g.options.UserModel.UsernameField)
		tags["Password"] = fmt.Sprintf(`gorm:"not null;column:%s" json:"-"`, 
			g.options.UserModel.PasswordField)
	case "mongodb":
		// MongoDB tags
		tags["Primary"] = `bson:"_id,omitempty" json:"id,omitempty"`
		tags["Email"] = `bson:"email" json:"email"`
		tags["Username"] = `bson:"username" json:"username"`
		tags["Password"] = `bson:"password_hash" json:"-"`
	default:
		// Default JSON tags
		tags["Primary"] = `json:"id"`
		tags["Email"] = `json:"email"`
		tags["Username"] = `json:"username"`
		tags["Password"] = `json:"-"`
	}

	return tags
}

// getAuthEndpoints returns authentication endpoints
func (g *AuthGenerator) getAuthEndpoints() []models.AuthEndpoint {
	endpoints := models.GetDefaultAuthEndpoints()
	
	if len(g.options.AuthConfig.OAuth2Providers) > 0 {
		endpoints = append(endpoints, models.GetOAuth2Endpoints()...)
	}
	
	if g.options.AuthConfig.EnableRBAC {
		endpoints = append(endpoints, models.GetRoleEndpoints()...)
	}
	
	return endpoints
}

// executeTemplate executes template and writes to file
func (g *AuthGenerator) executeTemplate(templateStr, outputPath string, data interface{}) error {
	tmpl := template.New("auth").Funcs(g.templateFuncs)
	tmpl, err := tmpl.Parse(templateStr)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outputPath, err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// Template helper functions
func (g *AuthGenerator) containsMethod(methods []models.AuthMethod, method models.AuthMethod) bool {
	for _, m := range methods {
		if m == method {
			return true
		}
	}
	return false
}

func (g *AuthGenerator) toJSON(v interface{}) string {
	// Simple JSON serialization for templates
	switch val := v.(type) {
	case []string:
		result := "["
		for i, item := range val {
			if i > 0 {
				result += ", "
			}
			result += fmt.Sprintf("\"%s\"", item)
		}
		result += "]"
		return result
	default:
		return fmt.Sprintf("%v", v)
	}
}

// Repository generation methods
func (g *AuthGenerator) generateGORMRepository(repoPath string, data map[string]interface{}) error {
	// TODO: Generate GORM-based repository
	repositoryTemplate := `package repositories

import (
	"{{.Module}}/internal/models"
	"gorm.io/gorm"
)

type {{.UserModel.StructName}}Repository struct {
	db *gorm.DB
}

func New{{.UserModel.StructName}}Repository(db *gorm.DB) *{{.UserModel.StructName}}Repository {
	return &{{.UserModel.StructName}}Repository{db: db}
}

func (r *{{.UserModel.StructName}}Repository) Create(user *models.{{.UserModel.StructName}}) error {
	return r.db.Create(user).Error
}

func (r *{{.UserModel.StructName}}Repository) FindByID(id {{.UserModel.PrimaryKeyType}}) (*models.{{.UserModel.StructName}}, error) {
	var user models.{{.UserModel.StructName}}
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *{{.UserModel.StructName}}Repository) FindByEmail(email string) (*models.{{.UserModel.StructName}}, error) {
	var user models.{{.UserModel.StructName}}
	err := r.db.Where("{{.UserModel.EmailField}} = ?", email).First(&user).Error
	return &user, err
}

func (r *{{.UserModel.StructName}}Repository) Update(user *models.{{.UserModel.StructName}}) error {
	return r.db.Save(user).Error
}

func (r *{{.UserModel.StructName}}Repository) Delete(id {{.UserModel.PrimaryKeyType}}) error {
	return r.db.Delete(&models.{{.UserModel.StructName}}{}, id).Error
}
`

	return g.executeTemplate(
		repositoryTemplate,
		filepath.Join(repoPath, "user.go"),
		data,
	)
}

func (g *AuthGenerator) generateMongoRepository(repoPath string, data map[string]interface{}) error {
	// TODO: Generate MongoDB repository
	return fmt.Errorf("MongoDB repository generation not yet implemented")
}

func (g *AuthGenerator) generateSupabaseRepository(repoPath string, data map[string]interface{}) error {
	// TODO: Generate Supabase repository
	return fmt.Errorf("Supabase repository generation not yet implemented")
}

// Migration generation methods
func (g *AuthGenerator) generateUserMigration(migrationsPath string) error {
	// TODO: Generate user table migration
	return nil
}

func (g *AuthGenerator) generateRBACMigrations(migrationsPath string) error {
	// TODO: Generate RBAC table migrations
	return nil
}

func (g *AuthGenerator) generateRouteTemplate() string {
	return `package routes

import (
	"net/http"

	"{{.Module}}/internal/handlers"
	"{{.Module}}/internal/middleware"
)

func SetupAuthRoutes(
	mux *http.ServeMux,
	authHandler *handlers.AuthHandler,
	jwtMiddleware *middleware.JWTMiddleware,
	{{if len .OAuth2Providers}}oauth2Handler *handlers.OAuth2Handler,{{end}}
) {
	// Setup auth routes
	authHandler.SetupAuthRoutes(mux, jwtMiddleware)
	
	{{if len .OAuth2Providers}}
	// Setup OAuth2 routes
	oauth2Handler.SetupOAuth2Routes(mux, jwtMiddleware)
	{{end}}
}
`
}

// DefaultAuthGeneratorOptions returns default options for auth generation
func DefaultAuthGeneratorOptions(projectName, outputPath string) models.AuthGeneratorOptions {
	return models.AuthGeneratorOptions{
		ProjectName:      projectName,
		OutputPath:       outputPath,
		AuthConfig:       models.DefaultAuthConfig(),
		UserModel:        models.DefaultUserModel(),
		DatabaseProvider: "postgres",
		Endpoints:        models.GetDefaultAuthEndpoints(),
	}
}

// WithRBAC adds RBAC configuration to auth generator options
func WithRBAC(options models.AuthGeneratorOptions) models.AuthGeneratorOptions {
	options.AuthConfig.EnableRBAC = true
	options.AuthConfig.EnablePermissions = true
	
	roleModel := models.DefaultRoleModel()
	permissionModel := models.DefaultPermissionModel()
	userRoleModel := models.DefaultUserRoleModel()
	rolePermModel := models.DefaultRolePermissionModel()
	
	options.RoleModel = &roleModel
	options.PermissionModel = &permissionModel
	options.UserRoleModel = &userRoleModel
	options.RolePermModel = &rolePermModel
	
	return options
}

// WithOAuth2 adds OAuth2 configuration to auth generator options
func WithOAuth2(options models.AuthGeneratorOptions, providers []models.OAuth2Provider) models.AuthGeneratorOptions {
	options.AuthConfig.OAuth2Providers = providers
	return options
}

// WithSupabase configures Supabase authentication
func WithSupabase(options models.AuthGeneratorOptions, supabaseConfig models.SupabaseAuthConfig) models.AuthGeneratorOptions {
	options.AuthConfig.Provider = models.AuthProviderSupabase
	options.AuthConfig.Supabase = supabaseConfig
	options.DatabaseProvider = "supabase"
	return options
}