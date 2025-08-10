# Contributing to ViberCode CLI

🎉 Thank you for considering contributing to ViberCode CLI! This document provides guidelines and information for contributors.

## 🚀 Getting Started

### Prerequisites

- **Go 1.19+** for the CLI development
- **Node.js 16+** for the React editor
- **Git** for version control
- **pnpm** (preferred) or npm/yarn for JavaScript dependencies

### Development Setup

1. **Fork and clone the repository**

   ```bash
   git clone https://github.com/YOUR_USERNAME/cli.git
   cd cli
   ```

2. **Install Go dependencies**

   ```bash
   go mod download
   ```

3. **Build the CLI**

   ```bash
   go build -o vibercode .
   ```

4. **Run tests**
   ```bash
   go test ./...
   ```

## 🛠️ Development Workflow

### 1. Create a Feature Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

### 2. Make Your Changes

- Write clean, readable code
- Follow Go conventions and best practices
- Add tests for new functionality
- Update documentation as needed

### 3. Test Your Changes

```bash
# Run all tests
go test ./...

# Run tests with race detection
go test -race ./...

# Test CLI functionality
./vibercode --help
./vibercode generate api
```

### 4. Commit Your Changes

Follow conventional commit format:

```bash
git commit -m "feat: add new database provider support"
git commit -m "fix: resolve WebSocket connection issue"
git commit -m "docs: update installation instructions"
```

Commit types:

- `feat`: New features
- `fix`: Bug fixes
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Adding tests
- `chore`: Maintenance tasks

### 5. Submit a Pull Request

1. Push your branch to your fork
2. Create a pull request from your fork to the main repository
3. Fill out the PR template completely
4. Wait for review and address feedback

## 📋 Coding Standards

### Go Code Style

- Follow `gofmt` formatting
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions small and focused
- Handle errors appropriately

```go
// Good example
func GenerateAPI(config *Config) error {
    if config == nil {
        return errors.New("config cannot be nil")
    }

    // Implementation...
    return nil
}
```

### Project Structure

```
cmd/                    # CLI command definitions
internal/
├── generator/          # Code generation logic
├── models/            # Data structures
├── templates/         # Go template strings
├── mcp/              # MCP server implementation
├── vibe/             # Vibe mode implementation
└── websocket/        # WebSocket server
pkg/                   # Public packages
docs/                  # Documentation
examples/              # Usage examples
```

### Testing Guidelines

- Write unit tests for all new functions
- Use table-driven tests where appropriate
- Mock external dependencies
- Aim for good test coverage

```go
func TestGenerateAPI(t *testing.T) {
    tests := []struct {
        name    string
        config  *Config
        wantErr bool
    }{
        {
            name:    "valid config",
            config:  &Config{Name: "test"},
            wantErr: false,
        },
        {
            name:    "nil config",
            config:  nil,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := GenerateAPI(tt.config)
            if (err != nil) != tt.wantErr {
                t.Errorf("GenerateAPI() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## 🧩 Contributing Areas

### 1. Code Generation

- Adding new database providers
- Improving template systems
- Enhancing generated code quality

### 2. Visual Editor

- React component improvements
- UI/UX enhancements
- WebSocket communication

### 3. AI Integration

- Chat functionality improvements
- MCP server enhancements
- AI prompt optimization

### 4. CLI Commands

- New command implementations
- Command-line interface improvements
- Interactive prompts

### 5. Documentation

- Tutorial creation
- API documentation
- Example projects

## 🚨 Reporting Issues

### Bug Reports

Use the bug report template and include:

- Clear description of the issue
- Steps to reproduce
- Expected vs actual behavior
- Environment information
- Relevant logs/error messages

### Feature Requests

Use the feature request template and include:

- Clear description of the proposed feature
- Use case and motivation
- Proposed implementation approach
- Any alternatives considered

## 🔄 Release Process

1. **Version Bumping**: Follow semantic versioning (SemVer)
2. **Changelog**: Update CHANGELOG.md with new features/fixes
3. **Tagging**: Create Git tags for releases
4. **Automated Release**: GitHub Actions handles binary builds and releases

## 📞 Getting Help

- **Discord**: Join our community server
- **GitHub Discussions**: For questions and discussions
- **GitHub Issues**: For bug reports and feature requests
- **Email**: team@vibercode.com for sensitive matters

## 🏆 Recognition

Contributors will be:

- Listed in the README and release notes
- Mentioned in our Discord community
- Given credit in the project documentation

## 📄 License

By contributing to ViberCode CLI, you agree that your contributions will be licensed under the MIT License.

---

**Thank you for contributing to ViberCode CLI! 🚀**
