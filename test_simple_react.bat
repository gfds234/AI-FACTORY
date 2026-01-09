@echo off
echo ========================================
echo Phase 7 Test: Simple React Counter
echo ========================================
echo.

echo Prerequisites Check:
echo 1. Ollama must be running (ollama serve)
echo 2. AI Factory must be running (./ai_factory.exe -mode=server)
echo.
echo Press Ctrl+C to cancel, or
pause

echo.
echo Sending code generation request...
echo.

curl -X POST http://localhost:8080/task ^
  -H "Content-Type: application/json" ^
  -d "{\"input\":\"Create a simple React counter component with Vite. The app should display a counter starting at 0, with buttons to increment (+1) and decrement (-1). Use modern React hooks. Include proper component structure with src/App.jsx, and add at least one Vitest test for the counter functionality.\",\"task_type\":\"code\"}"

echo.
echo.
echo ========================================
echo Test request sent!
echo ========================================
echo.
echo Next steps:
echo 1. Check the response above for the project path
echo 2. Look in projects/ for the newest generated_TIMESTAMP folder
echo 3. Verify it contains:
echo    - package.json (with vite, react, vitest)
echo    - vite.config.js
echo    - src/App.jsx (with REAL counter code, not placeholders)
echo    - NO files named "Prerequisites" or "Installing"
echo.
echo 4. If successful, run:
echo    cd projects/generated_TIMESTAMP
echo    npm install
echo    npm run dev
echo    npm run test
echo.
pause
