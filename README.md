# AI Resource Compiler (Go)

Transform AI resources into tool-specific formats for Cursor, Kiro, Claude, Copilot, and more.

## Overview

The AI Resource Compiler takes validated AI resources (rules and prompts) and compiles them into formats that AI coding tools understand. It preserves context through metadata blocks and provides flexible output options.

**Supported Targets:**
- **Markdown** - Generic markdown for any tool
- **Kiro** - Kiro CLI steering rules and prompts
- **Cursor** - Cursor IDE rules (.mdc) and commands (.md)
- **Claude** - Claude Code rules and skills
- **Copilot** - GitHub Copilot instructions and prompts

## Design Philosophy

1. **Pure Transformation** - Compiler produces path + content pairs, doesn't perform I/O
2. **Modular Output** - Users control where and how to write files
3. **Context Preservation** - Metadata blocks maintain ruleset/rule relationships
4. **Extensible** - Add new targets via TargetCompiler interface

## Installation

```bash
go get github.com/jomadu/ai-resource-compiler-go
```

Install CLI tool:

```bash
go install github.com/jomadu/ai-resource-compiler-go/cmd/arc@latest
```

## Usage

### Library API

```go
import (
    "github.com/jomadu/ai-resource-compiler-go/pkg/compiler"
)

// Create compiler
c := compiler.NewCompiler()

// Compile to single target
opts := compiler.CompileOptions{
    Targets: []compiler.Target{compiler.TargetMarkdown},
}
results, err := c.Compile(resource, opts)

// Compile to multiple targets
opts := compiler.CompileOptions{
    Targets: []compiler.Target{
        compiler.TargetKiro,
        compiler.TargetCursor,
        compiler.TargetClaude,
    },
}
results, err := c.Compile(resource, opts)

// Handle results
for _, result := range results {
    fmt.Printf("Path: %s\n", result.Path)
    fmt.Printf("Content:\n%s\n", result.Content)
}

// Compile multiple resources from file
// Note: Compile() accepts a single resource. Iterate for multiple resources.
resources, err := core.LoadResources("resources.yaml")
if err != nil {
    log.Fatal(err)
}

for _, resource := range resources {
    results, err := c.Compile(resource, opts)
    if err != nil {
        log.Printf("Failed to compile resource: %v", err)
        continue
    }
    // Handle results for each resource
}
```

### CLI

Compile to markdown, print to stdout:

```bash
arc compile resource.yaml --target markdown
```

Compile to multiple targets, print to stdout:

```bash
arc compile resource.yaml --target markdown --target kiro --target cursor
```

Compile to cursor, write to directory:

```bash
arc compile resource.yaml --target cursor --output .cursor/rules
```

Compile to all targets, write to directory:

```bash
arc compile resource.yaml \
  --target cursor \
  --target kiro \
  --target claude \
  --target copilot \
  --target markdown \
  --output ./output
```

## Supported Targets

| Target | Rules Extension | Prompts Extension | Frontmatter | Metadata Block | Notes |
|--------|----------------|-------------------|-------------|----------------|-------|
| markdown | .md | .md | None | Rules only | Generic markdown |
| kiro | .md | .md | None | Rules only | Kiro CLI format |
| cursor | .mdc | .md | Rules only | Rules only | MDC frontmatter |
| claude | .md | SKILL.md | Rules only (optional) | Rules only | Directory for prompts |
| copilot | .instructions.md | .prompt.md | Both | Rules only | applyTo frontmatter |

## Compilation Results

The compiler returns `CompilationResult` structs with path and content:

```go
type CompilationResult struct {
    Path    string  // e.g., "cleanCode_meaningfulNames.md"
    Content string  // Compiled content
}
```

**Path Structure:**
- Rules: `{ruleset-id}_{rule-id}.{ext}`
- Prompts: `{promptset-id}_{prompt-id}.{ext}`
- Claude prompts: `{promptset-id}_{prompt-id}/SKILL.md`

**Your Responsibility:**
- Decide where to write files
- Create directories as needed
- Handle file conflicts
- Manage permissions

**Our Responsibility:**
- Generate correct paths
- Produce valid content
- Preserve context via metadata

## Recommended Locations

Where to install compiled files for each tool:

| Target | Rules | Prompts |
|--------|-------|---------|
| kiro | `.kiro/steering/` | `.kiro/prompts/` |
| cursor | `.cursor/rules/` | `.cursor/commands/` |
| claude | `.claude/rules/` | `.claude/skills/` |
| copilot | `.github/instructions/` | `.github/prompts/` |
| markdown | User choice | User choice |

## Metadata Block Structure

Rules include YAML metadata blocks that preserve context:

```yaml
---
ruleset:
  id: cleanCode
  name: Clean Code
  description: Clean code practices
  rules:
    - meaningfulNames
    - smallFunctions
rule:
  id: meaningfulNames
  name: Use Meaningful Names
  description: Variables and functions should have descriptive names
  enforcement: must
  scope:
    files:
      - "**/*.ts"
      - "**/*.js"
---

# Use Meaningful Names (MUST)

Variables and functions should have descriptive names that reveal intent.
```

**Key Features:**
- Ruleset context (id, name, description, rules list)
- Rule context (id, name, description, enforcement, scope)
- Enforcement header (`# {Name} ({ENFORCEMENT})`)
- Optional fields omitted when not present

**Prompts:**
Prompts do NOT include metadata blocks - just body content (except Copilot prompts have frontmatter).

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         User Code                           │
│                                                             │
│  compiler.Compile(resource, opts) → []CompilationResult    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Compiler (pkg/compiler)                  │
│                                                             │
│  • Validates options                                        │
│  • Routes to target compilers                               │
│  • Aggregates results                                       │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│              Target Compilers (pkg/targets)                 │
│                                                             │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  │
│  │ Markdown │  │   Kiro   │  │  Cursor  │  │  Claude  │  │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘  │
│                                                             │
│  ┌──────────┐                                              │
│  │ Copilot  │                                              │
│  └──────────┘                                              │
│                                                             │
│  Each implements: TargetCompiler interface                 │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│           Metadata Generation (internal/format)             │
│                                                             │
│  • Builds YAML metadata blocks                              │
│  • Generates enforcement headers                            │
│  • Handles optional fields                                  │
└─────────────────────────────────────────────────────────────┘
```

**Extension Points:**
- Implement `TargetCompiler` interface for new targets
- Register custom compilers via `RegisterTarget()`
- Reuse metadata generation for consistency

## Development

This project depends on [ai-resource-core-go](https://github.com/jomadu/ai-resource-core-go) for resource loading and validation.

**Project Structure:**
```
ai-resource-compiler-go/
├── cmd/arc/              # CLI tool
├── pkg/
│   ├── compiler/         # Public API
│   └── targets/          # Target compilers
├── internal/format/      # Metadata generation
├── specs/                # Specifications
└── README.md
```

**Build:**
```bash
go build ./...
```

**Test:**
```bash
go test ./...
```

**Install CLI:**
```bash
go install ./cmd/arc
```

## Specifications

Detailed specifications are in the [specs/](specs/) directory:

- **Foundation:** [Metadata Block](specs/metadata-block.md), [Compiler Architecture](specs/compiler-architecture.md)
- **Targets:** [Markdown](specs/markdown-compiler.md), [Kiro](specs/kiro-compiler.md), [Cursor](specs/cursor-compiler.md), [Claude](specs/claude-compiler.md), [Copilot](specs/copilot-compiler.md)
- **Interface:** [CLI Design](specs/cli-design.md)

See [specs/README.md](specs/README.md) for reading order and key concepts.

## License

Apache License 2.0 - See [LICENSE](LICENSE) for details.
