package supervisor

import (
	"ai-studio/orchestrator/llm"
	"fmt"
	"strings"
	"time"
)

// ScopeAgent validates project scope
type ScopeAgent struct {
	client *llm.Client
	model  string
}

func NewScopeAgent(client *llm.Client, model string) *ScopeAgent {
	return &ScopeAgent{client: client, model: model}
}

func (a *ScopeAgent) Name() string {
	return "scope"
}

func (a *ScopeAgent) RequiresInput() bool {
	return true
}

func (a *ScopeAgent) Execute(taskType, input string, context map[string]interface{}) (*AgentOutput, error) {
	start := time.Now()

	prompt := a.buildPrompt(taskType, input)
	response, err := a.client.Generate(a.model, prompt)
	if err != nil {
		return nil, fmt.Errorf("scope agent failed: %w", err)
	}

	status := a.parseStatus(response)

	return &AgentOutput{
		AgentType: "scope",
		Status:    status,
		Output:    response,
		Duration:  time.Since(start).Seconds(),
		Timestamp: time.Now(),
	}, nil
}

func (a *ScopeAgent) buildPrompt(taskType, input string) string {
	return fmt.Sprintf(`You are a project scoping expert. Evaluate if this request has reasonable scope for a single AI-assisted task.

Task Type: %s
User Request:
%s

Analyze and respond:

## Scope Analysis
[How big is this project?]
- Estimated complexity: [Simple/Medium/Large/Too Large]
- Estimated files: [Number]
- Estimated LOC: [Rough estimate]

## Feasibility Assessment
[Can this be completed in a single AI generation?]
- Single iteration: [Yes/No]
- Needs breaking down: [Yes/No]

## Scope Issues
[Any scope creep or unrealistic expectations?]

## Recommendation
[How to scope this properly]

## Verdict
[One word: APPROPRIATE, TOO_BROAD, or TOO_NARROW]

Be realistic about what can be accomplished in a single AI task.`, taskType, input)
}

func (a *ScopeAgent) parseStatus(response string) string {
	if strings.Contains(response, "APPROPRIATE") {
		return "passed"
	}
	if strings.Contains(response, "TOO_BROAD") {
		return "failed"
	}
	return "warning" // TOO_NARROW
}
