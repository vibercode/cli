# Comandos CLI - Referencia Completa

ViberCode CLI ofrece varios comandos para generar y gestionar proyectos Go API completos con arquitectura limpia, middleware de seguridad, tests comprehensivos, configuraciones de deployment y componentes UI modernos.

## 📋 Tabla de Contenidos

- [Comandos Principales](#comandos-principales)
- [Generación de Código](#generación-de-código)
  - [APIs Completas](#apis-completas)
  - [Recursos CRUD](#recursos-crud)
  - [Middleware](#middleware)
  - [Componentes UI](#componentes-ui)
  - [Tests](#tests)
  - [Deployment](#deployment)
- [Generación de Esquemas](#generación-de-esquemas)
- [Servidor MCP](#servidor-mcp)
- [Opciones Globales](#opciones-globales)
- [Ejemplos Detallados](#ejemplos-detallados)

## 🎯 Comandos Principales

### `vibercode help`

Muestra la ayuda general del CLI.

```bash
vibercode help
vibercode --help
vibercode -h
```

### `vibercode version`

Muestra la versión actual del CLI.

```bash
vibercode version
```

## 🏗️ Generación de Código

### `vibercode generate api`

Genera una API Go completa con arquitectura limpia, configuración de Docker y documentación.

#### Sintaxis

```bash
vibercode generate api
```

#### Funcionalidades

- **Arquitectura Limpia**: Handlers, services, repositories separados
- **Configuración de Base de Datos**: PostgreSQL, MySQL, SQLite, MongoDB
- **Setup Docker**: Dockerfile y docker-compose.yml
- **Variables de Entorno**: Configuración completa
- **Documentación**: README y guías de uso

#### Ejemplo

```bash
vibercode generate api
# Sigue las instrucciones interactivas
```

---

### `vibercode generate resource`

Genera un recurso CRUD completo con todas las capas de la arquitectura.

#### Sintaxis

```bash
vibercode generate resource
```

#### Funcionalidades

- **Modelo GORM**: Con validaciones y relaciones
- **Handlers HTTP**: Endpoints completos (GET, POST, PUT, DELETE)
- **Capa de Servicio**: Lógica de negocio
- **Repositorio**: Operaciones de base de datos
- **Validación**: Validación de datos y manejo de errores

#### Ejemplo

```bash
vibercode generate resource
# Interactivamente define campos y relaciones
```

---

### `vibercode generate middleware`

Genera componentes de middleware para autenticación, logging, CORS y rate limiting.

#### Sintaxis

```bash
vibercode generate middleware [flags]
```

#### Flags Disponibles

| Flag | Descripción | Valores | Por Defecto |
|------|-------------|---------|-------------|
| `--type` | Tipo de middleware | `auth`, `logging`, `cors`, `rate-limit` | - |
| `--name` | Nombre personalizado | string | - |
| `--custom` | Generar middleware personalizado | boolean | `false` |
| `--preset` | Preset predefinido | `api-security`, `web-app`, `microservice`, `public-api` | - |

#### Tipos de Middleware

- **Auth**: JWT, API Key, Session, Basic Auth
- **Logging**: Structured logging con niveles
- **CORS**: Configuración Cross-Origin Resource Sharing
- **Rate Limit**: Limitación de velocidad por IP/usuario
- **Custom**: Plantilla para middleware personalizado

#### Ejemplos

```bash
# Middleware de autenticación JWT
vibercode generate middleware --type auth

# Preset de seguridad para API
vibercode generate middleware --preset api-security

# Middleware personalizado
vibercode generate middleware --name CustomValidator --custom

# CORS con configuración específica
vibercode generate middleware --type cors
```

---

### `vibercode generate ui`

Genera componentes UI siguiendo la metodología Atomic Design.

#### Sintaxis

```bash
vibercode generate ui [flags]
```

#### Flags Disponibles

| Flag | Descripción | Valores | Por Defecto |
|------|-------------|---------|-------------|
| `--atomic-design` | Estructura Atomic Design completa | boolean | `false` |
| `--framework` | Framework frontend | `react`, `vue`, `angular` | `react` |
| `--typescript` | Generar componentes TypeScript | boolean | `true` |
| `--storybook` | Incluir historias Storybook | boolean | `false` |

#### Estructura Atomic Design

- **Atoms**: Botones, inputs, labels básicos
- **Molecules**: Formularios, cards, navegación
- **Organisms**: Headers, sidebars, secciones
- **Templates**: Layouts de página y grids
- **Pages**: Páginas completas con datos

#### Ejemplos

```bash
# Estructura completa con React y TypeScript
vibercode generate ui --atomic-design --typescript

# Componentes Vue con Storybook
vibercode generate ui --framework vue --storybook

# Setup básico Angular
vibercode generate ui --framework angular
```

---

### `vibercode generate test`

Genera suite de tests completa con utilidades y mocks.

#### Sintaxis

```bash
vibercode generate test [flags]
```

#### Flags Disponibles

| Flag | Descripción | Valores | Por Defecto |
|------|-------------|---------|-------------|
| `--type` | Tipo de test | `unit`, `integration`, `benchmark`, `mock`, `utils` | - |
| `--framework` | Framework de testing | `testify`, `ginkgo`, `goconvey` | `testify` |
| `--target` | Objetivo del test | `handler`, `service`, `repository`, `middleware`, `api` | - |
| `--name` | Nombre del componente | string | - |
| `--full-suite` | Suite completa de tests | boolean | `false` |
| `--with-mocks` | Incluir generación de mocks | boolean | `false` |
| `--with-utils` | Incluir utilidades de test | boolean | `false` |
| `--with-bench` | Incluir tests de benchmark | boolean | `false` |
| `--bdd` | Estilo BDD (Ginkgo/Gomega) | boolean | `false` |

#### Tipos de Tests

- **Unit**: Tests unitarios para componentes individuales
- **Integration**: Tests de integración API end-to-end
- **Benchmark**: Tests de rendimiento y benchmarking
- **Mock**: Generación de mocks para dependencias
- **Utils**: Utilidades de testing (database, server, client)

#### Frameworks Soportados

- **Testify**: Framework estándar con assertions y mocks
- **Ginkgo/Gomega**: Framework BDD con sintaxis expresiva
- **GoConvey**: Framework con interfaz web para testing

#### Ejemplos

```bash
# Suite completa de tests
vibercode generate test --full-suite

# Test unitario para handler específico
vibercode generate test --type unit --target handler --name User

# Tests de integración con mocks
vibercode generate test --type integration --name User --with-mocks

# Tests BDD con Ginkgo
vibercode generate test --framework ginkgo --bdd --target service --name User

# Tests de benchmark
vibercode generate test --type benchmark --target repository --name Product
```

---

### `vibercode generate deployment`

Genera configuraciones completas de deployment para Docker, Kubernetes y cloud.

#### Sintaxis

```bash
vibercode generate deployment [flags]
```

#### Flags Disponibles

| Flag | Descripción | Valores | Por Defecto |
|------|-------------|---------|-------------|
| `--type` | Tipo de deployment | `docker`, `kubernetes`, `cloud`, `cicd` | - |
| `--provider` | Proveedor cloud/CI | `aws`, `gcp`, `azure`, `github-actions`, `gitlab-ci` | - |
| `--service` | Servicio cloud específico | `ecs`, `run`, `containers`, etc. | - |
| `--namespace` | Namespace Kubernetes | string | `default` |
| `--environment` | Entorno objetivo | `development`, `staging`, `production` | `production` |
| `--multi-stage` | Build Docker multi-stage | boolean | `false` |
| `--optimize` | Optimización de imagen | boolean | `false` |
| `--security` | Hardening de seguridad | boolean | `false` |
| `--with-ingress` | Incluir Ingress Kubernetes | boolean | `false` |
| `--with-secrets` | Incluir Secrets y ConfigMaps | boolean | `false` |
| `--with-hpa` | Incluir Horizontal Pod Autoscaler | boolean | `false` |
| `--full-suite` | Suite completa de deployment | boolean | `false` |

#### Tipos de Deployment

- **Docker**: Dockerfiles multi-stage con optimización y seguridad
- **Kubernetes**: Manifests completos con scaling y monitoring
- **Cloud**: Configuraciones para AWS ECS/Fargate, GCP Cloud Run, Azure Container Instances
- **CI/CD**: Pipelines para GitHub Actions, GitLab CI, Jenkins, CircleCI

#### Proveedores Cloud

- **AWS**: ECS, Fargate, EKS, Lambda + Terraform/CloudFormation
- **GCP**: Cloud Run, GKE, App Engine + Terraform
- **Azure**: Container Instances, AKS, Web Apps + ARM Templates

#### Ejemplos

```bash
# Suite completa para AWS
vibercode generate deployment --full-suite --provider aws

# Docker multi-stage con seguridad
vibercode generate deployment --type docker --multi-stage --security

# Kubernetes con Ingress y HPA
vibercode generate deployment --type kubernetes --with-ingress --with-hpa

# Cloud Run en GCP
vibercode generate deployment --type cloud --provider gcp --service run

# Pipeline GitHub Actions
vibercode generate deployment --type cicd --provider github-actions
```

---

## 🔌 Servidor MCP

### `vibercode mcp`

Inicia el servidor MCP (Model Context Protocol) para integración con agentes de IA.

```bash
vibercode mcp
```

#### Variables de Entorno

```bash
export ANTHROPIC_API_KEY="tu-api-key"
export VIBE_DEBUG="true"
```

#### Funcionalidades del Servidor MCP

- **Edición de componentes en vivo**
- **Chat integrado con IA**
- **Generación automática de código**
- **Gestión de estado del proyecto**

## ⚙️ Opciones Globales

### Flags Comunes

| Flag | Descripción | Ejemplo |
|------|-------------|---------|
| `--verbose` | Salida detallada | `vibercode schema generate User --verbose` |
| `--config` | Archivo de configuración | `vibercode --config ./custom.json` |
| `--dry-run` | Simular sin ejecutar | `vibercode schema generate --dry-run` |

## 📊 Ejemplos Detallados

### Ejemplo 1: API de Blog Completa

```bash
# 1. Crear API base
vibercode generate api
# Seleccionar: blog-api, PostgreSQL, puerto 8080

# 2. Generar recursos principales
vibercode generate resource
# Crear: User (name, email, password, role)
vibercode generate resource
# Crear: Post (title, content, author_id, published)
vibercode generate resource
# Crear: Comment (content, post_id, user_id)

# 3. Agregar middleware de seguridad
vibercode generate middleware --preset api-security

# 4. Generar tests completos
vibercode generate test --full-suite --with-mocks

# 5. Setup deployment con Docker
vibercode generate deployment --type docker --multi-stage --security
```

### Ejemplo 2: Microservicio E-commerce

```bash
# API base con MongoDB
vibercode generate api
# Seleccionar: ecommerce-service, MongoDB, puerto 3000

# Recursos principales
vibercode generate resource  # Product
vibercode generate resource  # Category
vibercode generate resource  # Order
vibercode generate resource  # Customer

# Middleware personalizado para e-commerce
vibercode generate middleware --preset microservice
vibercode generate middleware --type rate-limit

# Tests de integración
vibercode generate test --type integration --with-utils

# Deployment en Kubernetes
vibercode generate deployment --type kubernetes --with-ingress --with-hpa
```

### Ejemplo 3: API Pública con Documentación

```bash
# API pública con documentación completa
vibercode generate api

# Middleware para API pública
vibercode generate middleware --preset public-api
vibercode generate middleware --type cors

# UI para documentación
vibercode generate ui --framework react --typescript

# Tests y benchmarks
vibercode generate test --full-suite --with-bench

# Deployment completo en AWS
vibercode generate deployment --full-suite --provider aws
```

### Ejemplo 4: Desarrollo Rápido con SQLite

```bash
# Setup rápido para desarrollo
vibercode generate api
# Seleccionar: dev-api, SQLite, puerto 8000

# Recursos de desarrollo
vibercode generate resource  # Task
vibercode generate resource  # Project

# Middleware básico
vibercode generate middleware --type logging

# Tests unitarios simples
vibercode generate test --type unit --framework testify

# Docker para desarrollo
vibercode generate deployment --type docker
```

## 🏗️ Arquitectura Generada

### Estructura de Proyecto API

```
mi-api/
├── cmd/
│   └── server/
│       └── main.go              # Punto de entrada
├── internal/
│   ├── handlers/                # Capa HTTP (Gin)
│   │   ├── user_handler.go
│   │   └── middleware/
│   ├── services/                # Lógica de negocio
│   │   └── user_service.go
│   ├── repositories/            # Acceso a datos
│   │   └── user_repository.go
│   └── models/                  # Modelos de dominio
│       └── user.go
├── pkg/
│   ├── database/                # Conexión DB
│   ├── config/                  # Configuración
│   └── utils/                   # Utilidades
├── tests/                       # Tests organizados
│   ├── unit/
│   ├── integration/
│   └── utils/
├── deployment/                  # Configuraciones deployment
│   ├── docker/
│   ├── kubernetes/
│   └── cloud/
├── docs/                        # Documentación
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md
```

### Tipos de Campo Soportados

| Tipo | Go Type | Validación | Base de Datos |
|------|---------|------------|---------------|
| `string` | `string` | length, required | VARCHAR |
| `text` | `string` | length | TEXT |
| `number` | `int` | min, max, required | INTEGER |
| `float` | `float64` | min, max, precision | DECIMAL |
| `boolean` | `bool` | - | BOOLEAN |
| `date` | `time.Time` | format | TIMESTAMP |
| `uuid` | `uuid.UUID` | format | UUID |
| `json` | `json.RawMessage` | valid JSON | JSON/JSONB |
| `relation` | `uint` | foreign key | FOREIGN KEY |
| `relation-array` | `[]Model` | - | JOIN TABLE |

### Middleware Incluido

- **JWT Authentication**: Validación de tokens
- **API Key Auth**: Autenticación por clave API
- **Request Logging**: Log estructurado de requests
- **CORS**: Cross-Origin Resource Sharing
- **Rate Limiting**: Límites por IP/usuario
- **Error Handling**: Manejo centralizado de errores
- **Request Validation**: Validación automática de payloads

## 🚨 Manejo de Errores

### Errores Comunes

```bash
# Error: Módulo requerido
vibercode schema generate User
# Error: flag needs an argument: -m

# Solución: Especificar módulo
vibercode schema generate User -m mi-api

# Error: Base de datos no soportada
vibercode schema generate User -m test -d oracle
# Error: unsupported database type: oracle

# Solución: Usar base de datos soportada
vibercode schema generate User -m test -d postgres
```

### Códigos de Salida

| Código | Descripción |
|--------|-------------|
| 0 | Éxito |
| 1 | Error general |
| 2 | Error de argumentos |
| 3 | Error de archivo/directorio |
| 4 | Error de base de datos |
| 5 | Error de template |
| 6 | Error de compilación |

### Mejores Prácticas

#### Generación de Código

1. **Usa nombres descriptivos** en singular para recursos (User, Product, Order)
2. **Sigue convenciones Go** para naming (PascalCase para exports, camelCase para locales)
3. **Define relaciones claramente** especificando foreign keys y joins
4. **Incluye validaciones** apropiadas para cada tipo de campo

#### Middleware

1. **Combina presets** para configuraciones complejas de seguridad
2. **Personaliza rate limiting** según tu caso de uso específico
3. **Configura CORS** correctamente para tu frontend
4. **Usa JWT** para APIs stateless, sessions para aplicaciones web

#### Testing

1. **Genera suite completa** para proyectos de producción
2. **Usa mocks** para aislar tests unitarios de dependencias externas
3. **Incluye tests de integración** para validar el flujo completo
4. **Agrega benchmarks** para componentes críticos de performance

#### Deployment

1. **Usa multi-stage builds** para imágenes Docker optimizadas
2. **Habilita security hardening** para entornos de producción
3. **Configura monitoring** y health checks apropiados
4. **Implementa auto-scaling** con HPA en Kubernetes

## 📝 Notas Importantes

1. **Arquitectura Limpia**: Todo el código generado sigue principios de clean architecture
2. **Directorios**: El CLI crea automáticamente la estructura de directorios necesaria
3. **Sobrescritura**: Los archivos existentes se sobrescriben sin confirmación
4. **Dependencias**: Se agregan automáticamente las dependencias Go necesarias
5. **Configuración**: Variables de entorno y archivos de config se generan automáticamente
6. **Documentación**: Se incluye documentación completa y ejemplos de uso
7. **Tests**: Los tests generados incluyen casos de éxito y error
8. **Security**: El código incluye mejores prácticas de seguridad por defecto
9. **Performance**: Templates optimizados para mejor rendimiento
10. **Compatibility**: Compatible con las últimas versiones de Go (1.21+)

## 🔗 Enlaces Relacionados

- [**Generación de APIs**](api-generation.md) - Guía detallada de generación de APIs
- [**Middleware Guide**](middleware-guide.md) - Configuración avanzada de middleware
- [**Testing Guide**](testing-guide.md) - Guía completa de testing
- [**Deployment Guide**](deployment-guide.md) - Estrategias de deployment
- [**UI Components**](ui-components.md) - Componentes UI con Atomic Design
- [**Configuración**](configuration.md) - Opciones de configuración global
- [**Solución de Problemas**](../troubleshooting/common-errors.md) - Errores comunes y soluciones
- [**Ejemplos Avanzados**](../examples/advanced-examples.md) - Casos de uso complejos
- [**API Reference**](../api/cli-reference.md) - Referencia completa de comandos

---

*Para más ayuda, ejecuta `vibercode help` o visita la documentación completa.*