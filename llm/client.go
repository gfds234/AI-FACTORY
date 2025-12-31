package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client interfaces with Ollama API
type Client struct {
	baseURL string
	client  *http.Client
}

// NewClient creates a new Ollama client
func NewClient(baseURL string, timeoutSeconds int) *Client {
	return &Client{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: time.Duration(timeoutSeconds) * time.Second,
		},
	}
}

// GenerateRequest represents an Ollama generation request
type GenerateRequest struct {
	Model   string `json:"model"`
	Prompt  string `json:"prompt"`
	Stream  bool   `json:"stream"`
	Context []int  `json:"context,omitempty"`
}

// GenerateResponse represents an Ollama generation response
type GenerateResponse struct {
	Model     string `json:"model"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
	Context   []int  `json:"context,omitempty"`
	TotalDuration int64 `json:"total_duration,omitempty"`
}

// Generate sends a prompt to Ollama and returns the response
func (c *Client) Generate(model, prompt string) (string, error) {
	req := GenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false, // Non-streaming for simplicity in Phase 1
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.client.Post(
		c.baseURL+"/api/generate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("ollama request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, string(body))
	}

	var genResp GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return genResp.Response, nil
}

// ListModels returns available models from Ollama
func (c *Client) ListModels() ([]string, error) {
	resp, err := c.client.Get(c.baseURL + "/api/tags")
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode models: %w", err)
	}

	models := make([]string, len(result.Models))
	for i, m := range result.Models {
		models[i] = m.Name
	}

	return models, nil
}

// GenerateWithContext sends a prompt to Ollama with conversation context and returns both response and new context
func (c *Client) GenerateWithContext(model, prompt string, context []int) (string, []int, error) {
	req := GenerateRequest{
		Model:   model,
		Prompt:  prompt,
		Stream:  false,
		Context: context,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.client.Post(
		c.baseURL+"/api/generate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", nil, fmt.Errorf("ollama request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", nil, fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, string(body))
	}

	var genResp GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		return "", nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return genResp.Response, genResp.Context, nil
}

// Ping checks if Ollama is accessible
func (c *Client) Ping() error {
	resp, err := c.client.Get(c.baseURL + "/api/tags")
	if err != nil {
		return fmt.Errorf("ollama not accessible: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ollama returned status %d", resp.StatusCode)
	}

	return nil
}
