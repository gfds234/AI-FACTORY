# Hunter Maze Game - Roadmap

## Completed Features

- ✅ Core game loop with 60 FPS rendering
- ✅ Procedural maze generation (DFS algorithm)
- ✅ Player movement (8-directional, WASD/Arrow keys)
- ✅ Ghost AI (4 personalities: fast, random, smart, balanced)
- ✅ Collision detection (player vs walls, player vs ghosts)
- ✅ Tagging mechanic (hunt ghosts to win)
- ✅ Maze shrinking every 20 seconds
- ✅ Screen centering when maze shrinks
- ✅ Path connectivity after shrinking (cross-pattern)
- ✅ Entity repositioning after shrink (never stuck in walls)
- ✅ Timer system (120 seconds)
- ✅ Score tracking with combo multipliers
- ✅ Win/lose conditions
- ✅ Game restart (SPACE key)
- ✅ Visual effects (glow, flash, wobble animation)
- ✅ Test infrastructure (Vitest + 16 unit tests)

## Planned Features

### High Priority

- [ ] **Wrap-around boundaries**: When entities exit one side of the maze, they teleport to the opposite side
  - Player wraps horizontally and vertically at maze edges
  - Ghosts wrap around as well
  - Maintains position offset when wrapping
  - Enables strategic chasing/escaping across boundaries

### Medium Priority

- [ ] Sound effects and background music
- [ ] Particle effects when tagging ghosts
- [ ] Power-ups (speed boost, freeze ghosts, extra time)
- [ ] Multiple difficulty levels
- [ ] High score persistence (local storage)
- [ ] Ghost AI improvements (pathfinding algorithms)
- [ ] Minimap overlay
- [ ] Settings menu (sound volume, controls customization)

### Low Priority

- [ ] Multiple maze themes (ice, lava, forest)
- [ ] Achievements system
- [ ] Leaderboard integration
- [ ] Multiplayer mode (local co-op)
- [ ] Custom maze editor
- [ ] Replay system
- [ ] Mobile/touch controls support

## Bug Tracking

All bugs are tracked via Zero-Bug TDD Workflow with automated tests.

See [BUG_FIXES.md](./BUG_FIXES.md) for detailed bug fix history.
