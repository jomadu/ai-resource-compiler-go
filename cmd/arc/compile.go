package main

import (
	"fmt"
	"os"

	"github.com/jomadu/ai-resource-compiler-go/pkg/compiler"
	"gopkg.in/yaml.v3"
)

type targetResults struct {
	target  string
	results []compiler.CompilationResult
}

func compile(resourceFile string, targets []string, output string, flat bool) error {
	data, err := os.ReadFile(resourceFile)
	if err != nil {
		return fmt.Errorf("failed to read resource file: %w", err)
	}

	var resource compiler.Resource
	if err := yaml.Unmarshal(data, &resource); err != nil {
		return fmt.Errorf("failed to parse resource file: %w", err)
	}

	targetEnums := make([]compiler.Target, len(targets))
	for i, t := range targets {
		switch t {
		case "markdown":
			targetEnums[i] = compiler.TargetMarkdown
		case "kiro":
			targetEnums[i] = compiler.TargetKiro
		case "cursor":
			targetEnums[i] = compiler.TargetCursor
		case "claude":
			targetEnums[i] = compiler.TargetClaude
		case "copilot":
			targetEnums[i] = compiler.TargetCopilot
		default:
			return fmt.Errorf("unknown target: %s", t)
		}
	}

	c := compiler.NewCompiler()
	
	// Compile each target separately to track which results belong to which target
	var allResults []targetResults
	
	for i, targetEnum := range targetEnums {
		opts := compiler.CompileOptions{Targets: []compiler.Target{targetEnum}}
		results, err := c.Compile(&resource, opts)
		if err != nil {
			return fmt.Errorf("compilation failed for target %s: %w", targets[i], err)
		}
		allResults = append(allResults, targetResults{target: targets[i], results: results})
	}

	if output == "stdout" {
		return outputStdout(allResults)
	}
	return outputFiles(allResults, output, flat)
}
