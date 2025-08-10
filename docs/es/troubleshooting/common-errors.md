# Errores Comunes - ViberCode CLI

Esta guía te ayudará a resolver los problemas más frecuentes al usar ViberCode CLI.

## 🚨 Errores de Template

### Error: `function "sortField" not defined`

**Problema**: Error en el parsing de templates durante la generación de esquemas.

```bash
Error: failed to generate code: failed to generate repository: 
failed to parse template: template: generator:95: function "sortField" not defined
```

**Solución**: Este error ha sido corregido en versiones recientes. Actualiza tu CLI:

```bash
# Recompilar desde código fuente
go build -o vibercode main.go

# O reinstalar
go install github.com/vibercode/cli@latest
```

## 🗄️ Errores de Base de Datos

### Error: `unsupported database type`

**Problema**: Base de datos no soportada especificada.

```bash
Error: unsupported database type: oracle
```

**Solución**: Usa una base de datos soportada:

```bash
# Bases de datos soportadas
vibercode schema generate User -m test -d postgres   # PostgreSQL
vibercode schema generate User -m test -d mysql     # MySQL
vibercode schema generate User -m test -d sqlite    # SQLite
vibercode schema generate User -m test -d mongodb   # MongoDB
```

### Error: Dependencias MongoDB faltantes

**Problema**: Error al compilar con esquemas MongoDB.

```bash
no required module provides package go.mongodb.org/mongo-driver/bson/primitive
```

**Solución**: Instalar dependencias de MongoDB:

```bash
go get go.mongodb.org/mongo-driver/bson/primitive
go mod tidy
```

## 📦 Errores de Módulos Go

### Error: `flag needs an argument: -m`

**Problema**: Módulo no especificado.

```bash
vibercode schema generate User
Error: flag needs an argument: -m
```

**Solución**: Especificar el nombre del módulo:

```bash
vibercode schema generate User -m mi-proyecto
```

### Error: `go.mod` no encontrado

**Problema**: Directorio sin inicializar como módulo Go.

```bash
Error: go.mod file not found
```

**Solución**: Inicializar módulo Go:

```bash
go mod init mi-proyecto
vibercode schema generate User -m mi-proyecto
```

## 🔧 Errores de Compilación

### Error: Imports no utilizados

**Problema**: Imports automáticos no utilizados en el código generado.

```bash
"encoding/json" imported and not used
undefined: fmt
```

**Solución**: El CLI limpia automáticamente los imports. Si persiste:

```bash
# Limpiar imports manualmente
go mod tidy
gofmt -w .
```

### Error: `go.sum` inconsistente

**Problema**: Checksums incorrectos en go.sum.

```bash
missing go.sum entry for module
```

**Solución**: Limpiar y regenerar go.sum:

```bash
go clean -modcache
go mod download
go mod tidy
```

## 🖥️ Errores del Servidor MCP

### Error: `ANTHROPIC_API_KEY` no configurada

**Problema**: Variable de entorno faltante para MCP.

```bash
Error: ANTHROPIC_API_KEY environment variable not set
```

**Solución**: Configurar la variable de entorno:

```bash
export ANTHROPIC_API_KEY="tu-api-key-aqui"
vibercode mcp
```

### Error: Puerto en uso

**Problema**: Puerto del servidor MCP ocupado.

```bash
Error: bind: address already in use
```

**Solución**: Cambiar puerto o cerrar proceso:

```bash
# Encontrar proceso usando el puerto
lsof -i :8080

# Cerrar proceso
kill -9 <PID>

# O usar puerto diferente
vibercode mcp --port 8081
```

## 🏗️ Errores de Estructura de Proyecto

### Error: Directorio de salida no existe

**Problema**: Directorio especificado con `-o` no existe.

```bash
Error: output directory does not exist: ./non-existent
```

**Solución**: Crear directorio o usar uno existente:

```bash
mkdir -p ./mi-directorio
vibercode schema generate User -m test -o ./mi-directorio
```

### Error: Permisos insuficientes

**Problema**: Sin permisos para escribir archivos.

```bash
Error: permission denied: cannot create file
```

**Solución**: Verificar permisos del directorio:

```bash
# Verificar permisos
ls -la

# Cambiar propietario si es necesario
sudo chown -R $USER:$USER .

# O usar directorio con permisos
vibercode schema generate User -m test -o ~/mi-proyecto
```

## 🔍 Errores de Parsing

### Error: JSON de configuración inválido

**Problema**: Archivo `.vibercode-config.json` mal formado.

```bash
Error: invalid character '}' looking for beginning of object key string
```

**Solución**: Validar sintaxis JSON:

```bash
# Validar JSON
cat .vibercode-config.json | jq .

# O recrear archivo
rm .vibercode-config.json
# ViberCode creará uno nuevo con valores por defecto
```

## 🛠️ Herramientas de Diagnóstico

### Verificar instalación

```bash
# Verificar versión
vibercode version

# Verificar ayuda
vibercode help

# Verificar Go
go version
```

### Verificar dependencias

```bash
# Verificar módulos
go list -m all

# Verificar imports
go mod why golang.org/x/crypto
```

### Logs de debugging

```bash
# Activar modo verbose
vibercode schema generate User -m test --verbose

# Variables de debug
export VIBE_DEBUG=true
export GO_DEBUG=1
```

## 📝 Reportar Problemas

Si encuentras un error no documentado aquí:

1. **Recopila información**:
   ```bash
   vibercode version
   go version
   echo $GOOS $GOARCH
   ```

2. **Reproduce el error**:
   ```bash
   vibercode schema generate Test -m example --verbose
   ```

3. **Reporta en GitHub**: [Issues](https://github.com/vibercode/cli/issues)

## 🔗 Enlaces Útiles

- [**FAQ**](faq.md) - Preguntas frecuentes
- [**Debugging**](debugging.md) - Herramientas de depuración
- [**Configuración**](../user-guide/configuration.md) - Opciones de configuración

---

*¿No encontraste tu error? Consulta la [documentación completa](../README.md) o reporta el problema.*