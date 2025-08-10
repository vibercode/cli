# 🚀 VibeCode Vector/Graph Storage

Integración de **Qdrant** (base de datos vectorial) y **Neo4j** (base de datos de grafos) al modo "vibe" de VibeCode para crear relaciones semánticas inteligentes entre componentes, esquemas y conversaciones.

## ✨ Características

- 🔍 **Búsqueda semántica** de componentes por similitud
- 🔗 **Relaciones automáticas** entre elementos del proyecto
- 💾 **Almacenamiento persistente** de vectores y grafos
- 🤖 **Embeddings generados con Claude AI**
- 📊 **Analytics y insights** de patrones de uso
- ⚡ **Integración transparente** con el modo vibe existente

## 🏗️ Arquitectura

```
┌────────────────┐    ┌─────────────────┐    ┌──────────────────┐
│   Claude AI    │───▶│ Embedding       │───▶│   Qdrant DB     │
│   (Embeddings) │    │ Service         │    │   (Vectors)      │
└────────────────┘    └─────────────────┘    └──────────────────┘
                             │
                             ▼
┌────────────────┐    ┌─────────────────┐    ┌──────────────────┐
│   Chat Mode    │───▶│ VectorGraph     │───▶│   Neo4j DB      │
│   (Vibe)       │    │ Service         │    │   (Graphs)       │
└────────────────┘    └─────────────────┘    └──────────────────┘
```

## 🚀 Inicio Rápido

### 1. Instalación de Dependencias

```bash
# Clonar el repositorio
git clone <repository-url>
cd vibercode-cli-go

# Instalar dependencias Go
go mod tidy

# Levantar Qdrant y Neo4j con Docker
docker-compose up -d
```

### 2. Configuración

```bash
# Variables de entorno requeridas
export QDRANT_ENABLED="true"
export NEO4J_ENABLED="true"
export ANTHROPIC_API_KEY="tu_clave_aqui"

# Variables opcionales (con valores por defecto)
export QDRANT_URL="localhost"
export NEO4J_URI="bolt://localhost:7687"
export NEO4J_USERNAME="neo4j"
export NEO4J_PASSWORD="password"
export NEO4J_DATABASE="neo4j"
```

### 3. Uso Básico

```bash
# Modo vibe normal
./vibercode vibe

# El modo vibe detectará automáticamente si Qdrant y Neo4j están disponibles
# y habilitará las funcionalidades de vector/graph storage
```

## 🔧 Servicios Disponibles

### 📦 VectorStorage (Qdrant)

- Almacena embeddings de 256 dimensiones
- Búsqueda por similitud semántica
- Filtrado por proyecto
- Operaciones CRUD optimizadas

### 🔗 GraphStorage (Neo4j)

- Relaciones tipadas entre nodos
- Consultas Cypher avanzadas
- Análisis de patrones y paths
- Constraints y índices automáticos

### 🧠 EmbeddingService (Claude)

- Generación de vectores semánticos
- Cache inteligente para optimización
- Fallback determinístico sin Claude
- Soporte para componentes, esquemas y conversaciones

### 🎯 VectorGraphService (Unificado)

- Coordinación entre vectores y grafos
- Almacenamiento automático
- Búsqueda híbrida
- Analytics y métricas

## 📡 API Endpoints

### Búsqueda Semántica

```bash
POST /api/semantic-search
Content-Type: application/json

{
  "query": "componente de autenticación",
  "limit": 10
}
```

### Componentes Relacionados

```bash
GET /api/related-components/button_123
```

### Insights de Relaciones

```bash
GET /api/relationship-insights
```

### Estadísticas del Proyecto

```bash
GET /api/project-stats
```

## 💻 Ejemplo de Código

```go
// Configuración del servicio
config := &storage.VectorGraphConfig{
    QdrantEnabled: true,
    Neo4jEnabled:  true,
    EmbeddingSize: 256,
}

// Inicializar servicio
service, err := storage.NewVectorGraphService(config, claudeClient, "mi-proyecto")

// Almacenar un componente
component := &prompts.ComponentState{
    ID:       "button_login",
    Type:     "button",
    Category: "atom",
    Properties: map[string]interface{}{
        "text": "Login",
        "variant": "primary",
    },
}

err = service.StoreComponent(ctx, component)

// Búsqueda semántica
results, err := service.SemanticSearch(ctx, "botón de autenticación", 5)

// Encontrar relacionados
related, err := service.FindRelatedComponents(ctx, "button_login", 2)
```

## 🔍 Casos de Uso

### 1. **Sugerencias Inteligentes**

Cuando el usuario arrastra un componente, el sistema sugiere elementos relacionados basados en patrones similares.

### 2. **Detección de Patrones**

Identifica automáticamente patrones comunes de diseño y sugiere mejores prácticas.

### 3. **Búsqueda Contextual**

Encuentra componentes por descripción semántica, no solo por nombre exacto.

### 4. **Análisis de Coherencia**

Detecta inconsistencias entre esquemas de datos y componentes de UI.

### 5. **Historial Inteligente**

Recupera conversaciones y decisiones de diseño relevantes al contexto actual.

## 📊 Métricas y Analytics

El sistema recopila automáticamente:

- **Componentes más utilizados** y patrones frecuentes
- **Eficiencia de búsquedas** semánticas
- **Calidad de sugerencias** y aceptación
- **Densidad de relaciones** en el grafo del proyecto
- **Evolución temporal** de la arquitectura

## 🛠️ Desarrollo

### Estructura de Archivos

```
internal/storage/
├── vector_storage.go      # Cliente Qdrant
├── graph_storage.go       # Cliente Neo4j
├── embedding_service.go   # Generación de embeddings
└── vector_graph_service.go # Servicio unificado

examples/
└── vector_graph_example.go # Ejemplo completo

docker-compose.yml         # Servicios de BD
```

### Ejecutar Ejemplo

```bash
# Levantar servicios
docker-compose up -d

# Configurar variables
export ANTHROPIC_API_KEY="tu_clave"
export QDRANT_ENABLED="true"
export NEO4J_ENABLED="true"

# Ejecutar ejemplo
go run examples/vector_graph_example.go
```

### Tests

```bash
# Tests unitarios
go test ./internal/storage/...

# Test de integración (requiere Docker)
docker-compose up -d
go test -tags integration ./internal/storage/...
```

## 🐳 Docker

### Servicios Incluidos

- **Qdrant**: Puerto 6333 (HTTP), 6334 (gRPC)
- **Neo4j**: Puerto 7474 (HTTP), 7687 (Bolt)

### Datos Persistentes

Los datos se almacenan en:

- `./data/qdrant_data/` - Vectores de Qdrant
- `./data/neo4j_data/` - Grafo de Neo4j
- `./data/neo4j_logs/` - Logs de Neo4j

### Comandos Útiles

```bash
# Levantar servicios
docker-compose up -d

# Ver logs
docker-compose logs -f

# Parar servicios
docker-compose down

# Limpiar datos (⚠️ DESTRUCTIVO)
docker-compose down -v
rm -rf ./data/
```

## 🔧 Configuración Avanzada

### Variables de Entorno Completas

```bash
# Habilitación de servicios
QDRANT_ENABLED="true"
NEO4J_ENABLED="true"

# Configuración de Qdrant
QDRANT_URL="localhost"
QDRANT_PORT="6333"

# Configuración de Neo4j
NEO4J_URI="bolt://localhost:7687"
NEO4J_USERNAME="neo4j"
NEO4J_PASSWORD="password"
NEO4J_DATABASE="neo4j"

# Configuración de Claude
ANTHROPIC_API_KEY="tu_clave_aqui"

# Configuración de embeddings
EMBEDDING_SIZE="256"
EMBEDDING_CACHE_ENABLED="true"
```

### Configuración de Producción

```yaml
# docker-compose.prod.yml
version: "3.8"
services:
  qdrant:
    image: qdrant/qdrant:latest
    restart: unless-stopped
    volumes:
      - qdrant_data:/qdrant/storage
    environment:
      - QDRANT__STORAGE__MEMORY_THRESHOLD=1073741824
    deploy:
      resources:
        limits:
          memory: 2G
          cpus: "1"

  neo4j:
    image: neo4j:latest
    restart: unless-stopped
    volumes:
      - neo4j_data:/data
    environment:
      - NEO4J_AUTH=neo4j/strong_password_here
      - NEO4J_dbms_memory_heap_max__size=1G
      - NEO4J_dbms_memory_pagecache_size=512m
    deploy:
      resources:
        limits:
          memory: 2G
          cpus: "1"
```

## 🔍 Troubleshooting

### Problemas Comunes

**1. Error de conexión a Qdrant**

```bash
# Verificar que Qdrant esté corriendo
curl http://localhost:6333/health

# Verificar logs
docker logs vibercode-qdrant
```

**2. Error de conexión a Neo4j**

```bash
# Verificar conectividad
docker exec vibercode-neo4j cypher-shell -u neo4j -p password "RETURN 1"

# Verificar logs
docker logs vibercode-neo4j
```

**3. Error de Claude API**

```bash
# Verificar clave de API
echo $ANTHROPIC_API_KEY

# Test básico
curl -H "x-api-key: $ANTHROPIC_API_KEY" \
     -H "anthropic-version: 2023-06-01" \
     https://api.anthropic.com/v1/messages
```

### Logs de Debug

```go
// Habilitar logs detallados
log.SetLevel(log.DebugLevel)

// En el código
log.Printf("🔍 Debug: %+v", data)
```

## 📈 Roadmap

### v1.0 ✅

- [x] Integración básica con Qdrant y Neo4j
- [x] Embeddings con Claude API
- [x] Búsqueda semántica
- [x] Relaciones automáticas

### v1.1 🔄

- [ ] Interfaz web para explorar grafos
- [ ] Métricas avanzadas y dashboards
- [ ] Export/import de proyectos
- [ ] Optimizaciones de rendimiento

### v1.2 🚀

- [ ] Clustering automático de componentes
- [ ] Recomendaciones proactivas
- [ ] Colaboración en tiempo real
- [ ] Integración con Git

## 🤝 Contribuciones

¡Las contribuciones son bienvenidas! Por favor:

1. Fork el repositorio
2. Crea una rama para tu feature (`git checkout -b feature/amazing-feature`)
3. Commit tus cambios (`git commit -m 'Add amazing feature'`)
4. Push a la rama (`git push origin feature/amazing-feature`)
5. Abre un Pull Request

## 📄 Licencia

Este proyecto está bajo la licencia MIT. Ver `LICENSE` para más detalles.

## 🙏 Agradecimientos

- [Qdrant](https://qdrant.tech/) - Vector database
- [Neo4j](https://neo4j.com/) - Graph database
- [Claude AI](https://anthropic.com/) - Embedding generation
- [Go](https://golang.org/) - Programming language

---

**¿Preguntas?** Abre un [issue](https://github.com/tu-repo/vibercode/issues) o únete a nuestro [Discord](https://discord.gg/vibercode).
