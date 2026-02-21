# AI Resource Compiler (Go)

Compile [AI Resource Specification](https://github.com/jomadu/ai-resource-spec) resources to tool-specific formats.

## Overview

`ai-resource-compiler-go` transforms validated AI resources into formats consumed by AI coding tools:

- **Kiro CLI** - AWS AI assistant format
- **Cursor** - `.cursor/rules/*.mdc` formats
- **Claude Code** - `.claude/rules/*.md` and `.claude/skills/*/SKILL.md` formats
- **GitHub Copilot** - `.github/instructions/*.instructions.md` and `.github/prompts/*.prompt.md` formats

Built on [ai-resource-core-go](https://github.com/jomadu/ai-resource-core-go) for parsing and validation.

## Design Philosophy

The compiler is a **pure transformation tool**:
- Transforms resource content to target-specific formats
- Returns relative paths and compiled content
- **Does not** manage file I/O or directory structures
- Users control where and how files are written

This design enables flexible integration with build tools, CI/CD pipelines, and custom workflows.

## Installation

```bash
go get github.com/jomadu/ai-resource-compiler-go
```

## Usage

```go
import (
    "github.com/jomadu/ai-resource-core-go/pkg/airesource"
    "github.com/jomadu/ai-resource-compiler-go/pkg/compiler"
)

// Load resource
prompt, _ := airesource.LoadPrompt("prompt.yml")

// Compile to Cursor format
c := compiler.New()
results, _ := c.Compile(prompt, compiler.Options{
    Targets: []string{"cursor"},
})

// User decides where to write
for _, result := range results {
    // Prepend recommended directory for target
    fullPath := filepath.Join(".cursor/rules", result.Path)
    os.MkdirAll(filepath.Dir(fullPath), 0755)
    os.WriteFile(fullPath, result.Content, 0644)
}
```

## CLI

```bash
# Install
go install github.com/jomadu/ai-resource-compiler-go/cmd/arc@latest

# Compile to Cursor
arc compile --target cursor -o .cursor/rules prompts.yml

# Compile to multiple targets
arc compile --target kiro,cursor,claude -o ./build prompts.yml

# Output to stdout
arc compile --target cursor prompts.yml
```

## Supported Targets

| Target | Rules | Prompts | Notes |
|--------|-------|---------|-------|
| `kiro` | `{id}.md` | `{id}.md` | Plain markdown |
| `cursor` | `{id}.mdc` | `{id}.mdc` | MDC format with frontmatter |
| `claude` | `{id}.md` | `{id}/SKILL.md` | Rules as .md, prompts as directories |
| `copilot` | `{id}.instructions.md` | `{id}.prompt.md` | Markdown with frontmatter |

## Compilation Results

The compiler returns relative paths and content for modular output:

```go
type CompilationResult struct {
    Path    string  // Relative path: "api-standards.md", "deploy/SKILL.md", etc.
    Content []byte  // Compiled content
}
```

**Path Examples:**
- Single rule `id: api-standards` → `api-standards.md` (Kiro/Claude), `api-standards.mdc` (Cursor), `api-standards.instructions.md` (Copilot)
- Ruleset `id: backend` with rule `id: api` → `backend_api.md` (Kiro/Claude), `backend_api.mdc` (Cursor), `backend_api.instructions.md` (Copilot)
- Single prompt `id: deploy` → `deploy.md` (Kiro/Cursor), `deploy/SKILL.md` (Claude), `deploy.prompt.md` (Copilot)
- Promptset `id: ci-workflows` with prompt `id: deploy` → `ci-workflows_deploy.md` (Kiro/Cursor), `ci-workflows_deploy/SKILL.md` (Claude), `ci-workflows_deploy.prompt.md` (Copilot)

**Responsibility Split:**
- Compiler returns relative paths (e.g., `api-standards.md`)
- Users prepend target-specific directories when writing files

## Recommended Installation Locations

| Target   | Rules                  | Prompts               |
|----------|------------------------|-----------------------|
| Kiro     | `.kiro/steering/`      | `.kiro/prompts/`      |
| Cursor   | `.cursor/rules/`       | `.cursor/commands/`   |
| Claude   | `.claude/rules/`       | `.claude/skills/`     |
| Copilot  | `.github/instructions/`| `.github/prompts/`    |

## Architecture

- **Core** ([ai-resource-core-go](https://github.com/jomadu/ai-resource-core-go)) - Parsing and validation
- **Compiler** (this repo) - Format transformation
- **Target Compilers** - Tool-specific output generation

Clean separation: validation → transformation → user-controlled I/O.

## License

GNU General Public License v2.0