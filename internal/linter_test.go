package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestLintWorkflows is an integration test that tests the entire linting pipeline
// It uses test fixtures (YAML files) and compares the output to expected results
func TestLintWorkflows(t *testing.T) {
	tests := []struct {
		name               string
		workflowFile       string
		expectedViolations int
		expectError        bool
		checkViolations    func(t *testing.T, violations []Violation)
	}{
		{
			name:               "valid workflow with all emojis",
			workflowFile:       "testdata/valid_workflow.yml",
			expectedViolations: 0,
			expectError:        false,
			checkViolations: func(t *testing.T, violations []Violation) {
				if len(violations) != 0 {
					t.Errorf("Expected no violations, got %d", len(violations))
					for _, v := range violations {
						t.Logf("  - [%s] %s: %s", v.Type, v.Identifier, v.Msg)
					}
				}
			},
		},
		{
			name:               "invalid workflow missing emojis",
			workflowFile:       "testdata/invalid_workflow.yml",
			expectedViolations: 7, // 1 workflow + 2 jobs + 4 steps without emojis
			expectError:        false,
			checkViolations: func(t *testing.T, violations []Violation) {
				if len(violations) < 7 {
					t.Errorf("Expected at least 7 violations, got %d", len(violations))
				}

				// Count by type
				workflowViolations := 0
				jobViolations := 0
				stepViolations := 0
				for _, v := range violations {
					if v.Type == Workflow {
						workflowViolations++
					} else if v.Type == Job {
						jobViolations++
					} else if v.Type == Step {
						stepViolations++
					}
				}

				if workflowViolations < 1 {
					t.Errorf("Expected at least 1 workflow violation, got %d", workflowViolations)
				}
				if jobViolations < 2 {
					t.Errorf("Expected at least 2 job violations, got %d", jobViolations)
				}
				if stepViolations < 4 {
					t.Errorf("Expected at least 4 step violations, got %d", stepViolations)
				}
			},
		},
		{
			name:               "workflow with missing job name",
			workflowFile:       "testdata/missing_job_name.yml",
			expectedViolations: 1, // 1 job without name field
			expectError:        false,
			checkViolations: func(t *testing.T, violations []Violation) {
				if len(violations) == 0 {
					t.Error("Expected at least 1 violation for missing job name")
					return
				}

				// Check that we have a violation for missing name
				foundMissingName := false
				for _, v := range violations {
					if v.Type == Job && v.Identifier == "build" {
						foundMissingName = true
						if v.Msg != "Missing display name. Please add a 'name:' field starting with an emoji." {
							t.Errorf("Unexpected message: %s", v.Msg)
						}
					}
				}

				if !foundMissingName {
					t.Error("Expected violation for 'build' job with missing name")
				}
			},
		},
		{
			name:               "workflow with mixed emoji types",
			workflowFile:       "testdata/mixed_emojis.yml",
			expectedViolations: 0,
			expectError:        false,
		},
		{
			name:               "workflow without name field",
			workflowFile:       "testdata/missing_workflow_name.yml",
			expectedViolations: 1, // 1 workflow without name
			expectError:        false,
			checkViolations: func(t *testing.T, violations []Violation) {
				if len(violations) != 1 {
					t.Errorf("Expected 1 violation, got %d", len(violations))
					return
				}

				v := violations[0]
				if v.Type != Workflow {
					t.Errorf("Expected Workflow violation, got %s", v.Type)
				}
				if v.Identifier != "workflow" {
					t.Errorf("Expected identifier 'workflow', got %s", v.Identifier)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the workflow file
			node, err := ParseYAML(tt.workflowFile)
			if err != nil {
				if !tt.expectError {
					t.Fatalf("Failed to parse workflow: %v", err)
				}
				return
			}

			// Run the linter
			violations, err := LintWorkflow(node)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
				return
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Custom validation if provided
			if tt.checkViolations != nil {
				tt.checkViolations(t, violations)
			} else {
				// Default: just check count
				if len(violations) != tt.expectedViolations {
					t.Errorf("Expected %d violations, got %d", tt.expectedViolations, len(violations))
					for _, v := range violations {
						t.Logf("  - [%s] %s: %s", v.Type, v.Identifier, v.Msg)
					}
				}
			}
		})
	}
}

// TestLintWorkflows_Snapshots uses golden files to compare exact violation outputs
// This is useful for regression testing - if output changes, you review and update the golden file
func TestLintWorkflows_Snapshots(t *testing.T) {
	tests := []struct {
		name         string
		workflowFile string
		goldenFile   string
	}{
		{
			name:         "invalid workflow snapshot",
			workflowFile: "testdata/invalid_workflow.yml",
			goldenFile:   "testdata/invalid_workflow.golden.json",
		},
		{
			name:         "missing job name snapshot",
			workflowFile: "testdata/missing_job_name.yml",
			goldenFile:   "testdata/missing_job_name.golden.json",
		},
		{
			name:         "missing workflow name snapshot",
			workflowFile: "testdata/missing_workflow_name.yml",
			goldenFile:   "testdata/missing_workflow_name.golden.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse and lint
			node, err := ParseYAML(tt.workflowFile)
			if err != nil {
				t.Fatalf("Failed to parse workflow: %v", err)
			}

			violations, err := LintWorkflow(node)
			if err != nil {
				t.Fatalf("Linting failed: %v", err)
			}

			// Serialize violations to JSON
			actual, err := json.MarshalIndent(violations, "", "  ")
			if err != nil {
				t.Fatalf("Failed to marshal violations: %v", err)
			}

			// Check if we should update golden files
			if os.Getenv("UPDATE_GOLDEN") == "1" {
				err := os.WriteFile(tt.goldenFile, actual, 0644)
				if err != nil {
					t.Fatalf("Failed to write golden file: %v", err)
				}
				t.Logf("Updated golden file: %s", tt.goldenFile)
				return
			}

			// Read golden file
			expected, err := os.ReadFile(tt.goldenFile)
			if err != nil {
				t.Fatalf("Failed to read golden file %s: %v\nRun with UPDATE_GOLDEN=1 to create it", tt.goldenFile, err)
			}

			// Compare
			if string(actual) != string(expected) {
				t.Errorf("Output differs from golden file %s\n\nActual:\n%s\n\nExpected:\n%s\n\nTo update: UPDATE_GOLDEN=1 go test",
					tt.goldenFile, string(actual), string(expected))
			}
		})
	}
}

// TestStartsWithEmoji tests only the emoji detection logic (kept minimal for edge cases)
func TestStartsWithEmoji(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"🚀 Deploy", true},
		{"✅ Check", true},
		{"☀️ Weather", true},
		{"⭐ Star", true},
		{"🐹 Gopher", true},
		{"Deploy", false},
		{"", false},
		{" 🚀 Deploy", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := startsWithEmoji(tt.input)
			if result != tt.expected {
				t.Errorf("startsWithEmoji(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestEndToEnd simulates the full application flow
func TestEndToEnd(t *testing.T) {
	// Create a temporary workflow file
	tmpDir := t.TempDir()
	workflowPath := filepath.Join(tmpDir, "test-workflow.yml")

	workflow := `name: 🧪 Test Workflow
on: push
jobs:
  test:
    name: 🧪 Run Tests
    runs-on: ubuntu-latest
    steps:
      - name: 📥 Checkout
        uses: actions/checkout@v3
      - name: ✅ Test
        run: go test ./...
`

	err := os.WriteFile(workflowPath, []byte(workflow), 0644)
	if err != nil {
		t.Fatalf("Failed to write test workflow: %v", err)
	}

	// Parse it
	node, err := ParseYAML(workflowPath)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	// Lint it
	violations, err := LintWorkflow(node)
	if err != nil {
		t.Fatalf("Linting failed: %v", err)
	}

	// Should have no violations
	if len(violations) != 0 {
		t.Errorf("Expected no violations for valid workflow, got %d:", len(violations))
		for _, v := range violations {
			t.Logf("  - [%s] %s: %s", v.Type, v.Identifier, v.Msg)
		}
	}
}
