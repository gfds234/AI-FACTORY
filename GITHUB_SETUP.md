# GitHub Setup Instructions

## Create Private Repository

Since GitHub CLI (`gh`) is not installed, follow these manual steps:

### Step 1: Create Repository on GitHub

1. Go to https://github.com/new
2. Fill in:
   - **Repository name:** `AI-FACTORY`
   - **Description:** `AI-assisted game development studio using local LLMs`
   - **Visibility:** ✓ **Private**
   - **DO NOT** initialize with README, .gitignore, or license (we already have these)
3. Click **"Create repository"**

### Step 2: Push to GitHub

After creating the repo, GitHub will show you commands. Use these:

```bash
git remote add origin https://github.com/YOUR_USERNAME/AI-FACTORY.git
git branch -M main
git push -u origin main
```

**Or if you prefer SSH:**

```bash
git remote add origin git@github.com:YOUR_USERNAME/AI-FACTORY.git
git branch -M main
git push -u origin main
```

Replace `YOUR_USERNAME` with your actual GitHub username.

### Step 3: Verify

After pushing, visit:
```
https://github.com/YOUR_USERNAME/AI-FACTORY
```

You should see all your files there.

---

## Quick Command Reference

**Already done:**
- ✓ Git initialized
- ✓ Files cleaned up
- ✓ .gitignore created
- ✓ Initial commit made

**What's committed:**
- All Go source code (api/, config/, llm/, task/, main.go)
- Web UI (web/index.html)
- Documentation (README.md, guides, status files)
- Launcher script (start-manager.bat)
- Example files (idea_space_station.txt, arch_strategy_game.txt)

**What's ignored (not in Git):**
- orchestrator binary
- config.json (runtime generated)
- artifacts/ (your results)
- .claude/ (Claude Code internal)

---

## Future Updates

After making changes, use:

```bash
# See what changed
git status

# Stage all changes
git add .

# Commit
git commit -m "Your commit message here"

# Push to GitHub
git push
```

---

## Alternative: Install GitHub CLI (Optional)

If you want easier GitHub integration:

```bash
# Install via winget (Windows)
winget install GitHub.cli

# Or via scoop
scoop install gh

# Then authenticate
gh auth login

# Then you can create repos with:
gh repo create AI-FACTORY --private --source=. --push
```

But manual method works fine too.
