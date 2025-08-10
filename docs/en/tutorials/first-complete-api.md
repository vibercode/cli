# Your First Complete API - Step by Step Tutorial

This tutorial will guide you through creating a complete task management API (ToDo App) using ViberCode CLI.

## 🎯 What We'll Build

A complete REST API featuring:
- ✅ **User management** with authentication
- ✅ **Task management** with full CRUD
- ✅ **Task categories** with relationships
- ✅ **PostgreSQL database** with Docker
- ✅ **Clean architecture** Go with Gin and GORM

## 🚦 Prerequisites

Before starting, make sure you have:

```bash
# Check Go
go version  # Should show 1.19+

# Check Docker
docker --version

# Check ViberCode CLI
vibercode --help
```

## 📁 Final Project Structure

```
todo-api/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── handlers/                # HTTP controllers
│   │   ├── user_handler.go
│   │   ├── task_handler.go
│   │   └── category_handler.go
│   ├── services/                # Business logic
│   │   ├── user_service.go
│   │   ├── task_service.go
│   │   └── category_service.go
│   ├── repositories/            # Data access
│   │   ├── user_repository.go
│   │   ├── task_repository.go
│   │   └── category_repository.go
│   └── models/                  # Data models
│       ├── user.go
│       ├── task.go
│       └── category.go
├── pkg/
│   ├── database/                # DB connection
│   │   └── postgres.go
│   └── config/                  # Configuration
│       └── config.go
├── docker-compose.yml
├── .env
├── go.mod
└── README.md
```

## 🚀 Step 1: Create Base Project

### 1.1 Create Project Directory

```bash
# Create and enter project directory
mkdir todo-api
cd todo-api

# Initialize Go module
go mod init todo-api
```

### 1.2 Verify Structure

```bash
ls -la
# Should show:
# - go.mod
```

## 👤 Step 2: Generate User Schema

### 2.1 Generate User Model

```bash
vibercode schema generate User -m todo-api -d postgres -o .
```

**Expected output:**
```
🌐 Schema: User
📦 Module: todo-api
🗄️ Database: postgres
⚙️ Output: /path/to/todo-api
Generate code with these settings? [y/N]: y
```

**Answer:** `y`

### 2.2 Verify Generated Files

```bash
find . -name "*.go" -type f
```

**Expected files:**
- `internal/models/user.go`
- `internal/handlers/user_handler.go`
- `internal/services/user_service.go`
- `internal/repositories/user_repository.go`

### 2.3 Examine User Model

```bash
cat internal/models/user.go
```

**Expected content:**
```go
package models

import (
    "fmt"
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents the User model
type User struct {
    ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Email     string            `json:"email" bson:"email"`
    Password  string            `json:"password" bson:"password"`
    FirstName string            `json:"first_name" bson:"first_name"`
    LastName  string            `json:"last_name" bson:"last_name"`
    Avatar    string            `json:"avatar" bson:"avatar"`
    Active    bool              `json:"active" bson:"active"`
    CreatedAt time.Time         `json:"created_at" bson:"created_at"`
    UpdatedAt time.Time         `json:"updated_at" bson:"updated_at"`
}
```

## 📋 Step 3: Generate Task Schema

### 3.1 Generate Task Model

```bash
vibercode schema generate Task -m todo-api -d postgres -o .
```

**Answer:** `y`

### 3.2 Verify Generation

```bash
ls internal/models/
```

**Should show:**
- `user.go`
- `task.go`

## 🏷️ Step 4: Generate Category Schema

### 4.1 Generate Category Model

```bash
vibercode schema generate Category -m todo-api -d postgres -o .
```

**Answer:** `y`

### 4.2 Complete Generated Structure

```bash
tree internal/ || find internal/ -type f
```

**Expected structure:**
```
internal/
├── handlers/
│   ├── user_handler.go
│   ├── task_handler.go
│   └── category_handler.go
├── services/
│   ├── user_service.go
│   ├── task_service.go
│   └── category_service.go
├── repositories/
│   ├── user_repository.go
│   ├── task_repository.go
│   └── category_repository.go
└── models/
    ├── user.go
    ├── task.go
    └── category.go
```

## 🐳 Step 5: Setup Database

### 5.1 Create docker-compose.yml

```bash
cat > docker-compose.yml << 'EOF'
version: '3.8'
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: todo_api
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  postgres_data:
EOF
```

### 5.2 Create .env file

```bash
cat > .env << 'EOF'
# Database
DATABASE_URL=postgres://postgres:postgres@localhost:5432/todo_api?sslmode=disable
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=todo_api

# Server
PORT=8080
GIN_MODE=release

# JWT
JWT_SECRET=your-super-secure-secret-key-here-2024

# Redis
REDIS_URL=redis://localhost:6379
EOF
```

### 5.3 Start Database

```bash
# Start PostgreSQL and Redis
docker-compose up -d

# Verify they're running
docker-compose ps
```

**Expected output:**
```
NAME          SERVICE   STATUS    PORTS
todo-api-postgres-1  postgres  running  0.0.0.0:5432->5432/tcp
todo-api-redis-1     redis     running  0.0.0.0:6379->6379/tcp
```

## 🔧 Step 6: Install Dependencies

### 6.1 Install Main Dependencies

```bash
# Web framework
go get github.com/gin-gonic/gin

# ORM for database
go get gorm.io/gorm
go get gorm.io/driver/postgres

# Utilities
go get github.com/joho/godotenv
go get github.com/google/uuid
go get golang.org/x/crypto/bcrypt

# JWT for authentication
go get github.com/golang-jwt/jwt/v4

# Validation
go get github.com/go-playground/validator/v10
```

### 6.2 Clean Dependencies

```bash
go mod tidy
```

### 6.3 Verify go.mod

```bash
cat go.mod
```

## 🏗️ Step 7: Create Entry Point

### 7.1 Create main.go

```bash
mkdir -p cmd/server
cat > cmd/server/main.go << 'EOF'
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "todo-api/internal/handlers"
    "todo-api/pkg/database"
)

func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    // Connect to database
    db, err := database.Connect()
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Configure Gin
    if os.Getenv("GIN_MODE") == "release" {
        gin.SetMode(gin.ReleaseMode)
    }

    // Create router
    r := gin.Default()

    // Middleware
    r.Use(gin.Logger())
    r.Use(gin.Recovery())

    // Health routes
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    // API route groups
    api := r.Group("/api/v1")
    {
        // User routes
        userHandler := handlers.NewUserHandler(db)
        users := api.Group("/users")
        {
            users.GET("", userHandler.GetUsers)
            users.POST("", userHandler.CreateUser)
            users.GET("/:id", userHandler.GetUser)
            users.PUT("/:id", userHandler.UpdateUser)
            users.DELETE("/:id", userHandler.DeleteUser)
        }

        // Task routes
        taskHandler := handlers.NewTaskHandler(db)
        tasks := api.Group("/tasks")
        {
            tasks.GET("", taskHandler.GetTasks)
            tasks.POST("", taskHandler.CreateTask)
            tasks.GET("/:id", taskHandler.GetTask)
            tasks.PUT("/:id", taskHandler.UpdateTask)
            tasks.DELETE("/:id", taskHandler.DeleteTask)
        }

        // Category routes
        categoryHandler := handlers.NewCategoryHandler(db)
        categories := api.Group("/categories")
        {
            categories.GET("", categoryHandler.GetCategories)
            categories.POST("", categoryHandler.CreateCategory)
            categories.GET("/:id", categoryHandler.GetCategory)
            categories.PUT("/:id", categoryHandler.UpdateCategory)
            categories.DELETE("/:id", categoryHandler.DeleteCategory)
        }
    }

    // Get port
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    // Start server
    fmt.Printf("🚀 Server started on port %s\n", port)
    fmt.Printf("📚 API docs: http://localhost:%s/api/v1\n", port)
    fmt.Printf("❤️  Health check: http://localhost:%s/health\n", port)
    
    log.Fatal(r.Run(":" + port))
}
EOF
```

### 7.2 Create Database Connection

```bash
mkdir -p pkg/database
cat > pkg/database/postgres.go << 'EOF'
package database

import (
    "fmt"
    "os"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "todo-api/internal/models"
)

func Connect() (*gorm.DB, error) {
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
            os.Getenv("DB_HOST"),
            os.Getenv("DB_PORT"),
            os.Getenv("DB_USER"),
            os.Getenv("DB_PASSWORD"),
            os.Getenv("DB_NAME"),
        )
    }

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    // Auto migrate schemas
    err = db.AutoMigrate(
        &models.User{},
        &models.Task{},
        &models.Category{},
    )
    if err != nil {
        return nil, fmt.Errorf("failed to migrate database: %w", err)
    }

    fmt.Println("✅ Database connected and migrated successfully")
    return db, nil
}
EOF
```

## 🚀 Step 8: Run the Application

### 8.1 Build and Run

```bash
# Build
go build -o todo-api cmd/server/main.go

# Run
./todo-api
```

**Expected output:**
```
✅ Database connected and migrated successfully
🚀 Server started on port 8080
📚 API docs: http://localhost:8080/api/v1
❤️  Health check: http://localhost:8080/health
```

### 8.2 Test the API

In another terminal:

```bash
# Health check
curl http://localhost:8080/health

# List users
curl http://localhost:8080/api/v1/users

# Create user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

## ✅ Step 9: Verify Functionality

### 9.1 Available Endpoints

**Users:**
- `GET /api/v1/users` - List users
- `POST /api/v1/users` - Create user
- `GET /api/v1/users/:id` - Get user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

**Tasks:**
- `GET /api/v1/tasks` - List tasks
- `POST /api/v1/tasks` - Create task
- `GET /api/v1/tasks/:id` - Get task
- `PUT /api/v1/tasks/:id` - Update task
- `DELETE /api/v1/tasks/:id` - Delete task

**Categories:**
- `GET /api/v1/categories` - List categories
- `POST /api/v1/categories` - Create category
- `GET /api/v1/categories/:id` - Get category
- `PUT /api/v1/categories/:id` - Update category
- `DELETE /api/v1/categories/:id` - Delete category

### 9.2 Verify Database

```bash
# Connect to PostgreSQL
docker exec -it todo-api-postgres-1 psql -U postgres -d todo_api

# List tables
\dt

# View users table structure
\d users

# Exit
\q
```

## 🎉 Congratulations!

You've successfully created your first complete API with ViberCode CLI:

✅ **3 models** (User, Task, Category) with full CRUD  
✅ **Clean architecture** with handlers, services and repositories  
✅ **PostgreSQL database** with Docker  
✅ **Fully functional REST API**  
✅ **Automatic migrations** with GORM  

## 🚀 Next Steps

1. **[Add JWT Authentication](auth-tutorial.md)** - Secure your API
2. **[Frontend Integration](frontend-integration.md)** - Connect with React/Vue
3. **[Production Deployment](deployment-guide.md)** - Deploy to the cloud
4. **[Testing](testing-guide.md)** - Unit and integration tests

## 🆘 Issues?

- Check [**Common Errors**](../troubleshooting/common-errors.md)
- Review [**application logs**](../troubleshooting/debugging.md)
- Report issues on [**GitHub**](https://github.com/vibercode/cli/issues)

---

*Tutorial completed with ViberCode CLI 🚀*