package compiler

// Target represents a compilation target format.
type Target string

const (
	TargetCursor   Target = "cursor"
	TargetKiro     Target = "kiro"
	TargetClaude   Target = "claude"
	TargetCopilot  Target = "copilot"
	TargetMarkdown Target = "markdown"
)

// CompileOptions configures compilation behavior.
type CompileOptions struct {
	Targets []Target
}

// CompilationResult contains compiled output.
type CompilationResult struct {
	Path    string
	Content string
}
