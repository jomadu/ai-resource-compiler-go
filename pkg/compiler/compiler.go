package compiler

import "fmt"

// Compiler orchestrates compilation across multiple target formats.
type Compiler struct {
	targets map[Target]TargetCompiler
}

// NewCompiler creates a new compiler instance.
// Built-in targets will be registered here once implemented.
func NewCompiler() *Compiler {
	return &Compiler{
		targets: make(map[Target]TargetCompiler),
	}
}

// RegisterTarget adds or replaces a target compiler.
func (c *Compiler) RegisterTarget(target Target, compiler TargetCompiler) error {
	if compiler == nil {
		return fmt.Errorf("compiler cannot be nil")
	}
	c.targets[target] = compiler
	return nil
}

// Compile transforms a resource into one or more target formats.
func (c *Compiler) Compile(resource *Resource, opts CompileOptions) ([]CompilationResult, error) {
	// Step 1: Validate resource
	if resource.APIVersion == "" {
		return nil, fmt.Errorf("missing apiVersion")
	}
	if resource.Kind == "" {
		return nil, fmt.Errorf("missing kind")
	}
	if resource.Metadata.ID == "" {
		return nil, fmt.Errorf("missing metadata.id")
	}

	// Step 2: Validate options
	if len(opts.Targets) == 0 {
		return nil, fmt.Errorf("no targets specified")
	}

	// Step 3: Compile for each target
	var results []CompilationResult
	for _, target := range opts.Targets {
		compiler, ok := c.targets[target]
		if !ok {
			return nil, fmt.Errorf("unknown target: %s", target)
		}

		// Check version compatibility
		supported := false
		for _, version := range compiler.SupportedVersions() {
			if version == resource.APIVersion {
				supported = true
				break
			}
		}
		if !supported {
			return nil, fmt.Errorf("target %s does not support apiVersion: %s", target, resource.APIVersion)
		}

		// Compile resource
		targetResults, err := compiler.Compile(resource)
		if err != nil {
			return nil, err
		}
		results = append(results, targetResults...)
	}

	// Step 4: Return aggregated results
	return results, nil
}
