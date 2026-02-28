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
    SupportedVersions() []string
    Compile(resource *airesource.Resource) ([]CompilationResult, error)
}
```

**Methods:**
- `Name()` - Returns target identifier (matches Target enum value)
- `SupportedVersions()` - Returns list of supported API versions (e.g., ["ai-resource/draft"])
- `Compile()` - Transforms resource into target-specific format(s)
  - Handles Rule, Ruleset, Prompt, Promptset kinds
  - Expands collections (Ruleset/Promptset) into multiple results
  - Resolves Body (string or array with fragments)
  - Returns one result per rule/prompt

### Compiler
```go
type Compiler struct {
    targets map[Target]TargetCompiler
}

func NewCompiler() *Compiler
func (c *Compiler) RegisterTarget(target Target, compiler TargetCompiler) error
func (c *Compiler) Compile(resource *airesource.Resource, opts CompileOptions) ([]CompilationResult, error)
```

**Methods:**
- `NewCompiler()` - Creates compiler with all built-in targets registered
- `RegisterTarget()` - Adds or replaces target compiler
- `Compile()` - Compiles resource for all requested targets
  - Validates resource structure (apiVersion, kind, metadata.id)
  - For each target: calls target.Compile(resource)
  - Aggregates results from all targets

## Shared Functions

Target compilers use these shared functions to generate consistent file paths.

### BuildCollectionPath

```go
func BuildCollectionPath(collectionID, itemID, extension string) string
```

**Parameters:**
- `collectionID` - Collection identifier (Ruleset or Promptset ID)
- `itemID` - Item identifier (rule or prompt map key)
- `extension` - File extension (e.g., ".md", ".mdc", ".instructions.md")

**Returns:**
- File path string in format `{collectionID}_{itemID}{extension}`

**Algorithm:**
1. Concatenate collectionID + "_" + itemID
2. Append extension
3. Return formatted path

**Example:**
```go
path := BuildCollectionPath("cleanCode", "meaningfulNames", ".md")
// Returns: "cleanCode_meaningfulNames.md"

path := BuildCollectionPath("codeReview", "reviewPR", ".prompt.md")
// Returns: "codeReview_reviewPR.prompt.md"
```

### BuildStandalonePath

```go
func BuildStandalonePath(resourceID, extension string) string
```

**Parameters:**
- `resourceID` - Resource identifier (from Metadata.ID)
- `extension` - File extension (e.g., ".md", ".mdc", ".instructions.md")

**Returns:**
- File path string in format `{resourceID}{extension}`

**Algorithm:**
1. Concatenate resourceID + extension
2. Return formatted path

**Example:**
```go
path := BuildStandalonePath("meaningfulNames", ".md")
// Returns: "meaningfulNames.md"

path := BuildStandalonePath("reviewPR", ".prompt.md")
// Returns: "reviewPR.prompt.md"
```

### BuildClaudeCollectionPath

```go
func BuildClaudeCollectionPath(collectionID, itemID string) string
```

**Parameters:**
- `collectionID` - Collection identifier (Promptset ID)
- `itemID` - Item identifier (prompt map key)

**Returns:**
- Directory path string in format `{collectionID}_{itemID}/SKILL.md`

**Algorithm:**
1. Concatenate collectionID + "_" + itemID
2. Append "/SKILL.md"
3. Return formatted path

**Example:**
```go
path := BuildClaudeCollectionPath("codeReview", "reviewPR")
// Returns: "codeReview_reviewPR/SKILL.md"
```

### BuildClaudeStandalonePath

```go
func BuildClaudeStandalonePath(resourceID string) string
```

**Parameters:**
- `resourceID` - Resource identifier (from Metadata.ID)

**Returns:**
- Directory path string in format `{resourceID}/SKILL.md`

**Algorithm:**
1. Concatenate resourceID + "/SKILL.md"
2. Return formatted path

**Example:**
```go
path := BuildClaudeStandalonePath("reviewPR")
// Returns: "reviewPR/SKILL.md"
```

## Algorithm

### Compilation Pipeline

1. Validate resource structure (apiVersion, kind, metadata.id)
2. Validate CompileOptions (non-empty targets list)
3. For each requested target:
   - Look up registered TargetCompiler
   - Check target supports resource.APIVersion
   - Call target.Compile(resource)
   - Collect CompilationResults
4. Return aggregated results from all targets

**Pseudocode:**
```
function Compile(resource, opts):
    // Step 1: Validate resource
    if resource.APIVersion is empty:
        return error("missing apiVersion")
    
    if resource.Kind is empty:
        return error("missing kind")
    
    if resource.Metadata.ID is empty:
        return error("missing metadata.id")
    
    // Step 2: Validate options
    if opts.Targets is empty:
        return error("no targets specified")
    
    // Step 3: Compile for each target
    results = []
    for target in opts.Targets:
        compiler = lookupCompiler(target)
        if compiler is nil:
            return error("unknown target: " + target)
        
        // Check version compatibility
        supportedVersions = compiler.SupportedVersions()
        if resource.APIVersion not in supportedVersions:
            return error("target " + target + " does not support apiVersion: " + resource.APIVersion)
        
        // Compile resource
        targetResults = compiler.Compile(resource)
        results.append(targetResults)
    
    // Step 4: Return results
    return results
```

### Path Generation

Target compilers use shared path generation functions to ensure consistency.

**Collection items (Ruleset/Promptset):**
- Use `BuildCollectionPath(collectionID, itemID, extension)`
- Pattern: `{collection-id}_{item-id}.{ext}`
- Example: `cleanCode_meaningfulNames.mdc`

**Standalone resources (Rule/Prompt):**
- Use `BuildStandalonePath(resourceID, extension)`
- Pattern: `{resource-id}.{ext}`
- Example: `meaningfulNames.mdc`

**Claude collection items (special case):**
- Use `BuildClaudeCollectionPath(collectionID, itemID)`
- Pattern: `{collection-id}_{item-id}/SKILL.md`
- Example: `codeReview_reviewPR/SKILL.md`

**Claude standalone resources (special case):**
- Use `BuildClaudeStandalonePath(resourceID)`
- Pattern: `{resource-id}/SKILL.md`
- Example: `reviewPR/SKILL.md`

**Extension by Target:**
- Cursor rules: `.mdc`
- Cursor prompts: `.md`
- Kiro: `.md`
- Claude: `.md`
- Copilot instructions: `.instructions.md`
- Copilot prompts: `.prompt.md`
- Markdown: `.md`

## Version Handling

The compiler supports multiple spec versions by inspecting `resource.APIVersion`.

### Strategy

Each target compiler checks the resource version and handles version-specific structures:

**Pseudocode:**
```
function Compile(resource):
    switch resource.APIVersion:
    case "ai-resource/v1":
        return compileV1(resource)
    case "ai-resource/v2":
        return compileV2(resource)
    default:
        return error("unsupported apiVersion: " + resource.APIVersion)
```

### Version-Specific Compilation

When spec versions introduce breaking changes (e.g., field structure changes), target compilers must:

1. Check `resource.APIVersion` 
2. Type assert `resource.Spec` to version-specific structure
3. Access fields according to that version's schema
4. Generate output appropriate for that version

**Example - Enforcement field change:**
```go
func (c *CursorCompiler) Compile(resource Resource) ([]CompilationResult, error) {
    switch resource.APIVersion {
    case "ai-resource/v1":
        // v1: enforcement is string
        spec := resource.Spec.(RuleSpec)
        enforcement := spec.Enforcement.(string)
        
    case "ai-resource/v2":
        // v2: enforcement is object
        spec := resource.Spec.(RuleSpecV2)
        enforcement := spec.Enforcement.Level
    }
}
```

### Backward Compatibility

- Compilers SHOULD support all non-deprecated spec versions
- Compilers MUST return clear errors for unsupported versions
- Version support is independent per target (Cursor may support v1-v2, Kiro may support v1-v3)

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Empty targets list | Return error "no targets specified" |
| Unknown target | Return error "unknown target: {name}" |
| Unsupported apiVersion | Return error "target {name} does not support apiVersion: {version}" |
| Target compiler returns empty results | Include empty array in aggregated results |
| Target compiler returns error | Propagate error, stop compilation |
| Multiple targets requested | Compile independently, aggregate results |
| Resource with special characters in ID | Sanitize IDs for filesystem safety |

## Dependencies

- Resource model from ai-resource-core-go (Resource, Ruleset, Rule, Promptset, Prompt, RuleItem, PromptItem, Metadata)
  - **Note:** Resource includes `APIVersion` field for version detection
  - **Note:** Collections (Ruleset, Promptset) contain maps of items
  - **Note:** Body is union type (string or array) that may need resolution
  - **Note:** Scope is []ScopeEntry, each with Files []string
- Target-specific compilers (implement TargetCompiler interface)
- Metadata block generation (from metadata-block.md spec)
- Path generation functions (shared across all target compilers)
- Resource validation (from validation-rules.md spec)
- Body resolution from ai-resource-core-go (ResolveBody function)

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
- `internal/format/paths.go` - Implements `BuildCollectionPath()`, `BuildStandalonePath()`, `BuildClaudeCollectionPath()`, and `BuildClaudeStandalonePath()` functions
- `internal/format/validation.go` - Implements `ValidateResourceIDs()` and `ValidateRuleForCompilation()` functions

**Related specs:**
- `validation-rules.md` - Resource validation before compilation
- `metadata-block.md` - Metadata embedding used by all target compilers
- `markdown-compiler.md` - Markdown target implementation
- `kiro-compiler.md` - Kiro target implementation
- `cursor-compiler.md` - Cursor target implementation
- `claude-compiler.md` - Claude target implementation
- `copilot-compiler.md` - Copilot target implementation

## Examples

### Example 1: Single Rule Compilation

**Input:**
```go
compiler := NewCompiler()

// Load a standalone Rule resource
resource, _ := airesource.LoadResource("rule.yaml")
// resource.Kind == "Rule"
// resource.Metadata.ID == "meaningfulNames"

opts := CompileOptions{Targets: []Target{TargetMarkdown}}

results, err := compiler.Compile(resource, opts)
```

**Expected Output:**
```go
[]CompilationResult{
    {
        Path: "meaningfulNames.md",
        Content: "---\nruleset:\n  id: meaningfulNames\n...",
    },
}
```

**Verification:**
- Single result returned
- Path is just `{resource-id}.{ext}` (no collection prefix)
- Content includes metadata block and rule body

### Example 2: Ruleset Expansion

**Input:**
```go
compiler := NewCompiler()

// Load a Ruleset resource
resource, _ := airesource.LoadResource("ruleset.yaml")
// resource.Kind == "Ruleset"
// resource.Metadata.ID == "cleanCode"
// resource.Spec.Rules == map[string]RuleItem{
//     "meaningfulNames": {...},
//     "smallFunctions": {...},
// }

opts := CompileOptions{Targets: []Target{TargetMarkdown}}

results, err := compiler.Compile(resource, opts)
```

**Expected Output:**
```go
[]CompilationResult{
    {
        Path: "cleanCode_meaningfulNames.md",
        Content: "...",
    },
    {
        Path: "cleanCode_smallFunctions.md",
        Content: "...",
    },
}
```

**Verification:**
- Two results returned (one per rule in ruleset)
- Each path uses ruleset ID as collection ID
- Each path uses rule map key as item ID
- Compiler expanded Ruleset into individual rules

### Example 3: Multi-Target Compilation

**Input:**
```go
compiler := NewCompiler()

resource, _ := airesource.LoadResource("rule.yaml")

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

### Example 4: Claude Prompt (Special Path)

**Input:**
```go
compiler := NewCompiler()

resource, _ := airesource.LoadResource("prompt.yaml")
// resource.Kind == "Prompt"
// resource.Metadata.ID == "reviewPR"

opts := CompileOptions{Targets: []Target{TargetClaude}}

results, err := compiler.Compile(resource, opts)
```

**Expected Output:**
```go
[]CompilationResult{
    {
        Path: "reviewPR/SKILL.md",
        Content: "Review this PR.",
    },
}
```

**Verification:**
- Path uses directory structure with SKILL.md
- Content is plain body (no metadata for prompts)
- Path is just `{resource-id}/SKILL.md` (no collection prefix)

### Example 5: Error Handling

**Input:**
```go
compiler := NewCompiler()
resource, _ := airesource.LoadResource("rule.yaml")
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

### Example 6: Multi-Resource Compilation

**Input:**
```go
import (
    "github.com/aws/ai-resource-core-go/pkg/airesource"
    "github.com/jomadu/ai-resource-compiler-go/pkg/compiler"
)

// Load multi-document YAML
resources, err := airesource.LoadResources("resources.yaml")
if err != nil {
    log.Fatal(err)
}

// Compile each resource
c := compiler.NewCompiler()
opts := compiler.CompileOptions{
    Targets: []compiler.Target{compiler.TargetMarkdown},
}

allResults := []compiler.CompilationResult{}
for _, resource := range resources {
    results, err := c.Compile(resource, opts)
    if err != nil {
        log.Fatal(err)
    }
    allResults = append(allResults, results...)
}
```

**Expected Output:**
```go
// If resources.yaml contains 1 ruleset with 2 rules and 1 prompt:
len(allResults) == 3
allResults[0].Path == "cleanCode_meaningfulNames.md"
allResults[1].Path == "cleanCode_smallFunctions.md"
allResults[2].Path == "codeReview_reviewPR.md"
```

**Verification:**
- Each resource from multi-document YAML compiled independently
- Ruleset expanded into 2 individual rules
- Results aggregated across all resources
- Demonstrates integration with ai-resource-core-go

### Example 7: Version Handling

### Example 7: Version Handling

**Input:**
```go
// Load resources with different versions
draftResource, _ := airesource.LoadResource("draft-rule.yaml")  // apiVersion: ai-resource/draft
v1Resource, _ := airesource.LoadResource("v1-rule.yaml")        // apiVersion: ai-resource/v1 (future)

compiler := compiler.NewCompiler()
opts := compiler.CompileOptions{
    Targets: []compiler.Target{compiler.TargetMarkdown},
}

// Compile both versions
results1, _ := compiler.Compile(draftResource, opts)
results2, err := compiler.Compile(v1Resource, opts)
// err != nil if target doesn't support v1 yet
```

**Expected Output:**
```go
// Draft version compiles successfully
len(results1) == 1

// V1 version may fail if not yet supported
err.Error() == "target markdown does not support apiVersion: ai-resource/v1"
```

**Verification:**
- Compiler checks version compatibility per target
- Targets can support different version ranges
- Clear error when version not supported

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
