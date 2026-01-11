# AI Factory - Claude Configuration

## Project Context

**What:** AI-powered MVP generation system with automated quality validation (Triple Guarantee)
**Tech Stack:** Go (orchestrator), React/Vite (templates), Node.js, Ollama (local LLM), optional Claude API
**Goal:** Generate production-ready MVPs with zero-cost infrastructure, charge $800-2,500 per project

## Core Principles (Priority Order)

1. **TEST-FIRST DEVELOPMENT** - Write failing tests BEFORE writing code (MANDATORY)
2. **ZERO-COST FIRST** - Always prefer free solutions (Ollama, templates) over paid (Claude API)
3. **SIMPLICITY** - Every change impacts minimum code possible
4. **VALIDATION** - All changes must pass Triple Guarantee (Build + Runtime + Test)
5. **DOCUMENTATION** - Keep PRODUCTION_GUIDE.md, TEST_RESULTS.md current
6. **NO SPECULATION** - Always read files before making claims

## Workflows

### Before Any Code Change:
1. **Investigate First** - Read relevant files, understand current state
2. **Plan** - Explain what will change and why (high-level)
3. **Execute** - Make minimal, focused changes
4. **Validate** - Test with test_react_counter if relevant
5. **Document** - Update docs if architecture changes

### For New Features:
1. Check if it can be done with existing templates (free)
2. If needs AI generation, estimate Claude API cost
3. Implement simplest version that works
4. Test with Triple Guarantee validation
5. Update PRODUCTION_GUIDE.md with usage instructions

### For Bug Fixes:
1. Reproduce the bug with test_react_counter
2. Fix in one file if possible
3. Verify fix doesn't break existing tests
4. Document in TEST_RESULTS.md if significant

## ZERO-BUG WORKFLOW (MANDATORY)

**Status**: Enforced for ALL projects as of 2026-01-11

### Test-Driven Development (TDD) Process

**Rule**: NO code without tests. Tests written FIRST, then code makes them pass.

#### Step 1: Write Failing Tests FIRST
```bash
# Create test file BEFORE implementation
src/game/__tests__/feature.test.js

# Run tests - should FAIL
npm run test
# ❌ Expected: Tests fail (no code yet)
```

#### Step 2: Write Minimum Code to Pass Tests
```javascript
// Implement ONLY what makes tests pass
// No extra features, no over-engineering
```

#### Step 3: Verify Tests Pass
```bash
npm run test
# ✅ Expected: All tests pass
```

#### Step 4: Refactor (Optional)
```javascript
// Clean up code while keeping tests green
// Tests prevent regression
```

### Mandatory Test Coverage

Every new project MUST include:

**Unit Tests** (`src/**/__tests__/*.test.js`)
- Component/module logic
- Edge cases and error handling
- Input validation

**Integration Tests** (`tests/integration/*.test.js`)
- Feature workflows
- API contracts
- Data flows

**E2E Tests** (`e2e/*.spec.js`)
- User journeys
- Critical paths
- Cross-browser (if web)

**Minimum Coverage**: 80% for core logic, 100% for critical paths

### Validation Checklist (BEFORE Marking Task Complete)

❌ **DO NOT mark task complete unless ALL pass:**

1. ✅ Unit tests pass (`npm run test`)
2. ✅ Integration tests pass (`npm run test:integration`)
3. ✅ E2E tests pass (`npm run test:e2e`)
4. ✅ Build succeeds (`npm run build`)
5. ✅ Manual playtest completed (if UI)
6. ✅ No console errors in browser/terminal
7. ✅ Performance acceptable (no lag/stuttering)

### Testing Commands (All Projects)

```bash
# Unit tests (fast, isolated)
npm run test              # Run once
npm run test:watch        # Live watching
npm run test:coverage     # Coverage report

# Integration tests
npm run test:integration

# E2E tests (requires dev server running)
npm run dev               # Terminal 1
npm run test:e2e          # Terminal 2

# Full validation
npm run validate          # Runs ALL checks
```

### Test Template Structure

**New projects automatically get**:
```
project/
├── src/
│   └── __tests__/         # Unit tests
├── tests/
│   └── integration/       # Integration tests
├── e2e/                   # E2E tests
├── vitest.config.js       # Unit test config
├── playwright.config.js   # E2E test config
└── scripts/
    └── validate.js        # Full validation
```

### When Tests Fail

**DO NOT skip failing tests. Fix the code or fix the test.**

1. Read error message carefully
2. Debug with `npm run test:watch`
3. Fix root cause (not symptoms)
4. Re-run until green
5. Commit only when ALL tests pass

### Example: Adding New Feature

```javascript
// ❌ OLD WAY (Bug-Prone):
// 1. Write code
// 2. Hope it works
// 3. Ship broken code

// ✅ NEW WAY (Zero-Bug):
// 1. Write test that fails
describe('New Feature', () => {
  it('should do the thing', () => {
    const result = doTheThing();
    expect(result).toBe('expected value');
  });
});

// 2. Run test - FAILS ❌
// 3. Write code to make it pass
function doTheThing() {
  return 'expected value';
}

// 4. Run test - PASSES ✅
// 5. Ship with confidence
```

### Enforcement

**package.json**:
```json
{
  "scripts": {
    "prepackage": "npm run validate",
    "prebuild": "npm run test"
  }
}
```

**Result**: Cannot build/package without passing tests

## What NOT to Do

❌ **Don't write code without tests FIRST** - Zero-Bug Workflow is MANDATORY
❌ Don't mark tasks complete with failing tests - Fix or remove the test
❌ Don't skip validation checklist - All 7 items required before delivery
❌ Don't commit broken code - Tests must be green
❌ Don't ask for permission - just do it (user confirmed they want autonomy)
❌ Don't add paid services without showing cost/benefit
❌ Don't create new files when editing existing ones works
❌ Don't modify config.json without backing it up first
❌ Don't change package.json dependencies without testing
❌ Don't speculate about code you haven't read
❌ Don't create README/docs unless explicitly requested
❌ Don't use complex solutions when simple ones exist
❌ Don't break the Triple Guarantee validation system

## What TO Do

✅ **Write tests BEFORE code** - Test-Driven Development is mandatory
✅ Run `npm run validate` before marking tasks complete
✅ Use `npm run test:watch` while developing (live feedback)
✅ Keep test coverage above 80% for core logic
✅ Run tests after changes: `cd projects/test_react_counter && npm run test`
✅ Use test_react_counter as validation reference
✅ Keep changes focused and minimal
✅ Prioritize template-based solutions (zero cost)
✅ Update documentation when changing architecture
✅ Show code examples in explanations
✅ Think through the problem before coding
✅ Make autonomous decisions when requirements are clear
✅ Use TodoWrite for multi-step tasks

## Common Commands

```bash
# Start system
./START_AI_FACTORY.bat

# Build orchestrator
go build -o ai_factory_working.exe .

# Test validation system
cd projects/test_react_counter
npm install
npm run build  # Should succeed in <500ms
npm run test   # Should pass 3/3
npm run dev    # Should start on port 5173

# Start Ollama
ollama serve

# Check Ollama models
curl http://localhost:11434/api/tags
```

## File Structure

```
AI FACTORY/
├── main.go                    - Entry point (DON'T MODIFY without reason)
├── config.json               - Settings (BACKUP before changes)
├── .env                      - API keys (NEVER commit)
├── project/                  - Orchestrator logic
├── supervisor/               - Multi-agent system
├── validation/               - Triple Guarantee code
├── task/                     - Task routing
├── projects/
│   └── test_react_counter/  - Reference template (DON'T BREAK)
├── artifacts/               - Generated code artifacts
└── web/                     - UI dashboard
```

## Architecture Rules

- **Orchestrator (Go):** Manages 6-phase workflow, don't add phases
- **Validation:** Build → Runtime → Test, all must pass
- **Templates:** React/Vite is reference, can add more templates
- **LLM Routing:** Complexity >7 uses advanced models
- **Quality Score:** 0-100 based on validation results

## Testing Strategy

**MANDATORY: Test-Driven Development (TDD)**

### For Every New Feature/Change:

1. **Write Tests FIRST** (that fail)
   - Unit tests: `src/__tests__/*.test.js`
   - Integration tests: `tests/integration/*.test.js`
   - E2E tests: `e2e/*.spec.js`

2. **Write Code** (minimum to pass tests)
   - Implement only what's needed
   - No over-engineering

3. **Verify Tests Pass**
   - `npm run test` - Unit tests ✅
   - `npm run test:integration` - Integration ✅
   - `npm run test:e2e` - End-to-end ✅

4. **Refactor** (optional, keep tests green)

5. **Full Validation**
   - `npm run validate` - All checks ✅
   - Manual playtest if UI
   - Document in TEST_RESULTS.md

### For Go Code:
1. Write tests FIRST: `*_test.go`
2. Run: `go test ./...`
3. Implement code
4. Verify: `go test ./... -v`

### Integration Validation:
1. Build and run test_react_counter
2. Start orchestrator and verify UI works
3. Run full Triple Guarantee validation

## Git Commit Style

```bash
# Auto-commit after successful changes
git add .
git commit -m "feat: <what changed> - <why it matters>

- Specific change 1
- Specific change 2

Validated: <test results>
Impact: <files changed>

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

## Cost Awareness

Always specify cost impact:
- **FREE:** Template changes, Ollama usage, local validation
- **PAID:** Claude API calls (~$0.50-2/MVP), external services
- **ROI:** If paid solution saves >30min, it's worth it (1 client pays for year)

## Error Handling

When something breaks:
1. Check Ollama is running: `ollama serve`
2. Check .env file exists and is valid
3. Check logs in terminal output
4. Test with test_react_counter to isolate issue
5. Document fix in PRODUCTION_GUIDE.md

## Success Metrics

Track these in TEST_RESULTS.md:
- Build time for test_react_counter (<500ms target)
- Test pass rate (100% target)
- Quality score (95+ target)
- Time to generate MVP (<2 hours manual, <10min auto)

## Communication Style

- **Be autonomous** - Make decisions, don't ask for confirmation
- **Be concise** - High-level explanations unless asked for details
- **Show impact** - What files changed, what tests ran
- **Use examples** - Show code snippets when explaining
- **Track tasks** - Use TodoWrite for multi-step work

## Priority Matrix

**P0 (Critical):** Breaks Triple Guarantee, blocks MVP generation
**P1 (High):** Affects quality score, slows down workflow
**P2 (Medium):** Nice-to-have features, minor improvements
**P3 (Low):** Documentation updates, refactoring

Always fix P0 immediately. Ask user for P1+.

---

**Last Updated:** 2026-01-11
**Optimized for:** AI Factory MVP generation system
**Based on:** Anthropic best practices + developer community research
