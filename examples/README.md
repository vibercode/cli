# ViberCode CLI Examples

This directory contains examples and sample projects to help you get started with ViberCode CLI.

## 📁 Directory Structure

```
examples/
├── schemas/                 # Example schema definitions
│   ├── blog-api.json       # Simple blog API schema
│   ├── ecommerce.json      # E-commerce API schema
│   └── user-management.json # User management schema
├── generated-projects/      # Complete generated project examples
│   ├── simple-blog/        # Generated blog API
│   ├── task-manager/       # Generated task management API
│   └── inventory-system/   # Generated inventory API
├── scripts/                # Helper scripts and automation
│   ├── quick-start.sh      # Quick project setup script
│   └── demo-setup.sh       # Demo environment setup
└── vector_graph_example.go # Vector graph integration example
```

## 🚀 Quick Start Examples

### 1. Simple Blog API

```bash
# Generate a blog API with posts and comments
vibercode generate api --schema examples/schemas/blog-api.json
```

### 2. E-commerce API

```bash
# Generate a complete e-commerce API
vibercode generate api --schema examples/schemas/ecommerce.json
```

### 3. User Management System

```bash
# Generate user management with authentication
vibercode generate api --schema examples/schemas/user-management.json
```

## 📋 Schema Examples

### Blog API Schema (`blog-api.json`)

Simple blog with posts, comments, and categories:

- **Resources**: Posts, Comments, Categories, Authors
- **Features**: CRUD operations, relationships, validation
- **Database**: PostgreSQL with GORM

### E-commerce Schema (`ecommerce.json`)

Complete e-commerce platform:

- **Resources**: Products, Orders, Users, Inventory
- **Features**: Authentication, payment processing, inventory management
- **Database**: PostgreSQL with advanced relationships

### User Management Schema (`user-management.json`)

User system with roles and permissions:

- **Resources**: Users, Roles, Permissions, Sessions
- **Features**: JWT authentication, role-based access, password reset
- **Database**: PostgreSQL with security features

## 🎨 Generated Project Examples

Each generated project in `generated-projects/` includes:

- ✅ Complete Go API with clean architecture
- ✅ Docker and docker-compose configuration
- ✅ Database migrations
- ✅ API documentation (Swagger)
- ✅ Unit tests
- ✅ README with setup instructions

## 🔧 Helper Scripts

### Quick Start Script (`scripts/quick-start.sh`)

```bash
# Run the quick start script
./examples/scripts/quick-start.sh

# This will:
# 1. Ask for your project preferences
# 2. Generate the API
# 3. Set up the database
# 4. Start the development server
```

### Demo Setup Script (`scripts/demo-setup.sh`)

```bash
# Set up demo environment
./examples/scripts/demo-setup.sh

# This will:
# 1. Generate multiple example APIs
# 2. Set up demo databases
# 3. Populate with sample data
# 4. Start all services
```

## 🧪 Testing Examples

Each example includes comprehensive tests:

```bash
# Run tests for a generated project
cd examples/generated-projects/simple-blog
go test ./...

# Run integration tests
make test-integration

# Run with coverage
make test-coverage
```

## 📚 Learning Path

### Beginner

1. Start with **Simple Blog API**
2. Explore the generated code structure
3. Run the API and test endpoints
4. Modify the schema and regenerate

### Intermediate

1. Try **User Management System**
2. Understand authentication flow
3. Customize middleware and validation
4. Add custom business logic

### Advanced

1. Build **E-commerce API**
2. Implement complex relationships
3. Add custom generators
4. Integrate with external services

## 🔗 Related Documentation

- [CLI Commands](../docs/en/user-guide/cli-commands.md)
- [Schema Definition Guide](../docs/en/api/schema-format.md)
- [Generated Code Structure](../docs/en/development/code-structure.md)
- [Customization Guide](../docs/en/development/customization.md)

## 💡 Tips

- **Start Simple**: Begin with the blog example to understand basics
- **Read Generated Code**: Explore generated files to understand patterns
- **Modify & Regenerate**: Change schemas and see how code updates
- **Test Everything**: Use provided tests as learning material
- **Join Community**: Share your examples and get help from others

---

**Happy coding with ViberCode CLI! 🚀**
