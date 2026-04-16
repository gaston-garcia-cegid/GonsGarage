# Roadmap: documentación, configuración y alineación tipo Arnela

Este roadmap ordena el trabajo en fases. La parte “como Arnela” depende de completar las especificaciones reales en [arnela-specs.md](./arnela-specs.md).

## Estado alineado con MVP (A–D)

El MVP mínimo descrito en [mvp-minimum-phases.md](./mvp-minimum-phases.md) está **cerrado en repo** (compose, CI/TDD, identidad y coches/citas, **repairs**, Docker + deploy smoke + guía de staging). Correspondencia aproximada con este roadmap:

| MVP        | Roadmap (este archivo)                          |
| ---------- | ----------------------------------------------- |
| A–B        | Fases 1–2 (docs + compose + guía) y gran parte de Fase 3 |
| C          | Fase 3 (repairs + contrato API / Swagger)       |
| D          | Fases 4–5 (CI/deploy, Docker, secretos, CORS)  |

**Tras MVP A–D — orden sugerido:** cerrar pendientes de **Fase 0** (priorización matriz Arnela), luego **Fase 4** (versionado API, plantillas, observabilidad) y **Fase 5** (endurecimiento continuo de seguridad y despliegue real con URL/HTTPS operados por el equipo).

## Fase 0 — Especificaciones Arnela (bloqueante para paridad)

- [x] Ruta local confirmada: **`D:\Repos\Arnela`** (README, `arnela-rules/`, `docs/DOCUMENTATION_INDEX.md`).
- [x] Resumen en GonsGarage: [specs/arnela/ARNELA_SYNOPSIS.md](./specs/arnela/ARNELA_SYNOPSIS.md) e índice [specs/arnela/README.md](./specs/arnela/README.md); matriz en [arnela-specs.md](./arnela-specs.md).
- [ ] Revisión conjunta con el equipo: priorizar filas de la matriz (compose, CI, auth, migraciones) y convertir en issues/tareas. **Checklist:** [phase-0-team-review-checklist.md](./specs/arnela/phase-0-team-review-checklist.md).

## Fase 1 — Documentación y descubribilidad (este repo)

- [x] Carpeta `docs/` con índice, análisis de aplicación y guía de desarrollo.
- [x] README raíz y Swagger: comandos alineados con `cmd/api` y `docker-compose.yml` en raíz; `deployment/README.md` apunta al compose unificado.
- [x] `frontend/README.md`: enlaces a guía de desarrollo y `.env.local.example`.
- [x] Notas sueltas: `frontend_migration_plan.md` archivado en [docs/history/](./history/README.md) con índice breve.

## Fase 2 — Configuración unificada (objetivo: un solo “camino feliz”)

- [x] `docker-compose.yml` en raíz: **PostgreSQL + Redis** alineados con defaults de `cmd/api` y `backend/.env.example`.
- [x] `backend/.env.example` y `frontend/.env.local.example`.
- [x] [development-guide.md](./development-guide.md) actualizado al flujo actual.
- [x] Compose **prod** + **smoke** + Dockerfiles API/web y ejemplo de env: `docker-compose.prod.yml`, `docker-compose.smoke.yml`, `deploy/.env.production.example`, [deployment-staging.md](./deployment-staging.md).

## Fase 3 — Paridad funcional y API

- [x] Rutas REST de **repairs** (`/api/v1/repairs`, `/api/v1/cars/:id/repairs`) y UI de historial por coche (ver MVP Fase C).
- [x] Parte de permisos por rol: **`/employees/*`** solo **admin/manager**; **`/auth/me`** para sesión; registro con rol por defecto **client** (ver [mvp-minimum-phases.md](./mvp-minimum-phases.md) Fase B).
- [x] Swagger regenerado con **swag** y contratos **camelCase** alineados en los flujos activos (coches, citas, repairs, auth); el frontend usa varios clientes (`api-client`, `api.ts`); mantener coherencia en cada PR que toque API.

## Fase 4 — Prácticas de ingeniería (espejo de Arnela, según Fase 0)

Ajustar según lo que marque la matriz Arnela vs GonsGarage. Candidatos típicos:

- [x] Pipeline CI: `.github/workflows/ci.yml` (Go + vet + test con CGO; pnpm install + lint + typecheck + test + build frontend con `NEXT_STANDALONE` en Linux).
- [x] Workflow deploy: `.github/workflows/deploy.yml` — smoke Docker Compose + `/health`, push opcional a GHCR (`push_ghcr`); ver [deployment-staging.md](./deployment-staging.md).
- [x] Política de versionado de API y changelog: [docs/api/versioning.md](./api/versioning.md) y [CHANGELOG.md](../CHANGELOG.md) en la raíz; `apiVersion` en `GET /health`.
- [x] Plantilla de PR en [.github/pull_request_template.md](../.github/pull_request_template.md).
- [x] Plantillas de issues: [.github/ISSUE_TEMPLATE/](../.github/ISSUE_TEMPLATE/) (`bug_report`, `feature_request`).
- [x] Observabilidad mínima: `GET /health` (liveness + `apiVersion`) y `GET /ready` (readiness Postgres).
- [x] Observabilidad ampliada (base): logs estructurados con **`log/slog`** (`LOG_FORMAT` / `LOG_LEVEL`, JSON en `release`), middleware de acceso, **`GET /metrics`** (Prometheus). Guía: [observability.md](./observability.md). OpenTelemetry / métricas de negocio: pendiente de priorizar.

## Fase 5 — Producción y seguridad

- [x] Secretos por variables de entorno; `JWT_SECRET` obligatorio (no default de desarrollo) cuando `GIN_MODE=release` (`cmd/api/main.go`).
- [x] CORS restrictivo en `release`: lista blanca vía `CORS_ALLOWED_ORIGINS` (coma-separada); en modo no `release` se mantiene comportamiento permisivo para desarrollo local.
- [x] Imágenes Docker de backend/frontend y documentación de despliegue: [deployment-staging.md](./deployment-staging.md), `docker-compose.prod.yml`, Dockerfiles en `backend/` y `frontend/`.

---

**Prioridad sugerida:** ejecutar la **Fase 0** con la [checklist de reunión](./specs/arnela/phase-0-team-review-checklist.md); después, paridad visual/stack Arnela según [arnela-specs.md](./arnela-specs.md) y el código en `D:\Repos\Arnela`. La base de **observabilidad** (logs + `/metrics`) ya está en [observability.md](./observability.md); OpenTelemetry y métricas de negocio cuando el equipo las priorice.
