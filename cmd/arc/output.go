package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func outputStdout(allResults []targetResults) error {
	for _, tr := range allResults {
		for _, result := range tr.results {
			fmt.Printf("=== %s/%s ===\n", tr.target, result.Path)
			fmt.Println(result.Content)
			fmt.Println()
		}
	}
	return nil
}

func outputFiles(allResults []targetResults, outputDir string, flat bool) error {
	for _, tr := range allResults {
		for _, result := range tr.results {
			var filePath string
			if flat {
				filePath = filepath.Join(outputDir, result.Path)
			} else {
				filePath = filepath.Join(outputDir, tr.target, result.Path)
			}

			dir := filepath.Dir(filePath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dir, err)
			}

			if err := os.WriteFile(filePath, []byte(result.Content), 0644); err != nil {
				return fmt.Errorf("failed to write file %s: %w", filePath, err)
			}

			fmt.Fprintf(os.Stderr, "Wrote %s\n", filePath)
		}
	}
	return nil
}
