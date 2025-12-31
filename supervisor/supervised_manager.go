package supervisor

import (
	"ai-studio/orchestrator/config"
	"ai-studio/orchestrator/llm"
	"ai-studio/orchestrator/task"
	"fmt"
	"log"
	"time"
)

// SupervisedTaskManager wraps task.Manager with agent supervision
type SupervisedTaskManager struct {
	baseManager      *task.Manager
	cfg              *SupervisorConfig
	complexityScorer *ComplexityScorer
	claudeCodeClient *ClaudeCodeClient

	// Agents
	requirementsAgent *RequirementsAgent
	techStackAgent    *TechStackAgent
	scopeAgent        *ScopeAgent
	qaAgent           *QAAgent
	testingAgent      *TestingAgent
	docsAgent         *DocumentationAgent
}

// NewSupervisedTaskManager creates a supervised task manager
func NewSupervisedTaskManager(baseMgr *task.Manager, cfg *config.Config, supervisorCfg *SupervisorConfig) *SupervisedTaskManager {
	client := baseMgr.GetClient()

	stm := &SupervisedTaskManager{
		baseManager:      baseMgr,
		cfg:              supervisorCfg,
		complexityScorer: NewComplexityScorer(supervisorCfg),
	}

	// Initialize Claude Code client if endpoint configured
	if supervisorCfg.ClaudeCodeEndpoint != "" {
		stm.claudeCodeClient = NewClaudeCodeClient(supervisorCfg.ClaudeCodeEndpoint)
	}

	// Initialize agents
	if supervisorCfg.Agents.Requirements.Enabled {
		stm.requirementsAgent = NewRequirementsAgent(client, supervisorCfg.Agents.Requirements.Model)
	}
	if supervisorCfg.Agents.QA.Enabled {
		stm.qaAgent = NewQAAgent(client, supervisorCfg.Agents.QA.Model)
	}
	if supervisorCfg.Agents.Testing.Enabled {
		stm.testingAgent = NewTestingAgent(client, supervisorCfg.Agents.Testing.Model)
	}
	if supervisorCfg.Agents.Documentation.Enabled {
		stm.docsAgent = NewDocumentationAgent(client, supervisorCfg.Agents.Documentation.Model)
	}

	// Tech stack and scope agents (created on-demand for code tasks)
	stm.techStackAgent = NewTechStackAgent(client, "llama3:8b")
	stm.scopeAgent = NewScopeAgent(client, "mistral:7b-instruct-v0.2-q4_K_M")

	return stm
}

// ExecuteTask runs the full supervised execution pipeline
func (stm *SupervisedTaskManager) ExecuteTask(taskType, input string) (interface{}, error) {
	startTime := time.Now()

	result := &SupervisedResult{
		Result:         &task.Result{TaskType: taskType, Input: input, Timestamp: startTime},
		AgentDurations: make(map[string]float64),
	}

	// Phase 1: Pre-execution quality gates (only if enabled)
	if stm.cfg.Enabled {
		if err := stm.runQualityGates(taskType, input, result); err != nil {
			result.Error = err.Error()
			result.TotalDuration = time.Since(startTime).Seconds()
			// Return full SupervisedResult even on error
			return result, err
		}
	}

	// Phase 2: Complexity scoring and routing
	complexity := stm.complexityScorer.Score(taskType, input)
	result.ComplexityScore = complexity.Score
	result.ExecutionRoute = complexity.RecommendedRoute

	log.Printf("Complexity score: %d, route: %s", complexity.Score, complexity.RecommendedRoute)

	// Phase 3: Execute task (Ollama or Claude Code)
	var baseResult *task.Result
	var err error

	if complexity.RecommendedRoute == "claude_code" && stm.claudeCodeClient != nil {
		baseResult, err = stm.executeWithClaudeCode(taskType, input)
	} else {
		execResult, execErr := stm.baseManager.ExecuteTask(taskType, input)
		err = execErr
		if execErr == nil {
			// Type assert the result back to *task.Result
			var ok bool
			baseResult, ok = execResult.(*task.Result)
			if !ok {
				err = fmt.Errorf("unexpected result type from base manager")
			}
		}
	}

	if err != nil {
		result.Error = err.Error()
		result.TotalDuration = time.Since(startTime).Seconds()
		return result, err
	}

	// Merge base result
	result.Result = baseResult

	// Phase 4: Post-execution agents (only if enabled)
	if stm.cfg.Enabled {
		stm.runPostExecutionAgents(taskType, input, baseResult.Output, result)
	}

	result.TotalDuration = time.Since(startTime).Seconds()

	// Return full SupervisedResult with all metadata
	return result, nil
}

// runQualityGates executes pre-execution quality gates
func (stm *SupervisedTaskManager) runQualityGates(taskType, input string, result *SupervisedResult) error {
	context := make(map[string]interface{})

	// Gate 1: Requirements check
	if stm.cfg.QualityGates.RequirementsCheck && stm.requirementsAgent != nil {
		log.Printf("Running requirements check...")
		reqOutput, err := stm.requirementsAgent.Execute(taskType, input, context)
		if err != nil {
			return fmt.Errorf("requirements check failed: %w", err)
		}
		result.RequirementsAnalysis = reqOutput
		result.AgentDurations["requirements"] = reqOutput.Duration

		if reqOutput.Status == "failed" {
			return fmt.Errorf("requirements incomplete - cannot proceed")
		}
	}

	// Gate 2: Tech stack approval (code tasks only)
	if stm.cfg.QualityGates.TechStackApproval && taskType == "code" {
		log.Printf("Running tech stack approval...")
		tsOutput, err := stm.techStackAgent.Execute(taskType, input, context)
		if err != nil {
			return fmt.Errorf("tech stack approval failed: %w", err)
		}
		result.TechStackApproval = tsOutput
		result.AgentDurations["techstack"] = tsOutput.Duration

		if tsOutput.Status == "failed" {
			return fmt.Errorf("tech stack rejected - cannot proceed")
		}
	}

	// Gate 3: Scope validation
	if stm.cfg.QualityGates.ScopeValidation {
		log.Printf("Running scope validation...")
		scopeOutput, err := stm.scopeAgent.Execute(taskType, input, context)
		if err != nil {
			return fmt.Errorf("scope validation failed: %w", err)
		}
		result.ScopeValidation = scopeOutput
		result.AgentDurations["scope"] = scopeOutput.Duration

		if scopeOutput.Status == "failed" {
			return fmt.Errorf("scope too broad - break into smaller tasks")
		}
	}

	return nil
}

// runPostExecutionAgents executes QA, testing, and documentation agents
func (stm *SupervisedTaskManager) runPostExecutionAgents(taskType, input, output string, result *SupervisedResult) {
	context := map[string]interface{}{
		"output": output,
	}

	// QA Review
	if stm.qaAgent != nil {
		log.Printf("Running QA review...")
		qaOutput, err := stm.qaAgent.Execute(taskType, input, context)
		if err != nil {
			log.Printf("QA agent failed: %v", err)
		} else {
			result.QAReview = qaOutput
			result.AgentDurations["qa"] = qaOutput.Duration
		}
	}

	// Testing
	if stm.testingAgent != nil && taskType == "code" {
		log.Printf("Generating test plan...")
		testOutput, err := stm.testingAgent.Execute(taskType, input, context)
		if err != nil {
			log.Printf("Testing agent failed: %v", err)
		} else {
			result.TestPlan = testOutput
			result.AgentDurations["testing"] = testOutput.Duration
		}
	}

	// Documentation
	if stm.docsAgent != nil {
		log.Printf("Generating documentation...")
		docsOutput, err := stm.docsAgent.Execute(taskType, input, context)
		if err != nil {
			log.Printf("Documentation agent failed: %v", err)
		} else {
			result.Documentation = docsOutput
			result.AgentDurations["documentation"] = docsOutput.Duration
		}
	}
}

// executeWithClaudeCode routes execution to Claude Code
func (stm *SupervisedTaskManager) executeWithClaudeCode(taskType, input string) (*task.Result, error) {
	start := time.Now()

	output, err := stm.claudeCodeClient.Generate(taskType, input)
	if err != nil {
		return nil, fmt.Errorf("claude code execution failed: %w", err)
	}

	return &task.Result{
		TaskType:  taskType,
		Input:     input,
		Output:    output,
		Model:     "claude-code",
		Duration:  time.Since(start).Seconds(),
		Timestamp: start,
	}, nil
}

// Ping checks all backends
func (stm *SupervisedTaskManager) Ping() error {
	// Ping base manager (Ollama)
	if err := stm.baseManager.Ping(); err != nil {
		return err
	}

	// Ping Claude Code if configured
	if stm.claudeCodeClient != nil {
		if err := stm.claudeCodeClient.Ping(); err != nil {
			log.Printf("Warning: Claude Code not accessible: %v", err)
			// Don't fail - we can still work with Ollama
		}
	}

	return nil
}

// GetHistory delegates to base manager
func (stm *SupervisedTaskManager) GetHistory(taskType string) []task.Result {
	return stm.baseManager.GetHistory(taskType)
}

// GetClient returns LLM client for chat
func (stm *SupervisedTaskManager) GetClient() *llm.Client {
	return stm.baseManager.GetClient()
}
