@echo off
SETLOCAL EnableDelayedExpansion

:: AI FACTORY - One-Click Quickstart
:: This script starts the orchestrator and ngrok tunnel automatically

echo ========================================
echo    AI FACTORY - Quickstart Launcher
echo ========================================
echo.

:: Check if orchestrator.exe exists
if not exist "orchestrator.exe" (
    echo ERROR: orchestrator.exe not found!
    echo Please run: go build -o orchestrator.exe .
    echo.
    pause
    exit /b 1
)

:: Check for ANTHROPIC_API_KEY in environment
set "API_KEY_FOUND=0"
if defined ANTHROPIC_API_KEY set "API_KEY_FOUND=1"

:: Check for .env file
if exist ".env" (
    echo Found .env file - loading configuration...
    for /f "usebackq tokens=1,* delims==" %%a in (".env") do (
        if "%%a"=="ANTHROPIC_API_KEY" (
            set "ANTHROPIC_API_KEY=%%b"
            set "API_KEY_FOUND=1"
        )
    )
)

:: Prompt for API key if not found
if !API_KEY_FOUND!==0 (
    echo.
    echo ANTHROPIC_API_KEY not found in environment or .env file
    echo.
    set /p "ANTHROPIC_API_KEY=Enter your Anthropic API Key: "

    :: Save to .env for future use
    echo ANTHROPIC_API_KEY=!ANTHROPIC_API_KEY! > .env
    echo API_KEY saved to .env file
    echo.
)

:: Check if ngrok is installed
set "NGROK_PATH="
if exist "C:\Users\lampr\Downloads\ngrok\ngrok.exe" (
    set "NGROK_PATH=C:\Users\lampr\Downloads\ngrok\ngrok.exe"
) else if exist "C:\ngrok\ngrok.exe" (
    set "NGROK_PATH=C:\ngrok\ngrok.exe"
) else (
    where ngrok >nul 2>&1
    if !ERRORLEVEL!==0 (
        set "NGROK_PATH=ngrok"
    )
)

if "!NGROK_PATH!"=="" (
    echo.
    echo WARNING: ngrok not found. Only local access will be available.
    echo Download ngrok from: https://ngrok.com/download
    echo.
    set "USE_NGROK=0"
) else (
    set "USE_NGROK=1"
)

:: Check ngrok authentication
if !USE_NGROK!==1 (
    "!NGROK_PATH!" config check >nul 2>&1
    if !ERRORLEVEL! NEQ 0 (
        echo.
        echo ngrok is not configured with an auth token.
        echo.
        set /p "NGROK_TOKEN=Enter your ngrok authtoken (or press Enter to skip): "

        if not "!NGROK_TOKEN!"=="" (
            "!NGROK_PATH!" config add-authtoken !NGROK_TOKEN!
            echo ngrok configured successfully!
        ) else (
            set "USE_NGROK=0"
            echo Skipping ngrok - only local access will be available.
        )
        echo.
    )
)

:: Start orchestrator in new window
echo [1/3] Starting AI FACTORY Orchestrator...
start "AI FACTORY Orchestrator" /MIN cmd /c "orchestrator.exe"

:: Wait for orchestrator to start
echo Waiting for orchestrator to initialize...
timeout /t 3 /nobreak >nul

:: Verify orchestrator is running
powershell -Command "$response = try { Invoke-RestMethod -Uri 'http://localhost:8080/health' -TimeoutSec 2 } catch { $null }; if ($response) { exit 0 } else { exit 1 }"

if !ERRORLEVEL! NEQ 0 (
    echo ERROR: Orchestrator failed to start!
    echo Check the Orchestrator window for errors.
    pause
    exit /b 1
)

echo Orchestrator is running! [OK]
echo.

:: Start ngrok if available
if !USE_NGROK!==1 (
    echo [2/3] Starting ngrok tunnel...
    start "ngrok Tunnel" /MIN cmd /c ""!NGROK_PATH!" http 8080"

    :: Wait for ngrok to start
    timeout /t 5 /nobreak >nul

    :: Get public URL
    for /f "delims=" %%i in ('curl -s http://localhost:4040/api/tunnels ^| findstr /r "https://.*ngrok.*\.dev"') do (
        set "NGROK_LINE=%%i"
    )

    :: Parse URL from JSON
    for /f "tokens=2 delims=:," %%a in ("!NGROK_LINE!") do (
        set "PUBLIC_URL=%%a"
        set "PUBLIC_URL=!PUBLIC_URL:"=!"
        set "PUBLIC_URL=!PUBLIC_URL: =!"
    )

    if not "!PUBLIC_URL!"=="" (
        echo ngrok tunnel active! [OK]
        echo.
        echo ========================================
        echo    AI FACTORY is now LIVE!
        echo ========================================
        echo.
        echo Local Access:  http://localhost:8080
        echo Public Access: !PUBLIC_URL!
        echo.
        echo Opening browser...
        timeout /t 2 /nobreak >nul
        start !PUBLIC_URL!
    ) else (
        echo Warning: Could not retrieve ngrok URL
        echo Check the ngrok window or visit: http://localhost:4040
        echo.
        echo Opening local URL...
        start http://localhost:8080
    )
) else (
    echo [2/3] Skipping ngrok (not available)
    echo.
    echo ========================================
    echo    AI FACTORY is running locally!
    echo ========================================
    echo.
    echo Access at: http://localhost:8080
    echo.
    echo Opening browser...
    timeout /t 2 /nobreak >nul
    start http://localhost:8080
)

echo.
echo [3/3] Setup complete!
echo.
echo ========================================
echo  Quick Tips:
echo ========================================
echo.
echo - Go to "Projects" tab to create a new MVP
echo - Toggle between List and Kanban views
echo - Keep this window and terminals open
echo - Press Ctrl+C in each window to stop
echo.
echo ========================================

:: Keep window open
echo Press any key to open ngrok dashboard (optional)...
pause >nul
start http://localhost:4040

ENDLOCAL
