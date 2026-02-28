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
func (c *ClaudeCompiler) SupportedVersions() []string
func (c *ClaudeCompiler) Compile(resource *airesource.Resource) ([]CompilationResult, error)
```

**Methods:**
- `Name()` - Returns "claude"
- `SupportedVersions()` - Returns `["ai-resource/draft"]`
- `Compile()` - Transforms resource into Claude format
  - Handles Rule, Ruleset, Prompt, Promptset kinds
  - Returns one result per rule/prompt

### Paths Frontmatter (Rules Only, Optional)
```yaml
---
paths:
  - string  # File patterns from scope
---
```

**Fields:**
- `paths` - File patterns extracted from Scope []ScopeEntry (omit frontmatter if no scope)

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

1. Check resource kind (Rule, Ruleset, Prompt, Promptset)
2. Expand collections into individual items
3. For each item:
   - Resolve body using `airesource.ResolveBody(body, fragments)`
   - Validate IDs using `ValidateID()`
   - For rules: validate name using `ValidateRuleName()`
   - Extract scope files from `[]ScopeEntry` using `extractScopeFiles()`
   - Generate optional paths frontmatter (rules only, if scope defined)
   - Generate path using shared path functions
   - Generate content (optional frontmatter + metadata + header + body for rules, body only for prompts)
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
    
    // Prepend optional paths frontmatter
    if len(scopeFiles) > 0:
        frontmatter = GeneratePathsFrontmatter(scopeFiles)
        content = frontmatter + "\n" + content
    
    return CompilationResult{Path: path, Content: content}

function compilePrompt(metadata, promptID, promptSpec, fragments):
    // Resolve body
    resolvedBody = airesource.ResolveBody(promptSpec.Body, fragments)
    
    // Validate
    ValidateID(metadata.ID)
    ValidateID(promptID)
    
    // Generate path (special SKILL.md structure)
    if resource.Kind == "Promptset":
        path = BuildClaudeCollectionPath(metadata.ID, promptID)
    else:  // resource.Kind == "Prompt"
        path = BuildClaudeStandalonePath(metadata.ID)
    
    // Use body only
    content = resolvedBody
    
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
| Ruleset with multiple rules | Return one CompilationResult per rule |
| Promptset with multiple prompts | Return one CompilationResult per prompt |

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
