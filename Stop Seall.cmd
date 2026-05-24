@echo off
set "SCRIPT=%~dp0Stop Seall.ps1"
powershell.exe -NoProfile -ExecutionPolicy Bypass -WindowStyle Hidden -File "%SCRIPT%"
