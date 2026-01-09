package supervisor

import (
	"fmt"
	"strings"
)

// ComplexityScorer analyzes task complexity
type ComplexityScorer struct {
	cfg *SupervisorConfig
}

// NewComplexityScorer creates a complexity scorer
func NewComplexityScorer(cfg *SupervisorConfig) *ComplexityScorer {
	return &ComplexityScorer{cfg: cfg}
}

// Score analyzes input and returns complexity score (1-10)
func (cs *ComplexityScorer) Score(taskType, input string) *ComplexityAnalysis {
	analysis := &ComplexityAnalysis{
		Indicators: make(map[string]int),
		Score:      1,
	}

	// Only code generation tasks get complexity scoring
	if taskType != "code" {
		analysis.RecommendedRoute = "ollama"
		analysis.Reasoning = "Non-code task - use Ollama"
		analysis.ThinkingMode = "fast" // Non-code tasks use fast mode
		return analysis
	}

	// Indicator 1: Multi-file project (score: +3)
	multiFileKeywords := []string{
		"multiple files", "project structure", "folder structure",
		"separate files", "modular", "microservice", "architecture",
	}
	if containsAny(input, multiFileKeywords) {
		analysis.Indicators["multi_file"] = 3
		analysis.Score += 3
	}

	// Indicator 2: Database/persistence (score: +2)
	dbKeywords := []string{
		"database", "postgres", "mysql", "mongodb", "sqlite",
		"sql", "orm", "migration", "persistence", "storage",
	}
	if containsAny(input, dbKeywords) {
		analysis.Indicators["database"] = 2
		analysis.Score += 2
	}

	// Indicator 3: Authentication/security (score: +2)
	authKeywords := []string{
		"authentication", "authorization", "oauth", "jwt", "login",
		"security", "encryption", "token", "session", "password",
	}
	if containsAny(input, authKeywords) {
		analysis.Indicators["auth_security"] = 2
		analysis.Score += 2
	}

	// Indicator 4: External integrations (score: +2)
	integrationKeywords := []string{
		"api integration", "third-party", "webhook", "rest api",
		"graphql", "payment", "stripe", "aws", "google cloud",
	}
	if containsAny(input, integrationKeywords) {
		analysis.Indicators["integrations"] = 2
		analysis.Score += 2
	}

	// Indicator 5: Advanced algorithms (score: +2)
	algorithmKeywords := []string{
		"algorithm", "optimization", "machine learning", "ai",
		"neural", "pathfinding", "graph", "sorting", "search algorithm",
	}
	if containsAny(input, algorithmKeywords) {
		analysis.Indicators["algorithms"] = 2
		analysis.Score += 2
	}

	// Indicator 6: Real-time features (score: +1)
	realtimeKeywords := []string{
		"real-time", "websocket", "socket.io", "streaming",
		"live updates", "push notifications", "pubsub",
	}
	if containsAny(input, realtimeKeywords) {
		analysis.Indicators["realtime"] = 1
		analysis.Score += 1
	}

	// Indicator 7: Testing requirements (score: +1)
	testKeywords := []string{
		"unit test", "integration test", "test coverage",
		"testing framework", "e2e test", "test suite",
	}
	if containsAny(input, testKeywords) {
		analysis.Indicators["testing"] = 1
		analysis.Score += 1
	}

	// Indicator 8: Large word count (score: +1 if > 200 words)
	wordCount := len(strings.Fields(input))
	if wordCount > 200 {
		analysis.Indicators["word_count"] = 1
		analysis.Score += 1
	}

	// Cap score at 10
	if analysis.Score > 10 {
		analysis.Score = 10
	}

	// Determine thinking mode based on complexity score
	// Low complexity (1-3): Fast mode - quick, direct responses
	// Medium complexity (4-6): Normal mode - balanced reasoning
	// High complexity (7-10): Extended mode - deep thinking
	if analysis.Score <= 3 {
		analysis.ThinkingMode = "fast"
	} else if analysis.Score <= 6 {
		analysis.ThinkingMode = "normal"
	} else {
		analysis.ThinkingMode = "extended"
	}

	// Determine route based on threshold
	if analysis.Score >= cs.cfg.ComplexityThreshold {
		analysis.RecommendedRoute = "claude_code"
		analysis.Reasoning = fmt.Sprintf(
			"Complexity score %d >= threshold %d - route to Claude Code with %s thinking mode",
			analysis.Score, cs.cfg.ComplexityThreshold, analysis.ThinkingMode,
		)
	} else {
		analysis.RecommendedRoute = "ollama"
		analysis.Reasoning = fmt.Sprintf(
			"Complexity score %d < threshold %d - handle with Ollama using %s thinking mode",
			analysis.Score, cs.cfg.ComplexityThreshold, analysis.ThinkingMode,
		)
	}

	return analysis
}

// containsAny checks if input contains any of the keywords (case-insensitive)
func containsAny(input string, keywords []string) bool {
	lower := strings.ToLower(input)
	for _, kw := range keywords {
		if strings.Contains(lower, strings.ToLower(kw)) {
			return true
		}
	}
	return false
}
