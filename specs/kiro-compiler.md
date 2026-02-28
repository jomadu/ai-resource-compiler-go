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
func (k *KiroCompiler) SupportedVersions() []string
func (k *KiroCompiler) Compile(resource *airesource.Resource) ([]CompilationResult, error)
```

**Methods:**
- `Name()` - Returns "kiro"
- `SupportedVersions()` - Returns `["ai-resource/draft"]`
- `Compile()` - Transforms resource into Kiro format
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
| Unsupported apiVersion | Return error "unsupported apiVersion: {version} for kiro" |

## Dependencies

- Resource model from ai-resource-core-go
- Metadata block generation from metadata-block.md spec
- TargetCompiler interface from compiler-architecture.md spec

## Implementation Mapping

**Source files:**
- `pkg/targets/kiro.go` - KiroCompiler implementation
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

### Example 4: Ruleset Expansion

**Input:**
```go
resource := Resource{
    Kind: "Ruleset",
    Metadata: Metadata{
        ID: "cleanCode",
        Name: "Clean Code",
    },
    Spec: RulesetSpec{
        Rules: map[string]RuleItem{
            "meaningfulNames": {
                Name: "Use Meaningful Names",
                Enforcement: "must",
                Body: "Use descriptive names.",
            },
            "smallFunctions": {
                Name: "Keep Functions Small",
                Enforcement: "should",
                Body: "Functions should do one thing.",
            },
        },
    },
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
    - smallFunctions
rule:
  id: meaningfulNames
  name: Use Meaningful Names
  enforcement: must
---

# Use Meaningful Names (MUST)

Use descriptive names.`,
    },
    {
        Path: "cleanCode_smallFunctions.md",
        Content: `---
ruleset:
  id: cleanCode
  name: Clean Code
  rules:
    - meaningfulNames
    - smallFunctions
rule:
  id: smallFunctions
  name: Keep Functions Small
  enforcement: should
---

# Keep Functions Small (SHOULD)

Functions should do one thing.`,
    },
}
```

**Verification:**
- Ruleset expanded into two separate files
- Each file includes full ruleset context
- Rules list shows all rules in ruleset
- Paths follow {collection-id}_{item-id}.md pattern

### Example 5: Body with Fragments

**Input:**
```go
resource := Resource{
    Kind: "Rule",
    Metadata: Metadata{
        ID: "errorHandling",
        Name: "Error Handling",
    },
    Spec: RuleSpec{
        Name: "Handle All Errors",
        Enforcement: "must",
        Body: "{{intro}}\n\n{{examples}}",
        Fragments: map[string]string{
            "intro": "Always handle errors explicitly.",
            "examples": "Use if err != nil checks in Go.",
        },
    },
}

compiler := &KiroCompiler{}
results, err := compiler.Compile(resource)
```

**Expected Output:**
```go
[]CompilationResult{
    {
        Path: "errorHandling.md",
        Content: `---
rule:
  id: errorHandling
  name: Handle All Errors
  enforcement: must
---

# Handle All Errors (MUST)

Always handle errors explicitly.

Use if err != nil checks in Go.`,
    },
}
```

**Verification:**
- Body fragments resolved using `airesource.ResolveBody()`
- Fragments replaced with actual content
- Standalone rule (no ruleset context)
- Path uses rule ID only

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
