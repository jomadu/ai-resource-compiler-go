package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type arrayFlags []string

func (a *arrayFlags) String() string {
	return strings.Join(*a, ", ")
}

func (a *arrayFlags) Set(value string) error {
	*a = append(*a, value)
	return nil
}

func main() {
	var targets arrayFlags
	flag.Var(&targets, "target", "Target format to compile to (repeatable)")
	
	output := flag.String("output", "stdout", "Output mode: stdout or directory path")
	flat := flag.Bool("flat", false, "Disable target subdirectories in file output mode")
	help := flag.Bool("help", false, "Show help information")

	flag.Parse()

	if *help {
		printHelp()
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: resource file required")
		printUsage()
		os.Exit(1)
	}

	resourceFile := args[0]

	if len(targets) == 0 {
		fmt.Fprintln(os.Stderr, "Error: at least one target required")
		printUsage()
		os.Exit(1)
	}

	if _, err := os.Stat(resourceFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: resource file not found: %s\n", resourceFile)
		os.Exit(1)
	}

	validTargets := map[string]bool{
		"markdown": true,
		"kiro":     true,
		"cursor":   true,
		"claude":   true,
		"copilot":  true,
	}

	for _, target := range targets {
		if !validTargets[target] {
			fmt.Fprintf(os.Stderr, "Error: unknown target: %s\n\n", target)
			fmt.Fprintln(os.Stderr, "Valid targets: cursor, kiro, claude, copilot, markdown")
			os.Exit(1)
		}
	}

	if err := compile(resourceFile, targets, *output, *flat); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "\nUsage:")
	fmt.Fprintln(os.Stderr, "  arc [flags] <resource-file>")
	fmt.Fprintln(os.Stderr, "\nFlags:")
	fmt.Fprintln(os.Stderr, "  -target string   Target format (cursor, kiro, claude, copilot, markdown)")
	fmt.Fprintln(os.Stderr, "  -output string   Output mode: stdout or directory path (default \"stdout\")")
	fmt.Fprintln(os.Stderr, "  -flat            Disable target subdirectories in file output mode")
	fmt.Fprintln(os.Stderr, "  -help            Show help")
}

func printHelp() {
	fmt.Println("Compile AI resources to target-specific formats")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  arc [flags] <resource-file>")
	fmt.Println()
	fmt.Println("Arguments:")
	fmt.Println("  resource-file    Path to resource file (YAML or JSON)")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -target string   Target format to compile to (repeatable)")
	fmt.Println("                   Valid targets: cursor, kiro, claude, copilot, markdown")
	fmt.Println("  -output string   Output mode: \"stdout\" or directory path (default \"stdout\")")
	fmt.Println("  -flat            Disable target subdirectories in file output mode")
	fmt.Println("  -help            Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Compile to markdown, print to stdout")
	fmt.Println("  arc -target markdown resource.yaml")
	fmt.Println()
	fmt.Println("  # Compile to multiple targets, print to stdout")
	fmt.Println("  arc -target markdown -target kiro resource.yaml")
	fmt.Println()
	fmt.Println("  # Compile to cursor, write to target subdirectory")
	fmt.Println("  arc -target cursor -output .cursor/rules resource.yaml")
	fmt.Println()
	fmt.Println("  # Compile to cursor, write directly to directory (no subdirectory)")
	fmt.Println("  arc -target cursor -output .cursor/rules -flat resource.yaml")
	fmt.Println()
	fmt.Println("  # Compile to all targets, write to separate subdirectories")
	fmt.Println("  arc -target cursor -target kiro -target claude -target copilot -target markdown -output ./output resource.yaml")
}


