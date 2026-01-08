package project

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

// ProjectManager manages project lifecycle and persistence
type ProjectManager struct {
	projectsDir string
	projects    map[string]*Project // In-memory cache
	projectsMux sync.RWMutex
}

// NewProjectManager creates a new project manager
func NewProjectManager(projectsDir string) (*ProjectManager, error) {
	// Create projects directory if it doesn't exist
	if err := os.MkdirAll(projectsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create projects directory: %w", err)
	}

	pm := &ProjectManager{
		projectsDir: projectsDir,
		projects:    make(map[string]*Project),
	}

	// Load existing projects into cache
	if err := pm.loadAllProjects(); err != nil {
		return nil, fmt.Errorf("failed to load existing projects: %w", err)
	}

	return pm, nil
}

// CreateProject creates a new project
func (pm *ProjectManager) CreateProject(name, description string) (*Project, error) {
	pm.projectsMux.Lock()
	defer pm.projectsMux.Unlock()

	now := time.Now()

	project := &Project{
		ID:            uuid.New().String(),
		Name:          name,
		Description:   description,
		CurrentPhase:  PhaseDiscovery,
		Phases:        []PhaseExecution{},
		Tasks:         []TaskExecution{},
		ArtifactPaths: []string{},
		Metadata: ProjectMetadata{
			TechStack: []string{},
		},
		CreatedAt: now,
		UpdatedAt: now,
		Status:    ProjectStatusActive,
	}

	// Initialize first phase (Discovery) as pending
	project.Phases = append(project.Phases, PhaseExecution{
		Phase:        PhaseDiscovery,
		Status:       PhaseStatusPending,
		StartedAt:    now,
		AgentOutputs: make(map[string]string),
	})

	// Save to disk
	if err := pm.saveProjectToDisk(project); err != nil {
		return nil, fmt.Errorf("failed to save project: %w", err)
	}

	// Add to cache
	pm.projects[project.ID] = project

	return project, nil
}

// GetProject retrieves a project by ID
func (pm *ProjectManager) GetProject(id string) (*Project, error) {
	pm.projectsMux.RLock()
	defer pm.projectsMux.RUnlock()

	project, exists := pm.projects[id]
	if !exists {
		return nil, fmt.Errorf("project not found: %s", id)
	}

	return project, nil
}

// ListProjects returns all projects
func (pm *ProjectManager) ListProjects() []*Project {
	pm.projectsMux.RLock()
	defer pm.projectsMux.RUnlock()

	projects := make([]*Project, 0, len(pm.projects))
	for _, project := range pm.projects {
		projects = append(projects, project)
	}

	return projects
}

// SaveProject saves a project to disk and updates cache
func (pm *ProjectManager) SaveProject(project *Project) error {
	pm.projectsMux.Lock()
	defer pm.projectsMux.Unlock()

	project.UpdatedAt = time.Now()

	if err := pm.saveProjectToDisk(project); err != nil {
		return fmt.Errorf("failed to save project: %w", err)
	}

	pm.projects[project.ID] = project

	return nil
}

// LoadProject loads a project from disk by ID
func (pm *ProjectManager) LoadProject(id string) (*Project, error) {
    pm.projectsMux.Lock()
    defer pm.projectsMux.Unlock()

    projectPath := pm.getProjectPath(id)

    data, err := os.ReadFile(projectPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read project file: %w", err)
    }

    var project Project
    if err := json.Unmarshal(data, &project); err != nil {
        return nil, fmt.Errorf("failed to parse project JSON: %w", err)
    }

    // Validate schema after loading
    if err := pm.ValidateProjectSchema(&project); err != nil {
        fmt.Printf("Warning: project %s failed schema validation: %v\n", id, err)
        // Continue loading but mark as potentially problematic
    }

    pm.projects[project.ID] = &project

    return &project, nil
}

// DeleteProject deletes a project from disk and cache
func (pm *ProjectManager) DeleteProject(id string) error {
	pm.projectsMux.Lock()
	defer pm.projectsMux.Unlock()

	projectPath := pm.getProjectPath(id)

	if err := os.Remove(projectPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete project file: %w", err)
	}

	delete(pm.projects, id)

	return nil
}

// UpdateProjectPhase updates the current phase of a project
func (pm *ProjectManager) UpdateProjectPhase(project *Project, phase Phase, status PhaseStatus) error {
	// Update current phase
	project.CurrentPhase = phase

	// Find or create phase execution record
	var phaseExec *PhaseExecution
	for i := range project.Phases {
		if project.Phases[i].Phase == phase {
			phaseExec = &project.Phases[i]
			break
		}
	}

	if phaseExec == nil {
		// Create new phase execution
		newPhaseExec := PhaseExecution{
			Phase:        phase,
			Status:       status,
			StartedAt:    time.Now(),
			AgentOutputs: make(map[string]string),
		}
		project.Phases = append(project.Phases, newPhaseExec)
	} else {
		// Update existing phase execution
		phaseExec.Status = status
		if status == PhaseStatusComplete {
			now := time.Now()
			phaseExec.CompletedAt = &now
		}
	}

	return pm.SaveProject(project)
}

// AddTaskExecution adds a task execution to a project
func (pm *ProjectManager) AddTaskExecution(project *Project, task TaskExecution) error {
	project.Tasks = append(project.Tasks, task)
	return pm.SaveProject(project)
}

// AddArtifactPath adds an artifact path to a project
func (pm *ProjectManager) AddArtifactPath(project *Project, artifactPath string) error {
	project.ArtifactPaths = append(project.ArtifactPaths, artifactPath)
	return pm.SaveProject(project)
}

// saveProjectToDisk saves a project to disk (assumes lock is held)
func (pm *ProjectManager) saveProjectToDisk(project *Project) error {
	projectPath := pm.getProjectPath(project.ID)

	data, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal project JSON: %w", err)
	}

	// Atomic write: write to temp file, then rename
	tempPath := projectPath + ".tmp"
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write project file: %w", err)
	}

	if err := os.Rename(tempPath, projectPath); err != nil {
		os.Remove(tempPath) // Clean up temp file
		return fmt.Errorf("failed to rename project file: %w", err)
	}

	return nil
}

// loadAllProjects loads all projects from disk into cache
func (pm *ProjectManager) loadAllProjects() error {
    entries, err := os.ReadDir(pm.projectsDir)
    if err != nil {
        if os.IsNotExist(err) {
            return nil // Directory doesn't exist yet, no projects to load
        }
        return fmt.Errorf("failed to read projects directory: %w", err)
    }

    for _, entry := range entries {
        if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
            continue
        }

        projectPath := filepath.Join(pm.projectsDir, entry.Name())
        data, err := os.ReadFile(projectPath)
        if err != nil {
            // Log error but continue loading other projects
            fmt.Printf("Warning: failed to read project file %s: %v\n", entry.Name(), err)
            continue
        }

        var project Project
        if err := json.Unmarshal(data, &project); err != nil {
            fmt.Printf("Warning: failed to parse project file %s: %v\n", entry.Name(), err)
            continue
        }

        // Validate schema after unmarshalling
        if err := pm.ValidateProjectSchema(&project); err != nil {
            fmt.Printf("Warning: project %s failed schema validation: %v\n", project.ID, err)
            // Continue loading but keep project in cache for debugging
        }

        pm.projects[project.ID] = &project
    }

    return nil
}

// getProjectPath returns the file path for a project
func (pm *ProjectManager) getProjectPath(id string) string {
	return filepath.Join(pm.projectsDir, fmt.Sprintf("project_%s.json", id))
}

// ValidateProjectSchema performs basic validation on a Project struct
func (pm *ProjectManager) ValidateProjectSchema(project *Project) error {
	if project.ID == "" {
		return fmt.Errorf("project ID cannot be empty")
	}
	if project.Name == "" {
		return fmt.Errorf("project name cannot be empty")
	}
	// Add more validation rules as needed
	return nil
}
