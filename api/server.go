package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"ai-studio/orchestrator/task"
)

// Server handles HTTP API requests
type Server struct {
	taskMgr *task.Manager
	port    int
	mux     *http.ServeMux
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
		taskMgr: taskMgr,
		port:    port,
		mux:     http.NewServeMux(),
	}

	s.registerRoutes()
	return s
}

// registerRoutes sets up HTTP handlers
func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/health", s.handleHealth)
	s.mux.HandleFunc("/task", s.handleTask)
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
