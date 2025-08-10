package vibe

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/vibercode/cli/internal/vibe/prompts"
	"github.com/vibercode/cli/pkg/ui"
)

type ChatManager struct {
	client       *ClaudeClient
	context      *ProjectContext
	preview      *PreviewServer
	conversation []Message
	promptLoader *prompts.PromptLoader
	mode         string // "general" or "component"
	logger       *VibeLogger
	chatBridge   *ChatBridge
	ctx          context.Context
	cancel       context.CancelFunc
}

type Message struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"` // "text", "code", "ui_update"
}

type ProjectContext struct {
	Templates   map[string]string `json:"templates"`
	CurrentView *UIView           `json:"current_view"`
	ProjectInfo ProjectInfo       `json:"project_info"`
}

type UIView struct {
	Components []UIComponent `json:"components"`
	Layout     Layout        `json:"layout"`
	Theme      Theme         `json:"theme"`
	Timestamp  string        `json:"timestamp"`
}

type UIComponent struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Name        string                 `json:"name"`
	Position    Position               `json:"position"`
	Size        Size                   `json:"size"`
	Properties  map[string]interface{} `json:"properties"`
	Constraints Constraints            `json:"constraints"`
}

type Layout struct {
	Grid             int    `json:"grid"`
	RowHeight        int    `json:"row_height"`
	Margin           [2]int `json:"margin"`
	ContainerPadding [2]int `json:"container_padding"`
}

type Theme struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Colors  map[string]string `json:"colors"`
	Effects Effects           `json:"effects"`
}

type Effects struct {
	Glow       bool `json:"glow"`
	Gradients  bool `json:"gradients"`
	Animations bool `json:"animations"`
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

type ProjectInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Module      string `json:"module"`
	Port        int    `json:"port"`
}

func NewChatManager(mode string) *ChatManager {
	// Initialize logger for chat mode
	logger := NewVibeLogger(true) // true = chat mode

	// Initialize Claude client
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		fmt.Println(ui.Warning.Sprint("⚠️  ANTHROPIC_API_KEY not found. Set it to enable AI chat:"))
		fmt.Println(ui.Muted.Sprint("   export ANTHROPIC_API_KEY=your_api_key"))
	}

	client := NewClaudeClient(apiKey)

	// Initialize prompt loader
	promptLoader, err := prompts.NewPromptLoader()
	if err != nil {
		fmt.Printf(ui.Error.Sprint("❌ Failed to load prompts: %v\n"), err)
		return nil
	}

	// Create cancelable context
	ctx, cancel := context.WithCancel(context.Background())

	return &ChatManager{
		client:       client,
		context:      &ProjectContext{},
		conversation: make([]Message, 0),
		promptLoader: promptLoader,
		mode:         mode,
		logger:       logger,
		ctx:          ctx,
		cancel:       cancel,
	}
}

func (cm *ChatManager) StartChat() {
	modeText := "VibeCode Chat Mode"
	if cm.mode == "component" {
		modeText = "VibeCode Component Mode"
	}

	fmt.Println(ui.Primary.Sprint(fmt.Sprintf("\n🎨 Welcome to %s!", modeText)))
	fmt.Println(ui.Muted.Sprint("Type 'exit' to quit, 'clear' to clear conversation"))

	if cm.mode == "component" {
		fmt.Println(ui.Info.Sprint("💡 Component Mode: All interactions focused on UI components"))
		fmt.Println(ui.Info.Sprint("   Commands: 'agregar botón', 'cambiar tema', 'estado del canvas'"))
	}

	fmt.Println(ui.Muted.Sprint("──────────────────────────────────────────────────────"))

	// Load project context
	cm.loadProjectContext()

	// Initialize default UI view
	cm.initializeDefaultView()

	// Start preview server with component-focused mode
	cm.preview = NewPreviewServer("3001")
	if cm.mode == "component" {
		cm.preview.SetComponentMode(true)
	}

	// Initialize chat bridge
	cm.chatBridge = cm.preview.chatBridge

	go cm.preview.Start()

	// Show preview URL
	fmt.Println(ui.Success.Sprint("🔥 Preview server started at http://localhost:3001"))
	fmt.Println(ui.Info.Sprint("📡 WebSocket available at ws://localhost:3001/ws"))
	fmt.Println("")

	reader := bufio.NewReader(os.Stdin)

	// Create a channel to handle input in a non-blocking way
	inputChan := make(chan string, 1)
	errorChan := make(chan error, 1)

	// Helper function to read input safely
	readInput := func() {
		fmt.Print(ui.Primary.Sprint("\n💬 You: "))
		input, err := reader.ReadString('\n')
		if err != nil {
			errorChan <- err
			return
		}
		inputChan <- input
	}

	for {
		// Start reading input in a goroutine
		go readInput()

		var input string
		var err error

		// Wait for input or context cancellation
		select {
		case input = <-inputChan:
			// Continue with normal processing
		case err = <-errorChan:
			fmt.Printf("Error reading input: %v\n", err)
			continue
		case <-time.After(100 * time.Millisecond):
			// Check if context is cancelled
			select {
			case <-cm.ctx.Done():
				fmt.Println(ui.Info.Sprint("\n🛑 Chat interrupted"))
				return
			default:
				continue
			}
		}

		input = strings.TrimSpace(input)

		switch input {
		case "exit", "quit", "q":
			fmt.Println(ui.Success.Sprint("👋 Goodbye!"))
			return
		case "clear":
			cm.conversation = make([]Message, 0)
			fmt.Println(ui.Info.Sprint("🧹 Conversation cleared"))
			continue
		case "":
			continue
		}

		// In component mode, add context to make everything component-focused
		if cm.mode == "component" {
			input = cm.enhanceComponentInput(input)
		}

		// Add user message
		userMsg := Message{
			Role:      "user",
			Content:   input,
			Timestamp: time.Now(),
			Type:      "text",
		}
		cm.conversation = append(cm.conversation, userMsg)

		// Send through chat bridge
		if cm.chatBridge != nil {
			cm.chatBridge.SendTerminalMessage(input)
		} else {
			// Fallback to direct processing
			response, err := cm.getAIResponse(input)
			if err != nil {
				fmt.Printf(ui.Error.Sprint("❌ Error: %v\n"), err)
				continue
			}

			// Add AI response
			aiMsg := Message{
				Role:      "assistant",
				Content:   response,
				Timestamp: time.Now(),
				Type:      "text",
			}
			cm.conversation = append(cm.conversation, aiMsg)

			// Display response
			fmt.Printf(ui.Secondary.Sprint("\n🤖 Viber: ")+"%s\n", response)

			// Check if response contains UI updates
			cm.processUIUpdates(response)
		}
	}
}

// enhanceComponentInput adds context to make inputs component-focused
func (cm *ChatManager) enhanceComponentInput(input string) string {
	// Add component context to general queries
	lowerInput := strings.ToLower(input)

	// If the user is asking general questions, redirect to component context
	if !strings.Contains(lowerInput, "componente") && !strings.Contains(lowerInput, "component") &&
		!strings.Contains(lowerInput, "botón") && !strings.Contains(lowerInput, "button") &&
		!strings.Contains(lowerInput, "agregar") && !strings.Contains(lowerInput, "add") &&
		!strings.Contains(lowerInput, "tema") && !strings.Contains(lowerInput, "theme") &&
		!strings.Contains(lowerInput, "canvas") && !strings.Contains(lowerInput, "estado") {

		// Add component context to generic requests
		return fmt.Sprintf("En el contexto de componentes UI: %s", input)
	}

	return input
}

func (cm *ChatManager) loadProjectContext() {
	// Load templates from internal/templates
	templates := map[string]string{
		"model":      "SchemaModelTemplate",
		"repository": "SchemaRepositoryTemplate",
		"service":    "SchemaServiceTemplate",
		"handler":    "SchemaHandlerTemplate",
	}

	cm.context.Templates = templates

	// Set default project info
	cm.context.ProjectInfo = ProjectInfo{
		Name:        "vibercode-project",
		Description: "Generated with VibeCode CLI",
		Module:      "github.com/example/project",
		Port:        8080,
	}
}

func (cm *ChatManager) initializeDefaultView() {
	defaultView := &UIView{
		Components: []UIComponent{
			{
				ID:       "button-" + fmt.Sprintf("%d", time.Now().UnixNano()),
				Type:     "atom",
				Name:     "Button",
				Position: Position{X: 362.125, Y: 193},
				Size:     Size{W: 160, H: 40},
				Properties: map[string]interface{}{
					"text":    "Click me",
					"variant": "primary",
					"size":    "medium",
				},
				Constraints: Constraints{
					W: 160, H: 40,
					MinW: 80, MinH: 32,
					MaxW: 320, MaxH: 60,
				},
			},
		},
		Layout: Layout{
			Grid:             12,
			RowHeight:        60,
			Margin:           [2]int{12, 12},
			ContainerPadding: [2]int{16, 16},
		},
		Theme: Theme{
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
			Effects: Effects{
				Glow:       true,
				Gradients:  true,
				Animations: true,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339Nano),
	}

	cm.context.CurrentView = defaultView
}

func (cm *ChatManager) getAIResponse(input string) (string, error) {
	if cm.client == nil {
		return "AI chat is not available. Please set ANTHROPIC_API_KEY environment variable.", nil
	}

	// Build context for Claude using prompt loader
	contextStr := cm.buildContextString()

	// Prepare conversation history for prompt
	conversationHistory := make([]prompts.ConversationMessage, len(cm.conversation))
	for i, msg := range cm.conversation {
		conversationHistory[i] = prompts.ConversationMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// Build prompt data
	promptData := prompts.PromptData{
		ProjectContext:      contextStr,
		Templates:           cm.context.Templates,
		CurrentView:         cm.convertUIViewToCurrentViewState(),
		UserInput:           input,
		ConversationHistory: conversationHistory,
		Mode:                cm.mode, // Pass the mode to the prompt
	}

	// Generate full prompt using prompt loader
	fullPrompt, err := cm.promptLoader.BuildChatPrompt(promptData)
	if err != nil {
		return "", fmt.Errorf("failed to build prompt: %w", err)
	}

	// Create Claude message
	messages := []ClaudeMessage{{
		Role:    "user",
		Content: fullPrompt,
	}}

	response, err := cm.client.CreateMessage(messages)
	if err != nil {
		return "", err
	}

	return response, nil
}

func (cm *ChatManager) buildContextString() string {
	contextBytes, _ := json.MarshalIndent(cm.context, "", "  ")
	return string(contextBytes)
}

func (cm *ChatManager) processUIUpdates(response string) {
	// Extract JSON from response using prompt loader
	jsonStr, hasJSON := prompts.ExtractJSONFromResponse(response)
	if !hasJSON {
		return
	}

	// Validate the JSON structure
	if err := prompts.ValidateUIUpdateJSON(jsonStr); err != nil {
		fmt.Printf(ui.Warning.Sprint("⚠️  Invalid UI update JSON: %v\n"), err)
		return
	}

	// Parse and process the update
	var update map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &update); err != nil {
		fmt.Printf(ui.Warning.Sprint("⚠️  Failed to parse UI update: %v\n"), err)
		return
	}

	// Process the UI update
	cm.handleUIUpdate(update)

	// Show confirmation of update
	if explanation, ok := update["explanation"].(string); ok {
		fmt.Printf(ui.Success.Sprint("✨ UI Update: ")+"%s\n", explanation)
	}
}

func (cm *ChatManager) handleUIUpdate(update map[string]interface{}) {
	if cm.preview != nil {
		// Send update to preview server via WebSocket
		// cm.preview.BroadcastChatResponse(update) // TODO: Fix this call

		// Update local context
		if action, ok := update["action"].(string); ok {
			switch action {
			case "update_theme":
				if data, ok := update["data"].(map[string]interface{}); ok {
					cm.updateTheme(data)
				}
			case "update_component":
				if data, ok := update["data"].(map[string]interface{}); ok {
					cm.updateComponent(data)
				}
			case "add_component":
				if data, ok := update["data"].(map[string]interface{}); ok {
					cm.addComponent(data)
				}
			case "remove_component":
				if data, ok := update["data"].(map[string]interface{}); ok {
					cm.removeComponent(data)
				}
			}
		}
	}
}

func (cm *ChatManager) updateTheme(data map[string]interface{}) {
	// Update theme in current view
	if colors, ok := data["colors"].(map[string]interface{}); ok {
		for key, value := range colors {
			if colorStr, ok := value.(string); ok {
				cm.context.CurrentView.Theme.Colors[key] = colorStr
			}
		}
	}

	if effects, ok := data["effects"].(map[string]interface{}); ok {
		if glow, ok := effects["glow"].(bool); ok {
			cm.context.CurrentView.Theme.Effects.Glow = glow
		}
		if gradients, ok := effects["gradients"].(bool); ok {
			cm.context.CurrentView.Theme.Effects.Gradients = gradients
		}
		if animations, ok := effects["animations"].(bool); ok {
			cm.context.CurrentView.Theme.Effects.Animations = animations
		}
	}
}

func (cm *ChatManager) updateComponent(data map[string]interface{}) {
	if id, ok := data["id"].(string); ok {
		// Find and update component
		for i, comp := range cm.context.CurrentView.Components {
			if comp.ID == id {
				// Update properties
				if props, ok := data["properties"].(map[string]interface{}); ok {
					for key, value := range props {
						cm.context.CurrentView.Components[i].Properties[key] = value
					}
				}
				break
			}
		}
	}
}

func (cm *ChatManager) addComponent(data map[string]interface{}) {
	// Convert data to UIComponent
	comp := UIComponent{
		ID:         fmt.Sprintf("comp-%d", time.Now().UnixNano()),
		Type:       "atom",
		Properties: make(map[string]interface{}),
	}

	if name, ok := data["name"].(string); ok {
		comp.Name = name
	}
	if compType, ok := data["type"].(string); ok {
		comp.Type = compType
	}
	if props, ok := data["properties"].(map[string]interface{}); ok {
		comp.Properties = props
	}

	cm.context.CurrentView.Components = append(cm.context.CurrentView.Components, comp)
}

func (cm *ChatManager) removeComponent(data map[string]interface{}) {
	if id, ok := data["id"].(string); ok {
		// Find and remove component by ID
		for i, comp := range cm.context.CurrentView.Components {
			if comp.ID == id {
				// Remove component from slice
				cm.context.CurrentView.Components = append(
					cm.context.CurrentView.Components[:i],
					cm.context.CurrentView.Components[i+1:]...,
				)
				fmt.Printf(ui.Info.Sprint("🗑️  Removed component: ")+"%s (%s)\n", comp.Name, comp.ID)
				return
			}
		}
		fmt.Printf(ui.Warning.Sprint("⚠️  Component not found: ")+"%s\n", id)
	}
}

// convertUIViewToCurrentViewState converts the current UIView to prompts.CurrentViewState
func (cm *ChatManager) convertUIViewToCurrentViewState() *prompts.CurrentViewState {
	if cm.context == nil || cm.context.CurrentView == nil {
		// Return default view state if no current view
		return &prompts.CurrentViewState{
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
	}

	uiView := cm.context.CurrentView

	// Convert components
	components := make([]prompts.ComponentState, len(uiView.Components))
	for i, comp := range uiView.Components {
		// Determine category based on component type
		category := "atom" // default
		switch comp.Type {
		case "button", "text", "input", "image":
			category = "atom"
		case "card", "form", "navigation":
			category = "molecule"
		case "hero", "gallery", "section":
			category = "organism"
		}

		components[i] = prompts.ComponentState{
			ID:         comp.ID,
			Type:       comp.Type,
			Name:       comp.Name,
			Category:   category,
			Properties: comp.Properties,
			Position: prompts.Position{
				X: int(comp.Position.X),
				Y: int(comp.Position.Y),
			},
			Size: prompts.Size{
				W: int(comp.Size.W),
				H: int(comp.Size.H),
			},
			Style: make(map[string]interface{}), // Empty style for now
		}
	}

	// Convert theme
	themeState := prompts.ThemeState{
		ID:     uiView.Theme.ID,
		Name:   uiView.Theme.Name,
		Colors: uiView.Theme.Colors,
		Effects: map[string]interface{}{
			"glow":       uiView.Theme.Effects.Glow,
			"gradients":  uiView.Theme.Effects.Gradients,
			"animations": uiView.Theme.Effects.Animations,
		},
	}

	// Convert layout
	layoutState := prompts.LayoutState{
		Grid:             uiView.Layout.Grid,
		RowHeight:        uiView.Layout.RowHeight,
		Margin:           uiView.Layout.Margin,
		ContainerPadding: uiView.Layout.ContainerPadding,
		ShowGrid:         true, // Default value
		SnapToGrid:       true, // Default value
	}

	// Create default canvas state
	canvasState := prompts.CanvasState{
		Viewport:     "desktop",
		Zoom:         1.0,
		PanOffset:    prompts.Position{X: 0, Y: 0},
		SelectedItem: "",
	}

	return &prompts.CurrentViewState{
		Components: components,
		Theme:      themeState,
		Layout:     layoutState,
		Canvas:     canvasState,
	}
}
