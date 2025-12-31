package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"ai-studio/orchestrator/api"
	"ai-studio/orchestrator/config"
	"ai-studio/orchestrator/task"
)

func main() {
	// CLI flags
	mode := flag.String("mode", "server", "Run mode: server or cli")
	taskType := flag.String("task", "", "Task type for CLI mode: validate or review")
	input := flag.String("input", "", "Input file path for CLI mode")
	port := flag.Int("port", 8080, "Server port")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize task manager
	taskMgr := task.NewManager(cfg)

	switch *mode {
	case "server":
		// Start HTTP server
		server := api.NewServer(taskMgr, *port)
		log.Printf("Starting AI Studio Orchestrator on port %d", *port)
		log.Printf("Models: %v", cfg.Models)
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

		fmt.Printf("\n=== Task Result ===\n%s\n", result.Output)
		fmt.Printf("\nArtifact saved: %s\n", result.ArtifactPath)

	default:
		log.Fatalf("Unknown mode: %s", *mode)
	}
}
