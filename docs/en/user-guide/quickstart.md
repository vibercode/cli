# Quick Start - ViberCode CLI

This guide will help you get started with ViberCode CLI in just a few minutes.

## 🚦 Prerequisites

- Go 1.19 or higher
- Git
- A code editor (VS Code, GoLand, etc.)

## 📦 Installation

### Option 1: Install from source

```bash
# Clone the repository
git clone https://github.com/vibercode/cli.git
cd vibercode-cli-go

# Build the binary
go build -o vibercode main.go

# Make executable (Linux/macOS)
chmod +x vibercode

# Move to PATH (optional)
sudo mv vibercode /usr/local/bin/
```

### Option 2: Direct installation with Go

```bash
go install github.com/vibercode/cli@latest
```

## 🎯 Your First API

### 1. Create an API Project

```bash
# Create project directory
mkdir my-first-api
cd my-first-api

# Generate base project
vibercode generate api my-first-api
```

### 2. Generate a User Schema

```bash
# Generate User model with PostgreSQL
vibercode schema generate User -m my-first-api -d postgres

# The command will ask for confirmation
# Answer 'y' to continue
```

### 3. Generated Project Structure

```
my-first-api/
├── cmd/
│   └── server/
│       └── main.go          # Entry point
├── internal/
│   ├── handlers/            # HTTP controllers
│   ├── services/            # Business logic
│   ├── repositories/        # Data access
│   └── models/              # Data models
├── pkg/
│   ├── database/            # DB connection
│   └── config/              # Configuration
├── go.mod                   # Go dependencies
└── docker-compose.yml       # Docker config
```

### 4. Run the Project

```bash
# Install dependencies
go mod tidy

# Setup database (PostgreSQL)
docker-compose up -d

# Run the server
go run cmd/server/main.go
```

## 🔧 Basic Commands

### Generate API Project

```bash
vibercode generate api [project-name]
```

### Generate Schemas

```bash
# PostgreSQL
vibercode schema generate [SchemaName] -d postgres -m [module]

# MySQL  
vibercode schema generate [SchemaName] -d mysql -m [module]

# SQLite
vibercode schema generate [SchemaName] -d sqlite -m [module]

# MongoDB
vibercode schema generate [SchemaName] -d mongodb -m [module]
```

### MCP Server (AI Integration)

```bash
# Start MCP server
vibercode mcp

# In another terminal, test integration
./test-mcp-server.sh
```

## 🎨 Customization

### Project Configuration

Edit the `.vibercode-config.json` file:

```json
{
  "project_name": "my-api",
  "database_type": "postgres",
  "port": 8080,
  "enable_auth": true,
  "enable_swagger": true
}
```

### Environment Variables

Copy and configure the environment file:

```bash
cp .env.example .env
```

Edit `.env`:

```env
DATABASE_URL=postgres://user:pass@localhost:5432/mydb
JWT_SECRET=your-secret-key-here
PORT=8080
```

## 🚀 Next Steps

1. **[CLI Commands](cli-commands.md)** - Learn all available commands
2. **[Schema Generation](schema-generation.md)** - Master model creation
3. **[Configuration](configuration.md)** - Customize your project
4. **[Complete Tutorial](../tutorials/first-api.md)** - Build a complete API

## ❓ Issues?

- Check [**Troubleshooting**](../troubleshooting/common-errors.md)
- See [**FAQ**](../troubleshooting/faq.md)
- Report issues on [**GitHub**](https://github.com/vibercode/cli/issues)

---

Congratulations! You now have ViberCode CLI up and running. 🎉