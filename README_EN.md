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

🚀 **ViberCode CLI** is a command-line tool for generating Go APIs with clean architecture, including a **visual React editor** and integrated AI chat.

## 🌟 **Main Features**

### 🎨 **Full Vibe Mode** (NEW!)

```bash
vibercode vibe
```

**A single command that starts:**

- 📡 **Real-time WebSocket Server**
- 🎨 **Visual React Editor** with modern interface
- 💬 **Integrated AI Chat** with Claude
- 🔄 **Live synchronization** between components
- 🌐 **Automatic browser opening**

### 🔌 **MCP Server** (NEW!)

```bash
vibercode mcp
```

**AI agent integration:**

- 🤖 Compatible with Claude Desktop and other MCP clients
- 🎨 Remote control of visual editor
- ⚡ Code generation via agents
- 🔄 Real-time component updates

### ⚡ **Code Generation**

- 🏗️ **Complete Go APIs** with clean architecture
- 📊 **CRUD Resources** with models, handlers, services, and repositories
- 🗄️ **Multi-database support** (PostgreSQL, MySQL, SQLite, MongoDB)
- 🐳 **Docker-ready** with auto-generated docker-compose

## 🚀 **Quick Installation**

### **Option 1: Automatic Script**

```bash
git clone <repo>
cd vibercode-cli-go
./install.sh
```

### **Option 2: Manual**

```bash
go build -o vibercode .
sudo cp ./vibercode /usr/local/bin/
```

### **Verify Installation**

```bash
vibercode --help
which vibercode  # Should show: /usr/local/bin/vibercode
```

## 🎯 **Quick Usage**

### **Full Mode (Recommended)**

```bash
# All-in-one! - Editor + WebSocket + AI Chat
vibercode vibe

# Component-focused mode
vibercode vibe component
```

### **Individual Services**

```bash
# WebSocket server only
vibercode ws

# MCP server only
vibercode mcp

# HTTP API server only
vibercode serve
```

### **Code Generation**

```bash
# Complete API
vibercode generate api

# Individual resource
vibercode generate resource
```

## 📋 **Available Commands**

| Command                       | Description                                  | Status     |
| ----------------------------- | -------------------------------------------- | ---------- |
| `vibercode vibe`              | 🎨 **Full mode** - Editor + Chat + WebSocket | ✅ **New** |
| `vibercode mcp`               | 🔌 **MCP Server** for AI agents              | ✅ **New** |
| `vibercode ws`                | 📡 WebSocket server for React Editor         | ✅         |
| `vibercode serve`             | 🌐 HTTP API server                           | ✅         |
| `vibercode generate api`      | ⚡ Generate complete Go API                  | ✅         |
| `vibercode generate resource` | 📦 Generate CRUD resource                    | ✅         |
| `vibercode schema`            | 📋 Manage resource schemas                   | ✅         |
| `vibercode run`               | 🚀 Run generated project                     | ✅         |

## 🛠️ **Development Workflow**

### **1. Full Development Mode**

```bash
$ vibercode vibe

🎨 Welcome to VibeCode Full Mode
📡 Starting WebSocket server on port 3001...
🎨 Starting React Editor...
📂 Found editor at: /path/to/vibercode/editor
🌐 Opening browser...
✅ VibeCode is ready!

💬 Viber AI: Hello! How can I help you?
```

### **2. Visual Development + Chat**

- 🎨 **Drag components** in the visual editor
- 💬 **Chat with AI**: "Add a blue button here"
- 🔄 **See real-time changes** in the browser
- ⚡ **Generate Go code** from visual schema

### **3. AI Agent Integration**

```bash
# Terminal 1: MCP Server
vibercode mcp

# Terminal 2: Vibe mode
vibercode vibe

# Now Claude Desktop can control your editor
```

## 🎨 **Visual Editor**

The React editor includes:

- 🧩 **Atomic components** (Button, Text, Input, etc.)
- 🏗️ **Molecular components** (Card, Form, Navigation)
- 🌊 **Organizational components** (Hero, Layout, Dashboard)
- 🎨 **Dynamic theme system**
- 📱 **Responsive view** (Desktop, Tablet, Mobile)
- 🔄 **Real-time synchronization** with WebSocket

## 🤖 **AI Integration**

### **Interactive Chat**

```
💬 User: "Add a red button in the top right corner"
🤖 Viber AI: Perfect! I've added a red button at position (500, 50).
```

### **MCP Agents**

```
🔌 Claude Desktop → MCP → ViberCode → React Editor
                    ↓
                   AI Chat ← WebSocket ← Live Updates
```

## 📊 **Architecture**

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

## 🔧 **Configuration**

### **Environment Variables**

```bash
# For full AI functionality
export ANTHROPIC_API_KEY=your_api_key

# For detailed debug
export VIBE_DEBUG=true

# Custom ports (optional)
export VIBE_WS_PORT=3001
export VIBE_EDITOR_PORT=5173
```

### **System Requirements**

- ✅ **Go 1.19+** for the CLI
- ✅ **Node.js 16+** for the React editor
- ✅ **pnpm/npm/yarn** for dependencies
- 🎯 **ANTHROPIC_API_KEY** for AI (optional)

## 🧪 **Testing**

```bash
# Test full mode
./test-vibe-full.sh

# Test MCP server
./test-mcp-server.sh

# Verify connections
curl http://localhost:3001/health   # WebSocket
curl http://localhost:5173         # React Editor
```

## 🐛 **Troubleshooting**

### **Editor won't start**

```bash
⚠️  Could not start React editor
💡 You can manually start it with: cd vibercode/editor && pnpm dev
```

### **Port already in use**

```bash
❌ WebSocket server error: address already in use
# Solution:
lsof -ti:3001 | xargs kill -9
```

### **Missing dependencies**

```bash
📦 Installing dependencies...
❌ failed to install dependencies
# Solution:
npm install -g pnpm
```

## 📚 **Complete Documentation**

### 🇺🇸 **English**

- 📖 [**Complete Documentation**](docs/en/README.md) - Full English guide
- 🚀 [**Quick Start**](docs/en/user-guide/quickstart.md) - Installation and first steps
- 💻 [**CLI Commands**](docs/en/user-guide/cli-commands.md) - Complete reference
- 🚨 [**Troubleshooting**](docs/en/troubleshooting/common-errors.md) - Common errors

### 🇪🇸 **Español**

- 📖 [**Documentación Completa**](docs/es/README.md) - Guía completa en español
- 🚀 [**Inicio Rápido**](docs/es/user-guide/quickstart.md) - Instalación y primeros pasos
- 💻 [**Comandos CLI**](docs/es/user-guide/cli-commands.md) - Referencia completa
- 🚨 [**Solución de Problemas**](docs/es/troubleshooting/common-errors.md) - Errores comunes

### 📂 **Technical Files**

- 🎨 [**CLAUDE.md**](CLAUDE.md) - Technical documentation for Claude Code
- 📋 [**Archived Documents**](archive/old-docs/) - Previous documentation

## 🤝 **Contributing**

We love receiving contributions from the community!

### 🚀 **Ways to Contribute**

- 🐛 **Report bugs** using [issue templates](.github/ISSUE_TEMPLATE/)
- ✨ **Propose new features**
- 📝 **Improve documentation**
- 🧹 **Code cleanup and optimizations**

### 📋 **Quick Process**

1. Fork the repository
2. Create branch: `git checkout -b feature/amazing-feature`
3. Commit: `git commit -m 'feat: add amazing feature'`
4. Push: `git push origin feature/amazing-feature`
5. Create Pull Request

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

## 🌟 **Contributors**

<a href="https://github.com/vibercode/cli/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=vibercode/cli" />
</a>

## 📄 **License**

This project is under the MIT license. See [LICENSE](LICENSE) for more details.

## 🔒 **Security**

To report security vulnerabilities, see [SECURITY.md](SECURITY.md).

---

**🚀 Build Go APIs with visual superpowers!**
