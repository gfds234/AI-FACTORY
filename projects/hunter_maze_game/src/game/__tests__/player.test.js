// Test player movement and collision detection
import { describe, it, expect, beforeEach } from 'vitest';
import { Player } from '../player.js';
import { Maze } from '../maze.js';

describe('Player Movement', () => {
  let maze, player;

  beforeEach(() => {
    maze = new Maze(15);
    player = new Player(maze);
  });

  it('should initialize at center of maze', () => {
    expect(player.x).toBeGreaterThan(0);
    expect(player.y).toBeGreaterThan(0);
    expect(player.size).toBe(12);
  });

  it('should move down when ArrowDown key is pressed', () => {
    const initialY = player.y;
    const keys = { 'ArrowDown': true };

    player.update(keys);

    expect(player.y).toBeGreaterThan(initialY);
  });

  it('should move right when ArrowRight key is pressed', () => {
    const initialX = player.x;
    const keys = { 'ArrowRight': true };

    player.update(keys);

    expect(player.x).toBeGreaterThan(initialX);
  });

  it('should move with WASD keys', () => {
    const initialX = player.x;
    const keys = { 'd': true };

    player.update(keys);

    expect(player.x).toBeGreaterThan(initialX);
  });

  it('should NOT pass through walls', () => {
    // Create a simple test maze with known wall positions
    maze.grid = Array(15).fill(null).map(() => Array(15).fill(0));
    // Place walls
    for (let i = 0; i < 15; i++) {
      maze.grid[5][i] = 1; // Horizontal wall at row 5
    }

    // Position player just below the wall
    player.x = 7 * maze.cellSize;
    player.y = 6 * maze.cellSize;

    // Try to move up into the wall
    for (let i = 0; i < 20; i++) {
      player.update({ 'ArrowUp': true });
    }

    // Player should NOT be inside the wall
    const playerGridY = Math.floor(player.y / maze.cellSize);
    expect(playerGridY).toBeGreaterThan(5); // Should not be in row 5 or above
    expect(!maze.isWall(player.x, player.y)).toBe(true);
  });

  it('should stay within maze bounds', () => {
    // Try to move far left
    for (let i = 0; i < 100; i++) {
      player.update({ 'ArrowLeft': true });
    }

    expect(player.x).toBeGreaterThanOrEqual(player.size);
    expect(player.x).toBeLessThanOrEqual(maze.getWidth() - player.size);
  });

  it('should have velocity when moving', () => {
    const keys = { 'ArrowRight': true };
    player.update(keys);

    // Should have horizontal velocity
    expect(Math.abs(player.vx)).toBeGreaterThan(0);
  });

  it('should normalize diagonal movement speed', () => {
    const keys = { 'ArrowRight': true, 'ArrowDown': true };
    player.update(keys);

    // Diagonal speed should be normalized (not 2x faster)
    const speed = Math.sqrt(player.vx ** 2 + player.vy ** 2);
    expect(speed).toBeCloseTo(player.speed, 1);
  });
});
