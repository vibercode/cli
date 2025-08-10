# Changelog

All notable changes to ViberCode CLI will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Initial release of ViberCode CLI
- Go API generation with clean architecture
- Visual React editor integration
- AI chat functionality with Claude
- MCP server for AI agent integration
- WebSocket server for real-time communication
- Multi-database support (PostgreSQL, MySQL, SQLite, MongoDB)
- Docker and docker-compose generation
- Comprehensive test generation
- API documentation with Swagger
- Authentication and middleware systems
- Interactive CLI commands
- Example schemas and projects

### Features

#### 🎨 Visual Development

- **Full Vibe Mode**: Integrated visual editor with AI chat
- **Component-based UI**: Drag-and-drop interface components
- **Real-time sync**: Live updates between editor and generated code
- **Theme system**: Dynamic theme management
- **Responsive design**: Multi-device preview

#### ⚡ Code Generation

- **Clean Architecture**: Well-structured Go APIs
- **CRUD Operations**: Complete resource management
- **Database Integration**: Multiple provider support
- **Authentication**: JWT and session-based auth
- **API Documentation**: Auto-generated Swagger docs
- **Testing**: Unit and integration tests
- **Docker Support**: Container-ready projects

#### 🤖 AI Integration

- **Interactive Chat**: Claude-powered development assistant
- **MCP Protocol**: AI agent compatibility
- **Code Suggestions**: Context-aware recommendations
- **Template Enhancement**: AI-optimized code templates

#### 🔧 CLI Tools

- **Interactive Prompts**: User-friendly command interface
- **Schema Management**: JSON-based API definitions
- **Project Templates**: Quick-start examples
- **Development Server**: Hot-reload development mode

### Technical Details

#### Supported Databases

- PostgreSQL with advanced features
- MySQL with full compatibility
- SQLite for lightweight projects
- MongoDB for document-based APIs

#### Generated Architecture

```
cmd/server/          # Application entry point
internal/
├── handlers/        # HTTP layer (Gin framework)
├── services/        # Business logic layer
├── repositories/    # Data access layer
└── models/         # Domain models and DTOs
pkg/
├── database/       # Database utilities
├── config/         # Configuration management
└── utils/          # Shared utilities
```

#### Development Features

- Hot reload development server
- Comprehensive error handling
- Logging and monitoring setup
- Environment-based configuration
- Health check endpoints
- Graceful shutdown handling

### Commands Added

- `vibercode vibe` - Full development mode with visual editor
- `vibercode mcp` - MCP server for AI agents
- `vibercode ws` - WebSocket server
- `vibercode serve` - HTTP API server
- `vibercode generate api` - Complete API generation
- `vibercode generate resource` - CRUD resource generation
- `vibercode schema` - Schema management
- `vibercode run` - Project execution

### Documentation

- Complete English and Spanish documentation
- Interactive tutorials and examples
- API reference documentation
- Development guidelines
- Troubleshooting guides

## [1.0.0] - 2024-01-XX

### Added

- Initial stable release
- All core features implemented
- Production-ready code generation
- Community documentation
- Example projects and schemas

---

## Release Notes Format

Each release includes:

- **Added**: New features
- **Changed**: Changes in existing functionality
- **Deprecated**: Soon-to-be removed features
- **Removed**: Now removed features
- **Fixed**: Bug fixes
- **Security**: Security improvements

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for details on how to contribute to ViberCode CLI.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
