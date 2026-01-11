// Test ghost AI and movement
import { describe, it, expect, beforeEach } from 'vitest';
import { Ghost, createGhosts } from '../ghost.js';
import { Maze } from '../maze.js';

describe('Ghost AI', () => {
  let maze;

  beforeEach(() => {
    maze = new Maze(15);
  });

  it('should create 4 ghosts with unique colors', () => {
    const ghosts = createGhosts(maze);

    expect(ghosts.length).toBe(4);
    expect(ghosts[0].color).toBe('#FF0000'); // Red
    expect(ghosts[1].color).toBe('#00BFFF'); // Blue
    expect(ghosts[2].color).toBe('#00FF00'); // Green
    expect(ghosts[3].color).toBe('#FFD700'); // Yellow
  });

  it('should initialize ghosts in different corners', () => {
    const ghosts = createGhosts(maze);

    // Each ghost should be in a different position
    const positions = ghosts.map(g => `${g.x},${g.y}`);
    const uniquePositions = new Set(positions);

    expect(uniquePositions.size).toBe(4);
  });

  it('should move away from player when player is close', () => {
    const ghost = new Ghost(maze, '#FF0000', 'fast', { x: 200, y: 200 });
    const player = { x: 150, y: 150 }; // Close to ghost (within 150px)

    const initialX = ghost.x;
    const initialY = ghost.y;

    // Update ghost multiple times
    for (let i = 0; i < 10; i++) {
      ghost.update(player);
    }

    // Ghost should have moved from initial position
    const hasMoved = (ghost.x !== initialX || ghost.y !== initialY);
    expect(hasMoved).toBe(true);

    // Calculate distance before and after
    const distBefore = Math.sqrt((initialX - player.x) ** 2 + (initialY - player.y) ** 2);
    const distAfter = Math.sqrt((ghost.x - player.x) ** 2 + (ghost.y - player.y) ** 2);

    // Ghost should be moving away (distance should increase or stay similar due to walls)
    expect(distAfter).toBeGreaterThanOrEqual(distBefore - 20);
  });

  it('should NOT move when tagged', () => {
    const ghost = new Ghost(maze, '#FF0000', 'fast', { x: 200, y: 200 });
    const player = { x: 150, y: 150 };

    ghost.tag();

    const initialX = ghost.x;
    const initialY = ghost.y;

    ghost.update(player);

    expect(ghost.x).toBe(initialX);
    expect(ghost.y).toBe(initialY);
    expect(ghost.isTagged).toBe(true);
  });

  it('should have different speeds based on personality', () => {
    const fastGhost = new Ghost(maze, '#FF0000', 'fast', { x: 200, y: 200 });
    const randomGhost = new Ghost(maze, '#0000FF', 'random', { x: 200, y: 200 });

    expect(fastGhost.speed).toBe(4);
    expect(randomGhost.speed).toBe(2.5);
  });

  it('should exhibit random behavior for random personality', () => {
    const ghost = new Ghost(maze, '#0000FF', 'random', { x: 200, y: 200 });
    const player = { x: 500, y: 500 }; // Far from ghost

    const positions = [];

    // Update multiple times and track positions
    for (let i = 0; i < 50; i++) {
      ghost.update(player);
      positions.push({ x: ghost.x, y: ghost.y });
    }

    // Ghost should have moved from initial position eventually
    const lastPos = positions[positions.length - 1];
    const hasMoved = (lastPos.x !== 200 || lastPos.y !== 200);

    expect(hasMoved).toBe(true);
  });

  it('should wobble animation counter increase', () => {
    const ghost = new Ghost(maze, '#FF0000', 'fast', { x: 200, y: 200 });
    const initialWobble = ghost.wobble;

    ghost.update({ x: 500, y: 500 });

    expect(ghost.wobble).toBeGreaterThan(initialWobble);
  });

  it('should stay within maze bounds', () => {
    const ghost = new Ghost(maze, '#FF0000', 'fast', { x: 50, y: 50 });
    const player = { x: 0, y: 0 }; // Push ghost towards edge

    // Update many times
    for (let i = 0; i < 100; i++) {
      ghost.update(player);
    }

    const maxX = maze.getWidth();
    const maxY = maze.getHeight();

    expect(ghost.x).toBeGreaterThanOrEqual(ghost.size);
    expect(ghost.x).toBeLessThanOrEqual(maxX - ghost.size);
    expect(ghost.y).toBeGreaterThanOrEqual(ghost.size);
    expect(ghost.y).toBeLessThanOrEqual(maxY - ghost.size);
  });
});
