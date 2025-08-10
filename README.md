# ViberCode CLI

<div align="center">

![ViberCode CLI Logo](https://via.placeholder.com/200x200/0066CC/FFFFFF?text=ViberCode)

**🚀 Generate Go APIs with Clean Architecture + Visual Editor + AI Chat**

[![Go Version](https://img.shields.io/badge/Go-1.19+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Release](https://img.shields.io/github/v/release/vibercode/cli?style=flat&logo=github)](https://github.com/vibercode/cli/releases)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat)](LICENSE)
[![CI](https://github.com/vibercode/cli/workflows/CI/badge.svg)](https://github.com/vibercode/cli/actions)

[English](README_EN.md) • [Español](README.md) • [Documentation](docs/INDEX.md)

</div>

---

🚀 **ViberCode CLI** es una herramienta de línea de comandos para generar APIs Go con arquitectura limpia, incluyendo un **editor visual React** y chat AI integrado.

## 🌟 **Características Principales**

### 🎨 **Modo Vibe Completo** (¡NUEVO!)

```bash
vibercode vibe
```

**Un solo comando que inicia:**

- 📡 **Servidor WebSocket** en tiempo real
- 🎨 **Editor React Visual** con interfaz moderna
- 💬 **Chat AI integrado** con Claude
- 🔄 **Sincronización en vivo** entre componentes
- 🌐 **Apertura automática** del navegador

### 🔌 **Servidor MCP** (¡NUEVO!)

```bash
vibercode mcp
```

**Integración con agentes IA:**

- 🤖 Compatible con Claude Desktop y otros clientes MCP
- 🎨 Control remoto del editor visual
- ⚡ Generación de código via agentes
- 🔄 Actualización de componentes en tiempo real

### ⚡ **Generación de Código**

- 🏗️ **APIs Go completas** con arquitectura limpia
- 📊 **Recursos CRUD** con modelos, handlers, servicios y repositorios
- 🗄️ **Soporte multi-database** (PostgreSQL, MySQL, SQLite, MongoDB)
- 🐳 **Docker-ready** con auto-generación de docker-compose

## 🚀 **Instalación Rápida**

### **Opción 1: Script Automático**

```bash
git clone <repo>
cd vibercode-cli-go
./install.sh
```

### **Opción 2: Manual**

```bash
go build -o vibercode .
sudo cp ./vibercode /usr/local/bin/
```

### **Verificar Instalación**

```bash
vibercode --help
which vibercode  # Debería mostrar: /usr/local/bin/vibercode
```

## 🎯 **Uso Rápido**

### **Modo Completo (Recomendado)**

```bash
# ¡Todo en uno! - Editor + WebSocket + Chat AI
vibercode vibe

# Modo enfocado en componentes
vibercode vibe component
```

### **Servicios Individuales**

```bash
# Solo servidor WebSocket
vibercode ws

# Solo servidor MCP
vibercode mcp

# Solo servidor HTTP API
vibercode serve
```

### **Generación de Código**

```bash
# API completa
vibercode generate api

# Recurso individual
vibercode generate resource
```

## 📋 **Comandos Disponibles**

| Comando                       | Descripción                                      | Estado       |
| ----------------------------- | ------------------------------------------------ | ------------ |
| `vibercode vibe`              | 🎨 **Modo completo** - Editor + Chat + WebSocket | ✅ **Nuevo** |
| `vibercode mcp`               | 🔌 **Servidor MCP** para agentes IA              | ✅ **Nuevo** |
| `vibercode ws`                | 📡 Servidor WebSocket para React Editor          | ✅           |
| `vibercode serve`             | 🌐 Servidor HTTP API                             | ✅           |
| `vibercode generate api`      | ⚡ Generar API Go completa                       | ✅           |
| `vibercode generate resource` | 📦 Generar recurso CRUD                          | ✅           |
| `vibercode schema`            | 📋 Gestionar esquemas de recursos                | ✅           |
| `vibercode run`               | 🚀 Ejecutar proyecto generado                    | ✅           |

## 🛠️ **Flujo de Desarrollo**

### **1. Modo Desarrollo Completo**

```bash
$ vibercode vibe

🎨 Welcome to VibeCode Full Mode
📡 Starting WebSocket server on port 3001...
🎨 Starting React Editor...
📂 Found editor at: /path/to/vibercode/editor
🌐 Opening browser...
✅ VibeCode is ready!

💬 Viber AI: ¡Hola! ¿En qué puedo ayudarte?
```

### **2. Desarrollo Visual + Chat**

- 🎨 **Arrastra componentes** en el editor visual
- 💬 **Chatea con AI**: "Agrega un botón azul aquí"
- 🔄 **Ve cambios en tiempo real** en el navegador
- ⚡ **Genera código Go** desde el esquema visual

### **3. Integración con Agentes IA**

```bash
# Terminal 1: Servidor MCP
vibercode mcp

# Terminal 2: Modo vibe
vibercode vibe

# Ahora Claude Desktop puede controlar tu editor
```

## 🎨 **Editor Visual**

El editor React incluye:

- 🧩 **Componentes atomicos** (Button, Text, Input, etc.)
- 🏗️ **Componentes moleculares** (Card, Form, Navigation)
- 🌊 **Componentes organizacionales** (Hero, Layout, Dashboard)
- 🎨 **Sistema de temas** dinámico
- 📱 **Vista responsive** (Desktop, Tablet, Mobile)
- 🔄 **Sincronización en tiempo real** con WebSocket

## 🤖 **Integración IA**

### **Chat Interactivo**

```
💬 Usuario: "Agrega un botón rojo en la esquina superior derecha"
🤖 Viber AI: ¡Perfecto! He agregado un botón rojo en la posición (500, 50).
```

### **Agentes MCP**

```
🔌 Claude Desktop → MCP → ViberCode → Editor React
                    ↓
                   Chat AI ← WebSocket ← Live Updates
```

## 📊 **Arquitectura**

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   AI Agents     │    │   ViberCode     │    │   React Editor  │
│   (Claude MCP)  │◄──►│   CLI Server    │◄──►│   (localhost)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │   Generated     │
                       │   Go API Code   │
                       └─────────────────┘
```

## 🔧 **Configuración**

### **Variables de Entorno**

```bash
# Para funcionalidad AI completa
export ANTHROPIC_API_KEY=your_api_key

# Para debug detallado
export VIBE_DEBUG=true

# Puertos personalizados (opcional)
export VIBE_WS_PORT=3001
export VIBE_EDITOR_PORT=5173
```

### **Requisitos del Sistema**

- ✅ **Go 1.19+** para el CLI
- ✅ **Node.js 16+** para el editor React
- ✅ **pnpm/npm/yarn** para dependencias
- 🎯 **ANTHROPIC_API_KEY** para AI (opcional)

## 🧪 **Testing**

```bash
# Probar modo completo
./test-vibe-full.sh

# Probar servidor MCP
./test-mcp-server.sh

# Verificar conexiones
curl http://localhost:3001/health   # WebSocket
curl http://localhost:5173         # React Editor
```

## 🐛 **Troubleshooting**

### **Editor no inicia**

```bash
⚠️  Could not start React editor
💡 You can manually start it with: cd vibercode/editor && pnpm dev
```

### **Puerto ocupado**

```bash
❌ WebSocket server error: address already in use
# Solución:
lsof -ti:3001 | xargs kill -9
```

### **Dependencias faltantes**

```bash
📦 Installing dependencies...
❌ failed to install dependencies
# Solución:
npm install -g pnpm
```

## 📚 **Documentación Completa**

### 🇪🇸 **Español**

- 📖 [**Documentación Completa**](docs/es/README.md) - Guía completa en español
- 🚀 [**Inicio Rápido**](docs/es/user-guide/quickstart.md) - Instalación y primeros pasos
- 💻 [**Comandos CLI**](docs/es/user-guide/cli-commands.md) - Referencia completa
- 🚨 [**Solución de Problemas**](docs/es/troubleshooting/common-errors.md) - Errores comunes

### 🇺🇸 **English**

- 📖 [**Complete Documentation**](docs/en/README.md) - Full English guide
- 🚀 [**Quick Start**](docs/en/user-guide/quickstart.md) - Installation and first steps
- 💻 [**CLI Commands**](docs/en/user-guide/cli-commands.md) - Complete reference
- 🚨 [**Troubleshooting**](docs/en/troubleshooting/common-errors.md) - Common errors

### 📂 **Archivos Técnicos**

- 🎨 [**CLAUDE.md**](CLAUDE.md) - Documentación técnica para Claude Code
- 📋 [**Documentos Archivados**](archive/old-docs/) - Documentación anterior

## 🤝 **Contribuir**

¡Nos encanta recibir contribuciones de la comunidad!

### 🚀 **Formas de Contribuir**

- 🐛 **Reportar bugs** usando los [issue templates](.github/ISSUE_TEMPLATE/)
- ✨ **Proponer nuevas funcionalidades**
- 📝 **Mejorar documentación**
- 🧹 **Limpiar código y optimizaciones**

### 📋 **Proceso Rápido**

1. Fork el repositorio
2. Crear rama: `git checkout -b feature/amazing-feature`
3. Commit: `git commit -m 'feat: add amazing feature'`
4. Push: `git push origin feature/amazing-feature`
5. Crear Pull Request

Ver [CONTRIBUTING.md](CONTRIBUTING.md) para guías detalladas.

## 🌟 **Contributors**

<a href="https://github.com/vibercode/cli/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=vibercode/cli" />
</a>

## 📄 **Licencia**

Este proyecto está bajo la licencia MIT. Ver [LICENSE](LICENSE) para más detalles.

## 🔒 **Seguridad**

Para reportar vulnerabilidades de seguridad, ver [SECURITY.md](SECURITY.md).

---

**🚀 ¡Construye APIs Go con superpoderes visuales!**
