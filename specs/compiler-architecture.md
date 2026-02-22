# Compiler Architecture

## Job to be Done
Provide an extensible architecture for compiling AI resources into multiple target-specific formats through a unified interface and pluggable target compilers.

## Activities
1. Define TargetCompiler interface for implementing target-specific compilation
2. Provide Target enum for type-safe target selection
3. Define CompileOptions for configuring compilation behavior
4. Define CompilationResult structure for modular output (path + content)
5. Implement compilation pipeline that orchestrates target compilers
6. Support target compiler registration and discovery

## Acceptance Criteria
- [ ] TargetCompiler interface has single Compile method
- [ ] Target enum includes all supported targets (cursor, kiro, claude, copilot, markdown)
- [ ] CompileOptions accepts list of targets
- [ ] CompilationResult contains path and content fields
- [ ] Compiler.Compile method returns results for all requested targets
- [ ] Each target compiler operates independently
- [ ] Path structure follows {collection-id}_{item-id}.{ext} pattern
- [ ] Claude prompts use {collection-id}_{item-id}/SKILL.md pattern

## Data Structures

### Target
```go
type Target string

const (
    TargetCursor   Target = "cursor"
    TargetKiro     Target = "kiro"
    TargetClaude   Target = "claude"
    TargetCopilot  Target = "copilot"
    TargetMarkdown Target = "markdown"
)
```

**Values:**
- `cursor` - Cursor IDE rules (.mdc) and commands (.md)
- `kiro` - Kiro CLI steering rules and prompts (.md)
- `claude` - Claude Code rules (.md) and skills (SKILL.md)
- `copilot` - GitHub Copilot instructions and prompts (.md)
- `markdown` - Generic markdown output (.md)

### CompileOptions
```go
type CompileOptions struct {
    Targets []Target
}
```

**Fields:**
- `Targets` - List of target formats to compile to

### CompilationResult
```go
type CompilationResult struct {
    Path    string
    Content string
}
```

**Fields:**
- `Path` - Relative path where content should be written (e.g., "cleanCode_meaningfulNames.md")
- `Content` - Compiled content ready to write

### TargetCompiler Interface
```go
type TargetCompiler interface {
    Name() string
    Compile(resource Resource) ([]CompilationResult, error)
}
```

**Methods:**
- `Name()` - Returns target identifier (matches Target enum value)
- `Compile()` - Transforms resource into target-specific format(s)

### Compiler
```go
type Compiler struct {
    targets map[Target]TargetCompiler
}

func NewCompiler() *Compiler
func (c *Compiler) RegisterTarget(target Target, compiler TargetCompiler) error
func (c *Compiler) Compile(resource Resource, opts CompileOptions) ([]CompilationResult, error)
```

**Methods:**
- `NewCompiler()` - Creates compiler with all built-in targets registered
- `RegisterTarget()` - Adds or replaces target compiler
- `Compile()` - Compiles resource for all requested targets

## Algorithm

### Compilation Pipeline

1. Validate CompileOptions (non-empty targets list)
2. For each requested target:
   - Look up registered TargetCompiler
   - Call compiler.Compile(resource)
   - Collect CompilationResults
3. Return aggregated results from all targets

**Pseudocode:**
```
function Compile(resource, opts):
    if opts.Targets is empty:
        return error("no targets specified")
    
    results = []
    for target in opts.Targets:
        compiler = lookupCompiler(target)
        if compiler is nil:
            return error("unknown target: " + target)
        
        targetResults = compiler.Compile(resource)
        results.append(targetResults)
    
    return results
```

### Path Generation

**Rules:**
- Pattern: `{collection-id}_{item-id}.{ext}`
- Example: `cleanCode_meaningfulNames.mdc`

**Prompts (most targets):**
- Pattern: `{collection-id}_{item-id}.{ext}`
- Example: `codeReview_reviewPR.md`

**Claude Prompts (special case):**
- Pattern: `{collection-id}_{item-id}/SKILL.md`
- Example: `codeReview_reviewPR/SKILL.md`

**Extension by Target:**
- Cursor rules: `.mdc`
- Cursor prompts: `.md`
- Kiro: `.md`
- Claude: `.md`
- Copilot instructions: `.instructions.md`
- Copilot prompts: `.prompt.md`
- Markdown: `.md`

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Empty targets list | Return error "no targets specified" |
| Unknown target | Return error "unknown target: {name}" |
| Target compiler returns empty results | Include empty array in aggregated results |
| Target compiler returns error | Propagate error, stop compilation |
| Multiple targets requested | Compile independently, aggregate results |
| Resource with special characters in ID | Sanitize IDs for filesystem safety |

## Dependencies

- Resource model from ai-resource-core-go (Ruleset, Rule, Promptset, Prompt)
- Target-specific compilers (implement TargetCompiler interface)
- Metadata block generation (from metadata-block.md spec)

## Implementation Mapping

**Source files:**
- `pkg/compiler/compiler.go` - Compiler struct, registration, pipeline
- `pkg/compiler/types.go` - Target enum, CompileOptions, CompilationResult
- `pkg/compiler/interface.go` - TargetCompiler interface
- `pkg/targets/cursor.go` - Cursor target compiler
- `pkg/targets/kiro.go` - Kiro target compiler
- `pkg/targets/claude.go` - Claude target compiler
- `pkg/targets/copilot.go` - Copilot target compiler
- `pkg/targets/markdown.go` - Markdown target compiler

**Related specs:**
- `metadata-block.md` - Metadata embedding used by all target compilers
- `markdown-compiler.md` - Markdown target implementation
- `kiro-compiler.md` - Kiro target implementation
- `cursor-compiler.md` - Cursor target implementation
- `claude-compiler.md` - Claude target implementation
- `copilot-compiler.md` - Copilot target implementation

## Examples

### Example 1: Single Target Compilation

**Input:**
```go
compiler := NewCompiler()
resource := Resource{
    Type: "rule",
    Ruleset: Ruleset{ID: "cleanCode", Name: "Clean Code"},
    Rule: Rule{ID: "meaningfulNames", Name: "Use Meaningful Names", Enforcement: "must"},
    Body: "Use descriptive names.",
}
opts := CompileOptions{Targets: []Target{TargetMarkdown}}

results, err := compiler.Compile(resource, opts)
```

**Expected Output:**
```go
[]CompilationResult{
    {
        Path: "cleanCode_meaningfulNames.md",
        Content: "---\nruleset:\n  id: cleanCode\n...",
    },
}
```

**Verification:**
- Single result returned
- Path follows {collection-id}_{item-id}.{ext} pattern
- Content includes metadata block and rule body

### Example 2: Multi-Target Compilation

**Input:**
```go
compiler := NewCompiler()
resource := Resource{
    Type: "rule",
    Ruleset: Ruleset{ID: "cleanCode", Name: "Clean Code"},
    Rule: Rule{ID: "meaningfulNames", Name: "Use Meaningful Names", Enforcement: "must"},
    Body: "Use descriptive names.",
}
opts := CompileOptions{Targets: []Target{TargetMarkdown, TargetKiro, TargetCursor}}

results, err := compiler.Compile(resource, opts)
```

**Expected Output:**
```go
[]CompilationResult{
    {Path: "cleanCode_meaningfulNames.md", Content: "..."},      // Markdown
    {Path: "cleanCode_meaningfulNames.md", Content: "..."},      // Kiro
    {Path: "cleanCode_meaningfulNames.mdc", Content: "..."},     // Cursor
}
```

**Verification:**
- Three results returned (one per target)
- Each target produces appropriate extension
- Content formatted per target requirements

### Example 3: Claude Prompt (Special Path)

**Input:**
```go
compiler := NewCompiler()
resource := Resource{
    Type: "prompt",
    Promptset: Promptset{ID: "codeReview", Name: "Code Review"},
    Prompt: Prompt{ID: "reviewPR", Name: "Review Pull Request"},
    Body: "Review this PR.",
}
opts := CompileOptions{Targets: []Target{TargetClaude}}

results, err := compiler.Compile(resource, opts)
```

**Expected Output:**
```go
[]CompilationResult{
    {
        Path: "codeReview_reviewPR/SKILL.md",
        Content: "Review this PR.",
    },
}
```

**Verification:**
- Path uses directory structure with SKILL.md
- Content is plain body (no metadata for prompts)

### Example 4: Error Handling

**Input:**
```go
compiler := NewCompiler()
resource := Resource{...}
opts := CompileOptions{Targets: []Target{"invalid"}}

results, err := compiler.Compile(resource, opts)
```

**Expected Output:**
```go
results == nil
err.Error() == "unknown target: invalid"
```

**Verification:**
- Error returned for unknown target
- No results produced

### Example 5: Multi-Resource Compilation

**Input:**
```go
import (
    "github.com/jomadu/ai-resource-core-go/pkg/core"
    "github.com/jomadu/ai-resource-compiler-go/pkg/compiler"
)

// Load multi-document YAML
resources, err := core.LoadResources("resources.yaml")
if err != nil {
    log.Fatal(err)
}

// Compile each resource
compiler := compiler.NewCompiler()
opts := compiler.CompileOptions{
    Targets: []compiler.Target{compiler.TargetMarkdown},
}

allResults := []compiler.CompilationResult{}
for _, resource := range resources {
    results, err := compiler.Compile(resource, opts)
    if err != nil {
        log.Fatal(err)
    }
    allResults = append(allResults, results...)
}
```

**Expected Output:**
```go
// If resources.yaml contains 2 rules and 1 prompt:
len(allResults) == 3
allResults[0].Path == "cleanCode_meaningfulNames.md"
allResults[1].Path == "security_noSecrets.md"
allResults[2].Path == "codeReview_reviewPR.md"
```

**Verification:**
- Each resource from multi-document YAML compiled independently
- Results aggregated across all resources
- Demonstrates integration with ai-resource-core-go

## Notes

**Design Philosophy:**
- **Pure transformation** - Compiler produces path + content pairs, doesn't perform I/O
- **Modular output** - Users control where and how to write files
- **Extensible** - New targets can be added via TargetCompiler interface
- **Independent targets** - Each target compiler operates without knowledge of others

**Path Structure Rationale:**
- Simple `{collection-id}_{item-id}` pattern avoids namespace complexity
- Extensions differentiate target formats
- Claude's directory structure accommodates tool-specific requirements

**Registration Pattern:**
- NewCompiler() pre-registers all built-in targets
- RegisterTarget() allows custom target compilers
- Map-based lookup enables O(1) target resolution

## Known Issues

None - this is a new specification.

## Areas for Improvement

- Consider adding validation for filesystem-safe IDs (sanitization)
- Evaluate parallel compilation for performance with many targets
- Explore caching compiled results for repeated compilations
- Consider adding CompileOptions.OutputDir for path prefixing
