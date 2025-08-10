# Task 11: IDE Integration

## Overview
Develop comprehensive IDE integrations and extensions that enhance the developer experience when working with ViberCode CLI-generated projects. This includes VS Code extensions, IntelliJ/GoLand plugins, and Language Server Protocol implementation.

## Objectives
- Create VS Code extension for ViberCode CLI integration
- Develop IntelliJ/GoLand plugin for Go development workflow
- Implement Language Server Protocol for advanced code intelligence
- Provide project templates and scaffolding integration
- Add debugging and testing integration
- Create code navigation and refactoring tools

## Implementation Details

### VS Code Extension

#### Extension Manifest
```json
{
  "name": "vibercode",
  "displayName": "ViberCode CLI Integration",
  "description": "Full integration with ViberCode CLI for Go API development",
  "version": "1.0.0",
  "publisher": "vibercode",
  "engines": {
    "vscode": "^1.74.0"
  },
  "categories": ["Other"],
  "activationEvents": [
    "workspaceContains:**/vibercode.yaml",
    "workspaceContains:**/go.mod",
    "onCommand:vibercode.generate"
  ],
  "main": "./out/extension.js",
  "contributes": {
    "commands": [
      {
        "command": "vibercode.generate.api",
        "title": "Generate API Project",
        "category": "ViberCode"
      },
      {
        "command": "vibercode.generate.resource",
        "title": "Generate Resource",
        "category": "ViberCode"
      },
      {
        "command": "vibercode.generate.middleware",
        "title": "Generate Middleware",
        "category": "ViberCode"
      }
    ],
    "menus": {
      "explorer/context": [
        {
          "submenu": "vibercode.generate",
          "group": "navigation"
        }
      ]
    },
    "configuration": {
      "title": "ViberCode",
      "properties": {
        "vibercode.cliPath": {
          "type": "string",
          "default": "vibercode",
          "description": "Path to ViberCode CLI executable"
        }
      }
    }
  }
}
```

#### Core Extension Features

##### 1. Project Scaffolding
```typescript
import * as vscode from 'vscode';
import { ViberCodeCLI } from './cli';

export class ProjectScaffolder {
    constructor(private cli: ViberCodeCLI) {}

    async generateAPI() {
        const options = await this.promptForAPIOptions();
        
        const workspaceFolder = vscode.workspace.workspaceFolders?.[0];
        if (!workspaceFolder) {
            vscode.window.showErrorMessage('No workspace folder open');
            return;
        }

        await vscode.window.withProgress({
            location: vscode.ProgressLocation.Notification,
            title: "Generating API project...",
            cancellable: false
        }, async (progress) => {
            progress.report({ increment: 0 });
            
            await this.cli.generateAPI(options, workspaceFolder.uri.fsPath);
            
            progress.report({ increment: 100 });
            vscode.window.showInformationMessage('API project generated successfully!');
        });
    }

    private async promptForAPIOptions(): Promise<APIGenerationOptions> {
        const projectName = await vscode.window.showInputBox({
            prompt: 'Project name',
            validateInput: (value) => {
                if (!value || value.trim().length === 0) {
                    return 'Project name is required';
                }
                return null;
            }
        });

        const port = await vscode.window.showInputBox({
            prompt: 'Server port',
            value: '8080',
            validateInput: (value) => {
                const port = parseInt(value);
                if (isNaN(port) || port < 1 || port > 65535) {
                    return 'Please enter a valid port number (1-65535)';
                }
                return null;
            }
        });

        const database = await vscode.window.showQuickPick([
            { label: 'PostgreSQL', value: 'postgresql' },
            { label: 'MySQL', value: 'mysql' },
            { label: 'SQLite', value: 'sqlite' },
            { label: 'MongoDB', value: 'mongodb' }
        ], {
            placeHolder: 'Select database type'
        });

        return {
            name: projectName!,
            port: parseInt(port!),
            database: database!.value
        };
    }
}
```

##### 2. Resource Generation
```typescript
export class ResourceGenerator {
    async generateResource() {
        const editor = vscode.window.activeTextEditor;
        if (!editor) {
            vscode.window.showErrorMessage('No active editor');
            return;
        }

        const resourceName = await vscode.window.showInputBox({
            prompt: 'Resource name (e.g., User, Product, Order)',
            validateInput: (value) => {
                if (!value || !value.match(/^[A-Z][a-zA-Z0-9]*$/)) {
                    return 'Resource name must start with uppercase letter';
                }
                return null;
            }
        });

        const fields = await this.promptForFields();
        
        await this.cli.generateResource({
            name: resourceName!,
            fields: fields,
            outputPath: editor.document.uri.fsPath
        });

        vscode.window.showInformationMessage(`Resource ${resourceName} generated successfully!`);
    }

    private async promptForFields(): Promise<Field[]> {
        const fields: Field[] = [];
        
        while (true) {
            const fieldName = await vscode.window.showInputBox({
                prompt: `Field name (or press Enter to finish)`,
                validateInput: (value) => {
                    if (value && !value.match(/^[a-z][a-zA-Z0-9]*$/)) {
                        return 'Field name must start with lowercase letter';
                    }
                    return null;
                }
            });

            if (!fieldName) break;

            const fieldType = await vscode.window.showQuickPick([
                { label: 'String', value: 'string' },
                { label: 'Number', value: 'number' },
                { label: 'Boolean', value: 'boolean' },
                { label: 'Date', value: 'date' },
                { label: 'UUID', value: 'uuid' },
                { label: 'JSON', value: 'json' }
            ], {
                placeHolder: 'Select field type'
            });

            fields.push({
                name: fieldName,
                type: fieldType!.value,
                required: true
            });
        }

        return fields;
    }
}
```

##### 3. Code Navigation and Intelligence
```typescript
export class ViberCodeDefinitionProvider implements vscode.DefinitionProvider {
    provideDefinition(
        document: vscode.TextDocument,
        position: vscode.Position,
        token: vscode.CancellationToken
    ): vscode.ProviderResult<vscode.Definition | vscode.LocationLink[]> {
        
        const range = document.getWordRangeAtPosition(position);
        const word = document.getText(range);
        
        // Check if it's a ViberCode-generated structure
        if (this.isViberCodeStruct(word)) {
            return this.findStructDefinition(word, document.uri);
        }
        
        return null;
    }

    private isViberCodeStruct(word: string): boolean {
        // Logic to identify ViberCode-generated structures
        return word.match(/^(.*Handler|.*Service|.*Repository)$/) !== null;
    }

    private async findStructDefinition(structName: string, currentUri: vscode.Uri): Promise<vscode.Location[]> {
        // Search for struct definition in generated files
        const workspaceFolder = vscode.workspace.getWorkspaceFolder(currentUri);
        if (!workspaceFolder) return [];

        const pattern = new vscode.RelativePattern(workspaceFolder, '**/*.go');
        const files = await vscode.workspace.findFiles(pattern);
        
        for (const file of files) {
            const document = await vscode.workspace.openTextDocument(file);
            const text = document.getText();
            
            const regex = new RegExp(`type\\s+${structName}\\s+struct`, 'g');
            const match = regex.exec(text);
            
            if (match) {
                const position = document.positionAt(match.index);
                return [new vscode.Location(file, position)];
            }
        }
        
        return [];
    }
}
```

##### 4. Testing Integration
```typescript
export class TestRunner {
    async runTests() {
        const workspaceFolder = vscode.workspace.workspaceFolders?.[0];
        if (!workspaceFolder) return;

        const terminal = vscode.window.createTerminal('ViberCode Tests');
        terminal.show();
        
        // Run tests with coverage
        terminal.sendText('go test -v -cover ./...');
        
        // Parse test results and show in problems panel
        this.parseTestResults();
    }

    async generateTests() {
        const editor = vscode.window.activeTextEditor;
        if (!editor) return;

        const currentFile = editor.document.uri.fsPath;
        const testType = await vscode.window.showQuickPick([
            { label: 'Unit Tests', value: 'unit' },
            { label: 'Integration Tests', value: 'integration' },
            { label: 'Benchmark Tests', value: 'benchmark' }
        ]);

        await this.cli.generateTests({
            type: testType!.value,
            sourceFile: currentFile
        });
    }

    private parseTestResults() {
        // Implementation to parse go test output and populate problems panel
    }
}
```

### IntelliJ/GoLand Plugin

#### Plugin Structure
```
vibercode-intellij-plugin/
├── src/main/
│   ├── java/com/vibercode/intellij/
│   │   ├── ViberCodePlugin.java
│   │   ├── actions/
│   │   │   ├── GenerateAPIAction.java
│   │   │   ├── GenerateResourceAction.java
│   │   │   └── GenerateMiddlewareAction.java
│   │   ├── ui/
│   │   │   ├── dialogs/
│   │   │   └── panels/
│   │   └── services/
│   │       ├── ViberCodeService.java
│   │       └── CLIExecutor.java
│   └── resources/
│       ├── META-INF/plugin.xml
│       └── icons/
└── build.gradle
```

#### Plugin Configuration
```xml
<!-- plugin.xml -->
<idea-plugin>
    <id>com.vibercode.intellij</id>
    <name>ViberCode</name>
    <version>1.0.0</version>
    <vendor>ViberCode</vendor>
    
    <description><![CDATA[
        ViberCode CLI integration for IntelliJ IDEA and GoLand.
        Generate Go APIs, resources, middleware, and tests directly from the IDE.
    ]]></description>
    
    <depends>com.intellij.modules.platform</depends>
    <depends>org.jetbrains.plugins.go</depends>
    
    <extensions defaultExtensionNs="com.intellij">
        <applicationService serviceImplementation="com.vibercode.intellij.services.ViberCodeService"/>
        
        <projectService serviceImplementation="com.vibercode.intellij.services.CLIExecutor"/>
        
        <toolWindow id="ViberCode" secondary="true" anchor="right"
                    factoryClass="com.vibercode.intellij.ui.ViberCodeToolWindowFactory"/>
    </extensions>
    
    <actions>
        <group id="ViberCode.GenerateGroup" text="ViberCode" popup="true">
            <add-to-group group-id="NewGroup" anchor="after" relative-to-action="NewFile"/>
            
            <action id="ViberCode.GenerateAPI" 
                    class="com.vibercode.intellij.actions.GenerateAPIAction"
                    text="API Project" description="Generate new API project"/>
                    
            <action id="ViberCode.GenerateResource"
                    class="com.vibercode.intellij.actions.GenerateResourceAction"
                    text="Resource" description="Generate CRUD resource"/>
        </group>
    </actions>
</idea-plugin>
```

### Language Server Protocol Implementation

#### LSP Server Structure
```go
// main.go
package main

import (
    "context"
    "github.com/sourcegraph/jsonrpc2"
    "github.com/sourcegraph/go-lsp"
)

type ViberCodeLSP struct {
    conn jsonrpc2.JSONRPC2
}

func (s *ViberCodeLSP) Initialize(ctx context.Context, params *lsp.InitializeParams) (*lsp.InitializeResult, error) {
    return &lsp.InitializeResult{
        Capabilities: lsp.ServerCapabilities{
            TextDocumentSync: &lsp.TextDocumentSyncOptionsOrKind{
                Options: &lsp.TextDocumentSyncOptions{
                    OpenClose: true,
                    Change:    lsp.TDSKFull,
                },
            },
            CompletionProvider: &lsp.CompletionOptions{
                TriggerCharacters: []string{"."},
            },
            DefinitionProvider:     true,
            DocumentSymbolProvider: true,
            WorkspaceSymbolProvider: true,
            CodeActionProvider:     true,
        },
    }, nil
}

func (s *ViberCodeLSP) TextDocumentDefinition(ctx context.Context, params *lsp.TextDocumentPositionParams) ([]lsp.Location, error) {
    // Implementation for go-to-definition
    return s.findDefinitions(params)
}

func (s *ViberCodeLSP) TextDocumentCompletion(ctx context.Context, params *lsp.CompletionParams) (*lsp.CompletionList, error) {
    // Implementation for code completion
    return s.provideCompletions(params)
}
```

#### Code Actions
```go
func (s *ViberCodeLSP) TextDocumentCodeAction(ctx context.Context, params *lsp.CodeActionParams) ([]lsp.CodeAction, error) {
    var actions []lsp.CodeAction
    
    // Generate resource from struct
    if s.isStructDefinition(params) {
        actions = append(actions, lsp.CodeAction{
            Title: "Generate ViberCode Resource",
            Kind:  lsp.CAKRefactor,
            Command: &lsp.Command{
                Command:   "vibercode.generateResource",
                Arguments: []interface{}{params.TextDocument.URI},
            },
        })
    }
    
    // Generate tests
    if s.isFunctionDefinition(params) {
        actions = append(actions, lsp.CodeAction{
            Title: "Generate Tests",
            Kind:  lsp.CAKRefactor,
            Command: &lsp.Command{
                Command:   "vibercode.generateTests",
                Arguments: []interface{}{params.TextDocument.URI},
            },
        })
    }
    
    return actions, nil
}
```

### Project Templates Integration

#### Template Discovery
```typescript
export class TemplateManager {
    async discoverTemplates(): Promise<ProjectTemplate[]> {
        const templates: ProjectTemplate[] = [];
        
        // Built-in templates
        templates.push(...this.getBuiltInTemplates());
        
        // User templates
        templates.push(...await this.getUserTemplates());
        
        // Community templates
        templates.push(...await this.getCommunityTemplates());
        
        return templates;
    }

    private getBuiltInTemplates(): ProjectTemplate[] {
        return [
            {
                id: 'go-api-basic',
                name: 'Go API (Basic)',
                description: 'Basic Go API with clean architecture',
                category: 'API',
                framework: 'Go + Gin',
                features: ['REST API', 'Database', 'Docker'],
                command: 'vibercode generate api --template basic'
            },
            {
                id: 'go-microservice',
                name: 'Go Microservice',
                description: 'Microservice with monitoring and deployment',
                category: 'Microservice',
                framework: 'Go + Gin',
                features: ['gRPC', 'Monitoring', 'K8s', 'CI/CD'],
                command: 'vibercode generate api --template microservice'
            }
        ];
    }
}
```

#### Template Wizard
```typescript
export class TemplateWizard {
    async showTemplateWizard(): Promise<void> {
        const templates = await this.templateManager.discoverTemplates();
        
        const selected = await vscode.window.showQuickPick(
            templates.map(t => ({
                label: t.name,
                description: t.description,
                detail: `Framework: ${t.framework} | Features: ${t.features.join(', ')}`,
                template: t
            })),
            {
                placeHolder: 'Choose a project template',
                matchOnDescription: true,
                matchOnDetail: true
            }
        );

        if (selected) {
            await this.generateFromTemplate(selected.template);
        }
    }

    private async generateFromTemplate(template: ProjectTemplate): Promise<void> {
        const config = await this.promptForTemplateConfig(template);
        
        const targetFolder = await vscode.window.showOpenDialog({
            canSelectFolders: true,
            canSelectFiles: false,
            canSelectMany: false,
            openLabel: 'Select project folder'
        });

        if (targetFolder) {
            await this.cli.executeCommand(template.command, {
                ...config,
                outputPath: targetFolder[0].fsPath
            });
        }
    }
}
```

### Debugging Integration

#### Debug Configuration
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug ViberCode API",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/server",
            "env": {
                "GO_ENV": "development"
            },
            "args": [],
            "showLog": true
        }
    ]
}
```

#### Test Debugging
```typescript
export class TestDebugger {
    async debugTest(testName: string, filePath: string): Promise<void> {
        const config: vscode.DebugConfiguration = {
            name: `Debug Test: ${testName}`,
            type: 'go',
            request: 'launch',
            mode: 'test',
            program: filePath,
            args: ['-test.run', `^${testName}$`],
            showLog: true
        };

        await vscode.debug.startDebugging(undefined, config);
    }
}
```

## Dependencies
- Task 02: Template System Enhancement (for IDE templates)
- Task 08: Testing Framework Integration (for test debugging)

## Deliverables
1. VS Code extension with full CLI integration
2. IntelliJ/GoLand plugin for Go development
3. Language Server Protocol implementation
4. Project template integration and wizards
5. Code navigation and refactoring tools
6. Testing and debugging integration
7. Documentation and setup guides
8. Extension marketplace publishing

## Acceptance Criteria
- [ ] Create VS Code extension with CLI integration
- [ ] Develop IntelliJ/GoLand plugin
- [ ] Implement Language Server Protocol
- [ ] Provide project scaffolding and templates
- [ ] Add code navigation and intelligence
- [ ] Integrate testing and debugging tools
- [ ] Support multiple IDE platforms
- [ ] Include comprehensive documentation
- [ ] Publish to extension marketplaces
- [ ] Provide developer setup guides

## Implementation Priority
Medium - Enhances developer productivity and adoption

## Estimated Effort
7-8 days

## Notes
- Focus on seamless integration with existing workflows
- Ensure cross-platform compatibility
- Provide rich developer experience with IntelliSense
- Consider performance impact of extensions
- Plan for extension updates and maintenance
- Include telemetry for improvement insights