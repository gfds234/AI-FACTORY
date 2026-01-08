# 🚀 AI FACTORY - Quickstart Guide

Get AI FACTORY running in **under 2 minutes** with automated setup!

## One-Click Launch

### Windows

**Double-click:** `QUICKSTART.bat`

That's it! The script will:
1. ✅ Check for required files
2. ✅ Prompt for API key (if needed)
3. ✅ Start the orchestrator server
4. ✅ Launch ngrok tunnel (if available)
5. ✅ Open your browser to the web UI

---

## First Time Setup

### Prerequisites

Before running QUICKSTART.bat, ensure you have:

1. **Built the project:**
   ```bash
   go build -o orchestrator.exe .
   ```

2. **Anthropic API Key** - Get one from https://console.anthropic.com

3. **(Optional) ngrok** - Download from https://ngrok.com/download
   - Extract ngrok.exe to `C:\Users\lampr\Downloads\ngrok\` or `C:\ngrok\`
   - Sign up at https://ngrok.com for an authtoken

---

## What QUICKSTART.bat Does

### Automatic Detection
- ✅ Finds your API key (from environment or .env)
- ✅ Locates ngrok installation
- ✅ Verifies orchestrator binary exists
- ✅ Checks if services are running

### Smart Configuration
- ✅ Prompts for missing API key
- ✅ Saves API key to `.env` for future use
- ✅ Configures ngrok authtoken (if needed)
- ✅ Falls back to local-only mode if ngrok unavailable

### Automated Startup
- ✅ Starts orchestrator on port 8080
- ✅ Waits for health check confirmation
- ✅ Launches ngrok tunnel with public URL
- ✅ Retrieves and displays ngrok URL
- ✅ Opens browser automatically

---

## Access URLs

After running QUICKSTART.bat, you'll see:

**Local Access (always available):**
```
http://localhost:8080
```

**Public Access (if ngrok is configured):**
```
https://your-subdomain.ngrok-free.dev
```

---

## Troubleshooting

### "orchestrator.exe not found"
**Solution:** Build the project first
```bash
go build -o orchestrator.exe .
```

### "ngrok not found"
**Solution:** Either:
1. Download ngrok and place in `C:\ngrok\`
2. Or skip ngrok and use local access only

### "Orchestrator failed to start"
**Causes:**
- Port 8080 already in use
- Missing ANTHROPIC_API_KEY
- Ollama not running (if using local models)

**Solution:**
- Check the "AI FACTORY Orchestrator" window for error messages
- Verify API key is correct
- Stop other services using port 8080

### ngrok shows "authentication failed"
**Solution:** Get authtoken from https://dashboard.ngrok.com/get-started/your-authtoken
```bash
ngrok config add-authtoken YOUR_TOKEN_HERE
```

---

## Manual Setup (Alternative)

If you prefer manual control:

### 1. Set API Key
```bash
# Option A: Environment variable
set ANTHROPIC_API_KEY=sk-ant-xxxxx

# Option B: Create .env file
echo ANTHROPIC_API_KEY=sk-ant-xxxxx > .env
```

### 2. Start Orchestrator
```bash
orchestrator.exe
```

### 3. Start ngrok (Optional)
```bash
C:\Users\lampr\Downloads\ngrok\ngrok.exe http 8080
```

### 4. Open Browser
```
http://localhost:8080
```

---

## What's Next?

Once AI FACTORY is running:

1. **Go to "Projects" tab**
2. **Click "New Project"**
3. **Enter project details:**
   - Name: "My First MVP"
   - Description: "A simple todo app with React"
4. **Click "Create Project"**
5. **Watch the autonomous 6-phase generation!**

---

## Sharing with Partners

If you started ngrok, share the public URL:

```
Hi! Check out AI FACTORY:

🔗 https://your-subdomain.ngrok-free.dev

What to try:
- Go to "Projects" tab
- Create a new project
- Toggle between List and Kanban views
- Watch real-time phase updates
- Download completed projects as ZIP

Features:
✅ Autonomous MVP generation
✅ Triple Guarantee System (Build + Runtime + Test)
✅ Quality scoring and certificates
✅ 6-phase workflow visualization
```

---

## Stopping AI FACTORY

To stop all services:

1. Close the ngrok window (or press Ctrl+C)
2. Close the orchestrator window (or press Ctrl+C)
3. Close the quickstart window

Or run:
```bash
taskkill /F /IM orchestrator.exe
taskkill /F /IM ngrok.exe
```

---

## Advanced Configuration

### Custom Port
Edit `QUICKSTART.bat` or use flag:
```bash
orchestrator.exe -port=3000
```

### Environment Variables

Create `.env` file with:
```env
ANTHROPIC_API_KEY=sk-ant-xxxxx
OLLAMA_BASE_URL=http://localhost:11434
PORT=8080
PROJECTS_DIR=./projects
ARTIFACTS_DIR=./artifacts
```

### Config File

Create `config.json`:
```json
{
  "ollama_url": "http://localhost:11434",
  "models": {
    "validate": "mistral:7b-instruct-v0.2-q4_K_M",
    "review": "llama3:8b",
    "code": "deepseek-coder:6.7b-instruct"
  },
  "artifacts_dir": "./artifacts",
  "max_retries": 2,
  "timeout_seconds": 120,
  "project_orchestrator": {
    "enabled": true,
    "projects_dir": "./projects",
    "auto_transition": false,
    "require_human_approval": true,
    "lead_agent_model": "llama3:8b"
  }
}
```

---

## Support

**Issues?**
- Check orchestrator logs in the terminal window
- Visit: http://localhost:4040 for ngrok dashboard
- Review `PROGRESS_LOG.md` for recent changes
- Check `STATUS_SUMMARY.md` for feature overview

**Documentation:**
- [README.md](README.md) - Main documentation
- [PROJECT_ORCHESTRATOR_GUIDE.md](PROJECT_ORCHESTRATOR_GUIDE.md) - Phase details
- [SUPERVISOR_GUIDE.md](SUPERVISOR_GUIDE.md) - Multi-agent system
- [PROGRESS_LOG.md](PROGRESS_LOG.md) - Technical changelog

---

**Built with ❤️ using Claude Sonnet 4.5**
