package vibe

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/vibercode/cli/internal/websocket"
	"github.com/vibercode/cli/pkg/ui"
)

// VibeServices holds all the services for vibe mode
type VibeServices struct {
	wsServer    *websocket.Server
	httpServer  *http.Server
	editorCmd   *exec.Cmd
	chatManager *ChatManager
	cleanup     []func()
	ctx         context.Context
	cancel      context.CancelFunc
}

// StartVibeMode initializes and starts the interactive vibe mode with React editor
func StartVibeMode(mode string) {
	// Initialize logger for vibe mode
	InitVibeLogger(true) // true = chat mode, less verbose
	defer CloseVibeLogger()

	// Create context for cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// Display welcome banner
	showVibeBanner(mode)

	// Initialize services
	services := &VibeServices{
		cleanup: make([]func(), 0),
		ctx:     ctx,
		cancel:  cancel,
	}

	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		ui.PrintWarning("\n\n⚡ Shutting down VibeCode...")
		// Cancel context to stop all goroutines
		services.cancel()
		// Stop chat manager if it exists
		if services.chatManager != nil {
			services.chatManager.cancel()
		}
		services.shutdown()
		os.Exit(0)
	}()

	// Start all services with timeout
	if err := services.startAllWithTimeout(mode, 30*time.Second); err != nil {
		ui.PrintError("Failed to start VibeCode services: " + err.Error())
		ui.PrintInfo("💡 Try running: ./debug-paths.sh to find your React editor")
		return
	}

	// Create chat manager with the specified mode
	chatManager := NewChatManager(mode)
	if chatManager == nil {
		ui.PrintError("❌ Failed to initialize ChatManager")
		return
	}

	services.chatManager = chatManager

	// Start chat interface (this blocks)
	ui.PrintInfo("💬 Starting chat interface...")
	ui.PrintInfo("🛑 Press Ctrl+C to stop all services")

	// Start chat with proper context handling
	go func() {
		defer func() {
			if r := recover(); r != nil {
				ui.PrintError("Chat panic recovered: " + fmt.Sprintf("%v", r))
				services.cancel()
			}
		}()

		select {
		case <-ctx.Done():
			ui.PrintInfo("🛑 Chat startup cancelled")
			return
		case <-time.After(2 * time.Second):
			// Start chat with context cancellation support
			chatManager.StartChat()
		}
	}()

	// Keep main thread alive with timeout fallback
	select {
	case <-ctx.Done():
		ui.PrintInfo("🛑 Context cancelled, shutting down...")
	case <-time.After(5 * time.Minute):
		ui.PrintWarning("⏰ Timeout reached, shutting down...")
		services.cancel()
	}
}

// startAllWithTimeout starts all required services for vibe mode with a timeout
func (s *VibeServices) startAllWithTimeout(mode string, timeout time.Duration) error {
	ui.PrintInfo("🚀 Starting VibeCode services...")

	// Create a timeout context
	ctx, cancel := context.WithTimeout(s.ctx, timeout)
	defer cancel()

	// Channel to receive startup result
	result := make(chan error, 1)

	go func() {
		result <- s.startAll(mode)
	}()

	select {
	case err := <-result:
		return err
	case <-ctx.Done():
		return fmt.Errorf("startup timeout after %v", timeout)
	}
}

// startAll starts all required services for vibe mode
func (s *VibeServices) startAll(mode string) error {
	// 1. Start WebSocket server
	if err := s.startWebSocketServer(); err != nil {
		return fmt.Errorf("failed to start WebSocket server: %w", err)
	}

	// 2. Find and start React editor (non-blocking if not found)
	s.startReactEditorNonBlocking()

	// 3. Wait a moment for services to start
	time.Sleep(1 * time.Second)

	// 4. Open browser
	s.openBrowser()

	ui.PrintSuccess("✅ VibeCode is ready!")
	ui.PrintInfo("🌐 WebSocket: ws://localhost:3001/ws")
	ui.PrintInfo("🎨 Editor: http://localhost:5173 (if React editor started)")
	ui.PrintSeparator()

	return nil
}

// startWebSocketServer starts the WebSocket server for real-time communication
func (s *VibeServices) startWebSocketServer() error {
	ui.PrintInfo("📡 Starting WebSocket server on port 3001...")

	s.wsServer = websocket.NewServer("localhost", 3001)

	// Create HTTP server for WebSocket
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.wsServer.HandleWebSocket)

	s.httpServer = &http.Server{
		Addr:    ":3001",
		Handler: mux,
	}

	// Start server in goroutine
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ui.PrintError("WebSocket server error: " + err.Error())
		}
	}()

	// Add cleanup
	s.cleanup = append(s.cleanup, func() {
		ui.PrintInfo("🔌 Stopping WebSocket server...")
		if s.httpServer != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			s.httpServer.Shutdown(ctx)
		}
	})

	return nil
}

// startReactEditorNonBlocking tries to start the React editor without blocking
func (s *VibeServices) startReactEditorNonBlocking() {
	ui.PrintInfo("🎨 Looking for React Editor...")

	// Find editor directory without prompting
	editorPath, err := s.findEditorPathQuiet()
	if err != nil {
		ui.PrintWarning("⚠️  React Editor not found automatically")
		ui.PrintInfo("💡 You can start it manually with:")
		ui.PrintInfo("   cd ../../../vibercode/editor && pnpm dev")
		ui.PrintInfo("💡 Or run: ./debug-paths.sh to find the correct path")
		return
	}

	ui.PrintSuccess("📂 Found editor at: " + editorPath)

	// Check dependencies in background
	go func() {
		nodeModulesPath := filepath.Join(editorPath, "node_modules")
		if _, err := os.Stat(nodeModulesPath); os.IsNotExist(err) {
			ui.PrintInfo("📦 Installing dependencies in background...")
			if err := s.installDependencies(editorPath); err != nil {
				ui.PrintWarning("⚠️  Failed to install dependencies: " + err.Error())
				ui.PrintInfo("💡 Run manually: cd " + editorPath + " && pnpm install")
				return
			}
		}

		// Start development server
		cmd := s.createEditorCommand(editorPath)
		if err := cmd.Start(); err != nil {
			ui.PrintWarning("⚠️  Failed to start editor: " + err.Error())
			ui.PrintInfo("💡 Run manually: cd " + editorPath + " && pnpm dev")
			return
		}

		s.editorCmd = cmd
		ui.PrintSuccess("✅ React Editor started successfully!")

		// Add cleanup
		s.cleanup = append(s.cleanup, func() {
			if s.editorCmd != nil && s.editorCmd.Process != nil {
				ui.PrintInfo("🎨 Stopping React Editor...")
				s.editorCmd.Process.Kill()
				s.editorCmd.Wait()
			}
		})
	}()
}

// findEditorPathQuiet finds the editor path without user interaction
func (s *VibeServices) findEditorPathQuiet() (string, error) {
	// Get current working directory for reference
	cwd, _ := os.Getwd()

	// Possible paths to check (based on actual project structure)
	possiblePaths := []string{
		"../../../vibercode/editor",                // ✅ FOUND: Correct path for this project structure
		"../vibercode/editor",                      // From vibercode-cli-go to vibercode/editor
		"../../vibercode/editor",                   // From nested structure
		"./vibercode/editor",                       // Current directory
		"vibercode/editor",                         // Direct subdirectory
		"../editor",                                // From vibercode-cli-go to editor
		"editor",                                   // Current directory
		"../../blackbox-cli/vibercode/editor",      // Absolute structure
		"../blackbox-cli/vibercode/editor",         // Alternative structure
		"/Users/jaambee/Projects/vibercode/editor", // Direct absolute path as fallback
	}

	for _, path := range possiblePaths {
		if absPath, err := filepath.Abs(path); err == nil {
			if s.isValidEditorPath(absPath) {
				return absPath, nil
			}
		}
	}

	// Try extensive search without user interaction
	return s.findEditorPathExtensiveQuiet(cwd)
}

// findEditorPathExtensiveQuiet performs a more thorough search without prompts
func (s *VibeServices) findEditorPathExtensiveQuiet(startDir string) (string, error) {
	// Walk up directories to find the project structure
	currentDir := startDir
	for i := 0; i < 5; i++ { // Limit to 5 levels up
		// Check common patterns from this level
		testPaths := []string{
			filepath.Join(currentDir, "vibercode", "editor"),
			filepath.Join(currentDir, "editor"),
			filepath.Join(currentDir, "..", "vibercode", "editor"),
		}

		for _, testPath := range testPaths {
			if absPath, err := filepath.Abs(testPath); err == nil {
				if s.isValidEditorPath(absPath) {
					return absPath, nil
				}
			}
		}

		// Move up one directory
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break // Reached filesystem root
		}
		currentDir = parentDir
	}

	return "", fmt.Errorf("could not find React editor directory")
}

// isValidEditorPath checks if a path contains a valid React editor
func (s *VibeServices) isValidEditorPath(path string) bool {
	// Check for package.json
	packagePath := filepath.Join(path, "package.json")
	if _, err := os.Stat(packagePath); os.IsNotExist(err) {
		return false
	}

	// Check for src directory
	srcPath := filepath.Join(path, "src")
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return false
	}

	return true
}

// installDependencies installs npm/pnpm dependencies
func (s *VibeServices) installDependencies(editorPath string) error {
	// Try pnpm first, then npm
	packageManagers := []string{"pnpm", "npm"}

	for _, pm := range packageManagers {
		if _, err := exec.LookPath(pm); err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()

			cmd := exec.CommandContext(ctx, pm, "install")
			cmd.Dir = editorPath

			if err := cmd.Run(); err == nil {
				return nil
			}
		}
	}

	return fmt.Errorf("no package manager found (tried pnpm, npm)")
}

// createEditorCommand creates the appropriate command to start the editor
func (s *VibeServices) createEditorCommand(editorPath string) *exec.Cmd {
	// Try pnpm first, then npm
	packageManagers := []struct {
		cmd  string
		args []string
	}{
		{"pnpm", []string{"dev"}},
		{"npm", []string{"run", "dev"}},
		{"yarn", []string{"dev"}},
	}

	for _, pm := range packageManagers {
		if _, err := exec.LookPath(pm.cmd); err == nil {
			cmd := exec.Command(pm.cmd, pm.args...)
			cmd.Dir = editorPath
			// Don't capture output to avoid blocking
			return cmd
		}
	}

	// Fallback to npm
	cmd := exec.Command("npm", "run", "dev")
	cmd.Dir = editorPath
	return cmd
}

// openBrowser opens the default browser with the editor URL
func (s *VibeServices) openBrowser() {
	editorURL := "http://localhost:5173"

	ui.PrintInfo("🌐 Browser will open in 3 seconds...")

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("open", editorURL)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", editorURL)
	case "linux":
		cmd = exec.Command("xdg-open", editorURL)
	default:
		ui.PrintInfo("📋 Please open: " + editorURL)
		return
	}

	go func() {
		time.Sleep(3 * time.Second) // Wait for server to start
		if err := cmd.Run(); err != nil {
			ui.PrintInfo("📋 Please manually open: " + editorURL)
		}
	}()
}

// shutdown gracefully stops all services
func (s *VibeServices) shutdown() {
	ui.PrintInfo("🛑 Shutting down services...")

	// Cancel context first
	if s.cancel != nil {
		s.cancel()
	}

	// Run all cleanup functions in reverse order
	for i := len(s.cleanup) - 1; i >= 0; i-- {
		s.cleanup[i]()
	}

	ui.PrintSuccess("✅ Shutdown complete")
}

func showVibeBanner(mode string) {
	modeDescription := "General Mode - API + UI + Editor"
	if mode == "component" {
		modeDescription = "Component Mode - UI Editor Focus"
	}

	banner := `
    ⚡ ⚡ ⚡ ⚡ ⚡ ⚡ ⚡ ⚡ ⚡ ⚡ ⚡ ⚡ ⚡ ⚡ ⚡
    ⚡                                           ⚡
    ⚡     🎨 Welcome to VibeCode Full Mode     ⚡
    ⚡                                           ⚡
    ⚡ ⚡ ⚡ ⚡ ⚡ ⚡ ⚡ ⚡ ⚡ ⚡ ⚡ ⚡ ⚡ ⚡ ⚡

    🤖 AI-Powered Go API Builder
    🎨 React Visual Editor
    🔥 Real-time UI Preview
    ⚡ Live Style Editor
    💬 Conversational Interface
    🎯 Context-Aware Suggestions
    🌐 WebSocket Integration

    `

	fmt.Print(ui.Primary.Sprint(banner))
	fmt.Println(ui.Bold.Sprint("    Mode: ") + ui.Info.Sprint(modeDescription))
	fmt.Println("")
	fmt.Println(ui.Bold.Sprint("    Services:"))
	fmt.Println(ui.Info.Sprint("    • 📡 WebSocket Server (localhost:3001)"))
	fmt.Println(ui.Info.Sprint("    • 🎨 React Editor (localhost:5173)"))
	fmt.Println(ui.Info.Sprint("    • 💬 AI Chat Interface"))
	fmt.Println(ui.Info.Sprint("    • 🔄 Live Component Sync"))

	if mode == "component" {
		fmt.Println("")
		fmt.Println(ui.Bold.Sprint("    Focus:"))
		fmt.Println(ui.Info.Sprint("    • Component-focused AI chat"))
		fmt.Println(ui.Info.Sprint("    • Real-time component preview"))
		fmt.Println(ui.Info.Sprint("    • Live component editing"))
		fmt.Println(ui.Info.Sprint("    • Interactive UI builder"))
	} else {
		fmt.Println("")
		fmt.Println(ui.Bold.Sprint("    Features:"))
		fmt.Println(ui.Info.Sprint("    • Chat with Viber AI about your Go API"))
		fmt.Println(ui.Info.Sprint("    • Visual component editor"))
		fmt.Println(ui.Info.Sprint("    • Real-time UI component editing"))
		fmt.Println(ui.Info.Sprint("    • Template-based code generation"))
		fmt.Println(ui.Info.Sprint("    • Interactive style customization"))
	}
	fmt.Println("")
}
