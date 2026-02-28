# Metadata Block Embedding

## Job to be Done
Preserve ruleset and rule context in compiled rule output so that AI coding tools understand the origin, purpose, and enforcement level of each rule.

## Activities
1. Embed YAML metadata block at the start of compiled rule files
2. Include ruleset context (id, name, description, rules list)
3. Include rule context (id, name, description, enforcement, scope)
4. Generate enforcement header from rule name and enforcement level
5. Omit metadata blocks from compiled prompts

## Acceptance Criteria
- [ ] Compiled rules include YAML metadata block with ruleset and rule sections
- [ ] Metadata block appears before rule body content
- [ ] Enforcement header follows format "# {Rule Name} ({ENFORCEMENT})"
- [ ] Compiled prompts contain only body content (no metadata)
- [ ] Minimal metadata blocks omit optional fields (description, scope)
- [ ] Metadata block uses standard YAML frontmatter delimiters (---)

## Data Structures

### Metadata Block
```yaml
---
ruleset:
  id: string
  name: string (optional)
  description: string (optional)
  rules: []string
rule:
  id: string
  name: string (optional)
  description: string (optional)
  enforcement: string
  scope: object (optional)
    files: []string
---
```

**For rules in rulesets:**
- `ruleset.id` - Collection metadata ID
- `ruleset.name` - Collection metadata name (optional, from Metadata.Name)
- `ruleset.description` - Collection metadata description (optional, from Metadata.Description)
- `ruleset.rules` - List of rule IDs from collection (map keys)
- `rule.id` - Item ID (map key from Ruleset.Spec.Rules)
- `rule.name` - Item name (optional, from RuleItem.Name)
- `rule.description` - Item description (optional, from RuleItem.Description)
- `rule.enforcement` - Enforcement level (may, should, must)
- `rule.scope.files` - File patterns where rule applies (extracted from []ScopeEntry)

**For standalone rules:**
```yaml
---
id: string
name: string (optional)
description: string (optional)
enforcement: string
scope: object (optional)
  files: []string
---
```

- No nesting - fields at root level
- `id` - Resource metadata ID
- `name` - Resource metadata name (optional, from Metadata.Name)
- `description` - Resource metadata description (optional, from Metadata.Description)
- `enforcement` - Enforcement level (may, should, must)
- `scope.files` - File patterns where rule applies (extracted from []ScopeEntry)

### Enforcement Header
```
# {Rule Name} ({ENFORCEMENT})
```

**Format:**
- Rule name from `rule.name` field
- Enforcement level uppercased (MUST, SHOULD, MAY)
- Follows metadata block, precedes rule body

## Shared Functions

Target compilers use these shared functions to generate consistent metadata blocks for rules.

### GenerateRuleMetadataBlockFromRuleset

```go
func GenerateRuleMetadataBlockFromRuleset(
    ruleset *airesource.Ruleset,
    ruleID string,
) string
```

**Parameters:**
- `ruleset` - Ruleset resource
- `ruleID` - Rule ID (map key from Ruleset.Spec.Rules)

**Returns:**
- Complete compiled rule content: metadata block + enforcement header + resolved body

**Algorithm:**
1. Get ruleSpec from ruleset.Spec.Rules[ruleID]
2. Resolve body using airesource.ResolveBody(ruleSpec.Body, ruleset.Spec.Fragments)
3. Create YAML metadata block with ruleset and rule sections
4. Generate enforcement header: `# {Name} ({ENFORCEMENT})`
5. Concatenate: metadata block + enforcement header + resolved body
6. Return complete content

**Example:**
```go
ruleset := &airesource.Ruleset{
    Metadata: airesource.Metadata{
        ID: "cleanCode",
        Name: "Clean Code",
    },
    Spec: airesource.RulesetSpec{
        Rules: map[string]airesource.RuleItem{
            "meaningfulNames": {
                Name: "Use Meaningful Names",
                Enforcement: airesource.EnforcementMust,
                Body: airesource.Body{String: &bodyStr},
            },
        },
    },
}

content := GenerateRuleMetadataBlockFromRuleset(ruleset, "meaningfulNames")
// Returns:
// ---
// ruleset:
//   id: cleanCode
//   name: Clean Code
//   rules:
//     - meaningfulNames
// rule:
//   id: meaningfulNames
//   name: Use Meaningful Names
//   enforcement: must
// ---
//
// # Use Meaningful Names (MUST)
//
// [resolved body content]
```

### GenerateRuleMetadataBlockFromRule

```go
func GenerateRuleMetadataBlockFromRule(
    rule *airesource.Rule,
) string
```

**Parameters:**
- `rule` - Standalone Rule resource

**Returns:**
- Complete compiled rule content: metadata block + enforcement header + resolved body

**Algorithm:**
1. Resolve body using airesource.ResolveBody(rule.Spec.Body, rule.Spec.Fragments)
2. Create YAML metadata block with fields at root level (no rule wrapper)
3. Add fields from rule.Metadata and rule.Spec (id, name, description, enforcement, scope)
4. Generate enforcement header: `# {Name} ({ENFORCEMENT})`
5. Concatenate: metadata block + enforcement header + resolved body
6. Return complete content

**Example:**
```go
rule := &airesource.Rule{
    Metadata: airesource.Metadata{
        ID: "meaningfulNames",
        Name: "Use Meaningful Names",
    },
    Spec: airesource.RuleSpec{
        Enforcement: airesource.EnforcementMust,
        Body: airesource.Body{String: &bodyStr},
    },
}

content := GenerateRuleMetadataBlockFromRule(rule)
// Returns:
// ---
// id: meaningfulNames
// name: Use Meaningful Names
// enforcement: must
// ---
//
// # Use Meaningful Names (MUST)
//
// [resolved body content]
```

## Algorithm

1. Check resource kind (Rule, Ruleset, Prompt, Promptset)
2. If prompt: resolve body and return (no metadata block)
3. If ruleset: call `GenerateRuleMetadataBlockFromRuleset(ruleset, ruleID)` for each rule
4. If standalone rule: call `GenerateRuleMetadataBlockFromRule(rule)`

**Pseudocode:**
```
function CompileRule(resource, ruleID):
    if resource.Kind == "Ruleset":
        return GenerateRuleMetadataBlockFromRuleset(resource, ruleID)
    else:  // resource.Kind == "Rule"
        return GenerateRuleMetadataBlockFromRule(resource)

function CompilePrompt(resource, promptID):
    if resource.Kind == "Promptset":
        promptset = resource.AsPromptset()
        promptSpec = promptset.Spec.Prompts[promptID]
        return airesource.ResolveBody(promptSpec.Body, promptset.Spec.Fragments)
    else:  // resource.Kind == "Prompt"
        prompt = resource.AsPrompt()
        return airesource.ResolveBody(prompt.Spec.Body, prompt.Spec.Fragments)
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Prompt resource | Return body only, no metadata block |
| Missing optional fields | Omit from metadata block |
| Empty rules list | Include empty array `rules: []` |
| No scope defined | Omit scope section entirely |
| Enforcement level lowercase | Uppercase in header (must → MUST) |

## Dependencies

- Resource model from ai-resource-core-go (Metadata, RuleItem, PromptItem, Enforcement)
- YAML serialization library for metadata block generation

## Implementation Mapping

**Source files:**
- `internal/format/metadata.go` - Implements `GenerateRuleMetadataBlockFromRuleset()` and `GenerateRuleMetadataBlockFromRule()` functions
- `pkg/targets/*.go` - Target compilers call these functions when compiling rules

**Related specs:**
- `compiler-architecture.md` - Defines where metadata embedding fits in compilation pipeline
- All target compiler specs - Each target uses these functions for rule compilation

## Examples

### Example 1: Full Metadata Block

**Input:**
```go
resource := Resource{
    Type: "rule",
    Ruleset: Ruleset{
        ID: "cleanCode",
        Name: "Clean Code",
        Description: "Clean code practices",
        Rules: []string{"meaningfulNames", "smallFunctions"},
    },
    Rule: Rule{
        ID: "meaningfulNames",
        Name: "Use Meaningful Names",
        Description: "Variables and functions should have descriptive names",
        Enforcement: "must",
        Scope: Scope{
            Files: []string{"**/*.ts", "**/*.js"},
        },
    },
    Body: "Variables and functions should have descriptive names that reveal intent.",
}
```

**Expected Output:**
```yaml
---
ruleset:
  id: cleanCode
  name: Clean Code
  description: Clean code practices
  rules:
    - meaningfulNames
    - smallFunctions
rule:
  id: meaningfulNames
  name: Use Meaningful Names
  description: Variables and functions should have descriptive names
  enforcement: must
  scope:
    files:
      - "**/*.ts"
      - "**/*.js"
---

# Use Meaningful Names (MUST)

Variables and functions should have descriptive names that reveal intent.
```

**Verification:**
- Metadata block present with all fields
- Enforcement header shows "MUST" (uppercased)
- Body content follows header
- Uses `GenerateRuleMetadataBlockFromRuleset()` for Ruleset

### Example 2: Minimal Metadata Block (Standalone Rule)

**Input:**
```go
Resource{
    Type: "rule",
    Ruleset: Ruleset{
        ID: "cleanCode",
        Name: "Clean Code",
        Rules: []string{"meaningfulNames"},
    },
    Rule: Rule{
        ID: "meaningfulNames",
        Name: "Use Meaningful Names",
        Enforcement: "should",
    },
    Body: "Use descriptive names.",
}
```

**Expected Output:**
```yaml
---
id: meaningfulNames
name: Use Meaningful Names
enforcement: should
---

# Use Meaningful Names (SHOULD)

Use descriptive names.
```

**Verification:**
- No ruleset section (standalone rule)
- No rule wrapper (fields at root level)
- Optional fields (description, scope) omitted
- Enforcement header shows "SHOULD"
- Minimal valid metadata block
- Uses `GenerateRuleMetadataBlockFromRule()` for standalone Rule

### Example 3: Prompt (No Metadata)

**Input:**
```go
Resource{
    Type: "prompt",
    Promptset: Promptset{
        ID: "codeReview",
        Name: "Code Review",
    },
    Prompt: Prompt{
        ID: "reviewPR",
        Name: "Review Pull Request",
    },
    Body: "Review this pull request for code quality issues.",
}
```

**Expected Output:**
```
Review this pull request for code quality issues.
```

**Verification:**
- No metadata block present
- No enforcement header
- Body content only

## Notes

The metadata block format is designed to be:
- **Human-readable** - YAML is easy to scan and understand
- **Machine-parseable** - Tools can extract context programmatically
- **Consistent** - Same structure across all target formats
- **Minimal** - Optional fields can be omitted to reduce noise

The enforcement header (generated internally) provides immediate visual feedback about rule importance without requiring tools to parse YAML.

Prompts intentionally exclude metadata because they represent reusable instructions rather than constraints with enforcement levels.

The two functions handle the distinct cases:
- **FromRuleset** - Rule is part of a collection, needs full context
- **FromRule** - Standalone rule, uses its own metadata for both ruleset and rule sections

## Known Issues

None - this is a new specification.

## Areas for Improvement

- Consider adding metadata version field for future schema evolution
- Evaluate whether prompts should include minimal metadata (promptset/prompt IDs)
- Explore compact metadata format for tools with strict size limits
