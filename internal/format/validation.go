package format

import (
	"fmt"
	"strings"
)

// ValidateID checks if an ID contains only allowed characters.
// Allowed: a-z, A-Z, 0-9, -, _
func ValidateID(id string) error {
	if id == "" {
		return fmt.Errorf("ID cannot be empty")
	}

	for _, char := range id {
		if !isValidIDChar(char) {
			return fmt.Errorf("ID contains invalid character '%c' in '%s'", char, id)
		}
	}

	return nil
}

// ValidateRuleName checks if a rule name contains parentheses.
func ValidateRuleName(name string) error {
	if strings.ContainsAny(name, "()") {
		return fmt.Errorf("rule name cannot contain parentheses: '%s'", name)
	}
	return nil
}

func isValidIDChar(char rune) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		(char >= '0' && char <= '9') ||
		char == '-' ||
		char == '_'
}
