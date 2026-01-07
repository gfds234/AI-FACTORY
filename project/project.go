package project

import (
	"time"
)

// Project represents a project in the AI Factory workflow
type Project struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	CurrentPhase  Phase             `json:"current_phase"`
	Phases        []PhaseExecution  `json:"phases"`
	Tasks         []TaskExecution   `json:"tasks"`
	ArtifactPaths     []string           `json:"artifact_paths"`
	Metadata          ProjectMetadata    `json:"metadata"`
	ValidationResults *ValidationResults `json:"validation_results,omitempty"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
	CompletedAt       *time.Time         `json:"completed_at,omitempty"`
	Status            ProjectStatus      `json:"status"`
}

// Phase represents a project phase
type Phase string

const (
	PhaseDiscovery   Phase = "discovery"
	PhaseValidation  Phase = "validation"
	PhasePlanning    Phase = "planning"
	PhaseCodeGen     Phase = "codegen"
	PhaseReview      Phase = "review"
	PhaseQA          Phase = "qa"
	PhaseDocs        Phase = "docs"
	PhaseComplete    Phase = "complete"
)

// PhaseExecution tracks execution of a single phase
type PhaseExecution struct {
	Phase             Phase             `json:"phase"`
	Status            PhaseStatus       `json:"status"`
	StartedAt         time.Time         `json:"started_at"`
	CompletedAt       *time.Time        `json:"completed_at,omitempty"`
	LeadAgentInput    string            `json:"lead_agent_input"`
	LeadAgentDecision string            `json:"lead_agent_decision"` // PROCEED, REFINE, BLOCK
	AgentOutputs      map[string]string `json:"agent_outputs"`
	HumanApproval     bool              `json:"human_approval"`
	Notes             string            `json:"notes"`
}

// PhaseStatus represents the status of a phase
type PhaseStatus string

const (
	PhaseStatusPending    PhaseStatus = "pending"
	PhaseStatusInProgress PhaseStatus = "in_progress"
	PhaseStatusBlocked    PhaseStatus = "blocked"
	PhaseStatusComplete   PhaseStatus = "complete"
)

// TaskExecution tracks a single task execution within a project
type TaskExecution struct {
	TaskID          string                 `json:"task_id"`
	Phase           Phase                  `json:"phase"`
	TaskType        string                 `json:"task_type"` // code, validate, review
	Input           string                 `json:"input"`
	Output          string                 `json:"output"`
	ArtifactPath    string                 `json:"artifact_path"`
	ComplexityScore int                    `json:"complexity_score"`
	ExecutionRoute  string                 `json:"execution_route"` // ollama or claude_code
	AgentMetadata   map[string]interface{} `json:"agent_metadata"`
	CreatedAt       time.Time              `json:"created_at"`
}

// ProjectMetadata holds additional project information
type ProjectMetadata struct {
	ProjectType       string   `json:"project_type"`        // game, web_app, mobile_app, saas
	TechStack         []string `json:"tech_stack"`
	TargetPlatform    string   `json:"target_platform"`
	EstimatedDuration string   `json:"estimated_duration"`
	ComplexityRating  int      `json:"complexity_rating"`
}

// ProjectStatus represents the overall status of a project
type ProjectStatus string

const (
	ProjectStatusActive    ProjectStatus = "active"
	ProjectStatusBlocked   ProjectStatus = "blocked"
	ProjectStatusComplete  ProjectStatus = "complete"
	ProjectStatusArchived  ProjectStatus = "archived"
)

// CompletionMetrics tracks hand-off ready criteria
type CompletionMetrics struct {
	// Original fields
	HasRunnableBuild bool     `json:"has_runnable_build"`
	HasTests         bool     `json:"has_tests"`
	HasReadme        bool     `json:"has_readme"`
	CompletionPct    float64  `json:"completion_pct"`
	BlockingIssues   []string `json:"blocking_issues"`

	// Day 1: Build Verification
	SyntaxValid    bool   `json:"syntax_valid"`     // Code parses without errors
	DependenciesOK bool   `json:"dependencies_ok"`  // All packages installable
	BuildLog       string `json:"build_log"`        // Full build output

	// Day 2: Runtime & Testing (for future implementation)
	RuntimeVerified bool `json:"runtime_verified"` // App starts successfully
	TestsExecuted   bool `json:"tests_executed"`   // Tests were run
	TestsPassed     int  `json:"tests_passed"`     // Number of passing tests
	TestsFailed     int  `json:"tests_failed"`     // Number of failing tests

	// Day 3: Deployment & Quality (for future implementation)
	DeploymentReady   bool `json:"deployment_ready"`    // Docker builds successfully
	DockerBuilds      bool `json:"docker_builds"`       // Dockerfile valid
	EnvVarsDocumented bool `json:"env_vars_documented"` // .env.example complete
	QualityScore      int  `json:"quality_score"`       // 0-100 overall score
}

// ValidationResults stores Triple Guarantee System results
type ValidationResults struct {
	// Build verification
	BuildVerified   bool     `json:"build_verified"`
	SyntaxValid     bool     `json:"syntax_valid"`
	DependenciesOK  bool     `json:"dependencies_ok"`
	EntryPointValid bool     `json:"entry_point_valid"`
	BuildErrors     []string `json:"build_errors"`

	// Runtime verification
	RuntimeVerified   bool     `json:"runtime_verified"`
	ApplicationStarts bool     `json:"application_starts"`
	HealthCheckPassed bool     `json:"health_check_passed"`
	RuntimeErrors     []string `json:"runtime_errors"`
	RuntimeWarnings   []string `json:"runtime_warnings"`

	// Test execution
	TestsExecuted bool   `json:"tests_executed"`
	TestsPassed   int    `json:"tests_passed"`
	TestsFailed   int    `json:"tests_failed"`
	TestsSkipped  int    `json:"tests_skipped"`
	TotalTests    int    `json:"total_tests"`
	TestFramework string `json:"test_framework"`
	TestErrors    []string `json:"test_errors"`

	LastValidated time.Time `json:"last_validated"`
}

// PhaseTransitions defines valid phase transitions
var PhaseTransitions = map[Phase][]Phase{
	PhaseDiscovery:   {PhaseValidation},
	PhaseValidation:  {PhasePlanning},
	PhasePlanning:    {PhaseCodeGen},
	PhaseCodeGen:     {PhaseReview},
	PhaseReview:      {PhaseQA, PhaseCodeGen}, // Can re-generate code
	PhaseQA:          {PhaseDocs, PhaseReview}, // Can re-review
	PhaseDocs:        {PhaseComplete},
	PhaseComplete:    {}, // Terminal state
}

// CanTransition checks if a phase transition is valid
func (p Phase) CanTransition(to Phase) bool {
	validTransitions, exists := PhaseTransitions[p]
	if !exists {
		return false
	}

	for _, validPhase := range validTransitions {
		if validPhase == to {
			return true
		}
	}

	return false
}

// CanGoBackTo checks if we can revert to a previous phase
func (p Phase) CanGoBackTo(targetPhase Phase) bool {
	// Define phase order
	phaseOrder := []Phase{
		PhaseDiscovery,
		PhaseValidation,
		PhasePlanning,
		PhaseCodeGen,
		PhaseReview,
		PhaseQA,
		PhaseDocs,
		PhaseComplete,
	}

	currentIdx := -1
	targetIdx := -1

	for i, phase := range phaseOrder {
		if phase == p {
			currentIdx = i
		}
		if phase == targetPhase {
			targetIdx = i
		}
	}

	// Can go back if target is before current
	return targetIdx >= 0 && currentIdx >= 0 && targetIdx < currentIdx
}

// PhaseWeights defines completion percentage weights for each phase
var PhaseWeights = map[Phase]float64{
	PhaseDiscovery:   10.0,
	PhaseValidation:  10.0,
	PhasePlanning:    10.0,
	PhaseCodeGen:     30.0, // Biggest weight - code is critical
	PhaseReview:      20.0,
	PhaseQA:          10.0,
	PhaseDocs:        10.0,
	PhaseComplete:    0.0, // Bonus weight, not base
}
