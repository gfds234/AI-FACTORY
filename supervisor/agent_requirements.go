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
	return fmt.Sprintf(`<system>
You are a senior requirements analyst specializing in MVP validation for AI-generated software projects.
Your goal: Ensure we have enough information to build a working MVP, not a perfect spec.
</system>

<context>
<task_type>%s</task_type>
<user_request>
%s
</user_request>
<market_context>
- Target: Solo developers or small teams
- Timeline: MVP should be buildable in under 1 week
- Stack: Modern, lightweight frameworks (avoid enterprise complexity)
- Deployment: Should be deployable on free tiers (Vercel, Railway, Fly.io)
</market_context>
</context>

<instructions>
Think step by step before providing your analysis:

<thinking>
1. What is the core user problem being solved?
2. What is the absolute minimum feature set for a usable MVP?
3. Are there any scope creep indicators (too many features, enterprise requirements)?
4. What assumptions can we safely make vs. what MUST be clarified?
</thinking>

Now provide your structured analysis:
</instructions>

<output_format>
## Completeness Score
[Rate 1-10: 7+ means we can proceed with reasonable assumptions]

## Core Problem Identified
[One sentence describing the user's actual need]

## MVP Feature Set
[Bullet list of ONLY the essential features - trim anything that isn't launch-critical]

## Missing Information
[List ONLY blockers - things we genuinely cannot assume]
- [Blocker 1]
- [Blocker 2]

## Safe Assumptions
[Things we can reasonably decide without asking]
- [Assumption 1]
- [Assumption 2]

## Clarifying Questions
[ONLY if truly blocking - keep to max 2 questions]
1. [Question 1]
2. [Question 2]

## Status
COMPLETE | NEEDS_CLARIFICATION | INCOMPLETE

## Recommendation
[One sentence: proceed, clarify, or reject with reason]
</output_format>

<rules>
- Bias toward COMPLETE - MVPs should ship fast
- Assume modern defaults (REST API, SQLite/PostgreSQL, JWT auth if needed)
- Flag scope creep aggressively
- Never ask for details that can be decided during implementation
</rules>`, taskType, input)
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
