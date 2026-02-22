# Kiro Compiler

## Job to be Done
Generate Kiro CLI steering rules and prompts in the format expected by Kiro's context system.

## Activities
1. Compile rules with metadata block, enforcement header, and body
2. Compile prompts with body content only
3. Generate paths following {collection-id}_{item-id}.md pattern
4. Produce CompilationResult with path and content
5. Document recommended installation directories

## Acceptance Criteria
- [ ] Rules include metadata block from metadata-block.md spec
- [ ] Rules include enforcement header (# {Name} ({ENFORCEMENT}))
- [ ] Prompts include body content only (no metadata, no header)
- [ ] All outputs use .md extension
- [ ] Paths follow {collection-id}_{item-id}.md pattern
- [ ] No frontmatter added (metadata block only for rules)
- [ ] Implements TargetCompiler interface
- [ ] Recommended installation: .kiro/steering/ (rules), .kiro/prompts/ (prompts)

## Data Structures

### KiroCompiler
```go
type KiroCompiler struct{}

func (k *KiroCompiler) Name() string
func (k *KiroCompiler) Compile(resource Resource) ([]CompilationResult, error)
```

**Methods:**
- `Name()` - Returns "kiro"
- `Compile()` - Transforms resource into Kiro format

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

1. Determine resource type (rule vs prompt)
2. Generate path: `{collection-id}_{item-id}.md`
3. If rule:
   - Generate metadata block (ruleset + rule sections)
   - Generate enforcement header
   - Concatenate: metadata + header + body
4. If prompt:
   - Use body content only
5. Return CompilationResult with path and content

**Pseudocode:**
```
function Compile(resource):
    path = resource.collectionID + "_" + resource.itemID + ".md"
    
    if resource.type == "rule":
        metadata = GenerateMetadataBlock(resource)
        header = "# " + resource.rule.name + " (" + uppercase(resource.rule.enforcement) + ")"
        content = metadata + "\n" + header + "\n\n" + resource.body
    else:
        content = resource.body
    
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

## Dependencies

- Resource model from ai-resource-core-go
- Metadata block generation from metadata-block.md spec
- TargetCompiler interface from compiler-architecture.md spec

## Implementation Mapping

**Source files:**
- `pkg/targets/kiro.go` - KiroCompiler implementation

**Related specs:**
- `metadata-block.md` - Metadata block structure and generation
- `compiler-architecture.md` - TargetCompiler interface and CompilationResult

## Examples

### Example 1: Rule Compilation

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
        Enforcement: "must",
        Scope: Scope{
            Files: []string{"**/*.go"},
        },
    },
    Body: "Use descriptive variable and function names.",
}

compiler := &KiroCompiler{}
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
  enforcement: must
  scope:
    files:
      - "**/*.go"
---

# Use Meaningful Names (MUST)

Use descriptive variable and function names.`,
    },
}
```

**Verification:**
- Path follows {collection-id}_{item-id}.md pattern
- Metadata block present
- Enforcement header shows "MUST"
- Body content preserved

**Installation:**
Write to `.kiro/steering/cleanCode_meaningfulNames.md`

### Example 2: Prompt Compilation

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

compiler := &KiroCompiler{}
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
- No metadata block
- No enforcement header
- Body content only
- Path follows {collection-id}_{item-id}.md pattern

**Installation:**
Write to `.kiro/prompts/codeReview_reviewPR.md`

### Example 3: Minimal Rule

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

compiler := &KiroCompiler{}
results, err := compiler.Compile(resource)
```

**Expected Output:**
```go
[]CompilationResult{
    {
        Path: "security_noHardcodedSecrets.md",
        Content: `---
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
- Minimal metadata (no optional fields)
- Enforcement header present
- Body preserved

**Installation:**
Write to `.kiro/steering/security_noHardcodedSecrets.md`

## Notes

**Design Rationale:**
- **Identical to markdown target** - Kiro accepts standard markdown with metadata blocks
- **Metadata block provides context** - Kiro can parse ruleset/rule relationships
- **Enforcement headers** - Visual indication of rule importance
- **Separate directories** - Rules in steering/, prompts in prompts/

**Kiro Integration:**
- Kiro CLI reads steering rules from `.kiro/steering/`
- Kiro CLI reads prompts from `.kiro/prompts/`
- Metadata blocks allow Kiro to understand rule context
- Enforcement levels guide Kiro's behavior

**Installation Directories:**
- **Rules:** `.kiro/steering/` - Steering rules that guide Kiro's behavior
- **Prompts:** `.kiro/prompts/` - Reusable prompts for user invocation

## Known Issues

None - this is a new specification.

## Areas for Improvement

- Consider Kiro-specific metadata extensions for advanced features
- Evaluate integration with Kiro's context management system
- Explore dynamic rule loading and hot-reloading
