# TODO

## TASK-001
- Priority: 1
- Status: DONE
- Dependencies: []
- Issue: #1 (Metadata Block Duplication Across Targets)
- Description: Add shared metadata generation functions to metadata-block.md spec. Document `GenerateMetadataBlock(resource)` and `GenerateEnforcementHeader(rule)` functions that all target compilers will use. Include function signatures, parameters, return values, and examples.

## TASK-002
- Priority: 1
- Status: DONE
- Dependencies: []
- Issue: #2 (Path Generation Logic Scattered Across Targets)
- Description: Add centralized path generation functions to compiler-architecture.md spec. Document `BuildRulePath(rulesetID, ruleID, extension)`, `BuildPromptPath(promptsetID, promptID, extension)`, and `BuildClaudePromptPath(promptsetID, promptID)` in the Implementation Mapping section. Specify these live in `internal/format/paths.go`.

## TASK-003
- Priority: 1
- Status: DONE
- Dependencies: []
- Issue: #3 (No Validation of Resource IDs for Filesystem Safety)
- Description: Create validation-rules.md spec in specs/ directory. Define filesystem-safe ID requirements (allowed characters: a-z, A-Z, 0-9, -, _). Document `ValidateResourceIDs(resource)` function that checks collection and item IDs. Include edge cases for invalid characters (/, \, :, *, ?, ", <, >, |) and error message format.

## TASK-004
- Priority: 1
- Status: DONE
- Dependencies: [TASK-003]
- Issue: #7 (Enforcement Header Formatting Inconsistency Risk)
- Description: Add rule name validation to validation-rules.md spec. Document `ValidateRuleForCompilation(rule)` function that rejects rule names containing parentheses. Specify error message format and rationale (enforcement header parsing).

## TASK-005
- Priority: 1
- Status: DONE
- Dependencies: [TASK-003]
- Issue: #3 (No Validation of Resource IDs for Filesystem Safety)
- Description: Update compiler-architecture.md spec to include validation step in compilation pipeline. Add validation as step 1 in Algorithm section: "Validate resource IDs and rule names (if rule type)". Reference validation-rules.md spec in Dependencies section.

## TASK-006
- Priority: 2
- Status: TODO
- Dependencies: []
- Issue: #4 (Version Handling Strategy Incomplete)
- Description: Add `SupportedVersions() []string` method to TargetCompiler interface in compiler-architecture.md spec. Update interface definition in Data Structures section and add version compatibility check to compilation pipeline algorithm (step 1.5: "Check target supports resource.APIVersion").

## TASK-007
- Priority: 2
- Status: TODO
- Dependencies: [TASK-006]
- Issue: #4 (Version Handling Strategy Incomplete)
- Description: Update all target compiler specs (kiro-compiler.md, cursor-compiler.md, claude-compiler.md, copilot-compiler.md, markdown-compiler.md) to include `SupportedVersions()` method in their compiler struct definitions. Document which versions each target supports (start with ["ai-resource/v1"]).

## TASK-008
- Priority: 2
- Status: TODO
- Dependencies: []
- Issue: #5 (CLI Output Mode Ambiguity)
- Description: Update cli-design.md spec to document output format for stdout mode. Specify format: `=== {target}/{path} ===\n{content}\n\n`. Add examples showing multi-target output with path prefixes.

## TASK-009
- Priority: 2
- Status: TODO
- Dependencies: [TASK-008]
- Issue: #5 (CLI Output Mode Ambiguity)
- Description: Update cli-design.md spec to document directory output mode behavior. Specify that `--output <dir>` creates target subdirectories (e.g., `<dir>/markdown/`, `<dir>/kiro/`). Add `--flat` flag to disable subdirectories for single-target use. Include examples and edge cases.

## TASK-010
- Priority: 3
- Status: TODO
- Dependencies: []
- Issue: #6 (No Guidance on Multi-Resource Compilation)
- Description: Update README.md to clarify multi-resource compilation pattern. Add note in "Library API" section that `Compile()` accepts single resource and users should iterate for multiple resources. Show error handling pattern in multi-resource example.

## TASK-011
- Priority: 3
- Status: TODO
- Dependencies: [TASK-001, TASK-002]
- Issue: #1 (Metadata Block Duplication), #2 (Path Generation Logic Scattered)
- Description: Update all target compiler specs (kiro-compiler.md, cursor-compiler.md, claude-compiler.md, copilot-compiler.md, markdown-compiler.md) to reference shared metadata and path generation functions. Remove duplicated algorithm details and reference `internal/format/metadata.go` and `internal/format/paths.go` instead.

## TASK-012
- Priority: 3
- Status: TODO
- Dependencies: [TASK-001, TASK-002, TASK-003, TASK-004, TASK-005, TASK-006, TASK-007, TASK-008, TASK-009, TASK-010, TASK-011]
- Issue: Meta-task (consolidates all spec updates)
- Description: Update specs/README.md to reference new validation-rules.md spec in Foundation Layer section. Add validation to the reading order for understanding the system (after metadata-block.md, before compiler-architecture.md).
