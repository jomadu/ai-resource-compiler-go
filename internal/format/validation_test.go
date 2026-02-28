package format

import "testing"

func TestValidateID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		// Valid cases
		{name: "alphanumeric", id: "cleanCode123", wantErr: false},
		{name: "with hyphens", id: "clean-code", wantErr: false},
		{name: "with underscores", id: "clean_code", wantErr: false},
		{name: "mixed valid chars", id: "clean-Code_123", wantErr: false},
		{name: "uppercase", id: "CLEANCODE", wantErr: false},
		{name: "lowercase", id: "cleancode", wantErr: false},
		{name: "numbers only", id: "123456", wantErr: false},
		
		// Invalid cases
		{name: "empty", id: "", wantErr: true},
		{name: "forward slash", id: "clean/code", wantErr: true},
		{name: "backslash", id: "clean\\code", wantErr: true},
		{name: "colon", id: "clean:code", wantErr: true},
		{name: "asterisk", id: "clean*code", wantErr: true},
		{name: "question mark", id: "clean?code", wantErr: true},
		{name: "double quote", id: "clean\"code", wantErr: true},
		{name: "less than", id: "clean<code", wantErr: true},
		{name: "greater than", id: "clean>code", wantErr: true},
		{name: "pipe", id: "clean|code", wantErr: true},
		{name: "space", id: "clean code", wantErr: true},
		{name: "dot", id: "clean.code", wantErr: true},
		{name: "parentheses", id: "clean(code)", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateID(%q) error = %v, wantErr %v", tt.id, err, tt.wantErr)
			}
		})
	}
}

func TestValidateRuleName(t *testing.T) {
	tests := []struct {
		name    string
		ruleName string
		wantErr bool
	}{
		// Valid cases
		{name: "simple name", ruleName: "Use Meaningful Names", wantErr: false},
		{name: "with hyphens", ruleName: "Clean-Code-Rules", wantErr: false},
		{name: "with underscores", ruleName: "Clean_Code_Rules", wantErr: false},
		{name: "with numbers", ruleName: "Rule 123", wantErr: false},
		{name: "with special chars", ruleName: "Rule: Clean Code!", wantErr: false},
		{name: "empty", ruleName: "", wantErr: false},
		
		// Invalid cases
		{name: "with open paren", ruleName: "Rule (MUST)", wantErr: true},
		{name: "with close paren", ruleName: "Rule MUST)", wantErr: true},
		{name: "with both parens", ruleName: "Rule (MUST) Apply", wantErr: true},
		{name: "only parens", ruleName: "()", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRuleName(tt.ruleName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRuleName(%q) error = %v, wantErr %v", tt.ruleName, err, tt.wantErr)
			}
		})
	}
}
