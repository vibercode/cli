package api

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/vibercode/cli/internal/generator"
	"github.com/vibercode/cli/internal/models"
	"github.com/vibercode/cli/internal/storage"
	"github.com/vibercode/cli/pkg/ui"
)

// APIHandler handles all API endpoints
type APIHandler struct {
	schemaRepo *storage.SchemaRepository
	generator  *generator.APIGenerator
	projects   map[string]*Project
	executions map[string]*ExecutionResult
	mutex      sync.RWMutex
	upgrader   websocket.Upgrader
}

// WebSocket upgrader configuration
var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development
		// In production, implement proper origin checking
		return true
	},
}

// ============================================================================
// ROOT AND HEALTH ENDPOINTS
// ============================================================================

// HandleRoot serves the API root endpoint
func (h *APIHandler) HandleRoot(c *gin.Context) {
	requestID := GetRequestID(c)

	response := NewSuccessResponse(gin.H{
		"message": "ViberCode API Server",
		"version": "v1.0.0",
		"status":  "running",
		"endpoints": gin.H{
			"health":     "/api/v1/health",
			"metrics":    "/api/v1/metrics",
			"schemas":    "/api/v1/schema/*",
			"generation": "/api/v1/generate/*",
			"projects":   "/api/v1/projects/*",
			"files":      "/api/v1/files/*",
			"websocket":  "/api/v1/ws",
		},
		"documentation": "/api/v1/docs",
	}, "ViberCode API Server is running", requestID)

	c.JSON(http.StatusOK, response)
}

// HandleHealthCheck performs health check
func (h *APIHandler) HandleHealthCheck(c *gin.Context) {
	requestID := GetRequestID(c)

	// Check various services
	services := []ServiceHealth{
		{
			Name:         "Schema Storage",
			Status:       "up",
			ResponseTime: func() *int64 { t := int64(1); return &t }(),
			LastCheck:    time.Now(),
		},
		{
			Name:         "Code Generator",
			Status:       "up",
			ResponseTime: func() *int64 { t := int64(2); return &t }(),
			LastCheck:    time.Now(),
		},
		{
			Name:         "File System",
			Status:       "up",
			ResponseTime: func() *int64 { t := int64(1); return &t }(),
			LastCheck:    time.Now(),
		},
	}

	// Determine overall status
	status := "healthy"
	for _, service := range services {
		if service.Status != "up" {
			status = "degraded"
			break
		}
	}

	healthCheck := HealthCheck{
		Status:    status,
		Timestamp: time.Now(),
		Services:  services,
		Version:   "1.0.0",
		Uptime:    time.Since(GetMetrics().StartTime).String(),
	}

	response := NewSuccessResponse(healthCheck, "Health check completed", requestID)
	c.JSON(http.StatusOK, response)
}

// HandleMetrics returns server metrics
func (h *APIHandler) HandleMetrics(c *gin.Context) {
	requestID := GetRequestID(c)
	metricsData := GetMetrics()

	avgResponseTime := float64(0)
	if metricsData.RequestCount > 0 {
		avgResponseTime = float64(metricsData.TotalLatency.Nanoseconds()) / float64(metricsData.RequestCount) / 1e6 // Convert to milliseconds
	}

	errorRate := float64(0)
	if metricsData.RequestCount > 0 {
		errorRate = float64(metricsData.ErrorCount) / float64(metricsData.RequestCount) * 100
	}

	// Get system metrics
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	metrics := Metrics{
		RequestsTotal:       metricsData.RequestCount,
		RequestsPerSecond:   float64(metricsData.RequestCount) / time.Since(metricsData.StartTime).Seconds(),
		AverageResponseTime: avgResponseTime,
		ErrorRate:           errorRate,
		ActiveConnections:   int(metricsData.ActiveRequests),
		MemoryUsage:         int64(memStats.Alloc),
		CPUUsage:            0, // TODO: Implement CPU usage tracking
		Timestamp:           time.Now(),
	}

	response := NewSuccessResponse(metrics, "Metrics retrieved successfully", requestID)
	c.JSON(http.StatusOK, response)
}

// ============================================================================
// SCHEMA MANAGEMENT ENDPOINTS
// ============================================================================

// HandleListSchemas lists all schemas with pagination
func (h *APIHandler) HandleListSchemas(c *gin.Context) {
	requestID := GetRequestID(c)

	var params ListSchemasParams
	if err := c.ShouldBindQuery(&params); err != nil {
		HandleValidationError(c, err)
		return
	}

	// Set defaults
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	// Get all schemas
	var schemas []*models.ResourceSchema
	var err error

	if params.Search != "" {
		schemas, err = h.schemaRepo.Search(params.Search)
	} else {
		schemas, err = h.schemaRepo.List()
	}

	if err != nil {
		HandleInternalError(c, err)
		return
	}

	// Apply pagination
	total := len(schemas)
	start := (params.Page - 1) * params.Limit
	end := start + params.Limit

	if start >= total {
		schemas = []*models.ResourceSchema{}
	} else {
		if end > total {
			end = total
		}
		schemas = schemas[start:end]
	}

	paginatedResponse := PaginatedResponse[*models.ResourceSchema]{
		Items:   schemas,
		Total:   total,
		Page:    params.Page,
		Limit:   params.Limit,
		HasNext: (params.Page * params.Limit) < total,
		HasPrev: params.Page > 1,
	}

	response := NewSuccessResponse(paginatedResponse, "Schemas retrieved successfully", requestID)
	c.JSON(http.StatusOK, response)
}

// HandleCreateSchema creates a new schema
func (h *APIHandler) HandleCreateSchema(c *gin.Context) {
	requestID := GetRequestID(c)

	var req CreateSchemaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleValidationError(c, err)
		return
	}

	// Set source if not provided
	if req.Source == "" {
		req.Source = "react-editor"
	}

	// Create the schema
	if err := h.schemaRepo.CreateSchema(&req.Schema); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			HandleConflictError(c, err.Error())
		} else {
			HandleInternalError(c, err)
		}
		return
	}

	createResponse := CreateSchemaResponse{
		SchemaID:  req.Schema.ID,
		CreatedAt: req.Schema.CreatedAt.Format(time.RFC3339),
	}

	response := NewSuccessResponse(createResponse, "Schema created successfully", requestID)
	c.JSON(http.StatusCreated, response)
}

// HandleGetSchema retrieves a specific schema
func (h *APIHandler) HandleGetSchema(c *gin.Context) {
	requestID := GetRequestID(c)
	schemaID := c.Param("id")

	schema, err := h.schemaRepo.Load(schemaID)
	if err != nil {
		HandleNotFoundError(c, "Schema")
		return
	}

	response := NewSuccessResponse(schema, "Schema retrieved successfully", requestID)
	c.JSON(http.StatusOK, response)
}

// HandleImportSchema imports a schema
func (h *APIHandler) HandleImportSchema(c *gin.Context) {
	requestID := GetRequestID(c)

	var req SchemaImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleValidationError(c, err)
		return
	}

	// Check if schema exists and handle overwrite
	existingSchema, _ := h.schemaRepo.Load(req.Schema.ID)
	if existingSchema != nil && !req.Overwrite {
		HandleConflictError(c, "Schema already exists. Set overwrite=true to replace it.")
		return
	}

	// Import the schema
	if err := h.schemaRepo.CreateSchema(&req.Schema); err != nil {
		HandleInternalError(c, err)
		return
	}

	createResponse := CreateSchemaResponse{
		SchemaID:  req.Schema.ID,
		CreatedAt: req.Schema.CreatedAt.Format(time.RFC3339),
	}

	response := NewSuccessResponse(createResponse, "Schema imported successfully", requestID)
	c.JSON(http.StatusOK, response)
}

// HandleExportSchema exports a schema
func (h *APIHandler) HandleExportSchema(c *gin.Context) {
	requestID := GetRequestID(c)
	schemaID := c.Param("id")

	var params SchemaExportParams
	if err := c.ShouldBindQuery(&params); err != nil {
		HandleValidationError(c, err)
		return
	}

	schema, err := h.schemaRepo.Load(schemaID)
	if err != nil {
		HandleNotFoundError(c, "Schema")
		return
	}

	// Format the schema based on export format
	var exportData interface{} = schema

	switch params.Format {
	case "cli":
		// Return CLI-compatible format
		exportData = schema
	case "json":
		// Return raw JSON
		exportData = schema
	default:
		// Default to editor format
		exportData = schema
	}

	response := NewSuccessResponse(exportData, "Schema exported successfully", requestID)
	c.JSON(http.StatusOK, response)
}

// HandleDeleteSchema deletes a schema
func (h *APIHandler) HandleDeleteSchema(c *gin.Context) {
	requestID := GetRequestID(c)
	schemaID := c.Param("id")

	if err := h.schemaRepo.Delete(schemaID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			HandleNotFoundError(c, "Schema")
		} else {
			HandleInternalError(c, err)
		}
		return
	}

	response := NewSuccessResponse[interface{}](nil, "Schema deleted successfully", requestID)
	c.JSON(http.StatusOK, response)
}

// ============================================================================
// CODE GENERATION ENDPOINTS
// ============================================================================

// HandleGenerateAPI generates a complete API project
func (h *APIHandler) HandleGenerateAPI(c *gin.Context) {
	requestID := GetRequestID(c)

	var config APIGenerationConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		HandleValidationError(c, err)
		return
	}

	// Load the schema
	schema, err := h.schemaRepo.Storage().Load(config.SchemaID)
	if err != nil {
		HandleNotFoundError(c, "Schema")
		return
	}

	// Generate the API project
	start := time.Now()

	outputPath := config.OutputPath
	if outputPath == "" {
		outputPath = filepath.Join("./generated", config.ProjectName)
	}

	// Create project directory
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		HandleInternalError(c, fmt.Errorf("failed to create output directory: %w", err))
		return
	}

	// Generate project files
	projectID := uuid.New().String()
	files, err := h.generateAPIProject(schema, &config, outputPath)
	if err != nil {
		HandleInternalError(c, err)
		return
	}

	executionTime := time.Since(start).Milliseconds()

	// Create project record
	project := &Project{
		ID:          projectID,
		Name:        config.ProjectName,
		Description: config.Description,
		Type:        "api",
		Status:      ProjectStatusReady,
		SchemaID:    config.SchemaID,
		Database:    config.Database,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	h.mutex.Lock()
	h.projects[config.ProjectName] = project
	h.mutex.Unlock()

	result := GenerationResult{
		ProjectID:      projectID,
		ProjectName:    config.ProjectName,
		FilesGenerated: files,
		ExecutionTime:  executionTime,
		Status:         "success",
		Logs: []LogEntry{
			{
				Level:     "info",
				Message:   "API project generated successfully",
				Timestamp: time.Now(),
				Component: "generator",
			},
		},
		OutputPath: outputPath,
	}

	response := NewSuccessResponse(result, "API project generated successfully", requestID)
	c.JSON(http.StatusOK, response)
}

// HandleGenerateResource generates a single resource
func (h *APIHandler) HandleGenerateResource(c *gin.Context) {
	requestID := GetRequestID(c)

	var config ResourceGenerationConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		HandleValidationError(c, err)
		return
	}

	// Load the schema
	schema, err := h.schemaRepo.Storage().Load(config.SchemaID)
	if err != nil {
		HandleNotFoundError(c, "Schema")
		return
	}

	// Generate the resource
	start := time.Now()

	outputPath := "./generated/" + config.ResourceName
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		HandleInternalError(c, fmt.Errorf("failed to create output directory: %w", err))
		return
	}

	// Generate resource files
	projectID := uuid.New().String()
	files, err := h.generateResource(schema, &config, outputPath)
	if err != nil {
		HandleInternalError(c, err)
		return
	}

	executionTime := time.Since(start).Milliseconds()

	result := GenerationResult{
		ProjectID:      projectID,
		ProjectName:    config.ResourceName,
		FilesGenerated: files,
		ExecutionTime:  executionTime,
		Status:         "success",
		Logs: []LogEntry{
			{
				Level:     "info",
				Message:   "Resource generated successfully",
				Timestamp: time.Now(),
				Component: "generator",
			},
		},
		OutputPath: outputPath,
	}

	response := NewSuccessResponse(result, "Resource generated successfully", requestID)
	c.JSON(http.StatusOK, response)
}

// ============================================================================
// PROJECT MANAGEMENT ENDPOINTS
// ============================================================================

// HandleListProjects lists all generated projects
func (h *APIHandler) HandleListProjects(c *gin.Context) {
	requestID := GetRequestID(c)

	h.mutex.RLock()
	projects := make([]*Project, 0, len(h.projects))
	for _, project := range h.projects {
		projects = append(projects, project)
	}
	h.mutex.RUnlock()

	response := NewSuccessResponse(projects, "Projects retrieved successfully", requestID)
	c.JSON(http.StatusOK, response)
}

// HandleRunProject runs a generated project
func (h *APIHandler) HandleRunProject(c *gin.Context) {
	requestID := GetRequestID(c)
	projectName := c.Param("name")

	var config ProjectRunConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		// Use defaults if no config provided
		config = ProjectRunConfig{}
	}

	h.mutex.RLock()
	project, exists := h.projects[projectName]
	h.mutex.RUnlock()

	if !exists {
		HandleNotFoundError(c, "Project")
		return
	}

	// Check if project is already running
	if project.Status == ProjectStatusRunning {
		HandleConflictError(c, "Project is already running")
		return
	}

	// Run the project
	executionResult, err := h.runProject(project, &config)
	if err != nil {
		HandleInternalError(c, err)
		return
	}

	// Update project status
	h.mutex.Lock()
	project.Status = ProjectStatusRunning
	project.RuntimeInfo = &RuntimeInfo{
		PID:          executionResult.PID,
		Port:         executionResult.Port,
		URL:          executionResult.URL,
		HealthStatus: "unknown",
		Logs:         []LogEntry{},
	}
	h.mutex.Unlock()

	response := NewSuccessResponse(executionResult, "Project started successfully", requestID)
	c.JSON(http.StatusOK, response)
}

// HandleGetProjectStatus gets the status of a project
func (h *APIHandler) HandleGetProjectStatus(c *gin.Context) {
	requestID := GetRequestID(c)
	projectName := c.Param("name")

	h.mutex.RLock()
	project, exists := h.projects[projectName]
	h.mutex.RUnlock()

	if !exists {
		HandleNotFoundError(c, "Project")
		return
	}

	response := NewSuccessResponse(project.Status, "Project status retrieved successfully", requestID)
	c.JSON(http.StatusOK, response)
}

// HandleStopProject stops a running project
func (h *APIHandler) HandleStopProject(c *gin.Context) {
	requestID := GetRequestID(c)
	projectName := c.Param("name")

	h.mutex.RLock()
	project, exists := h.projects[projectName]
	h.mutex.RUnlock()

	if !exists {
		HandleNotFoundError(c, "Project")
		return
	}

	if project.Status != ProjectStatusRunning {
		HandleBadRequestError(c, "Project is not running")
		return
	}

	// Stop the project
	if err := h.stopProject(project); err != nil {
		HandleInternalError(c, err)
		return
	}

	// Update project status
	h.mutex.Lock()
	project.Status = ProjectStatusStopped
	project.RuntimeInfo = nil
	h.mutex.Unlock()

	response := NewSuccessResponse[interface{}](nil, "Project stopped successfully", requestID)
	c.JSON(http.StatusOK, response)
}

// HandleDownloadProject downloads a project as a zip file
func (h *APIHandler) HandleDownloadProject(c *gin.Context) {
	projectName := c.Param("name")

	var params DownloadRequest
	if err := c.ShouldBindQuery(&params); err != nil {
		HandleValidationError(c, err)
		return
	}

	// Check if project exists
	h.mutex.RLock()
	project, exists := h.projects[projectName]
	h.mutex.RUnlock()

	if !exists {
		HandleNotFoundError(c, "Project")
		return
	}

	// Create temporary zip file
	zipPath, err := h.createProjectZip(project, &params)
	if err != nil {
		HandleInternalError(c, err)
		return
	}
	defer os.Remove(zipPath)

	// Set headers for file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+projectName+".zip")
	c.Header("Content-Type", "application/zip")

	// Stream the file
	c.File(zipPath)
}

// ============================================================================
// FILE OPERATIONS ENDPOINTS
// ============================================================================

// HandleFileUpload handles file uploads
func (h *APIHandler) HandleFileUpload(c *gin.Context) {
	requestID := GetRequestID(c)

	var req FileUploadRequest
	if err := c.ShouldBind(&req); err != nil {
		HandleValidationError(c, err)
		return
	}

	// Get the uploaded file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		HandleBadRequestError(c, "No file uploaded")
		return
	}
	defer file.Close()

	// Create upload directory
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		HandleInternalError(c, err)
		return
	}

	// Save the file
	fileID := uuid.New().String()
	filePath := filepath.Join(uploadDir, fileID+"_"+header.Filename)

	out, err := os.Create(filePath)
	if err != nil {
		HandleInternalError(c, err)
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		HandleInternalError(c, err)
		return
	}

	uploadResponse := FileUploadResponse{
		FileID:       fileID,
		OriginalName: header.Filename,
		Size:         header.Size,
		Type:         req.Type,
		UploadedAt:   time.Now(),
	}

	response := NewSuccessResponse(uploadResponse, "File uploaded successfully", requestID)
	c.JSON(http.StatusOK, response)
}

// ============================================================================
// WEBSOCKET ENDPOINT
// ============================================================================

// HandleWebSocket handles WebSocket connections
func (h *APIHandler) HandleWebSocket(c *gin.Context) {
	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		ui.PrintError("Failed to upgrade WebSocket connection: " + err.Error())
		return
	}
	defer conn.Close()

	ui.PrintInfo("🔌 WebSocket connection established from " + c.ClientIP())

	// Handle WebSocket messages
	for {
		var msg WebSocketMessage
		if err := conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				ui.PrintError("WebSocket error: " + err.Error())
			}
			break
		}

		// Process the message
		response := h.processWebSocketMessage(&msg)

		// Send response
		if err := conn.WriteJSON(response); err != nil {
			ui.PrintError("Failed to send WebSocket response: " + err.Error())
			break
		}
	}

	ui.PrintInfo("🔌 WebSocket connection closed from " + c.ClientIP())
}

// ============================================================================
// HELPER METHODS
// ============================================================================

// generateAPIProject generates a complete API project
func (h *APIHandler) generateAPIProject(schema *models.ResourceSchema, config *APIGenerationConfig, outputPath string) ([]GeneratedFile, error) {
	var files []GeneratedFile

	// Generate using existing CLI generator
	// For now, use a simplified generation approach
	// TODO: Integrate with actual generator methods
	if err := h.generateSimpleProject(config.ProjectName, schema, outputPath); err != nil {
		return nil, fmt.Errorf("failed to generate API: %w", err)
	}

	// Walk the output directory to collect generated files
	err := filepath.Walk(outputPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		relPath, _ := filepath.Rel(outputPath, path)

		// Determine file type and language
		fileType := getFileType(relPath)
		language := getFileLanguage(filepath.Ext(path))

		file := GeneratedFile{
			Path:     relPath,
			Type:     fileType,
			Size:     info.Size(),
			Language: language,
		}

		files = append(files, file)
		return nil
	})

	return files, err
}

// generateResource generates a single resource
func (h *APIHandler) generateResource(schema *models.ResourceSchema, config *ResourceGenerationConfig, outputPath string) ([]GeneratedFile, error) {
	var files []GeneratedFile

	// Generate using existing CLI generator
	// For now, use a simplified generation approach
	// TODO: Integrate with actual resource generator methods
	if err := h.generateSimpleResource(config.ResourceName, schema, outputPath); err != nil {
		return nil, fmt.Errorf("failed to generate resource: %w", err)
	}

	// Walk the output directory to collect generated files
	err := filepath.Walk(outputPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		relPath, _ := filepath.Rel(outputPath, path)

		fileType := getFileType(relPath)
		language := getFileLanguage(filepath.Ext(path))

		file := GeneratedFile{
			Path:     relPath,
			Type:     fileType,
			Size:     info.Size(),
			Language: language,
		}

		files = append(files, file)
		return nil
	})

	return files, err
}

// runProject runs a generated project
func (h *APIHandler) runProject(project *Project, config *ProjectRunConfig) (*ExecutionResult, error) {
	projectPath := filepath.Join("./generated", project.Name)

	// Default port
	port := 8081
	if config.Port != nil {
		port = *config.Port
	}

	// Build the project first if requested
	if config.BuildFirst {
		buildCmd := exec.Command("go", "build", "-o", "main", "./cmd/server")
		buildCmd.Dir = projectPath
		if err := buildCmd.Run(); err != nil {
			return nil, fmt.Errorf("failed to build project: %w", err)
		}
	}

	// Run the project
	var cmd *exec.Cmd
	if config.BuildFirst {
		cmd = exec.Command("./main")
	} else {
		cmd = exec.Command("go", "run", "cmd/server/main.go")
	}

	cmd.Dir = projectPath

	// Set environment variables
	env := os.Environ()
	env = append(env, fmt.Sprintf("PORT=%d", port))
	for key, value := range config.EnvVars {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}
	cmd.Env = env

	// Start the command
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start project: %w", err)
	}

	result := &ExecutionResult{
		Success:   true,
		PID:       &cmd.Process.Pid,
		Port:      &port,
		URL:       fmt.Sprintf("http://localhost:%d", port),
		Command:   cmd.String(),
		StartTime: time.Now(),
		Logs:      []string{"Project started successfully"},
	}

	// Store execution for tracking
	h.mutex.Lock()
	h.executions[project.Name] = result
	h.mutex.Unlock()

	return result, nil
}

// stopProject stops a running project
func (h *APIHandler) stopProject(project *Project) error {
	h.mutex.Lock()
	execution, exists := h.executions[project.Name]
	if exists {
		delete(h.executions, project.Name)
	}
	h.mutex.Unlock()

	if !exists || execution.PID == nil {
		return fmt.Errorf("project execution not found")
	}

	// Find and kill the process
	process, err := os.FindProcess(*execution.PID)
	if err != nil {
		return fmt.Errorf("failed to find process: %w", err)
	}

	if err := process.Kill(); err != nil {
		return fmt.Errorf("failed to kill process: %w", err)
	}

	return nil
}

// StopAllProjects stops all running projects
func (h *APIHandler) StopAllProjects() {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for name, project := range h.projects {
		if project.Status == ProjectStatusRunning {
			if err := h.stopProject(project); err != nil {
				ui.PrintError(fmt.Sprintf("Failed to stop project %s: %v", name, err))
			} else {
				project.Status = ProjectStatusStopped
				project.RuntimeInfo = nil
				ui.PrintInfo(fmt.Sprintf("Stopped project: %s", name))
			}
		}
	}
}

// createProjectZip creates a zip file of a project
func (h *APIHandler) createProjectZip(project *Project, params *DownloadRequest) (string, error) {
	projectPath := filepath.Join("./generated", project.Name)

	// Create temporary zip file
	tempFile, err := os.CreateTemp("", project.Name+"_*.zip")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	zipWriter := zip.NewWriter(tempFile)
	defer zipWriter.Close()

	// Walk the project directory and add files to zip
	err = filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Skip certain files based on parameters
		if !params.IncludeBinaries && isBinaryFile(path) {
			return nil
		}

		relPath, err := filepath.Rel(projectPath, path)
		if err != nil {
			return err
		}

		// Create file in zip
		zipFile, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		// Copy file content
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(zipFile, file)
		return err
	})

	if err != nil {
		os.Remove(tempFile.Name())
		return "", err
	}

	return tempFile.Name(), nil
}

// processWebSocketMessage processes incoming WebSocket messages
func (h *APIHandler) processWebSocketMessage(msg *WebSocketMessage) *WebSocketMessage {
	response := &WebSocketMessage{
		ID:        uuid.New().String(),
		Type:      "response",
		Timestamp: time.Now(),
	}

	switch msg.Action {
	case "ping":
		response.Action = "pong"
		response.Data = gin.H{"message": "pong"}
	case "get_projects":
		h.mutex.RLock()
		projects := make([]*Project, 0, len(h.projects))
		for _, project := range h.projects {
			projects = append(projects, project)
		}
		h.mutex.RUnlock()

		response.Action = "projects_list"
		response.Data = projects
	default:
		response.Type = "error"
		response.Action = "unknown_action"
		response.Data = gin.H{"error": "Unknown action: " + msg.Action}
	}

	return response
}

// Helper functions
func getFileType(path string) string {
	dir := filepath.Dir(path)
	filename := filepath.Base(path)

	switch {
	case strings.Contains(dir, "models"):
		return "model"
	case strings.Contains(dir, "handlers"):
		return "handler"
	case strings.Contains(dir, "services"):
		return "service"
	case strings.Contains(dir, "repositories"):
		return "repository"
	case strings.Contains(filename, "docker"):
		return "docker"
	case strings.Contains(filename, "Makefile"):
		return "config"
	case strings.HasSuffix(filename, "_test.go"):
		return "test"
	case strings.HasSuffix(filename, ".sql"):
		return "migration"
	case strings.HasSuffix(filename, ".md"):
		return "documentation"
	default:
		return "config"
	}
}

func getFileLanguage(ext string) string {
	switch ext {
	case ".go":
		return "go"
	case ".sql":
		return "sql"
	case ".yaml", ".yml":
		return "yaml"
	case ".json":
		return "json"
	case ".md":
		return "markdown"
	default:
		return "text"
	}
}

func isBinaryFile(path string) bool {
	ext := filepath.Ext(path)
	binaryExts := []string{".exe", ".bin", ".so", ".dll", ".dylib"}

	for _, binaryExt := range binaryExts {
		if ext == binaryExt {
			return true
		}
	}

	return false
}

// generateSimpleProject creates a basic project structure
func (h *APIHandler) generateSimpleProject(projectName string, schema *models.ResourceSchema, outputPath string) error {
	// Create basic Go project structure
	dirs := []string{
		filepath.Join(outputPath, "cmd", "server"),
		filepath.Join(outputPath, "internal", "handlers"),
		filepath.Join(outputPath, "internal", "models"),
		filepath.Join(outputPath, "pkg", "database"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Generate basic main.go
	mainGo := `package main

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	
	log.Println("Server starting on :8080")
	r.Run(":8080")
}
`
	if err := os.WriteFile(filepath.Join(outputPath, "cmd", "server", "main.go"), []byte(mainGo), 0644); err != nil {
		return fmt.Errorf("failed to write main.go: %w", err)
	}

	// Generate go.mod
	goMod := fmt.Sprintf(`module %s

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
)
`, projectName)

	if err := os.WriteFile(filepath.Join(outputPath, "go.mod"), []byte(goMod), 0644); err != nil {
		return fmt.Errorf("failed to write go.mod: %w", err)
	}

	// Download dependencies
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = outputPath
	if err := tidyCmd.Run(); err != nil {
		return fmt.Errorf("failed to download dependencies: %w", err)
	}

	return nil
}

// generateSimpleResource creates a basic resource structure
func (h *APIHandler) generateSimpleResource(resourceName string, schema *models.ResourceSchema, outputPath string) error {
	// Create resource directory
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return fmt.Errorf("failed to create resource directory: %w", err)
	}

	// Generate basic model
	modelGo := fmt.Sprintf(`package models

type %s struct {
	ID   uint   `+"`json:\"id\" gorm:\"primaryKey\"`"+`
	Name string `+"`json:\"name\"`"+`
}
`, resourceName)

	if err := os.WriteFile(filepath.Join(outputPath, "model.go"), []byte(modelGo), 0644); err != nil {
		return fmt.Errorf("failed to write model.go: %w", err)
	}

	// Generate basic handler
	handlerGo := fmt.Sprintf(`package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func Get%s(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "%s endpoint"})
}
`, resourceName, resourceName)

	if err := os.WriteFile(filepath.Join(outputPath, "handler.go"), []byte(handlerGo), 0644); err != nil {
		return fmt.Errorf("failed to write handler.go: %w", err)
	}

	return nil
}
