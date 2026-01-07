package api

import (
	"net/http"
	"os"
	"strings"
)

// authMiddleware validates API key for protected endpoints
func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for public endpoints
		if r.URL.Path == "/" || r.URL.Path == "/health" || r.URL.Path == "/api" {
			next(w, r)
			return
		}

		expectedKey := os.Getenv("API_KEY")
		if expectedKey == "" {
			// No key configured - dev mode
			next(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			s.respondError(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		providedKey := strings.TrimPrefix(authHeader, "Bearer ")
		providedKey = strings.TrimSpace(providedKey)

		if providedKey != expectedKey {
			s.respondError(w, "Invalid API key", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

// corsMiddleware adds CORS headers for external access
func (s *Server) corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// wrapMiddleware chains middleware
func (s *Server) wrapMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return s.corsMiddleware(s.authMiddleware(handler))
}
