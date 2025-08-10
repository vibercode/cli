package templates

import (
	"fmt"
	"strings"

	"github.com/vibercode/cli/internal/models"
)

// GetComponentTest generates test file for component
func GetComponentTest(component models.UIComponent) string {
	switch component.Framework {
	case models.ReactFramework:
		return getReactTest(component)
	case models.VueFramework:
		return getVueTest(component)
	case models.AngularFramework:
		return getAngularTest(component)
	default:
		return getReactTest(component)
	}
}

// GetComponentStory generates Storybook story for component
func GetComponentStory(component models.UIComponent) string {
	return getStorybookStory(component)
}

// GetComponentIndex generates index file for component
func GetComponentIndex(component models.UIComponent) string {
	return fmt.Sprintf(`export { default } from './%s';
export type { %sProps } from './%s';
`, component.GetFileName(), component.Name, component.GetFileName())
}

// GetLevelIndex generates index file for atomic design level
func GetLevelIndex(level string, config *models.AtomicDesignStructure) string {
	var exports []string
	
	switch level {
	case "atoms":
		atoms := models.GetDefaultAtoms(config.Framework, config.TypeScript)
		for _, atom := range atoms {
			exports = append(exports, fmt.Sprintf("export { default as %s } from './%s';", atom.Name, atom.Name))
		}
	case "molecules":
		molecules := models.GetDefaultMolecules(config.Framework, config.TypeScript)
		for _, molecule := range molecules {
			exports = append(exports, fmt.Sprintf("export { default as %s } from './%s';", molecule.Name, molecule.Name))
		}
	case "organisms":
		organisms := models.GetDefaultOrganisms(config.Framework, config.TypeScript)
		for _, organism := range organisms {
			exports = append(exports, fmt.Sprintf("export { default as %s } from './%s';", organism.Name, organism.Name))
		}
	case "templates":
		exports = append(exports, "// Export your template components here")
	}
	
	return strings.Join(exports, "\n") + "\n"
}

// GetMainComponentsIndex generates main components index
func GetMainComponentsIndex(config *models.AtomicDesignStructure) string {
	return `// Atomic Design Components
export * from './atoms';
export * from './molecules';
export * from './organisms';
export * from './templates';

// Re-export everything for convenience
export * as Atoms from './atoms';
export * as Molecules from './molecules';
export * as Organisms from './organisms';
export * as Templates from './templates';
`
}

// getReactTest generates React component test
func getReactTest(component models.UIComponent) string {
	testContent := generateTestContent(component)
	
	if component.TypeScript {
		return fmt.Sprintf(`import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import %s from './%s';

describe('%s', () => {
%s
});
`, component.Name, component.GetFileName(), component.Name, testContent)
	}

	return fmt.Sprintf(`import React from 'react';
import { render, screen } from '@testing-library/react';
import %s from './%s';

describe('%s', () => {
%s
});
`, component.Name, component.GetFileName(), component.Name, testContent)
}

// getVueTest generates Vue component test
func getVueTest(component models.UIComponent) string {
	return fmt.Sprintf(`import { describe, it, expect } from 'vitest';
import { mount } from '@vue/test-utils';
import %s from './%s';

describe('%s', () => {
  it('renders properly', () => {
    const wrapper = mount(%s, {
      props: {
        // Add required props here
      }
    });
    
    expect(wrapper.text()).toContain('%s');
  });
  
  it('handles props correctly', () => {
    const wrapper = mount(%s, {
      props: {
        // Add test props here
      }
    });
    
    // Add assertions
    expect(wrapper.exists()).toBe(true);
  });
});
`, component.Name, component.GetFileName(), component.Name, 
component.Name, component.Name, component.Name)
}

// getAngularTest generates Angular component test
func getAngularTest(component models.UIComponent) string {
	componentName := strings.ToLower(component.Name)
	return fmt.Sprintf(`import { ComponentFixture, TestBed } from '@angular/core/testing';
import { %sComponent } from './%s.component';

describe('%sComponent', () => {
  let component: %sComponent;
  let fixture: ComponentFixture<%sComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [%sComponent]
    }).compileComponents();

    fixture = TestBed.createComponent(%sComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should render component content', () => {
    const compiled = fixture.nativeElement as HTMLElement;
    expect(compiled.querySelector('h2')?.textContent).toContain('%s Component');
  });
});
`, component.Name, componentName, component.Name, component.Name, 
component.Name, component.Name, component.Name, component.Name)
}

// generateTestContent generates test cases based on component type
func generateTestContent(component models.UIComponent) string {
	var tests []string
	
	// Basic render test
	tests = append(tests, `  it('renders without crashing', () => {
    render(<%s />);
  });`)
	
	// Component-specific tests
	switch component.Name {
	case "Button":
		tests = append(tests, 
			`  it('displays children content', () => {
    render(<Button>Click me</Button>);
    expect(screen.getByRole('button')).toHaveTextContent('Click me');
  });`,
			`  it('handles click events', () => {
    const handleClick = vi.fn();
    render(<Button onClick={handleClick}>Click me</Button>);
    
    const button = screen.getByRole('button');
    button.click();
    
    expect(handleClick).toHaveBeenCalledTimes(1);
  });`,
			`  it('applies variant classes', () => {
    render(<Button variant="secondary">Button</Button>);
    expect(screen.getByRole('button')).toHaveClass('secondary');
  });`)
		
	case "Input":
		tests = append(tests,
			`  it('displays placeholder text', () => {
    render(<Input placeholder="Enter text" />);
    expect(screen.getByPlaceholderText('Enter text')).toBeInTheDocument();
  });`,
			`  it('shows error message when provided', () => {
    render(<Input error="This field is required" />);
    expect(screen.getByText('This field is required')).toBeInTheDocument();
  });`)
		
	case "Label":
		tests = append(tests,
			`  it('displays label text', () => {
    render(<Label>Username</Label>);
    expect(screen.getByText('Username')).toBeInTheDocument();
  });`,
			`  it('shows required indicator when required', () => {
    render(<Label required>Username</Label>);
    expect(screen.getByText('*')).toBeInTheDocument();
  });`)
		
	case "Card":
		tests = append(tests,
			`  it('displays title and subtitle', () => {
    render(<Card title="Card Title" subtitle="Card Subtitle">Content</Card>);
    expect(screen.getByText('Card Title')).toBeInTheDocument();
    expect(screen.getByText('Card Subtitle')).toBeInTheDocument();
  });`)
	}
	
	// Format tests with component name
	var formattedTests []string
	for _, test := range tests {
		formattedTests = append(formattedTests, fmt.Sprintf(test, component.Name))
	}
	
	return strings.Join(formattedTests, "\n\n")
}

// getStorybookStory generates Storybook story
func getStorybookStory(component models.UIComponent) string {
	args := generateStoryArgs(component)
	variants := generateStoryVariants(component)
	
	return fmt.Sprintf(`import type { Meta, StoryObj } from '@storybook/react';
import %s from './%s';

const meta: Meta<typeof %s> = {
  title: '%s/%s',
  component: %s,
  parameters: {
    layout: 'centered',
    docs: {
      description: {
        component: '%s'
      }
    }
  },
  tags: ['autodocs'],
  argTypes: {
%s
  }
};

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
%s
  }
};

%s
`, component.Name, component.GetFileName(), component.Name, 
strings.Title(string(component.Type)), component.Name, component.Name,
component.Description, generateArgTypes(component), args, variants)
}

// generateStoryArgs generates default args for story
func generateStoryArgs(component models.UIComponent) string {
	var args []string
	
	for _, prop := range component.Props {
		if prop.Default != "" {
			args = append(args, fmt.Sprintf("    %s: %s", prop.Name, prop.Default))
		} else {
			switch prop.Type {
			case "string":
				if prop.Name == "children" {
					args = append(args, fmt.Sprintf("    %s: '%s content'", prop.Name, component.Name))
				} else {
					args = append(args, fmt.Sprintf("    %s: 'Sample %s'", prop.Name, prop.Name))
				}
			case "boolean":
				args = append(args, fmt.Sprintf("    %s: false", prop.Name))
			case "number":
				args = append(args, fmt.Sprintf("    %s: 0", prop.Name))
			}
		}
	}
	
	return strings.Join(args, ",\n")
}

// generateArgTypes generates argTypes for Storybook controls
func generateArgTypes(component models.UIComponent) string {
	var argTypes []string
	
	for _, prop := range component.Props {
		var control string
		switch prop.Type {
		case "string":
			if strings.Contains(prop.Name, "variant") || strings.Contains(prop.Name, "size") {
				control = "{ control: 'select', options: ['primary', 'secondary', 'danger'] }"
			} else {
				control = "{ control: 'text' }"
			}
		case "boolean":
			control = "{ control: 'boolean' }"
		case "number":
			control = "{ control: 'number' }"
		case "function":
			control = "{ action: 'clicked' }"
		default:
			control = "{ control: 'object' }"
		}
		
		argType := fmt.Sprintf("    %s: %s", prop.Name, control)
		if prop.Description != "" {
			argType = fmt.Sprintf("    %s: {\n      ...%s,\n      description: '%s'\n    }", prop.Name, control, prop.Description)
		}
		
		argTypes = append(argTypes, argType)
	}
	
	return strings.Join(argTypes, ",\n")
}

// generateStoryVariants generates additional story variants
func generateStoryVariants(component models.UIComponent) string {
	var variants []string
	
	switch component.Name {
	case "Button":
		variants = append(variants,
			`export const Primary: Story = {
  args: {
    variant: 'primary',
    children: 'Primary Button'
  }
};

export const Secondary: Story = {
  args: {
    variant: 'secondary',
    children: 'Secondary Button'
  }
};

export const Large: Story = {
  args: {
    size: 'large',
    children: 'Large Button'
  }
};

export const Disabled: Story = {
  args: {
    disabled: true,
    children: 'Disabled Button'
  }
};`)
		
	case "Input":
		variants = append(variants,
			`export const WithPlaceholder: Story = {
  args: {
    placeholder: 'Enter your text here...'
  }
};

export const WithError: Story = {
  args: {
    error: 'This field is required',
    value: ''
  }
};

export const Disabled: Story = {
  args: {
    disabled: true,
    value: 'Disabled input'
  }
};`)
		
	case "Card":
		variants = append(variants,
			`export const WithHeader: Story = {
  args: {
    title: 'Card Title',
    subtitle: 'Card subtitle',
    children: 'This is the card content area.'
  }
};

export const SimpleCard: Story = {
  args: {
    children: 'Simple card with just content.'
  }
};`)
	}
	
	return strings.Join(variants, "\n\n")
}