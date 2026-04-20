# Proposal: MVP post‑v1 — estado y siguientes pasos

## Intent

El **MVP v1 operativo** está cerrado en checklist (fases **1–5** hechas; **6** con repairs staff + política deploy doc). Falta una **foto clara** de qué queda como **deuda consciente**, qué es **opcional**, y en qué **orden** abordarlo para no mezclar “prod hardening” con “paridad Arnela” ni con **accounting** ya diferido.

## Estado resumido (abril 2026)

| Área | Estado |
|------|--------|
| Alcance MVP v1 | Congelado en `docs/mvp-solo-checklist.md` (invoices/billing fuera; repairs staff dentro). |
| Contrato API + docs | Alineados (`application-analysis.md`, Swagger, seeds `seed-test-client` + `seed-mvp-users`). |
| Verificación por rol | Spec **`mvp-role-access`** en `openspec/specs/mvp-role-access/spec.md`; tests HTTP en `handler/mvp_role_access_test.go`. |
| Deploy remoto | Documentado (`deploy/README.md`); **CI → servidor** sigue manual u opcional (checklist Fase 6). |
| Accounting | Diferido (`openspec/specs/p1-accounting-defer/spec.md`). |
| Roadmap repo | Fases 1–3 mayormente [x]; quedan ítems opcionales Fase 0–2–4–5 (matriz Arnela, compose “todo en contenedor”, changelog API, observabilidad, endurecimiento prod). |

## Scope

### In Scope

- **Documentar** en este change (y enlazar desde `mvp-solo-checklist` o `roadmap.md`) una **tabla priorizada** de siguientes pasos: P0 (rápido / riesgo), P1 (valor producto), P2 (paridad Arnela / nice‑to‑have).
- **Cerrar deuda SDD**: archivar changes activos **`readme-verified-refresh`** y **`mvp-funcionando-plan`** si siguen solo como auditoría.
- **Issues sugeridos**: test `GET /employees` con JWT **admin**; deploy opcional; revisión matriz Arnela (`roadmap.md` Fase 0).

### Out of Scope

- Invoices/billing HTTP/UI (defer). Paridad Arnela completa en un sprint.

## Capabilities

### New Capabilities

- **None** (planificación y docs; sin nuevo dominio normativo obligatorio).

### Modified Capabilities

- **None** (si más adelante se promueve una fila a spec, será un change dedicado).

## Approach

1. Subsección **«Siguientes pasos (post‑MVP v1)»** en checklist **o** `docs/mvp-next-steps.md` enlazado.
2. Pasada **`docs/roadmap.md`**: pendientes reales vs ya cubierto (MVP + `mvp-role-access`).
3. Archivar OpenSpec **`readme-verified-refresh`** y **`mvp-funcionando-plan`** cuando corresponda.

## Affected Areas

| Area | Impact |
|------|--------|
| `docs/mvp-solo-checklist.md` o `docs/mvp-next-steps.md` | Modified / Create |
| `docs/roadmap.md` | Modified (aclarar pendientes reales) |
| `openspec/changes/archive/2026-04-20-readme-verified-refresh`, `.../2026-04-20-mvp-funcionando-plan` | Archivados (apply 2026-04-20) |

## Risks

| Risk | L | Mitigation |
|------|---|------------|
| Lista de pasos demasiado vaga | M | Cada fila con criterio “hecho” observable |

## Rollback Plan

Revertir commits de docs; no hay migración.

## Dependencies

- Ninguna externa.

## Success Criteria

- [x] Tabla o doc **«Siguientes pasos»** con ≥5 ítems priorizados (P0/P1/P2) enlazada desde el checklist o roadmap.
- [x] `readme-verified-refresh` y `mvp-funcionando-plan` archivados **o** justificación explícita en el doc si se mantienen activos.
- [x] Issue o tarea suelta creada (en tracker del equipo) para **admin + GET /employees** si no existe (checklist en `docs/mvp-next-steps.md`; test de regresión en repo).
