package plugin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vibercode/cli/internal/models"
)

// DefaultRegistryURL is the default plugin registry URL
const DefaultRegistryURL = "https://registry.vibercode.dev"

// Registry represents a plugin registry
type Registry struct {
	URL         string    `json:"url"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	LastSync    time.Time `json:"last_sync"`
	Enabled     bool      `json:"enabled"`
}

// RegistryManager manages plugin registries
type RegistryManager struct {
	registries []Registry
	cacheDir   string
	client     *http.Client
}

// NewRegistryManager creates a new registry manager
func NewRegistryManager(cacheDir string) *RegistryManager {
	return &RegistryManager{
		registries: []Registry{
			{
				URL:         DefaultRegistryURL,
				Name:        "Official ViberCode Registry",
				Description: "Official plugin registry for ViberCode CLI",
				Enabled:     true,
			},
		},
		cacheDir: cacheDir,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// AddRegistry adds a new registry
func (rm *RegistryManager) AddRegistry(url, name, description string) error {
	// Validate URL
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return fmt.Errorf("invalid registry URL: %s", url)
	}

	// Check if registry already exists
	for _, reg := range rm.registries {
		if reg.URL == url {
			return fmt.Errorf("registry %s already exists", url)
		}
	}

	registry := Registry{
		URL:         url,
		Name:        name,
		Description: description,
		Enabled:     true,
	}

	// Test registry connectivity
	if err := rm.testRegistry(registry); err != nil {
		return fmt.Errorf("failed to connect to registry %s: %v", url, err)
	}

	rm.registries = append(rm.registries, registry)
	return rm.saveRegistriesConfig()
}

// RemoveRegistry removes a registry
func (rm *RegistryManager) RemoveRegistry(url string) error {
	for i, reg := range rm.registries {
		if reg.URL == url {
			rm.registries = append(rm.registries[:i], rm.registries[i+1:]...)
			return rm.saveRegistriesConfig()
		}
	}
	return fmt.Errorf("registry %s not found", url)
}

// ListRegistries returns all registered registries
func (rm *RegistryManager) ListRegistries() []Registry {
	return rm.registries
}

// EnableRegistry enables a registry
func (rm *RegistryManager) EnableRegistry(url string) error {
	for i, reg := range rm.registries {
		if reg.URL == url {
			rm.registries[i].Enabled = true
			return rm.saveRegistriesConfig()
		}
	}
	return fmt.Errorf("registry %s not found", url)
}

// DisableRegistry disables a registry
func (rm *RegistryManager) DisableRegistry(url string) error {
	for i, reg := range rm.registries {
		if reg.URL == url {
			rm.registries[i].Enabled = false
			return rm.saveRegistriesConfig()
		}
	}
	return fmt.Errorf("registry %s not found", url)
}

// RefreshAll refreshes all enabled registries
func (rm *RegistryManager) RefreshAll() error {
	var errors []string

	for i, reg := range rm.registries {
		if !reg.Enabled {
			continue
		}

		if err := rm.refreshRegistry(&rm.registries[i]); err != nil {
			errors = append(errors, fmt.Sprintf("failed to refresh %s: %v", reg.URL, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("refresh errors: %s", strings.Join(errors, "; "))
	}

	return rm.saveRegistriesConfig()
}

// SearchPlugins searches for plugins across all enabled registries
func (rm *RegistryManager) SearchPlugins(query models.PluginSearchQuery) (*models.PluginSearchResult, error) {
	var allPlugins []models.PluginRegistryEntry
	
	for _, reg := range rm.registries {
		if !reg.Enabled {
			continue
		}

		plugins, err := rm.searchInRegistry(reg, query)
		if err != nil {
			// Log error but continue with other registries
			continue
		}

		allPlugins = append(allPlugins, plugins...)
	}

	// Filter and sort results
	filteredPlugins := rm.filterPlugins(allPlugins, query)
	sortedPlugins := rm.sortPlugins(filteredPlugins, query)

	// Apply pagination
	start := query.Offset
	end := start + query.Limit
	if start > len(sortedPlugins) {
		start = len(sortedPlugins)
	}
	if end > len(sortedPlugins) {
		end = len(sortedPlugins)
	}

	result := &models.PluginSearchResult{
		Plugins:    sortedPlugins[start:end],
		TotalCount: len(sortedPlugins),
		Limit:      query.Limit,
		Offset:     query.Offset,
		Query:      query,
	}

	return result, nil
}

// GetPlugin gets detailed information about a specific plugin
func (rm *RegistryManager) GetPlugin(name string) (*models.PluginRegistryEntry, error) {
	for _, reg := range rm.registries {
		if !reg.Enabled {
			continue
		}

		plugin, err := rm.getPluginFromRegistry(reg, name)
		if err == nil {
			return plugin, nil
		}
	}

	return nil, fmt.Errorf("plugin %s not found in any registry", name)
}

// DownloadPlugin downloads a plugin from a registry
func (rm *RegistryManager) DownloadPlugin(name, version string) (string, error) {
	plugin, err := rm.GetPlugin(name)
	if err != nil {
		return "", err
	}

	// Create download directory
	downloadDir := filepath.Join(rm.cacheDir, "downloads")
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create download directory: %v", err)
	}

	// Download plugin
	downloadURL := plugin.DownloadURL
	if version != "" && version != plugin.Version {
		// Construct version-specific URL
		downloadURL = strings.Replace(downloadURL, plugin.Version, version, 1)
	}

	resp, err := rm.client.Get(downloadURL)
	if err != nil {
		return "", fmt.Errorf("failed to download plugin: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// Save to file
	filename := fmt.Sprintf("%s-%s.tar.gz", name, plugin.Version)
	if version != "" {
		filename = fmt.Sprintf("%s-%s.tar.gz", name, version)
	}
	
	filePath := filepath.Join(downloadDir, filename)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read download data: %v", err)
	}

	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to save plugin file: %v", err)
	}

	return filePath, nil
}

// testRegistry tests connectivity to a registry
func (rm *RegistryManager) testRegistry(reg Registry) error {
	url := reg.URL + "/health"
	resp, err := rm.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("registry health check failed with status: %d", resp.StatusCode)
	}

	return nil
}

// refreshRegistry refreshes a single registry
func (rm *RegistryManager) refreshRegistry(reg *Registry) error {
	url := reg.URL + "/plugins"
	resp, err := rm.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch plugins: status %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	// Save to cache
	cacheFile := filepath.Join(rm.cacheDir, rm.getRegistryCacheFilename(reg.URL))
	if err := os.MkdirAll(filepath.Dir(cacheFile), 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %v", err)
	}

	if err := ioutil.WriteFile(cacheFile, data, 0644); err != nil {
		return fmt.Errorf("failed to save cache: %v", err)
	}

	reg.LastSync = time.Now()
	return nil
}

// searchInRegistry searches for plugins in a specific registry
func (rm *RegistryManager) searchInRegistry(reg Registry, query models.PluginSearchQuery) ([]models.PluginRegistryEntry, error) {
	// Load from cache first
	cacheFile := filepath.Join(rm.cacheDir, rm.getRegistryCacheFilename(reg.URL))
	data, err := ioutil.ReadFile(cacheFile)
	if err != nil {
		// Try to refresh if cache doesn't exist
		if os.IsNotExist(err) {
			if refreshErr := rm.refreshRegistry(&reg); refreshErr != nil {
				return nil, refreshErr
			}
			data, err = ioutil.ReadFile(cacheFile)
			if err != nil {
				return nil, fmt.Errorf("failed to read cache after refresh: %v", err)
			}
		} else {
			return nil, fmt.Errorf("failed to read cache: %v", err)
		}
	}

	var registry models.PluginRegistry
	if err := json.Unmarshal(data, &registry); err != nil {
		return nil, fmt.Errorf("failed to parse registry data: %v", err)
	}

	return registry.Plugins, nil
}

// getPluginFromRegistry gets a specific plugin from a registry
func (rm *RegistryManager) getPluginFromRegistry(reg Registry, name string) (*models.PluginRegistryEntry, error) {
	url := fmt.Sprintf("%s/plugins/%s", reg.URL, name)
	resp, err := rm.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("plugin not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch plugin: status %d", resp.StatusCode)
	}

	var plugin models.PluginRegistryEntry
	if err := json.NewDecoder(resp.Body).Decode(&plugin); err != nil {
		return nil, fmt.Errorf("failed to decode plugin data: %v", err)
	}

	return &plugin, nil
}

// filterPlugins filters plugins based on search criteria
func (rm *RegistryManager) filterPlugins(plugins []models.PluginRegistryEntry, query models.PluginSearchQuery) []models.PluginRegistryEntry {
	var filtered []models.PluginRegistryEntry

	for _, plugin := range plugins {
		// Text search
		if query.Query != "" {
			searchText := strings.ToLower(query.Query)
			pluginText := strings.ToLower(plugin.Name + " " + plugin.Description)
			if !strings.Contains(pluginText, searchText) {
				continue
			}
		}

		// Category filter
		if query.Category != "" && plugin.Category != query.Category {
			continue
		}

		// Type filter
		if query.Type != "" && string(query.Type) != plugin.Category {
			continue
		}

		// Author filter
		if query.Author != "" && plugin.Author != query.Author {
			continue
		}

		// Tags filter
		if len(query.Tags) > 0 {
			hasTag := false
			for _, queryTag := range query.Tags {
				for _, pluginTag := range plugin.Tags {
					if strings.EqualFold(queryTag, pluginTag) {
						hasTag = true
						break
					}
				}
				if hasTag {
					break
				}
			}
			if !hasTag {
				continue
			}
		}

		filtered = append(filtered, plugin)
	}

	return filtered
}

// sortPlugins sorts plugins by relevance
func (rm *RegistryManager) sortPlugins(plugins []models.PluginRegistryEntry, query models.PluginSearchQuery) []models.PluginRegistryEntry {
	// Simple sorting by downloads and rating for now
	// TODO: Implement more sophisticated relevance scoring
	sorted := make([]models.PluginRegistryEntry, len(plugins))
	copy(sorted, plugins)

	// Sort by downloads (descending) then by rating (descending)
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i].Downloads < sorted[j].Downloads ||
				(sorted[i].Downloads == sorted[j].Downloads && sorted[i].Rating < sorted[j].Rating) {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	return sorted
}

// getRegistryCacheFilename generates a cache filename for a registry
func (rm *RegistryManager) getRegistryCacheFilename(url string) string {
	// Simple URL to filename conversion
	filename := strings.ReplaceAll(url, "://", "_")
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, ".", "_")
	return filename + ".json"
}

// saveRegistriesConfig saves the registries configuration
func (rm *RegistryManager) saveRegistriesConfig() error {
	configFile := filepath.Join(rm.cacheDir, "registries.json")
	data, err := json.MarshalIndent(rm.registries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal registries config: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(configFile), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	return ioutil.WriteFile(configFile, data, 0644)
}

// loadRegistriesConfig loads the registries configuration
func (rm *RegistryManager) loadRegistriesConfig() error {
	configFile := filepath.Join(rm.cacheDir, "registries.json")
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			// Use default registries if config doesn't exist
			return nil
		}
		return fmt.Errorf("failed to read registries config: %v", err)
	}

	return json.Unmarshal(data, &rm.registries)
}

// PluginDiscovery handles plugin discovery across multiple sources
type PluginDiscovery struct {
	registryManager *RegistryManager
	localPlugins    []models.PluginInfo
}

// NewPluginDiscovery creates a new plugin discovery service
func NewPluginDiscovery(cacheDir string) *PluginDiscovery {
	return &PluginDiscovery{
		registryManager: NewRegistryManager(cacheDir),
		localPlugins:    []models.PluginInfo{},
	}
}

// DiscoverPlugins discovers plugins from all sources
func (pd *PluginDiscovery) DiscoverPlugins(query models.PluginSearchQuery) (*models.PluginSearchResult, error) {
	// Search in registries
	registryResult, err := pd.registryManager.SearchPlugins(query)
	if err != nil {
		return nil, fmt.Errorf("failed to search registries: %v", err)
	}

	// TODO: Add local plugin discovery
	// TODO: Add Git repository discovery
	// TODO: Merge results from all sources

	return registryResult, nil
}

// GetRegistryManager returns the registry manager
func (pd *PluginDiscovery) GetRegistryManager() *RegistryManager {
	return pd.registryManager
}

// AddLocalPlugin adds a local plugin to the discovery service
func (pd *PluginDiscovery) AddLocalPlugin(plugin models.PluginInfo) {
	pd.localPlugins = append(pd.localPlugins, plugin)
}

// RemoveLocalPlugin removes a local plugin from the discovery service
func (pd *PluginDiscovery) RemoveLocalPlugin(name string) {
	for i, plugin := range pd.localPlugins {
		if plugin.Manifest.Name == name {
			pd.localPlugins = append(pd.localPlugins[:i], pd.localPlugins[i+1:]...)
			break
		}
	}
}

// GetLocalPlugins returns all local plugins
func (pd *PluginDiscovery) GetLocalPlugins() []models.PluginInfo {
	return pd.localPlugins
}