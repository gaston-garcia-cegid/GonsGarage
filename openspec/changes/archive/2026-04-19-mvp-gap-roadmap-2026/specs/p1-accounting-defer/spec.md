# Delta spec: P1 accounting — decisión de aplazamiento

## Context

La proposal [`mvp-gap-roadmap-2026`](../../proposal.md) prevé P1 como **rutas Gin + Swagger + UI mínima** para invoices/billing **o** una **enmienda explícita** al checklist Entra 1.1.

## Requirements

### R-1 Enmienda de alcance documentada

**The system documentation MUST** reflejar que **CRUD invoices** y **CRUD billing** quedan **fuera de la primera entrega operativa** del MVP v1 (2026-04-17), manteniendo la visión de producto en backlog hasta un change dedicado (`mvp-accounting` o equivalente).

### R-2 Tareas condicionales 3.2 / 3.3

**The OpenSpec tasks MUST** cerrarse como **N/A** respecto a implementación HTTP/UI de accounting mientras R-1 aplique (no se registran handlers nuevos de invoices en `cmd/api` por este change).

## Scenarios

- **S-1** — Un lector del checklist ve los bullets Entra 1.1 con invoices/billing marcados como diferidos y enlace al spec delta.
- **S-2** — `tasks.md` marca 3.1 hecha y 3.2–3.3 como completadas con nota N/A alineada a este spec.
