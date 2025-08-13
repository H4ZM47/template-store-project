@echo off
setlocal enabledelayedexpansion

REM TemplateStore Development Startup Script for Windows
REM This script starts the database, backend, and frontend servers

title TemplateStore Development Environment

REM Configuration
set PROJECT_ROOT=%~dp0..
set BACKEND_PORT=8080
set FRONTEND_PORT=3000
set LOG_DIR=%PROJECT_ROOT%\tmp\logs
set PID_DIR=%PROJECT_ROOT%\tmp\pids

REM Create directories for logs and PIDs
if not exist "%LOG_DIR%" mkdir "%LOG_DIR%"
if not exist "%PID_DIR%" mkdir "%PID_DIR%"

REM Colors for Windows (using PowerShell for colored output)
set "GREEN=[32m"
set "RED=[31m"
set "YELLOW=[33m"
set "BLUE=[34m"
set "NC=[0m"

echo.
echo %BLUE%================================================%NC%
echo %BLUE%  TemplateStore Development Environment%NC%
echo %BLUE%================================================%NC%
echo.

REM Function to print status with timestamp
call :print_status "Starting TemplateStore Development Environment..."

REM Check dependencies
call :print_status "Checking dependencies..."

where go >nul 2>nul
if %errorlevel% neq 0 (
    call :print_error "Go is not installed or not in PATH"
    call :print_error "Please install Go 1.21 or higher and add it to your PATH"
    pause
    exit /b 1
)

where psql >nul 2>nul
if %errorlevel% neq 0 (
    call :print_warning "PostgreSQL client tools not found in PATH"
    call :print_warning "Database operations may not work properly"
)

call :print_success "Dependencies check completed"
echo.

REM Check if PostgreSQL is running
call :print_status "Checking PostgreSQL status..."
netstat -an | findstr ":5432" >nul 2>nul
if %errorlevel% equ 0 (
    call :print_success "PostgreSQL appears to be running on port 5432"
) else (
    call :print_warning "PostgreSQL may not be running on port 5432"
    call :print_status "Attempting to start PostgreSQL..."

    REM Try to start PostgreSQL service (Windows service)
    net start postgresql-x64-13 >nul 2>nul
    if !errorlevel! equ 0 (
        call :print_success "PostgreSQL service started"
    ) else (
        net start postgresql >nul 2>nul
        if !errorlevel! equ 0 (
            call :print_success "PostgreSQL service started"
        ) else (
            call :print_warning "Could not start PostgreSQL service automatically"
            call :print_warning "Please start PostgreSQL manually"
        )
    )
)

REM Setup database
call :print_status "Setting up database..."
createdb -h localhost -U postgres template_store >nul 2>nul
if %errorlevel% equ 0 (
    call :print_success "Database 'template_store' created"
) else (
    call :print_success "Database 'template_store' already exists or creation skipped"
)

REM Check if backend port is in use
call :print_status "Checking if port %BACKEND_PORT% is available..."
netstat -an | findstr ":%BACKEND_PORT%" >nul 2>nul
if %errorlevel% equ 0 (
    call :print_warning "Port %BACKEND_PORT% is already in use"
    call :print_status "Checking if it's our backend server..."

    powershell -Command "try { Invoke-WebRequest -Uri http://localhost:%BACKEND_PORT%/health -UseBasicParsing | Out-Null; Write-Host 'Backend already running' } catch { Write-Host 'Port in use by different service' }" >nul 2>nul
    if !errorlevel! equ 0 (
        call :print_success "Backend server is already running"
        set BACKEND_RUNNING=true
    ) else (
        call :print_error "Port %BACKEND_PORT% is in use by another service"
        pause
        exit /b 1
    )
) else (
    set BACKEND_RUNNING=false
)

REM Start backend server if not running
if "%BACKEND_RUNNING%"=="false" (
    call :print_status "Starting backend server..."

    REM Set environment variables
    set DB_HOST=localhost
    set DB_PORT=5432
    set DB_USER=postgres
    set DB_NAME=template_store
    set DB_PASSWORD=
    set DB_SSLMODE=disable
    set PORT=%BACKEND_PORT%
    set GIN_MODE=debug

    REM Start backend server
    cd /d "%PROJECT_ROOT%"
    start "TemplateStore Backend" /MIN cmd /c "go run cmd/server/main.go > %LOG_DIR%\backend.log 2>&1"

    call :print_status "Waiting for backend server to start..."
    timeout /t 3 /nobreak >nul

    REM Check if backend is responding
    for /l %%i in (1,1,15) do (
        powershell -Command "try { Invoke-WebRequest -Uri http://localhost:%BACKEND_PORT%/health -UseBasicParsing | Out-Null; exit 0 } catch { exit 1 }" >nul 2>nul
        if !errorlevel! equ 0 (
            call :print_success "Backend server is running on http://localhost:%BACKEND_PORT%"
            goto :backend_started
        )
        timeout /t 1 /nobreak >nul
    )

    call :print_error "Backend server failed to start or is not responding"
    call :print_error "Check the log: %LOG_DIR%\backend.log"
    pause
    exit /b 1

    :backend_started
)

REM Check if frontend port is in use
call :print_status "Checking if port %FRONTEND_PORT% is available..."
netstat -an | findstr ":%FRONTEND_PORT%" >nul 2>nul
if %errorlevel% equ 0 (
    call :print_warning "Port %FRONTEND_PORT% is already in use"
    powershell -Command "try { Invoke-WebRequest -Uri http://localhost:%FRONTEND_PORT% -UseBasicParsing | Out-Null; Write-Host 'Frontend already running' } catch { Write-Host 'Port in use by different service' }" >nul 2>nul
    if !errorlevel! equ 0 (
        call :print_success "Frontend server is already running"
        set FRONTEND_RUNNING=true
    ) else (
        call :print_error "Port %FRONTEND_PORT% is in use by another service"
        pause
        exit /b 1
    )
) else (
    set FRONTEND_RUNNING=false
)

REM Start frontend server if not running
if "%FRONTEND_RUNNING%"=="false" (
    call :print_status "Starting frontend server..."

    REM Start frontend server
    cd /d "%PROJECT_ROOT%"
    start "TemplateStore Frontend" /MIN cmd /c "go run cmd/web/main.go > %LOG_DIR%\frontend.log 2>&1"

    call :print_status "Waiting for frontend server to start..."
    timeout /t 2 /nobreak >nul

    REM Check if frontend is responding
    for /l %%i in (1,1,10) do (
        powershell -Command "try { Invoke-WebRequest -Uri http://localhost:%FRONTEND_PORT% -UseBasicParsing | Out-Null; exit 0 } catch { exit 1 }" >nul 2>nul
        if !errorlevel! equ 0 (
            call :print_success "Frontend server is running on http://localhost:%FRONTEND_PORT%"
            goto :frontend_started
        )
        timeout /t 1 /nobreak >nul
    )

    call :print_error "Frontend server failed to start or is not responding"
    call :print_error "Check the log: %LOG_DIR%\frontend.log"
    pause
    exit /b 1

    :frontend_started
)

REM Show final status
echo.
call :print_success "🚀 TemplateStore Development Environment Started!"
echo.
echo %GREEN%📊 Service Status:%NC%
echo   • Database (PostgreSQL): Running
echo   • Backend API: Running on http://localhost:%BACKEND_PORT%
echo   • Frontend Web: Running on http://localhost:%FRONTEND_PORT%
echo.
echo %BLUE%🌐 URLs:%NC%
echo   • Main Application: http://localhost:%FRONTEND_PORT%
echo   • Test Page: http://localhost:%FRONTEND_PORT%/test.html
echo   • Backend API: http://localhost:%BACKEND_PORT%/api/v1/
echo   • Health Check: http://localhost:%BACKEND_PORT%/health
echo.
echo %BLUE%📝 Logs:%NC%
echo   • Backend: %LOG_DIR%\backend.log
echo   • Frontend: %LOG_DIR%\frontend.log
echo.
echo %YELLOW%🛑 To stop all services: scripts\stop-dev.bat%NC%
echo.

REM Open browser to the application
call :print_status "Opening browser..."
timeout /t 2 /nobreak >nul
start http://localhost:%FRONTEND_PORT%

echo Press any key to exit...
pause >nul
goto :eof

REM Function definitions
:print_status
echo %BLUE%[%time%] %~1%NC%
goto :eof

:print_success
echo %GREEN%[%time%] ✅ %~1%NC%
goto :eof

:print_error
echo %RED%[%time%] ❌ %~1%NC%
goto :eof

:print_warning
echo %YELLOW%[%time%] ⚠️  %~1%NC%
goto :eof
