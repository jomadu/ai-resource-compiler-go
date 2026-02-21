# GitHub Copilot Compiler

## Job to be Done
Compile AI Resources to GitHub Copilot's modular format, producing `.instructions.md` and `.prompt.md` files with frontmatter that GitHub Copilot can interpret.

## Activities
- Transform Rule resources to .instructions.md with frontmatter
- Transform Prompt resources to .prompt.md with frontmatter
- Generate frontmatter (applyTo, excludeAgent)
- Generate paths using resource ID
- Handle collection items with {collection-id}_{item-id} naming
- Extract scope patterns for applyTo field
- Output markdown format

## Acceptance Criteria
- [x] Rules output as .instructions.md with frontmatter
- [x] Prompts output as .prompt.md with frontmatter
- [x] One file per resource or collection item
- [x] Frontmatter includes applyTo and excludeAgent
- [x] Paths use resource ID: {id}.instructions.md or {id}.prompt.md
- [x] Collection items use {collection-id}_{item-id} naming
- [x] applyTo extracted from scope.include
- [x] excludeAgent optional (empty array if not specified)
- [x] Multi-line bodies preserve formatting
- [x] Empty bodies are skipped

## Data Structures

### CopilotCompiler
```go
type CopilotCompiler struct{}

func (c *CopilotCompiler) Name() string {
    return "copilot"
}

func (c *CopilotCompiler) Compile(resource *airesource.Resource) ([]compiler.CompilationResult, error)
```

**Methods:**
- `Name()` - Returns "copilot"
- `Compile()` - Transforms resource to Copilot format, returns relative paths and content

### Frontmatter
```yaml
applyTo: []string      # File patterns from scope.include
excludeAgent: []string # Optional agent exclusions
```

**Fields:**
- `applyTo` - File patterns from spec.scope.include (empty array if no scope)
- `excludeAgent` - Optional array of agent names to exclude (empty array by default)

## Algorithm

### Compilation Process

1. Check resource kind
2. Extract body content and metadata
3. Generate frontmatter
4. Format as markdown with frontmatter
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
    
    frontmatter = generate_frontmatter(resource)
    content = format_with_frontmatter(frontmatter, body)
    path = resource.Metadata.ID + ".prompt.md"
    
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
        frontmatter = generate_frontmatter_for_item(prompt)
        content = format_with_frontmatter(frontmatter, prompt.Body)
        path = item_id + ".prompt.md"
        
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
    
    frontmatter = generate_frontmatter(resource)
    content = format_with_frontmatter(frontmatter, body)
    path = resource.Metadata.ID + ".instructions.md"
    
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
        frontmatter = generate_frontmatter_for_item(rule)
        content = format_with_frontmatter(frontmatter, rule.Body)
        path = item_id + ".instructions.md"
        
        results.append(CompilationResult{
            Path: path,
            Content: content,
        })
    
    return results
```

### Frontmatter Generation

```
function generate_frontmatter(resource):
    applyTo = []
    if resource.Spec.Scope and resource.Spec.Scope.Include:
        applyTo = resource.Spec.Scope.Include
    
    return {
        applyTo: applyTo,
        excludeAgent: [],
    }

function generate_frontmatter_for_item(item):
    applyTo = []
    if item.Scope and item.Scope.Include:
        applyTo = item.Scope.Include
    
    return {
        applyTo: applyTo,
        excludeAgent: [],
    }
```

### Frontmatter Formatting

```
function format_with_frontmatter(frontmatter, body):
    yaml = serialize_yaml(frontmatter)
    return "---\n" + yaml + "---\n\n" + body
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Empty body | Return empty slice |
| Empty collection | Return empty slice |
| Multi-line body | Preserve all newlines |
| No scope | Empty applyTo array in frontmatter |
| Fragments | Already resolved before compilation |
| Special characters in body | No escaping needed (markdown) |
| Special characters in frontmatter | YAML escaping applied |

## Dependencies

- `compiler-architecture.md` - TargetCompiler interface
- `target-formats.md` - Copilot format specification
- `ai-resource-core-go` - Resource types

## Implementation Mapping

**Source files:**
- `pkg/targets/copilot/compiler.go` - CopilotCompiler implementation
- `pkg/targets/copilot/format.go` - Formatting utilities

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
spec:
  body: "Deploy the application to production"
```

**Expected Output:**
```go
len(results) == 1
results[0].Path == "deploy.prompt.md"
string(results[0].Content) == `---
applyTo: []
excludeAgent: []
---

Deploy the application to production`
```

**Verification:**
- Returns slice with one CompilationResult
- Path is `deploy.prompt.md`
- Frontmatter includes empty applyTo and excludeAgent
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
results[0].Path == "review.prompt.md"
string(results[0].Content) == `---
applyTo: ["src/**/*.ts", "lib/**/*.ts"]
excludeAgent: []
---

Review this code for issues`
```

**Verification:**
- applyTo populated from scope.include
- excludeAgent is empty array

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
      scope:
        include: ["**/*.yml"]
      body: "Run test suite"
```

**Expected Output:**
```go
len(results) == 2

results[0].Path == "ci-workflows_deploy.prompt.md"
string(results[0].Content) == `---
applyTo: []
excludeAgent: []
---

Deploy the application`

results[1].Path == "ci-workflows_test.prompt.md"
string(results[1].Content) == `---
applyTo: ["**/*.yml"]
excludeAgent: []
---

Run test suite`
```

**Verification:**
- Two separate .prompt.md files
- Paths use {collection-id}_{item-id}.prompt.md naming
- Each has own frontmatter with per-item scope

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
results[0].Path == "api-standards.instructions.md"
string(results[0].Content) == `---
applyTo: ["src/**/*.ts"]
excludeAgent: []
---

Follow RESTful API design principles`
```

**Verification:**
- Path is `api-standards.instructions.md`
- applyTo from scope.include
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

results[0].Path == "backend_api.instructions.md"
string(results[0].Content) == `---
applyTo: ["src/api/**/*.ts"]
excludeAgent: []
---

Follow REST principles`

results[1].Path == "backend_security.instructions.md"
string(results[1].Content) == `---
applyTo: []
excludeAgent: []
---

Validate all inputs`
```

**Verification:**
- Two separate .instructions.md files
- Paths use {collection-id}_{item-id}.instructions.md naming
- Each rule has own frontmatter
- applyTo per-item from item scope

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
results[0].Path == "multi.prompt.md"
string(results[0].Content) == `---
applyTo: []
excludeAgent: []
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

- GitHub Copilot uses markdown format with YAML frontmatter
- Different file extensions for rules (.instructions.md) vs prompts (.prompt.md)
- Modular approach: one file per resource or collection item
- Frontmatter fields: applyTo, excludeAgent
- applyTo extracted from scope.include (empty array if no scope)
- excludeAgent defaults to empty array (can be customized in future)
- Collection items use {collection-id}_{item-id} naming with appropriate extension
- Compiler returns relative paths (e.g., "deploy.prompt.md", "api-standards.instructions.md")
- Users prepend `.github/prompts/` or `.github/instructions/` when writing files
- Fragments must be resolved before compilation
- Empty bodies result in empty slice
- YAML escaping applied to frontmatter values
- Body content is markdown (no escaping needed)

## Known Issues

None.

## Areas for Improvement

- Could support custom excludeAgent values from resource metadata
- Could add validation that applyTo patterns are valid globs
- Could support additional Copilot-specific frontmatter fields as they're added
- Could add option to customize frontmatter field names
