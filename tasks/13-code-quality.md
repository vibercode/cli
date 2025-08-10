# Task 13: Code Quality Tools Integration

## Overview
Integrate comprehensive code quality tools and automation into ViberCode CLI, including linting, formatting, static analysis, security scanning, and code metrics. This system ensures generated code follows best practices and maintains high quality standards.

## Objectives
- Integrate Go linting tools (golangci-lint, staticcheck, vet)
- Implement code formatting and style enforcement (gofmt, goimports)
- Add static analysis and security scanning (gosec, govulncheck)
- Create code metrics and complexity analysis
- Provide pre-commit hooks and automation
- Generate quality reports and dashboards

## Implementation Details

### Command Structure
```bash
# Code quality commands
vibercode quality check                     # Run all quality checks
vibercode quality lint                      # Run linters only
vibercode quality format                    # Format code
vibercode quality security                  # Security analysis
vibercode quality metrics                   # Generate code metrics

# Configuration and setup
vibercode quality init                      # Initialize quality tools
vibercode quality config                    # Show current configuration
vibercode quality install-hooks            # Install pre-commit hooks

# Reporting
vibercode quality report                    # Generate quality report
vibercode quality dashboard                 # Open quality dashboard
vibercode quality trends                    # Show quality trends
```

### Quality Tools Integration

#### Core Linting System
```go
package quality

import (
    "context"
    "fmt"
    "os/exec"
    "path/filepath"
)

type QualityChecker struct {
    config     *QualityConfig
    linters    []Linter
    formatters []Formatter
    scanners   []SecurityScanner
    metrics    *MetricsCollector
}

type QualityConfig struct {
    Enabled         bool                   `yaml:"enabled" json:"enabled"`
    LintOnGenerate  bool                   `yaml:"lint_on_generate" json:"lint_on_generate"`
    FormatOnSave    bool                   `yaml:"format_on_save" json:"format_on_save"`
    Linters         map[string]LinterConfig `yaml:"linters" json:"linters"`
    Formatters      map[string]bool         `yaml:"formatters" json:"formatters"`
    SecurityScan    SecurityConfig          `yaml:"security" json:"security"`
    Metrics         MetricsConfig           `yaml:"metrics" json:"metrics"`
    PreCommitHooks  bool                   `yaml:"pre_commit_hooks" json:"pre_commit_hooks"`
}

type LinterConfig struct {
    Enabled bool     `yaml:"enabled" json:"enabled"`
    Args    []string `yaml:"args" json:"args"`
    Exclude []string `yaml:"exclude" json:"exclude"`
}

func NewQualityChecker(config *QualityConfig) *QualityChecker {
    return &QualityChecker{
        config: config,
        linters: []Linter{
            NewGolangCILint(),
            NewStaticCheck(),
            NewGoVet(),
            NewGoFmt(),
            NewGoImports(),
        },
        formatters: []Formatter{
            NewGoFormatter(),
            NewGoImportsFormatter(),
        },
        scanners: []SecurityScanner{
            NewGosecScanner(),
            NewGovulncheckScanner(),
        },
        metrics: NewMetricsCollector(),
    }
}

func (qc *QualityChecker) RunAllChecks(ctx context.Context, projectPath string) (*QualityReport, error) {
    report := &QualityReport{
        ProjectPath: projectPath,
        Timestamp:   time.Now(),
    }
    
    // Run linting
    lintResults, err := qc.runLinters(ctx, projectPath)
    if err != nil {
        return nil, err
    }
    report.LintResults = lintResults
    
    // Run formatting check
    formatResults, err := qc.checkFormatting(ctx, projectPath)
    if err != nil {
        return nil, err
    }
    report.FormatResults = formatResults
    
    // Run security scans
    securityResults, err := qc.runSecurityScans(ctx, projectPath)
    if err != nil {
        return nil, err
    }
    report.SecurityResults = securityResults
    
    // Collect metrics
    metrics, err := qc.collectMetrics(ctx, projectPath)
    if err != nil {
        return nil, err
    }
    report.Metrics = metrics
    
    return report, nil
}
```

#### Linter Integration
```go
type Linter interface {
    Name() string
    Run(ctx context.Context, path string, config LinterConfig) (*LintResult, error)
    IsAvailable() bool
    Install() error
}

type GolangCILint struct {
    executable string
}

func NewGolangCILint() *GolangCILint {
    return &GolangCILint{
        executable: "golangci-lint",
    }
}

func (g *GolangCILint) Run(ctx context.Context, path string, config LinterConfig) (*LintResult, error) {
    args := []string{"run", "--out-format", "json"}
    args = append(args, config.Args...)
    
    cmd := exec.CommandContext(ctx, g.executable, args...)
    cmd.Dir = path
    
    output, err := cmd.CombinedOutput()
    if err != nil && cmd.ProcessState.ExitCode() != 1 {
        return nil, fmt.Errorf("golangci-lint failed: %v", err)
    }
    
    result := &LintResult{
        Linter:   g.Name(),
        Output:   string(output),
        Issues:   g.parseOutput(output),
        ExitCode: cmd.ProcessState.ExitCode(),
    }
    
    return result, nil
}

func (g *GolangCILint) parseOutput(output []byte) []LintIssue {
    var result struct {
        Issues []struct {
            Pos struct {
                Filename string `json:"filename"`
                Line     int    `json:"line"`
                Column   int    `json:"column"`
            } `json:"pos"`
            Text        string `json:"text"`
            FromLinter  string `json:"fromLinter"`
            Severity    string `json:"severity"`
        } `json:"Issues"`
    }
    
    if err := json.Unmarshal(output, &result); err != nil {
        return nil
    }
    
    issues := make([]LintIssue, len(result.Issues))
    for i, issue := range result.Issues {
        issues[i] = LintIssue{
            File:     issue.Pos.Filename,
            Line:     issue.Pos.Line,
            Column:   issue.Pos.Column,
            Message:  issue.Text,
            Linter:   issue.FromLinter,
            Severity: issue.Severity,
        }
    }
    
    return issues
}
```

#### Code Formatting System
```go
type Formatter interface {
    Name() string
    Format(ctx context.Context, path string) error
    Check(ctx context.Context, path string) (*FormatResult, error)
}

type GoFormatter struct {
    executable string
}

func NewGoFormatter() *GoFormatter {
    return &GoFormatter{
        executable: "gofmt",
    }
}

func (f *GoFormatter) Format(ctx context.Context, path string) error {
    cmd := exec.CommandContext(ctx, f.executable, "-w", path)
    return cmd.Run()
}

func (f *GoFormatter) Check(ctx context.Context, path string) (*FormatResult, error) {
    cmd := exec.CommandContext(ctx, f.executable, "-d", path)
    output, err := cmd.CombinedOutput()
    
    result := &FormatResult{
        Formatter: f.Name(),
        HasIssues: len(output) > 0,
        Diff:      string(output),
    }
    
    if err != nil && cmd.ProcessState.ExitCode() != 0 {
        result.Error = err.Error()
    }
    
    return result, nil
}

type GoImportsFormatter struct {
    executable string
}

func NewGoImportsFormatter() *GoImportsFormatter {
    return &GoImportsFormatter{
        executable: "goimports",
    }
}

func (f *GoImportsFormatter) Format(ctx context.Context, path string) error {
    cmd := exec.CommandContext(ctx, f.executable, "-w", path)
    return cmd.Run()
}
```

#### Security Scanning
```go
type SecurityScanner interface {
    Name() string
    Scan(ctx context.Context, path string) (*SecurityResult, error)
    IsAvailable() bool
}

type GosecScanner struct {
    executable string
}

func NewGosecScanner() *GosecScanner {
    return &GosecScanner{
        executable: "gosec",
    }
}

func (s *GosecScanner) Scan(ctx context.Context, path string) (*SecurityResult, error) {
    cmd := exec.CommandContext(ctx, s.executable, "-fmt", "json", "./...")
    cmd.Dir = path
    
    output, err := cmd.CombinedOutput()
    if err != nil && cmd.ProcessState.ExitCode() != 1 {
        return nil, fmt.Errorf("gosec scan failed: %v", err)
    }
    
    result := &SecurityResult{
        Scanner: s.Name(),
        Issues:  s.parseOutput(output),
    }
    
    return result, nil
}

type GovulncheckScanner struct {
    executable string
}

func NewGovulncheckScanner() *GovulncheckScanner {
    return &GovulncheckScanner{
        executable: "govulncheck",
    }
}

func (s *GovulncheckScanner) Scan(ctx context.Context, path string) (*SecurityResult, error) {
    cmd := exec.CommandContext(ctx, s.executable, "-json", "./...")
    cmd.Dir = path
    
    output, err := cmd.CombinedOutput()
    if err != nil {
        return nil, fmt.Errorf("govulncheck scan failed: %v", err)
    }
    
    result := &SecurityResult{
        Scanner:         s.Name(),
        Vulnerabilities: s.parseVulnerabilities(output),
    }
    
    return result, nil
}
```

#### Code Metrics Collection
```go
type MetricsCollector struct {
    tools []MetricsTool
}

type MetricsTool interface {
    Name() string
    Collect(ctx context.Context, path string) (*Metrics, error)
}

type CodeMetrics struct {
    LinesOfCode       int                    `json:"lines_of_code"`
    CyclomaticComplexity int                 `json:"cyclomatic_complexity"`
    TestCoverage      float64                `json:"test_coverage"`
    TechnicalDebt     time.Duration          `json:"technical_debt"`
    DuplicationRatio  float64                `json:"duplication_ratio"`
    Functions         int                    `json:"functions"`
    Packages          int                    `json:"packages"`
    Dependencies      []string               `json:"dependencies"`
    Maintainability   MaintainabilityIndex   `json:"maintainability"`
}

type MaintainabilityIndex struct {
    Score       float64 `json:"score"`
    Grade       string  `json:"grade"`
    Description string  `json:"description"`
}

func (mc *MetricsCollector) CollectAll(ctx context.Context, path string) (*CodeMetrics, error) {
    metrics := &CodeMetrics{}
    
    // Count lines of code
    loc, err := mc.countLinesOfCode(path)
    if err != nil {
        return nil, err
    }
    metrics.LinesOfCode = loc
    
    // Calculate cyclomatic complexity
    complexity, err := mc.calculateComplexity(path)
    if err != nil {
        return nil, err
    }
    metrics.CyclomaticComplexity = complexity
    
    // Get test coverage
    coverage, err := mc.getTestCoverage(ctx, path)
    if err != nil {
        return nil, err
    }
    metrics.TestCoverage = coverage
    
    // Calculate maintainability index
    maintainability := mc.calculateMaintainability(metrics)
    metrics.Maintainability = maintainability
    
    return metrics, nil
}

func (mc *MetricsCollector) getTestCoverage(ctx context.Context, path string) (float64, error) {
    cmd := exec.CommandContext(ctx, "go", "test", "-cover", "-coverprofile=coverage.out", "./...")
    cmd.Dir = path
    
    output, err := cmd.CombinedOutput()
    if err != nil {
        return 0, err
    }
    
    // Parse coverage from output
    return mc.parseCoverage(string(output)), nil
}
```

### Quality Report Generation

#### Report Structures
```go
type QualityReport struct {
    ProjectPath     string              `json:"project_path"`
    Timestamp       time.Time           `json:"timestamp"`
    LintResults     []LintResult        `json:"lint_results"`
    FormatResults   []FormatResult      `json:"format_results"`
    SecurityResults []SecurityResult    `json:"security_results"`
    Metrics         *CodeMetrics        `json:"metrics"`
    Summary         QualitySummary      `json:"summary"`
}

type QualitySummary struct {
    OverallScore    float64     `json:"overall_score"`
    Grade          string      `json:"grade"`
    TotalIssues    int         `json:"total_issues"`
    CriticalIssues int         `json:"critical_issues"`
    SecurityIssues int         `json:"security_issues"`
    Recommendations []string   `json:"recommendations"`
}

type LintResult struct {
    Linter   string      `json:"linter"`
    Output   string      `json:"output"`
    Issues   []LintIssue `json:"issues"`
    ExitCode int         `json:"exit_code"`
}

type LintIssue struct {
    File     string `json:"file"`
    Line     int    `json:"line"`
    Column   int    `json:"column"`
    Message  string `json:"message"`
    Linter   string `json:"linter"`
    Severity string `json:"severity"`
}
```

#### HTML Report Template
```html
<!DOCTYPE html>
<html>
<head>
    <title>ViberCode Quality Report</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, sans-serif; margin: 0; padding: 20px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; }
        .header { background: white; padding: 30px; border-radius: 8px; margin-bottom: 20px; text-align: center; }
        .score { font-size: 3em; font-weight: bold; margin: 10px 0; }
        .score.excellent { color: #4CAF50; }
        .score.good { color: #8BC34A; }
        .score.fair { color: #FFC107; }
        .score.poor { color: #FF5722; }
        .metrics-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; margin: 20px 0; }
        .metric-card { background: white; padding: 20px; border-radius: 8px; text-align: center; }
        .metric-value { font-size: 2em; font-weight: bold; color: #2196F3; }
        .metric-label { color: #666; margin-top: 5px; }
        .issues-section { background: white; padding: 20px; border-radius: 8px; margin: 20px 0; }
        .issue { padding: 10px; border-left: 4px solid #ccc; margin: 10px 0; background: #f9f9f9; }
        .issue.critical { border-left-color: #F44336; }
        .issue.warning { border-left-color: #FF9800; }
        .issue.info { border-left-color: #2196F3; }
        .recommendations { background: #E3F2FD; padding: 20px; border-radius: 8px; margin: 20px 0; }
        .recommendations ul { margin: 0; padding-left: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔍 Code Quality Report</h1>
            <div class="score {{.Summary.Grade | lower}}">{{.Summary.OverallScore | printf "%.1f"}}</div>
            <div class="grade">Grade: {{.Summary.Grade}}</div>
            <div class="timestamp">Generated: {{.Timestamp.Format "2006-01-02 15:04:05"}}</div>
        </div>

        <div class="metrics-grid">
            <div class="metric-card">
                <div class="metric-value">{{.Metrics.LinesOfCode}}</div>
                <div class="metric-label">Lines of Code</div>
            </div>
            <div class="metric-card">
                <div class="metric-value">{{.Metrics.TestCoverage | printf "%.1f"}}%</div>
                <div class="metric-label">Test Coverage</div>
            </div>
            <div class="metric-card">
                <div class="metric-value">{{.Metrics.CyclomaticComplexity}}</div>
                <div class="metric-label">Complexity</div>
            </div>
            <div class="metric-card">
                <div class="metric-value">{{len .Summary.TotalIssues}}</div>
                <div class="metric-label">Total Issues</div>
            </div>
        </div>

        {{if .Summary.TotalIssues}}
        <div class="issues-section">
            <h2>Issues Found</h2>
            {{range .LintResults}}
            {{range .Issues}}
            <div class="issue {{.Severity}}">
                <strong>{{.File}}:{{.Line}}:{{.Column}}</strong> - {{.Message}}
                <div style="font-size: 0.9em; color: #666;">{{.Linter}}</div>
            </div>
            {{end}}
            {{end}}
        </div>
        {{end}}

        {{if .Summary.Recommendations}}
        <div class="recommendations">
            <h2>💡 Recommendations</h2>
            <ul>
            {{range .Summary.Recommendations}}
                <li>{{.}}</li>
            {{end}}
            </ul>
        </div>
        {{end}}
    </div>
</body>
</html>
```

### Pre-commit Hooks Integration

#### Git Hooks Setup
```go
type PreCommitHookManager struct {
    projectPath string
    config      *QualityConfig
}

func (pm *PreCommitHookManager) InstallHooks() error {
    hookPath := filepath.Join(pm.projectPath, ".git", "hooks", "pre-commit")
    
    hookScript := `#!/bin/bash
# ViberCode pre-commit hook
set -e

echo "🔍 Running ViberCode quality checks..."

# Run quality checks
vibercode quality check --pre-commit

if [ $? -ne 0 ]; then
    echo "❌ Quality checks failed. Commit aborted."
    echo "Run 'vibercode quality check' to see details."
    exit 1
fi

echo "✅ Quality checks passed!"
`
    
    err := ioutil.WriteFile(hookPath, []byte(hookScript), 0755)
    if err != nil {
        return fmt.Errorf("failed to install pre-commit hook: %v", err)
    }
    
    return nil
}
```

#### Quality Gates
```go
type QualityGate struct {
    MinScore            float64 `yaml:"min_score" json:"min_score"`
    MaxCriticalIssues   int     `yaml:"max_critical_issues" json:"max_critical_issues"`
    MinTestCoverage     float64 `yaml:"min_test_coverage" json:"min_test_coverage"`
    MaxComplexity       int     `yaml:"max_complexity" json:"max_complexity"`
    AllowSecurityIssues bool    `yaml:"allow_security_issues" json:"allow_security_issues"`
}

func (qg *QualityGate) Check(report *QualityReport) error {
    var violations []string
    
    if report.Summary.OverallScore < qg.MinScore {
        violations = append(violations, fmt.Sprintf("Overall score %.1f is below minimum %.1f", 
            report.Summary.OverallScore, qg.MinScore))
    }
    
    if report.Summary.CriticalIssues > qg.MaxCriticalIssues {
        violations = append(violations, fmt.Sprintf("Critical issues %d exceed maximum %d", 
            report.Summary.CriticalIssues, qg.MaxCriticalIssues))
    }
    
    if report.Metrics.TestCoverage < qg.MinTestCoverage {
        violations = append(violations, fmt.Sprintf("Test coverage %.1f%% is below minimum %.1f%%", 
            report.Metrics.TestCoverage, qg.MinTestCoverage))
    }
    
    if !qg.AllowSecurityIssues && report.Summary.SecurityIssues > 0 {
        violations = append(violations, fmt.Sprintf("Security issues found: %d", 
            report.Summary.SecurityIssues))
    }
    
    if len(violations) > 0 {
        return fmt.Errorf("Quality gate violations:\n%s", strings.Join(violations, "\n"))
    }
    
    return nil
}
```

### Configuration Management

#### Quality Configuration File
```yaml
# vibercode.quality.yaml
quality:
  # Global settings
  enabled: true
  lint_on_generate: true
  format_on_save: true
  pre_commit_hooks: true
  
  # Linter configuration
  linters:
    golangci-lint:
      enabled: true
      args: ["--timeout=5m", "--enable-all"]
      exclude: ["test"]
    
    staticcheck:
      enabled: true
      args: []
    
    govet:
      enabled: true
      args: []
    
    gofmt:
      enabled: true
      args: []
    
    goimports:
      enabled: true
      args: []
  
  # Formatter settings
  formatters:
    gofmt: true
    goimports: true
  
  # Security scanning
  security:
    gosec:
      enabled: true
      severity: "medium"
      exclude_rules: []
    
    govulncheck:
      enabled: true
      db_update: true
  
  # Code metrics
  metrics:
    complexity_threshold: 10
    duplication_threshold: 0.1
    maintainability_threshold: 70
  
  # Quality gates
  gates:
    min_score: 8.0
    max_critical_issues: 0
    min_test_coverage: 80.0
    max_complexity: 15
    allow_security_issues: false
  
  # Reporting
  reports:
    format: ["html", "json", "text"]
    output_dir: ".vibercode/reports"
    history: true
```

## Dependencies
- Task 02: Template System Enhancement (for quality templates)
- Task 08: Testing Framework Integration (for test coverage analysis)

## Deliverables
1. Comprehensive linting tool integration
2. Code formatting and style enforcement
3. Static analysis and security scanning
4. Code metrics collection and analysis
5. Quality report generation system
6. Pre-commit hooks and automation
7. Quality gates and compliance checking
8. Configuration management system

## Acceptance Criteria
- [ ] Integrate multiple Go linting tools (golangci-lint, staticcheck, vet)
- [ ] Implement automated code formatting (gofmt, goimports)
- [ ] Add security scanning (gosec, govulncheck)
- [ ] Collect comprehensive code metrics
- [ ] Generate HTML/JSON quality reports
- [ ] Support pre-commit hooks and automation
- [ ] Implement quality gates and compliance checking
- [ ] Provide configurable quality standards
- [ ] Include trend analysis and history tracking
- [ ] Support CI/CD integration

## Implementation Priority
Low - Improves code quality and maintainability

## Estimated Effort
4-5 days

## Notes
- Focus on actionable feedback and recommendations
- Ensure integration with existing development workflows
- Provide clear configuration options for different teams
- Consider performance impact of quality checks
- Support gradual adoption and customizable standards
- Include educational content about best practices