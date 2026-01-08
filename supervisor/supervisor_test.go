package supervisor

import (
	"strings"
	"testing"
)

// TestRequirementsAgentPromptBuild tests the requirements agent prompt construction
func TestRequirementsAgentPromptBuild(t *testing.T) {
	agent := &RequirementsAgent{}
	
	prompt := agent.buildPrompt("code", "Build a todo app with React")
	
	// Verify XML structure
	if !strings.Contains(prompt, "<system>") {
		t.Error("Prompt should contain <system> tag")
	}
	if !strings.Contains(prompt, "<context>") {
		t.Error("Prompt should contain <context> tag")
	}
	if !strings.Contains(prompt, "<instructions>") {
		t.Error("Prompt should contain <instructions> tag")
	}
	if !strings.Contains(prompt, "<thinking>") {
		t.Error("Prompt should contain <thinking> tag for chain-of-thought")
	}
	if !strings.Contains(prompt, "<output_format>") {
		t.Error("Prompt should contain <output_format> tag")
	}
	
	// Verify content inclusion
	if !strings.Contains(prompt, "todo app") {
		t.Error("Prompt should contain user input")
	}
	if !strings.Contains(prompt, "code") {
		t.Error("Prompt should contain task type")
	}
}

// TestRequirementsAgentParseStatus tests status parsing
func TestRequirementsAgentParseStatus(t *testing.T) {
	agent := &RequirementsAgent{}
	
	tests := []struct {
		name     string
		response string
		expected string
	}{
		{
			name:     "Complete status",
			response: "## Status\nCOMPLETE",
			expected: "passed",
		},
		{
			name:     "Incomplete status",
			response: "The requirements are INCOMPLETE",
			expected: "failed",
		},
		{
			name:     "Needs clarification",
			response: "Status: NEEDS_CLARIFICATION",
			expected: "warning",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := agent.parseStatus(tt.response)
			if result != tt.expected {
				t.Errorf("parseStatus() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestTechStackAgentPromptBuild tests tech stack prompt construction
func TestTechStackAgentPromptBuild(t *testing.T) {
	agent := &TechStackAgent{}
	
	prompt := agent.buildPrompt("Build an API with user authentication")
	
	// Verify 2025 tech recommendations are included
	if !strings.Contains(prompt, "Hono") {
		t.Error("Prompt should mention Hono as a modern framework")
	}
	if !strings.Contains(prompt, "SQLite") {
		t.Error("Prompt should mention SQLite as preferred database")
	}
	if !strings.Contains(prompt, "Vercel") || !strings.Contains(prompt, "Railway") {
		t.Error("Prompt should mention free tier deployment options")
	}
	
	// Verify structure
	if !strings.Contains(prompt, "<tech_landscape_2025>") {
		t.Error("Prompt should contain tech landscape section")
	}
}

// TestTechStackAgentParseStatus tests status parsing
func TestTechStackAgentParseStatus(t *testing.T) {
	agent := &TechStackAgent{}
	
	tests := []struct {
		name     string
		response string
		expected string
	}{
		{
			name:     "Approved",
			response: "## Verdict\nAPPROVED",
			expected: "passed",
		},
		{
			name:     "Rejected",
			response: "This stack is REJECTED due to complexity",
			expected: "failed",
		},
		{
			name:     "Needs revision",
			response: "NEEDS_REVISION",
			expected: "warning",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := agent.parseStatus(tt.response)
			if result != tt.expected {
				t.Errorf("parseStatus() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestScopeAgentPromptBuild tests scope agent prompt construction
func TestScopeAgentPromptBuild(t *testing.T) {
	agent := &ScopeAgent{}
	
	prompt := agent.buildPrompt("code", "Build an e-commerce platform")
	
	// Verify constraints are present
	if !strings.Contains(prompt, "1 week") {
		t.Error("Prompt should contain 1-week timeline constraint")
	}
	if !strings.Contains(prompt, "2000 lines") {
		t.Error("Prompt should contain LOC constraint")
	}
	if !strings.Contains(prompt, "Fibonacci") {
		t.Error("Prompt should mention Fibonacci estimation")
	}
	
	// Verify structure
	if !strings.Contains(prompt, "<constraints>") {
		t.Error("Prompt should contain constraints section")
	}
}

// TestScopeAgentParseStatus tests status parsing
func TestScopeAgentParseStatus(t *testing.T) {
	agent := &ScopeAgent{}
	
	tests := []struct {
		name     string
		response string
		expected string
	}{
		{
			name:     "Appropriate scope",
			response: "This scope is APPROPRIATE",
			expected: "passed",
		},
		{
			name:     "Too broad",
			response: "The scope is TOO_BROAD",
			expected: "failed",
		},
		{
			name:     "Too narrow",
			response: "TOO_NARROW - add more features",
			expected: "warning",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := agent.parseStatus(tt.response)
			if result != tt.expected {
				t.Errorf("parseStatus() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestQAAgentPromptBuild tests QA agent prompt construction
func TestQAAgentPromptBuild(t *testing.T) {
	agent := &QAAgent{}
	
	prompt := agent.buildPrompt("code", "original request", "generated code output")
	
	// Verify OWASP references
	if !strings.Contains(prompt, "OWASP") {
		t.Error("Prompt should reference OWASP")
	}
	if !strings.Contains(prompt, "Injection") {
		t.Error("Prompt should mention Injection vulnerability")
	}
	if !strings.Contains(prompt, "XSS") {
		t.Error("Prompt should mention XSS")
	}
	
	// Verify severity levels
	if !strings.Contains(prompt, "CRITICAL") {
		t.Error("Prompt should define CRITICAL severity")
	}
}

// TestTestingAgentPromptBuild tests testing agent prompt construction
func TestTestingAgentPromptBuild(t *testing.T) {
	agent := &TestingAgent{}
	
	prompt := agent.buildPrompt("code", "build an api", "const express = require('express')")
	
	// Verify framework references
	if !strings.Contains(prompt, "Jest") {
		t.Error("Prompt should mention Jest")
	}
	if !strings.Contains(prompt, "Vitest") {
		t.Error("Prompt should mention Vitest")
	}
	if !strings.Contains(prompt, "pytest") {
		t.Error("Prompt should mention pytest")
	}
	
	// Verify coverage targets
	if !strings.Contains(prompt, "70%") {
		t.Error("Prompt should mention 70% MVP coverage target")
	}
}

// TestDocumentationAgentPromptBuild tests documentation agent prompt construction
func TestDocumentationAgentPromptBuild(t *testing.T) {
	agent := &DocumentationAgent{}
	
	prompt := agent.buildPrompt("code", "build a cli tool", "package main")
	
	// Verify README structure
	if !strings.Contains(prompt, "shields.io") {
		t.Error("Prompt should mention shields.io badges")
	}
	if !strings.Contains(prompt, "Quick Start") {
		t.Error("Prompt should mention Quick Start section")
	}
	if !strings.Contains(prompt, "Troubleshooting") {
		t.Error("Prompt should mention Troubleshooting section")
	}
}

// TestComplexityScorer tests the complexity scorer
func TestComplexityScorer(t *testing.T) {
	cfg := &SupervisorConfig{
		ComplexityThreshold: 5,
	}
	scorer := NewComplexityScorer(cfg)
	
	tests := []struct {
		name     string
		taskType string
		input    string
		minScore int
		maxScore int
	}{
		{
			name:     "Simple task",
			taskType: "code",
			input:    "hello world",
			minScore: 1,
			maxScore: 3,
		},
		{
			name:     "Complex task with keywords",
			taskType: "code",
			input:    "Build a full-stack application with authentication, database, real-time updates, and deployment",
			minScore: 5,
			maxScore: 10,
		},
		{
			name:     "Non-code task",
			taskType: "review",
			input:    "Review this code",
			minScore: 1,
			maxScore: 1,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := scorer.Score(tt.taskType, tt.input)
			if analysis.Score < tt.minScore || analysis.Score > tt.maxScore {
				t.Errorf("Score() = %v, want between %v and %v", analysis.Score, tt.minScore, tt.maxScore)
			}
		})
	}
}
