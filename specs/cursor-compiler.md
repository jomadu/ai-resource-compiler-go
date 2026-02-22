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
func (c *CursorCompiler) Compile(resource Resource) ([]CompilationResult, error)
```

**Methods:**
- `Name()` - Returns "cursor"
- `Compile()` - Transforms resource into Cursor format

### MDC Frontmatter (Rules Only)
```yaml
---
description: string       # Rule description
globs: []string          # File patterns from scope
alwaysApply: bool        # true for must enforcement
---
```

**Fields:**
- `description` - Rule description (from rule.description or rule.name)
- `globs` - File patterns from rule.scope.files (empty array if no scope)
- `alwaysApply` - true if enforcement is "must", false otherwise

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

1. Determine resource type (rule vs prompt)
2. If rule:
   - Generate path: `{collection-id}_{item-id}.mdc`
   - Generate MDC frontmatter
   - Generate metadata block
   - Generate enforcement header
   - Concatenate: frontmatter + metadata + header + body
3. If prompt:
   - Generate path: `{collection-id}_{item-id}.md`
   - Use body content only
4. Return CompilationResult with path and content

**Pseudocode:**
```
function Compile(resource):
    if resource.type == "rule":
        path = resource.collectionID + "_" + resource.itemID + ".mdc"
        
        frontmatter = GenerateMDCFrontmatter(resource.rule)
        metadata = GenerateMetadataBlock(resource)
        header = "# " + resource.rule.name + " (" + uppercase(resource.rule.enforcement) + ")"
        content = frontmatter + "\n" + metadata + "\n" + header + "\n\n" + resource.body
    else:
        path = resource.collectionID + "_" + resource.itemID + ".md"
        content = resource.body
    
    return CompilationResult{Path: path, Content: content}
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Rule without description | Use rule.name as description |
| Rule without scope | Set globs to empty array [] |
| Rule with "should" or "may" | Set alwaysApply to false |
| Rule with "must" | Set alwaysApply to true |
| Prompt resource | Return body only, no frontmatter or metadata |
| Empty body | Return frontmatter + metadata + header with empty body |

## Dependencies

- Resource model from ai-resource-core-go
- Metadata block generation from metadata-block.md spec
- TargetCompiler interface from compiler-architecture.md spec

## Implementation Mapping

**Source files:**
- `pkg/targets/cursor.go` - CursorCompiler implementation

**Related specs:**
- `metadata-block.md` - Metadata block structure and generation
- `compiler-architecture.md` - TargetCompiler interface and CompilationResult

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
