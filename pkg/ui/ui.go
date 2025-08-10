package ui

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/pterm/pterm"
)

// Color definitions
var (
	// Main colors
	Primary   = color.New(color.FgCyan, color.Bold)
	Secondary = color.New(color.FgMagenta)
	Success   = color.New(color.FgGreen, color.Bold)
	Warning   = color.New(color.FgYellow, color.Bold)
	Error     = color.New(color.FgRed, color.Bold)
	Info      = color.New(color.FgBlue)
	Muted     = color.New(color.FgHiBlack)
	Dim       = color.New(color.FgHiBlack)
	
	// Text styles
	Bold      = color.New(color.Bold)
	Underline = color.New(color.Underline)
)

// Icons
const (
	IconCheck    = "✅"
	IconCross    = "❌"
	IconWarning  = "⚠️"
	IconInfo     = "ℹ️"
	IconRocket   = "🚀"
	IconPackage  = "📦"
	IconDatabase = "🗄️"
	IconDocker   = "🐳"
	IconGear     = "⚙️"
	IconFolder   = "📁"
	IconFile     = "📄"
	IconArrow    = "➤"
	IconStar     = "⭐"
	IconCode     = "💻"
	IconBuild    = "🔨"
	IconTest     = "🧪"
	IconDoc      = "📋"
	IconAPI      = "🌐"
	IconCLI      = "🖥️"
	IconMagic    = "✨"
	IconSpeed    = "⚡"
	IconReact    = "⚛️"
	IconCORS     = "🔀"
	IconHealth   = "💚"
	IconLogs     = "📝"
	IconWebSocket = "🔌"
)

// ShowBanner displays the Vibercode CLI banner
func ShowBanner() {
	banner := `
██╗   ██╗██╗██████╗ ███████╗██████╗  ██████╗ ██████╗ ██████╗ ███████╗
██║   ██║██║██╔══██╗██╔════╝██╔══██╗██╔════╝██╔═══██╗██╔══██╗██╔════╝
██║   ██║██║██████╔╝█████╗  ██████╔╝██║     ██║   ██║██║  ██║█████╗  
╚██╗ ██╔╝██║██╔══██╗██╔══╝  ██╔══██╗██║     ██║   ██║██║  ██║██╔══╝  
 ╚████╔╝ ██║██████╔╝███████╗██║  ██║╚██████╗╚██████╔╝██████╔╝███████╗
  ╚═══╝  ╚═╝╚═════╝ ╚══════╝╚═╝  ╚═╝ ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝`

	Primary.Println(banner)
	fmt.Println()
	
	pterm.DefaultBox.WithTitle("Vibercode CLI").WithTitleTopCenter().WithBoxStyle(&pterm.Style{
		pterm.FgCyan,
	}).Println("Generate Go web APIs with clean architecture " + IconRocket + "\n" +
		"Built with " + IconMagic + " and powered by Go templates")
	fmt.Println()
}

// PrintHeader prints a section header
func PrintHeader(title string) {
	pterm.DefaultSection.WithLevel(1).Println(title)
	fmt.Println()
}

// PrintSubHeader prints a subsection header
func PrintSubHeader(title string) {
	pterm.DefaultSection.WithLevel(2).Println(title)
}

// PrintSuccess prints a success message with icon
func PrintSuccess(message string) {
	Success.Printf("%s %s\n", IconCheck, message)
}

// PrintError prints an error message with icon
func PrintError(message string) {
	Error.Printf("%s %s\n", IconCross, message)
}

// PrintWarning prints a warning message with icon
func PrintWarning(message string) {
	Warning.Printf("%s %s\n", IconWarning, message)
}

// PrintInfo prints an info message with icon
func PrintInfo(message string) {
	Info.Printf("%s %s\n", IconInfo, message)
}

// PrintStep prints a step in the process
func PrintStep(step int, total int, message string) {
	Primary.Printf("Step %d/%d: %s %s\n", step, total, IconArrow, message)
}

// PrintFileCreated prints a file creation message
func PrintFileCreated(filename string) {
	Success.Printf("  %s Created: %s\n", IconFile, Muted.Sprint(filename))
}

// PrintFeature prints a feature with icon
func PrintFeature(icon, feature, description string) {
	fmt.Printf("  %s %s: %s\n", 
		Primary.Sprint(icon), 
		Bold.Sprint(feature), 
		Muted.Sprint(description))
}

// PrintCommand prints a command example
func PrintCommand(command string) {
	pterm.DefaultBox.WithBoxStyle(&pterm.Style{pterm.FgGreen}).Println(command)
}

// ShowSpinner creates and starts a spinner with message
func ShowSpinner(message string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " " + message
	s.Color("cyan")
	s.Start()
	return s
}

// PrintProjectStructure prints the project structure
func PrintProjectStructure(projectName string) {
	structure := fmt.Sprintf(`%s/
├── cmd/server/         # %s Application entry point
├── internal/
│   ├── handlers/       # %s HTTP handlers
│   ├── services/       # %s Business logic
│   ├── repositories/   # %s Data access
│   ├── models/         # %s Data models
│   └── middleware/     # %s HTTP middleware
├── pkg/
│   ├── database/       # %s Database connection
│   ├── config/         # %s Configuration
│   └── utils/          # %s Utilities
├── docker-compose.yml  # %s Docker configuration
├── Dockerfile          # %s Docker image
├── Makefile           # %s Build commands
└── .env.example       # %s Environment variables`, 
		Bold.Sprint(projectName),
		IconAPI, IconCode, IconGear, IconDatabase, IconPackage, IconGear,
		IconDatabase, IconGear, IconGear, IconDocker, IconDocker, IconBuild, IconGear)

	pterm.DefaultBox.WithTitle("Project Structure").WithTitleTopLeft().Println(structure)
}

// PrintNextSteps prints the next steps after project generation
func PrintNextSteps(projectName string) {
	fmt.Println()
	PrintHeader("Next Steps")
	
	steps := []struct {
		icon    string
		title   string
		command string
		desc    string
	}{
		{IconFolder, "Navigate to project", fmt.Sprintf("cd %s", projectName), "Enter the project directory"},
		{IconDocker, "Start with Docker", "docker-compose up --build", "Start API and database (recommended)"},
		{IconTest, "Run tests", "make test", "Execute test suite"},
		{IconBuild, "Build locally", "make build", "Compile the application"},
		{IconDoc, "View documentation", "cat README.md", "Read project documentation"},
	}

	for i, step := range steps {
		fmt.Printf("\n%d. %s %s\n", i+1, step.icon, Bold.Sprint(step.title))
		PrintCommand(step.command)
		fmt.Printf("   %s\n", Muted.Sprint(step.desc))
	}
}

// PrintDatabaseInfo prints database configuration info
func PrintDatabaseInfo(dbType, projectName string) {
	fmt.Println()
	PrintSubHeader("Database Configuration")
	
	switch dbType {
	case "postgres":
		PrintFeature(IconDatabase, "PostgreSQL", "Production-ready relational database")
		PrintFeature(IconDocker, "Docker Image", "postgres:15-alpine")
		PrintFeature(IconGear, "Connection", fmt.Sprintf("postgres://postgres:postgres@db:5432/%s", projectName))
	case "mysql":
		PrintFeature(IconDatabase, "MySQL", "Popular relational database")
		PrintFeature(IconDocker, "Docker Image", "mysql:8.0")
		PrintFeature(IconGear, "Connection", fmt.Sprintf("root:password@tcp(db:3306)/%s", projectName))
	case "sqlite":
		PrintFeature(IconDatabase, "SQLite", "Lightweight file-based database")
		PrintFeature(IconFile, "Database File", fmt.Sprintf("%s.db", projectName))
	case "mongodb":
		PrintFeature(IconDatabase, "MongoDB", "Document-oriented NoSQL database")
		PrintFeature(IconDocker, "Docker Image", "mongo:7.0")
		PrintFeature(IconGear, "Connection", fmt.Sprintf("mongodb://db:27017/%s", projectName))
	}
}

// PrintGeneratedFiles prints the list of generated files
func PrintGeneratedFiles(files []string) {
	fmt.Println()
	PrintSubHeader("Generated Files")
	
	for _, file := range files {
		PrintFileCreated(file)
	}
}

// ConfirmAction asks for user confirmation with styling
func ConfirmAction(message string) bool {
	result, _ := pterm.DefaultInteractiveConfirm.WithDefaultValue(false).Show(message)
	return result
}

// SelectOption shows a styled select menu
func SelectOption(message string, options []string) (string, error) {
	return pterm.DefaultInteractiveSelect.WithOptions(options).Show(message)
}

// TextInput shows a styled text input
func TextInput(message string, defaultValue ...string) (string, error) {
	prompt := pterm.DefaultInteractiveTextInput.WithTextStyle(&pterm.Style{pterm.FgCyan})
	
	if len(defaultValue) > 0 && defaultValue[0] != "" {
		prompt = prompt.WithDefaultValue(defaultValue[0])
	}
	
	return prompt.Show(message)
}

// PrintResourceSummary prints a summary of the generated resource
func PrintResourceSummary(resourceName string, fields []string) {
	fmt.Println()
	PrintHeader(fmt.Sprintf("Resource Summary: %s", resourceName))
	
	PrintFeature(IconPackage, "Resource Name", resourceName)
	PrintFeature(IconCode, "Generated Files", "Model, Handler, Service, Repository")
	
	if len(fields) > 0 {
		fmt.Printf("\n  %s %s:\n", IconAPI, Bold.Sprint("Fields"))
		for _, field := range fields {
			fmt.Printf("    • %s\n", field)
		}
	}
}

// ExitWithError prints an error and exits
func ExitWithError(message string) {
	PrintError(message)
	os.Exit(1)
}

// PrintSeparator prints a visual separator line
func PrintSeparator() {
	fmt.Println(Muted.Sprint("────────────────────────────────────────────────────────────────"))
}

// PrintKeyValue prints a key-value pair with formatting
func PrintKeyValue(key, value string) {
	fmt.Printf("  %s %s\n", Bold.Sprint(key+":"), value)
}