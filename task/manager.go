package task

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"ai-studio/orchestrator/config"
	"ai-studio/orchestrator/llm"
)

// Manager handles task execution and routing
type Manager struct {
	cfg        *config.Config
	client     *llm.Client
	history    []Result
	historyMux sync.RWMutex
	maxHistory int
}

// Result represents the output of a task execution
type Result struct {
	TaskType     string    `json:"task_type"`
	Input        string    `json:"input"`
	Output       string    `json:"output"`
	Model        string    `json:"model"`
	ArtifactPath string    `json:"artifact_path"`
	Duration     float64   `json:"duration_seconds"`
	Timestamp    time.Time `json:"timestamp"`
	Error        string    `json:"error,omitempty"`
}

// NewManager creates a new task manager
func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		cfg:        cfg,
		client:     llm.NewClient(cfg.OllamaURL, cfg.Timeout),
		history:    make([]Result, 0, 20),
		maxHistory: 20,
	}
}

// ExecuteTask routes and executes a task
func (m *Manager) ExecuteTask(taskType, input string) (interface{}, error) {
	return m.ExecuteTaskWithThinking(taskType, input, "normal")
}

// ExecuteTaskWithThinking routes and executes a task with specified thinking mode
func (m *Manager) ExecuteTaskWithThinking(taskType, input, thinkingMode string) (interface{}, error) {
	start := time.Now()
	result := &Result{
		TaskType:  taskType,
		Input:     input,
		Timestamp: start,
	}

	// Get model for this task type
	model, ok := m.cfg.Models[taskType]
	if !ok {
		result.Error = fmt.Sprintf("unknown task type: %s", taskType)
		return result, fmt.Errorf(result.Error)
	}
	result.Model = model

	// Build prompt based on task type
	prompt, err := m.buildPrompt(taskType, input)
	if err != nil {
		result.Error = err.Error()
		return result, err
	}

	// Execute with retries and thinking mode
	var output string
	var lastErr error

	log.Printf("Executing task with %s thinking mode", thinkingMode)

	for attempt := 0; attempt <= m.cfg.MaxRetries; attempt++ {
		output, lastErr = m.client.GenerateWithThinking(model, prompt, thinkingMode)
		if lastErr == nil {
			break
		}

		if attempt < m.cfg.MaxRetries {
			time.Sleep(time.Second * time.Duration(attempt+1)) // Exponential backoff
		}
	}

	if lastErr != nil {
		result.Error = lastErr.Error()
		return result, lastErr
	}

	result.Output = output
	result.Duration = time.Since(start).Seconds()

	// Save artifact
	artifactPath, err := m.saveArtifact(result)
	if err != nil {
		// Non-fatal - still return the result
		result.Error = fmt.Sprintf("artifact save failed: %v", err)
	} else {
		result.ArtifactPath = artifactPath
	}

	// Add to history
	m.addToHistory(result)

	return result, nil
}

// buildPrompt constructs the prompt for each task type
func (m *Manager) buildPrompt(taskType, input string) (string, error) {
	switch taskType {
	case "validate":
		return m.buildValidationPrompt(input), nil
	case "review":
		return m.buildReviewPrompt(input), nil
	case "code":
		return m.buildCodePrompt(input), nil
	default:
		return "", fmt.Errorf("unknown task type: %s", taskType)
	}
}

// buildValidationPrompt creates prompt for idea validation
func (m *Manager) buildValidationPrompt(input string) string {
	return fmt.Sprintf(`You are a game design consultant. Analyze the following game idea and provide structured feedback.

Game Idea:
%s

Provide your analysis in the following format:

## Core Concept
[1-2 sentence summary of the idea]

## Strengths
- [List 2-3 key strengths]

## Potential Issues
- [List 2-3 concerns or challenges]

## Market Viability
[Brief assessment of target audience and market fit]

## Recommendation
[Clear recommendation: Proceed, Revise, or Reconsider]

## Next Steps
[2-3 concrete next steps if proceeding]

Be direct and honest. Focus on actionable insights.`, input)
}

// buildReviewPrompt creates prompt for architecture review
func (m *Manager) buildReviewPrompt(input string) string {
	return fmt.Sprintf(`You are a senior software architect and technical reviewer. Review the following architecture, code, or technical proposal.

Document to Review:
%s

Provide your review in the following format:

## Summary
[1-2 sentence overview of what's being reviewed]

## Tech Stack Analysis
- [Evaluate technology choices - are they appropriate?]
- [Are there better alternatives for this use case?]
- [Consider: performance, maintainability, ecosystem, learning curve]

## Architecture Assessment
- [Evaluate overall design and structure]
- [Identify architectural strengths]

## Risk Assessment
- [Technical risks, performance concerns, scalability issues]
- [Security considerations]
- [Maintenance and long-term viability]

## Code Quality (if applicable)
- [Code structure, readability, best practices]
- [Potential bugs or issues]

## Recommendations
- [Specific improvements with rationale]
- [Alternative approaches to consider]

## Standards Compliance
- [Does this meet production standards?]
- [What needs to change before approval?]

## Verdict
[Overall assessment: Approved, Approved with Changes, Needs Revision, or Rejected]

Be specific and technical. Focus on practical implications and real-world viability.`, input)
}

// buildCodePrompt creates prompt for code generation with intelligent tech stack selection
func (m *Manager) buildCodePrompt(input string) string {
	return fmt.Sprintf(`You are an expert full-stack software architect with deep knowledge across multiple domains:

**Game Development:**
- Unity/C#, Godot/GDScript, Unreal/C++
- 2D/3D engines, game mechanics, physics

**Mobile Development:**
- Flutter/Dart (cross-platform)
- React Native, Swift (iOS), Kotlin (Android)

**Web Development:**
- React/TypeScript, Vue, Next.js
- Node.js, Python FastAPI, Go backends

**Desktop Applications:**
- Electron, Python/Tkinter, C#/.NET, Rust

**Backend Services:**
- Go, Python, Node.js
- REST APIs, databases, authentication

Project Request:
%s

Your task:
1. **Analyze the request** to determine:
   - Project type (game, mobile app, web app, backend, desktop tool, etc.)
   - Complexity level (prototype, MVP, production)
   - Target platform(s)
   - Key requirements

2. **Select optimal tech stack:**
   - Choose the BEST language and framework for this specific use case
   - Prioritize: solo developer friendliness, modern ecosystem, cross-platform when beneficial
   - Consider: performance needs, learning curve, maintenance

3. **Generate production-quality code:**
   - Follow best practices for chosen language
   - Include comments explaining key decisions
   - Structure code clearly and maintainably
   - Include error handling where appropriate

**CRITICAL REQUIREMENTS - READ CAREFULLY:**
1. DO NOT generate a README template with placeholders like "Give examples" or "Add examples"
2. DO NOT output generic instructions or placeholder text
3. YOU MUST generate ACTUAL, COMPLETE, RUNNABLE source code
4. Every file must contain real implementation code, NOT TODOs or placeholders
5. The code must be production-quality and immediately executable
6. If you generate a README, it must have ACTUAL setup instructions, not placeholder text

4. **Provide complete output in this format:**

## Tech Stack Decision
**Project Type:** [Game/Mobile/Web/Backend/Desktop/etc.]
**Language:** [Chosen language]
**Framework/Engine:** [Chosen framework]
**Rationale:** [2-3 sentences explaining why this stack is optimal for this request]

## Implementation

**CRITICAL: You MUST use this EXACT format for EVERY file:**

### filename.ext
` + "```" + `[language]
[COMPLETE file content - NO PLACEHOLDERS, NO TODOS, ACTUAL WORKING CODE]
` + "```" + `

**REQUIREMENTS FOR EVERY FILE:**
- Use ### followed by the filename with extension (e.g., ### src/App.jsx)
- Wrap code in triple backticks with language specified
- Include COMPLETE, WORKING code - not comments like "// Add implementation here"
- Every function must have a real implementation, not just a comment
- If generating React components, write the FULL component with actual JSX and logic

**For web projects, you MUST create separate files:**
- index.html (main HTML structure)
- css/styles.css (all styling)
- js/app.js (all JavaScript logic)
- README.md (setup instructions)

**For full-stack web projects, ALSO include backend:**
- backend/server.js (or server.py, main.go) - Main server file
- backend/routes/ - API route handlers
- backend/models/ - Data models (if using database)
- backend/.env.example - Environment variable template
- backend/package.json (or requirements.txt, go.mod) - Dependencies

**For backend/API projects:**
- server.js (or main.go, app.py) - Main entry point
- routes/ - API endpoints
- controllers/ - Business logic
- models/ - Data models
- middleware/ - Authentication, error handling
- config/ - Configuration files
- .env.example - Environment variables
- README.md - Setup and API documentation

**For projects requiring a database, ALSO include:**
- database/schema.sql (or schema.prisma) - Database schema definition
- database/migrations/ - Migration files for schema changes
- models/ - ORM models (Sequelize, Prisma, TypeORM, GORM)
- database/seeds/ - Initial data/fixtures (optional)
- database/connection.js (or db.js, database.go) - Database connection setup
- Include database URL in .env.example

**For React/Vite projects (MUST include all these files with REAL code):**
- package.json (with vite, react, vitest dependencies)
- vite.config.js (Vite configuration)
- index.html (entry HTML file)
- src/main.jsx (React entry point)
- src/App.jsx (main App component with REAL functionality)
- src/App.css (actual styles)
- src/index.css (global styles)
- src/App.test.jsx (Vitest tests)
- README.md (ACTUAL setup instructions, not placeholders)

**For other projects, organize logically:**
- Separate concerns (UI, logic, data, config)
- Follow the chosen framework's best practices
- Include README.md with setup instructions

**Example multi-file output:**

### index.html
` + "```" + `html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>App Name</title>
    <link rel="stylesheet" href="css/styles.css">
</head>
<body>
    <!-- HTML content -->
    <script src="js/app.js"></script>
</body>
</html>
` + "```" + `

### css/styles.css
` + "```" + `css
/* Stylesheet content */
body {
    margin: 0;
    padding: 0;
}
` + "```" + `

### js/app.js
` + "```" + `javascript
// JavaScript logic
console.log('App initialized');
` + "```" + `

### README.md
` + "```" + `markdown
# Project Name

## Setup Instructions
1. [Steps to run]
2. [Dependencies needed]

## Usage
[How to use the application]
` + "```" + `

**COMPLETE React/Vite Project Example (use this structure for React projects):**

### package.json
` + "```" + `json
{
  "name": "react-app",
  "version": "1.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vite build",
    "preview": "vite preview",
    "test": "vitest"
  },
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0"
  },
  "devDependencies": {
    "@vitejs/plugin-react": "^4.0.0",
    "vite": "^4.3.9",
    "vitest": "^0.32.0",
    "@testing-library/react": "^14.0.0",
    "@testing-library/jest-dom": "^6.1.0"
  }
}
` + "```" + `

### vite.config.js
` + "```" + `javascript
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: { port: 5173 }
})
` + "```" + `

### index.html
` + "```" + `html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>React App</title>
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.jsx"></script>
  </body>
</html>
` + "```" + `

### src/main.jsx
` + "```" + `javascript
import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'
import './index.css'

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
)
` + "```" + `

### src/App.jsx
` + "```" + `javascript
import { useState } from 'react'
import './App.css'

function App() {
  const [count, setCount] = useState(0)

  return (
    <div className="App">
      <h1>React App</h1>
      <button onClick={() => setCount(count + 1)}>
        Count: {count}
      </button>
    </div>
  )
}

export default App
` + "```" + `

### src/App.test.jsx
` + "```" + `javascript
import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import '@testing-library/jest-dom'
import App from './App'

describe('App', () => {
  it('renders without crashing', () => {
    render(<App />)
    expect(screen.getByText('React App')).toBeInTheDocument()
  })

  it('displays count button', () => {
    render(<App />)
    expect(screen.getByRole('button')).toBeInTheDocument()
  })
})
` + "```" + `

### src/App.css
` + "```" + `css
.App {
  text-align: center;
  padding: 2rem;
}

button {
  padding: 0.5rem 1rem;
  font-size: 1rem;
  cursor: pointer;
}
` + "```" + `

### src/index.css
` + "```" + `css
body {
  margin: 0;
  padding: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', sans-serif;
}

#root {
  min-height: 100vh;
}
` + "```" + `

**For projects requiring a backend, also include backend files:**

### backend/server.js
` + "```" + `javascript
const express = require('express');
const app = express();

app.use(express.json());

app.get('/api/data', (req, res) => {
    res.json({ message: 'API response' });
});

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => console.log(\` + "`" + `Server running on port ${PORT}\` + "`" + `));
` + "```" + `

### backend/package.json
` + "```" + `json
{
  "name": "backend",
  "version": "1.0.0",
  "main": "server.js",
  "dependencies": {
    "express": "^4.18.0"
  }
}
` + "```" + `

### backend/.env.example
` + "```" + `
PORT=3000
DATABASE_URL=your_database_url_here
` + "```" + `

**For projects with databases, also include schema and models:**

### database/schema.sql
` + "```" + `sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
` + "```" + `

### models/User.js
` + "```" + `javascript
const { DataTypes } = require('sequelize');

module.exports = (sequelize) => {
    return sequelize.define('User', {
        username: {
            type: DataTypes.STRING,
            allowNull: false,
            unique: true
        },
        email: {
            type: DataTypes.STRING,
            allowNull: false,
            unique: true
        }
    });
};
` + "```" + `

### database/connection.js
` + "```" + `javascript
const { Sequelize } = require('sequelize');
require('dotenv').config();

const sequelize = new Sequelize(process.env.DATABASE_URL, {
    dialect: 'postgres',
    logging: false
});

module.exports = sequelize;
` + "```" + `

**For projects requiring authentication, ALSO include:**

### middleware/auth.js
` + "```" + `javascript
const jwt = require('jsonwebtoken');

module.exports = (req, res, next) => {
    const token = req.header('Authorization')?.replace('Bearer ', '');

    if (!token) {
        return res.status(401).json({ error: 'Access denied' });
    }

    try {
        const verified = jwt.verify(token, process.env.JWT_SECRET);
        req.user = verified;
        next();
    } catch (err) {
        res.status(400).json({ error: 'Invalid token' });
    }
};
` + "```" + `

### routes/auth.js
` + "```" + `javascript
const express = require('express');
const bcrypt = require('bcryptjs');
const jwt = require('jsonwebtoken');
const User = require('../models/User');

const router = express.Router();

router.post('/register', async (req, res) => {
    try {
        const { username, email, password } = req.body;

        const hashedPassword = await bcrypt.hash(password, 10);
        const user = await User.create({
            username,
            email,
            password: hashedPassword
        });

        res.status(201).json({ message: 'User created', userId: user.id });
    } catch (err) {
        res.status(400).json({ error: err.message });
    }
});

router.post('/login', async (req, res) => {
    try {
        const { email, password } = req.body;
        const user = await User.findOne({ where: { email } });

        if (!user) {
            return res.status(400).json({ error: 'Invalid credentials' });
        }

        const validPassword = await bcrypt.compare(password, user.password);
        if (!validPassword) {
            return res.status(400).json({ error: 'Invalid credentials' });
        }

        const token = jwt.sign({ id: user.id }, process.env.JWT_SECRET);
        res.json({ token });
    } catch (err) {
        res.status(500).json({ error: err.message });
    }
});

module.exports = router;
` + "```" + `

**For production deployment, ALSO include:**

### Dockerfile
` + "```" + `dockerfile
FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY . .
EXPOSE 3000
CMD ["node", "server.js"]
` + "```" + `

### docker-compose.yml
` + "```" + `yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "3000:3000"
    environment:
      - NODE_ENV=production
      - DATABASE_URL=postgresql://user:password@db:5432/dbname
      - JWT_SECRET=your-secret-key
    depends_on:
      - db

  db:
    image: postgres:15-alpine
    environment:
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=password
      - POSTGRES_DB=dbname
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
` + "```" + `

### .dockerignore
` + "```" + `
node_modules
npm-debug.log
.env
.git
.gitignore
` + "```" + `

### DEPLOYMENT.md
` + "```" + `markdown
# Deployment Guide

## Option 1: Docker (Recommended)

1. Build and run with Docker Compose:
   \` + "`" + `\` + "`" + `\` + "`" + `bash
   docker-compose up -d
   \` + "`" + `\` + "`" + `\` + "`" + `

2. Check logs:
   \` + "`" + `\` + "`" + `\` + "`" + `bash
   docker-compose logs -f
   \` + "`" + `\` + "`" + `\` + "`" + `

## Option 2: Railway

1. Install Railway CLI: \` + "`" + `npm i -g @railway/cli\` + "`" + `
2. Login: \` + "`" + `railway login\` + "`" + `
3. Initialize: \` + "`" + `railway init\` + "`" + `
4. Add PostgreSQL: \` + "`" + `railway add\` + "`" + `
5. Deploy: \` + "`" + `railway up\` + "`" + `

## Option 3: Vercel (Frontend) + Railway (Backend)

**Frontend (Vercel):**
1. Push to GitHub
2. Import project on vercel.com
3. Set environment variables

**Backend (Railway):**
1. Connect GitHub repo
2. Add PostgreSQL database
3. Set environment variables
4. Deploy automatically on push

## Environment Variables

Required variables:
- \` + "`" + `PORT\` + "`" + ` - Server port (default: 3000)
- \` + "`" + `DATABASE_URL\` + "`" + ` - PostgreSQL connection string
- \` + "`" + `JWT_SECRET\` + "`" + ` - Secret key for JWT tokens
- \` + "`" + `NODE_ENV\` + "`" + ` - Environment (production/development)
` + "```" + `

**For Python/FastAPI projects, use similar structure:**

### main.py
` + "```" + `python
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel

app = FastAPI()

class Task(BaseModel):
    title: str
    description: str

@app.get("/api/tasks")
async def get_tasks():
    return {"tasks": []}

@app.post("/api/tasks")
async def create_task(task: Task):
    return {"id": 1, **task.dict()}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
` + "```" + `

**For Go projects, use similar structure:**

### main.go
` + "```" + `go
package main

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
)

type Task struct {
    ID          int    \` + "`" + `json:"id"\` + "`" + `
    Title       string \` + "`" + `json:"title"\` + "`" + `
    Description string \` + "`" + `json:"description"\` + "`" + `
}

func getTasks(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode([]Task{})
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/api/tasks", getTasks).Methods("GET")
    http.ListenAndServe(":8000", r)
}
` + "```" + `

## Setup Instructions
[Brief summary - detailed instructions should be in README.md]

## Next Steps
[What to implement next to expand this project]

Focus on practical, working code with proper file organization that a solo developer can immediately use and understand. Always separate HTML, CSS, and JavaScript into different files for web projects. For Python projects, use virtual environments. For Go projects, use Go modules.`, input)
}

// FileContent represents a parsed file from LLM output
type FileContent struct {
	Path    string
	Content string
}

// ParseFilesFromOutput extracts multiple files from LLM output
// Looks for patterns like "### filename.ext" or "#### filename.ext" followed by code blocks
func ParseFilesFromOutput(output string) []FileContent {
	files := make([]FileContent, 0)

	// Pattern: ### filename.ext or #### filename.ext
	// We want to be strict to avoid matching generic headers. 
	// We look for a filename that HAS an extension.
	fileMarkerRegex := regexp.MustCompile(`###+\s+([^\s\n(]+\.[a-zA-Z0-9]+).*\n`)
	// Use (?s) flag for DOTALL mode - allows . to match newlines
	codeBlockRegex := regexp.MustCompile("(?s)```[a-zA-Z0-9]*\\n(.*?)```")

	// Find all file markers
	fileMarkers := fileMarkerRegex.FindAllStringSubmatchIndex(output, -1)

	if len(fileMarkers) == 0 {
		// No multi-file format found, return empty
		return files
	}

	for i, match := range fileMarkers {
		// Extract filename from capture group
		filenameStart := match[2]
		filenameEnd := match[3]
		filename := output[filenameStart:filenameEnd]

		// Find the content after this marker (until next marker or end)
		contentStart := match[1] // End of the marker line
		contentEnd := len(output)

		// If there's a next marker, content ends there
		if i+1 < len(fileMarkers) {
			contentEnd = fileMarkers[i+1][0]
		}

		contentSection := output[contentStart:contentEnd]

		// Extract code from code block
		codeMatch := codeBlockRegex.FindStringSubmatch(contentSection)
		if len(codeMatch) > 1 {
			content := strings.TrimSpace(codeMatch[1])
			files = append(files, FileContent{
				Path:    filename,
				Content: content,
			})
		}
	}

	return files
}

// saveArtifact saves the task result to a file
// For code generation tasks, it detects multi-file projects and saves them properly
func (m *Manager) saveArtifact(result *Result) (string, error) {
	filename := fmt.Sprintf("%s_%d.md",
		result.TaskType,
		result.Timestamp.Unix())

	path := m.cfg.GetArtifactPath(filename)

	// Create artifact content
	content := fmt.Sprintf(`# Task Result: %s

**Timestamp:** %s
**Model:** %s
**Duration:** %.2fs

## Input
%s

## Output
%s
`,
		result.TaskType,
		result.Timestamp.Format(time.RFC3339),
		result.Model,
		result.Duration,
		result.Input,
		result.Output,
	)

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Write artifact file
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write artifact: %w", err)
	}

	// For code generation tasks, check if multi-file format exists
	if result.TaskType == "code" {
		files := ParseFilesFromOutput(result.Output)

		// Validate parsed files - detect README template errors
		if len(files) == 0 {
			log.Printf("[WARN] No files parsed from output. First 500 chars: %s",
				truncateString(result.Output, 500))
		} else {
			// Check for README template indicators
			templateIndicators := []string{
				"Give examples",
				"Add examples",
				"Add_Names",
				"Add_inspiration",
				"your-repo-link",
				"your-directory-name",
			}

			hasTemplateError := false
			for _, file := range files {
				for _, indicator := range templateIndicators {
					if strings.Contains(file.Content, indicator) {
						hasTemplateError = true
						log.Printf("[ERROR] Detected README template in file %s: contains '%s'",
							file.Path, indicator)
					}
				}
			}

			if hasTemplateError {
				log.Printf("[ERROR] Code generation produced README template instead of actual code")
				// Clear files array to prevent saving template fragments
				files = []FileContent{}
			}
		}

		if len(files) > 0 {
			// Multi-file project detected - save to projects directory
			projectDir := filepath.Join("projects", fmt.Sprintf("generated_%d", result.Timestamp.Unix()))

			if err := m.saveMultiFileProject(projectDir, files); err != nil {
				// Non-fatal - artifact is already saved
				return path + fmt.Sprintf(" (multi-file save failed: %v)", err), nil
			}

			// Return path indicating multi-file project was created
			return fmt.Sprintf("%s (project: %s)", path, projectDir), nil
		} else if result.TaskType == "code" {
			// No files extracted - likely a template or error
			return path + " (no files extracted - check artifact for template errors)", nil
		}
	}

	return path, nil
}

// saveMultiFileProject saves parsed files to a project directory
func (m *Manager) saveMultiFileProject(projectDir string, files []FileContent) error {
	// Create project directory
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Save each file
	for _, file := range files {
		fullPath := filepath.Join(projectDir, file.Path)

		// Create subdirectories if needed (e.g., css/, js/)
		fileDir := filepath.Dir(fullPath)
		if err := os.MkdirAll(fileDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", file.Path, err)
		}

		// Write file
		if err := os.WriteFile(fullPath, []byte(file.Content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", file.Path, err)
		}
	}

	return nil
}

// truncateString truncates string to maxLen characters
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// Ping checks if the LLM backend is accessible
func (m *Manager) Ping() error {
	return m.client.Ping()
}

// addToHistory adds a task result to the history (ring buffer)
func (m *Manager) addToHistory(result *Result) {
	m.historyMux.Lock()
	defer m.historyMux.Unlock()

	m.history = append(m.history, *result)
	if len(m.history) > m.maxHistory {
		m.history = m.history[1:] // Remove oldest
	}
}

// GetHistory returns task history, optionally filtered by task type
func (m *Manager) GetHistory(taskType string) []Result {
	m.historyMux.RLock()
	defer m.historyMux.RUnlock()

	// Return in reverse chronological order
	results := make([]Result, 0, len(m.history))
	for i := len(m.history) - 1; i >= 0; i-- {
		result := m.history[i]
		if taskType == "all" || result.TaskType == taskType {
			results = append(results, result)
		}
	}
	return results
}

// GetClient returns the LLM client (for chat functionality)
func (m *Manager) GetClient() *llm.Client {
	return m.client
}

// GetWebSocketHub returns nil (Manager doesn't use WebSocket)
func (m *Manager) GetWebSocketHub() interface{} {
	return nil
}
