# CLI Commands - Complete Reference

ViberCode CLI offers several commands for generating and managing Go API projects.

## 📋 Table of Contents

- [Main Commands](#main-commands)
- [Schema Generation](#schema-generation)
- [MCP Server](#mcp-server)
- [Global Options](#global-options)
- [Detailed Examples](#detailed-examples)

## 🎯 Main Commands

### `vibercode help`

Shows general CLI help.

```bash
vibercode help
vibercode --help
vibercode -h
```

### `vibercode version`

Shows the current CLI version.

```bash
vibercode version
```

## 🗄️ Schema Generation

### `vibercode schema generate`

Generates a complete schema with model, repository, service, and handler.

#### Syntax

```bash
vibercode schema generate <schema-name> [flags]
```

#### Available Flags

| Flag | Description | Values | Default |
|------|-------------|---------|---------|
| `-d, --database` | Database provider | `postgres`, `mysql`, `sqlite`, `mongodb` | `postgres` |
| `-m, --module` | Go module name | string | required |
| `-o, --output` | Output directory | path | `.` (current) |
| `-h, --help` | Command help | - | - |

#### Examples

```bash
# Basic schema with PostgreSQL
vibercode schema generate User -m my-api -d postgres

# Schema with MySQL in specific directory
vibercode schema generate Product -m ecommerce -d mysql -o ./src

# Schema with MongoDB
vibercode schema generate Order -m store -d mongodb

# Schema with SQLite
vibercode schema generate Category -m blog -d sqlite
```

## 🔌 MCP Server

### `vibercode mcp`

Starts the MCP (Model Context Protocol) server for AI agent integration.

```bash
vibercode mcp
```

#### Environment Variables

```bash
export ANTHROPIC_API_KEY="your-api-key"
export VIBE_DEBUG="true"
```

#### MCP Server Features

- **Live component editing**
- **Integrated AI chat**
- **Automatic code generation**
- **Project state management**

## ⚙️ Global Options

### Common Flags

| Flag | Description | Example |
|------|-------------|---------|
| `--verbose` | Detailed output | `vibercode schema generate User --verbose` |
| `--config` | Configuration file | `vibercode --config ./custom.json` |
| `--dry-run` | Simulate without executing | `vibercode schema generate --dry-run` |

## 📊 Detailed Examples

### Example 1: Blog API

```bash
# 1. Create project directory
mkdir blog-api && cd blog-api

# 2. Generate user schema
vibercode schema generate User -m blog-api -d postgres

# 3. Generate post schema
vibercode schema generate Post -m blog-api -d postgres

# 4. Generate comment schema
vibercode schema generate Comment -m blog-api -d postgres
```

### Example 2: E-commerce with MongoDB

```bash
# Project directory
mkdir ecommerce-api && cd ecommerce-api

# Main schemas
vibercode schema generate Product -m ecommerce -d mongodb
vibercode schema generate Category -m ecommerce -d mongodb
vibercode schema generate Order -m ecommerce -d mongodb
vibercode schema generate Customer -m ecommerce -d mongodb
```

### Example 3: API with SQLite (development)

```bash
# For quick development with SQLite
mkdir dev-api && cd dev-api

vibercode schema generate Task -m dev-api -d sqlite
vibercode schema generate Project -m dev-api -d sqlite
```

## 🔍 Predefined Schemas

ViberCode includes some common predefined schemas:

### User

Auto-generated fields:
- `ID` (UUID/ObjectID)
- `Email` (string, unique)
- `Password` (string, hashed)
- `FirstName` (string)
- `LastName` (string)
- `Avatar` (string, optional)
- `Active` (boolean)
- `CreatedAt` (timestamp)
- `UpdatedAt` (timestamp)

### Product

Auto-generated fields:
- `ID` (UUID/ObjectID)
- `Name` (string)
- `Description` (text)
- `Price` (decimal)
- `Stock` (integer)
- `CategoryID` (relation)
- `Images` (array)
- `Active` (boolean)
- `CreatedAt` (timestamp)
- `UpdatedAt` (timestamp)

## 🚨 Error Handling

### Common Errors

```bash
# Error: Module required
vibercode schema generate User
# Error: flag needs an argument: -m

# Solution: Specify module
vibercode schema generate User -m my-api

# Error: Unsupported database
vibercode schema generate User -m test -d oracle
# Error: unsupported database type: oracle

# Solution: Use supported database
vibercode schema generate User -m test -d postgres
```

### Exit Codes

| Code | Description |
|------|-------------|
| 0 | Success |
| 1 | General error |
| 2 | Arguments error |
| 3 | File/directory error |
| 4 | Database error |

## 📝 Important Notes

1. **Go Modules**: The `-m` flag must match the name in `go.mod`
2. **Directories**: CLI creates necessary structure automatically
3. **Overwriting**: Existing files are overwritten without confirmation
4. **Dependencies**: Required dependencies are added automatically

## 🔗 Related Links

- [**Schema Generation**](schema-generation.md) - Detailed guide
- [**Configuration**](configuration.md) - Configuration options
- [**Troubleshooting**](../troubleshooting/common-errors.md) - Common errors

---

*For more help, run `vibercode help` or visit the complete documentation.*