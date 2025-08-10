package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vibercode/cli/pkg/ui"
)

// ============================================================================
// REQUEST ID MIDDLEWARE
// ============================================================================

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// ============================================================================
// CORS MIDDLEWARE
// ============================================================================

// CORSMiddleware configures CORS for the API
func CORSMiddleware(origins []string) gin.HandlerFunc {
	config := cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID", "X-Total-Count"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	// If origins contain "*", allow all origins
	for _, origin := range origins {
		if origin == "*" {
			config.AllowAllOrigins = true
			config.AllowOrigins = nil
			break
		}
	}

	return cors.New(config)
}

// ============================================================================
// LOGGING MIDDLEWARE
// ============================================================================

// LoggingMiddleware logs HTTP requests with colors and formatting
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		var statusColor, methodColor, resetColor string
		if param.IsOutputColor() {
			statusColor = getStatusColor(param.StatusCode)
			methodColor = getMethodColor(param.Method)
			resetColor = "\033[0m"
		}

		return ui.Dim.Sprintf("[API] %v | %s%3d%s | %13v | %15s | %s%-7s%s %s\n",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.ClientIP,
			methodColor, param.Method, resetColor,
			param.Path,
		)
	})
}

// getStatusColor returns the appropriate color for HTTP status codes
func getStatusColor(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "\033[97;42m" // White text on green background
	case code >= 300 && code < 400:
		return "\033[90;47m" // Dark gray text on white background
	case code >= 400 && code < 500:
		return "\033[97;43m" // White text on yellow background
	default:
		return "\033[97;41m" // White text on red background
	}
}

// getMethodColor returns the appropriate color for HTTP methods
func getMethodColor(method string) string {
	switch method {
	case "GET":
		return "\033[97;44m" // White text on blue background
	case "POST":
		return "\033[97;42m" // White text on green background
	case "PUT":
		return "\033[97;43m" // White text on yellow background
	case "DELETE":
		return "\033[97;41m" // White text on red background
	case "PATCH":
		return "\033[97;45m" // White text on magenta background
	case "HEAD":
		return "\033[97;46m" // White text on cyan background
	case "OPTIONS":
		return "\033[97;47m" // White text on white background
	default:
		return "\033[97;40m" // White text on black background
	}
}

// ============================================================================
// ERROR HANDLING MIDDLEWARE
// ============================================================================

// ErrorHandlingMiddleware handles panics and errors
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		requestID := c.GetString("request_id")
		
		ui.PrintError("🚨 Panic recovered in API handler")
		ui.PrintKeyValue("Request ID", requestID)
		ui.PrintKeyValue("Error", fmt.Sprintf("%v", recovered))
		ui.PrintKeyValue("Path", c.Request.URL.Path)
		ui.PrintKeyValue("Method", c.Request.Method)

		response := NewErrorResponse[interface{}](
			ErrorCodeInternalServer,
			"Internal server error occurred",
			requestID,
		)

		c.JSON(http.StatusInternalServerError, response)
		c.Abort()
	})
}

// ============================================================================
// VALIDATION MIDDLEWARE
// ============================================================================

// ValidationMiddleware handles request validation
func ValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set JSON binding to validate requests
		c.Next()
	}
}

// ============================================================================
// SECURITY MIDDLEWARE
// ============================================================================

// SecurityHeadersMiddleware adds security headers
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// API-specific headers
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		
		c.Next()
	}
}

// ============================================================================
// RATE LIMITING MIDDLEWARE (Basic Implementation)
// ============================================================================

// RateLimitMiddleware provides basic rate limiting
func RateLimitMiddleware() gin.HandlerFunc {
	// Simple in-memory rate limiter
	// In production, consider using Redis or a more sophisticated solution
	return func(c *gin.Context) {
		// For now, just pass through
		// TODO: Implement proper rate limiting
		c.Next()
	}
}

// ============================================================================
// METRICS MIDDLEWARE
// ============================================================================

// MetricsData holds basic metrics
type MetricsData struct {
	RequestCount    int64
	ErrorCount      int64
	TotalLatency    time.Duration
	ActiveRequests  int64
	StartTime       time.Time
}

var metrics = &MetricsData{
	StartTime: time.Now(),
}

// MetricsMiddleware collects basic metrics
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		metrics.ActiveRequests++
		
		c.Next()
		
		metrics.ActiveRequests--
		metrics.RequestCount++
		metrics.TotalLatency += time.Since(start)
		
		if c.Writer.Status() >= 400 {
			metrics.ErrorCount++
		}
	}
}

// GetMetrics returns current metrics
func GetMetrics() *MetricsData {
	return metrics
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// HandleValidationError creates a standardized validation error response
func HandleValidationError(c *gin.Context, err error) {
	requestID := c.GetString("request_id")
	
	response := NewErrorResponse[interface{}](
		ErrorCodeValidation,
		err.Error(),
		requestID,
	)
	
	c.JSON(http.StatusBadRequest, response)
}

// HandleNotFoundError creates a standardized not found error response
func HandleNotFoundError(c *gin.Context, resource string) {
	requestID := c.GetString("request_id")
	
	response := NewErrorResponse[interface{}](
		ErrorCodeNotFound,
		resource+" not found",
		requestID,
	)
	
	c.JSON(http.StatusNotFound, response)
}

// HandleInternalError creates a standardized internal error response
func HandleInternalError(c *gin.Context, err error) {
	requestID := c.GetString("request_id")
	
	ui.PrintError("⚠️  Internal API Error")
	ui.PrintKeyValue("Request ID", requestID)
	ui.PrintKeyValue("Error", err.Error())
	ui.PrintKeyValue("Path", c.Request.URL.Path)
	ui.PrintKeyValue("Method", c.Request.Method)
	
	response := NewErrorResponse[interface{}](
		ErrorCodeInternalServer,
		"An internal error occurred",
		requestID,
	)
	
	c.JSON(http.StatusInternalServerError, response)
}

// HandleConflictError creates a standardized conflict error response
func HandleConflictError(c *gin.Context, message string) {
	requestID := c.GetString("request_id")
	
	response := NewErrorResponse[interface{}](
		ErrorCodeConflict,
		message,
		requestID,
	)
	
	c.JSON(http.StatusConflict, response)
}

// HandleBadRequestError creates a standardized bad request error response
func HandleBadRequestError(c *gin.Context, message string) {
	requestID := c.GetString("request_id")
	
	response := NewErrorResponse[interface{}](
		ErrorCodeBadRequest,
		message,
		requestID,
	)
	
	c.JSON(http.StatusBadRequest, response)
}