package plugin

import (
	"context"
	"time"

	"github.com/vibercode/cli/internal/models"
)

// Plugin represents the core plugin interface that all plugins must implement
type Plugin interface {
	// Metadata methods
	Name() string
	Version() string
	Description() string
	Author() string

	// Lifecycle methods
	Initialize(ctx PluginContext) error
	Execute(args []string) error
	Cleanup() error

	// Capability methods
	Commands() []Command
	Generators() []Generator
	Templates() []Template
}

// PluginContext provides access to CLI internals and services
type PluginContext interface {
	Config() map[string]interface{}
	Logger() Logger
	FileSystem() FileSystem
	TemplateEngine() TemplateEngine
	UserInterface() UserInterface
	ProjectPath() string
	PluginPath() string
	TempDir() string
	SecurityPolicy() models.SecurityPolicy
}

// Logger provides logging capabilities to plugins
type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warning(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
}

// FileSystem provides file system operations
type FileSystem interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte) error
	CreateDir(path string) error
	RemoveFile(path string) error
	RemoveDir(path string) error
	Exists(path string) bool
	ListFiles(dir string) ([]string, error)
	CopyFile(src, dst string) error
	MoveFile(src, dst string) error
}

// TemplateEngine provides template rendering capabilities
type TemplateEngine interface {
	Render(templateContent string, data interface{}) (string, error)
	RenderFile(templatePath string, data interface{}) (string, error)
	RegisterHelper(name string, helper interface{}) error
	LoadTemplate(path string) (Template, error)
}

// UserInterface provides user interaction capabilities
type UserInterface interface {
	Print(message string)
	Printf(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warning(format string, args ...interface{})
	Error(format string, args ...interface{})
	Success(format string, args ...interface{})
	Prompt(message string) (string, error)
	PromptConfirm(message string) (bool, error)
	PromptSelect(message string, options []string) (string, error)
}

// Command represents a CLI command provided by a plugin
type Command struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Usage       string   `json:"usage"`
	Aliases     []string `json:"aliases"`
	Flags       []Flag   `json:"flags"`
}

// Flag represents a command line flag
type Flag struct {
	Name        string      `json:"name"`
	ShortName   string      `json:"short_name"`
	Description string      `json:"description"`
	Type        string      `json:"type"`
	Required    bool        `json:"required"`
	Default     interface{} `json:"default"`
}

// Generator represents a code generator provided by a plugin
type Generator struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Extensions  []string `json:"extensions"`
	Templates   []string `json:"templates"`
}

// Template represents a template provided by a plugin
type Template struct {
	Name        string            `json:"name"`
	Path        string            `json:"path"`
	Description string            `json:"description"`
	Variables   map[string]string `json:"variables"`
	Helpers     []string          `json:"helpers"`
}

// ExecutionResult represents the result of plugin execution
type ExecutionResult struct {
	Success       bool          `json:"success"`
	Message       string        `json:"message"`
	Output        string        `json:"output"`
	Error         string        `json:"error"`
	Duration      time.Duration `json:"duration"`
	ExitCode      int           `json:"exit_code"`
	FilesCreated  []string      `json:"files_created"`
	FilesModified []string      `json:"files_modified"`
}

// GeneratorPlugin interface for generator-type plugins
type GeneratorPlugin interface {
	Plugin
	Generate(options map[string]interface{}) (*ExecutionResult, error)
	GetTemplates() ([]string, error)
	ValidateOptions(options map[string]interface{}) error
	GetSchema() (map[string]interface{}, error)
}

// TemplatePlugin interface for template-type plugins
type TemplatePlugin interface {
	Plugin
	GetTemplates() (map[string]string, error)
	RenderTemplate(name string, data interface{}) (string, error)
	ValidateTemplate(name string) error
	GetTemplateVars(name string) ([]string, error)
}

// CommandPlugin interface for command-type plugins
type CommandPlugin interface {
	Plugin
	ExecuteCommand(command string, args []string) (*ExecutionResult, error)
	GetCommandHelp(command string) string
	ValidateArgs(command string, args []string) error
	GetCommands() []Command
}

// IntegrationPlugin interface for integration-type plugins
type IntegrationPlugin interface {
	Plugin
	Initialize(config map[string]interface{}) error
	Connect() error
	Disconnect() error
	IsConnected() bool
	GetStatus() string
	Sync() error
	GetData(query map[string]interface{}) (interface{}, error)
	SendData(data interface{}) error
}

// PluginFactory creates plugin instances
type PluginFactory interface {
	CreatePlugin(pluginType string, config map[string]interface{}) (Plugin, error)
	GetSupportedTypes() []string
}

// PluginRegistry manages plugin registration and discovery
type PluginRegistry interface {
	Register(plugin Plugin) error
	Unregister(name string) error
	Get(name string) (Plugin, error)
	List() []Plugin
	Search(query string) []Plugin
}

// PluginValidator validates plugin implementations
type PluginValidator interface {
	Validate(plugin Plugin) (*models.PluginValidationResult, error)
	ValidateManifest(manifest *models.PluginManifest) (*models.PluginValidationResult, error)
	ValidateSecurity(plugin Plugin, policy models.SecurityPolicy) error
}

// PluginExecutor executes plugins with security constraints
type PluginExecutor interface {
	Execute(ctx context.Context, plugin Plugin, args []string) (*ExecutionResult, error)
	ExecuteWithTimeout(ctx context.Context, plugin Plugin, args []string, timeout time.Duration) (*ExecutionResult, error)
	SetSecurityPolicy(policy models.SecurityPolicy)
}

// PluginLoader loads plugins from various sources
type PluginLoader interface {
	LoadFromPath(path string) (Plugin, error)
	LoadFromManifest(manifest *models.PluginManifest) (Plugin, error)
	LoadFromBinary(binaryPath string) (Plugin, error)
	Unload(plugin Plugin) error
}

// Event represents a plugin system event
type Event struct {
	Type      string                 `json:"type"`
	Plugin    string                 `json:"plugin"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// EventHandler handles plugin system events
type EventHandler interface {
	OnPluginInstalled(plugin Plugin) error
	OnPluginUninstalled(name string) error
	OnPluginEnabled(plugin Plugin) error
	OnPluginDisabled(plugin Plugin) error
	OnPluginExecuted(plugin Plugin, result *ExecutionResult) error
	OnPluginError(plugin Plugin, err error) error
}

// PluginManager orchestrates all plugin operations
type PluginManager interface {
	// Plugin lifecycle
	Install(source string, options models.PluginInstallOptions) error
	Uninstall(name string) error
	Enable(name string) error
	Disable(name string) error
	Update(name string) error

	// Plugin information
	List() ([]models.PluginInfo, error)
	Get(name string) (*models.PluginInfo, error)
	Search(query models.PluginSearchQuery) (*models.PluginSearchResult, error)

	// Plugin execution
	Execute(name string, args []string) (*ExecutionResult, error)
	ExecuteGenerator(name string, options map[string]interface{}) (*ExecutionResult, error)
	ExecuteCommand(pluginName string, command string, args []string) (*ExecutionResult, error)

	// Configuration
	GetConfig() *models.PluginConfigManager
	SetConfig(config *models.PluginConfigManager) error
	ValidateConfig(config *models.PluginConfigManager) error

	// Registry operations
	AddRegistry(url string) error
	RemoveRegistry(url string) error
	RefreshRegistry() error

	// Security
	SetSecurityPolicy(policy models.SecurityPolicy) error
	ValidatePlugin(plugin Plugin) (*models.PluginValidationResult, error)

	// Events
	RegisterEventHandler(handler EventHandler) error
	UnregisterEventHandler(handler EventHandler) error
	EmitEvent(event Event) error

	// Development tools
	DevLink(path string) error
	DevUnlink(name string) error
	Package(path string, output string) error
	Publish(packagePath string, registry string) error
}