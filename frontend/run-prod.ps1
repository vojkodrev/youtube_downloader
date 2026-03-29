$ErrorActionPreference = "Stop"

$NginxExe = "C:\tools\nginx-1.29.7\nginx.exe"
$ConfFile = "$PSScriptRoot\nginx\nginx.conf"

# Build
Write-Host "Building frontend..." -ForegroundColor Cyan
Push-Location $PSScriptRoot
npm run build
Pop-Location

# Stop existing nginx if running
$existing = Get-Process nginx -ErrorAction SilentlyContinue
if ($existing) {
    Write-Host "Stopping existing nginx..." -ForegroundColor Yellow
    & $NginxExe -p $PSScriptRoot -c $ConfFile -e nginx/nginx.log -s stop
    $existing | Wait-Process -Timeout 10
}

# Start nginx (foreground — Ctrl+C to stop)
Write-Host "Starting nginx on http://localhost ..." -ForegroundColor Green
Start-Process -NoNewWindow $NginxExe -ArgumentList "-p `"$PSScriptRoot`" -c `"$ConfFile`" -e nginx/nginx.log"
Write-Host "Finished" -ForegroundColor Green
