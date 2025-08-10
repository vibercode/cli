# ViberCode Full Mode - Vibe + React Editor

## 🎨 Descripción

El modo **ViberCode Full** combina el poder del chat AI con el editor visual React para una experiencia completa de desarrollo. Con un solo comando, obtienes:

- 📡 **Servidor WebSocket** en tiempo real
- 🎨 **Editor React Visual** con interfaz moderna
- 💬 **Chat AI integrado** con Claude
- 🔄 **Sincronización en vivo** entre todos los componentes
- 🌐 **Apertura automática** del navegador

## 🚀 Uso Rápido

```bash
# Iniciar modo completo
vibercode vibe

# Modo enfocado en componentes
vibercode vibe component
```

## 🛠️ Lo que hace automáticamente

### 1. **Detección Inteligente del Editor**

```
Busca el editor React en:
├── ../vibercode/editor          # Estructura estándar
├── ../../vibercode/editor       # Estructura anidada
├── vibercode/editor             # Directorio actual
├── ../editor                    # Directorio hermano
└── editor                       # Subdirectorio
```

### 2. **Instalación Automática de Dependencias**

```bash
# Intenta automáticamente:
pnpm install    # Preferido
npm install     # Fallback
yarn install    # Alternativo
```

### 3. **Inicio de Servicios en Paralelo**

```
📡 WebSocket Server → localhost:3001
🎨 React Editor     → localhost:5173
💬 AI Chat          → Terminal interactivo
🌐 Browser          → Apertura automática
```

## 📋 Requisitos

### Requeridos

- ✅ **ViberCode CLI** instalado globalmente
- ✅ **Node.js** (v16+)
- ✅ **pnpm/npm/yarn** para gestión de paquetes

### Opcionales

- 🎯 **ANTHROPIC_API_KEY** para funcionalidad AI completa
- 🌐 **Navegador moderno** para la mejor experiencia

## 🎯 Estructura del Proyecto

```
proyecto/
├── vibercode-cli-go/           # CLI en Go
│   ├── cmd/vibe.go            # Comando vibe
│   ├── internal/vibe/         # Lógica vibe
│   └── internal/websocket/    # Servidor WebSocket
└── vibercode/editor/          # Editor React
    ├── src/                   # Código fuente
    ├── package.json          # Dependencias
    └── vite.config.ts        # Configuración Vite
```

## 🔧 Configuración Personalizada

### Variables de Entorno

```bash
# Para chat AI completo
export ANTHROPIC_API_KEY=your_api_key

# Para debug detallado
export VIBE_DEBUG=true

# Puerto personalizado del WebSocket (opcional)
export VIBE_WS_PORT=3001

# Puerto personalizado del editor (opcional)
export VIBE_EDITOR_PORT=5173
```

### Archivo de Configuración

Crea `.vibercode/config.json`:

```json
{
  "vibe": {
    "auto_open_browser": true,
    "auto_install_deps": true,
    "editor_port": 5173,
    "websocket_port": 3001,
    "default_mode": "general"
  },
  "editor": {
    "theme": "dark",
    "auto_save": true,
    "live_reload": true
  }
}
```

## 💡 Modos de Operación

### Modo General (Por defecto)

```bash
vibercode vibe
```

- 🎯 **Enfoque**: Desarrollo completo de APIs Go + UI
- 🔧 **Características**: Generación de código, chat AI, editor visual
- 🎨 **UI**: Editor completo con todas las funcionalidades

### Modo Componente

```bash
vibercode vibe component
```

- 🎯 **Enfoque**: Diseño y edición de componentes UI
- 🔧 **Características**: Editor visual, chat enfocado en UI
- 🎨 **UI**: Interfaz optimizada para componentes

## 🔄 Flujo de Trabajo

### 1. **Inicio**

```bash
$ vibercode vibe

🎨 Welcome to VibeCode Full Mode
📡 Starting WebSocket server on port 3001...
🎨 Starting React Editor...
📂 Found editor at: /path/to/vibercode/editor
📦 Dependencies already installed
🌐 Opening browser...
✅ VibeCode is ready!
💬 Viber AI: ¡Hola! ¿En qué puedo ayudarte hoy?
```

### 2. **Desarrollo**

- 🎨 **Editor Visual**: Arrastra y suelta componentes
- 💬 **Chat AI**: "Agrega un botón azul en la esquina superior"
- 🔄 **Sincronización**: Los cambios se reflejan instantáneamente
- ⚡ **Live Reload**: El navegador se actualiza automáticamente

### 3. **Finalización**

```bash
Ctrl+C

🛑 Shutting down services...
🔌 Stopping WebSocket server...
🎨 Stopping React Editor...
✅ Shutdown complete
👋 ¡Hasta luego!
```

## 🐛 Troubleshooting

### Editor React no encontrado

```bash
⚠️  Could not start React editor: could not find React editor directory
💡 You can manually start it with: cd vibercode/editor && pnpm dev
```

**Solución**: Asegúrate de que el editor esté en una de las rutas esperadas.

### Puerto ocupado

```bash
❌ WebSocket server error: listen tcp :3001: bind: address already in use
```

**Solución**: Mata el proceso que usa el puerto o cambia el puerto:

```bash
# Matar proceso en puerto 3001
lsof -ti:3001 | xargs kill -9

# O cambiar puerto
export VIBE_WS_PORT=3002
```

### Dependencias no instaladas

```bash
📦 Installing dependencies...
❌ failed to install dependencies: no package manager found
```

**Solución**: Instala un gestor de paquetes:

```bash
# Instalar pnpm (recomendado)
npm install -g pnpm

# O usar npm que ya tienes
cd vibercode/editor && npm install
```

### Chat AI no responde

```bash
💬 Viber AI: [Error: ANTHROPIC_API_KEY not set]
```

**Solución**: Configura tu API key:

```bash
export ANTHROPIC_API_KEY=your_key_here
```

## 🎯 Comandos Útiles

### Desarrollo

```bash
# Inicio rápido
vibercode vibe

# Solo WebSocket (sin editor)
vibercode ws

# Solo chat (sin servicios)
vibercode vibe --chat-only

# Modo debug
VIBE_DEBUG=true vibercode vibe
```

### Testing

```bash
# Probar modo completo
./test-vibe-full.sh

# Probar solo MCP
./test-mcp-server.sh

# Probar conexión
curl -i -N -H "Connection: Upgrade" \
     -H "Upgrade: websocket" \
     -H "Sec-WebSocket-Key: test" \
     -H "Sec-WebSocket-Version: 13" \
     http://localhost:3001/ws
```

### Mantenimiento

```bash
# Actualizar binario global
sudo cp ./vibercode /usr/local/bin/

# Reinstalar dependencias del editor
cd vibercode/editor && rm -rf node_modules && pnpm install

# Limpiar puertos
lsof -ti:3001,5173 | xargs kill -9
```

## 🔗 Integración con MCP

El modo vibe funciona perfectamente con el servidor MCP:

```bash
# Terminal 1: Servidor MCP
vibercode mcp

# Terminal 2: Modo vibe completo
vibercode vibe

# Ahora puedes usar agentes IA que controlen el editor via MCP
```

## 📊 Métricas y Monitoreo

### Logs

```bash
# Ver logs en tiempo real
tail -f vibe.log

# Logs con debug
VIBE_DEBUG=true vibercode vibe 2>&1 | tee vibe-debug.log
```

### Estado de Servicios

```bash
# WebSocket
curl http://localhost:3001/health

# Editor React
curl http://localhost:5173

# Procesos activos
ps aux | grep -E "(vibercode|node.*vite)"
```

## 🎉 Próximas Características

- 🔄 **Hot Module Replacement** mejorado
- 🎨 **Temas personalizables** del editor
- 📱 **Preview móvil** integrado
- 🚀 **Deploy automático** a producción
- 🔌 **Plugins** de terceros
- 💾 **Auto-save** de proyectos

---

**¿Problemas?** Abre un issue en GitHub o consulta la documentación completa en `CLAUDE.md`.
