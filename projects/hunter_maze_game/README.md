# Hunter Maze Game

A fast-paced maze game where YOU are the hunter! Chase down and tag all 4 ghosts before time runs out.

## Unique Twist: Hunter Mode

Unlike traditional maze games, you're not running FROM the ghosts - you're hunting THEM! The ghosts flee from you using different AI behaviors, and the maze progressively shrinks every 20 seconds to force confrontations.

## Features

- **Dynamic Maze**: Procedurally generated maze that shrinks over time
- **Smart AI**: 4 ghosts with unique personalities (Fast, Random, Smart, Balanced)
- **Time Pressure**: 120-second timer with increasing difficulty
- **Combo System**: Chain tags within 5 seconds for score multipliers
- **Smooth Controls**: WASD or Arrow keys for 8-directional movement

## How to Play

1. Use **Arrow Keys** or **WASD** to move
2. Chase and **tag all 4 ghosts** before time runs out
3. **Watch the timer!** The maze shrinks every 20 seconds
4. Press **SPACE** to restart after game over

## Quick Play

**Double-click `PLAY_GAME.bat`** to start the game immediately!

Or run directly: `release\HunterMazeGame-win32-x64\HunterMazeGame.exe`

## Development

### Run in Browser
```bash
npm install
npm run dev
# Open http://localhost:5173
```

### Build Standalone .exe
```bash
npm run package
```
The executable will be in `release/HunterMazeGame-win32-x64/HunterMazeGame.exe`

## Tech Stack

- React + Vite
- HTML5 Canvas
- Electron (for .exe packaging)
- Vanilla JavaScript game engine

## Gameplay Tips

- **Red Ghost (Fast)**: Hard to catch, prioritize early
- **Green Ghost (Smart)**: Uses perpendicular fleeing, corner it
- **Blue Ghost (Random)**: Unpredictable, be patient
- **Yellow Ghost (Balanced)**: Easiest target, good for combo practice

**Score Bonus**: Tag ghosts quickly for time bonus and combo multipliers!

---

Made with AI Factory - Zero-cost game development
