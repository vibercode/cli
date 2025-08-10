package models

import (
	"time"
)

// PluginType represents the type of plugin
type PluginType string

const (
	PluginTypeGenerator   PluginType = "generator"
	PluginTypeTemplate    PluginType = "template"
	PluginTypeCommand     PluginType = "command"
	PluginTypeIntegration PluginType = "integration"
)

// PluginStatus represents the status of a plugin
type PluginStatus string

const (
	PluginStatusInstalled PluginStatus = "installed"
	PluginStatusEnabled   PluginStatus = "enabled"
	PluginStatusDisabled  PluginStatus = "disabled"
	PluginStatusUpdating  PluginStatus = "updating"
	PluginStatusError     PluginStatus = "error"
)

// PluginOptions represents options for plugin generation
type PluginOptions struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Template    string `json:"template"`
	Author      string `json:"author"`
	Description string `json:"description"`
}

// PluginManifest represents the plugin.yaml file structure
type PluginManifest struct {
	Name         string            `yaml:"name" json:"name"`
	Version      string            `yaml:"version" json:"version"`
	Description  string            `yaml:"description" json:"description"`
	Author       string            `yaml:"author" json:"author"`
	License      string            `yaml:"license" json:"license"`
	Homepage     string            `yaml:"homepage" json:"homepage"`
	Type         PluginType        `yaml:"type" json:"type"`
	Capabilities []string          `yaml:"capabilities" json:"capabilities"`
	Dependencies PluginDependencies `yaml:"dependencies" json:"dependencies"`
	Commands     []PluginCommand   `yaml:"commands" json:"commands"`
	Config       PluginConfig      `yaml:"config" json:"config"`
	Main         string            `yaml:"main" json:"main"`
}

// PluginDependencies represents plugin dependencies
type PluginDependencies struct {
	ViberCode string `yaml:"vibercode" json:"vibercode"`
	Go        string `yaml:"go" json:"go"`
}

// PluginCommand represents a command provided by the plugin
type PluginCommand struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
	Usage       string `yaml:"usage" json:"usage"`
}

// PluginConfig represents plugin configuration schema
type PluginConfig struct {
	Properties map[string]PluginConfigProperty `yaml:"properties" json:"properties"`
}

// PluginConfigProperty represents a single configuration property
type PluginConfigProperty struct {
	Type        string      `yaml:"type" json:"type"`
	Default     interface{} `yaml:"default" json:"default"`
	Description string      `yaml:"description" json:"description"`
	Required    bool        `yaml:"required" json:"required"`
}

// PluginInfo represents information about a plugin
type PluginInfo struct {
	Manifest     PluginManifest `json:"manifest"`
	Status       PluginStatus   `json:"status"`
	InstallPath  string         `json:"install_path"`
	LastUpdated  time.Time      `json:"last_updated"`
	Downloads    int            `json:"downloads"`
	Rating       float64        `json:"rating"`
	Compatibility PluginCompatibility `json:"compatibility"`
}

// PluginCompatibility represents compatibility information
type PluginCompatibility struct {
	ViberCode string `json:"vibercode"`
	Go        string `json:"go"`
}

// PluginRegistry represents plugin registry information
type PluginRegistry struct {
	Plugins []PluginRegistryEntry `json:"plugins"`
}

// PluginRegistryEntry represents a plugin entry in the registry
type PluginRegistryEntry struct {
	Name          string              `json:"name"`
	Version       string              `json:"version"`
	Description   string              `json:"description"`
	Author        string              `json:"author"`
	Category      string              `json:"category"`
	Tags          []string            `json:"tags"`
	DownloadURL   string              `json:"download_url"`
	Homepage      string              `json:"homepage"`
	License       string              `json:"license"`
	Downloads     int                 `json:"downloads"`
	Rating        float64             `json:"rating"`
	LastUpdated   time.Time           `json:"last_updated"`
	Compatibility PluginCompatibility `json:"compatibility"`
}

// SecurityPolicy represents plugin security configuration
type SecurityPolicy struct {
	AllowFileSystemAccess bool          `yaml:"allow_filesystem_access" json:"allow_filesystem_access"`
	AllowNetworkAccess    bool          `yaml:"allow_network_access" json:"allow_network_access"`
	AllowShellExecution   bool          `yaml:"allow_shell_execution" json:"allow_shell_execution"`
	AllowedDirectories    []string      `yaml:"allowed_directories" json:"allowed_directories"`
	BlockedCommands       []string      `yaml:"blocked_commands" json:"blocked_commands"`
	ResourceLimits        ResourceLimits `yaml:"resource_limits" json:"resource_limits"`
}

// ResourceLimits represents resource usage limits for plugins
type ResourceLimits struct {
	MaxMemoryMB      int           `yaml:"max_memory_mb" json:"max_memory_mb"`
	MaxExecutionTime time.Duration `yaml:"max_execution_time" json:"max_execution_time"`
	MaxFileSize      int64         `yaml:"max_file_size" json:"max_file_size"`
}

// PluginContext represents context passed to plugins
type PluginContext struct {
	Config         map[string]interface{} `json:"config"`
	ProjectPath    string                 `json:"project_path"`
	PluginPath     string                 `json:"plugin_path"`
	TempDir        string                 `json:"temp_dir"`
	SecurityPolicy SecurityPolicy         `json:"security_policy"`
}

// PluginExecutionResult represents the result of plugin execution
type PluginExecutionResult struct {
	Success    bool      `json:"success"`
	Message    string    `json:"message"`
	Output     string    `json:"output"`
	Error      string    `json:"error"`
	Duration   time.Duration `json:"duration"`
	ExitCode   int       `json:"exit_code"`
	FilesCreated []string `json:"files_created"`
	FilesModified []string `json:"files_modified"`
}

// PluginValidationResult represents plugin validation results
type PluginValidationResult struct {
	Valid    bool     `json:"valid"`
	Errors   []string `json:"errors"`
	Warnings []string `json:"warnings"`
}

// PluginInstallOptions represents options for plugin installation
type PluginInstallOptions struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Source       string `json:"source"` // registry, file, git
	URL          string `json:"url"`
	Force        bool   `json:"force"`
	SkipSecurity bool   `json:"skip_security"`
}

// PluginSearchQuery represents a plugin search query
type PluginSearchQuery struct {
	Query    string     `json:"query"`
	Category string     `json:"category"`
	Type     PluginType `json:"type"`
	Tags     []string   `json:"tags"`
	Author   string     `json:"author"`
	Limit    int        `json:"limit"`
	Offset   int        `json:"offset"`
}

// PluginSearchResult represents plugin search results
type PluginSearchResult struct {
	Plugins    []PluginRegistryEntry `json:"plugins"`
	TotalCount int                   `json:"total_count"`
	Limit      int                   `json:"limit"`
	Offset     int                   `json:"offset"`
	Query      PluginSearchQuery     `json:"query"`
}

// GeneratorPlugin represents a generator plugin interface
type GeneratorPlugin interface {
	Generate(options map[string]interface{}) (*PluginExecutionResult, error)
	GetTemplates() ([]string, error)
	ValidateOptions(options map[string]interface{}) error
}

// TemplatePlugin represents a template plugin interface
type TemplatePlugin interface {
	GetTemplates() (map[string]string, error)
	RenderTemplate(name string, data interface{}) (string, error)
	ValidateTemplate(name string) error
}

// CommandPlugin represents a command plugin interface
type CommandPlugin interface {
	ExecuteCommand(args []string) (*PluginExecutionResult, error)
	GetCommandHelp() string
	ValidateArgs(args []string) error
}

// IntegrationPlugin represents an integration plugin interface
type IntegrationPlugin interface {
	Initialize(config map[string]interface{}) error
	Connect() error
	Disconnect() error
	IsConnected() bool
	GetStatus() string
}

// PluginConfigManager represents plugin configuration management
type PluginConfigManager struct {
	PluginsDir   string            `json:"plugins_dir"`
	CacheDir     string            `json:"cache_dir"`
	ConfigDir    string            `json:"config_dir"`
	Registries   []string          `json:"registries"`
	Security     SecurityPolicy    `json:"security"`
	GlobalConfig map[string]interface{} `json:"global_config"`
}

// GetDefaultConfig returns default plugin configuration
func GetDefaultConfig() *PluginConfigManager {
	return &PluginConfigManager{
		PluginsDir: "~/.vibercode/plugins/installed",
		CacheDir:   "~/.vibercode/plugins/cache",
		ConfigDir:  "~/.vibercode/plugins/config",
		Registries: []string{
			"https://registry.vibercode.dev",
		},
		Security: SecurityPolicy{
			AllowFileSystemAccess: true,
			AllowNetworkAccess:    false,
			AllowShellExecution:   false,
			AllowedDirectories:    []string{"."},
			BlockedCommands:       []string{"rm", "sudo", "su"},
			ResourceLimits: ResourceLimits{
				MaxMemoryMB:      256,
				MaxExecutionTime: 5 * time.Minute,
				MaxFileSize:      10 * 1024 * 1024, // 10MB
			},
		},
		GlobalConfig: make(map[string]interface{}),
	}
}

// ValidatePluginType validates if the plugin type is supported
func (pt PluginType) IsValid() bool {
	switch pt {
	case PluginTypeGenerator, PluginTypeTemplate, PluginTypeCommand, PluginTypeIntegration:
		return true
	default:
		return false
	}
}

// GetCapabilities returns the capabilities for a plugin type
func (pt PluginType) GetCapabilities() []string {
	switch pt {
	case PluginTypeGenerator:
		return []string{"code-generation", "template-rendering", "file-management"}
	case PluginTypeTemplate:
		return []string{"template-rendering", "file-management"}
	case PluginTypeCommand:
		return []string{"command-execution", "cli-integration"}
	case PluginTypeIntegration:
		return []string{"third-party-integration", "api-communication"}
	default:
		return []string{}
	}
}

// Validate validates the plugin manifest
func (pm *PluginManifest) Validate() *PluginValidationResult {
	result := &PluginValidationResult{
		Valid:    true,
		Errors:   []string{},
		Warnings: []string{},
	}

	// Required fields validation
	if pm.Name == "" {
		result.Valid = false
		result.Errors = append(result.Errors, "Plugin name is required")
	}

	if pm.Version == "" {
		result.Valid = false
		result.Errors = append(result.Errors, "Plugin version is required")
	}

	if pm.Author == "" {
		result.Warnings = append(result.Warnings, "Plugin author is recommended")
	}

	if pm.Description == "" {
		result.Warnings = append(result.Warnings, "Plugin description is recommended")
	}

	// Type validation
	if !pm.Type.IsValid() {
		result.Valid = false
		result.Errors = append(result.Errors, "Invalid plugin type: "+string(pm.Type))
	}

	// Main file validation
	if pm.Main == "" {
		result.Valid = false
		result.Errors = append(result.Errors, "Plugin main entry point is required")
	}

	// Dependencies validation
	if pm.Dependencies.ViberCode == "" {
		result.Warnings = append(result.Warnings, "ViberCode version dependency is recommended")
	}

	if pm.Dependencies.Go == "" {
		result.Warnings = append(result.Warnings, "Go version dependency is recommended")
	}

	return result
}