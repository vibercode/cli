package models

import (
	"fmt"
	"strings"
)

// UIComponentType represents the atomic design level
type UIComponentType string

const (
	AtomType     UIComponentType = "atom"
	MoleculeType UIComponentType = "molecule"
	OrganismType UIComponentType = "organism"
	TemplateType UIComponentType = "template"
	PageType     UIComponentType = "page"
)

// Framework represents supported frontend frameworks
type Framework string

const (
	ReactFramework   Framework = "react"
	VueFramework     Framework = "vue"
	AngularFramework Framework = "angular"
)

// UIComponent represents a UI component definition
type UIComponent struct {
	Name        string          `json:"name"`
	Type        UIComponentType `json:"type"`
	Framework   Framework       `json:"framework"`
	TypeScript  bool            `json:"typescript"`
	Props       []UIProp        `json:"props"`
	Children    []UIComponent   `json:"children,omitempty"`
	HasStory    bool            `json:"hasStory"`
	HasTest     bool            `json:"hasTest"`
	HasStyles   bool            `json:"hasStyles"`
	Description string          `json:"description"`
}

// UIProp represents a component property
type UIProp struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Optional    bool   `json:"optional"`
	Default     string `json:"default,omitempty"`
	Description string `json:"description"`
}

// GetFileName returns the component file name based on framework
func (c *UIComponent) GetFileName() string {
	switch c.Framework {
	case ReactFramework:
		if c.TypeScript {
			return fmt.Sprintf("%s.tsx", c.Name)
		}
		return fmt.Sprintf("%s.jsx", c.Name)
	case VueFramework:
		if c.TypeScript {
			return fmt.Sprintf("%s.vue", c.Name)
		}
		return fmt.Sprintf("%s.vue", c.Name)
	case AngularFramework:
		if c.TypeScript {
			return fmt.Sprintf("%s.component.ts", strings.ToLower(c.Name))
		}
		return fmt.Sprintf("%s.component.js", strings.ToLower(c.Name))
	default:
		return fmt.Sprintf("%s.tsx", c.Name)
	}
}

// GetStyleFileName returns the style file name
func (c *UIComponent) GetStyleFileName() string {
	switch c.Framework {
	case ReactFramework:
		return fmt.Sprintf("%s.module.scss", c.Name)
	case VueFramework:
		return "" // Styles are inside .vue files
	case AngularFramework:
		return fmt.Sprintf("%s.component.scss", strings.ToLower(c.Name))
	default:
		return fmt.Sprintf("%s.module.scss", c.Name)
	}
}

// GetTestFileName returns the test file name
func (c *UIComponent) GetTestFileName() string {
	switch c.Framework {
	case ReactFramework:
		if c.TypeScript {
			return fmt.Sprintf("%s.test.tsx", c.Name)
		}
		return fmt.Sprintf("%s.test.jsx", c.Name)
	case VueFramework:
		return fmt.Sprintf("%s.spec.ts", c.Name)
	case AngularFramework:
		return fmt.Sprintf("%s.component.spec.ts", strings.ToLower(c.Name))
	default:
		return fmt.Sprintf("%s.test.tsx", c.Name)
	}
}

// GetStoryFileName returns the Storybook story file name
func (c *UIComponent) GetStoryFileName() string {
	if c.TypeScript {
		return fmt.Sprintf("%s.stories.tsx", c.Name)
	}
	return fmt.Sprintf("%s.stories.jsx", c.Name)
}

// GetDirectoryPath returns the component directory path based on atomic design
func (c *UIComponent) GetDirectoryPath() string {
	switch c.Type {
	case AtomType:
		return "src/components/atoms"
	case MoleculeType:
		return "src/components/molecules"
	case OrganismType:
		return "src/components/organisms"
	case TemplateType:
		return "src/components/templates"
	case PageType:
		return "src/pages"
	default:
		return "src/components"
	}
}

// AtomicDesignStructure defines the complete atomic design structure
type AtomicDesignStructure struct {
	ProjectName string        `json:"projectName"`
	Framework   Framework     `json:"framework"`
	TypeScript  bool          `json:"typescript"`
	Storybook   bool          `json:"storybook"`
	Atoms       []UIComponent `json:"atoms"`
	Molecules   []UIComponent `json:"molecules"`
	Organisms   []UIComponent `json:"organisms"`
	Templates   []UIComponent `json:"templates"`
	Pages       []UIComponent `json:"pages"`
}

// GetDefaultAtoms returns default atom components
func GetDefaultAtoms(framework Framework, typescript bool) []UIComponent {
	return []UIComponent{
		{
			Name:        "Button",
			Type:        AtomType,
			Framework:   framework,
			TypeScript:  typescript,
			HasStory:    true,
			HasTest:     true,
			HasStyles:   true,
			Description: "Reusable button component with variants",
			Props: []UIProp{
				{Name: "variant", Type: "string", Optional: true, Default: "primary", Description: "Button variant (primary, secondary, danger)"},
				{Name: "size", Type: "string", Optional: true, Default: "medium", Description: "Button size (small, medium, large)"},
				{Name: "disabled", Type: "boolean", Optional: true, Default: "false", Description: "Disable button interaction"},
				{Name: "onClick", Type: "function", Optional: true, Description: "Click handler function"},
				{Name: "children", Type: "ReactNode", Optional: false, Description: "Button content"},
			},
		},
		{
			Name:        "Input",
			Type:        AtomType,
			Framework:   framework,
			TypeScript:  typescript,
			HasStory:    true,
			HasTest:     true,
			HasStyles:   true,
			Description: "Text input component with validation",
			Props: []UIProp{
				{Name: "type", Type: "string", Optional: true, Default: "text", Description: "Input type"},
				{Name: "placeholder", Type: "string", Optional: true, Description: "Input placeholder text"},
				{Name: "value", Type: "string", Optional: true, Description: "Input value"},
				{Name: "onChange", Type: "function", Optional: true, Description: "Change handler function"},
				{Name: "error", Type: "string", Optional: true, Description: "Error message"},
				{Name: "disabled", Type: "boolean", Optional: true, Default: "false", Description: "Disable input interaction"},
			},
		},
		{
			Name:        "Label",
			Type:        AtomType,
			Framework:   framework,
			TypeScript:  typescript,
			HasStory:    true,
			HasTest:     true,
			HasStyles:   true,
			Description: "Form label component",
			Props: []UIProp{
				{Name: "htmlFor", Type: "string", Optional: true, Description: "Associated input ID"},
				{Name: "required", Type: "boolean", Optional: true, Default: "false", Description: "Show required indicator"},
				{Name: "children", Type: "ReactNode", Optional: false, Description: "Label content"},
			},
		},
	}
}

// GetDefaultMolecules returns default molecule components
func GetDefaultMolecules(framework Framework, typescript bool) []UIComponent {
	return []UIComponent{
		{
			Name:        "FormField",
			Type:        MoleculeType,
			Framework:   framework,
			TypeScript:  typescript,
			HasStory:    true,
			HasTest:     true,
			HasStyles:   true,
			Description: "Form field combining label and input",
			Props: []UIProp{
				{Name: "label", Type: "string", Optional: false, Description: "Field label"},
				{Name: "name", Type: "string", Optional: false, Description: "Field name"},
				{Name: "type", Type: "string", Optional: true, Default: "text", Description: "Input type"},
				{Name: "placeholder", Type: "string", Optional: true, Description: "Input placeholder"},
				{Name: "required", Type: "boolean", Optional: true, Default: "false", Description: "Required field"},
				{Name: "error", Type: "string", Optional: true, Description: "Error message"},
			},
		},
		{
			Name:        "Card",
			Type:        MoleculeType,
			Framework:   framework,
			TypeScript:  typescript,
			HasStory:    true,
			HasTest:     true,
			HasStyles:   true,
			Description: "Reusable card component",
			Props: []UIProp{
				{Name: "title", Type: "string", Optional: true, Description: "Card title"},
				{Name: "subtitle", Type: "string", Optional: true, Description: "Card subtitle"},
				{Name: "actions", Type: "ReactNode", Optional: true, Description: "Card action buttons"},
				{Name: "children", Type: "ReactNode", Optional: false, Description: "Card content"},
			},
		},
	}
}

// GetDefaultOrganisms returns default organism components
func GetDefaultOrganisms(framework Framework, typescript bool) []UIComponent {
	return []UIComponent{
		{
			Name:        "Header",
			Type:        OrganismType,
			Framework:   framework,
			TypeScript:  typescript,
			HasStory:    true,
			HasTest:     true,
			HasStyles:   true,
			Description: "Application header with navigation",
			Props: []UIProp{
				{Name: "logo", Type: "string", Optional: true, Description: "Logo URL or text"},
				{Name: "navigation", Type: "NavItem[]", Optional: true, Description: "Navigation items"},
				{Name: "user", Type: "User", Optional: true, Description: "Current user info"},
			},
		},
		{
			Name:        "Sidebar",
			Type:        OrganismType,
			Framework:   framework,
			TypeScript:  typescript,
			HasStory:    true,
			HasTest:     true,
			HasStyles:   true,
			Description: "Navigation sidebar component",
			Props: []UIProp{
				{Name: "items", Type: "SidebarItem[]", Optional: false, Description: "Sidebar navigation items"},
				{Name: "collapsed", Type: "boolean", Optional: true, Default: "false", Description: "Sidebar collapsed state"},
				{Name: "onToggle", Type: "function", Optional: true, Description: "Toggle handler"},
			},
		},
	}
}