package targets

import (
	"strings"
	"testing"

	"github.com/jomadu/ai-resource-compiler-go/internal/format"
	"github.com/jomadu/ai-resource-compiler-go/pkg/compiler"
)

func TestClaudeCompiler_Name(t *testing.T) {
	c := &ClaudeCompiler{}
	if got := c.Name(); got != "claude" {
		t.Errorf("Name() = %v, want claude", got)
	}
}

func TestClaudeCompiler_SupportedVersions(t *testing.T) {
	c := &ClaudeCompiler{}
	versions := c.SupportedVersions()
	if len(versions) != 1 || versions[0] != "ai-resource/draft" {
		t.Errorf("SupportedVersions() = %v, want [ai-resource/draft]", versions)
	}
}

func TestClaudeCompiler_CompileRuleWithScope(t *testing.T) {
	c := &ClaudeCompiler{}
	resource := &compiler.Resource{
		APIVersion: "ai-resource/draft",
		Kind:       "Rule",
		Spec: &format.Rule{
			Metadata: format.Metadata{
				ID:          "testRule",
				Name:        "Test Rule",
				Description: "A test rule",
			},
			Spec: format.RuleSpec{
				Enforcement: "must",
				Scope: []format.ScopeEntry{
					{Files: []string{"**/*.ts", "**/*.js"}},
				},
				Body: format.Body{String: strPtr("Rule body content")},
			},
		},
	}

	results, err := c.Compile(resource)
	if err != nil {
		t.Fatalf("Compile() error = %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Compile() returned %d results, want 1", len(results))
	}

	result := results[0]
	if result.Path != "testRule.md" {
		t.Errorf("Path = %v, want testRule.md", result.Path)
	}

	if !strings.Contains(result.Content, "paths:") {
		t.Error("Content missing paths frontmatter")
	}
	if !strings.Contains(result.Content, "**/*.ts") {
		t.Error("Content missing scope files in frontmatter")
	}
	if !strings.Contains(result.Content, "# Test Rule (MUST)") {
		t.Error("Content missing enforcement header")
	}
	if !strings.Contains(result.Content, "Rule body content") {
		t.Error("Content missing body")
	}
}

func TestClaudeCompiler_CompileRuleWithoutScope(t *testing.T) {
	c := &ClaudeCompiler{}
	resource := &compiler.Resource{
		APIVersion: "ai-resource/draft",
		Kind:       "Rule",
		Spec: &format.Rule{
			Metadata: format.Metadata{
				ID:          "testRule",
				Name:        "Test Rule",
				Description: "A test rule",
			},
			Spec: format.RuleSpec{
				Enforcement: "must",
				Body:        format.Body{String: strPtr("Rule body content")},
			},
		},
	}

	results, err := c.Compile(resource)
	if err != nil {
		t.Fatalf("Compile() error = %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Compile() returned %d results, want 1", len(results))
	}

	result := results[0]
	
	// Should not have paths frontmatter
	lines := strings.Split(result.Content, "\n")
	if len(lines) > 0 && lines[0] == "---" {
		// Check if it's paths frontmatter (not metadata block)
		if strings.Contains(result.Content, "paths:") && strings.Index(result.Content, "paths:") < strings.Index(result.Content, "rule:") {
			t.Error("Content should not have paths frontmatter when no scope")
		}
	}
}

func TestClaudeCompiler_CompileRuleset(t *testing.T) {
	c := &ClaudeCompiler{}
	resource := &compiler.Resource{
		APIVersion: "ai-resource/draft",
		Kind:       "Ruleset",
		Spec: &format.Ruleset{
			Metadata: format.Metadata{
				ID:          "testRuleset",
				Name:        "Test Ruleset",
				Description: "A test ruleset",
			},
			Spec: struct {
				Rules     map[string]format.RuleItem
				Fragments map[string]string
			}{
				Rules: map[string]format.RuleItem{
					"rule1": {
						Name:        "Rule One",
						Enforcement: "should",
						Body:        format.Body{String: strPtr("First rule")},
					},
					"rule2": {
						Name:        "Rule Two",
						Enforcement: "must",
						Body:        format.Body{String: strPtr("Second rule")},
					},
				},
			},
		},
	}

	results, err := c.Compile(resource)
	if err != nil {
		t.Fatalf("Compile() error = %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Compile() returned %d results, want 2", len(results))
	}

	paths := []string{results[0].Path, results[1].Path}
	if !contains(paths, "testRuleset_rule1.md") {
		t.Error("Missing testRuleset_rule1.md")
	}
	if !contains(paths, "testRuleset_rule2.md") {
		t.Error("Missing testRuleset_rule2.md")
	}
}

func TestClaudeCompiler_CompilePrompt(t *testing.T) {
	c := &ClaudeCompiler{}
	resource := &compiler.Resource{
		APIVersion: "ai-resource/draft",
		Kind:       "Prompt",
		Spec: &format.Prompt{
			Metadata: format.Metadata{
				ID:          "testPrompt",
				Name:        "Test Prompt",
				Description: "A test prompt",
			},
			Spec: format.PromptSpec{
				Body: format.Body{String: strPtr("Prompt body content")},
			},
		},
	}

	results, err := c.Compile(resource)
	if err != nil {
		t.Fatalf("Compile() error = %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Compile() returned %d results, want 1", len(results))
	}

	result := results[0]
	if result.Path != "testPrompt/SKILL.md" {
		t.Errorf("Path = %v, want testPrompt/SKILL.md", result.Path)
	}

	if result.Content != "Prompt body content" {
		t.Errorf("Content = %v, want 'Prompt body content'", result.Content)
	}
}

func TestClaudeCompiler_CompilePromptset(t *testing.T) {
	c := &ClaudeCompiler{}
	resource := &compiler.Resource{
		APIVersion: "ai-resource/draft",
		Kind:       "Promptset",
		Spec: &format.Promptset{
			Metadata: format.Metadata{
				ID:          "testPromptset",
				Name:        "Test Promptset",
				Description: "A test promptset",
			},
			Spec: struct {
				Prompts   map[string]format.PromptItem
				Fragments map[string]string
			}{
				Prompts: map[string]format.PromptItem{
					"prompt1": {
						Body: format.Body{String: strPtr("First prompt")},
					},
					"prompt2": {
						Body: format.Body{String: strPtr("Second prompt")},
					},
				},
			},
		},
	}

	results, err := c.Compile(resource)
	if err != nil {
		t.Fatalf("Compile() error = %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Compile() returned %d results, want 2", len(results))
	}

	paths := []string{results[0].Path, results[1].Path}
	if !contains(paths, "testPromptset_prompt1/SKILL.md") {
		t.Error("Missing testPromptset_prompt1/SKILL.md")
	}
	if !contains(paths, "testPromptset_prompt2/SKILL.md") {
		t.Error("Missing testPromptset_prompt2/SKILL.md")
	}
}
