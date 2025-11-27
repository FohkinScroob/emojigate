package internal

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type GithubActionType string

const (
	Workflow GithubActionType = "Workflow"
	Job      GithubActionType = "Job"
	Step     GithubActionType = "Step"
)

const yamlKeyValuePairSize = 2

type Violation struct {
	Type       GithubActionType
	Identifier string
	Msg        string
}

func LintWorkflow(root *yaml.Node) ([]Violation, error) {
	violations := []Violation{}

	if root.Kind != yaml.DocumentNode || len(root.Content) == 0 {
		return nil, fmt.Errorf("invalid workflow: expected document node")
	}

	workflowRoot := root.Content[0]
	if workflowRoot.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("invalid workflow: root must be a mapping")
	}

	lintWorkflowName(workflowRoot, &violations)

	jobsNode, err := findJobsNode(workflowRoot)
	if err != nil {
		return nil, fmt.Errorf("jobs section not found in workflow")
	}

	if err := lintJobs(jobsNode, &violations); err != nil {
		return nil, err
	}

	return violations, nil
}

func lintWorkflowName(workflowRoot *yaml.Node, out *[]Violation) {
	name, err := getName(workflowRoot)
	if err != nil {
		*out = append(*out, Violation{
			Type:       Workflow,
			Identifier: "workflow",
			Msg:        "Missing display name. Please add a 'name:' field starting with an emoji.",
		})
		return
	}

	if !startsWithEmoji(name) {
		*out = append(*out, Violation{
			Type:       Workflow,
			Identifier: name,
			Msg:        "Name must start with an emoji. Example: '🚀 Deploy'",
		})
	}
}

func findJobsNode(workflowRoot *yaml.Node) (*yaml.Node, error) {
	for i := 0; i < len(workflowRoot.Content); i += yamlKeyValuePairSize {
		key := workflowRoot.Content[i]
		if key.Kind == yaml.ScalarNode && key.Value == "jobs" {
			if i+1 >= len(workflowRoot.Content) {
				return nil, fmt.Errorf("jobs key has no value")
			}
			return workflowRoot.Content[i+1], nil
		}
	}
	return nil, fmt.Errorf("jobs key not found")
}

func lintJobs(jobsNode *yaml.Node, out *[]Violation) error {
	if len(jobsNode.Content)%yamlKeyValuePairSize != 0 {
		return fmt.Errorf("jobs node content is not even. Every job must have configuration set")
	}

	for i := 0; i < len(jobsNode.Content); i += yamlKeyValuePairSize {
		jobName := jobsNode.Content[i].Value
		jobConfig := jobsNode.Content[i+1]

		name, err := getName(jobConfig)
		if err != nil {
			*out = append(*out, Violation{
				Type:       Job,
				Identifier: jobName,
				Msg:        "Missing display name. Please add a 'name:' field starting with an emoji.",
			})
			continue
		}

		if !startsWithEmoji(name) {
			*out = append(*out, Violation{
				Type:       Job,
				Identifier: jobName,
				Msg:        "Name must start with an emoji. Example: '🚀 Deploy'",
			})
		}

		if err := lintSteps(jobConfig, out); err != nil {
			return err
		}
	}
	return nil
}

func findStepsIndex(nodes []*yaml.Node) int {
	for index, node := range nodes {
		if node.Kind == yaml.ScalarNode && node.Value == "steps" {
			return index
		}
	}

	return -1
}

func lintSteps(stepsNode *yaml.Node, out *[]Violation) error {
	index := findStepsIndex(stepsNode.Content)
	if index == -1 {
		return nil
	}

	if len(stepsNode.Content) < index+yamlKeyValuePairSize {
		return fmt.Errorf("steps node content is not even. Every step must have configuration set")
	}

	for _, step := range stepsNode.Content[index+1].Content {
		err := lintStep(step, out)
		if err != nil {
			return err
		}
	}

	return nil
}

func lintStep(stepNode *yaml.Node, out *[]Violation) error {
	index, err := findNameIndex(stepNode.Content)
	if err != nil {
		return fmt.Errorf("step does not contain a name property")
	}

	nameNode := stepNode.Content[index+1]

	name := nameNode.Value
	if !startsWithEmoji(name) {
		*out = append(*out, Violation{
			Msg:        "Name must start with an emoji. Example: '🚀 Deploy'",
			Type:       Step,
			Identifier: name,
		})
	}

	return nil
}

func findNameIndex(nodes []*yaml.Node) (int, error) {
	for index, node := range nodes {
		if node.Kind == yaml.ScalarNode && node.Value == "name" {
			return index, nil
		}
	}

	return -1, fmt.Errorf("name field not found")
}

func getName(configNode *yaml.Node) (string, error) {
	index, err := findNameIndex(configNode.Content)
	if err != nil {
		return "", err
	}

	if len(configNode.Content) < index+yamlKeyValuePairSize {
		return "", fmt.Errorf("name field has no value")
	}

	nameNode := configNode.Content[index+1]
	return nameNode.Value, nil
}

func startsWithEmoji(s string) bool {
	r := []rune(s)
	if len(r) == 0 {
		return false
	}

	emoji := r[0]
	return (emoji >= 0x1F600 && emoji <= 0x1F64F) ||
		(emoji >= 0x1F300 && emoji <= 0x1F5FF) ||
		(emoji >= 0x1F680 && emoji <= 0x1F6FF) ||
		(emoji >= 0x1F1E0 && emoji <= 0x1F1FF) ||
		(emoji >= 0x2600 && emoji <= 0x26FF) ||
		(emoji >= 0x2700 && emoji <= 0x27BF) ||
		(emoji >= 0x1F900 && emoji <= 0x1F9FF) ||
		(emoji >= 0x1FA00 && emoji <= 0x1FA6F) ||
		(emoji >= 0x1F004 && emoji <= 0x1F0CF) ||
		(emoji >= 0x2300 && emoji <= 0x23FF) ||
		(emoji >= 0x2B00 && emoji <= 0x2BFF)
}
