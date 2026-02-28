package targets

import (
	"fmt"
	"strings"

	"github.com/jomadu/ai-resource-compiler-go/internal/format"
	"github.com/jomadu/ai-resource-compiler-go/pkg/compiler"
	"gopkg.in/yaml.v3"
)

type CopilotCompiler struct{}

func init() {
	compiler.RegisterDefaultTarget(compiler.TargetCopilot, &CopilotCompiler{})
}

func (c *CopilotCompiler) Name() string {
	return "copilot"
}

func (c *CopilotCompiler) SupportedVersions() []string {
	return []string{"ai-resource/draft"}
}

func (c *CopilotCompiler) Compile(resource *compiler.Resource) ([]compiler.CompilationResult, error) {
	if resource.APIVersion != "ai-resource/draft" {
		return nil, fmt.Errorf("unsupported apiVersion: %s for copilot", resource.APIVersion)
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

func (c *CopilotCompiler) compileRule(resource *compiler.Resource) ([]compiler.CompilationResult, error) {
	rule := resource.Spec.(*format.Rule)

	if err := format.ValidateID(rule.Metadata.ID); err != nil {
		return nil, err
	}
	if err := format.ValidateRuleName(rule.Metadata.Name); err != nil {
		return nil, err
	}

	scopeFiles := extractScopeFiles(rule.Spec.Scope)
	frontmatter := generateApplyToFrontmatter(scopeFiles)
	path := format.BuildStandalonePath(rule.Metadata.ID, ".instructions.md")
	metadataBlock := format.GenerateRuleMetadataBlockFromRule(rule)
	content := frontmatter + "\n" + metadataBlock

	return []compiler.CompilationResult{{Path: path, Content: content}}, nil
}

func (c *CopilotCompiler) compileRuleset(resource *compiler.Resource) ([]compiler.CompilationResult, error) {
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
		frontmatter := generateApplyToFrontmatter(scopeFiles)
		path := format.BuildCollectionPath(ruleset.Metadata.ID, ruleID, ".instructions.md")
		metadataBlock := format.GenerateRuleMetadataBlockFromRuleset(ruleset, ruleID)
		content := frontmatter + "\n" + metadataBlock

		results = append(results, compiler.CompilationResult{Path: path, Content: content})
	}

	return results, nil
}

func (c *CopilotCompiler) compilePrompt(resource *compiler.Resource) ([]compiler.CompilationResult, error) {
	prompt := resource.Spec.(*format.Prompt)

	if err := format.ValidateID(prompt.Metadata.ID); err != nil {
		return nil, err
	}

	frontmatter := generateApplyToFrontmatter([]string{})
	path := format.BuildStandalonePath(prompt.Metadata.ID, ".prompt.md")
	body := format.ResolveBody(prompt.Spec.Body, prompt.Spec.Fragments)
	content := frontmatter + "\n" + body

	return []compiler.CompilationResult{{Path: path, Content: content}}, nil
}

func (c *CopilotCompiler) compilePromptset(resource *compiler.Resource) ([]compiler.CompilationResult, error) {
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
		frontmatter := generateApplyToFrontmatter([]string{})
		path := format.BuildCollectionPath(promptset.Metadata.ID, promptID, ".prompt.md")
		body := format.ResolveBody(promptSpec.Body, promptset.Spec.Fragments)
		content := frontmatter + "\n" + body

		results = append(results, compiler.CompilationResult{Path: path, Content: content})
	}

	return results, nil
}

func generateApplyToFrontmatter(files []string) string {
	frontmatter := map[string]interface{}{
		"applyTo": files,
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
