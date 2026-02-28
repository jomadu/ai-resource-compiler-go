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
func (c *CopilotCompiler) SupportedVersions() []string
func (c *CopilotCompiler) Compile(resource *airesource.Resource) ([]CompilationResult, error)
```

**Methods:**
- `Name()` - Returns "copilot"
- `SupportedVersions()` - Returns `["ai-resource/draft"]`
- `Compile()` - Transforms resource into Copilot format
  - Handles Rule, Ruleset, Prompt, Promptset kinds
  - Returns one result per rule/prompt

### ApplyTo Frontmatter (Rules and Prompts)
```yaml
---
applyTo: []string        # File patterns from scope
---
```

**Fields:**
- `applyTo` - File patterns extracted from Scope []ScopeEntry (empty array if no scope)

**Note:** excludeAgent field is omitted (not populated by compiler)

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

1. Check resource kind (Rule, Ruleset, Prompt, Promptset)
2. Expand collections into individual items
3. For each item:
   - Resolve body using `airesource.ResolveBody(body, fragments)`
   - Validate IDs using `ValidateID()`
   - For rules: validate name using `ValidateRuleName()`
   - Extract scope files from `[]ScopeEntry` using `extractScopeFiles()`
   - Generate applyTo frontmatter (both rules and prompts)
   - Generate path using shared path functions
   - Generate content (frontmatter + metadata + header + body for rules, frontmatter + body for prompts)
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
    
    // Generate applyTo frontmatter
    frontmatter = GenerateApplyToFrontmatter(scopeFiles)
    
    // Generate path
    if resource.Kind == "Ruleset":
        path = BuildCollectionPath(metadata.ID, ruleID, ".instructions.md")
    else:  // resource.Kind == "Rule"
        path = BuildStandalonePath(metadata.ID, ".instructions.md")
    
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
    
    // Extract scope files
    scopeFiles = extractScopeFiles(promptSpec.Scope)
    
    // Generate applyTo frontmatter
    frontmatter = GenerateApplyToFrontmatter(scopeFiles)
    
    // Generate path
    if resource.Kind == "Promptset":
        path = BuildCollectionPath(metadata.ID, promptID, ".prompt.md")
    else:  // resource.Kind == "Prompt"
        path = BuildStandalonePath(metadata.ID, ".prompt.md")
    
    // Use frontmatter + body
    content = frontmatter + "\n" + resolvedBody
    
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
| Unsupported apiVersion | Return error "unsupported apiVersion: {version} for copilot" |
| Ruleset with multiple rules | Return one CompilationResult per rule |
| Promptset with multiple prompts | Return one CompilationResult per prompt |

## Dependencies

- Resource model from ai-resource-core-go
- Metadata block generation from metadata-block.md spec
- TargetCompiler interface from compiler-architecture.md spec

## Implementation Mapping

**Source files:**
- `pkg/targets/copilot.go` - CopilotCompiler implementation
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
