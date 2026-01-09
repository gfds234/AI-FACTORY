package validation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TestExecutor executes tests and parses results
type TestExecutor struct {
	timeout time.Duration
}

// NewTestExecutor creates a new test executor
func NewTestExecutor() *TestExecutor {
	return &TestExecutor{
		timeout: 30 * time.Second, // 30 seconds for test execution
	}
}

// TestResult contains test execution results
type TestResult struct {
	TestsExecuted bool     `json:"tests_executed"` // true if tests were found and run
	TestsPassed   int      `json:"tests_passed"`   // number of passing tests
	TestsFailed   int      `json:"tests_failed"`   // number of failing tests
	TestsSkipped  int      `json:"tests_skipped"`  // number of skipped tests
	TotalTests    int      `json:"total_tests"`    // total number of tests
	TestFramework string   `json:"test_framework"` // detected test framework
	TestOutput    string   `json:"test_output"`    // full test output
	Errors        []string `json:"errors"`         // execution errors
	Warnings      []string `json:"warnings"`       // non-fatal issues
}

// ExecuteTests runs tests and returns parsed results
func (te *TestExecutor) ExecuteTests(projectPath string, projectType string) (*TestResult, error) {
	result := &TestResult{
		Errors:   make([]string, 0),
		Warnings: make([]string, 0),
	}

	switch projectType {
	case "nodejs":
		return te.executeNodeJSTests(projectPath)
	case "python":
		return te.executePythonTests(projectPath)
	case "go":
		return te.executeGoTests(projectPath)
	case "frontend":
		result.Warnings = append(result.Warnings, "Frontend projects typically don't have backend tests")
		return result, nil
	default:
		result.Warnings = append(result.Warnings, fmt.Sprintf("Unknown project type %s - skipping test execution", projectType))
		return result, nil
	}
}

// executeNodeJSTests runs Node.js tests
func (te *TestExecutor) executeNodeJSTests(projectPath string) (*TestResult, error) {
	result := &TestResult{
		Errors:   make([]string, 0),
		Warnings: make([]string, 0),
	}

	// Check for package.json
	packagePath := filepath.Join(projectPath, "package.json")
	if _, err := os.Stat(packagePath); os.IsNotExist(err) {
		// Try backend subdirectory
		packagePath = filepath.Join(projectPath, "backend", "package.json")
		if _, err := os.Stat(packagePath); os.IsNotExist(err) {
			result.Errors = append(result.Errors, "No package.json found")
			return result, nil
		}
		projectPath = filepath.Join(projectPath, "backend")
	}

	// Read package.json to detect test framework
	packageData, err := os.ReadFile(packagePath)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to read package.json: %v", err))
		return result, nil
	}

	var packageJSON map[string]interface{}
	if err := json.Unmarshal(packageData, &packageJSON); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to parse package.json: %v", err))
		return result, nil
	}

	// Check if test script exists
	scripts, ok := packageJSON["scripts"].(map[string]interface{})
	if !ok || scripts["test"] == nil {
		result.Warnings = append(result.Warnings, "No test script found in package.json")
		return result, nil
	}

	// Detect test framework from dependencies
	deps := make(map[string]interface{})
	if devDeps, ok := packageJSON["devDependencies"].(map[string]interface{}); ok {
		for k, v := range devDeps {
			deps[k] = v
		}
	}
	if regularDeps, ok := packageJSON["dependencies"].(map[string]interface{}); ok {
		for k, v := range regularDeps {
			deps[k] = v
		}
	}

	framework := "unknown"
	if _, hasJest := deps["jest"]; hasJest {
		framework = "jest"
	} else if _, hasMocha := deps["mocha"]; hasMocha {
		framework = "mocha"
	} else if _, hasVitest := deps["vitest"]; hasVitest {
		framework = "vitest"
	}
	result.TestFramework = framework

	// Run npm test
	ctx, cancel := context.WithTimeout(context.Background(), te.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "npm", "test")
	cmd.Dir = projectPath

	var outputBuf bytes.Buffer
	cmd.Stdout = &outputBuf
	cmd.Stderr = &outputBuf

	err = cmd.Run()
	output := outputBuf.String()
	result.TestOutput = output

	// Parse test output based on framework
	if framework == "jest" {
		te.parseJestOutput(output, result)
	} else if framework == "mocha" {
		te.parseMochaOutput(output, result)
	} else if framework == "vitest" {
		te.parseVitestOutput(output, result)
	} else {
		// Try to parse generic output
		te.parseGenericNodeOutput(output, result)
	}

	// If no tests were detected but command ran, mark as executed
	if result.TotalTests > 0 {
		result.TestsExecuted = true
	} else if err == nil {
		result.Warnings = append(result.Warnings, "Tests ran but no test results detected in output")
	}

	// If command failed due to test failures, that's not an error (tests just failed)
	if err != nil && ctx.Err() == context.DeadlineExceeded {
		result.Errors = append(result.Errors, "Test execution timed out after 30 seconds")
	} else if err != nil && result.TotalTests == 0 {
		result.Errors = append(result.Errors, fmt.Sprintf("Test execution failed: %v", err))
	}

	return result, nil
}

// executePythonTests runs Python tests
func (te *TestExecutor) executePythonTests(projectPath string) (*TestResult, error) {
	result := &TestResult{
		Errors:   make([]string, 0),
		Warnings: make([]string, 0),
	}

	// Check for pytest
	ctx, cancel := context.WithTimeout(context.Background(), te.timeout)
	defer cancel()

	// Try pytest first
	cmd := exec.CommandContext(ctx, "pytest", "--tb=short", "-v")
	cmd.Dir = projectPath

	var outputBuf bytes.Buffer
	cmd.Stdout = &outputBuf
	cmd.Stderr = &outputBuf

	err := cmd.Run()
	output := outputBuf.String()
	result.TestOutput = output

	if err != nil && strings.Contains(output, "no tests ran") {
		// Try unittest instead
		cmd = exec.CommandContext(ctx, "python", "-m", "unittest", "discover", "-v")
		cmd.Dir = projectPath

		outputBuf.Reset()
		cmd.Stdout = &outputBuf
		cmd.Stderr = &outputBuf

		err = cmd.Run()
		output = outputBuf.String()
		result.TestOutput = output
		result.TestFramework = "unittest"
		te.parseUnittestOutput(output, result)
	} else {
		result.TestFramework = "pytest"
		te.parsePytestOutput(output, result)
	}

	// If no tests were detected but command ran, mark as executed
	if result.TotalTests > 0 {
		result.TestsExecuted = true
	} else if err == nil {
		result.Warnings = append(result.Warnings, "Tests ran but no test results detected")
	}

	if err != nil && ctx.Err() == context.DeadlineExceeded {
		result.Errors = append(result.Errors, "Test execution timed out after 30 seconds")
	} else if err != nil && result.TotalTests == 0 {
		result.Errors = append(result.Errors, fmt.Sprintf("Test execution failed: %v", err))
	}

	return result, nil
}

// executeGoTests runs Go tests
func (te *TestExecutor) executeGoTests(projectPath string) (*TestResult, error) {
	result := &TestResult{
		Errors:       make([]string, 0),
		Warnings:     make([]string, 0),
		TestFramework: "go test",
	}

	// Check for go.mod
	var workDir string
	if _, err := os.Stat(filepath.Join(projectPath, "go.mod")); err == nil {
		workDir = projectPath
	} else if _, err := os.Stat(filepath.Join(projectPath, "backend", "go.mod")); err == nil {
		workDir = filepath.Join(projectPath, "backend")
	} else {
		result.Errors = append(result.Errors, "No go.mod found")
		return result, nil
	}

	// Run go test
	ctx, cancel := context.WithTimeout(context.Background(), te.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "test", "-v", "./...")
	cmd.Dir = workDir

	var outputBuf bytes.Buffer
	cmd.Stdout = &outputBuf
	cmd.Stderr = &outputBuf

	err := cmd.Run()
	output := outputBuf.String()
	result.TestOutput = output

	// Parse go test output
	te.parseGoTestOutput(output, result)

	// If no tests were detected but command ran, mark as executed
	if result.TotalTests > 0 {
		result.TestsExecuted = true
	} else if err == nil {
		result.Warnings = append(result.Warnings, "Tests ran but no test results detected")
	}

	if err != nil && ctx.Err() == context.DeadlineExceeded {
		result.Errors = append(result.Errors, "Test execution timed out after 30 seconds")
	} else if err != nil && result.TotalTests == 0 {
		result.Errors = append(result.Errors, fmt.Sprintf("Test execution failed: %v", err))
	}

	return result, nil
}

// parseJestOutput parses Jest test output
func (te *TestExecutor) parseJestOutput(output string, result *TestResult) {
	// Jest output format: "Tests: 2 failed, 3 passed, 5 total"
	re := regexp.MustCompile(`Tests:\s+(?:(\d+)\s+failed,\s*)?(?:(\d+)\s+passed,\s*)?(\d+)\s+total`)
	matches := re.FindStringSubmatch(output)

	if len(matches) > 0 {
		if matches[1] != "" {
			result.TestsFailed, _ = strconv.Atoi(matches[1])
		}
		if matches[2] != "" {
			result.TestsPassed, _ = strconv.Atoi(matches[2])
		}
		if matches[3] != "" {
			result.TotalTests, _ = strconv.Atoi(matches[3])
		}
	}
}

// parseMochaOutput parses Mocha test output
func (te *TestExecutor) parseMochaOutput(output string, result *TestResult) {
	// Mocha output format: "5 passing" and "2 failing"
	passingRe := regexp.MustCompile(`(\d+)\s+passing`)
	failingRe := regexp.MustCompile(`(\d+)\s+failing`)

	if matches := passingRe.FindStringSubmatch(output); len(matches) > 1 {
		result.TestsPassed, _ = strconv.Atoi(matches[1])
	}
	if matches := failingRe.FindStringSubmatch(output); len(matches) > 1 {
		result.TestsFailed, _ = strconv.Atoi(matches[1])
	}

	result.TotalTests = result.TestsPassed + result.TestsFailed
}

// parseVitestOutput parses Vitest test output
func (te *TestExecutor) parseVitestOutput(output string, result *TestResult) {
	// Vitest output formats:
	// "Test Files  1 passed (1)"
	// "Tests  5 passed (5)"
	// "Tests  3 passed | 2 failed (5)"

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Match lines starting with "Tests " or "Test Files "
		if strings.HasPrefix(line, "Tests ") || strings.HasPrefix(line, "Test Files ") {
			// Extract numbers using regex
			passedRe := regexp.MustCompile(`(\d+)\s+passed`)
			failedRe := regexp.MustCompile(`(\d+)\s+failed`)
			skippedRe := regexp.MustCompile(`(\d+)\s+skipped`)
			totalRe := regexp.MustCompile(`\((\d+)\)`)

			if match := passedRe.FindStringSubmatch(line); len(match) > 1 {
				if count, err := strconv.Atoi(match[1]); err == nil {
					result.TestsPassed = count
				}
			}

			if match := failedRe.FindStringSubmatch(line); len(match) > 1 {
				if count, err := strconv.Atoi(match[1]); err == nil {
					result.TestsFailed = count
				}
			}

			if match := skippedRe.FindStringSubmatch(line); len(match) > 1 {
				if count, err := strconv.Atoi(match[1]); err == nil {
					result.TestsSkipped = count
				}
			}

			if match := totalRe.FindStringSubmatch(line); len(match) > 1 {
				if count, err := strconv.Atoi(match[1]); err == nil {
					result.TotalTests = count
				}
			}
		}
	}

	// Fallback: if total not detected, calculate from passed + failed + skipped
	if result.TotalTests == 0 && (result.TestsPassed > 0 || result.TestsFailed > 0 || result.TestsSkipped > 0) {
		result.TotalTests = result.TestsPassed + result.TestsFailed + result.TestsSkipped
	}
}

// parseGenericNodeOutput tries to parse generic test output
func (te *TestExecutor) parseGenericNodeOutput(output string, result *TestResult) {
	// Try to count "✓" or "✔" for passed tests
	passCount := strings.Count(output, "✓") + strings.Count(output, "✔")
	// Try to count "✗" or "×" for failed tests
	failCount := strings.Count(output, "✗") + strings.Count(output, "×")

	if passCount > 0 || failCount > 0 {
		result.TestsPassed = passCount
		result.TestsFailed = failCount
		result.TotalTests = passCount + failCount
	}
}

// parsePytestOutput parses pytest output
func (te *TestExecutor) parsePytestOutput(output string, result *TestResult) {
	// Pytest output format: "5 passed, 2 failed in 1.23s"
	re := regexp.MustCompile(`(?:(\d+)\s+failed)?(?:,\s*)?(?:(\d+)\s+passed)?(?:,\s*)?(?:(\d+)\s+skipped)?`)
	matches := re.FindStringSubmatch(output)

	if len(matches) > 0 {
		if matches[1] != "" {
			result.TestsFailed, _ = strconv.Atoi(matches[1])
		}
		if matches[2] != "" {
			result.TestsPassed, _ = strconv.Atoi(matches[2])
		}
		if matches[3] != "" {
			result.TestsSkipped, _ = strconv.Atoi(matches[3])
		}
		result.TotalTests = result.TestsFailed + result.TestsPassed + result.TestsSkipped
	}
}

// parseUnittestOutput parses unittest output
func (te *TestExecutor) parseUnittestOutput(output string, result *TestResult) {
	// Unittest output format: "Ran 5 tests in 0.123s" followed by "OK" or "FAILED"
	ranRe := regexp.MustCompile(`Ran\s+(\d+)\s+test`)
	if matches := ranRe.FindStringSubmatch(output); len(matches) > 1 {
		result.TotalTests, _ = strconv.Atoi(matches[1])
	}

	// If "OK" appears, all tests passed
	if strings.Contains(output, "OK") {
		result.TestsPassed = result.TotalTests
		result.TestsFailed = 0
	} else if strings.Contains(output, "FAILED") {
		// Try to find failure count
		failRe := regexp.MustCompile(`failures=(\d+)`)
		if matches := failRe.FindStringSubmatch(output); len(matches) > 1 {
			result.TestsFailed, _ = strconv.Atoi(matches[1])
			result.TestsPassed = result.TotalTests - result.TestsFailed
		}
	}
}

// parseGoTestOutput parses go test output
func (te *TestExecutor) parseGoTestOutput(output string, result *TestResult) {
	// Go test output format: "PASS" or "FAIL" for each test
	// Count PASS and FAIL lines
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "--- PASS:") {
			result.TestsPassed++
		} else if strings.HasPrefix(strings.TrimSpace(line), "--- FAIL:") {
			result.TestsFailed++
		} else if strings.HasPrefix(strings.TrimSpace(line), "--- SKIP:") {
			result.TestsSkipped++
		}
	}

	result.TotalTests = result.TestsPassed + result.TestsFailed + result.TestsSkipped
}
