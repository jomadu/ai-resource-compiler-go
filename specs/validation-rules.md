# Resource Validation Rules

## Job to be Done
Ensure resource IDs are filesystem-safe before compilation to prevent path traversal, invalid filenames, and cross-platform compatibility issues.

## Activities
1. Validate collection IDs (ruleset, promptset) for filesystem safety
2. Validate item IDs (rule, prompt) for filesystem safety
3. Reject resources with invalid characters in IDs
4. Return clear error messages identifying invalid IDs and characters

## Acceptance Criteria
- [ ] ValidateResourceIDs function checks all collection and item IDs
- [ ] Only alphanumeric characters, hyphens, and underscores allowed (a-z, A-Z, 0-9, -, _)
- [ ] Invalid characters rejected: / \ : * ? " < > |
- [ ] Error messages identify which ID is invalid and which character caused rejection
- [ ] Empty IDs rejected with clear error message
- [ ] Validation occurs before any compilation steps

## Data Structures

### ValidationError
```go
type ValidationError struct {
    Field      string   // e.g., "ruleset.id", "rule.id"
    Value      string   // The invalid ID value
    InvalidChar string  // The character that caused rejection
    Message    string   // Human-readable error message
}
```

**Fields:**
- `Field` - Dot-notation path to invalid field (ruleset.id, rule.id, promptset.id, prompt.id)
- `Value` - The ID value that failed validation
- `InvalidChar` - The specific character that violated rules (empty if ID is empty)
- `Message` - Formatted error message for user display

## Algorithm

### ValidateResourceIDs

```go
func ValidateResourceIDs(resource Resource) error
```

**Parameters:**
- `resource` - Resource to validate (contains ruleset/promptset and rule/prompt)

**Returns:**
- `nil` if all IDs are valid
- `ValidationError` if any ID is invalid

**Algorithm:**
1. Determine resource type (rule vs prompt)
2. If rule type:
   - Validate `resource.Ruleset.ID`
   - Validate `resource.Rule.ID`
3. If prompt type:
   - Validate `resource.Promptset.ID`
   - Validate `resource.Prompt.ID`
4. For each ID:
   - Check if empty → return error
   - Check each character against allowed set
   - If invalid character found → return error with details
5. Return nil if all IDs valid

**Pseudocode:**
```
function ValidateResourceIDs(resource):
    if resource.Type == "rule":
        if err := validateID(resource.Ruleset.ID, "ruleset.id"):
            return err
        if err := validateID(resource.Rule.ID, "rule.id"):
            return err
    else if resource.Type == "prompt":
        if err := validateID(resource.Promptset.ID, "promptset.id"):
            return err
        if err := validateID(resource.Prompt.ID, "prompt.id"):
            return err
    
    return nil

function validateID(id, fieldName):
    if id is empty:
        return ValidationError{
            Field: fieldName,
            Value: "",
            Message: fieldName + " cannot be empty"
        }
    
    for each character in id:
        if character not in [a-z, A-Z, 0-9, -, _]:
            return ValidationError{
                Field: fieldName,
                Value: id,
                InvalidChar: character,
                Message: fieldName + " contains invalid character '" + character + "' in '" + id + "'"
            }
    
    return nil
```

### ValidateRuleForCompilation

```go
func ValidateRuleForCompilation(rule Rule) error
```

**Parameters:**
- `rule` - Rule to validate for compilation

**Returns:**
- `nil` if rule name is valid
- `ValidationError` if rule name contains parentheses

**Algorithm:**
1. Check if rule.Name contains '(' or ')'
2. If found → return error
3. Return nil if valid

**Pseudocode:**
```
function ValidateRuleForCompilation(rule):
    if rule.Name contains '(' or ')':
        return ValidationError{
            Field: "rule.name",
            Value: rule.Name,
            InvalidChar: "(" or ")",
            Message: "rule.name cannot contain parentheses: '" + rule.Name + "'"
        }
    
    return nil
```

**Rationale:**
- Enforcement headers use format: `# {Name} ({ENFORCEMENT})`
- Parentheses in rule name would break header parsing
- Example: `# Use (Smart) Names (MUST)` is ambiguous - is "(Smart)" part of name or enforcement?

### Allowed Characters

**Valid:** `a-z A-Z 0-9 - _`

**Invalid:** `/ \ : * ? " < > |` (and any other special characters)

**Rationale:**
- **Forward slash (/)** - Directory separator on Unix/Linux/macOS
- **Backslash (\\)** - Directory separator on Windows, escape character
- **Colon (:)** - Drive separator on Windows, reserved on macOS
- **Asterisk (*)** - Wildcard character in shells
- **Question mark (?)** - Wildcard character in shells
- **Double quote (")** - String delimiter in shells
- **Less than (<)** - Redirection operator in shells
- **Greater than (>)** - Redirection operator in shells
- **Pipe (|)** - Pipe operator in shells

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Empty ID | Return error "{field} cannot be empty" |
| ID with forward slash | Return error "{field} contains invalid character '/' in '{id}'" |
| ID with backslash | Return error "{field} contains invalid character '\\' in '{id}'" |
| ID with colon | Return error "{field} contains invalid character ':' in '{id}'" |
| ID with asterisk | Return error "{field} contains invalid character '*' in '{id}'" |
| ID with question mark | Return error "{field} contains invalid character '?' in '{id}'" |
| ID with double quote | Return error "{field} contains invalid character '"' in '{id}'" |
| ID with angle brackets | Return error "{field} contains invalid character '<' or '>' in '{id}'" |
| ID with pipe | Return error "{field} contains invalid character '\|' in '{id}'" |
| ID with spaces | Return error "{field} contains invalid character ' ' in '{id}'" |
| ID with only valid chars | Return nil (success) |
| Multiple invalid chars | Return error for first invalid character encountered |
| Rule name with opening paren | Return error "rule.name cannot contain parentheses: '{name}'" |
| Rule name with closing paren | Return error "rule.name cannot contain parentheses: '{name}'" |
| Rule name with both parens | Return error "rule.name cannot contain parentheses: '{name}'" |
| Rule name without parens | Return nil (success) |

## Dependencies

- Resource model from ai-resource-core-go (Ruleset, Rule, Promptset, Prompt structures)
- No external validation libraries required (simple character checking)

## Implementation Mapping

**Source files:**
- `internal/format/validation.go` - Implements `ValidateResourceIDs()` and helper functions
- `pkg/compiler/compiler.go` - Calls validation before compilation pipeline

**Related specs:**
- `compiler-architecture.md` - References validation as step 1 in compilation pipeline
- `metadata-block.md` - Assumes valid IDs for metadata generation
- All target compiler specs - Assume IDs are pre-validated

## Examples

### Example 1: Valid Rule IDs

**Input:**
```go
resource := Resource{
    Type: "rule",
    Ruleset: Ruleset{ID: "cleanCode"},
    Rule: Rule{ID: "meaningfulNames"},
}

err := ValidateResourceIDs(resource)
```

**Expected Output:**
```go
err == nil
```

**Verification:**
- Both IDs contain only lowercase letters
- No error returned

### Example 2: Valid IDs with Numbers and Hyphens

**Input:**
```go
resource := Resource{
    Type: "rule",
    Ruleset: Ruleset{ID: "clean-code-v2"},
    Rule: Rule{ID: "meaningful_names_2024"},
}

err := ValidateResourceIDs(resource)
```

**Expected Output:**
```go
err == nil
```

**Verification:**
- IDs contain letters, numbers, hyphens, underscores
- All characters valid

### Example 3: Invalid Character - Forward Slash

**Input:**
```go
resource := Resource{
    Type: "rule",
    Ruleset: Ruleset{ID: "clean/code"},
    Rule: Rule{ID: "meaningfulNames"},
}

err := ValidateResourceIDs(resource)
```

**Expected Output:**
```go
err.Error() == "ruleset.id contains invalid character '/' in 'clean/code'"
```

**Verification:**
- Forward slash detected in ruleset ID
- Error message identifies field, character, and value

### Example 4: Invalid Character - Asterisk

**Input:**
```go
resource := Resource{
    Type: "prompt",
    Promptset: Promptset{ID: "codeReview"},
    Prompt: Prompt{ID: "review*PR"},
}

err := ValidateResourceIDs(resource)
```

**Expected Output:**
```go
err.Error() == "prompt.id contains invalid character '*' in 'review*PR'"
```

**Verification:**
- Asterisk detected in prompt ID
- Error identifies prompt field specifically

### Example 5: Empty ID

**Input:**
```go
resource := Resource{
    Type: "rule",
    Ruleset: Ruleset{ID: "cleanCode"},
    Rule: Rule{ID: ""},
}

err := ValidateResourceIDs(resource)
```

**Expected Output:**
```go
err.Error() == "rule.id cannot be empty"
```

**Verification:**
- Empty ID detected
- Clear error message without character reference

### Example 6: Multiple Invalid Characters

**Input:**
```go
resource := Resource{
    Type: "rule",
    Ruleset: Ruleset{ID: "clean:code?"},
    Rule: Rule{ID: "meaningfulNames"},
}

err := ValidateResourceIDs(resource)
```

**Expected Output:**
```go
err.Error() == "ruleset.id contains invalid character ':' in 'clean:code?'"
```

**Verification:**
- First invalid character (colon) reported
- Error stops at first violation (doesn't report question mark)

### Example 7: Space Character

**Input:**
```go
resource := Resource{
    Type: "rule",
    Ruleset: Ruleset{ID: "clean code"},
    Rule: Rule{ID: "meaningfulNames"},
}

err := ValidateResourceIDs(resource)
```

**Expected Output:**
```go
err.Error() == "ruleset.id contains invalid character ' ' in 'clean code'"
```

**Verification:**
- Space character detected and rejected
- Error message shows space in quotes

### Example 8: Rule Name with Parentheses

**Input:**
```go
rule := Rule{
    ID: "meaningfulNames",
    Name: "Use (Smart) Names",
    Enforcement: "must",
}

err := ValidateRuleForCompilation(rule)
```

**Expected Output:**
```go
err.Error() == "rule.name cannot contain parentheses: 'Use (Smart) Names'"
```

**Verification:**
- Parentheses detected in rule name
- Clear error message explaining restriction

### Example 9: Valid Rule Name

**Input:**
```go
rule := Rule{
    ID: "meaningfulNames",
    Name: "Use Meaningful Names",
    Enforcement: "must",
}

err := ValidateRuleForCompilation(rule)
```

**Expected Output:**
```go
err == nil
```

**Verification:**
- No parentheses in rule name
- Validation passes

## Notes

**Design Rationale:**
- **Fail fast** - Validation occurs before any compilation work
- **Clear errors** - Users know exactly which ID and character is problematic
- **Cross-platform** - Rules ensure compatibility across Windows, macOS, Linux
- **Shell-safe** - Avoids characters with special meaning in shells

**Character Set Choice:**
- Alphanumeric + hyphen + underscore is widely compatible
- Hyphens common in kebab-case identifiers
- Underscores common in snake_case identifiers
- No dots to avoid confusion with file extensions
- No spaces to avoid quoting issues in shells

**Error Handling Strategy:**
- Return first error encountered (fail fast)
- Could be extended to collect all errors for batch reporting
- ValidationError struct allows programmatic error handling

## Known Issues

None - this is a new specification.

## Areas for Improvement

- Consider adding `ValidateAllResourceIDs([]Resource)` for batch validation
- Evaluate whether to allow dots (.) for namespacing (e.g., "org.cleanCode")
- Consider maximum length validation for IDs (filesystem limits)
- Explore Unicode support for international characters (requires careful testing)
