package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/vibercode/cli/internal/templates"
	"github.com/vibercode/cli/pkg/ui"
)

var (
	templateDir    string
	templateOutput string
	templateVars   []string
	showAll        bool
	categoryFilter string
	typeFilter     string
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "🎨 Manage and generate from templates",
	Long: ui.Primary.Sprint("ViberCode Template Management") + "\n\n" +
		"Manage templates for both backend (Go) and frontend (React/Vue/Angular) code generation.\n" +
		"Templates allow you to generate consistent, high-quality code across different technologies.\n\n" +
		ui.Bold.Sprint("📦 Template Categories:") + "\n" +
		"  " + ui.IconAPI + " fullstack-ai - Complete full-stack applications with AI integrations\n" +
		"  " + ui.IconReact + " frontend - React, Vue, Angular components and applications\n" +
		"  " + ui.IconCode + " backend - API servers, microservices, serverless functions\n" +
		"  " + ui.IconDatabase + " api - RESTful APIs with Clean Architecture\n" +
		"  " + ui.IconPackage + " resource - CRUD resources with validation and relationships\n\n" +
		ui.Bold.Sprint("🔧 Template Types:") + "\n" +
		"  • go - Go language templates\n" +
		"  • react - React TypeScript templates\n" +
		"  • vue - Vue.js templates\n" +
		"  • angular - Angular TypeScript templates\n" +
		"  • nextjs - Next.js templates\n" +
		"  • svelte - Svelte templates",
	Example: ui.Dim.Sprint("  # List all available templates\n") +
		"  vibercode template list\n\n" +
		ui.Dim.Sprint("  # List only React templates\n") +
		"  vibercode template list --type react\n\n" +
		ui.Dim.Sprint("  # Generate from a template\n") +
		"  vibercode template generate go-api-resource --output ./my-project\n\n" +
		ui.Dim.Sprint("  # Generate with custom variables\n") +
		"  vibercode template generate react-crud-component --output ./components --var schema=user.json\n\n" +
		ui.Dim.Sprint("  # Show template details\n") +
		"  vibercode template show react-crud-component",
}

var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "📋 List available templates",
	Long:  "Display all available templates with their metadata and descriptions.",
	RunE:  runTemplateListCommand,
}

var templateShowCmd = &cobra.Command{
	Use:   "show [template-id]",
	Short: "👁️  Show template details",
	Long:  "Display detailed information about a specific template including variables and files.",
	Args:  cobra.ExactArgs(1),
	RunE:  runTemplateShowCommand,
}

var templateGenerateCmd = &cobra.Command{
	Use:   "generate [template-id]",
	Short: "🚀 Generate code from template",
	Long:  "Generate code files from a template with optional variable substitution.",
	Args:  cobra.ExactArgs(1),
	RunE:  runTemplateGenerateCommand,
}

var templateValidateCmd = &cobra.Command{
	Use:   "validate [template-id]",
	Short: "✅ Validate template",
	Long:  "Validate a template's structure, dependencies, and syntax.",
	Args:  cobra.ExactArgs(1),
	RunE:  runTemplateValidateCommand,
}

func runTemplateListCommand(cmd *cobra.Command, args []string) error {
	registry := templates.NewTemplateRegistry(templateDir)
	if err := registry.LoadTemplates(); err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	allTemplates := registry.ListTemplates()
	if len(allTemplates) == 0 {
		ui.PrintWarning("No templates found")
		return nil
	}

	// Filter templates
	var filteredTemplates []*templates.TemplateMetadata
	for _, tmpl := range allTemplates {
		if categoryFilter != "" && string(tmpl.Category) != categoryFilter {
			continue
		}
		if typeFilter != "" && string(tmpl.Type) != typeFilter {
			continue
		}
		if !showAll && tmpl.Deprecated {
			continue
		}
		filteredTemplates = append(filteredTemplates, tmpl)
	}

	// Sort by category, then by name
	sort.Slice(filteredTemplates, func(i, j int) bool {
		if filteredTemplates[i].Category != filteredTemplates[j].Category {
			return filteredTemplates[i].Category < filteredTemplates[j].Category
		}
		return filteredTemplates[i].Name < filteredTemplates[j].Name
	})

	// Display templates
	ui.PrintSeparator()
	ui.PrintInfo(fmt.Sprintf("Found %d templates", len(filteredTemplates)))
	ui.PrintSeparator()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "ID\tNAME\tCATEGORY\tTYPE\tVERSION\tDESCRIPTION\n")
	fmt.Fprintf(w, "──\t────\t────────\t────\t───────\t───────────\n")

	var currentCategory templates.TemplateCategory
	for _, tmpl := range filteredTemplates {
		if tmpl.Category != currentCategory {
			if currentCategory != "" {
				fmt.Fprintf(w, "\n")
			}
			currentCategory = tmpl.Category
		}

		status := ""
		if tmpl.Deprecated {
			status = ui.Warning.Sprint(" [DEPRECATED]")
		} else if tmpl.Experimental {
			status = ui.Info.Sprint(" [EXPERIMENTAL]")
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s%s\n",
			ui.Primary.Sprint(tmpl.ID),
			tmpl.DisplayName,
			string(tmpl.Category),
			string(tmpl.Type),
			tmpl.Version,
			tmpl.Description,
			status,
		)
	}

	w.Flush()
	ui.PrintSeparator()

	// Show usage hint
	ui.PrintInfo("Use 'vibercode template show <template-id>' for detailed information")
	ui.PrintInfo("Use 'vibercode template generate <template-id>' to generate code")

	return nil
}

func runTemplateShowCommand(cmd *cobra.Command, args []string) error {
	templateID := args[0]

	registry := templates.NewTemplateRegistry(templateDir)
	if err := registry.LoadTemplates(); err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	tmpl, err := registry.GetTemplate(templateID)
	if err != nil {
		return fmt.Errorf("template not found: %s", templateID)
	}

	// Display template details
	ui.PrintSeparator()
	ui.PrintInfo(fmt.Sprintf("Template: %s", tmpl.DisplayName))
	ui.PrintSeparator()

	ui.PrintKeyValue("ID", tmpl.ID)
	ui.PrintKeyValue("Name", tmpl.Name)
	ui.PrintKeyValue("Category", string(tmpl.Category))
	ui.PrintKeyValue("Type", string(tmpl.Type))
	ui.PrintKeyValue("Version", tmpl.Version)
	ui.PrintKeyValue("Author", tmpl.Author)
	ui.PrintKeyValue("Framework", tmpl.Framework)
	ui.PrintKeyValue("Language", tmpl.Language)

	if tmpl.Deprecated {
		ui.PrintWarning("⚠️  This template is deprecated")
	}
	if tmpl.Experimental {
		ui.PrintWarning("🧪 This template is experimental")
	}

	fmt.Println()
	ui.PrintSubHeader("Description")
	fmt.Printf("  %s\n", tmpl.Description)

	if len(tmpl.Tags) > 0 {
		fmt.Println()
		ui.PrintSubHeader("Tags")
		fmt.Printf("  %s\n", strings.Join(tmpl.Tags, ", "))
	}

	if len(tmpl.Dependencies) > 0 {
		fmt.Println()
		ui.PrintSubHeader("Dependencies")
		for _, dep := range tmpl.Dependencies {
			fmt.Printf("  • %s\n", dep)
		}
	}

	if len(tmpl.Requirements) > 0 {
		fmt.Println()
		ui.PrintSubHeader("Requirements")
		for _, req := range tmpl.Requirements {
			fmt.Printf("  • %s\n", req)
		}
	}

	if len(tmpl.Variables) > 0 {
		fmt.Println()
		ui.PrintSubHeader("Variables")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "  NAME\tTYPE\tREQUIRED\tDEFAULT\tDESCRIPTION\n")
		fmt.Fprintf(w, "  ────\t────\t────────\t───────\t───────────\n")
		for _, variable := range tmpl.Variables {
			required := "No"
			if variable.Required {
				required = "Yes"
			}
			defaultVal := "-"
			if variable.Default != nil {
				defaultVal = fmt.Sprintf("%v", variable.Default)
			}
			fmt.Fprintf(w, "  %s\t%s\t%s\t%s\t%s\n",
				variable.Name,
				variable.Type,
				required,
				defaultVal,
				variable.Description,
			)
		}
		w.Flush()
	}

	if len(tmpl.Files) > 0 {
		fmt.Println()
		ui.PrintSubHeader("Generated Files")
		for _, file := range tmpl.Files {
			fmt.Printf("  📄 %s\n", file.Path)
			if file.Condition != "" {
				fmt.Printf("     Condition: %s\n", file.Condition)
			}
		}
	}

	if len(tmpl.Examples) > 0 {
		fmt.Println()
		ui.PrintSubHeader("Examples")
		for _, example := range tmpl.Examples {
			fmt.Printf("  %s\n", example)
		}
	}

	ui.PrintSeparator()
	return nil
}

func runTemplateGenerateCommand(cmd *cobra.Command, args []string) error {
	templateID := args[0]

	if templateOutput == "" {
		templateOutput = "."
	}

	registry := templates.NewTemplateRegistry(templateDir)
	if err := registry.LoadTemplates(); err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	// Parse variables
	variables := make(map[string]interface{})
	for _, varStr := range templateVars {
		parts := strings.SplitN(varStr, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid variable format: %s (expected key=value)", varStr)
		}
		key := parts[0]
		value := parts[1]

		// Try to parse as JSON first, then fall back to string
		var parsedValue interface{}
		if err := json.Unmarshal([]byte(value), &parsedValue); err != nil {
			// If it's a file path ending in .json, try to load it
			if strings.HasSuffix(value, ".json") {
				if data, err := os.ReadFile(value); err == nil {
					if err := json.Unmarshal(data, &parsedValue); err == nil {
						variables[key] = parsedValue
						continue
					}
				}
			}
			// Fall back to string value
			parsedValue = value
		}
		variables[key] = parsedValue
	}

	// Generate from template
	ui.PrintInfo(fmt.Sprintf("Generating from template: %s", templateID))
	ui.PrintInfo(fmt.Sprintf("Output directory: %s", templateOutput))

	if err := registry.GenerateFromTemplate(templateID, templateOutput, variables); err != nil {
		return fmt.Errorf("failed to generate from template: %w", err)
	}

	ui.PrintSuccess("Template generation completed successfully!")
	return nil
}

func runTemplateValidateCommand(cmd *cobra.Command, args []string) error {
	templateID := args[0]

	registry := templates.NewTemplateRegistry(templateDir)
	if err := registry.LoadTemplates(); err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	ui.PrintInfo(fmt.Sprintf("Validating template: %s", templateID))

	if err := registry.ValidateTemplate(templateID); err != nil {
		ui.PrintError(fmt.Sprintf("Template validation failed: %v", err))
		return err
	}

	ui.PrintSuccess("Template validation passed!")
	return nil
}

func init() {
	// Template management flags
	templateCmd.PersistentFlags().StringVar(&templateDir, "template-dir", getDefaultTemplateDir(), "Directory containing templates")

	// List command flags
	templateListCmd.Flags().BoolVar(&showAll, "all", false, "Show all templates including deprecated")
	templateListCmd.Flags().StringVar(&categoryFilter, "category", "", "Filter by category (fullstack-ai, frontend, backend, api, resource)")
	templateListCmd.Flags().StringVar(&typeFilter, "type", "", "Filter by type (go, react, vue, angular, nextjs, svelte)")

	// Generate command flags
	templateGenerateCmd.Flags().StringVarP(&templateOutput, "output", "o", "", "Output directory for generated files")
	templateGenerateCmd.Flags().StringSliceVarP(&templateVars, "var", "v", []string{}, "Template variables (key=value or key=file.json)")

	// Add subcommands
	templateCmd.AddCommand(templateListCmd)
	templateCmd.AddCommand(templateShowCmd)
	templateCmd.AddCommand(templateGenerateCmd)
	templateCmd.AddCommand(templateValidateCmd)
}

func getDefaultTemplateDir() string {
	// Try to find templates directory relative to executable
	if exe, err := os.Executable(); err == nil {
		if dir := filepath.Join(filepath.Dir(exe), "templates"); dir != "" {
			if _, err := os.Stat(dir); err == nil {
				return dir
			}
		}
	}

	// Try current directory
	if dir := "./templates"; dir != "" {
		if _, err := os.Stat(dir); err == nil {
			return dir
		}
	}

	// Default to home directory
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".vibercode", "templates")
	}

	return "./templates"
}