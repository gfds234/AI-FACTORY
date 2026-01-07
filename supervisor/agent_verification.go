package supervisor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// VerificationAgent validates generated code (syntax, dependencies, build)
type VerificationAgent struct {
	timeout time.Duration
}

// NewVerificationAgent creates a new verification agent
func NewVerificationAgent() *VerificationAgent {
	return &VerificationAgent{
		timeout: 120 * time.Second, // 2 minutes for dependency installation
	}
}

// VerificationResult contains verification results
type VerificationResult struct {
	ProjectType     string   `json:"project_type"`      // "nodejs", "python", "go", "frontend"
	SyntaxValid     bool     `json:"syntax_valid"`      // true if code parses without errors
	DependenciesOK  bool     `json:"dependencies_ok"`   // true if all imports/requires work
	EntryPointValid bool     `json:"entry_point_valid"` // true if main entry point exists
	Errors          []string `json:"errors"`            // compilation/syntax errors
	Warnings        []string `json:"warnings"`          // non-fatal issues
	BuildLog        string   `json:"build_log"`         // full build output
}

// ProjectType represents detected project type
type ProjectType string

const (
	ProjectTypeNodeJS   ProjectType = "nodejs"
	ProjectTypePython   ProjectType = "python"
	ProjectTypeGo       ProjectType = "go"
	ProjectTypeFrontend ProjectType = "frontend"
	ProjectTypeUnknown  ProjectType = "unknown"
)

// VerifyProject verifies a generated project
func (va *VerificationAgent) VerifyProject(projectPath string) (*VerificationResult, error) {
	result := &VerificationResult{
		Errors:   make([]string, 0),
		Warnings: make([]string, 0),
	}

	// Detect project type
	projectType, err := va.detectProjectType(projectPath)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to detect project type: %v", err))
		return result, nil
	}
	result.ProjectType = string(projectType)

	if projectType == ProjectTypeUnknown {
		result.Warnings = append(result.Warnings, "Could not determine project type - skipping validation")
		return result, nil
	}

	// Check entry point exists
	entryPointValid, entryPointErr := va.checkEntryPoint(projectPath, projectType)
	result.EntryPointValid = entryPointValid
	if entryPointErr != nil {
		result.Errors = append(result.Errors, entryPointErr.Error())
	}

	// Install dependencies
	depsOK, depsLog, depsErr := va.installDependencies(projectPath, projectType)
	result.DependenciesOK = depsOK
	result.BuildLog += depsLog

	if depsErr != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Dependency installation failed: %v", depsErr))
	}

	// Validate syntax
	syntaxValid, syntaxLog, syntaxErr := va.validateSyntax(projectPath, projectType)
	result.SyntaxValid = syntaxValid
	result.BuildLog += syntaxLog

	if syntaxErr != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Syntax validation failed: %v", syntaxErr))
	}

	return result, nil
}

// detectProjectType detects the type of project based on files present
func (va *VerificationAgent) detectProjectType(projectPath string) (ProjectType, error) {
	// Check for package.json (Node.js)
	if _, err := os.Stat(filepath.Join(projectPath, "package.json")); err == nil {
		return ProjectTypeNodeJS, nil
	}

	// Check for backend/package.json (Node.js backend)
	if _, err := os.Stat(filepath.Join(projectPath, "backend", "package.json")); err == nil {
		return ProjectTypeNodeJS, nil
	}

	// Check for requirements.txt or setup.py (Python)
	if _, err := os.Stat(filepath.Join(projectPath, "requirements.txt")); err == nil {
		return ProjectTypePython, nil
	}
	if _, err := os.Stat(filepath.Join(projectPath, "setup.py")); err == nil {
		return ProjectTypePython, nil
	}
	if _, err := os.Stat(filepath.Join(projectPath, "backend", "requirements.txt")); err == nil {
		return ProjectTypePython, nil
	}

	// Check for go.mod (Go)
	if _, err := os.Stat(filepath.Join(projectPath, "go.mod")); err == nil {
		return ProjectTypeGo, nil
	}
	if _, err := os.Stat(filepath.Join(projectPath, "backend", "go.mod")); err == nil {
		return ProjectTypeGo, nil
	}

	// Check for index.html (Frontend-only)
	if _, err := os.Stat(filepath.Join(projectPath, "index.html")); err == nil {
		return ProjectTypeFrontend, nil
	}

	return ProjectTypeUnknown, nil
}

// checkEntryPoint verifies the entry point exists
func (va *VerificationAgent) checkEntryPoint(projectPath string, projectType ProjectType) (bool, error) {
	switch projectType {
	case ProjectTypeNodeJS:
		// Check for server.js, index.js, app.js, or main.js
		entryPoints := []string{"server.js", "backend/server.js", "index.js", "app.js", "main.js"}
		for _, ep := range entryPoints {
			if _, err := os.Stat(filepath.Join(projectPath, ep)); err == nil {
				return true, nil
			}
		}
		return false, fmt.Errorf("no Node.js entry point found (server.js, index.js, app.js, or main.js)")

	case ProjectTypePython:
		// Check for main.py, app.py, server.py
		entryPoints := []string{"main.py", "backend/main.py", "app.py", "server.py"}
		for _, ep := range entryPoints {
			if _, err := os.Stat(filepath.Join(projectPath, ep)); err == nil {
				return true, nil
			}
		}
		return false, fmt.Errorf("no Python entry point found (main.py, app.py, or server.py)")

	case ProjectTypeGo:
		// Check for main.go
		entryPoints := []string{"main.go", "backend/main.go"}
		for _, ep := range entryPoints {
			if _, err := os.Stat(filepath.Join(projectPath, ep)); err == nil {
				return true, nil
			}
		}
		return false, fmt.Errorf("no Go entry point found (main.go)")

	case ProjectTypeFrontend:
		// Check for index.html
		if _, err := os.Stat(filepath.Join(projectPath, "index.html")); err == nil {
			return true, nil
		}
		return false, fmt.Errorf("no HTML entry point found (index.html)")
	}

	return false, fmt.Errorf("unknown project type")
}

// installDependencies installs project dependencies
func (va *VerificationAgent) installDependencies(projectPath string, projectType ProjectType) (bool, string, error) {
	var cmd *exec.Cmd
	var workDir string

	switch projectType {
	case ProjectTypeNodeJS:
		// Try both root and backend directories
		if _, err := os.Stat(filepath.Join(projectPath, "package.json")); err == nil {
			workDir = projectPath
		} else if _, err := os.Stat(filepath.Join(projectPath, "backend", "package.json")); err == nil {
			workDir = filepath.Join(projectPath, "backend")
		} else {
			return false, "", fmt.Errorf("no package.json found")
		}
		cmd = exec.Command("npm", "install")
		cmd.Dir = workDir

	case ProjectTypePython:
		// Try both root and backend directories
		if _, err := os.Stat(filepath.Join(projectPath, "requirements.txt")); err == nil {
			workDir = projectPath
		} else if _, err := os.Stat(filepath.Join(projectPath, "backend", "requirements.txt")); err == nil {
			workDir = filepath.Join(projectPath, "backend")
		} else {
			// No requirements.txt - not necessarily an error for simple scripts
			return true, "No requirements.txt found - assuming no dependencies\n", nil
		}
		cmd = exec.Command("pip", "install", "-r", "requirements.txt")
		cmd.Dir = workDir

	case ProjectTypeGo:
		// Try both root and backend directories
		if _, err := os.Stat(filepath.Join(projectPath, "go.mod")); err == nil {
			workDir = projectPath
		} else if _, err := os.Stat(filepath.Join(projectPath, "backend", "go.mod")); err == nil {
			workDir = filepath.Join(projectPath, "backend")
		} else {
			return false, "", fmt.Errorf("no go.mod found")
		}
		cmd = exec.Command("go", "mod", "download")
		cmd.Dir = workDir

	case ProjectTypeFrontend:
		// Frontend-only projects don't need dependency installation
		return true, "Frontend-only project - no dependencies to install\n", nil

	default:
		return false, "", fmt.Errorf("unsupported project type for dependency installation")
	}

	// Run installation with timeout
	output, err := va.runCommandWithTimeout(cmd, va.timeout)
	if err != nil {
		return false, string(output), err
	}

	return true, string(output), nil
}

// validateSyntax validates code syntax
func (va *VerificationAgent) validateSyntax(projectPath string, projectType ProjectType) (bool, string, error) {
	switch projectType {
	case ProjectTypeNodeJS:
		return va.validateNodeJSSyntax(projectPath)

	case ProjectTypePython:
		return va.validatePythonSyntax(projectPath)

	case ProjectTypeGo:
		return va.validateGoSyntax(projectPath)

	case ProjectTypeFrontend:
		return va.validateHTMLSyntax(projectPath)

	default:
		return false, "", fmt.Errorf("unsupported project type for syntax validation")
	}
}

// validateNodeJSSyntax validates Node.js syntax
func (va *VerificationAgent) validateNodeJSSyntax(projectPath string) (bool, string, error) {
	// Find all .js files
	jsFiles, err := va.findFiles(projectPath, ".js")
	if err != nil {
		return false, "", err
	}

	if len(jsFiles) == 0 {
		return false, "", fmt.Errorf("no JavaScript files found")
	}

	var allOutput strings.Builder
	allValid := true

	// Validate each JS file
	for _, jsFile := range jsFiles {
		cmd := exec.Command("node", "--check", jsFile)
		output, err := va.runCommandWithTimeout(cmd, 10*time.Second)
		allOutput.WriteString(string(output))

		if err != nil {
			allValid = false
			allOutput.WriteString(fmt.Sprintf("Syntax error in %s: %v\n", jsFile, err))
		}
	}

	if allValid {
		allOutput.WriteString("All JavaScript files have valid syntax\n")
	}

	return allValid, allOutput.String(), nil
}

// validatePythonSyntax validates Python syntax
func (va *VerificationAgent) validatePythonSyntax(projectPath string) (bool, string, error) {
	// Find all .py files
	pyFiles, err := va.findFiles(projectPath, ".py")
	if err != nil {
		return false, "", err
	}

	if len(pyFiles) == 0 {
		return false, "", fmt.Errorf("no Python files found")
	}

	var allOutput strings.Builder
	allValid := true

	// Validate each Python file
	for _, pyFile := range pyFiles {
		cmd := exec.Command("python", "-m", "py_compile", pyFile)
		output, err := va.runCommandWithTimeout(cmd, 10*time.Second)
		allOutput.WriteString(string(output))

		if err != nil {
			allValid = false
			allOutput.WriteString(fmt.Sprintf("Syntax error in %s: %v\n", pyFile, err))
		}
	}

	if allValid {
		allOutput.WriteString("All Python files have valid syntax\n")
	}

	return allValid, allOutput.String(), nil
}

// validateGoSyntax validates Go syntax
func (va *VerificationAgent) validateGoSyntax(projectPath string) (bool, string, error) {
	// Try both root and backend directories
	var workDir string
	if _, err := os.Stat(filepath.Join(projectPath, "go.mod")); err == nil {
		workDir = projectPath
	} else if _, err := os.Stat(filepath.Join(projectPath, "backend", "go.mod")); err == nil {
		workDir = filepath.Join(projectPath, "backend")
	} else {
		return false, "", fmt.Errorf("no go.mod found")
	}

	cmd := exec.Command("go", "build", "-o", "nul") // Build to nul (discard output on Windows)
	cmd.Dir = workDir

	output, err := va.runCommandWithTimeout(cmd, 30*time.Second)
	if err != nil {
		return false, string(output), err
	}

	return true, string(output) + "Go code builds successfully\n", nil
}

// validateHTMLSyntax validates HTML syntax (basic check)
func (va *VerificationAgent) validateHTMLSyntax(projectPath string) (bool, string, error) {
	// For frontend projects, check if index.html exists and is readable
	indexPath := filepath.Join(projectPath, "index.html")
	content, err := os.ReadFile(indexPath)
	if err != nil {
		return false, "", fmt.Errorf("failed to read index.html: %w", err)
	}

	// Basic HTML validation - check for required tags
	htmlContent := string(content)
	warnings := make([]string, 0)

	if !strings.Contains(htmlContent, "<html") {
		warnings = append(warnings, "Missing <html> tag")
	}
	if !strings.Contains(htmlContent, "<head") {
		warnings = append(warnings, "Missing <head> tag")
	}
	if !strings.Contains(htmlContent, "<body") {
		warnings = append(warnings, "Missing <body> tag")
	}

	if len(warnings) > 0 {
		return false, strings.Join(warnings, "\n"), nil
	}

	return true, "HTML structure is valid\n", nil
}

// findFiles finds all files with a given extension
func (va *VerificationAgent) findFiles(root string, ext string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip node_modules, venv, __pycache__, .git directories
		if info.IsDir() {
			name := info.Name()
			if name == "node_modules" || name == "venv" || name == "__pycache__" || name == ".git" {
				return filepath.SkipDir
			}
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ext) {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// runCommandWithTimeout runs a command with a timeout
func (va *VerificationAgent) runCommandWithTimeout(cmd *exec.Cmd, timeout time.Duration) ([]byte, error) {
	// Create a context with timeout
	done := make(chan error, 1)
	var output []byte

	go func() {
		var err error
		output, err = cmd.CombinedOutput()
		done <- err
	}()

	select {
	case <-time.After(timeout):
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return output, fmt.Errorf("command timed out after %v", timeout)
	case err := <-done:
		return output, err
	}
}
