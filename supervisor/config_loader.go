package supervisor

import (
	"ai-studio/orchestrator/config"
	"encoding/json"
	"fmt"
	"os"
)

// ExtendedConfig includes supervisor settings alongside base config
type ExtendedConfig struct {
	OllamaURL      string            `json:"ollama_url"`
	Models         map[string]string `json:"models"`
	ArtifactsDir   string            `json:"artifacts_dir"`
	MaxRetries     int               `json:"max_retries"`
	TimeoutSeconds int               `json:"timeout_seconds"`
	Supervisor     *SupervisorConfig `json:"supervisor"`
}

// LoadConfig loads both base and supervisor configuration
func LoadConfig() (*config.Config, *SupervisorConfig, error) {
	// Load base config first
	baseConfig, err := config.Load()
	if err != nil {
		return nil, nil, err
	}

	// Try to load extended config with supervisor settings
	configPath := "config.json"
	data, err := os.ReadFile(configPath)
	if err != nil {
		// Return base config with default supervisor config
		return baseConfig, DefaultSupervisorConfig(), nil
	}

	var extConfig ExtendedConfig
	if err := json.Unmarshal(data, &extConfig); err != nil {
		return nil, nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// If supervisor config not present, use defaults
	if extConfig.Supervisor == nil {
		extConfig.Supervisor = DefaultSupervisorConfig()
	}

	return baseConfig, extConfig.Supervisor, nil
}

// DefaultSupervisorConfig returns default supervisor settings
func DefaultSupervisorConfig() *SupervisorConfig {
	return &SupervisorConfig{
		Enabled: false, // Start disabled - opt-in
		QualityGates: QualityGatesConfig{
			RequirementsCheck: false,
			TechStackApproval: false,
			ScopeValidation:   false,
		},
		Agents: AgentsConfig{
			Requirements: AgentConfig{
				Enabled: false,
				Model:   "mistral:7b-instruct-v0.2-q4_K_M",
			},
			QA: AgentConfig{
				Enabled: false,
				Model:   "llama3:8b",
			},
			Testing: AgentConfig{
				Enabled: false,
				Model:   "deepseek-coder:6.7b-instruct",
			},
			Documentation: AgentConfig{
				Enabled: false,
				Model:   "mistral:7b-instruct-v0.2-q4_K_M",
			},
		},
		ComplexityThreshold: 7, // 1-10 scale
		ClaudeCodeEndpoint:  "", // Not configured by default
	}
}
