package supervisor

import (
	"ai-studio/orchestrator/task"
	"time"
)

// SupervisorConfig holds supervisor-specific configuration
type SupervisorConfig struct {
	Enabled             bool               `json:"enabled"`
	QualityGates        QualityGatesConfig `json:"quality_gates"`
	Agents              AgentsConfig       `json:"agents"`
	ComplexityThreshold int                `json:"complexity_threshold"` // 1-10 scale
	ClaudeCodeEndpoint  string             `json:"claude_code_endpoint"`
}

// QualityGatesConfig controls which gates are enforced
type QualityGatesConfig struct {
	RequirementsCheck bool `json:"requirements_check"`
	TechStackApproval bool `json:"tech_stack_approval"`
	ScopeValidation   bool `json:"scope_validation"`
}

// AgentsConfig controls which agents are active
type AgentsConfig struct {
	Requirements  AgentConfig `json:"requirements"`
	QA            AgentConfig `json:"qa"`
	Testing       AgentConfig `json:"testing"`
	Documentation AgentConfig `json:"documentation"`
}

// AgentConfig for individual agent settings
type AgentConfig struct {
	Enabled bool   `json:"enabled"`
	Model   string `json:"model"` // Ollama model to use
}

// SupervisedResult extends task.Result with agent outputs
type SupervisedResult struct {
	*task.Result
	RequirementsAnalysis *AgentOutput       `json:"requirements_analysis,omitempty"`
	TechStackApproval    *AgentOutput       `json:"tech_stack_approval,omitempty"`
	ScopeValidation      *AgentOutput       `json:"scope_validation,omitempty"`
	ComplexityScore      int                `json:"complexity_score"`
	ExecutionRoute       string             `json:"execution_route"` // "ollama" or "claude_code"
	QAReview             *AgentOutput       `json:"qa_review,omitempty"`
	TestPlan             *AgentOutput       `json:"test_plan,omitempty"`
	Documentation        *AgentOutput       `json:"documentation,omitempty"`
	TotalDuration        float64            `json:"total_duration_seconds"`
	AgentDurations       map[string]float64 `json:"agent_durations"`
}

// AgentOutput represents the result from a single agent
type AgentOutput struct {
	AgentType string                 `json:"agent_type"`
	Status    string                 `json:"status"` // "passed", "failed", "warning"
	Output    string                 `json:"output"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Duration  float64                `json:"duration_seconds"`
	Timestamp time.Time              `json:"timestamp"`
}

// ComplexityAnalysis holds scoring details
type ComplexityAnalysis struct {
	Score            int            `json:"score"` // 1-10
	Indicators       map[string]int `json:"indicators"`
	Reasoning        string         `json:"reasoning"`
	RecommendedRoute string         `json:"recommended_route"`
}

// Agent interface that all agents must implement
type Agent interface {
	Execute(taskType, input string, context map[string]interface{}) (*AgentOutput, error)
	Name() string
	RequiresInput() bool // Some agents run pre-execution, some post
}
