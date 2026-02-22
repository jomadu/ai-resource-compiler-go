# Copilot Compiler

## Job to be Done
Generate GitHub Copilot instructions and prompts in the format expected by Copilot's context system.

## Activities
1. Compile rules with applyTo frontmatter, metadata block, enforcement header, and body
2. Compile prompts with applyTo frontmatter and body content (no metadata)
3. Generate paths following {collection-id}_{item-id}.{ext} pattern
4. Produce CompilationResult with path and content
5. Document recommended installation directories

## Acceptance Criteria
- [ ] Rules include applyTo frontmatter with file patterns
- [ ] Rules include metadata block from metadata-block.md spec
- [ ] Rules include enforcement header (# {Name} ({ENFORCEMENT}))
- [ ] Rules use .instructions.md extension
- [ ] Prompts include applyTo frontmatter with file patterns
- [ ] Prompts include body content only (no metadata, no header)
- [ ] Prompts use .prompt.md extension
- [ ] Paths follow {collection-id}_{item-id}.{ext} pattern
- [ ] Implements TargetCompiler interface
- [ ] Recommended installation: .github/instructions/ (rules), .github/prompts/ (prompts)

## Data Structures

### CopilotCompiler
```go
type CopilotCompiler struct{}

func (c *CopilotCompiler) Name() string
func (c *CopilotCompiler) Compile(resource Resource) ([]CompilationResult, error)
```

**Methods:**
- `Name()` - Returns "copilot"
- `Compile()` - Transforms resource into Copilot format

### ApplyTo Frontmatter (Rules and Prompts)
```yaml
---
applyTo: []string        # File patterns from scope
---
```

**Fields:**
- `applyTo` - File patterns from scope.files (empty array if no scope)

**Note:** excludeAgent field is omitted (not populated by compiler)

### Output Structure

**Rules (.instructions.md):**
```
---
applyTo: []string
---

---
{metadata block}
---

# {Rule Name} ({ENFORCEMENT})

{rule body}
```

**Prompts (.prompt.md):**
```
---
applyTo: []string
---

{prompt body}
```

## Algorithm

1. Determine resource type (rule vs prompt)
2. If rule:
   - Generate path: `{collection-id}_{item-id}.instructions.md`
   - Generate applyTo frontmatter from rule.scope.files
   - Generate metadata block
   - Generate enforcement header
   - Concatenate: frontmatter + metadata + header + body
3. If prompt:
   - Generate path: `{collection-id}_{item-id}.prompt.md`
   - Generate applyTo frontmatter from prompt.scope.files
   - Use body content only (no metadata, no header)
4. Return CompilationResult with path and content

**Pseudocode:**
```
function Compile(resource):
    if resource.type == "rule":
        path = resource.collectionID + "_" + resource.itemID + ".instructions.md"
        
        frontmatter = GenerateApplyToFrontmatter(resource.rule.scope.files)
        metadata = GenerateMetadataBlock(resource)
        header = "# " + resource.rule.name + " (" + uppercase(resource.rule.enforcement) + ")"
        content = frontmatter + "\n" + metadata + "\n" + header + "\n\n" + resource.body
    else:
        path = resource.collectionID + "_" + resource.itemID + ".prompt.md"
        
        frontmatter = GenerateApplyToFrontmatter(resource.prompt.scope.files)
        content = frontmatter + "\n" + resource.body
    
    return CompilationResult{Path: path, Content: content}
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Rule without scope | Set applyTo to empty array [] |
| Prompt without scope | Set applyTo to empty array [] |
| Empty body | Return frontmatter + [metadata + header] with empty body |
| Special characters in IDs | Use IDs as-is in path (sanitization handled by caller) |
| Multi-line body | Preserve formatting and line breaks |

## Dependencies

- Resource model from ai-resource-core-go
- Metadata block generation from metadata-block.md spec
- TargetCompiler interface from compiler-architecture.md spec

## Implementation Mapping

**Source files:**
- `pkg/targets/copilot.go` - CopilotCompiler implementation

**Related specs:**
- `metadata-block.md` - Metadata block structure and generation
- `compiler-architecture.md` - TargetCompiler interface and CompilationResult

## Examples

### Example 1: Rule with Scope

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
            Files: []string{"**/*.go", "**/*.py"},
        },
    },
    Body: "Use descriptive variable and function names.",
}

compiler := &CopilotCompiler{}
results, err := compiler.Compile(resource)
```

**Expected Output:**
```go
[]CompilationResult{
    {
        Path: "cleanCode_meaningfulNames.instructions.md",
        Content: `---
applyTo:
  - "**/*.go"
  - "**/*.py"
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
  enforcement: must
  scope:
    files:
      - "**/*.go"
      - "**/*.py"
---

# Use Meaningful Names (MUST)

Use descriptive variable and function names.`,
    },
}
```

**Verification:**
- Path uses .instructions.md extension
- applyTo frontmatter present with file patterns
- Metadata block present
- Enforcement header shows "MUST"
- Body content preserved

**Installation:**
Write to `.github/instructions/cleanCode_meaningfulNames.instructions.md`

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

compiler := &CopilotCompiler{}
results, err := compiler.Compile(resource)
```

**Expected Output:**
```go
[]CompilationResult{
    {
        Path: "security_noHardcodedSecrets.instructions.md",
        Content: `---
applyTo: []
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
- applyTo is empty array (no scope)
- Metadata block present
- Enforcement header shows "MUST"

**Installation:**
Write to `.github/instructions/security_noHardcodedSecrets.instructions.md`

### Example 3: Prompt with Scope

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
        Scope: Scope{
            Files: []string{"**/*.ts", "**/*.tsx"},
        },
    },
    Body: "Review this pull request for code quality and security issues.",
}

compiler := &CopilotCompiler{}
results, err := compiler.Compile(resource)
```

**Expected Output:**
```go
[]CompilationResult{
    {
        Path: "codeReview_reviewPR.prompt.md",
        Content: `---
applyTo:
  - "**/*.ts"
  - "**/*.tsx"
---

Review this pull request for code quality and security issues.`,
    },
}
```

**Verification:**
- Path uses .prompt.md extension
- applyTo frontmatter present
- No metadata block
- No enforcement header
- Body content only

**Installation:**
Write to `.github/prompts/codeReview_reviewPR.prompt.md`

### Example 4: Prompt without Scope

**Input:**
```go
resource := Resource{
    Type: "prompt",
    Promptset: Promptset{
        ID: "general",
        Name: "General",
    },
    Prompt: Prompt{
        ID: "explainCode",
        Name: "Explain Code",
    },
    Body: "Explain what this code does in simple terms.",
}

compiler := &CopilotCompiler{}
results, err := compiler.Compile(resource)
```

**Expected Output:**
```go
[]CompilationResult{
    {
        Path: "general_explainCode.prompt.md",
        Content: `---
applyTo: []
---

Explain what this code does in simple terms.`,
    },
}
```

**Verification:**
- applyTo is empty array (no scope)
- No metadata block
- Body content only

**Installation:**
Write to `.github/prompts/general_explainCode.prompt.md`

## Notes

**Design Rationale:**
- **applyTo frontmatter** - Copilot-specific format for file pattern matching
- **Metadata block for rules** - Preserves context for other tools and human readers
- **No metadata for prompts** - Prompts are simpler, just frontmatter + body
- **Extension differentiation** - .instructions.md for rules, .prompt.md for prompts
- **excludeAgent omitted** - Not populated by compiler (user can add manually if needed)

**Copilot Integration:**
- GitHub Copilot reads instructions from `.github/instructions/` directory
- GitHub Copilot reads prompts from `.github/prompts/` directory
- applyTo frontmatter restricts when instructions/prompts are active
- Empty applyTo array means applies to all files

**Installation Directories:**
- **Rules:** `.github/instructions/` - Instructions that guide Copilot's behavior
- **Prompts:** `.github/prompts/` - Reusable prompts for user invocation

## Known Issues

None - this is a new specification.

## Areas for Improvement

- Consider supporting excludeAgent field for advanced filtering
- Evaluate additional Copilot-specific frontmatter fields
- Explore integration with GitHub Copilot's context management features
