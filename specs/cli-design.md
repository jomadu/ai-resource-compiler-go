# CLI Design

## Job to be Done
Provide a command-line interface for compiling AI resources to target-specific formats with flexible output options.

## Activities
1. Define `arc compile` command with target selection
2. Support multiple targets in single invocation
3. Support stdout output mode (print to console)
4. Support file output mode (write to filesystem)
5. Provide clear error messages and usage help

## Acceptance Criteria
- [ ] `arc compile` command accepts resource file path
- [ ] `--target` flag accepts multiple values (cursor, kiro, claude, copilot, markdown)
- [ ] `--output` flag accepts "stdout" or directory path
- [ ] Default output is stdout
- [ ] Stdout mode prints path and content for each result
- [ ] File mode writes results to specified directory
- [ ] Error messages are clear and actionable
- [ ] `--help` flag shows usage information
- [ ] Exit code 0 on success, non-zero on error

## Data Structures

### Command Structure
```
arc compile <resource-file> [flags]
```

**Arguments:**
- `<resource-file>` - Path to resource file (YAML or JSON)

**Flags:**
- `--target, -t` - Target format(s) to compile to (repeatable)
- `--output, -o` - Output mode: "stdout" or directory path (default: "stdout")
- `--flat` - Disable target subdirectories in file output mode
- `--help, -h` - Show help information

### Output Modes

**Stdout Mode:**
```
=== {target}/{path} ===
{content}

=== {target}/{path} ===
{content}
```

**File Mode:**
- Default: Writes each CompilationResult to `{output-dir}/{target}/{path}`
- With `--flat`: Writes to `{output-dir}/{path}` (no target subdirectories)
- Creates directories as needed
- Reports files written to stderr
- Target subdirectories prevent filename collisions when compiling to multiple targets

## Algorithm

1. Parse command-line arguments
2. Validate resource file exists
3. Validate at least one target specified
4. Load resource from file
5. Create compiler with requested targets
6. Compile resource
7. If stdout mode:
   - Print each result with path separator
8. If file mode:
   - Write each result to {output-dir}/{path}
   - Create directories as needed
   - Report files written
9. Exit with appropriate code

**Pseudocode:**
```
function main():
    args = parseArgs()
    
    if args.help:
        printHelp()
        exit(0)
    
    if args.resourceFile is empty:
        printError("resource file required")
        exit(1)
    
    if args.targets is empty:
        printError("at least one target required")
        exit(1)
    
    resource = loadResource(args.resourceFile)
    compiler = NewCompiler()
    
    opts = CompileOptions{Targets: args.targets}
    results = compiler.Compile(resource, opts)
    
    if args.output == "stdout":
        for result in results:
            print("=== " + result.Target + "/" + result.Path + " ===")
            print(result.Content)
            print()
    else:
        for result in results:
            if args.flat:
                path = args.output + "/" + result.Path
            else:
                path = args.output + "/" + result.Target + "/" + result.Path
            writeFile(path, result.Content)
            printStderr("Wrote " + path)
    
    exit(0)
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| No resource file specified | Print error, show usage, exit 1 |
| Resource file not found | Print error with path, exit 1 |
| No targets specified | Print error, show usage, exit 1 |
| Invalid target name | Print error with valid options, exit 1 |
| Compilation error | Print error message, exit 1 |
| Output directory doesn't exist | Create directory (including target subdirs), write files |
| File write permission error | Print error, exit 1 |
| Multiple targets, stdout mode | Print all results sequentially |
| Multiple targets, file mode (no --flat) | Write to separate target subdirectories |
| Multiple targets, file mode (--flat) | Write to same directory (last target wins on collision) |
| Single target, file mode (no --flat) | Create target subdirectory |
| Single target, file mode (--flat) | Write directly to output directory |

## Dependencies

- Compiler from compiler-architecture.md spec
- Resource loading from ai-resource-core-go
- CLI framework (e.g., cobra, flag)

## Implementation Mapping

**Source files:**
- `cmd/arc/main.go` - CLI entry point
- `cmd/arc/compile.go` - Compile command implementation
- `cmd/arc/output.go` - Output mode handling

**Related specs:**
- `compiler-architecture.md` - Compiler interface and types
- All target compiler specs - Available compilation targets

## Examples

### Example 1: Single Target, Stdout

**Command:**
```bash
arc compile resource.yaml --target markdown
```

**Output:**
```
=== markdown/cleanCode_meaningfulNames.md ===
---
ruleset:
  id: cleanCode
  name: Clean Code
  rules:
    - meaningfulNames
rule:
  id: meaningfulNames
  name: Use Meaningful Names
  enforcement: must
---

# Use Meaningful Names (MUST)

Use descriptive variable and function names.
```

**Verification:**
- Path printed with separator
- Content printed below
- Exit code 0

### Example 2: Multiple Targets, Stdout

**Command:**
```bash
arc compile resource.yaml --target markdown --target kiro --target cursor
```

**Output:**
```
=== markdown/cleanCode_meaningfulNames.md ===
[markdown content]

=== kiro/cleanCode_meaningfulNames.md ===
[kiro content]

=== cursor/cleanCode_meaningfulNames.mdc ===
[cursor content]
```

**Verification:**
- Three results printed
- Each with path separator
- Exit code 0

### Example 3: Single Target, File Output

**Command:**
```bash
arc compile resource.yaml --target cursor --output .cursor/rules
```

**Output (stderr):**
```
Wrote .cursor/rules/cursor/cleanCode_meaningfulNames.mdc
```

**Verification:**
- File written to target subdirectory
- Path reported to stderr
- Exit code 0

### Example 4: Multiple Targets, File Output

**Command:**
```bash
arc compile resource.yaml --target markdown --target kiro --output ./output
```

**Output (stderr):**
```
Wrote ./output/markdown/cleanCode_meaningfulNames.md
Wrote ./output/kiro/cleanCode_meaningfulNames.md
```

**Verification:**
- Two files written to separate target subdirectories
- No filename collision
- Paths reported to stderr
- Exit code 0

### Example 4a: Single Target, File Output with --flat

**Command:**
```bash
arc compile resource.yaml --target cursor --output .cursor/rules --flat
```

**Output (stderr):**
```
Wrote .cursor/rules/cleanCode_meaningfulNames.mdc
```

**Verification:**
- File written directly to output directory (no target subdirectory)
- Path reported to stderr
- Exit code 0

### Example 4b: Multiple Targets, File Output with --flat

**Command:**
```bash
arc compile resource.yaml --target markdown --target kiro --output ./output --flat
```

**Output (stderr):**
```
Wrote ./output/cleanCode_meaningfulNames.md
Wrote ./output/cleanCode_meaningfulNames.md
```

**Verification:**
- Two files written to same directory
- Same filename causes overwrite (kiro overwrites markdown)
- Paths reported to stderr
- Exit code 0

### Example 5: Error - No Targets

**Command:**
```bash
arc compile resource.yaml
```

**Output (stderr):**
```
Error: at least one target required

Usage:
  arc compile <resource-file> [flags]

Flags:
  -t, --target string   Target format (cursor, kiro, claude, copilot, markdown)
  -o, --output string   Output mode: stdout or directory path (default "stdout")
  -h, --help           Show help
```

**Verification:**
- Error message printed
- Usage shown
- Exit code 1

### Example 6: Error - Invalid Target

**Command:**
```bash
arc compile resource.yaml --target invalid
```

**Output (stderr):**
```
Error: unknown target: invalid

Valid targets: cursor, kiro, claude, copilot, markdown
```

**Verification:**
- Error message with invalid target
- Valid options listed
- Exit code 1

### Example 7: Help

**Command:**
```bash
arc compile --help
```

**Output:**
```
Compile AI resources to target-specific formats

Usage:
  arc compile <resource-file> [flags]

Arguments:
  resource-file    Path to resource file (YAML or JSON)

Flags:
  -t, --target string   Target format to compile to (repeatable)
                        Valid targets: cursor, kiro, claude, copilot, markdown
  -o, --output string   Output mode: "stdout" or directory path (default "stdout")
      --flat            Disable target subdirectories in file output mode
  -h, --help           Show this help message

Examples:
  # Compile to markdown, print to stdout
  arc compile resource.yaml --target markdown

  # Compile to multiple targets, print to stdout
  arc compile resource.yaml --target markdown --target kiro

  # Compile to cursor, write to target subdirectory
  arc compile resource.yaml --target cursor --output .cursor/rules

  # Compile to cursor, write directly to directory (no subdirectory)
  arc compile resource.yaml --target cursor --output .cursor/rules --flat

  # Compile to all targets, write to separate subdirectories
  arc compile resource.yaml -t cursor -t kiro -t claude -t copilot -t markdown -o ./output
```

**Verification:**
- Usage information shown
- Examples provided
- Exit code 0

## Notes

**Design Rationale:**
- **Stdout default** - Enables piping and inspection without filesystem changes
- **File mode** - Convenient for direct installation to tool directories
- **Target subdirectories** - Prevents filename collisions when compiling to multiple targets
- **--flat flag** - Allows direct writes for single-target workflows (e.g., `.cursor/rules/`)
- **Multiple targets** - Compile once, output to multiple formats
- **Simple flags** - Minimal learning curve, standard CLI conventions
- **Clear errors** - Help users fix issues quickly

**Use Cases:**
- **Development** - Stdout mode for quick inspection
- **CI/CD** - File mode for automated deployment
- **Multi-tool** - Compile to multiple targets in one command
- **Testing** - Stdout mode for validation and diffing

**Future Enhancements:**
- Batch compilation (multiple resource files)
- Watch mode (recompile on file changes)
- Validation mode (check without compiling)
- Template support (custom output formats)

## Known Issues

None - this is a new specification.

## Areas for Improvement

- Consider adding `--watch` flag for development workflow
- Evaluate `--validate` flag for syntax checking
- Explore `--format` flag for custom output templates
- Consider `--quiet` flag to suppress progress messages
