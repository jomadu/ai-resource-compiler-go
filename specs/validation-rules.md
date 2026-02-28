# Resource Validation Rules

## Job to be Done
Ensure resource IDs are filesystem-safe before compilation to prevent path traversal, invalid filenames, and cross-platform compatibility issues.

## Activities
1. Validate collection IDs (metadata.id) for filesystem safety
2. Validate item IDs (map keys in rules/prompts) for filesystem safety
3. Validate rule names for compilation compatibility (no parentheses)
4. Reject resources with invalid characters in IDs
5. Return clear error messages identifying invalid IDs and characters

## Acceptance Criteria
- [ ] ValidateID function checks individual ID strings
- [ ] ValidateRuleName function checks rule names for parentheses
- [ ] Only alphanumeric characters, hyphens, and underscores allowed (a-z, A-Z, 0-9, -, _)
- [ ] Invalid characters rejected: / \ : * ? " < > |
- [ ] Error messages identify which ID is invalid and which character caused rejection
- [ ] Empty IDs rejected with clear error message
- [ ] Validation occurs during compilation pipeline

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

### ValidateID

```go
func ValidateID(id string) error
```

**Parameters:**
- `id` - ID string to validate (from metadata.id or map key)

**Returns:**
- `nil` if ID is valid
- `error` if ID is invalid

**Algorithm:**
1. Check if empty → return error
2. Check each character against allowed set
3. If invalid character found → return error with details
4. Return nil if valid

**Pseudocode:**
```
function ValidateID(id):
    if id is empty:
        return error("ID cannot be empty")
    
    for each character in id:
        if character not in [a-z, A-Z, 0-9, -, _]:
            return error("ID contains invalid character '" + character + "' in '" + id + "'")
    
    return nil
```

### ValidateRuleName

```go
func ValidateRuleName(name string) error
```

**Parameters:**
- `name` - Rule name to validate

**Returns:**
- `nil` if rule name is valid
- `error` if rule name contains parentheses

**Algorithm:**
1. Check if name contains '(' or ')'
2. If found → return error
3. Return nil if valid

**Pseudocode:**
```
function ValidateRuleName(name):
    if name contains '(' or ')':
        return error("rule name cannot contain parentheses: '" + name + "'")
    
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
| Empty ID | Return error "ID cannot be empty" |
| ID with forward slash | Return error "ID contains invalid character '/' in '{id}'" |
| ID with backslash | Return error "ID contains invalid character '\\' in '{id}'" |
| ID with colon | Return error "ID contains invalid character ':' in '{id}'" |
| ID with asterisk | Return error "ID contains invalid character '*' in '{id}'" |
| ID with question mark | Return error "ID contains invalid character '?' in '{id}'" |
| ID with double quote | Return error "ID contains invalid character '"' in '{id}'" |
| ID with angle brackets | Return error "ID contains invalid character '<' or '>' in '{id}'" |
| ID with pipe | Return error "ID contains invalid character '\|' in '{id}'" |
| ID with spaces | Return error "ID contains invalid character ' ' in '{id}'" |
| ID with only valid chars | Return nil (success) |
| Multiple invalid chars | Return error for first invalid character encountered |
| Rule name with opening paren | Return error "rule name cannot contain parentheses: '{name}'" |
| Rule name with closing paren | Return error "rule name cannot contain parentheses: '{name}'" |
| Rule name with both parens | Return error "rule name cannot contain parentheses: '{name}'" |
| Rule name without parens | Return nil (success) |

## Dependencies

- No external dependencies (simple character checking)

## Implementation Mapping

**Source files:**
- `internal/format/validation.go` - Implements `ValidateID()` and `ValidateRuleName()` functions
- `pkg/compiler/compiler.go` - Calls validation during compilation pipeline

**Related specs:**
- `compiler-architecture.md` - References validation in compilation pipeline
- `metadata-block.md` - Assumes valid IDs for metadata generation
- All target compiler specs - Assume IDs are validated during compilation

## Examples

### Example 1: Valid ID

**Input:**
```go
err := ValidateID("cleanCode")
```

**Expected Output:**
```go
err == nil
```

**Verification:**
- ID contains only lowercase letters
- No error returned

### Example 2: Valid ID with Numbers and Hyphens

**Input:**
```go
err := ValidateID("clean-code-v2")
```

**Expected Output:**
```go
err == nil
```

**Verification:**
- ID contains letters, numbers, hyphens
- All characters valid

### Example 3: Invalid Character - Forward Slash

**Input:**
```go
err := ValidateID("clean/code")
```

**Expected Output:**
```go
err.Error() == "ID contains invalid character '/' in 'clean/code'"
```

**Verification:**
- Forward slash detected
- Error message identifies character and value

### Example 4: Invalid Character - Asterisk

**Input:**
```go
err := ValidateID("review*PR")
```

**Expected Output:**
```go
err.Error() == "ID contains invalid character '*' in 'review*PR'"
```

**Verification:**
- Asterisk detected
- Error identifies character

### Example 5: Empty ID

**Input:**
```go
err := ValidateID("")
```

**Expected Output:**
```go
err.Error() == "ID cannot be empty"
```

**Verification:**
- Empty ID detected
- Clear error message

### Example 6: Multiple Invalid Characters

**Input:**
```go
err := ValidateID("clean:code?")
```

**Expected Output:**
```go
err.Error() == "ID contains invalid character ':' in 'clean:code?'"
```

**Verification:**
- First invalid character (colon) reported
- Error stops at first violation

### Example 7: Space Character

**Input:**
```go
err := ValidateID("clean code")
```

**Expected Output:**
```go
err.Error() == "ID contains invalid character ' ' in 'clean code'"
```

**Verification:**
- Space character detected and rejected
- Error message shows space

### Example 8: Rule Name with Parentheses

**Input:**
```go
err := ValidateRuleName("Use (Smart) Names")
```

**Expected Output:**
```go
err.Error() == "rule name cannot contain parentheses: 'Use (Smart) Names'"
```

**Verification:**
- Parentheses detected in rule name
- Clear error message explaining restriction

### Example 9: Valid Rule Name

**Input:**
```go
err := ValidateRuleName("Use Meaningful Names")
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
- **Simple functions** - Validate individual IDs and names, not entire resources
- **Called during compilation** - Validation happens as items are processed
- **Clear errors** - Users know exactly which ID or name is problematic
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
- Simple error strings, not complex error types
- Validation integrated into compilation pipeline

## Known Issues

None - this is a new specification.

## Areas for Improvement

- Consider maximum length validation for IDs (filesystem limits)
- Explore Unicode support for international characters (requires careful testing)
- Add validation for reserved names (e.g., "CON", "PRN" on Windows)
