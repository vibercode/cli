# Inicio Rápido - ViberCode CLI

Esta guía te ayudará a comenzar con ViberCode CLI en pocos minutos.

## 🚦 Prerequisitos

- Go 1.19 o superior
- Git
- Un editor de código (VS Code, GoLand, etc.)

## 📦 Instalación

### Opción 1: Instalación desde código fuente

```bash
# Clonar el repositorio
git clone https://github.com/vibercode/cli.git
cd vibercode-cli-go

# Compilar el binario
go build -o vibercode main.go

# Hacer ejecutable (Linux/macOS)
chmod +x vibercode

# Mover a PATH (opcional)
sudo mv vibercode /usr/local/bin/
```

### Opción 2: Instalación directa con Go

```bash
go install github.com/vibercode/cli@latest
```

## 🎯 Tu Primera API

### 1. Crear un Proyecto API

```bash
# Crear directorio del proyecto
mkdir mi-primera-api
cd mi-primera-api

# Generar proyecto base
vibercode generate api mi-primera-api
```

### 2. Generar un Esquema de Usuario

```bash
# Generar modelo User con PostgreSQL
vibercode schema generate User -m mi-primera-api -d postgres

# El comando te preguntará por confirmación
# Responde 'y' para continuar
```

### 3. Estructura del Proyecto Generado

```
mi-primera-api/
├── cmd/
│   └── server/
│       └── main.go          # Punto de entrada
├── internal/
│   ├── handlers/            # Controladores HTTP
│   ├── services/            # Lógica de negocio
│   ├── repositories/        # Acceso a datos
│   └── models/              # Modelos de datos
├── pkg/
│   ├── database/            # Conexión DB
│   └── config/              # Configuración
├── go.mod                   # Dependencias Go
└── docker-compose.yml       # Configuración Docker
```

### 4. Ejecutar el Proyecto

```bash
# Instalar dependencias
go mod tidy

# Configurar base de datos (PostgreSQL)
docker-compose up -d

# Ejecutar el servidor
go run cmd/server/main.go
```

## 🔧 Comandos Básicos

### Generar Proyecto API

```bash
vibercode generate api [nombre-proyecto]
```

### Generar Esquemas

```bash
# PostgreSQL
vibercode schema generate [NombreEsquema] -d postgres -m [modulo]

# MySQL  
vibercode schema generate [NombreEsquema] -d mysql -m [modulo]

# SQLite
vibercode schema generate [NombreEsquema] -d sqlite -m [modulo]

# MongoDB
vibercode schema generate [NombreEsquema] -d mongodb -m [modulo]
```

### Servidor MCP (Integración IA)

```bash
# Iniciar servidor MCP
vibercode mcp

# En otra terminal, probar integración
./test-mcp-server.sh
```

## 🎨 Personalización

### Configuración del Proyecto

Edita el archivo `.vibercode-config.json`:

```json
{
  "project_name": "mi-api",
  "database_type": "postgres",
  "port": 8080,
  "enable_auth": true,
  "enable_swagger": true
}
```

### Variables de Entorno

Copia y configura el archivo de entorno:

```bash
cp .env.example .env
```

Edita `.env`:

```env
DATABASE_URL=postgres://user:pass@localhost:5432/mydb
JWT_SECRET=tu-clave-secreta-aqui
PORT=8080
```

## 🚀 Próximos Pasos

1. **[Comandos CLI](cli-commands.md)** - Aprende todos los comandos disponibles
2. **[Generación de Esquemas](schema-generation.md)** - Domina la creación de modelos
3. **[Configuración](configuration.md)** - Personaliza tu proyecto
4. **[Tutorial Completo](../tutorials/first-api.md)** - Construye una API completa

## ❓ ¿Problemas?

- Revisa la [**Solución de Problemas**](../troubleshooting/common-errors.md)
- Consulta las [**Preguntas Frecuentes**](../troubleshooting/faq.md)
- Reporta issues en [**GitHub**](https://github.com/vibercode/cli/issues)

---

¡Felicidades! Ya tienes ViberCode CLI funcionando. 🎉