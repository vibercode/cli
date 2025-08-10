package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vibercode/cli/internal/generator"
	"github.com/vibercode/cli/pkg/ui"
)

var runCmd = &cobra.Command{
	Use:   "run [project_name]",
	Short: "🚀 Run a generated Vibercode project",
	Long: ui.Bold.Sprint("Run a generated Vibercode project") + "\n\n" +
		"This command runs a previously generated project by reading its manifest file.\n" +
		"The command will look for the project in the current directory or the specified path.\n\n" +
		ui.Bold.Sprint("Examples:") + "\n" +
		"  vibercode run my-api     # Run project in ./my-api\n" +
		"  vibercode run .          # Run project in current directory\n",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectPath := args[0]
		
		// If projectPath is ".", use current directory
		if projectPath == "." {
			pwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			projectPath = pwd
		}

		return runProject(projectPath)
	},
}

func runProject(projectPath string) error {
	ui.PrintHeader("Running Vibercode Project")

	// Check if project directory exists
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return fmt.Errorf("project directory '%s' does not exist", projectPath)
	}

	// Look for manifest file
	manifestPath := filepath.Join(projectPath, ".vibercode", "manifest.vibe")
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		return fmt.Errorf("manifest file not found at '%s'. This doesn't appear to be a Vibercode project", manifestPath)
	}

	// Read and parse manifest
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to read manifest file: %w", err)
	}

	var manifest generator.VibercodeManifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return fmt.Errorf("failed to parse manifest file: %w", err)
	}

	// Show project info
	ui.PrintSubHeader("Project Information")
	ui.PrintFeature(ui.IconAPI, "Project", manifest.Name)
	ui.PrintFeature(ui.IconGear, "Type", manifest.ProjectType)
	ui.PrintFeature(ui.IconGear, "Port", manifest.Port)
	ui.PrintFeature(ui.IconDatabase, "Database", manifest.Database.GetDisplayName())
	ui.PrintFeature(ui.IconGear, "Generated", manifest.GeneratedAt)
	fmt.Println()

	// Check if docker-compose.yml exists
	dockerComposePath := filepath.Join(projectPath, "docker-compose.yml")
	if _, err := os.Stat(dockerComposePath); err == nil {
		// Run with Docker
		ui.PrintInfo("Running project with Docker Compose...")
		return runWithDocker(projectPath)
	}

	// Run locally
	ui.PrintInfo("Running project locally...")
	return runLocally(projectPath)
}

func runWithDocker(projectPath string) error {
	ui.PrintStep(1, 2, "Building and starting containers")
	
	cmd := exec.Command("docker-compose", "up", "--build")
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}

func runLocally(projectPath string) error {
	ui.PrintStep(1, 2, "Installing dependencies")
	
	// Run go mod tidy
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = projectPath
	if err := tidyCmd.Run(); err != nil {
		ui.PrintWarning("Failed to run 'go mod tidy': " + err.Error())
	}

	ui.PrintStep(2, 2, "Starting application")
	
	// Check if main.go exists in cmd/server/
	mainPath := filepath.Join(projectPath, "cmd", "server", "main.go")
	if _, err := os.Stat(mainPath); os.IsNotExist(err) {
		// Try alternative main.go location
		mainPath = filepath.Join(projectPath, "main.go")
		if _, err := os.Stat(mainPath); os.IsNotExist(err) {
			return fmt.Errorf("main.go not found in expected locations")
		}
	}

	// Run the application
	runCmd := exec.Command("go", "run", "cmd/server/main.go")
	if mainPath == filepath.Join(projectPath, "main.go") {
		runCmd = exec.Command("go", "run", "main.go")
	}
	
	runCmd.Dir = projectPath
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	
	ui.PrintSuccess("Application started successfully!")
	return runCmd.Run()
}