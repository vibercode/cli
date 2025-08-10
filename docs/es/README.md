# ViberCode CLI - Documentación en Español

Bienvenido a la documentación completa de ViberCode CLI, una herramienta de línea de comandos para generar APIs Go con arquitectura limpia.

## 📚 Contenido de la Documentación

### Guía del Usuario
- [**Inicio Rápido**](user-guide/quickstart.md) - Instalación y primeros pasos
- [**Comandos CLI**](user-guide/cli-commands.md) - Referencia completa de comandos
- [**Generación de Esquemas**](user-guide/schema-generation.md) - Crear modelos y APIs
- [**Configuración**](user-guide/configuration.md) - Configurar el proyecto

### Tutoriales
- [**Tu Primera API Completa**](tutorials/primera-api-completa.md) - Tutorial paso a paso completo ⭐
- [**Integración con Bases de Datos**](tutorials/database-integration.md) - Conectar con diferentes DBs
- [**Autenticación y Autorización**](tutorials/auth-tutorial.md) - Implementar seguridad

### API y Desarrollo
- [**Arquitectura del Proyecto**](api/architecture.md) - Estructura y patrones
- [**Templates y Generadores**](api/templates.md) - Sistema de plantillas
- [**Extensiones**](api/extensions.md) - Crear extensiones personalizadas

### Solución de Problemas
- [**Errores Comunes**](troubleshooting/common-errors.md) - Problemas frecuentes
- [**Depuración**](troubleshooting/debugging.md) - Herramientas de debug
- [**FAQ**](troubleshooting/faq.md) - Preguntas frecuentes

## 🚀 Inicio Rápido

```bash
# Instalar ViberCode CLI
go install github.com/vibercode/cli@latest

# Generar un nuevo proyecto API
vibercode generate api mi-proyecto

# Generar un esquema de usuario
vibercode schema generate User -m mi-modulo -d postgres
```

## 🔧 Características Principales

- ✅ **Generación automática de código** - APIs completas con CRUD
- ✅ **Arquitectura limpia** - Separación clara de responsabilidades  
- ✅ **Múltiples bases de datos** - PostgreSQL, MySQL, SQLite, MongoDB
- ✅ **Templates personalizables** - Adapta el código a tus necesidades
- ✅ **Integración MCP** - Servidor Model Context Protocol
- ✅ **Chat AI interactivo** - Asistente de desarrollo integrado

## 📖 Enlaces Útiles

- [**English Documentation**](../en/README.md) - Documentación en inglés
- [**GitHub Repository**](https://github.com/vibercode/cli) - Código fuente
- [**Issues & Support**](https://github.com/vibercode/cli/issues) - Reportar problemas

---

*Generado con ViberCode CLI 🚀*