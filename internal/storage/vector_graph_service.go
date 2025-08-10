package storage

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/vibercode/cli/internal/models"
	"github.com/vibercode/cli/internal/vibe"
	"github.com/vibercode/cli/internal/vibe/prompts"
)

// VectorGraphService integrates vector storage (Qdrant) with graph storage (Neo4j)
type VectorGraphService struct {
	vectorStorage    *VectorStorage
	graphStorage     *GraphStorage
	embeddingService *EmbeddingService
	projectID        string
	enabled          bool
}

// VectorGraphConfig holds configuration for the vector graph service
type VectorGraphConfig struct {
	QdrantURL     string
	QdrantEnabled bool
	Neo4jURI      string
	Neo4jUsername string
	Neo4jPassword string
	Neo4jDatabase string
	Neo4jEnabled  bool
	EmbeddingSize int
}

// SemanticSearchResult combines vector and graph search results
type SemanticSearchResult struct {
	VectorResults []*SearchResult  `json:"vector_results"`
	GraphResults  []*ComponentNode `json:"graph_results"`
	Combined      []interface{}    `json:"combined"`
	Query         string           `json:"query"`
	Score         float64          `json:"score"`
	Explanation   string           `json:"explanation"`
	Timestamp     time.Time        `json:"timestamp"`
}

// RelationshipInsight represents insights from graph relationships
type RelationshipInsight struct {
	Type        string                 `json:"type"`
	Source      string                 `json:"source"`
	Target      string                 `json:"target"`
	Weight      float64                `json:"weight"`
	Properties  map[string]interface{} `json:"properties"`
	Explanation string                 `json:"explanation"`
	Confidence  float64                `json:"confidence"`
}

// NewVectorGraphService creates a new vector graph service
func NewVectorGraphService(config *VectorGraphConfig, claudeClient *vibe.ClaudeClient, projectID string) (*VectorGraphService, error) {
	service := &VectorGraphService{
		projectID: projectID,
		enabled:   config.QdrantEnabled || config.Neo4jEnabled,
	}

	// Initialize vector storage if enabled
	if config.QdrantEnabled {
		vectorStorage, err := NewVectorStorage(config.QdrantURL, "vibercode_vectors", config.EmbeddingSize)
		if err != nil {
			log.Printf("⚠️ Failed to initialize vector storage: %v", err)
		} else {
			service.vectorStorage = vectorStorage
			log.Printf("✅ Vector storage initialized")
		}
	}

	// Initialize graph storage if enabled
	if config.Neo4jEnabled {
		graphStorage, err := NewGraphStorage(config.Neo4jURI, config.Neo4jUsername, config.Neo4jPassword, config.Neo4jDatabase)
		if err != nil {
			log.Printf("⚠️ Failed to initialize graph storage: %v", err)
		} else {
			service.graphStorage = graphStorage
			log.Printf("✅ Graph storage initialized")
		}
	}

	// Initialize embedding service
	if claudeClient != nil {
		service.embeddingService = NewEmbeddingService(claudeClient, config.EmbeddingSize)
		log.Printf("✅ Embedding service initialized")
	}

	return service, nil
}

// StoreComponent stores a component in both vector and graph storage
func (vgs *VectorGraphService) StoreComponent(ctx context.Context, component *prompts.ComponentState) error {
	if !vgs.enabled {
		return nil
	}

	// Generate embedding if embedding service is available
	if vgs.embeddingService != nil {
		embedding, err := vgs.embeddingService.GenerateComponentEmbedding(ctx, component)
		if err != nil {
			log.Printf("⚠️ Failed to generate component embedding: %v", err)
		} else {
			// Store in vector storage
			if vgs.vectorStorage != nil {
				componentVector := &ComponentVector{
					ID:          component.ID,
					Type:        component.Type,
					Category:    component.Category,
					Properties:  component.Properties,
					Description: fmt.Sprintf("%s - %s", component.Name, component.Type),
					Tags:        []string{component.Type, component.Category},
					ProjectID:   vgs.projectID,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					Embedding:   embedding.Vector,
					Position:    component.Position,
					Size:        component.Size,
				}

				if err := vgs.vectorStorage.StoreComponentVector(ctx, componentVector); err != nil {
					log.Printf("⚠️ Failed to store component vector: %v", err)
				}
			}
		}
	}

	// Store in graph storage
	if vgs.graphStorage != nil {
		componentNode := &ComponentNode{
			ID:          component.ID,
			Type:        component.Type,
			Category:    component.Category,
			Name:        component.Name,
			Properties:  component.Properties,
			ProjectID:   vgs.projectID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Position:    component.Position,
			Size:        component.Size,
			Tags:        []string{component.Type, component.Category},
			Description: fmt.Sprintf("%s component", component.Type),
		}

		if err := vgs.graphStorage.StoreComponent(ctx, componentNode); err != nil {
			log.Printf("⚠️ Failed to store component node: %v", err)
		}
	}

	log.Printf("✅ Stored component in vector/graph storage: %s", component.ID)
	return nil
}

// StoreSchema stores a schema in both vector and graph storage
func (vgs *VectorGraphService) StoreSchema(ctx context.Context, schema *models.ResourceSchema) error {
	if !vgs.enabled {
		return nil
	}

	// Generate embedding if embedding service is available
	if vgs.embeddingService != nil {
		embedding, err := vgs.embeddingService.GenerateSchemaEmbedding(ctx, schema)
		if err != nil {
			log.Printf("⚠️ Failed to generate schema embedding: %v", err)
		} else {
			// Store in vector storage
			if vgs.vectorStorage != nil {
				schemaVector := &SchemaVector{
					ID:          schema.ID,
					Name:        schema.Name,
					Description: schema.Description,
					Fields:      schema.Fields,
					Tags:        []string{"schema", schema.Name},
					ProjectID:   vgs.projectID,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					Embedding:   embedding.Vector,
					Metadata:    schema.Metadata,
				}

				if err := vgs.vectorStorage.StoreSchemaVector(ctx, schemaVector); err != nil {
					log.Printf("⚠️ Failed to store schema vector: %v", err)
				}
			}
		}
	}

	// Store in graph storage
	if vgs.graphStorage != nil {
		schemaNode := &SchemaNode{
			ID:          schema.ID,
			Name:        schema.Name,
			Description: schema.Description,
			ProjectID:   vgs.projectID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Fields:      schema.Fields,
			Tags:        []string{"schema", schema.Name},
			Metadata:    schema.Metadata,
		}

		if err := vgs.graphStorage.StoreSchema(ctx, schemaNode); err != nil {
			log.Printf("⚠️ Failed to store schema node: %v", err)
		}
	}

	log.Printf("✅ Stored schema in vector/graph storage: %s", schema.ID)
	return nil
}

// StoreConversation stores a conversation in both vector and graph storage
func (vgs *VectorGraphService) StoreConversation(ctx context.Context, conversation *vibe.Message, sessionID string) error {
	if !vgs.enabled {
		return nil
	}

	// Generate embedding if embedding service is available
	if vgs.embeddingService != nil {
		embedding, err := vgs.embeddingService.GenerateConversationEmbedding(ctx, conversation)
		if err != nil {
			log.Printf("⚠️ Failed to generate conversation embedding: %v", err)
		} else {
			// Store in vector storage
			if vgs.vectorStorage != nil {
				conversationVector := &ConversationVector{
					ID:        embedding.ID,
					SessionID: sessionID,
					Content:   conversation.Content,
					Context:   fmt.Sprintf("Role: %s, Type: %s", conversation.Role, conversation.Type),
					Intent:    vgs.extractIntent(conversation.Content),
					Tags:      []string{conversation.Role, conversation.Type},
					ProjectID: vgs.projectID,
					CreatedAt: conversation.Timestamp,
					Embedding: embedding.Vector,
					Metadata: map[string]interface{}{
						"role":      conversation.Role,
						"type":      conversation.Type,
						"timestamp": conversation.Timestamp,
					},
				}

				if err := vgs.vectorStorage.StoreConversationVector(ctx, conversationVector); err != nil {
					log.Printf("⚠️ Failed to store conversation vector: %v", err)
				}
			}
		}
	}

	// Store in graph storage
	if vgs.graphStorage != nil {
		conversationNode := &ConversationNode{
			ID:        uuid.New().String(),
			SessionID: sessionID,
			Content:   conversation.Content,
			Context:   fmt.Sprintf("Role: %s, Type: %s", conversation.Role, conversation.Type),
			Intent:    vgs.extractIntent(conversation.Content),
			ProjectID: vgs.projectID,
			CreatedAt: conversation.Timestamp,
			Tags:      []string{conversation.Role, conversation.Type},
			Metadata: map[string]interface{}{
				"role":      conversation.Role,
				"type":      conversation.Type,
				"timestamp": conversation.Timestamp,
			},
		}

		if err := vgs.graphStorage.StoreConversation(ctx, conversationNode); err != nil {
			log.Printf("⚠️ Failed to store conversation node: %v", err)
		}
	}

	log.Printf("✅ Stored conversation in vector/graph storage")
	return nil
}

// CreateComponentRelationship creates a relationship between components
func (vgs *VectorGraphService) CreateComponentRelationship(ctx context.Context, fromComponentID, toComponentID string, relationType string, weight float64) error {
	if !vgs.enabled || vgs.graphStorage == nil {
		return nil
	}

	relationship := &Relationship{
		FromID:    fromComponentID,
		ToID:      toComponentID,
		Type:      relationType,
		Weight:    weight,
		CreatedAt: time.Now(),
		Properties: map[string]interface{}{
			"auto_generated": true,
			"project_id":     vgs.projectID,
		},
	}

	if err := vgs.graphStorage.CreateRelationship(ctx, relationship); err != nil {
		return fmt.Errorf("failed to create component relationship: %w", err)
	}

	log.Printf("✅ Created component relationship: %s -[%s]-> %s", fromComponentID, relationType, toComponentID)
	return nil
}

// CreateSchemaComponentRelationship creates a relationship between schema and component
func (vgs *VectorGraphService) CreateSchemaComponentRelationship(ctx context.Context, schemaID, componentID string, relationType string, weight float64) error {
	if !vgs.enabled || vgs.graphStorage == nil {
		return nil
	}

	relationship := &Relationship{
		FromID:    schemaID,
		ToID:      componentID,
		Type:      relationType,
		Weight:    weight,
		CreatedAt: time.Now(),
		Properties: map[string]interface{}{
			"auto_generated": true,
			"project_id":     vgs.projectID,
		},
	}

	if err := vgs.graphStorage.CreateRelationship(ctx, relationship); err != nil {
		return fmt.Errorf("failed to create schema-component relationship: %w", err)
	}

	log.Printf("✅ Created schema-component relationship: %s -[%s]-> %s", schemaID, relationType, componentID)
	return nil
}

// SemanticSearch performs semantic search across both vector and graph storage
func (vgs *VectorGraphService) SemanticSearch(ctx context.Context, query string, limit int) (*SemanticSearchResult, error) {
	if !vgs.enabled {
		return &SemanticSearchResult{
			Query:     query,
			Timestamp: time.Now(),
		}, nil
	}

	result := &SemanticSearchResult{
		Query:     query,
		Timestamp: time.Now(),
		Combined:  make([]interface{}, 0),
	}

	// Search in vector storage
	if vgs.vectorStorage != nil && vgs.embeddingService != nil {
		queryEmbedding, err := vgs.embeddingService.SearchSimilarText(ctx, query, "query")
		if err != nil {
			log.Printf("⚠️ Failed to generate query embedding: %v", err)
		} else {
			// Search components
			componentResults, err := vgs.vectorStorage.SearchSimilarComponents(ctx, queryEmbedding.Vector, limit, vgs.projectID)
			if err != nil {
				log.Printf("⚠️ Failed to search components: %v", err)
			} else {
				result.VectorResults = append(result.VectorResults, componentResults...)
			}

			// Search schemas
			schemaResults, err := vgs.vectorStorage.SearchSimilarSchemas(ctx, queryEmbedding.Vector, limit, vgs.projectID)
			if err != nil {
				log.Printf("⚠️ Failed to search schemas: %v", err)
			} else {
				result.VectorResults = append(result.VectorResults, schemaResults...)
			}
		}
	}

	// Search in graph storage
	if vgs.graphStorage != nil {
		// Search conversations for context
		conversations, err := vgs.graphStorage.FindConversationContext(ctx, query, vgs.projectID, limit)
		if err != nil {
			log.Printf("⚠️ Failed to search conversation context: %v", err)
		} else {
			for _, conv := range conversations {
				result.GraphResults = append(result.GraphResults, &ComponentNode{
					ID:          conv.ID,
					Type:        "conversation",
					Name:        conv.Content,
					Description: conv.Context,
					ProjectID:   conv.ProjectID,
					CreatedAt:   conv.CreatedAt,
				})
			}
		}
	}

	// Combine results
	for _, vr := range result.VectorResults {
		result.Combined = append(result.Combined, vr)
	}
	for _, gr := range result.GraphResults {
		result.Combined = append(result.Combined, gr)
	}

	// Calculate overall score
	result.Score = vgs.calculateSemanticScore(result.VectorResults, result.GraphResults)
	result.Explanation = vgs.generateSearchExplanation(query, result)

	log.Printf("✅ Semantic search completed: %d vector results, %d graph results", len(result.VectorResults), len(result.GraphResults))
	return result, nil
}

// FindRelatedComponents finds components related to a given component
func (vgs *VectorGraphService) FindRelatedComponents(ctx context.Context, componentID string, maxDepth int) ([]*ComponentNode, error) {
	if !vgs.enabled || vgs.graphStorage == nil {
		return nil, nil
	}

	components, err := vgs.graphStorage.FindRelatedComponents(ctx, componentID, vgs.projectID, maxDepth)
	if err != nil {
		return nil, fmt.Errorf("failed to find related components: %w", err)
	}

	log.Printf("✅ Found %d related components for %s", len(components), componentID)
	return components, nil
}

// GetRelationshipInsights analyzes relationships in the graph
func (vgs *VectorGraphService) GetRelationshipInsights(ctx context.Context) ([]*RelationshipInsight, error) {
	if !vgs.enabled || vgs.graphStorage == nil {
		return nil, nil
	}

	// Get project graph
	graph, err := vgs.graphStorage.GetProjectGraph(ctx, vgs.projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project graph: %w", err)
	}

	var insights []*RelationshipInsight

	// Analyze relationships
	if relationships, ok := graph["relationships"].([]interface{}); ok {
		for _, rel := range relationships {
			if relMap, ok := rel.(map[string]interface{}); ok {
				insight := &RelationshipInsight{
					Type:   relMap["type"].(string),
					Weight: 0.8, // Default weight
					Properties: map[string]interface{}{
						"source": "graph_analysis",
					},
					Confidence: 0.75,
				}

				// Generate explanation
				insight.Explanation = fmt.Sprintf("Relationship of type '%s' detected in project graph", insight.Type)
				insights = append(insights, insight)
			}
		}
	}

	log.Printf("✅ Generated %d relationship insights", len(insights))
	return insights, nil
}

// GetProjectStats returns statistics about the project's vector and graph data
func (vgs *VectorGraphService) GetProjectStats(ctx context.Context) (map[string]interface{}, error) {
	if !vgs.enabled {
		return map[string]interface{}{"enabled": false}, nil
	}

	stats := map[string]interface{}{
		"enabled":   true,
		"timestamp": time.Now(),
	}

	// Vector storage stats
	if vgs.vectorStorage != nil {
		vectorInfo, err := vgs.vectorStorage.GetCollectionInfo(ctx)
		if err != nil {
			log.Printf("⚠️ Failed to get vector storage info: %v", err)
		} else {
			stats["vector_storage"] = vectorInfo
		}
	}

	// Graph storage stats
	if vgs.graphStorage != nil {
		graphInfo, err := vgs.graphStorage.GetProjectGraph(ctx, vgs.projectID)
		if err != nil {
			log.Printf("⚠️ Failed to get graph storage info: %v", err)
		} else {
			stats["graph_storage"] = graphInfo["stats"]
		}
	}

	// Embedding service stats
	if vgs.embeddingService != nil {
		stats["embedding_service"] = map[string]interface{}{
			"cache_size": vgs.embeddingService.GetCacheSize(),
		}
	}

	return stats, nil
}

// extractIntent extracts intent from conversation content
func (vgs *VectorGraphService) extractIntent(content string) string {
	// Simple intent extraction - in production, this could use NLP
	content = strings.ToLower(content)

	if strings.Contains(content, "agregar") || strings.Contains(content, "add") {
		return "add_component"
	} else if strings.Contains(content, "cambiar") || strings.Contains(content, "change") {
		return "modify_component"
	} else if strings.Contains(content, "eliminar") || strings.Contains(content, "delete") {
		return "delete_component"
	} else if strings.Contains(content, "tema") || strings.Contains(content, "theme") {
		return "change_theme"
	} else if strings.Contains(content, "ayuda") || strings.Contains(content, "help") {
		return "help"
	} else if strings.Contains(content, "estado") || strings.Contains(content, "status") {
		return "get_status"
	}

	return "general"
}

// calculateSemanticScore calculates a score for semantic search results
func (vgs *VectorGraphService) calculateSemanticScore(vectorResults []*SearchResult, graphResults []*ComponentNode) float64 {
	if len(vectorResults) == 0 && len(graphResults) == 0 {
		return 0.0
	}

	// Calculate average vector score
	vectorScore := 0.0
	if len(vectorResults) > 0 {
		for _, result := range vectorResults {
			vectorScore += float64(result.Score)
		}
		vectorScore /= float64(len(vectorResults))
	}

	// Graph results contribute to overall score
	graphScore := 0.0
	if len(graphResults) > 0 {
		graphScore = 0.7 // Base score for having graph results
	}

	// Combine scores with weights
	return (vectorScore * 0.6) + (graphScore * 0.4)
}

// generateSearchExplanation generates an explanation for search results
func (vgs *VectorGraphService) generateSearchExplanation(query string, result *SemanticSearchResult) string {
	if len(result.VectorResults) == 0 && len(result.GraphResults) == 0 {
		return fmt.Sprintf("No results found for query: '%s'", query)
	}

	explanation := fmt.Sprintf("Found %d vector results and %d graph results for query: '%s'",
		len(result.VectorResults), len(result.GraphResults), query)

	if len(result.VectorResults) > 0 {
		explanation += fmt.Sprintf(" Vector search found semantic matches with similarity scores.")
	}

	if len(result.GraphResults) > 0 {
		explanation += fmt.Sprintf(" Graph search found related components and conversations.")
	}

	return explanation
}

// Close closes all storage connections
func (vgs *VectorGraphService) Close(ctx context.Context) error {
	if vgs.graphStorage != nil {
		if err := vgs.graphStorage.Close(ctx); err != nil {
			log.Printf("⚠️ Failed to close graph storage: %v", err)
		}
	}

	log.Printf("✅ Vector graph service closed")
	return nil
}

// IsEnabled returns whether the service is enabled
func (vgs *VectorGraphService) IsEnabled() bool {
	return vgs.enabled
}
