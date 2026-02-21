# Target Formats

## Job to be Done
Define the output format specifications for each compilation target so that compiled resources integrate correctly with their respective AI coding tools.

## Activities
- Document Kiro CLI modular format (.kiro/steering/, .kiro/prompts/)
- Document Cursor modular MDC format (.cursor/rules/, .cursor/commands/)
- Document Claude Code modular format (.claude/rules/, .claude/skills/)
- Document GitHub Copilot modular format (.github/instructions/, .github/prompts/)
- Define file extension conventions for each target
- Define frontmatter specifications for each target
- Define naming conventions for single resources and collection items
- Specify how resource kinds map to target formats

## Acceptance Criteria
- [x] Each target format is fully specified
- [x] Prompt formatting is defined for each target
- [x] Rule formatting is defined for each target
- [x] Collection handling (Promptset/Ruleset) is defined
- [x] Fragment resolution behavior is specified
- [x] Examples provided for each target
- [x] Format differences are clearly documented
- [x] File extension conventions documented
- [x] Frontmatter specifications documented
- [x] Naming conventions documented
- [x] Claude prompt directory structure documented

## Naming Conventions

**Single Resources:**
- Use resource ID as filename: `{id}.{ext}`
- Example: `id: api-standards` → `api-standards.md`

**Collection Items:**
- Combine collection ID and item ID: `{collection-id}_{item-id}.{ext}`
- Example: Ruleset `id: backend` with rule `id: api` → `backend_api.md`

**Extensions by Target:**
- Kiro: `.md`
- Cursor: `.mdc`
- Claude: `.md` (rules), directory with `SKILL.md` (prompts)
- Copilot: `.instructions.md` (rules), `.prompt.md` (prompts)

## Target Format Specifications

### Kiro CLI

**Rules:** `.kiro/steering/{id}.md`  
**Prompts:** `.kiro/prompts/{id}.md`  
**Format:** Plain markdown

**Single prompt (`id: deploy`):**
```markdown
[Prompt body content]
```
**Path:** `deploy.md`

**Single rule (`id: api-standards`):**
```markdown
[Rule body content]
```
**Path:** `api-standards.md`

**Promptset (`id: ci-workflows`) with prompt (`id: deploy`):**
```markdown
[Prompt body content]
```
**Path:** `ci-workflows_deploy.md`

**Ruleset (`id: backend`) with rule (`id: api`):**
```markdown
[Rule body content]
```
**Path:** `backend_api.md`

**Notes:**
- Plain markdown, no frontmatter
- One file per resource or collection item
- Collection items use `{collection-id}_{item-id}.md` naming

### Cursor

**Rules:** `.cursor/rules/{id}.mdc`  
**Prompts:** `.cursor/commands/{id}.mdc`  
**Format:** MDC with frontmatter

**Single prompt (`id: deploy`):**
```markdown
---
description: Deploy application
globs: ["**/*.yml"]
alwaysApply: true
---

[Prompt body content]
```
**Path:** `deploy.mdc`

**Single rule (`id: api-standards`):**
```markdown
---
description: API design standards
globs: ["src/**/*.ts"]
alwaysApply: false
---

[Rule body content]
```
**Path:** `api-standards.mdc`

**Promptset (`id: ci-workflows`) with prompt (`id: deploy`):**
```markdown
---
description: Deploy workflow
globs: ["**/*.yml"]
alwaysApply: true
---

[Prompt body content]
```
**Path:** `ci-workflows_deploy.mdc`

**Ruleset (`id: backend`) with rule (`id: api`):**
```markdown
---
description: Backend API standards
globs: ["src/**/*.ts"]
alwaysApply: false
---

[Rule body content]
```
**Path:** `backend_api.mdc`

**Frontmatter Fields:**
- `description` - Resource name or description
- `globs` - File patterns (from scope if present)
- `alwaysApply` - Boolean (true for prompts, false for rules by default)

**Notes:**
- MDC format with YAML frontmatter
- One file per resource or collection item
- Collection items use `{collection-id}_{item-id}.mdc` naming

### Claude Code

**Rules:** `.claude/rules/{id}.md`  
**Prompts:** `.claude/skills/{id}/SKILL.md`  
**Format:** Markdown with optional frontmatter

**Single prompt (`id: deploy`):**
```markdown
---
name: Deploy Application
description: Deploy to production
---

[Prompt body content]
```
**Path:** `deploy/SKILL.md`

**Single rule (`id: api-standards`):**
```markdown
[Rule body content]
```
**Path:** `api-standards.md`

**Promptset (`id: ci-workflows`) with prompt (`id: deploy`):**
```markdown
---
name: Deploy Workflow
description: CI/CD deployment
---

[Prompt body content]
```
**Path:** `ci-workflows_deploy/SKILL.md`

**Ruleset (`id: backend`) with rule (`id: api`):**
```markdown
[Rule body content]
```
**Path:** `backend_api.md`

**Notes:**
- Rules: plain markdown files
- Prompts: directory with `SKILL.md` file
- Prompt frontmatter optional (name, description)
- Collection items use `{collection-id}_{item-id}` naming

### GitHub Copilot

**Rules:** `.github/instructions/{id}.instructions.md`  
**Prompts:** `.github/prompts/{id}.prompt.md`  
**Format:** Markdown with frontmatter

**Single prompt (`id: deploy`):**
```markdown
---
applyTo: ["**/*.yml"]
---

[Prompt body content]
```
**Path:** `deploy.prompt.md`

**Single rule (`id: api-standards`):**
```markdown
---
applyTo: ["src/**/*.ts"]
excludeAgent: ["copilot-chat"]
---

[Rule body content]
```
**Path:** `api-standards.instructions.md`

**Promptset (`id: ci-workflows`) with prompt (`id: deploy`):**
```markdown
---
applyTo: ["**/*.yml"]
---

[Prompt body content]
```
**Path:** `ci-workflows_deploy.prompt.md`

**Ruleset (`id: backend`) with rule (`id: api`):**
```markdown
---
applyTo: ["src/**/*.ts"]
excludeAgent: []
---

[Rule body content]
```
**Path:** `backend_api.instructions.md`

**Frontmatter Fields:**
- `applyTo` - File glob patterns (from scope if present)
- `excludeAgent` - Agents to exclude (optional)

**Notes:**
- Markdown with YAML frontmatter
- Different extensions for rules vs prompts
- Collection items use `{collection-id}_{item-id}` naming

## Format Comparison

| Feature | Kiro | Cursor | Claude | Copilot |
|---------|------|--------|--------|---------|
| File format | Markdown | MDC | Markdown | Markdown |
| Output | One file per item | One file per item | One file per item | One file per item |
| Frontmatter | None | YAML | YAML (prompts) | YAML |
| Prompt location | `.kiro/prompts/` | `.cursor/commands/` | `.claude/skills/{id}/` | `.github/prompts/` |
| Rule location | `.kiro/steering/` | `.cursor/rules/` | `.claude/rules/` | `.github/instructions/` |
| Prompt extension | `.md` | `.mdc` | `SKILL.md` | `.prompt.md` |
| Rule extension | `.md` | `.mdc` | `.md` | `.instructions.md` |
| Collection naming | `{coll}_{item}.md` | `{coll}_{item}.mdc` | `{coll}_{item}.md` | `{coll}_{item}.instructions.md` |

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Empty body | Skip output (no file created) |
| Missing name | Use ID as fallback |
| Missing description | Omit from frontmatter |
| Fragment not resolved | Error (fragments must be resolved) |
| Multi-line body | Preserve formatting |
| Special characters | Escape per target format |
| Empty collection | No files created |
| Missing scope | Omit globs/applyTo from frontmatter |

## Dependencies

- `compiler-architecture.md` - Compilation pipeline
- `kiro-compiler.md` - Kiro-specific compilation
- `cursor-compiler.md` - Cursor-specific compilation
- `claude-compiler.md` - Claude-specific compilation
- `copilot-compiler.md` - Copilot-specific compilation

## Implementation Mapping

**Source files:**
- `pkg/targets/kiro/compiler.go` - Kiro compiler implementation
- `pkg/targets/cursor/compiler.go` - Cursor compiler implementation
- `pkg/targets/claude/compiler.go` - Claude compiler implementation
- `pkg/targets/copilot/compiler.go` - Copilot compiler implementation

**Related specs:**
- `kiro-compiler.md` - Kiro-specific compilation
- `cursor-compiler.md` - Cursor-specific compilation
- `claude-compiler.md` - Claude-specific compilation
- `copilot-compiler.md` - Copilot-specific compilation

## Examples

### Example 1: Single Prompt to All Targets

**Input:**
```yaml
apiVersion: ai-resource/draft
kind: Prompt
metadata:
  id: deploy
  name: Deploy Application
  description: Deploy to production
spec:
  body: "Deploy the application to production"
```

**Kiro Output:**
- Path: `deploy.md`
- Content:
```markdown
Deploy the application to production
```

**Cursor Output:**
- Path: `deploy.mdc`
- Content:
```markdown
---
description: Deploy Application
alwaysApply: true
---

Deploy the application to production
```

**Claude Output:**
- Path: `deploy/SKILL.md`
- Content:
```markdown
---
name: Deploy Application
description: Deploy to production
---

Deploy the application to production
```

**Copilot Output:**
- Path: `deploy.prompt.md`
- Content:
```markdown
---
applyTo: []
---

Deploy the application to production
```

### Example 2: Single Rule to All Targets

**Input:**
```yaml
apiVersion: ai-resource/draft
kind: Rule
metadata:
  id: api-standards
  name: API Design Standards
spec:
  enforcement: must
  scope:
    include: ["src/**/*.ts"]
  body: "Follow RESTful API design principles"
```

**Kiro Output:**
- Path: `api-standards.md`
- Content:
```markdown
Follow RESTful API design principles
```

**Cursor Output:**
- Path: `api-standards.mdc`
- Content:
```markdown
---
description: API Design Standards
globs: ["src/**/*.ts"]
alwaysApply: false
---

Follow RESTful API design principles
```

**Claude Output:**
- Path: `api-standards.md`
- Content:
```markdown
Follow RESTful API design principles
```

**Copilot Output:**
- Path: `api-standards.instructions.md`
- Content:
```markdown
---
applyTo: ["src/**/*.ts"]
---

Follow RESTful API design principles
```

### Example 3: Promptset to All Targets

**Input:**
```yaml
apiVersion: ai-resource/draft
kind: Promptset
metadata:
  id: ci-workflows
spec:
  prompts:
    deploy:
      name: Deploy
      body: "Deploy the application"
    test:
      name: Test
      body: "Run test suite"
```

**Kiro Output:**
- Path: `ci-workflows_deploy.md` - Content: `Deploy the application`
- Path: `ci-workflows_test.md` - Content: `Run test suite`

**Cursor Output:**
- Path: `ci-workflows_deploy.mdc` - Content with frontmatter
- Path: `ci-workflows_test.mdc` - Content with frontmatter

**Claude Output:**
- Path: `ci-workflows_deploy/SKILL.md` - Content with frontmatter
- Path: `ci-workflows_test/SKILL.md` - Content with frontmatter

**Copilot Output:**
- Path: `ci-workflows_deploy.prompt.md` - Content with frontmatter
- Path: `ci-workflows_test.prompt.md` - Content with frontmatter

### Example 4: Ruleset to All Targets

**Input:**
```yaml
apiVersion: ai-resource/draft
kind: Ruleset
metadata:
  id: backend
spec:
  rules:
    api:
      name: API Standards
      body: "Follow REST principles"
    security:
      name: Security
      body: "Validate all inputs"
```

**Kiro Output:**
- Path: `backend_api.md` - Content: `Follow REST principles`
- Path: `backend_security.md` - Content: `Validate all inputs`

**Cursor Output:**
- Path: `backend_api.mdc` - Content with frontmatter
- Path: `backend_security.mdc` - Content with frontmatter

**Claude Output:**
- Path: `backend_api.md` - Content: `Follow REST principles`
- Path: `backend_security.md` - Content: `Validate all inputs`

**Copilot Output:**
- Path: `backend_api.instructions.md` - Content with frontmatter
- Path: `backend_security.instructions.md` - Content with frontmatter

## Notes

- All targets receive fully resolved resources (fragments already rendered)
- Format differences reflect each tool's conventions and capabilities
- Modular approach: one file per resource or collection item
- Compiler returns relative paths; users prepend target-specific directories
- Collections are flattened: each item becomes a separate file
- Naming convention for collection items: `{collection-id}_{item-id}.{ext}`
- Target compilers are responsible for escaping special characters
- Frontmatter usage varies by target (Cursor and Copilot require it, Claude optional, Kiro none)
- Claude prompts use directory structure with `SKILL.md` file

## Known Issues

None.

## Areas for Improvement

- Could support custom templates per target
- Could add configuration for frontmatter field customization
- Could support target-specific extensions
- Could add validation that output conforms to target tool's requirements
- Could support custom naming patterns for collection items
