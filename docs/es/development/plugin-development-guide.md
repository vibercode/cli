# Guía de Desarrollo de Plugins para ViberCode CLI

Esta guía te enseñará cómo crear, desarrollar y distribuir plugins para ViberCode CLI.

## 📋 Tabla de Contenidos

1. [Introducción](#introducción)
2. [Tipos de Plugins](#tipos-de-plugins)
3. [Configuración del Entorno](#configuración-del-entorno)
4. [Crear tu Primer Plugin](#crear-tu-primer-plugin)
5. [Desarrollo Avanzado](#desarrollo-avanzado)
6. [Testing y Validación](#testing-y-validación)
7. [Distribución](#distribución)
8. [Ejemplos Prácticos](#ejemplos-prácticos)
9. [Mejores Prácticas](#mejores-prácticas)
10. [Solución de Problemas](#solución-de-problemas)

## 🎯 Introducción

Los plugins de ViberCode CLI te permiten extender la funcionalidad del CLI con generadores personalizados, plantillas, comandos y integraciones con servicios externos.

### ¿Qué puedes hacer con un plugin?

- **Generadores**: Crear código Go automáticamente con tus patrones específicos
- **Plantillas**: Proporcionar plantillas reutilizables para diferentes arquitecturas
- **Comandos**: Añadir nuevos comandos al CLI
- **Integraciones**: Conectar con servicios externos (APIs, bases de datos, etc.)

## 🔧 Tipos de Plugins

### 1. Plugin Generador
Genera código automáticamente basado en plantillas y configuraciones.

**Casos de uso:**
- Generadores de microservicios
- Creadores de APIs específicas (GraphQL, gRPC)
- Generadores de código para frameworks específicos

### 2. Plugin de Plantillas
Proporciona plantillas reutilizables para diferentes propósitos.

**Casos de uso:**
- Plantillas para industrias específicas (fintech, healthcare)
- Plantillas de empresa con estándares específicos
- Plantillas para diferentes patrones arquitectónicos

### 3. Plugin de Comandos
Añade nuevos comandos al CLI de ViberCode.

**Casos de uso:**
- Comandos de utilidades personalizadas
- Integraciones con herramientas de desarrollo
- Comandos de automatización de flujos de trabajo

### 4. Plugin de Integración
Conecta ViberCode con servicios externos.

**Casos de uso:**
- Integraciones con bases de datos (Supabase, Firebase)
- Conectores con servicios en la nube (AWS, GCP, Azure)
- Integraciones con herramientas de CI/CD

## ⚙️ Configuración del Entorno

### Requisitos Previos

```bash
# Instalar ViberCode CLI
go install github.com/vibercode/cli@latest

# Verificar instalación
vibercode --version

# Instalar herramientas de desarrollo
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/tools/cmd/goimports@latest
```

### Configurar Plugin SDK

El SDK de plugins está incluido en ViberCode CLI. No necesitas instalaciones adicionales.

## 🚀 Crear tu Primer Plugin

### Paso 1: Generar el Esqueleto del Plugin

```bash
# Crear un plugin generador
vibercode generate plugin --name mi-generador --type generator --author tu-nombre --description "Mi primer generador personalizado"

# Esto creará la estructura:
mi-generador-plugin/
├── plugin.yaml           # Manifiesto del plugin
├── cmd/plugin/main.go    # Punto de entrada
├── internal/plugin.go    # Implementación principal
├── templates/            # Plantillas del generador
├── examples/             # Ejemplos de uso
├── docs/                 # Documentación
├── go.mod               # Módulo Go
├── Makefile             # Comandos de construcción
└── README.md            # Documentación del plugin
```

### Paso 2: Explorar la Estructura Generada

#### `plugin.yaml` - Manifiesto del Plugin
```yaml
name: "mi-generador"
version: "1.0.0"
description: "Mi primer generador personalizado"
author: "tu-nombre"
license: "MIT"
homepage: "https://github.com/tu-nombre/mi-generador"

type: "generator"
capabilities:
  - "code-generation"
  - "template-rendering"
  - "file-management"

dependencies:
  vibercode: ">=1.0.0"
  go: ">=1.19"

commands:
  - name: "mi-generador-command"
    description: "Mi generador command"
    usage: "vibercode mi-generador --help"

config:
  properties:
    default_option:
      type: "string"
      default: "default_value"
      description: "Opción de configuración por defecto"
      required: false

main: "./cmd/plugin/main.go"
```

#### `internal/plugin.go` - Implementación Principal
```go
package internal

import (
    "fmt"
    "path/filepath"
    "github.com/vibercode/plugin-sdk/api"
    "github.com/vibercode/plugin-sdk/utils"
)

type MiGeneradorPlugin struct {
    ctx api.PluginContext
}

func NewMiGeneradorPlugin() *MiGeneradorPlugin {
    return &MiGeneradorPlugin{}
}

// Implementa la interfaz Plugin
func (p *MiGeneradorPlugin) Name() string {
    return "mi-generador"
}

func (p *MiGeneradorPlugin) Version() string {
    return "1.0.0"
}

func (p *MiGeneradorPlugin) Initialize(ctx api.PluginContext) error {
    p.ctx = ctx
    p.ctx.Logger().Info("Inicializando plugin %s", p.Name())
    return nil
}

func (p *MiGeneradorPlugin) Execute(args []string) error {
    // Parsear argumentos
    config, err := p.parseArgs(args)
    if err != nil {
        return fmt.Errorf("error al parsear argumentos: %v", err)
    }

    // Generar código
    return p.generateCode(config)
}

// Implementa la interfaz GeneratorPlugin
func (p *MiGeneradorPlugin) Generate(options map[string]interface{}) (*api.ExecutionResult, error) {
    p.ctx.Logger().Info("Generando código con opciones: %+v", options)

    // Extraer opciones
    name, ok := options["name"].(string)
    if !ok || name == "" {
        return nil, fmt.Errorf("el nombre es requerido")
    }

    // Generar archivos
    files, err := p.generateFiles(name, options)
    if err != nil {
        return &api.ExecutionResult{
            Success: false,
            Error:   fmt.Sprintf("Error en la generación: %v", err),
        }, err
    }

    return &api.ExecutionResult{
        Success:      true,
        Message:      "Generación completada exitosamente",
        FilesCreated: files,
        Output:       fmt.Sprintf("Generados %d archivos para %s", len(files), name),
    }, nil
}
```

### Paso 3: Desarrollar la Lógica del Plugin

#### Crear Plantillas

Crea plantillas en el directorio `templates/`:

**`templates/service.tmpl`**
```go
package {{.Package}}

import (
    "context"
    "fmt"
)

// {{.Name}}Service maneja la lógica de negocio para {{.Name}}
type {{.Name}}Service struct {
    // TODO: Añadir dependencias
}

// New{{.Name}}Service crea una nueva instancia del servicio
func New{{.Name}}Service() *{{.Name}}Service {
    return &{{.Name}}Service{}
}

// Create{{.Name}} crea un nuevo {{.Name}}
func (s *{{.Name}}Service) Create{{.Name}}(ctx context.Context, req *Create{{.Name}}Request) (*{{.Name}}, error) {
    // TODO: Implementar lógica de creación
    return nil, fmt.Errorf("no implementado")
}

// Get{{.Name}} obtiene un {{.Name}} por ID
func (s *{{.Name}}Service) Get{{.Name}}(ctx context.Context, id string) (*{{.Name}}, error) {
    // TODO: Implementar lógica de obtención
    return nil, fmt.Errorf("no implementado")
}
```

#### Implementar la Generación de Archivos

```go
func (p *MiGeneradorPlugin) generateFiles(name string, options map[string]interface{}) ([]string, error) {
    var generatedFiles []string

    // Preparar datos para la plantilla
    data := map[string]interface{}{
        "Name":      name,
        "Package":   strings.ToLower(name),
        "CamelName": utils.ToCamel(name),
        "SnakeName": utils.ToSnake(name),
        "Author":    p.ctx.Config()["author"],
        "Timestamp": time.Now().Format("2006-01-02 15:04:05"),
    }

    // Cargar y renderizar plantilla de servicio
    serviceContent, err := p.ctx.TemplateEngine().RenderFile("templates/service.tmpl", data)
    if err != nil {
        return nil, fmt.Errorf("error al renderizar plantilla de servicio: %v", err)
    }

    // Escribir archivo de servicio
    serviceFile := filepath.Join(".", strings.ToLower(name)+"_service.go")
    if err := p.ctx.FileSystem().WriteFile(serviceFile, []byte(serviceContent)); err != nil {
        return nil, fmt.Errorf("error al escribir archivo de servicio: %v", err)
    }

    generatedFiles = append(generatedFiles, serviceFile)

    // Generar archivo de pruebas si se especifica
    if generateTests, ok := options["with_tests"].(bool); ok && generateTests {
        testFile := filepath.Join(".", strings.ToLower(name)+"_service_test.go")
        testContent := p.generateTestFile(data)
        if err := p.ctx.FileSystem().WriteFile(testFile, []byte(testContent)); err != nil {
            return nil, fmt.Errorf("error al escribir archivo de pruebas: %v", err)
        }
        generatedFiles = append(generatedFiles, testFile)
    }

    return generatedFiles, nil
}
```

### Paso 4: Desarrollo y Testing Local

```bash
# Navegar al directorio del plugin
cd mi-generador-plugin

# Instalar dependencias
go mod tidy

# Construir el plugin
make build

# Enlazar para desarrollo (crea un symlink)
make dev-link

# Verificar que el plugin está disponible
vibercode plugins list

# Probar el plugin
vibercode mi-generador --name TestService --with-tests
```

## 🔬 Desarrollo Avanzado

### Configuración Avanzada del Plugin

#### Esquema de Configuración JSON
```go
func (p *MiGeneradorPlugin) GetSchema() (map[string]interface{}, error) {
    schema := map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "name": map[string]interface{}{
                "type":        "string",
                "description": "Nombre del componente a generar",
                "pattern":     "^[a-zA-Z][a-zA-Z0-9_]*$",
            },
            "package": map[string]interface{}{
                "type":        "string",
                "description": "Nombre del paquete Go",
                "default":     "main",
            },
            "with_tests": map[string]interface{}{
                "type":        "boolean",
                "description": "Generar archivos de pruebas",
                "default":     true,
            },
            "output_dir": map[string]interface{}{
                "type":        "string",
                "description": "Directorio de salida",
                "default":     ".",
            },
        },
        "required": []string{"name"},
    }

    return schema, nil
}
```

#### Validación de Opciones
```go
func (p *MiGeneradorPlugin) ValidateOptions(options map[string]interface{}) error {
    // Validar nombre requerido
    name, ok := options["name"].(string)
    if !ok || name == "" {
        return fmt.Errorf("el nombre es requerido")
    }

    // Validar formato del nombre
    if !utils.IsValidIdentifier(name) {
        return fmt.Errorf("formato de nombre inválido: %s", name)
    }

    // Validar directorio de salida
    if outputDir, ok := options["output_dir"].(string); ok {
        if !filepath.IsAbs(outputDir) && !strings.HasPrefix(outputDir, ".") {
            return fmt.Errorf("directorio de salida debe ser absoluto o relativo: %s", outputDir)
        }
    }

    return nil
}
```

### Plantillas Avanzadas con Helpers

#### Registrar Helpers Personalizados
```go
func (p *MiGeneradorPlugin) Initialize(ctx api.PluginContext) error {
    p.ctx = ctx

    // Registrar helpers personalizados
    err := p.ctx.TemplateEngine().RegisterHelper("pluralize", func(word string) string {
        // Lógica simple de pluralización
        if strings.HasSuffix(word, "y") {
            return strings.TrimSuffix(word, "y") + "ies"
        }
        if strings.HasSuffix(word, "s") {
            return word + "es"
        }
        return word + "s"
    })
    if err != nil {
        return fmt.Errorf("error registrando helper pluralize: %v", err)
    }

    return nil
}
```

#### Usar Helpers en Plantillas
```go
// En la plantilla
type {{.Name}}Repository interface {
    Create{{.Name}}(ctx context.Context, {{.SnakeName}} *{{.Name}}) error
    Get{{.Name}}(ctx context.Context, id string) (*{{.Name}}, error)
    List{{pluralize .Name}}(ctx context.Context) ([]*{{.Name}}, error)
    Update{{.Name}}(ctx context.Context, {{.SnakeName}} *{{.Name}}) error
    Delete{{.Name}}(ctx context.Context, id string) error
}
```

## 🧪 Testing y Validación

### Crear Tests para tu Plugin

**`internal/plugin_test.go`**
```go
package internal

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Mock del contexto del plugin
type MockPluginContext struct {
    mock.Mock
}

func (m *MockPluginContext) Config() map[string]interface{} {
    args := m.Called()
    return args.Get(0).(map[string]interface{})
}

func (m *MockPluginContext) Logger() api.Logger {
    args := m.Called()
    return args.Get(0).(api.Logger)
}

// Test básico del plugin
func TestMiGeneradorPlugin_Generate(t *testing.T) {
    // Preparar
    plugin := NewMiGeneradorPlugin()
    mockCtx := &MockPluginContext{}
    
    mockCtx.On("Config").Return(map[string]interface{}{
        "author": "test-author",
    })
    
    plugin.ctx = mockCtx

    // Opciones de prueba
    options := map[string]interface{}{
        "name":       "TestService",
        "package":    "testpkg",
        "with_tests": true,
    }

    // Ejecutar
    result, err := plugin.Generate(options)

    // Verificar
    assert.NoError(t, err)
    assert.True(t, result.Success)
    assert.Contains(t, result.FilesCreated, "testservice_service.go")
    assert.Contains(t, result.FilesCreated, "testservice_service_test.go")
}

func TestMiGeneradorPlugin_ValidateOptions(t *testing.T) {
    plugin := NewMiGeneradorPlugin()

    tests := []struct {
        name    string
        options map[string]interface{}
        wantErr bool
    }{
        {
            name: "opciones válidas",
            options: map[string]interface{}{
                "name": "ValidService",
            },
            wantErr: false,
        },
        {
            name:    "nombre faltante",
            options: map[string]interface{}{},
            wantErr: true,
        },
        {
            name: "nombre inválido",
            options: map[string]interface{}{
                "name": "123InvalidName",
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := plugin.ValidateOptions(tt.options)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Ejecutar Tests

```bash
# Ejecutar todos los tests
make test

# Ejecutar tests con cobertura
make test-coverage

# Ver reporte de cobertura
open coverage.html
```

### Validar el Plugin

```bash
# Validar la estructura y configuración del plugin
vibercode plugins validate .

# Ejecutar linting
make lint

# Formatear código
make format
```

## 📦 Distribución

### Paso 1: Preparar para Distribución

```bash
# Validar el plugin completamente
make validate

# Ejecutar todos los tests
make test

# Linting y formateo
make lint
make format

# Construir versión final
make build
```

### Paso 2: Empaquetar el Plugin

```bash
# Crear paquete distributable
make package

# Esto creará: mi-generador-1.0.0.tar.gz
```

### Paso 3: Publicar en el Registry (Futuro)

```bash
# Publicar en el registry oficial (cuando esté disponible)
make publish

# O subir manualmente a tu repositorio
git tag v1.0.0
git push origin v1.0.0
```

### Distribución Manual

Mientras tanto, puedes distribuir tu plugin de varias formas:

1. **GitHub Releases**
```bash
# Crear release en GitHub con el archivo .tar.gz
gh release create v1.0.0 mi-generador-1.0.0.tar.gz
```

2. **Instalación desde URL**
```bash
# Los usuarios pueden instalar desde URL
vibercode plugins install https://github.com/tu-usuario/mi-generador/releases/download/v1.0.0/mi-generador-1.0.0.tar.gz
```

## 💡 Ejemplos Prácticos

### Ejemplo 1: Plugin Generador de APIs REST

**Estructura del proyecto:**
```
rest-api-generator/
├── plugin.yaml
├── templates/
│   ├── handler.tmpl       # Plantilla para handlers HTTP
│   ├── service.tmpl       # Plantilla para servicios
│   ├── repository.tmpl    # Plantilla para repositorios
│   ├── model.tmpl         # Plantilla para modelos
│   └── main.tmpl          # Plantilla para main.go
└── internal/plugin.go
```

**Comando de uso:**
```bash
vibercode rest-api-generator --name User --fields "name:string,email:string,age:int" --with-crud
```

### Ejemplo 2: Plugin de Integración con Supabase

**Funcionalidades:**
- Configurar conexión con Supabase
- Generar modelos basados en tablas de Supabase
- Crear operaciones CRUD con Supabase Go client

**Comando de uso:**
```bash
vibercode supabase-integration --project-url https://xxx.supabase.co --anon-key xxx --table users
```

### Ejemplo 3: Plugin de Comandos de Utilidades

**Comandos añadidos:**
- `vibercode db-migrate` - Ejecutar migraciones
- `vibercode db-seed` - Poblar base de datos
- `vibercode api-docs` - Generar documentación

## 🎯 Mejores Prácticas

### 1. Estructura del Código
```go
// ✅ Bueno: Funciones pequeñas y enfocadas
func (p *MyPlugin) generateService(name string, data map[string]interface{}) error {
    // Lógica específica para generar servicio
}

func (p *MyPlugin) generateModel(name string, fields []Field) error {
    // Lógica específica para generar modelo
}

// ❌ Malo: Función muy grande que hace todo
func (p *MyPlugin) generateEverything(options map[string]interface{}) error {
    // 200 líneas de código mezclando responsabilidades
}
```

### 2. Manejo de Errores
```go
// ✅ Bueno: Errores específicos y útiles
func (p *MyPlugin) validateName(name string) error {
    if name == "" {
        return fmt.Errorf("el nombre no puede estar vacío")
    }
    if !regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`).MatchString(name) {
        return fmt.Errorf("el nombre '%s' debe comenzar con una letra y contener solo letras, números y guiones bajos", name)
    }
    return nil
}

// ❌ Malo: Errores genéricos
func (p *MyPlugin) validateName(name string) error {
    if name == "" || !isValid(name) {
        return fmt.Errorf("nombre inválido")
    }
    return nil
}
```

### 3. Logging Efectivo
```go
// ✅ Bueno: Logs informativos y estructurados
func (p *MyPlugin) Generate(options map[string]interface{}) (*api.ExecutionResult, error) {
    p.ctx.Logger().Info("Iniciando generación con opciones: %+v", options)
    
    name := options["name"].(string)
    p.ctx.Logger().Info("Generando archivos para: %s", name)
    
    files, err := p.generateFiles(name, options)
    if err != nil {
        p.ctx.Logger().Error("Error generando archivos: %v", err)
        return nil, err
    }
    
    p.ctx.Logger().Info("Generación completada exitosamente. Archivos creados: %d", len(files))
    return &api.ExecutionResult{Success: true, FilesCreated: files}, nil
}
```

### 4. Plantillas Mantenibles
```go
// ✅ Bueno: Plantillas modulares
const ServiceTemplate = `
{{- template "header" . }}

package {{.Package}}

{{- template "imports" . }}

// {{.Name}}Service handles business logic for {{.Name}}
type {{.Name}}Service struct {
    {{- template "service-fields" . }}
}

{{- template "service-methods" . }}
`

// Plantillas parciales reutilizables
const HeaderTemplate = `
// Code generated by {{.PluginName}} v{{.PluginVersion}}. DO NOT EDIT.
// Generated at: {{.Timestamp}}
`
```

### 5. Testing Completo
```go
func TestPlugin_Integration(t *testing.T) {
    // Setup temporal directory
    tmpDir := t.TempDir()
    
    // Crear plugin con contexto mock
    plugin := setupTestPlugin(t, tmpDir)
    
    // Test completo de extremo a extremo
    options := map[string]interface{}{
        "name":      "TestService",
        "output":    tmpDir,
        "with_crud": true,
    }
    
    result, err := plugin.Generate(options)
    require.NoError(t, err)
    assert.True(t, result.Success)
    
    // Verificar que los archivos se crearon correctamente
    for _, file := range result.FilesCreated {
        assert.FileExists(t, file)
        
        // Verificar que el contenido es válido Go
        content, err := ioutil.ReadFile(file)
        require.NoError(t, err)
        
        // Parse como código Go para verificar sintaxis
        _, err = parser.ParseFile(token.NewFileSet(), file, content, parser.ParseComments)
        assert.NoError(t, err, "Generated file %s should be valid Go code", file)
    }
}
```

## 🔧 Solución de Problemas

### Problemas Comunes

#### 1. Plugin no se carga
```bash
# Verificar que el plugin está bien enlazado
vibercode plugins list

# Verificar el manifiesto
vibercode plugins validate ./mi-plugin

# Revisar logs
vibercode --debug plugins list
```

#### 2. Errores de plantillas
```go
// Debug: Imprimir datos antes de renderizar
fmt.Printf("Template data: %+v\n", data)

// Verificar que la plantilla existe
templatePath := "templates/service.tmpl"
if !p.ctx.FileSystem().Exists(templatePath) {
    return fmt.Errorf("plantilla no encontrada: %s", templatePath)
}
```

#### 3. Problemas de permisos
```bash
# Verificar permisos del directorio de plugins
ls -la ~/.vibercode/plugins/

# Reparar permisos si es necesario
chmod -R 755 ~/.vibercode/plugins/
```

### Debug Avanzado

#### Habilitar Logs de Debug
```bash
export VIBERCODE_DEBUG=true
vibercode mi-plugin --name Test
```

#### Usar el Debugger de Go
```go
// En tu plugin, añadir puntos de breakpoint
import "runtime/debug"

func (p *MyPlugin) Generate(options map[string]interface{}) (*api.ExecutionResult, error) {
    debug.PrintStack() // Imprimir stack trace
    
    // Tu lógica aquí
    return result, nil
}
```

## 📚 Recursos Adicionales

### Documentación de Referencia
- [API del Plugin SDK](../api/plugin-sdk.md)
- [Referencia de Plantillas](../templates/template-guide.md)
- [Guía de Seguridad](../security/plugin-security.md)

### Ejemplos en GitHub
- [Plugin Generador REST API](https://github.com/vibercode/plugins/tree/main/rest-api-generator)
- [Plugin de Integración Supabase](https://github.com/vibercode/plugins/tree/main/supabase-integration)
- [Plugin de Comandos de Utilidades](https://github.com/vibercode/plugins/tree/main/dev-utilities)

### Comunidad
- [Discord de ViberCode](https://discord.gg/vibercode)
- [GitHub Discussions](https://github.com/vibercode/cli/discussions)
- [Stack Overflow Tag: vibercode](https://stackoverflow.com/questions/tagged/vibercode)

---

¡Feliz desarrollo de plugins! 🚀

Si tienes preguntas o necesitas ayuda, no dudes en contactar a la comunidad de ViberCode.