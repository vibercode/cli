package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/vibercode/cli/internal/models"
	"github.com/vibercode/cli/internal/templates"
	"github.com/vibercode/cli/pkg/ui"
)

// DeploymentOptions contains configuration for deployment generation
type DeploymentOptions struct {
	Type         string
	Provider     string
	Service      string
	Namespace    string
	MultiStage   bool
	Optimize     bool
	Security     bool
	WithIngress  bool
	WithSecrets  bool
	WithHPA      bool
	FullSuite    bool
	Environment  string
}

// DeploymentGenerator handles deployment generation
type DeploymentGenerator struct {
	options DeploymentOptions
}

// NewDeploymentGenerator creates a new deployment generator
func NewDeploymentGenerator() *DeploymentGenerator {
	return &DeploymentGenerator{}
}

// Generate generates deployment configurations based on options
func (g *DeploymentGenerator) Generate(options DeploymentOptions) error {
	g.options = options

	ui.PrintStep(1, 1, "Starting deployment generation...")

	// Handle full suite generation
	if options.FullSuite {
		return g.generateFullDeploymentSuite()
	}

	// Handle specific deployment type
	if options.Type != "" {
		return g.generateSpecificDeployment()
	}

	// Interactive mode
	return g.generateInteractiveDeployment()
}

// generateFullDeploymentSuite generates complete deployment suite
func (g *DeploymentGenerator) generateFullDeploymentSuite() error {
	ui.PrintStep(1, 6, "Generating full deployment suite...")

	// Get deployment configuration
	config, err := g.getFullSuiteConfig()
	if err != nil {
		return fmt.Errorf("failed to get deployment config: %w", err)
	}

	// Create deployment directory structure
	ui.PrintStep(2, 6, "Creating deployment directories...")
	if err := g.createDeploymentDirectories(); err != nil {
		return fmt.Errorf("failed to create deployment directories: %w", err)
	}

	// Generate Docker configurations
	ui.PrintStep(3, 6, "Generating Docker configurations...")
	if err := g.generateDockerFiles(config); err != nil {
		return fmt.Errorf("failed to generate Docker files: %w", err)
	}

	// Generate Kubernetes manifests
	ui.PrintStep(4, 6, "Generating Kubernetes manifests...")
	if err := g.generateKubernetesFiles(config); err != nil {
		return fmt.Errorf("failed to generate Kubernetes files: %w", err)
	}

	// Generate cloud deployment files
	ui.PrintStep(5, 6, "Generating cloud deployment files...")
	if err := g.generateCloudFiles(config); err != nil {
		return fmt.Errorf("failed to generate cloud files: %w", err)
	}

	// Generate CI/CD pipelines
	ui.PrintStep(6, 6, "Generating CI/CD pipelines...")
	if err := g.generateCICDFiles(config); err != nil {
		return fmt.Errorf("failed to generate CI/CD files: %w", err)
	}

	ui.PrintSuccess("Full deployment suite generated successfully!")
	g.showDeploymentSummary(config)

	return nil
}

// generateSpecificDeployment generates specific deployment type
func (g *DeploymentGenerator) generateSpecificDeployment() error {
	deploymentType := models.DeploymentType(g.options.Type)
	
	ui.PrintStep(1, 2, fmt.Sprintf("Generating %s deployment...", deploymentType))

	config, err := g.getSpecificDeploymentConfig()
	if err != nil {
		return fmt.Errorf("failed to get deployment config: %w", err)
	}

	// Create directories
	if err := g.createDeploymentDirectories(); err != nil {
		return fmt.Errorf("failed to create deployment directories: %w", err)
	}

	ui.PrintStep(2, 2, "Creating deployment files...")

	switch deploymentType {
	case models.DockerDeployment:
		return g.generateDockerFiles(config)
	case models.KubernetesDeployment:
		return g.generateKubernetesFiles(config)
	case models.CloudDeployment:
		return g.generateCloudFiles(config)
	case models.CICDDeployment:
		return g.generateCICDFiles(config)
	}

	return fmt.Errorf("unsupported deployment type: %s", deploymentType)
}

// generateInteractiveDeployment handles interactive deployment generation
func (g *DeploymentGenerator) generateInteractiveDeployment() error {
	ui.PrintStep(1, 1, "Interactive deployment generation...")

	// Ask what to generate
	actionPrompt := promptui.Select{
		Label: ui.IconDocker + " What would you like to generate?",
		Items: []string{
			"Full deployment suite",
			"Docker configurations",
			"Kubernetes manifests",
			"Cloud deployment",
			"CI/CD pipelines",
		},
	}

	_, action, err := actionPrompt.Run()
	if err != nil {
		return err
	}

	switch action {
	case "Full deployment suite":
		g.options.FullSuite = true
		return g.generateFullDeploymentSuite()
	case "Docker configurations":
		g.options.Type = string(models.DockerDeployment)
		return g.generateInteractiveDocker()
	case "Kubernetes manifests":
		g.options.Type = string(models.KubernetesDeployment)
		return g.generateInteractiveKubernetes()
	case "Cloud deployment":
		g.options.Type = string(models.CloudDeployment)
		return g.generateInteractiveCloud()
	case "CI/CD pipelines":
		g.options.Type = string(models.CICDDeployment)
		return g.generateInteractiveCICD()
	}

	return nil
}

// generateInteractiveDocker handles interactive Docker generation
func (g *DeploymentGenerator) generateInteractiveDocker() error {
	// Multi-stage build
	multiStagePrompt := promptui.Select{
		Label: ui.IconBuild + " Use multi-stage build?",
		Items: []string{"Yes", "No"},
	}
	_, multiStageAnswer, _ := multiStagePrompt.Run()
	g.options.MultiStage = multiStageAnswer == "Yes"

	// Security hardening
	securityPrompt := promptui.Select{
		Label: ui.IconGear + " Enable security hardening?",
		Items: []string{"Yes", "No"},
	}
	_, securityAnswer, _ := securityPrompt.Run()
	g.options.Security = securityAnswer == "Yes"

	// Optimization
	optimizePrompt := promptui.Select{
		Label: ui.IconSpeed + " Enable image optimization?",
		Items: []string{"Yes", "No"},
	}
	_, optimizeAnswer, _ := optimizePrompt.Run()
	g.options.Optimize = optimizeAnswer == "Yes"

	return g.generateSpecificDeployment()
}

// generateInteractiveKubernetes handles interactive Kubernetes generation
func (g *DeploymentGenerator) generateInteractiveKubernetes() error {
	// Namespace
	namespacePrompt := promptui.Prompt{
		Label:   ui.IconPackage + " Kubernetes namespace",
		Default: "default",
	}
	g.options.Namespace, _ = namespacePrompt.Run()

	// Ingress
	ingressPrompt := promptui.Select{
		Label: ui.IconAPI + " Include Ingress configuration?",
		Items: []string{"Yes", "No"},
	}
	_, ingressAnswer, _ := ingressPrompt.Run()
	g.options.WithIngress = ingressAnswer == "Yes"

	// Secrets
	secretsPrompt := promptui.Select{
		Label: ui.IconGear + " Include Secrets and ConfigMaps?",
		Items: []string{"Yes", "No"},
	}
	_, secretsAnswer, _ := secretsPrompt.Run()
	g.options.WithSecrets = secretsAnswer == "Yes"

	// HPA
	hpaPrompt := promptui.Select{
		Label: ui.IconSpeed + " Include Horizontal Pod Autoscaler?",
		Items: []string{"Yes", "No"},
	}
	_, hpaAnswer, _ := hpaPrompt.Run()
	g.options.WithHPA = hpaAnswer == "Yes"

	return g.generateSpecificDeployment()
}

// generateInteractiveCloud handles interactive cloud deployment generation
func (g *DeploymentGenerator) generateInteractiveCloud() error {
	// Cloud provider
	providerPrompt := promptui.Select{
		Label: ui.IconDocker + " Select cloud provider",
		Items: []string{"aws", "gcp", "azure"},
	}
	_, g.options.Provider, _ = providerPrompt.Run()

	// Cloud service based on provider
	var services []string
	switch g.options.Provider {
	case "aws":
		services = []string{"ecs", "fargate", "eks", "lambda"}
	case "gcp":
		services = []string{"run", "gke", "appengine"}
	case "azure":
		services = []string{"containers", "aks", "webapp"}
	}

	servicePrompt := promptui.Select{
		Label: ui.IconAPI + " Select cloud service",
		Items: services,
	}
	_, g.options.Service, _ = servicePrompt.Run()

	return g.generateSpecificDeployment()
}

// generateInteractiveCICD handles interactive CI/CD generation
func (g *DeploymentGenerator) generateInteractiveCICD() error {
	// CI/CD provider
	providerPrompt := promptui.Select{
		Label: ui.IconBuild + " Select CI/CD provider",
		Items: []string{"github-actions", "gitlab-ci", "jenkins", "circleci"},
	}
	_, provider, _ := providerPrompt.Run()
	g.options.Provider = provider

	return g.generateSpecificDeployment()
}

// getFullSuiteConfig creates configuration for full deployment suite
func (g *DeploymentGenerator) getFullSuiteConfig() (models.DeploymentConfig, error) {
	var config models.DeploymentConfig

	// Get app name
	namePrompt := promptui.Prompt{
		Label:   ui.IconPackage + " Application name",
		Default: "my-app",
	}
	appName, err := namePrompt.Run()
	if err != nil {
		return config, err
	}

	// Get cloud provider if not specified
	provider := g.options.Provider
	if provider == "" {
		providerPrompt := promptui.Select{
			Label: ui.IconDocker + " Primary cloud provider",
			Items: []string{"aws", "gcp", "azure"},
		}
		_, provider, err = providerPrompt.Run()
		if err != nil {
			return config, err
		}
	}

	// Get environment
	environment := g.options.Environment
	if environment == "" {
		envPrompt := promptui.Select{
			Label: ui.IconGear + " Target environment",
			Items: []string{"development", "staging", "production"},
		}
		_, environment, err = envPrompt.Run()
		if err != nil {
			return config, err
		}
	}

	config = models.GetDefaultDeploymentSuite(models.CloudProvider(provider))
	config.AppName = appName
	config.Environment = environment

	return config, nil
}

// getSpecificDeploymentConfig creates configuration for specific deployment
func (g *DeploymentGenerator) getSpecificDeploymentConfig() (models.DeploymentConfig, error) {
	config := models.DeploymentConfig{
		Type:        models.DeploymentType(g.options.Type),
		AppName:     "my-app",
		Version:     "v1.0.0",
		Environment: g.options.Environment,
		Port:        8080,
		Namespace:   g.options.Namespace,
		MultiStage:  g.options.MultiStage,
		Optimize:    g.options.Optimize,
		Security:    g.options.Security,
		WithIngress: g.options.WithIngress,
		WithSecrets: g.options.WithSecrets,
		WithHPA:     g.options.WithHPA,
		Replicas:    3,
	}

	if g.options.Provider != "" {
		config.CloudProvider = models.CloudProvider(g.options.Provider)
		config.Region = getDefaultRegionForProvider(config.CloudProvider)
	}

	if g.options.Service != "" {
		config.CloudService = models.CloudService(g.options.Service)
	}

	if g.options.Environment == "" {
		config.Environment = "production"
	}

	return config, nil
}

// createDeploymentDirectories creates deployment directory structure
func (g *DeploymentGenerator) createDeploymentDirectories() error {
	dirs := []string{
		"deployment",
		"deployment/docker",
		"deployment/kubernetes",
		"deployment/cloud/aws",
		"deployment/cloud/gcp",
		"deployment/cloud/azure",
		"deployment/cicd",
		"deployment/monitoring",
		"deployment/terraform",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// generateDockerFiles generates Docker configuration files
func (g *DeploymentGenerator) generateDockerFiles(config models.DeploymentConfig) error {
	files := []string{
		config.GetDockerFileName(),
		"docker-compose.yml",
		"docker-compose.production.yml",
		".dockerignore",
	}

	for _, file := range files {
		content := templates.GetDockerTemplate(file, config)
		filePath := filepath.Join("deployment/docker", file)

		if err := g.writeFile(filePath, content); err != nil {
			return fmt.Errorf("failed to write Docker file %s: %w", file, err)
		}

		ui.PrintFileCreated(filePath)
	}

	return nil
}

// generateKubernetesFiles generates Kubernetes manifest files
func (g *DeploymentGenerator) generateKubernetesFiles(config models.DeploymentConfig) error {
	files := config.GetKubernetesFiles()

	for _, file := range files {
		content := templates.GetKubernetesTemplate(file, config)
		filePath := filepath.Join("deployment/kubernetes", file)

		if err := g.writeFile(filePath, content); err != nil {
			return fmt.Errorf("failed to write Kubernetes file %s: %w", file, err)
		}

		ui.PrintFileCreated(filePath)
	}

	return nil
}

// generateCloudFiles generates cloud-specific deployment files
func (g *DeploymentGenerator) generateCloudFiles(config models.DeploymentConfig) error {
	files := config.GetCloudFiles()
	providerDir := string(config.CloudProvider)

	for _, file := range files {
		content := templates.GetCloudTemplate(file, config)
		filePath := filepath.Join("deployment/cloud", providerDir, file)

		if err := g.writeFile(filePath, content); err != nil {
			return fmt.Errorf("failed to write cloud file %s: %w", file, err)
		}

		ui.PrintFileCreated(filePath)
	}

	return nil
}

// generateCICDFiles generates CI/CD pipeline files
func (g *DeploymentGenerator) generateCICDFiles(config models.DeploymentConfig) error {
	files := config.GetCICDFiles()

	for _, file := range files {
		content := templates.GetCICDTemplate(file, config)
		
		// Handle special directory structure for GitHub Actions
		var filePath string
		if strings.HasPrefix(file, ".github/") {
			filePath = file
		} else {
			filePath = filepath.Join("deployment/cicd", file)
		}

		if err := g.writeFile(filePath, content); err != nil {
			return fmt.Errorf("failed to write CI/CD file %s: %w", file, err)
		}

		ui.PrintFileCreated(filePath)
	}

	return nil
}

// writeFile writes content to a file
func (g *DeploymentGenerator) writeFile(filePath, content string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Write file
	return os.WriteFile(filePath, []byte(content), 0644)
}

// showDeploymentSummary shows summary of generated deployment files
func (g *DeploymentGenerator) showDeploymentSummary(config models.DeploymentConfig) {
	ui.PrintInfo("Generated deployment suite:")

	components := []struct {
		icon        string
		name        string
		description string
	}{
		{ui.IconDocker, "Docker", "Multi-stage Dockerfile with security hardening"},
		{ui.IconPackage, "Kubernetes", "Complete K8s manifests with scaling and monitoring"},
		{ui.IconAPI, "Cloud", fmt.Sprintf("%s deployment with %s", strings.ToUpper(string(config.CloudProvider)), string(config.CloudService))},
		{ui.IconBuild, "CI/CD", fmt.Sprintf("%s pipeline with testing and security", string(config.CICDProvider))},
	}

	if config.WithMonitoring {
		components = append(components, struct {
			icon        string
			name        string
			description string
		}{ui.IconHealth, "Monitoring", "Prometheus and Grafana configurations"})
	}

	for _, comp := range components {
		fmt.Printf("  %s %s: %s\n",
			comp.icon,
			ui.Bold.Sprint(comp.name),
			ui.Muted.Sprint(comp.description))
	}

	ui.PrintSeparator()
	ui.PrintInfo("Next steps:")
	fmt.Println("  1. Review and customize generated configurations")
	fmt.Println("  2. Set up environment variables and secrets")
	fmt.Println("  3. Test deployment in staging environment")
	fmt.Println("  4. Configure monitoring and alerting")
}

// Helper function
func getDefaultRegionForProvider(provider models.CloudProvider) string {
	switch provider {
	case models.AWSProvider:
		return "us-east-1"
	case models.GCPProvider:
		return "us-central1"
	case models.AzureProvider:
		return "East US"
	default:
		return "us-east-1"
	}
}