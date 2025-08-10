package generator

import (
	"testing"

	"github.com/vibercode/cli/internal/models"
)

func TestTestingGenerator_Creation(t *testing.T) {
	gen := NewTestingGenerator()
	if gen == nil {
		t.Fatal("NewTestingGenerator() returned nil")
	}
}

func TestTestingOptions(t *testing.T) {
	options := TestingOptions{
		Type:      "unit",
		Framework: "testify",
		Target:    "handler",
		Name:      "User",
		FullSuite: false,
		WithMocks: true,
		WithUtils: true,
		WithBench: false,
		BDDStyle:  false,
	}

	if options.Type != "unit" {
		t.Errorf("Expected Type to be 'unit', got %s", options.Type)
	}

	if options.Framework != "testify" {
		t.Errorf("Expected Framework to be 'testify', got %s", options.Framework)
	}
}

func TestTestingModels_ValidTestTypes(t *testing.T) {
	validTypes := []models.TestType{
		models.UnitTest,
		models.IntegrationTest,
		models.BenchmarkTest,
		models.MockTest,
		models.UtilsTest,
	}

	expectedValues := []string{"unit", "integration", "benchmark", "mock", "utils"}

	for i, testType := range validTypes {
		if string(testType) != expectedValues[i] {
			t.Errorf("Expected %s, got %s", expectedValues[i], string(testType))
		}
	}
}

func TestTestingModels_ValidFrameworks(t *testing.T) {
	frameworks := []models.TestFramework{
		models.TestifyFramework,
		models.GinkgoFramework,
		models.GoConveyFramework,
	}

	expectedValues := []string{"testify", "ginkgo", "goconvey"}

	for i, framework := range frameworks {
		if string(framework) != expectedValues[i] {
			t.Errorf("Expected %s, got %s", expectedValues[i], string(framework))
		}
	}
}

func TestTestingModels_ValidTargets(t *testing.T) {
	targets := []models.TestTarget{
		models.ResourceTarget,
		models.HandlerTarget,
		models.ServiceTarget,
		models.RepositoryTarget,
		models.MiddlewareTarget,
		models.APITarget,
		models.AllTargets,
	}

	expectedValues := []string{"resource", "handler", "service", "repository", "middleware", "api", "all"}

	for i, target := range targets {
		if string(target) != expectedValues[i] {
			t.Errorf("Expected %s, got %s", expectedValues[i], string(target))
		}
	}
}