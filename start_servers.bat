@echo off
echo ========================================
echo AI FACTORY - Server Startup Script
echo ========================================
echo.
echo This script will help you start both required services:
echo 1. Ollama (LLM backend)
echo 2. AI Factory Server
echo.
echo ========================================
echo.

echo Step 1: Starting Ollama...
echo.
echo Opening new window for Ollama...
start "Ollama Server" cmd /k "echo Ollama Server && echo. && ollama serve"

timeout /t 3 /nobreak > nul

echo.
echo Step 2: Starting AI Factory...
echo.
echo Opening new window for AI Factory...
start "AI Factory Server" cmd /k "cd /d "%~dp0" && echo AI Factory Server && echo. && ai_factory.exe -mode=server"

timeout /t 5 /nobreak > nul

echo.
echo ========================================
echo Both servers should be starting now!
echo ========================================
echo.
echo Check the two new windows that opened:
echo 1. Ollama Server window - should show Ollama running
echo 2. AI Factory Server window - should show "Server running at http://localhost:8080"
echo.
echo Once both are running, you can:
echo 1. Run test_simple_react.bat to test code generation
echo 2. Open http://localhost:8080 in your browser for the Web UI
echo 3. Run demo_generator.exe for full demo generation
echo.
echo Press any key to close this window...
pause > nul
