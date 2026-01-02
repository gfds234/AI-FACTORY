package project

import (
	"fmt"
	"os"
	"strings"
)

// CompletionValidator validates hand-off ready criteria
type CompletionValidator struct {
	artifactsDir string
}

// NewCompletionValidator creates a new completion validator
func NewCompletionValidator(artifactsDir string) *CompletionValidator {
	return &CompletionValidator{
		artifactsDir: artifactsDir,
	}
}

// ValidateHandoffReady validates if a project meets hand-off ready criteria
func (cv *CompletionValidator) ValidateHandoffReady(project *Project) (*CompletionMetrics, error) {
	metrics := &CompletionMetrics{
		BlockingIssues: []string{},
	}

	// Check criteria
	hasCode, err := cv.checkRunnableBuild(project)
	if err != nil {
		metrics.BlockingIssues = append(metrics.BlockingIssues, fmt.Sprintf("Build check error: %v", err))
	}
	metrics.HasRunnableBuild = hasCode

	hasTests, err := cv.checkTestsExist(project)
	if err != nil {
		metrics.BlockingIssues = append(metrics.BlockingIssues, fmt.Sprintf("Test check error: %v", err))
	}
	metrics.HasTests = hasTests

	hasReadme, err := cv.checkReadmeExists(project)
	if err != nil {
		metrics.BlockingIssues = append(metrics.BlockingIssues, fmt.Sprintf("README check error: %v", err))
	}
	metrics.HasReadme = hasReadme

	// Calculate completion percentage
	metrics.CompletionPct = cv.calculateCompletionPercentage(project, metrics)

	// Add blocking issues if criteria not met
	if !metrics.HasRunnableBuild {
		metrics.BlockingIssues = append(metrics.BlockingIssues, "No runnable build detected")
	}
	if !metrics.HasTests {
		metrics.BlockingIssues = append(metrics.BlockingIssues, "No tests detected")
	}
	if !metrics.HasReadme {
		metrics.BlockingIssues = append(metrics.BlockingIssues, "No README documentation detected")
	}

	return metrics, nil
}

// checkRunnableBuild checks if runnable code exists
func (cv *CompletionValidator) checkRunnableBuild(project *Project) (bool, error) {
	// Entry point markers for different languages
	entryPointMarkers := []string{
		// Go
		"func main()",
		"package main",

		// Python
		"if __name__ == \"__main__\":",
		"if __name__ == '__main__':",

		// JavaScript/Node
		"function main()",
		"export default",
		"module.exports",

		// C#
		"static void Main(",
		"static async Task Main(",

		// Java
		"public static void main(",

		// Rust
		"fn main()",

		// C/C++
		"int main(",
	}

	// Scan all artifacts for entry point markers
	for _, artifactPath := range project.ArtifactPaths {
		fullPath := cv.getFullArtifactPath(artifactPath)

		content, err := os.ReadFile(fullPath)
		if err != nil {
			continue // Skip unreadable artifacts
		}

		contentStr := string(content)

		for _, marker := range entryPointMarkers {
			if strings.Contains(contentStr, marker) {
				return true, nil
			}
		}
	}

	return false, nil
}

// checkTestsExist checks if tests exist
func (cv *CompletionValidator) checkTestsExist(project *Project) (bool, error) {
	// Test markers for different languages/frameworks
	testMarkers := []string{
		// Go
		"func Test",
		"func Benchmark",
		"testing.T",

		// Python
		"def test_",
		"class Test",
		"import unittest",
		"import pytest",

		// JavaScript/TypeScript
		"it(\"",
		"it('",
		"describe(\"",
		"describe('",
		"test(\"",
		"test('",

		// C#
		"[Test]",
		"[Fact]",
		"[Theory]",

		// Java
		"@Test",
		"@TestCase",

		// Rust
		"#[test]",
	}

	// Scan all artifacts for test markers
	for _, artifactPath := range project.ArtifactPaths {
		fullPath := cv.getFullArtifactPath(artifactPath)

		content, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}

		contentStr := string(content)

		for _, marker := range testMarkers {
			if strings.Contains(contentStr, marker) {
				return true, nil
			}
		}
	}

	return false, nil
}

// checkReadmeExists checks if README documentation exists
func (cv *CompletionValidator) checkReadmeExists(project *Project) (bool, error) {
	// README markers
	readmeMarkers := []string{
		"# Setup",
		"## Setup",
		"# Installation",
		"## Installation",
		"# Usage",
		"## Usage",
		"# Getting Started",
		"## Getting Started",
		"### Prerequisites",
		"How to run",
		"How to use",
		"Quick Start",
	}

	markerCount := 0
	requiredMarkers := 2 // At least 2 README markers needed

	// Scan all artifacts for README markers
	for _, artifactPath := range project.ArtifactPaths {
		fullPath := cv.getFullArtifactPath(artifactPath)

		content, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}

		contentStr := string(content)

		for _, marker := range readmeMarkers {
			if strings.Contains(contentStr, marker) {
				markerCount++
				if markerCount >= requiredMarkers {
					return true, nil
				}
			}
		}
	}

	return markerCount >= requiredMarkers, nil
}

// calculateCompletionPercentage calculates overall project completion percentage
func (cv *CompletionValidator) calculateCompletionPercentage(project *Project, metrics *CompletionMetrics) float64 {
	// Phase weights (total 80%)
	var phaseTotal float64
	for _, phaseExec := range project.Phases {
		if phaseExec.Status == PhaseStatusComplete {
			weight, exists := PhaseWeights[phaseExec.Phase]
			if exists {
				phaseTotal += weight
			}
		}
	}

	// Criteria bonus (max 20%)
	criteriaBonus := 0.0
	if metrics.HasRunnableBuild {
		criteriaBonus += 7.0
	}
	if metrics.HasTests {
		criteriaBonus += 7.0
	}
	if metrics.HasReadme {
		criteriaBonus += 6.0
	}

	total := phaseTotal + criteriaBonus

	// Cap at 100%
	if total > 100.0 {
		total = 100.0
	}

	return total
}

// getFullArtifactPath returns the full path to an artifact
func (cv *CompletionValidator) getFullArtifactPath(artifactPath string) string {
	// If already absolute path, return as-is
	if strings.HasPrefix(artifactPath, cv.artifactsDir) {
		return artifactPath
	}

	// Otherwise, join with artifacts directory
	return artifactPath
}
