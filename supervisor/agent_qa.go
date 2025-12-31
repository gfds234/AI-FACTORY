package supervisor

import (
	"ai-studio/orchestrator/llm"
	"fmt"
	"time"
)

// QAAgent reviews code quality, bugs, and security
type QAAgent struct {
	client *llm.Client
	model  string
}

func NewQAAgent(client *llm.Client, model string) *QAAgent {
	return &QAAgent{client: client, model: model}
}

func (a *QAAgent) Name() string {
	return "qa"
}

func (a *QAAgent) RequiresInput() bool {
	return false // Runs after execution
}

func (a *QAAgent) Execute(taskType, input string, context map[string]interface{}) (*AgentOutput, error) {
	start := time.Now()

	// Get original output from context
	output, ok := context["output"].(string)
	if !ok {
		return nil, fmt.Errorf("no output in context")
	}

	prompt := a.buildPrompt(taskType, input, output)
	response, err := a.client.Generate(a.model, prompt)
	if err != nil {
		return nil, fmt.Errorf("qa agent failed: %w", err)
	}

	return &AgentOutput{
		AgentType: "qa",
		Status:    "passed", // QA always runs, doesn't block
		Output:    response,
		Duration:  time.Since(start).Seconds(),
		Timestamp: time.Now(),
	}, nil
}

func (a *QAAgent) buildPrompt(taskType, input, output string) string {
	return fmt.Sprintf(`You are a senior QA engineer. Review the following AI-generated output for quality issues.

Original Request:
%s

Generated Output:
%s

Review and report:

## Code Quality Issues
[List any code smells, anti-patterns, or bad practices]

## Potential Bugs
[Identify possible bugs or edge cases not handled]

## Security Concerns
[Any security vulnerabilities or risks?]

## Best Practices
[What best practices are missing?]

## Performance Considerations
[Any performance issues?]

## Overall Assessment
[Rate quality 1-10 and provide summary]

Be thorough but practical. Focus on issues that would impact production use.`, input, output)
}
