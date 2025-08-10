package vibe

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// VibeLogger maneja el logging en el modo vibe
type VibeLogger struct {
	enabled   bool
	debugMode bool
	chatMode  bool
	logFile   *os.File
	prefix    string
}

// LogLevel representa el nivel de log
type LogLevel int

const (
	LogLevelInfo LogLevel = iota
	LogLevelWarning
	LogLevelError
	LogLevelDebug
)

// NewVibeLogger crea un nuevo logger para el modo vibe
func NewVibeLogger(chatMode bool) *VibeLogger {
	// En modo chat, desactivamos logs por defecto
	enabled := !chatMode

	// Si está en modo debug, siempre habilitamos
	if os.Getenv("VIBE_DEBUG") == "true" {
		enabled = true
	}

	var logFile *os.File
	if enabled {
		// Crear archivo de logs para no interferir con el chat
		logFile, _ = os.OpenFile("vibe.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	}

	return &VibeLogger{
		enabled:   enabled,
		debugMode: os.Getenv("VIBE_DEBUG") == "true",
		chatMode:  chatMode,
		logFile:   logFile,
		prefix:    "VibeCode",
	}
}

// SetChatMode activa o desactiva el modo chat
func (vl *VibeLogger) SetChatMode(enabled bool) {
	vl.chatMode = enabled
	if enabled {
		vl.enabled = false // Desactivar logs en modo chat
	}
}

// Info registra información
func (vl *VibeLogger) Info(message string, args ...interface{}) {
	vl.log(LogLevelInfo, message, args...)
}

// Warning registra advertencias
func (vl *VibeLogger) Warning(message string, args ...interface{}) {
	vl.log(LogLevelWarning, message, args...)
}

// Error registra errores
func (vl *VibeLogger) Error(message string, args ...interface{}) {
	vl.log(LogLevelError, message, args...)
}

// Debug registra información de debug
func (vl *VibeLogger) Debug(message string, args ...interface{}) {
	if vl.debugMode {
		vl.log(LogLevelDebug, message, args...)
	}
}

// ChatInfo registra información relevante para el chat (siempre visible)
func (vl *VibeLogger) ChatInfo(message string, args ...interface{}) {
	formattedMessage := fmt.Sprintf(message, args...)

	// En modo chat, solo mostramos información crítica
	if vl.chatMode {
		// Solo mostramos errores críticos y mensajes de estado importantes
		if strings.Contains(formattedMessage, "❌") ||
			strings.Contains(formattedMessage, "🚀") ||
			strings.Contains(formattedMessage, "✅ Preview server") {
			fmt.Printf("  %s\n", formattedMessage)
		}
	} else {
		fmt.Printf("[%s] %s\n", vl.prefix, formattedMessage)
	}
}

// log es el método interno para escribir logs
func (vl *VibeLogger) log(level LogLevel, message string, args ...interface{}) {
	if !vl.enabled {
		return
	}

	formattedMessage := fmt.Sprintf(message, args...)
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	levelStr := vl.getLevelString(level)

	logEntry := fmt.Sprintf("[%s] [%s] %s", timestamp, levelStr, formattedMessage)

	// Escribir al archivo si está disponible
	if vl.logFile != nil {
		vl.logFile.WriteString(logEntry + "\n")
		vl.logFile.Sync()
	}

	// En modo debug, también mostrar en consola
	if vl.debugMode {
		log.Print(logEntry)
	}
}

// getLevelString devuelve la string del nivel de log
func (vl *VibeLogger) getLevelString(level LogLevel) string {
	switch level {
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarning:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	case LogLevelDebug:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}

// Close cierra el logger
func (vl *VibeLogger) Close() {
	if vl.logFile != nil {
		vl.logFile.Close()
	}
}

// Logger global para el modo vibe
var GlobalVibeLogger *VibeLogger

// InitVibeLogger inicializa el logger global
func InitVibeLogger(chatMode bool) {
	GlobalVibeLogger = NewVibeLogger(chatMode)
}

// CloseVibeLogger cierra el logger global
func CloseVibeLogger() {
	if GlobalVibeLogger != nil {
		GlobalVibeLogger.Close()
	}
}
