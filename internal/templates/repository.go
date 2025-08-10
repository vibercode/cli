package templates

const RepositoryTemplate = `package repositories

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"{{.Module}}/internal/models"
)

// {{.NameVariations.Pascal}}Repository handles database operations for {{.NameVariations.PluralLower}}
type {{.NameVariations.Pascal}}Repository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// New{{.NameVariations.Pascal}}Repository creates a new {{.NameVariations.Pascal}}Repository
func New{{.NameVariations.Pascal}}Repository(db *mongo.Database) *{{.NameVariations.Pascal}}Repository {
	return &{{.NameVariations.Pascal}}Repository{
		db:         db,
		collection: db.Collection("{{.TableName}}"),
	}
}

// Create creates a new {{.NameVariations.Lower}}
func (r *{{.NameVariations.Pascal}}Repository) Create(ctx context.Context, {{.NameVariations.Camel}} *models.{{.NameVariations.Pascal}}) error {
	{{.NameVariations.Camel}}.ID = primitive.NewObjectID()
	{{.NameVariations.Camel}}.CreatedAt = time.Now()
	{{.NameVariations.Camel}}.UpdatedAt = time.Now()
	
	_, err := r.collection.InsertOne(ctx, {{.NameVariations.Camel}})
	return err
}

// GetByID retrieves a {{.NameVariations.Lower}} by ID
func (r *{{.NameVariations.Pascal}}Repository) GetByID(ctx context.Context, idStr string) (*models.{{.NameVariations.Pascal}}, error) {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}

	var {{.NameVariations.Camel}} models.{{.NameVariations.Pascal}}
	err = r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&{{.NameVariations.Camel}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &{{.NameVariations.Camel}}, nil
}

// GetAll retrieves all {{.NameVariations.PluralLower}} with pagination
func (r *{{.NameVariations.Pascal}}Repository) GetAll(ctx context.Context, limit, offset int) ([]models.{{.NameVariations.Pascal}}, error) {
	var {{.NameVariations.PluralCamel}} []models.{{.NameVariations.Pascal}}
	
	opts := options.Find()
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(offset))
	opts.SetSort(bson.D{{"created_at", -1}})
	
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	if err = cursor.All(ctx, &{{.NameVariations.PluralCamel}}); err != nil {
		return nil, err
	}
	
	return {{.NameVariations.PluralCamel}}, nil
}

// Update updates a {{.NameVariations.Lower}}
func (r *{{.NameVariations.Pascal}}Repository) Update(ctx context.Context, {{.NameVariations.Camel}} *models.{{.NameVariations.Pascal}}) error {
	{{.NameVariations.Camel}}.UpdatedAt = time.Now()
	
	filter := bson.M{"_id": {{.NameVariations.Camel}}.ID}
	update := bson.M{"$set": {{.NameVariations.Camel}}}
	
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete deletes a {{.NameVariations.Lower}} by ID
func (r *{{.NameVariations.Pascal}}Repository) Delete(ctx context.Context, idStr string) error {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}
	
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// Count returns the total number of {{.NameVariations.PluralLower}}
func (r *{{.NameVariations.Pascal}}Repository) Count(ctx context.Context) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	return count, err
}

// Search searches {{.NameVariations.PluralLower}} by query
func (r *{{.NameVariations.Pascal}}Repository) Search(ctx context.Context, query string, limit, offset int) ([]models.{{.NameVariations.Pascal}}, error) {
	var {{.NameVariations.PluralCamel}} []models.{{.NameVariations.Pascal}}
	
	// Build search filter
	filter := bson.M{}
	if query != "" {
		searchConditions := bson.A{}
		
		// Add search conditions for string fields
		{{- range .Fields}}
		{{- if or (eq .Type "string") (eq .Type "text")}}
		searchConditions = append(searchConditions, bson.M{"{{.Name | ToSnake}}": bson.M{"$regex": query, "$options": "i"}})
		{{- end}}
		{{- end}}
		
		if len(searchConditions) > 0 {
			filter["$or"] = searchConditions
		}
	}
	
	opts := options.Find()
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(offset))
	opts.SetSort(bson.D{{"created_at", -1}})
	
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	if err = cursor.All(ctx, &{{.NameVariations.PluralCamel}}); err != nil {
		return nil, err
	}
	
	return {{.NameVariations.PluralCamel}}, nil
}

// CountBySearch returns the count of {{.NameVariations.PluralLower}} matching the search query
func (r *{{.NameVariations.Pascal}}Repository) CountBySearch(ctx context.Context, query string) (int64, error) {
	filter := bson.M{}
	if query != "" {
		searchConditions := bson.A{}
		
		// Add search conditions for string fields
		{{- range .Fields}}
		{{- if or (eq .Type "string") (eq .Type "text")}}
		searchConditions = append(searchConditions, bson.M{"{{.Name | ToSnake}}": bson.M{"$regex": query, "$options": "i"}})
		{{- end}}
		{{- end}}
		
		if len(searchConditions) > 0 {
			filter["$or"] = searchConditions
		}
	}
	
	count, err := r.collection.CountDocuments(ctx, filter)
	return count, err
}

// GetByField retrieves {{.NameVariations.PluralLower}} by a specific field value
{{- range .Fields}}
{{- if or (eq .Type "string") (eq .Type "text")}}
func (r *{{$.NameVariations.Pascal}}Repository) GetBy{{.Name | ToCamel}}(ctx context.Context, {{.Name | ToLowerCamel}} {{.GoType}}) ([]models.{{$.NameVariations.Pascal}}, error) {
	var {{$.NameVariations.PluralCamel}} []models.{{$.NameVariations.Pascal}}
	
	filter := bson.M{"{{.Name | ToSnake}}": {{.Name | ToLowerCamel}}}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	if err = cursor.All(ctx, &{{$.NameVariations.PluralCamel}}); err != nil {
		return nil, err
	}
	
	return {{$.NameVariations.PluralCamel}}, nil
}
{{- end}}
{{- end}}

// Exists checks if a {{.NameVariations.Lower}} exists by ID
func (r *{{.NameVariations.Pascal}}Repository) Exists(ctx context.Context, idStr string) (bool, error) {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return false, fmt.Errorf("invalid ID format: %w", err)
	}
	
	count, err := r.collection.CountDocuments(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}
`