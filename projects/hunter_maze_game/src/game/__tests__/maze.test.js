// Test maze generation and collision
import { describe, it, expect, beforeEach } from 'vitest';
import { Maze } from '../maze.js';

describe('Maze', () => {
  let maze;

  beforeEach(() => {
    maze = new Maze(15);
  });

  it('should initialize with correct size', () => {
    expect(maze.size).toBe(15);
    expect(maze.cellSize).toBe(40);
    expect(maze.grid.length).toBe(15);
  });

  it('should have paths (0) and walls (1)', () => {
    let hasWalls = false;
    let hasPaths = false;

    for (let y = 0; y < maze.size; y++) {
      for (let x = 0; x < maze.size; x++) {
        if (maze.grid[y][x] === 1) hasWalls = true;
        if (maze.grid[y][x] === 0) hasPaths = true;
      }
    }

    expect(hasWalls).toBe(true);
    expect(hasPaths).toBe(true);
  });

  it('should correctly detect walls at grid positions', () => {
    // Manually set a wall
    maze.grid[5][5] = 1;

    // Center of that cell
    const x = 5 * maze.cellSize + maze.cellSize / 2;
    const y = 5 * maze.cellSize + maze.cellSize / 2;

    expect(maze.isWall(x, y)).toBe(true);
  });

  it('should correctly detect paths', () => {
    // Manually set a path
    maze.grid[5][5] = 0;

    const x = 5 * maze.cellSize + maze.cellSize / 2;
    const y = 5 * maze.cellSize + maze.cellSize / 2;

    expect(maze.isWall(x, y)).toBe(false);
  });

  it('should shrink when called', () => {
    const initialSize = maze.size;

    const didShrink = maze.shrink();

    expect(didShrink).toBe(true);
    expect(maze.size).toBe(initialSize - 2);
  });

  it('should not shrink below minimum size', () => {
    // Shrink multiple times
    while (maze.size > 6) {
      maze.shrink();
    }

    const result = maze.shrink();

    expect(result).toBe(false);
    expect(maze.size).toBe(6);
  });

  it('should update cellSize when shrinking', () => {
    const initialCellSize = maze.cellSize;

    // Shrink to force cellSize adjustment
    while (maze.size > 8) {
      maze.shrink();
    }

    // CellSize should have changed
    expect(maze.cellSize).toBeGreaterThanOrEqual(initialCellSize);
  });

  it('should treat out-of-bounds as walls', () => {
    expect(maze.isWall(-10, -10)).toBe(true);
    expect(maze.isWall(10000, 10000)).toBe(true);
  });

  it('should calculate correct width and height', () => {
    expect(maze.getWidth()).toBe(maze.size * maze.cellSize);
    expect(maze.getHeight()).toBe(maze.size * maze.cellSize);
  });
});
