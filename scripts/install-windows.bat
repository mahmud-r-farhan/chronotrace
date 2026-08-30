@echo off
title ChronoTrace Installer
echo Installing ChronoTrace...
powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0install-windows.ps1"
if %errorlevel% neq 0 (
    echo Installation encountered an error.
    pause
)
