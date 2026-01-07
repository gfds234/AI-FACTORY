package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"ai-studio/orchestrator/api"
	"ai-studio/orchestrator/project"
	"ai-studio/orchestrator/supervisor"
	"ai-studio/orchestrator/task"
)

func main() {
	// CLI flags
	mode := flag.String("mode", "server", "Run mode: server or cli")
	taskType := flag.String("task", "", "Task type for CLI mode: validate or review")
	input := flag.String("input", "", "Input file path for CLI mode")

	// Use Railway PORT if available
	defaultPort := 8080
	if envPort := os.Getenv("PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			defaultPort = p
		}
	}

	port := flag.Int("port", defaultPort, "Server port")
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
	var supervisedMgr *supervisor.SupervisedTaskManager

	if supervisorConfig.Enabled {
		supervisedMgr = supervisor.NewSupervisedTaskManager(baseMgr, baseConfig, supervisorConfig)
		taskMgr = supervisedMgr
		log.Printf("✓ Supervisor enabled (complexity threshold: %d)", supervisorConfig.ComplexityThreshold)
	} else {
		taskMgr = baseMgr
		log.Printf("Supervisor disabled - using standard task manager")
	}

	// Wrap with ProjectOrchestrator if enabled
	if baseConfig.ProjectOrchestrator.Enabled && supervisedMgr != nil {
		orchestrator, err := project.NewProjectOrchestrator(
			supervisedMgr,
			baseConfig.ProjectOrchestrator.ProjectsDir,
			baseConfig.ArtifactsDir,
			baseMgr.GetClient(),
			supervisedMgr.GetRequirementsAgent(),
			supervisedMgr.GetTechStackAgent(),
			supervisedMgr.GetScopeAgent(),
			supervisedMgr.GetQAAgent(),
			supervisedMgr.GetTestingAgent(),
			supervisedMgr.GetDocsAgent(),
		)
		if err != nil {
			log.Fatalf("Failed to create ProjectOrchestrator: %v", err)
		}
		taskMgr = orchestrator
		log.Printf("✓ ProjectOrchestrator enabled (projects dir: %s)", baseConfig.ProjectOrchestrator.ProjectsDir)
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
