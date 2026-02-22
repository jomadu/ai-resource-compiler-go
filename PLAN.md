# Add Metadata Block and Markdown Compiler

## Goal
1. Add metadata block specification for compiled rule files
2. Add markdown compiler target
3. Update compiler architecture to accept namespace parameter
4. Update all compiler specs to include metadata block

## Metadata Block Examples

### Full Metadata Block (Rules)
```yaml
---
namespace: registry/package@1.0.0
ruleset:
    id: "cleanCode"
    name: "Clean Code"
    description: "Clean code practices"
    rules:
        - meaningfulNames
        - smallFunctions
        - singleResponsibility
rule:
    id: meaningfulNames
    name: "Use Meaningful Names"
    description: "Variables and functions should have descriptive names"
    enforcement: must
    scope:
        - files: ["**/*.ts", "**/*.js"]
---

# Use Meaningful Names (MUST)

Variables and functions should have descriptive names that reveal intent.
```

### Minimal Metadata Block (Rule without optional fields)
```yaml
---
namespace: registry/package@1.0.0
ruleset:
    id: "cleanCode"
    name: "Clean Code"
    rules:
        - meaningfulNames
        - smallFunctions
rule:
    id: meaningfulNames
    name: "Use Meaningful Names"
    enforcement: must
---

# Use Meaningful Names (MUST)

Rule body content.
```

### No Metadata Block (Prompts)
```markdown
Deploy the application to production environment.
```

**Note:** Promptsets do NOT include metadata blocks - just the prompt body.

## Changes

### 1. Create `specs/metadata-block.md`
New spec defining the metadata block structure embedded in compiled rule files.

**Content:**
- YAML structure with namespace, ruleset, and rule sections
- Namespace: Package identifier (e.g., "registry/package@1.0.0" or custom string)
- Ruleset section includes: id, name (required), description (optional), rules array (required)
- Rule section includes: id, name (required), description (optional), enforcement (optional), scope (optional)
- Used by: ALL targets for rules (kiro, cursor, claude, copilot, markdown)
- NOT used by: Prompts (any target)
- Enforcement levels in headers: "# Rule Name (MUST)", "# Rule Name (SHOULD)", "# Rule Name (MAY)", or "# Rule Name" (no enforcement)

### 2. Create `specs/markdown-compiler.md`
New compiler spec for vanilla markdown target.

**Format:**
- Extension: `.md`
- Rules: Metadata block frontmatter + body with enforcement in header
- Prompts: Plain body (no frontmatter)
- Installation: User-defined (no standard location)
- Naming: Same as other targets (`{id}.md`, `{collection-id}_{item-id}.md`)

### 3. Update `specs/compiler-architecture.md`
- Add namespace parameter to CompileOptions:
  ```go
  type CompileOptions struct {
      Targets          []string
      Namespace        string  // Package identifier for metadata block
      ResolveFragments bool
  }
  ```
- Add "markdown" to supported targets table
- Add reference to `metadata-block.md` in related specs
- Update examples to show namespace parameter
- Note: Namespace is required for rule compilation, optional for prompts (ignored)

### 4. Update `specs/kiro-compiler.md`
- Add metadata block to rule compilation (after frontmatter if any)
- Format:
  ```yaml
  ---
  namespace: registry/package@1.0.0
  ruleset: {...}
  rule: {...}
  ---
  
  # Rule Name (MUST)
  
  Rule body
  ```
- Reference `metadata-block.md` for structure
- Add enforcement in headers for rules
- Prompts remain plain markdown (no metadata)

### 5. Update `specs/cursor-compiler.md`
- Add metadata block after MDC frontmatter for rules
- Format:
  ```yaml
  ---
  description: "Rule description"
  globs: ["**/*.ts"]
  alwaysApply: true
  ---
  
  ---
  namespace: registry/package@1.0.0
  ruleset: {...}
  rule: {...}
  ---
  
  # Rule Name (MUST)
  
  Rule body
  ```
- Reference `metadata-block.md` for structure
- Add enforcement in headers for rules
- Prompts remain plain markdown (no metadata)

### 6. Update `specs/claude-compiler.md`
- Add metadata block after optional paths frontmatter for rules
- Format:
  ```yaml
  ---
  paths:
    - "src/**/*.ts"
  ---
  
  ---
  namespace: registry/package@1.0.0
  ruleset: {...}
  rule: {...}
  ---
  
  # Rule Name (MUST)
  
  Rule body
  ```
- Prompts: No metadata block (just optional name/description frontmatter + body)
- Reference `metadata-block.md` for structure
- Add enforcement in headers for rules

### 7. Update `specs/copilot-compiler.md`
- Add metadata block after applyTo frontmatter for rules
- Format:
  ```yaml
  ---
  applyTo: ["**/*.ts"]
  ---
  
  ---
  namespace: registry/package@1.0.0
  ruleset: {...}
  rule: {...}
  ---
  
  # Rule Name (MUST)
  
  Rule body
  ```
- Reference `metadata-block.md` for structure
- Add enforcement in headers for rules
- Prompts remain plain markdown (no metadata)

### 8. Update `specs/README.md`
- Add `metadata-block.md` to foundation section
- Add `markdown-compiler.md` to target compilers section
- Update JTBDs to mention metadata block embedding

### 9. Update `REVIEW_CHECKLIST.md`
- Add `metadata-block.md` to foundation
- Add `markdown-compiler.md` to target compilers

### 10. Update main `README.md`
- Add markdown to supported targets list in overview

### 11. Update `specs/cli-design.md`
- Add `--namespace` flag to compile command
- Required for rule compilation
- Optional for prompt compilation (ignored)
- Example: `arc compile --target cursor --namespace myorg/rules@1.0.0 --output .cursor/rules rules.yml`

## Key Changes Summary

**Separate Compilation Methods by Resource Type:**
```go
type TargetCompiler interface {
    Name() string
    CompileRule(rule *airesource.Rule, namespace string) (CompilationResult, error)
    CompileRuleset(ruleset *airesource.Ruleset, namespace string) ([]CompilationResult, error)
    CompilePrompt(prompt *airesource.Prompt) (CompilationResult, error)
    CompilePromptset(promptset *airesource.Promptset) ([]CompilationResult, error)
}
```

**All Targets Include Metadata Block for Rules:**
- Kiro: Metadata block + body with enforcement header
- Cursor: MDC frontmatter + metadata block + body with enforcement header
- Claude: Paths frontmatter (optional) + metadata block + body with enforcement header
- Copilot: ApplyTo frontmatter + metadata block + body with enforcement header
- Markdown: Metadata block + body with enforcement header

**Namespace Parameter:**
- Added to CompileOptions in compiler architecture
- Required for rule/ruleset compilation (type-safe via separate methods)
- Not used for prompt/promptset compilation
- Passed by CLI via `--namespace` flag

**Enforcement in Headers:**
- All targets include enforcement level in rule headers
- Format: "# Rule Name (MUST/SHOULD/MAY)" or "# Rule Name" (no enforcement)

## Order of Execution
1. Create `specs/metadata-block.md`
2. Create `specs/markdown-compiler.md`
3. Update `specs/compiler-architecture.md` (add namespace parameter)
4. Update `specs/kiro-compiler.md` (add metadata block)
5. Update `specs/cursor-compiler.md` (add metadata block)
6. Update `specs/claude-compiler.md` (add metadata block)
7. Update `specs/copilot-compiler.md` (add metadata block)
8. Update `specs/cli-design.md` (add --namespace flag)
9. Update `specs/README.md`
10. Update `REVIEW_CHECKLIST.md`
11. Update main `README.md`
