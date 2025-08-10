package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vibercode/cli/pkg/ui"
)

var rootCmd = &cobra.Command{
	Use:   "vibercode",
	Short: "🚀 Vibercode CLI - Generate Go web APIs with clean architecture",
	Long: ui.Primary.Sprint("Vibercode CLI") + " is a command-line tool for generating Go web APIs\n" +
		"with complete CRUD operations, following clean architecture principles.\n\n" +
		ui.Bold.Sprint("🌟 Features:") + "\n" +
		"  " + ui.IconRocket + " Generate complete Go API projects\n" +
		"  " + ui.IconCode + " Create CRUD resources with models, handlers, services, and repositories\n" +
		"  " + ui.IconCLI + " Interactive prompts for easy configuration\n" +
		"  " + ui.IconDatabase + " Built-in support for GORM and Gin framework\n" +
		"  " + ui.IconPackage + " Clean architecture with separation of concerns\n" +
		"  " + ui.IconDocker + " Docker-ready with auto-generated docker-compose\n" +
		"  " + ui.IconSpeed + " Zero configuration - works out of the box\n",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Name() != "help" && cmd.Name() != "completion" {
			ui.ShowBanner()
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(schemaCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(wsCmd)
	rootCmd.AddCommand(templateCmd)
	rootCmd.AddCommand(mcpCmd)
	rootCmd.AddCommand(vibeCmd)
}
