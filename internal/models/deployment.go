package models

import (
	"fmt"
)

// DeploymentType represents different deployment types
type DeploymentType string

const (
	DockerDeployment     DeploymentType = "docker"
	KubernetesDeployment DeploymentType = "kubernetes"
	CloudDeployment      DeploymentType = "cloud"
	CICDDeployment       DeploymentType = "cicd"
)

// CloudProvider represents supported cloud providers
type CloudProvider string

const (
	AWSProvider   CloudProvider = "aws"
	GCPProvider   CloudProvider = "gcp"
	AzureProvider CloudProvider = "azure"
)

// CloudService represents cloud services
type CloudService string

const (
	// AWS Services
	ECSService     CloudService = "ecs"
	FargateService CloudService = "fargate"
	EKSService     CloudService = "eks"
	LambdaService  CloudService = "lambda"
	
	// GCP Services
	CloudRunService    CloudService = "run"
	GKEService         CloudService = "gke"
	AppEngineService   CloudService = "appengine"
	
	// Azure Services
	ContainerInstances CloudService = "containers"
	AKSService         CloudService = "aks"
	WebAppService      CloudService = "webapp"
)

// CICDProvider represents CI/CD providers
type CICDProvider string

const (
	GitHubActions CICDProvider = "github-actions"
	GitLabCI      CICDProvider = "gitlab-ci"
	Jenkins       CICDProvider = "jenkins"
	CircleCI      CICDProvider = "circleci"
)

// DeploymentConfig represents deployment configuration
type DeploymentConfig struct {
	Type              DeploymentType   `json:"type"`
	AppName           string           `json:"app_name"`
	Version           string           `json:"version"`
	Environment       string           `json:"environment"`
	Port              int              `json:"port"`
	
	// Docker specific
	MultiStage        bool             `json:"multi_stage"`
	Optimize          bool             `json:"optimize"`
	Security          bool             `json:"security"`
	BaseImage         string           `json:"base_image"`
	
	// Kubernetes specific
	Namespace         string           `json:"namespace"`
	Replicas          int              `json:"replicas"`
	WithIngress       bool             `json:"with_ingress"`
	WithSecrets       bool             `json:"with_secrets"`
	WithHPA           bool             `json:"with_hpa"`
	
	// Cloud specific
	CloudProvider     CloudProvider    `json:"cloud_provider"`
	CloudService      CloudService     `json:"cloud_service"`
	Region            string           `json:"region"`
	
	// CI/CD specific
	CICDProvider      CICDProvider     `json:"cicd_provider"`
	WithSecurity      bool             `json:"with_security_scan"`
	WithTesting       bool             `json:"with_testing"`
	
	// Resource limits
	CPURequest        string           `json:"cpu_request"`
	CPULimit          string           `json:"cpu_limit"`
	MemoryRequest     string           `json:"memory_request"`
	MemoryLimit       string           `json:"memory_limit"`
	
	// Monitoring
	WithMonitoring    bool             `json:"with_monitoring"`
	WithLogging       bool             `json:"with_logging"`
	
	// Full suite
	FullSuite         bool             `json:"full_suite"`
}

// DeploymentFile represents a deployment file to be generated
type DeploymentFile struct {
	Name        string            `json:"name"`
	Path        string            `json:"path"`
	Type        DeploymentType    `json:"type"`
	Content     string            `json:"content"`
	Variables   map[string]string `json:"variables"`
}

// GetDockerFileName returns Docker file name based on configuration
func (dc *DeploymentConfig) GetDockerFileName() string {
	if dc.MultiStage {
		return "Dockerfile.multi-stage"
	}
	if dc.Environment == "production" {
		return "Dockerfile.production"
	}
	return "Dockerfile"
}

// GetKubernetesFiles returns list of Kubernetes files to generate
func (dc *DeploymentConfig) GetKubernetesFiles() []string {
	files := []string{"deployment.yaml", "service.yaml"}
	
	if dc.Namespace != "" {
		files = append(files, "namespace.yaml")
	}
	if dc.WithIngress {
		files = append(files, "ingress.yaml")
	}
	if dc.WithSecrets {
		files = append(files, "secret.yaml", "configmap.yaml")
	}
	if dc.WithHPA {
		files = append(files, "hpa.yaml")
	}
	
	return files
}

// GetCloudFiles returns list of cloud-specific files
func (dc *DeploymentConfig) GetCloudFiles() []string {
	var files []string
	
	switch dc.CloudProvider {
	case AWSProvider:
		switch dc.CloudService {
		case ECSService, FargateService:
			files = append(files, "ecs-task-definition.json", "ecs-service.json")
		case EKSService:
			files = append(files, "eks-cluster.yaml", "eks-nodegroup.yaml")
		case LambdaService:
			files = append(files, "lambda-function.json", "api-gateway.json")
		}
		files = append(files, "terraform/main.tf", "cloudformation/template.yaml")
		
	case GCPProvider:
		switch dc.CloudService {
		case CloudRunService:
			files = append(files, "cloud-run.yaml")
		case GKEService:
			files = append(files, "gke-cluster.yaml")
		case AppEngineService:
			files = append(files, "app.yaml")
		}
		files = append(files, "terraform/main.tf")
		
	case AzureProvider:
		switch dc.CloudService {
		case ContainerInstances:
			files = append(files, "container-instances.json")
		case AKSService:
			files = append(files, "aks-cluster.json")
		case WebAppService:
			files = append(files, "webapp.json")
		}
		files = append(files, "arm-templates/template.json")
	}
	
	return files
}

// GetCICDFiles returns list of CI/CD files
func (dc *DeploymentConfig) GetCICDFiles() []string {
	var files []string
	
	switch dc.CICDProvider {
	case GitHubActions:
		files = append(files, ".github/workflows/ci.yml")
		if dc.WithSecurity {
			files = append(files, ".github/workflows/security.yml")
		}
		files = append(files, ".github/workflows/cd.yml")
		
	case GitLabCI:
		files = append(files, ".gitlab-ci.yml")
		
	case Jenkins:
		files = append(files, "Jenkinsfile")
		
	case CircleCI:
		files = append(files, ".circleci/config.yml")
	}
	
	return files
}

// GetBaseImage returns appropriate base image
func (dc *DeploymentConfig) GetBaseImage() string {
	if dc.BaseImage != "" {
		return dc.BaseImage
	}
	
	if dc.MultiStage {
		return "golang:1.21-alpine"
	}
	
	if dc.Security {
		return "alpine:latest"
	}
	
	return "ubuntu:22.04"
}

// GetFinalImage returns final stage image for multi-stage builds
func (dc *DeploymentConfig) GetFinalImage() string {
	if dc.Security {
		return "scratch"
	}
	return "alpine:latest"
}

// GetResourceRequests returns resource requests for Kubernetes
func (dc *DeploymentConfig) GetResourceRequests() map[string]string {
	requests := make(map[string]string)
	
	if dc.CPURequest != "" {
		requests["cpu"] = dc.CPURequest
	} else {
		requests["cpu"] = "50m"
	}
	
	if dc.MemoryRequest != "" {
		requests["memory"] = dc.MemoryRequest
	} else {
		requests["memory"] = "64Mi"
	}
	
	return requests
}

// GetResourceLimits returns resource limits for Kubernetes
func (dc *DeploymentConfig) GetResourceLimits() map[string]string {
	limits := make(map[string]string)
	
	if dc.CPULimit != "" {
		limits["cpu"] = dc.CPULimit
	} else {
		limits["cpu"] = "100m"
	}
	
	if dc.MemoryLimit != "" {
		limits["memory"] = dc.MemoryLimit
	} else {
		limits["memory"] = "128Mi"
	}
	
	return limits
}

// GetEnvironmentVariables returns environment variables for deployment
func (dc *DeploymentConfig) GetEnvironmentVariables() map[string]string {
	envVars := make(map[string]string)
	
	envVars["ENV"] = dc.Environment
	envVars["PORT"] = fmt.Sprintf("%d", dc.Port)
	envVars["APP_NAME"] = dc.AppName
	envVars["VERSION"] = dc.Version
	
	return envVars
}

// GetLabels returns labels for Kubernetes resources
func (dc *DeploymentConfig) GetLabels() map[string]string {
	labels := make(map[string]string)
	
	labels["app"] = dc.AppName
	labels["version"] = dc.Version
	labels["environment"] = dc.Environment
	
	return labels
}

// GetDefaultDeploymentSuite returns default deployment suite configuration
func GetDefaultDeploymentSuite(provider CloudProvider) DeploymentConfig {
	config := DeploymentConfig{
		Type:          DockerDeployment,
		AppName:       "my-app",
		Version:       "v1.0.0",
		Environment:   "production",
		Port:          8080,
		MultiStage:    true,
		Optimize:      true,
		Security:      true,
		Namespace:     "default",
		Replicas:      3,
		WithIngress:   true,
		WithSecrets:   true,
		WithHPA:       true,
		CloudProvider: provider,
		Region:        getDefaultRegion(provider),
		WithSecurity:  true,
		WithTesting:   true,
		CPURequest:    "50m",
		CPULimit:      "100m",
		MemoryRequest: "64Mi",
		MemoryLimit:   "128Mi",
		WithMonitoring: true,
		WithLogging:   true,
		FullSuite:     true,
	}
	
	// Set default cloud service based on provider
	switch provider {
	case AWSProvider:
		config.CloudService = ECSService
		config.CICDProvider = GitHubActions
	case GCPProvider:
		config.CloudService = CloudRunService
		config.CICDProvider = GitHubActions
	case AzureProvider:
		config.CloudService = ContainerInstances
		config.CICDProvider = GitHubActions
	}
	
	return config
}

// Helper functions

func getDefaultRegion(provider CloudProvider) string {
	switch provider {
	case AWSProvider:
		return "us-east-1"
	case GCPProvider:
		return "us-central1"
	case AzureProvider:
		return "East US"
	default:
		return "us-east-1"
	}
}

// IsValidDeploymentType checks if deployment type is valid
func IsValidDeploymentType(deploymentType string) bool {
	validTypes := []string{
		string(DockerDeployment),
		string(KubernetesDeployment),
		string(CloudDeployment),
		string(CICDDeployment),
	}
	
	for _, validType := range validTypes {
		if deploymentType == validType {
			return true
		}
	}
	return false
}

// IsValidCloudProvider checks if cloud provider is valid
func IsValidCloudProvider(provider string) bool {
	validProviders := []string{
		string(AWSProvider),
		string(GCPProvider),
		string(AzureProvider),
	}
	
	for _, validProvider := range validProviders {
		if provider == validProvider {
			return true
		}
	}
	return false
}

// IsValidCloudService checks if cloud service is valid
func IsValidCloudService(service string) bool {
	validServices := []string{
		string(ECSService),
		string(FargateService),
		string(EKSService),
		string(LambdaService),
		string(CloudRunService),
		string(GKEService),
		string(AppEngineService),
		string(ContainerInstances),
		string(AKSService),
		string(WebAppService),
	}
	
	for _, validService := range validServices {
		if service == validService {
			return true
		}
	}
	return false
}

// IsValidCICDProvider checks if CI/CD provider is valid
func IsValidCICDProvider(provider string) bool {
	validProviders := []string{
		string(GitHubActions),
		string(GitLabCI),
		string(Jenkins),
		string(CircleCI),
	}
	
	for _, validProvider := range validProviders {
		if provider == validProvider {
			return true
		}
	}
	return false
}