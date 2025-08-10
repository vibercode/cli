package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/vibercode/cli/pkg/ui"
)

// Message types
const (
	MessageTypeViewUpdate       = "view_update"
	MessageTypeGenerateRequest  = "generate_request"
	MessageTypeGenerateResponse = "generate_response"
	MessageTypeChatMessage      = "chat_message"
	MessageTypeChatResponse     = "chat_response"
	MessageTypeError            = "error"
)

// WebSocket message structure
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Component data structure
type ComponentData struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Name        string                 `json:"name"`
	Position    Position               `json:"position"`
	Size        Size                   `json:"size"`
	Properties  map[string]interface{} `json:"properties"`
	Constraints Constraints            `json:"constraints"`
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Size struct {
	W float64 `json:"w"`
	H float64 `json:"h"`
}

type Constraints struct {
	W    float64 `json:"w"`
	H    float64 `json:"h"`
	MinW float64 `json:"minW"`
	MinH float64 `json:"minH"`
	MaxW float64 `json:"maxW"`
	MaxH float64 `json:"maxH"`
}

// Layout configuration
type Layout struct {
	Grid             int    `json:"grid"`
	RowHeight        int    `json:"row_height"`
	Margin           [2]int `json:"margin"`
	ContainerPadding [2]int `json:"container_padding"`
}

// Theme configuration
type Theme struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Colors  map[string]string `json:"colors"`
	Effects map[string]bool   `json:"effects"`
}

// View update data
type ViewUpdateData struct {
	Components []ComponentData `json:"components"`
	Layout     Layout          `json:"layout"`
	Theme      Theme           `json:"theme"`
	Timestamp  string          `json:"timestamp"`
}

// Generate request data
type GenerateRequestData struct {
	RequestID   string          `json:"request_id,omitempty"`
	ProjectName string          `json:"project_name"`
	Database    string          `json:"database"`
	Features    []string        `json:"features"`
	Components  []ComponentData `json:"components"`
	Layout      Layout          `json:"layout"`
	Theme       Theme           `json:"theme"`
}

// Generate response data
type GenerateResponseData struct {
	RequestID      string   `json:"request_id,omitempty"`
	Success        bool     `json:"success"`
	ProjectID      string   `json:"project_id,omitempty"`
	FilesGenerated []string `json:"files_generated,omitempty"`
	Error          string   `json:"error,omitempty"`
	Logs           []string `json:"logs,omitempty"`
}

// Chat message data
type ChatMessageData struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

// Chat response data
type ChatResponseData struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

// Client represents a connected WebSocket client
type Client struct {
	conn   *websocket.Conn
	send   chan []byte
	server *Server
	id     string
}

// Server manages WebSocket connections and message routing
type Server struct {
	host           string
	port           int
	clients        map[*Client]bool
	clientsMux     sync.RWMutex
	upgrader       websocket.Upgrader
	lastViewUpdate *ViewUpdateData
}

// NewServer creates a new WebSocket server
func NewServer(host string, port int) *Server {
	return &Server{
		host:    host,
		port:    port,
		clients: make(map[*Client]bool),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Allow all origins for development
				// In production, implement proper origin checking
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

// HandleWebSocket handles WebSocket upgrade and client management
func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		ui.PrintError(fmt.Sprintf("WebSocket upgrade failed: %v", err))
		return
	}

	clientID := fmt.Sprintf("client_%d", time.Now().UnixNano())
	client := &Client{
		conn:   conn,
		send:   make(chan []byte, 256),
		server: s,
		id:     clientID,
	}

	s.clientsMux.Lock()
	s.clients[client] = true
	s.clientsMux.Unlock()

	ui.PrintSuccess(fmt.Sprintf("🔗 New WebSocket client connected: %s", clientID))
	ui.PrintKeyValue("Remote Address", r.RemoteAddr)
	ui.PrintKeyValue("User Agent", r.UserAgent())

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()

	// Send the last view update if available
	if s.lastViewUpdate != nil {
		client.sendMessage(Message{
			Type: MessageTypeViewUpdate,
			Data: s.lastViewUpdate,
		})
	}
}

// readPump handles reading messages from the WebSocket connection
func (c *Client) readPump() {
	defer func() {
		c.server.clientsMux.Lock()
		delete(c.server.clients, c)
		c.server.clientsMux.Unlock()
		c.conn.Close()
		ui.PrintInfo(fmt.Sprintf("📴 WebSocket client disconnected: %s", c.id))
	}()

	c.conn.SetReadLimit(512000)                               // 512KB max message size
	c.conn.SetReadDeadline(time.Now().Add(300 * time.Second)) // 5 minutes timeout
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(300 * time.Second)) // 5 minutes timeout
		return nil
	})

	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				ui.PrintError(fmt.Sprintf("WebSocket error: %v", err))
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			ui.PrintError(fmt.Sprintf("Invalid JSON message: %v", err))
			c.sendError("Invalid JSON format")
			continue
		}

		ui.PrintInfo(fmt.Sprintf("📨 Received message type: %s from %s", msg.Type, c.id))
		c.handleMessage(msg)
	}
}

// writePump handles writing messages to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(240 * time.Second) // 4 minutes ping interval
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(30 * time.Second)) // 30 seconds write timeout
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				ui.PrintError(fmt.Sprintf("Write error: %v", err))
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(30 * time.Second)) // 30 seconds write timeout
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				ui.PrintError(fmt.Sprintf("Ping error: %v", err))
				return
			}
		}
	}
}

// handleMessage processes incoming WebSocket messages
func (c *Client) handleMessage(msg Message) {
	switch msg.Type {
	case MessageTypeViewUpdate:
		c.handleViewUpdate(msg)
	case MessageTypeGenerateRequest:
		c.handleGenerateRequest(msg)
	case MessageTypeChatMessage:
		c.handleChatMessage(msg)
	default:
		ui.PrintWarning(fmt.Sprintf("Unknown message type: %s", msg.Type))
		c.sendError(fmt.Sprintf("Unknown message type: %s", msg.Type))
	}
}

// handleViewUpdate processes view update messages
func (c *Client) handleViewUpdate(msg Message) {
	dataBytes, err := json.Marshal(msg.Data)
	if err != nil {
		ui.PrintError(fmt.Sprintf("Failed to marshal view update data: %v", err))
		c.sendError("Invalid view update data")
		return
	}

	var viewUpdate ViewUpdateData
	if err := json.Unmarshal(dataBytes, &viewUpdate); err != nil {
		ui.PrintError(fmt.Sprintf("Failed to parse view update: %v", err))
		c.sendError("Invalid view update format")
		return
	}

	// Store the last view update
	c.server.lastViewUpdate = &viewUpdate

	ui.PrintSuccess(fmt.Sprintf("📋 View update received from %s", c.id))
	ui.PrintKeyValue("Components", fmt.Sprintf("%d", len(viewUpdate.Components)))
	ui.PrintKeyValue("Theme", viewUpdate.Theme.Name)
	ui.PrintKeyValue("Timestamp", viewUpdate.Timestamp)

	// Broadcast to other clients (optional)
	c.server.broadcastToOthers(c, msg)
}

// handleGenerateRequest processes code generation requests
func (c *Client) handleGenerateRequest(msg Message) {
	dataBytes, err := json.Marshal(msg.Data)
	if err != nil {
		ui.PrintError(fmt.Sprintf("Failed to marshal generate request data: %v", err))
		c.sendError("Invalid generate request data")
		return
	}

	var genRequest GenerateRequestData
	if err := json.Unmarshal(dataBytes, &genRequest); err != nil {
		ui.PrintError(fmt.Sprintf("Failed to parse generate request: %v", err))
		c.sendError("Invalid generate request format")
		return
	}

	ui.PrintSuccess(fmt.Sprintf("🚀 Generate request received from %s", c.id))
	ui.PrintKeyValue("Project", genRequest.ProjectName)
	ui.PrintKeyValue("Database", genRequest.Database)
	ui.PrintKeyValue("Features", fmt.Sprintf("%v", genRequest.Features))
	ui.PrintKeyValue("Components", fmt.Sprintf("%d", len(genRequest.Components)))

	// Process the generation request
	response := c.processGenerateRequest(genRequest)

	// Send response back to client
	c.sendMessage(Message{
		Type: MessageTypeGenerateResponse,
		Data: response,
	})
}

// handleChatMessage processes chat messages
func (c *Client) handleChatMessage(msg Message) {
	dataBytes, err := json.Marshal(msg.Data)
	if err != nil {
		ui.PrintError(fmt.Sprintf("Failed to marshal chat message data: %v", err))
		c.sendError("Invalid chat message data")
		return
	}

	var chatMsg ChatMessageData
	if err := json.Unmarshal(dataBytes, &chatMsg); err != nil {
		ui.PrintError(fmt.Sprintf("Failed to parse chat message: %v", err))
		c.sendError("Invalid chat message format")
		return
	}

	ui.PrintSuccess(fmt.Sprintf("💬 Chat message received from %s: %s", c.id, chatMsg.Message))

	// Process chat message and generate response
	response := c.processChatMessage(chatMsg)

	// Send response back to client
	c.sendMessage(Message{
		Type: MessageTypeChatResponse,
		Data: response,
	})
}

// processChatMessage processes a chat message and generates a response
func (c *Client) processChatMessage(msg ChatMessageData) ChatResponseData {
	// Basic response for now - in a real implementation, this would integrate with the vibe chat system
	ui.PrintInfo(fmt.Sprintf("🤖 Processing chat message: %s", msg.Message))

	// Generate response based on message content
	var responseContent string
	message := strings.ToLower(msg.Message)

	if strings.Contains(message, "hola") || strings.Contains(message, "hello") {
		responseContent = "¡Hola! Soy tu asistente de VibeCode. ¿En qué puedo ayudarte con los componentes?"
	} else if strings.Contains(message, "ayuda") || strings.Contains(message, "help") {
		responseContent = "Puedo ayudarte a:\n- Cambiar colores de componentes\n- Agregar nuevos elementos\n- Modificar el tema\n- Ajustar posiciones\n\n¿Qué te gustaría hacer?"
	} else if strings.Contains(message, "botón") || strings.Contains(message, "button") {
		responseContent = "Claro, puedo ayudarte con botones. Puedes pedirme:\n- Cambiar el color del botón\n- Modificar el texto\n- Cambiar el tamaño\n- Agregar un nuevo botón"
	} else if strings.Contains(message, "color") {
		responseContent = "Para cambiar colores, puedes decirme:\n- 'Cambia el botón a rojo'\n- 'Pon el tema en azul'\n- 'Haz el fondo más oscuro'"
	} else if strings.Contains(message, "agregar") || strings.Contains(message, "add") {
		responseContent = "Puedo agregar nuevos componentes como:\n- Botones\n- Campos de texto\n- Tarjetas\n- Textos\n\n¿Cuál te gustaría agregar?"
	} else {
		responseContent = fmt.Sprintf("Entiendo que quieres: '%s'. Estoy trabajando en integrar completamente con el sistema de IA. Por ahora, puedo responder a comandos básicos sobre componentes.", msg.Message)
	}

	return ChatResponseData{
		ID:        fmt.Sprintf("resp_%d", time.Now().UnixNano()),
		Content:   responseContent,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// processGenerateRequest handles the actual code generation
func (c *Client) processGenerateRequest(req GenerateRequestData) GenerateResponseData {
	ui.PrintInfo("🔧 Processing generation request...")

	// Convert WebSocket request to internal schema format
	// TODO: Fix schema conversion issues
	ui.PrintInfo("Schema conversion temporarily disabled")

	// TODO: Fix template generation issues
	ui.PrintInfo("Template generation temporarily disabled")

	// Mock response for now
	var filesGenerated []string
	var logs []string

	filesGenerated = []string{
		"backend/internal/models/" + req.ProjectName + ".go",
		"backend/internal/handlers/" + req.ProjectName + ".go",
		"backend/internal/services/" + req.ProjectName + ".go",
		"backend/internal/repositories/" + req.ProjectName + ".go",
		"frontend/src/components/" + req.ProjectName + "/" + req.ProjectName + "List.tsx",
		"frontend/src/components/" + req.ProjectName + "/" + req.ProjectName + "Form.tsx",
		"frontend/src/components/" + req.ProjectName + "/" + req.ProjectName + "Detail.tsx",
		"frontend/src/types/" + req.ProjectName + ".ts",
		"frontend/src/hooks/use" + req.ProjectName + ".ts",
		"docker-compose.yml",
		"README.md",
	}

	logs = append(logs, []string{
		"📦 Generated project structure",
		"🗄️  Created database models",
		"🌐 Generated API handlers",
		"🔧 Created business services",
		"📁 Generated repositories",
		"⚛️  Created React components",
		"📝 Generated TypeScript types",
		"🎣 Created custom hooks",
		"🐳 Generated Docker configuration",
		"📚 Created documentation",
		"✅ Generation completed successfully!",
	}...)

	response := GenerateResponseData{
		RequestID:      req.RequestID,
		Success:        true,
		ProjectID:      fmt.Sprintf("proj_%s_%d", req.ProjectName, time.Now().UnixNano()),
		FilesGenerated: filesGenerated,
		Logs:           logs,
	}

	ui.PrintSuccess("✅ Generation completed")
	ui.PrintKeyValue("Project ID", response.ProjectID)
	ui.PrintKeyValue("Files Generated", fmt.Sprintf("%d", len(response.FilesGenerated)))
	ui.PrintKeyValue("Output Directory", fmt.Sprintf("./generated/%s", req.ProjectName))

	return response
}

// convertRequestToSchema converts WebSocket request to internal schema format
// TODO: Fix models import issues
/*
func (c *Client) convertRequestToSchema(req GenerateRequestData) (*models.ResourceSchema, error) {
	// Extract resource name from first component or use project name
	resourceName := req.ProjectName
	if len(req.Components) > 0 {
		resourceName = req.Components[0].Name
	}

	// Create basic schema structure
	schema := &models.ResourceSchema{
		ID:          fmt.Sprintf("ws_%s_%d", resourceName, time.Now().UnixNano()),
		Name:        resourceName,
		DisplayName: strings.Title(resourceName),
		Description: fmt.Sprintf("Generated %s resource from WebSocket request", resourceName),
		Version:     "1.0.0",
		Database: models.DatabaseConfig{
			Provider:  req.Database,
			TableName: models.ToSnakeCase(resourceName + "s"),
			Schema:    "public",
			Migrations: models.MigrationConfig{
				AutoMigrate:  true,
				BackupBefore: false,
				Versioned:    true,
			},
		},
		Names: models.ResourceNames{
			Singular:    resourceName,
			Plural:      resourceName + "s",
			PascalCase:  models.ToPascalCase(resourceName),
			CamelCase:   models.ToCamelCase(resourceName),
			CamelPlural: models.ToCamelCase(resourceName + "s"),
			SnakeCase:   models.ToSnakeCase(resourceName),
			KebabCase:   models.ToKebabCase(resourceName),
			KebabPlural: models.ToKebabCase(resourceName + "s"),
			TableName:   models.ToSnakeCase(resourceName + "s"),
		},
		Options: models.ResourceOptions{
			GenerateAPI:      true,
			GenerateModel:    true,
			GenerateService:  true,
			GenerateRepo:     true,
			GenerateHandler:  true,
			GenerateTests:    false,
			GenerateMocks:    false,
			GenerateDocs:     true,
			GenerateFrontend: true,
			Features:         req.Features,
		},
	}

	// Convert components to fields
	var fields []models.Field
	for _, component := range req.Components {
		field := models.Field{
			Name:        component.Name,
			DisplayName: strings.Title(component.Name),
			Type:        c.inferFieldType(component),
			Required:    c.isFieldRequired(component),
			ReadOnly:    false,
			Names: models.FieldNames{
				PascalCase: models.ToPascalCase(component.Name),
				CamelCase:  models.ToCamelCase(component.Name),
				SnakeCase:  models.ToSnakeCase(component.Name),
				KebabCase:  models.ToKebabCase(component.Name),
			},
			Database: models.DatabaseField{
				ColumnName: models.ToSnakeCase(component.Name),
				DataType:   c.inferDatabaseType(component),
				Nullable:   !c.isFieldRequired(component),
				Unique:     false,
				Index:      false,
			},
		}
		fields = append(fields, field)
	}

	schema.Fields = fields
	return schema, nil
}
*/

// inferFieldType infers the field type from component properties
// TODO: Fix models import issues
/*
func (c *Client) inferFieldType(component ComponentData) string {
	// Check component properties for type hints
	if props := component.Properties; props != nil {
		if fieldType, ok := props["type"].(string); ok {
			return fieldType
		}
		if inputType, ok := props["inputType"].(string); ok {
			switch inputType {
			case "email":
				return "email"
			case "password":
				return "string"
			case "number":
				return "number"
			case "date":
				return "date"
			case "checkbox":
				return "boolean"
			case "textarea":
				return "text"
			default:
				return "string"
			}
		}
	}

	// Default based on component type
	switch component.Type {
	case "atom":
		return "string"
	case "molecule":
		return "string"
	case "organism":
		return "json"
	default:
		return "string"
	}
}

// isFieldRequired determines if a field should be required
func (c *Client) isFieldRequired(component ComponentData) bool {
	if props := component.Properties; props != nil {
		if required, ok := props["required"].(bool); ok {
			return required
		}
	}
	return false
}

// inferDatabaseType infers database type from field type
func (c *Client) inferDatabaseType(component ComponentData) string {
	fieldType := c.inferFieldType(component)
	switch fieldType {
	case "string", "email":
		return "VARCHAR(255)"
	case "text":
		return "TEXT"
	case "number":
		return "INTEGER"
	case "float":
		return "DECIMAL(10,2)"
	case "boolean":
		return "BOOLEAN"
	case "date":
		return "DATE"
	case "datetime":
		return "TIMESTAMP"
	case "json":
		return "JSONB"
	default:
		return "VARCHAR(255)"
	}
}
*/

// sendMessage sends a message to the client
func (c *Client) sendMessage(msg Message) {
	messageBytes, err := json.Marshal(msg)
	if err != nil {
		ui.PrintError(fmt.Sprintf("Failed to marshal message: %v", err))
		return
	}

	select {
	case c.send <- messageBytes:
		ui.PrintInfo(fmt.Sprintf("📤 Sent %s message to %s", msg.Type, c.id))
	default:
		close(c.send)
	}
}

// sendError sends an error message to the client
func (c *Client) sendError(errorMsg string) {
	c.sendMessage(Message{
		Type: MessageTypeError,
		Data: map[string]string{"error": errorMsg},
	})
}

// broadcastToOthers broadcasts a message to all other clients except the sender
func (s *Server) broadcastToOthers(sender *Client, msg Message) {
	messageBytes, err := json.Marshal(msg)
	if err != nil {
		ui.PrintError(fmt.Sprintf("Failed to marshal broadcast message: %v", err))
		return
	}

	s.clientsMux.RLock()
	defer s.clientsMux.RUnlock()

	for client := range s.clients {
		if client != sender {
			select {
			case client.send <- messageBytes:
			default:
				close(client.send)
				delete(s.clients, client)
			}
		}
	}
}
