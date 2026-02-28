# Cursor Compiler

## Job to be Done
Generate Cursor IDE rules (.mdc) and commands (.md) in the format expected by Cursor's context system.

## Activities
1. Compile rules with MDC frontmatter, metadata block, enforcement header, and body
2. Compile prompts with body content only (no frontmatter, no metadata)
3. Generate paths following {collection-id}_{item-id}.{ext} pattern
4. Produce CompilationResult with path and content
5. Document recommended installation directories

## Acceptance Criteria
- [ ] Rules include MDC frontmatter with description, globs, alwaysApply
- [ ] Rules include metadata block from metadata-block.md spec
- [ ] Rules include enforcement header (# {Name} ({ENFORCEMENT}))
- [ ] Rules use .mdc extension
- [ ] Prompts include body content only (no frontmatter, no metadata)
- [ ] Prompts use .md extension
- [ ] Paths follow {collection-id}_{item-id}.{ext} pattern
- [ ] Implements TargetCompiler interface
- [ ] Recommended installation: .cursor/rules/ (rules), .cursor/commands/ (prompts)

## Data Structures

### CursorCompiler
```go
type CursorCompiler struct{}

func (c *CursorCompiler) Name() string
func (c *CursorCompiler) SupportedVersions() []string
func (c *CursorCompiler) Compile(resource *airesource.Resource) ([]CompilationResult, error)
```

**Methods:**
- `Name()` - Returns "cursor"
- `SupportedVersions()` - Returns `["ai-resource/draft"]`
- `Compile()` - Transforms resource into Cursor format
  - Handles Rule, Ruleset, Prompt, Promptset kinds
  - Returns one result per rule/prompt

### MDC Frontmatter (Rules Only)
```yaml
---
description: string       # Rule description
globs: []string          # File patterns from scope
alwaysApply: bool        # true for must enforcement
---
```

**Fields:**
- `description` - Rule description (from RuleItem.Description or RuleItem.Name)
- `globs` - File patterns extracted from Scope []ScopeEntry (empty array if no scope)
- `alwaysApply` - true if enforcement is "must", false otherwise

**Scope Extraction:**
```go
func extractScopeFiles(scope []airesource.ScopeEntry) []string {
    var files []string
    for _, entry := range scope {
        files = append(files, entry.Files...)
    }
    return files
}
```

### Output Structure

**Rules (.mdc):**
```
---
description: string
globs: []string
alwaysApply: bool
---

---
{metadata block}
---

# {Rule Name} ({ENFORCEMENT})

{rule body}
```

**Prompts (.md):**
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
   - Extract scope files from `[]ScopeEntry` using `extractScopeFiles()`
   - Generate MDC frontmatter (rules only)
   - Generate path using shared path functions
   - Generate content (frontmatter + metadata + header + body for rules, body only for prompts)
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
    
    // Extract scope files
    scopeFiles = extractScopeFiles(ruleSpec.Scope)
    
    // Generate MDC frontmatter
    frontmatter = GenerateMDCFrontmatter(ruleSpec, scopeFiles)
    
    // Generate path
    if resource.Kind == "Ruleset":
        path = BuildCollectionPath(metadata.ID, ruleID, ".mdc")
    else:  // resource.Kind == "Rule"
        path = BuildStandalonePath(metadata.ID, ".mdc")
    
    // Generate complete content
    if resource.Kind == "Ruleset":
        content = GenerateRuleMetadataBlockFromRuleset(resource, ruleID)
    else:  // resource.Kind == "Rule"
        content = GenerateRuleMetadataBlockFromRule(resource)
    
    // Prepend frontmatter
    content = frontmatter + "\n" + content
    
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
| Rule without description | Use RuleItem.Name as description |
| Rule without scope | Set globs to empty array [] |
| Rule with "should" or "may" | Set alwaysApply to false |
| Rule with "must" | Set alwaysApply to true |
| Prompt resource | Return body only, no frontmatter or metadata |
| Empty body | Return frontmatter + metadata + header with empty body |
| Unsupported apiVersion | Return error "unsupported apiVersion: {version} for cursor" |
| Ruleset with multiple rules | Return one CompilationResult per rule |
| Promptset with multiple prompts | Return one CompilationResult per prompt |

## Dependencies

- Resource model from ai-resource-core-go
- Metadata block generation from metadata-block.md spec
- TargetCompiler interface from compiler-architecture.md spec

## Implementation Mapping

**Source files:**
- `pkg/targets/cursor.go` - CursorCompiler implementation
- `internal/format/metadata.go` - Shared metadata generation functions
- `internal/format/paths.go` - Shared path generation functions

**Related specs:**
- `metadata-block.md` - Metadata block structure and shared functions
- `compiler-architecture.md` - TargetCompiler interface, CompilationResult, and shared path functions

## Examples

### Example 1: Rule with Full Metadata

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
        Description: "Variables and functions should have descriptive names",
        Enforcement: "must",
        Scope: Scope{
            Files: []string{"**/*.ts", "**/*.js"},
        },
    },
    Body: "Use descriptive variable and function names that reveal intent.",
}

compiler := &CursorCompiler{}
results, err := compiler.Compile(resource)
```

**Expected Output:**
```go
[]CompilationResult{
    {
        Path: "cleanCode_meaningfulNames.mdc",
        Content: `---
description: Variables and functions should have descriptive names
globs:
  - "**/*.ts"
  - "**/*.js"
alwaysApply: true
---

---
ruleset:
  id: cleanCode
  name: Clean Code
  rules:
    - meaningfulNames
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

Use descriptive variable and function names that reveal intent.`,
    },
}
```

**Verification:**
- Path uses .mdc extension
- MDC frontmatter present with description, globs, alwaysApply
- Metadata block present
- Enforcement header shows "MUST"
- Body content preserved

**Installation:**
Write to `.cursor/rules/cleanCode_meaningfulNames.mdc`

### Example 2: Rule without Scope

**Input:**
```go
resource := Resource{
    Type: "rule",
    Ruleset: Ruleset{
        ID: "security",
        Name: "Security",
        Rules: []string{"noHardcodedSecrets"},
    },
    Rule: Rule{
        ID: "noHardcodedSecrets",
        Name: "No Hardcoded Secrets",
        Enforcement: "must",
    },
    Body: "Never commit secrets to version control.",
}

compiler := &CursorCompiler{}
results, err := compiler.Compile(resource)
```

**Expected Output:**
```go
[]CompilationResult{
    {
        Path: "security_noHardcodedSecrets.mdc",
        Content: `---
description: No Hardcoded Secrets
globs: []
alwaysApply: true
---

---
ruleset:
  id: security
  name: Security
  rules:
    - noHardcodedSecrets
rule:
  id: noHardcodedSecrets
  name: No Hardcoded Secrets
  enforcement: must
---

# No Hardcoded Secrets (MUST)

Never commit secrets to version control.`,
    },
}
```

**Verification:**
- Description uses rule.name (no description field)
- globs is empty array
- alwaysApply is true (must enforcement)

**Installation:**
Write to `.cursor/rules/security_noHardcodedSecrets.mdc`

### Example 3: Rule with "should" Enforcement

**Input:**
```go
resource := Resource{
    Type: "rule",
    Ruleset: Ruleset{
        ID: "style",
        Name: "Style Guide",
        Rules: []string{"preferConst"},
    },
    Rule: Rule{
        ID: "preferConst",
        Name: "Prefer Const",
        Enforcement: "should",
        Scope: Scope{
            Files: []string{"**/*.ts"},
        },
    },
    Body: "Use const instead of let when variables are not reassigned.",
}

compiler := &CursorCompiler{}
results, err := compiler.Compile(resource)
```

**Expected Output:**
```go
[]CompilationResult{
    {
        Path: "style_preferConst.mdc",
        Content: `---
description: Prefer Const
globs:
  - "**/*.ts"
alwaysApply: false
---

---
ruleset:
  id: style
  name: Style Guide
  rules:
    - preferConst
rule:
  id: preferConst
  name: Prefer Const
  enforcement: should
  scope:
    files:
      - "**/*.ts"
---

# Prefer Const (SHOULD)

Use const instead of let when variables are not reassigned.`,
    },
}
```

**Verification:**
- alwaysApply is false (should enforcement)
- Enforcement header shows "SHOULD"

**Installation:**
Write to `.cursor/rules/style_preferConst.mdc`

### Example 4: Prompt Compilation

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
    Body: "Review this pull request for code quality and security issues.",
}

compiler := &CursorCompiler{}
results, err := compiler.Compile(resource)
```

**Expected Output:**
```go
[]CompilationResult{
    {
        Path: "codeReview_reviewPR.md",
        Content: "Review this pull request for code quality and security issues.",
    },
}
```

**Verification:**
- Path uses .md extension (not .mdc)
- No MDC frontmatter
- No metadata block
- Body content only

**Installation:**
Write to `.cursor/commands/codeReview_reviewPR.md`

## Notes

**Design Rationale:**
- **MDC frontmatter** - Cursor-specific format for rule configuration
- **Metadata block** - Preserves context for other tools and human readers
- **alwaysApply mapping** - "must" enforcement → always apply, others → conditional
- **Extension differentiation** - .mdc for rules, .md for prompts
- **Prompts are simple** - No frontmatter or metadata, just body content

**Cursor Integration:**
- Cursor IDE reads rules from `.cursor/rules/` directory
- Cursor IDE reads commands from `.cursor/commands/` directory
- MDC frontmatter controls when rules are applied
- alwaysApply=true means rule is always active
- globs restrict rule to specific file patterns

**Installation Directories:**
- **Rules:** `.cursor/rules/` - Rules that guide Cursor's behavior
- **Prompts:** `.cursor/commands/` - Reusable commands for user invocation

## Known Issues

None - this is a new specification.

## Areas for Improvement

- Consider supporting additional MDC frontmatter fields as Cursor evolves
- Evaluate integration with Cursor's context management features
- Explore dynamic rule loading and hot-reloading
