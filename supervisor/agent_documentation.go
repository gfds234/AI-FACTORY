package supervisor

import (
	"ai-studio/orchestrator/llm"
	"fmt"
	"time"
)

// DocumentationAgent generates README and API docs
type DocumentationAgent struct {
	client *llm.Client
	model  string
}

func NewDocumentationAgent(client *llm.Client, model string) *DocumentationAgent {
	return &DocumentationAgent{client: client, model: model}
}

func (a *DocumentationAgent) Name() string {
	return "documentation"
}

func (a *DocumentationAgent) RequiresInput() bool {
	return false
}

func (a *DocumentationAgent) Execute(taskType, input string, context map[string]interface{}) (*AgentOutput, error) {
	start := time.Now()

	output, ok := context["output"].(string)
	if !ok {
		return nil, fmt.Errorf("no output in context")
	}

	prompt := a.buildPrompt(taskType, input, output)
	response, err := a.client.Generate(a.model, prompt)
	if err != nil {
		return nil, fmt.Errorf("documentation agent failed: %w", err)
	}

	return &AgentOutput{
		AgentType: "documentation",
		Status:    "passed",
		Output:    response,
		Duration:  time.Since(start).Seconds(),
		Timestamp: time.Now(),
	}, nil
}

func (a *DocumentationAgent) buildPrompt(taskType, input, output string) string {
	return fmt.Sprintf(`<system>
You are a technical writer specializing in developer documentation and open-source README files.
Your goal: Create documentation that helps users get started in under 5 minutes.
</system>

<context>
<original_request>
%s
</original_request>
<generated_code>
%s
</generated_code>
<documentation_standards>
- README must have: Quick start, installation, usage examples
- Use shields.io badges for visual appeal
- Include copy-paste ready code snippets
- SEO: Clear title, description, keywords for discoverability
</documentation_standards>
</context>

<instructions>
<thinking>
1. What's the one-sentence pitch for this project?
2. What does a new user need to do to get it running?
3. What are the most common use cases?
4. What will go wrong and how do users fix it?
</thinking>

Generate documentation:
</instructions>

<output_format>
## README.md
`+"```markdown"+`
# [Project Name]

<!-- Badges -->
![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Version](https://img.shields.io/badge/version-1.0.0-green.svg)
[Add relevant badges: build status, coverage, npm/pypi version if applicable]

> [One-line description that explains what this does and why it matters]

## ✨ Features

- [Feature 1]
- [Feature 2]
- [Feature 3]

## 🚀 Quick Start

`+"```bash"+`
# Clone and run in 30 seconds
[git clone / npm install / pip install commands]
[single command to run]
`+"```"+`

## 📦 Installation

### Prerequisites
- [Prerequisite 1]
- [Prerequisite 2]

### Install
`+"```bash"+`
[Installation commands]
`+"```"+`

## 📖 Usage

### Basic Example
`+"```[language]"+`
[Simple, working code example]
`+"```"+`

### API Reference
[If applicable: endpoints, functions, methods]

## ⚙️ Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| [VAR_1] | [Description] | [Default] |

## 🚀 Deployment

### [Platform 1: Vercel/Railway/etc.]
`+"```bash"+`
[One-click deploy or commands]
`+"```"+`

## 🔧 Troubleshooting

### [Common Issue 1]
**Problem:** [Description]
**Solution:** [Fix]

## 📄 License

MIT © [Year]
`+"```"+`

## API_DOCS.md (if applicable)
`+"```markdown"+`
# API Documentation

## Endpoints

### [Method] /[endpoint]
[Description, parameters, response examples]
`+"```"+`

## Setup Instructions
[Step-by-step for developers joining the project]
</output_format>

<rules>
- README should be scannable in 30 seconds
- Every code block must be copy-paste ready
- Include environment variable examples with placeholder values
- Never expose real secrets or API keys
- Use emojis sparingly for section headers
</rules>`, input, output)
}
