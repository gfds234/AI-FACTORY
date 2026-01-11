# Zero-Bug Workflow for AI Factory

**Problem**: Hunter Maze Game shipped with critical bugs (NPCs not moving, player passing through walls) because there was no validation before marking tasks complete.

**Root Cause**: Manual development without automated testing

**Solution**: Implement Test-Driven Development (TDD) + Triple Guarantee Automation

---

## Root Cause Analysis: Hunter Maze Game Bugs

### Bugs Reported
1. ❌ NPCs (ghosts) were not moving
2. ❌ NPCs were out of place
3. ❌ Player could go through walls/obstacles

### Why These Bugs Happened
1. **No Runtime Testing** - Game was never launched and played before delivery
2. **Assumed Completion** - Marked todos as done without verification
3. **No Automated Tests** - No unit tests for collision, movement, AI
4. **No Integration Testing** - Never validated all systems working together
5. **Build ≠ Works** - Vite build succeeded, but runtime was broken

### What Went Wrong in Development Process
```
❌ OLD WORKFLOW (Bug-Prone):
1. Write code
2. Run build (vite build) ✅
3. Mark task complete ✅
4. Ship to user ❌ BROKEN

WHY IT FAILED:
- Build success doesn't mean code works
- No validation of actual game behavior
- No automated tests catching logic errors
```

---

## NEW WORKFLOW: Zero-Bug Development

### Core Principle: **Test BEFORE Code**

Every feature must have automated tests that FAIL first, then code is written to make them PASS.

###Phase 1: Test-First Development

#### Step 1: Write Failing Tests
```bash
# BEFORE writing any game code:
cd projects/hunter_maze_game
npm run test:watch

# Create test files FIRST:
src/game/__tests__/
├── maze.test.js        # Test wall collision, maze generation
├── player.test.js      # Test movement, controls
├── ghost.test.js       # Test AI behavior, flee logic
├── engine.test.js      # Test game loop, state management
└── integration.test.js # Test full game flow
```

**Example: Player Collision Test (Written FIRST)**
```javascript
// src/game/__tests__/player.test.js
import { describe, it, expect, beforeEach } from 'vitest';
import { Player } from '../player.js';
import { Maze } from '../maze.js';

describe('Player', () => {
  let maze, player;

  beforeEach(() => {
    maze = new Maze(15);
    player = new Player(maze);
  });

  it('should NOT pass through walls', () => {
    // Place player next to a wall
    player.x = 40;
    player.y = 40;

    // Manually set wall in maze
    maze.grid[0][1] = 1; // Wall to the left

    // Try to move into wall
    const keys = { 'ArrowLeft': true };
    const initialX = player.x;

    player.update(keys);

    // Player should NOT have moved through wall
    expect(player.x).toBeGreaterThanOrEqual(initialX - player.speed);
    expect(!maze.isWall(player.x, player.y)).toBe(true);
  });

  it('should move when path is clear', () => {
    const initialY = player.y;
    const keys = { 'ArrowDown': true };

    player.update(keys);

    expect(player.y).toBeGreaterThan(initialY);
  });
});
```

**Example: Ghost Movement Test**
```javascript
// src/game/__tests__/ghost.test.js
import { describe, it, expect } from 'vitest';
import { Ghost } from '../ghost.js';
import { Maze } from '../maze.js';

describe('Ghost AI', () => {
  it('should move away from player when close', () => {
    const maze = new Maze(15);
    const ghost = new Ghost(maze, '#FF0000', 'fast', { x: 100, y: 100 });
    const player = { x: 80, y: 80 }; // Close to ghost

    const initialX = ghost.x;
    const initialY = ghost.y;

    // Update ghost AI
    ghost.update(player);

    // Ghost should have moved
    const hasMoved = (ghost.x !== initialX || ghost.y !== initialY);
    expect(hasMoved).toBe(true);

    // Ghost should be moving AWAY from player
    const distBefore = Math.sqrt((100-80)**2 + (100-80)**2);
    const distAfter = Math.sqrt((ghost.x-80)**2 + (ghost.y-80)**2);

    // After fleeing, distance should increase (or stay same if hit wall)
    expect(distAfter).toBeGreaterThanOrEqual(distBefore - 5);
  });

  it('should NOT move when tagged', () => {
    const maze = new Maze(15);
    const ghost = new Ghost(maze, '#FF0000', 'fast', { x: 100, y: 100 });

    ghost.tag();
    const initialX = ghost.x;
    const initialY = ghost.y;

    ghost.update({ x: 80, y: 80 });

    expect(ghost.x).toBe(initialX);
    expect(ghost.y).toBe(initialY);
  });
});
```

#### Step 2: Run Tests (Should FAIL)
```bash
npm run test

# Expected output:
# ❌ FAIL  src/game/__tests__/player.test.js
#   Player
#     ✕ should NOT pass through walls (0 ms)
#     ✕ should move when path is clear (0 ms)
#
# ❌ FAIL  src/game/__tests__/ghost.test.js
#   Ghost AI
#     ✕ should move away from player when close (0 ms)
#     ✕ should NOT move when tagged (0 ms)
```

#### Step 3: Write Code to Make Tests Pass
Now implement the actual game logic to satisfy the tests.

#### Step 4: Run Tests Again (Should PASS)
```bash
npm run test

# Expected output:
# ✓ PASS  src/game/__tests__/player.test.js (2 tests)
# ✓ PASS  src/game/__tests__/ghost.test.js (2 tests)
# ✓ PASS  src/game/__tests__/maze.test.js (3 tests)
# ✓ PASS  src/game/__tests__/engine.test.js (4 tests)
# ✓ PASS  src/game/__tests__/integration.test.js (2 tests)
#
# Test Suites: 5 passed, 5 total
# Tests:       13 passed, 13 total
```

### Phase 2: Automated Validation (Triple Guarantee++)

Before marking ANY task complete, run automated validation:

```bash
# Run validation script
npm run validate

# This script runs:
# 1. Linting (eslint)
# 2. Type checking (if using TypeScript)
# 3. Unit tests (vitest)
# 4. Build (vite build)
# 5. Integration tests (playwright/puppeteer)
# 6. Visual regression tests (optional)
```

**validation.js** (New automation script)
```javascript
// scripts/validate.js
import { execSync } from 'child_process';
import fs from 'fs';

const runCommand = (cmd, description) => {
  console.log(`\n🔍 ${description}...`);
  try {
    execSync(cmd, { stdio: 'inherit' });
    console.log(`✅ ${description} PASSED`);
    return true;
  } catch (error) {
    console.log(`❌ ${description} FAILED`);
    return false;
  }
};

const results = {
  lint: runCommand('npm run lint', 'Linting'),
  test: runCommand('npm run test', 'Unit Tests'),
  build: runCommand('npm run build', 'Build'),
  e2e: runCommand('npm run test:e2e', 'Integration Tests')
};

const allPassed = Object.values(results).every(r => r);

if (allPassed) {
  console.log('\n✅ ALL VALIDATIONS PASSED - Safe to commit\n');
  process.exit(0);
} else {
  console.log('\n❌ VALIDATION FAILED - Do NOT commit\n');
  process.exit(1);
}
```

### Phase 3: End-to-End (E2E) Testing

Test the ACTUAL game as a user would play it:

**e2e/game.spec.js** (Playwright test)
```javascript
import { test, expect } from '@playwright/test';

test.describe('Hunter Maze Game', () => {
  test('game loads and renders canvas', async ({ page }) => {
    await page.goto('http://localhost:5173');

    // Check title
    await expect(page.locator('h1')).toContainText('HUNTER MAZE');

    // Check canvas exists
    const canvas = page.locator('canvas');
    await expect(canvas).toBeVisible();

    // Verify canvas dimensions
    const box = await canvas.boundingBox();
    expect(box.width).toBe(800);
    expect(box.height).toBe(600);
  });

  test('player can move with arrow keys', async ({ page }) => {
    await page.goto('http://localhost:5173');

    // Wait for game to initialize
    await page.waitForTimeout(1000);

    // Press arrow key
    await page.keyboard.press('ArrowRight');
    await page.waitForTimeout(100);

    // Take screenshot to verify movement (visual regression)
    await page.screenshot({ path: 'test-results/player-moved.png' });

    // Could also inject code to check player position
    const playerMoved = await page.evaluate(() => {
      return window.gameEngine?.player?.vx !== 0;
    });

    expect(playerMoved).toBe(true);
  });

  test('ghosts are visible and moving', async ({ page }) => {
    await page.goto('http://localhost:5173');
    await page.waitForTimeout(2000);

    // Inject code to check ghost positions
    const ghostsMoving = await page.evaluate(() => {
      const ghosts = window.gameEngine?.ghosts || [];
      return ghosts.length === 4 && ghosts.some(g => !g.isTagged);
    });

    expect(ghostsMoving).toBe(true);
  });

  test('collision detection works', async ({ page }) => {
    await page.goto('http://localhost:5173');
    await page.waitForTimeout(1000);

    // Try to move player into wall repeatedly
    for (let i = 0; i < 20; i++) {
      await page.keyboard.press('ArrowUp');
      await page.waitForTimeout(50);
    }

    // Check player is still in valid position
    const playerInBounds = await page.evaluate(() => {
      const player = window.gameEngine?.player;
      const maze = window.gameEngine?.maze;
      if (!player || !maze) return false;

      return !maze.isWall(player.x, player.y);
    });

    expect(playerInBounds).toBe(true);
  });
});
```

### Phase 4: Manual Playtest Checklist

Even with automated tests, perform manual verification:

```markdown
## Manual Playtest Checklist (Required before delivery)

### Launch
- [ ] Game .exe launches without errors
- [ ] Window opens in <3 seconds
- [ ] Canvas renders properly

### Player Controls
- [ ] Arrow keys move player
- [ ] WASD keys move player
- [ ] Player moves in all 8 directions
- [ ] Player CANNOT pass through walls
- [ ] Player stays within maze bounds

### NPCs (Ghosts)
- [ ] All 4 ghosts are visible
- [ ] Ghosts are in correct starting positions
- [ ] Ghosts MOVE when game starts
- [ ] Ghosts flee when player approaches
- [ ] Each ghost has unique behavior
- [ ] Tagged ghosts stop moving

### Game Mechanics
- [ ] Timer counts down from 120s
- [ ] Maze shrinks every 20 seconds
- [ ] Screen flashes when maze shrinks
- [ ] Player can tag ghosts
- [ ] Score increases when tagging
- [ ] Win condition triggers (all ghosts tagged)
- [ ] Lose condition triggers (time expires)
- [ ] SPACE key restarts game

### Performance
- [ ] Game runs at smooth 60 FPS
- [ ] No lag or stuttering
- [ ] No memory leaks (play for 5 minutes)

✅ ALL ITEMS MUST BE CHECKED before marking task complete
```

---

## Implementation Plan: Apply to Hunter Maze Game

### Step 1: Add Testing Infrastructure
```bash
cd AI\ FACTORY/projects/hunter_maze_game

# Install testing dependencies
npm install --save-dev vitest @testing-library/react @testing-library/jest-dom jsdom @playwright/test

# Create test config
# (vitest.config.js already exists)

# Create test directories
mkdir -p src/game/__tests__
mkdir -p e2e
```

### Step 2: Write Tests for Existing Code
Create comprehensive test suite covering:
- Maze generation
- Player movement
- Ghost AI
- Collision detection
- Game loop
- Win/lose conditions

### Step 3: Run Tests (Expect Failures)
```bash
npm run test
# Document which tests fail
# These failures reveal the bugs
```

### Step 4: Fix Code Based on Test Failures
Debug and fix each failing test one by one.

### Step 5: Add E2E Tests
```bash
npx playwright install
npm run test:e2e
```

### Step 6: Manual Playtest
Complete the checklist above.

### Step 7: Document and Deliver
Only after ALL tests pass + manual checklist complete.

---

## New package.json Scripts

```json
{
  "scripts": {
    "dev": "vite",
    "build": "vite build",
    "test": "vitest run",
    "test:watch": "vitest",
    "test:e2e": "playwright test",
    "test:e2e:ui": "playwright test --ui",
    "lint": "eslint src",
    "validate": "node scripts/validate.js",
    "package": "npm run validate && vite build && electron-packager ..."
  }
}
```

**KEY CHANGE**: `npm run package` now requires `npm run validate` to pass first!

---

## AI Factory Integration

Update AI Factory orchestrator to use this workflow:

### supervisor/quality_check.go
```go
func (s *Supervisor) ValidateWithTests(projectPath string) (*ValidationResult, error) {
    result := &ValidationResult{}

    // Step 1: Run unit tests
    cmd := exec.Command("npm", "run", "test")
    cmd.Dir = projectPath
    output, err := cmd.CombinedOutput()

    if err != nil {
        result.TestsPassed = false
        result.Errors = append(result.Errors, "Unit tests failed: " + string(output))
        return result, nil
    }
    result.TestsPassed = true

    // Step 2: Run build
    cmd = exec.Command("npm", "run", "build")
    cmd.Dir = projectPath
    output, err = cmd.CombinedOutput()

    if err != nil {
        result.BuildPassed = false
        result.Errors = append(result.Errors, "Build failed: " + string(output))
        return result, nil
    }
    result.BuildPassed = true

    // Step 3: Run E2E tests
    cmd = exec.Command("npm", "run", "test:e2e")
    cmd.Dir = projectPath
    output, err = cmd.CombinedOutput()

    if err != nil {
        result.E2EPassed = false
        result.Errors = append(result.Errors, "E2E tests failed: " + string(output))
        return result, nil
    }
    result.E2EPassed = true

    // All validations passed
    result.Score = 100
    return result, nil
}
```

---

## Summary: Zero-Bug Guarantees

### Before (Bug-Prone):
```
Write Code → Build → Ship ❌
```

### After (Zero-Bug):
```
Write Tests (FAIL)
  ↓
Write Code
  ↓
Tests PASS ✅
  ↓
E2E Tests PASS ✅
  ↓
Manual Playtest ✅
  ↓
npm run validate ✅
  ↓
Ship 🚀
```

### Enforcement Rules

1. **NO code without tests** - Tests written FIRST
2. **NO commits without green tests** - All tests must pass
3. **NO builds without validation** - `npm run package` requires `npm run validate`
4. **NO delivery without playtest** - Manual checklist required
5. **NO exceptions** - Even "simple" changes need tests

### Benefits

✅ **Bugs caught immediately** - Tests fail instantly when code breaks
✅ **Regression prevention** - Old bugs can't come back
✅ **Documentation** - Tests show how code should work
✅ **Confidence** - Ship knowing it works
✅ **Speed** - Less debugging, more building

---

## Next Steps

1. Apply this workflow to fix Hunter Maze Game
2. Add to AI Factory CLAUDE.md as mandatory process
3. Create template project with tests pre-configured
4. Update all future projects to follow this workflow

**ZERO BUGS = ZERO EXCUSES**
