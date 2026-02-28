package format

import (
	"strings"
	"testing"
)

func TestGenerateRuleMetadataBlockFromRuleset(t *testing.T) {
	tests := []struct {
		name    string
		ruleset *Ruleset
		ruleID  string
		want    []string // Expected substrings in output
	}{
		{
			name: "full metadata",
			ruleset: &Ruleset{
				Metadata: Metadata{
					ID:          "cleanCode",
					Name:        "Clean Code",
					Description: "Clean code practices",
				},
				Spec: struct {
					Rules     map[string]RuleItem
					Fragments map[string]string
				}{
					Rules: map[string]RuleItem{
						"meaningfulNames": {
							Name:        "Use Meaningful Names",
							Description: "Variables should have descriptive names",
							Enforcement: "must",
							Scope: []ScopeEntry{
								{Files: []string{"**/*.ts", "**/*.js"}},
							},
							Body: Body{String: strPtr("Use descriptive variable names.")},
						},
						"smallFunctions": {
							Name:        "Keep Functions Small",
							Enforcement: "should",
							Body:        Body{String: strPtr("Functions should be small.")},
						},
					},
					Fragments: map[string]string{},
				},
			},
			ruleID: "meaningfulNames",
			want: []string{
				"---",
				"ruleset:",
				"  id: cleanCode",
				"  name: Clean Code",
				"  description: Clean code practices",
				"  rules:",
				"    - meaningfulNames",
				"    - smallFunctions",
				"rule:",
				"  id: meaningfulNames",
				"  name: Use Meaningful Names",
				"  description: Variables should have descriptive names",
				"  enforcement: must",
				"  scope:",
				"    files:",
				"      - \"**/*.ts\"",
				"      - \"**/*.js\"",
				"# Use Meaningful Names (MUST)",
				"Use descriptive variable names.",
			},
		},
		{
			name: "minimal metadata",
			ruleset: &Ruleset{
				Metadata: Metadata{
					ID: "simple",
				},
				Spec: struct {
					Rules     map[string]RuleItem
					Fragments map[string]string
				}{
					Rules: map[string]RuleItem{
						"rule1": {
							Name:        "Rule One",
							Enforcement: "may",
							Body:        Body{String: strPtr("Rule body.")},
						},
					},
					Fragments: map[string]string{},
				},
			},
			ruleID: "rule1",
			want: []string{
				"---",
				"ruleset:",
				"  id: simple",
				"  rules:",
				"    - rule1",
				"rule:",
				"  id: rule1",
				"  name: Rule One",
				"  enforcement: may",
				"# Rule One (MAY)",
				"Rule body.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateRuleMetadataBlockFromRuleset(tt.ruleset, tt.ruleID)
			for _, substr := range tt.want {
				if !strings.Contains(got, substr) {
					t.Errorf("GenerateRuleMetadataBlockFromRuleset() missing substring:\n%q\nGot:\n%s", substr, got)
				}
			}
		})
	}
}

func TestGenerateRuleMetadataBlockFromRule(t *testing.T) {
	tests := []struct {
		name string
		rule *Rule
		want []string // Expected substrings in output
	}{
		{
			name: "full metadata",
			rule: &Rule{
				Metadata: Metadata{
					ID:          "standaloneRule",
					Name:        "Standalone Rule",
					Description: "A standalone rule",
				},
				Spec: RuleSpec{
					Enforcement: "should",
					Scope: []ScopeEntry{
						{Files: []string{"**/*.go"}},
					},
					Body:      Body{String: strPtr("Follow this rule.")},
					Fragments: map[string]string{},
				},
			},
			want: []string{
				"---",
				"id: standaloneRule",
				"name: Standalone Rule",
				"description: A standalone rule",
				"enforcement: should",
				"scope:",
				"  files:",
				"    - \"**/*.go\"",
				"# Standalone Rule (SHOULD)",
				"Follow this rule.",
			},
		},
		{
			name: "minimal metadata",
			rule: &Rule{
				Metadata: Metadata{
					ID:   "minimal",
					Name: "Minimal",
				},
				Spec: RuleSpec{
					Enforcement: "must",
					Body:        Body{String: strPtr("Minimal rule.")},
					Fragments:   map[string]string{},
				},
			},
			want: []string{
				"---",
				"id: minimal",
				"name: Minimal",
				"enforcement: must",
				"# Minimal (MUST)",
				"Minimal rule.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateRuleMetadataBlockFromRule(tt.rule)
			for _, substr := range tt.want {
				if !strings.Contains(got, substr) {
					t.Errorf("GenerateRuleMetadataBlockFromRule() missing substring:\n%q\nGot:\n%s", substr, got)
				}
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}
