package templates

import (
	"fmt"
	"strings"

	"github.com/vibercode/cli/internal/models"
)

// GetPackageJSON generates package.json for UI project
func GetPackageJSON(config *models.AtomicDesignStructure) string {
	dependencies := getFrameworkDependencies(config.Framework, config.TypeScript)
	devDependencies := getDevDependencies(config.Framework, config.TypeScript, config.Storybook)
	scripts := getScripts(config.Framework, config.Storybook)

	return fmt.Sprintf(`{
  "name": "%s",
  "version": "1.0.0",
  "description": "UI components following Atomic Design methodology",
  "main": "src/components/index.ts",
  "scripts": {
%s
  },
  "dependencies": {
%s
  },
  "devDependencies": {
%s
  },
  "peerDependencies": {
    "react": "^18.0.0",
    "react-dom": "^18.0.0"
  }
}`, config.ProjectName, scripts, dependencies, devDependencies)
}

// GetTSConfig generates TypeScript configuration
func GetTSConfig(config *models.AtomicDesignStructure) string {
	return `{
  "compilerOptions": {
    "target": "ES2020",
    "lib": ["DOM", "DOM.Iterable", "ES6"],
    "allowJs": true,
    "skipLibCheck": true,
    "esModuleInterop": true,
    "allowSyntheticDefaultImports": true,
    "strict": true,
    "forceConsistentCasingInFileNames": true,
    "noFallthroughCasesInSwitch": true,
    "module": "ESNext",
    "moduleResolution": "node",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "react-jsx",
    "declaration": true,
    "outDir": "dist",
    "baseUrl": "src",
    "paths": {
      "@/*": ["*"],
      "@/components/*": ["components/*"],
      "@/styles/*": ["styles/*"],
      "@/types/*": ["types/*"]
    }
  },
  "include": [
    "src/**/*",
    "src/**/*.stories.*"
  ],
  "exclude": [
    "node_modules",
    "dist",
    "build"
  ]
}`
}

// GetStorybookConfig generates Storybook configuration  
func GetStorybookConfig(config *models.AtomicDesignStructure) string {
	return `module.exports = {
  stories: [
    "../src/**/*.stories.@(js|jsx|ts|tsx|mdx)"
  ],
  addons: [
    "@storybook/addon-essentials",
    "@storybook/addon-interactions",
    "@storybook/addon-a11y",
    "@storybook/addon-docs"
  ],
  framework: {
    name: "@storybook/react-webpack5",
    options: {}
  },
  features: {
    buildStoriesJson: true
  },
  webpackFinal: async (config) => {
    // Add SCSS support
    config.module.rules.push({
      test: /\.scss$/,
      use: [
        'style-loader',
        {
          loader: 'css-loader',
          options: {
            modules: {
              localIdentName: '[name]__[local]--[hash:base64:5]'
            }
          }
        },
        'sass-loader'
      ]
    });
    
    return config;
  }
}`
}

// GetGlobalStyles generates global SCSS styles
func GetGlobalStyles(config *models.AtomicDesignStructure) string {
	return `// Global Variables
:root {
  // Colors
  --primary-color: #007bff;
  --secondary-color: #6c757d;
  --success-color: #28a745;
  --danger-color: #dc3545;
  --warning-color: #ffc107;
  --info-color: #17a2b8;
  --light-color: #f8f9fa;
  --dark-color: #343a40;
  
  // Typography
  --font-family-base: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', sans-serif;
  --font-size-base: 1rem;
  --font-weight-normal: 400;
  --font-weight-bold: 700;
  --line-height-base: 1.5;
  
  // Spacing
  --spacing-xs: 0.25rem;
  --spacing-sm: 0.5rem;
  --spacing-md: 1rem;
  --spacing-lg: 1.5rem;
  --spacing-xl: 3rem;
  
  // Border radius
  --border-radius-sm: 0.125rem;
  --border-radius: 0.25rem;
  --border-radius-lg: 0.375rem;
  --border-radius-xl: 0.5rem;
  
  // Shadows
  --shadow-sm: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  --shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06);
  --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
}

// Reset
*,
*::before,
*::after {
  box-sizing: border-box;
}

body {
  margin: 0;
  font-family: var(--font-family-base);
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-normal);
  line-height: var(--line-height-base);
  color: var(--dark-color);
  background-color: var(--light-color);
}

// Utility Classes
.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border: 0;
}

.text-center { text-align: center; }
.text-left { text-align: left; }
.text-right { text-align: right; }

.mb-0 { margin-bottom: 0; }
.mb-1 { margin-bottom: var(--spacing-xs); }
.mb-2 { margin-bottom: var(--spacing-sm); }
.mb-3 { margin-bottom: var(--spacing-md); }
.mb-4 { margin-bottom: var(--spacing-lg); }
.mb-5 { margin-bottom: var(--spacing-xl); }

.d-none { display: none; }
.d-block { display: block; }
.d-flex { display: flex; }
.d-inline-block { display: inline-block; }`
}

// GetComponentStyles generates component-specific styles
func GetComponentStyles(component models.UIComponent) string {
	switch component.Name {
	case "Button":
		return getButtonStyles()
	case "Input":
		return getInputStyles()
	case "Label":
		return getLabelStyles()
	case "FormField":
		return getFormFieldStyles()
	case "Card":
		return getCardStyles()
	case "Header":
		return getHeaderStyles()
	case "Sidebar":
		return getSidebarStyles()
	default:
		return getDefaultComponentStyles(component)
	}
}

// getButtonStyles returns Button component styles
func getButtonStyles() string {
	return `.button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-sm) var(--spacing-md);
  border: 1px solid transparent;
  border-radius: var(--border-radius);
  font-family: var(--font-family-base);
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-normal);
  line-height: var(--line-height-base);
  text-decoration: none;
  cursor: pointer;
  transition: all 0.15s ease-in-out;
  
  &:focus {
    outline: 0;
    box-shadow: 0 0 0 0.2rem rgba(0, 123, 255, 0.25);
  }
  
  &:disabled {
    opacity: 0.65;
    cursor: not-allowed;
  }
  
  // Variants
  &.primary {
    color: white;
    background-color: var(--primary-color);
    border-color: var(--primary-color);
    
    &:hover:not(:disabled) {
      background-color: #0056b3;
      border-color: #004085;
    }
  }
  
  &.secondary {
    color: white;
    background-color: var(--secondary-color);
    border-color: var(--secondary-color);
    
    &:hover:not(:disabled) {
      background-color: #545b62;
      border-color: #4e555b;
    }
  }
  
  &.danger {
    color: white;
    background-color: var(--danger-color);
    border-color: var(--danger-color);
    
    &:hover:not(:disabled) {
      background-color: #c82333;
      border-color: #bd2130;
    }
  }
  
  // Sizes
  &.small {
    padding: var(--spacing-xs) var(--spacing-sm);
    font-size: 0.875rem;
  }
  
  &.large {
    padding: var(--spacing-md) var(--spacing-lg);
    font-size: 1.125rem;
  }
}`
}

// getInputStyles returns Input component styles
func getInputStyles() string {
	return `.inputContainer {
  position: relative;
}

.input {
  display: block;
  width: 100%;
  padding: var(--spacing-sm) var(--spacing-md);
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-normal);
  line-height: var(--line-height-base);
  color: var(--dark-color);
  background-color: white;
  border: 1px solid #ced4da;
  border-radius: var(--border-radius);
  transition: border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out;
  
  &:focus {
    color: var(--dark-color);
    background-color: white;
    border-color: #80bdff;
    outline: 0;
    box-shadow: 0 0 0 0.2rem rgba(0, 123, 255, 0.25);
  }
  
  &:disabled {
    background-color: #e9ecef;
    opacity: 1;
    cursor: not-allowed;
  }
  
  &.error {
    border-color: var(--danger-color);
    
    &:focus {
      border-color: var(--danger-color);
      box-shadow: 0 0 0 0.2rem rgba(220, 53, 69, 0.25);
    }
  }
}

.errorMessage {
  display: block;
  width: 100%;
  margin-top: var(--spacing-xs);
  font-size: 0.875rem;
  color: var(--danger-color);
}`
}

// getLabelStyles returns Label component styles
func getLabelStyles() string {
	return `.label {
  display: inline-block;
  margin-bottom: var(--spacing-xs);
  font-weight: var(--font-weight-bold);
  color: var(--dark-color);
  
  &.required {
    .asterisk {
      color: var(--danger-color);
      margin-left: var(--spacing-xs);
    }
  }
}`
}

// getFormFieldStyles returns FormField component styles
func getFormFieldStyles() string {
	return `.formField {
  margin-bottom: var(--spacing-md);
}`
}

// getCardStyles returns Card component styles
func getCardStyles() string {
	return `.card {
  position: relative;
  display: flex;
  flex-direction: column;
  min-width: 0;
  word-wrap: break-word;
  background-color: white;
  border: 1px solid rgba(0, 0, 0, 0.125);
  border-radius: var(--border-radius);
  box-shadow: var(--shadow-sm);
}

.header {
  padding: var(--spacing-md);
  border-bottom: 1px solid rgba(0, 0, 0, 0.125);
}

.title {
  margin: 0 0 var(--spacing-xs) 0;
  font-size: 1.25rem;
  font-weight: var(--font-weight-bold);
}

.subtitle {
  margin: 0;
  color: var(--secondary-color);
  font-size: 0.875rem;
}

.content {
  flex: 1 1 auto;
  padding: var(--spacing-md);
}

.actions {
  padding: var(--spacing-md);
  border-top: 1px solid rgba(0, 0, 0, 0.125);
  display: flex;
  gap: var(--spacing-sm);
}`
}

// getHeaderStyles returns Header component styles
func getHeaderStyles() string {
	return `.header {
  background-color: white;
  border-bottom: 1px solid #e9ecef;
  box-shadow: var(--shadow-sm);
}

.container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  max-width: 1200px;
  margin: 0 auto;
  padding: var(--spacing-md);
}

.logo {
  img {
    height: 40px;
    width: auto;
  }
}

.navigation {
  display: flex;
  gap: var(--spacing-lg);
}

.navItem {
  color: var(--dark-color);
  text-decoration: none;
  font-weight: var(--font-weight-normal);
  transition: color 0.15s ease-in-out;
  
  &:hover {
    color: var(--primary-color);
  }
}

.userInfo {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}`
}

// getSidebarStyles returns Sidebar component styles
func getSidebarStyles() string {
	return `.sidebar {
  width: 250px;
  height: 100vh;
  background-color: var(--dark-color);
  color: white;
  transition: width 0.3s ease;
  position: fixed;
  left: 0;
  top: 0;
  overflow-x: hidden;
  
  &.collapsed {
    width: 60px;
  }
}

.toggleButton {
  width: 100%;
  padding: var(--spacing-md);
  background: none;
  border: none;
  color: white;
  font-size: 1.2rem;
  cursor: pointer;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  
  &:hover {
    background-color: rgba(255, 255, 255, 0.1);
  }
}

.navigation {
  padding: var(--spacing-md) 0;
}

.navItem {
  display: flex;
  align-items: center;
  padding: var(--spacing-md);
  color: white;
  text-decoration: none;
  transition: background-color 0.15s ease;
  
  &:hover {
    background-color: rgba(255, 255, 255, 0.1);
  }
  
  &.active {
    background-color: var(--primary-color);
  }
}

.icon {
  width: 20px;
  margin-right: var(--spacing-md);
  text-align: center;
}

.label {
  white-space: nowrap;
  overflow: hidden;
}`
}

// getDefaultComponentStyles returns default component styles
func getDefaultComponentStyles(component models.UIComponent) string {
	className := strings.ToLower(component.Name)
	return fmt.Sprintf(`.%s {
  /* Add your %s component styles here */
  padding: var(--spacing-md);
  border: 1px solid #e9ecef;
  border-radius: var(--border-radius);
  background-color: white;
}`, className, component.Name)
}

// Helper functions for package.json generation

func getFrameworkDependencies(framework models.Framework, typescript bool) string {
	switch framework {
	case models.ReactFramework:
		deps := []string{
			`    "react": "^18.2.0"`,
			`    "react-dom": "^18.2.0"`,
		}
		if typescript {
			deps = append(deps, `    "clsx": "^2.0.0"`)
		}
		return strings.Join(deps, ",\n")
	case models.VueFramework:
		return `    "vue": "^3.3.0"`
	case models.AngularFramework:
		return `    "@angular/core": "^16.0.0",
    "@angular/common": "^16.0.0"`
	default:
		return `    "react": "^18.2.0",
    "react-dom": "^18.2.0"`
	}
}

func getDevDependencies(framework models.Framework, typescript bool, storybook bool) string {
	var deps []string

	// Common dev dependencies
	deps = append(deps, `    "sass": "^1.69.0"`)

	if typescript {
		deps = append(deps, 
			`    "typescript": "^5.2.0"`,
			`    "@types/node": "^20.8.0"`)
	}

	// Framework-specific dev dependencies
	switch framework {
	case models.ReactFramework:
		deps = append(deps,
			`    "@vitejs/plugin-react": "^4.0.0"`,
			`    "vite": "^4.4.0"`)
		if typescript {
			deps = append(deps, 
				`    "@types/react": "^18.2.0"`,
				`    "@types/react-dom": "^18.2.0"`)
		}
	case models.VueFramework:
		deps = append(deps, `    "@vitejs/plugin-vue": "^4.0.0"`)
	}

	// Storybook dependencies
	if storybook {
		deps = append(deps,
			`    "@storybook/react": "^7.5.0"`,
			`    "@storybook/addon-essentials": "^7.5.0"`,
			`    "@storybook/addon-interactions": "^7.5.0"`,
			`    "@storybook/addon-a11y": "^7.5.0"`)
	}

	return strings.Join(deps, ",\n")
}

func getScripts(framework models.Framework, storybook bool) string {
	scripts := []string{
		`    "build": "tsc && vite build"`,
		`    "dev": "vite"`,
		`    "preview": "vite preview"`,
		`    "lint": "eslint . --ext ts,tsx --report-unused-disable-directives --max-warnings 0"`,
	}

	if storybook {
		scripts = append(scripts,
			`    "storybook": "storybook dev -p 6006"`,
			`    "build-storybook": "storybook build"`)
	}

	return strings.Join(scripts, ",\n")
}