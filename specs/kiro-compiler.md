# Kiro Compiler

## Job to be Done
Compile AI Resources to Kiro CLI's plain markdown format, producing `.md` files that Kiro can interpret as steering rules and prompts.

## Activities
- Transform Prompt resources to .md files
- Transform Rule resources to .md files
- Generate paths using resource ID
- Handle collection items with {collection-id}_{item-id} naming
- Output plain markdown (no frontmatter)

## Acceptance Criteria
- [ ] Output is plain markdown format
- [ ] One file per resource or collection item
- [ ] No frontmatter
- [ ] Paths use resource ID: {id}.md
- [ ] Collection items use {collection-id}_{item-id}.md naming
- [ ] Multi-line bodies preserve formatting
- [ ] Empty bodies are skipped

## Target Format

**File Extension:** `.md`  
**Format:** Plain markdown (no frontmatter)  
**Installation Locations:**
- Rules: `.kiro/steering/`
- Prompts: `.kiro/prompts/`

**Naming Conventions:**
- Single resources: `{id}.md`
  - Example: `id: api-standards` → `api-standards.md`
- Collection items: `{collection-id}_{item-id}.md`
  - Example: Ruleset `id: backend` with rule `id: api` → `backend_api.md`

## Data Structures

### KiroCompiler
```go
type KiroCompiler struct{}

func (c *KiroCompiler) Name() string {
    return "kiro"
}

func (c *KiroCompiler) Compile(resource *airesource.Resource) ([]compiler.CompilationResult, error)
```

**Methods:**
- `Name()` - Returns "kiro"
- `Compile()` - Transforms resource to plain markdown, returns relative paths and content

## Algorithm

### Compilation Process

1. Check resource kind
2. Extract body content
3. Return CompilationResult with relative path and content

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
    
    path = resource.Metadata.ID + ".md"
    
    return [CompilationResult{
        Path: path,
        Content: body,
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
        path = item_id + ".md"
        
        results.append(CompilationResult{
            Path: path,
            Content: prompt.Body,
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

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Empty body | Return empty slice |
| Empty collection | Return empty slice |
| Multi-line body | Preserve all newlines |
| Fragments | Already resolved before compilation |
| Special characters in body | No escaping needed (markdown) |

## Dependencies

- `compiler-architecture.md` - TargetCompiler interface
- `ai-resource-core-go` - Resource types

## Implementation Mapping

**Source files:**
- `pkg/targets/kiro/compiler.go` - KiroCompiler implementation

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
results[0].Path == "deploy.md"
string(results[0].Content) == "Deploy the application to production"
```

**Verification:**
- Returns slice with one CompilationResult
- Path is `deploy.md`
- Content is plain markdown body

### Example 2: Promptset

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

results[0].Path == "ci-workflows_deploy.md"
string(results[0].Content) == "Deploy the application"

results[1].Path == "ci-workflows_test.md"
string(results[1].Content) == "Run test suite"
```

**Verification:**
- Two separate .md files
- Paths use {collection-id}_{item-id}.md naming
- Plain markdown content

### Example 3: Rule

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
- Path is `api-standards.md`
- Content is plain markdown body
- Scope and enforcement not in output

### Example 4: Ruleset

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
- Plain markdown content

### Example 5: Multi-line Body

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
results[0].Path == "multi.md"
string(results[0].Content) == "First line\nSecond line\nThird line"
```

**Verification:**
- All newlines preserved in body

### Example 6: Empty Body

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

- Kiro uses plain markdown format (no frontmatter)
- Simplest target compiler
- Modular approach: one .md file per resource or collection item
- Collection items use {collection-id}_{item-id}.md naming
- Compiler returns relative paths (e.g., "deploy.md")
- Users prepend `.kiro/steering/` or `.kiro/prompts/` when writing files
- Fragments must be resolved before compilation
- Empty bodies result in empty slice
- No escaping needed (markdown)
- Metadata fields (name, description, scope, enforcement) not included in output

## Known Issues

None.

## Areas for Improvement

- Could support optional frontmatter for metadata
- Could add option to include resource metadata as markdown comments
- Could support custom templates for body formatting
