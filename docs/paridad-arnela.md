# Paridad explícita con el stack Arnela (referencia)

Esta tabla resume qué replica GonsGarage frente al README/stack Arnela y dónde vivir en el repo.

| Arnela (referencia) | GonsGarage |
| --- | --- |
| Next.js 16 + TypeScript | `frontend/package.json` (`next@16`, `pnpm`) |
| Tailwind CSS v4 | `frontend/postcss.config.mjs`, `tailwindcss@4`, `@import "tailwindcss"` en `src/app/globals.css` |
| Shadcn / Radix (base) | `frontend/components.json`, `src/components/ui/button.tsx`, `@radix-ui/react-slot`, `cva`, `cn` en `src/lib/utils.ts` |
| Zustand | Sin cambio (`zustand` en frontend) |
| PostgreSQL 16 | `docker-compose*.yml` (`postgres:16-alpine`) |
| golang-migrate (15 revisiones SQL de ejemplo) | `backend/db/migrations/*.sql`, `internal/sqlmigrate`, flags `ENABLE_SQL_MIGRATIONS`, `MIGRATIONS_DIR`, `SKIP_GORM_AUTOMIGRATE` |
| Redis + worker async (BRPOP) | `backend/cmd/worker`, `Dockerfile.worker`, cola `WORKER_QUEUE_KEY` |
| Docker Compose + Nginx | `docker-compose.prod.yml` servicio `edge`, `deploy/nginx-gonsgarage.conf` |
| JWT + roles (API) | Handlers + middleware existentes |
| Rate limit rutas públicas de auth | `internal/adapters/http/middleware/ratelimit.go`, `AUTH_RATE_LIMIT_*` |
| OpenAPI / Swagger | `swag` + `/swagger/*` (proxy en Nginx) |

Producción Compose: el host publica solo **Nginx** (`HTTP_PUBLISH_PORT`). Configurar `NEXT_PUBLIC_API_URL` con el mismo origen que el navegador (p. ej. `http://localhost:8080/api/v1` si el edge escucha en 8080) y `CORS_ALLOWED_ORIGINS` con el origen del SPA.
