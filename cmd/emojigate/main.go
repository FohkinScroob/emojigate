package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/FohkinScroob/emojigate/internal"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "workflows":
		lintWorkflowsDirectory()
	case "lint":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Error: 'lint' command requires at least one file argument")
			printUsage()
			os.Exit(1)
		}
		lintFiles(os.Args[2:])
	case "help", "-h", "--help":
		printUsage()
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command '%s'\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`emojigate - Lint GitHub Actions workflows for emoji usage

Usage:
  emojigate workflows          Lint all workflow files in .github/workflows/
  emojigate lint <file>...     Lint specific workflow file(s)
  emojigate help               Show this help message

Examples:
  emojigate workflows
  emojigate lint .github/workflows/ci.yml
  emojigate lint .github/workflows/ci.yml .github/workflows/release.yml

Pre-commit Hook:
  Add to .pre-commit-config.yaml:
    - repo: local
      hooks:
        - id: emojigate
          name: Lint workflow emojis
          entry: emojigate lint
          language: system
          files: ^\.github/workflows/.*\.ya?ml$`)
}

func lintWorkflowsDirectory() {
	workflowsDir := ".github/workflows"

	if _, err := os.Stat(workflowsDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: directory '%s' does not exist\n", workflowsDir)
		os.Exit(1)
	}

	entries, err := os.ReadDir(workflowsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading directory '%s': %v\n", workflowsDir, err)
		os.Exit(1)
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".yml") || strings.HasSuffix(name, ".yaml") {
			files = append(files, filepath.Join(workflowsDir, name))
		}
	}

	if len(files) == 0 {
		fmt.Printf("No workflow files found in %s\n", workflowsDir)
		os.Exit(0)
	}

	lintFiles(files)
}

func lintFiles(files []string) {
	type fileViolations struct {
		file       string
		violations []internal.Violation
	}

	var allViolations []fileViolations
	totalViolations := 0

	for _, file := range files {
		node, err := internal.ParseYAML(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", file, err)
			os.Exit(1)
		}

		violations, err := internal.LintWorkflow(node)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error linting %s: %v\n", file, err)
			os.Exit(1)
		}

		if len(violations) > 0 {
			allViolations = append(allViolations, fileViolations{
				file:       file,
				violations: violations,
			})
			totalViolations += len(violations)
		}
	}

	if totalViolations == 0 {
		fmt.Printf("✅ All %d workflow(s) passed!\n", len(files))
		os.Exit(0)
	}

	fmt.Fprintf(os.Stderr, "❌ Found %d violation(s) across %d file(s):\n\n", totalViolations, len(allViolations))

	for _, fv := range allViolations {
		fmt.Fprintf(os.Stderr, "File: %s\n", fv.file)
		for _, v := range fv.violations {
			fmt.Fprintf(os.Stderr, "  [%s] %s\n", v.Type, v.Identifier)
			fmt.Fprintf(os.Stderr, "    → %s\n", v.Msg)
		}
		fmt.Fprintln(os.Stderr)
	}

	fmt.Fprintln(os.Stderr, "❗ Please add an emoji at the beginning of each workflow, job, and step name.")
	os.Exit(1)
}
