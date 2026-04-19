# Spec: P1 accounting — decisión de aplazamiento (MVP v1)

> **Promoción:** incorporado al catálogo principal desde el change archivado `openspec/changes/archive/2026-04-19-mvp-gap-roadmap-2026/` (2026-04-19). Historial: ver `proposal.md` en esa carpeta.

## Context

El roadmap `mvp-gap-roadmap-2026` prevé P1 como **rutas Gin + Swagger + UI mínima** para invoices/billing **o** una **enmienda explícita** al checklist Entra 1.1. La decisión tomada fue **aplazar** accounting HTTP/UI y documentarlo en el checklist del MVP y en tareas OpenSpec.

## Requirements

### R-1 Enmienda de alcance documentada

**The system documentation MUST** reflejar que **CRUD invoices** y **CRUD billing** quedan **fuera de la primera entrega operativa** del MVP v1 (2026-04-17), manteniendo la visión de producto en backlog hasta un change dedicado (`mvp-accounting` o equivalente).

### R-2 Tareas condicionales 3.2 / 3.3

**The OpenSpec tasks** (en el change archivado) **MUST** haberse cerrado como **N/A** respecto a implementación HTTP/UI de accounting mientras R-1 aplique (sin handlers nuevos de invoices en `cmd/api` por ese change).

## Scenarios

- **S-1** — Un lector del checklist ve los bullets Entra 1.1 con invoices/billing marcados como diferidos y enlace al spec / archivo de change archivado.
- **S-2** — Las tareas 3.1–3.3 del change archivado reflejan 3.1 hecha y 3.2–3.3 N/A alineadas a este spec.
