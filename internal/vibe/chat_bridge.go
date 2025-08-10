package vibe

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/vibercode/cli/pkg/ui"
)

// ChatBridge conecta el chat terminal con el WebSocket
type ChatBridge struct {
	previewServer    *PreviewServer
	websocketClients map[*websocket.Conn]bool
	terminalChannel  chan ChatMessage
	websocketChannel chan ChatMessage
	mu               sync.RWMutex
	isActive         bool
	logger           *VibeLogger
}

// ChatMessage representa un mensaje de chat
type ChatMessage struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"` // "user", "assistant", "system"
	Content   string                 `json:"content"`
	Timestamp time.Time              `json:"timestamp"`
	Source    string                 `json:"source"` // "terminal", "websocket"
	Data      map[string]interface{} `json:"data,omitempty"`
}

// ChatResponse representa una respuesta del chat
type ChatResponse struct {
	ID        string                 `json:"id"`
	Content   string                 `json:"content"`
	Action    string                 `json:"action,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// NewChatBridge crea un nuevo bridge de chat
func NewChatBridge(previewServer *PreviewServer, logger *VibeLogger) *ChatBridge {
	return &ChatBridge{
		previewServer:    previewServer,
		websocketClients: make(map[*websocket.Conn]bool),
		terminalChannel:  make(chan ChatMessage, 100),
		websocketChannel: make(chan ChatMessage, 100),
		isActive:         true,
		logger:           logger,
	}
}

// Start inicia el bridge de chat
func (cb *ChatBridge) Start() {
	cb.logger.Debug("Starting chat bridge")

	// Goroutine para manejar mensajes del terminal
	go cb.handleTerminalMessages()

	// Goroutine para manejar mensajes del WebSocket
	go cb.handleWebSocketMessages()

	cb.logger.Debug("Chat bridge started successfully")
}

// Stop detiene el bridge de chat
func (cb *ChatBridge) Stop() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.isActive = false
	close(cb.terminalChannel)
	close(cb.websocketChannel)

	cb.logger.Debug("Chat bridge stopped")
}

// RegisterWebSocketClient registra un cliente WebSocket
func (cb *ChatBridge) RegisterWebSocketClient(conn *websocket.Conn) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.websocketClients[conn] = true
	cb.logger.Debug("WebSocket client registered. Total clients: %d", len(cb.websocketClients))
}

// UnregisterWebSocketClient desregistra un cliente WebSocket
func (cb *ChatBridge) UnregisterWebSocketClient(conn *websocket.Conn) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	delete(cb.websocketClients, conn)
	cb.logger.Debug("WebSocket client unregistered. Total clients: %d", len(cb.websocketClients))
}

// SendToTerminal envía un mensaje al terminal
func (cb *ChatBridge) SendToTerminal(message ChatMessage) {
	if !cb.isActive {
		return
	}

	select {
	case cb.terminalChannel <- message:
		cb.logger.Debug("Message sent to terminal channel")
	default:
		cb.logger.Warning("Terminal channel is full, message dropped")
	}
}

// SendToWebSocket envía un mensaje al WebSocket
func (cb *ChatBridge) SendToWebSocket(message ChatMessage) {
	if !cb.isActive {
		return
	}

	select {
	case cb.websocketChannel <- message:
		cb.logger.Debug("Message sent to websocket channel")
	default:
		cb.logger.Warning("WebSocket channel is full, message dropped")
	}
}

// ProcessChatMessage procesa un mensaje de chat desde cualquier fuente
func (cb *ChatBridge) ProcessChatMessage(message ChatMessage) *ChatResponse {
	cb.logger.Debug("Processing chat message from %s: %s", message.Source, message.Content)

	// Procesar con Viber AI a través del preview server
	if cb.previewServer != nil {
		claudeResponse := cb.previewServer.processChatMessageWithContext(message.Content)

		// Convertir Data a map[string]interface{} si es necesario
		var responseData map[string]interface{}
		if claudeResponse.Data != nil {
			if dataMap, ok := claudeResponse.Data.(map[string]interface{}); ok {
				responseData = dataMap
			}
		}

		response := &ChatResponse{
			ID:        generateMessageID(),
			Content:   claudeResponse.Response,
			Action:    claudeResponse.Action,
			Data:      responseData,
			Timestamp: time.Now(),
		}

		// Enviar respuesta a todos los clientes
		cb.broadcastResponse(response)

		return response
	}

	// Respuesta de fallback
	return &ChatResponse{
		ID:        generateMessageID(),
		Content:   "Chat bridge no disponible",
		Timestamp: time.Now(),
	}
}

// handleTerminalMessages maneja los mensajes del terminal
func (cb *ChatBridge) handleTerminalMessages() {
	for message := range cb.terminalChannel {
		if !cb.isActive {
			return
		}

		// Procesar el mensaje
		response := cb.ProcessChatMessage(message)

		// Mostrar respuesta en terminal
		cb.displayTerminalResponse(response)

		// Enviar también a WebSocket clients
		cb.broadcastToWebSocket(message, response)
	}
}

// handleWebSocketMessages maneja los mensajes del WebSocket
func (cb *ChatBridge) handleWebSocketMessages() {
	for message := range cb.websocketChannel {
		if !cb.isActive {
			return
		}

		// Procesar el mensaje
		response := cb.ProcessChatMessage(message)

		// Enviar respuesta de vuelta al WebSocket
		cb.broadcastToWebSocket(message, response)

		// Mostrar también en terminal si no es del terminal
		if message.Source != "terminal" {
			cb.displayTerminalNotification(message, response)
		}
	}
}

// displayTerminalResponse muestra la respuesta en el terminal
func (cb *ChatBridge) displayTerminalResponse(response *ChatResponse) {
	fmt.Printf(ui.Secondary.Sprint("\n🤖 Viber: ")+"%s\n", response.Content)

	// Mostrar información adicional si hay acción
	if response.Action != "" {
		fmt.Printf(ui.Success.Sprint("✨ Acción: ")+"%s\n", response.Action)
	}
}

// displayTerminalNotification muestra una notificación en el terminal
func (cb *ChatBridge) displayTerminalNotification(message ChatMessage, response *ChatResponse) {
	fmt.Printf(ui.Info.Sprint("\n📱 Chat Web: ")+"%s\n", message.Content)
	fmt.Printf(ui.Secondary.Sprint("🤖 Viber: ")+"%s\n", response.Content)
}

// broadcastResponse envía la respuesta a todos los clientes
func (cb *ChatBridge) broadcastResponse(response *ChatResponse) {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	for conn := range cb.websocketClients {
		cb.sendToWebSocketClient(conn, response)
	}
}

// broadcastToWebSocket envía mensaje y respuesta a los clientes WebSocket
func (cb *ChatBridge) broadcastToWebSocket(message ChatMessage, response *ChatResponse) {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	// Crear mensaje de respuesta para WebSocket
	wsMessage := map[string]interface{}{
		"type": "chat_response",
		"data": map[string]interface{}{
			"user_message":       message,
			"assistant_response": response,
		},
		"timestamp": time.Now(),
	}

	messageBytes, err := json.Marshal(wsMessage)
	if err != nil {
		cb.logger.Error("Failed to marshal WebSocket message: %v", err)
		return
	}

	for conn := range cb.websocketClients {
		if err := conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
			cb.logger.Error("Failed to send message to WebSocket client: %v", err)
			// Remover cliente problemático
			delete(cb.websocketClients, conn)
			conn.Close()
		}
	}
}

// sendToWebSocketClient envía un mensaje a un cliente WebSocket específico
func (cb *ChatBridge) sendToWebSocketClient(conn *websocket.Conn, response *ChatResponse) {
	wsMessage := map[string]interface{}{
		"type":      "chat_response",
		"data":      response,
		"timestamp": time.Now(),
	}

	messageBytes, err := json.Marshal(wsMessage)
	if err != nil {
		cb.logger.Error("Failed to marshal WebSocket response: %v", err)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
		cb.logger.Error("Failed to send response to WebSocket client: %v", err)
		// Remover cliente problemático
		delete(cb.websocketClients, conn)
		conn.Close()
	}
}

// SendTerminalMessage envía un mensaje desde el terminal
func (cb *ChatBridge) SendTerminalMessage(content string) {
	message := ChatMessage{
		ID:        generateMessageID(),
		Type:      "user",
		Content:   content,
		Timestamp: time.Now(),
		Source:    "terminal",
	}

	cb.SendToTerminal(message)
}

// HandleWebSocketMessage maneja un mensaje desde WebSocket
func (cb *ChatBridge) HandleWebSocketMessage(conn *websocket.Conn, data []byte) {
	var wsMessage map[string]interface{}
	if err := json.Unmarshal(data, &wsMessage); err != nil {
		cb.logger.Error("Failed to unmarshal WebSocket message: %v", err)
		return
	}

	// Extraer el contenido del mensaje
	content := ""
	if msgData, ok := wsMessage["data"].(map[string]interface{}); ok {
		if msgContent, ok := msgData["message"].(string); ok {
			content = msgContent
		}
	}

	if content == "" {
		cb.logger.Warning("Empty message content from WebSocket")
		return
	}

	message := ChatMessage{
		ID:        generateMessageID(),
		Type:      "user",
		Content:   content,
		Timestamp: time.Now(),
		Source:    "websocket",
		Data:      wsMessage,
	}

	cb.SendToWebSocket(message)
}

// GetStats devuelve estadísticas del bridge
func (cb *ChatBridge) GetStats() map[string]interface{} {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	return map[string]interface{}{
		"active":            cb.isActive,
		"websocket_clients": len(cb.websocketClients),
		"terminal_queue":    len(cb.terminalChannel),
		"websocket_queue":   len(cb.websocketChannel),
	}
}

// generateMessageID genera un ID único para el mensaje
func generateMessageID() string {
	return fmt.Sprintf("msg_%d", time.Now().UnixNano())
}
