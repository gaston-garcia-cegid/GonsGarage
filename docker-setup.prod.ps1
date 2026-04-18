# GonsGarage — gestión del stack prod en la máquina actual (compose prod).
# Uso: .\docker-setup.prod.ps1 [build|up|down|restart|logs|status]
# Requiere `.env.prod` en la raíz del repo (copiar de `.env.prod.example`).

param(
    [Parameter(Mandatory = $false)]
    [ValidateSet('build', 'up', 'down', 'restart', 'logs', 'status')]
    [string]$Command = 'up'
)

$ErrorActionPreference = "Stop"
$ComposeFile = "docker-compose.prod.yml"
$EnvFile = ".env.prod"

if (-not (Test-Path $EnvFile)) {
    Write-Host "Falta $EnvFile — copiá .env.prod.example y completá DATABASE_URL / CORS / JWT." -ForegroundColor Red
    exit 1
}

function Invoke-Compose {
    param([string[]]$Args)
    docker compose -f $ComposeFile --env-file $EnvFile @Args
}

Write-Host "=== GonsGarage Docker (prod compose) ===" -ForegroundColor Green

switch ($Command) {
    'build' {
        Write-Host "Building..." -ForegroundColor Cyan
        Invoke-Compose @("build")
    }
    'up' {
        Write-Host "Starting..." -ForegroundColor Cyan
        Invoke-Compose @("up", "-d")
        Start-Sleep -Seconds 3
        Invoke-Compose @("ps")
    }
    'down' {
        Invoke-Compose @("down")
    }
    'restart' {
        Invoke-Compose @("restart")
    }
    'logs' {
        Invoke-Compose @("logs", "-f")
    }
    'status' {
        Invoke-Compose @("ps")
        try {
            $r = Invoke-RestMethod -Uri "http://localhost:8102/health" -Method Get -TimeoutSec 5
            $r | ConvertTo-Json
        }
        catch {
            Write-Host "Health en :8102 no accesible (¿nginx arriba?)." -ForegroundColor Yellow
        }
    }
}

Write-Host "=== Hecho ===" -ForegroundColor Green
