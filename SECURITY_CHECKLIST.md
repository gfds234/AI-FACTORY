# Security Checklist for AI FACTORY

## Before Every GitHub Push

### Environment Files
- [ ] `.env` file is listed in `.gitignore`
- [ ] `.env.example` contains NO real API keys (only placeholders)
- [ ] No `.env` files committed to Git history

### API Keys & Secrets
- [ ] No `ANTHROPIC_API_KEY` hardcoded in source files
- [ ] No ngrok authtokens in source files
- [ ] No database credentials in source files
- [ ] All secrets loaded from environment variables

### Sensitive Data
- [ ] No production URLs or IPs hardcoded
- [ ] No user data or test data with PII
- [ ] No internal documentation paths exposed

## Files That Should NEVER Be Committed

```
.env
.env.local
.env.production
config.json (may contain secrets)
artifacts/ (may contain generated code with keys)
projects/ (may contain user data)
*.pem
*.key
*_credentials.json
```

## Safe to Commit

```
.env.example
.gitignore
*.md (documentation)
Source code (*.go, *.js, *.html)
QUICKSTART.bat
```

## Verification Commands

```bash
# Check if .env is ignored
git check-ignore .env
# Should output: .env

# Search for potential API keys in staged files
git diff --cached | grep -i "api.*key"
git diff --cached | grep "sk-ant-"

# List all files that will be committed
git status
```

## If You Accidentally Committed Secrets

### Option 1: Not Yet Pushed
```bash
# Remove from last commit
git reset --soft HEAD~1
# Remove the sensitive file
git rm --cached .env
# Recommit without the file
git commit -m "Your message"
```

### Option 2: Already Pushed (NUCLEAR OPTION)
```bash
# WARNING: This rewrites history!
# Only use if absolutely necessary

# Remove file from all history
git filter-branch --force --index-filter \
  "git rm --cached --ignore-unmatch .env" \
  --prune-empty --tag-name-filter cat -- --all

# Force push (coordinate with team first!)
git push origin --force --all
```

### Option 3: Secrets Already Exposed
1. **Immediately rotate the exposed secret:**
   - Anthropic API: Generate new key at console.anthropic.com
   - ngrok: Generate new authtoken at dashboard.ngrok.com
   - Custom API_KEY: Generate new random key
2. **Update .env with new secrets**
3. **Update deployment configurations** (Fly.io, etc.)
4. **Clean Git history** using Option 1 or 2 above

## Additional Security Measures

### Production Deployment
- [ ] Use environment variables on hosting platform (Fly.io secrets)
- [ ] Enable API key authentication for web endpoints
- [ ] Use HTTPS only (force_https = true in fly.toml)
- [ ] Implement rate limiting for public endpoints
- [ ] Log access attempts for monitoring

### Development
- [ ] Never share .env files via Slack/email
- [ ] Use password manager for storing API keys
- [ ] Rotate API keys periodically
- [ ] Use different keys for dev/staging/prod

## Current Security Status

✅ `.gitignore` includes .env files
✅ `.env.example` has placeholder values only
✅ No hardcoded API keys found in source code
✅ ANTHROPIC_API_KEY loaded from environment only
✅ All sensitive paths excluded from Git

## Last Security Audit

Date: 2026-01-08
Audited by: Claude Sonnet 4.5
Status: PASSED ✅
