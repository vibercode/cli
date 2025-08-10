package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vibercode/cli/internal/generator"
	"github.com/vibercode/cli/pkg/ui"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "🔨 Generate Go code components",
	Long: ui.Bold.Sprint("Generate") + " various Go code components like APIs, resources, models, etc.\n\n" +
		ui.Bold.Sprint("Available generators:") + "\n" +
		"  " + ui.IconAPI + " api        - Complete Go API project with Docker setup\n" +
		"  " + ui.IconCode + " resource   - CRUD resource with clean architecture\n" +
		"  " + ui.IconGear + " middleware - Authentication, logging, CORS, rate limiting\n" +
		"  " + ui.IconReact + " ui         - Frontend components with Atomic Design\n" +
		"  " + ui.IconTest + " test       - Unit, integration, and benchmark tests\n" +
		"  " + ui.IconDocker + " deployment - Docker, Kubernetes, cloud deployment\n" +
		"  " + ui.IconCode + " plugin     - Plugin scaffolding and templates\n",
}

var generateAPICmd = &cobra.Command{
	Use:   "api",
	Short: "🌐 Generate a complete Go API project",
	Long: ui.Bold.Sprint("Generate a complete Go API project") + "\n\n" +
		"This command creates a production-ready Go API with:\n" +
		"  " + ui.IconPackage + " Clean architecture (handlers, services, repositories)\n" +
		"  " + ui.IconDatabase + " Database integration (PostgreSQL, MySQL, SQLite)\n" +
		"  " + ui.IconDocker + " Docker setup with docker-compose\n" +
		"  " + ui.IconGear + " Environment configuration\n" +
		"  " + ui.IconBuild + " Makefile with common commands\n" +
		"  " + ui.IconDoc + " Complete documentation\n",
	RunE: func(cmd *cobra.Command, args []string) error {
		gen := generator.NewAPIGenerator()
		return gen.Generate()
	},
}

var generateResourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "📝 Generate a CRUD resource",
	Long: ui.Bold.Sprint("Generate a complete CRUD resource") + "\n\n" +
		"This command creates a full CRUD resource with:\n" +
		"  " + ui.IconCode + " Model with GORM annotations\n" +
		"  " + ui.IconAPI + " HTTP handlers with all endpoints\n" +
		"  " + ui.IconGear + " Service layer with business logic\n" +
		"  " + ui.IconDatabase + " Repository with database operations\n" +
		"  " + ui.IconTest + " Validation and error handling\n",
	RunE: func(cmd *cobra.Command, args []string) error {
		gen := generator.NewResourceGenerator()
		return gen.Generate()
	},
}

var generateUICmd = &cobra.Command{
	Use:   "ui",
	Short: "🎨 Generate UI components with Atomic Design",
	Long: ui.Bold.Sprint("Generate UI components following Atomic Design methodology") + "\n\n" +
		"This command creates frontend components with:\n" +
		"  " + ui.IconCode + " Atoms (buttons, inputs, labels)\n" +
		"  " + ui.IconPackage + " Molecules (forms, cards, navigation)\n" +
		"  " + ui.IconBuild + " Organisms (headers, sidebars, sections)\n" +
		"  " + ui.IconDoc + " Templates (page layouts, grids)\n" +
		"  " + ui.IconGear + " TypeScript definitions and CSS modules\n" +
		"  " + ui.IconTest + " Storybook stories and unit tests\n\n" +
		ui.Bold.Sprint("Flags:") + "\n" +
		"  --atomic-design    Generate complete Atomic Design structure\n" +
		"  --framework        Choose framework (react, vue, angular)\n" +
		"  --typescript       Generate TypeScript components\n" +
		"  --storybook        Include Storybook stories\n",
	RunE: func(cmd *cobra.Command, args []string) error {
		atomicDesign, _ := cmd.Flags().GetBool("atomic-design")
		framework, _ := cmd.Flags().GetString("framework")
		typescript, _ := cmd.Flags().GetBool("typescript")
		storybook, _ := cmd.Flags().GetBool("storybook")

		gen := generator.NewUIGenerator()
		return gen.Generate(generator.UIOptions{
			AtomicDesign: atomicDesign,
			Framework:    framework,
			TypeScript:   typescript,
			Storybook:    storybook,
		})
	},
}

var generateMiddlewareCmd = &cobra.Command{
	Use:   "middleware",
	Short: "🔧 Generate middleware components",
	Long: ui.Bold.Sprint("Generate middleware components") + "\n\n" +
		"This command creates middleware with:\n" +
		"  " + ui.IconGear + " Authentication (JWT, API Key)\n" +
		"  " + ui.IconDoc + " Logging and monitoring\n" +
		"  " + ui.IconCORS + " CORS configuration\n" +
		"  " + ui.IconSpeed + " Rate limiting\n" +
		"  " + ui.IconCode + " Custom middleware templates\n\n" +
		ui.Bold.Sprint("Examples:") + "\n" +
		"  vibercode generate middleware --type auth\n" +
		"  vibercode generate middleware --preset api-security\n" +
		"  vibercode generate middleware --name CustomValidator --custom\n",
	RunE: func(cmd *cobra.Command, args []string) error {
		middlewareType, _ := cmd.Flags().GetString("type")
		customName, _ := cmd.Flags().GetString("name")
		isCustom, _ := cmd.Flags().GetBool("custom")
		preset, _ := cmd.Flags().GetString("preset")

		gen := generator.NewMiddlewareGenerator()
		return gen.Generate(generator.MiddlewareOptions{
			Type:   middlewareType,
			Name:   customName,
			Custom: isCustom,
			Preset: preset,
		})
	},
}

var generateTestCmd = &cobra.Command{
	Use:   "test",
	Short: "🧪 Generate test files and utilities",
	Long: ui.Bold.Sprint("Generate comprehensive test suite") + "\n\n" +
		"This command creates tests with:\n" +
		"  " + ui.IconTest + " Unit tests (handlers, services, repositories)\n" +
		"  " + ui.IconAPI + " Integration tests (API endpoints)\n" +
		"  " + ui.IconGear + " Test utilities (database, server, client)\n" +
		"  " + ui.IconCode + " Mock generation\n" +
		"  " + ui.IconSpeed + " Benchmark tests\n\n" +
		ui.Bold.Sprint("Examples:") + "\n" +
		"  vibercode generate test --full-suite\n" +
		"  vibercode generate test --type unit --target handler --name User\n" +
		"  vibercode generate test --type integration --name User\n" +
		"  vibercode generate test --type mock --target service --name User\n" +
		"  vibercode generate test --framework ginkgo --bdd\n",
	RunE: func(cmd *cobra.Command, args []string) error {
		testType, _ := cmd.Flags().GetString("type")
		framework, _ := cmd.Flags().GetString("framework")
		target, _ := cmd.Flags().GetString("target")
		name, _ := cmd.Flags().GetString("name")
		fullSuite, _ := cmd.Flags().GetBool("full-suite")
		withMocks, _ := cmd.Flags().GetBool("with-mocks")
		withUtils, _ := cmd.Flags().GetBool("with-utils")
		withBench, _ := cmd.Flags().GetBool("with-bench")
		bddStyle, _ := cmd.Flags().GetBool("bdd")

		gen := generator.NewTestingGenerator()
		return gen.Generate(generator.TestingOptions{
			Type:      testType,
			Framework: framework,
			Target:    target,
			Name:      name,
			FullSuite: fullSuite,
			WithMocks: withMocks,
			WithUtils: withUtils,
			WithBench: withBench,
			BDDStyle:  bddStyle,
		})
	},
}

var generateDeploymentCmd = &cobra.Command{
	Use:   "deployment",
	Short: "🚀 Generate deployment configurations",
	Long: ui.Bold.Sprint("Generate deployment configurations") + "\n\n" +
		"This command creates deployment files with:\n" +
		"  " + ui.IconDocker + " Multi-stage Docker builds with security\n" +
		"  " + ui.IconPackage + " Kubernetes manifests with scaling\n" +
		"  " + ui.IconAPI + " Cloud deployment (AWS, GCP, Azure)\n" +
		"  " + ui.IconBuild + " CI/CD pipelines\n" +
		"  " + ui.IconHealth + " Monitoring and logging\n\n" +
		ui.Bold.Sprint("Examples:") + "\n" +
		"  vibercode generate deployment --full-suite --provider aws\n" +
		"  vibercode generate deployment --type docker --multi-stage --security\n" +
		"  vibercode generate deployment --type kubernetes --with-ingress\n" +
		"  vibercode generate deployment --type cloud --provider gcp --service run\n" +
		"  vibercode generate deployment --type cicd --provider github-actions\n",
	RunE: func(cmd *cobra.Command, args []string) error {
		deploymentType, _ := cmd.Flags().GetString("type")
		provider, _ := cmd.Flags().GetString("provider")
		service, _ := cmd.Flags().GetString("service")
		namespace, _ := cmd.Flags().GetString("namespace")
		multiStage, _ := cmd.Flags().GetBool("multi-stage")
		optimize, _ := cmd.Flags().GetBool("optimize")
		security, _ := cmd.Flags().GetBool("security")
		withIngress, _ := cmd.Flags().GetBool("with-ingress")
		withSecrets, _ := cmd.Flags().GetBool("with-secrets")
		withHPA, _ := cmd.Flags().GetBool("with-hpa")
		fullSuite, _ := cmd.Flags().GetBool("full-suite")
		environment, _ := cmd.Flags().GetString("environment")

		gen := generator.NewDeploymentGenerator()
		return gen.Generate(generator.DeploymentOptions{
			Type:        deploymentType,
			Provider:    provider,
			Service:     service,
			Namespace:   namespace,
			MultiStage:  multiStage,
			Optimize:    optimize,
			Security:    security,
			WithIngress: withIngress,
			WithSecrets: withSecrets,
			WithHPA:     withHPA,
			FullSuite:   fullSuite,
			Environment: environment,
		})
	},
}

var generatePluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "🔌 Generate plugin scaffolding and templates",
	Long: ui.Bold.Sprint("Generate plugin scaffolding and templates") + "\n\n" +
		"This command creates plugin projects with:\n" +
		"  " + ui.IconCode + " Plugin interface implementation\n" +
		"  " + ui.IconGear + " Configuration and manifest files\n" +
		"  " + ui.IconBuild + " Build and packaging setup\n" +
		"  " + ui.IconDoc + " Documentation and examples\n" +
		"  " + ui.IconTest + " Test framework integration\n\n" +
		ui.Bold.Sprint("Examples:") + "\n" +
		"  vibercode generate plugin --name my-generator --type generator\n" +
		"  vibercode generate plugin --name custom-templates --type template\n" +
		"  vibercode generate plugin --name deploy-helper --type command\n" +
		"  vibercode generate plugin --template generator-example\n",
	RunE: func(cmd *cobra.Command, args []string) error {
		pluginName, _ := cmd.Flags().GetString("name")
		pluginType, _ := cmd.Flags().GetString("type")
		template, _ := cmd.Flags().GetString("template")
		author, _ := cmd.Flags().GetString("author")
		description, _ := cmd.Flags().GetString("description")

		gen := generator.NewPluginGenerator()
		return gen.Generate(generator.PluginOptions{
			Name:        pluginName,
			Type:        pluginType,
			Template:    template,
			Author:      author,
			Description: description,
		})
	},
}

func init() {
	generateCmd.AddCommand(generateAPICmd)
	generateCmd.AddCommand(generateResourceCmd)
	generateCmd.AddCommand(generateUICmd)
	generateCmd.AddCommand(generateMiddlewareCmd)
	generateCmd.AddCommand(generateTestCmd)
	generateCmd.AddCommand(generateDeploymentCmd)
	generateCmd.AddCommand(generatePluginCmd)

	// UI command flags
	generateUICmd.Flags().Bool("atomic-design", false, "Generate complete Atomic Design structure")
	generateUICmd.Flags().String("framework", "react", "Choose framework (react, vue, angular)")
	generateUICmd.Flags().Bool("typescript", true, "Generate TypeScript components")
	generateUICmd.Flags().Bool("storybook", false, "Include Storybook stories")

	// Middleware command flags
	generateMiddlewareCmd.Flags().String("type", "", "Middleware type (auth, logging, cors, rate-limit)")
	generateMiddlewareCmd.Flags().String("name", "", "Custom middleware name")
	generateMiddlewareCmd.Flags().Bool("custom", false, "Generate custom middleware")
	generateMiddlewareCmd.Flags().String("preset", "", "Middleware preset (api-security, web-app, microservice, public-api)")

	// Test command flags
	generateTestCmd.Flags().String("type", "", "Test type (unit, integration, benchmark, mock, utils)")
	generateTestCmd.Flags().String("framework", "testify", "Testing framework (testify, ginkgo, goconvey)")
	generateTestCmd.Flags().String("target", "", "Test target (handler, service, repository, middleware, api)")
	generateTestCmd.Flags().String("name", "", "Component name to test")
	generateTestCmd.Flags().Bool("full-suite", false, "Generate complete test suite")
	generateTestCmd.Flags().Bool("with-mocks", false, "Include mock generation")
	generateTestCmd.Flags().Bool("with-utils", false, "Include test utilities")
	generateTestCmd.Flags().Bool("with-bench", false, "Include benchmark tests")
	generateTestCmd.Flags().Bool("bdd", false, "Use BDD style tests (Ginkgo/Gomega)")

	// Deployment command flags
	generateDeploymentCmd.Flags().String("type", "", "Deployment type (docker, kubernetes, cloud, cicd)")
	generateDeploymentCmd.Flags().String("provider", "", "Cloud provider (aws, gcp, azure)")
	generateDeploymentCmd.Flags().String("service", "", "Cloud service (ecs, run, containers, etc.)")
	generateDeploymentCmd.Flags().String("namespace", "default", "Kubernetes namespace")
	generateDeploymentCmd.Flags().String("environment", "production", "Target environment")
	generateDeploymentCmd.Flags().Bool("multi-stage", false, "Use multi-stage Docker build")
	generateDeploymentCmd.Flags().Bool("optimize", false, "Enable image optimization")
	generateDeploymentCmd.Flags().Bool("security", false, "Enable security hardening")
	generateDeploymentCmd.Flags().Bool("with-ingress", false, "Include Kubernetes Ingress")
	generateDeploymentCmd.Flags().Bool("with-secrets", false, "Include Secrets and ConfigMaps")
	generateDeploymentCmd.Flags().Bool("with-hpa", false, "Include Horizontal Pod Autoscaler")
	generateDeploymentCmd.Flags().Bool("full-suite", false, "Generate complete deployment suite")

	// Plugin command flags
	generatePluginCmd.Flags().String("name", "", "Plugin name (required)")
	generatePluginCmd.Flags().String("type", "generator", "Plugin type (generator, template, command, integration)")
	generatePluginCmd.Flags().String("template", "", "Plugin template to use")
	generatePluginCmd.Flags().String("author", "", "Plugin author")
	generatePluginCmd.Flags().String("description", "", "Plugin description")
	generatePluginCmd.MarkFlagRequired("name")
}