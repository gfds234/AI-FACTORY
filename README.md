# AI FACTORY - Autonomous MVP Generation with Quality Guarantees

**Status: Phase 4+ Complete** ✓ (Triple Guarantee System Implemented)

AI FACTORY is an autonomous system for generating production-ready MVPs with automated quality verification. Unlike competitors that only generate code, AI FACTORY **proves every deliverable works** through automated Build + Runtime + Test verification.

## Quick Links

- **[Status Summary](STATUS_SUMMARY.md)** - Current features and capabilities
- **[Project Orchestrator Guide](PROJECT_ORCHESTRATOR_GUIDE.md)** - Phase 4 implementation details
- **[Supervisor Guide](SUPERVISOR_GUIDE.md)** - Multi-agent system architecture
- **[Progress Log](PROGRESS_LOG.md)** - Detailed technical changelog
- **[User Guide](README_MANAGER.md)** - Non-technical daily usage

## Triple Guarantee System™

AI FACTORY is the **only** MVP generation tool that provides automated quality guarantees:

1. **Build Guarantee** ✅
   - Syntax validation
   - Dependency resolution
   - Entry point verification

2. **Runtime Guarantee** ✅
   - Application startup verification
   - Health check validation
   - Port binding confirmation

3. **Test Guarantee** ✅
   - Test execution with framework detection (Jest, pytest, go test)
   - Pass/fail count reporting
   - Test coverage tracking

**Every project includes:**
- Quality Score (0-100) based on all verification phases
- Professional Quality Certificate markdown report
- Visual dashboard showing verification status
- API access to quality data (`GET /project/quality`)

**Business Impact:**
- Justify premium pricing ($800-2500/MVP vs $20/month competitors)
- "Money-back if code doesn't work" guarantee
- Client deliverables include verification proof

## Overview

AI FACTORY autonomously generates complete MVPs through a 6-phase process: Planning → Code Generation → QA → Documentation → Finalization → Completion. Each phase uses specialized agents coordinated by a multi-agent supervisor system.

## Project Structure

```
AI FACTORY/
├── api/              # REST API server
├── project/          # Project orchestrator and phase management
├── supervisor/       # Multi-agent coordination system
├── validation/       # Triple Guarantee verification
├── task/             # Task execution and routing
├── llm/              # Claude API and Ollama LLM clients
├── web/              # Web UI dashboard
├── artifacts/        # Generated project artifacts
└── projects/         # Generated MVP projects
```

## Key Features

- **Autonomous MVP Generation**: 6-phase process from idea to production-ready code
- **Multi-Agent Coordination**: Specialized agents for planning, coding, QA, and documentation
- **Triple Guarantee Verification**: Automated build, runtime, and test validation
- **Web Dashboard**: Visual project management and quality reporting
- **Quality Certificates**: Professional markdown reports for client deliverables
- **Project Export**: Comprehensive documentation bundles
- **Phase Reversion**: Go back to previous phases while preserving data

## Environment Requirements

**Core:**
- Go 1.22+
- Git

**LLM Providers:**
- Claude API (Anthropic) - Primary agent for project generation
- Ollama (optional) - Local models for task routing

**Language Runtimes (for validation):**
- Node.js (for JavaScript/TypeScript projects)
- Python 3.10+ (for Python projects)
- Go 1.22+ (for Go projects)

## Quick Start

1. **Clone and build:**
```bash
git clone https://github.com/gfds234/AI-FACTORY
cd AI-FACTORY
go build -o orchestrator.exe .
```

2. **Configure environment:**
```bash
# Create .env file
ANTHROPIC_API_KEY=your_claude_api_key_here
OLLAMA_BASE_URL=http://localhost:11434
```

3. **Run the system:**
```bash
./orchestrator.exe
# Opens web UI at http://localhost:8080
```

4. **Create your first MVP:**
- Open web dashboard
- Click "New Project"
- Enter project idea
- Watch autonomous generation with real-time phase updates
- Get quality-verified deliverable

## Architecture Highlights

1. **Phase-Based Generation**: Structured 6-phase workflow ensures quality at each step
2. **Multi-Agent Supervision**: Supervisor coordinates specialized agents (Planning, Code, QA, Docs)
3. **File-Based Artifacts**: Project files stored in structured directories for easy inspection
4. **Triple Guarantee Verification**: Automated validation proves code works before delivery
5. **Progressive Web UI**: Real-time phase updates with quality dashboards
6. **REST API**: Full programmatic access to all features

## Implementation Phases

- **Phase 1:** Foundation - Basic orchestration and task routing ✅
- **Phase 2:** Task expansion - Additional task types and monitoring ✅
- **Phase 3:** Multi-agent supervision - Coordinated specialist agents ✅
- **Phase 4:** Project orchestration - 6-phase MVP generation ✅
- **Phase 4+:** Triple Guarantee System - Automated quality verification ✅

**Current Focus:** Client acquisition and demo preparation

## Business Model

**Target Market:** Agencies, freelancers, and technical founders needing rapid MVP development

**Pricing Strategy:**
- $800-2500 per MVP (vs competitors at $20/month with no guarantees)
- Premium justified by Triple Guarantee System
- "Money-back if code doesn't work" guarantee differentiates from competition

**Competitive Advantage:**
- Only tool with automated quality verification
- Professional quality certificates included
- Visual proof that deliverables actually work

## License

Private project - not open source.

---

*Last updated: 2026-01-07*
