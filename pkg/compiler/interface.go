package compiler

// Resource is a placeholder for ai-resource-core-go Resource type.
// TODO: Replace with actual import when ai-resource-core-go is implemented.
type Resource struct {
	APIVersion string
	Kind       string
	Metadata   struct {
		ID string
	}
	Spec interface{}
}

// TargetCompiler transforms resources into target-specific formats.
type TargetCompiler interface {
	// Name returns the target identifier (matches Target enum value).
	Name() string

	// SupportedVersions returns list of supported API versions (e.g., ["ai-resource/draft"]).
	SupportedVersions() []string

	// Compile transforms a resource into target-specific format(s).
	// Handles Rule, Ruleset, Prompt, Promptset kinds.
	// Expands collections (Ruleset/Promptset) into multiple results.
	// Returns one result per rule/prompt.
	Compile(resource *Resource) ([]CompilationResult, error)
}
