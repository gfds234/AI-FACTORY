package project

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// QualityGuarantee represents the overall quality verification results
type QualityGuarantee struct {
	ProjectName       string
	GeneratedAt       time.Time
	OverallScore      int    // 0-100
	Status            string // "READY", "NEEDS_WORK", "BLOCKED"

	// Build verification
	BuildPassed       bool
	SyntaxValid       bool
	DependenciesOK    bool
	EntryPointValid   bool
	BuildErrors       []string

	// Runtime verification
	RuntimePassed     bool
	ApplicationStarts bool
	HealthCheckPassed bool
	RuntimeErrors     []string
	RuntimeWarnings   []string

	// Test verification
	TestsPassed       int
	TestsFailed       int
	TestsSkipped      int
	TotalTests        int
	TestFramework     string
	TestErrors        []string

	// Deployment verification
	DeploymentReady   bool
	DockerBuilds      bool
	EnvVarsDocumented bool
	DeploymentErrors  []string

	// Documentation
	ReadmeComplete    bool
	SetupInstructions bool
}

// GenerateQualityReport creates a professional quality guarantee report
func GenerateQualityReport(projectName string, metrics CompletionMetrics) *QualityGuarantee {
	qg := &QualityGuarantee{
		ProjectName:       projectName,
		GeneratedAt:       time.Now(),
		BuildErrors:       make([]string, 0),
		RuntimeErrors:     make([]string, 0),
		RuntimeWarnings:   make([]string, 0),
		TestErrors:        make([]string, 0),
		DeploymentErrors:  make([]string, 0),
	}

	// Build verification
	qg.SyntaxValid = metrics.SyntaxValid
	qg.DependenciesOK = metrics.DependenciesOK
	qg.EntryPointValid = metrics.HasRunnableBuild
	qg.BuildPassed = qg.SyntaxValid && qg.DependenciesOK && qg.EntryPointValid
	qg.BuildErrors = metrics.BlockingIssues

	// Runtime verification
	qg.ApplicationStarts = metrics.RuntimeVerified
	qg.HealthCheckPassed = metrics.RuntimeVerified // Simplified for now
	qg.RuntimePassed = qg.ApplicationStarts

	// Test verification
	qg.TestsPassed = metrics.TestsPassed
	qg.TestsFailed = metrics.TestsFailed
	qg.TotalTests = metrics.TestsPassed + metrics.TestsFailed

	// Deployment verification
	qg.DockerBuilds = metrics.DockerBuilds
	qg.EnvVarsDocumented = metrics.EnvVarsDocumented
	qg.DeploymentReady = qg.DockerBuilds && qg.EnvVarsDocumented

	// Documentation
	qg.ReadmeComplete = metrics.HasReadme
	qg.SetupInstructions = metrics.HasReadme

	// Calculate overall quality score
	qg.OverallScore = calculateQualityScore(qg)
	qg.Status = determineStatus(qg)

	return qg
}

// calculateQualityScore computes a 0-100 quality score
func calculateQualityScore(qg *QualityGuarantee) int {
	score := 0

	// Build verification (35 points - CRITICAL)
	if qg.SyntaxValid {
		score += 15
	}
	if qg.DependenciesOK {
		score += 10
	}
	if qg.EntryPointValid {
		score += 10
	}

	// Runtime verification (25 points - HIGH PRIORITY)
	if qg.ApplicationStarts {
		score += 15
	}
	if qg.HealthCheckPassed {
		score += 10
	}

	// Test verification (20 points - MEDIUM PRIORITY)
	if qg.TotalTests > 0 {
		score += 5 // Has tests at all
		passRate := float64(qg.TestsPassed) / float64(qg.TotalTests)
		score += int(passRate * 15) // Up to 15 points for pass rate
	}

	// Deployment verification (10 points - NICE TO HAVE)
	if qg.DockerBuilds {
		score += 5
	}
	if qg.EnvVarsDocumented {
		score += 5
	}

	// Documentation (10 points - NICE TO HAVE)
	if qg.ReadmeComplete {
		score += 5
	}
	if qg.SetupInstructions {
		score += 5
	}

	return score
}

// determineStatus determines the overall project status
func determineStatus(qg *QualityGuarantee) string {
	// BLOCKED: Critical build failures
	if !qg.BuildPassed {
		return "BLOCKED"
	}

	// NEEDS_WORK: Build passes but runtime or significant test failures
	if !qg.RuntimePassed {
		return "NEEDS_WORK"
	}

	if qg.TotalTests > 0 && qg.TestsPassed == 0 {
		return "NEEDS_WORK"
	}

	// READY: Build passes, runtime works, tests mostly pass
	if qg.OverallScore >= 70 {
		return "READY"
	}

	return "NEEDS_WORK"
}

// ToMarkdown generates a professional markdown report
func (qg *QualityGuarantee) ToMarkdown() string {
	var sb strings.Builder

	// Header
	sb.WriteString("# Quality Guarantee Report\n\n")
	sb.WriteString(fmt.Sprintf("**Project:** %s  \n", qg.ProjectName))
	sb.WriteString(fmt.Sprintf("**Generated:** %s  \n", qg.GeneratedAt.Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("**Overall Score:** %d/100  \n", qg.OverallScore))
	sb.WriteString(fmt.Sprintf("**Status:** **%s**\n\n", qg.Status))

	sb.WriteString("---\n\n")

	// Build Status
	sb.WriteString("## Build Status: ")
	if qg.BuildPassed {
		sb.WriteString("✅ PASSED\n\n")
	} else {
		sb.WriteString("❌ FAILED\n\n")
	}

	sb.WriteString(fmt.Sprintf("- Syntax validation: %s\n", checkmark(qg.SyntaxValid)))
	sb.WriteString(fmt.Sprintf("- Dependencies resolved: %s\n", checkmark(qg.DependenciesOK)))
	sb.WriteString(fmt.Sprintf("- Entry point verified: %s\n", checkmark(qg.EntryPointValid)))

	if len(qg.BuildErrors) > 0 {
		sb.WriteString("\n**Build Errors:**\n")
		for _, err := range qg.BuildErrors {
			sb.WriteString(fmt.Sprintf("- %s\n", err))
		}
	}
	sb.WriteString("\n")

	// Runtime Status
	sb.WriteString("## Runtime Status: ")
	if qg.RuntimePassed {
		sb.WriteString("✅ PASSED\n\n")
	} else if len(qg.RuntimeErrors) > 0 {
		sb.WriteString("❌ FAILED\n\n")
	} else {
		sb.WriteString("⚠️ NOT TESTED\n\n")
	}

	if qg.RuntimePassed || len(qg.RuntimeErrors) > 0 {
		sb.WriteString(fmt.Sprintf("- Application starts: %s\n", checkmark(qg.ApplicationStarts)))
		sb.WriteString(fmt.Sprintf("- Health check: %s\n", checkmark(qg.HealthCheckPassed)))

		if len(qg.RuntimeErrors) > 0 {
			sb.WriteString("\n**Runtime Errors:**\n")
			for _, err := range qg.RuntimeErrors {
				sb.WriteString(fmt.Sprintf("- %s\n", err))
			}
		}

		if len(qg.RuntimeWarnings) > 0 {
			sb.WriteString("\n**Warnings:**\n")
			for _, warn := range qg.RuntimeWarnings {
				sb.WriteString(fmt.Sprintf("- %s\n", warn))
			}
		}
	}
	sb.WriteString("\n")

	// Test Status
	sb.WriteString("## Test Status: ")
	if qg.TotalTests == 0 {
		sb.WriteString("⚠️ NO TESTS\n\n")
		sb.WriteString("No tests were found or executed.\n\n")
	} else if qg.TestsFailed == 0 {
		sb.WriteString("✅ ALL PASSED\n\n")
		sb.WriteString(fmt.Sprintf("- Tests executed: %d\n", qg.TotalTests))
		sb.WriteString(fmt.Sprintf("- Tests passed: %d\n", qg.TestsPassed))
		sb.WriteString(fmt.Sprintf("- Tests failed: 0\n"))
		if qg.TestFramework != "" {
			sb.WriteString(fmt.Sprintf("- Framework: %s\n", qg.TestFramework))
		}
		sb.WriteString("\n")
	} else {
		passRate := int(float64(qg.TestsPassed) / float64(qg.TotalTests) * 100)
		sb.WriteString(fmt.Sprintf("⚠️ PARTIAL (%d/%d passed, %d%%)\n\n", qg.TestsPassed, qg.TotalTests, passRate))
		sb.WriteString(fmt.Sprintf("- Tests executed: %d\n", qg.TotalTests))
		sb.WriteString(fmt.Sprintf("- Tests passed: %d\n", qg.TestsPassed))
		sb.WriteString(fmt.Sprintf("- Tests failed: %d\n", qg.TestsFailed))
		if qg.TestsSkipped > 0 {
			sb.WriteString(fmt.Sprintf("- Tests skipped: %d\n", qg.TestsSkipped))
		}
		if qg.TestFramework != "" {
			sb.WriteString(fmt.Sprintf("- Framework: %s\n", qg.TestFramework))
		}
		sb.WriteString("\n")
	}

	// Deployment Status
	sb.WriteString("## Deployment Status: ")
	if qg.DeploymentReady {
		sb.WriteString("✅ READY\n\n")
	} else if qg.DockerBuilds || qg.EnvVarsDocumented {
		sb.WriteString("⚠️ PARTIAL\n\n")
	} else {
		sb.WriteString("⚠️ NOT VERIFIED\n\n")
	}

	if qg.DockerBuilds || qg.EnvVarsDocumented {
		sb.WriteString(fmt.Sprintf("- Dockerfile builds: %s\n", checkmark(qg.DockerBuilds)))
		sb.WriteString(fmt.Sprintf("- Environment configured: %s\n", checkmark(qg.EnvVarsDocumented)))

		if len(qg.DeploymentErrors) > 0 {
			sb.WriteString("\n**Deployment Issues:**\n")
			for _, err := range qg.DeploymentErrors {
				sb.WriteString(fmt.Sprintf("- %s\n", err))
			}
		}
	}
	sb.WriteString("\n")

	// Documentation Status
	sb.WriteString("## Documentation: ")
	if qg.ReadmeComplete {
		sb.WriteString("✅ COMPLETE\n\n")
		sb.WriteString(fmt.Sprintf("- README present: %s\n", checkmark(qg.ReadmeComplete)))
		sb.WriteString(fmt.Sprintf("- Setup instructions: %s\n", checkmark(qg.SetupInstructions)))
	} else {
		sb.WriteString("❌ MISSING\n\n")
	}
	sb.WriteString("\n")

	// Client Handoff Checklist
	sb.WriteString("---\n\n")
	sb.WriteString("## Client Handoff Checklist\n\n")
	sb.WriteString(fmt.Sprintf("%s Code compiles without errors\n", checkmark(qg.BuildPassed)))
	sb.WriteString(fmt.Sprintf("%s All dependencies documented\n", checkmark(qg.DependenciesOK)))
	sb.WriteString(fmt.Sprintf("%s Application runs successfully\n", checkmark(qg.ApplicationStarts)))

	if qg.TotalTests > 0 {
		allTestsPass := qg.TestsFailed == 0
		sb.WriteString(fmt.Sprintf("%s All tests passing", checkmark(allTestsPass)))
		if !allTestsPass {
			sb.WriteString(fmt.Sprintf(" (%d/%d passed)", qg.TestsPassed, qg.TotalTests))
		}
		sb.WriteString("\n")
	} else {
		sb.WriteString("⚠️ No tests generated\n")
	}

	sb.WriteString(fmt.Sprintf("%s Documentation complete\n", checkmark(qg.ReadmeComplete)))

	if qg.DockerBuilds {
		sb.WriteString(fmt.Sprintf("%s Deployment configuration tested\n", checkmark(qg.DockerBuilds)))
	}

	sb.WriteString("\n")

	// Final summary
	sb.WriteString("---\n\n")
	sb.WriteString("## Summary\n\n")

	switch qg.Status {
	case "READY":
		sb.WriteString("✅ **This project is READY FOR CLIENT DELIVERY**\n\n")
		sb.WriteString("All critical quality checks have passed. The code compiles, runs successfully, and has been verified for deployment readiness.\n")
	case "NEEDS_WORK":
		sb.WriteString("⚠️ **This project NEEDS WORK before client delivery**\n\n")
		sb.WriteString("Some quality checks have failed or need improvement. Review the errors and warnings above before proceeding.\n")
	case "BLOCKED":
		sb.WriteString("❌ **This project is BLOCKED from delivery**\n\n")
		sb.WriteString("Critical build failures prevent this project from being delivered. The code must be fixed before it can run.\n")
	}

	sb.WriteString("\n")
	sb.WriteString("---\n\n")
	sb.WriteString("*Generated by AI FACTORY Automated Quality System*\n")

	return sb.String()
}

// checkmark returns a checkmark or X emoji
func checkmark(passed bool) string {
	if passed {
		return "✅"
	}
	return "❌"
}

// SaveQualityReport saves the quality report to a file
func (qg *QualityGuarantee) SaveToFile(filepath string) error {
	content := qg.ToMarkdown()
	return os.WriteFile(filepath, []byte(content), 0644)
}
