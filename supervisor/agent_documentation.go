package supervisor

import (
	"ai-studio/orchestrator/llm"
	"fmt"
	"time"
)

// DocumentationAgent generates README and API docs
type DocumentationAgent struct {
	client *llm.Client
	model  string
}

func NewDocumentationAgent(client *llm.Client, model string) *DocumentationAgent {
	return &DocumentationAgent{client: client, model: model}
}

func (a *DocumentationAgent) Name() string {
	return "documentation"
}

func (a *DocumentationAgent) RequiresInput() bool {
	return false
}

func (a *DocumentationAgent) Execute(taskType, input string, context map[string]interface{}) (*AgentOutput, error) {
	start := time.Now()

	output, ok := context["output"].(string)
	if !ok {
		return nil, fmt.Errorf("no output in context")
	}

	prompt := a.buildPrompt(taskType, input, output)
	response, err := a.client.Generate(a.model, prompt)
	if err != nil {
		return nil, fmt.Errorf("documentation agent failed: %w", err)
	}

	return &AgentOutput{
		AgentType: "documentation",
		Status:    "passed",
		Output:    response,
		Duration:  time.Since(start).Seconds(),
		Timestamp: time.Now(),
	}, nil
}

func (a *DocumentationAgent) buildPrompt(taskType, input, output string) string {
	return fmt.Sprintf(`You are a technical writer. Create comprehensive documentation for this code.

Original Request:
%s

Generated Code:
%s

Create:

## README.md
` + "```markdown" + `
[Full README with: project overview, installation, usage, examples, configuration]
` + "```" + `

## API Documentation
[If applicable - document all public functions/endpoints]

## Code Comments
[Suggest inline comments for complex sections]

## Setup Guide
[Step-by-step setup instructions]

## Troubleshooting
[Common issues and solutions]

Make documentation clear, complete, and beginner-friendly.`, input, output)
}
