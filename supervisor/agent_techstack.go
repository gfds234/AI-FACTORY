package supervisor

import (
	"ai-studio/orchestrator/llm"
	"fmt"
	"strings"
	"time"
)

// TechStackAgent validates technology choices
type TechStackAgent struct {
	client *llm.Client
	model  string
}

func NewTechStackAgent(client *llm.Client, model string) *TechStackAgent {
	return &TechStackAgent{client: client, model: model}
}

func (a *TechStackAgent) Name() string {
	return "techstack"
}

func (a *TechStackAgent) RequiresInput() bool {
	return true
}

func (a *TechStackAgent) Execute(taskType, input string, context map[string]interface{}) (*AgentOutput, error) {
	start := time.Now()

	// Only run for code generation tasks
	if taskType != "code" {
		return &AgentOutput{
			AgentType: "techstack",
			Status:    "passed",
			Output:    "N/A - not a code generation task",
			Duration:  0,
			Timestamp: time.Now(),
		}, nil
	}

	prompt := a.buildPrompt(input)
	response, err := a.client.Generate(a.model, prompt)
	if err != nil {
		return nil, fmt.Errorf("tech stack agent failed: %w", err)
	}

	status := a.parseStatus(response)

	return &AgentOutput{
		AgentType: "techstack",
		Status:    status,
		Output:    response,
		Duration:  time.Since(start).Seconds(),
		Timestamp: time.Now(),
	}, nil
}

func (a *TechStackAgent) buildPrompt(input string) string {
	return fmt.Sprintf(`You are a senior tech lead reviewing technology choices. Analyze this project request and pre-approve or reject the tech stack.

Project Request:
%s

Evaluate and respond in this format:

## Inferred Tech Stack
[What technologies/frameworks would you use for this?]
- Language: [X]
- Framework: [Y]
- Database: [Z] (if needed)

## Tech Stack Assessment
[Is this stack appropriate for the requirements?]
- Pros: [List benefits]
- Cons: [List drawbacks]
- Alternatives: [Better options?]

## Solo Developer Friendliness
[Can one person build and maintain this? Rate 1-10]

## Production Readiness
[Will this stack produce production-quality code?]

## Verdict
[One word: APPROVED, NEEDS_REVISION, or REJECTED]

## Recommendation
[Final advice on tech stack choice]

Focus on practical, modern stacks that a solo developer can handle.`, input)
}

func (a *TechStackAgent) parseStatus(response string) string {
	if strings.Contains(response, "APPROVED") {
		return "passed"
	}
	if strings.Contains(response, "REJECTED") {
		return "failed"
	}
	return "warning" // NEEDS_REVISION
}
