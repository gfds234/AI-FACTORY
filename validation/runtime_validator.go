package validation

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// RuntimeValidator validates that generated code actually runs
type RuntimeValidator struct {
	timeout time.Duration
}

// NewRuntimeValidator creates a new runtime validator
func NewRuntimeValidator() *RuntimeValidator {
	return &RuntimeValidator{
		timeout: 10 * time.Second, // 10 seconds for server startup
	}
}

// RuntimeResult contains runtime validation results
type RuntimeResult struct {
	ApplicationStarts bool     `json:"application_starts"` // true if app starts without crashing
	HealthCheckPassed bool     `json:"health_check_passed"` // true if server responds to health check
	Port              int      `json:"port"`                // port the server is listening on
	Errors            []string `json:"errors"`              // runtime errors
	Warnings          []string `json:"warnings"`            // non-fatal issues
	RuntimeLog        string   `json:"runtime_log"`         // execution output
}

// ValidateRuntime validates that a project runs without crashing
func (rv *RuntimeValidator) ValidateRuntime(projectPath string, projectType string) (*RuntimeResult, error) {
	result := &RuntimeResult{
		Errors:   make([]string, 0),
		Warnings: make([]string, 0),
	}

	switch projectType {
	case "nodejs":
		return rv.validateNodeJS(projectPath)
	case "python":
		return rv.validatePython(projectPath)
	case "go":
		return rv.validateGo(projectPath)
	case "frontend":
		return rv.validateFrontend(projectPath)
	default:
		result.Warnings = append(result.Warnings, fmt.Sprintf("Unknown project type %s - skipping runtime validation", projectType))
		return result, nil
	}
}

// validateNodeJS validates Node.js runtime
func (rv *RuntimeValidator) validateNodeJS(projectPath string) (*RuntimeResult, error) {
	result := &RuntimeResult{
		Errors:   make([]string, 0),
		Warnings: make([]string, 0),
		Port:     3000, // Default Node.js port
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

	// Start the server
	cmd := exec.Command("node", "server.js")
	cmd.Dir = projectPath

	// Capture output
	outputBuilder := &strings.Builder{}
	cmd.Stdout = outputBuilder
	cmd.Stderr = outputBuilder

	// Start server
	if err := cmd.Start(); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to start server: %v", err))
		result.RuntimeLog = outputBuilder.String()
		return result, nil
	}

	result.ApplicationStarts = true

	// Wait for server to start
	time.Sleep(3 * time.Second)

	// Check if process is still running
	if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
		result.ApplicationStarts = false
		result.Errors = append(result.Errors, "Server exited immediately after start")
		result.RuntimeLog = outputBuilder.String()
		return result, nil
	}

	// Try health check on common endpoints
	endpoints := []string{
		fmt.Sprintf("http://localhost:%d/", result.Port),
		fmt.Sprintf("http://localhost:%d/health", result.Port),
		fmt.Sprintf("http://localhost:%d/api", result.Port),
	}

	healthCheckPassed := false
	for _, endpoint := range endpoints {
		client := &http.Client{Timeout: 2 * time.Second}
		resp, err := client.Get(endpoint)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode < 500 {
				healthCheckPassed = true
				break
			}
		}
	}

	result.HealthCheckPassed = healthCheckPassed
	if !healthCheckPassed {
		result.Warnings = append(result.Warnings, "Server started but health check failed")
	}

	// Kill the server
	if cmd.Process != nil {
		cmd.Process.Kill()
	}

	result.RuntimeLog = outputBuilder.String()
	return result, nil
}

// validatePython validates Python runtime
func (rv *RuntimeValidator) validatePython(projectPath string) (*RuntimeResult, error) {
	result := &RuntimeResult{
		Errors:   make([]string, 0),
		Warnings: make([]string, 0),
		Port:     8000, // Default Python/FastAPI port
	}

	// Find entry point
	entryPoints := []string{"main.py", "app.py", "server.py", "backend/main.py"}
	var entryPoint string
	var workDir string

	for _, ep := range entryPoints {
		fullPath := filepath.Join(projectPath, ep)
		if _, err := os.Stat(fullPath); err == nil {
			entryPoint = filepath.Base(ep)
			workDir = filepath.Dir(fullPath)
			if workDir == "." {
				workDir = projectPath
			}
			break
		}
	}

	if entryPoint == "" {
		result.Errors = append(result.Errors, "No Python entry point found (main.py, app.py, server.py)")
		return result, nil
	}

	// Start the server
	cmd := exec.Command("python", entryPoint)
	cmd.Dir = workDir

	// Capture output
	outputBuilder := &strings.Builder{}
	cmd.Stdout = outputBuilder
	cmd.Stderr = outputBuilder

	// Start server
	if err := cmd.Start(); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to start: %v", err))
		result.RuntimeLog = outputBuilder.String()
		return result, nil
	}

	result.ApplicationStarts = true

	// Wait for server to start
	time.Sleep(3 * time.Second)

	// Check if process is still running
	if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
		result.ApplicationStarts = false
		result.Errors = append(result.Errors, "Application exited immediately")
		result.RuntimeLog = outputBuilder.String()
		return result, nil
	}

	// Try health check
	endpoints := []string{
		fmt.Sprintf("http://localhost:%d/", result.Port),
		fmt.Sprintf("http://localhost:%d/docs", result.Port), // FastAPI docs
	}

	healthCheckPassed := false
	for _, endpoint := range endpoints {
		client := &http.Client{Timeout: 2 * time.Second}
		resp, err := client.Get(endpoint)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode < 500 {
				healthCheckPassed = true
				break
			}
		}
	}

	result.HealthCheckPassed = healthCheckPassed

	// Kill the server
	if cmd.Process != nil {
		cmd.Process.Kill()
	}

	result.RuntimeLog = outputBuilder.String()
	return result, nil
}

// validateGo validates Go runtime
func (rv *RuntimeValidator) validateGo(projectPath string) (*RuntimeResult, error) {
	result := &RuntimeResult{
		Errors:   make([]string, 0),
		Warnings: make([]string, 0),
		Port:     8080, // Default Go port
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

	// Start the server
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = workDir

	// Capture output
	outputBuilder := &strings.Builder{}
	cmd.Stdout = outputBuilder
	cmd.Stderr = outputBuilder

	// Start server
	if err := cmd.Start(); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to start: %v", err))
		result.RuntimeLog = outputBuilder.String()
		return result, nil
	}

	result.ApplicationStarts = true

	// Wait for server to start
	time.Sleep(3 * time.Second)

	// Check if process is still running
	if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
		result.ApplicationStarts = false
		result.Errors = append(result.Errors, "Server exited immediately")
		result.RuntimeLog = outputBuilder.String()
		return result, nil
	}

	// Try health check
	endpoints := []string{
		fmt.Sprintf("http://localhost:%d/", result.Port),
		fmt.Sprintf("http://localhost:%d/health", result.Port),
	}

	healthCheckPassed := false
	for _, endpoint := range endpoints {
		client := &http.Client{Timeout: 2 * time.Second}
		resp, err := client.Get(endpoint)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode < 500 {
				healthCheckPassed = true
				break
			}
		}
	}

	result.HealthCheckPassed = healthCheckPassed

	// Kill the server
	if cmd.Process != nil {
		cmd.Process.Kill()
	}

	result.RuntimeLog = outputBuilder.String()
	return result, nil
}

// validateFrontend validates frontend runtime (basic check)
func (rv *RuntimeValidator) validateFrontend(projectPath string) (*RuntimeResult, error) {
	result := &RuntimeResult{
		ApplicationStarts: true, // Frontend is static, no server to start
		HealthCheckPassed: true, // If index.html exists and is valid, we consider it healthy
		Errors:            make([]string, 0),
		Warnings:          make([]string, 0),
	}

	// Check if index.html exists and is readable
	indexPath := filepath.Join(projectPath, "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		result.ApplicationStarts = false
		result.HealthCheckPassed = false
		result.Errors = append(result.Errors, "index.html not found")
		return result, nil
	}

	// For frontend, we could potentially use a headless browser here
	// For now, we just verify the file exists
	result.RuntimeLog = "Frontend project - static files verified\n"
	return result, nil
}
