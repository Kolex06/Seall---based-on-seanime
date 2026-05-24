$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent $PSCommandPath
$DenshiDir = Join-Path $Root "seall-denshi"
$DesktopExe = Join-Path $Root "seall-denshi\dist\win-unpacked\Seall.exe"
$SourceElectronExe = Join-Path $Root "seall-denshi\node_modules\electron\dist\electron.exe"

function Show-SeallMessage {
    param(
        [string] $Message,
        [int] $Style = 48
    )

    try {
        $shell = New-Object -ComObject WScript.Shell
        $null = $shell.Popup($Message, 0, "Seall", $Style)
    } catch {
        Write-Host $Message
        Start-Sleep -Seconds 10
    }
}

try {
    if (Test-Path -LiteralPath $DesktopExe) {
        Start-Process -FilePath $DesktopExe -WorkingDirectory (Split-Path -Parent $DesktopExe)
        exit 0
    }

    if (Test-Path -LiteralPath $SourceElectronExe) {
        Start-Process -FilePath $SourceElectronExe -ArgumentList "`"$DenshiDir`"" -WorkingDirectory $DenshiDir
        exit 0
    }

    Show-SeallMessage "Seall could not start because the desktop app file was not found.`n`nMissing:`n$DesktopExe`n$SourceElectronExe"
    exit 1
} catch {
    Show-SeallMessage "Seall could not start.`n`n$($_.Exception.Message)"
    exit 1
}
