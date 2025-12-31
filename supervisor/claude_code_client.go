package supervisor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ClaudeCodeClient handles communication with Claude Code endpoint
type ClaudeCodeClient struct {
	endpoint string
	client   *http.Client
}

// NewClaudeCodeClient creates a Claude Code client
func NewClaudeCodeClient(endpoint string) *ClaudeCodeClient {
	return &ClaudeCodeClient{
		endpoint: endpoint,
		client: &http.Client{
			Timeout: 300 * time.Second, // 5 min for complex tasks
		},
	}
}

// ClaudeCodeRequest represents a request to Claude Code
type ClaudeCodeRequest struct {
	Task     string                 `json:"task"`
	Input    string                 `json:"input"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ClaudeCodeResponse represents Claude Code's response
type ClaudeCodeResponse struct {
	Output   string  `json:"output"`
	Files    []File  `json:"files,omitempty"`
	Duration float64 `json:"duration"`
	Error    string  `json:"error,omitempty"`
}

// File represents a generated file
type File struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// Generate sends a task to Claude Code and returns the result
func (cc *ClaudeCodeClient) Generate(taskType, input string) (string, error) {
	req := ClaudeCodeRequest{
		Task:  taskType,
		Input: input,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := cc.client.Post(
		cc.endpoint,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("claude code request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("claude code returned status %d: %s", resp.StatusCode, string(body))
	}

	var ccResp ClaudeCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&ccResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if ccResp.Error != "" {
		return "", fmt.Errorf("claude code error: %s", ccResp.Error)
	}

	return ccResp.Output, nil
}

// Ping checks if Claude Code endpoint is accessible
func (cc *ClaudeCodeClient) Ping() error {
	resp, err := cc.client.Get(cc.endpoint + "/health")
	if err != nil {
		return fmt.Errorf("claude code not accessible: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("claude code returned status %d", resp.StatusCode)
	}

	return nil
}
