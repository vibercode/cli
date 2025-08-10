package templates

// JWTMiddlewareTemplate generates JWT authentication middleware
const JWTMiddlewareTemplate = `package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"{{.Module}}/internal/models"
	"{{.Module}}/pkg/config"
	{{range .Imports}}
	"{{.}}"{{end}}
)

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID    {{.UserModel.PrimaryKeyType}} ` + "`json:\"user_id\"`" + `
	Email     string ` + "`json:\"email\"`" + `
	Username  string ` + "`json:\"username\"`" + `
	Roles     []string ` + "`json:\"roles\"`" + `
	jwt.RegisteredClaims
}

// JWTMiddleware provides JWT authentication middleware
type JWTMiddleware struct {
	secretKey []byte
	config    *config.AuthConfig
}

// NewJWTMiddleware creates a new JWT middleware instance
func NewJWTMiddleware(secretKey string, authConfig *config.AuthConfig) *JWTMiddleware {
	return &JWTMiddleware{
		secretKey: []byte(secretKey),
		config:    authConfig,
	}
}

// Authenticate middleware verifies JWT tokens
func (m *JWTMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Check Bearer format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := tokenParts[1]

		// Parse and validate token
		claims := &JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return m.secretKey, nil
		})

		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Token is not valid", http.StatusUnauthorized)
			return
		}

		// Check token expiration
		if claims.ExpiresAt.Time.Before(time.Now()) {
			http.Error(w, "Token has expired", http.StatusUnauthorized)
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "email", claims.Email)
		ctx = context.WithValue(ctx, "username", claims.Username)
		ctx = context.WithValue(ctx, "roles", claims.Roles)

		// Continue to next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole middleware checks if user has required role
func (m *JWTMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roles, ok := r.Context().Value("roles").([]string)
			if !ok {
				http.Error(w, "No roles found in context", http.StatusForbidden)
				return
			}

			// Check if user has required role
			hasRole := false
			for _, userRole := range roles {
				if userRole == role || userRole == "admin" { // admin has all permissions
					hasRole = true
					break
				}
			}

			if !hasRole {
				http.Error(w, fmt.Sprintf("Role '%s' required", role), http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequirePermission middleware checks if user has required permission
func (m *JWTMiddleware) RequirePermission(resource, action string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := r.Context().Value("user_id").({{.UserModel.PrimaryKeyType}})
			if !ok {
				http.Error(w, "User ID not found in context", http.StatusForbidden)
				return
			}

			// TODO: Check user permissions from database
			// This should query user permissions through roles
			hasPermission := m.checkUserPermission(userID, resource, action)
			
			if !hasPermission {
				http.Error(w, fmt.Sprintf("Permission '%s:%s' required", resource, action), http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// checkUserPermission checks if user has specific permission
func (m *JWTMiddleware) checkUserPermission(userID {{.UserModel.PrimaryKeyType}}, resource, action string) bool {
	// TODO: Implement permission checking logic
	// This should query the database to check user permissions
	return true // Placeholder
}

// GenerateToken generates a new JWT token for user
func (m *JWTMiddleware) GenerateToken(userID {{.UserModel.PrimaryKeyType}}, email, username string, roles []string) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Email:    email,
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.config.TokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "{{.ProjectName}}",
			Subject:   fmt.Sprintf("%v", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
}

// GenerateRefreshToken generates a refresh token
func (m *JWTMiddleware) GenerateRefreshToken(userID {{.UserModel.PrimaryKeyType}}) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.config.RefreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "{{.ProjectName}}",
			Subject:   fmt.Sprintf("%v", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
}

// ValidateRefreshToken validates and extracts claims from refresh token
func (m *JWTMiddleware) ValidateRefreshToken(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, fmt.Errorf("token has expired")
	}

	return claims, nil
}

// GetUserFromContext extracts user information from request context
func GetUserFromContext(r *http.Request) (userID {{.UserModel.PrimaryKeyType}}, email, username string, roles []string, err error) {
	if userIDVal := r.Context().Value("user_id"); userIDVal != nil {
		if uid, ok := userIDVal.({{.UserModel.PrimaryKeyType}}); ok {
			userID = uid
		} else {
			err = fmt.Errorf("invalid user ID type")
			return
		}
	} else {
		err = fmt.Errorf("user ID not found in context")
		return
	}

	if emailVal := r.Context().Value("email"); emailVal != nil {
		email, _ = emailVal.(string)
	}

	if usernameVal := r.Context().Value("username"); usernameVal != nil {
		username, _ = usernameVal.(string)
	}

	if rolesVal := r.Context().Value("roles"); rolesVal != nil {
		roles, _ = rolesVal.([]string)
	}

	return
}
`

// AuthHandlersTemplate generates authentication handlers
const AuthHandlersTemplate = `package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"strings"

	"{{.Module}}/internal/models"
	"{{.Module}}/internal/services"
	"{{.Module}}/internal/middleware"
	"{{.Module}}/pkg/config"
	{{range .Imports}}
	"{{.}}"{{end}}
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	authService *services.AuthService
	jwtMiddleware *middleware.JWTMiddleware
	config      *config.AuthConfig
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService *services.AuthService, jwtMiddleware *middleware.JWTMiddleware, config *config.AuthConfig) *AuthHandler {
	return &AuthHandler{
		authService:   authService,
		jwtMiddleware: jwtMiddleware,
		config:        config,
	}
}

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Email     string ` + "`json:\"email\" validate:\"required,email\"`" + `
	Password  string ` + "`json:\"password\" validate:\"required,min={{.AuthConfig.PasswordMinLength}}\"`" + `
	{{if .UserModel.UsernameField}}Username  string ` + "`json:\"username\" validate:\"required,min=3\"`" + `{{end}}
	{{if .UserModel.FirstNameField}}FirstName string ` + "`json:\"first_name\"`" + `{{end}}
	{{if .UserModel.LastNameField}}LastName  string ` + "`json:\"last_name\"`" + `{{end}}
	{{if .UserModel.PhoneField}}Phone     string ` + "`json:\"phone\"`" + `{{end}}
}

// LoginRequest represents user login request
type LoginRequest struct {
	{{if .AuthConfig.Methods | contains "email"}}Email    string ` + "`json:\"email\" validate:\"required,email\"`" + `{{end}}
	{{if .AuthConfig.Methods | contains "username"}}Username string ` + "`json:\"username\" validate:\"required\"`" + `{{end}}
	Password string ` + "`json:\"password\" validate:\"required\"`" + `
}

// AuthResponse represents authentication response
type AuthResponse struct {
	User         *models.{{.UserModel.StructName}} ` + "`json:\"user\"`" + `
	AccessToken  string                ` + "`json:\"access_token\"`" + `
	RefreshToken string                ` + "`json:\"refresh_token\"`" + `
	ExpiresIn    int64                 ` + "`json:\"expires_in\"`" + `
}

// RefreshRequest represents token refresh request
type RefreshRequest struct {
	RefreshToken string ` + "`json:\"refresh_token\" validate:\"required\"`" + `
}

// ForgotPasswordRequest represents forgot password request
type ForgotPasswordRequest struct {
	Email string ` + "`json:\"email\" validate:\"required,email\"`" + `
}

// ResetPasswordRequest represents password reset request
type ResetPasswordRequest struct {
	Token    string ` + "`json:\"token\" validate:\"required\"`" + `
	Password string ` + "`json:\"password\" validate:\"required,min={{.AuthConfig.PasswordMinLength}}\"`" + `
}

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	{{if .UserModel.FirstNameField}}FirstName *string ` + "`json:\"first_name\"`" + `{{end}}
	{{if .UserModel.LastNameField}}LastName  *string ` + "`json:\"last_name\"`" + `{{end}}
	{{if .UserModel.PhoneField}}Phone     *string ` + "`json:\"phone\"`" + `{{end}}
	{{if .UserModel.AvatarField}}AvatarURL *string ` + "`json:\"avatar_url\"`" + `{{end}}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if err := h.validateRegisterRequest(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if user already exists
	if exists, err := h.authService.UserExists(req.Email{{if .UserModel.UsernameField}}, req.Username{{end}}); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	} else if exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Create user
	user, err := h.authService.CreateUser(&req)
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	{{if .AuthConfig.EnableEmailVerify}}
	// Send verification email
	if err := h.authService.SendVerificationEmail(user.{{.UserModel.EmailField}}); err != nil {
		// Log error but don't fail the registration
		fmt.Printf("Failed to send verification email: %v\n", err)
	}
	{{end}}

	// Generate tokens
	roles, _ := h.authService.GetUserRoles(user.{{.UserModel.PrimaryKey}})
	accessToken, err := h.jwtMiddleware.GenerateToken(
		user.{{.UserModel.PrimaryKey}}, 
		user.{{.UserModel.EmailField}}, 
		{{if .UserModel.UsernameField}}user.{{.UserModel.UsernameField}}{{else}}"{{end}}", 
		roles,
	)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := h.jwtMiddleware.GenerateRefreshToken(user.{{.UserModel.PrimaryKey}})
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(h.config.TokenExpiry.Seconds()),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login handles user authentication
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Authenticate user
	user, err := h.authService.AuthenticateUser({{if .AuthConfig.Methods | contains "email"}}req.Email{{else}}req.Username{{end}}, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	{{if .AuthConfig.EnableEmailVerify}}
	// Check if email is verified
	if user.{{.UserModel.EmailVerifiedField}} == nil {
		http.Error(w, "Email not verified", http.StatusForbidden)
		return
	}
	{{end}}

	// Generate tokens
	roles, _ := h.authService.GetUserRoles(user.{{.UserModel.PrimaryKey}})
	accessToken, err := h.jwtMiddleware.GenerateToken(
		user.{{.UserModel.PrimaryKey}}, 
		user.{{.UserModel.EmailField}}, 
		{{if .UserModel.UsernameField}}user.{{.UserModel.UsernameField}}{{else}}""{{end}}, 
		roles,
	)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := h.jwtMiddleware.GenerateRefreshToken(user.{{.UserModel.PrimaryKey}})
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(h.config.TokenExpiry.Seconds()),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate refresh token
	claims, err := h.jwtMiddleware.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Get user
	user, err := h.authService.GetUserByID(claims.UserID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Generate new tokens
	roles, _ := h.authService.GetUserRoles(user.{{.UserModel.PrimaryKey}})
	accessToken, err := h.jwtMiddleware.GenerateToken(
		user.{{.UserModel.PrimaryKey}}, 
		user.{{.UserModel.EmailField}}, 
		{{if .UserModel.UsernameField}}user.{{.UserModel.UsernameField}}{{else}}""{{end}}, 
		roles,
	)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	newRefreshToken, err := h.jwtMiddleware.GenerateRefreshToken(user.{{.UserModel.PrimaryKey}})
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(h.config.TokenExpiry.Seconds()),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Logout handles user logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement token blacklisting if needed
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

// GetProfile returns current user profile
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, _, _, _, err := middleware.GetUserFromContext(r)
	if err != nil {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateProfile updates current user profile
func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, _, _, _, err := middleware.GetUserFromContext(r)
	if err != nil {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.authService.UpdateUserProfile(userID, &req)
	if err != nil {
		http.Error(w, "Failed to update profile: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

{{if .AuthConfig.EnablePasswordReset}}
// ForgotPassword handles forgot password requests
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.authService.SendPasswordResetEmail(req.Email)
	if err != nil {
		// Don't reveal if email exists or not
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "If email exists, reset link sent"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password reset email sent"})
}

// ResetPassword handles password reset
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.authService.ResetPassword(req.Token, req.Password)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password reset successfully"})
}
{{end}}

// validateRegisterRequest validates registration request
func (h *AuthHandler) validateRegisterRequest(req *RegisterRequest) error {
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	
	if len(req.Password) < h.config.PasswordMinLength {
		return fmt.Errorf("password must be at least %d characters", h.config.PasswordMinLength)
	}
	
	{{if .UserModel.UsernameField}}
	if req.Username == "" {
		return fmt.Errorf("username is required")
	}
	
	if len(req.Username) < 3 {
		return fmt.Errorf("username must be at least 3 characters")
	}
	{{end}}
	
	// Add more validation as needed
	return nil
}

// SetupAuthRoutes sets up authentication routes
func (h *AuthHandler) SetupAuthRoutes(mux *http.ServeMux, jwtMiddleware *middleware.JWTMiddleware) {
	// Public routes
	mux.HandleFunc("POST /auth/register", h.Register)
	mux.HandleFunc("POST /auth/login", h.Login)
	mux.HandleFunc("POST /auth/refresh", h.RefreshToken)
	{{if .AuthConfig.EnablePasswordReset}}
	mux.HandleFunc("POST /auth/forgot-password", h.ForgotPassword)
	mux.HandleFunc("POST /auth/reset-password", h.ResetPassword)
	{{end}}

	// Protected routes
	mux.Handle("POST /auth/logout", jwtMiddleware.Authenticate(http.HandlerFunc(h.Logout)))
	mux.Handle("GET /auth/profile", jwtMiddleware.Authenticate(http.HandlerFunc(h.GetProfile)))
	mux.Handle("PUT /auth/profile", jwtMiddleware.Authenticate(http.HandlerFunc(h.UpdateProfile)))
}
`

// AuthServiceTemplate generates authentication service
const AuthServiceTemplate = `package services

import (
	"fmt"
	"time"
	"crypto/rand"
	"encoding/base64"

	"{{.Module}}/internal/models"
	"{{.Module}}/internal/repositories"
	"{{.Module}}/internal/handlers"
	{{range .Imports}}
	"{{.}}"{{end}}
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo     *repositories.{{.UserModel.StructName}}Repository
	{{if .AuthConfig.EnableRBAC}}roleRepo     *repositories.RoleRepository{{end}}
	{{if .AuthConfig.EnablePermissions}}permRepo     *repositories.PermissionRepository{{end}}
	emailService *EmailService // Optional email service
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo *repositories.{{.UserModel.StructName}}Repository{{if .AuthConfig.EnableRBAC}}, roleRepo *repositories.RoleRepository{{end}}{{if .AuthConfig.EnablePermissions}}, permRepo *repositories.PermissionRepository{{end}}) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		{{if .AuthConfig.EnableRBAC}}roleRepo: roleRepo,{{end}}
		{{if .AuthConfig.EnablePermissions}}permRepo: permRepo,{{end}}
	}
}

// CreateUser creates a new user account
func (s *AuthService) CreateUser(req *handlers.RegisterRequest) (*models.{{.UserModel.StructName}}, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.{{.UserModel.StructName}}{
		{{.UserModel.EmailField}}:    req.Email,
		{{.UserModel.PasswordField}}: string(hashedPassword),
		{{if .UserModel.UsernameField}}{{.UserModel.UsernameField}}: req.Username,{{end}}
		{{if .UserModel.FirstNameField}}{{.UserModel.FirstNameField}}: req.FirstName,{{end}}
		{{if .UserModel.LastNameField}}{{.UserModel.LastNameField}}: req.LastName,{{end}}
		{{if .UserModel.PhoneField}}{{.UserModel.PhoneField}}: req.Phone,{{end}}
		{{if .UserModel.StatusField}}{{.UserModel.StatusField}}: "active",{{end}}
		{{if .UserModel.CreatedAtField}}{{.UserModel.CreatedAtField}}: time.Now(),{{end}}
		{{if .UserModel.UpdatedAtField}}{{.UserModel.UpdatedAtField}}: time.Now(),{{end}}
	}

	// Create user
	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	{{if .AuthConfig.EnableRBAC}}
	// Assign default role if configured
	if err := s.assignDefaultRole(user.{{.UserModel.PrimaryKey}}); err != nil {
		// Log error but don't fail user creation
		fmt.Printf("Warning: failed to assign default role: %v\n", err)
	}
	{{end}}

	return user, nil
}

// AuthenticateUser authenticates user with email/username and password
func (s *AuthService) AuthenticateUser(identifier, password string) (*models.{{.UserModel.StructName}}, error) {
	{{if .AuthConfig.Methods | contains "email"}}
	// Try to find user by email
	user, err := s.userRepo.FindBy{{.UserModel.EmailField | title}}(identifier)
	{{else}}
	// Try to find user by username
	user, err := s.userRepo.FindBy{{.UserModel.UsernameField | title}}(identifier)
	{{end}}
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.{{.UserModel.PasswordField}}), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	{{if .UserModel.StatusField}}
	// Check user status
	if user.{{.UserModel.StatusField}} != "active" {
		return nil, fmt.Errorf("user account is not active")
	}
	{{end}}

	return user, nil
}

// GetUserByID retrieves user by ID
func (s *AuthService) GetUserByID(id {{.UserModel.PrimaryKeyType}}) (*models.{{.UserModel.StructName}}, error) {
	return s.userRepo.FindByID(id)
}

// UserExists checks if user already exists by email or username
func (s *AuthService) UserExists(email{{if .UserModel.UsernameField}}, username{{end}} string) (bool, error) {
	// Check by email
	if user, _ := s.userRepo.FindBy{{.UserModel.EmailField | title}}(email); user != nil {
		return true, nil
	}

	{{if .UserModel.UsernameField}}
	// Check by username
	if user, _ := s.userRepo.FindBy{{.UserModel.UsernameField | title}}(username); user != nil {
		return true, nil
	}
	{{end}}

	return false, nil
}

{{if .AuthConfig.EnableRBAC}}
// GetUserRoles returns user roles
func (s *AuthService) GetUserRoles(userID {{.UserModel.PrimaryKeyType}}) ([]string, error) {
	roles, err := s.roleRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.{{.RoleModel.NameField | title}}
	}

	return roleNames, nil
}

// assignDefaultRole assigns default role to new user
func (s *AuthService) assignDefaultRole(userID {{.UserModel.PrimaryKeyType}}) error {
	defaultRole, err := s.roleRepo.FindDefault()
	if err != nil {
		return err // No default role configured
	}

	return s.roleRepo.AssignRoleToUser(userID, defaultRole.{{.RoleModel.PrimaryKey}})
}
{{else}}
// GetUserRoles returns empty roles (RBAC disabled)
func (s *AuthService) GetUserRoles(userID {{.UserModel.PrimaryKeyType}}) ([]string, error) {
	return []string{}, nil
}
{{end}}

// UpdateUserProfile updates user profile information
func (s *AuthService) UpdateUserProfile(userID {{.UserModel.PrimaryKeyType}}, req *handlers.UpdateProfileRequest) (*models.{{.UserModel.StructName}}, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// Update fields
	{{if .UserModel.FirstNameField}}
	if req.FirstName != nil {
		user.{{.UserModel.FirstNameField}} = *req.FirstName
	}
	{{end}}
	{{if .UserModel.LastNameField}}
	if req.LastName != nil {
		user.{{.UserModel.LastNameField}} = *req.LastName
	}
	{{end}}
	{{if .UserModel.PhoneField}}
	if req.Phone != nil {
		user.{{.UserModel.PhoneField}} = *req.Phone
	}
	{{end}}
	{{if .UserModel.AvatarField}}
	if req.AvatarURL != nil {
		user.{{.UserModel.AvatarField}} = *req.AvatarURL
	}
	{{end}}

	{{if .UserModel.UpdatedAtField}}
	user.{{.UserModel.UpdatedAtField}} = time.Now()
	{{end}}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

{{if .AuthConfig.EnableEmailVerify}}
// SendVerificationEmail sends email verification
func (s *AuthService) SendVerificationEmail(email string) error {
	// Generate verification token
	token, err := s.generateSecureToken()
	if err != nil {
		return err
	}

	// TODO: Store token in database with expiration
	// TODO: Send email with verification link

	return nil
}
{{end}}

{{if .AuthConfig.EnablePasswordReset}}
// SendPasswordResetEmail sends password reset email
func (s *AuthService) SendPasswordResetEmail(email string) error {
	user, err := s.userRepo.FindBy{{.UserModel.EmailField | title}}(email)
	if err != nil {
		return err // User not found
	}

	// Generate reset token
	token, err := s.generateSecureToken()
	if err != nil {
		return err
	}

	// TODO: Store token in database with expiration
	// TODO: Send email with reset link

	_ = user // Use user for sending email
	_ = token // Use token in reset link

	return nil
}

// ResetPassword resets user password with token
func (s *AuthService) ResetPassword(token, newPassword string) error {
	// TODO: Validate token from database
	// TODO: Find user by token
	// TODO: Check token expiration

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// TODO: Update user password
	// TODO: Invalidate reset token

	_ = hashedPassword // Use hashed password for update

	return nil
}
{{end}}

// generateSecureToken generates a secure random token
func (s *AuthService) generateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
`

// UserModelTemplate generates user model
const UserModelTemplate = `package models

import (
	"time"
	{{range .Imports}}
	"{{.}}"{{end}}
)

// {{.UserModel.StructName}} represents a user in the system
type {{.UserModel.StructName}} struct {
	{{.UserModel.PrimaryKey | title}} {{.UserModel.PrimaryKeyType}} ` + "`{{.DatabaseTags.Primary}}`" + `
	{{if .UserModel.UsernameField}}{{.UserModel.UsernameField | title}} string ` + "`{{.DatabaseTags.Username}}`" + `{{end}}
	{{.UserModel.EmailField | title}}    string    ` + "`{{.DatabaseTags.Email}}`" + `
	{{.UserModel.PasswordField | title}} string    ` + "`{{.DatabaseTags.Password}}`" + `
	{{if .UserModel.PhoneField}}{{.UserModel.PhoneField | title}}     *string   ` + "`{{.DatabaseTags.Phone}}`" + `{{end}}
	{{if .UserModel.FirstNameField}}{{.UserModel.FirstNameField | title}} string    ` + "`{{.DatabaseTags.FirstName}}`" + `{{end}}
	{{if .UserModel.LastNameField}}{{.UserModel.LastNameField | title}}  string    ` + "`{{.DatabaseTags.LastName}}`" + `{{end}}
	{{if .UserModel.AvatarField}}{{.UserModel.AvatarField | title}}   *string   ` + "`{{.DatabaseTags.Avatar}}`" + `{{end}}
	{{if .UserModel.EmailVerifiedField}}{{.UserModel.EmailVerifiedField | title}} *time.Time ` + "`{{.DatabaseTags.EmailVerified}}`" + `{{end}}
	{{if .UserModel.PhoneVerifiedField}}{{.UserModel.PhoneVerifiedField | title}} *time.Time ` + "`{{.DatabaseTags.PhoneVerified}}`" + `{{end}}
	{{if .UserModel.TwoFactorField}}{{.UserModel.TwoFactorField | title}} *string ` + "`{{.DatabaseTags.TwoFactor}}`" + `{{end}}
	{{if .UserModel.StatusField}}{{.UserModel.StatusField | title}}    string    ` + "`{{.DatabaseTags.Status}}`" + `{{end}}
	{{if .UserModel.CreatedAtField}}{{.UserModel.CreatedAtField | title}} time.Time ` + "`{{.DatabaseTags.CreatedAt}}`" + `{{end}}
	{{if .UserModel.UpdatedAtField}}{{.UserModel.UpdatedAtField | title}} time.Time ` + "`{{.DatabaseTags.UpdatedAt}}`" + `{{end}}
	
	{{range .UserModel.AdditionalFields}}
	{{.Name | title}} {{.GoType}} ` + "`{{.DatabaseTag}}`" + `{{end}}
	
	{{if .AuthConfig.EnableRBAC}}
	// Relationships
	Roles []{{.RoleModel.StructName}} ` + "`{{.DatabaseTags.UserRoles}}`" + `
	{{end}}
}

{{if .AuthConfig.EnableRBAC}}
// {{.RoleModel.StructName}} represents a role in the system
type {{.RoleModel.StructName}} struct {
	{{.RoleModel.PrimaryKey | title}}       {{.RoleModel.PrimaryKeyType}} ` + "`{{.DatabaseTags.RolePrimary}}`" + `
	{{.RoleModel.NameField | title}}        string     ` + "`{{.DatabaseTags.RoleName}}`" + `
	{{.RoleModel.DisplayField | title}}     string     ` + "`{{.DatabaseTags.RoleDisplay}}`" + `
	{{if .RoleModel.DescriptionField}}{{.RoleModel.DescriptionField | title}} *string    ` + "`{{.DatabaseTags.RoleDescription}}`" + `{{end}}
	{{if .RoleModel.ColorField}}{{.RoleModel.ColorField | title}}       *string    ` + "`{{.DatabaseTags.RoleColor}}`" + `{{end}}
	{{.RoleModel.IsDefaultField | title}}   bool       ` + "`{{.DatabaseTags.RoleIsDefault}}`" + `
	{{.RoleModel.IsSystemField | title}}    bool       ` + "`{{.DatabaseTags.RoleIsSystem}}`" + `
	{{.RoleModel.CreatedAtField | title}}   time.Time  ` + "`{{.DatabaseTags.RoleCreatedAt}}`" + `
	{{.RoleModel.UpdatedAtField | title}}   time.Time  ` + "`{{.DatabaseTags.RoleUpdatedAt}}`" + `
	
	{{if .AuthConfig.EnablePermissions}}
	// Relationships  
	Permissions []{{.PermissionModel.StructName}} ` + "`{{.DatabaseTags.RolePermissions}}`" + `
	{{end}}
	Users       []{{.UserModel.StructName}}      ` + "`{{.DatabaseTags.RoleUsers}}`" + `
}
{{end}}

{{if .AuthConfig.EnablePermissions}}
// {{.PermissionModel.StructName}} represents a permission in the system
type {{.PermissionModel.StructName}} struct {
	{{.PermissionModel.PrimaryKey | title}}         {{.PermissionModel.PrimaryKeyType}} ` + "`{{.DatabaseTags.PermissionPrimary}}`" + `
	{{.PermissionModel.NameField | title}}          string      ` + "`{{.DatabaseTags.PermissionName}}`" + `
	{{.PermissionModel.DisplayField | title}}       string      ` + "`{{.DatabaseTags.PermissionDisplay}}`" + `
	{{if .PermissionModel.DescriptionField}}{{.PermissionModel.DescriptionField | title}} *string     ` + "`{{.DatabaseTags.PermissionDescription}}`" + `{{end}}
	{{.PermissionModel.ResourceField | title}}      string      ` + "`{{.DatabaseTags.PermissionResource}}`" + `
	{{.PermissionModel.ActionField | title}}        string      ` + "`{{.DatabaseTags.PermissionAction}}`" + `
	{{.PermissionModel.CreatedAtField | title}}     time.Time   ` + "`{{.DatabaseTags.PermissionCreatedAt}}`" + `
	{{.PermissionModel.UpdatedAtField | title}}     time.Time   ` + "`{{.DatabaseTags.PermissionUpdatedAt}}`" + `
	
	// Relationships
	Roles []{{.RoleModel.StructName}} ` + "`{{.DatabaseTags.PermissionRoles}}`" + `
}
{{end}}

// TableName returns the table name for {{.UserModel.StructName}}
func ({{.UserModel.StructName}}) TableName() string {
	return "{{.UserModel.TableName}}"
}

{{if .AuthConfig.EnableRBAC}}
// TableName returns the table name for {{.RoleModel.StructName}}
func ({{.RoleModel.StructName}}) TableName() string {
	return "{{.RoleModel.TableName}}"
}
{{end}}

{{if .AuthConfig.EnablePermissions}}
// TableName returns the table name for {{.PermissionModel.StructName}}
func ({{.PermissionModel.StructName}}) TableName() string {
	return "{{.PermissionModel.TableName}}"
}
{{end}}

// BeforeCreate is called before creating a user
func (u *{{.UserModel.StructName}}) BeforeCreate() error {
	{{if .UserModel.CreatedAtField}}
	if u.{{.UserModel.CreatedAtField | title}}.IsZero() {
		u.{{.UserModel.CreatedAtField | title}} = time.Now()
	}
	{{end}}
	{{if .UserModel.UpdatedAtField}}
	if u.{{.UserModel.UpdatedAtField | title}}.IsZero() {
		u.{{.UserModel.UpdatedAtField | title}} = time.Now()
	}
	{{end}}
	return nil
}

// BeforeUpdate is called before updating a user
func (u *{{.UserModel.StructName}}) BeforeUpdate() error {
	{{if .UserModel.UpdatedAtField}}
	u.{{.UserModel.UpdatedAtField | title}} = time.Now()
	{{end}}
	return nil
}

// GetFullName returns user's full name
func (u *{{.UserModel.StructName}}) GetFullName() string {
	{{if and .UserModel.FirstNameField .UserModel.LastNameField}}
	return u.{{.UserModel.FirstNameField | title}} + " " + u.{{.UserModel.LastNameField | title}}
	{{else if .UserModel.FirstNameField}}
	return u.{{.UserModel.FirstNameField | title}}
	{{else if .UserModel.UsernameField}}
	return u.{{.UserModel.UsernameField | title}}
	{{else}}
	return u.{{.UserModel.EmailField | title}}
	{{end}}
}

{{if .AuthConfig.EnableRBAC}}
// HasRole checks if user has a specific role
func (u *{{.UserModel.StructName}}) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role.{{.RoleModel.NameField | title}} == roleName {
			return true
		}
	}
	return false
}

{{if .AuthConfig.EnablePermissions}}
// HasPermission checks if user has a specific permission
func (u *{{.UserModel.StructName}}) HasPermission(resource, action string) bool {
	for _, role := range u.Roles {
		for _, permission := range role.Permissions {
			if permission.{{.PermissionModel.ResourceField | title}} == resource && 
			   permission.{{.PermissionModel.ActionField | title}} == action {
				return true
			}
		}
	}
	return false
}
{{end}}
{{end}}

// IsActive returns true if user is active
func (u *{{.UserModel.StructName}}) IsActive() bool {
	{{if .UserModel.StatusField}}
	return u.{{.UserModel.StatusField | title}} == "active"
	{{else}}
	return true
	{{end}}
}

{{if .UserModel.EmailVerifiedField}}
// IsEmailVerified returns true if email is verified
func (u *{{.UserModel.StructName}}) IsEmailVerified() bool {
	return u.{{.UserModel.EmailVerifiedField | title}} != nil
}
{{end}}

{{if .UserModel.PhoneVerifiedField}}
// IsPhoneVerified returns true if phone is verified
func (u *{{.UserModel.StructName}}) IsPhoneVerified() bool {
	return u.{{.UserModel.PhoneVerifiedField | title}} != nil
}
{{end}}

{{if .UserModel.TwoFactorField}}
// IsTwoFactorEnabled returns true if 2FA is enabled
func (u *{{.UserModel.StructName}}) IsTwoFactorEnabled() bool {
	return u.{{.UserModel.TwoFactorField | title}} != nil && *u.{{.UserModel.TwoFactorField | title}} != ""
}
{{end}}
`