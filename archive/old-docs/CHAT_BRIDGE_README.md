# Chat Bridge - Sistema de Comunicación VibeCode

## Descripción General

El Chat Bridge es un sistema de comunicación bidireccional que conecta el chat terminal con el WebSocket del preview server, permitiendo una experiencia de chat unificada entre la terminal y el navegador.

## Características

### 🔇 Logging Limpio

- **Modo Silencioso**: En modo vibe, los logs se reducen para no interrumpir el flujo del chat
- **Archivo de Logs**: Los logs se escriben a `vibe.log` en lugar de interferir con la consola
- **Modo Debug**: Activar con `VIBE_DEBUG=true` para logs detallados

### 🌉 Chat Bridge

- **Comunicación Bidireccional**: Los mensajes del terminal se sincronizan con el WebSocket
- **Múltiples Clientes**: Soporte para múltiples clientes WebSocket simultáneos
- **Estado Compartido**: El estado del chat se mantiene consistente entre terminal y web

### 💬 Chat Web

- **Interfaz Integrada**: Panel de chat en el preview server
- **Tiempo Real**: Comunicación instantánea via WebSocket
- **Historial**: Mantiene historial de conversaciones

## Arquitectura

```
Terminal Chat ←→ Chat Bridge ←→ WebSocket Clients
                      ↓
                 Viber AI
                      ↓
                 Preview Server
```

## Componentes

### 1. VibeLogger (`internal/vibe/logger.go`)

- Maneja logging con diferentes niveles de verbosidad
- Modo chat reduce logs para no interferir con la experiencia
- Soporte para archivo de logs

### 2. ChatBridge (`internal/vibe/chat_bridge.go`)

- Conecta terminal con WebSocket
- Maneja cola de mensajes bidireccional
- Procesa mensajes con Viber AI
- Mantiene estado de clientes WebSocket

### 3. Preview Server Actualizado (`internal/vibe/preview.go`)

- Interfaz web con panel de chat
- Manejo de WebSocket mejorado
- Logging silencioso en modo chat

## Uso

### Iniciar Modo Vibe

```bash
vibercode vibe
```

### Comandos de Debug

```bash
# Activar logs detallados
export VIBE_DEBUG=true
vibercode vibe

# Usar modo componente
vibercode vibe --component
```

### Interfaz Web

- Abrir `http://localhost:3001` para ver el preview
- Panel de chat integrado en la derecha
- Comunicación en tiempo real con el terminal

## Flujo de Mensajes

### Terminal → Web

1. Usuario escribe en terminal
2. `ChatBridge.SendTerminalMessage()`
3. Viber AI procesa el mensaje
4. Respuesta se envía a todos los clientes WebSocket
5. Respuesta se muestra en terminal

### Web → Terminal

1. Usuario escribe en interfaz web
2. Mensaje se envía via WebSocket
3. `ChatBridge.HandleWebSocketMessage()`
4. Viber AI procesa el mensaje
5. Respuesta se envía a web y notificación en terminal

## Configuración

### Variables de Entorno

```bash
# Clave API de Viber (requerida)
export ANTHROPIC_API_KEY=your_api_key

# Activar modo debug
export VIBE_DEBUG=true
```

### Modos de Operación

- **General**: Chat sobre APIs Go y componentes UI
- **Component**: Enfoque exclusivo en componentes UI

## Estructura de Mensajes

### ChatMessage

```go
type ChatMessage struct {
    ID        string                 `json:"id"`
    Type      string                 `json:"type"` // "user", "assistant", "system"
    Content   string                 `json:"content"`
    Timestamp time.Time              `json:"timestamp"`
    Source    string                 `json:"source"` // "terminal", "websocket"
    Data      map[string]interface{} `json:"data,omitempty"`
}
```

### ChatResponse

```go
type ChatResponse struct {
    ID        string                 `json:"id"`
    Content   string                 `json:"content"`
    Action    string                 `json:"action,omitempty"`
    Data      map[string]interface{} `json:"data,omitempty"`
    Timestamp time.Time              `json:"timestamp"`
}
```

## Ejemplos de Uso

### Chat Terminal

```bash
# Iniciar vibe mode
vibercode vibe

# Chatear normalmente
💬 You: agregar un botón rojo
🤖 Viber: ¡Perfecto! Agregando un botón rojo...

# Ver mensajes del chat web
📱 Chat Web: cambiar tema a azul
🤖 Viber: Cambiando el tema a azul...
```

### Chat Web

- Escribir mensajes en el panel de chat
- Ver respuestas en tiempo real
- Sincronizado con el terminal

## Beneficios

1. **Experiencia Unificada**: Chat funciona tanto en terminal como en web
2. **Logs Limpios**: No hay interferencia de logs durante el chat
3. **Tiempo Real**: Comunicación instantánea
4. **Escalable**: Soporte para múltiples clientes
5. **Robusto**: Manejo de errores y reconexión automática

## Archivos Principales

- `internal/vibe/logger.go` - Sistema de logging
- `internal/vibe/chat_bridge.go` - Bridge de comunicación
- `internal/vibe/preview.go` - Servidor con interfaz web
- `internal/vibe/chat.go` - Manager del chat terminal
- `internal/vibe/vibe.go` - Punto de entrada del modo vibe

## Troubleshooting

### Problema: Logs interferiendo con el chat

**Solución**: El sistema automáticamente reduce logs en modo chat

### Problema: WebSocket no conecta

**Solución**: Verificar que el preview server esté corriendo en puerto 3001

### Problema: Chat no responde

**Solución**: Verificar que `ANTHROPIC_API_KEY` esté configurada

### Problema: Mensajes no se sincronizan

**Solución**: Verificar conexión WebSocket y reiniciar el servidor

## Desarrollo

### Agregar Nuevas Funcionalidades

1. Extender `ChatBridge` con nuevos métodos
2. Actualizar interfaz web en `preview.go`
3. Agregar manejo de mensajes en JavaScript
4. Documentar cambios

### Testing

```bash
# Ejecutar con debug
export VIBE_DEBUG=true
vibercode vibe

# Verificar logs
tail -f vibe.log

# Probar WebSocket
# Abrir navegador en localhost:3001
```
