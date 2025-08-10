package prompts

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
)

//go:embed *.md
var promptFiles embed.FS

// PromptTemplate represents a prompt template with variables
type PromptTemplate struct {
	content  string
	template *template.Template
}

// PromptLoader handles loading and rendering prompt templates
type PromptLoader struct {
	templates map[string]*PromptTemplate
}

// CurrentViewState represents the current state of the UI
type CurrentViewState struct {
	Components []ComponentState `json:"components"`
	Theme      ThemeState       `json:"theme"`
	Layout     LayoutState      `json:"layout"`
	Canvas     CanvasState      `json:"canvas"`
}

// ComponentState represents a component in the current view
type ComponentState struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Name       string                 `json:"name"`
	Category   string                 `json:"category"`
	Properties map[string]interface{} `json:"properties"`
	Position   Position               `json:"position"`
	Size       Size                   `json:"size"`
	Style      map[string]interface{} `json:"style,omitempty"`
}

// ThemeState represents the current theme
type ThemeState struct {
	ID      string                 `json:"id"`
	Name    string                 `json:"name"`
	Colors  map[string]string      `json:"colors"`
	Effects map[string]interface{} `json:"effects"`
}

// LayoutState represents the current layout configuration
type LayoutState struct {
	Grid             int    `json:"grid"`
	RowHeight        int    `json:"rowHeight"`
	Margin           [2]int `json:"margin"`
	ContainerPadding [2]int `json:"containerPadding"`
	ShowGrid         bool   `json:"showGrid"`
	SnapToGrid       bool   `json:"snapToGrid"`
}

// CanvasState represents the current canvas state
type CanvasState struct {
	Viewport     string   `json:"viewport"`
	Zoom         float64  `json:"zoom"`
	PanOffset    Position `json:"panOffset"`
	SelectedItem string   `json:"selectedItem,omitempty"`
}

// Position represents x,y coordinates
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Size represents width and height
type Size struct {
	W int `json:"w"`
	H int `json:"h"`
}

// PromptData contains data to be injected into prompts
type PromptData struct {
	ProjectContext      string
	Templates           map[string]string
	CurrentView         *CurrentViewState
	UserInput           string
	ConversationHistory []ConversationMessage
	Mode                string // "general" or "component"
}

// ConversationMessage represents a message in the conversation history
type ConversationMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// NewPromptLoader creates a new prompt loader
func NewPromptLoader() (*PromptLoader, error) {
	loader := &PromptLoader{
		templates: make(map[string]*PromptTemplate),
	}

	err := loader.loadPrompts()
	if err != nil {
		return nil, fmt.Errorf("failed to load prompts: %w", err)
	}

	return loader, nil
}

// loadPrompts loads all prompt templates from embedded files
func (pl *PromptLoader) loadPrompts() error {
	files := []string{
		"system.md",
		"ui_examples.md",
		"go_api_examples.md",
	}

	for _, file := range files {
		content, err := promptFiles.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", file, err)
		}

		// Create template name from filename
		name := strings.TrimSuffix(file, ".md")

		// Parse template
		tmpl, err := template.New(name).Parse(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", name, err)
		}

		pl.templates[name] = &PromptTemplate{
			content:  string(content),
			template: tmpl,
		}
	}

	return nil
}

// RenderSystemPrompt renders the system prompt with context
func (pl *PromptLoader) RenderSystemPrompt(data PromptData) (string, error) {
	tmpl, exists := pl.templates["system"]
	if !exists {
		return "", fmt.Errorf("system prompt template not found")
	}

	var buf strings.Builder
	err := tmpl.template.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("failed to render system prompt: %w", err)
	}

	return buf.String(), nil
}

// GetUIExamples returns UI examples for reference
func (pl *PromptLoader) GetUIExamples() (string, error) {
	tmpl, exists := pl.templates["ui_examples"]
	if !exists {
		return "", fmt.Errorf("ui_examples template not found")
	}

	return tmpl.content, nil
}

// GetGoAPIExamples returns Go API examples for reference
func (pl *PromptLoader) GetGoAPIExamples() (string, error) {
	tmpl, exists := pl.templates["go_api_examples"]
	if !exists {
		return "", fmt.Errorf("go_api_examples template not found")
	}

	return tmpl.content, nil
}

// BuildChatPrompt builds a complete prompt for the chat interface
func (pl *PromptLoader) BuildChatPrompt(data PromptData) (string, error) {
	// Render system prompt
	systemPrompt, err := pl.RenderSystemPrompt(data)
	if err != nil {
		return "", err
	}

	// Build current view context
	var currentViewContext strings.Builder
	if data.Mode == "component" {
		currentViewContext.WriteString("\n## Current Component Canvas State (Component Mode):\n")
	} else {
		currentViewContext.WriteString("\n## Current Canvas State:\n")
	}

	if data.CurrentView != nil {
		// Serialize current view state to JSON
		viewJSON, err := json.MarshalIndent(data.CurrentView, "", "  ")
		if err != nil {
			currentViewContext.WriteString("Error serializing current view state\n")
		} else {
			currentViewContext.WriteString("```json\n")
			currentViewContext.WriteString(string(viewJSON))
			currentViewContext.WriteString("\n```\n")
		}

		// Add human-readable summary
		currentViewContext.WriteString("\n### Summary:\n")
		currentViewContext.WriteString(fmt.Sprintf("- **Components**: %d total\n", len(data.CurrentView.Components)))
		currentViewContext.WriteString(fmt.Sprintf("- **Viewport**: %s\n", data.CurrentView.Canvas.Viewport))
		currentViewContext.WriteString(fmt.Sprintf("- **Theme**: %s\n", data.CurrentView.Theme.Name))
		currentViewContext.WriteString(fmt.Sprintf("- **Selected**: %s\n",
			func() string {
				if data.CurrentView.Canvas.SelectedItem != "" {
					return data.CurrentView.Canvas.SelectedItem
				}
				return "None"
			}()))

		// List components by type
		if len(data.CurrentView.Components) > 0 {
			componentsByType := make(map[string][]string)
			for _, comp := range data.CurrentView.Components {
				componentsByType[comp.Type] = append(componentsByType[comp.Type],
					fmt.Sprintf("%s (%s)", comp.Name, comp.ID))
			}

			currentViewContext.WriteString("\n### Components by Type:\n")
			for compType, comps := range componentsByType {
				currentViewContext.WriteString(fmt.Sprintf("- **%s**: %s\n",
					strings.Title(compType), strings.Join(comps, ", ")))
			}
		} else {
			currentViewContext.WriteString("\n**Canvas is empty** - No components added yet\n")
		}

		// Show available positions
		currentViewContext.WriteString(fmt.Sprintf("\n### Canvas Info:\n"))
		currentViewContext.WriteString(fmt.Sprintf("- **Grid**: %d columns, %d row height\n",
			data.CurrentView.Layout.Grid, data.CurrentView.Layout.RowHeight))
		currentViewContext.WriteString(fmt.Sprintf("- **Zoom**: %.1f%%\n",
			data.CurrentView.Canvas.Zoom*100))
		currentViewContext.WriteString(fmt.Sprintf("- **Grid Visible**: %t\n",
			data.CurrentView.Layout.ShowGrid))

		// Add component mode specific analysis
		if data.Mode == "component" {
			currentViewContext.WriteString(fmt.Sprintf("\n### Component Mode Analysis:\n"))
			currentViewContext.WriteString(fmt.Sprintf("- **Mode**: Component-focused design mode\n"))
			currentViewContext.WriteString(fmt.Sprintf("- **Focus**: All responses must be component-related\n"))
			currentViewContext.WriteString(fmt.Sprintf("- **Priority**: Visual design and component structure\n"))

			// Add design insights
			if len(data.CurrentView.Components) > 0 {
				currentViewContext.WriteString(fmt.Sprintf("- **Design Status**: %d components provide good foundation\n", len(data.CurrentView.Components)))
				currentViewContext.WriteString(fmt.Sprintf("- **Visual Hierarchy**: Review component positioning and relationships\n"))
			} else {
				currentViewContext.WriteString(fmt.Sprintf("- **Design Status**: Empty canvas - perfect for component creation\n"))
				currentViewContext.WriteString(fmt.Sprintf("- **Suggestion**: Start with fundamental components (button, text, card)\n"))
			}
		}
	} else {
		currentViewContext.WriteString("No current view state available\n")
	}

	// Build conversation context
	var conversationContext strings.Builder

	// Add recent conversation history (last 5 messages)
	if len(data.ConversationHistory) > 0 {
		conversationContext.WriteString("\n## Recent Conversation:\n")

		start := 0
		if len(data.ConversationHistory) > 5 {
			start = len(data.ConversationHistory) - 5
		}

		for _, msg := range data.ConversationHistory[start:] {
			conversationContext.WriteString(fmt.Sprintf("\n**%s**: %s\n",
				strings.Title(msg.Role), msg.Content))
		}
	}

	// Add mode-specific context
	var modeContext strings.Builder
	if data.Mode == "component" {
		modeContext.WriteString("\n## COMPONENT MODE CONTEXT:\n")
		modeContext.WriteString("- **CRITICAL**: Every response must be component-focused\n")
		modeContext.WriteString("- **REDIRECT**: Non-component questions should be redirected to component context\n")
		modeContext.WriteString("- **ANALYSIS**: Always provide visual design insights and component structure feedback\n")
		modeContext.WriteString("- **PRIORITY**: Visual hierarchy, design consistency, and component relationships\n")
		modeContext.WriteString("- **VALIDATION**: Ensure all component properties follow the defined structure\n")
	}

	// Combine system prompt with all contexts and user input
	fullPrompt := fmt.Sprintf(`%s

%s

%s

%s

## Current User Request:
%s

## Instructions:
- Consider the current canvas state when making decisions
- If adding components, choose appropriate positions that don't overlap existing ones
- If modifying components, reference them by their actual IDs from the current state
- Maintain design consistency with the current theme
- Be contextually aware of what's already on the canvas
- Always respond with valid JSON as specified in the system prompt
- Provide clear explanations of changes%s`,
		systemPrompt,
		currentViewContext.String(),
		conversationContext.String(),
		modeContext.String(),
		data.UserInput,
		func() string {
			if data.Mode == "component" {
				return "\n- **COMPONENT MODE**: Focus on visual design and component structure analysis"
			}
			return ""
		}())

	return fullPrompt, nil
}

// AnalyzeViewState provides analysis of the current view state
func (pl *PromptLoader) AnalyzeViewState(viewState *CurrentViewState) string {
	if viewState == nil {
		return "No view state available"
	}

	var analysis strings.Builder
	analysis.WriteString("## Canvas Analysis:\n")

	// Component analysis
	if len(viewState.Components) == 0 {
		analysis.WriteString("- Canvas is empty - ready for new components\n")
	} else {
		analysis.WriteString(fmt.Sprintf("- %d components present\n", len(viewState.Components)))

		// Find available space
		usedPositions := make(map[string]bool)
		for _, comp := range viewState.Components {
			key := fmt.Sprintf("%d,%d", comp.Position.X, comp.Position.Y)
			usedPositions[key] = true
		}

		// Suggest good positions for new components
		goodPositions := []Position{
			{X: 100, Y: 100}, {X: 300, Y: 100}, {X: 500, Y: 100},
			{X: 100, Y: 300}, {X: 300, Y: 300}, {X: 500, Y: 300},
			{X: 100, Y: 500}, {X: 300, Y: 500}, {X: 500, Y: 500},
		}

		availablePositions := []Position{}
		for _, pos := range goodPositions {
			key := fmt.Sprintf("%d,%d", pos.X, pos.Y)
			if !usedPositions[key] {
				availablePositions = append(availablePositions, pos)
			}
		}

		if len(availablePositions) > 0 {
			analysis.WriteString("- Good positions for new components: ")
			for i, pos := range availablePositions {
				if i > 0 {
					analysis.WriteString(", ")
				}
				analysis.WriteString(fmt.Sprintf("(%d,%d)", pos.X, pos.Y))
				if i >= 2 { // Only show first 3 positions
					analysis.WriteString("...")
					break
				}
			}
			analysis.WriteString("\n")
		}
	}

	// Theme analysis
	analysis.WriteString(fmt.Sprintf("- Current theme: %s\n", viewState.Theme.Name))
	analysis.WriteString(fmt.Sprintf("- Primary color: %s\n", viewState.Theme.Colors["primary"]))

	return analysis.String()
}

// ExtractJSONFromResponse extracts JSON from an AI response
func ExtractJSONFromResponse(response string) (string, bool) {
	// Look for JSON blocks in the response
	var start int
	var isCodeBlock bool

	// Check for code block format first
	codeBlockStart := strings.Index(response, "```json")
	if codeBlockStart != -1 {
		start = codeBlockStart + 7 // Skip "```json"
		isCodeBlock = true

		// Skip whitespace after ```json
		for start < len(response) && (response[start] == ' ' || response[start] == '\n' || response[start] == '\r') {
			start++
		}
	} else {
		// Look for direct JSON (first opening brace)
		start = strings.Index(response, "{")
		isCodeBlock = false
	}

	if start == -1 {
		return "", false
	}

	// Find the matching closing brace
	var end int
	braceCount := 0
	inString := false
	escaped := false

	for i := start; i < len(response); i++ {
		char := response[i]

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
			braceCount++
		} else if char == '}' {
			braceCount--
			if braceCount == 0 {
				end = i
				break
			}
		}
	}

	if braceCount != 0 {
		return "", false
	}

	// For code blocks, check if we should stop at ``` before the closing brace
	if isCodeBlock {
		codeBlockEnd := strings.Index(response[start:], "```")
		if codeBlockEnd != -1 && start+codeBlockEnd-1 < end {
			// Find the actual closing brace before the ```
			tempEnd := -1
			tempBraceCount := 0
			for i := start; i < start+codeBlockEnd; i++ {
				if response[i] == '{' {
					tempBraceCount++
				} else if response[i] == '}' {
					tempBraceCount--
					if tempBraceCount == 0 {
						tempEnd = i
						break
					}
				}
			}
			if tempEnd != -1 {
				end = tempEnd
			}
		}
	}

	if end == -1 || end <= start {
		return "", false
	}

	jsonStr := strings.TrimSpace(response[start : end+1])

	// Basic validation - check if it looks like JSON
	if !strings.HasPrefix(jsonStr, "{") || !strings.HasSuffix(jsonStr, "}") {
		return "", false
	}

	// Try to parse as JSON to ensure it's valid
	var temp interface{}
	if err := json.Unmarshal([]byte(jsonStr), &temp); err != nil {
		return "", false
	}

	return jsonStr, true
}

// ValidateUIUpdateJSON validates if a JSON string is a valid UI update
func ValidateUIUpdateJSON(jsonStr string) error {
	// This is a basic validation - in a real implementation,
	// you'd use a proper JSON schema validator

	required := []string{
		`"type"`,
		`"action"`,
		`"data"`,
		`"explanation"`,
	}

	for _, req := range required {
		if !strings.Contains(jsonStr, req) {
			return fmt.Errorf("missing required field: %s", req)
		}
	}

	// Check for valid action types
	validActions := []string{
		`"update_component"`,
		`"add_component"`,
		`"remove_component"`,
		`"update_theme"`,
		`"update_layout"`,
	}

	hasValidAction := false
	for _, action := range validActions {
		if strings.Contains(jsonStr, action) {
			hasValidAction = true
			break
		}
	}

	if !hasValidAction {
		return fmt.Errorf("invalid action type")
	}

	// Validate component types for add_component action
	if strings.Contains(jsonStr, `"add_component"`) {
		validComponentTypes := []string{
			"button", "text", "animated-text", "t-animated",
			"image", "input", "card", "form",
			"navigation", "hero", "gallery",
		}

		hasValidComponentId := false
		for _, componentType := range validComponentTypes {
			// Check if the JSON contains this component type in the type field (not id field)
			// This allows for "type": "hero", "type":"hero", etc. (flexible with spaces)
			if strings.Contains(jsonStr, `"type":"`+componentType+`"`) ||
				strings.Contains(jsonStr, `"type": "`+componentType+`"`) {
				hasValidComponentId = true
				break
			}
		}

		if !hasValidComponentId {
			return fmt.Errorf("invalid component type for add_component action - type must be one of: %v", validComponentTypes)
		}

		// Validate component categories (atom, molecule, organism)
		validComponentCategories := []string{
			`"atom"`, `"molecule"`, `"organism"`,
		}

		hasValidComponentCategory := false
		for _, category := range validComponentCategories {
			if strings.Contains(jsonStr, `"category":`+category) ||
				strings.Contains(jsonStr, `"category": `+category) {
				hasValidComponentCategory = true
				break
			}
		}

		if !hasValidComponentCategory {
			return fmt.Errorf("invalid component category for add_component action - category must be one of: atom, molecule, organism")
		}
	}

	// Validate theme colors for update_theme action
	if strings.Contains(jsonStr, `"update_theme"`) {
		// Check for valid hex colors (basic validation)
		if strings.Contains(jsonStr, `"colors"`) {
			// This is a simplified check - in production you'd want proper color validation
			if !strings.Contains(jsonStr, "#") {
				return fmt.Errorf("theme colors should contain valid hex color codes")
			}
		}
	}

	// Validate layout properties for update_layout action
	if strings.Contains(jsonStr, `"update_layout"`) {
		layoutProperties := []string{"grid", "rowHeight", "margin", "containerPadding"}
		hasLayoutProperty := false

		for _, prop := range layoutProperties {
			if strings.Contains(jsonStr, `"`+prop+`"`) {
				hasLayoutProperty = true
				break
			}
		}

		if !hasLayoutProperty {
			return fmt.Errorf("update_layout action should contain at least one valid layout property")
		}
	}

	return nil
}

// ValidateComponentProperties validates component-specific properties
func ValidateComponentProperties(componentId string, jsonStr string) error {
	switch componentId {
	case "button":
		requiredProps := []string{"text"}
		for _, prop := range requiredProps {
			if !strings.Contains(jsonStr, `"`+prop+`"`) {
				return fmt.Errorf("button component requires %s property", prop)
			}
		}
	case "text":
		requiredProps := []string{"content"}
		for _, prop := range requiredProps {
			if !strings.Contains(jsonStr, `"`+prop+`"`) {
				return fmt.Errorf("text component requires %s property", prop)
			}
		}
	case "animated-text":
		requiredProps := []string{"text", "effect"}
		for _, prop := range requiredProps {
			if !strings.Contains(jsonStr, `"`+prop+`"`) {
				return fmt.Errorf("animated-text component requires %s property", prop)
			}
		}
	case "t-animated":
		requiredProps := []string{"id"}
		for _, prop := range requiredProps {
			if !strings.Contains(jsonStr, `"`+prop+`"`) {
				return fmt.Errorf("t-animated component requires %s property", prop)
			}
		}
	case "image":
		requiredProps := []string{"src", "alt"}
		for _, prop := range requiredProps {
			if !strings.Contains(jsonStr, `"`+prop+`"`) {
				return fmt.Errorf("image component requires %s property", prop)
			}
		}
	case "input":
		requiredProps := []string{"placeholder", "type"}
		for _, prop := range requiredProps {
			if !strings.Contains(jsonStr, `"`+prop+`"`) {
				return fmt.Errorf("input component requires %s property", prop)
			}
		}
	case "card":
		requiredProps := []string{"title", "content"}
		for _, prop := range requiredProps {
			if !strings.Contains(jsonStr, `"`+prop+`"`) {
				return fmt.Errorf("card component requires %s property", prop)
			}
		}
	case "form":
		requiredProps := []string{"title", "fields", "submitText"}
		for _, prop := range requiredProps {
			if !strings.Contains(jsonStr, `"`+prop+`"`) {
				return fmt.Errorf("form component requires %s property", prop)
			}
		}
	case "navigation":
		requiredProps := []string{"items"}
		for _, prop := range requiredProps {
			if !strings.Contains(jsonStr, `"`+prop+`"`) {
				return fmt.Errorf("navigation component requires %s property", prop)
			}
		}
	case "hero":
		requiredProps := []string{"title", "subtitle", "ctaText"}
		for _, prop := range requiredProps {
			if !strings.Contains(jsonStr, `"`+prop+`"`) {
				return fmt.Errorf("hero component requires %s property", prop)
			}
		}
	case "gallery":
		requiredProps := []string{"title", "images", "columns"}
		for _, prop := range requiredProps {
			if !strings.Contains(jsonStr, `"`+prop+`"`) {
				return fmt.Errorf("gallery component requires %s property", prop)
			}
		}
	}

	return nil
}
