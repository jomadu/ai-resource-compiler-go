# Claude Compiler

## Job to be Done
Generate Claude Code rules and skills in the format expected by Claude's context system.

## Activities
1. Compile rules with optional paths frontmatter, metadata block, enforcement header, and body
2. Compile prompts as SKILL.md files with body content only (no frontmatter, no metadata)
3. Generate paths: {collection-id}_{item-id}.md (rules), {collection-id}_{item-id}/SKILL.md (prompts)
4. Produce CompilationResult with path and content
5. Document recommended installation directories

## Acceptance Criteria
- [ ] Rules include optional paths frontmatter (only if scope defined)
- [ ] Rules include metadata block from metadata-block.md spec
- [ ] Rules include enforcement header (# {Name} ({ENFORCEMENT}))
- [ ] Rules use .md extension
- [ ] Prompts use {collection-id}_{item-id}/SKILL.md path structure
- [ ] Prompts include body content only (no frontmatter, no metadata)
- [ ] Implements TargetCompiler interface
- [ ] Recommended installation: .claude/rules/ (rules), .claude/skills/ (prompts)

## Data Structures

### ClaudeCompiler
```go
type ClaudeCompiler struct{}

func (c *ClaudeCompiler) Name() string
func (c *ClaudeCompiler) Compile(resource Resource) ([]CompilationResult, error)
func (c *ClaudeCompiler) SupportedVersions() []string
```

**Methods:**
- `Name()` - Returns "claude"
- `Compile()` - Transforms resource into Claude format
- `SupportedVersions()` - Returns `["ai-resource/v1"]`

### Paths Frontmatter (Rules Only, Optional)
```yaml
---
paths:
  - string  # File patterns from scope
---
```

**Fields:**
- `paths` - File patterns from rule.scope.files (omit frontmatter if no scope)

### Output Structure

**Rules (.md) with scope:**
```
---
paths:
  - string
---

---
{metadata block}
---

# {Rule Name} ({ENFORCEMENT})

{rule body}
```

**Rules (.md) without scope:**
```
---
{metadata block}
---

# {Rule Name} ({ENFORCEMENT})

{rule body}
```

**Prompts (SKILL.md):**
```
{prompt body}
```

## Algorithm

1. Determine resource type (rule vs prompt)
2. Generate path using shared path functions
3. If rule:
   - If scope defined: Generate paths frontmatter
   - Call `GenerateMetadataBlock(ruleset, rule)` from `internal/format/metadata.go`
   - Call `GenerateEnforcementHeader(rule)` from `internal/format/metadata.go`
   - Concatenate: [frontmatter +] metadata + header + body
4. If prompt:
   - Use body content only
5. Return CompilationResult with path and content

**Pseudocode:**
```
function Compile(resource):
    if resource.type == "rule":
        path = BuildRulePath(resource.rulesetID, resource.ruleID, ".md")
        
        content = ""
        if resource.rule.scope.files is not empty:
            frontmatter = GeneratePathsFrontmatter(resource.rule.scope.files)
            content += frontmatter + "\n"
        
        metadata = GenerateMetadataBlock(resource.ruleset, resource.rule)
        header = GenerateEnforcementHeader(resource.rule)
        content += metadata + "\n" + header + "\n\n" + resource.body
    else:
        path = BuildClaudePromptPath(resource.promptsetID, resource.promptID)
        content = resource.body
    
    return CompilationResult{Path: path, Content: content}
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Rule without scope | Omit paths frontmatter entirely |
| Rule with scope | Include paths frontmatter |
| Prompt resource | Use directory/SKILL.md path, body only |
| Empty body | Return [frontmatter +] metadata + header with empty body |
| Special characters in IDs | Use IDs as-is in path (sanitization handled by caller) |
| Unsupported apiVersion | Return error "unsupported apiVersion: {version} for claude" |

## Dependencies

- Resource model from ai-resource-core-go
- Metadata block generation from metadata-block.md spec
- TargetCompiler interface from compiler-architecture.md spec

## Implementation Mapping

**Source files:**
- `pkg/targets/claude.go` - ClaudeCompiler implementation
- `internal/format/metadata.go` - Shared metadata generation functions
- `internal/format/paths.go` - Shared path generation functions

**Related specs:**
- `metadata-block.md` - Metadata block structure and shared functions
- `compiler-architecture.md` - TargetCompiler interface, CompilationResult, and shared path functions

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
            Files: []string{"**/*.py", "**/*.js"},
        },
    },
    Body: "Use descriptive variable and function names.",
}

compiler := &ClaudeCompiler{}
results, err := compiler.Compile(resource)
```

**Expected Output:**
```go
[]CompilationResult{
    {
        Path: "cleanCode_meaningfulNames.md",
        Content: `---
paths:
  - "**/*.py"
  - "**/*.js"
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
      - "**/*.py"
      - "**/*.js"
---

# Use Meaningful Names (MUST)

Use descriptive variable and function names.`,
    },
}
```

**Verification:**
- Path uses .md extension
- Paths frontmatter present with file patterns
- Metadata block present
- Enforcement header shows "MUST"
- Body content preserved

**Installation:**
Write to `.claude/rules/cleanCode_meaningfulNames.md`

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

compiler := &ClaudeCompiler{}
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
- No paths frontmatter (no scope defined)
- Metadata block present
- Enforcement header shows "MUST"

**Installation:**
Write to `.claude/rules/security_noHardcodedSecrets.md`

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
    Body: "Review this pull request for code quality and security issues.",
}

compiler := &ClaudeCompiler{}
results, err := compiler.Compile(resource)
```

**Expected Output:**
```go
[]CompilationResult{
    {
        Path: "codeReview_reviewPR/SKILL.md",
        Content: "Review this pull request for code quality and security issues.",
    },
}
```

**Verification:**
- Path uses directory/SKILL.md structure
- No paths frontmatter
- No metadata block
- Body content only

**Installation:**
Write to `.claude/skills/codeReview_reviewPR/SKILL.md`

### Example 4: Rule with "should" Enforcement

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

compiler := &ClaudeCompiler{}
results, err := compiler.Compile(resource)
```

**Expected Output:**
```go
[]CompilationResult{
    {
        Path: "style_preferConst.md",
        Content: `---
paths:
  - "**/*.ts"
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
- Paths frontmatter present
- Enforcement header shows "SHOULD"
- Metadata block includes enforcement level

**Installation:**
Write to `.claude/rules/style_preferConst.md`

## Notes

**Design Rationale:**
- **Optional paths frontmatter** - Only include when scope is defined
- **Metadata block** - Preserves context for other tools and human readers
- **SKILL.md convention** - Claude's standard for skill definitions
- **Directory structure for prompts** - Allows additional files per skill in future

**Claude Integration:**
- Claude Code reads rules from `.claude/rules/` directory
- Claude Code reads skills from `.claude/skills/` directory
- Paths frontmatter restricts rule to specific file patterns
- Skills are reusable prompts invoked by name

**Installation Directories:**
- **Rules:** `.claude/rules/` - Rules that guide Claude's behavior
- **Prompts:** `.claude/skills/` - Reusable skills for user invocation

## Known Issues

None - this is a new specification.

## Areas for Improvement

- Consider supporting additional Claude-specific frontmatter fields
- Evaluate integration with Claude's context management features
- Explore multi-file skills (additional files beyond SKILL.md)
