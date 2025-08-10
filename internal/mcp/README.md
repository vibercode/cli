# ViberCode MCP Server

## Descripción

El servidor MCP (Model Context Protocol) de ViberCode permite a agentes de IA interactuar directamente con ViberCode para edición de componentes en vivo, generación de código y gestión de proyectos.

## Características

### 🔌 Integración MCP

- Servidor MCP compatible con el estándar 2024-11-05
- Comunicación via stdin/stdout JSON-RPC
- Herramientas bien definidas con esquemas de validación

### 🎨 Edición de Componentes en Vivo

- Actualización de propiedades de componentes en tiempo real
- Modificación de posición y tamaño
- Gestión de temas y estilos

### 💬 Chat Integrado

- Envío de mensajes al asistente Viber AI
- Procesamiento de respuestas contextuales
- Integración con el chat bridge

### ⚡ Generación de Código

- Generación de APIs Go completas
- Soporte para múltiples bases de datos
- Configuración de características (auth, swagger, docker)

## Instalación

1. Compilar ViberCode CLI:

```bash
cd vibercode-cli-go
go build -o vibercode .
```

2. Configurar el servidor MCP en `.mcp.json`:

```json
{
  "mcpServers": {
    "vibercode": {
      "name": "ViberCode MCP Server",
      "description": "ViberCode CLI integration for live component editing",
      "command": "vibercode",
      "args": ["mcp"],
      "env": {
        "ANTHROPIC_API_KEY": "${ANTHROPIC_API_KEY}",
        "VIBE_DEBUG": "true"
      }
    }
  }
}
```

3. Configurar variables de entorno:

```bash
export ANTHROPIC_API_KEY=your_api_key
export VIBE_DEBUG=true  # Opcional para debug
```

## Herramientas Disponibles

### `vibe_start`

Inicia el modo vibe con chat AI y preview en vivo.

**Parámetros:**

- `mode` (string): "general" o "component" (default: "general")
- `port` (integer): Puerto para WebSocket (default: 3001)

**Ejemplo:**

```json
{
  "mode": "component",
  "port": 3001
}
```

### `component_update`

Actualiza las propiedades de un componente en tiempo real.

**Parámetros:**

- `componentId` (string): ID del componente
- `action` (string): "add", "update", o "remove"
- `properties` (object): Nuevas propiedades
- `position` (object): Nueva posición {x, y}
- `size` (object): Nuevo tamaño {w, h}

**Ejemplo:**

```json
{
  "componentId": "button_123",
  "action": "update",
  "properties": {
    "text": "Nuevo Texto",
    "variant": "primary",
    "color": "#FF0000"
  },
  "position": {
    "x": 100,
    "y": 200
  }
}
```

### `view_state_get`

Obtiene el estado actual de la vista y componentes.

**Parámetros:** Ninguno

**Respuesta:**

```json
{
  "components": [...],
  "theme": {...},
  "layout": {...},
  "canvas": {...},
  "status": "active"
}
```

### `chat_send`

Envía un mensaje al asistente Viber AI.

**Parámetros:**

- `message` (string): Mensaje para el asistente
- `context` (object): Contexto adicional

**Ejemplo:**

```json
{
  "message": "agregar un botón rojo en la esquina superior derecha",
  "context": {
    "current_theme": "dark",
    "active_components": ["button_1", "text_2"]
  }
}
```

### `generate_code`

Genera código Go API basado en un schema.

**Parámetros:**

- `project_name` (string): Nombre del proyecto
- `database` (string): "postgres", "mysql", "sqlite", "mongodb"
- `features` (array): ["auth", "swagger", "docker", "tests"]
- `schema` (object): Schema de recursos y modelos

**Ejemplo:**

```json
{
  "project_name": "mi_api",
  "database": "postgres",
  "features": ["auth", "swagger", "docker"],
  "schema": {
    "resources": [...],
    "models": [...]
  }
}
```

### `project_status`

Obtiene el estado del proyecto y servidores.

**Parámetros:** Ninguno

**Respuesta:**

```json
{
  "status": "running",
  "services": {
    "websocket": { "active": true, "port": 3001 },
    "http": { "active": false, "port": 8080 },
    "vibe": { "active": true }
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## Uso con Agentes de IA

### Claude con MCP

1. Configurar el servidor en el cliente MCP de Claude
2. El agente puede usar las herramientas directamente:

```
Usuario: "Agrega un botón azul y aumenta el tamaño del texto"

Claude: Voy a ayudarte a agregar un botón azul y modificar el texto.

[Usa herramienta component_update para agregar botón]
[Usa herramienta component_update para modificar texto]
```

### Flujo de Trabajo Típico

1. **Iniciar sesión:**

   ```json
   { "name": "vibe_start", "arguments": { "mode": "component" } }
   ```

2. **Obtener estado actual:**

   ```json
   { "name": "view_state_get", "arguments": {} }
   ```

3. **Actualizar componente:**

   ```json
   {
     "name": "component_update",
     "arguments": {
       "componentId": "button_1",
       "properties": { "color": "#0066FF" }
     }
   }
   ```

4. **Enviar mensaje al chat:**
   ```json
   {
     "name": "chat_send",
     "arguments": {
       "message": "El botón ahora es azul, ¿te gusta así?"
     }
   }
   ```

## Integración con React Editor

El servidor MCP se integra automáticamente con:

- **WebSocket Client** (`websocketClient.ts`): Para comunicación en tiempo real
- **API Client** (`apiClient.ts`): Para operaciones HTTP
- **Vibe Mode** (`vibe.go`): Para chat AI integrado
- **WebSocket Server** (`ws.go`): Para broadcasting de actualizaciones

## Debugging

### Logs Detallados

```bash
export VIBE_DEBUG=true
vibercode mcp
```

### Verificar Conexión

```bash
# Verificar que el servidor WebSocket esté activo
curl -i -N -H "Connection: Upgrade" \
     -H "Upgrade: websocket" \
     -H "Sec-WebSocket-Key: test" \
     -H "Sec-WebSocket-Version: 13" \
     http://localhost:3001/ws
```

### Monitorear Mensajes

Los mensajes MCP se pueden monitorear en tiempo real observando stdin/stdout del proceso.

## Limitaciones Actuales

- El broadcast directo desde MCP está en desarrollo (TODO)
- Algunas integraciones con vibe session necesitan refinamiento
- El manejo de errores se puede mejorar

## Contribuir

1. Fork el repositorio
2. Crear rama para feature: `git checkout -b feature/mcp-improvement`
3. Commit cambios: `git commit -am 'Mejora en servidor MCP'`
4. Push a la rama: `git push origin feature/mcp-improvement`
5. Crear Pull Request

## Licencia

Este proyecto sigue la misma licencia que ViberCode CLI.
