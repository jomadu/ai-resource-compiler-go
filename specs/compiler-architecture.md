# Compiler Architecture

## Job to be Done
Provide a clear, extensible architecture for compiling AI Resources to multiple target formats while maintaining separation of concerns and enabling easy addition of new targets.

## Activities
- Define compilation pipeline stages
- Establish target compiler interface
- Separate resource loading from compilation
- Enable multi-target compilation in single pass
- Provide extension points for new targets

## Acceptance Criteria
- [ ] Each target compiler is independent and self-contained
- [ ] Adding new target requires only implementing target interface
- [ ] Resource loading uses ai-resource-core-go (no duplication)
- [ ] Compilation errors are target-specific and actionable
- [ ] Multiple targets can be compiled from same resource
- [ ] Fragment resolution happens before target compilation
- [ ] Public API is minimal and stable

## Data Structures

### Compiler
```go
type Compiler struct {
    targets map[string]TargetCompiler
}
```

**Fields:**
- `targets` - Registry of available target compilers

### TargetCompiler
```go
type TargetCompiler interface {
    Name() string
    CompileRule(rule *airesource.Rule, namespace string) (CompilationResult, error)
    CompileRuleset(ruleset *airesource.Ruleset, namespace string) ([]CompilationResult, error)
    CompilePrompt(prompt *airesource.Prompt) (CompilationResult, error)
    CompilePromptset(promptset *airesource.Promptset) ([]CompilationResult, error)
}
```

**Methods:**
- `Name()` - Target identifier (e.g., "kiro", "cursor")
- `CompileRule()` - Transform single rule to target format, requires namespace for metadata block
- `CompileRuleset()` - Transform ruleset to target format, returns one result per rule, requires namespace
- `CompilePrompt()` - Transform single prompt to target format, no namespace needed
- `CompilePromptset()` - Transform promptset to target format, returns one result per prompt, no namespace needed

**Design Notes:**
- Separate methods for rules vs prompts reflect different compilation needs
- Namespace required for rules (metadata block), not used for prompts
- Type-safe: compiler signature enforces namespace where needed
- Single resources return one result, collections return multiple results

### CompilationResult
```go
type CompilationResult struct {
    Path    string
    Content []byte
}
```

**Fields:**
- `Path` - Relative path for output (e.g., "api-standards.md", "deploy.mdc", "backend_api.instructions.md")
- `Content` - Compiled content as bytes

## Naming Conventions

**Single Resources:**
- Use resource ID as filename: `{id}.{ext}`
- Example: `id: api-standards` → `api-standards.md` (Kiro/Claude), `api-standards.mdc` (Cursor), `api-standards.instructions.md` (Copilot)

**Collection Items:**
- Combine collection ID and item ID: `{collection-id}_{item-id}.{ext}`
- Example: Ruleset `id: backend` with rule `id: api` → `backend_api.md` (Kiro/Claude), `backend_api.mdc` (Cursor), `backend_api.instructions.md` (Copilot)

**Extensions by Target:**
- Kiro: `.md`
- Cursor: `.mdc`
- Claude: `.md` (rules), directory with `SKILL.md` (prompts)
- Copilot: `.instructions.md` (rules), `.prompt.md` (prompts)

## Recommended Installation Locations

| Target   | Rules                  | Prompts               |
|----------|------------------------|-----------------------|
| Kiro     | `.kiro/steering/`      | `.kiro/prompts/`      |
| Cursor   | `.cursor/rules/`       | `.cursor/commands/`   |
| Claude   | `.claude/rules/`       | `.claude/skills/`     |
| Copilot  | `.github/instructions/`| `.github/prompts/`    |

**Note:** The compiler returns relative paths (e.g., `api-standards.md`). Users prepend target-specific directories when writing files.

### CompileOptions
```go
type CompileOptions struct {
    Targets          []string
    Namespace        string  // Required for rules/rulesets, ignored for prompts
    ResolveFragments bool
}
```

**Fields:**
- `Targets` - List of target names to compile to
- `Namespace` - Package identifier for metadata block (e.g., "registry/package@1.0.0")
- `ResolveFragments` - Whether to resolve fragments before compilation (default: true)

**Note:** 
- Namespace is required when compiling rules/rulesets (used in metadata block)
- Namespace is ignored when compiling prompts/promptsets
- Output directory management is the caller's responsibility
- The compiler returns relative paths; users decide where to write files

### CompileResult
```go
type CompileResult struct {
    Target  string
    Results []CompilationResult
    Error   error
}
```

**Fields:**
- `Target` - Target name that was compiled
- `Results` - One or more compilation results (path + content pairs)
- `Error` - Compilation error if any

## Algorithm

### Compilation Pipeline

1. Load resource using ai-resource-core-go
2. Validate resource (schema + semantic)
3. Resolve fragments if enabled
4. Determine resource kind
5. For each target:
   - Get target compiler
   - Call appropriate compile method based on resource kind:
     - Rule → CompileRule(rule, namespace)
     - Ruleset → CompileRuleset(ruleset, namespace)
     - Prompt → CompilePrompt(prompt)
     - Promptset → CompilePromptset(promptset)
   - Collect results (one or more path/content pairs)
6. Return results for all targets

**Pseudocode:**
```
function Compile(resource, options):
    // Validate
    if not resource.IsValid():
        return error
    
    // Resolve fragments
    if options.ResolveFragments:
        resource = resolve_fragments(resource)
    
    // Validate namespace for rules
    if (resource.Kind == Rule or resource.Kind == Ruleset) and options.Namespace == "":
        return error("namespace required for rule compilation")
    
    // Compile to each target
    results = []
    for target_name in options.Targets:
        compiler = get_target_compiler(target_name)
        if not compiler:
            results.append(CompileResult{
                Target: target_name,
                Error: error("unknown target: {target_name}"),
            })
            continue
        
        // Call appropriate compile method based on resource kind
        var compilation_results []CompilationResult
        var err error
        
        switch resource.Kind:
            case Rule:
                result, err = compiler.CompileRule(resource, options.Namespace)
                compilation_results = []CompilationResult{result}
            case Ruleset:
                compilation_results, err = compiler.CompileRuleset(resource, options.Namespace)
            case Prompt:
                result, err = compiler.CompilePrompt(resource)
                compilation_results = []CompilationResult{result}
            case Promptset:
                compilation_results, err = compiler.CompilePromptset(resource)
            default:
                err = error("unsupported resource kind: {resource.Kind}")
        
        results.append(CompileResult{
            Target: target_name,
            Results: compilation_results,
            Error: err,
        })
    
    return results
```

### Target Registration

```
function RegisterTarget(compiler TargetCompiler):
    targets[compiler.Name()] = compiler

function GetTarget(name string) TargetCompiler:
    return targets[name]
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Unknown target name | Return CompileResult with error |
| Resource validation failure | Return error, don't attempt compilation |
| One target fails | Continue compiling other targets, collect error |
| Empty targets list | Error: "no targets specified" |
| Fragment resolution fails | Return error, don't attempt compilation |
| Empty body | Skip output (no file created) |
| Missing namespace for rule | Error: "namespace required for rule compilation" |
| Namespace provided for prompt | Ignored (prompts don't use namespace) |
| Single rule/prompt | Returns one result with `{id}.{ext}` path |
| Ruleset/promptset | Returns one result per item with `{collection-id}_{item-id}.{ext}` path |
| Empty collection | No files created |

## Dependencies

- `ai-resource-core-go` - Resource loading and validation
- Target compiler implementations (kiro, cursor, claude, copilot, markdown)

## Implementation Mapping

**Source files:**
- `pkg/compiler/compiler.go` - Main Compiler type and Compile function
- `pkg/compiler/target.go` - TargetCompiler interface
- `pkg/compiler/registry.go` - Target registration
- `pkg/targets/kiro/compiler.go` - Kiro target implementation
- `pkg/targets/cursor/compiler.go` - Cursor target implementation
- `pkg/targets/claude/compiler.go` - Claude target implementation
- `pkg/targets/copilot/compiler.go` - Copilot target implementation
- `pkg/targets/markdown/compiler.go` - Markdown target implementation

**Related specs:**
- `metadata-block.md` - Metadata block structure for rules
- `kiro-compiler.md` - Kiro target implementation
- `cursor-compiler.md` - Cursor target implementation
- `claude-compiler.md` - Claude target implementation
- `copilot-compiler.md` - Copilot target implementation
- `markdown-compiler.md` - Markdown target implementation

## Examples

### Example 1: Single Prompt Compilation

**Input:**
```go
import (
    "github.com/jomadu/ai-resource-compiler-go/pkg/compiler"
)

// Single prompt: id: deploy
c := compiler.New()
results, err := c.Compile(prompt, compiler.Options{
    Targets: []string{"cursor"},
})
```

**Expected Output:**
```go
len(results) == 1
results[0].Target == "cursor"
len(results[0].Results) == 1
results[0].Results[0].Path == "deploy.mdc"
results[0].Results[0].Content == []byte("Deploy the application")
results[0].Error == nil
```

**Verification:**
- Single CompileResult returned
- Path is `deploy.mdc` (resource ID + extension)
- Content includes frontmatter
- No errors

### Example 2: Single Rule Compilation with Namespace

**Input:**
```go
// Single rule: id: api-standards
results, err := c.Compile(rule, compiler.Options{
    Targets: []string{"kiro", "cursor", "claude"},
    Namespace: "myorg/standards@1.0.0",
})
```

**Expected Output:**
```go
len(results) == 3
results[0].Target == "kiro"
results[0].Results[0].Path == "api-standards.md"
results[1].Target == "cursor"
results[1].Results[0].Path == "api-standards.mdc"
results[2].Target == "claude"
results[2].Results[0].Path == "api-standards.md"
```

**Verification:**
- Three CompileResults returned
- Each target has different extension
- All include metadata block with namespace
- All successful

### Example 3: Missing Namespace Error

**Input:**
```go
// Rule without namespace
results, err := c.Compile(rule, compiler.Options{
    Targets: []string{"cursor"},
})
```

**Expected Output:**
```go
err.Error() == "namespace required for rule compilation"
```

**Verification:**
- Error returned before compilation
- No CompileResults generated

### Example 4: Unknown Target Error

**Input:**
```go
results, err := c.Compile(prompt, compiler.Options{
    Targets: []string{"unknown"},
})
```

**Expected Output:**
```go
len(results) == 1
results[0].Target == "unknown"
results[0].Error.Error() == "unknown target: unknown"
```

**Verification:**
- Error indicates unknown target
- Compilation not attempted

### Example 5: Partial Failure

**Input:**
```go
results, err := c.Compile(prompt, compiler.Options{
    Targets: []string{"cursor", "unknown", "kiro"},
})
```

**Expected Output:**
```go
len(results) == 3
results[0].Error == nil  // cursor succeeded
results[1].Error != nil  // unknown failed
results[2].Error == nil  // kiro succeeded
```

**Verification:**
- All targets attempted
- Failures don't stop other compilations
- Errors collected per target

### Example 6: Ruleset Compilation

**Input:**
```go
// Ruleset: id: backend with rules: api, security
results, err := c.Compile(ruleset, compiler.Options{
    Targets: []string{"cursor"},
    Namespace: "myorg/backend@1.0.0",
})
```

**Expected Output:**
```go
len(results) == 1
results[0].Target == "cursor"
len(results[0].Results) == 2
results[0].Results[0].Path == "backend_api.mdc"
results[0].Results[1].Path == "backend_security.mdc"
```

**Verification:**
- Single CompileResult for cursor target
- Two CompilationResults (one per rule)
- Paths use `{collection-id}_{item-id}.mdc` naming
- User decides base directory for writing

## Notes

- The compiler is a pure transformation tool - it does not manage file I/O
- Target compilers return relative paths; callers decide where to write files
- Fragment resolution happens once before target compilation to avoid duplication
- Each target compiler is completely independent and can be developed separately
- The TargetCompiler interface is the extension point for new targets
- Resource validation is handled by ai-resource-core-go, not duplicated here
- Multi-target compilation is efficient - resource processed once, compiled to multiple formats
- Target compilers receive fully validated, resolved resources
- Single resources output one file per target: `{id}.{ext}`
- Collection items output one file per item: `{collection-id}_{item-id}.{ext}`
- Users prepend recommended installation directories when writing files
- See "Recommended Installation Locations" table for target-specific directories

## Known Issues

None.

## Areas for Improvement

- Could add caching for repeated compilations
- Could support streaming compilation for large resource bundles
- Could add compilation middleware/hooks for custom transformations
- Could support target-specific options in CompileOptions
- Could add dry-run mode to preview output without writing files
