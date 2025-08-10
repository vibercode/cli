package storage

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/vibercode/cli/internal/models"
	"github.com/vibercode/cli/internal/vibe/prompts"
)

// GraphStorage handles graph operations with Neo4j
type GraphStorage struct {
	driver neo4j.DriverWithContext
	db     string
}

// ComponentNode represents a component node in the graph
type ComponentNode struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Category    string                 `json:"category"`
	Name        string                 `json:"name"`
	Properties  map[string]interface{} `json:"properties"`
	ProjectID   string                 `json:"project_id"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	Position    prompts.Position       `json:"position"`
	Size        prompts.Size           `json:"size"`
	Tags        []string               `json:"tags"`
	Description string                 `json:"description"`
}

// SchemaNode represents a schema node in the graph
type SchemaNode struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	ProjectID   string                 `json:"project_id"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	Fields      []models.SchemaField   `json:"fields"`
	Tags        []string               `json:"tags"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// ConversationNode represents a conversation node in the graph
type ConversationNode struct {
	ID        string                 `json:"id"`
	SessionID string                 `json:"session_id"`
	Content   string                 `json:"content"`
	Context   string                 `json:"context"`
	Intent    string                 `json:"intent"`
	ProjectID string                 `json:"project_id"`
	CreatedAt time.Time              `json:"created_at"`
	Tags      []string               `json:"tags"`
	Metadata  map[string]interface{} `json:"metadata"`
}

// ProjectNode represents a project node in the graph
type ProjectNode struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Theme       string                 `json:"theme"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// Relationship represents a relationship between nodes
type Relationship struct {
	FromID     string                 `json:"from_id"`
	ToID       string                 `json:"to_id"`
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	CreatedAt  time.Time              `json:"created_at"`
	Weight     float64                `json:"weight"`
	Direction  string                 `json:"direction"` // "incoming", "outgoing", "bidirectional"
}

// GraphPath represents a path in the graph
type GraphPath struct {
	Nodes         []interface{}  `json:"nodes"`
	Relationships []Relationship `json:"relationships"`
	Length        int            `json:"length"`
	Score         float64        `json:"score"`
}

// NewGraphStorage creates a new graph storage instance
func NewGraphStorage(neo4jURI string, username string, password string, database string) (*GraphStorage, error) {
	driver, err := neo4j.NewDriverWithContext(neo4jURI, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, fmt.Errorf("failed to create Neo4j driver: %w", err)
	}

	// Verify connectivity
	ctx := context.Background()
	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to verify Neo4j connectivity: %w", err)
	}

	storage := &GraphStorage{
		driver: driver,
		db:     database,
	}

	// Initialize constraints and indexes
	if err := storage.initializeSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize graph schema: %w", err)
	}

	return storage, nil
}

// initializeSchema creates necessary constraints and indexes
func (gs *GraphStorage) initializeSchema() error {
	ctx := context.Background()
	session := gs.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: gs.db})
	defer session.Close(ctx)

	// Create constraints
	constraints := []string{
		"CREATE CONSTRAINT component_id IF NOT EXISTS FOR (c:Component) REQUIRE c.id IS UNIQUE",
		"CREATE CONSTRAINT schema_id IF NOT EXISTS FOR (s:Schema) REQUIRE s.id IS UNIQUE",
		"CREATE CONSTRAINT conversation_id IF NOT EXISTS FOR (c:Conversation) REQUIRE c.id IS UNIQUE",
		"CREATE CONSTRAINT project_id IF NOT EXISTS FOR (p:Project) REQUIRE p.id IS UNIQUE",
	}

	for _, constraint := range constraints {
		_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
			return tx.Run(ctx, constraint, nil)
		})
		if err != nil {
			log.Printf("⚠️ Warning: Failed to create constraint: %v", err)
		}
	}

	// Create indexes
	indexes := []string{
		"CREATE INDEX component_type IF NOT EXISTS FOR (c:Component) ON (c.type)",
		"CREATE INDEX component_category IF NOT EXISTS FOR (c:Component) ON (c.category)",
		"CREATE INDEX component_project IF NOT EXISTS FOR (c:Component) ON (c.project_id)",
		"CREATE INDEX schema_name IF NOT EXISTS FOR (s:Schema) ON (s.name)",
		"CREATE INDEX schema_project IF NOT EXISTS FOR (s:Schema) ON (s.project_id)",
		"CREATE INDEX conversation_session IF NOT EXISTS FOR (c:Conversation) ON (c.session_id)",
		"CREATE INDEX conversation_project IF NOT EXISTS FOR (c:Conversation) ON (c.project_id)",
		"CREATE INDEX project_name IF NOT EXISTS FOR (p:Project) ON (p.name)",
	}

	for _, index := range indexes {
		_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
			return tx.Run(ctx, index, nil)
		})
		if err != nil {
			log.Printf("⚠️ Warning: Failed to create index: %v", err)
		}
	}

	log.Printf("✅ Graph schema initialized")
	return nil
}

// StoreComponent stores a component node in the graph
func (gs *GraphStorage) StoreComponent(ctx context.Context, component *ComponentNode) error {
	session := gs.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: gs.db})
	defer session.Close(ctx)

	query := `
		MERGE (c:Component {id: $id})
		SET c += $properties
		SET c.updated_at = $updated_at
		RETURN c
	`

	params := map[string]interface{}{
		"id":         component.ID,
		"updated_at": component.UpdatedAt,
		"properties": map[string]interface{}{
			"type":        component.Type,
			"category":    component.Category,
			"name":        component.Name,
			"properties":  component.Properties,
			"project_id":  component.ProjectID,
			"created_at":  component.CreatedAt,
			"position":    component.Position,
			"size":        component.Size,
			"tags":        component.Tags,
			"description": component.Description,
		},
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return tx.Run(ctx, query, params)
	})

	if err != nil {
		return fmt.Errorf("failed to store component: %w", err)
	}

	log.Printf("✅ Stored component node: %s", component.ID)
	return nil
}

// StoreSchema stores a schema node in the graph
func (gs *GraphStorage) StoreSchema(ctx context.Context, schema *SchemaNode) error {
	session := gs.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: gs.db})
	defer session.Close(ctx)

	query := `
		MERGE (s:Schema {id: $id})
		SET s += $properties
		SET s.updated_at = $updated_at
		RETURN s
	`

	params := map[string]interface{}{
		"id":         schema.ID,
		"updated_at": schema.UpdatedAt,
		"properties": map[string]interface{}{
			"name":        schema.Name,
			"description": schema.Description,
			"project_id":  schema.ProjectID,
			"created_at":  schema.CreatedAt,
			"fields":      schema.Fields,
			"tags":        schema.Tags,
			"metadata":    schema.Metadata,
		},
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return tx.Run(ctx, query, params)
	})

	if err != nil {
		return fmt.Errorf("failed to store schema: %w", err)
	}

	log.Printf("✅ Stored schema node: %s", schema.ID)
	return nil
}

// StoreConversation stores a conversation node in the graph
func (gs *GraphStorage) StoreConversation(ctx context.Context, conversation *ConversationNode) error {
	session := gs.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: gs.db})
	defer session.Close(ctx)

	query := `
		MERGE (c:Conversation {id: $id})
		SET c += $properties
		RETURN c
	`

	params := map[string]interface{}{
		"id": conversation.ID,
		"properties": map[string]interface{}{
			"session_id": conversation.SessionID,
			"content":    conversation.Content,
			"context":    conversation.Context,
			"intent":     conversation.Intent,
			"project_id": conversation.ProjectID,
			"created_at": conversation.CreatedAt,
			"tags":       conversation.Tags,
			"metadata":   conversation.Metadata,
		},
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return tx.Run(ctx, query, params)
	})

	if err != nil {
		return fmt.Errorf("failed to store conversation: %w", err)
	}

	log.Printf("✅ Stored conversation node: %s", conversation.ID)
	return nil
}

// StoreProject stores a project node in the graph
func (gs *GraphStorage) StoreProject(ctx context.Context, project *ProjectNode) error {
	session := gs.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: gs.db})
	defer session.Close(ctx)

	query := `
		MERGE (p:Project {id: $id})
		SET p += $properties
		SET p.updated_at = $updated_at
		RETURN p
	`

	params := map[string]interface{}{
		"id":         project.ID,
		"updated_at": project.UpdatedAt,
		"properties": map[string]interface{}{
			"name":        project.Name,
			"description": project.Description,
			"theme":       project.Theme,
			"created_at":  project.CreatedAt,
			"metadata":    project.Metadata,
		},
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return tx.Run(ctx, query, params)
	})

	if err != nil {
		return fmt.Errorf("failed to store project: %w", err)
	}

	log.Printf("✅ Stored project node: %s", project.ID)
	return nil
}

// CreateRelationship creates a relationship between two nodes
func (gs *GraphStorage) CreateRelationship(ctx context.Context, relationship *Relationship) error {
	session := gs.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: gs.db})
	defer session.Close(ctx)

	query := fmt.Sprintf(`
		MATCH (from {id: $from_id})
		MATCH (to {id: $to_id})
		MERGE (from)-[r:%s]->(to)
		SET r += $properties
		SET r.created_at = $created_at
		SET r.weight = $weight
		RETURN r
	`, strings.ToUpper(relationship.Type))

	params := map[string]interface{}{
		"from_id":    relationship.FromID,
		"to_id":      relationship.ToID,
		"properties": relationship.Properties,
		"created_at": relationship.CreatedAt,
		"weight":     relationship.Weight,
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return tx.Run(ctx, query, params)
	})

	if err != nil {
		return fmt.Errorf("failed to create relationship: %w", err)
	}

	log.Printf("✅ Created relationship: %s -[%s]-> %s", relationship.FromID, relationship.Type, relationship.ToID)
	return nil
}

// FindRelatedComponents finds components related to a given component
func (gs *GraphStorage) FindRelatedComponents(ctx context.Context, componentID string, projectID string, maxDepth int) ([]*ComponentNode, error) {
	session := gs.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: gs.db})
	defer session.Close(ctx)

	query := `
		MATCH (c:Component {id: $component_id})
		MATCH path = (c)-[*1..$max_depth]-(related:Component)
		WHERE related.project_id = $project_id
		RETURN DISTINCT related
		ORDER BY length(path) ASC
	`

	params := map[string]interface{}{
		"component_id": componentID,
		"project_id":   projectID,
		"max_depth":    maxDepth,
	}

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return tx.Run(ctx, query, params)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to find related components: %w", err)
	}

	// Cast to ResultWithContext
	resultWithContext, ok := result.(neo4j.ResultWithContext)
	if !ok {
		return nil, fmt.Errorf("unexpected result type")
	}

	var components []*ComponentNode

	for resultWithContext.Next(ctx) {
		record := resultWithContext.Record()
		nodeValue, found := record.Get("related")
		if !found {
			continue
		}

		node := nodeValue.(neo4j.Node)
		component := &ComponentNode{
			ID:          node.Props["id"].(string),
			Type:        node.Props["type"].(string),
			Category:    node.Props["category"].(string),
			Name:        node.Props["name"].(string),
			ProjectID:   node.Props["project_id"].(string),
			Description: node.Props["description"].(string),
		}

		if createdAt, ok := node.Props["created_at"].(time.Time); ok {
			component.CreatedAt = createdAt
		}

		if updatedAt, ok := node.Props["updated_at"].(time.Time); ok {
			component.UpdatedAt = updatedAt
		}

		components = append(components, component)
	}

	log.Printf("✅ Found %d related components for %s", len(components), componentID)
	return components, nil
}

// FindSchemaComponentRelations finds relationships between schemas and components
func (gs *GraphStorage) FindSchemaComponentRelations(ctx context.Context, schemaID string, projectID string) ([]*ComponentNode, error) {
	session := gs.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: gs.db})
	defer session.Close(ctx)

	query := `
		MATCH (s:Schema {id: $schema_id})
		MATCH (s)-[r:GENERATES|USES|REFERENCES]-(c:Component)
		WHERE c.project_id = $project_id
		RETURN DISTINCT c, r
		ORDER BY r.weight DESC
	`

	params := map[string]interface{}{
		"schema_id":  schemaID,
		"project_id": projectID,
	}

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return tx.Run(ctx, query, params)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to find schema component relations: %w", err)
	}

	// Cast to ResultWithContext
	resultWithContext, ok := result.(neo4j.ResultWithContext)
	if !ok {
		return nil, fmt.Errorf("unexpected result type")
	}

	var components []*ComponentNode

	for resultWithContext.Next(ctx) {
		record := resultWithContext.Record()
		nodeValue, found := record.Get("c")
		if !found {
			continue
		}

		node := nodeValue.(neo4j.Node)
		component := &ComponentNode{
			ID:          node.Props["id"].(string),
			Type:        node.Props["type"].(string),
			Category:    node.Props["category"].(string),
			Name:        node.Props["name"].(string),
			ProjectID:   node.Props["project_id"].(string),
			Description: node.Props["description"].(string),
		}

		components = append(components, component)
	}

	log.Printf("✅ Found %d schema-component relations for %s", len(components), schemaID)
	return components, nil
}

// FindConversationContext finds conversation context for a given topic
func (gs *GraphStorage) FindConversationContext(ctx context.Context, intent string, projectID string, limit int) ([]*ConversationNode, error) {
	session := gs.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: gs.db})
	defer session.Close(ctx)

	query := `
		MATCH (c:Conversation)
		WHERE c.intent = $intent OR c.content CONTAINS $intent
		AND c.project_id = $project_id
		RETURN c
		ORDER BY c.created_at DESC
		LIMIT $limit
	`

	params := map[string]interface{}{
		"intent":     intent,
		"project_id": projectID,
		"limit":      limit,
	}

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return tx.Run(ctx, query, params)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to find conversation context: %w", err)
	}

	// Cast to ResultWithContext
	resultWithContext, ok := result.(neo4j.ResultWithContext)
	if !ok {
		return nil, fmt.Errorf("unexpected result type")
	}

	var conversations []*ConversationNode

	for resultWithContext.Next(ctx) {
		record := resultWithContext.Record()
		nodeValue, found := record.Get("c")
		if !found {
			continue
		}

		node := nodeValue.(neo4j.Node)
		conversation := &ConversationNode{
			ID:        node.Props["id"].(string),
			SessionID: node.Props["session_id"].(string),
			Content:   node.Props["content"].(string),
			Context:   node.Props["context"].(string),
			Intent:    node.Props["intent"].(string),
			ProjectID: node.Props["project_id"].(string),
		}

		if createdAt, ok := node.Props["created_at"].(time.Time); ok {
			conversation.CreatedAt = createdAt
		}

		conversations = append(conversations, conversation)
	}

	log.Printf("✅ Found %d conversation contexts for intent: %s", len(conversations), intent)
	return conversations, nil
}

// FindShortestPath finds the shortest path between two nodes
func (gs *GraphStorage) FindShortestPath(ctx context.Context, fromID string, toID string, maxLength int) (*GraphPath, error) {
	session := gs.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: gs.db})
	defer session.Close(ctx)

	query := `
		MATCH (from {id: $from_id})
		MATCH (to {id: $to_id})
		MATCH path = shortestPath((from)-[*1..$max_length]-(to))
		RETURN path
	`

	params := map[string]interface{}{
		"from_id":    fromID,
		"to_id":      toID,
		"max_length": maxLength,
	}

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return tx.Run(ctx, query, params)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to find shortest path: %w", err)
	}

	// Cast to ResultWithContext
	resultWithContext, ok := result.(neo4j.ResultWithContext)
	if !ok {
		return nil, fmt.Errorf("unexpected result type")
	}

	if !resultWithContext.Next(ctx) {
		return nil, fmt.Errorf("no path found between %s and %s", fromID, toID)
	}

	record := resultWithContext.Record()
	pathValue, found := record.Get("path")
	if !found {
		return nil, fmt.Errorf("path not found in result")
	}

	path := pathValue.(neo4j.Path)
	graphPath := &GraphPath{
		Length: len(path.Relationships),
		Score:  1.0 / float64(len(path.Relationships)+1), // Higher score for shorter paths
	}

	// Extract nodes and relationships
	for _, node := range path.Nodes {
		graphPath.Nodes = append(graphPath.Nodes, map[string]interface{}{
			"id":     node.Props["id"],
			"labels": node.Labels,
			"props":  node.Props,
		})
	}

	for _, rel := range path.Relationships {
		relationship := Relationship{
			Type:       rel.Type,
			Properties: rel.Props,
		}
		graphPath.Relationships = append(graphPath.Relationships, relationship)
	}

	log.Printf("✅ Found shortest path between %s and %s (length: %d)", fromID, toID, graphPath.Length)
	return graphPath, nil
}

// GetProjectGraph returns the complete graph for a project
func (gs *GraphStorage) GetProjectGraph(ctx context.Context, projectID string) (map[string]interface{}, error) {
	session := gs.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: gs.db})
	defer session.Close(ctx)

	query := `
		MATCH (n)
		WHERE n.project_id = $project_id
		OPTIONAL MATCH (n)-[r]-(m)
		WHERE m.project_id = $project_id
		RETURN n, r, m
	`

	params := map[string]interface{}{
		"project_id": projectID,
	}

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		return tx.Run(ctx, query, params)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get project graph: %w", err)
	}

	// Cast to ResultWithContext
	resultWithContext, ok := result.(neo4j.ResultWithContext)
	if !ok {
		return nil, fmt.Errorf("unexpected result type")
	}

	nodes := make(map[string]interface{})
	relationships := make([]interface{}, 0)

	for resultWithContext.Next(ctx) {
		record := resultWithContext.Record()

		// Process node
		if nodeValue, found := record.Get("n"); found {
			node := nodeValue.(neo4j.Node)
			nodes[node.Props["id"].(string)] = map[string]interface{}{
				"id":     node.Props["id"],
				"labels": node.Labels,
				"props":  node.Props,
			}
		}

		// Process relationship
		if relValue, found := record.Get("r"); found && relValue != nil {
			rel := relValue.(neo4j.Relationship)
			relationships = append(relationships, map[string]interface{}{
				"type":  rel.Type,
				"start": rel.StartId,
				"end":   rel.EndId,
				"props": rel.Props,
			})
		}
	}

	graph := map[string]interface{}{
		"nodes":         nodes,
		"relationships": relationships,
		"stats": map[string]interface{}{
			"node_count":         len(nodes),
			"relationship_count": len(relationships),
		},
	}

	log.Printf("✅ Retrieved project graph for %s: %d nodes, %d relationships", projectID, len(nodes), len(relationships))
	return graph, nil
}

// Close closes the graph storage connection
func (gs *GraphStorage) Close(ctx context.Context) error {
	return gs.driver.Close(ctx)
}
