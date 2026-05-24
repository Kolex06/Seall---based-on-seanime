$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent $PSCommandPath
$BackendExe = Join-Path $Root "Seall.exe"
$DenshiBackendExe = Join-Path $Root "seall-denshi\binaries\seall-server-windows.exe"
$DesktopExe = Join-Path $Root "seall-denshi\dist\win-unpacked\Seall.exe"
$SourceElectronExe = Join-Path $Root "seall-denshi\node_modules\electron\dist\electron.exe"

function Show-SeallMessage {
    param(
        [string] $Message,
        [int] $Style = 64
    )

    try {
        $shell = New-Object -ComObject WScript.Shell
        $null = $shell.Popup($Message, 4, "Seall", $Style)
    } catch {
        Write-Host $Message
        Start-Sleep -Seconds 5
    }
}

try {
    $expectedPaths = @($BackendExe, $DenshiBackendExe, $DesktopExe, $SourceElectronExe)
    $processes = Get-Process -ErrorAction SilentlyContinue | Where-Object {
        try { $expectedPaths -contains $_.Path } catch { $false }
    }

    if (-not $processes) {
        Show-SeallMessage "Seall is not running."
        exit 0
    }

    $processes | Stop-Process -Force
    Show-SeallMessage "Seall has been stopped."
} catch {
    Show-SeallMessage "Seall could not be stopped.`n`n$($_.Exception.Message)" 48
    exit 1
}
