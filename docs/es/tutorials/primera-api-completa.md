# Tu Primera API Completa - Tutorial Paso a Paso

Este tutorial te guiará para crear una API completa de gestión de tareas (ToDo App) usando ViberCode CLI.

## 🎯 Lo que Construiremos

Una API REST completa con:
- ✅ **Gestión de usuarios** con autenticación
- ✅ **Gestión de tareas** con CRUD completo
- ✅ **Categorías de tareas** con relaciones
- ✅ **Base de datos PostgreSQL** con Docker
- ✅ **Arquitectura limpia** Go con Gin y GORM

## 🚦 Prerequisitos

Antes de comenzar, asegúrate de tener instalado:

```bash
# Verificar Go
go version  # Debe mostrar 1.19+

# Verificar Docker
docker --version

# Verificar ViberCode CLI
vibercode --help
```

## 📁 Estructura del Proyecto Final

```
todo-api/
├── cmd/
│   └── server/
│       └── main.go              # Punto de entrada
├── internal/
│   ├── handlers/                # Controladores HTTP
│   │   ├── user_handler.go
│   │   ├── task_handler.go
│   │   └── category_handler.go
│   ├── services/                # Lógica de negocio
│   │   ├── user_service.go
│   │   ├── task_service.go
│   │   └── category_service.go
│   ├── repositories/            # Acceso a datos
│   │   ├── user_repository.go
│   │   ├── task_repository.go
│   │   └── category_repository.go
│   └── models/                  # Modelos de datos
│       ├── user.go
│       ├── task.go
│       └── category.go
├── pkg/
│   ├── database/                # Conexión DB
│   │   └── postgres.go
│   └── config/                  # Configuración
│       └── config.go
├── docker-compose.yml
├── .env
├── go.mod
└── README.md
```

## 🚀 Paso 1: Crear el Proyecto Base

### 1.1 Crear Directorio del Proyecto

```bash
# Crear y entrar al directorio del proyecto
mkdir todo-api
cd todo-api

# Inicializar módulo Go
go mod init todo-api
```

### 1.2 Verificar Estructura

```bash
ls -la
# Debería mostrar:
# - go.mod
```

## 👤 Paso 2: Generar Esquema de Usuarios

### 2.1 Generar Modelo User

```bash
vibercode schema generate User -m todo-api -d postgres -o .
```

**Salida esperada:**
```
🌐 Schema: User
📦 Module: todo-api
🗄️ Database: postgres
⚙️ Output: /path/to/todo-api
Generate code with these settings? [y/N]: y
```

**Responder:** `y`

### 2.2 Verificar Archivos Generados

```bash
find . -name "*.go" -type f
```

**Archivos esperados:**
- `internal/models/user.go`
- `internal/handlers/user_handler.go`
- `internal/services/user_service.go`
- `internal/repositories/user_repository.go`

### 2.3 Examinar el Modelo User

```bash
cat internal/models/user.go
```

**Contenido esperado:**
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

## 📋 Paso 3: Generar Esquema de Tareas

### 3.1 Generar Modelo Task

```bash
vibercode schema generate Task -m todo-api -d postgres -o .
```

**Responder:** `y`

### 3.2 Verificar Generación

```bash
ls internal/models/
```

**Debería mostrar:**
- `user.go`
- `task.go`

## 🏷️ Paso 4: Generar Esquema de Categorías

### 4.1 Generar Modelo Category

```bash
vibercode schema generate Category -m todo-api -d postgres -o .
```

**Responder:** `y`

### 4.2 Estructura Completa Generada

```bash
tree internal/ || find internal/ -type f
```

**Estructura esperada:**
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

## 🐳 Paso 5: Configurar Base de Datos

### 5.1 Crear docker-compose.yml

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

### 5.2 Crear archivo .env

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
JWT_SECRET=tu-clave-secreta-super-segura-aqui-2024

# Redis
REDIS_URL=redis://localhost:6379
EOF
```

### 5.3 Iniciar Base de Datos

```bash
# Iniciar PostgreSQL y Redis
docker-compose up -d

# Verificar que estén ejecutándose
docker-compose ps
```

**Salida esperada:**
```
NAME          SERVICE   STATUS    PORTS
todo-api-postgres-1  postgres  running  0.0.0.0:5432->5432/tcp
todo-api-redis-1     redis     running  0.0.0.0:6379->6379/tcp
```

## 🔧 Paso 6: Instalar Dependencias

### 6.1 Instalar Dependencias Principales

```bash
# Framework web
go get github.com/gin-gonic/gin

# ORM para base de datos
go get gorm.io/gorm
go get gorm.io/driver/postgres

# Utilidades
go get github.com/joho/godotenv
go get github.com/google/uuid
go get golang.org/x/crypto/bcrypt

# JWT para autenticación
go get github.com/golang-jwt/jwt/v4

# Validación
go get github.com/go-playground/validator/v10
```

### 6.2 Limpiar Dependencias

```bash
go mod tidy
```

### 6.3 Verificar go.mod

```bash
cat go.mod
```

## 🏗️ Paso 7: Crear Punto de Entrada

### 7.1 Crear main.go

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
    // Cargar variables de entorno
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    // Conectar a la base de datos
    db, err := database.Connect()
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Configurar Gin
    if os.Getenv("GIN_MODE") == "release" {
        gin.SetMode(gin.ReleaseMode)
    }

    // Crear router
    r := gin.Default()

    // Middleware
    r.Use(gin.Logger())
    r.Use(gin.Recovery())

    // Rutas de salud
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    // Grupos de rutas API
    api := r.Group("/api/v1")
    {
        // Rutas de usuarios
        userHandler := handlers.NewUserHandler(db)
        users := api.Group("/users")
        {
            users.GET("", userHandler.GetUsers)
            users.POST("", userHandler.CreateUser)
            users.GET("/:id", userHandler.GetUser)
            users.PUT("/:id", userHandler.UpdateUser)
            users.DELETE("/:id", userHandler.DeleteUser)
        }

        // Rutas de tareas
        taskHandler := handlers.NewTaskHandler(db)
        tasks := api.Group("/tasks")
        {
            tasks.GET("", taskHandler.GetTasks)
            tasks.POST("", taskHandler.CreateTask)
            tasks.GET("/:id", taskHandler.GetTask)
            tasks.PUT("/:id", taskHandler.UpdateTask)
            tasks.DELETE("/:id", taskHandler.DeleteTask)
        }

        // Rutas de categorías
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

    // Obtener puerto
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    // Iniciar servidor
    fmt.Printf("🚀 Servidor iniciado en puerto %s\n", port)
    fmt.Printf("📚 API docs: http://localhost:%s/api/v1\n", port)
    fmt.Printf("❤️  Health check: http://localhost:%s/health\n", port)
    
    log.Fatal(r.Run(":" + port))
}
EOF
```

### 7.2 Crear Conexión de Base de Datos

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

    // Auto migrar esquemas
    err = db.AutoMigrate(
        &models.User{},
        &models.Task{},
        &models.Category{},
    )
    if err != nil {
        return nil, fmt.Errorf("failed to migrate database: %w", err)
    }

    fmt.Println("✅ Base de datos conectada y migrada exitosamente")
    return db, nil
}
EOF
```

## 🚀 Paso 8: Ejecutar la Aplicación

### 8.1 Compilar y Ejecutar

```bash
# Compilar
go build -o todo-api cmd/server/main.go

# Ejecutar
./todo-api
```

**Salida esperada:**
```
✅ Base de datos conectada y migrada exitosamente
🚀 Servidor iniciado en puerto 8080
📚 API docs: http://localhost:8080/api/v1
❤️  Health check: http://localhost:8080/health
```

### 8.2 Probar la API

En otra terminal:

```bash
# Health check
curl http://localhost:8080/health

# Listar usuarios
curl http://localhost:8080/api/v1/users

# Crear usuario
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@ejemplo.com",
    "password": "password123",
    "first_name": "Juan",
    "last_name": "Pérez"
  }'
```

## ✅ Paso 9: Verificar Funcionalidad

### 9.1 Endpoints Disponibles

**Usuarios:**
- `GET /api/v1/users` - Listar usuarios
- `POST /api/v1/users` - Crear usuario
- `GET /api/v1/users/:id` - Obtener usuario
- `PUT /api/v1/users/:id` - Actualizar usuario
- `DELETE /api/v1/users/:id` - Eliminar usuario

**Tareas:**
- `GET /api/v1/tasks` - Listar tareas
- `POST /api/v1/tasks` - Crear tarea
- `GET /api/v1/tasks/:id` - Obtener tarea
- `PUT /api/v1/tasks/:id` - Actualizar tarea
- `DELETE /api/v1/tasks/:id` - Eliminar tarea

**Categorías:**
- `GET /api/v1/categories` - Listar categorías
- `POST /api/v1/categories` - Crear categoría
- `GET /api/v1/categories/:id` - Obtener categoría
- `PUT /api/v1/categories/:id` - Actualizar categoría
- `DELETE /api/v1/categories/:id` - Eliminar categoría

### 9.2 Verificar Base de Datos

```bash
# Conectar a PostgreSQL
docker exec -it todo-api-postgres-1 psql -U postgres -d todo_api

# Listar tablas
\dt

# Ver estructura de tabla users
\d users

# Salir
\q
```

## 🎉 ¡Felicidades!

Has creado exitosamente tu primera API completa con ViberCode CLI:

✅ **3 modelos** (User, Task, Category) con CRUD completo  
✅ **Arquitectura limpia** con handlers, services y repositories  
✅ **Base de datos PostgreSQL** con Docker  
✅ **API REST** completamente funcional  
✅ **Migraciones automáticas** con GORM  

## 🚀 Próximos Pasos

1. **[Agregar Autenticación JWT](auth-tutorial.md)** - Seguridad para tu API
2. **[Integración con Frontend](frontend-integration.md)** - Conectar con React/Vue
3. **[Deploy en Producción](deployment-guide.md)** - Subir a la nube
4. **[Testing](testing-guide.md)** - Pruebas unitarias e integración

## 🆘 ¿Problemas?

- Consulta [**Errores Comunes**](../troubleshooting/common-errors.md)
- Revisa los [**logs de la aplicación**](../troubleshooting/debugging.md)
- Reporta issues en [**GitHub**](https://github.com/vibercode/cli/issues)

---

*Tutorial completado con ViberCode CLI 🚀*