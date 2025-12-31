package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the orchestrator configuration
type Config struct {
	OllamaURL    string            `json:"ollama_url"`
	Models       map[string]string `json:"models"`        // task_type -> model_name
	ArtifactsDir string            `json:"artifacts_dir"`
	MaxRetries   int               `json:"max_retries"`
	Timeout      int               `json:"timeout_seconds"`
}

// Default configuration
func defaultConfig() *Config {
	return &Config{
		OllamaURL: "http://localhost:11434",
		Models: map[string]string{
			"validate": "mistral:7b-instruct-v0.2-q4_K_M", // Idea validation
			"review":   "llama3:8b",                       // Architecture review
			"code":     "deepseek-coder:6.7b-instruct",    // Future: code generation
		},
		ArtifactsDir: "./artifacts",
		MaxRetries:   2,
		Timeout:      120, // 2 minutes per task
	}
}

// Load reads configuration from config.json or returns defaults
func Load() (*Config, error) {
	configPath := "config.json"
	
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		cfg := defaultConfig()
		
		// Ensure artifacts directory exists
		if err := os.MkdirAll(cfg.ArtifactsDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create artifacts dir: %w", err)
		}
		
		// Save default config
		if err := Save(cfg, configPath); err != nil {
			return nil, fmt.Errorf("failed to save default config: %w", err)
		}
		
		return cfg, nil
	}

	// Load existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Ensure artifacts directory exists
	if err := os.MkdirAll(cfg.ArtifactsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create artifacts dir: %w", err)
	}

	return &cfg, nil
}

// Save writes configuration to file
func Save(cfg *Config, path string) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// GetArtifactPath returns the full path for an artifact
func (c *Config) GetArtifactPath(filename string) string {
	return filepath.Join(c.ArtifactsDir, filename)
}
