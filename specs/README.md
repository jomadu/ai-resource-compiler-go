# AI Resource Compiler Specifications Index

## Jobs to be Done (JTBDs)

1. **Compile AI Resources to modular target formats** - Transform validated AI Resources into tool-specific modular file structures
2. **Support multiple compilation targets** - Generate output for Kiro CLI, Cursor, Claude Code, and GitHub Copilot
3. **Preserve semantic meaning** - Ensure compiled output maintains the intent of the original resource
4. **Provide CLI interface** - Command-line tool for compiling resources with target selection
5. **Handle multi-resource bundles** - Compile multiple resources from single or multi-document files
6. **Generate idiomatic output** - Produce modular output that follows each target's conventions and best practices

## Topics of Concern

### Compilation Pipeline
- **Resource Loading** - Load and validate resources using ai-resource-core-go
- **Target Selection** - Choose which target formats to compile to
- **Format Transformation** - Convert resource structure to target-specific format
- **Modular Output Generation** - Return relative paths and content for user-controlled file writing

### Target Formats
- **Kiro CLI** - Modular markdown format (.kiro/steering/, .kiro/prompts/)
- **Cursor** - Modular MDC format (.cursor/rules/, .cursor/commands/)
- **Claude Code** - Modular format (.claude/rules/, .claude/skills/)
- **GitHub Copilot** - Modular format (.github/instructions/, .github/prompts/)

### CLI Tool
- **Command Structure** - `arc compile` with target and file options
- **Flag Handling** - Target selection, output paths, multi-target compilation
- **Error Reporting** - Clear messages for compilation failures
- **Help System** - Usage documentation and examples

### Format Conventions
- **Prompt Formatting** - How prompts are represented in each target
- **Rule Formatting** - How rules are represented in each target
- **Fragment Resolution** - When to resolve fragments vs preserve structure
- **Collection Handling** - How Promptsets and Rulesets are named and organized

## Specification Documents
 
### Foundation
- [compiler-architecture.md](compiler-architecture.md) - Overall compilation system design

### CLI
- [cli-design.md](cli-design.md) - Command-line interface specification

### Target Compilers
- [kiro-compiler.md](kiro-compiler.md) - Kiro CLI format compilation
- [cursor-compiler.md](cursor-compiler.md) - Cursor format compilation
- [claude-compiler.md](claude-compiler.md) - Claude Code format compilation
- [copilot-compiler.md](copilot-compiler.md) - GitHub Copilot format compilation
