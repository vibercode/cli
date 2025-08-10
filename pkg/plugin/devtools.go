package plugin

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vibercode/cli/internal/models"
	"gopkg.in/yaml.v2"
)

// DevTools provides development tools for plugin development
type DevTools struct {
	projectDir string
	config     *models.PluginConfigManager
}

// NewDevTools creates a new dev tools instance
func NewDevTools(projectDir string, config *models.PluginConfigManager) *DevTools {
	return &DevTools{
		projectDir: projectDir,
		config:     config,
	}
}

// DevLink creates a development link for a plugin
func (dt *DevTools) DevLink(pluginPath string) error {
	// Validate plugin directory
	if err := dt.validatePluginDirectory(pluginPath); err != nil {
		return fmt.Errorf("invalid plugin directory: %v", err)
	}

	// Load plugin manifest
	manifest, err := dt.loadPluginManifest(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to load plugin manifest: %v", err)
	}

	// Create symlink in development plugins directory
	devPluginsDir := filepath.Join(dt.config.PluginsDir, "dev")
	if err := os.MkdirAll(devPluginsDir, 0755); err != nil {
		return fmt.Errorf("failed to create dev plugins directory: %v", err)
	}

	linkPath := filepath.Join(devPluginsDir, manifest.Name)
	
	// Remove existing link if it exists
	if _, err := os.Lstat(linkPath); err == nil {
		if err := os.Remove(linkPath); err != nil {
			return fmt.Errorf("failed to remove existing link: %v", err)
		}
	}

	// Create symlink
	absPluginPath, err := filepath.Abs(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}

	if err := os.Symlink(absPluginPath, linkPath); err != nil {
		return fmt.Errorf("failed to create symlink: %v", err)
	}

	// Update dev links registry
	if err := dt.updateDevLinksRegistry(manifest.Name, absPluginPath); err != nil {
		return fmt.Errorf("failed to update dev links registry: %v", err)
	}

	fmt.Printf("✅ Plugin %s linked for development\n", manifest.Name)
	fmt.Printf("📁 Linked from: %s\n", absPluginPath)
	fmt.Printf("🔗 Linked to: %s\n", linkPath)

	return nil
}

// DevUnlink removes a development link for a plugin
func (dt *DevTools) DevUnlink(pluginName string) error {
	devPluginsDir := filepath.Join(dt.config.PluginsDir, "dev")
	linkPath := filepath.Join(devPluginsDir, pluginName)

	// Check if link exists
	if _, err := os.Lstat(linkPath); os.IsNotExist(err) {
		return fmt.Errorf("plugin %s is not dev-linked", pluginName)
	}

	// Remove symlink
	if err := os.Remove(linkPath); err != nil {
		return fmt.Errorf("failed to remove dev link: %v", err)
	}

	// Update dev links registry
	if err := dt.removeFromDevLinksRegistry(pluginName); err != nil {
		return fmt.Errorf("failed to update dev links registry: %v", err)
	}

	fmt.Printf("✅ Plugin %s unlinked from development\n", pluginName)

	return nil
}

// Package creates a distributable package for a plugin
func (dt *DevTools) Package(pluginPath, outputPath string) error {
	// Validate plugin directory
	if err := dt.validatePluginDirectory(pluginPath); err != nil {
		return fmt.Errorf("invalid plugin directory: %v", err)
	}

	// Load plugin manifest
	manifest, err := dt.loadPluginManifest(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to load plugin manifest: %v", err)
	}

	// Validate plugin
	if err := dt.validatePlugin(pluginPath, manifest); err != nil {
		return fmt.Errorf("plugin validation failed: %v", err)
	}

	// Determine output file
	if outputPath == "" {
		outputPath = fmt.Sprintf("%s-%s.tar.gz", manifest.Name, manifest.Version)
	}

	// Create package
	if err := dt.createPackage(pluginPath, outputPath, manifest); err != nil {
		return fmt.Errorf("failed to create package: %v", err)
	}

	fmt.Printf("✅ Plugin %s packaged successfully\n", manifest.Name)
	fmt.Printf("📦 Package: %s\n", outputPath)

	// Show package info
	if info, err := os.Stat(outputPath); err == nil {
		fmt.Printf("📊 Size: %.2f MB\n", float64(info.Size())/(1024*1024))
	}

	return nil
}

// Test runs tests for a plugin
func (dt *DevTools) Test(pluginPath string, options TestOptions) (*TestResult, error) {
	// Validate plugin directory
	if err := dt.validatePluginDirectory(pluginPath); err != nil {
		return nil, fmt.Errorf("invalid plugin directory: %v", err)
	}

	result := &TestResult{
		StartTime: time.Now(),
		Tests:     []TestCase{},
	}

	// Run Go tests
	if options.RunGoTests {
		goResult, err := dt.runGoTests(pluginPath, options)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Go tests failed: %v", err))
		} else {
			result.Tests = append(result.Tests, goResult...)
		}
	}

	// Run validation tests
	if options.RunValidation {
		validationResult, err := dt.runValidationTests(pluginPath)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Validation tests failed: %v", err))
		} else {
			result.Tests = append(result.Tests, validationResult...)
		}
	}

	// Run security tests
	if options.RunSecurity {
		securityResult, err := dt.runSecurityTests(pluginPath)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Security tests failed: %v", err))
		} else {
			result.Tests = append(result.Tests, securityResult...)
		}
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	// Calculate summary
	for _, test := range result.Tests {
		switch test.Status {
		case "passed":
			result.Passed++
		case "failed":
			result.Failed++
		case "skipped":
			result.Skipped++
		}
	}

	result.Success = result.Failed == 0 && len(result.Errors) == 0

	return result, nil
}

// Validate validates a plugin
func (dt *DevTools) Validate(pluginPath string) (*ValidationResult, error) {
	// Validate plugin directory
	if err := dt.validatePluginDirectory(pluginPath); err != nil {
		return nil, fmt.Errorf("invalid plugin directory: %v", err)
	}

	// Load plugin manifest
	manifest, err := dt.loadPluginManifest(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load plugin manifest: %v", err)
	}

	result := &ValidationResult{
		Valid:    true,
		Issues:   []ValidationIssue{},
		Warnings: []string{},
	}

	// Validate manifest
	manifestResult := manifest.Validate()
	if !manifestResult.Valid {
		result.Valid = false
		for _, err := range manifestResult.Errors {
			result.Issues = append(result.Issues, ValidationIssue{
				Type:     "error",
				Category: "manifest",
				Message:  err,
			})
		}
	}
	for _, warning := range manifestResult.Warnings {
		result.Warnings = append(result.Warnings, warning)
	}

	// Validate file structure
	if err := dt.validateFileStructure(pluginPath, manifest); err != nil {
		result.Valid = false
		result.Issues = append(result.Issues, ValidationIssue{
			Type:     "error",
			Category: "structure",
			Message:  err.Error(),
		})
	}

	// Validate Go code
	if err := dt.validateGoCode(pluginPath); err != nil {
		result.Issues = append(result.Issues, ValidationIssue{
			Type:     "warning",
			Category: "code",
			Message:  err.Error(),
		})
	}

	// Validate dependencies
	if err := dt.validateDependencies(pluginPath); err != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Dependencies: %v", err))
	}

	return result, nil
}

// GetInfo returns information about a plugin
func (dt *DevTools) GetInfo(pluginPath string) (*PluginInfo, error) {
	// Validate plugin directory
	if err := dt.validatePluginDirectory(pluginPath); err != nil {
		return nil, fmt.Errorf("invalid plugin directory: %v", err)
	}

	// Load plugin manifest
	manifest, err := dt.loadPluginManifest(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load plugin manifest: %v", err)
	}

	info := &PluginInfo{
		Manifest:   *manifest,
		Path:       pluginPath,
		Size:       0,
		FileCount:  0,
		ModTime:    time.Time{},
		IsDevLinked: dt.isDevLinked(manifest.Name),
	}

	// Calculate directory info
	err = filepath.Walk(pluginPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() {
			info.Size += info.Size()
			info.FileCount++
			
			if info.ModTime().After(info.ModTime) {
				info.ModTime = info.ModTime()
			}
		}
		
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to calculate directory info: %v", err)
	}

	return info, nil
}

// validatePluginDirectory validates that a directory is a valid plugin directory
func (dt *DevTools) validatePluginDirectory(pluginPath string) error {
	// Check if directory exists
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", pluginPath)
	}

	// Check for plugin.yaml
	manifestPath := filepath.Join(pluginPath, "plugin.yaml")
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		return fmt.Errorf("plugin.yaml not found")
	}

	return nil
}

// loadPluginManifest loads and parses the plugin manifest
func (dt *DevTools) loadPluginManifest(pluginPath string) (*models.PluginManifest, error) {
	manifestPath := filepath.Join(pluginPath, "plugin.yaml")
	
	data, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin.yaml: %v", err)
	}

	var manifest models.PluginManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse plugin.yaml: %v", err)
	}

	return &manifest, nil
}

// validatePlugin validates a plugin
func (dt *DevTools) validatePlugin(pluginPath string, manifest *models.PluginManifest) error {
	// Validate manifest
	result := manifest.Validate()
	if !result.Valid {
		return fmt.Errorf("manifest validation failed: %v", result.Errors)
	}

	// Check main file exists
	mainPath := filepath.Join(pluginPath, manifest.Main)
	if _, err := os.Stat(mainPath); os.IsNotExist(err) {
		return fmt.Errorf("main file not found: %s", manifest.Main)
	}

	return nil
}

// createPackage creates a tar.gz package of the plugin
func (dt *DevTools) createPackage(pluginPath, outputPath string, manifest *models.PluginManifest) error {
	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	// Create gzip writer
	gzWriter := gzip.NewWriter(outFile)
	defer gzWriter.Close()

	// Create tar writer
	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	// Add files to package
	err = filepath.Walk(pluginPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip certain files and directories
		if dt.shouldSkipFile(path, pluginPath) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Create tar header
		relPath, err := filepath.Rel(pluginPath, path)
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = relPath

		// Write header
		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		// Write file content if it's a regular file
		if info.Mode().IsRegular() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(tarWriter, file)
			return err
		}

		return nil
	})

	return err
}

// shouldSkipFile determines if a file should be skipped during packaging
func (dt *DevTools) shouldSkipFile(path, basePath string) bool {
	relPath, _ := filepath.Rel(basePath, path)
	
	skipPatterns := []string{
		".git",
		".gitignore",
		"node_modules",
		"vendor",
		".vscode",
		".idea",
		"*.tmp",
		"*.log",
		".DS_Store",
		"Thumbs.db",
		"build",
		"dist",
		"coverage.out",
		"coverage.html",
	}

	for _, pattern := range skipPatterns {
		if matched, _ := filepath.Match(pattern, relPath); matched {
			return true
		}
		if strings.Contains(relPath, pattern) {
			return true
		}
	}

	return false
}

// updateDevLinksRegistry updates the development links registry
func (dt *DevTools) updateDevLinksRegistry(pluginName, pluginPath string) error {
	registryPath := filepath.Join(dt.config.ConfigDir, "dev-links.json")
	
	// Load existing registry
	var registry map[string]string
	if data, err := ioutil.ReadFile(registryPath); err == nil {
		json.Unmarshal(data, &registry)
	}
	if registry == nil {
		registry = make(map[string]string)
	}

	// Update registry
	registry[pluginName] = pluginPath

	// Save registry
	data, _ := json.MarshalIndent(registry, "", "  ")
	
	if err := os.MkdirAll(filepath.Dir(registryPath), 0755); err != nil {
		return err
	}
	
	return ioutil.WriteFile(registryPath, data, 0644)
}

// removeFromDevLinksRegistry removes a plugin from the development links registry
func (dt *DevTools) removeFromDevLinksRegistry(pluginName string) error {
	registryPath := filepath.Join(dt.config.ConfigDir, "dev-links.json")
	
	// Load existing registry
	var registry map[string]string
	if data, err := ioutil.ReadFile(registryPath); err == nil {
		json.Unmarshal(data, &registry)
	}
	if registry == nil {
		return nil
	}

	// Remove from registry
	delete(registry, pluginName)

	// Save registry
	data, _ := json.MarshalIndent(registry, "", "  ")
	return ioutil.WriteFile(registryPath, data, 0644)
}

// isDevLinked checks if a plugin is development linked
func (dt *DevTools) isDevLinked(pluginName string) bool {
	registryPath := filepath.Join(dt.config.ConfigDir, "dev-links.json")
	
	var registry map[string]string
	if data, err := ioutil.ReadFile(registryPath); err == nil {
		json.Unmarshal(data, &registry)
	}
	
	_, exists := registry[pluginName]
	return exists
}

// validateFileStructure validates the plugin file structure
func (dt *DevTools) validateFileStructure(pluginPath string, manifest *models.PluginManifest) error {
	requiredFiles := []string{
		"plugin.yaml",
		manifest.Main,
	}

	for _, file := range requiredFiles {
		filePath := filepath.Join(pluginPath, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fmt.Errorf("required file missing: %s", file)
		}
	}

	return nil
}

// validateGoCode validates Go code in the plugin
func (dt *DevTools) validateGoCode(pluginPath string) error {
	// Check if go.mod exists and is valid
	goModPath := filepath.Join(pluginPath, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		// Try to parse go.mod by running go mod verify
		// This is a simplified check
	}

	return nil
}

// validateDependencies validates plugin dependencies
func (dt *DevTools) validateDependencies(pluginPath string) error {
	// TODO: Implement dependency validation
	// - Check for known vulnerable dependencies
	// - Validate version constraints
	// - Check for circular dependencies
	return nil
}

// runGoTests runs Go tests for the plugin
func (dt *DevTools) runGoTests(pluginPath string, options TestOptions) ([]TestCase, error) {
	// TODO: Implement Go test execution
	return []TestCase{}, nil
}

// runValidationTests runs validation tests
func (dt *DevTools) runValidationTests(pluginPath string) ([]TestCase, error) {
	// TODO: Implement validation tests
	return []TestCase{}, nil
}

// runSecurityTests runs security tests
func (dt *DevTools) runSecurityTests(pluginPath string) ([]TestCase, error) {
	// TODO: Implement security tests
	return []TestCase{}, nil
}

// Supporting types

// TestOptions configures test execution
type TestOptions struct {
	RunGoTests    bool
	RunValidation bool
	RunSecurity   bool
	Verbose       bool
	Coverage      bool
}

// TestResult represents test execution results
type TestResult struct {
	StartTime time.Time   `json:"start_time"`
	EndTime   time.Time   `json:"end_time"`
	Duration  time.Duration `json:"duration"`
	Success   bool        `json:"success"`
	Passed    int         `json:"passed"`
	Failed    int         `json:"failed"`
	Skipped   int         `json:"skipped"`
	Tests     []TestCase  `json:"tests"`
	Errors    []string    `json:"errors"`
}

// TestCase represents a single test case
type TestCase struct {
	Name     string        `json:"name"`
	Status   string        `json:"status"` // passed, failed, skipped
	Duration time.Duration `json:"duration"`
	Error    string        `json:"error"`
	Output   string        `json:"output"`
}

// ValidationResult represents plugin validation results
type ValidationResult struct {
	Valid    bool               `json:"valid"`
	Issues   []ValidationIssue  `json:"issues"`
	Warnings []string           `json:"warnings"`
}

// ValidationIssue represents a validation issue
type ValidationIssue struct {
	Type     string `json:"type"`     // error, warning
	Category string `json:"category"` // manifest, structure, code, etc.
	Message  string `json:"message"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

// PluginInfo represents plugin information
type PluginInfo struct {
	Manifest    models.PluginManifest `json:"manifest"`
	Path        string                `json:"path"`
	Size        int64                 `json:"size"`
	FileCount   int                   `json:"file_count"`
	ModTime     time.Time             `json:"mod_time"`
	IsDevLinked bool                  `json:"is_dev_linked"`
}