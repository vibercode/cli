package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// Manager handles configuration loading, validation, and environment management
type Manager struct {
	config     *Config
	validator  *Validator
	configPath string
	envPath    string
	mu         sync.RWMutex
	callbacks  []ReloadCallback
}

// ReloadCallback is called when configuration is reloaded
type ReloadCallback func(oldConfig, newConfig *Config) error

// LoadOptions contains options for loading configuration
type LoadOptions struct {
	ConfigPath     string   // Path to main config file
	Environment    string   // Environment name (development, staging, production)
	EnvFiles       []string // Additional env files to load
	WatchForChanges bool     // Enable hot-reloading
	Validate       bool     // Validate configuration after loading
}

// NewManager creates a new configuration manager
func NewManager() *Manager {
	return &Manager{
		validator: NewValidator(),
		callbacks: make([]ReloadCallback, 0),
	}
}

// Load loads configuration from files and environment variables
func (m *Manager) Load(opts LoadOptions) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Set default environment if not provided
	if opts.Environment == "" {
		opts.Environment = os.Getenv("APP_ENV")
		if opts.Environment == "" {
			opts.Environment = "development"
		}
	}

	// Load base configuration
	config := DefaultConfig()
	config.Environment = opts.Environment

	// Load configuration file if provided
	if opts.ConfigPath != "" {
		if err := m.loadConfigFile(config, opts.ConfigPath); err != nil {
			return fmt.Errorf("failed to load config file: %w", err)
		}
		m.configPath = opts.ConfigPath
	}

	// Load environment-specific configuration
	envConfigPath := m.getEnvironmentConfigPath(opts.ConfigPath, opts.Environment)
	if m.fileExists(envConfigPath) {
		if err := m.loadConfigFile(config, envConfigPath); err != nil {
			return fmt.Errorf("failed to load environment config: %w", err)
		}
	}

	// Load environment files (simplified without godotenv for now)
	envFiles := m.getEnvironmentFiles(opts.EnvFiles, opts.Environment)
	for _, envFile := range envFiles {
		if m.fileExists(envFile) {
			if err := m.loadSimpleEnvFile(envFile); err != nil {
				return fmt.Errorf("failed to load env file %s: %w", envFile, err)
			}
		}
	}

	// Override with environment variables
	if err := m.loadFromEnv(config); err != nil {
		return fmt.Errorf("failed to load from environment: %w", err)
	}

	// Validate configuration if requested
	if opts.Validate {
		if err := m.validator.ValidateConfig(config); err != nil {
			return fmt.Errorf("configuration validation failed: %w", err)
		}
	}

	m.config = config

	// Note: File watching would be implemented with fsnotify when needed
	// For now, hot-reloading is disabled

	return nil
}

// GetConfig returns the current configuration (thread-safe)
func (m *Manager) GetConfig() *Config {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	// Return a copy to prevent external modifications
	configCopy := *m.config
	return &configCopy
}

// UpdateConfig updates the configuration and triggers validation
func (m *Manager) UpdateConfig(updater func(*Config) error) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	oldConfig := *m.config
	
	if err := updater(m.config); err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	if err := m.validator.ValidateConfig(m.config); err != nil {
		// Restore old config on validation failure
		*m.config = oldConfig
		return fmt.Errorf("config validation failed after update: %w", err)
	}

	// Notify callbacks
	for _, callback := range m.callbacks {
		if err := callback(&oldConfig, m.config); err != nil {
			return fmt.Errorf("reload callback failed: %w", err)
		}
	}

	return nil
}

// AddReloadCallback adds a callback that's called when configuration is reloaded
func (m *Manager) AddReloadCallback(callback ReloadCallback) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.callbacks = append(m.callbacks, callback)
}

// Reload reloads the configuration from files
func (m *Manager) Reload() error {
	if m.configPath == "" {
		return fmt.Errorf("no config path set, cannot reload")
	}

	opts := LoadOptions{
		ConfigPath:      m.configPath,
		Environment:     m.config.Environment,
		WatchForChanges: false, // Don't restart watcher
		Validate:        true,
	}

	oldConfig := *m.config
	
	if err := m.Load(opts); err != nil {
		return err
	}

	// Notify callbacks
	for _, callback := range m.callbacks {
		if err := callback(&oldConfig, m.config); err != nil {
			return fmt.Errorf("reload callback failed: %w", err)
		}
	}

	return nil
}

// Close cleans up resources (placeholder for future file watcher cleanup)
func (m *Manager) Close() error {
	// Future: cleanup file watcher resources
	return nil
}

// loadConfigFile loads configuration from a file (JSON or YAML)
func (m *Manager) loadConfigFile(config *Config, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".json":
		return json.Unmarshal(data, config)
	case ".yaml", ".yml":
		return yaml.Unmarshal(data, config)
	default:
		return fmt.Errorf("unsupported config file format: %s", ext)
	}
}

// loadSimpleEnvFile loads environment variables from a file (simple implementation)
func (m *Manager) loadSimpleEnvFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Remove quotes if present
			if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
			   (strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
				value = value[1 : len(value)-1]
			}
			os.Setenv(key, value)
		}
	}

	return nil
}

// loadFromEnv loads configuration from environment variables
func (m *Manager) loadFromEnv(config *Config) error {
	return m.mapEnvToStruct(config, "")
}

// mapEnvToStruct recursively maps environment variables to struct fields
func (m *Manager) mapEnvToStruct(v interface{}, prefix string) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("expected pointer to struct")
	}

	rv = rv.Elem()
	rt := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		fieldType := rt.Field(i)

		if !field.CanSet() {
			continue
		}

		// Get field name from JSON tag or use field name
		fieldName := fieldType.Tag.Get("json")
		if fieldName == "" || fieldName == "-" {
			fieldName = fieldType.Name
		} else {
			fieldName = strings.Split(fieldName, ",")[0]
		}

		envKey := m.buildEnvKey(prefix, fieldName)

		switch field.Kind() {
		case reflect.Struct:
			// Recursively handle nested structs
			if err := m.mapEnvToStruct(field.Addr().Interface(), envKey); err != nil {
				return err
			}
		case reflect.String:
			if envValue := os.Getenv(envKey); envValue != "" {
				field.SetString(envValue)
			}
		case reflect.Int, reflect.Int64:
			if envValue := os.Getenv(envKey); envValue != "" {
				if fieldType.Type.String() == "time.Duration" {
					if duration, err := time.ParseDuration(envValue); err == nil {
						field.SetInt(int64(duration))
					}
				} else {
					if intValue, err := parseInt64(envValue); err == nil {
						field.SetInt(intValue)
					}
				}
			}
		case reflect.Bool:
			if envValue := os.Getenv(envKey); envValue != "" {
				field.SetBool(parseBool(envValue))
			}
		case reflect.Slice:
			if envValue := os.Getenv(envKey); envValue != "" {
				switch field.Type().Elem().Kind() {
				case reflect.String:
					values := strings.Split(envValue, ",")
					for i, value := range values {
						values[i] = strings.TrimSpace(value)
					}
					field.Set(reflect.ValueOf(values))
				}
			}
		}
	}

	return nil
}

// buildEnvKey builds environment variable key from prefix and field name
func (m *Manager) buildEnvKey(prefix, fieldName string) string {
	key := strings.ToUpper(fieldName)
	if prefix != "" {
		key = strings.ToUpper(prefix) + "_" + key
	}
	return key
}

// getEnvironmentConfigPath returns path for environment-specific config
func (m *Manager) getEnvironmentConfigPath(basePath, environment string) string {
	if basePath == "" {
		return ""
	}

	dir := filepath.Dir(basePath)
	ext := filepath.Ext(basePath)
	name := strings.TrimSuffix(filepath.Base(basePath), ext)
	
	return filepath.Join(dir, fmt.Sprintf("%s.%s%s", name, environment, ext))
}

// getEnvironmentFiles returns list of environment files to load
func (m *Manager) getEnvironmentFiles(additional []string, environment string) []string {
	files := []string{
		".env",
		fmt.Sprintf(".env.%s", environment),
		".env.local",
	}
	
	files = append(files, additional...)
	return files
}

// fileExists checks if a file exists
func (m *Manager) fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Note: File watching functionality will be implemented when fsnotify dependency is added
// For now, configuration reloading is manual only

// Helper functions

func parseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func parseBool(s string) bool {
	switch strings.ToLower(s) {
	case "true", "1", "yes", "on", "enabled":
		return true
	default:
		return false
	}
}

// SaveConfig saves the current configuration to a file
func (m *Manager) SaveConfig(filePath string) error {
	m.mu.RLock()
	config := m.config
	m.mu.RUnlock()

	return m.saveConfigToFile(config, filePath)
}

// saveConfigToFile saves configuration to a file
func (m *Manager) saveConfigToFile(config *Config, filePath string) error {
	ext := strings.ToLower(filepath.Ext(filePath))
	
	var data []byte
	var err error

	switch ext {
	case ".json":
		data, err = json.MarshalIndent(config, "", "  ")
	case ".yaml", ".yml":
		data, err = yaml.Marshal(config)
	default:
		return fmt.Errorf("unsupported config file format: %s", ext)
	}

	if err != nil {
		return err
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// GenerateEnvFile generates an environment file with current configuration
func (m *Manager) GenerateEnvFile(filePath string) error {
	m.mu.RLock()
	config := m.config
	m.mu.RUnlock()

	var lines []string
	m.generateEnvLines(config, "", &lines)

	content := strings.Join(lines, "\n")
	return os.WriteFile(filePath, []byte(content), 0644)
}

// generateEnvLines recursively generates environment variable lines
func (m *Manager) generateEnvLines(v interface{}, prefix string, lines *[]string) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return
	}

	rt := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		fieldType := rt.Field(i)

		if !field.CanInterface() {
			continue
		}

		// Get field name from JSON tag
		fieldName := fieldType.Tag.Get("json")
		if fieldName == "" || fieldName == "-" {
			fieldName = fieldType.Name
		} else {
			fieldName = strings.Split(fieldName, ",")[0]
		}

		envKey := m.buildEnvKey(prefix, fieldName)

		switch field.Kind() {
		case reflect.Struct:
			m.generateEnvLines(field.Interface(), envKey, lines)
		case reflect.String:
			if value := field.String(); value != "" {
				*lines = append(*lines, fmt.Sprintf("%s=%s", envKey, value))
			}
		case reflect.Int, reflect.Int64:
			if fieldType.Type.String() == "time.Duration" {
				duration := time.Duration(field.Int())
				*lines = append(*lines, fmt.Sprintf("%s=%s", envKey, duration.String()))
			} else {
				*lines = append(*lines, fmt.Sprintf("%s=%d", envKey, field.Int()))
			}
		case reflect.Bool:
			*lines = append(*lines, fmt.Sprintf("%s=%t", envKey, field.Bool()))
		case reflect.Slice:
			if field.Type().Elem().Kind() == reflect.String {
				values := make([]string, field.Len())
				for j := 0; j < field.Len(); j++ {
					values[j] = field.Index(j).String()
				}
				*lines = append(*lines, fmt.Sprintf("%s=%s", envKey, strings.Join(values, ",")))
			}
		}
	}
}