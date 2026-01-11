// Main game engine
import { Maze } from './maze.js';
import { Player } from './player.js';
import { createGhosts } from './ghost.js';
import { Renderer } from './renderer.js';

export class GameEngine {
  constructor(canvas) {
    this.canvas = canvas;
    this.renderer = new Renderer(canvas);
    this.keys = {};
    this.isRunning = false;
    this.gameOver = false;
    this.won = false;

    this.maxTime = 120;
    this.shrinkInterval = 20;
    this.lastShrinkTime = 0;

    this.init();
  }

  init() {
    this.maze = new Maze(15);
    this.player = new Player(this.maze);
    this.ghosts = createGhosts(this.maze);

    this.time = this.maxTime;
    this.score = 0;
    this.combo = 0;
    this.lastTagTime = 0;
    this.gameOver = false;
    this.won = false;

    this.lastShrinkTime = 0;
  }

  start() {
    this.isRunning = true;
    this.lastTime = performance.now();
    this.loop();
  }

  restart() {
    this.init();
    this.start();
  }

  loop() {
    if (!this.isRunning) return;

    const currentTime = performance.now();
    const deltaTime = (currentTime - this.lastTime) / 1000;
    this.lastTime = currentTime;

    this.update(deltaTime);
    this.render();

    requestAnimationFrame(() => this.loop());
  }

  update(deltaTime) {
    if (this.gameOver) {
      // Check for restart
      if (this.keys[' ']) {
        this.restart();
      }
      return;
    }

    // Update timer
    this.time -= deltaTime;

    // Check for maze shrink
    const timeSinceShrink = this.maxTime - this.time - this.lastShrinkTime;
    if (timeSinceShrink >= this.shrinkInterval) {
      if (this.maze.shrink()) {
        this.renderer.flash();
        this.lastShrinkTime = this.maxTime - this.time;

        // Reposition player if now in wall
        if (this.maze.isWall(this.player.x, this.player.y)) {
          const newPos = this.maze.findNearestOpenCell(this.player.x, this.player.y);
          this.player.x = newPos.x;
          this.player.y = newPos.y;
        }

        // Reposition ghosts if now in walls
        for (const ghost of this.ghosts) {
          if (this.maze.isWall(ghost.x, ghost.y)) {
            const newPos = this.maze.findNearestOpenCell(ghost.x, ghost.y);
            ghost.x = newPos.x;
            ghost.y = newPos.y;
          }
        }
      }
    }

    // Update player
    this.player.update(this.keys);

    // Update ghosts
    for (const ghost of this.ghosts) {
      ghost.update(this.player);
    }

    // Check collisions (player tagging ghosts)
    for (const ghost of this.ghosts) {
      if (!ghost.isTagged && this.player.isNearGhost(ghost)) {
        ghost.tag();

        // Calculate combo
        const timeSinceLastTag = this.maxTime - this.time - this.lastTagTime;
        if (timeSinceLastTag < 5) {
          this.combo++;
        } else {
          this.combo = 0;
        }
        this.lastTagTime = this.maxTime - this.time;

        // Calculate score
        const baseScore = 100;
        const timeBonus = (this.time / this.maxTime) * 2;
        const comboMultiplier = 1 + (this.combo * 0.5);

        const points = Math.floor(baseScore * timeBonus * comboMultiplier);
        this.score += points;
      }
    }

    // Check win condition
    const allTagged = this.ghosts.every(g => g.isTagged);
    if (allTagged) {
      this.gameOver = true;
      this.won = true;
      this.isRunning = false;
    }

    // Check lose condition
    if (this.time <= 0) {
      this.time = 0;
      this.gameOver = true;
      this.won = false;
      this.isRunning = false;
    }
  }

  render() {
    this.renderer.clear();
    this.renderer.drawMaze(this.maze);

    // Draw all ghosts
    for (const ghost of this.ghosts) {
      this.renderer.drawGhost(ghost);
    }

    // Draw player on top
    this.renderer.drawPlayer(this.player);

    // Draw UI
    const ghostsRemaining = this.ghosts.filter(g => !g.isTagged).length;
    this.renderer.drawTimer(this.time, this.maxTime, this.canvas.width - 20);
    this.renderer.drawScore(this.score, ghostsRemaining);

    // Draw flash effect
    this.renderer.updateFlash();

    // Draw game over screen
    if (this.gameOver) {
      this.renderer.drawGameOver(this.won, this.score, this.maxTime - this.time);
    }
  }

  handleKeyDown(key) {
    this.keys[key] = true;
  }

  handleKeyUp(key) {
    this.keys[key] = false;
  }
}
