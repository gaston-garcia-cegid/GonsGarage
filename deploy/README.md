# Despliegue LAN / Docker (plantilla Arnela)

Archivos en la **raíz del repo**:

| Archivo | Uso |
|---------|-----|
| [`docker-compose.prod.yml`](../docker-compose.prod.yml) | Redis + API + Next (standalone) + nginx en **8102**. |
| [`.env.prod.example`](../.env.prod.example) | Plantilla; copiar a **`.env.prod`** en el servidor (no git). |
| [`deploy.ps1`](../deploy.ps1) | `scp` + `docker compose` remoto (ajustar `$REMOTE_DIR`, usuario SSH). |
| [`docker-setup.prod.ps1`](../docker-setup.prod.ps1) | Mismo compose en local para probar antes de subir. |
| [`nginx/default.conf`](../nginx/default.conf) | `/` → front, `/api/` y `/swagger/` → API. |
| [`backend/Dockerfile`](../backend/Dockerfile) | Binario `gonsgarage-api`. |
| [`frontend/Dockerfile`](../frontend/Dockerfile) | Activa `DOCKER_BUILD=1` para `output: "standalone"` (solo en build Linux/Docker; `pnpm build` local en Windows sigue sin standalone). |

## Postgres en el host (Arnela / otro)

1. Crear base `gonsgarage` y usuario con contraseña fuerte en tu Postgres del servidor.
2. En `.env.prod`, `DATABASE_URL` debe apuntar a **`host.docker.internal`** (Linux: ya está en `extra_hosts` del compose) si Postgres escucha en el mismo host que Docker.

## CORS y `GIN_MODE=release`

El API lee **`CORS_ORIGINS`** (coma-separado). Debe coincidir con la URL exacta del navegador (p. ej. `http://192.168.1.100:8102`, sin barra final).

## Checklist MVP

Marcar tareas 4.x en [`docs/mvp-solo-checklist.md`](../docs/mvp-solo-checklist.md) cuando verifiques URLs y secretos en el servidor real.
