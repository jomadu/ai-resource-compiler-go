# Cursor Compiler

## Job to be Done
Compile AI Resources to Cursor's modular MDC format, producing `.mdc` files with frontmatter that Cursor IDE can interpret as coding rules and commands.

## Activities
- Transform Prompt resources to .mdc with frontmatter
- Transform Rule resources to .mdc with frontmatter
- Generate frontmatter (description, globs, alwaysApply)
- Generate paths using resource ID
- Handle collection items with {collection-id}_{item-id} naming
- Extract scope patterns for globs field
- Set alwaysApply based on resource kind

## Acceptance Criteria
- [ ] Output is MDC format with YAML frontmatter
- [ ] One file per resource or collection item
- [ ] Frontmatter includes description, globs, alwaysApply
- [ ] Paths use resource ID: {id}.mdc
- [ ] Collection items use {collection-id}_{item-id}.mdc naming
- [ ] Prompts have alwaysApply: true
- [ ] Rules have alwaysApply: false by default
- [ ] Globs extracted from scope.include
- [ ] Multi-line bodies preserve formatting
- [ ] Empty bodies are skipped

## Target Format

**File Extension:** `.mdc`  
**Format:** MDC with YAML frontmatter  
**Installation Locations:**
- Rules: `.cursor/rules/`
- Prompts: `.cursor/commands/`

**Naming Conventions:**
- Single resources: `{id}.mdc`
  - Example: `id: api-standards` → `api-standards.mdc`
- Collection items: `{collection-id}_{item-id}.mdc`
  - Example: Ruleset `id: backend` with rule `id: api` → `backend_api.mdc`

## Data Structures

### CursorCompiler
```go
type CursorCompiler struct{}

func (c *CursorCompiler) Name() string {
    return "cursor"
}

func (c *CursorCompiler) Compile(resource *airesource.Resource) ([]compiler.CompilationResult, error)
```

**Methods:**
- `Name()` - Returns "cursor"
- `Compile()` - Transforms resource to Cursor MDC format, returns relative paths and content

### Frontmatter
```yaml
description: string  # Resource name or description
globs: []string      # File patterns from scope.include
alwaysApply: bool    # true for prompts, false for rules
```

**Fields:**
- `description` - Resource name (metadata.name) or description
- `globs` - File patterns from spec.scope.include (empty array if no scope)
- `alwaysApply` - true for Prompt/Promptset items, false for Rule/Ruleset items

## Algorithm

### Compilation Process

1. Check resource kind
2. Extract body content and metadata
3. Generate frontmatter
4. Format as MDC with frontmatter
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
    
    frontmatter = generate_frontmatter(resource, alwaysApply: true)
    content = format_mdc(frontmatter, body)
    path = resource.Metadata.ID + ".mdc"
    
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
        frontmatter = generate_frontmatter_for_item(prompt, alwaysApply: true)
        content = format_mdc(frontmatter, prompt.Body)
        path = item_id + ".mdc"
        
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
    
    frontmatter = generate_frontmatter(resource, alwaysApply: false)
    content = format_mdc(frontmatter, body)
    path = resource.Metadata.ID + ".mdc"
    
    return [CompilationResult{
        Path: path,
        Content: content,
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
        frontmatter = generate_frontmatter_for_item(rule, alwaysApply: false)
        content = format_mdc(frontmatter, rule.Body)
        path = item_id + ".mdc"
        
        results.append(CompilationResult{
            Path: path,
            Content: content,
        })
    
    return results
```

### Frontmatter Generation

```
function generate_frontmatter(resource, alwaysApply):
    description = resource.Metadata.Name
    if empty(description):
        description = resource.Metadata.ID
    
    globs = []
    if resource.Spec.Scope and resource.Spec.Scope.Include:
        globs = resource.Spec.Scope.Include
    
    return {
        description: description,
        globs: globs,
        alwaysApply: alwaysApply,
    }

function generate_frontmatter_for_item(item, alwaysApply):
    description = item.Name
    if empty(description):
        description = ""
    
    globs = []
    if item.Scope and item.Scope.Include:
        globs = item.Scope.Include
    
    return {
        description: description,
        globs: globs,
        alwaysApply: alwaysApply,
    }
```

### MDC Formatting

```
function format_mdc(frontmatter, body):
    yaml = serialize_yaml(frontmatter)
    return "---\n" + yaml + "---\n\n" + body
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Empty body | Return empty slice |
| Empty collection | Return empty slice |
| Multi-line body | Preserve all newlines |
| Missing name | Use ID as description |
| Missing description | Use empty string for collection items |
| No scope | Empty globs array in frontmatter |
| Fragments | Already resolved before compilation |
| Special characters in body | No escaping needed (markdown) |
| Special characters in frontmatter | YAML escaping applied |

## Dependencies

- `compiler-architecture.md` - TargetCompiler interface
- `ai-resource-core-go` - Resource types

## Implementation Mapping

**Source files:**
- `pkg/targets/cursor/compiler.go` - CursorCompiler implementation
- `pkg/targets/cursor/format.go` - Formatting utilities

**Related specs:**
- `compiler-architecture.md` - Compiler interface

## Examples

### Example 1: Simple Prompt

**Input:**
```yaml
apiVersion: ai-resource/draft
kind: Prompt
metadata:
  id: deploy
  name: Deploy Application
spec:
  body: "Deploy the application to production"
```

**Expected Output:**
```go
len(results) == 1
results[0].Path == "deploy.mdc"
string(results[0].Content) == `---
description: Deploy Application
globs: []
alwaysApply: true
---

Deploy the application to production`
```

**Verification:**
- Returns slice with one CompilationResult
- Path is `deploy.mdc`
- Frontmatter includes description, empty globs, alwaysApply true
- Body preserved after frontmatter

### Example 2: Prompt with Scope

**Input:**
```yaml
apiVersion: ai-resource/draft
kind: Prompt
metadata:
  id: review
  name: Code Review
spec:
  scope:
    include: ["src/**/*.ts", "lib/**/*.ts"]
  body: "Review this code for issues"
```

**Expected Output:**
```go
len(results) == 1
results[0].Path == "review.mdc"
string(results[0].Content) == `---
description: Code Review
globs: ["src/**/*.ts", "lib/**/*.ts"]
alwaysApply: true
---

Review this code for issues`
```

**Verification:**
- Globs populated from scope.include
- alwaysApply is true for prompts

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
      body: "Deploy the application"
    test:
      name: Test
      body: "Run test suite"
```

**Expected Output:**
```go
len(results) == 2

results[0].Path == "ci-workflows_deploy.mdc"
string(results[0].Content) == `---
description: Deploy
globs: []
alwaysApply: true
---

Deploy the application`

results[1].Path == "ci-workflows_test.mdc"
string(results[1].Content) == `---
description: Test
globs: []
alwaysApply: true
---

Run test suite`
```

**Verification:**
- Two separate .mdc files
- Paths use {collection-id}_{item-id}.mdc naming
- Each has own frontmatter

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
results[0].Path == "api-standards.mdc"
string(results[0].Content) == `---
description: API Design Standards
globs: ["src/**/*.ts"]
alwaysApply: false
---

Follow RESTful API design principles`
```

**Verification:**
- alwaysApply is false for rules
- Globs from scope.include
- Enforcement not in frontmatter

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

results[0].Path == "backend_api.mdc"
string(results[0].Content) == `---
description: API Standards
globs: ["src/api/**/*.ts"]
alwaysApply: false
---

Follow REST principles`

results[1].Path == "backend_security.mdc"
string(results[1].Content) == `---
description: Security
globs: []
alwaysApply: false
---

Validate all inputs`
```

**Verification:**
- Two separate .mdc files
- Paths use {collection-id}_{item-id}.mdc naming
- Each rule has own frontmatter
- Globs per-item from item scope

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
results[0].Path == "multi.mdc"
string(results[0].Content) == `---
description: Multi-line
globs: []
alwaysApply: true
---

First line
Second line
Third line`
```

**Verification:**
- All newlines preserved in body
- Frontmatter properly separated

## Notes

- Cursor uses MDC format with YAML frontmatter
- Modular approach: one .mdc file per resource or collection item
- Frontmatter fields: description, globs, alwaysApply
- Prompts have alwaysApply: true, rules have alwaysApply: false
- Globs extracted from scope.include (empty array if no scope)
- Collection items use {collection-id}_{item-id}.mdc naming
- Compiler returns relative paths (e.g., "deploy.mdc")
- Users prepend `.cursor/rules/` or `.cursor/commands/` when writing files
- Fragments must be resolved before compilation
- Empty bodies result in empty slice
- YAML escaping applied to frontmatter values
- Body content is markdown (no escaping needed)

## Known Issues

None.

## Areas for Improvement

- Could support custom frontmatter fields
- Could add validation that globs are valid patterns
- Could support priority-based ordering within collections
- Could add option to customize alwaysApply per resource
- Could support additional Cursor-specific frontmatter fields as they're added
