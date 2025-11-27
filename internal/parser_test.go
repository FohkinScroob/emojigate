package internal

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestParseYAML_ValidFile tests parsing a valid YAML file
func TestParseYAML_ValidFile(t *testing.T) {
	node, err := ParseYAML("testdata/valid_workflow.yml")
	if err != nil {
		t.Fatalf("ParseYAML() failed with valid file: %v", err)
	}

	if node == nil {
		t.Fatal("ParseYAML() returned nil node")
	}

	if node.Kind == 0 {
		t.Error("ParseYAML() returned empty node")
	}

	// Check that we have a document node
	if node.Kind != yaml.DocumentNode {
		t.Errorf("ParseYAML() node kind = %v, want %v", node.Kind, yaml.DocumentNode)
	}
}

// TestParseYAML_InvalidFile tests parsing a non-existent file
func TestParseYAML_InvalidFile(t *testing.T) {
	_, err := ParseYAML("testdata/non_existent_file.yml")
	if err == nil {
		t.Error("ParseYAML() should return error for non-existent file")
	}
}

// TestParseYAML_InvalidYAML tests parsing an invalid YAML file
func TestParseYAML_InvalidYAML(t *testing.T) {
	// Create a temporary invalid YAML file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "invalid.yml")

	// Write invalid YAML content
	invalidYAML := `
name: Test
jobs:
  build:
    name: Build
    invalid yaml structure here: [unclosed
`
	err := os.WriteFile(tmpFile, []byte(invalidYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	_, err = ParseYAML(tmpFile)
	if err == nil {
		t.Error("ParseYAML() should return error for invalid YAML")
	}
}

// TestParseYAML_EmptyFile tests parsing an empty file
func TestParseYAML_EmptyFile(t *testing.T) {
	// Create a temporary empty file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "empty.yml")

	err := os.WriteFile(tmpFile, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	node, err := ParseYAML(tmpFile)
	if err != nil {
		t.Fatalf("ParseYAML() failed with empty file: %v", err)
	}

	if node == nil {
		t.Fatal("ParseYAML() returned nil node")
	}

	// Empty file results in a zero-value node (Kind == 0)
	// This is expected behavior from yaml.Unmarshal with empty content
	if node.Kind != 0 {
		t.Errorf("ParseYAML() node kind = %v, want 0 (zero-value) for empty file", node.Kind)
	}
}

// TestParseYAML_WithComplexStructure tests parsing complex workflow structures
func TestParseYAML_WithComplexStructure(t *testing.T) {
	// Create a temporary file with complex structure
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "complex.yml")

	complexYAML := `
name: 🎯 Complex Workflow
on:
  push:
    branches: [main, develop]
  pull_request:
    types: [opened, synchronize]

env:
  GO_VERSION: '1.21'

jobs:
  build:
    name: 🔨 Build
    runs-on: ubuntu-latest
    strategy:
      matrix:
        os: [ubuntu-latest, windows-latest, macos-latest]
    steps:
      - name: 📥 Checkout
        uses: actions/checkout@v3
      - name: 🐹 Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: ${{ env.GO_VERSION }}
`

	err := os.WriteFile(tmpFile, []byte(complexYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	node, err := ParseYAML(tmpFile)
	if err != nil {
		t.Fatalf("ParseYAML() failed with complex structure: %v", err)
	}

	if node == nil {
		t.Fatal("ParseYAML() returned nil node")
	}

	if node.Kind != yaml.DocumentNode {
		t.Errorf("ParseYAML() node kind = %v, want %v", node.Kind, yaml.DocumentNode)
	}

	// Verify that the structure contains expected nodes
	if len(node.Content) == 0 {
		t.Error("ParseYAML() parsed complex YAML but has no content")
	}
}

// TestParseYAML_Permissions tests parsing with various file permissions
func TestParseYAML_Permissions(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping permission test in CI environment")
	}

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.yml")

	simpleYAML := `name: 🧪 Test`
	err := os.WriteFile(tmpFile, []byte(simpleYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Test with read permissions
	node, err := ParseYAML(tmpFile)
	if err != nil {
		t.Errorf("ParseYAML() failed with readable file: %v", err)
	}

	if node == nil {
		t.Fatal("ParseYAML() returned nil node")
	}

	if node.Kind == 0 {
		t.Error("ParseYAML() returned empty node")
	}
}
