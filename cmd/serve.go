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
	"github.com/vibercode/cli/internal/api"
	"github.com/vibercode/cli/pkg/ui"
)

var (
	servePort string
	serveMode string
	serveHost string
	corsOrigins []string
	enableSwagger bool
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "🌐 Start HTTP API server for React Editor integration",
	Long: ui.Primary.Sprint("ViberCode API Server") + "\n\n" +
		"Start the HTTP API server to enable integration with the React Editor.\n" +
		"The server provides RESTful endpoints for schema management, code generation,\n" +
		"and project management operations.\n\n" +
		ui.Bold.Sprint("📡 Features:") + "\n" +
		"  " + ui.IconAPI + " RESTful API for all CLI operations\n" +
		"  " + ui.IconReact + " React Editor integration\n" +
		"  " + ui.IconCORS + " CORS support for React/Vite development servers (ports 5173-5180)\n" +
		"  " + ui.IconHealth + " Health checks and monitoring\n" +
		"  " + ui.IconLogs + " Request logging and error handling\n" +
		"  " + ui.IconSpeed + " Real-time project execution\n\n" +
		ui.Bold.Sprint("🔗 Endpoints:") + "\n" +
		"  GET    /api/v1/health           - Health check\n" +
		"  GET    /api/v1/metrics          - Server metrics\n" +
		"  GET    /api/v1/schema/list      - List schemas\n" +
		"  POST   /api/v1/schema/create    - Create schema\n" +
		"  GET    /api/v1/schema/{id}      - Get schema\n" +
		"  DELETE /api/v1/schema/{id}      - Delete schema\n" +
		"  POST   /api/v1/generate/api     - Generate API project\n" +
		"  POST   /api/v1/generate/resource - Generate resource\n" +
		"  GET    /api/v1/projects/list    - List projects\n" +
		"  POST   /api/v1/projects/{name}/run - Run project\n" +
		"  GET    /api/v1/projects/{name}/status - Project status\n" +
		"  POST   /api/v1/projects/{name}/stop - Stop project\n" +
		"  GET    /api/v1/projects/{name}/download - Download project\n",
	Example: ui.Dim.Sprint("  # Start server on default port 8080 with React/Vite support\n") +
		"  vibercode serve\n\n" +
		ui.Dim.Sprint("  # Start server on custom port with additional CORS origins\n") +
		"  vibercode serve --port 3001 --cors-origins \"http://localhost:3000,https://app.vibercode.com\"\n\n" +
		ui.Dim.Sprint("  # Start in development mode with debugging\n") +
		"  vibercode serve --mode dev --host 0.0.0.0\n\n" +
		ui.Dim.Sprint("  # Production mode with restricted CORS\n") +
		"  vibercode serve --mode production --cors-origins \"https://vibercode.com\"",
	RunE: runServeCommand,
}

func runServeCommand(cmd *cobra.Command, args []string) error {
	// Parse and validate port
	port, err := strconv.Atoi(servePort)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("invalid port number: %s (must be 1-65535)", servePort)
	}

	// Display startup banner
	ui.PrintSeparator()
	ui.PrintInfo("🚀 Starting ViberCode API Server...")
	
	// Configuration summary
	ui.PrintKeyValue("🌐 Host", serveHost)
	ui.PrintKeyValue("🔌 Port", servePort)
	ui.PrintKeyValue("🎯 Mode", serveMode)
	ui.PrintKeyValue("🔀 CORS Origins", fmt.Sprintf("%v", corsOrigins))
	ui.PrintKeyValue("📚 API Documentation", fmt.Sprintf("http://%s:%s/api/v1/docs", serveHost, servePort))
	ui.PrintSeparator()

	// Create server configuration
	config := &api.ServerConfig{
		Host:          serveHost,
		Port:          port,
		Mode:          serveMode,
		CORSOrigins:   corsOrigins,
		EnableSwagger: enableSwagger,
	}

	// Create and configure the server
	server, err := api.NewServer(config)
	if err != nil {
		return fmt.Errorf("failed to create API server: %w", err)
	}

	// Start server in a goroutine
	go func() {
		ui.PrintSuccess(fmt.Sprintf("✅ Server started successfully on http://%s:%d", serveHost, port))
		ui.PrintInfo("📖 API Documentation: http://" + serveHost + ":" + servePort + "/api/v1/docs")
		ui.PrintInfo("💚 Health Check: http://" + serveHost + ":" + servePort + "/api/v1/health")
		ui.PrintInfo("📊 Metrics: http://" + serveHost + ":" + servePort + "/api/v1/metrics")
		ui.PrintInfo("🛑 Press Ctrl+C to stop the server")
		
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			ui.PrintError("Failed to start server: " + err.Error())
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ui.PrintWarning("🛑 Shutting down server...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		ui.PrintError("⚠️  Server forced to shutdown: " + err.Error())
		return err
	}

	ui.PrintSuccess("✅ Server gracefully stopped")
	return nil
}

func init() {
	// Server configuration flags
	serveCmd.Flags().StringVarP(&servePort, "port", "p", "8080", "Port to run the API server on")
	serveCmd.Flags().StringVarP(&serveMode, "mode", "m", "api", "Server mode (api, dev, production)")
	serveCmd.Flags().StringVarP(&serveHost, "host", "H", "localhost", "Host to bind the server to")
	
	// CORS configuration
	serveCmd.Flags().StringSliceVar(&corsOrigins, "cors-origins", []string{
		"http://localhost:3000", "http://localhost:3001", "http://127.0.0.1:3000",
		"http://localhost:5173", "http://localhost:5174", "http://localhost:5175", 
		"http://localhost:5176", "http://localhost:5177", "http://localhost:5178", 
		"http://localhost:5179", "http://localhost:5180",
	}, "Allowed CORS origins (comma-separated)")
	
	// Development options
	serveCmd.Flags().BoolVar(&enableSwagger, "swagger", true, "Enable Swagger API documentation")
	
	// Add descriptions for flags
	serveCmd.Flags().Lookup("port").Usage = "Port number for the HTTP server (1-65535)"
	serveCmd.Flags().Lookup("mode").Usage = "Server mode: 'api' (default), 'dev' (debug), 'production' (optimized)"
	serveCmd.Flags().Lookup("host").Usage = "Host address to bind to (localhost, 0.0.0.0, or specific IP)"
	serveCmd.Flags().Lookup("cors-origins").Usage = "Comma-separated list of allowed CORS origins (includes React/Vite dev servers by default)"
	serveCmd.Flags().Lookup("swagger").Usage = "Enable/disable Swagger API documentation endpoint"
}