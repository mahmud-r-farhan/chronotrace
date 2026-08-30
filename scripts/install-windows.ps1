# ChronoTrace Windows 1-Click Installer
# Installs ChronoTrace & ChronoTrace Daemon to %LOCALAPPDATA%\Programs\ChronoTrace

$ErrorActionPreference = "Stop"

$installDir = Join-Path $env:LOCALAPPDATA "Programs\ChronoTrace"
Write-Host "Installing ChronoTrace to $installDir ..." -ForegroundColor Cyan

# 1. Stop any running instances
Stop-Process -Name "ChronoTrace" -Force -ErrorAction SilentlyContinue
Stop-Process -Name "chronotrace-daemon" -Force -ErrorAction SilentlyContinue

# 2. Create destination directory
New-Item -ItemType Directory -Force -Path $installDir | Out-Null

# 3. Locate source binaries
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path

# Possible candidate paths
$candidateDaemons = @(
    (Join-Path $scriptDir "chronotrace-daemon.exe"),
    (Join-Path $scriptDir "..\build\chronotrace-daemon.exe"),
    (Join-Path (Get-Location) "build\chronotrace-daemon.exe")
)
$candidateGuis = @(
    (Join-Path $scriptDir "ChronoTrace.exe"),
    (Join-Path $scriptDir "..\build\ChronoTrace.exe"),
    (Join-Path (Get-Location) "build\ChronoTrace.exe")
)
$candidateIcons = @(
    (Join-Path $scriptDir "icon.ico"),
    (Join-Path $scriptDir "..\assets\icons\icon.ico"),
    (Join-Path (Get-Location) "assets\icons\icon.ico"),
    (Join-Path $scriptDir "..\gui\build\windows\icon.ico"),
    (Join-Path (Get-Location) "gui\build\windows\icon.ico")
)

$srcDaemon = $candidateDaemons | Where-Object { Test-Path $_ } | Select-Object -First 1
$srcGui    = $candidateGuis | Where-Object { Test-Path $_ } | Select-Object -First 1
$srcIcon   = $candidateIcons | Where-Object { Test-Path $_ } | Select-Object -First 1

if (-not $srcDaemon -or -not $srcGui) {
    Write-Error "Source binaries not found. Please build the project first or extract the full release package."
    exit 1
}

# 4. Copy binaries and assets
Copy-Item -Force $srcDaemon (Join-Path $installDir "chronotrace-daemon.exe")
Copy-Item -Force $srcGui (Join-Path $installDir "ChronoTrace.exe")
if ($srcIcon -and (Test-Path $srcIcon)) {
    Copy-Item -Force $srcIcon (Join-Path $installDir "icon.ico")
}

# 5. Create Start Menu & Desktop Shortcuts
$wsh = New-Object -ComObject WScript.Shell

$startMenuDir = Join-Path ([Environment]::GetFolderPath("Programs")) "ChronoTrace"
New-Item -ItemType Directory -Force -Path $startMenuDir | Out-Null

$shortcut = $wsh.CreateShortcut((Join-Path $startMenuDir "ChronoTrace.lnk"))
$shortcut.TargetPath = Join-Path $installDir "ChronoTrace.exe"
$shortcut.WorkingDirectory = $installDir
if (Test-Path (Join-Path $installDir "icon.ico")) {
    $shortcut.IconLocation = Join-Path $installDir "icon.ico"
}
$shortcut.Description = "ChronoTrace — Ultra-Lightweight Screen Time Tracker"
$shortcut.Save()

$desktopShortcut = $wsh.CreateShortcut((Join-Path ([Environment]::GetFolderPath("Desktop")) "ChronoTrace.lnk"))
$desktopShortcut.TargetPath = Join-Path $installDir "ChronoTrace.exe"
$desktopShortcut.WorkingDirectory = $installDir
if (Test-Path (Join-Path $installDir "icon.ico")) {
    $desktopShortcut.IconLocation = Join-Path $installDir "icon.ico"
}
$desktopShortcut.Description = "ChronoTrace — Ultra-Lightweight Screen Time Tracker"
$desktopShortcut.Save()

# 6. Register Background Daemon in Startup Registry
$regPath = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Run"
$daemonExe = Join-Path $installDir "chronotrace-daemon.exe"
Set-ItemProperty -Path $regPath -Name "ChronoTraceDaemon" -Value "`"$daemonExe`""

# 7. Start Background Daemon immediately
Start-Process -FilePath $daemonExe

# 8. Start GUI
Start-Process -FilePath (Join-Path $installDir "ChronoTrace.exe")

Write-Host "✓ ChronoTrace installed and running successfully!" -ForegroundColor Green
Write-Host "  - Location: $installDir"
Write-Host "  - Background Tracker registered to start with Windows"
Write-Host "  - Desktop & Start Menu shortcuts created"
