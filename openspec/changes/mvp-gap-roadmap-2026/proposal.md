# Proposal: MVP gap analysis, roadmap y fases (GonsGarage)

## Intent

Congelar en OpenSpec **qué falta** frente a **Entra MVP v1** ([`docs/mvp-solo-checklist.md`](../../../docs/mvp-solo-checklist.md)), **deuda técnica** observable, y un **roadmap ejecutable** con prioridades y fechas orientativas (abril–junio 2026). Complementa [`openspec/changes/mvp-funcionando-plan/proposal.md`](../mvp-funcionando-plan/proposal.md) (fases operativas ya avanzadas: local, contrato, deploy LAN).

## Scope

### In Scope

- Breve auditoría repo vs **Entra** (users, cars, appointments, repairs, **invoices**, **billing**).
- Roadmap por **prioridad** (P0–P3) con **ventanas de fecha**.
- Lista de **deuda técnica** accionable.

### Out of Scope

- Paridad Arnela/Tailwind v4; i18n; PSP; multi-tenant (ya “Fuera” en Fase 1).

## Capabilities

### New Capabilities

- `mvp-accounting`: API + UI mínima **invoices** (cliente; extensible proveedor/recibos) y **billing** cliente, alineado a Entra 1.1 — hoy dominio/servicios parciales **sin rutas en `cmd/api`**.

### Modified Capabilities

- None (sin cambiar requisitos de `client-auth-shell`; solo nuevos dominios).

## Análisis: qué falta al MVP (vs Entra 1.1)

| Entra | Estado repo |
|-------|-------------|
| CRUD users | Hecho (auth, roles, `/auth/me`). |
| CRUD cars / appointments | Hecho API + UI principal. |
| CRUD repairs | Lectura + dominio; **solo `GET` repairs por coche** en Gin; staff **POST/PATCH** pendiente (checklist 6.1 / Fase C opcional). |
| **CRUD invoices** | **Hueco mayor:** código en `internal/service/invoice`, **sin registro en `setupRoutes`**. Sin páginas `frontend` invoices. |
| **CRUD billing** | **Hueco:** sin evidencia de rutas/UI dedicadas en el monorepo actual. |

**Deploy / prod de pruebas:** plantilla Docker + Arnela-network documentada; Fase 4–5 del [checklist](../../../docs/mvp-solo-checklist.md) aún con smoke/backup/CORS a cerrar explícitamente.

## Roadmap (fases, fechas, prioridades)

| ID | Prioridad | Ventana (2026) | Entrega |
|----|-----------|------------------|---------|
| **P0** | Crítica | **18–25 abr** | Cerrar Fase 4–5 checklist: smoke login/coche/cita/repairs en servidor; backup `pg_dump`; revisar `CORS_ORIGINS` + `JWT_SECRET` en release; marcar 4.3–4.4 verificados. |
| **P1** | Alta | **28 abr – 16 may** | **Invoices + billing** primer corte: rutas Gin + repo + Swagger + pantallas mínimas **o** decisión documentada de **recortar** Entra 1.1 (enmienda checklist) si el negocio pospone contabilidad. |
| **P2** | Media | **19 may – 6 jun** | Repairs **staff** escritura (API + UI) si sigue “incluido” en Fase 1; si no, cerrar con “aplazado” + issue. |
| **P3** | Baja | **jun+** | `deploy.yml` real; observabilidad (logs estructurados); política versión API/changelog ([`docs/roadmap.md`](../../../docs/roadmap.md) Fase 4). |

## Deuda técnica (verificada / probable)

| Ítem | Impacto | Mitigación breve |
|------|---------|------------------|
| Cliente HTTP duplicado (`api-client.ts` vs `lib/api.ts`) | Medio | Unificar consumo o documentar “legacy solo dashboard” hasta migración. |
| `docs/roadmap.md` Fase 3 repairs | Bajo | Marcar hecho o enlazar a `application-analysis.md`. |
| `JWT` default débil en dev | Medio | Ya advertido; servidor: solo env. |
| UI mezcla PT/EN | Bajo | Fuera MVP i18n; unificar cuando haya scope. |
| Invoices sin HTTP | **Alto** | P1. |

## Approach

1. **P0:** ejecutar checklist servidor; sin código salvo ajustes CORS/env si falla smoke.  
2. **P1:** spec `mvp-accounting` → rutas REST + migraciones/AutoMigrate + UI listado/detalle mínimo.  
3. **P2/P3:** issues derivados de esta proposal; no bloquear P0.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `backend/cmd/api/main.go` | Modified (P1) | Registrar handlers invoices/billing. |
| `frontend/src/app/` | New (P1) | Rutas UI accounting. |
| `docs/mvp-solo-checklist.md` | Modified | Marcar fases servidor + eventual recorte Entra. |
| `docs/roadmap.md` | Modified | Alinear checkboxes con código. |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Scope invoices/billing subestimado | Med | MVP vertical: 1 lista + 1 detalle + 1 POST staff. |
| Arnela DB compartida | Med | DB dedicada `gonsgarage`; no mezclar datos Arnela. |

## Rollback Plan

Archivar o revertir `openspec/changes/mvp-gap-roadmap-2026/`; deshacer rutas nuevas por PR revert; BD: migraciones reversibles o backup previo a P1.

## Dependencies

- Decisión producto: **¿P1 full invoices/billing o recorte documentado de Entra 1.1?**

## Success Criteria

- [ ] P0: checklist Fase 4–5 verificado en servidor de pruebas con fechas.
- [ ] P1: primera ruta `GET/POST` invoices o documento de recorte firmado en repo.
- [ ] Deuda tabla arriba convertida en **≥3 issues** etiquetados `mvp` / `tech-debt`.
