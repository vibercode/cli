# 🧪 **Guía Completa de Pruebas - Enhanced Field Types**

## **Opción 1: Pruebas CLI Directas**

### **1.1 Generar API Completo**
```bash
# Ejecutar el CLI
./vibercode generate api

# Inputs de prueba:
# Project name: enhanced-test-api
# Port: 8080
# Database: postgres (o supabase para probar nuevas features)
# Module: github.com/test/enhanced-api
```

### **1.2 Generar Resource con Nuevos Field Types**
```bash
# Ejecutar resource generator
./vibercode generate resource

# Inputs de prueba para cada field type:
# Resource name: UserProfile
# Description: Testing enhanced field types

# Field 1: Email
# Name: email
# Type: email ← NUEVO!
# Required: yes

# Field 2: Password  
# Name: password
# Type: password ← NUEVO!
# Required: yes

# Field 3: Website
# Name: website
# Type: url ← NUEVO!
# Required: no

# Field 4: Phone
# Name: phone
# Type: phone ← NUEVO!
# Required: no

# Field 5: Color
# Name: theme_color
# Type: color ← NUEVO!
# Required: no

# Continuar con más fields...
```

### **1.3 Verificar Código Generado**
```bash
# Revisar el código generado
cd enhanced-test-api
find . -name "*.go" | head -10

# Verificar model con validaciones
cat internal/models/user_profile.go

# Verificar handlers con nuevos tipos
cat internal/handlers/user_profile_handler.go

# Verificar validation logic
grep -n "validation\|email\|url\|color" internal/models/user_profile.go
```

---

## **Opción 2: Pruebas Server Mode + Editor**

### **2.1 Iniciar CLI Server**
```bash
# Terminal 1: Iniciar server CLI
./vibercode serve

# Output esperado:
# 🌐 Server starting on port 8080
# 📡 API endpoints available at http://localhost:8080/api/v1
# 🔌 WebSocket server running on port 8081
```

### **2.2 Abrir React Editor**
```bash
# Terminal 2: Navegar al editor
cd /Users/jaambee/Projects/vibercode/editor

# Instalar dependencias si es necesario
npm install

# Iniciar editor
npm run dev

# Output esperado:
# ➜  Local:   http://localhost:5173/
```

### **2.3 Probar Integración en Browser**

1. **Abrir**: http://localhost:5173
2. **Verificar conexión**: Debe mostrar "Connected to Vibercode CLI"
3. **Probar nuevos componentes**:
   - Drag & drop "Email Input"
   - Drag & drop "Color Picker" 
   - Drag & drop "Currency Input"
   - Etc.

4. **Exportar Schema**:
   - Click "Export Schema"
   - Verificar JSON contiene nuevos field types
   - Guardar como `test-export.json`

5. **Importar Schema**:
   - Click "Import Schema"
   - Cargar `/Users/jaambee/Projects/vibercode/editor/example-enhanced-schema.json`
   - Verificar campos aparecen en canvas

---

## **Opción 3: Pruebas API Directas**

### **3.1 Test API Endpoints**
```bash
# Con CLI server corriendo, probar APIs:

# Health check
curl http://localhost:8080/api/v1/health

# Create schema con nuevos field types
curl -X POST http://localhost:8080/api/v1/schema/create \
  -H "Content-Type: application/json" \
  -d @test-enhanced-resource.json

# List schemas
curl http://localhost:8080/api/v1/schema/list

# Generate code
curl -X POST http://localhost:8080/api/v1/generate/resource \
  -H "Content-Type: application/json" \
  -d '{
    "schema_id": "user-profile-123",
    "output_path": "./generated-test",
    "database_provider": "postgres"
  }'
```

---

## **Opción 4: Pruebas Unit Tests**

### **4.1 Ejecutar Test Suite Completo**
```bash
# Todos los tests
go test ./...

# Solo tests de fields
go test ./internal/models/ -v

# Solo tests específicos
go test ./internal/models/ -run TestField_GoValidation -v

# Tests con coverage
go test ./internal/models/ -cover
```

### **4.2 Test Validación Manual**
```bash
# Crear test file temporal
cat > test_field_validation.go << 'EOF'
package main

import (
    "fmt"
    "github.com/vibercode/cli/internal/models"
)

func main() {
    // Test email field
    emailField := models.Field{
        Type: models.FieldTypeEmail,
        Name: "email",
        DisplayName: "Email",
        Required: true,
    }
    
    fmt.Println("Email GoType:", emailField.GoType())
    fmt.Println("Email Validation:", emailField.GoValidation())
    fmt.Println("Email Struct:", emailField.GoStructField())
    
    // Test enum field
    enumField := models.Field{
        Type: models.FieldTypeEnum,
        Name: "status", 
        DisplayName: "Status",
        EnumValues: []string{"active", "inactive"},
    }
    
    fmt.Println("\nEnum Type Generation:")
    fmt.Println(enumField.GenerateEnumType())
}
EOF

# Ejecutar test
go run test_field_validation.go

# Limpiar
rm test_field_validation.go
```

---

## **🎯 Casos de Prueba Específicos**

### **Caso 1: E-commerce Product**
```json
{
  "name": "Product",
  "fields": [
    {"name": "name", "type": "string", "required": true},
    {"name": "slug", "type": "slug", "required": true, "unique": true},
    {"name": "price", "type": "currency", "required": true, "min_value": 0},
    {"name": "color", "type": "color", "required": false},
    {"name": "website", "type": "url", "required": false},
    {"name": "status", "type": "enum", "enum_values": ["draft", "published", "archived"]}
  ]
}
```

### **Caso 2: User Registration**
```json
{
  "name": "User",
  "fields": [
    {"name": "email", "type": "email", "required": true, "unique": true},
    {"name": "password", "type": "password", "required": true, "min_length": 8},
    {"name": "phone", "type": "phone", "required": false},
    {"name": "website", "type": "url", "required": false},
    {"name": "location", "type": "coordinates", "required": false}
  ]
}
```

---

## **✅ Checklist de Validación**

### **CLI Features:**
- [ ] Nuevos field types aparecen en prompts
- [ ] Validación se genera correctamente
- [ ] Templates producen código válido
- [ ] Tests pasan al 100%
- [ ] Multiple database providers funcionan

### **Editor Features:**
- [ ] Nuevos componentes en librería
- [ ] Drag & drop funciona
- [ ] Propiedades se configuran correctamente
- [ ] Export/import preserve field types
- [ ] API integration funciona

### **Integration Features:**
- [ ] Schema round-trip CLI → Editor → CLI
- [ ] Code generation produce código compilable
- [ ] Validación funciona en runtime
- [ ] Database migrations son correctas
- [ ] API endpoints responden correctamente

---

## **🐛 Troubleshooting**

### **CLI No Compila:**
```bash
go mod tidy
go build -o vibercode main.go
```

### **Editor No Conecta:**
```bash
# Verificar CLI server corriendo
curl http://localhost:8080/api/v1/health

# Verificar puertos
lsof -i :8080
lsof -i :5173
```

### **Tests Fallan:**
```bash
# Verificar dependencies
go mod verify

# Run tests con verbose
go test ./internal/models/ -v -count=1
```

---

**¡Con esta guía puedes probar todas las funcionalidades nuevas paso a paso!** 🚀