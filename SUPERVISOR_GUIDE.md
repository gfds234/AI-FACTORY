# Multi-Agent Supervisor System - User Guide

## Overview

The AI FACTORY now includes an intelligent **Supervisor System** that manages quality gates, routes tasks based on complexity, and enriches outputs with automated QA, testing, and documentation.

**Key Benefits:**
- 💰 **Cost Savings** - Route simple tasks to free Ollama, complex tasks to Claude Code
- ✅ **Quality Control** - Enforce requirements, tech stack, and scope validation before execution
- 📚 **Auto-Documentation** - Generate README, tests, and API docs automatically
- 🎯 **Smart Routing** - AI analyzes complexity and chooses optimal execution path

---

## Quick Start

### Option 1: Use Without Supervisor (Default)

The system works exactly as before when supervisor is disabled:

```json
// config.json (no supervisor section needed)
{
  "ollama_url": "http://localhost:11434",
  "models": {
    "code": "deepseek-coder:6.7b-instruct",
    "review": "llama3:8b",
    "validate": "mistral:7b-instruct-v0.2-q4_K_M"
  },
  "artifacts_dir": "./artifacts",
  "max_retries": 2,
  "timeout_seconds": 120
}
```

**Start server:**
```bash
./orchestrator.exe -mode=server -port=8080
```

**Output:**
```
Supervisor disabled - using standard task manager
Starting AI Studio Orchestrator on port 8080
```

---

### Option 2: Enable Smart Routing Only

Add complexity-based routing without quality gates or agents:

```json
{
  "supervisor": {
    "enabled": true,
    "quality_gates": {
      "requirements_check": false,
      "tech_stack_approval": false,
      "scope_validation": false
    },
    "agents": {
      "qa": {"enabled": false},
      "testing": {"enabled": false},
      "documentation": {"enabled": false}
    },
    "complexity_threshold": 7,
    "claude_code_endpoint": ""
  }
}
```

**What this does:**
- Analyzes every code task for complexity (score 1-10)
- Routes simple tasks (score < 7) to Ollama (free, fast)
- Routes complex tasks (score >= 7) to Claude Code (if endpoint configured)
- No quality gates, no post-processing

**Use case:** Save costs by using free Ollama for simple tasks

---

### Option 3: Enable Quality Gates

Validate requests before execution:

```json
{
  "supervisor": {
    "enabled": true,
    "quality_gates": {
      "requirements_check": true,
      "tech_stack_approval": true,
      "scope_validation": true
    },
    "agents": {
      "requirements": {"enabled": true, "model": "mistral:7b-instruct-v0.2-q4_K_M"},
      "qa": {"enabled": false},
      "testing": {"enabled": false},
      "documentation": {"enabled": false}
    }
  }
}
```

**What this does:**
1. **Requirements Agent** - Checks if request has enough detail (COMPLETE/INCOMPLETE/NEEDS_CLARIFICATION)
2. **Tech Stack Agent** - Pre-approves technology choices for code tasks (APPROVED/REJECTED)
3. **Scope Agent** - Validates project size is appropriate (APPROPRIATE/TOO_BROAD)

**Use case:** Prevent bad requests from wasting time/tokens

---

### Option 4: Enable Post-Processing Agents

Enrich outputs with QA, tests, and docs:

```json
{
  "supervisor": {
    "enabled": true,
    "quality_gates": {
      "requirements_check": false,
      "tech_stack_approval": false,
      "scope_validation": false
    },
    "agents": {
      "qa": {"enabled": true, "model": "llama3:8b"},
      "testing": {"enabled": true, "model": "deepseek-coder:6.7b-instruct"},
      "documentation": {"enabled": true, "model": "mistral:7b-instruct-v0.2-q4_K_M"}
    }
  }
}
```

**What this does:**
1. **QA Agent** - Reviews code for bugs, security issues, best practices
2. **Testing Agent** - Generates unit tests and test plans
3. **Documentation Agent** - Creates README, API docs, setup guides

**Use case:** Get comprehensive deliverables automatically

---

### Option 5: Full Supervision (Production)

Enable all features:

```json
{
  "supervisor": {
    "enabled": true,
    "quality_gates": {
      "requirements_check": true,
      "tech_stack_approval": true,
      "scope_validation": true
    },
    "agents": {
      "requirements": {"enabled": true, "model": "mistral:7b-instruct-v0.2-q4_K_M"},
      "qa": {"enabled": true, "model": "llama3:8b"},
      "testing": {"enabled": true, "model": "deepseek-coder:6.7b-instruct"},
      "documentation": {"enabled": true, "model": "mistral:7b-instruct-v0.2-q4_K_M"}
    },
    "complexity_threshold": 7,
    "claude_code_endpoint": "http://localhost:8081/generate"
  }
}
```

**Complete workflow:**
```
User Request
  ↓
Requirements Agent → Validate completeness
  ↓
Tech Stack Agent → Pre-approve tech choices (code tasks)
  ↓
Scope Agent → Validate project size
  ↓
Complexity Scorer → Analyze indicators (1-10)
  ↓
Route Decision → Ollama (score < 7) or Claude Code (score >= 7)
  ↓
Execute Task → Generate code/validation/review
  ↓
QA Agent → Review quality, find bugs
  ↓
Testing Agent → Generate test plan
  ↓
Documentation Agent → Create README & docs
  ↓
Return Enriched Result
```

**Use case:** Maximum quality and automation for production use

---

## Complexity Scoring

The supervisor analyzes code generation requests using 8 indicators:

| Indicator | Weight | Examples |
|-----------|--------|----------|
| Multi-file project | +3 | "microservice architecture", "project structure" |
| Database/ORM | +2 | "PostgreSQL", "database", "SQL", "migration" |
| Authentication | +2 | "JWT", "OAuth", "login", "session" |
| External APIs | +2 | "REST API", "webhook", "third-party integration" |
| Advanced algorithms | +2 | "machine learning", "pathfinding", "optimization" |
| Real-time features | +1 | "WebSocket", "streaming", "live updates" |
| Testing requirements | +1 | "unit tests", "test coverage", "e2e tests" |
| Large input | +1 | > 200 words |

**Scoring Examples:**

| Request | Score | Route |
|---------|-------|-------|
| "Python Hello World script" | 1 | Ollama |
| "Todo app with React and localStorage" | 5 | Ollama |
| "Microservice with PostgreSQL, JWT auth, tests" | 10 | Claude Code |

**Threshold Recommendations:**
- **3** - Aggressive escalation (most code → Claude Code)
- **5** - Balanced (medium+ → Claude Code)
- **7** - Conservative (only truly complex → Claude Code) ⭐ **Recommended**
- **9** - Minimal (almost everything → Ollama)

---

## Configuration Reference

### supervisor.enabled
**Type:** boolean
**Default:** false
**Description:** Master switch for supervisor system

### supervisor.quality_gates.requirements_check
**Type:** boolean
**Default:** false
**Description:** Validate completeness before execution
**Agent:** Requirements Agent (mistral:7b)
**Duration:** ~5-15 seconds
**Blocks execution:** Yes (if status = INCOMPLETE)

### supervisor.quality_gates.tech_stack_approval
**Type:** boolean
**Default:** false
**Description:** Pre-approve technology choices (code tasks only)
**Agent:** Tech Stack Agent (llama3:8b)
**Duration:** ~5-15 seconds
**Blocks execution:** Yes (if status = REJECTED)

### supervisor.quality_gates.scope_validation
**Type:** boolean
**Default:** false
**Description:** Validate project scope appropriateness
**Agent:** Scope Agent (mistral:7b)
**Duration:** ~5-15 seconds
**Blocks execution:** Yes (if status = TOO_BROAD)

### supervisor.agents.qa.enabled
**Type:** boolean
**Default:** false
**Description:** Review code quality, bugs, security
**Agent:** QA Agent (llama3:8b)
**Duration:** ~10-30 seconds
**Blocks execution:** No (post-processing)

### supervisor.agents.testing.enabled
**Type:** boolean
**Default:** false
**Description:** Generate unit tests and test plans
**Agent:** Testing Agent (deepseek-coder:6.7b)
**Duration:** ~10-30 seconds
**Blocks execution:** No (post-processing)

### supervisor.agents.documentation.enabled
**Type:** boolean
**Default:** false
**Description:** Create README, API docs, setup guides
**Agent:** Documentation Agent (mistral:7b)
**Duration:** ~10-30 seconds
**Blocks execution:** No (post-processing)

### supervisor.complexity_threshold
**Type:** integer (1-10)
**Default:** 7
**Description:** Score threshold for Claude Code escalation

### supervisor.claude_code_endpoint
**Type:** string
**Default:** ""
**Description:** HTTP endpoint for Claude Code (e.g., "http://localhost:8081/generate")
**Required:** Only if you want to escalate complex tasks to Claude Code

---

## Performance Impact

### Supervisor Disabled (Default)
- **Overhead:** 0ms
- **Behavior:** Identical to Phase 2

### Supervisor Enabled (Complexity Scoring Only)
- **Overhead:** <1ms (in-memory analysis)
- **Behavior:** Smart routing, no agent calls

### Quality Gates Enabled (All 3)
- **Overhead:** 15-45 seconds total (3 agents × 5-15s each)
- **Behavior:** Validates before execution, may block

### Post-Processing Agents (All 3)
- **Overhead:** 30-90 seconds total (3 agents × 10-30s each)
- **Behavior:** Runs after task completion, doesn't block

### Full Supervision
- **Overhead:** 45-135 seconds total
- **Behavior:** Complete quality pipeline

**Performance Tips:**
- Start with complexity scoring only (minimal overhead)
- Enable agents selectively based on task type
- Use quality gates for untrusted inputs only
- Post-processing agents are non-blocking (safe to enable)

---

## Error Handling

### Quality Gate Failures

**Requirements Agent returns INCOMPLETE:**
```
Error: requirements incomplete - cannot proceed
```
**Action:** User should provide more detail

**Tech Stack Agent returns REJECTED:**
```
Error: tech stack rejected - cannot proceed
```
**Action:** User should revise technology choices

**Scope Agent returns TOO_BROAD:**
```
Error: scope too broad - break into smaller tasks
```
**Action:** User should split into multiple tasks

### Agent Failures

**Post-execution agents never block:**
- If QA Agent fails → task still succeeds, no QA review returned
- If Testing Agent fails → task still succeeds, no test plan returned
- If Documentation Agent fails → task still succeeds, no docs returned

**Logs show failures:**
```
QA agent failed: connection timeout
Testing agent failed: model not loaded
Documentation agent failed: invalid response
```

### Claude Code Unavailable

**If claude_code_endpoint configured but unreachable:**
```
Warning: Claude Code not accessible: connection refused
Falling back to Ollama for all tasks
```
**Behavior:** All tasks route to Ollama regardless of complexity score

---

## Troubleshooting

### Supervisor not activating

**Check config.json:**
```json
{
  "supervisor": {
    "enabled": true  // ← Must be true
  }
}
```

**Check server logs:**
```
✓ Supervisor enabled (complexity threshold: 7)  ← Should see this
```

### Quality gates always passing

**Likely causes:**
- Agent models not loaded in Ollama
- Agent status parsing failed (check agent output in artifacts)

**Debug:**
1. Check `artifacts/` folder for agent outputs
2. Look for status keywords: COMPLETE, APPROVED, APPROPRIATE
3. Verify Ollama has required models loaded

### Tasks routing incorrectly

**Check complexity score in logs:**
```
Complexity score: 5, route: ollama  ← Expected for simple tasks
Complexity score: 10, route: claude_code  ← Expected for complex tasks
```

**Adjust threshold if needed:**
```json
{
  "supervisor": {
    "complexity_threshold": 5  // Lower = more tasks to Claude Code
  }
}
```

### Slow performance

**Disable unnecessary agents:**
```json
{
  "supervisor": {
    "quality_gates": {
      "requirements_check": false,  // Disable if trusted input
      "tech_stack_approval": false,  // Disable for non-code tasks
      "scope_validation": false
    }
  }
}
```

---

## Configuration Presets

Use these preset files for common scenarios:

### config.example.supervisor.json
Full supervision with all features (for reference)

### To apply a preset:
```bash
cp config.example.supervisor.json config.json
# Edit config.json to customize
./orchestrator.exe -mode=server -port=8080
```

---

## Best Practices

### 1. Start Simple
- Begin with supervisor disabled
- Enable complexity scoring first
- Add agents incrementally

### 2. Monitor Performance
- Check `artifacts/` for agent outputs
- Review server logs for timing
- Adjust thresholds based on actual usage

### 3. Use Quality Gates Selectively
- Enable for user-facing APIs (untrusted input)
- Disable for internal/trusted workflows
- Consider overhead vs. benefit

### 4. Leverage Post-Processing
- QA Agent: Always useful for code generation
- Testing Agent: Valuable for production code
- Documentation Agent: Great for open-source projects

### 5. Optimize Costs
- Set complexity_threshold high (7+) to minimize Claude Code usage
- Use Ollama for most tasks (free, local)
- Reserve Claude Code for truly complex projects

---

## API Changes

### Request (Unchanged)
```json
POST /task
{
  "task_type": "code",
  "input": "Build a microservice..."
}
```

### Response (Enhanced if Supervisor Enabled)
```json
{
  "task_type": "code",
  "input": "Build a microservice...",
  "output": "Generated code here...",
  "model": "deepseek-coder:6.7b-instruct",
  "duration_seconds": 12.5,

  // New fields when supervisor enabled:
  "complexity_score": 10,
  "execution_route": "claude_code",
  "requirements_analysis": { "status": "passed", "output": "..." },
  "tech_stack_approval": { "status": "passed", "output": "..." },
  "scope_validation": { "status": "passed", "output": "..." },
  "qa_review": { "status": "passed", "output": "..." },
  "test_plan": { "status": "passed", "output": "..." },
  "documentation": { "status": "passed", "output": "..." },
  "total_duration_seconds": 45.2,
  "agent_durations": {
    "requirements": 5.1,
    "techstack": 4.8,
    "scope": 5.2,
    "qa": 12.3,
    "testing": 10.5,
    "documentation": 7.3
  }
}
```

**Note:** Extra fields only present when corresponding agents are enabled

---

## Migration Guide

### From Phase 2 to Supervised

**No changes required!** The system is backward compatible.

**To enable:**
1. Add supervisor section to config.json
2. Restart server
3. Verify logs show "✓ Supervisor enabled"

**Rollback:**
1. Remove supervisor section from config.json (or set `"enabled": false`)
2. Restart server
3. System reverts to Phase 2 behavior

**Zero downtime migration:**
1. Deploy new orchestrator.exe
2. Keep config.json unchanged (supervisor disabled)
3. Test in production
4. Enable supervisor when ready

---

## Support

**Issues:** Check server logs for detailed error messages
**Performance:** Review `artifacts/` folder for agent outputs
**Questions:** See PROGRESS_LOG.md for implementation details

**Common Issues:**
- Port conflicts → Use different port (`-port=8081`)
- Ollama not running → Start Ollama service
- Models not loaded → Run `ollama pull <model-name>`
- Config errors → Validate JSON syntax

---

## Summary

The supervisor system provides:
- ✅ Smart routing based on complexity
- ✅ Quality gates to prevent bad requests
- ✅ Automated QA, testing, and documentation
- ✅ Cost optimization (Ollama vs Claude Code)
- ✅ Zero breaking changes (fully backward compatible)
- ✅ Granular control (enable/disable individual components)

**Start using it today** by adding the supervisor section to your config.json!
