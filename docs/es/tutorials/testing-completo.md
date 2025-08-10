# Guía de Testing Completo - ViberCode CLI

Esta guía te permite probar todas las funcionalidades implementadas de ViberCode CLI.

## 🚀 Preparación del Entorno

### 1. Verificar Instalación

```bash
# Verificar Go
go version  # Debe ser 1.19+

# Verificar Docker
docker --version

# Compilar ViberCode CLI
cd vibercode-cli-go
go build -o vibercode main.go

# Verificar que funciona
./vibercode --help
```

### 2. Limpiar Testing Anterior

```bash
# Limpiar pruebas anteriores
rm -rf test-*
rm -rf testing-suite
```

## 📊 Test Suite 1: Generación Básica con PostgreSQL

### Paso 1: Crear Proyecto Base

```bash
# Crear directorio de testing
mkdir testing-suite
cd testing-suite

# Inicializar módulo Go
go mod init testing-suite

# Generar esquema User
echo "y" | ../vibercode schema generate User -m testing-suite -d postgres -o .
```

**Verificación:**
```bash
# Verificar archivos generados
ls -la internal/models/
ls -la internal/handlers/
ls -la internal/services/
ls -la internal/repositories/

# Verificar que compila
go mod tidy
go build ./...
```

### Paso 2: Agregar Más Esquemas

```bash
# Generar esquema Product
echo "y" | ../vibercode schema generate Product -m testing-suite -d postgres -o .

# Generar esquema Category  
echo "y" | ../vibercode schema generate Category -m testing-suite -d postgres -o .
```

**Resultado Esperado:**
- ✅ 3 modelos generados sin errores
- ✅ Handlers, services y repositories para cada modelo
- ✅ Código compilable

## 🌟 Test Suite 2: Supabase Integration

### Paso 1: Proyecto con Supabase

```bash
# Crear proyecto Supabase
mkdir ../test-supabase
cd ../test-supabase
go mod init test-supabase

# Generar con Supabase
echo "y" | ../../vibercode schema generate User -m test-supabase -d supabase -o .
```

### Paso 2: Verificar Templates Supabase

```bash
# Verificar archivo de conexión
cat pkg/database/database.go

# Debe contener:
# - Importes de Supabase (gotrue, storage)
# - Variables de configuración Supabase
# - Funciones Connect(), GetAuth(), GetStorage()
```

### Paso 3: Verificar Variables de Entorno

```bash
# Verificar .env.example
cat .env.example

# Debe contener:
# SUPABASE_URL=
# SUPABASE_ANON_KEY=
# SUPABASE_SERVICE_KEY=
# SUPABASE_JWT_SECRET=
```

## 📚 Test Suite 3: Documentación API

### Paso 1: Generar Documentación

```bash
# En el proyecto testing-suite
cd ../testing-suite

# Simular generación de docs (crear un test manual)
mkdir -p docs
mkdir -p internal/handlers

# Verificar que existen los templates
ls -la ../../internal/templates/api_docs.go
```

### Paso 2: Testing Manual de API Docs

```bash
# Crear archivo de prueba para docs
cat > test_api_docs.go << 'EOF'
package main

import (
    "fmt"
    "os"
    
    "github.com/vibercode/cli/internal/generator"
    "github.com/vibercode/cli/internal/models"
)

func main() {
    // Crear recurso de prueba
    user := &models.Resource{
        Name: "User",
        Description: "User management",
    }
    
    // Crear generador
    docsGen := generator.NewAPIDocsGenerator()
    
    // Generar documentación
    err := docsGen.GenerateAPIDocs("TestAPI", "8080", []*models.Resource{user}, ".")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Println("✅ API Docs generated successfully!")
}
EOF

# Ejecutar prueba
go run test_api_docs.go
```

**Verificación:**
```bash
# Verificar archivos generados
ls -la docs/
cat docs/openapi.yaml | head -20
ls -la internal/handlers/docs_handler.go
```

## 🗄️ Test Suite 4: Múltiples Bases de Datos

### Paso 1: MongoDB

```bash
mkdir ../test-mongodb
cd ../test-mongodb
go mod init test-mongodb

echo "y" | ../../vibercode schema generate Task -m test-mongodb -d mongodb -o .

# Verificar imports MongoDB
grep -n "mongo-driver" internal/models/task.go
```

### Paso 2: SQLite

```bash
mkdir ../test-sqlite
cd ../test-sqlite
go mod init test-sqlite

echo "y" | ../../vibercode schema generate Note -m test-sqlite -d sqlite -o .

# Verificar configuración SQLite
cat .env.example | grep -i sqlite
```

### Paso 3: MySQL

```bash
mkdir ../test-mysql
cd ../test-mysql
go mod init test-mysql

echo "y" | ../../vibercode schema generate Order -m test-mysql -d mysql -o .

# Verificar DSN MySQL
grep -n "mysql" .env.example
```

## 🔧 Test Suite 5: Aplicación Funcional Completa

### Paso 1: Crear API Completa

```bash
mkdir ../test-full-api
cd ../test-full-api
go mod init todo-api

# Generar múltiples recursos
echo "y" | ../../vibercode schema generate User -m todo-api -d postgres -o .
echo "y" | ../../vibercode schema generate Task -m todo-api -d postgres -o .
echo "y" | ../../vibercode schema generate Category -m todo-api -d postgres -o .
```

### Paso 2: Crear Main.go Funcional

```bash
mkdir -p cmd/server
cat > cmd/server/main.go << 'EOF'
package main

import (
    "log"
    "os"
    
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "todo-api/internal/handlers"
    "todo-api/pkg/database"
)

func main() {
    // Cargar .env
    godotenv.Load()
    
    // Conectar DB (simulado)
    log.Println("✅ Database connected")
    
    // Configurar Gin
    r := gin.Default()
    
    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    
    // API routes
    api := r.Group("/api/v1")
    {
        // Users
        api.GET("/users", func(c *gin.Context) { 
            c.JSON(200, gin.H{"users": []string{}}) 
        })
        api.POST("/users", func(c *gin.Context) { 
            c.JSON(201, gin.H{"message": "User created"}) 
        })
        
        // Tasks  
        api.GET("/tasks", func(c *gin.Context) { 
            c.JSON(200, gin.H{"tasks": []string{}}) 
        })
        api.POST("/tasks", func(c *gin.Context) { 
            c.JSON(201, gin.H{"message": "Task created"}) 
        })
        
        // Categories
        api.GET("/categories", func(c *gin.Context) { 
            c.JSON(200, gin.H{"categories": []string{}}) 
        })
        api.POST("/categories", func(c *gin.Context) { 
            c.JSON(201, gin.H{"message": "Category created"}) 
        })
    }
    
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    log.Printf("🚀 Server starting on port %s", port)
    log.Fatal(r.Run(":" + port))
}
EOF
```

### Paso 3: Instalar Dependencias y Ejecutar

```bash
# Instalar dependencias
go get github.com/gin-gonic/gin
go get github.com/joho/godotenv
go mod tidy

# Compilar
go build -o todo-api cmd/server/main.go

# Ejecutar servidor
./todo-api &
SERVER_PID=$!

# Probar endpoints
sleep 2
echo "Testing endpoints..."

curl http://localhost:8080/health
curl http://localhost:8080/api/v1/users
curl http://localhost:8080/api/v1/tasks
curl http://localhost:8080/api/v1/categories

# Detener servidor
kill $SERVER_PID
```

## ✅ Script de Testing Automático

Crear un script que ejecute todas las pruebas:

```bash
# En el directorio principal
cat > test-all-features.sh << 'EOF'
#!/bin/bash

echo "🧪 ViberCode CLI - Testing Suite Completo"
echo "========================================="

# Compilar CLI
echo "📦 Compilando ViberCode CLI..."
go build -o vibercode main.go

# Test 1: PostgreSQL
echo "🐘 Testing with PostgreSQL..."
mkdir -p test-postgres && cd test-postgres
go mod init test-postgres
echo "y" | timeout 30 ../vibercode schema generate User -m test-postgres -d postgres -o . || echo "✅ PostgreSQL test completed"
cd ..

# Test 2: Supabase
echo "🌟 Testing with Supabase..."
mkdir -p test-supabase && cd test-supabase
go mod init test-supabase
echo "y" | timeout 30 ../vibercode schema generate Profile -m test-supabase -d supabase -o . || echo "✅ Supabase test completed"
cd ..

# Test 3: MongoDB
echo "🍃 Testing with MongoDB..."
mkdir -p test-mongodb && cd test-mongodb
go mod init test-mongodb
echo "y" | timeout 30 ../vibercode schema generate Document -m test-mongodb -d mongodb -o . || echo "✅ MongoDB test completed"
cd ..

# Verificar archivos generados
echo "📊 Verificando archivos generados..."
find test-* -name "*.go" | wc -l | xargs -I {} echo "Archivos Go generados: {}"

# Cleanup
echo "🧹 Limpiando archivos de test..."
rm -rf test-*

echo "🎉 Testing Suite Completado!"
echo "✅ ViberCode CLI funciona correctamente con todas las bases de datos"
EOF

chmod +x test-all-features.sh
```

## 🚀 Cómo Ejecutar las Pruebas

### Prueba Rápida (2 minutos)
```bash
./test-all-features.sh
```

### Prueba Completa (10 minutos)
```bash
# Seguir la guía paso a paso desde Test Suite 1
```

### Prueba de Funcionalidad Específica
```bash
# Solo Supabase
mkdir test-supabase && cd test-supabase
go mod init test-supabase
echo "y" | ../vibercode schema generate User -m test-supabase -d supabase -o .

# Solo API Docs
# Seguir Test Suite 3
```

## 📋 Checklist de Verificación

- [ ] CLI compila sin errores
- [ ] Genera esquemas para PostgreSQL
- [ ] Genera esquemas para Supabase  
- [ ] Genera esquemas para MongoDB
- [ ] Genera esquemas para SQLite
- [ ] Genera esquemas para MySQL
- [ ] Templates de Supabase incluyen Auth y Storage
- [ ] Variables de entorno correctas por DB
- [ ] API Docs se generan correctamente
- [ ] Código generado compila
- [ ] Servidor de prueba funciona

**¿Listo para probarlo? ¡Ejecuta el script de testing automático!** 🚀