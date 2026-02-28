package compiler

import (
	"strings"
	"testing"

	"github.com/jomadu/ai-resource-compiler-go/internal/format"
)

func strPtr(s string) *string {
	return &s
}

// setupCompiler creates a compiler with a mock target registered
func setupCompiler() *Compiler {
	c := &Compiler{
		targets: make(map[Target]TargetCompiler),
	}
	c.RegisterTarget(TargetMarkdown, &mockMarkdownCompiler{})
	c.RegisterTarget(TargetKiro, &mockMarkdownCompiler{})
	c.RegisterTarget(TargetCursor, &mockCursorCompiler{})
	return c
}

// mockMarkdownCompiler for testing
type mockMarkdownCompiler struct{}

func (m *mockMarkdownCompiler) Name() string {
	return "markdown"
}

func (m *mockMarkdownCompiler) SupportedVersions() []string {
	return []string{"ai-resource/draft"}
}

func (m *mockMarkdownCompiler) Compile(resource *Resource) ([]CompilationResult, error) {
	switch resource.Kind {
	case "Rule":
		rule := resource.Spec.(*format.Rule)
		return []CompilationResult{
			{Path: rule.Metadata.ID + ".md", Content: "mock content"},
		}, nil
	case "Ruleset":
		ruleset := resource.Spec.(*format.Ruleset)
		var results []CompilationResult
		for id := range ruleset.Spec.Rules {
			results = append(results, CompilationResult{
				Path:    ruleset.Metadata.ID + "_" + id + ".md",
				Content: "mock content",
			})
		}
		return results, nil
	}
	return nil, nil
}

// mockCursorCompiler for testing
type mockCursorCompiler struct{}

func (m *mockCursorCompiler) Name() string {
	return "cursor"
}

func (m *mockCursorCompiler) SupportedVersions() []string {
	return []string{"ai-resource/draft"}
}

func (m *mockCursorCompiler) Compile(resource *Resource) ([]CompilationResult, error) {
	rule := resource.Spec.(*format.Rule)
	return []CompilationResult{
		{Path: rule.Metadata.ID + ".mdc", Content: "mock content"},
	}, nil
}

func TestCompiler_SingleTarget(t *testing.T) {
	c := setupCompiler()
	resource := &Resource{
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
				Body:        format.Body{String: strPtr("Rule body")},
			},
		},
	}
	resource.Metadata.ID = "testRule"

	opts := CompileOptions{
		Targets: []Target{TargetMarkdown},
	}

	results, err := c.Compile(resource, opts)
	if err != nil {
		t.Fatalf("Compile() error = %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Compile() returned %d results, want 1", len(results))
	}

	if !strings.HasSuffix(results[0].Path, ".md") {
		t.Errorf("Path = %v, want .md extension", results[0].Path)
	}
}

func TestCompiler_MultiTarget(t *testing.T) {
	c := setupCompiler()
	resource := &Resource{
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
				Body:        format.Body{String: strPtr("Rule body")},
			},
		},
	}
	resource.Metadata.ID = "testRule"

	opts := CompileOptions{
		Targets: []Target{TargetMarkdown, TargetKiro, TargetCursor},
	}

	results, err := c.Compile(resource, opts)
	if err != nil {
		t.Fatalf("Compile() error = %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("Compile() returned %d results, want 3", len(results))
	}

	// Verify different extensions
	extensions := make(map[string]bool)
	for _, r := range results {
		if strings.HasSuffix(r.Path, ".md") {
			extensions[".md"] = true
		} else if strings.HasSuffix(r.Path, ".mdc") {
			extensions[".mdc"] = true
		}
	}

	if !extensions[".md"] {
		t.Error("Expected at least one .md result")
	}
	if !extensions[".mdc"] {
		t.Error("Expected at least one .mdc result")
	}
}

func TestCompiler_Ruleset(t *testing.T) {
	c := setupCompiler()
	resource := &Resource{
		APIVersion: "ai-resource/draft",
		Kind:       "Ruleset",
		Spec: &format.Ruleset{
			Metadata: format.Metadata{
				ID:          "testRuleset",
				Name:        "Test Ruleset",
				Description: "A test ruleset",
			},
		},
	}
	resource.Metadata.ID = "testRuleset"
	
	// Initialize the Spec struct
	resource.Spec.(*format.Ruleset).Spec.Rules = map[string]format.RuleItem{
		"rule1": {
			Name:        "Rule 1",
			Enforcement: "must",
			Body:        format.Body{String: strPtr("Rule 1 body")},
		},
		"rule2": {
			Name:        "Rule 2",
			Enforcement: "should",
			Body:        format.Body{String: strPtr("Rule 2 body")},
		},
	}

	opts := CompileOptions{
		Targets: []Target{TargetMarkdown},
	}

	results, err := c.Compile(resource, opts)
	if err != nil {
		t.Fatalf("Compile() error = %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Compile() returned %d results, want 2", len(results))
	}

	// Verify paths contain ruleset and rule IDs
	for _, r := range results {
		if !strings.Contains(r.Path, "testRuleset") {
			t.Errorf("Path %v does not contain ruleset ID", r.Path)
		}
	}
}

func TestCompiler_MissingAPIVersion(t *testing.T) {
	c := NewCompiler()
	resource := &Resource{
		Kind: "Rule",
	}
	resource.Metadata.ID = "testRule"

	opts := CompileOptions{
		Targets: []Target{TargetMarkdown},
	}

	_, err := c.Compile(resource, opts)
	if err == nil {
		t.Fatal("Compile() expected error for missing apiVersion")
	}
	if !strings.Contains(err.Error(), "apiVersion") {
		t.Errorf("Error = %v, want apiVersion error", err)
	}
}

func TestCompiler_MissingKind(t *testing.T) {
	c := NewCompiler()
	resource := &Resource{
		APIVersion: "ai-resource/draft",
	}
	resource.Metadata.ID = "testRule"

	opts := CompileOptions{
		Targets: []Target{TargetMarkdown},
	}

	_, err := c.Compile(resource, opts)
	if err == nil {
		t.Fatal("Compile() expected error for missing kind")
	}
	if !strings.Contains(err.Error(), "kind") {
		t.Errorf("Error = %v, want kind error", err)
	}
}

func TestCompiler_MissingID(t *testing.T) {
	c := NewCompiler()
	resource := &Resource{
		APIVersion: "ai-resource/draft",
		Kind:       "Rule",
	}

	opts := CompileOptions{
		Targets: []Target{TargetMarkdown},
	}

	_, err := c.Compile(resource, opts)
	if err == nil {
		t.Fatal("Compile() expected error for missing metadata.id")
	}
	if !strings.Contains(err.Error(), "metadata.id") {
		t.Errorf("Error = %v, want metadata.id error", err)
	}
}

func TestCompiler_NoTargets(t *testing.T) {
	c := NewCompiler()
	resource := &Resource{
		APIVersion: "ai-resource/draft",
		Kind:       "Rule",
	}
	resource.Metadata.ID = "testRule"

	opts := CompileOptions{
		Targets: []Target{},
	}

	_, err := c.Compile(resource, opts)
	if err == nil {
		t.Fatal("Compile() expected error for no targets")
	}
	if !strings.Contains(err.Error(), "no targets") {
		t.Errorf("Error = %v, want no targets error", err)
	}
}

func TestCompiler_UnknownTarget(t *testing.T) {
	c := NewCompiler()
	resource := &Resource{
		APIVersion: "ai-resource/draft",
		Kind:       "Rule",
		Spec: &format.Rule{
			Metadata: format.Metadata{
				ID:   "testRule",
				Name: "Test Rule",
			},
			Spec: format.RuleSpec{
				Enforcement: "must",
				Body:        format.Body{String: strPtr("Rule body")},
			},
		},
	}
	resource.Metadata.ID = "testRule"

	opts := CompileOptions{
		Targets: []Target{"unknown"},
	}

	_, err := c.Compile(resource, opts)
	if err == nil {
		t.Fatal("Compile() expected error for unknown target")
	}
	if !strings.Contains(err.Error(), "unknown target") {
		t.Errorf("Error = %v, want unknown target error", err)
	}
}

func TestCompiler_UnsupportedVersion(t *testing.T) {
	c := setupCompiler()
	resource := &Resource{
		APIVersion: "ai-resource/v99",
		Kind:       "Rule",
		Spec: &format.Rule{
			Metadata: format.Metadata{
				ID:   "testRule",
				Name: "Test Rule",
			},
			Spec: format.RuleSpec{
				Enforcement: "must",
				Body:        format.Body{String: strPtr("Rule body")},
			},
		},
	}
	resource.Metadata.ID = "testRule"

	opts := CompileOptions{
		Targets: []Target{TargetMarkdown},
	}

	_, err := c.Compile(resource, opts)
	if err == nil {
		t.Fatal("Compile() expected error for unsupported version")
	}
	if !strings.Contains(err.Error(), "does not support apiVersion") {
		t.Errorf("Error = %v, want unsupported version error", err)
	}
}

func TestCompiler_RegisterTarget(t *testing.T) {
	c := NewCompiler()

	// Create mock compiler
	mock := &mockCompiler{}

	err := c.RegisterTarget("mock", mock)
	if err != nil {
		t.Fatalf("RegisterTarget() error = %v", err)
	}

	// Verify target is registered
	if _, ok := c.targets["mock"]; !ok {
		t.Error("Target not registered")
	}
}

func TestCompiler_RegisterTargetNil(t *testing.T) {
	c := NewCompiler()

	err := c.RegisterTarget("mock", nil)
	if err == nil {
		t.Fatal("RegisterTarget() expected error for nil compiler")
	}
}

// mockCompiler for testing
type mockCompiler struct{}

func (m *mockCompiler) Name() string {
	return "mock"
}

func (m *mockCompiler) SupportedVersions() []string {
	return []string{"ai-resource/draft"}
}

func (m *mockCompiler) Compile(resource *Resource) ([]CompilationResult, error) {
	return []CompilationResult{
		{Path: "mock.txt", Content: "mock content"},
	}, nil
}
