package targets

import (
	"fmt"

	"github.com/jomadu/ai-resource-compiler-go/internal/format"
	"github.com/jomadu/ai-resource-compiler-go/pkg/compiler"
)

type MarkdownCompiler struct{}

func (m *MarkdownCompiler) Name() string {
	return "markdown"
}

func (m *MarkdownCompiler) SupportedVersions() []string {
	return []string{"ai-resource/draft"}
}

func (m *MarkdownCompiler) Compile(resource *compiler.Resource) ([]compiler.CompilationResult, error) {
	if resource.APIVersion != "ai-resource/draft" {
		return nil, fmt.Errorf("unsupported apiVersion: %s for markdown", resource.APIVersion)
	}

	switch resource.Kind {
	case "Rule":
		return m.compileRule(resource)
	case "Ruleset":
		return m.compileRuleset(resource)
	case "Prompt":
		return m.compilePrompt(resource)
	case "Promptset":
		return m.compilePromptset(resource)
	default:
		return nil, fmt.Errorf("unsupported kind: %s", resource.Kind)
	}
}

func (m *MarkdownCompiler) compileRule(resource *compiler.Resource) ([]compiler.CompilationResult, error) {
	rule := resource.Spec.(*format.Rule)
	
	if err := format.ValidateID(rule.Metadata.ID); err != nil {
		return nil, err
	}
	if err := format.ValidateRuleName(rule.Metadata.Name); err != nil {
		return nil, err
	}

	path := format.BuildStandalonePath(rule.Metadata.ID, ".md")
	content := format.GenerateRuleMetadataBlockFromRule(rule)

	return []compiler.CompilationResult{{Path: path, Content: content}}, nil
}

func (m *MarkdownCompiler) compileRuleset(resource *compiler.Resource) ([]compiler.CompilationResult, error) {
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
		content := format.GenerateRuleMetadataBlockFromRuleset(ruleset, ruleID)

		results = append(results, compiler.CompilationResult{Path: path, Content: content})
	}

	return results, nil
}

func (m *MarkdownCompiler) compilePrompt(resource *compiler.Resource) ([]compiler.CompilationResult, error) {
	prompt := resource.Spec.(*format.Prompt)
	
	if err := format.ValidateID(prompt.Metadata.ID); err != nil {
		return nil, err
	}

	path := format.BuildStandalonePath(prompt.Metadata.ID, ".md")
	content := format.ResolveBody(prompt.Spec.Body, prompt.Spec.Fragments)

	return []compiler.CompilationResult{{Path: path, Content: content}}, nil
}

func (m *MarkdownCompiler) compilePromptset(resource *compiler.Resource) ([]compiler.CompilationResult, error) {
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
