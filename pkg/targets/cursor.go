package targets

import (
	"fmt"
	"strings"

	"github.com/jomadu/ai-resource-compiler-go/internal/format"
	"github.com/jomadu/ai-resource-compiler-go/pkg/compiler"
	"gopkg.in/yaml.v3"
)

type CursorCompiler struct{}

func init() {
	compiler.RegisterDefaultTarget(compiler.TargetCursor, &CursorCompiler{})
}

func (c *CursorCompiler) Name() string {
	return "cursor"
}

func (c *CursorCompiler) SupportedVersions() []string {
	return []string{"ai-resource/draft"}
}

func (c *CursorCompiler) Compile(resource *compiler.Resource) ([]compiler.CompilationResult, error) {
	if resource.APIVersion != "ai-resource/draft" {
		return nil, fmt.Errorf("unsupported apiVersion: %s for cursor", resource.APIVersion)
	}

	switch resource.Kind {
	case "Rule":
		return c.compileRule(resource)
	case "Ruleset":
		return c.compileRuleset(resource)
	case "Prompt":
		return c.compilePrompt(resource)
	case "Promptset":
		return c.compilePromptset(resource)
	default:
		return nil, fmt.Errorf("unsupported kind: %s", resource.Kind)
	}
}

func (c *CursorCompiler) compileRule(resource *compiler.Resource) ([]compiler.CompilationResult, error) {
	rule := resource.Spec.(*format.Rule)

	if err := format.ValidateID(rule.Metadata.ID); err != nil {
		return nil, err
	}
	if err := format.ValidateRuleName(rule.Metadata.Name); err != nil {
		return nil, err
	}

	scopeFiles := extractScopeFiles(rule.Spec.Scope)
	frontmatter := generateMDCFrontmatter(rule.Metadata.Description, rule.Metadata.Name, scopeFiles, rule.Spec.Enforcement)
	path := format.BuildStandalonePath(rule.Metadata.ID, ".mdc")
	metadataBlock := format.GenerateRuleMetadataBlockFromRule(rule)
	content := frontmatter + "\n" + metadataBlock

	return []compiler.CompilationResult{{Path: path, Content: content}}, nil
}

func (c *CursorCompiler) compileRuleset(resource *compiler.Resource) ([]compiler.CompilationResult, error) {
	ruleset := resource.Spec.(*format.Ruleset)

	if err := format.ValidateID(ruleset.Metadata.ID); err != nil {
		return nil, err
	}

	var results []compiler.CompilationResult
	for ruleID := range ruleset.Spec.Rules {
		if err := format.ValidateID(ruleID); err != nil {
			return nil, err
		}
		ruleSpec := ruleset.Spec.Rules[ruleID]
		if err := format.ValidateRuleName(ruleSpec.Name); err != nil {
			return nil, err
		}

		scopeFiles := extractScopeFiles(ruleSpec.Scope)
		frontmatter := generateMDCFrontmatter(ruleSpec.Description, ruleSpec.Name, scopeFiles, ruleSpec.Enforcement)
		path := format.BuildCollectionPath(ruleset.Metadata.ID, ruleID, ".mdc")
		metadataBlock := format.GenerateRuleMetadataBlockFromRuleset(ruleset, ruleID)
		content := frontmatter + "\n" + metadataBlock

		results = append(results, compiler.CompilationResult{Path: path, Content: content})
	}

	return results, nil
}

func (c *CursorCompiler) compilePrompt(resource *compiler.Resource) ([]compiler.CompilationResult, error) {
	prompt := resource.Spec.(*format.Prompt)

	if err := format.ValidateID(prompt.Metadata.ID); err != nil {
		return nil, err
	}

	path := format.BuildStandalonePath(prompt.Metadata.ID, ".md")
	content := format.ResolveBody(prompt.Spec.Body, prompt.Spec.Fragments)

	return []compiler.CompilationResult{{Path: path, Content: content}}, nil
}

func (c *CursorCompiler) compilePromptset(resource *compiler.Resource) ([]compiler.CompilationResult, error) {
	promptset := resource.Spec.(*format.Promptset)

	if err := format.ValidateID(promptset.Metadata.ID); err != nil {
		return nil, err
	}

	var results []compiler.CompilationResult
	for promptID := range promptset.Spec.Prompts {
		if err := format.ValidateID(promptID); err != nil {
			return nil, err
		}

		promptSpec := promptset.Spec.Prompts[promptID]
		path := format.BuildCollectionPath(promptset.Metadata.ID, promptID, ".md")
		content := format.ResolveBody(promptSpec.Body, promptset.Spec.Fragments)

		results = append(results, compiler.CompilationResult{Path: path, Content: content})
	}

	return results, nil
}

func extractScopeFiles(scope []format.ScopeEntry) []string {
	var files []string
	for _, entry := range scope {
		files = append(files, entry.Files...)
	}
	return files
}

func generateMDCFrontmatter(description, name string, globs []string, enforcement string) string {
	desc := description
	if desc == "" {
		desc = name
	}

	alwaysApply := enforcement == "must"

	frontmatter := map[string]interface{}{
		"description": desc,
		"globs":       globs,
		"alwaysApply": alwaysApply,
	}

	var b strings.Builder
	b.WriteString("---\n")
	encoder := yaml.NewEncoder(&b)
	encoder.SetIndent(2)
	encoder.Encode(frontmatter)
	encoder.Close()
	b.WriteString("---")

	return b.String()
}
