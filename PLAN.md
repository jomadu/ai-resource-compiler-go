# Implementation Plan

**Project:** AI Resource Compiler (Go)  
**Created:** 2026-02-28  
**Status:** Bootstrap Phase → Implementation  
**Audit Reference:** [AUDIT.md](AUDIT.md)

## Overview

Transform AI resources into tool-specific formats through a 4-phase implementation:
1. **Foundation** - Core infrastructure (P1)
2. **Compilation** - Target compilers (P2)
3. **Interface** - CLI tool (P3)
4. **Quality** - Tests and validation (P3)

## Phase 1: Foundation (P1 - Critical)

All subsequent work blocked without these components.

### Task 1.1: Initialize Go Module

**Context:**  
Project has no `go.mod` file. Cannot build, test, or import packages.

**Spec Reference:** AUDIT.md #1

**Actions:**
```bash
go mod init github.com/jomadu/ai-resource-compiler-go
go get github.com/jomadu/ai-resource-core-go
```

**Acceptance Criteria:**
- [ ] `go.mod` exists at repository root
- [ ] Module path is `github.com/jomadu/ai-resource-compiler-go`
- [ ] Dependency on `ai-resource-core-go` declared
- [ ] `go build ./...` succeeds (even with no code)

**Estimated Effort:** 5 minutes

---

### Task 1.2: Implement Core Compiler Types

**Context:**  
Define public API types that all components depend on.

**Spec Reference:**  
- AUDIT.md #2
- [specs/compiler-architecture.md](specs/compiler-architecture.md) - Data Structures section

**File:** `pkg/compiler/types.go`

**Implementation:**
```go
package compiler

type Target string

const (
    TargetCursor   Target = "cursor"
    TargetKiro     Target = "kiro"
    TargetClaude   Target = "claude"
    TargetCopilot  Target = "copilot"
    TargetMarkdown Target = "markdown"
)

type CompileOptions struct {
    Targets []Target
}

type CompilationResult struct {
    Path    string
    Content string
}
```

**Acceptance Criteria:**
- [ ] File `pkg/compiler/types.go` exists
- [ ] Target enum with 5 constants defined
- [ ] CompileOptions struct with Targets field
- [ ] CompilationResult struct with Path and Content fields
- [ ] `go build ./pkg/compiler` succeeds

**Estimated Effort:** 10 minutes

---

### Task 1.3: Implement TargetCompiler Interface

**Context:**  
Define contract that all target compilers must implement.

**Spec Reference:**  
- AUDIT.md #2
- [specs/compiler-architecture.md](specs/compiler-architecture.md) - TargetCompiler Interface section

**File:** `pkg/compiler/interface.go`

**Implementation:**
```go
package compiler

import "github.com/jomadu/ai-resource-core-go/pkg/airesource"

type TargetCompiler interface {
    Name() string
    SupportedVersions() []string
    Compile(resource *airesource.Resource) ([]CompilationResult, error)
}
```

**Acceptance Criteria:**
- [ ] File `pkg/compiler/interface.go` exists
- [ ] TargetCompiler interface with 3 methods defined
- [ ] Imports `ai-resource-core-go` types
- [ ] `go build ./pkg/compiler` succeeds

**Estimated Effort:** 5 minutes

---

### Task 1.4: Implement Path Generation Functions

**Context:**  
Shared functions for generating consistent file paths across all target compilers.

**Spec Reference:**  
- AUDIT.md #3
- [specs/compiler-architecture.md](specs/compiler-architecture.md) - Shared Functions section

**File:** `internal/format/paths.go`

**Functions to Implement:**
1. `BuildCollectionPath(collectionID, itemID, extension string) string`
2. `BuildStandalonePath(resourceID, extension string) string`
3. `BuildClaudeCollectionPath(collectionID, itemID string) string`
4. `BuildClaudeStandalonePath(resourceID string) string`

**Acceptance Criteria:**
- [ ] File `internal/format/paths.go` exists
- [ ] All 4 functions implemented
- [ ] `BuildCollectionPath("cleanCode", "meaningfulNames", ".md")` returns `"cleanCode_meaningfulNames.md"`
- [ ] `BuildStandalonePath("meaningfulNames", ".md")` returns `"meaningfulNames.md"`
- [ ] `BuildClaudeCollectionPath("codeReview", "reviewPR")` returns `"codeReview_reviewPR/SKILL.md"`
- [ ] `BuildClaudeStandalonePath("reviewPR")` returns `"reviewPR/SKILL.md"`
- [ ] `go build ./internal/format` succeeds

**Estimated Effort:** 15 minutes

---

### Task 1.5: Implement Validation Functions

**Context:**  
Validate IDs are filesystem-safe and rule names don't contain parentheses.

**Spec Reference:**  
- AUDIT.md #4
- [specs/validation-rules.md](specs/validation-rules.md) - Algorithm section

**File:** `internal/format/validation.go`

**Functions to Implement:**
1. `ValidateID(id string) error`
2. `ValidateRuleName(name string) error`

**Validation Rules:**
- IDs: Only `a-z A-Z 0-9 - _` allowed
- IDs: Cannot be empty
- Rule names: Cannot contain `(` or `)`

**Acceptance Criteria:**
- [ ] File `internal/format/validation.go` exists
- [ ] `ValidateID("cleanCode")` returns `nil`
- [ ] `ValidateID("clean/code")` returns error with message containing "invalid character '/'"
- [ ] `ValidateID("")` returns error "ID cannot be empty"
- [ ] `ValidateRuleName("Use Meaningful Names")` returns `nil`
- [ ] `ValidateRuleName("Use (Smart) Names")` returns error containing "cannot contain parentheses"
- [ ] `go build ./internal/format` succeeds

**Estimated Effort:** 20 minutes

---

### Task 1.6: Implement Core Compiler Pipeline

**Context:**  
Orchestrate compilation across multiple targets with validation and error handling.

**Spec Reference:**  
- AUDIT.md #2
- [specs/compiler-architecture.md](specs/compiler-architecture.md) - Algorithm section

**File:** `pkg/compiler/compiler.go`

**Implementation Requirements:**
- `Compiler` struct with target registry (map)
- `NewCompiler()` - Creates compiler, registers built-in targets (initially empty)
- `RegisterTarget(target Target, compiler TargetCompiler) error`
- `Compile(resource *airesource.Resource, opts CompileOptions) ([]CompilationResult, error)`

**Compile Algorithm:**
1. Validate resource (apiVersion, kind, metadata.id not empty)
2. Validate options (targets not empty)
3. For each target:
   - Look up compiler
   - Check version support
   - Call Compile()
   - Aggregate results
4. Return all results

**Acceptance Criteria:**
- [ ] File `pkg/compiler/compiler.go` exists
- [ ] `NewCompiler()` returns initialized Compiler
- [ ] `RegisterTarget()` adds/replaces target compilers
- [ ] `Compile()` validates resource structure
- [ ] `Compile()` validates options
- [ ] `Compile()` returns error for unknown target
- [ ] `Compile()` returns error for unsupported apiVersion
- [ ] `Compile()` aggregates results from all targets
- [ ] `go build ./pkg/compiler` succeeds

**Estimated Effort:** 45 minutes

---

## Phase 2: Compilation (P2 - High)

Depends on Phase 1 completion.

### Task 2.1: Implement Metadata Generation Functions

**Context:**  
Generate YAML metadata blocks and enforcement headers for rules.

**Spec Reference:**  
- AUDIT.md #5
- [specs/metadata-block.md](specs/metadata-block.md) - Shared Functions section

**File:** `internal/format/metadata.go`

**Functions to Implement:**
1. `GenerateRuleMetadataBlockFromRuleset(ruleset *airesource.Ruleset, ruleID string) string`
2. `GenerateRuleMetadataBlockFromRule(rule *airesource.Rule) string`

**Algorithm:**
- Resolve body using `airesource.ResolveBody()`
- Build YAML metadata block (ruleset + rule sections OR flat structure)
- Generate enforcement header: `# {Name} ({ENFORCEMENT})`
- Concatenate: metadata + header + body

**Acceptance Criteria:**
- [ ] File `internal/format/metadata.go` exists
- [ ] `GenerateRuleMetadataBlockFromRuleset()` returns complete rule content
- [ ] Output includes YAML metadata block with ruleset and rule sections
- [ ] Output includes enforcement header with uppercased enforcement
- [ ] Output includes resolved body content
- [ ] `GenerateRuleMetadataBlockFromRule()` returns complete rule content
- [ ] Standalone rule uses flat metadata structure (no nesting)
- [ ] Optional fields (description, scope) omitted when not present
- [ ] `go build ./internal/format` succeeds

**Estimated Effort:** 60 minutes

---

### Task 2.2: Implement Markdown Compiler

**Context:**  
Generate vanilla markdown output for generic use.

**Spec Reference:**  
- AUDIT.md #6
- [specs/markdown-compiler.md](specs/markdown-compiler.md)

**File:** `pkg/targets/markdown.go`

**Implementation:**
- Struct: `MarkdownCompiler`
- Methods: `Name()`, `SupportedVersions()`, `Compile()`
- Extension: `.md` for all outputs
- Rules: Use metadata generation functions
- Prompts: Body only (no metadata)

**Acceptance Criteria:**
- [ ] File `pkg/targets/markdown.go` exists
- [ ] Implements `TargetCompiler` interface
- [ ] `Name()` returns `"markdown"`
- [ ] `SupportedVersions()` returns `["ai-resource/draft"]`
- [ ] Handles Rule, Ruleset, Prompt, Promptset kinds
- [ ] Rules include metadata block and enforcement header
- [ ] Prompts include body only
- [ ] Paths use `.md` extension
- [ ] Validates IDs before compilation
- [ ] `go build ./pkg/targets` succeeds

**Estimated Effort:** 45 minutes

---

### Task 2.3: Implement Kiro Compiler

**Context:**  
Generate Kiro CLI steering rules and prompts.

**Spec Reference:**  
- AUDIT.md #7
- [specs/kiro-compiler.md](specs/kiro-compiler.md)

**File:** `pkg/targets/kiro.go`

**Implementation:**
- Struct: `KiroCompiler`
- Format: Identical to Markdown (no frontmatter)
- Extension: `.md` for all outputs

**Acceptance Criteria:**
- [ ] File `pkg/targets/kiro.go` exists
- [ ] Implements `TargetCompiler` interface
- [ ] `Name()` returns `"kiro"`
- [ ] `SupportedVersions()` returns `["ai-resource/draft"]`
- [ ] Output format matches Markdown compiler
- [ ] `go build ./pkg/targets` succeeds

**Estimated Effort:** 30 minutes

---

### Task 2.4: Implement Cursor Compiler

**Context:**  
Generate Cursor IDE rules (.mdc) and commands (.md).

**Spec Reference:**  
- AUDIT.md #8
- [specs/cursor-compiler.md](specs/cursor-compiler.md)

**File:** `pkg/targets/cursor.go`

**Implementation:**
- Struct: `CursorCompiler`
- Rules: `.mdc` extension with MDC frontmatter
- Prompts: `.md` extension (no frontmatter)
- MDC frontmatter: `alwaysApply: false` and optional `description`

**Acceptance Criteria:**
- [ ] File `pkg/targets/cursor.go` exists
- [ ] Implements `TargetCompiler` interface
- [ ] `Name()` returns `"cursor"`
- [ ] Rules use `.mdc` extension
- [ ] Rules include MDC frontmatter before metadata block
- [ ] Prompts use `.md` extension
- [ ] Prompts have no frontmatter
- [ ] `go build ./pkg/targets` succeeds

**Estimated Effort:** 45 minutes

---

### Task 2.5: Implement Claude Compiler

**Context:**  
Generate Claude Code rules and skills with special directory structure for prompts.

**Spec Reference:**  
- AUDIT.md #9
- [specs/claude-compiler.md](specs/claude-compiler.md)

**File:** `pkg/targets/claude.go`

**Implementation:**
- Struct: `ClaudeCompiler`
- Rules: `.md` extension, optional frontmatter with `paths`
- Prompts: `{id}/SKILL.md` directory structure
- Use `BuildClaudeCollectionPath()` and `BuildClaudeStandalonePath()`

**Acceptance Criteria:**
- [ ] File `pkg/targets/claude.go` exists
- [ ] Implements `TargetCompiler` interface
- [ ] `Name()` returns `"claude"`
- [ ] Rules use `.md` extension
- [ ] Prompts use `{id}/SKILL.md` path structure
- [ ] Prompts contain body only
- [ ] `go build ./pkg/targets` succeeds

**Estimated Effort:** 45 minutes

---

### Task 2.6: Implement Copilot Compiler

**Context:**  
Generate GitHub Copilot instructions and prompts with applyTo frontmatter.

**Spec Reference:**  
- AUDIT.md #10
- [specs/copilot-compiler.md](specs/copilot-compiler.md)

**File:** `pkg/targets/copilot.go`

**Implementation:**
- Struct: `CopilotCompiler`
- Rules: `.instructions.md` extension with `applyTo` frontmatter
- Prompts: `.prompt.md` extension with `applyTo` frontmatter
- Extract file patterns from scope for `applyTo`

**Acceptance Criteria:**
- [ ] File `pkg/targets/copilot.go` exists
- [ ] Implements `TargetCompiler` interface
- [ ] `Name()` returns `"copilot"`
- [ ] Rules use `.instructions.md` extension
- [ ] Prompts use `.prompt.md` extension
- [ ] Both include `applyTo` frontmatter with file patterns
- [ ] `go build ./pkg/targets` succeeds

**Estimated Effort:** 45 minutes

---

### Task 2.7: Register Built-in Target Compilers

**Context:**  
Wire up all target compilers in `NewCompiler()`.

**Spec Reference:**  
- [specs/compiler-architecture.md](specs/compiler-architecture.md) - Compiler section

**File:** `pkg/compiler/compiler.go` (modify)

**Implementation:**
Update `NewCompiler()` to register all 5 target compilers:
```go
func NewCompiler() *Compiler {
    c := &Compiler{targets: make(map[Target]TargetCompiler)}
    c.RegisterTarget(TargetMarkdown, &markdown.MarkdownCompiler{})
    c.RegisterTarget(TargetKiro, &kiro.KiroCompiler{})
    c.RegisterTarget(TargetCursor, &cursor.CursorCompiler{})
    c.RegisterTarget(TargetClaude, &claude.ClaudeCompiler{})
    c.RegisterTarget(TargetCopilot, &copilot.CopilotCompiler{})
    return c
}
```

**Acceptance Criteria:**
- [ ] `NewCompiler()` registers all 5 targets
- [ ] Imports all target packages
- [ ] `go build ./pkg/compiler` succeeds
- [ ] Can compile to any target without manual registration

**Estimated Effort:** 10 minutes

---

## Phase 3: Interface (P3 - Medium)

Depends on Phase 2 completion.

### Task 3.1: Implement CLI Main Entry Point

**Context:**  
Create `arc` command-line tool entry point.

**Spec Reference:**  
- AUDIT.md #11
- [specs/cli-design.md](specs/cli-design.md)

**File:** `cmd/arc/main.go`

**Implementation:**
- Use CLI framework (e.g., `cobra` or `flag`)
- Define `arc compile` command
- Parse flags: `--target`, `--output`, `--flat`, `--help`
- Call compile logic

**Acceptance Criteria:**
- [ ] File `cmd/arc/main.go` exists
- [ ] `arc --help` shows usage information
- [ ] `arc compile --help` shows compile command help
- [ ] Accepts resource file as positional argument
- [ ] Accepts `--target` flag (repeatable)
- [ ] Accepts `--output` flag (default: "stdout")
- [ ] Accepts `--flat` flag
- [ ] `go build ./cmd/arc` succeeds
- [ ] Binary `arc` can be executed

**Estimated Effort:** 45 minutes

---

### Task 3.2: Implement Compile Command Logic

**Context:**  
Load resource, compile, and handle results.

**Spec Reference:**  
- [specs/cli-design.md](specs/cli-design.md) - Algorithm section

**File:** `cmd/arc/compile.go`

**Implementation:**
1. Validate arguments (file exists, targets specified)
2. Load resource using `airesource.LoadResource()`
3. Create compiler and compile
4. Route to output handler based on `--output` flag

**Acceptance Criteria:**
- [ ] File `cmd/arc/compile.go` exists
- [ ] Validates resource file exists
- [ ] Validates at least one target specified
- [ ] Loads resource from file
- [ ] Creates compiler and compiles resource
- [ ] Returns error with exit code 1 on failure
- [ ] Returns success with exit code 0 on success
- [ ] `go build ./cmd/arc` succeeds

**Estimated Effort:** 30 minutes

---

### Task 3.3: Implement Output Handlers

**Context:**  
Handle stdout and file output modes.

**Spec Reference:**  
- [specs/cli-design.md](specs/cli-design.md) - Output Modes section

**File:** `cmd/arc/output.go`

**Implementation:**
- `outputStdout(results []CompilationResult)` - Print with separators
- `outputFiles(results []CompilationResult, dir string, flat bool)` - Write to filesystem

**Stdout Format:**
```
=== {target}/{path} ===
{content}

```

**File Mode:**
- Default: `{output-dir}/{target}/{path}`
- Flat: `{output-dir}/{path}`
- Create directories as needed
- Report files written to stderr

**Acceptance Criteria:**
- [ ] File `cmd/arc/output.go` exists
- [ ] `outputStdout()` prints results with separators
- [ ] `outputFiles()` writes files to correct paths
- [ ] `outputFiles()` creates directories as needed
- [ ] `outputFiles()` respects `--flat` flag
- [ ] `outputFiles()` reports written files to stderr
- [ ] `go build ./cmd/arc` succeeds

**Estimated Effort:** 30 minutes

---

### Task 3.4: Install CLI Tool

**Context:**  
Make `arc` command available system-wide.

**Actions:**
```bash
go install ./cmd/arc
```

**Acceptance Criteria:**
- [ ] `arc` command available in PATH
- [ ] `arc --help` works from any directory
- [ ] Can compile resources from command line

**Estimated Effort:** 2 minutes

---

## Phase 4: Quality (P3 - Medium)

Depends on Phase 3 completion.

### Task 4.1: Write Unit Tests for Path Generation

**Context:**  
Verify path generation functions produce correct output.

**File:** `internal/format/paths_test.go`

**Test Cases:**
- Collection paths with various IDs and extensions
- Standalone paths
- Claude special paths (directory structure)
- Edge cases (empty strings, special characters)

**Acceptance Criteria:**
- [ ] File `internal/format/paths_test.go` exists
- [ ] Tests for all 4 path functions
- [ ] All examples from spec verified
- [ ] `go test ./internal/format` passes

**Estimated Effort:** 30 minutes

---

### Task 4.2: Write Unit Tests for Validation

**Context:**  
Verify validation functions reject invalid IDs and names.

**File:** `internal/format/validation_test.go`

**Test Cases:**
- Valid IDs (alphanumeric, hyphens, underscores)
- Invalid IDs (special characters, empty)
- Valid rule names (no parentheses)
- Invalid rule names (with parentheses)

**Acceptance Criteria:**
- [ ] File `internal/format/validation_test.go` exists
- [ ] Tests for `ValidateID()` and `ValidateRuleName()`
- [ ] All edge cases from spec covered
- [ ] `go test ./internal/format` passes

**Estimated Effort:** 30 minutes

---

### Task 4.3: Write Unit Tests for Metadata Generation

**Context:**  
Verify metadata blocks are correctly formatted.

**File:** `internal/format/metadata_test.go`

**Test Cases:**
- Ruleset with full metadata
- Ruleset with minimal metadata
- Standalone rule with full metadata
- Standalone rule with minimal metadata
- Enforcement header formatting

**Acceptance Criteria:**
- [ ] File `internal/format/metadata_test.go` exists
- [ ] Tests for both metadata generation functions
- [ ] YAML structure validated
- [ ] Enforcement header format validated
- [ ] `go test ./internal/format` passes

**Estimated Effort:** 45 minutes

---

### Task 4.4: Write Unit Tests for Target Compilers

**Context:**  
Verify each target compiler produces correct output.

**Files:**
- `pkg/targets/markdown_test.go`
- `pkg/targets/kiro_test.go`
- `pkg/targets/cursor_test.go`
- `pkg/targets/claude_test.go`
- `pkg/targets/copilot_test.go`

**Test Cases (per target):**
- Compile standalone rule
- Compile ruleset (expansion)
- Compile standalone prompt
- Compile promptset (expansion)
- Verify extensions
- Verify frontmatter (if applicable)
- Verify metadata blocks

**Acceptance Criteria:**
- [ ] Test file exists for each target
- [ ] All resource kinds tested
- [ ] Path structure verified
- [ ] Content format verified
- [ ] `go test ./pkg/targets` passes

**Estimated Effort:** 2 hours (all targets)

---

### Task 4.5: Write Integration Tests

**Context:**  
End-to-end compilation workflow tests.

**File:** `pkg/compiler/compiler_test.go`

**Test Cases:**
- Compile to single target
- Compile to multiple targets
- Error handling (invalid resource, unknown target, unsupported version)
- Result aggregation

**Acceptance Criteria:**
- [ ] File `pkg/compiler/compiler_test.go` exists
- [ ] Tests cover happy path and error cases
- [ ] Multi-target compilation verified
- [ ] `go test ./pkg/compiler` passes

**Estimated Effort:** 45 minutes

---

### Task 4.6: Write CLI Integration Tests

**Context:**  
Test CLI tool end-to-end.

**File:** `cmd/arc/arc_test.go`

**Test Cases:**
- Compile with stdout output
- Compile with file output
- Compile with multiple targets
- Error handling (missing file, no targets)
- Flag parsing

**Acceptance Criteria:**
- [ ] File `cmd/arc/arc_test.go` exists
- [ ] Tests cover stdout and file modes
- [ ] Error cases verified
- [ ] `go test ./cmd/arc` passes

**Estimated Effort:** 45 minutes

---

## Summary

**Total Tasks:** 23  
**Total Estimated Effort:** ~12 hours

**Phase Breakdown:**
- Phase 1 (Foundation): 6 tasks, ~2 hours
- Phase 2 (Compilation): 7 tasks, ~4.5 hours
- Phase 3 (Interface): 4 tasks, ~2 hours
- Phase 4 (Quality): 6 tasks, ~3.5 hours

**Critical Path:**
1. Task 1.1 → 1.2 → 1.3 → 1.4 → 1.5 → 1.6 (Foundation)
2. Task 2.1 → 2.2-2.6 → 2.7 (Compilation)
3. Task 3.1 → 3.2 → 3.3 → 3.4 (Interface)
4. Task 4.1-4.6 (Quality - can be parallelized)

**Success Criteria:**
- [ ] All 23 tasks completed
- [ ] `go test ./...` passes
- [ ] `go build ./...` succeeds
- [ ] `arc compile` works end-to-end
- [ ] Can compile to all 5 targets
- [ ] Documentation matches implementation
