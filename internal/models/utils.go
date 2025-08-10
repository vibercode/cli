package models

import (
	"regexp"
	"strings"
	"unicode"
)

// ToPascalCase converts a string to PascalCase
func ToPascalCase(s string) string {
	words := splitIntoWords(s)
	var result []string
	for _, word := range words {
		if word != "" {
			result = append(result, strings.Title(strings.ToLower(word)))
		}
	}
	return strings.Join(result, "")
}

// ToCamelCase converts a string to camelCase
func ToCamelCase(s string) string {
	words := splitIntoWords(s)
	if len(words) == 0 {
		return ""
	}

	var result []string
	result = append(result, strings.ToLower(words[0]))

	for _, word := range words[1:] {
		if word != "" {
			result = append(result, strings.Title(strings.ToLower(word)))
		}
	}
	return strings.Join(result, "")
}

// ToSnakeCase converts a string to snake_case
func ToSnakeCase(s string) string {
	words := splitIntoWords(s)
	var result []string
	for _, word := range words {
		if word != "" {
			result = append(result, strings.ToLower(word))
		}
	}
	return strings.Join(result, "_")
}

// ToKebabCase converts a string to kebab-case
func ToKebabCase(s string) string {
	words := splitIntoWords(s)
	var result []string
	for _, word := range words {
		if word != "" {
			result = append(result, strings.ToLower(word))
		}
	}
	return strings.Join(result, "-")
}

// splitIntoWords splits a string into words, handling various separators and camelCase
func splitIntoWords(s string) []string {
	if s == "" {
		return []string{}
	}

	// First, replace common separators with spaces
	re := regexp.MustCompile(`[_\-\s]+`)
	s = re.ReplaceAllString(s, " ")

	// Split camelCase and PascalCase
	var words []string
	var currentWord strings.Builder

	runes := []rune(s)
	for i, r := range runes {
		if unicode.IsSpace(r) {
			if currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
			continue
		}

		// Check if this is the start of a new word in camelCase/PascalCase
		if i > 0 && unicode.IsUpper(r) && unicode.IsLower(runes[i-1]) {
			if currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
		}

		// Check for sequences of uppercase letters followed by lowercase
		if i > 0 && i < len(runes)-1 && unicode.IsUpper(r) && unicode.IsUpper(runes[i-1]) && unicode.IsLower(runes[i+1]) {
			if currentWord.Len() > 1 {
				// Keep the last character for the next word
				word := currentWord.String()
				words = append(words, word[:len(word)-1])
				currentWord.Reset()
				currentWord.WriteRune(runes[i-1])
			}
		}

		currentWord.WriteRune(r)
	}

	if currentWord.Len() > 0 {
		words = append(words, currentWord.String())
	}

	// Filter out empty words
	var filteredWords []string
	for _, word := range words {
		word = strings.TrimSpace(word)
		if word != "" {
			filteredWords = append(filteredWords, word)
		}
	}

	return filteredWords
}

// makePlural provides simple pluralization
func makePlural(word string) string {
	if word == "" {
		return ""
	}

	// Simple English pluralization rules
	word = strings.ToLower(word)

	// Special cases
	specialCases := map[string]string{
		"child":  "children",
		"foot":   "feet",
		"tooth":  "teeth",
		"goose":  "geese",
		"man":    "men",
		"woman":  "women",
		"mouse":  "mice",
		"person": "people",
	}

	if plural, exists := specialCases[word]; exists {
		return plural
	}

	// Words ending in s, ss, sh, ch, x, z
	if strings.HasSuffix(word, "s") || strings.HasSuffix(word, "ss") ||
		strings.HasSuffix(word, "sh") || strings.HasSuffix(word, "ch") ||
		strings.HasSuffix(word, "x") || strings.HasSuffix(word, "z") {
		return word + "es"
	}

	// Words ending in consonant + y
	if len(word) > 1 && strings.HasSuffix(word, "y") {
		penultimate := word[len(word)-2]
		if !isVowel(penultimate) {
			return word[:len(word)-1] + "ies"
		}
	}

	// Words ending in consonant + o
	if len(word) > 1 && strings.HasSuffix(word, "o") {
		penultimate := word[len(word)-2]
		if !isVowel(penultimate) {
			return word + "es"
		}
	}

	// Words ending in f or fe
	if strings.HasSuffix(word, "f") {
		return word[:len(word)-1] + "ves"
	}
	if strings.HasSuffix(word, "fe") {
		return word[:len(word)-2] + "ves"
	}

	// Default: just add s
	return word + "s"
}

// isVowel checks if a character is a vowel
func isVowel(c byte) bool {
	vowels := "aeiou"
	return strings.ContainsRune(vowels, rune(c))
}

// CreateResourceNames generates all naming conventions for a resource
func CreateResourceNames(name string) *NamingConventions {
	singular := strings.ToLower(strings.TrimSpace(name))

	// Simple pluralization (can be enhanced with proper pluralization library)
	plural := makePlural(singular)

	return &NamingConventions{
		Singular:       singular,
		Plural:         plural,
		PascalCase:     ToPascalCase(singular),
		PascalPlural:   ToPascalCase(plural),
		CamelCase:      ToCamelCase(singular),
		CamelPlural:    ToCamelCase(plural),
		SnakeCase:      ToSnakeCase(singular),
		SnakePlural:    ToSnakeCase(plural),
		KebabCase:      ToKebabCase(singular),
		KebabPlural:    ToKebabCase(plural),
		TableName:      ToSnakeCase(plural),
		CollectionName: ToSnakeCase(plural),
	}
}

// FieldNames holds various naming formats for fields (compatible with templates)
type FieldNames struct {
	PascalCase string `json:"pascal_case"` // UserName
	CamelCase  string `json:"camel_case"`  // userName
	SnakeCase  string `json:"snake_case"`  // user_name
	KebabCase  string `json:"kebab_case"`  // user-name
}

// CreateFieldNames generates all naming conventions for a field
func CreateFieldNames(name string) FieldNames {
	return FieldNames{
		PascalCase: ToPascalCase(name),
		CamelCase:  ToCamelCase(name),
		SnakeCase:  ToSnakeCase(name),
		KebabCase:  ToKebabCase(name),
	}
}
