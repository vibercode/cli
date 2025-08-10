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

// TestingOptions contains configuration for test generation
type TestingOptions struct {
	Type      string
	Framework string
	Target    string
	Name      string
	FullSuite bool
	WithMocks bool
	WithUtils bool
	WithBench bool
	BDDStyle  bool
}

// TestingGenerator handles test generation
type TestingGenerator struct {
	options TestingOptions
}

// NewTestingGenerator creates a new testing generator
func NewTestingGenerator() *TestingGenerator {
	return &TestingGenerator{}
}

// Generate generates tests based on options
func (g *TestingGenerator) Generate(options TestingOptions) error {
	g.options = options

	ui.PrintStep(1, 1, "Starting test generation...")

	// Handle full suite generation
	if options.FullSuite {
		return g.generateFullTestSuite()
	}

	// Handle specific test type generation
	if options.Type != "" {
		return g.generateSpecificTestType()
	}

	// Interactive mode
	return g.generateInteractiveTests()
}

// generateFullTestSuite generates a complete test suite
func (g *TestingGenerator) generateFullTestSuite() error {
	ui.PrintStep(1, 5, "Generating full test suite...")

	config := g.getFullSuiteConfig()

	// Create test directory structure
	ui.PrintStep(2, 5, "Creating test directory structure...")
	if err := g.createTestDirectories(); err != nil {
		return fmt.Errorf("failed to create test directories: %w", err)
	}

	// Generate unit tests
	ui.PrintStep(3, 5, "Generating unit tests...")
	if err := g.generateUnitTests(config); err != nil {
		return fmt.Errorf("failed to generate unit tests: %w", err)
	}

	// Generate integration tests
	ui.PrintStep(4, 5, "Generating integration tests...")
	if err := g.generateIntegrationTests(config); err != nil {
		return fmt.Errorf("failed to generate integration tests: %w", err)
	}

	// Generate test utilities
	ui.PrintStep(5, 5, "Generating test utilities...")
	if err := g.generateTestUtilities(config); err != nil {
		return fmt.Errorf("failed to generate test utilities: %w", err)
	}

	// Generate mocks if requested
	if g.options.WithMocks {
		if err := g.generateMocks(config); err != nil {
			return fmt.Errorf("failed to generate mocks: %w", err)
		}
	}

	// Generate benchmarks if requested
	if g.options.WithBench {
		if err := g.generateBenchmarks(config); err != nil {
			return fmt.Errorf("failed to generate benchmarks: %w", err)
		}
	}

	ui.PrintSuccess("Full test suite generated successfully!")
	g.showTestSuiteSummary()

	return nil
}

// generateSpecificTestType generates a specific type of test
func (g *TestingGenerator) generateSpecificTestType() error {
	testType := models.TestType(g.options.Type)
	
	ui.PrintStep(1, 2, fmt.Sprintf("Generating %s tests...", testType))

	config := models.TestConfig{
		Type:      testType,
		Framework: models.TestFramework(g.options.Framework),
		Target:    models.TestTarget(g.options.Target),
		Name:      g.options.Name,
	}

	// Create test directories
	if err := g.createTestDirectories(); err != nil {
		return fmt.Errorf("failed to create test directories: %w", err)
	}

	ui.PrintStep(2, 2, "Creating test files...")

	switch testType {
	case models.UnitTest:
		return g.generateUnitTestForTarget(config)
	case models.IntegrationTest:
		return g.generateIntegrationTestForTarget(config)
	case models.BenchmarkTest:
		return g.generateBenchmarkTestForTarget(config)
	case models.MockTest:
		return g.generateMockForTarget(config)
	case models.UtilsTest:
		return g.generateTestUtilitiesForConfig(config)
	}

	return fmt.Errorf("unsupported test type: %s", testType)
}

// generateInteractiveTests handles interactive test generation
func (g *TestingGenerator) generateInteractiveTests() error {
	ui.PrintStep(1, 1, "Interactive test generation...")

	// Ask what to generate
	actionPrompt := promptui.Select{
		Label: ui.IconTest + " What would you like to generate?",
		Items: []string{
			"Full test suite",
			"Unit tests",
			"Integration tests", 
			"Mocks",
			"Test utilities",
			"Benchmark tests",
		},
	}

	_, action, err := actionPrompt.Run()
	if err != nil {
		return err
	}

	switch action {
	case "Full test suite":
		g.options.FullSuite = true
		return g.generateFullTestSuite()
	case "Unit tests":
		g.options.Type = string(models.UnitTest)
		return g.generateInteractiveUnitTests()
	case "Integration tests":
		g.options.Type = string(models.IntegrationTest)
		return g.generateInteractiveIntegrationTests()
	case "Mocks":
		g.options.Type = string(models.MockTest)
		return g.generateInteractiveMocks()
	case "Test utilities":
		g.options.Type = string(models.UtilsTest)
		return g.generateSpecificTestType()
	case "Benchmark tests":
		g.options.Type = string(models.BenchmarkTest)
		return g.generateInteractiveBenchmarks()
	}

	return nil
}

// generateInteractiveUnitTests handles interactive unit test generation
func (g *TestingGenerator) generateInteractiveUnitTests() error {
	// Get target
	targetPrompt := promptui.Select{
		Label: ui.IconCode + " What to test?",
		Items: []string{"handler", "service", "repository", "middleware"},
	}

	_, target, err := targetPrompt.Run()
	if err != nil {
		return err
	}

	// Get name
	namePrompt := promptui.Prompt{
		Label: ui.IconPackage + " Component name",
	}
	name, err := namePrompt.Run()
	if err != nil {
		return err
	}

	// Get framework
	framework, err := g.getFrameworkChoice()
	if err != nil {
		return err
	}

	g.options.Target = target
	g.options.Name = name
	g.options.Framework = framework

	return g.generateSpecificTestType()
}

// generateInteractiveIntegrationTests handles interactive integration test generation
func (g *TestingGenerator) generateInteractiveIntegrationTests() error {
	// Get API name
	namePrompt := promptui.Prompt{
		Label: ui.IconAPI + " API/Resource name",
	}
	name, err := namePrompt.Run()
	if err != nil {
		return err
	}

	// Get framework
	framework, err := g.getFrameworkChoice()
	if err != nil {
		return err
	}

	g.options.Target = string(models.APITarget)
	g.options.Name = name
	g.options.Framework = framework

	return g.generateSpecificTestType()
}

// generateInteractiveMocks handles interactive mock generation
func (g *TestingGenerator) generateInteractiveMocks() error {
	// Get target
	targetPrompt := promptui.Select{
		Label: ui.IconGear + " Mock target",
		Items: []string{"service", "repository"},
	}

	_, target, err := targetPrompt.Run()
	if err != nil {
		return err
	}

	// Get name
	namePrompt := promptui.Prompt{
		Label: ui.IconCode + " Interface name",
	}
	name, err := namePrompt.Run()
	if err != nil {
		return err
	}

	g.options.Target = target
	g.options.Name = name
	g.options.Framework = string(models.TestifyFramework) // Mocks work best with testify

	return g.generateSpecificTestType()
}

// generateInteractiveBenchmarks handles interactive benchmark generation
func (g *TestingGenerator) generateInteractiveBenchmarks() error {
	// Get target
	targetPrompt := promptui.Select{
		Label: ui.IconSpeed + " Benchmark target",
		Items: []string{"handler", "service", "repository"},
	}

	_, target, err := targetPrompt.Run()
	if err != nil {
		return err
	}

	// Get name
	namePrompt := promptui.Prompt{
		Label: ui.IconPackage + " Component name",
	}
	name, err := namePrompt.Run()
	if err != nil {
		return err
	}

	g.options.Target = target
	g.options.Name = name
	g.options.Framework = string(models.TestifyFramework) // Benchmarks use standard testing

	return g.generateSpecificTestType()
}

// getFrameworkChoice gets framework choice from user
func (g *TestingGenerator) getFrameworkChoice() (string, error) {
	if g.options.Framework != "" {
		return g.options.Framework, nil
	}

	frameworkPrompt := promptui.Select{
		Label: ui.IconGear + " Testing framework",
		Items: []string{"testify", "ginkgo", "goconvey"},
	}

	_, framework, err := frameworkPrompt.Run()
	return framework, err
}

// getFullSuiteConfig creates configuration for full test suite
func (g *TestingGenerator) getFullSuiteConfig() models.TestConfig {
	framework := models.TestifyFramework
	if g.options.Framework != "" {
		framework = models.TestFramework(g.options.Framework)
	}

	return models.TestConfig{
		Type:      models.UnitTest, // Base type, will generate multiple types
		Framework: framework,
		Target:    models.AllTargets,
		Name:      g.options.Name,
		FullSuite: true,
		WithMocks: g.options.WithMocks,
		WithUtils: g.options.WithUtils,
		WithBench: g.options.WithBench,
		BDDStyle:  g.options.BDDStyle,
	}
}

// createTestDirectories creates the test directory structure
func (g *TestingGenerator) createTestDirectories() error {
	dirs := []string{
		"test",
		"test/unit/handlers",
		"test/unit/services", 
		"test/unit/repositories",
		"test/unit/middleware",
		"test/integration/api",
		"test/integration/database",
		"test/integration/e2e",
		"test/benchmark",
		"test/mocks",
		"test/utils",
		"test/fixtures",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// generateUnitTests generates unit tests for all targets
func (g *TestingGenerator) generateUnitTests(config models.TestConfig) error {
	targets := []models.TestTarget{
		models.HandlerTarget,
		models.ServiceTarget,
		models.RepositoryTarget,
		models.MiddlewareTarget,
	}

	for _, target := range targets {
		testConfig := config
		testConfig.Target = target
		testConfig.Name = "User" // Default example

		if err := g.generateUnitTestForTarget(testConfig); err != nil {
			return fmt.Errorf("failed to generate unit tests for %s: %w", target, err)
		}
	}

	return nil
}

// generateUnitTestForTarget generates unit test for specific target
func (g *TestingGenerator) generateUnitTestForTarget(config models.TestConfig) error {
	testFile := g.createTestFile(config)
	content := templates.GetTestTemplate(testFile)
	
	filePath := filepath.Join(testFile.GetTestPackagePath(), testFile.GetTestFileName())
	
	if err := g.writeFile(filePath, content); err != nil {
		return fmt.Errorf("failed to write test file: %w", err)
	}

	ui.PrintFileCreated(filePath)
	return nil
}

// generateIntegrationTests generates integration tests
func (g *TestingGenerator) generateIntegrationTests(config models.TestConfig) error {
	config.Type = models.IntegrationTest
	config.Target = models.APITarget
	config.Name = "User" // Default example

	testFile := g.createTestFile(config)
	content := templates.GetIntegrationTestTemplate(testFile)
	
	filePath := filepath.Join(testFile.GetTestPackagePath(), testFile.GetTestFileName())
	
	if err := g.writeFile(filePath, content); err != nil {
		return fmt.Errorf("failed to write integration test file: %w", err)
	}

	ui.PrintFileCreated(filePath)
	return nil
}

// generateIntegrationTestForTarget generates integration test for specific target
func (g *TestingGenerator) generateIntegrationTestForTarget(config models.TestConfig) error {
	testFile := g.createTestFile(config)
	content := templates.GetIntegrationTestTemplate(testFile)
	
	filePath := filepath.Join(testFile.GetTestPackagePath(), testFile.GetTestFileName())
	
	if err := g.writeFile(filePath, content); err != nil {
		return fmt.Errorf("failed to write integration test file: %w", err)
	}

	ui.PrintFileCreated(filePath)
	return nil
}

// generateTestUtilities generates test utility functions
func (g *TestingGenerator) generateTestUtilities(config models.TestConfig) error {
	utilities := []string{
		"test_database.go",
		"test_server.go", 
		"test_client.go",
		"test_factories.go",
	}

	for _, utility := range utilities {
		content := templates.GetTestUtilityTemplate(utility, config.Framework)
		filePath := filepath.Join("test/utils", utility)
		
		if err := g.writeFile(filePath, content); err != nil {
			return fmt.Errorf("failed to write utility file %s: %w", utility, err)
		}

		ui.PrintFileCreated(filePath)
	}

	return nil
}

// generateTestUtilitiesForConfig generates test utilities for specific config
func (g *TestingGenerator) generateTestUtilitiesForConfig(config models.TestConfig) error {
	return g.generateTestUtilities(config)
}

// generateMocks generates mock files
func (g *TestingGenerator) generateMocks(config models.TestConfig) error {
	targets := []models.TestTarget{
		models.ServiceTarget,
		models.RepositoryTarget,
	}

	for _, target := range targets {
		mockConfig := config
		mockConfig.Type = models.MockTest
		mockConfig.Target = target
		mockConfig.Name = "User" // Default example

		if err := g.generateMockForTarget(mockConfig); err != nil {
			return fmt.Errorf("failed to generate mock for %s: %w", target, err)
		}
	}

	return nil
}

// generateMockForTarget generates mock for specific target
func (g *TestingGenerator) generateMockForTarget(config models.TestConfig) error {
	testFile := g.createTestFile(config)
	content := templates.GetMockTemplate(testFile)
	
	filePath := filepath.Join(testFile.GetTestPackagePath(), testFile.GetTestFileName())
	
	if err := g.writeFile(filePath, content); err != nil {
		return fmt.Errorf("failed to write mock file: %w", err)
	}

	ui.PrintFileCreated(filePath)
	return nil
}

// generateBenchmarks generates benchmark tests
func (g *TestingGenerator) generateBenchmarks(config models.TestConfig) error {
	config.Type = models.BenchmarkTest
	targets := []models.TestTarget{
		models.HandlerTarget,
		models.ServiceTarget,
	}

	for _, target := range targets {
		benchConfig := config
		benchConfig.Target = target
		benchConfig.Name = "User" // Default example

		if err := g.generateBenchmarkTestForTarget(benchConfig); err != nil {
			return fmt.Errorf("failed to generate benchmark for %s: %w", target, err)
		}
	}

	return nil
}

// generateBenchmarkTestForTarget generates benchmark test for specific target
func (g *TestingGenerator) generateBenchmarkTestForTarget(config models.TestConfig) error {
	testFile := g.createTestFile(config)
	content := templates.GetBenchmarkTemplate(testFile)
	
	filePath := filepath.Join(testFile.GetTestPackagePath(), testFile.GetTestFileName())
	
	if err := g.writeFile(filePath, content); err != nil {
		return fmt.Errorf("failed to write benchmark file: %w", err)
	}

	ui.PrintFileCreated(filePath)
	return nil
}

// createTestFile creates a TestFile model from config
func (g *TestingGenerator) createTestFile(config models.TestConfig) models.TestFile {
	testFile := models.TestFile{
		Name:      config.Name,
		Package:   strings.ToLower(string(config.Target)),
		Type:      config.Type,
		Framework: config.Framework,
		Imports:   []string{},
	}

	// Add default imports
	testFile.Imports = testFile.GetDefaultImports()

	// Add test cases
	testFile.TestCases = models.GetDefaultTestCases(config.Target, config.Name)

	// Add mocks if needed
	if config.WithMocks || config.Type == models.MockTest {
		testFile.Mocks = models.GetDefaultMocks(config.Target, config.Name)
	}

	return testFile
}

// writeFile writes content to a file
func (g *TestingGenerator) writeFile(filePath, content string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Write file
	return os.WriteFile(filePath, []byte(content), 0644)
}

// showTestSuiteSummary shows a summary of the generated test suite
func (g *TestingGenerator) showTestSuiteSummary() {
	ui.PrintInfo("Generated test suite:")
	
	components := []struct {
		icon        string
		name        string
		description string
	}{
		{ui.IconTest, "Unit Tests", "Handler, service, repository, middleware tests"},
		{ui.IconAPI, "Integration Tests", "End-to-end API testing"},
		{ui.IconGear, "Test Utilities", "Database setup, test server, HTTP client"},
	}

	if g.options.WithMocks {
		components = append(components, struct {
			icon        string
			name        string
			description string
		}{ui.IconCode, "Mocks", "Service and repository mocks"})
	}

	if g.options.WithBench {
		components = append(components, struct {
			icon        string
			name        string
			description string
		}{ui.IconSpeed, "Benchmarks", "Performance testing"})
	}

	for _, comp := range components {
		fmt.Printf("  %s %s: %s\n", 
			comp.icon, 
			ui.Bold.Sprint(comp.name), 
			ui.Muted.Sprint(comp.description))
	}

	ui.PrintSeparator()
	ui.PrintInfo("Next steps:")
	fmt.Println("  1. Run tests with: go test ./test/...")
	fmt.Println("  2. Run specific tests: go test ./test/unit/handlers")
	fmt.Println("  3. Run benchmarks: go test -bench=. ./test/benchmark")
	fmt.Println("  4. Generate coverage: go test -cover ./test/...")
}