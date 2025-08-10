package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// Protocol constants
const (
	ProtocolVersion = "2024-11-05"
	Implementation  = "vibercode-cli"
	ImplVersion     = "1.0.0"
)

// JSON-RPC 2.0 message types
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      interface{}   `json:"id,omitempty"`
	Result  interface{}   `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
}

type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// MCP-specific types
type InitializeParams struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ClientCapabilities `json:"capabilities"`
	ClientInfo      ClientInfo         `json:"clientInfo"`
}

type ClientCapabilities struct {
	Roots    *RootsCapability    `json:"roots,omitempty"`
	Sampling *SamplingCapability `json:"sampling,omitempty"`
}

type RootsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

type SamplingCapability struct {
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResult struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      ServerInfo         `json:"serverInfo"`
}

type ServerCapabilities struct {
	Tools     *ToolsCapability     `json:"tools,omitempty"`
	Resources *ResourcesCapability `json:"resources,omitempty"`
}

type ToolsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

type ResourcesCapability struct {
	Subscribe   bool `json:"subscribe,omitempty"`
	ListChanged bool `json:"listChanged,omitempty"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema InputSchema `json:"inputSchema"`
}

type InputSchema struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Required   []string               `json:"required,omitempty"`
}

type CallToolParams struct {
	Name      string      `json:"name"`
	Arguments interface{} `json:"arguments,omitempty"`
}

type CallToolResult struct {
	Content []ToolContent `json:"content"`
}

type ToolContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// MCPProtocol handles the MCP communication protocol
type MCPProtocol struct {
	reader       *bufio.Scanner
	writer       io.Writer
	initialized  bool
	tools        map[string]Tool
	toolHandlers map[string]func(interface{}) (interface{}, error)
}

// NewMCPProtocol creates a new MCP protocol handler
func NewMCPProtocol() *MCPProtocol {
	return &MCPProtocol{
		reader:       bufio.NewScanner(os.Stdin),
		writer:       os.Stdout,
		initialized:  false,
		tools:        make(map[string]Tool),
		toolHandlers: make(map[string]func(interface{}) (interface{}, error)),
	}
}

// Start begins the MCP protocol communication loop
func (p *MCPProtocol) Start() error {
	for p.reader.Scan() {
		line := strings.TrimSpace(p.reader.Text())
		if line == "" {
			continue
		}

		var request JSONRPCRequest
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			p.sendError(nil, -32700, "Parse error", err.Error())
			continue
		}

		if err := p.handleRequest(&request); err != nil {
			// Don't print errors as they interfere with JSON-RPC protocol
			// Just continue processing
		}
	}

	// Return any scanner errors
	return p.reader.Err()
}

// handleRequest processes incoming JSON-RPC requests
func (p *MCPProtocol) handleRequest(request *JSONRPCRequest) error {
	switch request.Method {
	case "initialize":
		return p.handleInitialize(request)
	case "initialized":
		// Client notification that initialization is complete - no response needed
		return nil
	case "tools/list":
		return p.handleToolsList(request)
	case "tools/call":
		return p.handleToolsCall(request)
	default:
		return p.sendError(request.ID, -32601, "Method not found", request.Method)
	}
}

// handleInitialize handles the initialization handshake
func (p *MCPProtocol) handleInitialize(request *JSONRPCRequest) error {
	var params InitializeParams
	if request.Params != nil {
		paramsBytes, _ := json.Marshal(request.Params)
		if err := json.Unmarshal(paramsBytes, &params); err != nil {
			return p.sendError(request.ID, -32602, "Invalid params", err.Error())
		}
	}

	// Register tools
	p.registerTools()

	result := InitializeResult{
		ProtocolVersion: ProtocolVersion,
		Capabilities: ServerCapabilities{
			Tools: &ToolsCapability{
				ListChanged: false,
			},
		},
		ServerInfo: ServerInfo{
			Name:    Implementation,
			Version: ImplVersion,
		},
	}

	p.initialized = true
	return p.sendResponse(request.ID, result)
}

// handleToolsList returns the list of available tools
func (p *MCPProtocol) handleToolsList(request *JSONRPCRequest) error {
	if !p.initialized {
		return p.sendError(request.ID, -32002, "Server not initialized", nil)
	}

	tools := make([]Tool, 0, len(p.tools))
	for _, tool := range p.tools {
		tools = append(tools, tool)
	}

	result := map[string]interface{}{
		"tools": tools,
	}

	return p.sendResponse(request.ID, result)
}

// handleToolsCall executes a tool call
func (p *MCPProtocol) handleToolsCall(request *JSONRPCRequest) error {
	if !p.initialized {
		return p.sendError(request.ID, -32002, "Server not initialized", nil)
	}

	var params CallToolParams
	if request.Params != nil {
		paramsBytes, _ := json.Marshal(request.Params)
		if err := json.Unmarshal(paramsBytes, &params); err != nil {
			return p.sendError(request.ID, -32602, "Invalid params", err.Error())
		}
	}

	handler, exists := p.toolHandlers[params.Name]
	if !exists {
		return p.sendError(request.ID, -32601, "Tool not found", params.Name)
	}

	// Execute the tool
	result, err := handler(params.Arguments)
	if err != nil {
		return p.sendError(request.ID, -32603, "Tool execution failed", err.Error())
	}

	// Format result as MCP tool response
	toolResult := CallToolResult{
		Content: []ToolContent{
			{
				Type: "text",
				Text: fmt.Sprintf("%v", result),
			},
		},
	}

	return p.sendResponse(request.ID, toolResult)
}

// sendResponse sends a JSON-RPC response
func (p *MCPProtocol) sendResponse(id interface{}, result interface{}) error {
	response := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}

	return p.writeJSON(response)
}

// sendError sends a JSON-RPC error response
func (p *MCPProtocol) sendError(id interface{}, code int, message string, data interface{}) error {
	response := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &JSONRPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}

	return p.writeJSON(response)
}

// writeJSON writes a JSON message to stdout
func (p *MCPProtocol) writeJSON(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(p.writer, "%s\n", data)
	return err
}

// registerTools registers all available MCP tools
func (p *MCPProtocol) registerTools() {
	// Tool: vibe_start
	p.tools["vibe_start"] = Tool{
		Name:        "vibe_start",
		Description: "Inicia el modo vibe con chat AI y preview en vivo",
		InputSchema: InputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"mode": map[string]interface{}{
					"type":        "string",
					"description": "Modo de vibe: 'general' o 'component'",
					"enum":        []string{"general", "component"},
				},
				"port": map[string]interface{}{
					"type":        "integer",
					"description": "Puerto para el servidor WebSocket",
					"default":     3001,
				},
			},
		},
	}
	p.toolHandlers["vibe_start"] = p.handleVibeStart

	// Tool: component_update
	p.tools["component_update"] = Tool{
		Name:        "component_update",
		Description: "Actualiza las propiedades de un componente en tiempo real",
		InputSchema: InputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"componentId": map[string]interface{}{
					"type":        "string",
					"description": "ID del componente a actualizar",
				},
				"action": map[string]interface{}{
					"type":        "string",
					"description": "Acción a realizar",
					"enum":        []string{"add", "update", "remove"},
				},
				"properties": map[string]interface{}{
					"type":        "object",
					"description": "Nuevas propiedades del componente",
				},
			},
			Required: []string{"componentId"},
		},
	}
	p.toolHandlers["component_update"] = p.handleComponentUpdate

	// Tool: generate_code
	p.tools["generate_code"] = Tool{
		Name:        "generate_code",
		Description: "Genera código Go API basado en un schema",
		InputSchema: InputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"project_name": map[string]interface{}{
					"type":        "string",
					"description": "Nombre del proyecto a generar",
				},
				"database": map[string]interface{}{
					"type":        "string",
					"description": "Tipo de base de datos",
					"enum":        []string{"postgres", "mysql", "sqlite", "mongodb"},
				},
				"features": map[string]interface{}{
					"type":        "array",
					"description": "Características a incluir",
					"items": map[string]interface{}{
						"type": "string",
						"enum": []string{"auth", "swagger", "docker", "tests"},
					},
				},
				"schema": map[string]interface{}{
					"type":        "object",
					"description": "Schema de recursos y modelos",
				},
			},
			Required: []string{"project_name"},
		},
	}
	p.toolHandlers["generate_code"] = p.handleGenerateCode

	// Tool: project_status
	p.tools["project_status"] = Tool{
		Name:        "project_status",
		Description: "Obtiene el estado del proyecto actual y proyectos generados",
		InputSchema: InputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}
	p.toolHandlers["project_status"] = p.handleProjectStatus
}
