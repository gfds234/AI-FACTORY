package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Demo project specification
const demoSpec = `Create a modern React Task Manager app with Vite and the following features:

1. Add tasks with title and description
2. Mark tasks as complete/incomplete
3. Delete tasks
4. Filter tasks (all/active/completed)
5. Clean, professional UI
6. Local storage persistence
7. Responsive design

Tech requirements:
- React 18 with Vite
- Clean component architecture (TaskList, TaskItem, TaskForm components)
- Proper state management with useState
- Vitest for testing with at least 3 component tests
- Professional README with actual setup instructions

The app must be production-ready and immediately runnable.
Make this a portfolio-quality implementation with real, working code.`

func main() {
	fmt.Println("🚀 AI FACTORY Demo Generator")
	fmt.Println("=" + string(make([]byte, 58)))
	fmt.Println("Generating sellable-quality React Task Manager demo...")
	fmt.Println()

	baseURL := "http://localhost:8080"

	// Check if server is running
	if !checkServer(baseURL) {
		log.Fatal("❌ Server not running at " + baseURL + "\nPlease start with: go run main.go -mode=server")
	}

	fmt.Println("✓ Server is running")
	fmt.Println()

	// Step 1: Create project
	fmt.Println("Step 1: Creating project...")
	projectID, err := createProject(baseURL)
	if err != nil {
		log.Fatalf("❌ Failed to create project: %v", err)
	}
	fmt.Printf("✓ Project created: %s\n\n", projectID)

	// Step 2: Generate code directly (skip discovery/validation/planning for speed)
	fmt.Println("Step 2: Generating code with AI Factory...")
	codeResult, err := generateCode(baseURL, projectID)
	if err != nil {
		log.Fatalf("❌ Code generation failed: %v", err)
	}

	fmt.Printf("✓ Code generation complete\n")
	fmt.Printf("  Files generated: Check artifacts directory\n")
	fmt.Printf("  Artifact path: %s\n\n", codeResult["artifact_path"])

	// Step 3: Check results
	fmt.Println("Step 3: Checking generation results...")
	time.Sleep(2 * time.Second)

	// Try to find the generated project
	projectPath := findGeneratedProject()
	if projectPath == "" {
		fmt.Println("⚠️  No generated project found in projects/ directory")
		fmt.Println("   Check artifacts/ directory for LLM output")
		os.Exit(1)
	}

	fmt.Printf("✓ Found generated project: %s\n\n", projectPath)

	// Step 4: Verify project files
	fmt.Println("Step 4: Verifying project structure...")
	files := verifyProjectStructure(projectPath)
	fmt.Printf("✓ Project has %d files\n", len(files))
	for _, f := range files[:min(10, len(files))] {
		fmt.Printf("  - %s\n", f)
	}
	if len(files) > 10 {
		fmt.Printf("  ... and %d more files\n", len(files)-10)
	}
	fmt.Println()

	// Step 5: Final report
	fmt.Println("=" + string(make([]byte, 58)))
	fmt.Println("📊 DEMO GENERATION COMPLETE")
	fmt.Println("=" + string(make([]byte, 58)))
	fmt.Printf("Project Location: %s\n", projectPath)
	fmt.Printf("Total Files: %d\n", len(files))
	fmt.Println()
	fmt.Println("Next steps to test the demo:")
	fmt.Println("1. cd " + projectPath)
	fmt.Println("2. npm install")
	fmt.Println("3. npm run dev")
	fmt.Println("4. Open http://localhost:5173")
	fmt.Println()
	fmt.Println("To run tests:")
	fmt.Println("  npm run test")
	fmt.Println()
	fmt.Println("🎉 Demo ready to show customers!")
}

func checkServer(baseURL string) bool {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(baseURL + "/")
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == 200
}

func createProject(baseURL string) (string, error) {
	reqBody := map[string]string{
		"name":        "React Task Manager Demo",
		"description": demoSpec,
	}

	jsonData, _ := json.Marshal(reqBody)
	resp, err := http.Post(baseURL+"/project", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	projectID, ok := result["project_id"].(string)
	if !ok {
		return "", fmt.Errorf("no project_id in response")
	}

	return projectID, nil
}

func generateCode(baseURL string, projectID string) (map[string]interface{}, error) {
	reqBody := map[string]string{
		"input": demoSpec,
		"type":  "code",
	}

	jsonData, _ := json.Marshal(reqBody)
	resp, err := http.Post(baseURL+"/task", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func findGeneratedProject() string {
	projectsDir := "projects"
	entries, err := os.ReadDir(projectsDir)
	if err != nil {
		return ""
	}

	// Find the most recent generated_* directory
	var latestDir string
	var latestTime int64

	for _, entry := range entries {
		if entry.IsDir() && len(entry.Name()) > 10 && entry.Name()[:10] == "generated_" {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			if info.ModTime().Unix() > latestTime {
				latestTime = info.ModTime().Unix()
				latestDir = filepath.Join(projectsDir, entry.Name())
			}
		}
	}

	return latestDir
}

func verifyProjectStructure(projectPath string) []string {
	var files []string
	filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			relPath, _ := filepath.Rel(projectPath, path)
			files = append(files, relPath)
		}
		return nil
	})
	return files
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
