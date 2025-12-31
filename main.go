package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"ai-studio/orchestrator/api"
	"ai-studio/orchestrator/supervisor"
	"ai-studio/orchestrator/task"
)

func main() {
	// CLI flags
	mode := flag.String("mode", "server", "Run mode: server or cli")
	taskType := flag.String("task", "", "Task type for CLI mode: validate or review")
	input := flag.String("input", "", "Input file path for CLI mode")
	port := flag.Int("port", 8080, "Server port")
	flag.Parse()

	// Load configuration with supervisor settings
	baseConfig, supervisorConfig, err := supervisor.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize base task manager
	baseMgr := task.NewManager(baseConfig)

	// Wrap with supervisor if enabled
	var taskMgr api.TaskManager
	if supervisorConfig.Enabled {
		taskMgr = supervisor.NewSupervisedTaskManager(baseMgr, baseConfig, supervisorConfig)
		log.Printf("✓ Supervisor enabled (complexity threshold: %d)", supervisorConfig.ComplexityThreshold)
	} else {
		taskMgr = baseMgr
		log.Printf("Supervisor disabled - using standard task manager")
	}

	switch *mode {
	case "server":
		// Start HTTP server
		server := api.NewServer(taskMgr, *port)
		log.Printf("Starting AI Studio Orchestrator on port %d", *port)
		log.Printf("Models: %v", baseConfig.Models)
		if err := server.Start(); err != nil {
			log.Fatalf("Server failed: %v", err)
		}

	case "cli":
		// CLI mode for testing
		if *taskType == "" || *input == "" {
			fmt.Println("Usage: orchestrator -mode=cli -task=<validate|review> -input=<file>")
			os.Exit(1)
		}

		inputData, err := os.ReadFile(*input)
		if err != nil {
			log.Fatalf("Failed to read input: %v", err)
		}

		result, err := taskMgr.ExecuteTask(*taskType, string(inputData))
		if err != nil {
			log.Fatalf("Task execution failed: %v", err)
		}

		// Extract output from result (handle both standard and supervised results)
		if taskResult, ok := result.(*task.Result); ok {
			fmt.Printf("\n=== Task Result ===\n%s\n", taskResult.Output)
			fmt.Printf("\nArtifact saved: %s\n", taskResult.ArtifactPath)
		} else {
			fmt.Printf("\n=== Task Result ===\n%+v\n", result)
		}

	default:
		log.Fatalf("Unknown mode: %s", *mode)
	}
}
