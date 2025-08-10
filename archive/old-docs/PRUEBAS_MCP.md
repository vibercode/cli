# 🎯 Guía de Pruebas MCP - ViberCode

## **Paso 1: Probar MCP Localmente** (`vibercode mcp` + React Editor)

### 🛠️ **Prerrequisitos**

- ✅ **Node.js** instalado (`node --version`)
- ✅ **React Editor** en `vibercode/editor/`
- ✅ **Terminal** disponible

### 🚀 **Opción A: Prueba Rápida MCP (Solo Servidor)**

```bash
# En vibercode-cli-go/
chmod +x test-mcp-quick.sh
./test-mcp-quick.sh
```

**¿Qué hace?**

- ✅ Inicia MCP Server simulado en `http://localhost:3002`
- ✅ Ejecuta pruebas automáticas con `curl`
- ✅ Muestra respuestas JSON del servidor
- ✅ Disponible para pruebas manuales

**Resultado esperado:**

```
🚀 Prueba Rápida MCP - ViberCode
================================
🟣 Iniciando MCP Server rápido en puerto 3002...
✅ MCP Server iniciado (PID: 12345)

✅ MCP Server listo para pruebas!
🟣 URL: http://localhost:3002

🧪 Ejecutando pruebas automáticas...

1️⃣ Probando: GET /tools
{
  "tools": [
    { "name": "generate_code", "description": "Genera código Go" },
    { "name": "vibe_start", "description": "Inicia modo vibe" },
    { "name": "project_status", "description": "Estado del proyecto" }
  ],
  "status": "ready"
}

2️⃣ Probando: POST /call - generate_code
{
  "success": true,
  "tool": "generate_code",
  "data": {
    "message": "generate_code ejecutado exitosamente",
    "timestamp": "2024-01-15T10:30:00.000Z",
    "project_id": "mcp-1705320600000"
  }
}

✅ Pruebas automáticas completadas!
```

### 🎨 **Opción B: Prueba Completa (MCP + WebSocket + React Editor)**

```bash
# En vibercode-cli-go/
chmod +x test-mcp-local.sh
./test-mcp-local.sh
```

**¿Qué hace?**

- ✅ Inicia **MCP Server** en puerto 3002
- ✅ Inicia **WebSocket Server** en puerto 3001
- ✅ Inicia **React Editor** en puerto 5173
- ✅ Abre automáticamente el navegador
- ✅ Instala dependencias si es necesario

**Resultado esperado:**

```
🎯 Probando MCP Localmente - ViberCode
=====================================

🔍 Verificando prerrequisitos...
✅ React Editor encontrado en: /Users/.../vibercode/editor
✅ Node.js encontrado

🚀 Iniciando servicios de prueba MCP...
🟣 Iniciando MCP Server simulado en puerto 3002...
✅ MCP Server simulado iniciado (PID: 12345)
📡 Iniciando WebSocket server simulado en puerto 3001...
✅ WebSocket simulado iniciado (PID: 12346)
🎨 Iniciando React Editor...
✅ React Editor iniciado (PID: 12347)

✅ Entorno de prueba MCP listo!
🟣 MCP Server: http://localhost:3002
📡 WebSocket: ws://localhost:3001/ws
🎨 React Editor: http://localhost:5173
```

### 🧪 **Pasos de Prueba en React Editor**

1. **Abrir el Editor** (se abre automáticamente)

   - URL: `http://localhost:5173`

2. **Agregar Componentes al Canvas**

   - Arrastra algunos componentes (Button, Input, Card, etc.)
   - Configura propiedades básicas

3. **Abrir Generation Panel**

   - Click en el botón "Generate Code" o similar
   - Se abre el panel de generación

4. **Seleccionar MCP Mode**

   - En "Connection Mode", click en **"🟣 MCP Server"**
   - Debería mostrar: `localhost:3002 ✅`

5. **Configurar Proyecto**

   - Project Name: `mi-proyecto-mcp`
   - Database: `PostgreSQL`
   - Features: Seleccionar las que quieras

6. **Generar Código**
   - Click en **"Generate API"**
   - Verificar respuesta exitosa

### 🔍 **Verificación Manual**

**Terminal 1** (mientras corre el script):

```bash
# Ver logs del MCP Server
📤 MCP Call: generate_code { project_name: 'mi-proyecto-mcp', database: 'postgres', ... }
```

**Terminal 2** (comandos manuales):

```bash
# Probar herramientas disponibles
curl http://localhost:3002/tools

# Probar generación de código
curl -X POST http://localhost:3002/call \
  -H "Content-Type: application/json" \
  -d '{"tool":"generate_code","params":{"project_name":"test-manual"}}'

# Probar estado del proyecto
curl -X POST http://localhost:3002/call \
  -H "Content-Type: application/json" \
  -d '{"tool":"project_status","params":{}}'
```

### ✅ **Resultados Esperados**

#### **En el React Editor:**

- ✅ Conexión MCP muestra estado "conectado" (✅)
- ✅ Generation panel permite seleccionar MCP mode
- ✅ Al generar, aparece mensaje de éxito
- ✅ Se pueden ver logs en la consola del navegador

#### **En la Terminal:**

- ✅ MCP Server responde a requests HTTP
- ✅ WebSocket acepta conexiones
- ✅ React Editor se conecta sin errores

#### **Respuestas JSON del MCP:**

```json
{
  "success": true,
  "tool": "generate_code",
  "data": {
    "message": "generate_code ejecutado exitosamente",
    "project_id": "mcp-1705320600000"
  }
}
```

### 🛑 **Detener las Pruebas**

```bash
# En cualquier terminal donde corra el script
Ctrl + C

# Resultado:
🛑 Cerrando servicios de prueba...
🟣 Cerrando MCP Server...
🔌 Cerrando WebSocket server...
🎨 Cerrando React Editor...
✅ Limpieza completada
```

### 🚨 **Resolución de Problemas**

#### **Error: "React Editor no encontrado"**

```bash
# Verificar ruta
ls ../../../vibercode/editor/package.json

# Si no existe, actualizar EDITOR_PATH en el script
```

#### **Error: "Node.js no encontrado"**

```bash
# Instalar Node.js
# macOS: brew install node
# o descargar desde https://nodejs.org
```

#### **Error: "Puerto ocupado"**

```bash
# Matar procesos en puertos específicos
lsof -ti:3002 | xargs kill -9  # MCP
lsof -ti:3001 | xargs kill -9  # WebSocket
lsof -ti:5173 | xargs kill -9  # React Editor
```

#### **Error: "Dependencias de React faltantes"**

```bash
# Instalar manualmente
cd ../../../vibercode/editor
pnpm install  # o npm install
```

### 🎯 **Próximo Paso**

Una vez que **Paso 1** funcione correctamente:

- ✅ MCP Server responde
- ✅ React Editor se conecta
- ✅ Generación funciona sin errores

Pasaremos al **Paso 3**: Desarrollar funciones MCP reales (reemplazar simulaciones por comunicación real).

---

**¿Necesitas ayuda con algún error específico?** 🤔
