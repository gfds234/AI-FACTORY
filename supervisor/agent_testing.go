package supervisor

import (
	"ai-studio/orchestrator/llm"
	"fmt"
	"time"
)

// TestingAgent generates test plans and unit tests
type TestingAgent struct {
	client *llm.Client
	model  string
}

func NewTestingAgent(client *llm.Client, model string) *TestingAgent {
	return &TestingAgent{client: client, model: model}
}

func (a *TestingAgent) Name() string {
	return "testing"
}

func (a *TestingAgent) RequiresInput() bool {
	return false
}

func (a *TestingAgent) Execute(taskType, input string, context map[string]interface{}) (*AgentOutput, error) {
	start := time.Now()

	output, ok := context["output"].(string)
	if !ok {
		return nil, fmt.Errorf("no output in context")
	}

	prompt := a.buildPrompt(taskType, input, output)
	response, err := a.client.Generate(a.model, prompt)
	if err != nil {
		return nil, fmt.Errorf("testing agent failed: %w", err)
	}

	return &AgentOutput{
		AgentType: "testing",
		Status:    "passed",
		Output:    response,
		Duration:  time.Since(start).Seconds(),
		Timestamp: time.Now(),
	}, nil
}

func (a *TestingAgent) buildPrompt(taskType, input, output string) string {
	return fmt.Sprintf(`You are a test engineer. Create a comprehensive test plan for the generated code.

Original Request:
%s

Generated Code/Output:
%s

Provide:

## Test Strategy
[What types of tests are needed?]

## Unit Tests
[Generate sample unit tests for critical functions]
` + "```" + `
[Actual test code here]
` + "```" + `

## Integration Tests
[What integration tests would you recommend?]

## Edge Cases
[List edge cases to test]

## Test Data
[Sample test data needed]

## Testing Checklist
- [ ] [Test item 1]
- [ ] [Test item 2]

Generate practical, runnable tests.`, input, output)
}
