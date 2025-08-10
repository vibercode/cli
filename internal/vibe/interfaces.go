package vibe

import (
	"context"

	"github.com/vibercode/cli/internal/vibe/prompts"
)

// VectorGraphServiceInterface define las operaciones que necesita el preview server
type VectorGraphServiceInterface interface {
	IsEnabled() bool
	StoreComponent(ctx context.Context, component *prompts.ComponentState) error
	StoreConversation(ctx context.Context, message string, sessionID string) error
	CreateComponentRelationship(ctx context.Context, componentID, relatedID, relationshipType string, weight float64) error
	SemanticSearch(ctx context.Context, query string, limit int) (interface{}, error)
	FindRelatedComponents(ctx context.Context, componentID string, limit int) (interface{}, error)
	GetRelationshipInsights(ctx context.Context) (interface{}, error)
	GetProjectStats(ctx context.Context) (interface{}, error)
}
