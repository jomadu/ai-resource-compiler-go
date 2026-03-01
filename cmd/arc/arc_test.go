package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func createTestResource(t *testing.T, dir string) string {
	content := `apiVersion: ai-resource/draft
kind: Rule
metadata:
  id: testRule
  name: Test Rule
spec:
  enforcement: must
  body: Test rule body
`
	path := filepath.Join(dir, "test.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test resource: %v", err)
	}
	return path
}

func TestCompileStdoutSingleTarget(t *testing.T) {
	dir := t.TempDir()
	resourceFile := createTestResource(t, dir)

	err := compile(resourceFile, []string{"markdown"}, "stdout", false)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestCompileStdoutMultipleTargets(t *testing.T) {
	dir := t.TempDir()
	resourceFile := createTestResource(t, dir)

	err := compile(resourceFile, []string{"markdown", "kiro"}, "stdout", false)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestCompileFilesSingleTarget(t *testing.T) {
	dir := t.TempDir()
	resourceFile := createTestResource(t, dir)
	outputDir := filepath.Join(dir, "output")

	err := compile(resourceFile, []string{"markdown"}, outputDir, false)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expectedPath := filepath.Join(outputDir, "markdown", "testRule.md")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Expected file not created: %s", expectedPath)
	}
}

func TestCompileFilesMultipleTargets(t *testing.T) {
	dir := t.TempDir()
	resourceFile := createTestResource(t, dir)
	outputDir := filepath.Join(dir, "output")

	err := compile(resourceFile, []string{"markdown", "kiro"}, outputDir, false)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expectedPaths := []string{
		filepath.Join(outputDir, "markdown", "testRule.md"),
		filepath.Join(outputDir, "kiro", "testRule.md"),
	}

	for _, path := range expectedPaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected file not created: %s", path)
		}
	}
}

func TestCompileFilesFlat(t *testing.T) {
	dir := t.TempDir()
	resourceFile := createTestResource(t, dir)
	outputDir := filepath.Join(dir, "output")

	err := compile(resourceFile, []string{"markdown"}, outputDir, true)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	expectedPath := filepath.Join(outputDir, "testRule.md")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Expected file not created: %s", expectedPath)
	}

	unexpectedPath := filepath.Join(outputDir, "markdown", "testRule.md")
	if _, err := os.Stat(unexpectedPath); !os.IsNotExist(err) {
		t.Errorf("Unexpected target subdirectory created with --flat: %s", unexpectedPath)
	}
}

func TestCompileErrorMissingFile(t *testing.T) {
	err := compile("nonexistent.yaml", []string{"markdown"}, "stdout", false)
	if err == nil {
		t.Fatal("Expected error for missing file, got nil")
	}
	if !strings.Contains(err.Error(), "failed to read resource file") {
		t.Errorf("Expected 'failed to read resource file' error, got: %v", err)
	}
}

func TestCompileErrorInvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "invalid.yaml")
	if err := os.WriteFile(path, []byte("invalid: yaml: content:"), 0644); err != nil {
		t.Fatalf("Failed to create invalid YAML: %v", err)
	}

	err := compile(path, []string{"markdown"}, "stdout", false)
	if err == nil {
		t.Fatal("Expected error for invalid YAML, got nil")
	}
	if !strings.Contains(err.Error(), "failed to parse resource file") {
		t.Errorf("Expected 'failed to parse resource file' error, got: %v", err)
	}
}

func TestCompileErrorUnknownTarget(t *testing.T) {
	dir := t.TempDir()
	resourceFile := createTestResource(t, dir)

	err := compile(resourceFile, []string{"invalid"}, "stdout", false)
	if err == nil {
		t.Fatal("Expected error for unknown target, got nil")
	}
	if !strings.Contains(err.Error(), "unknown target") {
		t.Errorf("Expected 'unknown target' error, got: %v", err)
	}
}
