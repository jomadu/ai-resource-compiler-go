# Specifications

This directory contains specifications for the AI Resource Compiler. Specifications are organized by Jobs to be Done (user outcomes) and Topics of Concern (organizing principles).

## Jobs to be Done (User Outcomes)

1. **Transform AI Resources to tool-specific formats** - Convert validated resources into formats that AI coding tools understand
2. **Maintain context in compiled output** - Preserve ruleset/rule relationships and enforcement levels in compiled files
3. **Support multiple AI coding tools** - Generate output for Kiro, Cursor, Claude, Copilot, and generic markdown
4. **Enable flexible integration** - Return modular results that users control (paths + content, no I/O)
5. **Provide command-line workflow** - Compile resources via CLI with target selection

## Topics of Concern

### Foundation Layer
Shared formats and architecture that all targets depend on.

- **[Metadata Embedding](metadata-block.md)** - How ruleset/rule context is embedded in compiled rule content
- **[Compilation Architecture](compiler-architecture.md)** - Core interfaces, data flow, and extension points

### Target Layer
Tool-specific format transformations.

- **[Markdown Compiler](markdown-compiler.md)** - Vanilla markdown output for generic use
- **[Kiro Compiler](kiro-compiler.md)** - Kiro CLI steering rules and prompts
- **[Cursor Compiler](cursor-compiler.md)** - Cursor IDE rules (.mdc) and commands (.md)
- **[Claude Compiler](claude-compiler.md)** - Claude Code rules and skills
- **[Copilot Compiler](copilot-compiler.md)** - GitHub Copilot instructions and prompts

### Interface Layer
User-facing interfaces for compilation workflow.

- **[CLI Design](cli-design.md)** - Command-line interface for compilation

## Specification Format

All specifications follow the structure defined in [TEMPLATE.md](TEMPLATE.md):

1. **Job to be Done** - User outcome, not mechanism
2. **Activities** - Discrete operations to accomplish the JTBD
3. **Acceptance Criteria** - Observable, testable outcomes
4. **Data Structures** - Types, interfaces, and data formats
5. **Algorithm** - Step-by-step logic with pseudocode
6. **Edge Cases** - Boundary conditions and error scenarios
7. **Dependencies** - Prerequisites and related components
8. **Implementation Mapping** - Source files and related specs
9. **Examples** - Concrete usage scenarios with verification
10. **Notes** - Design decisions and rationale
11. **Known Issues** - Discovered bugs or limitations
12. **Areas for Improvement** - Future enhancements

## Reading Order

### For Understanding the System

1. Start with **[Metadata Embedding](metadata-block.md)** - Understand the shared format
2. Read **[Compilation Architecture](compiler-architecture.md)** - Understand the core contracts
3. Pick a target compiler to see how it all fits together:
   - **[Markdown Compiler](markdown-compiler.md)** - Simplest example
   - **[Kiro Compiler](kiro-compiler.md)** - Similar to markdown
   - **[Cursor Compiler](cursor-compiler.md)** - Adds frontmatter
   - **[Claude Compiler](claude-compiler.md)** - Optional frontmatter, special prompt paths
   - **[Copilot Compiler](copilot-compiler.md)** - Frontmatter for rules and prompts
4. Read **[CLI Design](cli-design.md)** - Understand the user interface

### For Implementing a Feature

1. Find the relevant spec by JTBD or topic
2. Read the spec's dependencies first
3. Review examples and edge cases
4. Check implementation mapping for affected files
5. Follow acceptance criteria for testing

### For Adding a New Target

1. Read **[Metadata Embedding](metadata-block.md)** - Understand metadata blocks
2. Read **[Compilation Architecture](compiler-architecture.md)** - Understand TargetCompiler interface
3. Review existing target specs for patterns
4. Create new spec following [TEMPLATE.md](TEMPLATE.md)
5. Implement TargetCompiler interface
6. Register target in compiler

## Key Concepts

### Metadata Block
YAML block embedded in compiled rules (not prompts) that preserves ruleset and rule context. Includes:
- Ruleset section: id, name, description (optional), rules list
- Rule section: id, name, description (optional), enforcement, scope (optional)

### Enforcement Header
Header line following metadata block in rules: `# {Rule Name} ({ENFORCEMENT})`
- Enforcement level uppercased (MUST, SHOULD, MAY)
- Provides visual indication of rule importance

### Path Structure
Simple pattern for output paths: `{collection-id}_{item-id}.{ext}`
- Rules: Use ruleset ID and rule ID
- Prompts: Use promptset ID and prompt ID
- Exception: Claude prompts use `{collection-id}_{item-id}/SKILL.md`

### Target Compiler
Interface for implementing target-specific compilation:
```go
type TargetCompiler interface {
    Name() string
    Compile(resource Resource) ([]CompilationResult, error)
}
```

### Compilation Result
Modular output (path + content) that users control:
```go
type CompilationResult struct {
    Path    string
    Content string
}
```

### Pure Transformation
Compiler produces path + content pairs, doesn't perform I/O. Users decide where and how to write files.

## Target Comparison

| Target | Rules Extension | Prompts Extension | Frontmatter | Metadata Block | Installation |
|--------|----------------|-------------------|-------------|----------------|--------------|
| Markdown | .md | .md | None | Rules only | User choice |
| Kiro | .md | .md | None | Rules only | .kiro/steering/, .kiro/prompts/ |
| Cursor | .mdc | .md | Rules only (MDC) | Rules only | .cursor/rules/, .cursor/commands/ |
| Claude | .md | SKILL.md | Rules only (paths, optional) | Rules only | .claude/rules/, .claude/skills/ |
| Copilot | .instructions.md | .prompt.md | Both (applyTo) | Rules only | .github/instructions/, .github/prompts/ |

## Contributing

When adding or modifying specifications:

1. Follow [TEMPLATE.md](TEMPLATE.md) structure
2. Lead with JTBD (user outcome)
3. Provide concrete, testable examples
4. Show metadata blocks in all rule examples
5. Show enforcement headers in all rule examples
6. Reference dependencies explicitly
7. Update this README if adding new specs

## Questions?

- **What's the difference between rules and prompts?** Rules have metadata blocks and enforcement headers. Prompts are just body content (except Copilot prompts have frontmatter).
- **Why metadata blocks?** Preserves context for tools and humans. Enables understanding of rule origin and importance.
- **Why different extensions per target?** Each tool has conventions. We follow them for seamless integration.
- **Why no I/O in compiler?** Pure transformation enables flexible integration. Users control where files go.
- **How do I add a new target?** Implement TargetCompiler interface, create spec, register in compiler.
