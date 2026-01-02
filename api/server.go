package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"ai-studio/orchestrator/llm"
	"ai-studio/orchestrator/project"
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

// TaskManager interface for both standard and supervised managers
type TaskManager interface {
	ExecuteTask(taskType, input string) (interface{}, error)
	Ping() error
	GetHistory(taskType string) []task.Result
	GetClient() *llm.Client
}

// Server handles HTTP API requests
type Server struct {
	taskMgr          TaskManager
	port             int
	mux              *http.ServeMux
	conversations    map[string]*Conversation
	convMux          sync.RWMutex
	discoverSessions *task.DiscoverManager
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
func NewServer(taskMgr TaskManager, port int) *Server {
	s := &Server{
		taskMgr:          taskMgr,
		port:             port,
		mux:              http.NewServeMux(),
		conversations:    make(map[string]*Conversation),
		discoverSessions: task.NewDiscoverManager(taskMgr.GetClient()),
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
	s.mux.HandleFunc("/discover", s.handleDiscover)
	s.mux.HandleFunc("/discover/history", s.handleDiscoverHistory)
	s.mux.HandleFunc("/discover/session", s.handleGetDiscoverSession)
	s.mux.HandleFunc("/api", s.handleAPIInfo)

	// Project orchestrator routes
	s.mux.HandleFunc("/project", s.handleProject)
	s.mux.HandleFunc("/project/list", s.handleProjectList)
	s.mux.HandleFunc("/project/phase", s.handleProjectPhase)
	s.mux.HandleFunc("/project/transition", s.handleProjectTransition)
	s.mux.HandleFunc("/project/approve", s.handleProjectApprove)
	s.mux.HandleFunc("/project/reject", s.handleProjectReject)
	s.mux.HandleFunc("/project/revert", s.handleProjectRevert)
	s.mux.HandleFunc("/project/metrics", s.handleProjectMetrics)

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

	// Log completion (extract duration/artifact from base task.Result)
	if taskResult, ok := result.(*task.Result); ok {
		log.Printf("Task completed: duration=%.2fs, artifact=%s", taskResult.Duration, taskResult.ArtifactPath)
	} else {
		log.Printf("Task completed (supervisor mode)")
	}

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

// saveDiscoveryArtifact saves the discovery session to a markdown file
func (s *Server) saveDiscoveryArtifact(session *task.DiscoverSession) {
	artifactPath := fmt.Sprintf("artifacts/discover_%d.md", time.Now().Unix())

	var content strings.Builder
	content.WriteString("# Discovery Session\n\n")
	content.WriteString(fmt.Sprintf("**ID:** %s\n", session.ID))
	content.WriteString(fmt.Sprintf("**Category:** %s\n", session.IdeaType))
	content.WriteString(fmt.Sprintf("**Status:** %s\n", session.Status))
	content.WriteString(fmt.Sprintf("**Verdict:** %s\n", session.Verdict))
	content.WriteString(fmt.Sprintf("**Started:** %s\n\n", session.CreatedAt.Format(time.RFC3339)))
	content.WriteString("---\n\n")

	content.WriteString("## Product Idea\n\n")
	content.WriteString(session.RawIdea)
	content.WriteString("\n\n")

	content.WriteString("## Validation Questions & Answers\n\n")
	for i, question := range session.Questions {
		content.WriteString(fmt.Sprintf("### Q%d: %s\n\n", i+1, question))
		if i < len(session.Answers) {
			content.WriteString(fmt.Sprintf("**Answer:** %s\n\n", session.Answers[i]))
		} else {
			content.WriteString("**Answer:** _(not answered)_\n\n")
		}
	}

	content.WriteString("---\n\n")
	content.WriteString("## Verdict\n\n")
	content.WriteString(fmt.Sprintf("**Decision:** %s\n\n", session.Verdict))
	content.WriteString(fmt.Sprintf("**Reasoning:** %s\n", session.Reasoning))

	os.WriteFile(artifactPath, []byte(content.String()), 0644)
	log.Printf("Discovery session saved: %s", artifactPath)
}

// handleDiscover handles discovery session requests
func (s *Server) handleDiscover(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		DiscoverID string `json:"discover_id"`
		Action     string `json:"action"`
		Input      string `json:"input"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Input) == "" {
		s.respondError(w, "Input cannot be empty", http.StatusBadRequest)
		return
	}

	if req.Action != "start" && req.Action != "answer" {
		s.respondError(w, "Action must be 'start' or 'answer'", http.StatusBadRequest)
		return
	}

	// Handle "start" action - create new session
	if req.Action == "start" {
		session := s.discoverSessions.StartSession(req.Input)
		questionNum, question := s.discoverSessions.GetCurrentQuestion(session)

		log.Printf("Started discovery session: id=%s", session.ID)

		s.respondJSON(w, map[string]interface{}{
			"discover_id":      session.ID,
			"status":           session.Status,
			"current_question": questionNum,
			"question":         question,
			"verdict":          nil,
			"reasoning":        nil,
			"timestamp":        session.CreatedAt,
			// Expose generated questions and detected category early
			"idea_type":        session.IdeaType,
			"questions":        session.Questions,
			"raw_idea":         session.RawIdea,
		})
		return
	}

	// Handle "answer" action - add answer to existing session
	if req.Action == "answer" {
		if req.DiscoverID == "" {
			s.respondError(w, "discover_id is required for 'answer' action", http.StatusBadRequest)
			return
		}

		session, err := s.discoverSessions.AddAnswer(req.DiscoverID, req.Input)
		if err != nil {
			s.respondError(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Added answer to session: id=%s, answers=%d", session.ID, len(session.Answers))

		// If session is complete, return verdict
		if session.Status == "complete" {
			log.Printf("Discovery complete: id=%s, verdict=%s", session.ID, session.Verdict)

			// Save discovery artifact
			s.saveDiscoveryArtifact(session)

			s.respondJSON(w, map[string]interface{}{
				"discover_id":      session.ID,
				"status":           session.Status,
				"current_question": 0,
				"question":         nil,
				"verdict":          session.Verdict,
				"reasoning":        session.Reasoning,
				"timestamp":        session.UpdatedAt,
				// Full session data for enriched prompts
				"raw_idea":         session.RawIdea,
				"idea_type":        session.IdeaType,
				"questions":        session.Questions,
				"answers":          session.Answers,
			})
			return
		}

		// Otherwise, return next question
		questionNum, question := s.discoverSessions.GetCurrentQuestion(session)

		s.respondJSON(w, map[string]interface{}{
			"discover_id":      session.ID,
			"status":           session.Status,
			"current_question": questionNum,
			"question":         question,
			"verdict":          nil,
			"reasoning":        nil,
			"timestamp":        session.UpdatedAt,
			// Include session context
			"idea_type":        session.IdeaType,
			"questions":        session.Questions,
			"answers":          session.Answers, // Accumulated answers so far
			"raw_idea":         session.RawIdea,
		})
	}
}

// handleDiscoverHistory returns all discovery sessions
func (s *Server) handleDiscoverHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessions := s.discoverSessions.GetAllSessions()

	// Sort by most recent first
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].CreatedAt.After(sessions[j].CreatedAt)
	})

	s.respondJSON(w, map[string]interface{}{
		"sessions": sessions,
		"count":    len(sessions),
	})
}

// handleGetDiscoverSession retrieves a specific discovery session by ID
func (s *Server) handleGetDiscoverSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionID := r.URL.Query().Get("discover_id")
	if sessionID == "" {
		s.respondError(w, "discover_id parameter required", http.StatusBadRequest)
		return
	}

	session := s.discoverSessions.GetSession(sessionID)
	if session == nil {
		s.respondError(w, "Session not found", http.StatusNotFound)
		return
	}

	s.respondJSON(w, session)
}

// Project Orchestrator Handlers

// handleProject handles creating and getting projects
func (s *Server) handleProject(w http.ResponseWriter, r *http.Request) {
	// Try to cast taskMgr to ProjectOrchestrator
	orchestrator, ok := s.taskMgr.(*project.ProjectOrchestrator)
	if !ok {
		s.respondError(w, "Project orchestrator not enabled", http.StatusNotImplemented)
		return
	}

	switch r.Method {
	case http.MethodPost:
		// Create new project
		var req struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.respondError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Description) == "" {
			s.respondError(w, "Name and description are required", http.StatusBadRequest)
			return
		}

		proj, err := orchestrator.CreateProject(req.Name, req.Description)
		if err != nil {
			s.respondError(w, fmt.Sprintf("Failed to create project: %v", err), http.StatusInternalServerError)
			return
		}

		s.respondJSON(w, proj)

	case http.MethodGet:
		// Get project by ID
		projectID := r.URL.Query().Get("id")
		if projectID == "" {
			s.respondError(w, "Project ID required", http.StatusBadRequest)
			return
		}

		proj, err := orchestrator.GetProject(projectID)
		if err != nil {
			s.respondError(w, fmt.Sprintf("Project not found: %v", err), http.StatusNotFound)
			return
		}

		s.respondJSON(w, proj)

	default:
		s.respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleProjectList lists all projects
func (s *Server) handleProjectList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orchestrator, ok := s.taskMgr.(*project.ProjectOrchestrator)
	if !ok {
		s.respondError(w, "Project orchestrator not enabled", http.StatusNotImplemented)
		return
	}

	projects := orchestrator.ListProjects()

	s.respondJSON(w, map[string]interface{}{
		"projects": projects,
		"count":    len(projects),
	})
}

// handleProjectPhase executes a project phase
func (s *Server) handleProjectPhase(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orchestrator, ok := s.taskMgr.(*project.ProjectOrchestrator)
	if !ok {
		s.respondError(w, "Project orchestrator not enabled", http.StatusNotImplemented)
		return
	}

	var req struct {
		ProjectID string         `json:"project_id"`
		Phase     project.Phase  `json:"phase"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ProjectID == "" || req.Phase == "" {
		s.respondError(w, "Project ID and phase are required", http.StatusBadRequest)
		return
	}

	result, err := orchestrator.ExecuteProjectPhase(req.ProjectID, req.Phase)
	if err != nil {
		s.respondError(w, fmt.Sprintf("Failed to execute phase: %v", err), http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, result)
}

// handleProjectTransition transitions a project to a new phase
func (s *Server) handleProjectTransition(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orchestrator, ok := s.taskMgr.(*project.ProjectOrchestrator)
	if !ok {
		s.respondError(w, "Project orchestrator not enabled", http.StatusNotImplemented)
		return
	}

	var req struct {
		ProjectID     string         `json:"project_id"`
		ToPhase       project.Phase  `json:"to_phase"`
		HumanApproval bool           `json:"human_approval"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ProjectID == "" || req.ToPhase == "" {
		s.respondError(w, "Project ID and to_phase are required", http.StatusBadRequest)
		return
	}

	err := orchestrator.TransitionPhase(req.ProjectID, req.ToPhase, req.HumanApproval)
	if err != nil {
		s.respondError(w, fmt.Sprintf("Failed to transition phase: %v", err), http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Transitioned to %s", req.ToPhase),
	})
}

// handleProjectApprove approves the current phase and transitions to next
func (s *Server) handleProjectApprove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orchestrator, ok := s.taskMgr.(*project.ProjectOrchestrator)
	if !ok {
		s.respondError(w, "Project orchestrator not enabled", http.StatusNotImplemented)
		return
	}

	var req struct {
		ProjectID string `json:"project_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ProjectID == "" {
		s.respondError(w, "Project ID is required", http.StatusBadRequest)
		return
	}

	err := orchestrator.ApprovePhase(req.ProjectID)
	if err != nil {
		s.respondError(w, fmt.Sprintf("Failed to approve phase: %v", err), http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, map[string]interface{}{
		"success": true,
		"message": "Phase approved",
	})
}

// handleProjectReject rejects the current phase
func (s *Server) handleProjectReject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orchestrator, ok := s.taskMgr.(*project.ProjectOrchestrator)
	if !ok {
		s.respondError(w, "Project orchestrator not enabled", http.StatusNotImplemented)
		return
	}

	var req struct {
		ProjectID string `json:"project_id"`
		Reason    string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ProjectID == "" {
		s.respondError(w, "Project ID is required", http.StatusBadRequest)
		return
	}

	err := orchestrator.RejectPhase(req.ProjectID, req.Reason)
	if err != nil {
		s.respondError(w, fmt.Sprintf("Failed to reject phase: %v", err), http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, map[string]interface{}{
		"success": true,
		"message": "Phase rejected",
	})
}

// handleProjectRevert reverts the project to a previous phase
func (s *Server) handleProjectRevert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orchestrator, ok := s.taskMgr.(*project.ProjectOrchestrator)
	if !ok {
		s.respondError(w, "Project orchestrator not enabled", http.StatusNotImplemented)
		return
	}

	var req struct {
		ProjectID   string `json:"project_id"`
		TargetPhase string `json:"target_phase"`
		Reason      string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ProjectID == "" {
		s.respondError(w, "Project ID is required", http.StatusBadRequest)
		return
	}

	if req.TargetPhase == "" {
		s.respondError(w, "Target phase is required", http.StatusBadRequest)
		return
	}

	// Convert string to Phase type
	targetPhase := project.Phase(req.TargetPhase)

	// Revert phase
	err := orchestrator.RevertPhase(req.ProjectID, targetPhase, req.Reason)
	if err != nil {
		s.respondError(w, fmt.Sprintf("Failed to revert phase: %v", err), http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Reverted to %s phase", targetPhase),
	})
}

// handleProjectMetrics gets completion metrics for a project
func (s *Server) handleProjectMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orchestrator, ok := s.taskMgr.(*project.ProjectOrchestrator)
	if !ok {
		s.respondError(w, "Project orchestrator not enabled", http.StatusNotImplemented)
		return
	}

	projectID := r.URL.Query().Get("project_id")
	if projectID == "" {
		s.respondError(w, "Project ID required", http.StatusBadRequest)
		return
	}

	metrics, err := orchestrator.GetCompletionMetrics(projectID)
	if err != nil {
		s.respondError(w, fmt.Sprintf("Failed to get metrics: %v", err), http.StatusInternalServerError)
		return
	}

	s.respondJSON(w, metrics)
}
