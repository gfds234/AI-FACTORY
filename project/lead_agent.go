package project

import (
	"ai-studio/orchestrator/llm"
	"ai-studio/orchestrator/supervisor"
	"fmt"
	"log"
	"strings"
)

// LeadAgent coordinates project workflow and makes phase decisions
type LeadAgent struct {
	llmClient *llm.Client
	model     string

	// Specialist agents
	requirementsAgent *supervisor.RequirementsAgent
	techStackAgent    *supervisor.TechStackAgent
	scopeAgent        *supervisor.ScopeAgent
	qaAgent           *supervisor.QAAgent
	testingAgent      *supervisor.TestingAgent
	docsAgent         *supervisor.DocumentationAgent
}

// PhaseResult contains the result of executing a phase
type PhaseResult struct {
	Phase             Phase
	Decision          string // PROCEED, REFINE, BLOCK
	Reasoning         string
	NextSteps         string
	AgentOutputs      map[string]*supervisor.AgentOutput
	RequiresApproval  bool
	RecommendedAction string
}

// NewLeadAgent creates a new lead agent
func NewLeadAgent(llmClient *llm.Client, model string,
	requirementsAgent *supervisor.RequirementsAgent,
	techStackAgent *supervisor.TechStackAgent,
	scopeAgent *supervisor.ScopeAgent,
	qaAgent *supervisor.QAAgent,
	testingAgent *supervisor.TestingAgent,
	docsAgent *supervisor.DocumentationAgent) *LeadAgent {

	return &LeadAgent{
		llmClient:         llmClient,
		model:             model,
		requirementsAgent: requirementsAgent,
		techStackAgent:    techStackAgent,
		scopeAgent:        scopeAgent,
		qaAgent:           qaAgent,
		testingAgent:      testingAgent,
		docsAgent:         docsAgent,
	}
}

// ExecutePhase executes a project phase and returns a decision
func (la *LeadAgent) ExecutePhase(project *Project, phase Phase) (*PhaseResult, error) {
	switch phase {
	case PhaseDiscovery:
		return la.executeDiscoveryPhase(project)
	case PhaseValidation:
		return la.executeValidationPhase(project)
	case PhasePlanning:
		return la.executePlanningPhase(project)
	case PhaseReview:
		return la.executeReviewPhase(project)
	case PhaseQA:
		return la.executeQAPhase(project)
	case PhaseDocs:
		return la.executeDocsPhase(project)
	default:
		return nil, fmt.Errorf("unsupported phase for lead agent: %s", phase)
	}
}

// executeDiscoveryPhase executes the Discovery phase
func (la *LeadAgent) executeDiscoveryPhase(project *Project) (*PhaseResult, error) {
	log.Printf("Lead Agent: Executing Discovery phase for project %s", project.Name)

	// Invoke Requirements Agent
	context := map[string]interface{}{
		"project_id": project.ID,
		"phase":      "discovery",
	}

	reqOutput, err := la.requirementsAgent.Execute("discovery", project.Description, context)
	if err != nil {
		return nil, fmt.Errorf("requirements agent failed: %w", err)
	}

	// Build prompt for Lead Agent decision
	prompt := la.buildDiscoveryPrompt(project, reqOutput)

	// Get Lead Agent decision
	response, err := la.llmClient.Generate(la.model, prompt)
	if err != nil {
		return nil, fmt.Errorf("lead agent LLM call failed: %w", err)
	}

	// Parse decision
	decision, reasoning, nextSteps := la.parseDecision(response)

	result := &PhaseResult{
		Phase:     PhaseDiscovery,
		Decision:  decision,
		Reasoning: reasoning,
		NextSteps: nextSteps,
		AgentOutputs: map[string]*supervisor.AgentOutput{
			"requirements": reqOutput,
		},
		RequiresApproval:  true,
		RecommendedAction: la.getRecommendedAction(decision),
	}

	return result, nil
}

// executeValidationPhase executes the Validation phase
func (la *LeadAgent) executeValidationPhase(project *Project) (*PhaseResult, error) {
	log.Printf("Lead Agent: Executing Validation phase for project %s", project.Name)

	context := map[string]interface{}{
		"project_id": project.ID,
		"phase":      "validation",
	}

	// Invoke TechStack and Scope agents in parallel
	techStackOutput, err := la.techStackAgent.Execute("code", project.Description, context)
	if err != nil {
		return nil, fmt.Errorf("tech stack agent failed: %w", err)
	}

	scopeOutput, err := la.scopeAgent.Execute("code", project.Description, context)
	if err != nil {
		return nil, fmt.Errorf("scope agent failed: %w", err)
	}

	// Build prompt for Lead Agent decision
	prompt := la.buildValidationPrompt(project, techStackOutput, scopeOutput)

	// Get Lead Agent decision
	response, err := la.llmClient.Generate(la.model, prompt)
	if err != nil {
		return nil, fmt.Errorf("lead agent LLM call failed: %w", err)
	}

	// Parse decision
	decision, reasoning, nextSteps := la.parseDecision(response)

	result := &PhaseResult{
		Phase:     PhaseValidation,
		Decision:  decision,
		Reasoning: reasoning,
		NextSteps: nextSteps,
		AgentOutputs: map[string]*supervisor.AgentOutput{
			"techstack": techStackOutput,
			"scope":     scopeOutput,
		},
		RequiresApproval:  true,
		RecommendedAction: la.getRecommendedAction(decision),
	}

	return result, nil
}

// executePlanningPhase executes the Planning phase
func (la *LeadAgent) executePlanningPhase(project *Project) (*PhaseResult, error) {
	log.Printf("Lead Agent: Executing Planning phase for project %s", project.Name)

	// Build prompt for plan generation
	prompt := la.buildPlanningPrompt(project)

	// Get Lead Agent plan
	response, err := la.llmClient.Generate(la.model, prompt)
	if err != nil {
		return nil, fmt.Errorf("lead agent LLM call failed: %w", err)
	}

	result := &PhaseResult{
		Phase:             PhasePlanning,
		Decision:          "PROCEED",
		Reasoning:         "Implementation roadmap created",
		NextSteps:         "Review plan and approve to proceed to code generation",
		AgentOutputs:      make(map[string]*supervisor.AgentOutput),
		RequiresApproval:  true,
		RecommendedAction: "Human should review the plan before proceeding",
	}

	// Store plan in agent outputs
	result.AgentOutputs["plan"] = &supervisor.AgentOutput{
		Output:   response,
		Status:   "passed",
		Duration: 0,
	}

	return result, nil
}

// executeReviewPhase executes the Review phase
func (la *LeadAgent) executeReviewPhase(project *Project) (*PhaseResult, error) {
	log.Printf("Lead Agent: Executing Review phase for project %s", project.Name)

	// Get the latest code artifact from project
	if len(project.ArtifactPaths) == 0 {
		return nil, fmt.Errorf("no artifacts found for review")
	}

	lastArtifact := project.ArtifactPaths[len(project.ArtifactPaths)-1]

	context := map[string]interface{}{
		"project_id":    project.ID,
		"phase":         "review",
		"artifact_path": lastArtifact,
	}

	// Invoke QA and Testing agents
	qaOutput, err := la.qaAgent.Execute("code", project.Description, context)
	if err != nil {
		log.Printf("Warning: QA agent failed: %v", err)
		qaOutput = &supervisor.AgentOutput{
			Output: "QA agent unavailable",
			Status: "warning",
		}
	}

	testingOutput, err := la.testingAgent.Execute("code", project.Description, context)
	if err != nil {
		log.Printf("Warning: Testing agent failed: %v", err)
		testingOutput = &supervisor.AgentOutput{
			Output: "Testing agent unavailable",
			Status: "warning",
		}
	}

	// Build prompt for Lead Agent decision
	prompt := la.buildReviewPrompt(project, qaOutput, testingOutput)

	// Get Lead Agent decision
	response, err := la.llmClient.Generate(la.model, prompt)
	if err != nil {
		return nil, fmt.Errorf("lead agent LLM call failed: %w", err)
	}

	// Parse decision
	decision, reasoning, nextSteps := la.parseDecision(response)

	result := &PhaseResult{
		Phase:     PhaseReview,
		Decision:  decision,
		Reasoning: reasoning,
		NextSteps: nextSteps,
		AgentOutputs: map[string]*supervisor.AgentOutput{
			"qa":      qaOutput,
			"testing": testingOutput,
		},
		RequiresApproval:  decision != "PROCEED",
		RecommendedAction: la.getRecommendedAction(decision),
	}

	return result, nil
}

// executeQAPhase executes the QA phase
func (la *LeadAgent) executeQAPhase(project *Project) (*PhaseResult, error) {
	log.Printf("Lead Agent: Executing QA phase for project %s", project.Name)

	// QA phase validation is handled by CompletionValidator
	// Lead Agent just generates a summary
	result := &PhaseResult{
		Phase:             PhaseQA,
		Decision:          "PROCEED",
		Reasoning:         "QA validation checks will be performed",
		NextSteps:         "Review hand-off ready criteria",
		AgentOutputs:      make(map[string]*supervisor.AgentOutput),
		RequiresApproval:  true,
		RecommendedAction: "Human should verify hand-off ready criteria",
	}

	return result, nil
}

// executeDocsPhase executes the Docs phase
func (la *LeadAgent) executeDocsPhase(project *Project) (*PhaseResult, error) {
	log.Printf("Lead Agent: Executing Docs phase for project %s", project.Name)

	context := map[string]interface{}{
		"project_id": project.ID,
		"phase":      "docs",
	}

	// Invoke Documentation agent
	docsOutput, err := la.docsAgent.Execute("code", project.Description, context)
	if err != nil {
		log.Printf("Warning: Documentation agent failed: %v", err)
		docsOutput = &supervisor.AgentOutput{
			Output: "Documentation agent unavailable",
			Status: "warning",
		}
	}

	result := &PhaseResult{
		Phase:     PhaseDocs,
		Decision:  "PROCEED",
		Reasoning: "Documentation generated",
		NextSteps: "Proceed to project completion",
		AgentOutputs: map[string]*supervisor.AgentOutput{
			"documentation": docsOutput,
		},
		RequiresApproval:  false,
		RecommendedAction: "Automated transition to Complete phase",
	}

	return result, nil
}

// buildDiscoveryPrompt builds the prompt for Discovery phase decision
func (la *LeadAgent) buildDiscoveryPrompt(project *Project, reqOutput *supervisor.AgentOutput) string {
	return fmt.Sprintf(`%s

Analyze this new project request in the Discovery phase.

Project: %s
Description: %s

Requirements Agent Output:
%s

Decide: PROCEED, REFINE, or BLOCK

Criteria:
- PROCEED: Requirements are clear and complete, score >= 7/10
- REFINE: Requirements need clarification, score 4-6/10
- BLOCK: Requirements incomplete or unclear, score < 4/10

Respond in this format:
DECISION: [PROCEED|REFINE|BLOCK]
REASONING: [2-3 sentences explaining decision]
NEXT_STEPS: [What needs to happen next]`,
		la.getBaseSystemPrompt(),
		project.Name,
		project.Description,
		reqOutput.Output,
	)
}

// buildValidationPrompt builds the prompt for Validation phase decision
func (la *LeadAgent) buildValidationPrompt(project *Project, techStack, scope *supervisor.AgentOutput) string {
	return fmt.Sprintf(`%s

Validate project feasibility in the Validation phase.

Project: %s

TechStack Agent Output:
%s

Scope Agent Output:
%s

Decide: PROCEED, REFINE, or BLOCK

Criteria:
- PROCEED: Both agents approved, no major concerns
- REFINE: Warnings present but addressable
- BLOCK: Tech stack rejected OR scope too broad

Respond in this format:
DECISION: [PROCEED|REFINE|BLOCK]
REASONING: [Why this decision]
NEXT_STEPS: [Required actions]`,
		la.getBaseSystemPrompt(),
		project.Name,
		techStack.Output,
		scope.Output,
	)
}

// buildPlanningPrompt builds the prompt for Planning phase
func (la *LeadAgent) buildPlanningPrompt(project *Project) string {
	return fmt.Sprintf(`%s

Create an implementation roadmap for this project in the Planning phase.

Project: %s
Description: %s

Generate a concise implementation plan that includes:
1. Key milestones (3-5 major deliverables)
2. Technical approach (architecture, patterns)
3. Estimated complexity (simple, medium, complex)
4. Potential risks or challenges
5. Success criteria

Keep it practical and shipping-focused. Avoid over-engineering.

Format as markdown with clear sections.`,
		la.getBaseSystemPrompt(),
		project.Name,
		project.Description,
	)
}

// buildReviewPrompt builds the prompt for Review phase decision
func (la *LeadAgent) buildReviewPrompt(project *Project, qa, testing *supervisor.AgentOutput) string {
	return fmt.Sprintf(`%s

Review code quality for this project in the Review phase.

Project: %s

QA Agent Output:
%s

Testing Agent Output:
%s

Decide: PROCEED, REFINE, or BLOCK

Criteria:
- PROCEED: QA score >= 7/10, no critical bugs, tests exist
- REFINE: QA score 5-6/10, some issues present
- BLOCK: QA score < 5/10, critical bugs detected

Respond in this format:
DECISION: [PROCEED|REFINE|BLOCK]
REASONING: [Assessment of code quality]
NEXT_STEPS: [What to do next]`,
		la.getBaseSystemPrompt(),
		project.Name,
		qa.Output,
		testing.Output,
	)
}

// getBaseSystemPrompt returns the base system prompt for Lead Agent
func (la *LeadAgent) getBaseSystemPrompt() string {
	return `You are the Lead Agent for an AI software factory. Your role is Product Producer + Tech Lead.

Core Principles:
1. SHIPPING MATTERS - Prefer working code over perfection
2. SCOPE CONTROL - Guard against feature creep aggressively
3. QUALITY GATES - Block on critical issues, warn on concerns
4. DELEGATION - Use specialist agents, don't do their jobs
5. CONSERVATIVE - Proven patterns over experimental approaches
6. TRANSPARENCY - Explain decisions clearly for human approval

Decision Framework:
- PROCEED: All criteria met, safe to continue
- REFINE: Concerns present, needs user clarification
- BLOCK: Critical issues, cannot proceed safely

Always:
- Explain your reasoning
- Cite specific agent outputs
- Provide actionable next steps
- Defer to human judgment on ambiguous cases`
}

// parseDecision parses the LLM response for DECISION, REASONING, NEXT_STEPS
func (la *LeadAgent) parseDecision(response string) (decision, reasoning, nextSteps string) {
	lines := strings.Split(response, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "DECISION:") {
			decision = strings.TrimSpace(strings.TrimPrefix(line, "DECISION:"))
			decision = strings.ToUpper(decision)
			// Extract just the decision word
			if strings.Contains(decision, "PROCEED") {
				decision = "PROCEED"
			} else if strings.Contains(decision, "REFINE") {
				decision = "REFINE"
			} else if strings.Contains(decision, "BLOCK") {
				decision = "BLOCK"
			}
		}

		if strings.HasPrefix(line, "REASONING:") {
			reasoning = strings.TrimSpace(strings.TrimPrefix(line, "REASONING:"))
		}

		if strings.HasPrefix(line, "NEXT_STEPS:") {
			nextSteps = strings.TrimSpace(strings.TrimPrefix(line, "NEXT_STEPS:"))
		}
	}

	// Default values if parsing failed
	if decision == "" {
		decision = "REFINE"
		reasoning = "Unable to parse decision from response"
		nextSteps = "Manual review required"
	}

	return decision, reasoning, nextSteps
}

// getRecommendedAction returns a human-readable recommended action based on decision
func (la *LeadAgent) getRecommendedAction(decision string) string {
	switch decision {
	case "PROCEED":
		return "Approve and proceed to next phase"
	case "REFINE":
		return "Review concerns and iterate before proceeding"
	case "BLOCK":
		return "Do not proceed - critical issues must be resolved"
	default:
		return "Manual review required"
	}
}

// GenerateProjectSummary generates a final project summary
func (la *LeadAgent) GenerateProjectSummary(project *Project) (string, error) {
	prompt := fmt.Sprintf(`Generate a comprehensive project completion summary for:

Project: %s
Description: %s
Status: %s
Phases Completed: %d

Include:
1. Project overview
2. Phase execution summary
3. Artifacts generated
4. Next steps for deployment/usage
5. Recommendations

Format as markdown.`,
		project.Name,
		project.Description,
		project.Status,
		len(project.Phases),
	)

	summary, err := la.llmClient.Generate(la.model, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate summary: %w", err)
	}

	return summary, nil
}
