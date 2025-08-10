package templates

import (
	"fmt"
	"strings"

	"github.com/vibercode/cli/internal/models"
)

// GetMiddlewareTemplate returns the middleware template based on type
func GetMiddlewareTemplate(config models.MiddlewareConfig) string {
	switch config.Type {
	case models.AuthMiddleware:
		return getAuthMiddlewareTemplate(config)
	case models.LoggingMiddleware:
		return getLoggingMiddlewareTemplate(config)
	case models.CORSMiddleware:
		return getCORSMiddlewareTemplate(config)
	case models.RateLimitMiddleware:
		return getRateLimitMiddlewareTemplate(config)
	case models.CustomMiddleware:
		return getCustomMiddlewareTemplate(config)
	default:
		return getBasicMiddlewareTemplate(config)
	}
}

// getAuthMiddlewareTemplate generates authentication middleware template
func getAuthMiddlewareTemplate(config models.MiddlewareConfig) string {
	imports := generateImports(config.GetImports())
	
	switch config.Options.AuthStrategy {
	case models.JWTAuth:
		return fmt.Sprintf(`package %s

%s

// %s handles JWT authentication
type %s struct {
	secretKey []byte
	issuer    string
}

// New%s creates a new authentication middleware
func New%s(secretKey, issuer string) *%s {
	return &%s{
		secretKey: []byte(secretKey),
		issuer:    issuer,
	}
}

// JWT validates JWT tokens
func (m *%s) JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return m.secretKey, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Validate issuer
			if iss, ok := claims["iss"].(string); ok {
				if iss != m.issuer {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token issuer"})
					c.Abort()
					return
				}
			}

			// Set user context
			c.Set("user_id", claims["sub"])
			c.Set("user_email", claims["email"])
			c.Set("user_role", claims["role"])
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
	}
}

// RequireRole checks if user has required role
func (m *%s) RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No user role found"})
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok || role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
`, config.GetPackageName(), imports, config.Description, config.GetStructName(),
			config.GetStructName(), config.GetStructName(), config.GetStructName(), config.GetStructName(),
			config.GetStructName(), config.GetStructName())

	case models.APIKeyAuth:
		return fmt.Sprintf(`package %s

%s

// %s handles API key authentication
type %s struct {
	validKeys map[string]bool
	header    string
}

// New%s creates a new API key authentication middleware
func New%s(validKeys []string, header string) *%s {
	keyMap := make(map[string]bool)
	for _, key := range validKeys {
		keyMap[key] = true
	}
	
	if header == "" {
		header = "X-API-Key"
	}
	
	return &%s{
		validKeys: keyMap,
		header:    header,
	}
}

// APIKey validates API keys
func (m *%s) APIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader(m.header)
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API key required"})
			c.Abort()
			return
		}

		// Use constant-time comparison to prevent timing attacks
		var isValid bool
		for validKey := range m.validKeys {
			if subtle.ConstantTimeCompare([]byte(apiKey), []byte(validKey)) == 1 {
				isValid = true
				break
			}
		}

		if !isValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		c.Set("api_key", apiKey)
		c.Next()
	}
}
`, config.GetPackageName(), imports, config.Description, config.GetStructName(),
			config.GetStructName(), config.GetStructName(), config.GetStructName(), config.GetStructName(),
			config.GetStructName())

	default:
		return getBasicAuthTemplate(config, imports)
	}
}

// getLoggingMiddlewareTemplate generates logging middleware template
func getLoggingMiddlewareTemplate(config models.MiddlewareConfig) string {
	imports := generateImports(config.GetImports())
	
	return fmt.Sprintf(`package %s

%s

// %s handles request/response logging
type %s struct {
	logger       *logrus.Logger
	excludePaths map[string]bool
}

// New%s creates a new logging middleware
func New%s(logger *logrus.Logger, excludePaths []string) *%s {
	pathMap := make(map[string]bool)
	for _, path := range excludePaths {
		pathMap[path] = true
	}
	
	return &%s{
		logger:       logger,
		excludePaths: pathMap,
	}
}

// RequestLogger logs HTTP requests and responses
func (m *%s) RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip logging for excluded paths
		if m.excludePaths[c.Request.URL.Path] {
			c.Next()
			return
		}

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Log request
		m.logger.WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"path":       path,
			"query":      raw,
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}).Info("Request started")

		// Process request
		c.Next()

		// Log response
		latency := time.Since(start)
		statusCode := c.Writer.Status()

		logLevel := logrus.InfoLevel
		if statusCode >= 500 {
			logLevel = logrus.ErrorLevel
		} else if statusCode >= 400 {
			logLevel = logrus.WarnLevel
		}

		m.logger.WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"path":       path,
			"status":     statusCode,
			"latency":    latency,
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}).Log(logLevel, "Request completed")
	}
}

// ErrorLogger logs errors with context
func (m *%s) ErrorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Log any errors that occurred
		for _, err := range c.Errors {
			m.logger.WithFields(logrus.Fields{
				"method": c.Request.Method,
				"path":   c.Request.URL.Path,
				"ip":     c.ClientIP(),
			}).Error(err.Error())
		}
	}
}
`, config.GetPackageName(), imports, config.Description, config.GetStructName(),
		config.GetStructName(), config.GetStructName(), config.GetStructName(), config.GetStructName(),
		config.GetStructName(), config.GetStructName())
}

// getCORSMiddlewareTemplate generates CORS middleware template
func getCORSMiddlewareTemplate(config models.MiddlewareConfig) string {
	imports := generateImports(config.GetImports())
	
	allowedOrigins := formatStringSlice(config.Options.AllowedOrigins)
	allowedMethods := formatStringSlice(config.Options.AllowedMethods)
	allowedHeaders := formatStringSlice(config.Options.AllowedHeaders)
	
	return fmt.Sprintf(`package %s

%s

// %s handles CORS configuration
type %s struct {
	allowedOrigins   []string
	allowedMethods   []string
	allowedHeaders   []string
	exposeHeaders    []string
	allowCredentials bool
	maxAge          int
}

// New%s creates a new CORS middleware
func New%s() *%s {
	return &%s{
		allowedOrigins:   %s,
		allowedMethods:   %s,
		allowedHeaders:   %s,
		exposeHeaders:    []string{},
		allowCredentials: %t,
		maxAge:          %d,
	}
}

// CORS handles Cross-Origin Resource Sharing
func (m *%s) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		if m.isOriginAllowed(origin) {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		// Set allowed methods
		c.Header("Access-Control-Allow-Methods", strings.Join(m.allowedMethods, ", "))

		// Set allowed headers
		c.Header("Access-Control-Allow-Headers", strings.Join(m.allowedHeaders, ", "))

		// Set exposed headers
		if len(m.exposeHeaders) > 0 {
			c.Header("Access-Control-Expose-Headers", strings.Join(m.exposeHeaders, ", "))
		}

		// Set credentials
		if m.allowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// Set max age
		if m.maxAge > 0 {
			c.Header("Access-Control-Max-Age", fmt.Sprintf("%%d", m.maxAge))
		}

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func (m *%s) isOriginAllowed(origin string) bool {
	for _, allowedOrigin := range m.allowedOrigins {
		if allowedOrigin == "*" || allowedOrigin == origin {
			return true
		}
	}
	return false
}
`, config.GetPackageName(), imports, config.Description, config.GetStructName(),
		config.GetStructName(), config.GetStructName(), config.GetStructName(), config.GetStructName(),
		allowedOrigins, allowedMethods, allowedHeaders, config.Options.AllowCredentials, config.Options.MaxAge,
		config.GetStructName(), config.GetStructName())
}

// getRateLimitMiddlewareTemplate generates rate limiting middleware template
func getRateLimitMiddlewareTemplate(config models.MiddlewareConfig) string {
	imports := generateImports(config.GetImports())
	
	if config.Options.UseRedis {
		return getRedisRateLimitTemplate(config, imports)
	}
	
	return fmt.Sprintf(`package %s

%s

// %s handles rate limiting
type %s struct {
	limiter *rate.Limiter
	mu      sync.RWMutex
	clients map[string]*rate.Limiter
}

// New%s creates a new rate limiting middleware
func New%s(rps int, burst int) *%s {
	return &%s{
		limiter: rate.NewLimiter(rate.Limit(rps), burst),
		clients: make(map[string]*rate.Limiter),
	}
}

// RateLimit applies rate limiting per IP
func (m *%s) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		m.mu.RLock()
		limiter, exists := m.clients[ip]
		m.mu.RUnlock()

		if !exists {
			m.mu.Lock()
			// Double-check after acquiring write lock
			if limiter, exists = m.clients[ip]; !exists {
				limiter = rate.NewLimiter(m.limiter.Limit(), m.limiter.Burst())
				m.clients[ip] = limiter
			}
			m.mu.Unlock()
		}

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"retry_after": "60s",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GlobalRateLimit applies global rate limiting
func (m *%s) GlobalRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !m.limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"retry_after": "60s",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
`, config.GetPackageName(), imports, config.Description, config.GetStructName(),
		config.GetStructName(), config.GetStructName(), config.GetStructName(), config.GetStructName(),
		config.GetStructName(), config.GetStructName())
}

// getCustomMiddlewareTemplate generates custom middleware template
func getCustomMiddlewareTemplate(config models.MiddlewareConfig) string {
	imports := generateImports([]string{
		"net/http",
		"github.com/gin-gonic/gin",
	})
	
	return fmt.Sprintf(`package %s

%s

// %s %s
type %s struct {
	// Add your configuration fields here
}

// New%s creates a new %s middleware
func New%s() *%s {
	return &%s{
		// Initialize your configuration here
	}
}

// Handle processes the request
func (m *%s) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add your middleware logic here
		
		// Example: Log request
		// log.Printf("Processing request: %%s %%s", c.Request.Method, c.Request.URL.Path)
		
		// Continue to next middleware/handler
		c.Next()
		
		// Example: Log response
		// log.Printf("Request completed with status: %%d", c.Writer.Status())
	}
}
`, config.GetPackageName(), imports, config.Name, config.Description, config.GetStructName(),
		config.GetStructName(), config.Name, config.GetStructName(), config.GetStructName(), config.GetStructName(),
		config.GetStructName())
}

// Helper templates

func getBasicAuthTemplate(config models.MiddlewareConfig, imports string) string {
	return fmt.Sprintf(`package %s

%s

// %s handles basic authentication
type %s struct{}

// New%s creates a new basic authentication middleware
func New%s() *%s {
	return &%s{}
}

// BasicAuth validates basic authentication
func (m *%s) BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implement basic authentication logic
		c.Next()
	}
}
`, config.GetPackageName(), imports, config.Description, config.GetStructName(),
		config.GetStructName(), config.GetStructName(), config.GetStructName(), config.GetStructName(),
		config.GetStructName())
}

func getRedisRateLimitTemplate(config models.MiddlewareConfig, imports string) string {
	return fmt.Sprintf(`package %s

%s

// %s handles distributed rate limiting with Redis
type %s struct {
	client *redis.Client
	rps    int
	burst  int
}

// New%s creates a new Redis-based rate limiting middleware
func New%s(redisClient *redis.Client, rps, burst int) *%s {
	return &%s{
		client: redisClient,
		rps:    rps,
		burst:  burst,
	}
}

// RateLimit applies distributed rate limiting
func (m *%s) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("rate_limit:%%s", ip)
		
		ctx := context.Background()
		
		// Get current count
		current, err := m.client.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			// If Redis is unavailable, allow the request
			c.Next()
			return
		}
		
		if current >= m.rps {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"retry_after": "60s",
			})
			c.Abort()
			return
		}
		
		// Increment counter
		pipe := m.client.Pipeline()
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, time.Minute)
		_, err = pipe.Exec(ctx)
		
		if err != nil {
			// If Redis operation fails, allow the request
			c.Next()
			return
		}
		
		c.Next()
	}
}
`, config.GetPackageName(), imports, config.Description, config.GetStructName(),
		config.GetStructName(), config.GetStructName(), config.GetStructName(), config.GetStructName(),
		config.GetStructName())
}

func getBasicMiddlewareTemplate(config models.MiddlewareConfig) string {
	imports := generateImports([]string{
		"net/http",
		"github.com/gin-gonic/gin",
	})
	
	return fmt.Sprintf(`package %s

%s

// %s %s
type %s struct{}

// New%s creates a new middleware
func New%s() *%s {
	return &%s{}
}

// Handle processes the request
func (m *%s) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add middleware logic here
		c.Next()
	}
}
`, config.GetPackageName(), imports, config.Name, config.Description, config.GetStructName(),
		config.GetStructName(), config.GetStructName(), config.GetStructName(), config.GetStructName(),
		config.GetStructName())
}

// Helper functions

func generateImports(imports []string) string {
	if len(imports) == 0 {
		return ""
	}
	
	var importLines []string
	for _, imp := range imports {
		importLines = append(importLines, fmt.Sprintf("\t\"%s\"", imp))
	}
	
	return fmt.Sprintf("import (\n%s\n)", strings.Join(importLines, "\n"))
}

func formatStringSlice(slice []string) string {
	if len(slice) == 0 {
		return "[]string{}"
	}
	
	var quoted []string
	for _, s := range slice {
		quoted = append(quoted, fmt.Sprintf("\"%s\"", s))
	}
	
	return fmt.Sprintf("[]string{%s}", strings.Join(quoted, ", "))
}