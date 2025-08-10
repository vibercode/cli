# VibeCode Component Mode Setup

## 🎨 Nueva Funcionalidad: Modo Component

Se ha implementado un nuevo modo especializado para el desarrollo de componentes UI que se activa con el parámetro `component`.

### Comandos

```bash
# Modo general (por defecto)
vibercode vibe

# Modo component (nuevo)
vibercode vibe component
```

### Características del Modo Component

#### 1. **Enfoque 100% en Componentes**

- Todas las respuestas de la IA están orientadas a componentes UI
- Análisis visual y de diseño en cada respuesta
- Validación automática de estructura de componentes
- Sugerencias de mejora de diseño

#### 2. **Preview Server Mejorado**

- Página web interactiva en `http://localhost:3001`
- Canvas visual para mostrar componentes
- Conexión WebSocket en tiempo real
- Información de estado del canvas

#### 3. **Respuestas Contextuales**

- Análisis del estado actual del canvas
- Sugerencias de posicionamiento inteligente
- Feedback sobre jerarquía visual
- Validación de propiedades de componentes

### Estructura del Mensaje

En modo component, todos los mensajes enviados al preview siguen esta estructura:

```json
{
  "type": "ui_update",
  "action": "add_component|update_component|update_theme|remove_component",
  "data": {
    "type": "button|text|card|form|hero|gallery|etc",
    "category": "atom|molecule|organism",
    "properties": {
      // Propiedades específicas del componente
    },
    "position": { "x": 200, "y": 200 },
    "size": { "w": 160, "h": 40 }
  },
  "explanation": "Descripción del cambio realizado"
}
```

### Componentes Disponibles

#### Atoms (Básicos)

- `button`: Botón interactivo
- `text`: Texto básico
- `animated-text`: Texto con animaciones
- `image`: Imagen con opciones
- `input`: Campo de entrada

#### Molecules (Compuestos)

- `card`: Tarjeta con imagen y contenido
- `form`: Formulario con campos
- `navigation`: Barra de navegación

#### Organisms (Complejos)

- `hero`: Sección hero completa
- `gallery`: Galería de imágenes

### Comandos de Prueba

```bash
# Agregar componentes
"agregar botón"
"agregar tarjeta con imagen"
"crear formulario de contacto"
"agregar hero section"

# Cambiar temas
"cambiar tema a azul"
"tema oceánico"
"colores profesionales"

# Información del canvas
"estado del canvas"
"qué componentes hay"
"análisis del diseño"
```

### Archivos Modificados

1. **`cmd/vibe.go`**: Acepta parámetro `component`
2. **`internal/vibe/vibe.go`**: Maneja modos diferentes
3. **`internal/vibe/chat.go`**: Enfoque contextual por modo
4. **`internal/vibe/preview.go`**: Servidor web con página HTML
5. **`internal/vibe/prompts/system.md`**: Instrucciones específicas por modo
6. **`internal/vibe/prompts/prompt_loader.go`**: Contexto de modo en prompts

### Flujo de Trabajo

1. **Inicio**: `vibercode vibe component`
2. **Preview**: Abrir `http://localhost:3001`
3. **Interacción**: Comandos en terminal
4. **Visualización**: Componentes aparecen en tiempo real
5. **Iteración**: Modificar y mejorar componentes

### Validación

El sistema valida automáticamente:

- Estructura JSON correcta
- Propiedades de componentes válidas
- Posicionamiento sin solapamiento
- Consistencia de temas
- Categorías de componentes (atom/molecule/organism)

### Características Técnicas

- **WebSocket**: Actualizaciones en tiempo real
- **CORS**: Configurado para desarrollo
- **Reconexión**: Automática en caso de pérdida de conexión
- **Validación**: JSON Schema para componentes
- **Fallbacks**: Respuestas offline si Claude no está disponible

### Próximos Pasos

1. Probar funcionalidad básica
2. Agregar más tipos de componentes
3. Mejorar validaciones
4. Implementar sistema de templates
5. Agregar modo colaborativo
