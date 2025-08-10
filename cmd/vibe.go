package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vibercode/cli/internal/vibe"
)

var vibeCmd = &cobra.Command{
	Use:   "vibe [mode]",
	Short: "🎨 Enter interactive vibe mode with AI chat and live preview",
	Long: `Vibe mode provides an interactive chat interface with Claude AI for building
and designing your Go APIs with real-time visual feedback.

Features:
• 💬 Conversational AI assistant with project context
• 🔥 Live preview of generated code and UI components
• ⚡ Real-time style editing with WebSocket updates
• 🎯 Context-aware suggestions based on your templates
• 🎨 Visual component builder with JSON schema

Available modes:
• vibe          - General mode (API + UI)
• vibe component - Component-focused mode (UI only)`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mode := "general"
		if len(args) > 0 {
			mode = args[0]
		}
		vibe.StartVibeMode(mode)
	},
}

// Note: vibeCmd is registered in root.go init function
