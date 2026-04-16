# Roadmap: documentación, configuración y alineación tipo Arnela

Este roadmap ordena el trabajo en fases. La parte “como Arnela” depende de completar las especificaciones reales en [arnela-specs.md](./arnela-specs.md).

## Fase 0 — Especificaciones Arnela (bloqueante para paridad)

- [x] Ruta local confirmada: **`D:\Repos\Arnela`** (README, `arnela-rules/`, `docs/DOCUMENTATION_INDEX.md`).
- [x] Resumen en GonsGarage: [specs/arnela/ARNELA_SYNOPSIS.md](./specs/arnela/ARNELA_SYNOPSIS.md) e índice [specs/arnela/README.md](./specs/arnela/README.md); matriz en [arnela-specs.md](./arnela-specs.md).
- [ ] Revisión conjunta con el equipo: priorizar filas de la matriz (compose, CI, auth, migraciones) y convertir en issues/tareas.

## Fase 1 — Documentación y descubribilidad (este repo)

- [x] Carpeta `docs/` con índice, análisis de aplicación y guía de desarrollo.
- [x] README raíz y Swagger: comandos alineados con `cmd/api` y `docker-compose.yml` en raíz; `deployment/README.md` apunta al compose unificado.
- [x] `frontend/README.md`: enlaces a guía de desarrollo y `.env.local.example`.
- [ ] Opcional: mover notas sueltas (`frontend_migration_plan.md`, etc.) a `docs/history/` con un índice breve.

## Fase 2 — Configuración unificada (objetivo: un solo “camino feliz”)

- [x] `docker-compose.yml` en raíz: **PostgreSQL + Redis** alineados con defaults de `cmd/api` y `backend/.env.example`.
- [x] `backend/.env.example` y `frontend/.env.local.example`.
- [x] [development-guide.md](./development-guide.md) actualizado al flujo actual.
- [ ] Opcional: perfiles compose para construir/ejecutar **backend + frontend** en contenedor (no requerido para dev local).

## Fase 3 — Paridad funcional y API

- [ ] Exponer rutas REST de **repairs** o documentar explícitamente el alcance MVP sin ellas.
- [x] Parte de permisos por rol: **`/employees/*`** solo **admin/manager**; **`/auth/me`** para sesión; registro con rol por defecto **client** (ver [mvp-minimum-phases.md](./mvp-minimum-phases.md) Fase B).
- [ ] Sincronizar Swagger / tipos del frontend con los endpoints reales.

## Fase 4 — Prácticas de ingeniería (espejo de Arnela, según Fase 0)

Ajustar según lo que marque la matriz Arnela vs GonsGarage. Candidatos típicos:

- [x] Pipeline CI: `.github/workflows/ci.yml` (Go + vet + test con CGO; pnpm install + lint + typecheck + test + build frontend).
- [x] Workflow deploy placeholder: `.github/workflows/deploy.yml` (`workflow_dispatch`).
- [ ] Política de versionado de API y changelog.
- [ ] Plantillas de PR / issues (si Arnela las usa).
- [ ] Observabilidad: logs estructurados, health agregado, métricas (si aplica).

## Fase 5 — Producción y seguridad

- [ ] Secretos solo por variables de entorno; eliminar defaults inseguros en producción.
- [ ] CORS restrictivo en `release` (revisar `corsMiddleware` en `main.go`).
- [ ] Imagen Docker de backend/frontend y documentación de despliegue acorde al entorno real.

---

**Prioridad sugerida**: Fase 1–2 en paralelo con el cierre de Fase 0 (priorización de la matriz). Las fases 3–5 se refinan con [arnela-specs.md](./arnela-specs.md) y el código en `D:\Repos\Arnela`.
