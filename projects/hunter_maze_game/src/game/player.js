// Player movement and controls
export class Player {
  constructor(maze) {
    this.maze = maze;
    this.size = 12;
    this.speed = 3;
    this.reset();
  }

  reset() {
    // Start in center of maze
    const centerCell = Math.floor(this.maze.size / 2);
    this.x = centerCell * this.maze.cellSize + this.maze.cellSize / 2;
    this.y = centerCell * this.maze.cellSize + this.maze.cellSize / 2;
    this.vx = 0;
    this.vy = 0;
    this.direction = 0;
  }

  update(keys) {
    // Update velocity based on keys
    this.vx = 0;
    this.vy = 0;

    if (keys['ArrowUp'] || keys['w'] || keys['W']) this.vy = -this.speed;
    if (keys['ArrowDown'] || keys['s'] || keys['S']) this.vy = this.speed;
    if (keys['ArrowLeft'] || keys['a'] || keys['A']) this.vx = -this.speed;
    if (keys['ArrowRight'] || keys['d'] || keys['D']) this.vx = this.speed;

    // Normalize diagonal movement
    if (this.vx !== 0 && this.vy !== 0) {
      this.vx *= 0.707;
      this.vy *= 0.707;
    }

    // Calculate direction for rendering
    if (this.vx !== 0 || this.vy !== 0) {
      this.direction = Math.atan2(this.vy, this.vx);
    }

    // Try to move, check collision
    const newX = this.x + this.vx;
    const newY = this.y + this.vy;

    // Check multiple points around player for better collision
    const points = [
      [newX, newY],
      [newX + this.size/2, newY],
      [newX - this.size/2, newY],
      [newX, newY + this.size/2],
      [newX, newY - this.size/2]
    ];

    let canMove = true;
    for (const [px, py] of points) {
      if (this.maze.isWall(px, py)) {
        canMove = false;
        break;
      }
    }

    if (canMove) {
      this.x = newX;
      this.y = newY;
    } else {
      // Try sliding along walls
      if (!this.maze.isWall(newX, this.y)) {
        this.x = newX;
      }
      if (!this.maze.isWall(this.x, newY)) {
        this.y = newY;
      }
    }

    // Keep player in bounds
    const maxX = this.maze.getWidth();
    const maxY = this.maze.getHeight();
    this.x = Math.max(this.size, Math.min(maxX - this.size, this.x));
    this.y = Math.max(this.size, Math.min(maxY - this.size, this.y));
  }

  isNearGhost(ghost, radius = 20) {
    const dx = this.x - ghost.x;
    const dy = this.y - ghost.y;
    return Math.sqrt(dx * dx + dy * dy) < radius;
  }
}
