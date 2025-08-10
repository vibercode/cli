# Authentication System Generator

## Overview

The Vibercode CLI Authentication System Generator provides comprehensive authentication and authorization code generation with support for JWT, OAuth2, RBAC, and multiple database providers including Supabase.

## Features

- ✅ **JWT Authentication**: Complete JWT-based auth with access and refresh tokens
- ✅ **OAuth2 Integration**: Support for Google, GitHub, and custom OAuth2 providers
- ✅ **RBAC (Role-Based Access Control)**: User roles and permissions system
- ✅ **Multiple Auth Methods**: Email, username, phone number authentication
- ✅ **Supabase Integration**: Native Supabase Auth integration
- ✅ **Database Agnostic**: Support for PostgreSQL, MySQL, SQLite, MongoDB
- ✅ **Security Features**: Password hashing, CSRF protection, rate limiting
- ✅ **Email Verification**: User email verification workflow
- ✅ **Password Reset**: Secure password reset functionality
- ✅ **Two-Factor Auth**: 2FA/TOTP support structure

## Quick Start

### Basic JWT Authentication

```go
package main

import (
    "github.com/vibercode/cli/internal/generator"
    "github.com/vibercode/cli/internal/models"
)

func main() {
    // Create default auth options
    opts := generator.DefaultAuthGeneratorOptions("my-api", "./output")
    
    // Generate authentication system
    authGen := generator.NewAuthGenerator(opts)
    if err := authGen.GenerateAuthSystem(); err != nil {
        panic(err)
    }
}
```

### With RBAC (Role-Based Access Control)

```go
// Enable RBAC with roles and permissions
opts := generator.DefaultAuthGeneratorOptions("my-api", "./output")
opts = generator.WithRBAC(opts)

authGen := generator.NewAuthGenerator(opts)
err := authGen.GenerateAuthSystem()
```

### With OAuth2 Providers

```go
// Configure OAuth2 providers
googleProvider := models.OAuth2Provider{
    Name:         "google",
    ClientID:     "your-google-client-id",
    ClientSecret: "your-google-client-secret",
    RedirectURL:  "http://localhost:8080/auth/google/callback",
    Scopes:       []string{"openid", "email", "profile"},
    AuthURL:      "https://accounts.google.com/o/oauth2/auth",
    TokenURL:     "https://oauth2.googleapis.com/token",
    UserInfoURL:  "https://www.googleapis.com/oauth2/v2/userinfo",
}

githubProvider := models.OAuth2Provider{
    Name:         "github",
    ClientID:     "your-github-client-id",
    ClientSecret: "your-github-client-secret",
    RedirectURL:  "http://localhost:8080/auth/github/callback",
    Scopes:       []string{"user:email"},
    AuthURL:      "https://github.com/login/oauth/authorize",
    TokenURL:     "https://github.com/login/oauth/access_token",
    UserInfoURL:  "https://api.github.com/user",
}

opts := generator.DefaultAuthGeneratorOptions("my-api", "./output")
opts = generator.WithOAuth2(opts, []models.OAuth2Provider{googleProvider, githubProvider})

authGen := generator.NewAuthGenerator(opts)
err := authGen.GenerateAuthSystem()
```

### With Supabase Authentication

```go
// Configure Supabase authentication
supabaseConfig := models.SupabaseAuthConfig{
    ProjectURL:      "https://your-project.supabase.co",
    APIKey:          "your-anon-key",
    ServiceKey:      "your-service-key",
    JWTSecret:       "your-jwt-secret",
    EnableAuth:      true,
    EnableSocial:    true,
    SocialProviders: []string{"google", "github"},
}

opts := generator.DefaultAuthGeneratorOptions("my-api", "./output")
opts = generator.WithSupabase(opts, supabaseConfig)

authGen := generator.NewAuthGenerator(opts)
err := authGen.GenerateAuthSystem()
```

## Generated Structure

The authentication generator creates a complete authentication system with the following structure:

```
output/
├── internal/
│   ├── middleware/
│   │   ├── jwt.go           # JWT authentication middleware
│   │   └── oauth2.go        # OAuth2 middleware (if enabled)
│   ├── handlers/
│   │   ├── auth.go          # Authentication handlers
│   │   └── oauth2.go        # OAuth2 handlers (if enabled)
│   ├── services/
│   │   ├── auth.go          # Authentication service
│   │   └── supabase_auth.go # Supabase integration (if enabled)
│   ├── models/
│   │   └── user.go          # User, Role, Permission models
│   ├── repositories/
│   │   └── user.go          # User repository (database-specific)
│   └── routes/
│       └── auth.go          # Route setup
├── migrations/
│   ├── 001_create_users.sql
│   ├── 002_create_roles.sql     # If RBAC enabled
│   └── 003_create_permissions.sql # If RBAC enabled
└── go.mod
```

## Authentication Endpoints

### Core Authentication

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/auth/register` | User registration | No |
| POST | `/auth/login` | User login | No |
| POST | `/auth/refresh` | Refresh access token | No |
| POST | `/auth/logout` | User logout | Yes |
| GET | `/auth/profile` | Get user profile | Yes |
| PUT | `/auth/profile` | Update user profile | Yes |
| POST | `/auth/forgot-password` | Request password reset | No |
| POST | `/auth/reset-password` | Reset password with token | No |

### OAuth2 Endpoints (if enabled)

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/auth/google` | Initiate Google OAuth | No |
| GET | `/auth/google/callback` | Google OAuth callback | No |
| GET | `/auth/github` | Initiate GitHub OAuth | No |
| GET | `/auth/github/callback` | GitHub OAuth callback | No |
| POST | `/auth/link` | Link OAuth account | Yes |
| DELETE | `/auth/unlink` | Unlink OAuth account | Yes |
| GET | `/auth/linked` | List linked accounts | Yes |

### RBAC Endpoints (if enabled)

| Method | Endpoint | Description | Auth Required | Role Required |
|--------|----------|-------------|---------------|---------------|
| GET | `/auth/roles` | List all roles | Yes | admin |
| POST | `/auth/roles` | Create new role | Yes | admin |
| POST | `/auth/users/{id}/roles` | Assign role to user | Yes | admin |
| DELETE | `/auth/users/{id}/roles/{role}` | Remove role from user | Yes | admin |

## Configuration Options

### AuthConfig

```go
type AuthConfig struct {
    Provider          AuthProvider   // jwt, oauth2, supabase
    Methods           []AuthMethod   // email, username, phone
    JWTSecret         string         // JWT signing secret
    TokenExpiry       time.Duration  // Access token expiry
    RefreshExpiry     time.Duration  // Refresh token expiry
    EnableRegistration bool          // Allow user registration
    EnableEmailVerify  bool          // Require email verification
    EnablePasswordReset bool         // Enable password reset
    EnableTwoFactor    bool          // Enable 2FA support
    PasswordMinLength  int           // Minimum password length
    EnableRBAC         bool          // Enable role-based access control
    EnablePermissions  bool          // Enable permission system
    OAuth2Providers    []OAuth2Provider // OAuth2 provider configurations
    Supabase          SupabaseAuthConfig // Supabase configuration
}
```

### UserModel Configuration

```go
type UserModel struct {
    TableName          string // Database table name
    StructName         string // Go struct name
    PrimaryKey         string // Primary key field name
    PrimaryKeyType     string // Primary key type (uint, string, etc.)
    EmailField         string // Email field name
    PasswordField      string // Password field name
    UsernameField      string // Username field name (optional)
    // ... other field configurations
}
```

## Database Support

### PostgreSQL/MySQL/SQLite (GORM)

Generated repositories use GORM for database operations:

```go
type UserRepository struct {
    db *gorm.DB
}

func (r *UserRepository) Create(user *models.User) error {
    return r.db.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
    var user models.User
    err := r.db.Where("email = ?", email).First(&user).Error
    return &user, err
}
```

### MongoDB

```go
// MongoDB repository implementation
type UserRepository struct {
    collection *mongo.Collection
}
```

### Supabase

Supabase integration provides direct API integration:

```go
type SupabaseAuthService struct {
    config     *config.SupabaseAuthConfig
    httpClient *http.Client
    baseURL    string
}

func (s *SupabaseAuthService) Register(email, password string, metadata map[string]interface{}) (*SupabaseSession, error)
func (s *SupabaseAuthService) Login(email, password string) (*SupabaseSession, error)
```

## Security Features

### Password Security

- **Bcrypt Hashing**: Secure password hashing using bcrypt with configurable cost
- **Minimum Length**: Configurable minimum password length requirements
- **Strength Validation**: Optional password strength validation

### JWT Security

- **HMAC Signing**: JWT tokens signed with HMAC-SHA256
- **Token Expiry**: Configurable access and refresh token expiration
- **Secure Claims**: User ID, email, roles included in JWT claims
- **Token Refresh**: Secure token refresh mechanism

### CSRF Protection

- **State Tokens**: OAuth2 flows protected with state tokens
- **CSRF Middleware**: Optional CSRF protection middleware generation

### Rate Limiting

- **Login Attempts**: Protection against brute force attacks
- **API Rate Limiting**: Configurable request rate limiting

## Middleware Usage

### JWT Authentication

```go
// Protect routes with JWT authentication
mux.Handle("GET /api/protected", jwtMiddleware.Authenticate(http.HandlerFunc(handler)))

// Require specific role
mux.Handle("GET /api/admin", jwtMiddleware.RequireRole("admin")(http.HandlerFunc(handler)))

// Require specific permission
mux.Handle("GET /api/users", jwtMiddleware.RequirePermission("users", "read")(http.HandlerFunc(handler)))
```

### Getting User from Context

```go
func protectedHandler(w http.ResponseWriter, r *http.Request) {
    userID, email, username, roles, err := middleware.GetUserFromContext(r)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    
    // Use user information
    fmt.Printf("User ID: %v, Email: %s, Roles: %v\n", userID, email, roles)
}
```

## Custom Field Types Integration

The authentication generator integrates with the enhanced field types system:

### Available Field Types for User Model

- **Basic Types**: `string`, `text`, `number`, `boolean`, `date`
- **Contact Types**: `email`, `phone`, `url`
- **Security Types**: `password` (with hashing)
- **Identity Types**: `slug` (for usernames)
- **Location Types**: `coordinates` (for user location)
- **Media Types**: `image` (for avatars), `file`
- **Selection Types**: `enum` (for user status, roles)
- **Financial Types**: `currency` (for user-related financial data)

### Example Custom User Model

```go
userModel := models.UserModel{
    TableName:      "users",
    StructName:     "User",
    EmailField:     "email",
    PasswordField:  "password_hash",
    AdditionalFields: []models.Field{
        {
            Name: "bio",
            Type: models.FieldTypeText,
            DisplayName: "User Biography",
            Required: false,
        },
        {
            Name: "avatar",
            Type: models.FieldTypeImage,
            DisplayName: "Profile Picture",
            Required: false,
        },
        {
            Name: "location",
            Type: models.FieldTypeCoordinates,
            DisplayName: "User Location",
            Required: false,
        },
        {
            Name: "subscription_tier",
            Type: models.FieldTypeEnum,
            DisplayName: "Subscription Tier",
            EnumValues: []string{"free", "premium", "enterprise"},
            Required: true,
            DefaultValue: "free",
        },
    },
}
```

## Testing

The authentication system includes comprehensive test coverage:

### Unit Tests

```bash
# Run authentication generator tests
go test ./internal/generator -run TestAuth -v

# Run authentication model tests  
go test ./internal/models -run TestAuth -v

# Run specific test suites
go test ./internal/generator -run TestAuthGeneratorWithOAuth2 -v
go test ./internal/generator -run TestAuthGeneratorWithRBAC -v
go test ./internal/generator -run TestAuthGeneratorWithSupabase -v
```

### Integration Tests

```bash
# Test complete authentication flow
go test ./internal/generator -run TestAuthGenerator -v

# Benchmark authentication generation
go test ./internal/generator -bench=BenchmarkAuthGenerator -v
```

### Manual Testing

```bash
# Generate test authentication system
./vibercode generate auth \
  --project-name test-auth \
  --output ./test-output \
  --provider jwt \
  --database postgres \
  --enable-rbac \
  --enable-oauth2

# Test generated system
cd test-output
go mod tidy
go test ./...
go run cmd/server/main.go
```

## Advanced Configuration

### Custom OAuth2 Provider

```go
customProvider := models.OAuth2Provider{
    Name:         "custom",
    ClientID:     os.Getenv("CUSTOM_CLIENT_ID"),
    ClientSecret: os.Getenv("CUSTOM_CLIENT_SECRET"),
    RedirectURL:  "http://localhost:8080/auth/custom/callback",
    Scopes:       []string{"read_user", "read_email"},
    AuthURL:      "https://auth.custom-provider.com/oauth/authorize",
    TokenURL:     "https://auth.custom-provider.com/oauth/token",
    UserInfoURL:  "https://api.custom-provider.com/user",
    ExtraConfig: map[string]string{
        "response_type": "code",
        "access_type":   "offline",
    },
}
```

### Custom User Fields

```go
// Add custom fields to user model
additionalFields := []models.Field{
    {
        Name:        "department",
        Type:        models.FieldTypeString,
        DisplayName: "Department",
        Required:    false,
    },
    {
        Name:        "employee_id",
        Type:        models.FieldTypeString,
        DisplayName: "Employee ID",
        Required:    true,
        Unique:      true,
    },
    {
        Name:        "hire_date",
        Type:        models.FieldTypeDate,
        DisplayName: "Hire Date",
        Required:    true,
    },
    {
        Name:        "salary",
        Type:        models.FieldTypeCurrency,
        DisplayName: "Salary",
        Required:    false,
        Validation: models.ValidationRule{
            Min: 0,
        },
    },
}

userModel := models.DefaultUserModel()
userModel.AdditionalFields = additionalFields
```

### Custom RBAC Configuration

```go
// Configure custom role and permission models
roleModel := models.RoleModel{
    TableName:        "user_roles",
    StructName:       "UserRole",
    NameField:        "role_name",
    DisplayField:     "display_name",
    DescriptionField: "description",
    IsDefaultField:   "is_default",
}

permissionModel := models.PermissionModel{
    TableName:     "permissions",
    StructName:    "Permission",
    NameField:     "permission_name",
    ResourceField: "resource",
    ActionField:   "action",
}

opts.RoleModel = &roleModel
opts.PermissionModel = &permissionModel
```

## Best Practices

### Security

1. **Use Strong JWT Secrets**: Generate cryptographically secure JWT signing keys
2. **Token Expiration**: Use short-lived access tokens (15-60 minutes) with longer refresh tokens
3. **Rate Limiting**: Implement rate limiting on authentication endpoints
4. **Input Validation**: Validate all user inputs on both client and server side
5. **HTTPS Only**: Always use HTTPS in production environments
6. **Secure Cookies**: Use secure, httpOnly cookies for sensitive data

### Database

1. **Index Critical Fields**: Index email, username, and other frequently queried fields
2. **Soft Deletes**: Consider soft deletes for user accounts for audit trails
3. **Data Encryption**: Encrypt sensitive data at rest
4. **Backup Strategy**: Implement regular database backups
5. **Connection Pooling**: Use appropriate database connection pool settings

### Performance

1. **Caching**: Cache user sessions and permissions for better performance
2. **Database Queries**: Optimize database queries and use appropriate indexes
3. **Token Validation**: Consider caching JWT validation results
4. **Load Balancing**: Design for horizontal scaling

### Development

1. **Environment Variables**: Use environment variables for all secrets and configuration
2. **Testing**: Write comprehensive tests for authentication flows
3. **Logging**: Implement proper security logging and monitoring
4. **Code Review**: Always review authentication-related code changes
5. **Documentation**: Keep authentication documentation up to date

## Migration from Other Systems

### From Firebase Auth

```go
// Configure similar to Firebase Auth
opts := generator.DefaultAuthGeneratorOptions("my-api", "./output")
opts.AuthConfig.EnableEmailVerify = true
opts.AuthConfig.EnablePasswordReset = true
opts.AuthConfig.OAuth2Providers = []models.OAuth2Provider{
    // Configure Google OAuth2 similar to Firebase
}
```

### From Auth0

```go
// Configure custom OAuth2 provider for Auth0
auth0Provider := models.OAuth2Provider{
    Name:        "auth0",
    AuthURL:     "https://your-domain.auth0.com/authorize",
    TokenURL:    "https://your-domain.auth0.com/oauth/token",
    UserInfoURL: "https://your-domain.auth0.com/userinfo",
    // ... other Auth0 specific configuration
}
```

## Troubleshooting

### Common Issues

1. **JWT Token Validation Errors**
   - Check JWT secret configuration
   - Verify token expiration settings
   - Ensure proper token format

2. **OAuth2 Callback Failures**
   - Verify redirect URLs match exactly
   - Check client ID and secret configuration
   - Validate OAuth2 provider endpoints

3. **Database Connection Issues**
   - Check database provider configuration
   - Verify connection string format
   - Ensure database server is accessible

4. **Permission Denied Errors**
   - Verify RBAC configuration
   - Check user role assignments
   - Validate permission definitions

### Debug Mode

Enable debug logging for troubleshooting:

```go
opts := generator.DefaultAuthGeneratorOptions("my-api", "./output")
opts.AuthConfig.Debug = true // Enable debug mode

// Generated code will include debug logging
```

## Contributing

To contribute to the authentication system generator:

1. **Add New Features**: Extend the generator for new authentication methods
2. **Database Support**: Add support for additional database providers
3. **Security Improvements**: Enhance security features and best practices
4. **Documentation**: Improve documentation and examples
5. **Testing**: Add more comprehensive test coverage

## API Reference

### Core Types

- `AuthConfig` - Authentication configuration
- `UserModel` - User model configuration
- `RoleModel` - Role model configuration (RBAC)
- `PermissionModel` - Permission model configuration (RBAC)
- `OAuth2Provider` - OAuth2 provider configuration
- `SupabaseAuthConfig` - Supabase integration configuration

### Key Functions

```go
// Create new auth generator
func NewAuthGenerator(options models.AuthGeneratorOptions) *AuthGenerator

// Generate complete auth system
func (g *AuthGenerator) GenerateAuthSystem() error

// Configuration helpers
func DefaultAuthGeneratorOptions(projectName, outputPath string) models.AuthGeneratorOptions
func WithRBAC(options models.AuthGeneratorOptions) models.AuthGeneratorOptions
func WithOAuth2(options models.AuthGeneratorOptions, providers []models.OAuth2Provider) models.AuthGeneratorOptions
func WithSupabase(options models.AuthGeneratorOptions, config models.SupabaseAuthConfig) models.AuthGeneratorOptions
```

This authentication system generator provides a complete, production-ready authentication solution that integrates seamlessly with the Vibercode CLI's enhanced field types and configuration management systems.