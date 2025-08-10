package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vibercode/cli/internal/generator"
	"github.com/vibercode/cli/internal/storage"
	"github.com/vibercode/cli/pkg/ui"
)

// Server represents the HTTP API server
type Server struct {
	config     *ServerConfig
	router     *gin.Engine
	httpServer *http.Server
	handler    *APIHandler
}

// NewServer creates a new API server instance
func NewServer(config *ServerConfig) (*Server, error) {
	if config == nil {
		return nil, fmt.Errorf("server config is required")
	}

	// Set Gin mode based on server mode
	switch config.Mode {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	case "dev":
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.New()

	// Initialize storage and handler
	schemaStorage := storage.NewFileSchemaStorage(storage.GetDefaultSchemaPath())
	schemaRepo := storage.NewSchemaRepository(schemaStorage)

	apiHandler := &APIHandler{
		schemaRepo: schemaRepo,
		generator:  generator.NewAPIGenerator(),
		projects:   make(map[string]*Project),
		executions: make(map[string]*ExecutionResult),
	}

	// Create server
	server := &Server{
		config:  config,
		router:  router,
		handler: apiHandler,
	}

	// Setup middleware
	server.setupMiddleware()

	// Setup routes
	server.setupRoutes()

	// Create HTTP server
	server.httpServer = &http.Server{
		Addr:           fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler:        router,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	return server, nil
}

// setupMiddleware configures all middleware
func (s *Server) setupMiddleware() {
	// Recovery middleware (must be first)
	s.router.Use(ErrorHandlingMiddleware())

	// Request ID middleware
	s.router.Use(RequestIDMiddleware())

	// CORS middleware
	s.router.Use(CORSMiddleware(s.config.CORSOrigins))

	// Security headers
	s.router.Use(SecurityHeadersMiddleware())

	// Logging middleware
	s.router.Use(LoggingMiddleware())

	// Metrics middleware
	s.router.Use(MetricsMiddleware())

	// Rate limiting (if enabled)
	if s.config.Mode != "dev" {
		s.router.Use(RateLimitMiddleware())
	}

	// Validation middleware
	s.router.Use(ValidationMiddleware())
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	// Root endpoint
	s.router.GET("/", s.handler.HandleRoot)

	// API version group
	v1 := s.router.Group("/api/v1")
	{
		// Health and monitoring endpoints
		v1.GET("/health", s.handler.HandleHealthCheck)
		v1.GET("/metrics", s.handler.HandleMetrics)

		// Schema management endpoints
		schemaGroup := v1.Group("/schema")
		{
			schemaGroup.GET("/list", s.handler.HandleListSchemas)
			schemaGroup.POST("/create", s.handler.HandleCreateSchema)
			schemaGroup.GET("/:id", s.handler.HandleGetSchema)
			schemaGroup.POST("/import", s.handler.HandleImportSchema)
			schemaGroup.GET("/:id/export", s.handler.HandleExportSchema)
			schemaGroup.DELETE("/:id", s.handler.HandleDeleteSchema)
		}

		// Code generation endpoints
		generateGroup := v1.Group("/generate")
		{
			generateGroup.POST("/api", s.handler.HandleGenerateAPI)
			generateGroup.POST("/resource", s.handler.HandleGenerateResource)
		}

		// Project management endpoints
		projectsGroup := v1.Group("/projects")
		{
			projectsGroup.GET("/list", s.handler.HandleListProjects)
			projectsGroup.POST("/:name/run", s.handler.HandleRunProject)
			projectsGroup.GET("/:name/status", s.handler.HandleGetProjectStatus)
			projectsGroup.POST("/:name/stop", s.handler.HandleStopProject)
			projectsGroup.GET("/:name/download", s.handler.HandleDownloadProject)
		}

		// File operations endpoints
		fileGroup := v1.Group("/files")
		{
			fileGroup.POST("/upload", s.handler.HandleFileUpload)
		}

		// WebSocket endpoints (for real-time updates)
		v1.GET("/ws", s.handler.HandleWebSocket)
	}

	// Swagger documentation (if enabled)
	if s.config.EnableSwagger {
		s.router.Static("/docs", "./docs")
		s.router.GET("/api/v1/docs", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "API Documentation",
				"version": "v1",
				"endpoints": map[string]interface{}{
					"health":     "/api/v1/health",
					"metrics":    "/api/v1/metrics",
					"schemas":    "/api/v1/schema/*",
					"generation": "/api/v1/generate/*",
					"projects":   "/api/v1/projects/*",
					"files":      "/api/v1/files/*",
					"websocket":  "/api/v1/ws",
				},
			})
		})
	}

	// 404 handler
	s.router.NoRoute(func(c *gin.Context) {
		requestID := c.GetString("request_id")
		response := NewErrorResponse[interface{}](
			ErrorCodeNotFound,
			"Endpoint not found",
			requestID,
		)
		c.JSON(http.StatusNotFound, response)
	})
}

// Start starts the HTTP server
func (s *Server) Start() error {
	ui.PrintInfo(fmt.Sprintf("🌍 Starting API server on %s", s.httpServer.Addr))

	// Log server configuration
	ui.PrintKeyValue("🎯 Mode", s.config.Mode)
	ui.PrintKeyValue("🔀 CORS Origins", fmt.Sprintf("%v", s.config.CORSOrigins))
	ui.PrintKeyValue("📚 Swagger Docs", fmt.Sprintf("%t", s.config.EnableSwagger))

	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	ui.PrintWarning("🛑 Shutting down API server...")

	// Stop any running projects
	s.handler.StopAllProjects()

	// Shutdown HTTP server
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server gracefully: %w", err)
	}

	ui.PrintSuccess("✅ API server stopped gracefully")
	return nil
}

// GetAddr returns the server address
func (s *Server) GetAddr() string {
	return s.httpServer.Addr
}

// GetConfig returns the server configuration
func (s *Server) GetConfig() *ServerConfig {
	return s.config
}

// GetMetrics returns current server metrics
func (s *Server) GetMetrics() *MetricsData {
	return GetMetrics()
}

// ============================================================================
// SERVER HEALTH INFORMATION
// ============================================================================

// ServerInfo represents information about the server
type ServerInfo struct {
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	Mode      string    `json:"mode"`
	StartTime time.Time `json:"start_time"`
	Uptime    string    `json:"uptime"`
	Address   string    `json:"address"`
	CORS      []string  `json:"cors_origins"`
	Features  []string  `json:"features"`
}

// GetServerInfo returns detailed server information
func (s *Server) GetServerInfo() *ServerInfo {
	uptime := time.Since(GetMetrics().StartTime)

	features := []string{"schema-management", "code-generation", "project-management"}
	if s.config.EnableSwagger {
		features = append(features, "swagger-docs")
	}

	return &ServerInfo{
		Name:      "ViberCode API Server",
		Version:   "1.0.0",
		Mode:      s.config.Mode,
		StartTime: GetMetrics().StartTime,
		Uptime:    uptime.String(),
		Address:   s.httpServer.Addr,
		CORS:      s.config.CORSOrigins,
		Features:  features,
	}
}

// ============================================================================
// MIDDLEWARE HELPERS
// ============================================================================

// RequireJSON middleware ensures request content type is JSON
func RequireJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			contentType := c.GetHeader("Content-Type")
			if contentType != "application/json" && contentType != "application/json; charset=utf-8" {
				requestID := c.GetString("request_id")
				response := NewErrorResponse[interface{}](
					ErrorCodeBadRequest,
					"Content-Type must be application/json",
					requestID,
				)
				c.JSON(http.StatusBadRequest, response)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// ============================================================================
// ERROR HANDLERS
// ============================================================================

// HandlePanic handles panics in HTTP handlers
func HandlePanic(c *gin.Context, recovered interface{}) {
	requestID := c.GetString("request_id")

	ui.PrintError("🚨 Panic in API handler")
	ui.PrintKeyValue("Request ID", requestID)
	ui.PrintKeyValue("Panic", fmt.Sprintf("%v", recovered))
	ui.PrintKeyValue("Path", c.Request.URL.Path)
	ui.PrintKeyValue("Method", c.Request.Method)
	ui.PrintKeyValue("IP", c.ClientIP())

	response := NewErrorResponse[interface{}](
		ErrorCodeInternalServer,
		"Internal server error",
		requestID,
	)

	c.JSON(http.StatusInternalServerError, response)
}

// ============================================================================
// UTILITY FUNCTIONS
// ============================================================================

// GetRequestID gets the request ID from context
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		return requestID.(string)
	}
	return ""
}

// LogRequest logs incoming requests in development mode
func LogRequest(c *gin.Context) {
	if gin.Mode() == gin.DebugMode {
		ui.PrintInfo(fmt.Sprintf("📨 %s %s from %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
		))
	}
}
