# 🚀 VibeCode CLI - Guía de Configuración Completa

## ✅ Sistema Compilado y Funcional

El proyecto ya está completamente funcional con todas las dependencias resueltas:

- ✅ Neo4j v5 API corregida
- ✅ Qdrant implementado en modo stub
- ✅ Ciclos de importación resueltos
- ✅ Preview Server Enhanced activo

## 🔧 Configuración Básica

### 1. Variables de Entorno

Copia `config.example.env` a `.env` y configura:

```bash
cp config.example.env .env
```

**Configuración mínima necesaria:**

```bash
# Solo para funcionalidad de AI Chat
ANTHROPIC_API_KEY=tu_clave_de_anthropic
```

**Configuración completa (opcional):**

```bash
# Para habilitar vector storage con Qdrant
QDRANT_ENABLED=true
QDRANT_URL=localhost:6333

# Para habilitar graph storage con Neo4j
NEO4J_ENABLED=true
NEO4J_URI=bolt://localhost:7687
NEO4J_USERNAME=neo4j
NEO4J_PASSWORD=tu_password
```

## 🎮 Comandos Disponibles

### 1. Modo Básico - Generación de Código

```bash
# Generar un nuevo proyecto
./vibercode generate project my-api

# Generar un recurso CRUD
./vibercode generate resource User

# Generar desde schema
./vibercode schema create user.json
./vibercode generate from-schema user.json
```

### 2. Modo Vibe - Chat AI Interactivo

```bash
# Iniciar chat AI con preview en tiempo real
./vibercode vibe

# Iniciar en puerto específico
./vibercode vibe --port 3001
```

### 3. Servidor de Desarrollo

```bash
# Iniciar servidor HTTP para integración con React Editor
./vibercode serve --port 8080

# Iniciar WebSocket server para tiempo real
./vibercode ws --port 8081
```

## 🎨 Funcionalidades del Modo Vibe

### Chat AI Mejorado

- 💬 Conversación con Claude AI
- 🧠 Contexto del proyecto mantenido
- 📊 Análisis de componentes en tiempo real
- 🔍 Búsqueda semántica (cuando vector storage está habilitado)

### Preview Server Enhanced

- ⚡ WebSocket para actualizaciones en tiempo real
- 🎯 APIs REST para integración con frontend
- 📈 Métricas y estadísticas del proyecto
- 🔗 Análisis de relaciones entre componentes

### Endpoints Disponibles

**WebSocket:** `ws://localhost:3001/ws`

- Actualizaciones de vista en tiempo real
- Chat bidireccional
- Estado sincronizado

**HTTP APIs:**

- `GET /api/status` - Estado del sistema
- `GET /api/view-state` - Estado actual de la vista
- `POST /api/chat` - Enviar mensajes al AI
- `GET /api/search?q=query` - Búsqueda semántica
- `GET /api/components/{id}/related` - Componentes relacionados
- `GET /api/insights` - Insights de relaciones
- `GET /api/stats` - Estadísticas del proyecto

## 🔮 Funcionalidades Avanzadas (Opcionales)

### Vector Storage con Qdrant

Si quieres habilitar búsqueda semántica:

```bash
# Instalar Qdrant con Docker
docker run -p 6333:6333 qdrant/qdrant

# Habilitar en .env
QDRANT_ENABLED=true
QDRANT_URL=localhost:6333
```

### Graph Storage con Neo4j

Para análisis avanzado de relaciones:

```bash
# Instalar Neo4j con Docker
docker run -p 7474:7474 -p 7687:7687 \
  -e NEO4J_AUTH=neo4j/password \
  neo4j:latest

# Habilitar en .env
NEO4J_ENABLED=true
NEO4J_PASSWORD=password
```

## 🚀 Inicio Rápido

### Solo Generación de Código

```bash
./vibercode generate project my-api
cd my-api
go run main.go
```

### Con Chat AI

```bash
export ANTHROPIC_API_KEY=tu_clave
./vibercode vibe
# Abre http://localhost:3001 en tu navegador
```

### Sistema Completo

```bash
# 1. Configurar variables
cp config.example.env .env
# Editar .env con tus configuraciones

# 2. Iniciar servicios opcionales (si los necesitas)
docker run -d -p 6333:6333 qdrant/qdrant
docker run -d -p 7687:7687 -e NEO4J_AUTH=neo4j/password neo4j

# 3. Iniciar VibeCode
./vibercode vibe --port 3001
```

## 🐛 Troubleshooting

### Error de compilación

```bash
# Re-sincronizar dependencias
go mod tidy
go build -o vibercode .
```

### Vector storage no disponible

- Normal si QDRANT_ENABLED=false
- El sistema funciona en modo stub sin problemas

### Graph storage no disponible

- Normal si NEO4J_ENABLED=false
- Todas las funciones devuelven datos mock

## 📝 Notas Importantes

1. **El sistema funciona perfectamente sin Qdrant o Neo4j** - están en modo stub
2. **Solo necesitas ANTHROPIC_API_KEY para funcionalidad de AI**
3. **Todos los endpoints funcionan correctamente** con datos mock si es necesario
4. **El preview server está totalmente funcional** para desarrollo

## 🎯 Próximos Pasos

1. Configura tu `ANTHROPIC_API_KEY`
2. Prueba `./vibercode vibe`
3. Experimenta con la generación de código
4. Si necesitas funcionalidades avanzadas, configura Qdrant/Neo4j
