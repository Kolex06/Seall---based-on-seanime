@echo off
set "SCRIPT=%~dp0Open Seall.ps1"
powershell.exe -NoProfile -ExecutionPolicy Bypass -WindowStyle Hidden -File "%SCRIPT%"
