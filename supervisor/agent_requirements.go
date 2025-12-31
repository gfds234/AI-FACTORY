package supervisor

import (
	"ai-studio/orchestrator/llm"
	"fmt"
	"strings"
	"time"
)

// RequirementsAgent validates requirement completeness
type RequirementsAgent struct {
	client *llm.Client
	model  string
}

// NewRequirementsAgent creates a new requirements agent
func NewRequirementsAgent(client *llm.Client, model string) *RequirementsAgent {
	return &RequirementsAgent{
		client: client,
		model:  model,
	}
}

func (a *RequirementsAgent) Name() string {
	return "requirements"
}

func (a *RequirementsAgent) RequiresInput() bool {
	return true // Runs before execution
}

func (a *RequirementsAgent) Execute(taskType, input string, context map[string]interface{}) (*AgentOutput, error) {
	start := time.Now()

	prompt := a.buildPrompt(taskType, input)
	response, err := a.client.Generate(a.model, prompt)
	if err != nil {
		return nil, fmt.Errorf("requirements agent failed: %w", err)
	}

	// Parse response to determine if requirements are complete
	status := a.parseStatus(response)

	return &AgentOutput{
		AgentType: "requirements",
		Status:    status,
		Output:    response,
		Duration:  time.Since(start).Seconds(),
		Timestamp: time.Now(),
	}, nil
}

func (a *RequirementsAgent) buildPrompt(taskType, input string) string {
	return fmt.Sprintf(`You are a requirements analysis agent. Evaluate if the following request has complete requirements for a %s task.

User Request:
%s

Analyze and respond in this format:

## Completeness Score
[Rate 1-10 how complete the requirements are]

## Missing Information
[List any critical missing details - be specific]
- [Item 1]
- [Item 2]

## Clarifying Questions
[Questions to ask the user to complete requirements]
1. [Question 1]
2. [Question 2]

## Status
[One word: COMPLETE, NEEDS_CLARIFICATION, or INCOMPLETE]

## Recommendation
[Should we proceed, or gather more information first?]

Be practical - don't demand perfection. Mark COMPLETE if we have enough to build an MVP.`, taskType, input)
}

func (a *RequirementsAgent) parseStatus(response string) string {
	// Simple status extraction
	if strings.Contains(response, "Status\nCOMPLETE") || strings.Contains(response, "## Status\nCOMPLETE") {
		return "passed"
	}
	if strings.Contains(response, "INCOMPLETE") {
		return "failed"
	}
	return "warning" // NEEDS_CLARIFICATION
}
