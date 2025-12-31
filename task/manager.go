package task

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"ai-studio/orchestrator/config"
	"ai-studio/orchestrator/llm"
)

// Manager handles task execution and routing
type Manager struct {
	cfg        *config.Config
	client     *llm.Client
	history    []Result
	historyMux sync.RWMutex
	maxHistory int
}

// Result represents the output of a task execution
type Result struct {
	TaskType     string    `json:"task_type"`
	Input        string    `json:"input"`
	Output       string    `json:"output"`
	Model        string    `json:"model"`
	ArtifactPath string    `json:"artifact_path"`
	Duration     float64   `json:"duration_seconds"`
	Timestamp    time.Time `json:"timestamp"`
	Error        string    `json:"error,omitempty"`
}

// NewManager creates a new task manager
func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		cfg:        cfg,
		client:     llm.NewClient(cfg.OllamaURL, cfg.Timeout),
		history:    make([]Result, 0, 20),
		maxHistory: 20,
	}
}

// ExecuteTask routes and executes a task
func (m *Manager) ExecuteTask(taskType, input string) (*Result, error) {
	start := time.Now()
	result := &Result{
		TaskType:  taskType,
		Input:     input,
		Timestamp: start,
	}

	// Get model for this task type
	model, ok := m.cfg.Models[taskType]
	if !ok {
		result.Error = fmt.Sprintf("unknown task type: %s", taskType)
		return result, fmt.Errorf(result.Error)
	}
	result.Model = model

	// Build prompt based on task type
	prompt, err := m.buildPrompt(taskType, input)
	if err != nil {
		result.Error = err.Error()
		return result, err
	}

	// Execute with retries
	var output string
	var lastErr error
	
	for attempt := 0; attempt <= m.cfg.MaxRetries; attempt++ {
		output, lastErr = m.client.Generate(model, prompt)
		if lastErr == nil {
			break
		}
		
		if attempt < m.cfg.MaxRetries {
			time.Sleep(time.Second * time.Duration(attempt+1)) // Exponential backoff
		}
	}

	if lastErr != nil {
		result.Error = lastErr.Error()
		return result, lastErr
	}

	result.Output = output
	result.Duration = time.Since(start).Seconds()

	// Save artifact
	artifactPath, err := m.saveArtifact(result)
	if err != nil {
		// Non-fatal - still return the result
		result.Error = fmt.Sprintf("artifact save failed: %v", err)
	} else {
		result.ArtifactPath = artifactPath
	}

	// Add to history
	m.addToHistory(result)

	return result, nil
}

// buildPrompt constructs the prompt for each task type
func (m *Manager) buildPrompt(taskType, input string) (string, error) {
	switch taskType {
	case "validate":
		return m.buildValidationPrompt(input), nil
	case "review":
		return m.buildReviewPrompt(input), nil
	case "code":
		return m.buildCodePrompt(input), nil
	default:
		return "", fmt.Errorf("unknown task type: %s", taskType)
	}
}

// buildValidationPrompt creates prompt for idea validation
func (m *Manager) buildValidationPrompt(input string) string {
	return fmt.Sprintf(`You are a game design consultant. Analyze the following game idea and provide structured feedback.

Game Idea:
%s

Provide your analysis in the following format:

## Core Concept
[1-2 sentence summary of the idea]

## Strengths
- [List 2-3 key strengths]

## Potential Issues
- [List 2-3 concerns or challenges]

## Market Viability
[Brief assessment of target audience and market fit]

## Recommendation
[Clear recommendation: Proceed, Revise, or Reconsider]

## Next Steps
[2-3 concrete next steps if proceeding]

Be direct and honest. Focus on actionable insights.`, input)
}

// buildReviewPrompt creates prompt for architecture review
func (m *Manager) buildReviewPrompt(input string) string {
	return fmt.Sprintf(`You are a senior software architect and technical reviewer. Review the following architecture, code, or technical proposal.

Document to Review:
%s

Provide your review in the following format:

## Summary
[1-2 sentence overview of what's being reviewed]

## Tech Stack Analysis
- [Evaluate technology choices - are they appropriate?]
- [Are there better alternatives for this use case?]
- [Consider: performance, maintainability, ecosystem, learning curve]

## Architecture Assessment
- [Evaluate overall design and structure]
- [Identify architectural strengths]

## Risk Assessment
- [Technical risks, performance concerns, scalability issues]
- [Security considerations]
- [Maintenance and long-term viability]

## Code Quality (if applicable)
- [Code structure, readability, best practices]
- [Potential bugs or issues]

## Recommendations
- [Specific improvements with rationale]
- [Alternative approaches to consider]

## Standards Compliance
- [Does this meet production standards?]
- [What needs to change before approval?]

## Verdict
[Overall assessment: Approved, Approved with Changes, Needs Revision, or Rejected]

Be specific and technical. Focus on practical implications and real-world viability.`, input)
}

// buildCodePrompt creates prompt for code generation with intelligent tech stack selection
func (m *Manager) buildCodePrompt(input string) string {
	return fmt.Sprintf(`You are an expert full-stack software architect with deep knowledge across multiple domains:

**Game Development:**
- Unity/C#, Godot/GDScript, Unreal/C++
- 2D/3D engines, game mechanics, physics

**Mobile Development:**
- Flutter/Dart (cross-platform)
- React Native, Swift (iOS), Kotlin (Android)

**Web Development:**
- React/TypeScript, Vue, Next.js
- Node.js, Python FastAPI, Go backends

**Desktop Applications:**
- Electron, Python/Tkinter, C#/.NET, Rust

**Backend Services:**
- Go, Python, Node.js
- REST APIs, databases, authentication

Project Request:
%s

Your task:
1. **Analyze the request** to determine:
   - Project type (game, mobile app, web app, backend, desktop tool, etc.)
   - Complexity level (prototype, MVP, production)
   - Target platform(s)
   - Key requirements

2. **Select optimal tech stack:**
   - Choose the BEST language and framework for this specific use case
   - Prioritize: solo developer friendliness, modern ecosystem, cross-platform when beneficial
   - Consider: performance needs, learning curve, maintenance

3. **Generate production-quality code:**
   - Follow best practices for chosen language
   - Include comments explaining key decisions
   - Structure code clearly and maintainably
   - Include error handling where appropriate

4. **Provide complete output in this format:**

## Tech Stack Decision
**Project Type:** [Game/Mobile/Web/Backend/Desktop/etc.]
**Language:** [Chosen language]
**Framework/Engine:** [Chosen framework]
**Rationale:** [2-3 sentences explaining why this stack is optimal for this request]

## Implementation

` + "```" + `[language]
[Your complete, production-quality code here]
[Include comments explaining architecture and key decisions]
` + "```" + `

## Project Structure
[Explain file/folder organization if this would be multi-file]

## Setup Instructions
1. [Step-by-step instructions to run this code]
2. [Required dependencies/tools]
3. [How to test/verify it works]

## Next Steps
[What to implement next to expand this project]

Focus on practical, working code that a solo developer can immediately use and understand.`, input)
}

// saveArtifact saves the task result to a file
func (m *Manager) saveArtifact(result *Result) (string, error) {
	filename := fmt.Sprintf("%s_%d.md", 
		result.TaskType, 
		result.Timestamp.Unix())
	
	path := m.cfg.GetArtifactPath(filename)

	// Create artifact content
	content := fmt.Sprintf(`# Task Result: %s

**Timestamp:** %s
**Model:** %s
**Duration:** %.2fs

## Input
%s

## Output
%s
`, 
		result.TaskType,
		result.Timestamp.Format(time.RFC3339),
		result.Model,
		result.Duration,
		result.Input,
		result.Output,
	)

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write artifact: %w", err)
	}

	return path, nil
}

// Ping checks if the LLM backend is accessible
func (m *Manager) Ping() error {
	return m.client.Ping()
}

// addToHistory adds a task result to the history (ring buffer)
func (m *Manager) addToHistory(result *Result) {
	m.historyMux.Lock()
	defer m.historyMux.Unlock()

	m.history = append(m.history, *result)
	if len(m.history) > m.maxHistory {
		m.history = m.history[1:] // Remove oldest
	}
}

// GetHistory returns task history, optionally filtered by task type
func (m *Manager) GetHistory(taskType string) []Result {
	m.historyMux.RLock()
	defer m.historyMux.RUnlock()

	// Return in reverse chronological order
	results := make([]Result, 0, len(m.history))
	for i := len(m.history) - 1; i >= 0; i-- {
		result := m.history[i]
		if taskType == "all" || result.TaskType == taskType {
			results = append(results, result)
		}
	}
	return results
}

// GetClient returns the LLM client (for chat functionality)
func (m *Manager) GetClient() *llm.Client {
	return m.client
}
