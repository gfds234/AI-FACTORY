// Maze generation and management
export class Maze {
  constructor(initialSize = 15) {
    this.size = initialSize;
    this.cellSize = 40;
    this.grid = [];
    this.generate();
  }

  generate() {
    // Initialize grid with walls
    this.grid = Array(this.size).fill(null).map(() =>
      Array(this.size).fill(1)
    );

    // Randomized DFS maze generation
    const stack = [];
    const startX = 1;
    const startY = 1;

    this.grid[startY][startX] = 0;
    stack.push([startX, startY]);

    const directions = [
      [0, -2], [2, 0], [0, 2], [-2, 0]
    ];

    while (stack.length > 0) {
      const [x, y] = stack[stack.length - 1];
      const neighbors = [];

      // Check all four directions
      for (const [dx, dy] of directions) {
        const nx = x + dx;
        const ny = y + dy;

        if (nx > 0 && nx < this.size - 1 && ny > 0 && ny < this.size - 1) {
          if (this.grid[ny][nx] === 1) {
            neighbors.push([nx, ny, x + dx/2, y + dy/2]);
          }
        }
      }

      if (neighbors.length > 0) {
        // Choose random neighbor
        const [nx, ny, wx, wy] = neighbors[Math.floor(Math.random() * neighbors.length)];
        this.grid[ny][nx] = 0;
        this.grid[wy][wx] = 0;
        stack.push([nx, ny]);
      } else {
        stack.pop();
      }
    }

    // Create some open areas for better gameplay
    for (let i = 0; i < Math.floor(this.size / 3); i++) {
      const x = 1 + Math.floor(Math.random() * (this.size - 2));
      const y = 1 + Math.floor(Math.random() * (this.size - 2));
      this.grid[y][x] = 0;
    }
  }

  shrink() {
    if (this.size < 8) return false;  // Don't shrink below 6 (size-2 >= 6)

    // Remove outer ring
    this.size -= 2;
    const newGrid = [];

    for (let y = 1; y < this.grid.length - 1; y++) {
      const row = [];
      for (let x = 1; x < this.grid[y].length - 1; x++) {
        row.push(this.grid[y][x]);
      }
      newGrid.push(row);
    }

    this.grid = newGrid;

    // Ensure connectivity after shrinking - carve paths through center
    this.ensureConnectivity();

    this.cellSize = Math.min(40, 600 / this.size);
    return true;
  }

  ensureConnectivity() {
    // Carve a cross pattern through the center to guarantee connectivity
    const centerX = Math.floor(this.size / 2);
    const centerY = Math.floor(this.size / 2);

    // Horizontal path through center
    for (let x = 1; x < this.size - 1; x++) {
      this.grid[centerY][x] = 0;
      if (centerY > 0) this.grid[centerY - 1][x] = 0;
      if (centerY < this.size - 1) this.grid[centerY + 1][x] = 0;
    }

    // Vertical path through center
    for (let y = 1; y < this.size - 1; y++) {
      this.grid[y][centerX] = 0;
      if (centerX > 0) this.grid[y][centerX - 1] = 0;
      if (centerX < this.size - 1) this.grid[y][centerX + 1] = 0;
    }
  }

  isWall(x, y) {
    const gridX = Math.floor(x / this.cellSize);
    const gridY = Math.floor(y / this.cellSize);

    if (gridX < 0 || gridX >= this.size || gridY < 0 || gridY >= this.size) {
      return true;
    }

    return this.grid[gridY][gridX] === 1;
  }

  getWidth() {
    return this.size * this.cellSize;
  }

  getHeight() {
    return this.size * this.cellSize;
  }

  findNearestOpenCell(x, y) {
    // Convert pixel coordinates to grid coordinates
    let gridX = Math.floor(x / this.cellSize);
    let gridY = Math.floor(y / this.cellSize);

    // Clamp to grid bounds
    gridX = Math.max(0, Math.min(this.size - 1, gridX));
    gridY = Math.max(0, Math.min(this.size - 1, gridY));

    // If current cell is open, return it
    if (this.grid[gridY][gridX] === 0) {
      return {
        x: (gridX + 0.5) * this.cellSize,
        y: (gridY + 0.5) * this.cellSize
      };
    }

    // BFS to find nearest open cell
    const queue = [[gridX, gridY, 0]];
    const visited = new Set();
    visited.add(`${gridX},${gridY}`);

    const directions = [[0, -1], [1, 0], [0, 1], [-1, 0]];

    while (queue.length > 0) {
      const [cx, cy, dist] = queue.shift();

      for (const [dx, dy] of directions) {
        const nx = cx + dx;
        const ny = cy + dy;
        const key = `${nx},${ny}`;

        if (nx >= 0 && nx < this.size && ny >= 0 && ny < this.size && !visited.has(key)) {
          visited.add(key);

          if (this.grid[ny][nx] === 0) {
            // Found open cell - return center of cell
            return {
              x: (nx + 0.5) * this.cellSize,
              y: (ny + 0.5) * this.cellSize
            };
          }

          queue.push([nx, ny, dist + 1]);
        }
      }
    }

    // Fallback to center of maze
    return {
      x: (this.size / 2) * this.cellSize,
      y: (this.size / 2) * this.cellSize
    };
  }
}
