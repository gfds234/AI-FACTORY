package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"ai-studio/orchestrator/task"
)

// gameDesignSystemPrompt provides expert context for chat conversations
const gameDesignSystemPrompt = `You are an expert game design consultant with deep knowledge in:

**Game Design Expertise:**
- Core game mechanics and systems design
- Player psychology and motivation (intrinsic vs extrinsic rewards)
- Flow theory and engagement loops
- Progression systems and difficulty curves
- Game economy and monetization ethics

**Industry Knowledge:**
- Current market trends and player preferences
- Successful game patterns across genres (roguelikes, strategy, RPGs, etc.)
- Common pitfalls and why games fail
- Platform-specific considerations (mobile, PC, console)

**Research-Backed Insights:**
- Player retention psychology (variable reward schedules, loss aversion)
- Addiction patterns vs healthy engagement
- Accessibility and inclusive design principles
- Data-driven design decisions

**Your Approach:**
- Provide specific, actionable advice with concrete examples
- Reference successful games that demonstrate concepts
- Highlight potential risks and challenges early
- Ask clarifying questions to better understand the vision
- Balance innovation with proven patterns
- Consider both player experience and development feasibility

When discussing addictive mechanics, emphasize ethical design that creates compelling experiences without exploiting psychological vulnerabilities.

Respond conversationally but with expertise. Keep answers focused and practical.`

// Server handles HTTP API requests
type Server struct {
	taskMgr       *task.Manager
	port          int
	mux           *http.ServeMux
	conversations map[string]*Conversation
	convMux       sync.RWMutex
}

// Conversation represents a chat conversation
type Conversation struct {
	ID        string
	Messages  []Message
	Context   []int
	Model     string
	CreatedAt time.Time
}

// Message represents a chat message
type Message struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// TaskRequest represents an incoming task request
type TaskRequest struct {
	TaskType string `json:"task_type"`
	Input    string `json:"input"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// NewServer creates a new API server
func NewServer(taskMgr *task.Manager, port int) *Server {
	s := &Server{
		taskMgr:       taskMgr,
		port:          port,
		mux:           http.NewServeMux(),
		conversations: make(map[string]*Conversation),
	}

	s.registerRoutes()
	return s
}

// registerRoutes sets up HTTP handlers
func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/health", s.handleHealth)
	s.mux.HandleFunc("/task", s.handleTask)
	s.mux.HandleFunc("/history", s.handleHistory)
	s.mux.HandleFunc("/export", s.handleExport)
	s.mux.HandleFunc("/chat", s.handleChat)
	s.mux.HandleFunc("/api", s.handleAPIInfo)
	s.mux.HandleFunc("/", s.handleRoot)
}

// Start begins listening for HTTP requests
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("AI Studio Manager Console")
	log.Printf("Open in browser: http://localhost:%d", s.port)
	log.Printf("API endpoints:")
	log.Printf("  GET  /       - Web UI")
	log.Printf("  GET  /health - Health check")
	log.Printf("  POST /task   - Execute task")
	return http.ListenAndServe(addr, s.mux)
}

// handleHealth returns service health status
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Ping Ollama
	if err := s.taskMgr.Ping(); err != nil {
		s.respondError(w, fmt.Sprintf("Ollama not accessible: %v", err), http.StatusServiceUnavailable)
		return
	}

	s.respondJSON(w, map[string]string{
		"status":  "healthy",
		"backend": "ollama",
	})
}

// handleTask processes task execution requests
func (s *Server) handleTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.respondError(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req TaskRequest
	if err := json.Unmarshal(body, &req); err != nil {
		s.respondError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.TaskType == "" {
		s.respondError(w, "task_type is required", http.StatusBadRequest)
		return
	}
	if req.Input == "" {
		s.respondError(w, "input is required", http.StatusBadRequest)
		return
	}

	// Execute task
	log.Printf("Executing task: type=%s, input_length=%d", req.TaskType, len(req.Input))
	result, err := s.taskMgr.ExecuteTask(req.TaskType, req.Input)
	
	if err != nil {
		log.Printf("Task execution failed: %v", err)
		s.respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Task completed: duration=%.2fs, artifact=%s", result.Duration, result.ArtifactPath)
	s.respondJSON(w, result)
}

// handleRoot serves the web UI
func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Serve the web UI HTML file
	htmlPath := "web/index.html"
	if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
		s.respondError(w, "Web UI not found. Make sure web/index.html exists.", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, htmlPath)
}

// handleAPIInfo provides basic API info
func (s *Server) handleAPIInfo(w http.ResponseWriter, r *http.Request) {
	s.respondJSON(w, map[string]interface{}{
		"service": "AI Studio Orchestrator",
		"version": "0.1.0",
		"endpoints": []string{
			"GET  /      - Web UI",
			"GET  /api   - API info",
			"GET  /health - Health check",
			"POST /task   - Execute task",
		},
	})
}

// respondJSON writes a JSON response
func (s *Server) respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

// respondError writes a JSON error response
func (s *Server) respondError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

// handleHistory returns task history
func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	taskType := r.URL.Query().Get("task_type")
	if taskType == "" {
		taskType = "all"
	}

	history := s.taskMgr.GetHistory(taskType)
	s.respondJSON(w, history)
}

// handleExport exports history as markdown
func (s *Server) handleExport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	history := s.taskMgr.GetHistory("all")

	// Generate markdown
	var md strings.Builder
	md.WriteString("# AI FACTORY - Task History Export\n\n")
	md.WriteString(fmt.Sprintf("**Exported:** %s\n\n", time.Now().Format(time.RFC3339)))
	md.WriteString(fmt.Sprintf("**Total Tasks:** %d\n\n", len(history)))
	md.WriteString("---\n\n")

	for i, task := range history {
		md.WriteString(fmt.Sprintf("## Task %d: %s\n\n", i+1, task.TaskType))
		md.WriteString(fmt.Sprintf("**Timestamp:** %s\n", task.Timestamp.Format(time.RFC3339)))
		md.WriteString(fmt.Sprintf("**Model:** %s\n", task.Model))
		md.WriteString(fmt.Sprintf("**Duration:** %.2fs\n\n", task.Duration))
		md.WriteString("### Input\n\n")
		md.WriteString(task.Input)
		md.WriteString("\n\n### Output\n\n")
		md.WriteString(task.Output)
		md.WriteString("\n\n---\n\n")
	}

	// Send as downloadable file
	w.Header().Set("Content-Type", "text/markdown")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"history_%d.md\"", time.Now().Unix()))
	w.Write([]byte(md.String()))
}

// handleChat handles chat conversation requests
func (s *Server) handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ConversationID string `json:"conversation_id"`
		Message        string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Message) == "" {
		s.respondError(w, "Message cannot be empty", http.StatusBadRequest)
		return
	}

	// Get or create conversation
	s.convMux.Lock()
	conv := s.conversations[req.ConversationID]
	if conv == nil {
		conv = &Conversation{
			ID:        generateID(),
			Messages:  []Message{},
			Context:   nil,
			Model:     "mistral:7b-instruct-v0.2-q4_K_M",
			CreatedAt: time.Now(),
		}
		s.conversations[conv.ID] = conv
	}
	s.convMux.Unlock()

	// Add user message
	userMsg := Message{
		Role:      "user",
		Content:   req.Message,
		Timestamp: time.Now(),
	}
	conv.Messages = append(conv.Messages, userMsg)

	// Build prompt - prepend system prompt on first message
	prompt := req.Message
	if conv.Context == nil {
		prompt = gameDesignSystemPrompt + "\n\n---\n\nUser: " + req.Message
	}

	// Generate response with context
	client := s.taskMgr.GetClient()
	response, newContext, err := client.GenerateWithContext(conv.Model, prompt, conv.Context)
	if err != nil {
		s.respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add assistant message
	assistantMsg := Message{
		Role:      "assistant",
		Content:   response,
		Timestamp: time.Now(),
	}
	conv.Messages = append(conv.Messages, assistantMsg)
	conv.Context = newContext

	// Save to artifact if conversation is substantial (>= 6 messages)
	if len(conv.Messages) >= 6 {
		s.saveConversationArtifact(conv)
	}

	s.respondJSON(w, map[string]interface{}{
		"conversation_id": conv.ID,
		"response":        response,
		"timestamp":       assistantMsg.Timestamp,
	})
}

// generateID creates a unique conversation ID
func generateID() string {
	return fmt.Sprintf("chat_%d", time.Now().Unix())
}

// saveConversationArtifact saves the conversation to a file
func (s *Server) saveConversationArtifact(conv *Conversation) {
	artifactPath := fmt.Sprintf("artifacts/chat_%d.md", time.Now().Unix())

	var content strings.Builder
	content.WriteString("# Chat Conversation\n\n")
	content.WriteString(fmt.Sprintf("**ID:** %s\n", conv.ID))
	content.WriteString(fmt.Sprintf("**Started:** %s\n", conv.CreatedAt.Format(time.RFC3339)))
	content.WriteString(fmt.Sprintf("**Model:** %s\n\n", conv.Model))
	content.WriteString("---\n\n")

	for _, msg := range conv.Messages {
		role := "User"
		if msg.Role == "assistant" {
			role = "AI"
		}
		content.WriteString(fmt.Sprintf("### %s (%s)\n\n", role, msg.Timestamp.Format("15:04:05")))
		content.WriteString(msg.Content)
		content.WriteString("\n\n")
	}

	os.WriteFile(artifactPath, []byte(content.String()), 0644)
}
