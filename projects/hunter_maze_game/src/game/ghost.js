// Ghost AI behavior
export class Ghost {
  constructor(maze, color, personality, startPos) {
    this.maze = maze;
    this.color = color;
    this.personality = personality;
    this.size = 14;
    this.isTagged = false;
    this.x = startPos.x;
    this.y = startPos.y;
    this.vx = 0;
    this.vy = 0;
    this.wobble = 0;

    // Different speeds based on personality
    this.speed = {
      'fast': 4,
      'random': 2.5,
      'smart': 3,
      'balanced': 3
    }[personality];
  }

  update(player) {
    if (this.isTagged) return;

    this.wobble += 0.2;

    // Calculate flee direction away from player
    const dx = this.x - player.x;
    const dy = this.y - player.y;
    const dist = Math.sqrt(dx * dx + dy * dy);

    if (dist < 150) {
      // Close to player - flee!
      switch (this.personality) {
        case 'fast':
          // Fast runner - direct flee
          this.vx = (dx / dist) * this.speed;
          this.vy = (dy / dist) * this.speed;
          break;

        case 'random':
          // Random direction changes
          if (Math.random() < 0.05) {
            const angle = Math.random() * Math.PI * 2;
            this.vx = Math.cos(angle) * this.speed;
            this.vy = Math.sin(angle) * this.speed;
          }
          break;

        case 'smart':
          // Smart fleeing - perpendicular movement
          const perpAngle = Math.atan2(dy, dx) + (Math.random() - 0.5) * Math.PI / 2;
          this.vx = Math.cos(perpAngle) * this.speed;
          this.vy = Math.sin(perpAngle) * this.speed;
          break;

        case 'balanced':
          // Balanced - direct flee with some randomness
          const fleeAngle = Math.atan2(dy, dx) + (Math.random() - 0.5) * 0.5;
          this.vx = Math.cos(fleeAngle) * this.speed;
          this.vy = Math.sin(fleeAngle) * this.speed;
          break;
      }
    } else {
      // Far from player - patrol randomly
      if (Math.random() < 0.02 || (this.vx === 0 && this.vy === 0)) {
        const angle = Math.random() * Math.PI * 2;
        this.vx = Math.cos(angle) * (this.speed * 0.5);
        this.vy = Math.sin(angle) * (this.speed * 0.5);
      }
    }

    // Ensure ghost always has some movement
    if (this.vx === 0 && this.vy === 0) {
      const angle = Math.random() * Math.PI * 2;
      this.vx = Math.cos(angle) * this.speed;
      this.vy = Math.sin(angle) * this.speed;
    }

    // Try to move
    const newX = this.x + this.vx;
    const newY = this.y + this.vy;

    // Check if new position is valid (not in wall)
    const centerOK = !this.maze.isWall(newX, newY);
    const rightOK = !this.maze.isWall(newX + this.size/2, newY);
    const leftOK = !this.maze.isWall(newX - this.size/2, newY);
    const bottomOK = !this.maze.isWall(newX, newY + this.size/2);
    const topOK = !this.maze.isWall(newX, newY - this.size/2);

    if (centerOK && rightOK && leftOK && bottomOK && topOK) {
      // Can move freely
      this.x = newX;
      this.y = newY;
    } else {
      // Try sliding along walls (move in one direction at a time)
      const tryX = this.x + this.vx;
      const tryY = this.y + this.vy;

      // Try X movement only
      if (!this.maze.isWall(tryX, this.y) &&
          !this.maze.isWall(tryX + this.size/2, this.y) &&
          !this.maze.isWall(tryX - this.size/2, this.y)) {
        this.x = tryX;
      } else {
        // Bounce X velocity
        this.vx *= -0.5;
      }

      // Try Y movement only
      if (!this.maze.isWall(this.x, tryY) &&
          !this.maze.isWall(this.x, tryY + this.size/2) &&
          !this.maze.isWall(this.x, tryY - this.size/2)) {
        this.y = tryY;
      } else {
        // Bounce Y velocity
        this.vy *= -0.5;
      }
    }

    // Keep in bounds
    const maxX = this.maze.getWidth();
    const maxY = this.maze.getHeight();
    this.x = Math.max(this.size, Math.min(maxX - this.size, this.x));
    this.y = Math.max(this.size, Math.min(maxY - this.size, this.y));
  }

  tag() {
    this.isTagged = true;
  }
}

export function createGhosts(maze) {
  const cellSize = maze.cellSize;
  const size = maze.size;

  return [
    new Ghost(maze, '#FF0000', 'fast', {
      x: 2 * cellSize,
      y: 2 * cellSize
    }),
    new Ghost(maze, '#00BFFF', 'random', {
      x: (size - 3) * cellSize,
      y: 2 * cellSize
    }),
    new Ghost(maze, '#00FF00', 'smart', {
      x: 2 * cellSize,
      y: (size - 3) * cellSize
    }),
    new Ghost(maze, '#FFD700', 'balanced', {
      x: (size - 3) * cellSize,
      y: (size - 3) * cellSize
    })
  ];
}
