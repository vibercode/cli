# Task 12: Live Reload Development

## Overview
Implement a comprehensive live reload development system that provides hot reloading for Go applications, real-time code compilation, automatic browser refresh, and integrated development server with file watching capabilities.

## Objectives
- Create development server with hot reload for Go applications
- Implement file watching and automatic recompilation
- Provide browser auto-refresh for web interfaces
- Build development proxy with API and frontend integration
- Add real-time error reporting and debugging
- Create development dashboard and monitoring

## Implementation Details

### Command Structure
```bash
# Development server commands
vibercode dev                           # Start development server
vibercode dev --port 3000              # Custom port
vibercode dev --proxy-api 8080         # Proxy to API server
vibercode dev --watch "**/*.go"        # Custom watch patterns
vibercode dev --verbose                # Verbose logging

# Development configuration
vibercode dev init                      # Initialize dev config
vibercode dev config                    # Show current config
vibercode dev clean                     # Clean build cache

# Development tools
vibercode dev logs                      # Show development logs
vibercode dev status                    # Show server status
vibercode dev restart                   # Restart development server
```

### Development Server Architecture

#### Core Development Server
```go
package dev

import (
    "context"
    "fmt"
    "net/http"
    "path/filepath"
    "time"
    
    "github.com/fsnotify/fsnotify"
    "github.com/gorilla/websocket"
)

type DevServer struct {
    config     *DevConfig
    watcher    *fsnotify.Watcher
    compiler   *GoCompiler
    proxy      *APIProxy
    wsHub      *WebSocketHub
    dashboard  *Dashboard
    logger     Logger
}

type DevConfig struct {
    Port         int      `yaml:"port" json:"port"`
    APIPort      int      `yaml:"api_port" json:"api_port"`
    WatchPaths   []string `yaml:"watch_paths" json:"watch_paths"`
    IgnorePaths  []string `yaml:"ignore_paths" json:"ignore_paths"`
    BuildDir     string   `yaml:"build_dir" json:"build_dir"`
    AutoReload   bool     `yaml:"auto_reload" json:"auto_reload"`
    ProxyEnabled bool     `yaml:"proxy_enabled" json:"proxy_enabled"`
    Dashboard    bool     `yaml:"dashboard" json:"dashboard"`
}

func NewDevServer(config *DevConfig) *DevServer {
    return &DevServer{
        config:    config,
        watcher:   setupFileWatcher(config.WatchPaths, config.IgnorePaths),
        compiler:  NewGoCompiler(config.BuildDir),
        proxy:     NewAPIProxy(config.APIPort),
        wsHub:     NewWebSocketHub(),
        dashboard: NewDashboard(),
        logger:    NewLogger("dev-server"),
    }
}

func (s *DevServer) Start(ctx context.Context) error {
    // Start file watcher
    go s.watchFiles(ctx)
    
    // Start WebSocket hub
    go s.wsHub.Run(ctx)
    
    // Start API proxy if enabled
    if s.config.ProxyEnabled {
        go s.proxy.Start(ctx)
    }
    
    // Setup HTTP handlers
    mux := s.setupHandlers()
    
    server := &http.Server{
        Addr:    fmt.Sprintf(":%d", s.config.Port),
        Handler: mux,
    }
    
    s.logger.Info("Development server starting on port %d", s.config.Port)
    return server.ListenAndServe()
}
```

#### File Watcher Implementation
```go
type FileWatcher struct {
    watcher     *fsnotify.Watcher
    debouncer   *Debouncer
    onChange    func([]string)
    watchPaths  []string
    ignorePaths []string
}

func (fw *FileWatcher) watchFiles(ctx context.Context) {
    defer fw.watcher.Close()
    
    for {
        select {
        case event, ok := <-fw.watcher.Events:
            if !ok {
                return
            }
            
            if fw.shouldIgnore(event.Name) {
                continue
            }
            
            if event.Op&fsnotify.Write == fsnotify.Write ||
               event.Op&fsnotify.Create == fsnotify.Create ||
               event.Op&fsnotify.Remove == fsnotify.Remove {
                
                fw.debouncer.Add(event.Name)
            }
            
        case err, ok := <-fw.watcher.Errors:
            if !ok {
                return
            }
            fw.logger.Error("File watcher error: %v", err)
            
        case <-ctx.Done():
            return
        }
    }
}

func (fw *FileWatcher) shouldIgnore(path string) bool {
    for _, ignorePath := range fw.ignorePaths {
        if matched, _ := filepath.Match(ignorePath, path); matched {
            return true
        }
    }
    
    // Ignore common non-source files
    switch filepath.Ext(path) {
    case ".tmp", ".log", ".swp", ".DS_Store":
        return true
    }
    
    return false
}
```

#### Go Compiler with Hot Reload
```go
type GoCompiler struct {
    buildDir    string
    lastBuild   time.Time
    buildCache  map[string]time.Time
    process     *os.Process
    logger      Logger
}

func (gc *GoCompiler) CompileAndReload(changedFiles []string) error {
    gc.logger.Info("Recompiling due to changes in: %v", changedFiles)
    
    // Stop current process
    if gc.process != nil {
        gc.stopProcess()
    }
    
    // Compile
    buildStart := time.Now()
    binary, err := gc.compile()
    if err != nil {
        gc.logger.Error("Compilation failed: %v", err)
        return err
    }
    
    buildDuration := time.Since(buildStart)
    gc.logger.Info("Compilation successful in %v", buildDuration)
    
    // Start new process
    return gc.startProcess(binary)
}

func (gc *GoCompiler) compile() (string, error) {
    binaryPath := filepath.Join(gc.buildDir, "main")
    
    cmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/server")
    cmd.Dir = "."
    
    var stderr bytes.Buffer
    cmd.Stderr = &stderr
    
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("build failed: %v\nOutput: %s", err, stderr.String())
    }
    
    return binaryPath, nil
}

func (gc *GoCompiler) startProcess(binaryPath string) error {
    cmd := exec.Command(binaryPath)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    
    if err := cmd.Start(); err != nil {
        return fmt.Errorf("failed to start process: %v", err)
    }
    
    gc.process = cmd.Process
    gc.logger.Info("Application restarted with PID %d", gc.process.Pid)
    
    return nil
}
```

#### WebSocket Hub for Real-time Updates
```go
type WebSocketHub struct {
    clients    map[*websocket.Conn]bool
    broadcast  chan []byte
    register   chan *websocket.Conn
    unregister chan *websocket.Conn
    upgrader   websocket.Upgrader
}

type ReloadMessage struct {
    Type      string    `json:"type"`
    Timestamp time.Time `json:"timestamp"`
    Files     []string  `json:"files,omitempty"`
    Error     string    `json:"error,omitempty"`
}

func NewWebSocketHub() *WebSocketHub {
    return &WebSocketHub{
        clients:    make(map[*websocket.Conn]bool),
        broadcast:  make(chan []byte),
        register:   make(chan *websocket.Conn),
        unregister: make(chan *websocket.Conn),
        upgrader: websocket.Upgrader{
            CheckOrigin: func(r *http.Request) bool {
                return true // Allow all origins in development
            },
        },
    }
}

func (h *WebSocketHub) Run(ctx context.Context) {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
            
        case client := <-h.unregister:
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                client.Close()
            }
            
        case message := <-h.broadcast:
            for client := range h.clients {
                select {
                case client <- message:
                default:
                    close(client)
                    delete(h.clients, client)
                }
            }
            
        case <-ctx.Done():
            return
        }
    }
}

func (h *WebSocketHub) NotifyReload(files []string) {
    message := ReloadMessage{
        Type:      "reload",
        Timestamp: time.Now(),
        Files:     files,
    }
    
    data, _ := json.Marshal(message)
    select {
    case h.broadcast <- data:
    default:
    }
}
```

### Browser Integration

#### Auto-reload Script
```javascript
// dev-reload.js
(function() {
    'use strict';
    
    class DevReloader {
        constructor() {
            this.ws = null;
            this.reconnectInterval = 1000;
            this.maxReconnectAttempts = 10;
            this.reconnectAttempts = 0;
            
            this.connect();
            this.setupUI();
        }
        
        connect() {
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const wsUrl = `${protocol}//${window.location.host}/_dev/ws`;
            
            this.ws = new WebSocket(wsUrl);
            
            this.ws.onopen = () => {
                console.log('🔄 Dev server connected');
                this.reconnectAttempts = 0;
                this.showStatus('connected');
            };
            
            this.ws.onmessage = (event) => {
                const message = JSON.parse(event.data);
                this.handleMessage(message);
            };
            
            this.ws.onclose = () => {
                console.log('🔌 Dev server disconnected');
                this.showStatus('disconnected');
                this.attemptReconnect();
            };
            
            this.ws.onerror = (error) => {
                console.error('❌ Dev server error:', error);
                this.showStatus('error');
            };
        }
        
        handleMessage(message) {
            switch (message.type) {
                case 'reload':
                    this.handleReload(message);
                    break;
                case 'error':
                    this.handleError(message);
                    break;
                case 'status':
                    this.handleStatus(message);
                    break;
            }
        }
        
        handleReload(message) {
            console.log('🔄 Reloading due to changes:', message.files);
            this.showReloadNotification(message.files);
            
            // Wait a bit for the server to be ready
            setTimeout(() => {
                window.location.reload();
            }, 500);
        }
        
        handleError(message) {
            console.error('🚨 Build error:', message.error);
            this.showErrorNotification(message.error);
        }
        
        setupUI() {
            // Create status indicator
            const statusDiv = document.createElement('div');
            statusDiv.id = 'dev-status';
            statusDiv.style.cssText = `
                position: fixed;
                top: 10px;
                right: 10px;
                padding: 8px 12px;
                border-radius: 4px;
                font-family: monospace;
                font-size: 12px;
                z-index: 10000;
                transition: all 0.3s ease;
            `;
            document.body.appendChild(statusDiv);
        }
        
        showStatus(status) {
            const statusDiv = document.getElementById('dev-status');
            if (!statusDiv) return;
            
            const styles = {
                connected: { background: '#4CAF50', color: 'white', text: '🟢 Dev Server' },
                disconnected: { background: '#FF9800', color: 'white', text: '🟡 Reconnecting...' },
                error: { background: '#F44336', color: 'white', text: '🔴 Error' }
            };
            
            const style = styles[status];
            statusDiv.style.background = style.background;
            statusDiv.style.color = style.color;
            statusDiv.textContent = style.text;
        }
        
        showReloadNotification(files) {
            const notification = document.createElement('div');
            notification.style.cssText = `
                position: fixed;
                top: 50px;
                right: 10px;
                padding: 12px 16px;
                background: #2196F3;
                color: white;
                border-radius: 4px;
                font-family: monospace;
                font-size: 14px;
                z-index: 10001;
                max-width: 300px;
            `;
            notification.innerHTML = `
                <div>🔄 Reloading...</div>
                <div style="font-size: 12px; opacity: 0.8; margin-top: 4px;">
                    Changed: ${files.slice(0, 3).join(', ')}${files.length > 3 ? `... +${files.length - 3}` : ''}
                </div>
            `;
            
            document.body.appendChild(notification);
            
            setTimeout(() => {
                if (notification.parentNode) {
                    notification.parentNode.removeChild(notification);
                }
            }, 2000);
        }
    }
    
    // Initialize when DOM is ready
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', () => new DevReloader());
    } else {
        new DevReloader();
    }
})();
```

### Development Dashboard

#### Dashboard Interface
```go
type Dashboard struct {
    server     *DevServer
    stats      *BuildStats
    logs       *LogBuffer
    clients    *ClientManager
}

type BuildStats struct {
    TotalBuilds     int           `json:"total_builds"`
    SuccessfulBuilds int          `json:"successful_builds"`
    FailedBuilds    int           `json:"failed_builds"`
    AverageBuildTime time.Duration `json:"average_build_time"`
    LastBuildTime   time.Time     `json:"last_build_time"`
    LastBuildDuration time.Duration `json:"last_build_duration"`
}

func (d *Dashboard) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch r.URL.Path {
    case "/_dev/dashboard":
        d.serveDashboard(w, r)
    case "/_dev/api/stats":
        d.serveStats(w, r)
    case "/_dev/api/logs":
        d.serveLogs(w, r)
    case "/_dev/api/clients":
        d.serveClients(w, r)
    default:
        http.NotFound(w, r)
    }
}

func (d *Dashboard) serveDashboard(w http.ResponseWriter, r *http.Request) {
    html := `
<!DOCTYPE html>
<html>
<head>
    <title>ViberCode Dev Dashboard</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, sans-serif; margin: 0; padding: 20px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; }
        .card { background: white; border-radius: 8px; padding: 20px; margin-bottom: 20px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; }
        .stat { text-align: center; }
        .stat-value { font-size: 2em; font-weight: bold; color: #2196F3; }
        .stat-label { color: #666; margin-top: 5px; }
        .logs { font-family: monospace; font-size: 14px; max-height: 400px; overflow-y: auto; background: #1e1e1e; color: #fff; padding: 15px; border-radius: 4px; }
        .log-entry { margin-bottom: 5px; }
        .log-info { color: #4CAF50; }
        .log-warn { color: #FF9800; }
        .log-error { color: #F44336; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚀 ViberCode Development Dashboard</h1>
        
        <div class="card">
            <h2>Build Statistics</h2>
            <div class="stats-grid" id="stats">
                <!-- Stats will be loaded here -->
            </div>
        </div>
        
        <div class="card">
            <h2>Recent Logs</h2>
            <div class="logs" id="logs">
                <!-- Logs will be loaded here -->
            </div>
        </div>
        
        <div class="card">
            <h2>Connected Clients</h2>
            <div id="clients">
                <!-- Client info will be loaded here -->
            </div>
        </div>
    </div>
    
    <script>
        // Dashboard JavaScript implementation
        class DevDashboard {
            constructor() {
                this.loadStats();
                this.loadLogs();
                this.loadClients();
                
                // Refresh every 5 seconds
                setInterval(() => {
                    this.loadStats();
                    this.loadLogs();
                    this.loadClients();
                }, 5000);
            }
            
            async loadStats() {
                try {
                    const response = await fetch('/_dev/api/stats');
                    const stats = await response.json();
                    this.renderStats(stats);
                } catch (error) {
                    console.error('Failed to load stats:', error);
                }
            }
            
            renderStats(stats) {
                const statsDiv = document.getElementById('stats');
                statsDiv.innerHTML = \`
                    <div class="stat">
                        <div class="stat-value">\${stats.total_builds}</div>
                        <div class="stat-label">Total Builds</div>
                    </div>
                    <div class="stat">
                        <div class="stat-value">\${stats.successful_builds}</div>
                        <div class="stat-label">Successful</div>
                    </div>
                    <div class="stat">
                        <div class="stat-value">\${stats.failed_builds}</div>
                        <div class="stat-label">Failed</div>
                    </div>
                    <div class="stat">
                        <div class="stat-value">\${this.formatDuration(stats.average_build_time)}</div>
                        <div class="stat-label">Avg Build Time</div>
                    </div>
                \`;
            }
            
            formatDuration(ns) {
                const ms = ns / 1000000;
                return ms < 1000 ? \`\${ms.toFixed(0)}ms\` : \`\${(ms/1000).toFixed(1)}s\`;
            }
        }
        
        new DevDashboard();
    </script>
</body>
</html>
    `
    
    w.Header().Set("Content-Type", "text/html")
    w.Write([]byte(html))
}
```

### Configuration System

#### Development Configuration
```yaml
# vibercode.dev.yaml
development:
  # Server configuration
  port: 3000
  api_port: 8080
  
  # File watching
  watch_paths:
    - "**/*.go"
    - "**/*.html"
    - "**/*.css"
    - "**/*.js"
    - "**/*.yaml"
  
  ignore_paths:
    - "**/*_test.go"
    - "**/node_modules/**"
    - "**/vendor/**"
    - "**/*.tmp"
    - "**/.git/**"
  
  # Build configuration
  build_dir: ".vibercode/build"
  build_flags: ["-race"]
  build_env:
    GO_ENV: "development"
    DEBUG: "true"
  
  # Features
  auto_reload: true
  proxy_enabled: true
  dashboard: true
  verbose_logging: true
  
  # Proxy configuration
  proxy:
    api_prefix: "/api"
    static_dir: "./web/static"
    spa_mode: true
  
  # Notifications
  notifications:
    desktop: true
    browser: true
    sound: false
```

#### Environment Detection
```go
type EnvironmentDetector struct {
    projectRoot string
}

func (ed *EnvironmentDetector) DetectProjectType() ProjectType {
    // Check for Go project
    if ed.fileExists("go.mod") {
        return ProjectTypeGo
    }
    
    // Check for Node.js project
    if ed.fileExists("package.json") {
        return ProjectTypeNode
    }
    
    // Check for ViberCode project
    if ed.fileExists("vibercode.yaml") {
        return ProjectTypeViberCode
    }
    
    return ProjectTypeUnknown
}

func (ed *EnvironmentDetector) GetDefaultConfig() *DevConfig {
    projectType := ed.DetectProjectType()
    
    config := &DevConfig{
        Port:        3000,
        AutoReload:  true,
        Dashboard:   true,
        BuildDir:    ".vibercode/build",
    }
    
    switch projectType {
    case ProjectTypeGo:
        config.WatchPaths = []string{"**/*.go", "**/*.html", "**/*.css"}
        config.IgnorePaths = []string{"**/*_test.go", "**/vendor/**"}
        config.APIPort = 8080
        
    case ProjectTypeViberCode:
        config.WatchPaths = []string{"**/*.go", "**/*.yaml", "**/*.html"}
        config.ProxyEnabled = true
    }
    
    return config
}
```

### Performance Optimization

#### Incremental Compilation
```go
type IncrementalCompiler struct {
    buildGraph  *DependencyGraph
    cache       *BuildCache
    lastBuild   time.Time
}

func (ic *IncrementalCompiler) NeedsRebuild(changedFiles []string) bool {
    for _, file := range changedFiles {
        if ic.affectsMainBinary(file) {
            return true
        }
    }
    return false
}

func (ic *IncrementalCompiler) affectsMainBinary(file string) bool {
    // Check if file is part of the main binary's dependency graph
    return ic.buildGraph.Affects(file, "main")
}

type BuildCache struct {
    entries map[string]*CacheEntry
    maxSize int64
    dir     string
}

type CacheEntry struct {
    Hash         string
    LastModified time.Time
    Dependencies []string
    Artifacts    []string
}
```

## Dependencies
- Task 02: Template System Enhancement (for development templates)
- Task 08: Testing Framework Integration (for test watching)

## Deliverables
1. Development server with hot reload
2. File watching and automatic compilation
3. WebSocket-based browser integration
4. Development dashboard and monitoring
5. API proxy and frontend integration
6. Configuration system for development
7. Performance optimization features
8. Documentation and setup guides

## Acceptance Criteria
- [ ] Implement development server with hot reload
- [ ] Create file watching system with debouncing
- [ ] Build WebSocket integration for real-time updates
- [ ] Provide browser auto-refresh functionality
- [ ] Add development dashboard with metrics
- [ ] Support API proxy and frontend integration
- [ ] Include configuration management
- [ ] Optimize for fast rebuild times
- [ ] Handle error reporting and recovery
- [ ] Support multiple project types

## Implementation Priority
Medium - Significantly improves development experience

## Estimated Effort
6-7 days

## Notes
- Focus on fast rebuild times and responsiveness
- Handle edge cases like compilation errors gracefully
- Provide clear feedback on development status
- Consider memory usage during long development sessions
- Support both standalone and proxy development modes
- Include comprehensive logging and debugging features