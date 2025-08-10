package storage

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/vibercode/cli/internal/models"
	"github.com/vibercode/cli/internal/vibe"
	"github.com/vibercode/cli/internal/vibe/prompts"
)

// EmbeddingService generates embeddings for components, schemas, and conversations
type EmbeddingService struct {
	claudeClient *vibe.ClaudeClient
	vectorSize   int
	cache        map[string][]float32 // Simple cache for embeddings
}

// EmbeddingRequest represents a request for generating embeddings
type EmbeddingRequest struct {
	Text      string            `json:"text"`
	Type      string            `json:"type"` // "component", "schema", "conversation"
	ID        string            `json:"id"`
	ProjectID string            `json:"project_id"`
	Metadata  map[string]string `json:"metadata"`
}

// EmbeddingResponse represents the response from embedding generation
type EmbeddingResponse struct {
	ID        string    `json:"id"`
	Vector    []float32 `json:"vector"`
	Text      string    `json:"text"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	Hash      string    `json:"hash"`
}

// NewEmbeddingService creates a new embedding service
func NewEmbeddingService(claudeClient *vibe.ClaudeClient, vectorSize int) *EmbeddingService {
	return &EmbeddingService{
		claudeClient: claudeClient,
		vectorSize:   vectorSize,
		cache:        make(map[string][]float32),
	}
}

// GenerateComponentEmbedding generates an embedding for a component
func (es *EmbeddingService) GenerateComponentEmbedding(ctx context.Context, component *prompts.ComponentState) (*EmbeddingResponse, error) {
	// Create text representation of the component
	componentText := es.componentToText(component)

	// Generate hash for caching
	hash := es.generateHash(componentText)

	// Check cache first
	if cachedVector, exists := es.cache[hash]; exists {
		log.Printf("✅ Using cached embedding for component: %s", component.ID)
		return &EmbeddingResponse{
			ID:        component.ID,
			Vector:    cachedVector,
			Text:      componentText,
			Type:      "component",
			CreatedAt: time.Now(),
			Hash:      hash,
		}, nil
	}

	// Generate embedding using Claude
	vector, err := es.generateEmbedding(componentText, "component")
	if err != nil {
		return nil, fmt.Errorf("failed to generate component embedding: %w", err)
	}

	// Cache the result
	es.cache[hash] = vector

	response := &EmbeddingResponse{
		ID:        component.ID,
		Vector:    vector,
		Text:      componentText,
		Type:      "component",
		CreatedAt: time.Now(),
		Hash:      hash,
	}

	log.Printf("✅ Generated embedding for component: %s", component.ID)
	return response, nil
}

// GenerateSchemaEmbedding generates an embedding for a schema
func (es *EmbeddingService) GenerateSchemaEmbedding(ctx context.Context, schema *models.ResourceSchema) (*EmbeddingResponse, error) {
	// Create text representation of the schema
	schemaText := es.schemaToText(schema)

	// Generate hash for caching
	hash := es.generateHash(schemaText)

	// Check cache first
	if cachedVector, exists := es.cache[hash]; exists {
		log.Printf("✅ Using cached embedding for schema: %s", schema.ID)
		return &EmbeddingResponse{
			ID:        schema.ID,
			Vector:    cachedVector,
			Text:      schemaText,
			Type:      "schema",
			CreatedAt: time.Now(),
			Hash:      hash,
		}, nil
	}

	// Generate embedding using Claude
	vector, err := es.generateEmbedding(schemaText, "schema")
	if err != nil {
		return nil, fmt.Errorf("failed to generate schema embedding: %w", err)
	}

	// Cache the result
	es.cache[hash] = vector

	response := &EmbeddingResponse{
		ID:        schema.ID,
		Vector:    vector,
		Text:      schemaText,
		Type:      "schema",
		CreatedAt: time.Now(),
		Hash:      hash,
	}

	log.Printf("✅ Generated embedding for schema: %s", schema.ID)
	return response, nil
}

// GenerateConversationEmbedding generates an embedding for a conversation
func (es *EmbeddingService) GenerateConversationEmbedding(ctx context.Context, conversation *vibe.Message) (*EmbeddingResponse, error) {
	// Create text representation of the conversation
	conversationText := es.conversationToText(conversation)

	// Generate hash for caching
	hash := es.generateHash(conversationText)

	// Check cache first
	if cachedVector, exists := es.cache[hash]; exists {
		log.Printf("✅ Using cached embedding for conversation")
		return &EmbeddingResponse{
			ID:        hash, // Use hash as ID for conversations
			Vector:    cachedVector,
			Text:      conversationText,
			Type:      "conversation",
			CreatedAt: time.Now(),
			Hash:      hash,
		}, nil
	}

	// Generate embedding using Claude
	vector, err := es.generateEmbedding(conversationText, "conversation")
	if err != nil {
		return nil, fmt.Errorf("failed to generate conversation embedding: %w", err)
	}

	// Cache the result
	es.cache[hash] = vector

	response := &EmbeddingResponse{
		ID:        hash,
		Vector:    vector,
		Text:      conversationText,
		Type:      "conversation",
		CreatedAt: time.Now(),
		Hash:      hash,
	}

	log.Printf("✅ Generated embedding for conversation")
	return response, nil
}

// generateEmbedding generates an embedding using Claude API
func (es *EmbeddingService) generateEmbedding(text string, embeddingType string) ([]float32, error) {
	// Since Claude doesn't have a direct embedding endpoint, we'll use a creative approach
	// by asking Claude to generate a semantic fingerprint and convert it to a vector

	prompt := fmt.Sprintf(`
Analiza el siguiente %s y genera un vector de características semánticas de %d dimensiones.
Cada dimensión debe representar un aspecto específico del contenido:
- Dimensiones 1-50: Características estructurales
- Dimensiones 51-100: Características funcionales  
- Dimensiones 101-150: Características de contexto
- Dimensiones 151-200: Características de relación
- Dimensiones 201-256: Características semánticas

Contenido a analizar:
%s

Responde SOLO con un array JSON de %d números flotantes entre -1.0 y 1.0:
`, embeddingType, es.vectorSize, text, es.vectorSize)

	messages := []vibe.ClaudeMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	response, err := es.claudeClient.CreateMessage(messages)
	if err != nil {
		return nil, fmt.Errorf("failed to call Claude API: %w", err)
	}

	// Parse the response to extract the vector
	vector, err := es.parseVectorResponse(response)
	if err != nil {
		// If parsing fails, generate a deterministic vector based on text hash
		log.Printf("⚠️ Vector parsing failed, using deterministic fallback")
		return es.generateDeterministicVector(text), nil
	}

	return vector, nil
}

// parseVectorResponse parses Claude's response to extract the vector
func (es *EmbeddingService) parseVectorResponse(response string) ([]float32, error) {
	// Look for JSON array in the response
	start := strings.Index(response, "[")
	end := strings.LastIndex(response, "]")

	if start == -1 || end == -1 {
		return nil, fmt.Errorf("no JSON array found in response")
	}

	jsonStr := response[start : end+1]

	var floats []float64
	if err := json.Unmarshal([]byte(jsonStr), &floats); err != nil {
		return nil, fmt.Errorf("failed to parse JSON array: %w", err)
	}

	if len(floats) != es.vectorSize {
		return nil, fmt.Errorf("expected %d dimensions, got %d", es.vectorSize, len(floats))
	}

	// Convert to float32
	vector := make([]float32, len(floats))
	for i, f := range floats {
		vector[i] = float32(f)
	}

	return vector, nil
}

// generateDeterministicVector generates a deterministic vector based on text hash
func (es *EmbeddingService) generateDeterministicVector(text string) []float32 {
	// Create a deterministic but distributed vector based on text hash
	hash := sha256.Sum256([]byte(text))

	vector := make([]float32, es.vectorSize)
	for i := 0; i < es.vectorSize; i++ {
		// Use different parts of the hash to create variation
		byteIndex := i % len(hash)
		hashByte := hash[byteIndex]

		// Create a value between -1.0 and 1.0
		value := (float32(hashByte)/255.0)*2.0 - 1.0

		// Add some variation based on position
		variation := float32(i) / float32(es.vectorSize)
		vector[i] = value * (0.5 + variation*0.5)
	}

	return vector
}

// componentToText converts a component to text representation
func (es *EmbeddingService) componentToText(component *prompts.ComponentState) string {
	var parts []string

	parts = append(parts, fmt.Sprintf("Component Type: %s", component.Type))
	parts = append(parts, fmt.Sprintf("Category: %s", component.Category))
	parts = append(parts, fmt.Sprintf("Name: %s", component.Name))

	if component.Position.X != 0 || component.Position.Y != 0 {
		parts = append(parts, fmt.Sprintf("Position: (%d, %d)", component.Position.X, component.Position.Y))
	}

	if component.Size.W != 0 || component.Size.H != 0 {
		parts = append(parts, fmt.Sprintf("Size: %dx%d", component.Size.W, component.Size.H))
	}

	// Add properties
	if len(component.Properties) > 0 {
		propJSON, _ := json.Marshal(component.Properties)
		parts = append(parts, fmt.Sprintf("Properties: %s", string(propJSON)))
	}

	// Add style if present
	if len(component.Style) > 0 {
		styleJSON, _ := json.Marshal(component.Style)
		parts = append(parts, fmt.Sprintf("Style: %s", string(styleJSON)))
	}

	return strings.Join(parts, " | ")
}

// schemaToText converts a schema to text representation
func (es *EmbeddingService) schemaToText(schema *models.ResourceSchema) string {
	var parts []string

	parts = append(parts, fmt.Sprintf("Schema Name: %s", schema.Name))
	parts = append(parts, fmt.Sprintf("Description: %s", schema.Description))

	// Add fields information
	var fieldDescriptions []string
	for _, field := range schema.Fields {
		fieldDesc := fmt.Sprintf("%s (%s)", field.Name, field.Type)
		if field.Description != "" {
			fieldDesc += fmt.Sprintf(" - %s", field.Description)
		}
		fieldDescriptions = append(fieldDescriptions, fieldDesc)
	}

	parts = append(parts, fmt.Sprintf("Fields: %s", strings.Join(fieldDescriptions, ", ")))

	// Add metadata
	if schema.Metadata != nil {
		metadataJSON, _ := json.Marshal(schema.Metadata)
		parts = append(parts, fmt.Sprintf("Metadata: %s", string(metadataJSON)))
	}

	return strings.Join(parts, " | ")
}

// conversationToText converts a conversation to text representation
func (es *EmbeddingService) conversationToText(conversation *vibe.Message) string {
	var parts []string

	parts = append(parts, fmt.Sprintf("Role: %s", conversation.Role))
	parts = append(parts, fmt.Sprintf("Content: %s", conversation.Content))
	parts = append(parts, fmt.Sprintf("Type: %s", conversation.Type))
	parts = append(parts, fmt.Sprintf("Timestamp: %s", conversation.Timestamp.Format(time.RFC3339)))

	return strings.Join(parts, " | ")
}

// generateHash generates a hash for caching purposes
func (es *EmbeddingService) generateHash(text string) string {
	hash := sha256.Sum256([]byte(text))
	return fmt.Sprintf("%x", hash)
}

// BatchGenerateEmbeddings generates embeddings for multiple items
func (es *EmbeddingService) BatchGenerateEmbeddings(ctx context.Context, requests []*EmbeddingRequest) ([]*EmbeddingResponse, error) {
	var responses []*EmbeddingResponse

	for _, req := range requests {
		var response *EmbeddingResponse

		switch req.Type {
		case "component":
			// Note: This would need additional logic to convert request to ComponentState
			log.Printf("⚠️ Batch component embedding not yet implemented")
			continue
		case "schema":
			// Note: This would need additional logic to convert request to ResourceSchema
			log.Printf("⚠️ Batch schema embedding not yet implemented")
			continue
		case "conversation":
			// Note: This would need additional logic to convert request to Message
			log.Printf("⚠️ Batch conversation embedding not yet implemented")
			continue
		default:
			// Generic text embedding
			vector, err := es.generateEmbedding(req.Text, req.Type)
			if err != nil {
				log.Printf("❌ Failed to generate embedding for %s: %v", req.ID, err)
				continue
			}

			response = &EmbeddingResponse{
				ID:        req.ID,
				Vector:    vector,
				Text:      req.Text,
				Type:      req.Type,
				CreatedAt: time.Now(),
				Hash:      es.generateHash(req.Text),
			}
		}

		if response != nil {
			responses = append(responses, response)
		}
	}

	log.Printf("✅ Generated %d embeddings from %d requests", len(responses), len(requests))
	return responses, nil
}

// ClearCache clears the embedding cache
func (es *EmbeddingService) ClearCache() {
	es.cache = make(map[string][]float32)
	log.Printf("✅ Embedding cache cleared")
}

// GetCacheSize returns the current cache size
func (es *EmbeddingService) GetCacheSize() int {
	return len(es.cache)
}

// SearchSimilarText finds similar text embeddings (utility function)
func (es *EmbeddingService) SearchSimilarText(ctx context.Context, queryText string, textType string) (*EmbeddingResponse, error) {
	// Generate embedding for the query text
	vector, err := es.generateEmbedding(queryText, textType)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	return &EmbeddingResponse{
		ID:        es.generateHash(queryText),
		Vector:    vector,
		Text:      queryText,
		Type:      textType,
		CreatedAt: time.Now(),
		Hash:      es.generateHash(queryText),
	}, nil
}
