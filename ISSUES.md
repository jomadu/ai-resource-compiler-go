# Design Issues Analysis

Generated: 2026-02-23

This document identifies design issues in the AI Resource Compiler specifications with options, tradeoffs, and recommendations for resolution.

---

## 1. Metadata Block Duplication Across Targets

**Issue:** Every target compiler will need to implement metadata block generation independently, leading to potential inconsistency.

**Options:**

A. **Shared metadata generator (current implicit design)**
   - Create `internal/format/metadata.go` with reusable functions
   - Each target calls the shared generator
   
B. **Pre-process metadata at compiler level**
   - Main compiler adds metadata before routing to targets
   - Targets receive pre-formatted content
   
C. **Metadata as separate compilation step**
   - Compiler pipeline: validate → add metadata → route to targets
   - Targets only handle target-specific formatting (frontmatter, extensions)

**Tradeoffs:**

| Approach | Pros | Cons |
|----------|------|------|
| A (Shared) | Flexible, targets control when to add metadata | Requires discipline, can drift |
| B (Pre-process) | Guaranteed consistency, simpler targets | Less flexible, harder to customize per target |
| C (Pipeline) | Clear separation of concerns | More complex architecture, potential over-engineering |

**Recommendation:** **Option A** - The specs already imply this with `internal/format/metadata.go`. Document the shared functions clearly and add integration tests that verify all targets produce consistent metadata blocks.

**Priority:** High  
**Effort:** Low

---

## 2. Path Generation Logic Scattered Across Targets

**Issue:** Each target compiler implements path generation (`{collection-id}_{item-id}.{ext}`), with Claude having a special case for prompts. This creates maintenance burden and inconsistency risk.

**Options:**

A. **Centralized path builder**
   - Create `internal/format/paths.go` with `BuildPath(resource, extension)` function
   - Special case for Claude prompts handled centrally
   
B. **Target-specific path generation (current design)**
   - Each target owns its path logic
   - Flexibility for future target-specific requirements
   
C. **Path templates in target registration**
   - Targets declare path templates during registration
   - Compiler fills in templates with resource data

**Tradeoffs:**

| Approach | Pros | Cons |
|----------|------|------|
| A (Centralized) | DRY, consistent, easier to change pattern | Less flexible for unusual targets |
| B (Scattered) | Maximum flexibility per target | Duplication, drift risk |
| C (Templates) | Declarative, flexible, testable | More complex, may not handle all cases |

**Recommendation:** **Option A** - Centralize the common pattern with an escape hatch for special cases. Implementation:

```go
// internal/format/paths.go
func BuildRulePath(rulesetID, ruleID, extension string) string
func BuildPromptPath(promptsetID, promptID, extension string) string
func BuildClaudePromptPath(promptsetID, promptID string) string // Special case
```

**Priority:** High  
**Effort:** Low

---

## 3. No Validation of Resource IDs for Filesystem Safety

**Issue:** Specs mention "sanitization handled by caller" but don't specify who the caller is or what sanitization means. IDs with `/`, `\`, or other special characters will break path generation.

**Options:**

A. **Validate at compiler entry point**
   - `Compiler.Compile()` validates IDs before routing to targets
   - Return error for invalid characters
   
B. **Sanitize at path generation**
   - Replace invalid characters with safe alternatives (e.g., `-`, `_`)
   - Log warnings about sanitization
   
C. **Assume valid IDs (current implicit design)**
   - Rely on ai-resource-core-go validation
   - Document requirement in specs

**Tradeoffs:**

| Approach | Pros | Cons |
|----------|------|------|
| A (Validate) | Fail fast, clear errors | Rejects potentially valid resources |
| B (Sanitize) | Permissive, always produces output | Silent data transformation, potential collisions |
| C (Assume) | Simple, trusts upstream | Fragile, poor error messages |

**Recommendation:** **Option A** - Add validation in `Compiler.Compile()` with clear error messages. Document the allowed character set in specs. This is a compiler responsibility, not a target responsibility.

```go
// pkg/compiler/validation.go
func ValidateResourceIDs(resource Resource) error {
    // Check for /, \, :, *, ?, ", <, >, |
    // Return descriptive error
}
```

**Priority:** High  
**Effort:** Low

---

## 4. Version Handling Strategy Incomplete

**Issue:** Specs mention version handling but don't specify:
- How to detect breaking changes between versions
- What happens when a target doesn't support a resource version
- Whether version support is documented/discoverable

**Options:**

A. **Version support matrix**
   - Each target declares supported versions
   - Compiler checks compatibility before routing
   - Clear error: "target X doesn't support version Y"
   
B. **Graceful degradation**
   - Targets attempt to compile any version
   - Fall back to common fields if version-specific fields missing
   
C. **Version-specific target compilers**
   - Register separate compilers per version (e.g., `cursor-v1`, `cursor-v2`)
   - Compiler routes based on version + target

**Tradeoffs:**

| Approach | Pros | Cons |
|----------|------|------|
| A (Matrix) | Clear contracts, fail fast | Maintenance burden, version explosion |
| B (Degradation) | Flexible, forward-compatible | Silent failures, unpredictable output |
| C (Separate) | Clean separation, testable | Duplication, registration complexity |

**Recommendation:** **Option A** - Add version support to the `TargetCompiler` interface:

```go
type TargetCompiler interface {
    Name() string
    SupportedVersions() []string
    Compile(resource Resource) ([]CompilationResult, error)
}
```

Check compatibility in `Compiler.Compile()` before routing. This makes version support explicit and discoverable.

**Priority:** Medium  
**Effort:** Medium

---

## 5. CLI Output Mode Ambiguity

**Issue:** CLI spec says "print to stdout" or "write to directory" but doesn't specify:
- How to handle multiple targets with same filename (e.g., markdown and kiro both produce `.md`)
- Whether stdout output includes paths or just content
- How to handle Claude's directory structure in stdout mode

**Options:**

A. **Target-prefixed paths in output mode**
   - When writing to directory: `--output ./out` creates `./out/markdown/`, `./out/kiro/`, etc.
   - Prevents filename collisions
   
B. **Require separate output dirs per target**
   - `arc compile --target markdown --output ./md --target kiro --output ./kiro`
   - Explicit but verbose
   
C. **Stdout shows path + content (current implicit)**
   - Format: `=== {path} ===\n{content}\n\n`
   - User can parse and write files themselves

**Tradeoffs:**

| Approach | Pros | Cons |
|----------|------|------|
| A (Prefixed) | Automatic collision avoidance | Opinionated directory structure |
| B (Separate) | Explicit control | Verbose, awkward for multi-target |
| C (Stdout) | Flexible, scriptable | Requires parsing, not beginner-friendly |

**Recommendation:** **Hybrid approach:**
- Stdout mode: Use format `=== {target}/{path} ===\n{content}\n\n` (Option C)
- Directory mode: Create target subdirectories automatically (Option A)
- Add `--flat` flag to disable subdirectories for single-target use

Update CLI spec to document this clearly.

**Priority:** Medium  
**Effort:** Medium

---

## 6. No Guidance on Multi-Resource Compilation

**Issue:** README shows multi-resource example but specs don't address:
- Should compiler accept `[]Resource` or is iteration the user's job?
- How to handle partial failures (resource 3 of 10 fails)
- Whether to aggregate results across resources

**Options:**

A. **Single resource only (current design)**
   - User iterates and aggregates
   - Simple, flexible
   
B. **Batch compilation API**
   - `CompileMany([]Resource, opts) ([]CompilationResult, error)`
   - All-or-nothing error handling
   
C. **Batch with partial results**
   - `CompileMany([]Resource, opts) ([]CompilationResult, []error)`
   - Continue on error, return what succeeded

**Tradeoffs:**

| Approach | Pros | Cons |
|----------|------|------|
| A (Single) | Simple, clear error handling | Boilerplate for users |
| B (Batch) | Convenient, atomic | Fails entire batch on one error |
| C (Partial) | Resilient, informative | Complex error handling |

**Recommendation:** **Option A** - Keep the API simple. Multi-resource compilation is a CLI concern, not a library concern. The CLI can implement retry logic, progress reporting, etc. Document the iteration pattern in README as the recommended approach.

**Priority:** Low  
**Effort:** Low

---

## 7. Enforcement Header Formatting Inconsistency Risk

**Issue:** Enforcement header format `# {Name} ({ENFORCEMENT})` is specified but:
- What if rule name contains parentheses?
- What if enforcement is custom (not must/should/may)?
- Should there be a space before the parentheses?

**Options:**

A. **Strict validation**
   - Reject rule names with parentheses
   - Reject non-standard enforcement values
   
B. **Escape special characters**
   - Replace `(` with `\(` in rule names
   - Allow any enforcement value
   
C. **Alternative format**
   - Use `# {Name} - {ENFORCEMENT}` or `# {Name} [{ENFORCEMENT}]`
   - Avoid ambiguity

**Tradeoffs:**

| Approach | Pros | Cons |
|----------|------|------|
| A (Strict) | Predictable, parseable | Restrictive, may reject valid rules |
| B (Escape) | Flexible | Ugly output, parsing complexity |
| C (Alternative) | Clean, unambiguous | Breaks existing format |

**Recommendation:** **Option A** - Validate at compile time. Rule names with parentheses are rare and can be rewritten. Enforcement values should be validated by ai-resource-core-go. Add validation to `Compiler.Compile()`:

```go
func ValidateRuleForCompilation(rule Rule) error {
    if strings.ContainsAny(rule.Name, "()") {
        return fmt.Errorf("rule name cannot contain parentheses: %s", rule.Name)
    }
    // Validate enforcement is must/should/may
}
```

**Priority:** Medium  
**Effort:** Low

---

## Summary

| Issue | Recommendation | Priority | Effort |
|-------|---------------|----------|--------|
| 1. Metadata duplication | Shared generator in `internal/format/metadata.go` | High | Low |
| 2. Path generation scattered | Centralized path builder with escape hatch | High | Low |
| 3. No ID validation | Validate IDs in `Compiler.Compile()` | High | Low |
| 4. Version handling incomplete | Add `SupportedVersions()` to interface | Medium | Medium |
| 5. CLI output ambiguity | Hybrid stdout/directory with `--flat` flag | Medium | Medium |
| 6. Multi-resource unclear | Keep single-resource API, document pattern | Low | Low |
| 7. Enforcement header risk | Validate rule names and enforcement values | Medium | Low |

## Next Steps

1. Update specs to address recommendations 1-3 (high priority, low effort)
2. Create validation spec for ID and rule name requirements
3. Update CLI spec with output format details
4. Add version support to compiler architecture spec
