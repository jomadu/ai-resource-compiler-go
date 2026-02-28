package targets

import (
	"fmt"
	"strings"

	"github.com/jomadu/ai-resource-compiler-go/internal/format"
	"github.com/jomadu/ai-resource-compiler-go/pkg/compiler"
	"gopkg.in/yaml.v3"
)

type ClaudeCompiler struct{}

func (c *ClaudeCompiler) Name() string {
	return "claude"
}

func (c *ClaudeCompiler) SupportedVersions() []string {
	return []string{"ai-resource/draft"}
}

func (c *ClaudeCompiler) Compile(resource *compiler.Resource) ([]compiler.CompilationResult, error) {
	if resource.APIVersion != "ai-resource/draft" {
		return nil, fmt.Errorf("unsupported apiVersion: %s for claude", resource.APIVersion)
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

func (c *ClaudeCompiler) compileRule(resource *compiler.Resource) ([]compiler.CompilationResult, error) {
	rule := resource.Spec.(*format.Rule)

	if err := format.ValidateID(rule.Metadata.ID); err != nil {
		return nil, err
	}
	if err := format.ValidateRuleName(rule.Metadata.Name); err != nil {
		return nil, err
	}

	path := format.BuildStandalonePath(rule.Metadata.ID, ".md")
	metadataBlock := format.GenerateRuleMetadataBlockFromRule(rule)
	
	var content strings.Builder
	if len(rule.Spec.Scope) > 0 {
		content.WriteString(generatePathsFrontmatter(rule.Spec.Scope))
		content.WriteString("\n\n")
	}
	content.WriteString(metadataBlock)

	return []compiler.CompilationResult{{Path: path, Content: content.String()}}, nil
}

func (c *ClaudeCompiler) compileRuleset(resource *compiler.Resource) ([]compiler.CompilationResult, error) {
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

		path := format.BuildCollectionPath(ruleset.Metadata.ID, ruleID, ".md")
		metadataBlock := format.GenerateRuleMetadataBlockFromRuleset(ruleset, ruleID)
		
		var content strings.Builder
		if len(ruleSpec.Scope) > 0 {
			content.WriteString(generatePathsFrontmatter(ruleSpec.Scope))
			content.WriteString("\n\n")
		}
		content.WriteString(metadataBlock)

		results = append(results, compiler.CompilationResult{Path: path, Content: content.String()})
	}

	return results, nil
}

func (c *ClaudeCompiler) compilePrompt(resource *compiler.Resource) ([]compiler.CompilationResult, error) {
	prompt := resource.Spec.(*format.Prompt)

	if err := format.ValidateID(prompt.Metadata.ID); err != nil {
		return nil, err
	}

	path := format.BuildClaudeStandalonePath(prompt.Metadata.ID)
	content := format.ResolveBody(prompt.Spec.Body, prompt.Spec.Fragments)

	return []compiler.CompilationResult{{Path: path, Content: content}}, nil
}

func (c *ClaudeCompiler) compilePromptset(resource *compiler.Resource) ([]compiler.CompilationResult, error) {
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
		path := format.BuildClaudeCollectionPath(promptset.Metadata.ID, promptID)
		content := format.ResolveBody(promptSpec.Body, promptset.Spec.Fragments)

		results = append(results, compiler.CompilationResult{Path: path, Content: content})
	}

	return results, nil
}

func generatePathsFrontmatter(scope []format.ScopeEntry) string {
	var files []string
	for _, entry := range scope {
		files = append(files, entry.Files...)
	}

	frontmatter := map[string]interface{}{
		"paths": files,
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
