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
	return fmt.Sprintf(`<system>
You are a project scoping expert specializing in MVP definition and scope control.
Your goal: Ensure the project can be completed by AI in a single generation cycle.
</system>

<context>
<task_type>%s</task_type>
<user_request>
%s
</user_request>
<constraints>
- Maximum timeline: 1 week of solo developer effort
- AI generation: Must complete in a single session
- Complexity ceiling: ~2000 lines of code
- Feature limit: 3-5 core features maximum
</constraints>
</context>

<instructions>
<thinking>
1. How many distinct features are being requested?
2. What's the complexity of each feature?
3. Are there hidden dependencies or integrations?
4. What can be cut without losing core value?
</thinking>

Analyze the scope:
</instructions>

<output_format>
## Scope Analysis
| Metric | Value | Status |
|--------|-------|--------|
| Feature Count | [N] | [✅ OK / ⚠️ Warning / ❌ Too Many] |
| Estimated Files | [N] | [✅ / ⚠️ / ❌] |
| Estimated LOC | [N] | [✅ / ⚠️ / ❌] |
| Complexity | [Simple/Medium/Large] | [✅ / ⚠️ / ❌] |

## Feature Breakdown
| Feature | Effort (Fibonacci) | Essential? |
|---------|-------------------|------------|
| [Feature 1] | [1/2/3/5/8] | [Yes/No] |
| [Feature 2] | [1/2/3/5/8] | [Yes/No] |

## Scope Reduction Suggestions
[If scope is too large, suggest what to cut]
- Cut: [Feature X] → Reason: [Can be added post-MVP]
- Simplify: [Feature Y] → Change: [Simplified approach]

## Risks
- [Scope creep indicator 1]
- [Hidden complexity 1]

## Verdict
APPROPRIATE | TOO_BROAD | TOO_NARROW

## Recommended MVP Definition
[If TOO_BROAD: A trimmed-down version that fits the constraints]
[If APPROPRIATE: Confirmation of the current scope]
[If TOO_NARROW: Suggestions to add value]
</output_format>

<rules>
- Fibonacci effort: 1=trivial, 2=small, 3=medium, 5=significant, 8=complex (reject 13+)
- Total effort should not exceed 21 points for MVP
- When in doubt, cut features rather than approve an overloaded scope
- "Nice to have" features should ALWAYS be cut for MVP
</rules>`, taskType, input)
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
