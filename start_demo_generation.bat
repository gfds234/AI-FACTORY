@echo off
echo.
echo ========================================
echo   AI FACTORY - Demo Generation Starter
echo ========================================
echo.
echo This will help you generate your first customer demo!
echo.
echo STEP 1: Checking if Ollama is running...
echo.

curl -s http://localhost:11434/api/version >nul 2>&1
if %errorlevel% neq 0 (
    echo [!] Ollama is NOT running
    echo.
    echo Please start Ollama in a separate terminal:
    echo   ^> ollama serve
    echo.
    echo Then run this script again.
    pause
    exit /b 1
)

echo [OK] Ollama is running!
echo.
echo STEP 2: Starting AI Factory server...
echo.
echo Opening server in new window...
start "AI Factory Server" cmd /k "cd /d %~dp0 && ai_factory.exe -mode=server"

timeout /t 8 /nobreak >nul
echo.
echo STEP 3: Testing server connection...
curl -s http://localhost:8080/ >nul 2>&1
if %errorlevel% neq 0 (
    echo [!] Server not responding yet, waiting a bit longer...
    timeout /t 5 /nobreak >nul
)

echo [OK] Server is running!
echo.
echo ========================================
echo   READY TO GENERATE DEMO!
echo ========================================
echo.
echo Option 1: Use Web UI (Easiest)
echo   1. Open: http://localhost:8080
echo   2. Click "Tasks" tab
echo   3. Select "Generate Code"
echo   4. Paste customer requirements
echo   5. Click Execute
echo.
echo Option 2: Use API (Automated)
echo   Run: curl -X POST http://localhost:8080/task -H "Content-Type: application/json" -d @customer_request.json
echo.
echo Your demo will be generated in: projects/generated_TIMESTAMP/
echo.
echo ========================================
echo.
echo Press any key to open the Web UI...
pause >nul
start http://localhost:8080
