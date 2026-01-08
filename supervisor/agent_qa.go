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
	return fmt.Sprintf(`<system>
You are a senior QA engineer and security reviewer.
Your goal: Identify issues that would cause production failures or security breaches.
</system>

<context>
<original_request>
%s
</original_request>
<generated_output>
%s
</generated_output>
<review_standards>
- OWASP Top 10 2024: Injection, Broken Auth, Sensitive Data Exposure, XXE, 
  Broken Access Control, Security Misconfiguration, XSS, Insecure Deserialization,
  Vulnerable Components, Insufficient Logging
- Code Quality: DRY, SOLID, error handling, input validation
- Performance: N+1 queries, memory leaks, unnecessary re-renders
</review_standards>
</context>

<instructions>
<thinking>
1. What could cause this code to crash in production?
2. What security vulnerabilities are present?
3. What would a code reviewer flag in a PR?
4. What edge cases aren't handled?
</thinking>

Provide your review:
</instructions>

<output_format>
## Critical Issues (Must Fix)
| Issue | Location | Severity | Fix |
|-------|----------|----------|-----|
| [Issue] | [File/Line] | CRITICAL | [How to fix] |

## Security Review
| OWASP Category | Status | Notes |
|----------------|--------|-------|
| Injection | ✅/⚠️/❌ | [Details] |
| Broken Auth | ✅/⚠️/❌ | [Details] |
| Sensitive Data | ✅/⚠️/❌ | [Details] |
| Access Control | ✅/⚠️/❌ | [Details] |
| XSS | ✅/⚠️/❌ | [Details] |

## Code Quality Issues
| Issue | Severity | Recommendation |
|-------|----------|----------------|
| [Issue] | HIGH/MEDIUM/LOW | [Fix] |

## Edge Cases Not Handled
- [Edge case 1]
- [Edge case 2]

## Dependency Concerns
[Any known vulnerabilities in dependencies? Outdated packages?]

## Performance Concerns
- [Performance issue 1]
- [Performance issue 2]

## Overall Assessment
- Quality Score: [1-10]
- Production Ready: [Yes/No/With Fixes]
- Security Score: [1-10]

## Summary
[2-3 sentence summary of the most important findings]
</output_format>

<rules>
- CRITICAL = will cause crashes or security breach (must fix before deploy)
- HIGH = significant issues (should fix)
- MEDIUM = code quality issues (nice to fix)
- LOW = minor suggestions (optional)
- Be specific about locations and fixes
- Don't flag style preferences, focus on real issues
</rules>`, input, output)
}
