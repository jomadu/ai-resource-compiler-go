package compiler

import (
	"fmt"

	"github.com/jomadu/ai-resource-compiler-go/internal/format"
	"gopkg.in/yaml.v3"
)

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

// UnmarshalYAML implements custom YAML unmarshaling for Resource.
// It unmarshals Spec into the appropriate type based on Kind.
func (r *Resource) UnmarshalYAML(node *yaml.Node) error {
	// First unmarshal into a temporary struct to get Kind
	type rawResource struct {
		APIVersion string    `yaml:"apiVersion"`
		Kind       string    `yaml:"kind"`
		Metadata   struct {
			ID          string `yaml:"id"`
			Name        string `yaml:"name"`
			Description string `yaml:"description,omitempty"`
		} `yaml:"metadata"`
		Spec yaml.Node `yaml:"spec"`
	}

	var raw rawResource
	if err := node.Decode(&raw); err != nil {
		return err
	}

	r.APIVersion = raw.APIVersion
	r.Kind = raw.Kind
	r.Metadata.ID = raw.Metadata.ID

	// Unmarshal Spec based on Kind
	switch raw.Kind {
	case "Rule":
		var rule format.Rule
		if err := raw.Spec.Decode(&rule); err != nil {
			return fmt.Errorf("failed to decode Rule spec: %w", err)
		}
		// Copy metadata from top level
		rule.Metadata.ID = raw.Metadata.ID
		rule.Metadata.Name = raw.Metadata.Name
		rule.Metadata.Description = raw.Metadata.Description
		r.Spec = &rule
	case "Ruleset":
		var ruleset format.Ruleset
		if err := raw.Spec.Decode(&ruleset); err != nil {
			return fmt.Errorf("failed to decode Ruleset spec: %w", err)
		}
		// Copy metadata from top level
		ruleset.Metadata.ID = raw.Metadata.ID
		ruleset.Metadata.Name = raw.Metadata.Name
		ruleset.Metadata.Description = raw.Metadata.Description
		r.Spec = &ruleset
	case "Prompt":
		var prompt format.Prompt
		if err := raw.Spec.Decode(&prompt); err != nil {
			return fmt.Errorf("failed to decode Prompt spec: %w", err)
		}
		// Copy metadata from top level
		prompt.Metadata.ID = raw.Metadata.ID
		prompt.Metadata.Name = raw.Metadata.Name
		prompt.Metadata.Description = raw.Metadata.Description
		r.Spec = &prompt
	case "Promptset":
		var promptset format.Promptset
		if err := raw.Spec.Decode(&promptset); err != nil {
			return fmt.Errorf("failed to decode Promptset spec: %w", err)
		}
		// Copy metadata from top level
		promptset.Metadata.ID = raw.Metadata.ID
		promptset.Metadata.Name = raw.Metadata.Name
		promptset.Metadata.Description = raw.Metadata.Description
		r.Spec = &promptset
	default:
		return fmt.Errorf("unsupported kind: %s", raw.Kind)
	}

	return nil
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
