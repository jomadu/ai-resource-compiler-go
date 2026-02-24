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
func (m *MarkdownCompiler) Compile(resource Resource) ([]CompilationResult, error)
func (m *MarkdownCompiler) SupportedVersions() []string
```

**Methods:**
- `Name()` - Returns "markdown"
- `Compile()` - Transforms resource into markdown format
- `SupportedVersions()` - Returns `["ai-resource/v1"]`

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
| Unsupported apiVersion | Return error "unsupported apiVersion: {version} for markdown" |

## Dependencies

- Resource model from ai-resource-core-go
- Metadata block generation from metadata-block.md spec
- TargetCompiler interface from compiler-architecture.md spec

## Implementation Mapping

**Source files:**
- `pkg/targets/markdown.go` - MarkdownCompiler implementation

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
