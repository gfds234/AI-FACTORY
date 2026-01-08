# AI Studio - Manager Quick Start

## How to Use (Non-Technical)

### Option 1: Double-Click to Start (Easiest)

1. **Double-click** `start-manager.bat`
2. Web interface opens automatically in your browser
3. Type or paste your game idea or architecture document
4. Click the button
5. Wait 4-50 seconds for AI analysis
6. Results appear on screen

**To stop:** Close the "AI Studio Server" window.

---

### Option 2: Manual Start

1. Open terminal in this folder
2. Run: `orchestrator -mode=server -port=8080`
3. Open `http://localhost:8080` in your browser
4. Use the interface

---

## What Each Task Does

### Validate Game Idea
**When to use:** You have a game concept and want to know if it's worth pursuing.

**You provide:** Plain text description of your game idea (mechanics, target audience, platform, etc.)

**You get back:**
- Core concept summary
- Strengths (2-3 points)
- Potential issues (2-3 concerns)
- Market viability assessment
- Recommendation (Proceed/Revise/Reconsider)
- Next steps (2-3 concrete actions)

**Model used:** Mistral 7B Instruct

---

### Review Architecture
**When to use:** You have technical architecture decisions and want expert feedback.

**You provide:** Architecture document (text or markdown) with technical details

**You get back:**
- Architecture summary
- Strengths
- Risk assessment (performance, scalability, technical risks)
- Recommendations
- Implementation notes
- Verdict (Approved/Needs Revision/Major Concerns)

**Model used:** Llama3 8B

---

## Performance Expectations

**First request of the day:** 30-50 seconds (loading AI model into GPU memory)
**Subsequent requests:** 4-5 seconds (model already loaded)
**Switching between validate/review:** 10-15 seconds (loading different model)

---

## Where Results Are Saved

All results automatically save to: `artifacts/`

File naming: `[task-type]_[timestamp].md`

Examples:
- `validate_1767149923.md`
- `review_1767149984.md`

These are markdown files - open with any text editor or markdown viewer.

---

## Troubleshooting

### "Server Offline" message in web UI

**Fix:**
1. Make sure you ran `start-manager.bat` or started the server manually
2. Check if Ollama is running (open new terminal, type: `ollama list`)
3. If Ollama isn't running, start it

### "Cannot connect" error

**Fix:**
1. Close the server (if running)
2. Re-run `start-manager.bat`

### Request takes forever

**Normal:** First request can take 50 seconds. This is the AI loading into GPU memory.
**If longer than 2 minutes:** Something's wrong. Check terminal for errors.

### Web page won't open

**Fix:**
1. Make sure the server is running (check for "AI Studio Server" window)
2. Manually open `http://localhost:8080` in Chrome, Firefox, or Edge
3. If you see "Server Offline", the server didn't start - check Ollama is running

---

## Daily Workflow

1. **Morning:** Run `start-manager.bat` once
2. **Throughout day:** Use web interface as needed
3. **End of day:** Close terminal window
4. **Review results:** Check `artifacts/` folder for all saved analyses

---

## Files You Care About

- `start-manager.bat` - Double-click this to start everything
- `web\index.html` - The web interface (opens automatically)
- `artifacts\` - All your results live here
- `config.json` - Settings (rarely need to touch)

---

## Tips

1. **Save good prompts:** If you get great results, save the input text for reuse
2. **Iterate:** Use validate → improve idea → validate again workflow
3. **Compare:** Run multiple ideas through validate to compare
4. **Chain tasks:** Validate idea first, then review the architecture you design
5. **Artifacts are yours:** Copy/paste from artifacts into design docs

---

## Advanced Features (Optional)

### Multi-Agent Supervisor System

**Status:** Implemented but disabled by default (Phase 3)

**What it does:**
- Automatically routes simple tasks to free AI (Ollama)
- Routes complex tasks to premium AI (Claude Code) - if configured
- Adds quality checks before running tasks
- Automatically generates test plans and documentation
- Reviews code quality and finds bugs

**Why disabled by default:**
- Adds 15-90 seconds of processing time
- Requires additional configuration
- Most users don't need it initially

**How to enable:**
See [SUPERVISOR_GUIDE.md](SUPERVISOR_GUIDE.md) for complete setup instructions.

**Recommended for:**
- Production code generation
- Teams needing quality gates
- Projects requiring documentation
- Cost optimization (80% free AI usage)

**Not recommended for:**
- Quick idea validation
- Rapid prototyping
- Learning the system
### Project-Based Workflow (Phase 4)
The AI Factory now follows a structured **8-Phase Lifecycle** for full MVP generation:
1. **Discovery** (Idea refinement)
2. **Validation** (Market fit/Tech stack)
3. **Planning** (Roadmap generation)
4. **CodeGen** (Automatic code generation) -> *Now with robust file parsing and context injection!*
5. **Review** (Quality verification)
6. **QA** (Success guarantee)
7. **Docs** (Manuals and API docs)
8. **Complete** (Final hand-off)

### Robustness & Reliability
- **Triple Guarantee System**: Every MVP is verified for Build, Runtime, and Test success.
- **Failover Logic**: If premium agents (Claude) are offline, the system automatically falls back to local models (Ollama).
- **Data Protection**: Automatic schema validation ensures project files never corrupt silently.

### Tips for Best Results
1. **Be Descriptive**: The more detail you provide in the Planning phase, the better the CodeGen results.
2. **Review the Roadmap**: Always check the Planning output before approving CodeGen.
3. **Use the "X"**: The welcome modal now has a close button if you need to skip the intro quickly.
