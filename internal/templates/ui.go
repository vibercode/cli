package templates

import (
	"fmt"
	"strings"

	"github.com/vibercode/cli/internal/models"
)

// GetComponentTemplate returns the component template based on framework
func GetComponentTemplate(component models.UIComponent) string {
	switch component.Framework {
	case models.ReactFramework:
		return getReactComponent(component)
	case models.VueFramework:
		return getVueComponent(component)
	case models.AngularFramework:
		return getAngularComponent(component)
	default:
		return getReactComponent(component)
	}
}

// getReactComponent generates React component template
func getReactComponent(component models.UIComponent) string {
	propsInterface := generatePropsInterface(component)
	propsDestructuring := generatePropsDestructuring(component)
	componentContent := generateReactComponentContent(component)

	if component.TypeScript {
		return fmt.Sprintf(`import React from 'react';
%s
import styles from './%s';

%s

interface %sProps {
%s
}

const %s: React.FC<%sProps> = ({
%s
}) => {
  return (
%s
  );
};

export default %s;
`, getImports(component), component.GetStyleFileName(), generateTypeDefinitions(component), 
component.Name, propsInterface, component.Name, component.Name, propsDestructuring, 
componentContent, component.Name)
	}

	return fmt.Sprintf(`import React from 'react';
%s
import styles from './%s';

const %s = ({
%s
}) => {
  return (
%s
  );
};

export default %s;
`, getImports(component), component.GetStyleFileName(), component.Name, 
propsDestructuring, componentContent, component.Name)
}

// getVueComponent generates Vue component template
func getVueComponent(component models.UIComponent) string {
	propsDefinition := generateVueProps(component)
	template := generateVueTemplate(component)
	
	if component.TypeScript {
		return fmt.Sprintf(`<template>
%s
</template>

<script setup lang="ts">
%s

interface Props {
%s
}

const props = defineProps<Props>();
</script>

<style scoped lang="scss">
%s
</style>
`, template, getVueImports(component), generatePropsInterface(component), generateVueStyles(component))
	}

	return fmt.Sprintf(`<template>
%s
</template>

<script setup>
%s

const props = defineProps({
%s
});
</script>

<style scoped lang="scss">
%s
</style>
`, template, getVueImports(component), propsDefinition, generateVueStyles(component))
}

// getAngularComponent generates Angular component template
func getAngularComponent(component models.UIComponent) string {
	componentName := strings.ToLower(component.Name)
	selector := fmt.Sprintf("app-%s", componentName)
	
	return fmt.Sprintf(`import { Component, Input } from '@angular/core';
%s

@Component({
  selector: '%s',
  templateUrl: './%s.component.html',
  styleUrls: ['./%s.component.scss']
})
export class %sComponent {
%s

  constructor() {}
}
`, getAngularImports(component), selector, componentName, componentName, 
component.Name, generateAngularInputs(component))
}

// generatePropsInterface generates TypeScript props interface
func generatePropsInterface(component models.UIComponent) string {
	var props []string
	for _, prop := range component.Props {
		optional := ""
		if prop.Optional {
			optional = "?"
		}
		propLine := fmt.Sprintf("  %s%s: %s;", prop.Name, optional, prop.Type)
		if prop.Description != "" {
			propLine = fmt.Sprintf("  /** %s */\n%s", prop.Description, propLine)
		}
		props = append(props, propLine)
	}
	return strings.Join(props, "\n")
}

// generatePropsDestructuring generates props destructuring for React
func generatePropsDestructuring(component models.UIComponent) string {
	var props []string
	for _, prop := range component.Props {
		if prop.Default != "" {
			props = append(props, fmt.Sprintf("%s = %s", prop.Name, prop.Default))
		} else {
			props = append(props, prop.Name)
		}
	}
	return "  " + strings.Join(props, ",\n  ")
}

// generateReactComponentContent generates React component JSX content
func generateReactComponentContent(component models.UIComponent) string {
	switch component.Type {
	case models.AtomType:
		return generateAtomContent(component)
	case models.MoleculeType:
		return generateMoleculeContent(component)
	case models.OrganismType:
		return generateOrganismContent(component)
	default:
		return generateDefaultContent(component)
	}
}

// generateAtomContent generates content for atom components
func generateAtomContent(component models.UIComponent) string {
	switch component.Name {
	case "Button":
		return `    <button 
      className={` + "`${styles.button} ${styles[variant]} ${styles[size]} ${disabled ? styles.disabled : ''}`" + `}
      disabled={disabled}
      onClick={onClick}
    >
      {children}
    </button>`
	case "Input":
		return `    <div className={styles.inputContainer}>
      <input
        type={type}
        placeholder={placeholder}
        value={value}
        onChange={onChange}
        disabled={disabled}
        className={` + "`${styles.input} ${error ? styles.error : ''}`" + `}
      />
      {error && <span className={styles.errorMessage}>{error}</span>}
    </div>`
	case "Label":
		return `    <label 
      htmlFor={htmlFor}
      className={` + "`${styles.label} ${required ? styles.required : ''}`" + `}
    >
      {children}
      {required && <span className={styles.asterisk}>*</span>}
    </label>`
	default:
		return `    <div className={styles.container}>
      {children}
    </div>`
	}
}

// generateMoleculeContent generates content for molecule components
func generateMoleculeContent(component models.UIComponent) string {
	switch component.Name {
	case "FormField":
		return `    <div className={styles.formField}>
      <Label htmlFor={name} required={required}>
        {label}
      </Label>
      <Input
        name={name}
        type={type}
        placeholder={placeholder}
        error={error}
      />
    </div>`
	case "Card":
		return `    <div className={styles.card}>
      {title && (
        <div className={styles.header}>
          <h3 className={styles.title}>{title}</h3>
          {subtitle && <p className={styles.subtitle}>{subtitle}</p>}
        </div>
      )}
      <div className={styles.content}>
        {children}
      </div>
      {actions && (
        <div className={styles.actions}>
          {actions}
        </div>
      )}
    </div>`
	default:
		return generateDefaultContent(component)
	}
}

// generateOrganismContent generates content for organism components
func generateOrganismContent(component models.UIComponent) string {
	switch component.Name {
	case "Header":
		return `    <header className={styles.header}>
      <div className={styles.container}>
        {logo && (
          <div className={styles.logo}>
            {typeof logo === 'string' ? <img src={logo} alt="Logo" /> : logo}
          </div>
        )}
        <nav className={styles.navigation}>
          {navigation?.map((item, index) => (
            <a key={index} href={item.href} className={styles.navItem}>
              {item.label}
            </a>
          ))}
        </nav>
        {user && (
          <div className={styles.userInfo}>
            <span>{user.name}</span>
          </div>
        )}
      </div>
    </header>`
	case "Sidebar":
		return `    <aside className={` + "`${styles.sidebar} ${collapsed ? styles.collapsed : ''}`" + `}>
      <button 
        className={styles.toggleButton}
        onClick={onToggle}
        aria-label="Toggle sidebar"
      >
        ☰
      </button>
      <nav className={styles.navigation}>
        {items.map((item, index) => (
          <a
            key={index}
            href={item.href}
            className={` + "`${styles.navItem} ${item.active ? styles.active : ''}`" + `}
          >
            {item.icon && <span className={styles.icon}>{item.icon}</span>}
            {!collapsed && <span className={styles.label}>{item.label}</span>}
          </a>
        ))}
      </nav>
    </aside>`
	default:
		return generateDefaultContent(component)
	}
}

// generateDefaultContent generates default component content
func generateDefaultContent(component models.UIComponent) string {
	return fmt.Sprintf(`    <div className={styles.%s}>
      <h2>%s Component</h2>
      <p>%s</p>
    </div>`, strings.ToLower(component.Name), component.Name, component.Description)
}

// generateVueProps generates Vue props definition
func generateVueProps(component models.UIComponent) string {
	var props []string
	for _, prop := range component.Props {
		required := "false"
		if !prop.Optional {
			required = "true"
		}
		
		propDef := fmt.Sprintf("  %s: {\n    type: %s,\n    required: %s", 
			prop.Name, getVueType(prop.Type), required)
		
		if prop.Default != "" {
			propDef += fmt.Sprintf(",\n    default: %s", prop.Default)
		}
		
		propDef += "\n  }"
		props = append(props, propDef)
	}
	return strings.Join(props, ",\n")
}

// generateVueTemplate generates Vue template
func generateVueTemplate(component models.UIComponent) string {
	return fmt.Sprintf(`  <div class="%s">
    <h2>%s Component</h2>
    <p>%s</p>
  </div>`, strings.ToLower(component.Name), component.Name, component.Description)
}

// generateAngularInputs generates Angular @Input properties
func generateAngularInputs(component models.UIComponent) string {
	var inputs []string
	for _, prop := range component.Props {
		input := fmt.Sprintf("  @Input() %s", prop.Name)
		if prop.Optional {
			input += "?"
		}
		input += fmt.Sprintf(": %s;", getAngularType(prop.Type))
		
		if prop.Description != "" {
			input = fmt.Sprintf("  /** %s */\n%s", prop.Description, input)
		}
		inputs = append(inputs, input)
	}
	return strings.Join(inputs, "\n")
}

// getImports returns necessary imports for React components
func getImports(component models.UIComponent) string {
	var imports []string
	
	// Add component-specific imports
	switch component.Type {
	case models.MoleculeType:
		if component.Name == "FormField" {
			imports = append(imports, "import { Label, Input } from '../atoms';")
		}
	case models.OrganismType:
		imports = append(imports, "import { Button } from '../atoms';")
	}
	
	if len(imports) > 0 {
		return strings.Join(imports, "\n")
	}
	return ""
}

// getVueImports returns necessary imports for Vue components
func getVueImports(component models.UIComponent) string {
	return ""
}

// getAngularImports returns necessary imports for Angular components
func getAngularImports(component models.UIComponent) string {
	return ""
}

// generateTypeDefinitions generates additional type definitions
func generateTypeDefinitions(component models.UIComponent) string {
	var types []string
	
	switch component.Name {
	case "Header":
		types = append(types, `interface NavItem {
  label: string;
  href: string;
}

interface User {
  name: string;
  email?: string;
  avatar?: string;
}`)
	case "Sidebar":
		types = append(types, `interface SidebarItem {
  label: string;
  href: string;
  icon?: React.ReactNode;
  active?: boolean;
}`)
	}
	
	return strings.Join(types, "\n\n")
}

// getVueType converts TypeScript type to Vue prop type
func getVueType(tsType string) string {
	switch tsType {
	case "string":
		return "String"
	case "number":
		return "Number"
	case "boolean":
		return "Boolean"
	case "function":
		return "Function"
	case "ReactNode":
		return "Object"  // Vue doesn't have ReactNode
	default:
		return "Object"
	}
}

// getAngularType converts TypeScript type to Angular type
func getAngularType(tsType string) string {
	switch tsType {
	case "ReactNode":
		return "any"
	case "function":
		return "() => void"
	default:
		return tsType
	}
}

// generateVueStyles generates Vue component styles
func generateVueStyles(component models.UIComponent) string {
	return fmt.Sprintf(`.%s {
  /* Component styles */
}`, strings.ToLower(component.Name))
}