package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/vibercode/cli/internal/websocket"
	"github.com/vibercode/cli/pkg/ui"
)

var (
	wsPort string
	wsHost string
)

var wsCmd = &cobra.Command{
	Use:   "ws",
	Short: "🌐 Start WebSocket server for real-time React Editor integration",
	Long: ui.Primary.Sprint("ViberCode WebSocket Server") + "\n\n" +
		"Start the WebSocket server to enable real-time integration with the React Editor.\n" +
		"The server provides real-time communication for view updates and code generation.\n\n" +
		ui.Bold.Sprint("📡 Features:") + "\n" +
		"  " + ui.IconWebSocket + " Real-time WebSocket communication\n" +
		"  " + ui.IconReact + " React Editor view updates\n" +
		"  " + ui.IconCode + " Live code generation\n" +
		"  " + ui.IconHealth + " Connection monitoring\n" +
		"  " + ui.IconLogs + " Message logging and debugging\n\n" +
		ui.Bold.Sprint("🔗 Messages:") + "\n" +
		"  view_update       - Receive component updates from editor\n" +
		"  generate_request  - Handle code generation requests\n" +
		"  generate_response - Send generation results back\n" +
		"  error            - Error notifications\n",
	Example: ui.Dim.Sprint("  # Start WebSocket server on default port 3001\n") +
		"  vibercode ws\n\n" +
		ui.Dim.Sprint("  # Start server on custom port\n") +
		"  vibercode ws --port 3002\n\n" +
		ui.Dim.Sprint("  # Bind to all interfaces\n") +
		"  vibercode ws --host 0.0.0.0\n\n" +
		ui.Dim.Sprint("  # Run both HTTP API and WebSocket servers\n") +
		"  vibercode serve --port 8080 & vibercode ws --port 3001",
	RunE: runWSCommand,
}

func runWSCommand(cmd *cobra.Command, args []string) error {
	// Parse and validate port
	port, err := strconv.Atoi(wsPort)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("invalid port number: %s (must be 1-65535)", wsPort)
	}

	// Display startup banner
	ui.PrintSeparator()
	ui.PrintInfo("🚀 Starting ViberCode WebSocket Server...")

	// Configuration summary
	ui.PrintKeyValue("🌐 Host", wsHost)
	ui.PrintKeyValue("🔌 Port", wsPort)
	ui.PrintKeyValue("📡 WebSocket URL", fmt.Sprintf("ws://%s:%s", wsHost, wsPort))
	ui.PrintSeparator()

	// Create WebSocket server
	server := websocket.NewServer(wsHost, port)

	// Create HTTP server for WebSocket endpoint
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", server.HandleWebSocket)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", wsHost, port),
		Handler: mux,
	}

	// Start server in a goroutine
	go func() {
		ui.PrintSuccess(fmt.Sprintf("✅ WebSocket server started on ws://%s:%d", wsHost, port))
		ui.PrintInfo("📖 Connect your React Editor to this WebSocket URL")
		ui.PrintInfo("🛑 Press Ctrl+C to stop the server")

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ui.PrintError("Failed to start WebSocket server: " + err.Error())
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ui.PrintWarning("🛑 Shutting down WebSocket server...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := httpServer.Shutdown(ctx); err != nil {
		ui.PrintError("⚠️  WebSocket server forced to shutdown: " + err.Error())
		return err
	}

	ui.PrintSuccess("✅ WebSocket server gracefully stopped")
	return nil
}

func init() {
	// Server configuration flags
	wsCmd.Flags().StringVarP(&wsPort, "port", "p", "3001", "Port to run the WebSocket server on")
	wsCmd.Flags().StringVarP(&wsHost, "host", "H", "localhost", "Host to bind the server to")

	// Add descriptions for flags
	wsCmd.Flags().Lookup("port").Usage = "Port number for the WebSocket server (1-65535)"
	wsCmd.Flags().Lookup("host").Usage = "Host address to bind to (localhost, 0.0.0.0, or specific IP)"
}
