# Project Orchestrator Guide

## Overview

The **AI Factory Project Orchestrator** transforms the system from a single-task execution engine into a comprehensive project-based workflow manager that guides ideas from conception to "hand-off ready" completion.

**Key Philosophy:** "Ollama leads, Claude escalates"
- Lead Agent uses Ollama (free, local) for coordination
- Claude Pro reserved for complex code generation tasks
- Human-in-the-loop for critical decisions

---

## Quick Start

### 1. Enable Project Orchestrator

Edit `config.json` and set:

```json
{
  "project_orchestrator": {
    "enabled": true,
    "projects_dir": "./projects",
    "auto_transition": false,
    "require_human_approval": true,
    "lead_agent_model": "llama3:8b"
  }
}
```

**Note:** Requires `supervisor.enabled: true` to be set as well.

### 2. Start the Server

```bash
./orchestrator.exe -mode=server -port=8080
```

You should see:
```
✓ Supervisor enabled (complexity threshold: 7)
✓ ProjectOrchestrator enabled (projects dir: ./projects)
```

### 3. Access the Web UI

Open http://localhost:8080 and click the **Projects** tab.

---

## Architecture

### Components

**ProjectOrchestrator**
- Main orchestration component
- Manages project lifecycle and phase execution
- Wraps SupervisedTaskManager (backward compatible)
- File: `project/orchestrator.go`

**Lead Agent (Ollama llama3:8b)**
- Conservative, shipping-focused decision-making
- Delegates to 6 specialist agents
- Makes PROCEED/REFINE/BLOCK recommendations
- File: `project/lead_agent.go`

**Specialist Agents (All use Ollama)**
- RequirementsAgent - Validates completeness
- TechStackAgent - Approves technology choices
- ScopeAgent - Prevents scope creep
- QAAgent - Code quality review
- TestingAgent - Test generation
- DocumentationAgent - Docs generation

**CompletionValidator**
- Validates "hand-off ready" criteria
- Checks: Runnable build + Tests + README
- Calculates completion percentage
- File: `project/completion_validator.go`

**ProjectManager**
- File-based JSON persistence (`./projects/`)
- In-memory caching
- CRUD operations
- File: `project/manager.go`

### Data Flow

```
User Request
    ↓
ProjectOrchestrator
    ↓
Lead Agent (Ollama) → PROCEED/REFINE/BLOCK
    ↓
Specialist Agents (Ollama)
    ↓
Human Approval (Required)
    ↓
SupervisedTaskManager (Code Execution)
    ↓
Complexity Scoring → Ollama (free) or Claude Code (paid)
    ↓
Artifacts Saved → Phase Completed
```

---

## Project Lifecycle

### 8-Phase Workflow

#### 1. Discovery Phase
**Purpose:** Validate we understand what to build

**Activities:**
- Lead Agent invokes Requirements Agent
- Requirements Agent analyzes completeness (score 1-10)
- Lead Agent decides: PROCEED (≥7), REFINE (4-6), BLOCK (<4)

**Entry Criteria:** User provides name + description (min 50 chars)
**Exit Criteria:** Requirements score ≥7/10
**Human Approval:** Required
**Outputs:** Requirements analysis artifact

---

#### 2. Validation Phase
**Purpose:** Confirm tech stack and scope are appropriate

**Activities:**
- Lead Agent invokes TechStack Agent (parallel)
- Lead Agent invokes Scope Agent (parallel)
- Both agents provide APPROVED/WARNING/REJECTED verdicts

**Entry Criteria:** Discovery phase complete
**Exit Criteria:** Both agents APPROVED or WARNING (not REJECTED)
**Human Approval:** Required
**Outputs:** Tech stack analysis, scope validation

---

#### 3. Planning Phase
**Purpose:** Create implementation roadmap

**Activities:**
- Lead Agent generates implementation plan
- Defines milestones and technical approach
- Estimates complexity and identifies risks

**Entry Criteria:** Validation phase complete
**Exit Criteria:** Plan complete
**Human Approval:** Required
**Outputs:** Project plan artifact

---

#### 4. CodeGen Phase
**Purpose:** Generate working code

**Activities:**
- Complexity scorer evaluates (existing logic)
- Routes to Ollama (free, score <7) or Claude Code (paid, score ≥7)
- Generates code artifacts

**Entry Criteria:** Planning phase complete
**Exit Criteria:** Code generation successful
**Human Approval:** Not required (automated)
**Outputs:** Generated code artifacts

---

#### 5. Review Phase
**Purpose:** Validate code quality

**Activities:**
- Lead Agent invokes QA Agent (parallel)
- Lead Agent invokes Testing Agent (parallel)
- Consolidates feedback and decides PROCEED/REFINE/BLOCK

**Entry Criteria:** CodeGen phase complete
**Exit Criteria:** QA score ≥5/10, no critical bugs
**Human Approval:** Required if QA score <7/10
**Outputs:** QA review, test plan

---

#### 6. QA Phase
**Purpose:** Verify hand-off ready criteria

**Activities:**
- CompletionValidator checks:
  - ✓ Runnable build exists
  - ✓ Tests exist
  - ✓ README exists
- Calculates completion percentage

**Entry Criteria:** Review phase complete
**Exit Criteria:** All criteria met, completion % ≥80%
**Human Approval:** Required
**Outputs:** QA validation report

---

#### 7. Docs Phase
**Purpose:** Generate final documentation

**Activities:**
- Lead Agent invokes Documentation Agent
- Generates README with setup instructions
- Creates API documentation if applicable

**Entry Criteria:** QA phase complete
**Exit Criteria:** README complete
**Human Approval:** Not required (automated)
**Outputs:** README, API documentation

---

#### 8. Complete Phase
**Purpose:** Finalize project

**Activities:**
- Generate project summary
- Archive artifacts
- Mark status = COMPLETE

**Entry Criteria:** Docs phase complete
**Exit Criteria:** Completion % = 100%
**Human Approval:** Not required (automated)
**Outputs:** Project completion summary

---

## Lead Agent Decision Framework

### Three-State Model

**PROCEED** - All criteria met, safe to continue
- Requirements score ≥7/10
- TechStack + Scope both APPROVED
- QA score ≥7/10, no critical bugs
- Hand-off criteria all met

**REFINE** - Concerns present, needs iteration
- Requirements score 4-6/10
- TechStack or Scope has WARNINGS
- QA score 5-6/10, issues present
- 1-2 hand-off criteria missing

**BLOCK** - Critical issues, cannot proceed
- Requirements score <4/10
- TechStack REJECTED or Scope TOO_BROAD
- QA score <5/10, critical bugs
- All hand-off criteria missing

### Lead Agent Principles

1. **SHIPPING MATTERS** - Prefer working code over perfection
2. **SCOPE CONTROL** - Guard against feature creep aggressively
3. **QUALITY GATES** - Block on critical issues, warn on concerns
4. **DELEGATION** - Use specialist agents, don't do their jobs
5. **CONSERVATIVE** - Proven patterns over experimental approaches
6. **TRANSPARENCY** - Explain decisions clearly for human approval

---

## Hand-Off Ready Criteria

### Three Required Checks

**1. Runnable Build**
- Scans artifacts for entry point markers:
  - Go: `func main()`, `package main`
  - Python: `if __name__ == "__main__":`
  - JavaScript: `function main()`, `export default`
  - C#: `static void Main(`
  - Java: `public static void main(`
  - Rust: `fn main()`
  - C/C++: `int main(`

**2. Tests Exist**
- Scans artifacts for test markers:
  - Go: `func Test`
  - Python: `def test_`, `import pytest`
  - JavaScript: `it("`, `describe("`
  - C#: `[Test]`, `[Fact]`
  - Java: `@Test`
  - Rust: `#[test]`

**3. README Documentation**
- Scans artifacts for README markers:
  - `# Setup`, `## Installation`, `## Usage`
  - `### Prerequisites`, `How to run`, `Getting Started`
- Requires ≥2 markers found

### Completion Percentage

**Formula:**
```
Completion % = (Phase Weights) + (Criteria Bonus)

Phase Weights (total 80%):
- Discovery:  10%
- Validation: 10%
- Planning:   10%
- CodeGen:    30%  (biggest weight)
- Review:     20%
- QA:         10%
- Docs:       10%

Criteria Bonus (max 20%):
- Runnable build: +7%
- Tests:          +7%
- README:         +6%

Total: Up to 100%
```

---

## API Reference

### Create Project

**Endpoint:** `POST /project`

**Request:**
```json
{
  "name": "2D Roguelike Game",
  "description": "A procedurally generated dungeon crawler with turn-based combat and permadeath"
}
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "2D Roguelike Game",
  "description": "...",
  "current_phase": "discovery",
  "status": "active",
  "created_at": "2026-01-01T12:00:00Z",
  "phases": [...]
}
```

---

### List Projects

**Endpoint:** `GET /project/list`

**Response:**
```json
{
  "projects": [...],
  "count": 5
}
```

---

### Get Project

**Endpoint:** `GET /project?id={project_id}`

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "2D Roguelike Game",
  "current_phase": "validation",
  "phases": [...],
  "tasks": [...],
  "artifact_paths": [...],
  "metadata": {...}
}
```

---

### Execute Phase

**Endpoint:** `POST /project/phase`

**Request:**
```json
{
  "project_id": "550e8400-e29b-41d4-a716-446655440000",
  "phase": "discovery"
}
```

**Response:**
```json
{
  "phase": "discovery",
  "decision": "PROCEED",
  "reasoning": "Requirements are clear and complete with score 8/10",
  "next_steps": "Proceed to validation phase",
  "requires_approval": true,
  "recommended_action": "Approve and proceed to next phase"
}
```

---

### Approve Phase

**Endpoint:** `POST /project/approve`

**Request:**
```json
{
  "project_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Phase approved"
}
```

---

### Reject Phase

**Endpoint:** `POST /project/reject`

**Request:**
```json
{
  "project_id": "550e8400-e29b-41d4-a716-446655440000",
  "reason": "Requirements need more detail"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Phase rejected"
}
```

---

### Get Completion Metrics

**Endpoint:** `GET /project/metrics?project_id={project_id}`

**Response:**
```json
{
  "has_runnable_build": true,
  "has_tests": true,
  "has_readme": false,
  "completion_pct": 87.0,
  "blocking_issues": ["No README documentation detected"]
}
```

---

## Web UI Guide

### Creating a Project

1. Click **Projects** tab
2. Fill in project name and description
3. Click **Create Project**
4. Project appears in list below

### Managing Projects

1. Click on a project in the list
2. Project dashboard opens showing:
   - Current phase
   - Completion percentage
   - Status (Active/Blocked/Complete)
   - Phase progress visualization

### Executing Phases

1. Click **Execute Current Phase**
2. Lead Agent analyzes and displays recommendation:
   - **Decision:** PROCEED/REFINE/BLOCK
   - **Reasoning:** Why this decision
   - **Next Steps:** What to do next
3. Review recommendation
4. Click **Approve & Continue** or **Reject Phase**

### Monitoring Progress

**Phase Indicators:**
- **Green:** Phase completed
- **Blue:** Current phase
- **Gray:** Pending phase

**Completion %:**
- Calculated from completed phases + hand-off criteria
- Updates automatically as phases complete

**Hand-Off Ready Criteria (QA Phase):**
- ✅ or ❌ indicators for each criterion
- Visible when project reaches QA phase

---

## Troubleshooting

### Project Orchestrator Not Enabled

**Symptom:** Projects tab shows "Project orchestrator not enabled"

**Solution:**
1. Edit `config.json`
2. Set `project_orchestrator.enabled: true`
3. Ensure `supervisor.enabled: true`
4. Restart server

### Lead Agent Not Responding

**Symptom:** Phase execution hangs or times out

**Solution:**
1. Check Ollama is running: `ollama list`
2. Verify `llama3:8b` model is installed
3. Check server logs for errors
4. Ensure Ollama URL is correct in config

### Phase Stuck in "In Progress"

**Symptom:** Cannot execute or approve phase

**Solution:**
1. Check project status: `GET /project?id={id}`
2. Look for blocking issues in phase execution
3. May need to manually edit project JSON file
4. Or reject phase and retry

### Completion Metrics Not Showing

**Symptom:** Hand-off criteria show all ❌

**Solution:**
1. Ensure CodeGen phase has completed
2. Check artifact paths are correct
3. Verify artifacts contain required markers
4. Review `CompletionValidator` logic for language support

---

## Best Practices

### 1. Start Small

**Do:**
- Create focused, well-defined projects
- Use Discovery phase to clarify requirements
- Let Scope Agent validate feasibility

**Don't:**
- Create massive multi-feature projects
- Skip validation phases
- Ignore REFINE recommendations

### 2. Trust the Lead Agent

**Do:**
- Review Lead Agent reasoning carefully
- Use REFINE feedback to improve requirements
- Approve when decision is PROCEED

**Don't:**
- Blindly approve without reading
- Reject without understanding reasoning
- Force transitions when blocked

### 3. Leverage Specialist Agents

**Do:**
- Let TechStack Agent guide technology choices
- Use QA Agent feedback to improve code
- Review Testing Agent test plans

**Don't:**
- Override specialist recommendations without reason
- Ignore warnings from agents
- Skip quality gates

### 4. Monitor Completion Metrics

**Do:**
- Check hand-off criteria in QA phase
- Ensure all criteria met before completion
- Review blocking issues

**Don't:**
- Mark project complete with missing criteria
- Ignore test requirements
- Skip documentation

---

## Configuration Reference

### project_orchestrator Section

```json
{
  "project_orchestrator": {
    "enabled": false,              // Enable/disable project orchestrator
    "projects_dir": "./projects",  // Directory for project JSON files
    "auto_transition": false,      // Auto-transition phases (not recommended)
    "require_human_approval": true,// Require human approval for transitions
    "lead_agent_model": "llama3:8b" // Ollama model for Lead Agent
  }
}
```

**Parameters:**

- **enabled** (bool): Master switch for project orchestrator
- **projects_dir** (string): Where to save project JSON files
- **auto_transition** (bool): Automatically transition phases without approval (NOT RECOMMENDED - defeats human-in-loop)
- **require_human_approval** (bool): Require approval for phase transitions
- **lead_agent_model** (string): Ollama model name for Lead Agent (recommended: llama3:8b)

---

## File Structure

```
AI FACTORY/
├── config.json                    # Configuration
├── projects/                      # Project JSON files
│   ├── project_{uuid}.json
│   └── project_{uuid}.json
├── artifacts/                     # Generated artifacts
│   ├── code_{timestamp}.md
│   ├── discover_{timestamp}.md
│   └── ...
├── project/                       # Project orchestrator code
│   ├── project.go                 # Data models
│   ├── manager.go                 # Persistence
│   ├── lead_agent.go              # Lead Agent
│   ├── completion_validator.go    # Hand-off validation
│   └── orchestrator.go            # Main orchestrator
├── supervisor/                    # Specialist agents
│   ├── supervised_manager.go
│   ├── agent_requirements.go
│   ├── agent_techstack.go
│   ├── agent_scope.go
│   ├── agent_qa.go
│   ├── agent_testing.go
│   └── agent_documentation.go
└── web/
    └── index.html                 # Web UI with Projects tab
```

---

## Cost Optimization

### Ollama (Free, Local)

**Used for:**
- Lead Agent decisions (all phases)
- All specialist agents
- 80% of code generation tasks (complexity <7)

**Cost:** $0

### Claude Code (Paid, Optional)

**Used for:**
- Complex code generation (complexity ≥7)
- Multi-file projects
- Advanced features (auth, database, APIs)

**Cost:** Based on Claude Pro usage

**Estimated Split:**
- 80% tasks: Ollama (free)
- 20% tasks: Claude Code (paid)

### Cost Control Tips

1. **Adjust Complexity Threshold**
   - Default: 7
   - Increase to 8 or 9 to use more Ollama
   - Decrease to 5 or 6 to use more Claude

2. **Keep Projects Focused**
   - Smaller scope = lower complexity
   - Lower complexity = free execution

3. **Review Before Approval**
   - Check if CodeGen really needs Claude
   - Can reject and simplify requirements

---

## Support & Feedback

**Issues:** https://github.com/anthropics/claude-code/issues
**Documentation:** This guide + STATUS_SUMMARY.md
**Configuration:** config.json

---

## Summary

The AI Factory Project Orchestrator provides a structured, human-guided workflow from idea to "hand-off ready" completion:

✓ **8-Phase Lifecycle** - Discovery → Validation → Planning → CodeGen → Review → QA → Docs → Complete
✓ **Lead Agent Coordination** - Ollama-based, conservative, shipping-focused
✓ **Human-in-the-Loop** - Approval required for critical transitions
✓ **Hand-Off Ready Validation** - Runnable build + Tests + README
✓ **Cost Optimized** - Ollama leads (free), Claude escalates (paid)
✓ **Backward Compatible** - Disabled by default, opt-in

**Result:** Projects that are ready to hand off, not perfect - but shippable.
