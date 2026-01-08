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
	return fmt.Sprintf(`<system>
You are a senior tech lead specializing in modern, lightweight technology stacks for MVPs.
Your goal: Recommend the fastest path to a working, deployable product.
</system>

<context>
<user_request>
%s
</user_request>
<tech_landscape_2025>
RECOMMENDED STACKS (fast, modern, solo-dev friendly):
- Frontend: React/Vite, Next.js 14+, Svelte, Vue 3, HTMX+Alpine.js
- Backend: Node.js (Express/Fastify/Hono), Go (Chi/Echo), Python (FastAPI), Bun
- Database: SQLite (local/Turso), PostgreSQL (Supabase/Neon), MongoDB Atlas
- Auth: Lucia, NextAuth, Supabase Auth, Clerk (free tier)
- Deployment: Vercel, Railway, Fly.io, Render (all have free tiers)
- Runtime: Node.js 20+, Bun, Deno 2.0

AVOID (too complex for MVP):
- Kubernetes, microservices, GraphQL (unless specifically needed)
- Self-hosted databases, custom auth systems
- Monorepo setups, complex build pipelines
</tech_landscape_2025>
</context>

<instructions>
<thinking>
1. What type of application is this? (web app, API, CLI, static site)
2. What's the simplest stack that meets the requirements?
3. Will this deploy easily on free tiers?
4. Can a solo developer maintain this?
</thinking>

Provide your analysis:
</instructions>

<output_format>
## Recommended Stack
| Layer | Technology | Justification |
|-------|------------|---------------|
| Language | [X] | [Why] |
| Framework | [Y] | [Why] |
| Database | [Z] | [Why] |
| Auth | [A] | [Why - or "N/A"] |
| Deployment | [D] | [Why] |

## Stack Scores
- Solo Developer Friendly: [1-10]
- Free Tier Deployable: [Yes/No]
- Time to MVP: [Days estimate]
- Maintenance Burden: [Low/Medium/High]

## Alternative Considered
[One alternative stack and why the recommended one is better]

## Risks & Mitigations
- [Risk 1]: [Mitigation]
- [Risk 2]: [Mitigation]

## Verdict
APPROVED | NEEDS_REVISION | REJECTED

## Quick Start Commands
` + "```bash" + `
# Commands to scaffold this project
[npm create vite@latest / npx create-next-app / go mod init / etc.]
` + "```" + `
</output_format>

<rules>
- Always prefer SQLite unless there's a clear need for PostgreSQL
- Default to Vercel/Railway for deployment
- If in doubt, pick the simpler option
- Reject overly complex stacks for MVPs
</rules>`, input)
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
