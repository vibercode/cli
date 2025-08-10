package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vibercode/cli/internal/mcp"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "🔌 Start MCP (Model Context Protocol) server for AI agent integration",
	Long: `MCP server provides tools for AI agents to interact with ViberCode.

Features:
• 🤖 AI agent integration via Model Context Protocol
• 🎨 Live component editing and updates
• 🔄 Real-time WebSocket communication
• 💬 Chat integration with Viber AI
• 📊 View state management
• ⚡ Code generation and project management

Available tools:
• vibe_start        - Start vibe mode with live preview
• component_update  - Update component properties in real-time
• view_state_get    - Get current view state and components
• chat_send         - Send message to Viber AI assistant
• generate_code     - Generate Go API code from schema
• project_status    - Get project and server status`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return mcp.StartMCPServer()
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}
