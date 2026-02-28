package compiler

import (
	"fmt"
	"sync"
)

var (
	defaultCompiler     *Compiler
	defaultCompilerOnce sync.Once
)

// Compiler orchestrates compilation across multiple target formats.
type Compiler struct {
	targets map[Target]TargetCompiler
}

// NewCompiler creates a new compiler instance with all built-in targets registered.
func NewCompiler() *Compiler {
	defaultCompilerOnce.Do(func() {
		defaultCompiler = &Compiler{
			targets: make(map[Target]TargetCompiler),
		}
	})
	// Return a copy with the same registered targets
	c := &Compiler{
		targets: make(map[Target]TargetCompiler),
	}
	for k, v := range defaultCompiler.targets {
		c.targets[k] = v
	}
	return c
}

// RegisterTarget adds or replaces a target compiler.
func (c *Compiler) RegisterTarget(target Target, compiler TargetCompiler) error {
	if compiler == nil {
		return fmt.Errorf("compiler cannot be nil")
	}
	c.targets[target] = compiler
	return nil
}

// RegisterDefaultTarget registers a target compiler in the default compiler instance.
// This is used by target packages to register themselves during initialization.
func RegisterDefaultTarget(target Target, compiler TargetCompiler) {
	defaultCompilerOnce.Do(func() {
		defaultCompiler = &Compiler{
			targets: make(map[Target]TargetCompiler),
		}
	})
	defaultCompiler.targets[target] = compiler
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
