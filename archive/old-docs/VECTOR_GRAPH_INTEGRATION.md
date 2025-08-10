# 🚀 Integración de Qdrant y Neo4j al Modo Vibe

## 📋 Resumen de la Propuesta

Esta propuesta describe cómo integrar **Qdrant** (base de datos vectorial) y **Neo4j** (base de datos de grafos) al modo "vibe" de VibeCode para crear un sistema inteligente de relaciones entre componentes, esquemas y conversaciones.

## 🎯 Objetivos

### 1. **Almacenamiento Vectorial con Qdrant**

- Generar embeddings semánticos para componentes UI
- Almacenar vectores de esquemas de datos
- Crear embeddings de conversaciones para contexto
- Búsqueda semántica inteligente

### 2. **Almacenamiento de Grafos con Neo4j**

- Modelar relaciones entre componentes
- Conectar esquemas con componentes generados
- Rastrear flujo de conversaciones
- Análisis de patrones de uso

### 3. **Integración Inteligente**

- Sugerencias contextualess basadas en similitud semántica
- Detección automática de patrones en diseños
- Recomendaciones de componentes relacionados
- Análisis de coherencia en el proyecto

## 🏗️ Arquitectura del Sistema

```
┌─────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                    MODO VIBE MEJORADO                                               │
├─────────────────────────┬─────────────────────────┬─────────────────────────┬─────────────────────────┤
│     Claude AI Chat      │    Component Canvas     │    Schema Designer      │    Project Manager      │
│                         │                         │                         │                         │
│  • Conversaciones       │  • Componentes UI       │  • Esquemas de datos    │  • Gestión de proyectos │
│  • Intents & Context    │  • Posicionamiento      │  • Campos & relaciones  │  • Sesiones de trabajo  │
│  • Embeddings de NLP    │  • Propiedades          │  • Validaciones         │  • Historial de cambios │
└─────────────────────────┴─────────────────────────┴─────────────────────────┴─────────────────────────┘
                                                    │
                                                    ▼
┌─────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                              VECTOR GRAPH SERVICE                                                   │
├─────────────────────────┬─────────────────────────┬─────────────────────────┬─────────────────────────┤
│   Embedding Service     │    Vector Storage       │    Graph Storage        │   Relationship Engine   │
│                         │                         │                         │                         │
│  • Claude-generated     │  • Qdrant Database      │  • Neo4j Database       │  • Spatial proximity    │
│  • Semantic vectors     │  • Component vectors    │  • Component nodes      │  • Functional relations │
│  • 256-dimensional     │  • Schema vectors       │  • Schema nodes         │  • Temporal connections │
│  • Cached embeddings   │  • Conversation vectors │  • Conversation nodes   │  • Semantic similarity  │
└─────────────────────────┴─────────────────────────┴─────────────────────────┴─────────────────────────┘
                                                    │
                                                    ▼
┌─────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                            STORAGE & ANALYSIS LAYER                                                 │
├─────────────────────────┬─────────────────────────┬─────────────────────────┬─────────────────────────┤
│    Qdrant Cluster      │    Neo4j Database       │   Analytics Engine      │    Export/Import        │
│                         │                         │                         │                         │
│  • Vector collections   │  • Graph database       │  • Pattern detection    │  • JSON/GraphML export │
│  • Similarity search   │  • Cypher queries       │  • Usage analytics      │  • Backup/restore       │
│  • Filtering & facets  │  • Relationship paths   │  • Recommendation AI    │  • Project templates    │
│  • Real-time updates   │  • Constraint validation│  • Performance metrics  │  • Collaboration sync   │
└─────────────────────────┴─────────────────────────┴─────────────────────────┴─────────────────────────┘
```

## 🔧 Implementación Técnica

### 1. **Servicios Principales Implementados**

#### 📁 `internal/storage/vector_storage.go`

- **VectorStorage**: Cliente para Qdrant con operaciones CRUD
- **ComponentVector**: Representación vectorial de componentes
- **SchemaVector**: Representación vectorial de esquemas
- **ConversationVector**: Representación vectorial de conversaciones
- **Búsqueda semántica** por similitud de coseno

#### 📁 `internal/storage/graph_storage.go`

- **GraphStorage**: Cliente para Neo4j con operaciones de grafos
- **ComponentNode**: Nodos de componentes en el grafo
- **SchemaNode**: Nodos de esquemas en el grafo
- **ConversationNode**: Nodos de conversaciones en el grafo
- **Relationship**: Relaciones tipadas entre nodos

#### 📁 `internal/storage/embedding_service.go`

- **EmbeddingService**: Generación de embeddings usando Claude
- **Vector de 256 dimensiones** con características semánticas
- **Cachéado inteligente** para evitar regeneración
- **Fallback determinístico** cuando Claude no está disponible

#### 📁 `internal/storage/vector_graph_service.go`

- **VectorGraphService**: Servicio unificado que coordina Qdrant y Neo4j
- **Almacenamiento automático** de componentes, esquemas y conversaciones
- **Búsqueda semántica híbrida** combinando vectores y grafos
- **Análisis de relaciones** y generación de insights

### 2. **Configuración del Sistema**

#### Variables de Entorno

```bash
# Configuración de Qdrant
export QDRANT_URL="localhost"
export QDRANT_ENABLED="true"

# Configuración de Neo4j
export NEO4J_URI="bolt://localhost:7687"
export NEO4J_USERNAME="neo4j"
export NEO4J_PASSWORD="password"
export NEO4J_DATABASE="neo4j"
export NEO4J_ENABLED="true"

# Configuración de Claude
export ANTHROPIC_API_KEY="tu_clave_aqui"
```

#### Dependencias en `go.mod`

```go
require (
    github.com/qdrant/go-client v1.7.0
    github.com/neo4j/neo4j-go-driver/v5 v5.15.0
    // ... existing dependencies
)
```

### 3. **Flujo de Datos Completo**

#### 🔄 Guardado Automático

```
Usuario interactúa con el sistema
           ↓
┌─────────────────────────────────────────────────────────────────────────────────────────────────────┐
│  1. CREAR/MODIFICAR COMPONENTE                                                                       │
│     • Usuario arrastra componente al canvas                                                         │
│     • Sistema captura: tipo, propiedades, posición, tamaño                                         │
│     • Genera embedding semántico del componente                                                     │
│     • Almacena en Qdrant como vector                                                               │
│     • Almacena en Neo4j como nodo                                                                  │
│     • Crea relaciones espaciales con componentes cercanos                                          │
└─────────────────────────────────────────────────────────────────────────────────────────────────────┘
           ↓
┌─────────────────────────────────────────────────────────────────────────────────────────────────────┐
│  2. GENERAR ESQUEMA                                                                                 │
│     • Usuario crea esquema de datos                                                                │
│     • Sistema genera embedding del esquema                                                         │
│     • Almacena en Qdrant como vector                                                               │
│     • Almacena en Neo4j como nodo                                                                  │
│     • Crea relaciones con componentes que usan el esquema                                          │
└─────────────────────────────────────────────────────────────────────────────────────────────────────┘
           ↓
┌─────────────────────────────────────────────────────────────────────────────────────────────────────┐
│  3. CONVERSACIÓN CON CLAUDE                                                                        │
│     • Usuario envía mensaje en el chat                                                             │
│     • Sistema genera embedding de la conversación                                                  │
│     • Almacena en Qdrant como vector                                                               │
│     • Almacena en Neo4j como nodo                                                                  │
│     • Crea relaciones con componentes/esquemas mencionados                                         │
└─────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

#### 🔍 Búsqueda Semántica

```
Usuario busca "botón de login"
           ↓
┌─────────────────────────────────────────────────────────────────────────────────────────────────────┐
│  1. GENERAR EMBEDDING DE CONSULTA                                                                   │
│     • Convierte consulta en vector de 256 dimensiones                                              │
│     • Usa Claude para generar embedding semántico                                                  │
└─────────────────────────────────────────────────────────────────────────────────────────────────────┘
           ↓
┌─────────────────────────────────────────────────────────────────────────────────────────────────────┐
│  2. BÚSQUEDA VECTORIAL EN QDRANT                                                                    │
│     • Busca componentes similares por coseno similarity                                            │
│     • Busca esquemas relacionados                                                                  │
│     • Busca conversaciones relevantes                                                              │
│     • Filtra por proyecto actual                                                                   │
└─────────────────────────────────────────────────────────────────────────────────────────────────────┘
           ↓
┌─────────────────────────────────────────────────────────────────────────────────────────────────────┐
│  3. BÚSQUEDA DE GRAFOS EN NEO4J                                                                     │
│     • Encuentra componentes conectados por relaciones                                              │
│     • Explora paths entre nodos relacionados                                                       │
│     • Analiza patrones de uso                                                                      │
└─────────────────────────────────────────────────────────────────────────────────────────────────────┘
           ↓
┌─────────────────────────────────────────────────────────────────────────────────────────────────────┐
│  4. COMBINAR RESULTADOS                                                                            │
│     • Merge resultados de vectores y grafos                                                        │
│     • Calcula score combinado (60% vectorial, 40% grafo)                                          │
│     • Genera explicación de los resultados                                                         │
│     • Retorna ranking de componentes/esquemas relevantes                                           │
└─────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

### 4. **Endpoints de la API**

#### 🔍 Búsqueda Semántica

```bash
POST /api/semantic-search
{
  "query": "componente de autenticación",
  "limit": 10
}
```

#### 🔗 Componentes Relacionados

```bash
GET /api/related-components/button_123
```

#### 📊 Insights de Relaciones

```bash
GET /api/relationship-insights
```

#### 📈 Estadísticas del Proyecto

```bash
GET /api/project-stats
```

## 🚀 Casos de Uso Prácticos

### 1. **Sugerencias Inteligentes**

- **Escenario**: Usuario arrastra un botón al canvas
- **Acción**: Sistema sugiere formulario relacionado basado en patrones similares
- **Beneficio**: Acelera el desarrollo con sugerencias contextuales

### 2. **Detección de Patrones**

- **Escenario**: Usuario está creando una página de login
- **Acción**: Sistema detecta patrón común y sugiere componentes complementarios
- **Beneficio**: Consistencia en diseños y mejores prácticas

### 3. **Búsqueda Semántica**

- **Escenario**: Usuario busca "componente para mostrar productos"
- **Acción**: Sistema encuentra cards, listas y grids relacionados
- **Beneficio**: Descubrimiento intuitivo de componentes relevantes

### 4. **Análisis de Coherencia**

- **Escenario**: Usuario tiene esquema de "User" pero no componentes relacionados
- **Acción**: Sistema sugiere crear componentes UserCard, UserForm, etc.
- **Beneficio**: Coherencia entre datos y UI

### 5. **Recomendaciones Contextuales**

- **Escenario**: Usuario conversa sobre "autenticación" en el chat
- **Acción**: Sistema resalta componentes y esquemas relacionados
- **Beneficio**: Conexión inteligente entre conversación y elementos del proyecto

## 🔧 Instalación y Configuración

### 1. **Instalación de Dependencias**

#### Docker Compose para Qdrant y Neo4j

```yaml
version: "3.8"
services:
  qdrant:
    image: qdrant/qdrant:latest
    ports:
      - "6333:6333"
    volumes:
      - ./qdrant_data:/qdrant/storage
    environment:
      - QDRANT__SERVICE__HTTP_PORT=6333
      - QDRANT__SERVICE__GRPC_PORT=6334

  neo4j:
    image: neo4j:latest
    ports:
      - "7474:7474"
      - "7687:7687"
    volumes:
      - ./neo4j_data:/data
    environment:
      - NEO4J_AUTH=neo4j/password
      - NEO4J_dbms_memory_heap_max__size=512m
```

#### Inicialización

```bash
# Levantar servicios
docker-compose up -d

# Instalar dependencias Go
go mod tidy

# Configurar variables de entorno
export QDRANT_ENABLED="true"
export NEO4J_ENABLED="true"
export ANTHROPIC_API_KEY="tu_clave_aqui"
```

### 2. **Integración con Modo Vibe**

#### Modificaciones en `internal/vibe/preview.go`

```go
// Agregar VectorGraphService al PreviewServer
type PreviewServer struct {
    // ... existing fields
    vectorGraphService *storage.VectorGraphService
}

// Inicializar en NewPreviewServer
func NewPreviewServer(port string) *PreviewServer {
    // ... existing code

    // Configurar vector/graph service
    config := &storage.VectorGraphConfig{
        QdrantEnabled: os.Getenv("QDRANT_ENABLED") == "true",
        Neo4jEnabled:  os.Getenv("NEO4J_ENABLED") == "true",
        // ... other config
    }

    vectorGraphService, err := storage.NewVectorGraphService(config, claudeClient, projectID)
    if err != nil {
        log.Printf("Vector/Graph service disabled: %v", err)
    }

    return &PreviewServer{
        // ... existing fields
        vectorGraphService: vectorGraphService,
    }
}
```

#### Hooks de Almacenamiento Automático

```go
// En handleWebSocketViewUpdate
func (ps *PreviewServer) handleWebSocketViewUpdate(conn *websocket.Conn, msg WebSocketMessage) {
    // ... existing code

    // Guardar componentes automáticamente
    if ps.vectorGraphService != nil && ps.vectorGraphService.IsEnabled() {
        for _, component := range updatedComponents {
            if err := ps.vectorGraphService.StoreComponent(ctx, &component); err != nil {
                log.Printf("Failed to store component: %v", err)
            }
        }
    }
}
```

### 3. **Comandos del CLI**

#### Nuevo Comando `vibe-enhanced`

```bash
# Modo vibe con vector/graph storage
./vibercode vibe-enhanced --project="mi-proyecto"

# Con configuración personalizada
./vibercode vibe-enhanced \
  --project="mi-proyecto" \
  --qdrant-url="localhost:6333" \
  --neo4j-uri="bolt://localhost:7687"
```

## 📊 Métricas y Análisis

### 1. **Métricas de Uso**

- Componentes más utilizados
- Patrones de diseño frecuentes
- Eficiencia en búsquedas semánticas
- Calidad de sugerencias

### 2. **Análisis de Relaciones**

- Densidad de conexiones en el grafo
- Componentes huérfanos (sin relaciones)
- Clusters de componentes relacionados
- Evolución temporal de relaciones

### 3. **Insights de Conversaciones**

- Temas más discutidos
- Efectividad de respuestas de Claude
- Patrones en consultas de usuarios
- Correlación entre chat y creación de componentes

## 🎯 Próximos Pasos

### Fase 1: Implementación Base ✅

- [x] Servicio de almacenamiento vectorial (Qdrant)
- [x] Servicio de almacenamiento de grafos (Neo4j)
- [x] Servicio de embeddings con Claude
- [x] Servicio unificado VectorGraphService

### Fase 2: Integración con Vibe Mode 🔄

- [ ] Modificar PreviewServer para usar VectorGraphService
- [ ] Hooks de almacenamiento automático en eventos
- [ ] Endpoints de API para búsqueda semántica
- [ ] Interfaz de usuario para explorar relaciones

### Fase 3: Funcionalidades Avanzadas 🚀

- [ ] Sugerencias proactivas basadas en contexto
- [ ] Análisis predictivo de patrones
- [ ] Exportación de grafos de proyecto
- [ ] Colaboración en tiempo real con sincronización

### Fase 4: Optimización y Escalabilidad 📈

- [ ] Optimización de rendimiento de vectores
- [ ] Clustering de componentes similares
- [ ] Compresión de embeddings
- [ ] Replicación y backup automático

## 🔍 Ejemplo de Uso Completo

```go
// Inicializar el sistema mejorado
server := NewEnhancedPreviewServer("3001", "mi-proyecto")

// El usuario arrastra un botón al canvas
component := &prompts.ComponentState{
    ID:       "button_123",
    Type:     "button",
    Category: "atom",
    Properties: map[string]interface{}{
        "text": "Login",
        "variant": "primary",
    },
    Position: prompts.Position{X: 100, Y: 100},
    Size:     prompts.Size{W: 120, H: 40},
}

// El sistema automáticamente:
// 1. Genera embedding semántico del botón
// 2. Almacena en Qdrant como vector
// 3. Almacena en Neo4j como nodo
// 4. Busca componentes cercanos y crea relaciones

// El usuario busca componentes relacionados
results, _ := server.vectorGraphService.SemanticSearch(ctx, "componente de autenticación", 5)

// El sistema devuelve:
// - Botones similares (Login, Register, etc.)
// - Formularios relacionados (LoginForm, etc.)
// - Esquemas relevantes (User, Session, etc.)
// - Conversaciones previas sobre autenticación
```

## 🎉 Beneficios del Sistema

### Para Desarrolladores

- **Desarrollo más rápido** con sugerencias inteligentes
- **Consistencia automática** en diseños
- **Descubrimiento fácil** de componentes existentes
- **Análisis visual** de relaciones en el proyecto

### Para Proyectos

- **Coherencia estructural** entre datos y UI
- **Patrones reutilizables** identificados automáticamente
- **Documentación visual** del proyecto
- **Análisis de complejidad** y refactoring

### Para Equipos

- **Conocimiento compartido** a través de patrones
- **Colaboración mejorada** con contexto común
- **Aprendizaje acelerado** para nuevos miembros
- **Estándares automáticos** en desarrollo

---

**Este sistema convierte el modo "vibe" en una herramienta inteligente que no solo facilita la creación de interfaces, sino que aprende de los patrones de uso y mejora continuamente las sugerencias y recomendaciones.**
