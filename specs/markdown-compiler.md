# Markdown Compiler

## Job to be Done
Generate vanilla markdown output for AI resources that can be used with any tool or system that accepts markdown files.

## Activities
1. Compile rules with metadata block, enforcement header, and body
2. Compile prompts with body content only
3. Generate paths following {collection-id}_{item-id}.md pattern
4. Produce CompilationResult with path and content

## Acceptance Criteria
- [ ] Rules include metadata block from metadata-block.md spec
- [ ] Rules include enforcement header (# {Name} ({ENFORCEMENT}))
- [ ] Prompts include body content only (no metadata, no header)
- [ ] All outputs use .md extension
- [ ] Paths follow {collection-id}_{item-id}.md pattern
- [ ] No frontmatter added (metadata block only for rules)
- [ ] Implements TargetCompiler interface

## Data Structures

### MarkdownCompiler
```go
type MarkdownCompiler struct{}

func (m *MarkdownCompiler) Name() string
func (m *MarkdownCompiler) SupportedVersions() []string
func (m *MarkdownCompiler) Compile(resource *airesource.Resource) ([]CompilationResult, error)
```

**Methods:**
- `Name()` - Returns "markdown"
- `SupportedVersions()` - Returns `["ai-resource/draft"]`
- `Compile()` - Transforms resource into markdown format
  - Handles Rule, Ruleset, Prompt, Promptset kinds
  - Returns one result per rule/prompt

### Output Structure

**Rules:**
```
---
{metadata block}
---

# {Rule Name} ({ENFORCEMENT})

{rule body}
```

**Prompts:**
```
{prompt body}
```

## Algorithm

1. Check resource kind (Rule, Ruleset, Prompt, Promptset)
2. Expand collections into individual items
3. For each item:
   - Resolve body using `airesource.ResolveBody(body, fragments)`
   - Validate IDs using `ValidateID()`
   - For rules: validate name using `ValidateRuleName()`
   - Generate path using shared path functions
   - Generate content (metadata + header + body for rules, body only for prompts)
4. Return array of CompilationResults

**Pseudocode:**
```
function Compile(resource):
    results = []
    
    switch resource.Kind:
    case "Rule":
        rule = resource.AsRule()
        result = compileRule(rule.Metadata, rule.Metadata.ID, rule.Spec, rule.Spec.Fragments)
        results.append(result)
    
    case "Ruleset":
        ruleset = resource.AsRuleset()
        for ruleID, ruleItem in ruleset.Spec.Rules:
            result = compileRule(ruleset.Metadata, ruleID, ruleItem, ruleset.Spec.Fragments)
            results.append(result)
    
    case "Prompt":
        prompt = resource.AsPrompt()
        result = compilePrompt(prompt.Metadata, prompt.Metadata.ID, prompt.Spec, prompt.Spec.Fragments)
        results.append(result)
    
    case "Promptset":
        promptset = resource.AsPromptset()
        for promptID, promptItem in promptset.Spec.Prompts:
            result = compilePrompt(promptset.Metadata, promptID, promptItem, promptset.Spec.Fragments)
            results.append(result)
    
    return results

function compileRule(metadata, ruleID, ruleSpec, fragments):
    // Resolve body
    resolvedBody = airesource.ResolveBody(ruleSpec.Body, fragments)
    
    // Validate
    ValidateID(metadata.ID)
    ValidateID(ruleID)
    ValidateRuleName(ruleSpec.Name)
    
    // Generate path
    if resource.Kind == "Ruleset":
        path = BuildCollectionPath(metadata.ID, ruleID, ".md")
    else:  // resource.Kind == "Rule"
        path = BuildStandalonePath(metadata.ID, ".md")
    
    // Generate complete content
    if resource.Kind == "Ruleset":
        content = GenerateRuleMetadataBlockFromRuleset(resource, ruleID)
    else:  // resource.Kind == "Rule"
        content = GenerateRuleMetadataBlockFromRule(resource)
    
    return CompilationResult{Path: path, Content: content}

function compilePrompt(metadata, promptID, promptSpec, fragments):
    // Resolve body
    resolvedBody = airesource.ResolveBody(promptSpec.Body, fragments)
    
    // Validate
    ValidateID(metadata.ID)
    ValidateID(promptID)
    
    // Generate path
    if resource.Kind == "Promptset":
        path = BuildCollectionPath(metadata.ID, promptID, ".md")
    else:  // resource.Kind == "Prompt"
        path = BuildStandalonePath(metadata.ID, ".md")
    
    // Use body only
    content = resolvedBody
    
    return CompilationResult{Path: path, Content: content}
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Rule with minimal metadata | Include only required fields in metadata block |
| Prompt resource | Return body only, no metadata or header |
| Empty body | Return metadata + header with empty body section |
| Special characters in IDs | Use IDs as-is in path (sanitization handled by caller) |
| Multi-line body | Preserve formatting and line breaks |
| Unsupported apiVersion | Return error "unsupported apiVersion: {version} for markdown" |

## Dependencies

- Resource model from ai-resource-core-go
- Metadata block generation from metadata-block.md spec
- TargetCompiler interface from compiler-architecture.md spec

## Implementation Mapping

**Source files:**
- `pkg/targets/markdown.go` - MarkdownCompiler implementation
- `internal/format/metadata.go` - Shared metadata generation functions
- `internal/format/paths.go` - Shared path generation functions

**Related specs:**
- `metadata-block.md` - Metadata block structure and shared functions
- `compiler-architecture.md` - TargetCompiler interface, CompilationResult, and shared path functions

## Examples

### Example 1: Rule Compilation

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

compiler := &MarkdownCompiler{}
results, err := compiler.Compile(resource)
```

**Expected Output:**
```go
[]CompilationResult{
    {
        Path: "cleanCode_meaningfulNames.md",
        Content: `---
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

Variables and functions should have descriptive names that reveal intent.`,
    },
}
```

**Verification:**
- Path follows {collection-id}_{item-id}.md pattern
- Metadata block includes all fields
- Enforcement header shows "MUST" (uppercased)
- Body content preserved

### Example 2: Minimal Rule

**Input:**
```go
resource := Resource{
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

compiler := &MarkdownCompiler{}
results, err := compiler.Compile(resource)
```

**Expected Output:**
```go
[]CompilationResult{
    {
        Path: "cleanCode_meaningfulNames.md",
        Content: `---
ruleset:
  id: cleanCode
  name: Clean Code
  rules:
    - meaningfulNames
rule:
  id: meaningfulNames
  name: Use Meaningful Names
  enforcement: should
---

# Use Meaningful Names (SHOULD)

Use descriptive names.`,
    },
}
```

**Verification:**
- Optional fields omitted from metadata
- Enforcement header shows "SHOULD"
- Minimal valid output

### Example 3: Prompt Compilation

**Input:**
```go
resource := Resource{
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

compiler := &MarkdownCompiler{}
results, err := compiler.Compile(resource)
```

**Expected Output:**
```go
[]CompilationResult{
    {
        Path: "codeReview_reviewPR.md",
        Content: "Review this pull request for code quality issues.",
    },
}
```

**Verification:**
- No metadata block
- No enforcement header
- Body content only
- Path follows {collection-id}_{item-id}.md pattern

## Notes

**Design Rationale:**
- **Simplest target** - Demonstrates metadata block usage without additional formatting
- **Universal compatibility** - Plain markdown works everywhere
- **No frontmatter** - Metadata block is the only structured data
- **Consistent with other targets** - Same metadata block format used by all targets

**Use Cases:**
- Generic AI tools that accept markdown
- Documentation generation
- Version control and review
- Testing and validation reference

**Installation:**
No specific installation directory - users decide where to place files.

## Known Issues

None - this is a new specification.

## Areas for Improvement

- Consider adding optional frontmatter for tools that support it
- Evaluate compact metadata format option for size-constrained scenarios
