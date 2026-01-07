package project

import (
	"ai-studio/orchestrator/supervisor"
	"fmt"
	"os"
	"path/filepath"
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

	// Run actual build verification if project directory exists
	projectDir := cv.extractProjectDirectory(project)
	if projectDir != "" {
		verifyAgent := supervisor.NewVerificationAgent()
		verifyResult, err := verifyAgent.VerifyProject(projectDir)
		if err == nil && verifyResult != nil {
			// Store verification results
			metrics.SyntaxValid = verifyResult.SyntaxValid
			metrics.DependenciesOK = verifyResult.DependenciesOK
			metrics.BuildLog = verifyResult.BuildLog
			metrics.HasRunnableBuild = verifyResult.EntryPointValid && verifyResult.SyntaxValid

			// Add errors as blocking issues
			for _, errMsg := range verifyResult.Errors {
				metrics.BlockingIssues = append(metrics.BlockingIssues, errMsg)
			}
		} else {
			// Fallback to marker-based detection if verification fails
			hasCode, err := cv.checkRunnableBuildMarkers(project)
			if err != nil {
				metrics.BlockingIssues = append(metrics.BlockingIssues, fmt.Sprintf("Build check error: %v", err))
			}
			metrics.HasRunnableBuild = hasCode
		}
	} else {
		// No project directory - use marker-based detection
		hasCode, err := cv.checkRunnableBuildMarkers(project)
		if err != nil {
			metrics.BlockingIssues = append(metrics.BlockingIssues, fmt.Sprintf("Build check error: %v", err))
		}
		metrics.HasRunnableBuild = hasCode
	}

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

	// Populate runtime and test metrics from stored validation results
	if project.ValidationResults != nil {
		vr := project.ValidationResults

		metrics.RuntimeVerified = vr.RuntimeVerified
		metrics.TestsExecuted = vr.TestsExecuted
		metrics.TestsPassed = vr.TestsPassed
		metrics.TestsFailed = vr.TestsFailed

		// Calculate quality score (0-100)
		metrics.QualityScore = cv.calculateQualityScore(metrics)
	}

	return metrics, nil
}

// extractProjectDirectory extracts project directory from artifact paths
func (cv *CompletionValidator) extractProjectDirectory(project *Project) string {
	// Look for project directory in artifact paths (e.g., "projects/generated_1234567890")
	for _, artifactPath := range project.ArtifactPaths {
		if strings.Contains(artifactPath, "projects/generated_") || strings.Contains(artifactPath, "projects\\generated_") {
			// Extract directory from path like "artifacts\code_123.md (project: projects/generated_456)"
			if strings.Contains(artifactPath, "(project: ") {
				start := strings.Index(artifactPath, "(project: ") + len("(project: ")
				end := strings.Index(artifactPath[start:], ")")
				if end > 0 {
					return artifactPath[start : start+end]
				}
			}
			// Direct project path
			if strings.HasPrefix(artifactPath, "projects/") || strings.HasPrefix(artifactPath, "projects\\") {
				// Return the base project directory
				parts := strings.Split(filepath.ToSlash(artifactPath), "/")
				if len(parts) >= 2 {
					return filepath.Join(parts[0], parts[1])
				}
			}
		}
	}
	return ""
}

// checkRunnableBuildMarkers checks if runnable code exists (marker-based fallback)
func (cv *CompletionValidator) checkRunnableBuildMarkers(project *Project) (bool, error) {
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

// calculateQualityScore calculates 0-100 quality score based on Triple Guarantee
func (cv *CompletionValidator) calculateQualityScore(metrics *CompletionMetrics) int {
	score := 0

	// Build verification (35 points)
	if metrics.SyntaxValid {
		score += 15
	}
	if metrics.DependenciesOK {
		score += 10
	}
	if metrics.HasRunnableBuild {
		score += 10
	}

	// Runtime verification (25 points)
	if metrics.RuntimeVerified {
		score += 25
	}

	// Test verification (20 points)
	if metrics.TestsExecuted && (metrics.TestsPassed+metrics.TestsFailed) > 0 {
		score += 5
		passRate := float64(metrics.TestsPassed) / float64(metrics.TestsPassed+metrics.TestsFailed)
		score += int(passRate * 15)
	}

	// Documentation (20 points)
	if metrics.HasReadme {
		score += 20
	}

	return score
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
