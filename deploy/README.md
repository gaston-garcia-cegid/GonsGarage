# Despliegue LAN / Docker (plantilla Arnela)

## Orden recomendado (desde tu PC con `deploy.ps1`)

1. **`git pull`** en el clon local (`D:\Repos\GonsGarage` o el que uses), para que lo que subís sea la última versión.
2. **`.\deploy.ps1`** en **PowerShell** desde la raíz del repo (no copia el `.git`; solo `backend/`, `frontend/`, `nginx/`, compose y example de env).

Si trabajás **solo en el servidor** con una carpeta que rellenás a mano o con `git clone` allí, entonces el equivalente es **`git pull`** en `/DATA/AppData/gonsgarage` (si usás git en el servidor) y después `docker compose … up -d --build` — no hace falta `deploy.ps1` en ese flujo.

---

Archivos en la **raíz del repo**:

| Archivo | Uso |
|---------|-----|
| [`docker-compose.prod.yml`](../docker-compose.prod.yml) | Redis + API + Next (standalone) + nginx en **8102**. |
| [`.env.prod.example`](../.env.prod.example) | Plantilla; copiar a **`.env.prod`** en el servidor (no git). |
| [`deploy.ps1`](../deploy.ps1) | `scp` + `docker compose` remoto. Por defecto `$REMOTE_DIR` = `/DATA/AppData/gonsgarage` (cambiar en el script si hace falta). |
| [`docker-setup.prod.ps1`](../docker-setup.prod.ps1) | Mismo compose en local para probar antes de subir. |
| [`nginx/default.conf`](../nginx/default.conf) | `/` → front, `/api/` y `/swagger/` → API. |
| [`backend/Dockerfile`](../backend/Dockerfile) | Binario `gonsgarage-api`. |
| [`frontend/Dockerfile`](../frontend/Dockerfile) | Activa `DOCKER_BUILD=1` para `output: "standalone"` (solo en build Linux/Docker; `pnpm build` local en Windows sigue sin standalone). |

## Postgres en el host (Arnela / otro)

1. Crear base `gonsgarage` y usuario con contraseña fuerte en tu Postgres del servidor.
2. En `.env.prod`, `DATABASE_URL` debe apuntar a **`host.docker.internal`** (Linux: ya está en `extra_hosts` del compose) si Postgres escucha en el mismo host que Docker.

## CORS y `GIN_MODE=release`

El API lee **`CORS_ORIGINS`** (coma-separado). Debe coincidir con la URL exacta del navegador (p. ej. `http://192.168.1.100:8102`, sin barra final).

## Secretos (checklist 4.2)

- En el servidor solo **`/DATA/AppData/gonsgarage/.env.prod`** (o la ruta que uses); **no** subir `.env.prod` al git (está en [`.gitignore`](../.gitignore)).
- Obligatorios: **`JWT_SECRET`** (largo, aleatorio), **`DATABASE_URL`** (usuario dedicado, no el superusuario de Postgres).
- **`REDIS_URL`** en el compose prod apunta al contenedor Redis interno; no hace falta secret extra para Redis salvo que cambies el diseño.

## Verificación tras el primer `up` (4.3 / 4.4)

Desde cualquier máquina en la LAN (sustituí host/puerto si cambiaron):

```bash
curl -sS "http://192.168.1.100:8102/health"
curl -sS "http://192.168.1.100:8102/ready"
```

- **4.3:** respuestas JSON OK del API vía nginx.
- **4.4:** abrir `http://192.168.1.100:8102` en el navegador, **login** contra el mismo origen (el front ya fue construido con `NEXT_PUBLIC_API_URL` igual a esa URL pública).

Si `/ready` falla, Postgres no es alcanzable desde el contenedor del API (revisá `DATABASE_URL`, firewall y `host.docker.internal`).

## HTTPS vs HTTP en LAN (4.5)

Para **solo red local** (`192.168.x.x`) es aceptable **HTTP** en el checklist: sin certificado, sin Let’s Encrypt. Si exponés el servidor a **internet**, añadí TLS (nginx + certbot, o TLS en el proxy) y actualizá **`CORS_ORIGINS`** / **`NEXT_PUBLIC_API_URL`** a `https://…`.

## Rollback (4.6)

1. **Antes de cada deploy:** copia del `.env.prod` vigente y, si aplica, volcado rápido de BD (`pg_dump -Fc gonsgarage > backup.dump`).
2. **Volver atrás en contenedores:** en el directorio del deploy, `docker compose -f docker-compose.prod.yml --env-file .env.prod down`, restaurá la carpeta/código anterior (o `git checkout` + rebuild) y `… up -d --build`.
3. **Solo API (binario):** sustituí el binario por la versión anterior y reiniciá el proceso; la base suele seguir compatible salvo que hayas dependido de migraciones destructivas (evitar `RESET_DATABASE` en servidor).

## Checklist MVP

Cerrar tareas **4.x** en [`docs/mvp-solo-checklist.md`](../docs/mvp-solo-checklist.md) cuando hayas hecho la **verificación en servidor** (curl + login). Este documento cubre criterios **4.2–4.6** a nivel de runbook; **4.1** son las URLs ya anotadas en el checklist.
