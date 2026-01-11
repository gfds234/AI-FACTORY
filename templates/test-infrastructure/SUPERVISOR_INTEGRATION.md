# Supervisor Integration - Zero-Bug Workflow

This document explains how to integrate automated testing into the AI Factory supervisor.

## Changes Required in supervisor/agent_verification.go

### 1. Update VerificationResult Struct

Add test-related fields to the `VerificationResult` struct:

```go
type VerificationResult struct {
	ProjectType     string   `json:"project_type"`
	SyntaxValid     bool     `json:"syntax_valid"`
	DependenciesOK  bool     `json:"dependencies_ok"`
	EntryPointValid bool     `json:"entry_point_valid"`

	// NEW: Test validation fields
	HasTests        bool     `json:"has_tests"`         // true if test files exist
	TestsPassed     bool     `json:"tests_passed"`      // true if all tests pass
	TestCoverage    float64  `json:"test_coverage"`     // percentage (0-100)
	TestOutput      string   `json:"test_output"`       // test run output

	Errors          []string `json:"errors"`
	Warnings        []string `json:"warnings"`
	BuildLog        string   `json:"build_log"`
}
```

### 2. Add Test Validation Method

Add new method to run tests:

```go
// RunTests executes the test suite and returns results
func (va *VerificationAgent) RunTests(projectPath string, projectType ProjectType) (bool, string, error) {
	switch projectType {
	case ProjectTypeVite, ProjectTypeNodeJS:
		return va.runNodeTests(projectPath)
	case ProjectTypePython:
		return va.runPythonTests(projectPath)
	case ProjectTypeGo:
		return va.runGoTests(projectPath)
	default:
		return false, "", fmt.Errorf("testing not supported for project type: %s", projectType)
	}
}

// runNodeTests runs npm test for Node.js/Vite projects
func (va *VerificationAgent) runNodeTests(projectPath string) (bool, string, error) {
	// Check if test script exists
	packageJSON := filepath.Join(projectPath, "package.json")
	content, err := ioutil.ReadFile(packageJSON)
	if err != nil {
		return false, "", fmt.Errorf("cannot read package.json: %v", err)
	}

	// Check for test script
	if !strings.Contains(string(content), `"test"`) {
		return false, "No test script found in package.json", nil
	}

	// Run tests
	cmd := exec.Command("npm", "run", "test")
	cmd.Dir = projectPath
	output, err := cmd.CombinedOutput()

	if err != nil {
		return false, string(output), fmt.Errorf("tests failed: %v", err)
	}

	return true, string(output), nil
}

// runPythonTests runs pytest for Python projects
func (va *VerificationAgent) runPythonTests(projectPath string) (bool, string, error) {
	// Check for test files
	testsExist := false
	testDirs := []string{"tests", "test", "backend/tests"}
	for _, dir := range testDirs {
		if _, err := os.Stat(filepath.Join(projectPath, dir)); err == nil {
			testsExist = true
			break
		}
	}

	if !testsExist {
		return false, "No test directory found", nil
	}

	// Run pytest
	cmd := exec.Command("pytest", "-v", "--tb=short")
	cmd.Dir = projectPath
	output, err := cmd.CombinedOutput()

	if err != nil {
		return false, string(output), fmt.Errorf("tests failed: %v", err)
	}

	return true, string(output), nil
}

// runGoTests runs go test for Go projects
func (va *VerificationAgent) runGoTests(projectPath string) (bool, string, error) {
	cmd := exec.Command("go", "test", "./...", "-v")
	cmd.Dir = projectPath
	output, err := cmd.CombinedOutput()

	if err != nil {
		return false, string(output), fmt.Errorf("tests failed: %v", err)
	}

	return true, string(output), nil
}
```

### 3. Update VerifyProject Method

Integrate test execution into the verification flow:

```go
func (va *VerificationAgent) VerifyProject(projectPath string) (*VerificationResult, error) {
	result := &VerificationResult{
		Errors:   make([]string, 0),
		Warnings: make([]string, 0),
	}

	// ... existing detection code ...

	// NEW: Check for test files
	hasTests := va.checkForTests(projectPath, projectType)
	result.HasTests = hasTests

	if !hasTests {
		result.Warnings = append(result.Warnings, "No test files found - consider adding tests")
	}

	// NEW: Run tests if they exist
	if hasTests {
		testsPassed, testOutput, testErr := va.RunTests(projectPath, projectType)
		result.TestsPassed = testsPassed
		result.TestOutput = testOutput

		if testErr != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Test execution failed: %v", testErr))
		}

		if !testsPassed {
			result.Errors = append(result.Errors, "Tests failed - fix failing tests before deploying")
		}
	}

	return result, nil
}

// checkForTests checks if test files exist
func (va *VerificationAgent) checkForTests(projectPath string, projectType ProjectType) bool {
	switch projectType {
	case ProjectTypeVite, ProjectTypeNodeJS:
		// Check for __tests__, *.test.js, *.spec.js
		patterns := []string{
			"src/__tests__",
			"tests",
			"__tests__",
		}
		for _, pattern := range patterns {
			if _, err := os.Stat(filepath.Join(projectPath, pattern)); err == nil {
				return true
			}
		}
		return false

	case ProjectTypePython:
		testDirs := []string{"tests", "test", "backend/tests"}
		for _, dir := range testDirs {
			if _, err := os.Stat(filepath.Join(projectPath, dir)); err == nil {
				return true
			}
		}
		return false

	case ProjectTypeGo:
		// Look for *_test.go files
		// (Implementation would use filepath.Walk to find test files)
		return false

	default:
		return false
	}
}
```

### 4. Update Quality Scoring

Modify the quality score calculation to include test results:

```go
func CalculateQualityScore(result *VerificationResult) int {
	score := 0

	// Syntax valid: 30 points
	if result.SyntaxValid {
		score += 30
	}

	// Dependencies OK: 20 points
	if result.DependenciesOK {
		score += 20
	}

	// Entry point valid: 20 points
	if result.EntryPointValid {
		score += 20
	}

	// NEW: Has tests: 10 points
	if result.HasTests {
		score += 10
	}

	// NEW: Tests pass: 20 points (CRITICAL)
	if result.TestsPassed {
		score += 20
	} else if result.HasTests {
		// Has tests but they failed: -10 penalty
		score -= 10
	}

	// Ensure score stays in 0-100 range
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}
```

## Enforcement Rules

### For supervisor/quality_check.go

Update quality check logic to REJECT projects with failing tests:

```go
func (s *Supervisor) QualityCheck(projectPath string) error {
	// ... existing verification ...

	// MANDATORY: Tests must pass
	if result.HasTests && !result.TestsPassed {
		return fmt.Errorf("QUALITY CHECK FAILED: Tests are failing. Fix tests before proceeding.\n\nTest Output:\n%s", result.TestOutput)
	}

	// RECOMMENDED: Warn if no tests exist
	if !result.HasTests {
		fmt.Println("WARNING: No tests found. Consider adding automated tests.")
	}

	// Calculate quality score
	score := CalculateQualityScore(result)

	// MANDATORY: Minimum quality score of 80 required
	if score < 80 {
		return fmt.Errorf("QUALITY CHECK FAILED: Quality score %d is below minimum threshold of 80", score)
	}

	return nil
}
```

## Project Template Integration

### Auto-copy test infrastructure to new projects

In `project/orchestrator.go` or wherever projects are created:

```go
func (o *Orchestrator) CreateProject(projectPath string) error {
	// ... existing project creation ...

	// Copy test infrastructure template
	templatePath := "templates/test-infrastructure"
	if err := copyTestInfrastructure(templatePath, projectPath); err != nil {
		return fmt.Errorf("failed to copy test infrastructure: %v", err)
	}

	return nil
}

func copyTestInfrastructure(templatePath, projectPath string) error {
	// Copy vitest.config.js
	copyFile(
		filepath.Join(templatePath, "vitest.config.js"),
		filepath.Join(projectPath, "vitest.config.js"),
	)

	// Copy playwright.config.js
	copyFile(
		filepath.Join(templatePath, "playwright.config.js"),
		filepath.Join(projectPath, "playwright.config.js"),
	)

	// Create test directories
	os.MkdirAll(filepath.Join(projectPath, "src", "__tests__"), 0755)
	os.MkdirAll(filepath.Join(projectPath, "e2e"), 0755)
	os.MkdirAll(filepath.Join(projectPath, "tests", "integration"), 0755)

	// Copy validation script
	os.MkdirAll(filepath.Join(projectPath, "scripts"), 0755)
	copyFile(
		filepath.Join(templatePath, "scripts", "validate.js"),
		filepath.Join(projectPath, "scripts", "validate.js"),
	)

	// Update package.json with test scripts
	updatePackageJSONWithTestScripts(projectPath)

	return nil
}
```

## Summary

The integration requires:

1. ✅ Update `VerificationResult` struct with test fields
2. ✅ Add `RunTests()` method and language-specific test runners
3. ✅ Add `checkForTests()` method
4. ✅ Update `VerifyProject()` to execute tests
5. ✅ Update quality scoring to include test results
6. ✅ Enforce test passing in quality check
7. ✅ Auto-copy test infrastructure to new projects

## Testing the Integration

```bash
# Create test project
cd AI\ FACTORY
go run . create-project test-zero-bug

# Verify test infrastructure was copied
ls test-zero-bug/
# Should see: vitest.config.js, playwright.config.js, src/__tests__/, scripts/validate.js

# Run verification
go run . verify test-zero-bug

# Should output:
# ✅ Tests found: 3 files
# ✅ All tests passed (3/3)
# ✅ Quality Score: 100
```

## Rollout Plan

1. **Phase 1**: Add test validation to supervisor (non-blocking warnings)
2. **Phase 2**: Make test infrastructure template available
3. **Phase 3**: Auto-copy template to new projects
4. **Phase 4**: Make failing tests BLOCK deployment (enforced)
5. **Phase 5**: Require minimum test coverage (80%)

**Status**: Template created, ready for integration into supervisor
