# VibeCode API Documentation

## Información General

**Servidor**: Preview Server con Chat Bridge integrado  
**Puerto**: 3001  
**Base URL**: `http://localhost:3001`  
**WebSocket**: `ws://localhost:3001/ws`  
**Versión**: 1.0.0

## Configuración

### Variables de Entorno

```bash
# Requerida para funcionalidad AI
export ANTHROPIC_API_KEY=your_api_key

# Opcional - para debug detallado
export VIBE_DEBUG=true
```

### Inicialización

```bash
# Iniciar modo vibe
vibercode vibe

# Modo componente (enfoque UI)
vibercode vibe --component
```

---

## HTTP API Endpoints

### 1. Página Principal

```http
GET /
```

**Descripción**: Sirve la interfaz web del preview con chat integrado  
**Respuesta**: HTML page con WebSocket client

---

### 2. Estado del Servidor

```http
GET /api/status
```

**Descripción**: Información del estado del servidor

**Respuesta**:

```json
{
  "status": "ok",
  "server": "vibe-preview",
  "version": "1.0.0",
  "connected_clients": 2,
  "timestamp": "2024-01-15T10:30:00Z",
  "prompts_loaded": true,
  "components_count": 5,
  "current_theme": "VibeCode",
  "current_viewport": "desktop"
}
```

---

### 3. Estado de la Vista

```http
GET /api/view-state
```

**Descripción**: Obtiene el estado actual de la vista

**Respuesta**:

```json
{
  "components": [
    {
      "id": "button_1642234567",
      "type": "button",
      "name": "Button_1642234567",
      "category": "Input",
      "properties": {
        "text": "Click me",
        "variant": "primary",
        "size": "medium"
      },
      "position": { "x": 100, "y": 100 },
      "size": { "w": 160, "h": 40 }
    }
  ],
  "theme": {
    "id": "vibercode",
    "name": "VibeCode",
    "colors": {
      "primary": "#8B5CF6",
      "secondary": "#06B6D4",
      "accent": "#F59E0B",
      "background": "#0F0F0F",
      "surface": "#1A1A1A",
      "text": "#FFFFFF"
    },
    "effects": {
      "glow": true,
      "gradients": true,
      "animations": true
    }
  },
  "layout": {
    "grid": 12,
    "row_height": 60,
    "margin": [16, 16],
    "container_padding": [24, 24],
    "show_grid": true,
    "snap_to_grid": true
  },
  "canvas": {
    "viewport": "desktop",
    "zoom": 1.0,
    "pan_offset": { "x": 0, "y": 0 },
    "selected_item": ""
  }
}
```

---

### 4. Actualizar Vista

```http
POST /api/view-update
```

**Descripción**: Actualiza componentes de la vista

**Payload**:

```json
{
  "components": [
    {
      "id": "button_1",
      "type": "button",
      "properties": {
        "text": "Updated Button",
        "variant": "secondary"
      }
    }
  ],
  "layout": {
    "grid": 12
  },
  "theme": {
    "name": "Updated Theme"
  }
}
```

**Respuesta**:

```json
{
  "status": "success"
}
```

---

### 5. Actualización en Vivo

```http
POST /api/live-update
```

**Descripción**: Actualiza un componente específico en tiempo real

**Payload**:

```json
{
  "componentId": "button_1",
  "action": "update",
  "changes": {
    "properties": {
      "text": "New Text"
    },
    "position": {
      "x": 200,
      "y": 150
    }
  }
}
```

**Respuesta**:

```json
{
  "status": "success"
}
```

---

### 6. Actualizar Estado Completo

```http
POST /api/view-state
```

**Descripción**: Actualiza el estado completo de la vista

**Payload**:

```json
{
  "components": [...],
  "theme": {...},
  "layout": {...},
  "canvas": {...}
}
```

**Respuesta**:

```json
{
  "status": "success"
}
```

---

### 7. Chat con Viber AI

```http
POST /api/chat
```

**Descripción**: Envía mensaje al asistente Viber

**Payload**:

```json
{
  "message": "agregar un botón rojo",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**Respuesta**:

```json
{
  "response": "¡Perfecto! He agregado un botón rojo en la posición (100, 100).",
  "action": "add_component",
  "data": {
    "id": "button_red_1642234567",
    "type": "button",
    "properties": {
      "text": "Botón Rojo",
      "variant": "danger",
      "color": "#EF4444"
    }
  },
  "timestamp": "2024-01-15T10:30:01Z"
}
```

---

## WebSocket API

### Conexión

```javascript
const ws = new WebSocket("ws://localhost:3001/ws");
```

### Configuración

- **Timeout de lectura**: 60 segundos
- **Timeout de escritura**: 10 segundos
- **Tamaño máximo de mensaje**: 512KB
- **Heartbeat**: Ping cada 30 segundos

---

## Tipos de Mensajes WebSocket

### 1. Estructura Base

```json
{
  "type": "message_type",
  "action": "optional_action",
  "data": {},
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### 2. Tipos de Mensaje

#### `view_state_update`

**Descripción**: Actualización del estado de la vista  
**Dirección**: Servidor → Cliente

```json
{
  "type": "view_state_update",
  "data": {
    "components": [...],
    "theme": {...},
    "layout": {...},
    "canvas": {...}
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

#### `view_update`

**Descripción**: Actualización de componentes  
**Dirección**: Cliente → Servidor

```json
{
  "type": "view_update",
  "data": {
    "components": [...],
    "layout": {...},
    "theme": {...}
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

#### `live_update`

**Descripción**: Actualización en tiempo real  
**Dirección**: Bidireccional

```json
{
  "type": "live_update",
  "data": {
    "componentId": "button_1",
    "action": "update",
    "changes": {
      "properties": {
        "text": "New Text"
      }
    }
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

#### `chat_message`

**Descripción**: Mensaje de chat  
**Dirección**: Cliente → Servidor

```json
{
  "type": "chat_message",
  "data": {
    "message": "agregar un botón azul",
    "timestamp": "2024-01-15T10:30:00Z"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

#### `chat_response`

**Descripción**: Respuesta del chat  
**Dirección**: Servidor → Cliente

```json
{
  "type": "chat_response",
  "data": {
    "user_message": {
      "id": "msg_1642234567",
      "type": "user",
      "content": "agregar un botón azul",
      "source": "websocket",
      "timestamp": "2024-01-15T10:30:00Z"
    },
    "assistant_response": {
      "id": "msg_1642234568",
      "content": "¡Perfecto! He agregado un botón azul.",
      "action": "add_component",
      "data": {...},
      "timestamp": "2024-01-15T10:30:01Z"
    }
  },
  "timestamp": "2024-01-15T10:30:01Z"
}
```

#### `ping` / `pong`

**Descripción**: Mensajes de heartbeat  
**Dirección**: Bidireccional

```json
{
  "type": "ping",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

## Tipos de Datos

### ComponentState

```json
{
  "id": "string",
  "type": "string",
  "name": "string",
  "category": "string",
  "properties": {
    "key": "value"
  },
  "position": {
    "x": 0,
    "y": 0
  },
  "size": {
    "w": 0,
    "h": 0
  }
}
```

### ThemeState

```json
{
  "id": "string",
  "name": "string",
  "colors": {
    "primary": "#8B5CF6",
    "secondary": "#06B6D4",
    "accent": "#F59E0B",
    "background": "#0F0F0F",
    "surface": "#1A1A1A",
    "text": "#FFFFFF"
  },
  "effects": {
    "glow": true,
    "gradients": true,
    "animations": true
  }
}
```

### LayoutState

```json
{
  "grid": 12,
  "row_height": 60,
  "margin": [16, 16],
  "container_padding": [24, 24],
  "show_grid": true,
  "snap_to_grid": true
}
```

### CanvasState

```json
{
  "viewport": "desktop",
  "zoom": 1.0,
  "pan_offset": {
    "x": 0,
    "y": 0
  },
  "selected_item": ""
}
```

---

## Ejemplos de Uso

### 1. Conectar WebSocket

```javascript
const ws = new WebSocket("ws://localhost:3001/ws");

ws.onopen = function () {
  console.log("Conectado al servidor VibeCode");
};

ws.onmessage = function (event) {
  const data = JSON.parse(event.data);
  console.log("Mensaje recibido:", data);

  switch (data.type) {
    case "view_state_update":
      updateUI(data.data);
      break;
    case "chat_response":
      displayChatResponse(data.data);
      break;
    case "live_update":
      updateComponent(data.data);
      break;
  }
};
```

### 2. Enviar Mensaje de Chat

```javascript
function sendChatMessage(message) {
  const chatData = {
    type: "chat_message",
    data: {
      message: message,
      timestamp: new Date().toISOString(),
    },
    timestamp: new Date().toISOString(),
  };

  ws.send(JSON.stringify(chatData));
}

// Usar
sendChatMessage("agregar un botón verde");
```

### 3. Actualizar Componente

```javascript
function updateComponent(componentId, changes) {
  const updateData = {
    type: "live_update",
    data: {
      componentId: componentId,
      action: "update",
      changes: changes,
    },
    timestamp: new Date().toISOString(),
  };

  ws.send(JSON.stringify(updateData));
}

// Usar
updateComponent("button_1", {
  properties: {
    text: "Texto Actualizado",
    variant: "success",
  },
});
```

### 4. Obtener Estado Actual

```javascript
async function getCurrentState() {
  const response = await fetch("http://localhost:3001/api/view-state");
  const state = await response.json();
  return state;
}
```

---

## Flujo de Comunicación

### 1. Inicialización

```
Cliente → Servidor: WebSocket connection
Servidor → Cliente: view_state_update (estado inicial)
```

### 2. Chat Bidireccional

```
Terminal Chat → Chat Bridge → WebSocket Clients
                    ↓
              Viber AI Processing
                    ↓
WebSocket Clients ← Chat Bridge ← Terminal Chat
```

### 3. Actualización de Componentes

```
Cliente → Servidor: view_update/live_update
Servidor: Procesa y actualiza estado
Servidor → Todos los Clientes: broadcast update
```

### 4. Heartbeat

```
Servidor → Cliente: ping (cada 30s)
Cliente → Servidor: pong
```

---

## Códigos de Respuesta HTTP

| Código | Descripción                |
| ------ | -------------------------- |
| 200    | Éxito                      |
| 400    | Solicitud incorrecta       |
| 404    | Endpoint no encontrado     |
| 405    | Método no permitido        |
| 500    | Error interno del servidor |

---

## CORS

Todos los endpoints HTTP incluyen headers CORS:

```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization, X-Requested-With
Access-Control-Max-Age: 86400
```

---

## Comandos de Chat Soportados

### Componentes Básicos

- `"agregar botón"` - Crea un botón
- `"agregar texto"` - Añade texto
- `"agregar imagen"` - Inserta imagen
- `"agregar input"` - Campo de entrada

### Componentes Compuestos

- `"agregar tarjeta"` - Crea una tarjeta
- `"agregar formulario"` - Formulario
- `"agregar navegación"` - Navegación

### Temas

- `"cambiar a azul/rojo/verde"`
- `"tema oscuro/claro"`
- `"océano/puesta de sol"`

### Información

- `"estado"` - Estado del canvas
- `"ayuda"` - Ayuda completa

---

## Manejo de Errores

### WebSocket

```javascript
ws.onerror = function (error) {
  console.error("WebSocket error:", error);
};

ws.onclose = function () {
  console.log("Conexión cerrada, reconectando...");
  setTimeout(connectWS, 3000);
};
```

### HTTP

```javascript
try {
  const response = await fetch("/api/endpoint");
  if (!response.ok) {
    throw new Error(`HTTP ${response.status}: ${response.statusText}`);
  }
  const data = await response.json();
} catch (error) {
  console.error("Error:", error);
}
```

---

## Características Avanzadas

### Chat Bridge

- Sincronización terminal ↔ web
- Múltiples clientes WebSocket
- Procesamiento unificado con Viber AI

### Logging

- Modo silencioso en chat
- Logs detallados en `vibe.log`
- Debug mode con `VIBE_DEBUG=true`

### Personalización

- Temas dinámicos
- Posicionamiento inteligente
- Responsive design

---

## Troubleshooting

### WebSocket no conecta

1. Verificar que el servidor esté en puerto 3001
2. Comprobar firewall/proxy
3. Verificar CORS en desarrollo

### Chat no responde

1. Verificar `ANTHROPIC_API_KEY`
2. Comprobar logs en `vibe.log`
3. Reiniciar servidor si es necesario

### Componentes no se actualizan

1. Verificar conexión WebSocket
2. Comprobar formato de mensajes
3. Verificar estado del servidor con `/api/status`

---

## Endpoints Adicionales (Enhanced Server)

Para proyectos con vector/graph storage habilitado:

```http
POST /api/semantic-search
GET /api/related-components/{id}
GET /api/relationship-insights
GET /api/project-stats
```

Estos endpoints requieren configuración adicional del vector graph service.

---

**Última actualización**: 2024-01-15  
**Versión de la API**: 1.0.0  
**Soporte**: VibeCode Team
