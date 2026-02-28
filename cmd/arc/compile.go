package main

import (
	"fmt"
	"os"

	"github.com/jomadu/ai-resource-compiler-go/pkg/compiler"
	"gopkg.in/yaml.v3"
)

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
	opts := compiler.CompileOptions{Targets: targetEnums}
	results, err := c.Compile(&resource, opts)
	if err != nil {
		return fmt.Errorf("compilation failed: %w", err)
	}

	if output == "stdout" {
		return outputStdout(results, targets)
	}
	return outputFiles(results, targets, output, flat)
}

func outputStdout(results []compiler.CompilationResult, targets []string) error {
	for _, result := range results {
		fmt.Printf("=== %s ===\n", result.Path)
		fmt.Println(result.Content)
		fmt.Println()
	}
	return nil
}

func outputFiles(results []compiler.CompilationResult, targets []string, outputDir string, flat bool) error {
	// Placeholder - will be implemented in TASK-016
	return fmt.Errorf("file output not yet implemented")
}
