package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/vibercode/cli/internal/models"
	"github.com/vibercode/cli/internal/templates"
	"github.com/vibercode/cli/pkg/ui"
)

// UIOptions contains configuration for UI generation
type UIOptions struct {
	AtomicDesign bool
	Framework    string
	TypeScript   bool
	Storybook    bool
}

// UIGenerator handles UI component generation
type UIGenerator struct {
	options UIOptions
}

// NewUIGenerator creates a new UI generator
func NewUIGenerator() *UIGenerator {
	return &UIGenerator{}
}

// Generate generates UI components based on options
func (g *UIGenerator) Generate(options UIOptions) error {
	g.options = options

	ui.PrintStep(1, 1, "Starting UI component generation...")

	// If atomic design flag is set, generate complete structure
	if options.AtomicDesign {
		return g.generateAtomicDesignStructure()
	}

	// Otherwise, generate individual component
	return g.generateSingleComponent()
}

// generateAtomicDesignStructure generates complete atomic design structure
func (g *UIGenerator) generateAtomicDesignStructure() error {
	ui.PrintStep(1, 6, "Generating Atomic Design structure...")

	// Get project configuration
	config, err := g.getProjectConfig()
	if err != nil {
		return fmt.Errorf("failed to get project config: %w", err)
	}

	// Create directory structure
	ui.PrintStep(2, 6, "Creating directory structure...")
	if err := g.createDirectoryStructure(); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}

	// Generate configuration files
	ui.PrintStep(3, 6, "Generating configuration files...")
	if err := g.generateConfigFiles(config); err != nil {
		return fmt.Errorf("failed to generate config files: %w", err)
	}

	// Generate atoms
	ui.PrintStep(4, 6, "Generating atoms...")
	atoms := models.GetDefaultAtoms(models.Framework(config.Framework), config.TypeScript)
	for _, atom := range atoms {
		if err := g.generateComponent(atom); err != nil {
			return fmt.Errorf("failed to generate atom %s: %w", atom.Name, err)
		}
	}

	// Generate molecules
	ui.PrintStep(5, 6, "Generating molecules...")
	molecules := models.GetDefaultMolecules(models.Framework(config.Framework), config.TypeScript)
	for _, molecule := range molecules {
		if err := g.generateComponent(molecule); err != nil {
			return fmt.Errorf("failed to generate molecule %s: %w", molecule.Name, err)
		}
	}

	// Generate organisms
	ui.PrintStep(6, 6, "Generating organisms...")
	organisms := models.GetDefaultOrganisms(models.Framework(config.Framework), config.TypeScript)
	for _, organism := range organisms {
		if err := g.generateComponent(organism); err != nil {
			return fmt.Errorf("failed to generate organism %s: %w", organism.Name, err)
		}
	}

	// Generate index files for each level
	if err := g.generateIndexFiles(config); err != nil {
		return fmt.Errorf("failed to generate index files: %w", err)
	}

	ui.PrintSuccess("Atomic Design structure generated successfully!")
	g.showGeneratedStructure()

	return nil
}

// generateSingleComponent generates a single component
func (g *UIGenerator) generateSingleComponent() error {
	ui.PrintStep(1, 2, "Generating single component...")

	// Get component configuration
	component, err := g.getComponentConfig()
	if err != nil {
		return fmt.Errorf("failed to get component config: %w", err)
	}

	// Generate the component
	ui.PrintStep(2, 2, "Creating component files...")
	if err := g.generateComponent(component); err != nil {
		return fmt.Errorf("failed to generate component: %w", err)
	}

	ui.PrintSuccess(fmt.Sprintf("Component %s generated successfully!", component.Name))
	return nil
}

// getProjectConfig gets project configuration from user
func (g *UIGenerator) getProjectConfig() (*models.AtomicDesignStructure, error) {
	var config models.AtomicDesignStructure

	// Project name
	namePrompt := promptui.Prompt{
		Label:   ui.IconPackage + " Project name",
		Default: "my-ui-components",
	}
	name, err := namePrompt.Run()
	if err != nil {
		return nil, err
	}
	config.ProjectName = name

	// Framework selection
	if g.options.Framework == "" {
		frameworkPrompt := promptui.Select{
			Label: ui.IconCode + " Select frontend framework",
			Items: []string{"react", "vue", "angular"},
		}
		_, framework, err := frameworkPrompt.Run()
		if err != nil {
			return nil, err
		}
		config.Framework = models.Framework(framework)
	} else {
		config.Framework = models.Framework(g.options.Framework)
	}

	config.TypeScript = g.options.TypeScript
	config.Storybook = g.options.Storybook

	return &config, nil
}

// getComponentConfig gets single component configuration
func (g *UIGenerator) getComponentConfig() (models.UIComponent, error) {
	var component models.UIComponent

	// Component name
	namePrompt := promptui.Prompt{
		Label: ui.IconCode + " Component name",
	}
	name, err := namePrompt.Run()
	if err != nil {
		return component, err
	}
	component.Name = name

	// Component type
	typePrompt := promptui.Select{
		Label: ui.IconBuild + " Component type",
		Items: []string{"atom", "molecule", "organism", "template", "page"},
	}
	_, componentType, err := typePrompt.Run()
	if err != nil {
		return component, err
	}
	component.Type = models.UIComponentType(componentType)

	// Framework
	if g.options.Framework == "" {
		frameworkPrompt := promptui.Select{
			Label: ui.IconGear + " Frontend framework",
			Items: []string{"react", "vue", "angular"},
		}
		_, framework, err := frameworkPrompt.Run()
		if err != nil {
			return component, err
		}
		component.Framework = models.Framework(framework)
	} else {
		component.Framework = models.Framework(g.options.Framework)
	}

	component.TypeScript = g.options.TypeScript
	component.HasStory = g.options.Storybook
	component.HasTest = true
	component.HasStyles = true

	// Description
	descPrompt := promptui.Prompt{
		Label:   ui.IconDoc + " Component description",
		Default: fmt.Sprintf("A reusable %s component", strings.ToLower(name)),
	}
	desc, err := descPrompt.Run()
	if err != nil {
		return component, err
	}
	component.Description = desc

	return component, nil
}

// createDirectoryStructure creates the atomic design directory structure
func (g *UIGenerator) createDirectoryStructure() error {
	dirs := []string{
		"src/components/atoms",
		"src/components/molecules",
		"src/components/organisms",
		"src/components/templates",
		"src/pages",
		"src/styles",
		"src/types",
	}

	if g.options.Storybook {
		dirs = append(dirs, ".storybook", "src/stories")
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// generateConfigFiles generates configuration files
func (g *UIGenerator) generateConfigFiles(config *models.AtomicDesignStructure) error {
	// Package.json
	packageContent := templates.GetPackageJSON(config)
	if err := g.writeFile("package.json", packageContent); err != nil {
		return fmt.Errorf("failed to write package.json: %w", err)
	}

	// TypeScript config
	if config.TypeScript {
		tsContent := templates.GetTSConfig(config)
		if err := g.writeFile("tsconfig.json", tsContent); err != nil {
			return fmt.Errorf("failed to write tsconfig.json: %w", err)
		}
	}

	// Storybook config
	if config.Storybook {
		storybookContent := templates.GetStorybookConfig(config)
		if err := g.writeFile(".storybook/main.js", storybookContent); err != nil {
			return fmt.Errorf("failed to write storybook config: %w", err)
		}
	}

	// Global styles
	stylesContent := templates.GetGlobalStyles(config)
	if err := g.writeFile("src/styles/globals.scss", stylesContent); err != nil {
		return fmt.Errorf("failed to write global styles: %w", err)
	}

	return nil
}

// generateComponent generates a single UI component
func (g *UIGenerator) generateComponent(component models.UIComponent) error {
	// Create component directory
	componentDir := filepath.Join(component.GetDirectoryPath(), component.Name)
	if err := os.MkdirAll(componentDir, 0755); err != nil {
		return fmt.Errorf("failed to create component directory: %w", err)
	}

	// Generate component file
	componentContent := templates.GetComponentTemplate(component)
	componentFile := filepath.Join(componentDir, component.GetFileName())
	if err := g.writeFile(componentFile, componentContent); err != nil {
		return fmt.Errorf("failed to write component file: %w", err)
	}

	// Generate styles file
	if component.HasStyles && component.GetStyleFileName() != "" {
		stylesContent := templates.GetComponentStyles(component)
		stylesFile := filepath.Join(componentDir, component.GetStyleFileName())
		if err := g.writeFile(stylesFile, stylesContent); err != nil {
			return fmt.Errorf("failed to write styles file: %w", err)
		}
	}

	// Generate test file
	if component.HasTest {
		testContent := templates.GetComponentTest(component)
		testFile := filepath.Join(componentDir, component.GetTestFileName())
		if err := g.writeFile(testFile, testContent); err != nil {
			return fmt.Errorf("failed to write test file: %w", err)
		}
	}

	// Generate story file
	if component.HasStory {
		storyContent := templates.GetComponentStory(component)
		storyFile := filepath.Join(componentDir, component.GetStoryFileName())
		if err := g.writeFile(storyFile, storyContent); err != nil {
			return fmt.Errorf("failed to write story file: %w", err)
		}
	}

	// Generate index file for component
	indexContent := templates.GetComponentIndex(component)
	indexFile := filepath.Join(componentDir, "index.ts")
	if err := g.writeFile(indexFile, indexContent); err != nil {
		return fmt.Errorf("failed to write component index: %w", err)
	}

	return nil
}

// generateIndexFiles generates index files for each atomic design level
func (g *UIGenerator) generateIndexFiles(config *models.AtomicDesignStructure) error {
	levels := []string{"atoms", "molecules", "organisms", "templates"}

	for _, level := range levels {
		indexContent := templates.GetLevelIndex(level, config)
		indexFile := filepath.Join("src/components", level, "index.ts")
		if err := g.writeFile(indexFile, indexContent); err != nil {
			return fmt.Errorf("failed to write %s index: %w", level, err)
		}
	}

	// Main components index
	mainIndexContent := templates.GetMainComponentsIndex(config)
	if err := g.writeFile("src/components/index.ts", mainIndexContent); err != nil {
		return fmt.Errorf("failed to write main components index: %w", err)
	}

	return nil
}

// writeFile writes content to a file
func (g *UIGenerator) writeFile(filePath, content string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Write file
	return os.WriteFile(filePath, []byte(content), 0644)
}

// showGeneratedStructure shows the generated directory structure
func (g *UIGenerator) showGeneratedStructure() {
	ui.PrintInfo("Generated structure:")
	structure := `
📁 src/
├── 📁 components/
│   ├── 📁 atoms/
│   │   ├── 📁 Button/
│   │   ├── 📁 Input/
│   │   └── 📁 Label/
│   ├── 📁 molecules/
│   │   ├── 📁 FormField/
│   │   └── 📁 Card/
│   ├── 📁 organisms/
│   │   ├── 📁 Header/
│   │   └── 📁 Sidebar/
│   └── 📁 templates/
├── 📁 pages/
├── 📁 styles/
└── 📁 types/`

	if g.options.Storybook {
		structure += `
📁 .storybook/
📁 src/stories/`
	}

	fmt.Println(structure)
}