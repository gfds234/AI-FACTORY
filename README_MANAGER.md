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
