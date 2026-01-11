// Canvas rendering functions
export class Renderer {
  constructor(canvas) {
    this.canvas = canvas;
    this.ctx = canvas.getContext('2d');
    this.flashAlpha = 0;
  }

  clear() {
    this.ctx.fillStyle = '#1a1a2e';
    this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);
  }

  drawMaze(maze) {
    // Calculate offset to center maze in canvas
    const mazeWidth = maze.size * maze.cellSize;
    const mazeHeight = maze.size * maze.cellSize;
    const offsetX = (this.canvas.width - mazeWidth) / 2;
    const offsetY = (this.canvas.height - mazeHeight) / 2;

    // Store offset for use in other draw methods
    this.offsetX = offsetX;
    this.offsetY = offsetY;

    this.ctx.strokeStyle = '#0f3460';
    this.ctx.lineWidth = 3;

    for (let y = 0; y < maze.size; y++) {
      for (let x = 0; x < maze.size; x++) {
        if (maze.grid[y][x] === 1) {
          const px = x * maze.cellSize + offsetX;
          const py = y * maze.cellSize + offsetY;

          this.ctx.fillStyle = '#16213e';
          this.ctx.fillRect(px, py, maze.cellSize, maze.cellSize);

          this.ctx.strokeStyle = '#0f3460';
          this.ctx.strokeRect(px, py, maze.cellSize, maze.cellSize);
        }
      }
    }
  }

  drawPlayer(player) {
    this.ctx.save();
    this.ctx.translate(player.x + (this.offsetX || 0), player.y + (this.offsetY || 0));
    this.ctx.rotate(player.direction);

    // Draw triangle (hunter)
    this.ctx.fillStyle = '#e94560';
    this.ctx.beginPath();
    this.ctx.moveTo(player.size, 0);
    this.ctx.lineTo(-player.size/2, -player.size/2);
    this.ctx.lineTo(-player.size/2, player.size/2);
    this.ctx.closePath();
    this.ctx.fill();

    // Glow effect
    this.ctx.shadowBlur = 15;
    this.ctx.shadowColor = '#e94560';
    this.ctx.fill();

    this.ctx.restore();
  }

  drawGhost(ghost) {
    const offsetX = this.offsetX || 0;
    const offsetY = this.offsetY || 0;

    if (ghost.isTagged) {
      // Draw X mark over tagged ghost
      this.ctx.strokeStyle = '#666';
      this.ctx.lineWidth = 3;
      this.ctx.beginPath();
      this.ctx.moveTo(ghost.x - ghost.size + offsetX, ghost.y - ghost.size + offsetY);
      this.ctx.lineTo(ghost.x + ghost.size + offsetX, ghost.y + ghost.size + offsetY);
      this.ctx.moveTo(ghost.x + ghost.size + offsetX, ghost.y - ghost.size + offsetY);
      this.ctx.lineTo(ghost.x - ghost.size + offsetX, ghost.y + ghost.size + offsetY);
      this.ctx.stroke();
      return;
    }

    const wobbleOffset = Math.sin(ghost.wobble) * 2;

    this.ctx.save();
    this.ctx.translate(ghost.x + offsetX, ghost.y + offsetY + wobbleOffset);

    // Ghost body
    this.ctx.fillStyle = ghost.color;
    this.ctx.beginPath();
    this.ctx.arc(0, 0, ghost.size, 0, Math.PI * 2);
    this.ctx.fill();

    // Ghost glow
    this.ctx.shadowBlur = 10;
    this.ctx.shadowColor = ghost.color;
    this.ctx.fill();

    // Eyes (scared look)
    this.ctx.shadowBlur = 0;
    this.ctx.fillStyle = 'white';
    this.ctx.beginPath();
    this.ctx.arc(-4, -3, 3, 0, Math.PI * 2);
    this.ctx.arc(4, -3, 3, 0, Math.PI * 2);
    this.ctx.fill();

    this.ctx.fillStyle = 'black';
    this.ctx.beginPath();
    this.ctx.arc(-4, -3, 1.5, 0, Math.PI * 2);
    this.ctx.arc(4, -3, 1.5, 0, Math.PI * 2);
    this.ctx.fill();

    this.ctx.restore();
  }

  drawTimer(time, maxTime, width) {
    const barWidth = width;
    const barHeight = 20;
    const y = 10;

    // Background
    this.ctx.fillStyle = '#16213e';
    this.ctx.fillRect(10, y, barWidth, barHeight);

    // Timer bar
    const ratio = time / maxTime;
    const color = ratio > 0.5 ? '#00ff00' : ratio > 0.25 ? '#ffd700' : '#ff0000';
    this.ctx.fillStyle = color;
    this.ctx.fillRect(10, y, barWidth * ratio, barHeight);

    // Border
    this.ctx.strokeStyle = '#0f3460';
    this.ctx.lineWidth = 2;
    this.ctx.strokeRect(10, y, barWidth, barHeight);

    // Time text
    this.ctx.fillStyle = 'white';
    this.ctx.font = '14px monospace';
    this.ctx.textAlign = 'center';
    this.ctx.fillText(Math.ceil(time) + 's', barWidth/2 + 10, y + 15);
  }

  drawScore(score, ghostsRemaining) {
    this.ctx.fillStyle = 'white';
    this.ctx.font = 'bold 18px monospace';
    this.ctx.textAlign = 'right';

    this.ctx.fillText('Score: ' + score, this.canvas.width - 10, 25);
    this.ctx.fillText('Ghosts: ' + ghostsRemaining, this.canvas.width - 10, 50);
  }

  drawGameOver(won, score, time) {
    this.ctx.fillStyle = 'rgba(0, 0, 0, 0.8)';
    this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);

    this.ctx.fillStyle = won ? '#00ff00' : '#ff0000';
    this.ctx.font = 'bold 48px monospace';
    this.ctx.textAlign = 'center';
    this.ctx.fillText(won ? 'VICTORY!' : 'TIME UP!', this.canvas.width/2, this.canvas.height/2 - 40);

    this.ctx.fillStyle = 'white';
    this.ctx.font = '24px monospace';
    this.ctx.fillText('Score: ' + score, this.canvas.width/2, this.canvas.height/2 + 10);
    this.ctx.fillText('Time: ' + time.toFixed(1) + 's', this.canvas.width/2, this.canvas.height/2 + 45);

    this.ctx.font = '18px monospace';
    this.ctx.fillStyle = '#aaa';
    this.ctx.fillText('Press SPACE to restart', this.canvas.width/2, this.canvas.height/2 + 90);
  }

  flash() {
    this.flashAlpha = 1;
  }

  updateFlash() {
    if (this.flashAlpha > 0) {
      this.ctx.fillStyle = `rgba(255, 255, 255, ${this.flashAlpha})`;
      this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);
      this.flashAlpha -= 0.05;
    }
  }
}
