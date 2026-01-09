@echo off
setlocal enabledelayedexpansion

echo ========================================
echo Verify Latest Generated Project
echo ========================================
echo.

REM Find the latest generated project directory
set "latest_dir="
set "latest_time=0"

for /d %%D in (projects\generated_*) do (
    set "dir_name=%%~nxD"
    set "timestamp=!dir_name:generated_=!"
    if !timestamp! gtr !latest_time! (
        set "latest_time=!timestamp!"
        set "latest_dir=%%D"
    )
)

if "!latest_dir!"=="" (
    echo ERROR: No generated projects found in projects/ directory
    echo.
    pause
    exit /b 1
)

echo Latest project: !latest_dir!
echo.
echo ========================================
echo Verification Checklist
echo ========================================
echo.

REM Check for template fragments (bad)
set "has_fragments=0"
if exist "!latest_dir!\Prerequisites" (
    echo [FAIL] Found "Prerequisites" file - this is a template fragment!
    set "has_fragments=1"
)
if exist "!latest_dir!\Installing" (
    echo [FAIL] Found "Installing" file - this is a template fragment!
    set "has_fragments=1"
)
if exist "!latest_dir!\Break" (
    echo [FAIL] Found "Break" file - this is a template fragment!
    set "has_fragments=1"
)

if !has_fragments!==1 (
    echo.
    echo [ERROR] Template fragments detected!
    echo This means the LLM generated a README template instead of real code.
    echo The fixes may not be working correctly.
    echo.
    pause
    exit /b 1
) else (
    echo [PASS] No template fragments detected
)

echo.
echo Checking for required React/Vite files...
echo.

REM Check for required files (good)
set "score=0"

if exist "!latest_dir!\package.json" (
    echo [PASS] package.json exists
    set /a score+=1
) else (
    echo [FAIL] package.json NOT found
)

if exist "!latest_dir!\vite.config.js" (
    echo [PASS] vite.config.js exists
    set /a score+=1
) else (
    echo [FAIL] vite.config.js NOT found
)

if exist "!latest_dir!\index.html" (
    echo [PASS] index.html exists
    set /a score+=1
) else (
    echo [FAIL] index.html NOT found
)

if exist "!latest_dir!\src\main.jsx" (
    echo [PASS] src/main.jsx exists
    set /a score+=1
) else (
    if exist "!latest_dir!\src\main.js" (
        echo [PASS] src/main.js exists
        set /a score+=1
    ) else (
        echo [FAIL] src/main.jsx NOT found
    )
)

if exist "!latest_dir!\src\App.jsx" (
    echo [PASS] src/App.jsx exists
    set /a score+=1
) else (
    if exist "!latest_dir!\src\App.js" (
        echo [PASS] src/App.js exists
        set /a score+=1
    ) else (
        echo [FAIL] src/App.jsx NOT found
    )
)

echo.
echo ========================================
echo Verification Score: !score!/5
echo ========================================
echo.

if !score! geq 4 (
    echo [SUCCESS] Project looks good!
    echo.
    echo Next steps:
    echo 1. cd !latest_dir!
    echo 2. npm install
    echo 3. npm run dev
    echo 4. Open http://localhost:5173
    echo.
) else (
    echo [WARNING] Project is incomplete
    echo Expected at least 4/5 files, got !score!/5
    echo.
    echo Check artifacts/code_!latest_time!.md to see what the LLM generated
    echo.
)

echo Full directory listing:
echo.
dir /b !latest_dir!
echo.

if exist "!latest_dir!\src" (
    echo src/ directory contents:
    dir /b !latest_dir!\src
    echo.
)

pause
