# TODO

## TASK-001
- Priority: 1
- Status: DONE
- Dependencies: []
- Description: Write specs/metadata-block.md - Define YAML metadata structure for preserving ruleset/rule context in compiled output

## TASK-002
- Priority: 1
- Status: TODO
- Dependencies: []
- Description: Write specs/compiler-architecture.md - Define TargetCompiler interface, Target enum, CompileOptions, and compilation pipeline

## TASK-003
- Priority: 2
- Status: TODO
- Dependencies: [TASK-001, TASK-002]
- Description: Write specs/markdown-compiler.md - Define vanilla markdown output format for rules and prompts

## TASK-004
- Priority: 2
- Status: TODO
- Dependencies: [TASK-001, TASK-002]
- Description: Write specs/kiro-compiler.md - Define Kiro CLI steering rules and prompts format

## TASK-005
- Priority: 2
- Status: TODO
- Dependencies: [TASK-001, TASK-002]
- Description: Write specs/cursor-compiler.md - Define Cursor IDE rules (.mdc) and commands (.md) format

## TASK-006
- Priority: 2
- Status: TODO
- Dependencies: [TASK-001, TASK-002]
- Description: Write specs/claude-compiler.md - Define Claude Code rules and skills format

## TASK-007
- Priority: 2
- Status: TODO
- Dependencies: [TASK-001, TASK-002]
- Description: Write specs/copilot-compiler.md - Define GitHub Copilot instructions and prompts format

## TASK-008
- Priority: 3
- Status: TODO
- Dependencies: [TASK-003, TASK-004, TASK-005, TASK-006, TASK-007]
- Description: Write specs/cli-design.md - Define arc compile command with target selection and output modes

## TASK-009
- Priority: 3
- Status: TODO
- Dependencies: [TASK-001, TASK-002, TASK-003, TASK-004, TASK-005, TASK-006, TASK-007, TASK-008]
- Description: Write specs/README.md - Create specification index organized by JTBD and topics of concern

## TASK-010
- Priority: 4
- Status: TODO
- Dependencies: [TASK-009]
- Description: Write README.md - Create user-facing documentation with overview, usage, CLI examples, and architecture
