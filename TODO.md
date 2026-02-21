# TODO

## TASK-001
- Priority: 1
- Status: DONE
- Dependencies: []
- Description: Update README.md with correct path conventions
  - Fix "Supported Targets" table with correct output path examples
  - Update "Compilation Results" section to explain modular output
  - Remove misleading note about single-file targets
  - Add "Recommended Installation Locations" table
  - Update code examples to show correct paths
  - Clarify compiler returns relative paths, users choose installation location
  - Note: Path conventions should match target-formats.md spec

## TASK-002
- Priority: 1
- Status: DONE
- Dependencies: []
- Description: Update specs/target-formats.md for modular approach
  - Replace single-file format specs with modular approach
  - Add file extension conventions per target (.mdc, .instructions.md, .prompt.md)
  - Add Claude prompt directory structure ({id}/SKILL.md)
  - Add frontmatter specifications for Cursor (.mdc) and Copilot (.instructions.md)
  - Update all examples to show correct paths
  - Add naming conventions section (single vs collection items)
  - Update format comparison table
  - Note: This spec becomes the source of truth for path conventions

## TASK-003
- Priority: 2
- Status: DONE
- Dependencies: [TASK-002]
- Description: Update specs/cursor-compiler.md for modular .mdc output
  - Change output from single .cursorrules to modular .cursor/rules/*.mdc
  - Add frontmatter generation (description, globs, alwaysApply)
  - Update path generation to use resource ID
  - Update naming for collection items ({collection-id}_{item-id})
  - Update all examples with correct paths
  - Add frontmatter examples

## TASK-004
- Priority: 2
- Status: DONE
- Dependencies: [TASK-002]
- Description: Create specs/kiro-compiler.md
  - Rules → .kiro/steering/{id}.md
  - Prompts → .kiro/prompts/{id}.md
  - Collection naming: {collection-id}_{item-id}.md
  - Plain markdown format, no frontmatter
  - Follow TEMPLATE.md structure
  - Include examples for single resources and collections

## TASK-005
- Priority: 2
- Status: TODO
- Dependencies: [TASK-002]
- Description: Create specs/claude-compiler.md
  - Rules → .claude/rules/{id}.md
  - Prompts → .claude/skills/{id}/SKILL.md
  - Collection naming: {collection-id}_{item-id}.md or {collection-id}_{item-id}/SKILL.md
  - Markdown format with optional frontmatter for skills (name, description)
  - Follow TEMPLATE.md structure
  - Include examples for single resources and collections

## TASK-006
- Priority: 2
- Status: TODO
- Dependencies: [TASK-002]
- Description: Create specs/copilot-compiler.md
  - Rules → .github/instructions/{id}.instructions.md
  - Prompts → .github/prompts/{id}.prompt.md
  - Collection naming: {collection-id}_{item-id}.instructions.md
  - Markdown format with frontmatter (applyTo globs, excludeAgent)
  - Follow TEMPLATE.md structure
  - Include examples for single resources and collections

## TASK-007
- Priority: 3
- Status: TODO
- Dependencies: [TASK-003, TASK-004, TASK-005, TASK-006]
- Description: Update specs/compiler-architecture.md with correct examples
  - Update all path examples to match target-formats.md conventions
  - Add naming conventions section
  - Update CompilationResult examples
  - Add note about recommended installation locations
  - Update edge cases to reflect modular approach
  - Update examples section with correct paths

## TASK-008
- Priority: 3
- Status: TODO
- Dependencies: [TASK-003, TASK-004, TASK-005, TASK-006]
- Description: Update specs/cli-design.md with correct examples
  - Update all output examples with correct paths
  - Add recommended installation locations guidance
  - Update help text examples
  - Clarify user responsibility for directory structure
  - Add example showing output directory usage

## TASK-009
- Priority: 4
- Status: TODO
- Dependencies: [TASK-003, TASK-004, TASK-005, TASK-006]
- Description: Update specs/README.md index
  - Remove references to non-existent specs (compilation-pipeline.md, fragment-handling.md, output-management.md)
  - Add actual spec list with current status
  - Update JTBDs to reflect modular approach
  - Document that path conventions are defined in target-formats.md
  - List all four target compiler specs (kiro, cursor, claude, copilot)
