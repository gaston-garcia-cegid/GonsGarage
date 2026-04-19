# Tasks: MVP gap roadmap 2026

## Phase 1: Checklist MVP Fase 1.1–1.3 (alcance congelado — trazabilidad)

> Criterios 1.1–1.3 ya **Hecho** en [`docs/mvp-solo-checklist.md`](../../../docs/mvp-solo-checklist.md) (2026-04-17). Esta fase evita **drift** entre checklist, gap proposal y fases mínimas.

- [x] 1.1 Auditar lista **Entra** / **Fuera** / **Repairs staff** en `docs/mvp-solo-checklist.md#decisiones-cerradas-mvp-v1` vs tabla §Análisis en `openspec/changes/mvp-gap-roadmap-2026/proposal.md`; sin divergencia de ítems Entra (2026-04-18).
- [x] 1.2 Añadir en `proposal.md` §Dependencies enlace explícito al ancla `#decisiones-cerradas-mvp-v1` del checklist.
- [x] 1.3 Referenciar `mvp-gap-roadmap-2026/proposal.md` desde `docs/mvp-minimum-phases.md` (línea bajo plan consolidado).

## Phase 2: P0 — Servidor de pruebas y endurecimiento (checklist Fase 4–5)

- [ ] 2.1 Smoke en URL real: login + coche + cita + repairs lectura; anotar fechas en `docs/mvp-solo-checklist.md` filas 4.3–4.4 si aún dicen “plantilla”.
- [ ] 2.2 `pg_dump` o política backup documentada (checklist 5.3); un comando en `deploy/README.md` o runbook propio.
- [ ] 2.3 Revisión CORS + `JWT_SECRET` en `.env.prod` servidor (`GIN_MODE=release`); checklist 5.1–5.2.

## Phase 3: P1 — Invoices / billing o recorte Entra

- [ ] 3.1 Decisión escrita: implementar `mvp-accounting` (rutas Gin + Swagger + UI mínima) **o** enmienda en `mvp-solo-checklist.md` recortando Entra 1.1.
- [ ] 3.2 Si implementa: registrar handlers en `backend/cmd/api/main.go`; spec delta bajo `openspec/changes/mvp-gap-roadmap-2026/specs/` (sdd-spec).
- [ ] 3.3 Si implementa: rutas App Router bajo `frontend/src/app/` (listado/detalle); `pnpm test` / smoke manual.

## Phase 4: P2 — Repairs staff (si sigue en alcance 1.3)

- [ ] 4.1 `POST`/`PATCH`/`DELETE` repairs en Gin + permisos staff; swag regenerado `backend/docs/`.
- [ ] 4.2 UI mínima staff (ruta existente empleados o nueva); alinear `docs/application-analysis.md`.

## Phase 5: P3 — CI deploy y docs

- [ ] 5.1 Sustituir placeholder `.github/workflows/deploy.yml` o documentar “solo manual”.
- [ ] 5.2 Actualizar `docs/roadmap.md` checkboxes obsoletos (repairs GET, Swagger).
