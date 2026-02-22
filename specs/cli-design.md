# CLI Design

## Job to be Done
Provide a command-line interface for compiling AI Resources to target formats with intuitive commands, clear output, and helpful error messages.

## Activities
- Define command structure and subcommands
- Specify flag handling and options
- Define output behavior (stdout vs files)
- Provide help and usage documentation
- Handle errors with actionable messages

## Acceptance Criteria
- [ ] `arc compile` command compiles resources to targets
- [ ] `--target` flag specifies compilation targets (comma-separated)
- [ ] `--output` flag specifies output directory
- [ ] Multiple files can be compiled in single invocation
- [ ] Help text is clear and includes examples
- [ ] Errors indicate what went wrong and how to fix
- [ ] Exit codes follow conventions (0=success, 1=error)
- [ ] Version information available via `--version`

## Data Structures

### Command Structure
```
arc [global-flags] <command> [command-flags] [args]

Global Flags:
  --version, -v    Show version information
  --help, -h       Show help

Commands:
  compile          Compile AI Resources to target formats
  version          Show version information
```

### Compile Command
```
arc compile [flags] <file>...

Flags:
  --target, -t     Target formats (comma-separated: kiro,cursor,claude,copilot)
  --output, -o     Output directory (required when writing files)
  --stdout         Write output to stdout instead of files
  --help, -h       Show help for compile command

Arguments:
  <file>...        One or more resource files to compile
```

### Recommended Installation Locations

The compiler returns relative paths. Users prepend target-specific directories:

| Target   | Rules                  | Prompts               |
|----------|------------------------|-----------------------|
| Kiro     | `.kiro/steering/`      | `.kiro/prompts/`      |
| Cursor   | `.cursor/rules/`       | `.cursor/commands/`   |
| Claude   | `.claude/rules/`       | `.claude/skills/`     |
| Copilot  | `.github/instructions/`| `.github/prompts/`    |

**User Responsibility:**
- Create output directories before compilation
- Choose appropriate installation location per target
- Manage file organization and cleanup

## Algorithm

### Command Execution

1. Parse global flags
2. Parse command and command-specific flags
3. Validate inputs
4. Load resources
5. Compile to targets
6. Write output (files or stdout)
7. Report results
8. Exit with appropriate code

**Pseudocode:**
```
function main():
    args = parse_args()
    
    if args.version:
        print_version()
        exit(0)
    
    if args.help or no_command:
        print_help()
        exit(0)
    
    switch args.command:
        case "compile":
            exit_code = run_compile(args)
            exit(exit_code)
        case "version":
            print_version()
            exit(0)
        default:
            print_error("unknown command: {args.command}")
            exit(1)
```

### Compile Command Execution

```
function run_compile(args):
    // Validate inputs
    if len(args.files) == 0:
        print_error("no input files specified")
        return 1
    
    if len(args.targets) == 0:
        print_error("no targets specified (use --target)")
        return 1
    
    if not args.stdout and args.output == "":
        print_error("output directory required (use --output)")
        return 1
    
    // Compile each file
    all_success = true
    for file in args.files:
        results = compiler.Compile(file, CompileOptions{
            Targets: args.targets,
        })
        
        for result in results:
            if result.Error:
                print_error("{file} -> {result.Target}: {result.Error}")
                all_success = false
            else:
                if args.stdout:
                    print(result.Content)
                else:
                    output_path = join(args.output, result.Path)
                    write_file(output_path, result.Content)
                    print_success("{file} -> {result.Path}")
    
    if all_success:
        return 0
    else:
        return 1
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| No files specified | Error: "no input files specified" |
| No targets specified | Error: "no targets specified" |
| No output dir specified | Error: "output directory required (use --output)" |
| File doesn't exist | Error: "file not found: {path}" |
| Invalid target name | Error: "unknown target: {name}" |
| Output dir doesn't exist | Error: "output directory not found: {path}" |
| Permission denied | Error: "cannot write to {path}" |
| Multiple files, one fails | Continue processing, exit 1 at end |
| Stdout with multiple targets | Separate with headers |
| Invalid flag | Error with usage hint |

## Dependencies

- `compiler-architecture.md` - Compiler API
- Flag parsing library (e.g., spf13/cobra or stdlib flag)

## Implementation Mapping

**Source files:**
- `cmd/arc/main.go` - Entry point and command routing
- `cmd/arc/compile.go` - Compile command implementation
- `cmd/arc/version.go` - Version command
- `cmd/arc/help.go` - Help text and usage

**Related specs:**
- `compiler-architecture.md` - Compilation API
- `kiro-compiler.md` - Kiro target format
- `cursor-compiler.md` - Cursor target format
- `claude-compiler.md` - Claude target format
- `copilot-compiler.md` - Copilot target format

## Examples

### Example 1: Compile to Single Target

**Input:**
```bash
arc compile --target cursor --output .cursor/rules prompt.yml
```

**Expected Output:**
```
✓ prompt.yml -> deploy.mdc
```

**Verification:**
- `.cursor/rules/deploy.mdc` file created
- Success message printed
- Exit code 0

### Example 2: Compile to Multiple Targets

**Input:**
```bash
arc compile --target kiro,cursor,claude --output ./build prompt.yml
```

**Expected Output:**
```
✓ prompt.yml -> deploy.md (kiro)
✓ prompt.yml -> deploy.mdc (cursor)
✓ prompt.yml -> deploy/SKILL.md (claude)
```

**Verification:**
- Three files created in `./build/`
- Three success messages
- Exit code 0

### Example 3: Compile Multiple Files

**Input:**
```bash
arc compile --target cursor --output .cursor/rules prompts.yml rules.yml
```

**Expected Output:**
```
✓ prompts.yml -> deploy.mdc
✓ rules.yml -> api-standards.mdc
```

**Verification:**
- Both files compiled to `.cursor/rules/`
- Two success messages
- Exit code 0

### Example 4: Output to Stdout

**Input:**
```bash
arc compile --target cursor --stdout prompt.yml
```

**Expected Output:**
```
---
description: Deploy Application
alwaysApply: true
---

Deploy the application to production
```

**Verification:**
- Content written to stdout
- No files created
- Exit code 0

### Example 5: Error Handling

**Input:**
```bash
arc compile --target unknown prompt.yml
```

**Expected Output:**
```
✗ prompt.yml -> unknown: unknown target: unknown
```

**Verification:**
- Error message printed
- Exit code 1

### Example 6: Help Text

**Input:**
```bash
arc compile --help
```

**Expected Output:**
```
Compile AI Resources to target formats

Usage:
  arc compile [flags] <file>...

Flags:
  -t, --target string   Target formats (comma-separated: kiro,cursor,claude,copilot)
  -o, --output string   Output directory (required when writing files)
      --stdout          Write output to stdout instead of files
  -h, --help            Show help for compile command

Recommended Installation Locations:
  Kiro     .kiro/steering/ (rules), .kiro/prompts/ (prompts)
  Cursor   .cursor/rules/ (rules), .cursor/commands/ (prompts)
  Claude   .claude/rules/ (rules), .claude/skills/ (prompts)
  Copilot  .github/instructions/ (rules), .github/prompts/ (prompts)

Examples:
  # Compile to Cursor format
  arc compile --target cursor --output .cursor/rules prompt.yml

  # Compile to multiple targets
  arc compile --target kiro,cursor,claude --output ./build prompt.yml

  # Compile multiple files
  arc compile --target cursor --output .cursor/rules prompts.yml rules.yml

  # Output to stdout
  arc compile --target cursor --stdout prompt.yml
```

**Verification:**
- Help text is clear
- Examples provided
- Exit code 0

## Notes

- CLI follows standard Unix conventions (flags, exit codes, stdout/stderr)
- Error messages go to stderr, success messages to stdout
- `--stdout` is useful for piping output or previewing
- Multiple targets create multiple output files
- File paths can be relative or absolute
- Users must specify output directory with `--output` flag
- Users are responsible for creating output directories
- Compiler returns relative paths; users choose installation location
- CLI is thin wrapper around compiler package (business logic in library)

## Known Issues

None.

## Areas for Improvement

- Could add `--watch` mode for continuous compilation
- Could add `--dry-run` to preview without writing
- Could add `--verbose` for detailed output
- Could add `--quiet` to suppress success messages
- Could support glob patterns for input files
- Could add `arc init` to create example resources
- Could add `arc validate` to check resources without compiling
