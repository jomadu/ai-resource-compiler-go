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
  name: string
  description: string (optional)
  rules: []string
rule:
  id: string
  name: string
  description: string (optional)
  enforcement: string
  scope: object (optional)
    files: []string
---
```

**Fields:**
- `ruleset.id` - Unique identifier for the ruleset collection
- `ruleset.name` - Human-readable ruleset name
- `ruleset.description` - Optional description of ruleset purpose
- `ruleset.rules` - List of rule IDs in this ruleset
- `rule.id` - Unique identifier for this specific rule
- `rule.name` - Human-readable rule name
- `rule.description` - Optional description of rule purpose
- `rule.enforcement` - Enforcement level (must, should, may)
- `rule.scope.files` - Optional file patterns where rule applies

### Enforcement Header
```
# {Rule Name} ({ENFORCEMENT})
```

**Format:**
- Rule name from `rule.name` field
- Enforcement level uppercased (MUST, SHOULD, MAY)
- Follows metadata block, precedes rule body

## Algorithm

1. Check resource type (rule vs prompt)
2. If prompt: return body content only
3. If rule: construct metadata block
4. Add ruleset section with required fields
5. Add optional ruleset.description if present
6. Add rule section with required fields
7. Add optional rule.description and rule.scope if present
8. Generate enforcement header from rule.name and rule.enforcement
9. Concatenate: metadata block + enforcement header + rule body

**Pseudocode:**
```
function EmbedMetadata(resource):
    if resource.type == "prompt":
        return resource.body
    
    metadata = BuildYAMLBlock(resource.ruleset, resource.rule)
    header = FormatEnforcementHeader(resource.rule.name, resource.rule.enforcement)
    
    return metadata + "\n" + header + "\n\n" + resource.body
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

- Resource model from ai-resource-core-go (provides ruleset, rule, prompt structures)
- YAML serialization library for metadata block generation

## Implementation Mapping

**Source files:**
- `internal/format/metadata.go` - Metadata block generation logic
- `pkg/compiler/compiler.go` - Integration point for metadata embedding

**Related specs:**
- `compiler-architecture.md` - Defines where metadata embedding fits in compilation pipeline
- All target compiler specs - Each target uses metadata blocks for rules

## Examples

### Example 1: Full Metadata Block

**Input:**
```go
Resource{
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

### Example 2: Minimal Metadata Block

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

Use descriptive names.
```

**Verification:**
- Optional fields (description, scope) omitted
- Enforcement header shows "SHOULD"
- Minimal valid metadata block

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

The enforcement header provides immediate visual feedback about rule importance without requiring tools to parse YAML.

Prompts intentionally exclude metadata because they represent reusable instructions rather than constraints with enforcement levels.

## Known Issues

None - this is a new specification.

## Areas for Improvement

- Consider adding metadata version field for future schema evolution
- Evaluate whether prompts should include minimal metadata (promptset/prompt IDs)
- Explore compact metadata format for tools with strict size limits
