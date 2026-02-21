# AGENTS.md

## Work Tracking System

Tasks are tracked in `TODO.md` at repository root.

Task format:
```markdown
## TASK-001
- Priority: 1-5 (1=highest)
- Status: TODO/IN_PROGRESS/BLOCKED/DONE
- Dependencies: [TASK-XXX, ...]
- Description: Task description
```

Manual editing. Tasks auto-increment. Keep all tasks (including DONE) in file.

## Feature Input

`TASK.md` contains feature requirements and specifications for the compiler.

## Quick Reference

- Edit `TODO.md` - Manage tasks
- `go test ./...` - Run tests (when initialized)
- `go build ./...` - Build packages (when initialized)
- `go install ./cmd/arc` - Install CLI tool

## Planning System

`PLAN.md` documents the current plan (not yet created).

## Build/Test/Lint Commands

Go project (not yet initialized):
- `go mod init` - Initialize module
- `go test ./...` - Run tests
- `go build ./...` - Build packages
- `go vet ./...` - Lint code
- `go install ./cmd/arc` - Install arc CLI

## Specification Definition

Specifications live in `specs/`. A spec template exists at `specs/TEMPLATE.md`.

Format: Markdown with structured sections following JTBD → Activities → Acceptance Criteria → Implementation pattern.

## Implementation Definition

Location: `pkg/compiler/` (public API), `internal/` (private implementation)

Patterns:
- `pkg/compiler/*.go` - Public compilation API
- `pkg/targets/*.go` - Target-specific compilers (kiro, cursor, claude, copilot)
- `internal/format/*.go` - Formatting utilities
- `cmd/arc/*.go` - CLI tool implementation

Excludes: `testdata/`, `.git/`

## Audit Output

Audit results written to `AUDIT.md` at repository root.

## Quality Criteria

**Specifications:**
- All requirements testable
- Examples provided for each target format
- Implementation notes clear
- Target format specifications documented

**Implementation:**
- Passes `go test ./...`
- Passes `go vet ./...`
- Public API minimal and documented
- Each target compiler is independent
- CLI follows standard conventions

**Refactoring triggers:**
- Spec/implementation divergence
- Test failures
- Target format changes
- New target addition

## Operational Learnings

Last verified: 2026-02-20

**Working:**
- README defines project vision
- Clear separation from ai-resource-core-go

**Not working:**
- Go module not initialized (no go.mod)
- No specs directory
- No implementation directories
- No CLI implementation

**Rationale:**
- Project in bootstrap phase
- Specifications needed before implementation
- Depends on ai-resource-core-go for resource loading
