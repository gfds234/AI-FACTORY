package project

import (
	"ai-studio/orchestrator/llm"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
)

// PlanGenerator generates structured implementation plans
type PlanGenerator struct {
	llmClient *llm.Client
	model     string
}

// NewPlanGenerator creates a new plan generator
func NewPlanGenerator(llmClient *llm.Client, model string) *PlanGenerator {
	return &PlanGenerator{
		llmClient: llmClient,
		model:     model,
	}
}

// GeneratePlan creates a structured implementation plan for a project
func (pg *PlanGenerator) GeneratePlan(project *Project) (*PlanDocument, error) {
	log.Printf("Plan Generator: Generating implementation plan for project %s", project.Name)

	// Build planning prompt
	prompt := pg.buildPlanningPrompt(project)

	// Determine thinking mode (use project metadata if set, default to normal)
	thinkingMode := string(project.Metadata.ThinkingMode)
	if thinkingMode == "" {
		thinkingMode = "normal"
	}

	log.Printf("Plan Generator: Using %s thinking mode for plan generation", thinkingMode)

	// Generate plan from LLM with appropriate thinking mode
	response, err := pg.llmClient.GenerateWithThinking(pg.model, prompt, thinkingMode)
	if err != nil {
		return nil, fmt.Errorf("failed to generate plan: %w", err)
	}

	// Parse the response into structured plan
	plan := pg.parsePlanResponse(response, project.ID)

	log.Printf("Plan Generator: Successfully generated plan with %d files to create, %d to modify",
		len(plan.FilesToCreate), len(plan.FilesToModify))

	return plan, nil
}

// buildPlanningPrompt creates the prompt for plan generation
func (pg *PlanGenerator) buildPlanningPrompt(project *Project) string {
	// Gather context from previous phases
	previousContext := ""
	for _, phaseExec := range project.Phases {
		if phaseExec.Status == PhaseStatusComplete {
			previousContext += fmt.Sprintf("\n### %s Phase Results:\n", phaseExec.Phase)
			for agentName, output := range phaseExec.AgentOutputs {
				previousContext += fmt.Sprintf("**%s:**\n%s\n\n", agentName, output)
			}
		}
	}

	prompt := fmt.Sprintf(`You are an expert software architect creating a detailed implementation plan.

PROJECT INFORMATION:
Name: %s
Description: %s

PREVIOUS ANALYSIS:
%s

TASK:
Create a comprehensive, structured implementation plan for this project. Your plan should be clear, actionable, and ready for a developer to execute.

OUTPUT FORMAT (use this EXACT structure):

## Implementation Approach
[Describe the overall strategy and architectural decisions in 2-3 paragraphs. Explain WHY you chose this approach.]

## Technical Stack
- Primary Language: [e.g., JavaScript, Python, Go]
- Framework: [e.g., React 18, Django, Gin]
- Build Tool: [e.g., Vite, Webpack, Go build]
- Testing Framework: [e.g., Vitest, pytest, go test]
- Additional Tools: [list any other key tools]

## Files to Create
[List each new file with a brief description of its purpose. Format: "- path/to/file.ext - Description"]
Example:
- src/App.jsx - Main React application component
- src/components/TaskList.jsx - Task list component with add/remove functionality
- package.json - Project dependencies and scripts
- vite.config.js - Vite build configuration
- src/App.test.jsx - Unit tests for App component

## Files to Modify
[If any existing files need modification, list them here. If starting from scratch, write "None - new project"]

## Testing Strategy
[Describe how the project will be tested. Include:
- What will be tested (components, functions, endpoints)
- Test coverage goals
- Testing approach (unit tests, integration tests)]

## Estimated Complexity
[Choose ONE: Low, Medium, or High]
Reasoning: [1-2 sentences explaining the complexity rating]

## Estimated Time
[e.g., "30 minutes", "2 hours", "4-6 hours"]

## Implementation Steps
1. [First step]
2. [Second step]
3. [Third step]
[List 5-8 concrete steps in order]

IMPORTANT:
- Be specific about file names and paths
- Ensure the plan is immediately actionable
- Consider build configuration, testing, and documentation
- For React/Vite projects, include: package.json, vite.config.js, index.html, src/main.jsx, src/App.jsx
- For Python projects, include: requirements.txt, main.py, tests/
- For Go projects, include: go.mod, main.go, *_test.go files

Generate the plan now:`, project.Name, project.Description, previousContext)

	return prompt
}

// parsePlanResponse extracts structured data from the LLM response
func (pg *PlanGenerator) parsePlanResponse(response string, projectID string) *PlanDocument {
	plan := &PlanDocument{
		ProjectID:   projectID,
		GeneratedAt: time.Now(),
		IsApproved:  false,
	}

	// Extract approach
	plan.Approach = pg.extractSection(response, "Implementation Approach")

	// Extract tech stack
	techStackSection := pg.extractSection(response, "Technical Stack")
	plan.TechStack = pg.parseListItems(techStackSection)

	// Extract files to create
	filesToCreateSection := pg.extractSection(response, "Files to Create")
	plan.FilesToCreate = pg.parseFileList(filesToCreateSection)

	// Extract files to modify
	filesToModifySection := pg.extractSection(response, "Files to Modify")
	if !strings.Contains(strings.ToLower(filesToModifySection), "none") {
		plan.FilesToModify = pg.parseFileList(filesToModifySection)
	}

	// Extract testing strategy
	plan.TestingStrategy = pg.extractSection(response, "Testing Strategy")

	// Extract complexity
	complexitySection := pg.extractSection(response, "Estimated Complexity")
	plan.Complexity = pg.extractComplexity(complexitySection)

	// Extract estimated time
	timeSection := pg.extractSection(response, "Estimated Time")
	plan.EstimatedTime = strings.TrimSpace(strings.Split(timeSection, "\n")[0])

	return plan
}

// extractSection extracts content between section headers
func (pg *PlanGenerator) extractSection(text string, sectionName string) string {
	// Match "## Section Name" or "**Section Name**" or "Section Name:"
	patterns := []string{
		`(?i)##\s*` + regexp.QuoteMeta(sectionName) + `\s*\n((?:(?!##).)*?)(?:\n##|$)`,
		`(?i)\*\*` + regexp.QuoteMeta(sectionName) + `\*\*\s*\n((?:(?!\*\*).)*?)(?:\n\*\*|$)`,
		`(?i)` + regexp.QuoteMeta(sectionName) + `:\s*\n((?:(?!\n[A-Z]).)*?)(?:\n[A-Z]|$)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}

	return ""
}

// parseListItems extracts list items from text
func (pg *PlanGenerator) parseListItems(text string) []string {
	lines := strings.Split(text, "\n")
	items := []string{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Remove list markers (-, *, •, 1., etc.)
		line = regexp.MustCompile(`^[-*•]\s*`).ReplaceAllString(line, "")
		line = regexp.MustCompile(`^\d+\.\s*`).ReplaceAllString(line, "")
		line = strings.TrimSpace(line)

		if line != "" {
			items = append(items, line)
		}
	}

	return items
}

// parseFileList extracts file paths from file listing
func (pg *PlanGenerator) parseFileList(text string) []string {
	lines := strings.Split(text, "\n")
	files := []string{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Remove list markers
		line = regexp.MustCompile(`^[-*•]\s*`).ReplaceAllString(line, "")
		line = regexp.MustCompile(`^\d+\.\s*`).ReplaceAllString(line, "")

		// Extract file path (before any " - " or " : " description)
		parts := regexp.MustCompile(`\s+-\s+|\s+:\s+`).Split(line, 2)
		if len(parts) > 0 {
			filePath := strings.TrimSpace(parts[0])
			if filePath != "" && !strings.HasPrefix(strings.ToLower(filePath), "none") {
				files = append(files, filePath)
			}
		}
	}

	return files
}

// extractComplexity extracts complexity rating
func (pg *PlanGenerator) extractComplexity(text string) string {
	text = strings.ToLower(text)

	if strings.Contains(text, "high") {
		return "High"
	} else if strings.Contains(text, "medium") {
		return "Medium"
	} else if strings.Contains(text, "low") {
		return "Low"
	}

	return "Medium" // Default
}
