# GonsGarage — despliegue remoto (plantilla Arnela).
# Ajustá $SERVER, $USER, $REMOTE_DIR antes de usar.

$SERVER = "192.168.1.100"
$USER = "root"
$REMOTE_DIR = "/DATA/AppData/gonsgarage"

$ErrorActionPreference = "Stop"

Write-Host "=== GonsGarage deploy ===" -ForegroundColor Cyan
Write-Host ""

$items = @(
    "docker-compose.prod.yml",
    ".env.prod.example",
    "backend",
    "frontend",
    "nginx"
)

Write-Host "[1/4] Crear directorio remoto..." -ForegroundColor Yellow
ssh "${USER}@${SERVER}" "mkdir -p ${REMOTE_DIR}"

Write-Host "[2/4] Copiar artefactos..." -ForegroundColor Yellow
foreach ($item in $items) {
    $localPath = Join-Path $PSScriptRoot $item
    if (Test-Path $localPath) {
        Write-Host "  -> $item" -ForegroundColor Gray
        scp -r $localPath "${USER}@${SERVER}:${REMOTE_DIR}/"
    }
    else {
        Write-Host "  [SKIP] $item no encontrado" -ForegroundColor Red
    }
}

Write-Host "[3/4] .env.prod en el servidor (si no existe)..." -ForegroundColor Yellow
ssh "${USER}@${SERVER}" "cd ${REMOTE_DIR} && if [ ! -f .env.prod ]; then cp .env.prod.example .env.prod && echo '.env.prod creado desde example — editar JWT_SECRET y DATABASE_URL antes del up'; else echo '.env.prod ya existe'; fi"

Write-Host "[4/4] build + up (Docker)..." -ForegroundColor Yellow
ssh "${USER}@${SERVER}" @"
cd ${REMOTE_DIR}
docker compose -f docker-compose.prod.yml --env-file .env.prod up -d --build
"@

Write-Host ""
Write-Host "=== Listo ===" -ForegroundColor Green
Write-Host "App:     http://${SERVER}:8102" -ForegroundColor Cyan
Write-Host "Health:  http://${SERVER}:8102/health" -ForegroundColor Cyan
Write-Host "Swagger: http://${SERVER}:8102/swagger/index.html" -ForegroundColor Cyan
