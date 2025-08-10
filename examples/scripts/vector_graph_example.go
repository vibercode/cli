package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/vibercode/cli/internal/models"
	"github.com/vibercode/cli/internal/storage"
	"github.com/vibercode/cli/internal/vibe"
	"github.com/vibercode/cli/internal/vibe/prompts"
)

func main() {
	fmt.Println("🚀 VibeCode Vector/Graph Storage Example")
	fmt.Println("=========================================")

	// Configuración del sistema
	config := &storage.VectorGraphConfig{
		QdrantURL:     "localhost",
		QdrantEnabled: true,
		Neo4jURI:      "bolt://localhost:7687",
		Neo4jUsername: "neo4j",
		Neo4jPassword: "password",
		Neo4jDatabase: "neo4j",
		Neo4jEnabled:  true,
		EmbeddingSize: 256,
	}

	// Inicializar cliente de Claude
	claudeClient := vibe.NewClaudeClient("your-anthropic-api-key")

	// Crear el servicio vector/graph
	service, err := storage.NewVectorGraphService(config, claudeClient, "example-project")
	if err != nil {
		log.Fatalf("❌ Error inicializando el servicio: %v", err)
	}

	ctx := context.Background()

	fmt.Printf("✅ Servicio inicializado. Habilitado: %t\n", service.IsEnabled())

	// Ejemplo 1: Almacenar componentes
	fmt.Println("\n📦 1. Almacenando componentes...")
	components := createExampleComponents()

	for _, component := range components {
		if err := service.StoreComponent(ctx, &component); err != nil {
			log.Printf("⚠️ Error almacenando componente %s: %v", component.ID, err)
		} else {
			fmt.Printf("✅ Componente almacenado: %s (%s)\n", component.ID, component.Type)
		}
	}

	// Ejemplo 2: Almacenar esquemas
	fmt.Println("\n🗂️ 2. Almacenando esquemas...")
	schemas := createExampleSchemas()

	for _, schema := range schemas {
		if err := service.StoreSchema(ctx, &schema); err != nil {
			log.Printf("⚠️ Error almacenando esquema %s: %v", schema.ID, err)
		} else {
			fmt.Printf("✅ Esquema almacenado: %s (%s)\n", schema.ID, schema.Name)
		}
	}

	// Ejemplo 3: Almacenar conversaciones
	fmt.Println("\n💬 3. Almacenando conversaciones...")
	conversations := createExampleConversations()

	for _, conversation := range conversations {
		if err := service.StoreConversation(ctx, &conversation, "example-session"); err != nil {
			log.Printf("⚠️ Error almacenando conversación: %v", err)
		} else {
			fmt.Printf("✅ Conversación almacenada: %s\n", conversation.Content[:50]+"...")
		}
	}

	// Ejemplo 4: Crear relaciones entre componentes
	fmt.Println("\n🔗 4. Creando relaciones...")

	// Relación espacial: botón cerca del formulario
	if err := service.CreateComponentRelationship(ctx, "button_login", "form_login", "NEAR", 0.9); err != nil {
		log.Printf("⚠️ Error creando relación: %v", err)
	} else {
		fmt.Println("✅ Relación creada: button_login -> form_login (NEAR)")
	}

	// Relación funcional: esquema genera componente
	if err := service.CreateSchemaComponentRelationship(ctx, "schema_user", "form_login", "GENERATES", 0.8); err != nil {
		log.Printf("⚠️ Error creando relación esquema-componente: %v", err)
	} else {
		fmt.Println("✅ Relación creada: schema_user -> form_login (GENERATES)")
	}

	// Ejemplo 5: Búsqueda semántica
	fmt.Println("\n🔍 5. Realizando búsqueda semántica...")

	searchQueries := []string{
		"componente de autenticación",
		"formulario de login",
		"botón para enviar",
		"esquema de usuario",
	}

	for _, query := range searchQueries {
		fmt.Printf("\n🔎 Buscando: '%s'\n", query)

		results, err := service.SemanticSearch(ctx, query, 3)
		if err != nil {
			log.Printf("⚠️ Error en búsqueda: %v", err)
			continue
		}

		fmt.Printf("📊 Resultados encontrados: %d vectoriales, %d grafos (Score: %.2f)\n",
			len(results.VectorResults), len(results.GraphResults), results.Score)
		fmt.Printf("💡 Explicación: %s\n", results.Explanation)

		// Mostrar algunos resultados
		for i, result := range results.VectorResults {
			if i >= 2 { // Mostrar solo los primeros 2
				break
			}
			fmt.Printf("   📌 Vector: %s (Score: %.3f)\n", result.ID, result.Score)
		}
	}

	// Ejemplo 6: Encontrar componentes relacionados
	fmt.Println("\n🔗 6. Encontrando componentes relacionados...")

	related, err := service.FindRelatedComponents(ctx, "button_login", 2)
	if err != nil {
		log.Printf("⚠️ Error encontrando relacionados: %v", err)
	} else {
		fmt.Printf("✅ Encontrados %d componentes relacionados con 'button_login':\n", len(related))
		for _, comp := range related {
			fmt.Printf("   🔗 %s (%s) - %s\n", comp.ID, comp.Type, comp.Description)
		}
	}

	// Ejemplo 7: Obtener insights de relaciones
	fmt.Println("\n📈 7. Obteniendo insights de relaciones...")

	insights, err := service.GetRelationshipInsights(ctx)
	if err != nil {
		log.Printf("⚠️ Error obteniendo insights: %v", err)
	} else {
		fmt.Printf("✅ Insights encontrados: %d\n", len(insights))
		for i, insight := range insights {
			if i >= 3 { // Mostrar solo los primeros 3
				break
			}
			fmt.Printf("   📊 Tipo: %s (Confianza: %.2f) - %s\n",
				insight.Type, insight.Confidence, insight.Explanation)
		}
	}

	// Ejemplo 8: Obtener estadísticas del proyecto
	fmt.Println("\n📊 8. Obteniendo estadísticas del proyecto...")

	stats, err := service.GetProjectStats(ctx)
	if err != nil {
		log.Printf("⚠️ Error obteniendo estadísticas: %v", err)
	} else {
		fmt.Printf("✅ Estadísticas del proyecto:\n")
		if enabled, ok := stats["enabled"].(bool); ok {
			fmt.Printf("   🟢 Habilitado: %t\n", enabled)
		}
		if vectorInfo, ok := stats["vector_storage"].(map[string]interface{}); ok {
			if pointsCount, ok := vectorInfo["points_count"]; ok {
				fmt.Printf("   📦 Vectores almacenados: %v\n", pointsCount)
			}
		}
		if graphInfo, ok := stats["graph_storage"].(map[string]interface{}); ok {
			if nodeCount, ok := graphInfo["node_count"]; ok {
				fmt.Printf("   🔗 Nodos en grafo: %v\n", nodeCount)
			}
			if relCount, ok := graphInfo["relationship_count"]; ok {
				fmt.Printf("   🔗 Relaciones en grafo: %v\n", relCount)
			}
		}
	}

	// Cerrar conexiones
	fmt.Println("\n🔒 Cerrando conexiones...")
	if err := service.Close(ctx); err != nil {
		log.Printf("⚠️ Error cerrando servicio: %v", err)
	} else {
		fmt.Println("✅ Servicio cerrado correctamente")
	}

	fmt.Println("\n🎉 Ejemplo completado exitosamente!")
}

// createExampleComponents crea componentes de ejemplo
func createExampleComponents() []prompts.ComponentState {
	return []prompts.ComponentState{
		{
			ID:       "button_login",
			Type:     "button",
			Category: "atom",
			Name:     "Login Button",
			Properties: map[string]interface{}{
				"text":    "Iniciar Sesión",
				"variant": "primary",
				"size":    "medium",
			},
			Position: prompts.Position{X: 100, Y: 100},
			Size:     prompts.Size{W: 120, H: 40},
			Style: map[string]interface{}{
				"backgroundColor": "#3B82F6",
				"color":           "#FFFFFF",
			},
		},
		{
			ID:       "form_login",
			Type:     "form",
			Category: "molecule",
			Name:     "Login Form",
			Properties: map[string]interface{}{
				"method": "POST",
				"action": "/api/login",
				"fields": []string{"email", "password"},
			},
			Position: prompts.Position{X: 50, Y: 200},
			Size:     prompts.Size{W: 300, H: 200},
			Style: map[string]interface{}{
				"padding": "20px",
				"border":  "1px solid #E5E7EB",
			},
		},
		{
			ID:       "text_welcome",
			Type:     "text",
			Category: "atom",
			Name:     "Welcome Text",
			Properties: map[string]interface{}{
				"content": "Bienvenido a VibeCode",
				"size":    "large",
				"weight":  "bold",
			},
			Position: prompts.Position{X: 50, Y: 50},
			Size:     prompts.Size{W: 300, H: 30},
			Style: map[string]interface{}{
				"fontSize":   "24px",
				"fontWeight": "bold",
				"color":      "#1F2937",
			},
		},
	}
}

// createExampleSchemas crea esquemas de ejemplo
func createExampleSchemas() []models.ResourceSchema {
	return []models.ResourceSchema{
		{
			ID:          "schema_user",
			Name:        "User",
			Description: "Esquema para usuarios del sistema",
			Fields: []models.SchemaField{
				{
					Name:        "id",
					Type:        "uuid",
					Description: "Identificador único del usuario",
					Required:    true,
					Constraints: map[string]interface{}{
						"primaryKey": true,
					},
				},
				{
					Name:        "email",
					Type:        "email",
					Description: "Correo electrónico del usuario",
					Required:    true,
					Constraints: map[string]interface{}{
						"unique": true,
					},
				},
				{
					Name:        "password",
					Type:        "string",
					Description: "Contraseña hasheada",
					Required:    true,
					Constraints: map[string]interface{}{
						"minLength": 8,
					},
				},
				{
					Name:        "name",
					Type:        "string",
					Description: "Nombre completo del usuario",
					Required:    true,
				},
				{
					Name:        "createdAt",
					Type:        "timestamp",
					Description: "Fecha de creación",
					Required:    true,
				},
			},
			Metadata: map[string]interface{}{
				"tableName": "users",
				"indexes":   []string{"email", "createdAt"},
				"relations": map[string]interface{}{
					"sessions": "hasMany",
					"profile":  "hasOne",
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:          "schema_session",
			Name:        "Session",
			Description: "Esquema para sesiones de usuario",
			Fields: []models.SchemaField{
				{
					Name:        "id",
					Type:        "uuid",
					Description: "Identificador único de la sesión",
					Required:    true,
				},
				{
					Name:        "userId",
					Type:        "uuid",
					Description: "ID del usuario propietario",
					Required:    true,
				},
				{
					Name:        "token",
					Type:        "string",
					Description: "Token de sesión",
					Required:    true,
				},
				{
					Name:        "expiresAt",
					Type:        "timestamp",
					Description: "Fecha de expiración",
					Required:    true,
				},
			},
			Metadata: map[string]interface{}{
				"tableName": "sessions",
				"relations": map[string]interface{}{
					"user": "belongsTo",
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

// createExampleConversations crea conversaciones de ejemplo
func createExampleConversations() []vibe.Message {
	return []vibe.Message{
		{
			Role:      "user",
			Content:   "Necesito crear un formulario de login para mi aplicación",
			Type:      "text",
			Timestamp: time.Now(),
		},
		{
			Role:      "assistant",
			Content:   "Perfecto, puedo ayudarte a crear un formulario de login. Incluirá campos para email y contraseña, con validación y un botón de envío.",
			Type:      "text",
			Timestamp: time.Now(),
		},
		{
			Role:      "user",
			Content:   "¿Puedes agregar un botón de login azul?",
			Type:      "text",
			Timestamp: time.Now(),
		},
		{
			Role:      "assistant",
			Content:   "¡Claro! Agregando un botón azul para el login. Será un botón primario con texto blanco.",
			Type:      "text",
			Timestamp: time.Now(),
		},
		{
			Role:      "user",
			Content:   "¿Qué esquema de datos necesito para autenticación?",
			Type:      "text",
			Timestamp: time.Now(),
		},
		{
			Role:      "assistant",
			Content:   "Para autenticación necesitarás al menos dos esquemas: User (con email, password, name) y Session (para manejar los tokens de sesión activa).",
			Type:      "text",
			Timestamp: time.Now(),
		},
	}
}
