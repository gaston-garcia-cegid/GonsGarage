# Despliegue staging / producción (Fase D)

Objetivo: **secretos fuera del repositorio**, artefactos **Docker** reproducibles y un workflow de **deploy** que al menos construya la API y ejecute un **smoke test** (`/health`).

## Secretos y variables

| Variable | Dónde | Notas |
|----------|--------|--------|
| `JWT_SECRET` | Solo entorno / secret manager | En `GIN_MODE=release` la API **falla al arrancar** si falta o sigue siendo la clave de desarrollo por defecto. |
| `DATABASE_URL` | Entorno del contenedor `api` | DSN PostgreSQL; en Compose el host suele ser el nombre del servicio (`postgres`). |
| `POSTGRES_*` | `docker-compose.prod.yml` | Definen la base del contenedor Postgres; deben coincidir con el usuario/clave del `DATABASE_URL`. |
| `NEXT_PUBLIC_API_URL` | Build args de `web` | URL **pública** del API (HTTPS recomendado detrás de proxy). |
| `CORS_ALLOWED_ORIGINS` | Contenedor `api` | Con `GIN_MODE=release`, lista separada por **comas** de orígenes exactos del SPA (ej. `https://app.tudominio.com`). Vacío: solo encaja bien con clientes sin cabecera `Origin` (curl, server-to-server) o mismo sitio. |

### CORS en `release`

El middleware en `cmd/api/main.go` solo envía `Access-Control-Allow-Origin` si el `Origin` del navegador coincide con una entrada de `CORS_ALLOWED_ORIGINS`. Fuera de `release` el comportamiento sigue siendo permisivo para desarrollo local.

Copia `deploy/.env.production.example` a **`.env` en la raíz del repo** (o al path que pases con `docker compose --env-file`) y rellena valores reales. **No subas `.env` a git.**

## TLS (HTTPS)

Los contenedores exponen HTTP. En producción se recomienda **Caddy**, **Traefik** o **nginx** delante como reverse proxy con certificados (Let’s Encrypt u otro). La “URL pública con HTTPS” es responsabilidad de ese proxy y del DNS.

## Backup de PostgreSQL

Ejemplo manual (desde el host con `docker compose` levantado):

```bash
docker compose -f docker-compose.prod.yml exec -T postgres \
  pg_dump -U "$POSTGRES_USER" "$POSTGRES_DB" > backup-$(date +%Y%m%d).sql
```

Restauración (esquema simplificado; probar antes en entorno aislado):

```bash
cat backup-YYYYMMDD.sql | docker compose -f docker-compose.prod.yml exec -T postgres \
  psql -U "$POSTGRES_USER" -d "$POSTGRES_DB"
```

Automatiza copias (cron, volúmenes snapshot, PaaS managed) según tu plataforma.

## Archivos relevantes

- `backend/Dockerfile` — imagen API.
- `frontend/Dockerfile` — Next.js **standalone** (`NEXT_STANDALONE=true` en la etapa de build; en Windows local sin symlinks dejar `pnpm build` sin esa variable, ver `next.config.ts`).
- `docker-compose.prod.yml` — API + Postgres + Redis + frontend.
- `docker-compose.smoke.yml` — mismo patrón con credenciales fijas solo para CI / pruebas locales.
- `.github/workflows/deploy.yml` — build + smoke; opcional push a GHCR.

## Rollback (básico)

1. Mantén etiquetas de imagen por versión (`:v1.2.3` o digest SHA).
2. Ante regresión: `docker compose ... pull` (si usas registry) o vuelve a desplegar la imagen anterior y `docker compose up -d`.
3. Si hubo migraciones destructivas, restaura backup de BD antes del deploy fallido.
