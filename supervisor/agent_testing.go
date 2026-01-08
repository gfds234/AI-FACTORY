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
	return fmt.Sprintf(`<system>
You are a senior test engineer specializing in automated testing strategies.
Your goal: Create runnable tests that verify the code actually works.
</system>

<context>
<original_request>
%s
</original_request>
<generated_code>
%s
</generated_code>
<testing_frameworks_2025>
JavaScript/TypeScript:
- Jest (most common), Vitest (Vite projects), Playwright (E2E)
- Run: npm test, npx vitest, npx playwright test

Python:
- pytest (modern standard), unittest (built-in)
- Run: pytest, python -m pytest

Go:
- testing package (built-in), testify (assertions)
- Run: go test ./...

Coverage targets: 70%% for MVP, 80%% for production
</testing_frameworks_2025>
</context>

<instructions>
<thinking>
1. What testing framework should be used based on the tech stack?
2. What are the critical paths that MUST be tested?
3. What edge cases could cause failures?
4. What's the minimum viable test suite?
</thinking>

Generate a test plan:
</instructions>

<output_format>
## Detected Stack & Framework
| Aspect | Value |
|--------|-------|
| Language | [JavaScript/Python/Go] |
| Test Framework | [Jest/Vitest/pytest/go test] |
| Test Runner Command | [exact command] |

## Test Strategy
| Type | Count | Priority |
|------|-------|----------|
| Unit Tests | [N] | HIGH |
| Integration Tests | [N] | MEDIUM |
| E2E Tests | [N] | LOW (for MVP) |

## Critical Test Cases
| Test | What It Verifies | Priority |
|------|------------------|----------|
| [test_name] | [What it tests] | HIGH/MED/LOW |

## Generated Tests
`+"```"+`[language]
// File: [test_file_name]

[COMPLETE, RUNNABLE TEST CODE]
[Include imports, setup, teardown]
[Tests should be copy-paste ready]
`+"```"+`

## Test Commands
`+"```bash"+`
# Install test dependencies
[npm install --save-dev jest / pip install pytest / etc.]

# Run tests
[npm test / pytest / go test ./...]

# Run with coverage
[npm test -- --coverage / pytest --cov / go test -cover ./...]
`+"```"+`

## Edge Cases to Test
- [Edge case 1]
- [Edge case 2]
- [Edge case 3]

## Coverage Recommendation
- MVP Target: 70%%
- Files to prioritize: [list critical files]
</output_format>

<rules>
- Tests MUST be runnable - no pseudocode
- Include all necessary imports and setup
- Focus on business logic, not boilerplate
- Minimum 3 test cases for any feature
- Always include at least one error/edge case test
</rules>`, input, output)
}
