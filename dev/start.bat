@echo off
title Todo App Development Environment

echo.
echo ========================================
echo    TODO APP DEVELOPMENT STARTUP
echo ========================================
echo.

echo [1/4] Checking Docker...
docker --version >nul 2>&1
if errorlevel 1 (
    echo Docker is not available
    pause
    exit /b 1
)
echo Docker is available

echo.
echo [2/4] Starting Cassandra...
docker start cassandra >nul 2>&1
if errorlevel 1 (
    echo Creating new Cassandra container...
    docker run --name cassandra -p 9042:9042 -d cassandra:latest
    timeout /t 30 /nobreak >nul
) else (
    echo Cassandra started
)

echo.
echo [3/4] Starting Todo App...
cd /d "%~dp0\.."
start /b cmd /c "go run . > app.log 2>&1"
timeout /t 8 /nobreak >nul

echo.
echo [4/4] Opening browser...
start http://localhost:8080

echo.
echo ========================================
echo    DEVELOPMENT ENVIRONMENT READY!
echo ========================================
echo.
echo Application URL: http://localhost:8080
echo Login Page: http://localhost:8080/login
echo Database: Cassandra on localhost:9042
echo Logs: app.log
echo.
echo Press Ctrl+C to stop...
echo ========================================

:: This pause keeps the app running until user interrupts
pause >nul

echo.
echo Stopping Todo App...
taskkill /f /im go.exe >nul 2>&1

echo Stopping Cassandra...
docker stop cassandra >nul 2>&1

echo All services stopped.
pause
