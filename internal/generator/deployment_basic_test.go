package generator

import (
	"testing"

	"github.com/vibercode/cli/internal/models"
)

func TestDeploymentGenerator_Creation(t *testing.T) {
	gen := NewDeploymentGenerator()
	if gen == nil {
		t.Fatal("NewDeploymentGenerator() returned nil")
	}
}

func TestDeploymentOptions(t *testing.T) {
	options := DeploymentOptions{
		Type:        "docker",
		Provider:    "aws",
		Service:     "ecs",
		Namespace:   "production",
		MultiStage:  true,
		Optimize:    true,
		Security:    true,
		WithIngress: false,
		WithSecrets: false,
		WithHPA:     false,
		FullSuite:   false,
		Environment: "production",
	}

	if options.Type != "docker" {
		t.Errorf("Expected Type to be 'docker', got %s", options.Type)
	}

	if options.Provider != "aws" {
		t.Errorf("Expected Provider to be 'aws', got %s", options.Provider)
	}

	if !options.MultiStage {
		t.Error("Expected MultiStage to be true")
	}

	if !options.Security {
		t.Error("Expected Security to be true")
	}
}

func TestDeploymentModels_ValidDeploymentTypes(t *testing.T) {
	validTypes := []models.DeploymentType{
		models.DockerDeployment,
		models.KubernetesDeployment,
		models.CloudDeployment,
		models.CICDDeployment,
	}

	expectedValues := []string{"docker", "kubernetes", "cloud", "cicd"}

	for i, deploymentType := range validTypes {
		if string(deploymentType) != expectedValues[i] {
			t.Errorf("Expected %s, got %s", expectedValues[i], string(deploymentType))
		}
	}
}

func TestDeploymentModels_ValidCloudProviders(t *testing.T) {
	providers := []models.CloudProvider{
		models.AWSProvider,
		models.GCPProvider,
		models.AzureProvider,
	}

	expectedValues := []string{"aws", "gcp", "azure"}

	for i, provider := range providers {
		if string(provider) != expectedValues[i] {
			t.Errorf("Expected %s, got %s", expectedValues[i], string(provider))
		}
	}
}

func TestDeploymentModels_ValidCloudServices(t *testing.T) {
	// Test AWS services
	awsServices := []models.CloudService{
		models.ECSService,
		models.FargateService,
		models.EKSService,
		models.LambdaService,
	}
	awsExpected := []string{"ecs", "fargate", "eks", "lambda"}

	for i, service := range awsServices {
		if string(service) != awsExpected[i] {
			t.Errorf("Expected AWS service %s, got %s", awsExpected[i], string(service))
		}
	}

	// Test GCP services
	gcpServices := []models.CloudService{
		models.CloudRunService,
		models.GKEService,
		models.AppEngineService,
	}
	gcpExpected := []string{"run", "gke", "appengine"}

	for i, service := range gcpServices {
		if string(service) != gcpExpected[i] {
			t.Errorf("Expected GCP service %s, got %s", gcpExpected[i], string(service))
		}
	}

	// Test Azure services
	azureServices := []models.CloudService{
		models.ContainerInstances,
		models.AKSService,
		models.WebAppService,
	}
	azureExpected := []string{"containers", "aks", "webapp"}

	for i, service := range azureServices {
		if string(service) != azureExpected[i] {
			t.Errorf("Expected Azure service %s, got %s", azureExpected[i], string(service))
		}
	}
}

func TestDeploymentConfig_BasicValidation(t *testing.T) {
	config := models.DeploymentConfig{
		AppName:       "test-app",
		Port:          8080,
		Environment:   "production",
		Version:       "1.0.0",
		MultiStage:    true,
		Security:      true,
		CloudProvider: models.AWSProvider,
		CloudService:  models.ECSService,
	}

	if config.AppName != "test-app" {
		t.Errorf("Expected AppName to be 'test-app', got %s", config.AppName)
	}

	if config.Port != 8080 {
		t.Errorf("Expected Port to be 8080, got %d", config.Port)
	}

	if config.CloudProvider != models.AWSProvider {
		t.Errorf("Expected CloudProvider to be %s, got %s", models.AWSProvider, config.CloudProvider)
	}

	if config.CloudService != models.ECSService {
		t.Errorf("Expected CloudService to be %s, got %s", models.ECSService, config.CloudService)
	}
}