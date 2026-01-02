package project

import (
	"ai-studio/orchestrator/llm"
	"ai-studio/orchestrator/supervisor"
	"ai-studio/orchestrator/task"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// ProjectOrchestrator manages project lifecycle and phase execution
type ProjectOrchestrator struct {
	supervisedMgr       *supervisor.SupervisedTaskManager
	projectMgr          *ProjectManager
	leadAgent           *LeadAgent
	completionValidator *CompletionValidator
}

// NewProjectOrchestrator creates a new project orchestrator
func NewProjectOrchestrator(
	supervisedMgr *supervisor.SupervisedTaskManager,
	projectsDir string,
	artifactsDir string,
	llmClient *llm.Client,
	requirementsAgent *supervisor.RequirementsAgent,
	techStackAgent *supervisor.TechStackAgent,
	scopeAgent *supervisor.ScopeAgent,
	qaAgent *supervisor.QAAgent,
	testingAgent *supervisor.TestingAgent,
	docsAgent *supervisor.DocumentationAgent,
) (*ProjectOrchestrator, error) {

	// Create ProjectManager
	projectMgr, err := NewProjectManager(projectsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create project manager: %w", err)
	}

	// Create LeadAgent
	leadAgent := NewLeadAgent(
		llmClient,
		"llama3:8b", // Ollama model for Lead Agent
		requirementsAgent,
		techStackAgent,
		scopeAgent,
		qaAgent,
		testingAgent,
		docsAgent,
	)

	// Create CompletionValidator
	completionValidator := NewCompletionValidator(artifactsDir)

	return &ProjectOrchestrator{
		supervisedMgr:       supervisedMgr,
		projectMgr:          projectMgr,
		leadAgent:           leadAgent,
		completionValidator: completionValidator,
	}, nil
}

// CreateProject creates a new project
func (po *ProjectOrchestrator) CreateProject(name, description string) (*Project, error) {
	log.Printf("ProjectOrchestrator: Creating project '%s'", name)

	project, err := po.projectMgr.CreateProject(name, description)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	log.Printf("ProjectOrchestrator: Project created with ID %s", project.ID)

	return project, nil
}

// GetProject retrieves a project by ID
func (po *ProjectOrchestrator) GetProject(id string) (*Project, error) {
	return po.projectMgr.GetProject(id)
}

// ListProjects returns all projects
func (po *ProjectOrchestrator) ListProjects() []*Project {
	return po.projectMgr.ListProjects()
}

// ExecuteProjectPhase executes a specific phase for a project
func (po *ProjectOrchestrator) ExecuteProjectPhase(projectID string, phase Phase) (*PhaseResult, error) {
	project, err := po.projectMgr.GetProject(projectID)
	if err != nil {
		return nil, err
	}

	log.Printf("ProjectOrchestrator: Executing %s phase for project %s", phase, project.Name)

	// Update phase status to in_progress
	if err := po.projectMgr.UpdateProjectPhase(project, phase, PhaseStatusInProgress); err != nil {
		return nil, fmt.Errorf("failed to update phase status: %w", err)
	}

	// Execute phase based on type
	var phaseResult *PhaseResult

	switch phase {
	case PhaseDiscovery, PhaseValidation, PhasePlanning, PhaseReview, PhaseQA, PhaseDocs:
		// Lead Agent handles these phases
		phaseResult, err = po.leadAgent.ExecutePhase(project, phase)
		if err != nil {
			// Revert phase status on error
			po.projectMgr.UpdateProjectPhase(project, phase, PhaseStatusPending)
			return nil, fmt.Errorf("lead agent execution failed: %w", err)
		}

	case PhaseCodeGen:
		// Delegate to SupervisedTaskManager for code generation
		phaseResult, err = po.executeCodeGenPhase(project)
		if err != nil {
			// Revert phase status on error
			po.projectMgr.UpdateProjectPhase(project, phase, PhaseStatusPending)
			return nil, fmt.Errorf("code generation failed: %w", err)
		}

	case PhaseComplete:
		// Finalize project
		phaseResult, err = po.executeCompletePhase(project)
		if err != nil {
			// Revert phase status on error
			po.projectMgr.UpdateProjectPhase(project, phase, PhaseStatusPending)
			return nil, fmt.Errorf("project completion failed: %w", err)
		}

	default:
		return nil, fmt.Errorf("unsupported phase: %s", phase)
	}

	// Store phase result in project
	err = po.storePhaseResult(project, phase, phaseResult)
	if err != nil {
		// Revert phase status on storage error
		po.projectMgr.UpdateProjectPhase(project, phase, PhaseStatusPending)
		return nil, fmt.Errorf("failed to store phase result: %w", err)
	}

	log.Printf("ProjectOrchestrator: %s phase completed with decision: %s", phase, phaseResult.Decision)

	return phaseResult, nil
}

// executeCodeGenPhase executes the code generation phase
func (po *ProjectOrchestrator) executeCodeGenPhase(project *Project) (*PhaseResult, error) {
	log.Printf("ProjectOrchestrator: Delegating to SupervisedTaskManager for code generation")

	// Execute code generation via SupervisedTaskManager
	result, err := po.supervisedMgr.ExecuteTask("code", project.Description)
	if err != nil {
		return nil, fmt.Errorf("supervised task execution failed: %w", err)
	}

	// Type assert result
	supervisedResult, ok := result.(*supervisor.SupervisedResult)
	if !ok {
		return nil, fmt.Errorf("unexpected result type from supervised manager")
	}

	// Store task execution in project
	taskExec := TaskExecution{
		TaskID:          uuid.New().String(),
		Phase:           PhaseCodeGen,
		TaskType:        "code",
		Input:           project.Description,
		Output:          supervisedResult.Result.Output,
		ArtifactPath:    supervisedResult.Result.ArtifactPath,
		ComplexityScore: supervisedResult.ComplexityScore,
		ExecutionRoute:  supervisedResult.ExecutionRoute,
		CreatedAt:       time.Now(),
	}

	if err := po.projectMgr.AddTaskExecution(project, taskExec); err != nil {
		return nil, fmt.Errorf("failed to add task execution: %w", err)
	}

	// Add artifact path to project
	if supervisedResult.Result.ArtifactPath != "" {
		if err := po.projectMgr.AddArtifactPath(project, supervisedResult.Result.ArtifactPath); err != nil {
			return nil, fmt.Errorf("failed to add artifact path: %w", err)
		}
	}

	phaseResult := &PhaseResult{
		Phase:             PhaseCodeGen,
		Decision:          "PROCEED",
		Reasoning:         fmt.Sprintf("Code generated via %s (complexity: %d)", supervisedResult.ExecutionRoute, supervisedResult.ComplexityScore),
		NextSteps:         "Proceed to Review phase",
		AgentOutputs:      make(map[string]*supervisor.AgentOutput),
		RequiresApproval:  false, // CodeGen is automated
		RecommendedAction: "Automated transition to Review phase",
	}

	return phaseResult, nil
}

// executeCompletePhase finalizes the project
func (po *ProjectOrchestrator) executeCompletePhase(project *Project) (*PhaseResult, error) {
	log.Printf("ProjectOrchestrator: Finalizing project %s", project.Name)

	// Validate hand-off ready criteria
	metrics, err := po.completionValidator.ValidateHandoffReady(project)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Generate project summary
	summary, err := po.leadAgent.GenerateProjectSummary(project)
	if err != nil {
		log.Printf("Warning: Failed to generate summary: %v", err)
		summary = "Summary generation unavailable"
	}

	// Update project status
	now := time.Now()
	project.Status = ProjectStatusComplete
	project.CompletedAt = &now

	if err := po.projectMgr.SaveProject(project); err != nil {
		return nil, fmt.Errorf("failed to save project: %w", err)
	}

	phaseResult := &PhaseResult{
		Phase:     PhaseComplete,
		Decision:  "PROCEED",
		Reasoning: fmt.Sprintf("Project complete - Completion: %.1f%%, Build: %t, Tests: %t, README: %t", metrics.CompletionPct, metrics.HasRunnableBuild, metrics.HasTests, metrics.HasReadme),
		NextSteps: summary,
		AgentOutputs: map[string]*supervisor.AgentOutput{
			"completion_metrics": {
				Output: fmt.Sprintf("Completion: %.1f%%\nBuild: %t\nTests: %t\nREADME: %t\nBlocking Issues: %v",
					metrics.CompletionPct, metrics.HasRunnableBuild, metrics.HasTests, metrics.HasReadme, metrics.BlockingIssues),
				Status:   "passed",
				Duration: 0,
			},
		},
		RequiresApproval:  false,
		RecommendedAction: "Project complete - ready for hand-off",
	}

	return phaseResult, nil
}

// TransitionPhase transitions a project to a new phase
func (po *ProjectOrchestrator) TransitionPhase(projectID string, toPhase Phase, humanApproval bool) error {
	project, err := po.projectMgr.GetProject(projectID)
	if err != nil {
		return err
	}

	// Validate transition is allowed
	if !project.CurrentPhase.CanTransition(toPhase) {
		return fmt.Errorf("invalid phase transition: %s -> %s", project.CurrentPhase, toPhase)
	}

	log.Printf("ProjectOrchestrator: Transitioning project %s from %s to %s (approved: %t)",
		project.Name, project.CurrentPhase, toPhase, humanApproval)

	// Update current phase execution with approval status
	for i := range project.Phases {
		if project.Phases[i].Phase == project.CurrentPhase {
			project.Phases[i].HumanApproval = humanApproval
			if humanApproval {
				now := time.Now()
				project.Phases[i].CompletedAt = &now
				project.Phases[i].Status = PhaseStatusComplete
			}
			break
		}
	}

	// Create new phase execution record
	newPhaseExec := PhaseExecution{
		Phase:        toPhase,
		Status:       PhaseStatusPending,
		StartedAt:    time.Now(),
		AgentOutputs: make(map[string]string),
	}

	project.Phases = append(project.Phases, newPhaseExec)
	project.CurrentPhase = toPhase

	return po.projectMgr.SaveProject(project)
}

// ApprovePhase approves the current phase and transitions to next
func (po *ProjectOrchestrator) ApprovePhase(projectID string) error {
	project, err := po.projectMgr.GetProject(projectID)
	if err != nil {
		return err
	}

	// Determine next phase based on current phase
	nextPhase, err := po.getNextPhase(project.CurrentPhase)
	if err != nil {
		return err
	}

	return po.TransitionPhase(projectID, nextPhase, true)
}

// RejectPhase rejects the current phase and blocks progress
func (po *ProjectOrchestrator) RejectPhase(projectID string, reason string) error {
	project, err := po.projectMgr.GetProject(projectID)
	if err != nil {
		return err
	}

	// Mark current phase as blocked
	for i := range project.Phases {
		if project.Phases[i].Phase == project.CurrentPhase {
			project.Phases[i].Status = PhaseStatusBlocked
			project.Phases[i].Notes = reason
			break
		}
	}

	project.Status = ProjectStatusBlocked

	return po.projectMgr.SaveProject(project)
}

// RevertPhase reverts the project to a previous phase
func (po *ProjectOrchestrator) RevertPhase(projectID string, targetPhase Phase, reason string) error {
	project, err := po.projectMgr.GetProject(projectID)
	if err != nil {
		return err
	}

	// Validate target phase exists in history and was completed
	targetPhaseExists := false
	for _, phaseExec := range project.Phases {
		if phaseExec.Phase == targetPhase && phaseExec.Status == PhaseStatusComplete {
			targetPhaseExists = true
			break
		}
	}

	if !targetPhaseExists {
		return fmt.Errorf("cannot revert to %s: phase not found or not completed", targetPhase)
	}

	// Validate backward transition is allowed
	if !project.CurrentPhase.CanGoBackTo(targetPhase) {
		return fmt.Errorf("cannot revert from %s to %s", project.CurrentPhase, targetPhase)
	}

	log.Printf("ProjectOrchestrator: Reverting project %s from %s to %s (reason: %s)",
		project.Name, project.CurrentPhase, targetPhase, reason)

	// Update current phase pointer (preserves all data)
	project.CurrentPhase = targetPhase

	// Mark phases after target as reverted (keeps data for audit)
	for i := range project.Phases {
		if project.Phases[i].Phase == targetPhase {
			// Keep this phase as-is
			continue
		}

		// Check if this phase comes after target in workflow
		if isAfter(project.Phases[i].Phase, targetPhase) {
			// Keep status but clear approval and add revert note
			project.Phases[i].HumanApproval = false
			if project.Phases[i].Notes == "" {
				project.Phases[i].Notes = fmt.Sprintf("Reverted on %s: %s",
					time.Now().Format("2006-01-02 15:04:05"), reason)
			} else {
				project.Phases[i].Notes = fmt.Sprintf("%s\n[Reverted on %s: %s]",
					project.Phases[i].Notes,
					time.Now().Format("2006-01-02 15:04:05"), reason)
			}
		}
	}

	// Restore project to active status if it was blocked
	if project.Status == ProjectStatusBlocked {
		project.Status = ProjectStatusActive
	}

	project.UpdatedAt = time.Now()

	return po.projectMgr.SaveProject(project)
}

// isAfter checks if phaseA comes after phaseB in the workflow
func isAfter(phaseA Phase, phaseB Phase) bool {
	phaseOrder := []Phase{
		PhaseDiscovery, PhaseValidation, PhasePlanning,
		PhaseCodeGen, PhaseReview, PhaseQA, PhaseDocs, PhaseComplete,
	}

	idxA, idxB := -1, -1
	for i, p := range phaseOrder {
		if p == phaseA {
			idxA = i
		}
		if p == phaseB {
			idxB = i
		}
	}

	return idxA > idxB && idxA >= 0 && idxB >= 0
}

// GetCompletionMetrics gets hand-off ready metrics for a project
func (po *ProjectOrchestrator) GetCompletionMetrics(projectID string) (*CompletionMetrics, error) {
	project, err := po.projectMgr.GetProject(projectID)
	if err != nil {
		return nil, err
	}

	return po.completionValidator.ValidateHandoffReady(project)
}

// storePhaseResult stores phase result in project
func (po *ProjectOrchestrator) storePhaseResult(project *Project, phase Phase, result *PhaseResult) error {
	// Find current phase execution
	for i := range project.Phases {
		if project.Phases[i].Phase == phase && project.Phases[i].Status == PhaseStatusInProgress {
			project.Phases[i].LeadAgentDecision = result.Decision
			project.Phases[i].LeadAgentInput = project.Description

			// Mark phase as complete (execution finished successfully)
			now := time.Now()
			project.Phases[i].CompletedAt = &now
			project.Phases[i].Status = PhaseStatusComplete

			// Store agent outputs
			for agentName, output := range result.AgentOutputs {
				project.Phases[i].AgentOutputs[agentName] = output.Output
			}

			break
		}
	}

	return po.projectMgr.SaveProject(project)
}

// getNextPhase determines the next phase based on current phase
func (po *ProjectOrchestrator) getNextPhase(current Phase) (Phase, error) {
	validTransitions, exists := PhaseTransitions[current]
	if !exists || len(validTransitions) == 0 {
		return "", fmt.Errorf("no valid transitions from phase: %s", current)
	}

	// Return first valid transition (for linear flow)
	return validTransitions[0], nil
}

// Implement task.Manager interface for backward compatibility

// ExecuteTask executes a task (backward compatibility)
func (po *ProjectOrchestrator) ExecuteTask(taskType, input string) (interface{}, error) {
	// Delegate to SupervisedTaskManager for backward compatibility
	return po.supervisedMgr.ExecuteTask(taskType, input)
}

// Ping checks system health
func (po *ProjectOrchestrator) Ping() error {
	return po.supervisedMgr.Ping()
}

// GetHistory gets task history
func (po *ProjectOrchestrator) GetHistory(taskType string) []task.Result {
	return po.supervisedMgr.GetHistory(taskType)
}

// GetClient gets LLM client
func (po *ProjectOrchestrator) GetClient() *llm.Client {
	return po.supervisedMgr.GetClient()
}

// DeleteProject deletes a project
func (po *ProjectOrchestrator) DeleteProject(id string) error {
	return po.projectMgr.DeleteProject(id)
}
