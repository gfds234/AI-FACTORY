@echo off
echo ========================================
echo AI Studio - Manager Console
echo ========================================
echo.
echo Starting orchestrator server...
echo Server will run on http://localhost:8080
echo.
echo Opening web interface in your browser...
echo.

REM Start the orchestrator server
start "AI Studio Server" "%~dp0orchestrator" -mode=server -port=8080

REM Wait 3 seconds for server to start
timeout /t 3 /nobreak >nul

REM Open web UI in default browser
start "" "http://localhost:8080"

echo.
echo ========================================
echo Server is running!
echo ========================================
echo Web UI should open automatically.
echo If not, open: http://localhost:8080
echo.
echo To stop the server, close the "AI Studio Server" window
echo ========================================
echo.

pause
