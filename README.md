# AI FACTORY - Game Development Decision Amplifier

**Status:** Phase 1 Complete ✓ - Fully Operational
**Hardware Target:** RTX 4070 Super (12GB VRAM)
**Philosophy:** Human leverage over automation percentage
**Repository:** https://github.com/gfds234/AI-FACTORY

## Overview

AI-assisted studio for game development using local LLMs. NOT an autonomous factory - it's a human-centered decision amplifier that helps a solo game developer achieve 2x creative output per day.

## Project Structure

```
ai-studio/
├── orchestrator/     # Go-based orchestration layer
│   ├── api/         # HTTP API handlers
│   ├── task/        # Task management
│   └── llm/         # LLM client interfaces
├── models/          # Downloaded GGUF models (gitignored)
├── docs/            # Architecture decisions and runbooks
└── verify-env.sh    # Environment verification script
```

## Current Phase: Phase 1 - Foundation

**Goal:** Basic orchestration + 2 working tasks (idea validation + architecture review)

**Scope:**
- Sequential model execution only (RTX 4070 Super constraint)
- GGUF-format models via llama.cpp or ollama
- REST API for task submission
- File-based artifact storage
- Basic CLI for testing

## Environment Requirements

**Core:**
- Go 1.22+
- Python 3.12+
- Git

**Models (to download):**
- Planning: Mistral-7B-Instruct (4-bit GGUF, ~4GB VRAM)
- Implementation: DeepSeek-Coder-6.7B (4-bit GGUF, ~4GB VRAM)
- Review: Llama-3-8B (4-bit GGUF, ~5GB VRAM)

**LLM Runtime (choose one):**
- llama.cpp (direct, more control)
- Ollama (simpler, abstracts model management)

## Quick Start (You've Already Done This ✓)

```powershell
# 1. Models installed ✓
ollama list
# Shows: mistral:7b-instruct-v0.2-q4_K_M, llama3:8b, deepseek-coder:6.7b-instruct

# 2. Build orchestrator
cd orchestrator
go build -o orchestrator.exe .

# 3. Run tests
cd ..
.\test.ps1
```

**That's it.** The orchestrator is ready to route tasks to your local models.

## Design Decisions

1. **Sequential execution:** No parallel models (hardware constraint)
2. **Local-first:** All models run on RTX 4070 Super
3. **File artifacts:** Easier to inspect/version than database
4. **Go orchestrator:** Performance + simplicity for task routing
5. **No enterprise tooling:** DataDog, K8s, etc. are overkill

## Risk Register

See project specification document for full risk analysis.

**Top risks:**
- VRAM exhaustion from model switching
- Model quality variance across specialized tasks
- Latency accumulation in sequential pipeline

## Phases

- **Phase 1 (Days 1-10):** Foundation - orchestration + 2 tasks
- **Phase 2 (Days 11-20):** Expansion - 2 more tasks + monitoring
- **Phase 3 (Days 21-30):** Integration - workflow templates + metrics

## License

Private project - not open source.

---

*Last updated: 2025-12-31*
