# §1 — Arquitectura

- **Backend:** monolito modular en **Go**, **Clean Architecture**: transporte (`internal/handler`, `cmd/`, `internal/middleware`), reglas (`internal/service`), entidades (`internal/domain`), persistencia (`internal/repository` + **PostgreSQL**), utilidades (`pkg/`). *Fase 1: fachadas sobre `internal/adapters` y `internal/core` — ver `docs/template-adoption-plan.md`.*
- **Frontend:** **Next.js (App Router)** + **TypeScript**, app monolítica por segmentos (público, cliente, staff). Estado global con **Zustand**.
- **Datos:** PostgreSQL como fuente de verdad; **Redis** para caché, rate limiting y/o cola async (emails, notificaciones, jobs).
- **API:** REST JSON bajo prefijo estable (p. ej. `/api/v1`); contrato documentado con **OpenAPI/Swagger** generado desde el backend.
- **Despliegue dev:** **Docker Compose** (Postgres + Redis mínimo); opcional `docker-compose.prod.yml` + reverse proxy (Nginx).
- **CI:** GitHub Actions — backend `vet` → `build` → `test -race`; frontend `lint` → `tsc --noEmit` → **Vitest** → `build`.

**Dominio de negocio (`{{BUSINESS_DOMAIN}}`):** auto repair shop management system (tabla en `template_project.md`).
