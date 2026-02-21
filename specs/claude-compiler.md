# Claude Compiler

## Job to be Done
Compile AI Resources to Claude Code's modular format, producing `.md` files for rules and directory-based `SKILL.md` files for prompts that Claude Code can interpret.

## Activities
- Transform Prompt resources to {id}/SKILL.md directories
- Transform Rule resources to .md files
- Generate optional frontmatter for prompts (name, description)
- Generate paths using resource ID
- Handle collection items with {collection-id}_{item-id} naming
- Output markdown format

## Acceptance Criteria
- [x] Rules output as plain markdown .md files
- [x] Prompts output as directories with SKILL.md file
- [x] One file per resource or collection item
- [x] Prompt frontmatter includes name and description (optional)
- [x] Rule files are plain markdown (no frontmatter)
- [x] Paths use resource ID: {id}.md or {id}/SKILL.md
- [x] Collection items use {collection-id}_{item-id} naming
- [x] Multi-line bodies preserve formatting
- [x] Empty bodies are skipped

## Data Structures

### ClaudeCompiler
```go
type ClaudeCompiler struct{}

func (c *ClaudeCompiler) Name() string {
    return "claude"
}

func (c *ClaudeCompiler) Compile(resource *airesource.Resource) ([]compiler.CompilationResult, error)
```

**Methods:**
- `Name()` - Returns "claude"
- `Compile()` - Transforms resource to Claude format, returns relative paths and content

### Prompt Frontmatter (Optional)
```yaml
name: string         # Resource name
description: string  # Resource description
```

**Fields:**
- `name` - Resource name from metadata.name (optional)
- `description` - Resource description from metadata.description (optional)

## Algorithm

### Compilation Process

1. Check resource kind
2. Extract body content and metadata
3. Generate frontmatter (prompts only, if name/description present)
4. Format as markdown
5. Return CompilationResult with relative path and content

**Pseudocode:**
```
function Compile(resource):
    switch resource.Kind:
        case Prompt:
            return compile_prompt(resource)
        case Promptset:
            return compile_promptset(resource)
        case Rule:
            return compile_rule(resource)
        case Ruleset:
            return compile_ruleset(resource)
        default:
            return error("unsupported kind: {resource.Kind}")
```

### Prompt Compilation

```
function compile_prompt(resource):
    body = resource.Spec.Body
    if empty(body):
        return []
    
    content = body
    if resource.Metadata.Name or resource.Metadata.Description:
        frontmatter = generate_frontmatter(resource)
        content = format_with_frontmatter(frontmatter, body)
    
    path = resource.Metadata.ID + "/SKILL.md"
    
    return [CompilationResult{
        Path: path,
        Content: content,
    }]
```

### Promptset Compilation

```
function compile_promptset(resource):
    results = []
    collection_id = resource.Metadata.ID
    
    for key, prompt in resource.Spec.Prompts:
        if empty(prompt.Body):
            continue
        
        item_id = collection_id + "_" + key
        content = prompt.Body
        
        if prompt.Name or prompt.Description:
            frontmatter = generate_frontmatter_for_item(prompt)
            content = format_with_frontmatter(frontmatter, prompt.Body)
        
        path = item_id + "/SKILL.md"
        
        results.append(CompilationResult{
            Path: path,
            Content: content,
        })
    
    return results
```

### Rule Compilation

```
function compile_rule(resource):
    body = resource.Spec.Body
    if empty(body):
        return []
    
    path = resource.Metadata.ID + ".md"
    
    return [CompilationResult{
        Path: path,
        Content: body,
    }]
```

### Ruleset Compilation

```
function compile_ruleset(resource):
    results = []
    collection_id = resource.Metadata.ID
    
    for key, rule in resource.Spec.Rules:
        if empty(rule.Body):
            continue
        
        item_id = collection_id + "_" + key
        path = item_id + ".md"
        
        results.append(CompilationResult{
            Path: path,
            Content: rule.Body,
        })
    
    return results
```

### Frontmatter Generation

```
function generate_frontmatter(resource):
    frontmatter = {}
    
    if resource.Metadata.Name:
        frontmatter["name"] = resource.Metadata.Name
    
    if resource.Metadata.Description:
        frontmatter["description"] = resource.Metadata.Description
    
    return frontmatter

function generate_frontmatter_for_item(item):
    frontmatter = {}
    
    if item.Name:
        frontmatter["name"] = item.Name
    
    if item.Description:
        frontmatter["description"] = item.Description
    
    return frontmatter
```

### Frontmatter Formatting

```
function format_with_frontmatter(frontmatter, body):
    if empty(frontmatter):
        return body
    
    yaml = serialize_yaml(frontmatter)
    return "---\n" + yaml + "---\n\n" + body
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Empty body | Return empty slice |
| Empty collection | Return empty slice |
| Multi-line body | Preserve all newlines |
| Missing name and description | No frontmatter for prompts |
| Fragments | Already resolved before compilation |
| Special characters in body | No escaping needed (markdown) |
| Special characters in frontmatter | YAML escaping applied |

## Dependencies

- `compiler-architecture.md` - TargetCompiler interface
- `target-formats.md` - Claude format specification
- `ai-resource-core-go` - Resource types

## Implementation Mapping

**Source files:**
- `pkg/targets/claude/compiler.go` - ClaudeCompiler implementation
- `pkg/targets/claude/format.go` - Formatting utilities

**Related specs:**
- `compiler-architecture.md` - Compiler interface
- `target-formats.md` - Format specification

## Examples

### Example 1: Simple Prompt

**Input:**
```yaml
apiVersion: ai-resource/draft
kind: Prompt
metadata:
  id: deploy
  name: Deploy Application
  description: Deploy to production
spec:
  body: "Deploy the application to production"
```

**Expected Output:**
```go
len(results) == 1
results[0].Path == "deploy/SKILL.md"
string(results[0].Content) == `---
name: Deploy Application
description: Deploy to production
---

Deploy the application to production`
```

**Verification:**
- Returns slice with one CompilationResult
- Path is `deploy/SKILL.md` (directory structure)
- Frontmatter includes name and description
- Body preserved after frontmatter

### Example 2: Prompt without Metadata

**Input:**
```yaml
apiVersion: ai-resource/draft
kind: Prompt
metadata:
  id: simple
spec:
  body: "Simple prompt body"
```

**Expected Output:**
```go
len(results) == 1
results[0].Path == "simple/SKILL.md"
string(results[0].Content) == "Simple prompt body"
```

**Verification:**
- No frontmatter when name/description absent
- Plain markdown body

### Example 3: Promptset

**Input:**
```yaml
apiVersion: ai-resource/draft
kind: Promptset
metadata:
  id: ci-workflows
spec:
  prompts:
    deploy:
      name: Deploy
      description: Deploy the application
      body: "Deploy the application"
    test:
      name: Test
      body: "Run test suite"
```

**Expected Output:**
```go
len(results) == 2

results[0].Path == "ci-workflows_deploy/SKILL.md"
string(results[0].Content) == `---
name: Deploy
description: Deploy the application
---

Deploy the application`

results[1].Path == "ci-workflows_test/SKILL.md"
string(results[1].Content) == `---
name: Test
---

Run test suite`
```

**Verification:**
- Two separate directories with SKILL.md files
- Paths use {collection-id}_{item-id}/SKILL.md naming
- Each has own frontmatter (description optional)

### Example 4: Rule

**Input:**
```yaml
apiVersion: ai-resource/draft
kind: Rule
metadata:
  id: api-standards
  name: API Design Standards
spec:
  enforcement: must
  scope:
    include: ["src/**/*.ts"]
  body: "Follow RESTful API design principles"
```

**Expected Output:**
```go
len(results) == 1
results[0].Path == "api-standards.md"
string(results[0].Content) == "Follow RESTful API design principles"
```

**Verification:**
- Path is `api-standards.md` (file, not directory)
- Plain markdown body (no frontmatter)
- Scope and enforcement not in output

### Example 5: Ruleset

**Input:**
```yaml
apiVersion: ai-resource/draft
kind: Ruleset
metadata:
  id: backend
spec:
  rules:
    api:
      name: API Standards
      scope:
        include: ["src/api/**/*.ts"]
      body: "Follow REST principles"
    security:
      name: Security
      body: "Validate all inputs"
```

**Expected Output:**
```go
len(results) == 2

results[0].Path == "backend_api.md"
string(results[0].Content) == "Follow REST principles"

results[1].Path == "backend_security.md"
string(results[1].Content) == "Validate all inputs"
```

**Verification:**
- Two separate .md files
- Paths use {collection-id}_{item-id}.md naming
- Plain markdown content (no frontmatter)

### Example 6: Multi-line Body

**Input:**
```yaml
apiVersion: ai-resource/draft
kind: Prompt
metadata:
  id: multi
  name: Multi-line
spec:
  body: |
    First line
    Second line
    Third line
```

**Expected Output:**
```go
len(results) == 1
results[0].Path == "multi/SKILL.md"
string(results[0].Content) == `---
name: Multi-line
---

First line
Second line
Third line`
```

**Verification:**
- All newlines preserved in body
- Frontmatter properly separated

### Example 7: Empty Body

**Input:**
```yaml
apiVersion: ai-resource/draft
kind: Prompt
metadata:
  id: empty
spec:
  body: ""
```

**Expected Output:**
```go
len(results) == 0
```

**Verification:**
- Empty slice returned
- No file created

## Notes

- Claude uses directory structure for prompts: `{id}/SKILL.md`
- Rules are plain .md files: `{id}.md`
- Prompt frontmatter is optional (only if name/description present)
- Rule files have no frontmatter
- Modular approach: one file per resource or collection item
- Collection items use {collection-id}_{item-id} naming
- Compiler returns relative paths (e.g., "deploy/SKILL.md", "api-standards.md")
- Users prepend `.claude/skills/` or `.claude/rules/` when writing files
- Fragments must be resolved before compilation
- Empty bodies result in empty slice
- YAML escaping applied to frontmatter values
- Body content is markdown (no escaping needed)
- Metadata fields (scope, enforcement) not included in output

## Known Issues

None.

## Areas for Improvement

- Could support additional frontmatter fields as Claude adds them
- Could add option to always include frontmatter (even when empty)
- Could support custom templates for SKILL.md formatting
- Could add validation that directory names are filesystem-safe
