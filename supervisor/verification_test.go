package supervisor

import (
	"os"
	"path/filepath"
	"testing"
)

// TestVerificationAgentDetectProjectType tests project type detection
func TestVerificationAgentDetectProjectType(t *testing.T) {
	agent := NewVerificationAgent()
	
	// Create temporary test directories
	tempDir := t.TempDir()
	
	tests := []struct {
		name         string
		setupFiles   map[string]string
		expectedType ProjectType
	}{
		{
			name: "Node.js project",
			setupFiles: map[string]string{
				"package.json": `{"name": "test", "version": "1.0.0"}`,
				"server.js":    "const express = require('express');",
			},
			expectedType: ProjectTypeNodeJS,
		},
		{
			name: "Python project",
			setupFiles: map[string]string{
				"requirements.txt": "flask==2.0.0",
				"main.py":          "print('hello')",
			},
			expectedType: ProjectTypePython,
		},
		{
			name: "Go project",
			setupFiles: map[string]string{
				"go.mod":  "module test\n\ngo 1.22",
				"main.go": "package main\n\nfunc main() {}",
			},
			expectedType: ProjectTypeGo,
		},
		{
			name: "Frontend project",
			setupFiles: map[string]string{
				"index.html": "<html><head></head><body></body></html>",
			},
			expectedType: ProjectTypeFrontend,
		},
		{
			name:         "Unknown project",
			setupFiles:   map[string]string{},
			expectedType: ProjectTypeUnknown,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test directory
			testDir := filepath.Join(tempDir, tt.name)
			os.MkdirAll(testDir, 0755)
			
			// Create test files
			for filename, content := range tt.setupFiles {
				filePath := filepath.Join(testDir, filename)
				os.WriteFile(filePath, []byte(content), 0644)
			}
			
			// Test detection
			projectType, err := agent.detectProjectType(testDir)
			if err != nil {
				t.Fatalf("detectProjectType() error = %v", err)
			}
			if projectType != tt.expectedType {
				t.Errorf("detectProjectType() = %v, want %v", projectType, tt.expectedType)
			}
		})
	}
}

// TestVerificationAgentCheckEntryPoint tests entry point validation
func TestVerificationAgentCheckEntryPoint(t *testing.T) {
	agent := NewVerificationAgent()
	
	tempDir := t.TempDir()
	
	tests := []struct {
		name        string
		projectType ProjectType
		files       []string
		expectValid bool
	}{
		{
			name:        "Node.js with server.js",
			projectType: ProjectTypeNodeJS,
			files:       []string{"server.js", "package.json"},
			expectValid: true,
		},
		{
			name:        "Node.js with index.js",
			projectType: ProjectTypeNodeJS,
			files:       []string{"index.js", "package.json"},
			expectValid: true,
		},
		{
			name:        "Node.js missing entry point",
			projectType: ProjectTypeNodeJS,
			files:       []string{"package.json"},
			expectValid: false,
		},
		{
			name:        "Python with main.py",
			projectType: ProjectTypePython,
			files:       []string{"main.py"},
			expectValid: true,
		},
		{
			name:        "Python with app.py",
			projectType: ProjectTypePython,
			files:       []string{"app.py"},
			expectValid: true,
		},
		{
			name:        "Python missing entry point",
			projectType: ProjectTypePython,
			files:       []string{"utils.py"},
			expectValid: false,
		},
		{
			name:        "Go with main.go",
			projectType: ProjectTypeGo,
			files:       []string{"main.go", "go.mod"},
			expectValid: true,
		},
		{
			name:        "Go missing main.go",
			projectType: ProjectTypeGo,
			files:       []string{"go.mod"},
			expectValid: false,
		},
		{
			name:        "Frontend with index.html",
			projectType: ProjectTypeFrontend,
			files:       []string{"index.html"},
			expectValid: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test directory
			testDir := filepath.Join(tempDir, tt.name)
			os.MkdirAll(testDir, 0755)
			
			// Create test files
			for _, filename := range tt.files {
				filePath := filepath.Join(testDir, filename)
				os.WriteFile(filePath, []byte("test content"), 0644)
			}
			
			// Test entry point check
			valid, err := agent.checkEntryPoint(testDir, tt.projectType)
			
			if tt.expectValid {
				if err != nil || !valid {
					t.Errorf("checkEntryPoint() = %v, %v; want true, nil", valid, err)
				}
			} else {
				if valid {
					t.Error("checkEntryPoint() = true; want false")
				}
			}
		})
	}
}

// TestVerificationAgentFindFiles tests file finding with extension filter
func TestVerificationAgentFindFiles(t *testing.T) {
	agent := NewVerificationAgent()
	
	tempDir := t.TempDir()
	
	// Create test files
	files := []string{
		"main.js",
		"utils.js",
		"styles.css",
		"index.html",
	}
	
	for _, f := range files {
		os.WriteFile(filepath.Join(tempDir, f), []byte("test"), 0644)
	}
	
	// Create node_modules (should be skipped)
	nodeModules := filepath.Join(tempDir, "node_modules")
	os.MkdirAll(nodeModules, 0755)
	os.WriteFile(filepath.Join(nodeModules, "lib.js"), []byte("test"), 0644)
	
	// Find JS files
	jsFiles, err := agent.findFiles(tempDir, ".js")
	if err != nil {
		t.Fatalf("findFiles() error = %v", err)
	}
	
	// Should find 2 JS files (not the one in node_modules)
	if len(jsFiles) != 2 {
		t.Errorf("findFiles() found %d files, want 2", len(jsFiles))
	}
	
	// Verify node_modules was skipped
	for _, f := range jsFiles {
		if filepath.Base(filepath.Dir(f)) == "node_modules" {
			t.Error("findFiles() should skip node_modules directory")
		}
	}
}

// TestVerificationAgentHTMLSyntax tests basic HTML validation
func TestVerificationAgentHTMLSyntax(t *testing.T) {
	agent := NewVerificationAgent()
	
	tempDir := t.TempDir()
	
	tests := []struct {
		name        string
		htmlContent string
		expectValid bool
	}{
		{
			name:        "Valid HTML",
			htmlContent: "<!DOCTYPE html><html><head><title>Test</title></head><body><h1>Hello</h1></body></html>",
			expectValid: true,
		},
		{
			name:        "Missing html tag",
			htmlContent: "<head></head><body></body>",
			expectValid: false,
		},
		{
			name:        "Missing body tag",
			htmlContent: "<html><head></head></html>",
			expectValid: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := filepath.Join(tempDir, tt.name)
			os.MkdirAll(testDir, 0755)
			
			indexPath := filepath.Join(testDir, "index.html")
			os.WriteFile(indexPath, []byte(tt.htmlContent), 0644)
			
			valid, _, err := agent.validateHTMLSyntax(testDir)
			if err != nil {
				t.Fatalf("validateHTMLSyntax() error = %v", err)
			}
			if valid != tt.expectValid {
				t.Errorf("validateHTMLSyntax() = %v, want %v", valid, tt.expectValid)
			}
		})
	}
}

// TestVerificationResultCreation tests verification result initialization
func TestVerificationResultCreation(t *testing.T) {
	result := &VerificationResult{
		ProjectType:     "nodejs",
		SyntaxValid:     true,
		DependenciesOK:  true,
		EntryPointValid: true,
		Errors:          make([]string, 0),
		Warnings:        make([]string, 0),
	}
	
	if result.ProjectType != "nodejs" {
		t.Error("ProjectType should be 'nodejs'")
	}
	if !result.SyntaxValid {
		t.Error("SyntaxValid should be true")
	}
	if len(result.Errors) != 0 {
		t.Error("Errors should be empty initially")
	}
}
