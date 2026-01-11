import { useEffect, useRef } from 'react';
import { GameEngine } from './game/engine.js';
import './Game.css';

function Game() {
  const canvasRef = useRef(null);
  const engineRef = useRef(null);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    canvas.width = 800;
    canvas.height = 600;

    // Initialize game engine
    const engine = new GameEngine(canvas);
    engineRef.current = engine;

    // Keyboard event handlers
    const handleKeyDown = (e) => {
      e.preventDefault();
      engine.handleKeyDown(e.key);
    };

    const handleKeyUp = (e) => {
      e.preventDefault();
      engine.handleKeyUp(e.key);
    };

    window.addEventListener('keydown', handleKeyDown);
    window.addEventListener('keyup', handleKeyUp);

    // Start game
    engine.start();

    // Cleanup
    return () => {
      engine.isRunning = false;
      window.removeEventListener('keydown', handleKeyDown);
      window.removeEventListener('keyup', handleKeyUp);
    };
  }, []);

  return (
    <div className="game-container">
      <div className="game-header">
        <h1>HUNTER MAZE</h1>
        <p>Tag all ghosts before time runs out!</p>
        <p className="controls">Controls: Arrow Keys or WASD</p>
      </div>
      <canvas ref={canvasRef} className="game-canvas" />
      <div className="game-footer">
        <p>The maze shrinks every 20 seconds - hunt them down quickly!</p>
      </div>
    </div>
  );
}

export default Game;
