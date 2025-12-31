package task

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"ai-studio/orchestrator/config"
	"ai-studio/orchestrator/llm"
)

// Manager handles task execution and routing
type Manager struct {
	cfg    *config.Config
	client *llm.Client
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
		cfg:    cfg,
		client: llm.NewClient(cfg.OllamaURL, cfg.Timeout),
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

	return result, nil
}

// buildPrompt constructs the prompt for each task type
func (m *Manager) buildPrompt(taskType, input string) (string, error) {
	switch taskType {
	case "validate":
		return m.buildValidationPrompt(input), nil
	case "review":
		return m.buildReviewPrompt(input), nil
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
	return fmt.Sprintf(`You are a senior game architect. Review the following technical architecture and provide detailed feedback.

Architecture Document:
%s

Provide your review in the following format:

## Architecture Summary
[1-2 sentence overview]

## Strengths
- [List key architectural strengths]

## Risk Assessment
- [Identify technical risks, performance concerns, or scalability issues]

## Recommendations
- [Specific improvements or alternatives]

## Implementation Notes
[Any important considerations for implementation]

## Verdict
[Overall assessment: Approved, Needs Revision, or Major Concerns]

Be specific and technical. Focus on practical implications.`, input)
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
