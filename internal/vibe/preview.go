package vibe

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/vibercode/cli/internal/vibe/prompts"
)

type PreviewServer struct {
	port          string
	clients       map[*websocket.Conn]bool
	upgrader      websocket.Upgrader
	promptLoader  *prompts.PromptLoader
	currentView   *prompts.CurrentViewState
	claudeClient  *ClaudeClient
	componentMode bool // true for component-focused mode
	logger        *VibeLogger
	chatBridge    *ChatBridge
}

type WebSocketMessage struct {
	Type      string      `json:"type"`
	Action    string      `json:"action,omitempty"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

type ViewUpdateData struct {
	Components []interface{} `json:"components"`
	Layout     interface{}   `json:"layout"`
	Theme      interface{}   `json:"theme"`
}

type LiveUpdateData struct {
	ComponentID string      `json:"componentId"`
	Changes     interface{} `json:"changes"`
	Action      string      `json:"action"` // "update", "add", "remove"
}

type ChatMessageData struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type ChatResponseData struct {
	Response  string      `json:"response"`
	Action    string      `json:"action,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

type UIUpdateRequest struct {
	Type        string      `json:"type"`
	Action      string      `json:"action"`
	Data        interface{} `json:"data"`
	Explanation string      `json:"explanation"`
}

// ViewStateUpdateRequest represents a request to update the current view state
type ViewStateUpdateRequest struct {
	Components []prompts.ComponentState `json:"components"`
	Theme      prompts.ThemeState       `json:"theme"`
	Layout     prompts.LayoutState      `json:"layout"`
	Canvas     prompts.CanvasState      `json:"canvas"`
}

const (
	MessageTypeViewUpdate   = "view_update"
	MessageTypeLiveUpdate   = "live_update"
	MessageTypeChatMessage  = "chat_message"
	MessageTypeChatResponse = "chat_response"
	MessageTypeViewState    = "view_state_update"

	// Write and read timeouts for WebSocket connections
	writeTimeout = 300 * time.Second
	readTimeout  = 300 * time.Second
)

func NewPreviewServer(port string) *PreviewServer {
	// Initialize logger for preview server
	logger := NewVibeLogger(true) // true = chat mode, less verbose

	promptLoader, err := prompts.NewPromptLoader()
	if err != nil {
		logger.Warning("Failed to load prompts: %v", err)
		promptLoader = nil
	}

	// Initialize Claude client
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	claudeClient := NewClaudeClient(apiKey)

	// Initialize with default view state
	defaultView := &prompts.CurrentViewState{
		Components: []prompts.ComponentState{},
		Theme: prompts.ThemeState{
			ID:   "vibercode",
			Name: "VibeCode",
			Colors: map[string]string{
				"primary":    "#8B5CF6",
				"secondary":  "#06B6D4",
				"accent":     "#F59E0B",
				"background": "#0F0F0F",
				"surface":    "#1A1A1A",
				"text":       "#FFFFFF",
			},
			Effects: map[string]interface{}{
				"glow":       true,
				"gradients":  true,
				"animations": true,
			},
		},
		Layout: prompts.LayoutState{
			Grid:             12,
			RowHeight:        60,
			Margin:           [2]int{16, 16},
			ContainerPadding: [2]int{24, 24},
			ShowGrid:         true,
			SnapToGrid:       true,
		},
		Canvas: prompts.CanvasState{
			Viewport:     "desktop",
			Zoom:         1.0,
			PanOffset:    prompts.Position{X: 0, Y: 0},
			SelectedItem: "",
		},
	}

	ps := &PreviewServer{
		port:         port,
		clients:      make(map[*websocket.Conn]bool),
		promptLoader: promptLoader,
		currentView:  defaultView,
		claudeClient: claudeClient,
		logger:       logger,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// Allow connections from any origin for development
				return true
			},
		},
	}

	// Initialize chat bridge
	ps.chatBridge = NewChatBridge(ps, logger)

	return ps
}

func (ps *PreviewServer) Start() error {
	r := mux.NewRouter()

	// Apply CORS middleware to all routes
	r.Use(ps.corsMiddleware)

	// Serve the preview HTML page at root - simplified
	r.HandleFunc("/", ps.handlePreviewPage)
	ps.logger.Debug("Root route registered: /")

	// WebSocket endpoint
	r.HandleFunc("/ws", ps.handleWebSocket)
	ps.logger.Debug("WebSocket route registered: /ws")

	// HTTP API endpoints
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/view-update", ps.handleViewUpdate).Methods("POST", "OPTIONS")
	api.HandleFunc("/live-update", ps.handleLiveUpdate).Methods("POST", "OPTIONS")
	api.HandleFunc("/chat", ps.handleChatMessage).Methods("POST", "OPTIONS")
	api.HandleFunc("/view-state", ps.handleViewStateUpdate).Methods("POST", "OPTIONS")
	api.HandleFunc("/view-state", ps.handleGetViewState).Methods("GET", "OPTIONS")
	api.HandleFunc("/status", ps.handleStatus).Methods("GET", "OPTIONS")
	ps.logger.Debug("API routes registered under /api")

	// Start chat bridge
	ps.chatBridge.Start()
	ps.logger.Debug("Chat bridge started")

	ps.logger.ChatInfo("🚀 Preview server starting on port %s", ps.port)
	ps.logger.Debug("📡 WebSocket: ws://localhost:%s/ws", ps.port)
	ps.logger.Debug("🌐 HTTP API: http://localhost:%s/api", ps.port)
	ps.logger.Debug("🏠 Preview page: http://localhost:%s/", ps.port)

	return http.ListenAndServe(":"+ps.port, r)
}

// Stop stops the preview server and cleans up resources
func (ps *PreviewServer) Stop() {
	if ps.chatBridge != nil {
		ps.chatBridge.Stop()
	}
	if ps.logger != nil {
		ps.logger.Close()
	}
}

// corsMiddleware adds CORS headers to HTTP responses
func (ps *PreviewServer) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (ps *PreviewServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := ps.upgrader.Upgrade(w, r, nil)
	if err != nil {
		ps.logger.Error("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	ps.logger.Debug("New WebSocket connection established")
	ps.clients[conn] = true
	defer delete(ps.clients, conn)

	// Register client in chat bridge
	ps.chatBridge.RegisterWebSocketClient(conn)
	defer ps.chatBridge.UnregisterWebSocketClient(conn)

	// Send current view state to new client
	ps.sendCurrentViewState(conn)

	// Configure connection settings
	conn.SetReadLimit(512 * 1024)                           // 512KB max message size
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))  // 60 seconds read timeout
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second)) // 10 seconds write timeout

	// Set up pong handler to handle heartbeat
	conn.SetPongHandler(func(string) error {
		ps.logger.Debug("Received pong from client")
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Start heartbeat goroutine
	pingTicker := time.NewTicker(30 * time.Second)
	defer pingTicker.Stop()

	// Channel to signal when connection should close
	done := make(chan bool)

	// Goroutine to handle heartbeat pings
	go func() {
		for {
			select {
			case <-pingTicker.C:
				ps.logger.Debug("Sending ping to client")
				conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					ps.logger.Error("Error sending ping: %v", err)
					done <- true
					return
				}
			case <-done:
				return
			}
		}
	}()

	// Main message loop
	for {
		select {
		case <-done:
			ps.logger.Debug("WebSocket connection closing due to heartbeat failure")
			return
		default:
			var msg WebSocketMessage
			err := conn.ReadJSON(&msg)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					ps.logger.Error("WebSocket error: %v", err)
				}
				done <- true
				return
			}

			// Reset read deadline on successful read
			conn.SetReadDeadline(time.Now().Add(60 * time.Second))

			ps.logger.Debug("Received WebSocket message: %s", msg.Type)

			// Handle ping messages from client
			if msg.Type == "ping" {
				ps.logger.Debug("Received ping from client, sending pong")
				pongMsg := WebSocketMessage{
					Type:      "pong",
					Timestamp: time.Now(),
				}
				conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err := conn.WriteJSON(pongMsg); err != nil {
					ps.logger.Error("Error sending pong: %v", err)
					done <- true
					return
				}
				continue
			}

			// Handle other message types
			switch msg.Type {
			case MessageTypeViewUpdate:
				ps.handleWebSocketViewUpdate(conn, msg)
			case MessageTypeLiveUpdate:
				ps.handleWebSocketLiveUpdate(conn, msg)
			case MessageTypeChatMessage:
				// Route chat messages through the bridge
				msgBytes, _ := json.Marshal(msg)
				ps.chatBridge.HandleWebSocketMessage(conn, msgBytes)
			case MessageTypeViewState:
				ps.handleWebSocketViewStateUpdate(conn, msg)
			default:
				ps.logger.Warning("Unknown message type: %s", msg.Type)
			}
		}
	}
}

func (ps *PreviewServer) sendCurrentViewState(conn *websocket.Conn) {
	msg := WebSocketMessage{
		Type:      MessageTypeViewState,
		Data:      ps.currentView,
		Timestamp: time.Now(),
	}

	conn.SetWriteDeadline(time.Now().Add(writeTimeout))
	if err := conn.WriteJSON(msg); err != nil {
		ps.logger.Error("Error sending current view state: %v", err)
	} else {
		ps.logger.Debug("Sent current view state to new client")
	}
}

func (ps *PreviewServer) handleWebSocketViewUpdate(conn *websocket.Conn, msg WebSocketMessage) {
	ps.logger.Debug("Processing view update via WebSocket")

	// Update current view state if the message contains view state data
	ps.updateViewStateFromMessage(msg)

	ps.broadcastToClients(msg)
}

func (ps *PreviewServer) handleWebSocketLiveUpdate(conn *websocket.Conn, msg WebSocketMessage) {
	log.Printf("⚡ Processing live update via WebSocket")

	// Update current view state if the message contains live update data
	ps.updateViewStateFromLiveUpdate(msg)

	ps.broadcastToClients(msg)
}

func (ps *PreviewServer) handleWebSocketViewStateUpdate(conn *websocket.Conn, msg WebSocketMessage) {
	log.Printf("🗂️ Processing view state update via WebSocket")

	// Update current view state from the message
	ps.updateViewStateFromMessage(msg)

	// Broadcast to all clients
	ps.broadcastToClients(msg)
}

func (ps *PreviewServer) handleWebSocketChatMessage(conn *websocket.Conn, msg WebSocketMessage) {
	log.Printf("💬 Processing chat message via WebSocket")

	// Extract message from data
	dataBytes, err := json.Marshal(msg.Data)
	if err != nil {
		log.Printf("❌ Error marshaling chat message data: %v", err)
		return
	}

	var chatData ChatMessageData
	if err := json.Unmarshal(dataBytes, &chatData); err != nil {
		log.Printf("❌ Error unmarshaling chat message: %v", err)
		return
	}

	// Process the chat message with current view context
	response := ps.processChatMessageWithContext(chatData.Message)

	// Send response back to client
	responseMsg := WebSocketMessage{
		Type: MessageTypeChatResponse,
		Data: ChatResponseData{
			Response:  response.Response,
			Action:    response.Action,
			Data:      response.Data,
			Timestamp: time.Now(),
		},
		Timestamp: time.Now(),
	}

	// Reset write deadline
	conn.SetWriteDeadline(time.Now().Add(writeTimeout))

	if err := conn.WriteJSON(responseMsg); err != nil {
		log.Printf("❌ Error sending chat response via WebSocket: %v", err)
	} else {
		log.Printf("✅ Chat response sent via WebSocket")
	}
}

func (ps *PreviewServer) updateViewStateFromMessage(msg WebSocketMessage) {
	// Try to extract view state from different message types
	dataBytes, err := json.Marshal(msg.Data)
	if err != nil {
		log.Printf("⚠️ Error marshaling message data for view state update: %v", err)
		return
	}

	// Try to parse as ViewStateUpdateRequest
	var viewStateUpdate ViewStateUpdateRequest
	if err := json.Unmarshal(dataBytes, &viewStateUpdate); err == nil {
		ps.currentView.Components = viewStateUpdate.Components
		ps.currentView.Theme = viewStateUpdate.Theme
		ps.currentView.Layout = viewStateUpdate.Layout
		ps.currentView.Canvas = viewStateUpdate.Canvas
		log.Printf("✅ Updated current view state from message")
		return
	}

	// Try to parse as ViewUpdateData
	var viewUpdate ViewUpdateData
	if err := json.Unmarshal(dataBytes, &viewUpdate); err == nil {
		// Update components if available
		if viewUpdate.Components != nil {
			components := []prompts.ComponentState{}
			for _, comp := range viewUpdate.Components {
				compBytes, _ := json.Marshal(comp)
				var componentState prompts.ComponentState
				if json.Unmarshal(compBytes, &componentState) == nil {
					components = append(components, componentState)
				}
			}
			ps.currentView.Components = components
		}

		// Update theme if available
		if viewUpdate.Theme != nil {
			themeBytes, _ := json.Marshal(viewUpdate.Theme)
			var themeState prompts.ThemeState
			if json.Unmarshal(themeBytes, &themeState) == nil {
				ps.currentView.Theme = themeState
			}
		}

		// Update layout if available
		if viewUpdate.Layout != nil {
			layoutBytes, _ := json.Marshal(viewUpdate.Layout)
			var layoutState prompts.LayoutState
			if json.Unmarshal(layoutBytes, &layoutState) == nil {
				ps.currentView.Layout = layoutState
			}
		}

		log.Printf("✅ Updated current view state from view update")
	}
}

func (ps *PreviewServer) updateViewStateFromLiveUpdate(msg WebSocketMessage) {
	// Try to extract component updates from live update messages
	dataBytes, err := json.Marshal(msg.Data)
	if err != nil {
		return
	}

	var liveUpdate LiveUpdateData
	if err := json.Unmarshal(dataBytes, &liveUpdate); err != nil {
		return
	}

	// Update specific component if it exists
	if liveUpdate.ComponentID != "" {
		for i, comp := range ps.currentView.Components {
			if comp.ID == liveUpdate.ComponentID {
				// Update component with changes
				changesBytes, _ := json.Marshal(liveUpdate.Changes)
				var changes map[string]interface{}
				if json.Unmarshal(changesBytes, &changes) == nil {
					// Apply changes to component
					if properties, ok := changes["properties"].(map[string]interface{}); ok {
						for key, value := range properties {
							ps.currentView.Components[i].Properties[key] = value
						}
					}
					if position, ok := changes["position"].(map[string]interface{}); ok {
						if x, ok := position["x"].(float64); ok {
							ps.currentView.Components[i].Position.X = int(x)
						}
						if y, ok := position["y"].(float64); ok {
							ps.currentView.Components[i].Position.Y = int(y)
						}
					}
					if size, ok := changes["size"].(map[string]interface{}); ok {
						if w, ok := size["w"].(float64); ok {
							ps.currentView.Components[i].Size.W = int(w)
						}
						if h, ok := size["h"].(float64); ok {
							ps.currentView.Components[i].Size.H = int(h)
						}
					}
				}
				break
			}
		}
	}
}

func (ps *PreviewServer) handleViewUpdate(w http.ResponseWriter, r *http.Request) {
	var data ViewUpdateData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	msg := WebSocketMessage{
		Type:      MessageTypeViewUpdate,
		Data:      data,
		Timestamp: time.Now(),
	}

	ps.updateViewStateFromMessage(msg)
	ps.broadcastToClients(msg)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (ps *PreviewServer) handleLiveUpdate(w http.ResponseWriter, r *http.Request) {
	var data LiveUpdateData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	msg := WebSocketMessage{
		Type:      MessageTypeLiveUpdate,
		Data:      data,
		Timestamp: time.Now(),
	}

	ps.updateViewStateFromLiveUpdate(msg)
	ps.broadcastToClients(msg)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (ps *PreviewServer) handleViewStateUpdate(w http.ResponseWriter, r *http.Request) {
	var data ViewStateUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update current view state
	ps.currentView.Components = data.Components
	ps.currentView.Theme = data.Theme
	ps.currentView.Layout = data.Layout
	ps.currentView.Canvas = data.Canvas

	msg := WebSocketMessage{
		Type:      MessageTypeViewState,
		Data:      ps.currentView,
		Timestamp: time.Now(),
	}

	ps.broadcastToClients(msg)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (ps *PreviewServer) handleGetViewState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ps.currentView)
}

func (ps *PreviewServer) handleChatMessage(w http.ResponseWriter, r *http.Request) {
	var data ChatMessageData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := ps.processChatMessageWithContext(data.Message)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (ps *PreviewServer) handleStatus(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"status":            "ok",
		"server":            "vibe-preview",
		"version":           "1.0.0",
		"connected_clients": len(ps.clients),
		"timestamp":         time.Now(),
		"prompts_loaded":    ps.promptLoader != nil,
		"components_count":  len(ps.currentView.Components),
		"current_theme":     ps.currentView.Theme.Name,
		"current_viewport":  ps.currentView.Canvas.Viewport,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (ps *PreviewServer) processChatMessageWithContext(message string) ChatResponseData {
	ps.logger.Debug("Processing chat message with Claude AI: %s", message)

	// Use Claude AI for all responses
	if ps.claudeClient != nil && ps.promptLoader != nil {
		// Determine mode based on componentMode setting
		mode := "general"
		if ps.componentMode {
			mode = "component"
		}

		promptData := prompts.PromptData{
			ProjectContext:      "VibeCode UI Editor - Real-time Component Builder",
			Templates:           make(map[string]string),
			CurrentView:         ps.currentView,
			UserInput:           message,
			ConversationHistory: []prompts.ConversationMessage{},
			Mode:                mode, // Pass the mode to the prompt
		}

		// Build contextual prompt
		fullPrompt, err := ps.promptLoader.BuildChatPrompt(promptData)
		if err != nil {
			ps.logger.Warning("Error building contextual prompt: %v", err)
			return ps.fallbackResponse(message)
		}

		ps.logger.Debug("Built contextual prompt with current view state (mode: %s)", mode)

		// Send to Claude AI
		messages := []ClaudeMessage{{
			Role:    "user",
			Content: fullPrompt,
		}}

		claudeResponse, err := ps.claudeClient.CreateMessage(messages)
		if err != nil {
			ps.logger.Error("Error getting Claude response: %v", err)
			return ps.fallbackResponse(message)
		}

		ps.logger.Debug("Got Claude response: %s", claudeResponse)

		// Process Claude response for UI updates
		response := ps.processClaudeResponse(claudeResponse)

		return ChatResponseData{
			Response:  response.Response,
			Action:    response.Action,
			Data:      response.Data,
			Timestamp: time.Now(),
		}
	}

	// Fallback if Claude is not available
	return ps.fallbackResponse(message)
}

func (ps *PreviewServer) fallbackResponse(message string) ChatResponseData {
	// Enhanced AI responses that consider current state
	lowerMessage := strings.ToLower(message)

	// Context-aware responses
	switch {
	case strings.Contains(lowerMessage, "hola") || strings.Contains(lowerMessage, "hello"):
		contextInfo := ps.getContextualGreeting()
		return ChatResponseData{
			Response:  fmt.Sprintf("¡Hola! Soy Viber, tu asistente de VibeCode. %s Puedo ayudarte a crear y modificar componentes UI en tiempo real.", contextInfo),
			Timestamp: time.Now(),
		}

	case strings.Contains(lowerMessage, "estado") || strings.Contains(lowerMessage, "status") || strings.Contains(lowerMessage, "qué hay"):
		return ps.getCanvasStatusResponse()

	case strings.Contains(lowerMessage, "ayuda") || strings.Contains(lowerMessage, "help"):
		return ps.getContextualHelpResponse()

	// Component creation commands with context-aware positioning
	case strings.Contains(lowerMessage, "agregar botón") || strings.Contains(lowerMessage, "add button"):
		return ps.createContextualComponentResponse("button", "atom", "Input", map[string]interface{}{
			"text":    "Nuevo Botón",
			"variant": "primary",
			"size":    "medium",
		}, 160, 40, "Agregué un botón interactivo")

	case strings.Contains(lowerMessage, "agregar texto") && !strings.Contains(lowerMessage, "animado"):
		return ps.createContextualComponentResponse("text", "atom", "Typography", map[string]interface{}{
			"content": "Texto de ejemplo",
			"size":    "medium",
			"weight":  "normal",
		}, 320, 32, "Agregué un texto simple")

	case strings.Contains(lowerMessage, "texto animado") || strings.Contains(lowerMessage, "animated text"):
		return ps.createContextualComponentResponse("animated-text", "atom", "Typography", map[string]interface{}{
			"text":      "Texto Animado",
			"effect":    "rotate3D",
			"className": "text-white text-2xl font-bold",
			"delay":     0,
		}, 400, 48, "Agregué texto con animación 3D")

	case strings.Contains(lowerMessage, "agregar imagen") || strings.Contains(lowerMessage, "add image"):
		return ps.createContextualComponentResponse("image", "atom", "Media", map[string]interface{}{
			"src":     "https://images.pexels.com/photos/3861969/pexels-photo-3861969.jpeg?auto=compress&cs=tinysrgb&w=400",
			"alt":     "Imagen de ejemplo",
			"rounded": true,
		}, 200, 200, "Agregué una imagen redondeada")

	case strings.Contains(lowerMessage, "agregar input") || strings.Contains(lowerMessage, "add input"):
		return ps.createContextualComponentResponse("input", "atom", "Input", map[string]interface{}{
			"placeholder": "Escribe aquí...",
			"type":        "text",
			"label":       "Campo de entrada",
		}, 280, 56, "Agregué un campo de entrada")

	case strings.Contains(lowerMessage, "agregar tarjeta") || strings.Contains(lowerMessage, "add card"):
		return ps.createContextualComponentResponse("card", "molecule", "Layout", map[string]interface{}{
			"title":    "Tarjeta de Ejemplo",
			"content":  "Esta es una tarjeta con imagen y descripción",
			"hasImage": true,
			"imageUrl": "https://images.pexels.com/photos/3861969/pexels-photo-3861969.jpeg?auto=compress&cs=tinysrgb&w=400",
		}, 320, 400, "Creé una tarjeta con imagen")

	case strings.Contains(lowerMessage, "agregar formulario") || strings.Contains(lowerMessage, "add form"):
		return ps.createContextualComponentResponse("form", "molecule", "Input", map[string]interface{}{
			"title":      "Formulario de Contacto",
			"fields":     []string{"nombre", "email", "mensaje"},
			"submitText": "Enviar",
		}, 400, 360, "Creé un formulario de contacto")

	case strings.Contains(lowerMessage, "agregar navegación") || strings.Contains(lowerMessage, "add navigation"):
		return ps.createContextualComponentResponse("navigation", "molecule", "Layout", map[string]interface{}{
			"items": []string{"Inicio", "Acerca de", "Servicios", "Contacto"},
			"style": "horizontal",
		}, 800, 60, "Agregué una barra de navegación")

	case strings.Contains(lowerMessage, "agregar hero") || strings.Contains(lowerMessage, "add hero"):
		return ps.createContextualComponentResponse("hero", "organism", "Layout", map[string]interface{}{
			"title":           "Bienvenido a VibeCode",
			"subtitle":        "Construye interfaces increíbles con IA",
			"ctaText":         "Comenzar",
			"backgroundImage": "https://images.pexels.com/photos/3861969/pexels-photo-3861969.jpeg?auto=compress&cs=tinysrgb&w=1200",
		}, 900, 480, "Creé una sección hero completa")

	case strings.Contains(lowerMessage, "agregar galería") || strings.Contains(lowerMessage, "add gallery"):
		return ps.createContextualComponentResponse("gallery", "organism", "Media", map[string]interface{}{
			"title": "Galería de Fotos",
			"images": []string{
				"https://images.pexels.com/photos/3861969/pexels-photo-3861969.jpeg?auto=compress&cs=tinysrgb&w=400",
				"https://images.pexels.com/photos/2662116/pexels-photo-2662116.jpeg?auto=compress&cs=tinysrgb&w=400",
				"https://images.pexels.com/photos/1181671/pexels-photo-1181671.jpeg?auto=compress&cs=tinysrgb&w=400",
			},
			"columns": 3,
		}, 640, 400, "Creé una galería de 3 columnas")

	// Theme commands
	case strings.Contains(lowerMessage, "azul") || strings.Contains(lowerMessage, "blue"):
		return ps.createThemeResponse(map[string]interface{}{
			"primary":   "#3B82F6",
			"secondary": "#06B6D4",
			"accent":    "#1D4ED8",
		}, "Cambié el tema a azul")

	case strings.Contains(lowerMessage, "rojo") || strings.Contains(lowerMessage, "red"):
		return ps.createThemeResponse(map[string]interface{}{
			"primary":   "#EF4444",
			"secondary": "#F87171",
			"accent":    "#DC2626",
		}, "Cambié el tema a rojo")

	case strings.Contains(lowerMessage, "verde") || strings.Contains(lowerMessage, "green"):
		return ps.createThemeResponse(map[string]interface{}{
			"primary":   "#10B981",
			"secondary": "#34D399",
			"accent":    "#059669",
		}, "Cambié el tema a verde")

	case strings.Contains(lowerMessage, "océano") || strings.Contains(lowerMessage, "ocean"):
		return ps.createThemeResponse(map[string]interface{}{
			"primary":    "#0EA5E9",
			"secondary":  "#06B6D4",
			"accent":     "#10B981",
			"background": "#0C1A2B",
			"surface":    "#1E3A5F",
			"text":       "#F0F9FF",
		}, "Apliqué el tema océano")

	case strings.Contains(lowerMessage, "puesta de sol") || strings.Contains(lowerMessage, "sunset"):
		return ps.createThemeResponse(map[string]interface{}{
			"primary":    "#F59E0B",
			"secondary":  "#EF4444",
			"accent":     "#EC4899",
			"background": "#1A0B0B",
			"surface":    "#2D1B1B",
			"text":       "#FFF7ED",
		}, "Apliqué el tema puesta de sol")

	default:
		return ChatResponseData{
			Response: fmt.Sprintf(`No entendí exactamente qué quieres hacer con "%s". 

**Estado actual del canvas:**
- %d componentes presentes
- Tema: %s
- Vista: %s

Prueba comandos como:
• "agregar botón" - para crear componentes
• "cambiar a azul" - para cambiar colores
• "estado" - para ver información del canvas
• "ayuda" - para ver todos los comandos

¿Qué te gustaría crear o modificar?`, message, len(ps.currentView.Components), ps.currentView.Theme.Name, ps.currentView.Canvas.Viewport),
			Timestamp: time.Now(),
		}
	}
}

// Helper methods for contextual responses
func (ps *PreviewServer) getContextualGreeting() string {
	if len(ps.currentView.Components) == 0 {
		return "Veo que tu canvas está vacío y listo para crear algo increíble."
	}
	return fmt.Sprintf("Veo que tienes %d componentes en tu canvas con el tema %s.",
		len(ps.currentView.Components), ps.currentView.Theme.Name)
}

func (ps *PreviewServer) getCanvasStatusResponse() ChatResponseData {
	var analysis strings.Builder
	analysis.WriteString("📊 **Estado actual del canvas:**\n\n")
	analysis.WriteString(fmt.Sprintf("• **Componentes**: %d total\n", len(ps.currentView.Components)))
	analysis.WriteString(fmt.Sprintf("• **Tema**: %s\n", ps.currentView.Theme.Name))
	analysis.WriteString(fmt.Sprintf("• **Vista**: %s\n", ps.currentView.Canvas.Viewport))
	analysis.WriteString(fmt.Sprintf("• **Zoom**: %.0f%%\n", ps.currentView.Canvas.Zoom*100))
	analysis.WriteString(fmt.Sprintf("• **Grid**: %d columnas\n", ps.currentView.Layout.Grid))

	if ps.currentView.Canvas.SelectedItem != "" {
		analysis.WriteString(fmt.Sprintf("• **Seleccionado**: %s\n", ps.currentView.Canvas.SelectedItem))
	}

	if len(ps.currentView.Components) > 0 {
		componentsByType := make(map[string]int)
		for _, comp := range ps.currentView.Components {
			componentsByType[comp.Type]++
		}

		analysis.WriteString("\n**Componentes por tipo:**\n")
		for compType, count := range componentsByType {
			analysis.WriteString(fmt.Sprintf("• %s: %d\n", strings.Title(compType), count))
		}
	} else {
		analysis.WriteString("\n*El canvas está vacío, listo para crear componentes.*\n")
	}

	return ChatResponseData{
		Response:  analysis.String(),
		Timestamp: time.Now(),
	}
}

func (ps *PreviewServer) getContextualHelpResponse() ChatResponseData {
	helpText := `🎨 **Comandos disponibles:**

**Componentes Básicos (Atoms):**
• "agregar botón" - Crea un botón interactivo
• "agregar texto" - Añade texto simple
• "agregar texto animado" - Texto con efectos
• "agregar imagen" - Inserta una imagen
• "agregar input" - Campo de entrada

**Componentes Compuestos (Molecules):**
• "agregar tarjeta" - Crea una tarjeta
• "agregar formulario" - Formulario de contacto
• "agregar navegación" - Barra de navegación

**Componentes Complejos (Organisms):**
• "agregar hero" - Sección hero completa
• "agregar galería" - Galería de imágenes

**Temas:**
• "cambiar a azul/rojo/verde"
• "tema oscuro/claro"
• "océano/puesta de sol"

**Información:**
• "estado" - Ver información del canvas actual
• "qué hay" - Analizar componentes presentes

**Posiciones inteligentes disponibles:**`

	// Add intelligent positioning based on current state
	positions := ps.findAvailablePositions()
	if len(positions) > 0 {
		helpText += "\n• Las nuevas componentes se colocarán automáticamente en posiciones libres"
	}

	helpText += fmt.Sprintf("\n\n**Tu canvas actual**: %d componentes, tema %s, vista %s",
		len(ps.currentView.Components), ps.currentView.Theme.Name, ps.currentView.Canvas.Viewport)

	return ChatResponseData{
		Response:  helpText,
		Timestamp: time.Now(),
	}
}

func (ps *PreviewServer) findAvailablePositions() []prompts.Position {
	// Find available positions based on current components
	usedPositions := make(map[string]bool)
	for _, comp := range ps.currentView.Components {
		key := fmt.Sprintf("%d,%d", comp.Position.X, comp.Position.Y)
		usedPositions[key] = true
	}

	// Suggest good positions
	goodPositions := []prompts.Position{
		{X: 100, Y: 100}, {X: 300, Y: 100}, {X: 500, Y: 100},
		{X: 100, Y: 300}, {X: 300, Y: 300}, {X: 500, Y: 300},
		{X: 100, Y: 500}, {X: 300, Y: 500}, {X: 500, Y: 500},
	}

	availablePositions := []prompts.Position{}
	for _, pos := range goodPositions {
		key := fmt.Sprintf("%d,%d", pos.X, pos.Y)
		if !usedPositions[key] {
			availablePositions = append(availablePositions, pos)
		}
	}

	return availablePositions
}

func (ps *PreviewServer) createContextualComponentResponse(id, componentType, category string, properties map[string]interface{}, width, height int, explanation string) ChatResponseData {
	// Find best available position
	availablePositions := ps.findAvailablePositions()
	position := prompts.Position{X: 200, Y: 200} // default

	if len(availablePositions) > 0 {
		position = availablePositions[0]
	}

	uiUpdate := UIUpdateRequest{
		Type:   "ui_update",
		Action: "add_component",
		Data: map[string]interface{}{
			"id":         id,
			"type":       componentType,
			"name":       fmt.Sprintf("%s_%d", strings.Title(id), time.Now().Unix()),
			"category":   category,
			"properties": properties,
			"position":   map[string]int{"x": position.X, "y": position.Y},
			"size":       map[string]int{"w": width, "h": height},
		},
		Explanation: explanation,
	}

	// Update current view state with new component
	newComponent := prompts.ComponentState{
		ID:         fmt.Sprintf("%s_%d", id, time.Now().Unix()),
		Type:       componentType,
		Name:       fmt.Sprintf("%s_%d", strings.Title(id), time.Now().Unix()),
		Category:   category,
		Properties: properties,
		Position:   position,
		Size:       prompts.Size{W: width, H: height},
	}
	ps.currentView.Components = append(ps.currentView.Components, newComponent)

	// Validate the component
	if jsonBytes, err := json.Marshal(uiUpdate); err == nil {
		if ps.promptLoader != nil {
			if err := prompts.ValidateUIUpdateJSON(string(jsonBytes)); err != nil {
				log.Printf("⚠️ Component validation failed: %v", err)
			} else {
				log.Printf("✅ Component validation passed")
			}
		}
	}

	// Broadcast the update
	ps.BroadcastChatResponse(uiUpdate)

	contextualMessage := fmt.Sprintf("%s en posición (%d, %d). Ahora tienes %d componentes en el canvas.",
		explanation, position.X, position.Y, len(ps.currentView.Components))

	return ChatResponseData{
		Response:  contextualMessage + " ✨",
		Action:    "add_component",
		Data:      uiUpdate.Data,
		Timestamp: time.Now(),
	}
}

func (ps *PreviewServer) createThemeResponse(colors map[string]interface{}, explanation string) ChatResponseData {
	uiUpdate := UIUpdateRequest{
		Type:   "ui_update",
		Action: "update_theme",
		Data: map[string]interface{}{
			"colors": colors,
		},
		Explanation: explanation,
	}

	// Update current view state theme
	for key, value := range colors {
		if strValue, ok := value.(string); ok {
			ps.currentView.Theme.Colors[key] = strValue
		}
	}

	// Broadcast the update
	ps.BroadcastChatResponse(uiUpdate)

	return ChatResponseData{
		Response:  explanation + " 🎨",
		Action:    "update_theme",
		Data:      uiUpdate.Data,
		Timestamp: time.Now(),
	}
}

func (ps *PreviewServer) BroadcastChatResponse(uiUpdate UIUpdateRequest) {
	// Frontend expects: { action: string, data: unknown }
	// So we create a simple structure that matches this expectation
	liveUpdateData := map[string]interface{}{
		"action": uiUpdate.Action,
		"data":   uiUpdate.Data,
	}

	msg := WebSocketMessage{
		Type:      MessageTypeLiveUpdate,
		Data:      liveUpdateData,
		Timestamp: time.Now(),
	}

	ps.broadcastToClients(msg)
}

func (ps *PreviewServer) broadcastToClients(msg WebSocketMessage) {
	for client := range ps.clients {
		// Set write deadline for each client
		client.SetWriteDeadline(time.Now().Add(writeTimeout))

		err := client.WriteJSON(msg)
		if err != nil {
			ps.logger.Error("Error sending message to client: %v", err)
			client.Close()
			delete(ps.clients, client)
		}
	}

	ps.logger.Debug("Broadcasted message to %d clients", len(ps.clients))
}

// broadcastViewStateUpdate sends the complete current view state to all connected clients
func (ps *PreviewServer) broadcastViewStateUpdate() {
	log.Printf("📡 Broadcasting complete view state update to all clients")

	msg := WebSocketMessage{
		Type:      MessageTypeViewState,
		Data:      ps.currentView,
		Timestamp: time.Now(),
	}

	ps.broadcastToClients(msg)
}

// processClaudeResponse processes Claude AI response for UI updates and actions
func (ps *PreviewServer) processClaudeResponse(claudeResponse string) ChatResponseData {
	log.Printf("🔄 Processing Claude response for UI updates")

	// Check for multiple JSONs first
	multipleJSONs := ps.detectMultipleJSONs(claudeResponse)
	if multipleJSONs > 1 {
		log.Printf("⚠️ Detected %d JSON objects in response, processing only the first one", multipleJSONs)
	}

	// Extract JSON from Claude response if present
	jsonStr, hasJSON := prompts.ExtractJSONFromResponse(claudeResponse)

	if hasJSON {
		log.Printf("📄 Found JSON in Claude response: %s", jsonStr)

		// Validate the JSON structure
		if err := prompts.ValidateUIUpdateJSON(jsonStr); err != nil {
			log.Printf("⚠️ Invalid UI update JSON: %v", err)
			// Return full response if JSON is invalid
			return ChatResponseData{
				Response:  claudeResponse,
				Timestamp: time.Now(),
			}
		}

		log.Printf("✅ Valid UI update JSON found")

		// Parse the update
		var update map[string]interface{}
		if err := json.Unmarshal([]byte(jsonStr), &update); err != nil {
			log.Printf("⚠️ Failed to parse UI update: %v", err)
			return ChatResponseData{
				Response:  claudeResponse,
				Timestamp: time.Now(),
			}
		}

		// Process the UI update
		ps.handleUIUpdateFromClaude(update)

		// Extract conversational text (everything before the JSON)
		conversationalText := ps.extractConversationalText(claudeResponse, jsonStr)

		// Add warning if multiple JSONs were detected
		if multipleJSONs > 1 {
			conversationalText += fmt.Sprintf("\n\n⚠️ Nota: Solo procesé el primer componente. Para agregar más componentes, pídeme cada uno por separado.")
		}

		// Return response with action data but only conversational text
		if action, ok := update["action"].(string); ok {
			return ChatResponseData{
				Response:  conversationalText,
				Action:    action,
				Data:      update["data"],
				Timestamp: time.Now(),
			}
		}

		// Return conversational text even if no action
		return ChatResponseData{
			Response:  conversationalText,
			Timestamp: time.Now(),
		}
	}

	// Return simple text response if no UI updates
	return ChatResponseData{
		Response:  claudeResponse,
		Timestamp: time.Now(),
	}
}

// detectMultipleJSONs counts how many JSON objects are in the response
func (ps *PreviewServer) detectMultipleJSONs(response string) int {
	count := 0
	inString := false
	escaped := false
	braceDepth := 0

	for i, char := range response {
		if escaped {
			escaped = false
			continue
		}

		if char == '\\' {
			escaped = true
			continue
		}

		if char == '"' {
			inString = !inString
			continue
		}

		if inString {
			continue
		}

		if char == '{' {
			if braceDepth == 0 {
				// Start of a new JSON object
			}
			braceDepth++
		} else if char == '}' {
			braceDepth--
			if braceDepth == 0 {
				// End of a JSON object
				count++
				// Look ahead to see if there's another JSON
				remaining := response[i+1:]
				if strings.TrimSpace(remaining) != "" && strings.Contains(remaining, "{") {
					// There might be another JSON, continue counting
				}
			}
		}
	}

	return count
}

// extractConversationalText extracts the conversational part of the response, removing JSON
func (ps *PreviewServer) extractConversationalText(fullResponse, jsonStr string) string {
	// Remove JSON from the response

	// First, try to find the JSON and remove it
	jsonStart := strings.Index(fullResponse, jsonStr)
	if jsonStart != -1 {
		// Take everything before the JSON
		textBefore := fullResponse[:jsonStart]

		// Also check if there's text after the JSON (though this should be rare)
		jsonEnd := jsonStart + len(jsonStr)
		textAfter := ""
		if jsonEnd < len(fullResponse) {
			textAfter = fullResponse[jsonEnd:]
		}

		conversationalText := strings.TrimSpace(textBefore + textAfter)

		// Remove any code block markers that might be left
		conversationalText = strings.Replace(conversationalText, "```json", "", -1)
		conversationalText = strings.Replace(conversationalText, "```", "", -1)

		// Clean up extra whitespace
		conversationalText = strings.TrimSpace(conversationalText)

		// If we have meaningful text, return it
		if len(conversationalText) > 0 {
			return conversationalText
		}
	}

	// Fallback: if we can't extract conversational text, provide a default based on the action
	// This should rarely happen with well-formed Claude responses
	return "✨ Actualización aplicada al canvas"
}

// handleUIUpdateFromClaude processes UI updates from Claude response
func (ps *PreviewServer) handleUIUpdateFromClaude(update map[string]interface{}) {
	log.Printf("🎨 Processing UI update from Claude: %v", update)

	if action, ok := update["action"].(string); ok {
		switch action {
		case "add_component":
			if data, ok := update["data"].(map[string]interface{}); ok {
				ps.addComponentFromClaude(data)
			}
		case "update_component":
			if data, ok := update["data"].(map[string]interface{}); ok {
				ps.updateComponentFromClaude(data)
			}
		case "update_theme":
			if data, ok := update["data"].(map[string]interface{}); ok {
				ps.updateThemeFromClaude(data)
			}
		case "update_layout":
			if data, ok := update["data"].(map[string]interface{}); ok {
				ps.updateLayoutFromClaude(data)
			}
		default:
			log.Printf("⚠️ Unknown UI update action: %s", action)
		}
	}

	// Broadcast the update to all clients
	uiUpdateRequest := UIUpdateRequest{
		Type:        "ui_update",
		Action:      update["action"].(string),
		Data:        update["data"],
		Explanation: fmt.Sprintf("Viber AI: %v", update["explanation"]),
	}

	ps.BroadcastChatResponse(uiUpdateRequest)

	// Also send the complete view state update after processing
	ps.broadcastViewStateUpdate()
}

// addComponentFromClaude adds a component from Claude AI response
func (ps *PreviewServer) addComponentFromClaude(data map[string]interface{}) {
	log.Printf("➕ Adding component from Claude: %v", data)

	// Find available position
	availablePositions := ps.findAvailablePositions()
	position := prompts.Position{X: 200, Y: 200} // default

	if len(availablePositions) > 0 {
		position = availablePositions[0]
	}

	// Extract and validate required fields
	componentType := data["type"].(string)
	componentID := fmt.Sprintf("%s_%d", componentType, time.Now().Unix())
	componentName := fmt.Sprintf("%s Component", strings.Title(componentType))

	// Create component state
	componentState := prompts.ComponentState{
		ID:         componentID,
		Type:       componentType,
		Name:       componentName,
		Category:   data["category"].(string),
		Properties: data["properties"].(map[string]interface{}),
		Position:   position,
		Size: prompts.Size{
			W: int(data["size"].(map[string]interface{})["w"].(float64)),
			H: int(data["size"].(map[string]interface{})["h"].(float64)),
		},
	}

	// Add to current view
	ps.currentView.Components = append(ps.currentView.Components, componentState)

	// Update the data being sent to frontend to include ID and name
	data["id"] = componentID
	data["name"] = componentName

	log.Printf("✅ Component added: %s (%s) at position (%d, %d)", componentState.ID, componentState.Name, position.X, position.Y)
}

// updateComponentFromClaude updates a component from Claude AI response
func (ps *PreviewServer) updateComponentFromClaude(data map[string]interface{}) {
	log.Printf("🔄 Updating component from Claude: %v", data)

	componentID := data["id"].(string)

	// Find and update component
	for i, comp := range ps.currentView.Components {
		if comp.ID == componentID {
			if properties, ok := data["properties"].(map[string]interface{}); ok {
				ps.currentView.Components[i].Properties = properties
			}
			log.Printf("✅ Component updated: %s", componentID)
			return
		}
	}

	log.Printf("⚠️ Component not found for update: %s", componentID)
}

// updateThemeFromClaude updates theme from Claude AI response
func (ps *PreviewServer) updateThemeFromClaude(data map[string]interface{}) {
	log.Printf("🎨 Updating theme from Claude: %v", data)

	if colors, ok := data["colors"].(map[string]interface{}); ok {
		for key, value := range colors {
			if strValue, ok := value.(string); ok {
				ps.currentView.Theme.Colors[key] = strValue
			}
		}
		log.Printf("✅ Theme updated with new colors")
	}
}

// updateLayoutFromClaude updates layout from Claude AI response
func (ps *PreviewServer) updateLayoutFromClaude(data map[string]interface{}) {
	log.Printf("📐 Updating layout from Claude: %v", data)

	if grid, ok := data["grid"].(float64); ok {
		ps.currentView.Layout.Grid = int(grid)
	}
	if rowHeight, ok := data["row_height"].(float64); ok {
		ps.currentView.Layout.RowHeight = int(rowHeight)
	}

	log.Printf("✅ Layout updated")
}

// SetComponentMode sets the component mode for the preview server
func (ps *PreviewServer) SetComponentMode(enabled bool) {
	ps.componentMode = enabled
	ps.logger.Debug("Component mode set to: %t", enabled)
}

// GetChatBridge returns the chat bridge for external access
func (ps *PreviewServer) GetChatBridge() *ChatBridge {
	return ps.chatBridge
}

// handlePreviewPage serves the preview HTML page
func (ps *PreviewServer) handlePreviewPage(w http.ResponseWriter, r *http.Request) {
	ps.logger.Debug("Serving preview page request: %s %s", r.Method, r.URL.Path)

	modeTitle := "VibeCode Preview"
	if ps.componentMode {
		modeTitle = "VibeCode Component Builder"
	}

	html := `<!DOCTYPE html>
<html>
<head>
    <title>` + modeTitle + `</title>
    <style>
        body { 
            font-family: Arial, sans-serif; 
            background: #0F0F0F; 
            color: white; 
            padding: 20px; 
            margin: 0;
            display: flex;
            flex-direction: column;
            height: 100vh;
        }
        .main-content {
            display: flex;
            flex: 1;
            gap: 20px;
            min-height: 0;
        }
        .header { 
            background: #8B5CF6; 
            padding: 20px; 
            border-radius: 8px; 
            text-align: center; 
            margin-bottom: 20px; 
        }
        .canvas { 
            border: 2px solid #8B5CF6; 
            border-radius: 8px; 
            padding: 20px; 
            background: #1A1A1A; 
            flex: 2;
            overflow-y: auto;
        }
        .chat-panel {
            border: 2px solid #06B6D4;
            border-radius: 8px;
            background: #1A1A1A;
            flex: 1;
            display: flex;
            flex-direction: column;
            min-width: 300px;
            max-width: 400px;
        }
        .chat-header {
            background: #06B6D4;
            color: white;
            padding: 15px;
            border-radius: 6px 6px 0 0;
            text-align: center;
            font-weight: bold;
        }
        .chat-messages {
            flex: 1;
            padding: 15px;
            overflow-y: auto;
            min-height: 200px;
            max-height: 400px;
        }
        .chat-input-container {
            padding: 15px;
            border-top: 1px solid #333;
        }
        .chat-input {
            width: 100%;
            padding: 10px;
            border: 1px solid #333;
            border-radius: 4px;
            background: #2A2A2A;
            color: white;
            font-size: 14px;
        }
        .chat-send {
            width: 100%;
            padding: 10px;
            margin-top: 10px;
            background: #06B6D4;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-size: 14px;
        }
        .chat-send:hover {
            background: #0891B2;
        }
        .chat-send:disabled {
            background: #333;
            cursor: not-allowed;
        }
        .message {
            margin-bottom: 15px;
            padding: 10px;
            border-radius: 6px;
            line-height: 1.4;
        }
        .message.user {
            background: #8B5CF6;
            margin-left: 20px;
            text-align: right;
        }
        .message.assistant {
            background: #333;
            margin-right: 20px;
        }
        .message.system {
            background: #F59E0B;
            text-align: center;
            font-style: italic;
        }
        .message-time {
            font-size: 11px;
            opacity: 0.7;
            margin-top: 5px;
        }
        .component { 
            border: 1px solid #06B6D4; 
            background: rgba(6, 182, 212, 0.1); 
            padding: 10px; 
            margin: 10px; 
            border-radius: 4px; 
        }
        .status { 
            position: fixed; 
            top: 10px; 
            right: 10px; 
            background: #10B981; 
            padding: 10px; 
            border-radius: 20px; 
            font-size: 12px; 
            z-index: 1000;
        }
    </style>
</head>
<body>
    <div class="status" id="status">🔌 Conectando...</div>
    <div class="header">
        <h1>⚡ ` + modeTitle + `</h1>
        <p>Preview Server - Puerto 3001</p>
    </div>
    <div class="main-content">
        <div class="canvas" id="canvas">
            <h2>🎨 Canvas de Componentes</h2>
            <div id="components">
                <p>Canvas vacío. Usa el terminal para agregar componentes.</p>
            </div>
        </div>
        <div class="chat-panel">
            <div class="chat-header">
                💬 Chat con Viber
            </div>
            <div class="chat-messages" id="chatMessages">
                <div class="message system">
                    <div>¡Hola! Soy Viber, tu asistente de VibeCode. Puedo ayudarte a crear componentes y modificar el diseño.</div>
                    <div class="message-time">Conectado</div>
                </div>
            </div>
            <div class="chat-input-container">
                <input type="text" class="chat-input" id="chatInput" placeholder="Escribe tu mensaje aquí..." autocomplete="off">
                <button class="chat-send" id="chatSend">Enviar</button>
            </div>
        </div>
    </div>
    <script>
        let ws;
        function connectWS() {
            ws = new WebSocket('ws://localhost:3001/ws');
            ws.onopen = function() {
                console.log('WebSocket conectado');
                document.getElementById('status').textContent = '🔌 Conectado';
                document.getElementById('status').style.background = '#10B981';
            };
            ws.onmessage = function(event) {
                console.log('Mensaje recibido:', event.data);
                try {
                    const data = JSON.parse(event.data);
                    if (data.type === 'view_state_update') {
                        updateComponents(data.data);
                    } else if (data.type === 'chat_response') {
                        handleChatResponse(data.data);
                    }
                } catch (e) {
                    console.error('Error:', e);
                }
            };
            ws.onclose = function() {
                console.log('WebSocket desconectado');
                document.getElementById('status').textContent = '🔌 Desconectado';
                document.getElementById('status').style.background = '#EF4444';
                setTimeout(connectWS, 3000);
            };
            ws.onerror = function(error) {
                console.error('WebSocket error:', error);
            };
        }
        
        function updateComponents(viewState) {
            const container = document.getElementById('components');
            if (!viewState || !viewState.components || viewState.components.length === 0) {
                container.innerHTML = '<p>Canvas vacío. Usa el terminal para agregar componentes.</p>';
                return;
            }
            
            container.innerHTML = '';
            viewState.components.forEach(component => {
                const div = document.createElement('div');
                div.className = 'component';
                div.innerHTML = '<strong>' + component.type + '</strong><br>' + (component.name || component.id);
                container.appendChild(div);
            });
        }
        
        // Chat functions
        function sendChatMessage() {
            const input = document.getElementById('chatInput');
            const message = input.value.trim();
            if (!message || !ws) return;
            
            // Add user message to chat
            addMessageToChat('user', message);
            
            // Send to WebSocket
            const chatData = {
                type: 'chat_message',
                data: {
                    message: message
                },
                timestamp: new Date().toISOString()
            };
            
            ws.send(JSON.stringify(chatData));
            input.value = '';
            
            // Disable send button temporarily
            const sendBtn = document.getElementById('chatSend');
            sendBtn.disabled = true;
            sendBtn.textContent = 'Enviando...';
        }
        
        function handleChatResponse(responseData) {
            const sendBtn = document.getElementById('chatSend');
            sendBtn.disabled = false;
            sendBtn.textContent = 'Enviar';
            
            // Handle different response structures
            let content = '';
            if (responseData.assistant_response) {
                content = responseData.assistant_response.content || responseData.assistant_response.response;
            } else if (responseData.response) {
                content = responseData.response;
            } else if (responseData.content) {
                content = responseData.content;
            }
            
            if (content) {
                addMessageToChat('assistant', content);
            }
        }
        
        function addMessageToChat(type, content) {
            const chatMessages = document.getElementById('chatMessages');
            const messageDiv = document.createElement('div');
            messageDiv.className = 'message ' + type;
            
            const contentDiv = document.createElement('div');
            contentDiv.textContent = content;
            
            const timeDiv = document.createElement('div');
            timeDiv.className = 'message-time';
            timeDiv.textContent = new Date().toLocaleTimeString();
            
            messageDiv.appendChild(contentDiv);
            messageDiv.appendChild(timeDiv);
            chatMessages.appendChild(messageDiv);
            
            // Scroll to bottom
            chatMessages.scrollTop = chatMessages.scrollHeight;
        }
        
        window.onload = function() {
            connectWS();
            
            // Setup chat event listeners
            const chatInput = document.getElementById('chatInput');
            const chatSend = document.getElementById('chatSend');
            
            chatSend.addEventListener('click', sendChatMessage);
            chatInput.addEventListener('keypress', function(e) {
                if (e.key === 'Enter') {
                    e.preventDefault();
                    sendChatMessage();
                }
            });
        };
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
