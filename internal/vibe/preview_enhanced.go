package vibe

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/vibercode/cli/internal/vibe/prompts"
)

// EnhancedPreviewServer is an enhanced version of PreviewServer with vector/graph storage
type EnhancedPreviewServer struct {
	port               string
	clients            map[*websocket.Conn]bool
	upgrader           websocket.Upgrader
	promptLoader       *prompts.PromptLoader
	currentView        *prompts.CurrentViewState
	claudeClient       *ClaudeClient
	vectorGraphService VectorGraphServiceInterface
	projectID          string
	sessionID          string
}

// NewEnhancedPreviewServer creates a new enhanced preview server with vector/graph storage
func NewEnhancedPreviewServer(port string, projectID string) *EnhancedPreviewServer {
	promptLoader, err := prompts.NewPromptLoader()
	if err != nil {
		log.Printf("Warning: Failed to load prompts: %v", err)
		promptLoader = nil
	}

	// Initialize Claude client
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	claudeClient := NewClaudeClient(apiKey)

	// Create a stub implementation for now
	var vectorGraphService VectorGraphServiceInterface
	vectorGraphService = nil // Will be set externally if needed

	return &EnhancedPreviewServer{
		port:               port,
		clients:            make(map[*websocket.Conn]bool),
		upgrader:           websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
		promptLoader:       promptLoader,
		currentView:        &prompts.CurrentViewState{},
		claudeClient:       claudeClient,
		vectorGraphService: vectorGraphService,
		projectID:          projectID,
		sessionID:          fmt.Sprintf("session_%d", time.Now().Unix()),
	}
}

// SetVectorGraphService sets the vector graph service (dependency injection)
func (eps *EnhancedPreviewServer) SetVectorGraphService(service VectorGraphServiceInterface) {
	eps.vectorGraphService = service
}

// broadcastToClients sends a message to all connected clients
func (eps *EnhancedPreviewServer) broadcastToClients(message interface{}) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	for client := range eps.clients {
		if err := client.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
			log.Printf("Error sending message to client: %v", err)
			client.Close()
			delete(eps.clients, client)
		}
	}
}

// handleViewUpdate handles view update HTTP requests
func (eps *EnhancedPreviewServer) handleViewUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"message": "View update endpoint (stub)",
	})
}

// handleLiveUpdate handles live update HTTP requests
func (eps *EnhancedPreviewServer) handleLiveUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"message": "Live update endpoint (stub)",
	})
}

// handleChatMessage handles chat message HTTP requests
func (eps *EnhancedPreviewServer) handleChatMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"message": "Chat message endpoint (stub)",
	})
}

// handleViewStateUpdate handles view state update HTTP requests
func (eps *EnhancedPreviewServer) handleViewStateUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"message": "View state update endpoint (stub)",
	})
}

// handleGetViewState handles get view state HTTP requests
func (eps *EnhancedPreviewServer) handleGetViewState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"data":   eps.currentView,
	})
}

// handleStatus handles status HTTP requests
func (eps *EnhancedPreviewServer) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status := map[string]interface{}{
		"connected":         true,
		"project_id":        eps.projectID,
		"session_id":        eps.sessionID,
		"vector_enabled":    eps.vectorGraphService != nil && eps.vectorGraphService.IsEnabled(),
		"clients_connected": len(eps.clients),
	}
	json.NewEncoder(w).Encode(status)
}

// Start starts the enhanced preview server
func (eps *EnhancedPreviewServer) Start() error {
	r := mux.NewRouter()

	// Apply CORS middleware
	r.Use(eps.corsMiddleware)

	// WebSocket endpoint
	r.HandleFunc("/ws", eps.handleWebSocket)

	// HTTP API endpoints
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/view-update", eps.handleViewUpdate).Methods("POST", "OPTIONS")
	api.HandleFunc("/live-update", eps.handleLiveUpdate).Methods("POST", "OPTIONS")
	api.HandleFunc("/chat", eps.handleChatMessage).Methods("POST", "OPTIONS")
	api.HandleFunc("/view-state", eps.handleViewStateUpdate).Methods("POST", "OPTIONS")
	api.HandleFunc("/view-state", eps.handleGetViewState).Methods("GET", "OPTIONS")
	api.HandleFunc("/status", eps.handleStatus).Methods("GET", "OPTIONS")

	// Enhanced endpoints for vector/graph functionality
	api.HandleFunc("/semantic-search", eps.handleSemanticSearch).Methods("POST", "OPTIONS")
	api.HandleFunc("/related-components/{id}", eps.handleRelatedComponents).Methods("GET", "OPTIONS")
	api.HandleFunc("/relationship-insights", eps.handleRelationshipInsights).Methods("GET", "OPTIONS")
	api.HandleFunc("/project-stats", eps.handleProjectStats).Methods("GET", "OPTIONS")

	log.Printf("🚀 Enhanced Preview server starting on port %s", eps.port)
	log.Printf("📡 WebSocket: ws://localhost:%s/ws", eps.port)
	log.Printf("🌐 HTTP API: http://localhost:%s/api", eps.port)
	log.Printf("🔍 Semantic Search: http://localhost:%s/api/semantic-search", eps.port)
	log.Printf("📊 Vector/Graph enabled: %t", eps.vectorGraphService != nil && eps.vectorGraphService.IsEnabled())

	return http.ListenAndServe(":"+eps.port, r)
}

// corsMiddleware adds CORS headers
func (eps *EnhancedPreviewServer) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// handleWebSocket handles WebSocket connections
func (eps *EnhancedPreviewServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := eps.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("❌ WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	eps.clients[conn] = true
	log.Printf("✅ New WebSocket connection established")

	// Send current view state to new client
	eps.sendCurrentViewState(conn)

	// Handle messages
	for {
		var msg WebSocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("❌ WebSocket read error: %v", err)
			break
		}

		switch msg.Type {
		case "view_update":
			eps.handleWebSocketViewUpdate(conn, msg)
		case "live_update":
			eps.handleWebSocketLiveUpdate(conn, msg)
		case "chat_message":
			eps.handleWebSocketChatMessage(conn, msg)
		case "view_state_update":
			eps.handleWebSocketViewStateUpdate(conn, msg)
		}
	}

	delete(eps.clients, conn)
	log.Printf("🔌 WebSocket connection closed")
}

// sendCurrentViewState sends the current view state to a client
func (eps *EnhancedPreviewServer) sendCurrentViewState(conn *websocket.Conn) {
	msg := WebSocketMessage{
		Type:      "view_state_update",
		Data:      eps.currentView,
		Timestamp: time.Now(),
	}

	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("❌ Error sending view state: %v", err)
	}
}

// handleWebSocketViewUpdate handles view updates from WebSocket
func (eps *EnhancedPreviewServer) handleWebSocketViewUpdate(conn *websocket.Conn, msg WebSocketMessage) {
	// Update current view and store components in vector/graph storage
	if data, ok := msg.Data.(map[string]interface{}); ok {
		eps.processViewUpdate(data)
	}

	// Broadcast to all clients
	eps.broadcastToClients(map[string]interface{}{
		"type": "view_updated",
		"data": msg.Data,
	})
}

// handleWebSocketLiveUpdate handles live updates from WebSocket
func (eps *EnhancedPreviewServer) handleWebSocketLiveUpdate(conn *websocket.Conn, msg WebSocketMessage) {
	// Process live updates and potentially create relationships
	if data, ok := msg.Data.(map[string]interface{}); ok {
		eps.processLiveUpdate(data)
	}

	// Broadcast to all clients
	eps.broadcastToClients(map[string]interface{}{
		"type": "live_updated",
		"data": msg.Data,
	})
}

// handleWebSocketChatMessage handles chat messages from WebSocket
func (eps *EnhancedPreviewServer) handleWebSocketChatMessage(conn *websocket.Conn, msg WebSocketMessage) {
	if data, ok := msg.Data.(map[string]interface{}); ok {
		if messageText, ok := data["message"].(string); ok {
			// Process with Claude and store conversation
			response := eps.processChatMessageWithStorage(messageText)

			// Send response back to client
			responseMsg := WebSocketMessage{
				Type:      "chat_response",
				Data:      response,
				Timestamp: time.Now(),
			}

			if err := conn.WriteJSON(responseMsg); err != nil {
				log.Printf("❌ Error sending chat response: %v", err)
			}
		}
	}
}

// handleWebSocketViewStateUpdate handles view state updates from WebSocket
func (eps *EnhancedPreviewServer) handleWebSocketViewStateUpdate(conn *websocket.Conn, msg WebSocketMessage) {
	if data, ok := msg.Data.(map[string]interface{}); ok {
		if components, ok := data["components"].([]interface{}); ok {
			eps.currentView.Components = make([]prompts.ComponentState, len(components))
			for i, comp := range components {
				if compMap, ok := comp.(map[string]interface{}); ok {
					eps.currentView.Components[i] = eps.mapToComponentState(compMap)
				}
			}
		}
	}
	// Broadcast to all clients
	eps.broadcastToClients(map[string]interface{}{
		"type": "view_state_updated",
		"data": eps.currentView,
	})
}

// processViewUpdate processes view updates and stores components
func (eps *EnhancedPreviewServer) processViewUpdate(data map[string]interface{}) {
	// Update current view
	if components, ok := data["components"].([]interface{}); ok {
		eps.currentView.Components = make([]prompts.ComponentState, 0)

		for _, comp := range components {
			if compMap, ok := comp.(map[string]interface{}); ok {
				component := eps.mapToComponentState(compMap)
				eps.currentView.Components = append(eps.currentView.Components, component)

				// Store in vector/graph storage
				if eps.vectorGraphService != nil && eps.vectorGraphService.IsEnabled() {
					ctx := context.Background()
					if err := eps.vectorGraphService.StoreComponent(ctx, &component); err != nil {
						log.Printf("⚠️ Failed to store component in vector/graph storage: %v", err)
					}
				}
			}
		}
	}
}

// processLiveUpdate processes live updates and creates relationships
func (eps *EnhancedPreviewServer) processLiveUpdate(data map[string]interface{}) {
	// Extract component information and create relationships
	if componentID, ok := data["componentId"].(string); ok {
		if action, ok := data["action"].(string); ok {
			switch action {
			case "update":
				// Find related components and create relationships
				eps.createComponentRelationships(componentID)
			case "add":
				// Create relationships for new components
				eps.createComponentRelationships(componentID)
			}
		}
	}
}

// processChatMessageWithStorage processes chat messages and stores conversations
func (eps *EnhancedPreviewServer) processChatMessageWithStorage(message string) map[string]interface{} {
	// Create conversation message
	userMessage := &Message{
		Role:      "user",
		Content:   message,
		Timestamp: time.Now(),
		Type:      "text",
	}

	// Store user message
	if eps.vectorGraphService != nil && eps.vectorGraphService.IsEnabled() {
		ctx := context.Background()
		if err := eps.vectorGraphService.StoreConversation(ctx, userMessage.Content, eps.sessionID); err != nil {
			log.Printf("⚠️ Failed to store user message: %v", err)
		}
	}

	// Get Claude response
	response := eps.processChatMessageWithStorage(message)

	// Create assistant message
	assistantMessage := &Message{
		Role:      "assistant",
		Content:   response["response"].(string),
		Timestamp: time.Now(),
		Type:      "text",
	}

	// Store assistant message
	if eps.vectorGraphService != nil && eps.vectorGraphService.IsEnabled() {
		ctx := context.Background()
		if err := eps.vectorGraphService.StoreConversation(ctx, assistantMessage.Content, eps.sessionID); err != nil {
			log.Printf("⚠️ Failed to store assistant message: %v", err)
		}
	}

	return map[string]interface{}{
		"response":  response["response"],
		"action":    response["action"],
		"data":      response["data"],
		"timestamp": response["timestamp"],
	}
}

// createComponentRelationships creates relationships between components
func (eps *EnhancedPreviewServer) createComponentRelationships(componentID string) {
	if eps.vectorGraphService == nil || !eps.vectorGraphService.IsEnabled() {
		return
	}

	ctx := context.Background()

	// Find components that are close to each other (spatial relationships)
	currentComponent := eps.findComponentByID(componentID)
	if currentComponent == nil {
		return
	}

	for _, otherComponent := range eps.currentView.Components {
		if otherComponent.ID == componentID {
			continue
		}

		// Calculate distance between components
		distance := eps.calculateComponentDistance(currentComponent, &otherComponent)

		// Create relationship if components are close
		if distance < 200 { // Threshold for proximity
			weight := 1.0 - (distance / 200.0) // Higher weight for closer components

			if err := eps.vectorGraphService.CreateComponentRelationship(ctx, componentID, otherComponent.ID, "NEAR", weight); err != nil {
				log.Printf("⚠️ Failed to create proximity relationship: %v", err)
			}
		}
	}
}

// New HTTP handlers for enhanced functionality

// handleSemanticSearch handles semantic search requests
func (eps *EnhancedPreviewServer) handleSemanticSearch(w http.ResponseWriter, r *http.Request) {
	if eps.vectorGraphService == nil || !eps.vectorGraphService.IsEnabled() {
		http.Error(w, "Vector/Graph storage not enabled", http.StatusServiceUnavailable)
		return
	}

	var request struct {
		Query string `json:"query"`
		Limit int    `json:"limit"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.Limit == 0 {
		request.Limit = 10
	}

	ctx := context.Background()
	results, err := eps.vectorGraphService.SemanticSearch(ctx, request.Query, request.Limit)
	if err != nil {
		log.Printf("❌ Semantic search error: %v", err)
		http.Error(w, "Search failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// handleRelatedComponents handles related components requests
func (eps *EnhancedPreviewServer) handleRelatedComponents(w http.ResponseWriter, r *http.Request) {
	if eps.vectorGraphService == nil || !eps.vectorGraphService.IsEnabled() {
		http.Error(w, "Vector/Graph storage not enabled", http.StatusServiceUnavailable)
		return
	}

	vars := mux.Vars(r)
	componentID := vars["id"]

	ctx := context.Background()
	components, err := eps.vectorGraphService.FindRelatedComponents(ctx, componentID, 3)
	if err != nil {
		log.Printf("❌ Find related components error: %v", err)
		http.Error(w, "Find related components failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Handle the interface{} type from the service
	componentCount := 0
	if componentsSlice, ok := components.([]interface{}); ok {
		componentCount = len(componentsSlice)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"component_id": componentID,
		"related":      components,
		"count":        componentCount,
	})
}

// handleRelationshipInsights handles relationship insights requests
func (eps *EnhancedPreviewServer) handleRelationshipInsights(w http.ResponseWriter, r *http.Request) {
	if eps.vectorGraphService == nil || !eps.vectorGraphService.IsEnabled() {
		http.Error(w, "Vector/Graph storage not enabled", http.StatusServiceUnavailable)
		return
	}

	ctx := context.Background()
	insights, err := eps.vectorGraphService.GetRelationshipInsights(ctx)
	if err != nil {
		log.Printf("❌ Get relationship insights error: %v", err)
		http.Error(w, "Get insights failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Handle the interface{} type from the service
	insightCount := 0
	if insightsSlice, ok := insights.([]interface{}); ok {
		insightCount = len(insightsSlice)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"insights": insights,
		"count":    insightCount,
	})
}

// handleProjectStats handles project statistics requests
func (eps *EnhancedPreviewServer) handleProjectStats(w http.ResponseWriter, r *http.Request) {
	if eps.vectorGraphService == nil || !eps.vectorGraphService.IsEnabled() {
		http.Error(w, "Vector/Graph storage not enabled", http.StatusServiceUnavailable)
		return
	}

	ctx := context.Background()
	stats, err := eps.vectorGraphService.GetProjectStats(ctx)
	if err != nil {
		log.Printf("❌ Get project stats error: %v", err)
		http.Error(w, "Get stats failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// Utility functions

// findComponentByID finds a component by its ID
func (eps *EnhancedPreviewServer) findComponentByID(id string) *prompts.ComponentState {
	for _, component := range eps.currentView.Components {
		if component.ID == id {
			return &component
		}
	}
	return nil
}

// calculateComponentDistance calculates distance between two components
func (eps *EnhancedPreviewServer) calculateComponentDistance(comp1, comp2 *prompts.ComponentState) float64 {
	dx := float64(comp1.Position.X - comp2.Position.X)
	dy := float64(comp1.Position.Y - comp2.Position.Y)
	return math.Sqrt(dx*dx + dy*dy)
}

// mapToComponentState converts a map to ComponentState
func (eps *EnhancedPreviewServer) mapToComponentState(data map[string]interface{}) prompts.ComponentState {
	component := prompts.ComponentState{
		Properties: make(map[string]interface{}),
		Style:      make(map[string]interface{}),
	}

	if id, ok := data["id"].(string); ok {
		component.ID = id
	}
	if compType, ok := data["type"].(string); ok {
		component.Type = compType
	}
	if name, ok := data["name"].(string); ok {
		component.Name = name
	}
	if category, ok := data["category"].(string); ok {
		component.Category = category
	}

	// Handle position
	if pos, ok := data["position"].(map[string]interface{}); ok {
		if x, ok := pos["x"].(float64); ok {
			component.Position.X = int(x)
		}
		if y, ok := pos["y"].(float64); ok {
			component.Position.Y = int(y)
		}
	}

	// Handle size
	if size, ok := data["size"].(map[string]interface{}); ok {
		if w, ok := size["w"].(float64); ok {
			component.Size.W = int(w)
		}
		if h, ok := size["h"].(float64); ok {
			component.Size.H = int(h)
		}
	}

	// Handle properties
	if props, ok := data["properties"].(map[string]interface{}); ok {
		component.Properties = props
	}

	// Handle style
	if style, ok := data["style"].(map[string]interface{}); ok {
		component.Style = style
	}

	return component
}

// Additional handler methods would be implemented here following the same pattern as the original PreviewServer...
