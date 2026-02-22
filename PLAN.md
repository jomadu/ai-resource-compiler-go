# Rewrite Outline

**Status:** Published to TODO.md on 2026-02-22  
**Tasks Created:** TASK-001 through TASK-010

## Instructions for Implementation

**This outline is the authoritative reference** for rewriting the specs and README.

**Purpose:** Construct new specification documents and README based on this outline.

**Reference materials:**
- `archive/` - Historical specs and README (possibly inaccurate, use for reference only)
- `specs/TEMPLATE.md` - Structure template for all specs
- `PLAN.md` - Original requirements that led to this outline

**Approach:**
1. Follow this outline's structure and API design exactly
2. Use `specs/TEMPLATE.md` for consistent spec formatting
3. Reference `archive/` for examples and patterns, but verify against this outline
4. When conflicts arise, this outline takes precedence

---

## Jobs to be Done (User Outcomes)

1. **Transform AI Resources to tool-specific formats** - Convert validated resources into formats that AI coding tools understand
2. **Maintain context in compiled output** - Preserve ruleset/rule relationships and enforcement levels in compiled files
3. **Support multiple AI coding tools** - Generate output for Kiro, Cursor, Claude, Copilot, and generic markdown
4. **Enable flexible integration** - Return modular results that users control (paths + content, no I/O)
5. **Provide command-line workflow** - Compile resources via CLI with target selection

## Topics of Concern (Organizing Principles)

### Foundation Layer
- **Metadata Embedding** - How ruleset/rule context is embedded in compiled rule content
- **Compilation Architecture** - Core interfaces, data flow, and extension points

### Target Layer  
- **Target Compilers** - Tool-specific format transformations (Kiro, Cursor, Claude, Copilot, Markdown)

### Interface Layer
- **CLI** - Command-line interface for compilation workflow

## Specification Structure

### Foundation Specs

#### `specs/metadata-block.md`
**JTBD:** Preserve ruleset/rule context in compiled output

**Key Activities:**
- Define YAML metadata structure (ruleset, rule sections - no namespace)
- Specify when metadata is included (rules yes, prompts no)
- Define enforcement header format ("# Rule Name (MUST)")

**Full Metadata Block Example:**
```yaml
---
ruleset:
  id: cleanCode
  name: Clean Code
  description: Clean code practices
  rules:
    - meaningfulNames
    - smallFunctions
    - singleResponsibility
rule:
  id: meaningfulNames
  name: Use Meaningful Names
  description: Variables and functions should have descriptive names
  enforcement: must
  scope:
    - files: ["**/*.ts", "**/*.js"]
---

# Use Meaningful Names (MUST)

Variables and functions should have descriptive names that reveal intent.
```

**Minimal Metadata Block Example (without optional fields):**
```yaml
---
ruleset:
  id: cleanCode
  name: Clean Code
  rules:
    - meaningfulNames
    - smallFunctions
rule:
  id: meaningfulNames
  name: Use Meaningful Names
  enforcement: must
---

# Use Meaningful Names (MUST)

Rule body content.
```

**Note:** Promptsets do NOT include metadata blocks - just the prompt body.

**Why First:** All target compilers depend on this shared format

#### `specs/compiler-architecture.md`
**JTBD:** Provide extensible architecture for multi-target compilation

**Key Activities:**
- Define TargetCompiler interface with single Compile method
- Define Target enum for type-safe target selection
- Specify CompileOptions with targets list
- Define CompilationResult structure (path + content)
- Establish compilation pipeline and target registration

**API:**
```go
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

type TargetCompiler interface {
    Name() string
    Compile(resource Resource) ([]CompilationResult, error)
}

func (c *Compiler) Compile(resource Resource, opts CompileOptions) ([]CompileResult, error)
```

**Path Structure:**
- Rules: `{collection-id}_{item-id}.{ext}`
- Prompts: `{collection-id}_{item-id}.{ext}`
- Claude prompts: `{collection-id}_{item-id}/SKILL.md`

**Why Second:** Defines contracts that all targets implement

### Target Compiler Specs

#### `specs/markdown-compiler.md`
**JTBD:** Generate vanilla markdown output for generic use

**Key Activities:**
- Rules: Metadata block + enforcement header + body
- Prompts: Plain body only
- Naming: `{collection-id}_{item-id}.md`

**Format:**
- No frontmatter
- Rules: Metadata block only
- Prompts: Plain body

**Why First Target:** Simplest format, demonstrates metadata block usage

#### `specs/kiro-compiler.md`
**JTBD:** Generate Kiro CLI steering rules and prompts

**Key Activities:**
- Rules: Metadata block + enforcement header + body
- Prompts: Plain body only
- Extension: `.md`
- Paths: `{collection-id}_{item-id}.md`
- Installation: `.kiro/steering/`, `.kiro/prompts/`

**Format:**
- No frontmatter
- Rules: Metadata block only
- Prompts: Plain body

#### `specs/cursor-compiler.md`
**JTBD:** Generate Cursor IDE rules and commands

**Key Activities:**
- Rules: MDC frontmatter + metadata block + enforcement header + body
- Prompts: Plain body (no frontmatter, no metadata)
- Extension: `.mdc` (rules), `.md` (prompts)
- Paths: `{collection-id}_{item-id}.{ext}`
- Installation: `.cursor/rules/`, `.cursor/commands/`

**Frontmatter (Rules only):**
```yaml
---
description: string       # Rule description
globs: []string          # File patterns from scope
alwaysApply: bool        # true for must enforcement
---
```

#### `specs/claude-compiler.md`
**JTBD:** Generate Claude Code rules and skills

**Key Activities:**
- Rules: Optional paths frontmatter + metadata block + enforcement header + body (`.md`)
- Prompts: Plain body (`{collection-id}_{item-id}/SKILL.md`)
- Paths: `{collection-id}_{item-id}.md` (rules), `{collection-id}_{item-id}/SKILL.md` (prompts)
- Installation: `.claude/rules/`, `.claude/skills/`

**Frontmatter (Rules only, optional):**
```yaml
---
paths:
  - string  # File patterns from scope
---
```

**Note:** Prompts have no frontmatter, no metadata - just body in SKILL.md

#### `specs/copilot-compiler.md`
**JTBD:** Generate GitHub Copilot instructions and prompts

**Key Activities:**
- Rules: applyTo frontmatter + metadata block + enforcement header + body (`.instructions.md`)
- Prompts: applyTo frontmatter + body (`.prompt.md`)
- Paths: `{collection-id}_{item-id}.{ext}`
- Installation: `.github/instructions/`, `.github/prompts/`

**Frontmatter (Rules and Prompts):**
```yaml
---
applyTo: []string        # File patterns from scope
---
```

**Note:** excludeAgent field is omitted (not populated by compiler)

### Interface Specs

#### `specs/cli-design.md`
**JTBD:** Provide command-line workflow for compilation

**Key Activities:**
- `arc compile` command with `--target` and `--output` flags
- Multi-file and multi-target support
- Stdout vs file output modes

**Why Last:** Depends on all compiler specs being defined

### Index Spec

#### `specs/README.md`
**JTBD:** Navigate specification documents by concern

**Structure:**
- JTBDs (user outcomes)
- Topics of Concern (organizing principles)
- Foundation specs
- Target compiler specs
- Interface specs

## README Structure

### `README.md`
**JTBD:** Help users understand and adopt the compiler

**Structure:**
1. **Overview** - What it does, supported targets
2. **Design Philosophy** - Pure transformation, modular output, user-controlled I/O
3. **Installation** - Go get, CLI install
4. **Usage** - Library API examples
5. **CLI** - Command examples
6. **Supported Targets** - Table with extensions and notes
7. **Compilation Results** - Path examples, responsibility split
8. **Recommended Locations** - Installation directories per target
9. **Architecture** - High-level component diagram
10. **License**

**Key Updates:**
- Add markdown to targets table
- Show metadata block structure (without namespace)
- Explain enforcement headers in rules
- Show simple path structure: `{collection-id}_{item-id}.{ext}`
- Emphasize modular output design

## Writing Principles

### Specs
- Follow TEMPLATE.md structure rigorously
- Lead with JTBD (user outcome, not mechanism)
- Provide concrete, testable examples
- Show metadata blocks in all rule examples (without namespace field)
- Show enforcement headers in all rule examples
- Show simple path structure: `{collection-id}_{item-id}.{ext}`
- Reference dependencies explicitly

### README
- User-focused, outcome-oriented
- Show practical examples first
- Explain design decisions briefly
- Link to specs for details
- Maintain "pure transformation" philosophy
- Keep it simple - no namespace complexity

## Execution Order

1. **Foundation**
   - `specs/metadata-block.md` (shared format)
   - `specs/compiler-architecture.md` (core contracts)

2. **Targets** (parallel, any order)
   - `specs/markdown-compiler.md` (simplest)
   - `specs/kiro-compiler.md`
   - `specs/cursor-compiler.md`
   - `specs/claude-compiler.md`
   - `specs/copilot-compiler.md`

3. **Interface**
   - `specs/cli-design.md` (depends on all targets)

4. **Navigation**
   - `specs/README.md` (index of all specs)

5. **User Documentation**
   - `README.md` (user-facing overview)
