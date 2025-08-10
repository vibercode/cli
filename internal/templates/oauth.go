package templates

// OAuth2MiddlewareTemplate generates OAuth2 middleware
const OAuth2MiddlewareTemplate = `package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"crypto/rand"
	"encoding/base64"
	"time"

	"{{.Module}}/internal/models"
	"{{.Module}}/internal/services"
	"{{.Module}}/pkg/config"
	{{range .Imports}}
	"{{.}}"{{end}}
)

// OAuth2Middleware handles OAuth2 authentication
type OAuth2Middleware struct {
	providers   map[string]*oauth2.Config
	authService *services.AuthService
	config      *config.AuthConfig
}

// NewOAuth2Middleware creates a new OAuth2 middleware
func NewOAuth2Middleware(authService *services.AuthService, config *config.AuthConfig) *OAuth2Middleware {
	m := &OAuth2Middleware{
		providers:   make(map[string]*oauth2.Config),
		authService: authService,
		config:      config,
	}

	// Setup OAuth2 providers
	for _, provider := range config.OAuth2Providers {
		m.setupProvider(provider)
	}

	return m
}

// setupProvider configures an OAuth2 provider
func (m *OAuth2Middleware) setupProvider(provider models.OAuth2Provider) {
	config := &oauth2.Config{
		ClientID:     provider.ClientID,
		ClientSecret: provider.ClientSecret,
		RedirectURL:  provider.RedirectURL,
		Scopes:       provider.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  provider.AuthURL,
			TokenURL: provider.TokenURL,
		},
	}

	m.providers[provider.Name] = config
}

// GetAuthURL returns OAuth2 authorization URL
func (m *OAuth2Middleware) GetAuthURL(provider string) (string, error) {
	config, exists := m.providers[provider]
	if !exists {
		return "", fmt.Errorf("provider %s not configured", provider)
	}

	// Generate state token for CSRF protection
	state, err := m.generateStateToken()
	if err != nil {
		return "", err
	}

	// TODO: Store state token with expiration
	// In production, store this in Redis or database with TTL

	return config.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

// HandleCallback handles OAuth2 callback
func (m *OAuth2Middleware) HandleCallback(provider, code, state string) (*models.{{.UserModel.StructName}}, error) {
	config, exists := m.providers[provider]
	if !exists {
		return nil, fmt.Errorf("provider %s not configured", provider)
	}

	// TODO: Validate state token to prevent CSRF attacks
	if !m.validateStateToken(state) {
		return nil, fmt.Errorf("invalid state token")
	}

	// Exchange code for token
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info from provider
	userInfo, err := m.getUserInfo(provider, token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Find or create user
	user, err := m.findOrCreateUser(provider, userInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to find or create user: %w", err)
	}

	return user, nil
}

// getUserInfo fetches user information from OAuth2 provider
func (m *OAuth2Middleware) getUserInfo(provider, accessToken string) (map[string]interface{}, error) {
	var userInfoURL string

	// Get user info URL for provider
	for _, p := range m.config.OAuth2Providers {
		if p.Name == provider {
			userInfoURL = p.UserInfoURL
			break
		}
	}

	if userInfoURL == "" {
		return nil, fmt.Errorf("user info URL not configured for provider %s", provider)
	}

	// Make request to user info endpoint
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: status %d", resp.StatusCode)
	}

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}

// findOrCreateUser finds existing user or creates new one from OAuth2 info
func (m *OAuth2Middleware) findOrCreateUser(provider string, userInfo map[string]interface{}) (*models.{{.UserModel.StructName}}, error) {
	// Extract email from user info
	email, ok := userInfo["email"].(string)
	if !ok || email == "" {
		return nil, fmt.Errorf("email not found in user info")
	}

	// Try to find existing user by email
	user, err := m.authService.GetUserByEmail(email)
	if err == nil {
		// User exists, update OAuth2 info if needed
		return user, nil
	}

	// Create new user
	newUser := &models.{{.UserModel.StructName}}{
		{{.UserModel.EmailField | title}}: email,
		{{if .UserModel.FirstNameField}}{{.UserModel.FirstNameField | title}}: m.extractStringField(userInfo, "given_name", "first_name"),{{end}}
		{{if .UserModel.LastNameField}}{{.UserModel.LastNameField | title}}: m.extractStringField(userInfo, "family_name", "last_name"),{{end}}
		{{if .UserModel.AvatarField}}{{.UserModel.AvatarField | title}}: m.extractStringPointer(userInfo, "picture", "avatar_url"),{{end}}
		{{if .UserModel.UsernameField}}{{.UserModel.UsernameField | title}}: m.generateUsernameFromEmail(email),{{end}}
		{{if .UserModel.EmailVerifiedField}}{{.UserModel.EmailVerifiedField | title}}: m.getEmailVerifiedTime(userInfo),{{end}}
		{{if .UserModel.StatusField}}{{.UserModel.StatusField | title}}: "active",{{end}}
		{{if .UserModel.CreatedAtField}}{{.UserModel.CreatedAtField | title}}: time.Now(),{{end}}
		{{if .UserModel.UpdatedAtField}}{{.UserModel.UpdatedAtField | title}}: time.Now(),{{end}}
	}

	// Generate a random password (OAuth2 users don't use password auth)
	randomPassword, err := m.generateRandomPassword()
	if err != nil {
		return nil, err
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(randomPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	newUser.{{.UserModel.PasswordField | title}} = string(hashedPassword)

	// Create user
	if err := m.authService.CreateOAuth2User(newUser, provider); err != nil {
		return nil, err
	}

	return newUser, nil
}

// extractStringField extracts string field from user info with fallbacks
func (m *OAuth2Middleware) extractStringField(userInfo map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if value, ok := userInfo[key].(string); ok && value != "" {
			return value
		}
	}
	return ""
}

// extractStringPointer extracts string field as pointer
func (m *OAuth2Middleware) extractStringPointer(userInfo map[string]interface{}, keys ...string) *string {
	value := m.extractStringField(userInfo, keys...)
	if value == "" {
		return nil
	}
	return &value
}

// generateUsernameFromEmail generates username from email
func (m *OAuth2Middleware) generateUsernameFromEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 0 {
		return parts[0]
	}
	return email
}

// getEmailVerifiedTime gets email verification time from user info
func (m *OAuth2Middleware) getEmailVerifiedTime(userInfo map[string]interface{}) *time.Time {
	if verified, ok := userInfo["email_verified"].(bool); ok && verified {
		now := time.Now()
		return &now
	}
	return nil
}

// generateRandomPassword generates a random password for OAuth2 users
func (m *OAuth2Middleware) generateRandomPassword() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// generateStateToken generates a secure state token
func (m *OAuth2Middleware) generateStateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// validateStateToken validates the state token (placeholder)
func (m *OAuth2Middleware) validateStateToken(state string) bool {
	// TODO: Implement proper state token validation
	// In production, check against stored tokens with expiration
	return len(state) > 0
}
`

// OAuth2HandlersTemplate generates OAuth2 handlers
const OAuth2HandlersTemplate = `package handlers

import (
	"fmt"
	"net/http"
	"encoding/json"

	"{{.Module}}/internal/middleware"
	"{{.Module}}/internal/services"
	"{{.Module}}/pkg/config"
	{{range .Imports}}
	"{{.}}"{{end}}
)

// OAuth2Handler handles OAuth2 authentication requests
type OAuth2Handler struct {
	oauth2Middleware *middleware.OAuth2Middleware
	authService      *services.AuthService
	jwtMiddleware    *middleware.JWTMiddleware
	config           *config.AuthConfig
}

// NewOAuth2Handler creates a new OAuth2 handler
func NewOAuth2Handler(
	oauth2Middleware *middleware.OAuth2Middleware,
	authService *services.AuthService,
	jwtMiddleware *middleware.JWTMiddleware,
	config *config.AuthConfig,
) *OAuth2Handler {
	return &OAuth2Handler{
		oauth2Middleware: oauth2Middleware,
		authService:      authService,
		jwtMiddleware:    jwtMiddleware,
		config:           config,
	}
}

{{range .OAuth2Providers}}
// {{.Name | title}}Login initiates {{.Name}} OAuth2 flow
func (h *OAuth2Handler) {{.Name | title}}Login(w http.ResponseWriter, r *http.Request) {
	authURL, err := h.oauth2Middleware.GetAuthURL("{{.Name}}")
	if err != nil {
		http.Error(w, "Failed to get auth URL: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// {{.Name | title}}Callback handles {{.Name}} OAuth2 callback
func (h *OAuth2Handler) {{.Name | title}}Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" {
		http.Error(w, "Authorization code not found", http.StatusBadRequest)
		return
	}

	if state == "" {
		http.Error(w, "State parameter not found", http.StatusBadRequest)
		return
	}

	// Handle OAuth2 callback
	user, err := h.oauth2Middleware.HandleCallback("{{.Name}}", code, state)
	if err != nil {
		http.Error(w, "OAuth2 callback failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Generate JWT tokens
	roles, _ := h.authService.GetUserRoles(user.{{$.UserModel.PrimaryKey}})
	accessToken, err := h.jwtMiddleware.GenerateToken(
		user.{{$.UserModel.PrimaryKey}},
		user.{{$.UserModel.EmailField}},
		{{if $.UserModel.UsernameField}}user.{{$.UserModel.UsernameField}}{{else}}""{{end}},
		roles,
	)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := h.jwtMiddleware.GenerateRefreshToken(user.{{$.UserModel.PrimaryKey}})
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
{{end}}

// LinkAccount links OAuth2 account to existing user
func (h *OAuth2Handler) LinkAccount(w http.ResponseWriter, r *http.Request) {
	// Get current user from JWT
	userID, _, _, _, err := middleware.GetUserFromContext(r)
	if err != nil {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	// Get provider from URL path or query parameter
	provider := r.URL.Query().Get("provider")
	if provider == "" {
		http.Error(w, "Provider parameter required", http.StatusBadRequest)
		return
	}

	// TODO: Implement account linking logic
	// This would store the OAuth2 account association with existing user

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Account linking for %s initiated", provider),
	})
}

// UnlinkAccount removes OAuth2 account link
func (h *OAuth2Handler) UnlinkAccount(w http.ResponseWriter, r *http.Request) {
	// Get current user from JWT
	userID, _, _, _, err := middleware.GetUserFromContext(r)
	if err != nil {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	// Get provider from URL path or query parameter
	provider := r.URL.Query().Get("provider")
	if provider == "" {
		http.Error(w, "Provider parameter required", http.StatusBadRequest)
		return
	}

	// TODO: Implement account unlinking logic
	// This would remove the OAuth2 account association

	_ = userID // Use userID for unlinking

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Account unlinked from %s", provider),
	})
}

// ListLinkedAccounts returns linked OAuth2 accounts
func (h *OAuth2Handler) ListLinkedAccounts(w http.ResponseWriter, r *http.Request) {
	// Get current user from JWT
	userID, _, _, _, err := middleware.GetUserFromContext(r)
	if err != nil {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	// TODO: Get linked accounts from database
	linkedAccounts := []map[string]interface{}{
		// Example structure:
		// {
		//   "provider": "google",
		//   "linked_at": "2023-01-01T00:00:00Z",
		//   "provider_user_id": "123456789",
		// }
	}

	_ = userID // Use userID to query linked accounts

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"linked_accounts": linkedAccounts,
	})
}

// SetupOAuth2Routes sets up OAuth2 routes
func (h *OAuth2Handler) SetupOAuth2Routes(mux *http.ServeMux, jwtMiddleware *middleware.JWTMiddleware) {
	{{range .OAuth2Providers}}
	// {{.Name | title}} OAuth2 routes
	mux.HandleFunc("GET /auth/{{.Name}}", h.{{.Name | title}}Login)
	mux.HandleFunc("GET /auth/{{.Name}}/callback", h.{{.Name | title}}Callback)
	{{end}}

	// Account linking routes (protected)
	mux.Handle("POST /auth/link", jwtMiddleware.Authenticate(http.HandlerFunc(h.LinkAccount)))
	mux.Handle("DELETE /auth/unlink", jwtMiddleware.Authenticate(http.HandlerFunc(h.UnlinkAccount)))
	mux.Handle("GET /auth/linked", jwtMiddleware.Authenticate(http.HandlerFunc(h.ListLinkedAccounts)))
}
`

// SupabaseAuthTemplate generates Supabase authentication integration
const SupabaseAuthTemplate = `package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"{{.Module}}/internal/models"
	"{{.Module}}/pkg/config"
	{{range .Imports}}
	"{{.}}"{{end}}
)

// SupabaseAuthService handles Supabase authentication
type SupabaseAuthService struct {
	config     *config.SupabaseAuthConfig
	httpClient *http.Client
	baseURL    string
}

// NewSupabaseAuthService creates a new Supabase auth service
func NewSupabaseAuthService(config *config.SupabaseAuthConfig) *SupabaseAuthService {
	return &SupabaseAuthService{
		config:     config,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    config.ProjectURL + "/auth/v1",
	}
}

// SupabaseUser represents Supabase user response
type SupabaseUser struct {
	ID                 string                 ` + "`json:\"id\"`" + `
	Email              string                 ` + "`json:\"email\"`" + `
	EmailConfirmedAt   *time.Time             ` + "`json:\"email_confirmed_at\"`" + `
	Phone              string                 ` + "`json:\"phone\"`" + `
	PhoneConfirmedAt   *time.Time             ` + "`json:\"phone_confirmed_at\"`" + `
	ConfirmationSentAt *time.Time             ` + "`json:\"confirmation_sent_at\"`" + `
	RecoverySentAt     *time.Time             ` + "`json:\"recovery_sent_at\"`" + `
	LastSignInAt       *time.Time             ` + "`json:\"last_sign_in_at\"`" + `
	AppMetadata        map[string]interface{} ` + "`json:\"app_metadata\"`" + `
	UserMetadata       map[string]interface{} ` + "`json:\"user_metadata\"`" + `
	CreatedAt          time.Time              ` + "`json:\"created_at\"`" + `
	UpdatedAt          time.Time              ` + "`json:\"updated_at\"`" + `
}

// SupabaseSession represents Supabase session response
type SupabaseSession struct {
	AccessToken  string       ` + "`json:\"access_token\"`" + `
	TokenType    string       ` + "`json:\"token_type\"`" + `
	ExpiresIn    int64        ` + "`json:\"expires_in\"`" + `
	RefreshToken string       ` + "`json:\"refresh_token\"`" + `
	User         SupabaseUser ` + "`json:\"user\"`" + `
}

// RegisterRequest represents Supabase registration request
type SupabaseRegisterRequest struct {
	Email    string                 ` + "`json:\"email\"`" + `
	Password string                 ` + "`json:\"password\"`" + `
	Data     map[string]interface{} ` + "`json:\"data,omitempty\"`" + `
}

// LoginRequest represents Supabase login request
type SupabaseLoginRequest struct {
	Email    string ` + "`json:\"email\"`" + `
	Password string ` + "`json:\"password\"`" + `
}

// Register creates a new user in Supabase
func (s *SupabaseAuthService) Register(email, password string, metadata map[string]interface{}) (*SupabaseSession, error) {
	request := SupabaseRegisterRequest{
		Email:    email,
		Password: password,
		Data:     metadata,
	}

	var session SupabaseSession
	err := s.makeRequest("POST", "/signup", request, &session)
	if err != nil {
		return nil, fmt.Errorf("registration failed: %w", err)
	}

	return &session, nil
}

// Login authenticates user with Supabase
func (s *SupabaseAuthService) Login(email, password string) (*SupabaseSession, error) {
	request := SupabaseLoginRequest{
		Email:    email,
		Password: password,
	}

	var session SupabaseSession
	err := s.makeRequest("POST", "/token?grant_type=password", request, &session)
	if err != nil {
		return nil, fmt.Errorf("login failed: %w", err)
	}

	return &session, nil
}

// RefreshToken refreshes the access token
func (s *SupabaseAuthService) RefreshToken(refreshToken string) (*SupabaseSession, error) {
	request := map[string]string{
		"refresh_token": refreshToken,
	}

	var session SupabaseSession
	err := s.makeRequest("POST", "/token?grant_type=refresh_token", request, &session)
	if err != nil {
		return nil, fmt.Errorf("token refresh failed: %w", err)
	}

	return &session, nil
}

// GetUser retrieves user information by access token
func (s *SupabaseAuthService) GetUser(accessToken string) (*SupabaseUser, error) {
	var user SupabaseUser
	err := s.makeAuthenticatedRequest("GET", "/user", nil, &user, accessToken)
	if err != nil {
		return nil, fmt.Errorf("get user failed: %w", err)
	}

	return &user, nil
}

// UpdateUser updates user metadata
func (s *SupabaseAuthService) UpdateUser(accessToken string, metadata map[string]interface{}) (*SupabaseUser, error) {
	request := map[string]interface{}{
		"data": metadata,
	}

	var user SupabaseUser
	err := s.makeAuthenticatedRequest("PUT", "/user", request, &user, accessToken)
	if err != nil {
		return nil, fmt.Errorf("update user failed: %w", err)
	}

	return &user, nil
}

// Logout signs out the user
func (s *SupabaseAuthService) Logout(accessToken string) error {
	err := s.makeAuthenticatedRequest("POST", "/logout", nil, nil, accessToken)
	if err != nil {
		return fmt.Errorf("logout failed: %w", err)
	}

	return nil
}

// ResetPassword sends password reset email
func (s *SupabaseAuthService) ResetPassword(email string) error {
	request := map[string]string{
		"email": email,
	}

	err := s.makeRequest("POST", "/recover", request, nil)
	if err != nil {
		return fmt.Errorf("password reset failed: %w", err)
	}

	return nil
}

// VerifyOTP verifies OTP for email/phone verification
func (s *SupabaseAuthService) VerifyOTP(email, token, otpType string) (*SupabaseSession, error) {
	request := map[string]string{
		"email": email,
		"token": token,
		"type":  otpType, // "signup", "recovery", "email_change"
	}

	var session SupabaseSession
	err := s.makeRequest("POST", "/verify", request, &session)
	if err != nil {
		return nil, fmt.Errorf("OTP verification failed: %w", err)
	}

	return &session, nil
}

{{if .SupabaseConfig.EnableSocial}}
// GetOAuthURL returns OAuth URL for social provider
func (s *SupabaseAuthService) GetOAuthURL(provider, redirectURL string) (string, error) {
	return fmt.Sprintf("%s/authorize?provider=%s&redirect_to=%s", s.baseURL, provider, redirectURL), nil
}
{{end}}

// ConvertToLocalUser converts Supabase user to local user model
func (s *SupabaseAuthService) ConvertToLocalUser(supabaseUser *SupabaseUser) *models.{{.UserModel.StructName}} {
	user := &models.{{.UserModel.StructName}}{
		{{.UserModel.EmailField | title}}: supabaseUser.Email,
		{{if .UserModel.EmailVerifiedField}}{{.UserModel.EmailVerifiedField | title}}: supabaseUser.EmailConfirmedAt,{{end}}
		{{if .UserModel.PhoneField}}{{.UserModel.PhoneField | title}}: &supabaseUser.Phone,{{end}}
		{{if .UserModel.PhoneVerifiedField}}{{.UserModel.PhoneVerifiedField | title}}: supabaseUser.PhoneConfirmedAt,{{end}}
		{{if .UserModel.StatusField}}{{.UserModel.StatusField | title}}: "active",{{end}}
		{{if .UserModel.CreatedAtField}}{{.UserModel.CreatedAtField | title}}: supabaseUser.CreatedAt,{{end}}
		{{if .UserModel.UpdatedAtField}}{{.UserModel.UpdatedAtField | title}}: supabaseUser.UpdatedAt,{{end}}
	}

	// Extract metadata
	if firstName, ok := supabaseUser.UserMetadata["first_name"].(string); ok && firstName != "" {
		{{if .UserModel.FirstNameField}}user.{{.UserModel.FirstNameField | title}} = firstName{{end}}
	}

	if lastName, ok := supabaseUser.UserMetadata["last_name"].(string); ok && lastName != "" {
		{{if .UserModel.LastNameField}}user.{{.UserModel.LastNameField | title}} = lastName{{end}}
	}

	if avatar, ok := supabaseUser.UserMetadata["avatar_url"].(string); ok && avatar != "" {
		{{if .UserModel.AvatarField}}user.{{.UserModel.AvatarField | title}} = &avatar{{end}}
	}

	return user
}

// makeRequest makes HTTP request to Supabase
func (s *SupabaseAuthService) makeRequest(method, endpoint string, body interface{}, result interface{}) error {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequest(method, s.baseURL+endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", s.config.APIKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			if msg, ok := errorResp["msg"].(string); ok {
				return fmt.Errorf("Supabase error: %s", msg)
			}
		}
		return fmt.Errorf("Supabase error: status %d", resp.StatusCode)
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return err
		}
	}

	return nil
}

// makeAuthenticatedRequest makes authenticated HTTP request
func (s *SupabaseAuthService) makeAuthenticatedRequest(method, endpoint string, body interface{}, result interface{}, accessToken string) error {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequest(method, s.baseURL+endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", s.config.APIKey)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errorResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			if msg, ok := errorResp["msg"].(string); ok {
				return fmt.Errorf("Supabase error: %s", msg)
			}
		}
		return fmt.Errorf("Supabase error: status %d", resp.StatusCode)
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return err
		}
	}

	return nil
}
`