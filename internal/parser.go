package internal

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func ParseYAML(path string) (*yaml.Node, error) {
	filedata, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var node yaml.Node
	err = yaml.Unmarshal(filedata, &node)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &node, nil
}
