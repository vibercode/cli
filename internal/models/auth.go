package models

import (
	"fmt"
	"time"
)

// AuthProvider represents different authentication providers
type AuthProvider string

const (
	AuthProviderJWT      AuthProvider = "jwt"
	AuthProviderOAuth2   AuthProvider = "oauth2"
	AuthProviderSupabase AuthProvider = "supabase"
	AuthProviderAPIKey   AuthProvider = "api_key"
	AuthProviderSession  AuthProvider = "session"
)

// AuthMethod represents authentication methods within a provider
type AuthMethod string

const (
	AuthMethodEmail    AuthMethod = "email"
	AuthMethodUsername AuthMethod = "username"
	AuthMethodPhone    AuthMethod = "phone"
	AuthMethodGoogle   AuthMethod = "google"
	AuthMethodGitHub   AuthMethod = "github"
	AuthMethodFacebook AuthMethod = "facebook"
	AuthMethodTwitter  AuthMethod = "twitter"
)

// AuthConfig represents authentication configuration for generation
type AuthConfig struct {
	Provider          AuthProvider   `json:"provider" yaml:"provider"`
	Methods           []AuthMethod   `json:"methods" yaml:"methods"`
	JWTSecret         string         `json:"jwt_secret" yaml:"jwt_secret"`
	TokenExpiry       time.Duration  `json:"token_expiry" yaml:"token_expiry"`
	RefreshExpiry     time.Duration  `json:"refresh_expiry" yaml:"refresh_expiry"`
	EnableRegistration bool          `json:"enable_registration" yaml:"enable_registration"`
	EnableEmailVerify  bool          `json:"enable_email_verify" yaml:"enable_email_verify"`
	EnablePasswordReset bool         `json:"enable_password_reset" yaml:"enable_password_reset"`
	EnableTwoFactor    bool          `json:"enable_two_factor" yaml:"enable_two_factor"`
	PasswordMinLength  int           `json:"password_min_length" yaml:"password_min_length"`
	EnableRBAC         bool          `json:"enable_rbac" yaml:"enable_rbac"`
	EnablePermissions  bool          `json:"enable_permissions" yaml:"enable_permissions"`
	OAuth2Providers    []OAuth2Provider `json:"oauth2_providers" yaml:"oauth2_providers"`
	Supabase          SupabaseAuthConfig `json:"supabase" yaml:"supabase"`
}

// OAuth2Provider represents OAuth2 provider configuration
type OAuth2Provider struct {
	Name         string            `json:"name" yaml:"name"`
	ClientID     string            `json:"client_id" yaml:"client_id"`
	ClientSecret string            `json:"client_secret" yaml:"client_secret"`
	RedirectURL  string            `json:"redirect_url" yaml:"redirect_url"`
	Scopes       []string          `json:"scopes" yaml:"scopes"`
	AuthURL      string            `json:"auth_url" yaml:"auth_url"`
	TokenURL     string            `json:"token_url" yaml:"token_url"`
	UserInfoURL  string            `json:"user_info_url" yaml:"user_info_url"`
	ExtraConfig  map[string]string `json:"extra_config" yaml:"extra_config"`
}

// SupabaseAuthConfig represents Supabase authentication configuration
type SupabaseAuthConfig struct {
	ProjectURL     string   `json:"project_url" yaml:"project_url"`
	APIKey         string   `json:"api_key" yaml:"api_key"`
	ServiceKey     string   `json:"service_key" yaml:"service_key"`
	JWTSecret      string   `json:"jwt_secret" yaml:"jwt_secret"`
	EnableAuth     bool     `json:"enable_auth" yaml:"enable_auth"`
	EnableSocial   bool     `json:"enable_social" yaml:"enable_social"`
	SocialProviders []string `json:"social_providers" yaml:"social_providers"`
}

// UserModel represents user model configuration for generation
type UserModel struct {
	TableName          string        `json:"table_name" yaml:"table_name"`
	StructName         string        `json:"struct_name" yaml:"struct_name"`
	PrimaryKey         string        `json:"primary_key" yaml:"primary_key"`
	PrimaryKeyType     string        `json:"primary_key_type" yaml:"primary_key_type"`
	UsernameField      string        `json:"username_field" yaml:"username_field"`
	EmailField         string        `json:"email_field" yaml:"email_field"`
	PasswordField      string        `json:"password_field" yaml:"password_field"`
	PhoneField         string        `json:"phone_field" yaml:"phone_field"`
	FirstNameField     string        `json:"first_name_field" yaml:"first_name_field"`
	LastNameField      string        `json:"last_name_field" yaml:"last_name_field"`
	AvatarField        string        `json:"avatar_field" yaml:"avatar_field"`
	EmailVerifiedField string        `json:"email_verified_field" yaml:"email_verified_field"`
	PhoneVerifiedField string        `json:"phone_verified_field" yaml:"phone_verified_field"`
	TwoFactorField     string        `json:"two_factor_field" yaml:"two_factor_field"`
	StatusField        string        `json:"status_field" yaml:"status_field"`
	CreatedAtField     string        `json:"created_at_field" yaml:"created_at_field"`
	UpdatedAtField     string        `json:"updated_at_field" yaml:"updated_at_field"`
	AdditionalFields   []Field       `json:"additional_fields" yaml:"additional_fields"`
}

// RoleModel represents role model configuration for RBAC
type RoleModel struct {
	TableName        string `json:"table_name" yaml:"table_name"`
	StructName       string `json:"struct_name" yaml:"struct_name"`
	PrimaryKey       string `json:"primary_key" yaml:"primary_key"`
	PrimaryKeyType   string `json:"primary_key_type" yaml:"primary_key_type"`
	NameField        string `json:"name_field" yaml:"name_field"`
	DisplayField     string `json:"display_field" yaml:"display_field"`
	DescriptionField string `json:"description_field" yaml:"description_field"`
	ColorField       string `json:"color_field" yaml:"color_field"`
	IsDefaultField   string `json:"is_default_field" yaml:"is_default_field"`
	IsSystemField    string `json:"is_system_field" yaml:"is_system_field"`
	CreatedAtField   string `json:"created_at_field" yaml:"created_at_field"`
	UpdatedAtField   string `json:"updated_at_field" yaml:"updated_at_field"`
}

// PermissionModel represents permission model configuration
type PermissionModel struct {
	TableName        string `json:"table_name" yaml:"table_name"`
	StructName       string `json:"struct_name" yaml:"struct_name"`
	PrimaryKey       string `json:"primary_key" yaml:"primary_key"`
	PrimaryKeyType   string `json:"primary_key_type" yaml:"primary_key_type"`
	NameField        string `json:"name_field" yaml:"name_field"`
	DisplayField     string `json:"display_field" yaml:"display_field"`
	DescriptionField string `json:"description_field" yaml:"description_field"`
	ResourceField    string `json:"resource_field" yaml:"resource_field"`
	ActionField      string `json:"action_field" yaml:"action_field"`
	CreatedAtField   string `json:"created_at_field" yaml:"created_at_field"`
	UpdatedAtField   string `json:"updated_at_field" yaml:"updated_at_field"`
}

// UserRoleModel represents user-role relationship for many-to-many
type UserRoleModel struct {
	TableName    string `json:"table_name" yaml:"table_name"`
	UserIDField  string `json:"user_id_field" yaml:"user_id_field"`
	RoleIDField  string `json:"role_id_field" yaml:"role_id_field"`
	AssignedAtField string `json:"assigned_at_field" yaml:"assigned_at_field"`
	AssignedByField string `json:"assigned_by_field" yaml:"assigned_by_field"`
}

// RolePermissionModel represents role-permission relationship
type RolePermissionModel struct {
	TableName      string `json:"table_name" yaml:"table_name"`
	RoleIDField    string `json:"role_id_field" yaml:"role_id_field"`
	PermissionIDField string `json:"permission_id_field" yaml:"permission_id_field"`
	GrantedAtField string `json:"granted_at_field" yaml:"granted_at_field"`
	GrantedByField string `json:"granted_by_field" yaml:"granted_by_field"`
}

// AuthEndpoint represents authentication endpoints to generate
type AuthEndpoint struct {
	Name        string            `json:"name" yaml:"name"`
	Method      string            `json:"method" yaml:"method"`
	Path        string            `json:"path" yaml:"path"`
	Handler     string            `json:"handler" yaml:"handler"`
	Middleware  []string          `json:"middleware" yaml:"middleware"`
	Parameters  []EndpointParam   `json:"parameters" yaml:"parameters"`
	Responses   []EndpointResponse `json:"responses" yaml:"responses"`
	RequiresAuth bool             `json:"requires_auth" yaml:"requires_auth"`
	RequiredRole string           `json:"required_role" yaml:"required_role"`
	RequiredPerms []string        `json:"required_permissions" yaml:"required_permissions"`
}

// EndpointParam represents endpoint parameter
type EndpointParam struct {
	Name        string `json:"name" yaml:"name"`
	Type        string `json:"type" yaml:"type"`
	Location    string `json:"location" yaml:"location"` // body, query, path, header
	Required    bool   `json:"required" yaml:"required"`
	Description string `json:"description" yaml:"description"`
	Validation  string `json:"validation" yaml:"validation"`
}

// EndpointResponse represents endpoint response
type EndpointResponse struct {
	StatusCode  int                    `json:"status_code" yaml:"status_code"`
	Description string                 `json:"description" yaml:"description"`
	Schema      map[string]interface{} `json:"schema" yaml:"schema"`
	Example     interface{}            `json:"example" yaml:"example"`
}

// AuthGeneratorOptions represents options for authentication generation
type AuthGeneratorOptions struct {
	ProjectName     string                 `json:"project_name" yaml:"project_name"`
	OutputPath      string                 `json:"output_path" yaml:"output_path"`
	AuthConfig      AuthConfig             `json:"auth_config" yaml:"auth_config"`
	UserModel       UserModel              `json:"user_model" yaml:"user_model"`
	RoleModel       *RoleModel             `json:"role_model" yaml:"role_model"`
	PermissionModel *PermissionModel       `json:"permission_model" yaml:"permission_model"`
	UserRoleModel   *UserRoleModel         `json:"user_role_model" yaml:"user_role_model"`
	RolePermModel   *RolePermissionModel   `json:"role_permission_model" yaml:"role_permission_model"`
	Endpoints       []AuthEndpoint         `json:"endpoints" yaml:"endpoints"`
	DatabaseProvider string                `json:"database_provider" yaml:"database_provider"`
	CustomFields    map[string]interface{} `json:"custom_fields" yaml:"custom_fields"`
}

// DefaultAuthConfig returns default authentication configuration
func DefaultAuthConfig() AuthConfig {
	return AuthConfig{
		Provider:          AuthProviderJWT,
		Methods:           []AuthMethod{AuthMethodEmail},
		TokenExpiry:       24 * time.Hour,
		RefreshExpiry:     7 * 24 * time.Hour,
		EnableRegistration: true,
		EnableEmailVerify: true,
		EnablePasswordReset: true,
		EnableTwoFactor:   false,
		PasswordMinLength: 8,
		EnableRBAC:        true,
		EnablePermissions: true,
		OAuth2Providers:   []OAuth2Provider{},
	}
}

// DefaultUserModel returns default user model configuration
func DefaultUserModel() UserModel {
	return UserModel{
		TableName:          "users",
		StructName:         "User",
		PrimaryKey:         "id",
		PrimaryKeyType:     "uint",
		UsernameField:      "username",
		EmailField:         "email",
		PasswordField:      "password_hash",
		PhoneField:         "phone",
		FirstNameField:     "first_name",
		LastNameField:      "last_name",
		AvatarField:        "avatar_url",
		EmailVerifiedField: "email_verified_at",
		PhoneVerifiedField: "phone_verified_at",
		TwoFactorField:     "two_factor_secret",
		StatusField:        "status",
		CreatedAtField:     "created_at",
		UpdatedAtField:     "updated_at",
		AdditionalFields:   []Field{},
	}
}

// DefaultRoleModel returns default role model configuration
func DefaultRoleModel() RoleModel {
	return RoleModel{
		TableName:        "roles",
		StructName:       "Role",
		PrimaryKey:       "id",
		PrimaryKeyType:   "uint",
		NameField:        "name",
		DisplayField:     "display_name",
		DescriptionField: "description",
		ColorField:       "color",
		IsDefaultField:   "is_default",
		IsSystemField:    "is_system",
		CreatedAtField:   "created_at",
		UpdatedAtField:   "updated_at",
	}
}

// DefaultPermissionModel returns default permission model configuration
func DefaultPermissionModel() PermissionModel {
	return PermissionModel{
		TableName:        "permissions",
		StructName:       "Permission",
		PrimaryKey:       "id",
		PrimaryKeyType:   "uint",
		NameField:        "name",
		DisplayField:     "display_name",
		DescriptionField: "description",
		ResourceField:    "resource",
		ActionField:      "action",
		CreatedAtField:   "created_at",
		UpdatedAtField:   "updated_at",
	}
}

// DefaultUserRoleModel returns default user-role relationship model
func DefaultUserRoleModel() UserRoleModel {
	return UserRoleModel{
		TableName:       "user_roles",
		UserIDField:     "user_id",
		RoleIDField:     "role_id",
		AssignedAtField: "assigned_at",
		AssignedByField: "assigned_by",
	}
}

// DefaultRolePermissionModel returns default role-permission relationship model
func DefaultRolePermissionModel() RolePermissionModel {
	return RolePermissionModel{
		TableName:         "role_permissions",
		RoleIDField:       "role_id",
		PermissionIDField: "permission_id",
		GrantedAtField:    "granted_at",
		GrantedByField:    "granted_by",
	}
}

// GetDefaultAuthEndpoints returns standard authentication endpoints
func GetDefaultAuthEndpoints() []AuthEndpoint {
	return []AuthEndpoint{
		{
			Name:    "Register",
			Method:  "POST",
			Path:    "/auth/register",
			Handler: "Register",
			Parameters: []EndpointParam{
				{Name: "email", Type: "string", Location: "body", Required: true, Validation: "email"},
				{Name: "password", Type: "string", Location: "body", Required: true, Validation: "min=8"},
				{Name: "first_name", Type: "string", Location: "body", Required: false},
				{Name: "last_name", Type: "string", Location: "body", Required: false},
			},
			Responses: []EndpointResponse{
				{StatusCode: 201, Description: "User registered successfully"},
				{StatusCode: 400, Description: "Invalid input"},
				{StatusCode: 409, Description: "User already exists"},
			},
		},
		{
			Name:    "Login",
			Method:  "POST",
			Path:    "/auth/login",
			Handler: "Login",
			Parameters: []EndpointParam{
				{Name: "email", Type: "string", Location: "body", Required: true, Validation: "email"},
				{Name: "password", Type: "string", Location: "body", Required: true},
			},
			Responses: []EndpointResponse{
				{StatusCode: 200, Description: "Login successful"},
				{StatusCode: 401, Description: "Invalid credentials"},
			},
		},
		{
			Name:         "RefreshToken",
			Method:       "POST",
			Path:         "/auth/refresh",
			Handler:      "RefreshToken",
			RequiresAuth: true,
			Parameters: []EndpointParam{
				{Name: "refresh_token", Type: "string", Location: "body", Required: true},
			},
			Responses: []EndpointResponse{
				{StatusCode: 200, Description: "Token refreshed successfully"},
				{StatusCode: 401, Description: "Invalid refresh token"},
			},
		},
		{
			Name:         "Logout",
			Method:       "POST",
			Path:         "/auth/logout",
			Handler:      "Logout",
			RequiresAuth: true,
			Responses: []EndpointResponse{
				{StatusCode: 200, Description: "Logged out successfully"},
			},
		},
		{
			Name:         "Profile",
			Method:       "GET",
			Path:         "/auth/profile",
			Handler:      "GetProfile",
			RequiresAuth: true,
			Responses: []EndpointResponse{
				{StatusCode: 200, Description: "User profile retrieved"},
				{StatusCode: 401, Description: "Unauthorized"},
			},
		},
		{
			Name:         "UpdateProfile",
			Method:       "PUT",
			Path:         "/auth/profile",
			Handler:      "UpdateProfile",
			RequiresAuth: true,
			Parameters: []EndpointParam{
				{Name: "first_name", Type: "string", Location: "body", Required: false},
				{Name: "last_name", Type: "string", Location: "body", Required: false},
				{Name: "phone", Type: "string", Location: "body", Required: false},
			},
			Responses: []EndpointResponse{
				{StatusCode: 200, Description: "Profile updated successfully"},
				{StatusCode: 401, Description: "Unauthorized"},
			},
		},
		{
			Name:    "ForgotPassword",
			Method:  "POST",
			Path:    "/auth/forgot-password",
			Handler: "ForgotPassword",
			Parameters: []EndpointParam{
				{Name: "email", Type: "string", Location: "body", Required: true, Validation: "email"},
			},
			Responses: []EndpointResponse{
				{StatusCode: 200, Description: "Password reset email sent"},
				{StatusCode: 404, Description: "User not found"},
			},
		},
		{
			Name:    "ResetPassword",
			Method:  "POST",
			Path:    "/auth/reset-password",
			Handler: "ResetPassword",
			Parameters: []EndpointParam{
				{Name: "token", Type: "string", Location: "body", Required: true},
				{Name: "password", Type: "string", Location: "body", Required: true, Validation: "min=8"},
			},
			Responses: []EndpointResponse{
				{StatusCode: 200, Description: "Password reset successfully"},
				{StatusCode: 400, Description: "Invalid or expired token"},
			},
		},
	}
}

// GetOAuth2Endpoints returns OAuth2-specific endpoints
func GetOAuth2Endpoints() []AuthEndpoint {
	return []AuthEndpoint{
		{
			Name:    "GoogleLogin",
			Method:  "GET",
			Path:    "/auth/google",
			Handler: "GoogleLogin",
			Responses: []EndpointResponse{
				{StatusCode: 302, Description: "Redirect to Google OAuth"},
			},
		},
		{
			Name:    "GoogleCallback",
			Method:  "GET",
			Path:    "/auth/google/callback",
			Handler: "GoogleCallback",
			Parameters: []EndpointParam{
				{Name: "code", Type: "string", Location: "query", Required: true},
				{Name: "state", Type: "string", Location: "query", Required: true},
			},
			Responses: []EndpointResponse{
				{StatusCode: 200, Description: "OAuth login successful"},
				{StatusCode: 400, Description: "OAuth error"},
			},
		},
		{
			Name:    "GitHubLogin",
			Method:  "GET",
			Path:    "/auth/github",
			Handler: "GitHubLogin",
			Responses: []EndpointResponse{
				{StatusCode: 302, Description: "Redirect to GitHub OAuth"},
			},
		},
		{
			Name:    "GitHubCallback",
			Method:  "GET",
			Path:    "/auth/github/callback",
			Handler: "GitHubCallback",
			Parameters: []EndpointParam{
				{Name: "code", Type: "string", Location: "query", Required: true},
				{Name: "state", Type: "string", Location: "query", Required: true},
			},
			Responses: []EndpointResponse{
				{StatusCode: 200, Description: "OAuth login successful"},
				{StatusCode: 400, Description: "OAuth error"},
			},
		},
	}
}

// GetRoleEndpoints returns role management endpoints
func GetRoleEndpoints() []AuthEndpoint {
	return []AuthEndpoint{
		{
			Name:         "ListRoles",
			Method:       "GET",
			Path:         "/auth/roles",
			Handler:      "ListRoles",
			RequiresAuth: true,
			RequiredRole: "admin",
			Responses: []EndpointResponse{
				{StatusCode: 200, Description: "List of roles"},
				{StatusCode: 401, Description: "Unauthorized"},
				{StatusCode: 403, Description: "Insufficient permissions"},
			},
		},
		{
			Name:         "CreateRole",
			Method:       "POST",
			Path:         "/auth/roles",
			Handler:      "CreateRole",
			RequiresAuth: true,
			RequiredRole: "admin",
			Parameters: []EndpointParam{
				{Name: "name", Type: "string", Location: "body", Required: true},
				{Name: "display_name", Type: "string", Location: "body", Required: true},
				{Name: "description", Type: "string", Location: "body", Required: false},
			},
			Responses: []EndpointResponse{
				{StatusCode: 201, Description: "Role created successfully"},
				{StatusCode: 400, Description: "Invalid input"},
				{StatusCode: 409, Description: "Role already exists"},
			},
		},
		{
			Name:         "AssignRole",
			Method:       "POST",
			Path:         "/auth/users/{user_id}/roles",
			Handler:      "AssignRole",
			RequiresAuth: true,
			RequiredRole: "admin",
			Parameters: []EndpointParam{
				{Name: "user_id", Type: "string", Location: "path", Required: true},
				{Name: "role_id", Type: "string", Location: "body", Required: true},
			},
			Responses: []EndpointResponse{
				{StatusCode: 200, Description: "Role assigned successfully"},
				{StatusCode: 404, Description: "User or role not found"},
			},
		},
	}
}

// ValidateAuthConfig validates authentication configuration
func ValidateAuthConfig(config AuthConfig) error {
	if config.Provider == "" {
		return fmt.Errorf("auth provider is required")
	}
	
	if len(config.Methods) == 0 {
		return fmt.Errorf("at least one auth method is required")
	}
	
	if config.Provider == AuthProviderJWT && config.JWTSecret == "" {
		return fmt.Errorf("JWT secret is required for JWT provider")
	}
	
	if config.PasswordMinLength < 6 {
		return fmt.Errorf("password minimum length must be at least 6")
	}
	
	return nil
}

// GetRequiredImports returns Go imports needed for authentication
func (config AuthConfig) GetRequiredImports() []string {
	imports := []string{
		"time",
		"fmt",
		"errors",
		"context",
		"crypto/rand",
		"golang.org/x/crypto/bcrypt",
	}
	
	if config.Provider == AuthProviderJWT {
		imports = append(imports, "github.com/golang-jwt/jwt/v4")
	}
	
	if len(config.OAuth2Providers) > 0 {
		imports = append(imports, "golang.org/x/oauth2")
		imports = append(imports, "encoding/json")
		imports = append(imports, "net/http")
		imports = append(imports, "io/ioutil")
	}
	
	if config.EnableTwoFactor {
		imports = append(imports, "github.com/pquerna/otp/totp")
	}
	
	return imports
}