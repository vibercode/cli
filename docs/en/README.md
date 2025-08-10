# ViberCode CLI - English Documentation

Welcome to the complete documentation for ViberCode CLI, a command-line tool for generating Go APIs with clean architecture.

## 📚 Documentation Contents

### User Guide
- [**Quick Start**](user-guide/quickstart.md) - Installation and first steps
- [**CLI Commands**](user-guide/cli-commands.md) - Complete command reference
- [**Schema Generation**](user-guide/schema-generation.md) - Create models and APIs
- [**Configuration**](user-guide/configuration.md) - Project configuration

### Tutorials
- [**Your First Complete API**](tutorials/first-complete-api.md) - Complete step-by-step tutorial ⭐
- [**Database Integration**](tutorials/database-integration.md) - Connect with different DBs
- [**Authentication & Authorization**](tutorials/auth-tutorial.md) - Implement security

### API & Development
- [**Project Architecture**](api/architecture.md) - Structure and patterns
- [**Templates & Generators**](api/templates.md) - Template system
- [**Extensions**](api/extensions.md) - Create custom extensions

### Troubleshooting
- [**Common Errors**](troubleshooting/common-errors.md) - Frequent issues
- [**Debugging**](troubleshooting/debugging.md) - Debug tools
- [**FAQ**](troubleshooting/faq.md) - Frequently asked questions

## 🚀 Quick Start

```bash
# Install ViberCode CLI
go install github.com/vibercode/cli@latest

# Generate a new API project
vibercode generate api my-project

# Generate a user schema
vibercode schema generate User -m my-module -d postgres
```

## 🔧 Key Features

- ✅ **Automatic code generation** - Complete APIs with CRUD
- ✅ **Clean architecture** - Clear separation of concerns
- ✅ **Multiple databases** - PostgreSQL, MySQL, SQLite, MongoDB
- ✅ **Customizable templates** - Adapt code to your needs
- ✅ **MCP Integration** - Model Context Protocol server
- ✅ **Interactive AI Chat** - Integrated development assistant

## 📖 Useful Links

- [**Documentación en Español**](../es/README.md) - Spanish documentation
- [**GitHub Repository**](https://github.com/vibercode/cli) - Source code
- [**Issues & Support**](https://github.com/vibercode/cli/issues) - Report issues

---

*Generated with ViberCode CLI 🚀*