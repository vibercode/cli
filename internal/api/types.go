package api

import (
	"time"
	"github.com/vibercode/cli/internal/models"
)

// ============================================================================
// CORE API TYPES
// ============================================================================

// ApiResponse represents a standard API response wrapper
type ApiResponse[T any] struct {
	Success   bool      `json:"success"`
	Data      *T        `json:"data,omitempty"`
	Error     string    `json:"error,omitempty"`
	Message   string    `json:"message,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	RequestID string    `json:"request_id,omitempty"`
}

// ApiError represents a detailed API error
type ApiError struct {
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Details   interface{} `json:"details,omitempty"`
	Field     string      `json:"field,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse[T any] struct {
	Items   []T  `json:"items"`
	Total   int  `json:"total"`
	Page    int  `json:"page"`
	Limit   int  `json:"limit"`
	HasNext bool `json:"has_next"`
	HasPrev bool `json:"has_prev"`
}

// ============================================================================
// SCHEMA API TYPES
// ============================================================================

// CreateSchemaRequest represents a request to create a new schema
type CreateSchemaRequest struct {
	Schema models.ResourceSchema `json:"schema"`
	Source string                `json:"source,omitempty"` // "react-editor", "manual", "import"
}

// CreateSchemaResponse represents the response after creating a schema
type CreateSchemaResponse struct {
	SchemaID            string   `json:"schema_id"`
	CreatedAt           string   `json:"created_at"`
	ValidationWarnings  []string `json:"validation_warnings,omitempty"`
}

// ListSchemasParams represents query parameters for listing schemas
type ListSchemasParams struct {
	Page      int      `form:"page,default=1"`
	Limit     int      `form:"limit,default=10"`
	Search    string   `form:"search"`
	Tags      []string `form:"tags"`
	CreatedBy string   `form:"created_by"`
	Sort      string   `form:"sort,default=created_at"` // "name", "created_at", "updated_at"
	Order     string   `form:"order,default=desc"`      // "asc", "desc"
}

// SchemaImportRequest represents a request to import a schema
type SchemaImportRequest struct {
	Schema    models.ResourceSchema `json:"schema"`
	Overwrite bool                  `json:"overwrite,omitempty"`
	Validate  bool                  `json:"validate,omitempty"`
}

// SchemaExportParams represents query parameters for exporting schemas
type SchemaExportParams struct {
	Format           string `form:"format,default=editor"`     // "editor", "cli", "json"
	IncludeMetadata  bool   `form:"include_metadata"`
	IncludeGenerated bool   `form:"include_generated"`
}

// ============================================================================
// CODE GENERATION TYPES
// ============================================================================

// APIGenerationConfig represents configuration for API generation
type APIGenerationConfig struct {
	SchemaID        string                     `json:"schema_id"`
	ProjectName     string                     `json:"project_name"`
	Description     string                     `json:"description,omitempty"`
	Database        models.DatabaseConfig      `json:"database"`
	Features        []string                   `json:"features"`
	OutputPath      string                     `json:"output_path,omitempty"`
	TemplateVersion string                     `json:"template_version,omitempty"`
	CustomTemplates map[string]string          `json:"custom_templates,omitempty"`
}

// ResourceGenerationConfig represents configuration for resource generation
type ResourceGenerationConfig struct {
	SchemaID       string                    `json:"schema_id"`
	ResourceName   string                    `json:"resource_name"`
	Fields         []models.SchemaField      `json:"fields"`
	Options        models.GenerationOptions  `json:"options"`
	TargetProject  string                    `json:"target_project,omitempty"`
	CustomHandlers bool                      `json:"custom_handlers,omitempty"`
}

// GenerationResult represents the result of a code generation operation
type GenerationResult struct {
	ProjectID     string          `json:"project_id"`
	ProjectName   string          `json:"project_name"`
	FilesGenerated []GeneratedFile `json:"files_generated"`
	ExecutionTime int64           `json:"execution_time"` // milliseconds
	Status        string          `json:"status"`         // "success", "error", "warning", "partial"
	Logs          []LogEntry      `json:"logs"`
	Warnings      []string        `json:"warnings,omitempty"`
	Errors        []string        `json:"errors,omitempty"`
	OutputPath    string          `json:"output_path"`
}

// GeneratedFile represents a file that was generated
type GeneratedFile struct {
	Path           string `json:"path"`
	Type           string `json:"type"`     // "model", "handler", "service", "repository", etc.
	Size           int64  `json:"size"`
	Language       string `json:"language"` // "go", "sql", "yaml", "markdown", "json"
	ContentPreview string `json:"content_preview,omitempty"`
	Checksum       string `json:"checksum,omitempty"`
}

// LogEntry represents a log entry
type LogEntry struct {
	Level     string      `json:"level"`     // "debug", "info", "warn", "error"
	Message   string      `json:"message"`
	Timestamp time.Time   `json:"timestamp"`
	Component string      `json:"component,omitempty"`
	Details   interface{} `json:"details,omitempty"`
}

// ============================================================================
// PROJECT MANAGEMENT TYPES
// ============================================================================

// Project represents a generated project
type Project struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Type        string               `json:"type"` // "api", "resource", "fullstack"
	Status      ProjectStatus        `json:"status"`
	SchemaID    string               `json:"schema_id,omitempty"`
	Database    models.DatabaseConfig `json:"database"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	LastBuild   *time.Time           `json:"last_build,omitempty"`
	BuildInfo   *BuildInfo           `json:"build_info,omitempty"`
	RuntimeInfo *RuntimeInfo         `json:"runtime_info,omitempty"`
}

// ProjectStatus represents the status of a project
type ProjectStatus string

const (
	ProjectStatusCreated   ProjectStatus = "created"
	ProjectStatusBuilding  ProjectStatus = "building"
	ProjectStatusReady     ProjectStatus = "ready"
	ProjectStatusRunning   ProjectStatus = "running"
	ProjectStatusStopped   ProjectStatus = "stopped"
	ProjectStatusError     ProjectStatus = "error"
	ProjectStatusDeploying ProjectStatus = "deploying"
	ProjectStatusDeployed  ProjectStatus = "deployed"
)

// BuildInfo represents build information for a project
type BuildInfo struct {
	Version      string       `json:"version"`
	GoVersion    string       `json:"go_version"`
	BuildTime    time.Time    `json:"build_time"`
	CommitHash   string       `json:"commit_hash,omitempty"`
	BuildFlags   []string     `json:"build_flags,omitempty"`
	Dependencies []Dependency `json:"dependencies"`
}

// Dependency represents a project dependency
type Dependency struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Type    string `json:"type"` // "direct", "indirect"
	License string `json:"license,omitempty"`
}

// RuntimeInfo represents runtime information for a running project
type RuntimeInfo struct {
	PID          *int    `json:"pid,omitempty"`
	Port         *int    `json:"port,omitempty"`
	URL          string  `json:"url,omitempty"`
	HealthStatus string  `json:"health_status,omitempty"` // "healthy", "unhealthy", "unknown"
	MemoryUsage  *int64  `json:"memory_usage,omitempty"`
	CPUUsage     *float64 `json:"cpu_usage,omitempty"`
	Uptime       string  `json:"uptime,omitempty"`
	Logs         []LogEntry `json:"logs"`
}

// ExecutionResult represents the result of project execution
type ExecutionResult struct {
	Success   bool     `json:"success"`
	PID       *int     `json:"pid,omitempty"`
	Port      *int     `json:"port,omitempty"`
	URL       string   `json:"url,omitempty"`
	Command   string   `json:"command"`
	StartTime time.Time `json:"start_time"`
	Logs      []string `json:"logs"`
	Error     string   `json:"error,omitempty"`
}

// ProjectRunConfig represents configuration for running a project
type ProjectRunConfig struct {
	Port       *int              `json:"port,omitempty"`
	EnvVars    map[string]string `json:"env_vars,omitempty"`
	BuildFirst bool              `json:"build_first,omitempty"`
	Detached   bool              `json:"detached,omitempty"`
	LogLevel   string            `json:"log_level,omitempty"` // "debug", "info", "warn", "error"
}

// ============================================================================
// HEALTH AND MONITORING TYPES
// ============================================================================

// HealthCheck represents the health status of the API server
type HealthCheck struct {
	Status    string          `json:"status"`    // "healthy", "degraded", "unhealthy"
	Timestamp time.Time       `json:"timestamp"`
	Services  []ServiceHealth `json:"services"`
	Version   string          `json:"version"`
	Uptime    string          `json:"uptime"`
}

// ServiceHealth represents the health of a specific service
type ServiceHealth struct {
	Name         string      `json:"name"`
	Status       string      `json:"status"` // "up", "down", "degraded"
	ResponseTime *int64      `json:"response_time,omitempty"`
	LastCheck    time.Time   `json:"last_check"`
	Error        string      `json:"error,omitempty"`
	Details      interface{} `json:"details,omitempty"`
}

// Metrics represents server metrics
type Metrics struct {
	RequestsTotal         int64     `json:"requests_total"`
	RequestsPerSecond     float64   `json:"requests_per_second"`
	AverageResponseTime   float64   `json:"average_response_time"`
	ErrorRate             float64   `json:"error_rate"`
	ActiveConnections     int       `json:"active_connections"`
	MemoryUsage           int64     `json:"memory_usage"`
	CPUUsage              float64   `json:"cpu_usage"`
	Timestamp             time.Time `json:"timestamp"`
}

// ============================================================================
// FILE OPERATIONS TYPES
// ============================================================================

// FileUploadRequest represents a file upload request
type FileUploadRequest struct {
	Type      string `form:"type"`      // "schema", "template", "config"
	Overwrite bool   `form:"overwrite"`
	Validate  bool   `form:"validate"`
}

// FileUploadResponse represents the response after uploading a file
type FileUploadResponse struct {
	FileID           string      `json:"file_id"`
	OriginalName     string      `json:"original_name"`
	Size             int64       `json:"size"`
	Type             string      `json:"type"`
	UploadedAt       time.Time   `json:"uploaded_at"`
	ValidationResult interface{} `json:"validation_result,omitempty"`
}

// DownloadRequest represents a project download request
type DownloadRequest struct {
	ProjectID       string `form:"project_id"`
	FilePath        string `form:"file_path"`
	Format          string `form:"format,default=zip"`      // "zip", "tar", "individual"
	IncludeSource   bool   `form:"include_source"`
	IncludeBinaries bool   `form:"include_binaries"`
	IncludeDocs     bool   `form:"include_docs"`
}

// ============================================================================
// STREAMING AND REAL-TIME TYPES
// ============================================================================

// StreamingResponse represents a streaming response
type StreamingResponse struct {
	Type        string      `json:"type"`        // "progress", "log", "error", "complete"
	Data        interface{} `json:"data"`
	Timestamp   time.Time   `json:"timestamp"`
	OperationID string      `json:"operation_id"`
}

// ProgressUpdate represents a progress update
type ProgressUpdate struct {
	Operation              string  `json:"operation"`
	Progress               int     `json:"progress"` // 0-100
	CurrentStep            string  `json:"current_step"`
	TotalSteps             int     `json:"total_steps"`
	CurrentStepIndex       int     `json:"current_step_index"`
	EstimatedTimeRemaining *int64  `json:"estimated_time_remaining,omitempty"`
	Details                string  `json:"details,omitempty"`
}

// WebSocketMessage represents a WebSocket message
type WebSocketMessage struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"` // "request", "response", "event", "error"
	Action    string      `json:"action"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// ============================================================================
// UTILITY TYPES
// ============================================================================

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Host          string
	Port          int
	Mode          string   // "api", "dev", "production"
	CORSOrigins   []string
	EnableSwagger bool
}

// RequestMetadata represents metadata about an API request
type RequestMetadata struct {
	RequestID string
	UserAgent string
	ClientIP  string
	StartTime time.Time
}

// ============================================================================
// ERROR CODES
// ============================================================================

const (
	// General errors
	ErrorCodeInternalServer = "INTERNAL_SERVER_ERROR"
	ErrorCodeBadRequest     = "BAD_REQUEST"
	ErrorCodeNotFound       = "NOT_FOUND"
	ErrorCodeValidation     = "VALIDATION_ERROR"
	ErrorCodeUnauthorized   = "UNAUTHORIZED"
	ErrorCodeForbidden      = "FORBIDDEN"
	ErrorCodeConflict       = "CONFLICT"
	ErrorCodeTooManyRequests = "TOO_MANY_REQUESTS"

	// Schema errors
	ErrorCodeSchemaNotFound     = "SCHEMA_NOT_FOUND"
	ErrorCodeSchemaExists       = "SCHEMA_EXISTS"
	ErrorCodeSchemaInvalid      = "SCHEMA_INVALID"
	ErrorCodeSchemaValidation   = "SCHEMA_VALIDATION_ERROR"

	// Generation errors
	ErrorCodeGenerationFailed   = "GENERATION_FAILED"
	ErrorCodeTemplateNotFound   = "TEMPLATE_NOT_FOUND"
	ErrorCodeOutputPathInvalid  = "OUTPUT_PATH_INVALID"

	// Project errors
	ErrorCodeProjectNotFound    = "PROJECT_NOT_FOUND"
	ErrorCodeProjectExists      = "PROJECT_EXISTS"
	ErrorCodeProjectRunning     = "PROJECT_RUNNING"
	ErrorCodeProjectNotRunning  = "PROJECT_NOT_RUNNING"
	ErrorCodeProjectBuildFailed = "PROJECT_BUILD_FAILED"

	// File errors
	ErrorCodeFileNotFound       = "FILE_NOT_FOUND"
	ErrorCodeFileUploadFailed   = "FILE_UPLOAD_FAILED"
	ErrorCodeFileValidationFailed = "FILE_VALIDATION_FAILED"
)

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// NewSuccessResponse creates a successful API response
func NewSuccessResponse[T any](data T, message string, requestID string) ApiResponse[T] {
	return ApiResponse[T]{
		Success:   true,
		Data:      &data,
		Message:   message,
		Timestamp: time.Now(),
		RequestID: requestID,
	}
}

// NewErrorResponse creates an error API response
func NewErrorResponse[T any](errorCode, errorMessage, requestID string) ApiResponse[T] {
	return ApiResponse[T]{
		Success:   false,
		Error:     errorCode + ": " + errorMessage,
		Timestamp: time.Now(),
		RequestID: requestID,
	}
}

// NewValidationErrorResponse creates a validation error response
func NewValidationErrorResponse[T any](errors []ValidationError, requestID string) ApiResponse[T] {
	return ApiResponse[T]{
		Success:   false,
		Error:     ErrorCodeValidation,
		Timestamp: time.Now(),
		RequestID: requestID,
	}
}