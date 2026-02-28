package format

import (
	"fmt"
	"strings"
)

// Placeholder types until ai-resource-core-go is implemented
type Metadata struct {
	ID          string
	Name        string
	Description string
}

type Body struct {
	String *string
	Array  []string
}

type ScopeEntry struct {
	Files []string
}

type RuleItem struct {
	Name        string
	Description string
	Enforcement string
	Scope       []ScopeEntry
	Body        Body
}

type RuleSpec struct {
	Enforcement string
	Scope       []ScopeEntry
	Body        Body
	Fragments   map[string]string
}

type Ruleset struct {
	Metadata Metadata
	Spec     struct {
		Rules     map[string]RuleItem
		Fragments map[string]string
	}
}

type Rule struct {
	Metadata Metadata
	Spec     RuleSpec
}

type PromptItem struct {
	Name string
	Body Body
}

type PromptSpec struct {
	Body      Body
	Fragments map[string]string
}

type Promptset struct {
	Metadata Metadata
	Spec     struct {
		Prompts   map[string]PromptItem
		Fragments map[string]string
	}
}

type Prompt struct {
	Metadata Metadata
	Spec     PromptSpec
}

// GenerateRuleMetadataBlockFromRuleset generates complete rule content from a ruleset.
// Returns: metadata block + enforcement header + resolved body
func GenerateRuleMetadataBlockFromRuleset(ruleset *Ruleset, ruleID string) string {
	ruleSpec := ruleset.Spec.Rules[ruleID]
	body := resolveBody(ruleSpec.Body, ruleset.Spec.Fragments)

	var sb strings.Builder
	sb.WriteString("---\n")
	sb.WriteString("ruleset:\n")
	sb.WriteString(fmt.Sprintf("  id: %s\n", ruleset.Metadata.ID))
	if ruleset.Metadata.Name != "" {
		sb.WriteString(fmt.Sprintf("  name: %s\n", ruleset.Metadata.Name))
	}
	if ruleset.Metadata.Description != "" {
		sb.WriteString(fmt.Sprintf("  description: %s\n", ruleset.Metadata.Description))
	}
	sb.WriteString("  rules:\n")
	for id := range ruleset.Spec.Rules {
		sb.WriteString(fmt.Sprintf("    - %s\n", id))
	}
	sb.WriteString("rule:\n")
	sb.WriteString(fmt.Sprintf("  id: %s\n", ruleID))
	if ruleSpec.Name != "" {
		sb.WriteString(fmt.Sprintf("  name: %s\n", ruleSpec.Name))
	}
	if ruleSpec.Description != "" {
		sb.WriteString(fmt.Sprintf("  description: %s\n", ruleSpec.Description))
	}
	sb.WriteString(fmt.Sprintf("  enforcement: %s\n", ruleSpec.Enforcement))
	if len(ruleSpec.Scope) > 0 {
		sb.WriteString("  scope:\n")
		sb.WriteString("    files:\n")
		for _, entry := range ruleSpec.Scope {
			for _, file := range entry.Files {
				sb.WriteString(fmt.Sprintf("      - \"%s\"\n", file))
			}
		}
	}
	sb.WriteString("---\n\n")

	header := generateEnforcementHeader(ruleSpec.Name, ruleSpec.Enforcement)
	sb.WriteString(header)
	sb.WriteString("\n\n")
	sb.WriteString(body)

	return sb.String()
}

// GenerateRuleMetadataBlockFromRule generates complete rule content from a standalone rule.
// Returns: metadata block + enforcement header + resolved body
func GenerateRuleMetadataBlockFromRule(rule *Rule) string {
	body := resolveBody(rule.Spec.Body, rule.Spec.Fragments)

	var sb strings.Builder
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("id: %s\n", rule.Metadata.ID))
	if rule.Metadata.Name != "" {
		sb.WriteString(fmt.Sprintf("name: %s\n", rule.Metadata.Name))
	}
	if rule.Metadata.Description != "" {
		sb.WriteString(fmt.Sprintf("description: %s\n", rule.Metadata.Description))
	}
	sb.WriteString(fmt.Sprintf("enforcement: %s\n", rule.Spec.Enforcement))
	if len(rule.Spec.Scope) > 0 {
		sb.WriteString("scope:\n")
		sb.WriteString("  files:\n")
		for _, entry := range rule.Spec.Scope {
			for _, file := range entry.Files {
				sb.WriteString(fmt.Sprintf("    - \"%s\"\n", file))
			}
		}
	}
	sb.WriteString("---\n\n")

	header := generateEnforcementHeader(rule.Metadata.Name, rule.Spec.Enforcement)
	sb.WriteString(header)
	sb.WriteString("\n\n")
	sb.WriteString(body)

	return sb.String()
}

func generateEnforcementHeader(name, enforcement string) string {
	return fmt.Sprintf("# %s (%s)", name, strings.ToUpper(enforcement))
}

// ResolveBody resolves body content with fragment substitution.
func ResolveBody(body Body, fragments map[string]string) string {
	return resolveBody(body, fragments)
}

func resolveBody(body Body, fragments map[string]string) string {
	if body.String != nil {
		return *body.String
	}
	if len(body.Array) > 0 {
		var parts []string
		for _, ref := range body.Array {
			if strings.HasPrefix(ref, "$") {
				fragmentKey := strings.TrimPrefix(ref, "$")
				if fragment, ok := fragments[fragmentKey]; ok {
					parts = append(parts, fragment)
				}
			} else {
				parts = append(parts, ref)
			}
		}
		return strings.Join(parts, "\n\n")
	}
	return ""
}
