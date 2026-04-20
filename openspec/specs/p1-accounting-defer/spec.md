# Spec: P1 accounting — decisión de aplazamiento (MVP v1) + programa P1

> **Promoción:** texto base desde `openspec/changes/archive/2026-04-19-mvp-gap-roadmap-2026/` (2026-04-19). **Actualización 2026-04-20:** incorporado delta do change archivado `openspec/changes/archive/2026-04-20-p1-invoices-billing-suppliers/` — glosario P1 e ligazón aos specs `invoices`, `billing`, `suppliers` no catálogo principal.

## Context

El roadmap `mvp-gap-roadmap-2026` prevé P1 como **rutas Gin + Swagger + UI mínima** para invoices/billing **o** una **enmienda explícita** al checklist Entra 1.1. A decisión tomada para o **MVP v1** foi **aplazar** accounting HTTP/UI e documentalo no checklist. A **implementación P1** (facturas recibidas, billing emitido, suppliers) realizouse no change archivado **`2026-04-20-p1-invoices-billing-suppliers`**; glosario: **invoices** = documentos **recibidos** polo taller; **billing** = documentos **emitidos** polo taller.

## Requirements

### Requirement: Programme P1 tracking

**The system documentation MUST** indicar que a implementación HTTP/UI de accounting en P1 (invoices recibidas, billing emitido, suppliers) segue os specs de catálogo principal en `openspec/specs/invoices/`, `openspec/specs/billing/`, `openspec/specs/suppliers/`, sen borrar a decisión de aplazamento do MVP v1.

#### Scenario: Trazabilidade lector

- GIVEN un lector do spec `p1-accounting-defer` no catálogo principal
- WHEN busca «que pasa despois do aplazamento MVP v1»
- THEN atopa referencia explícita ao change archivado e aos tres specs P1 no cartafol `openspec/specs/`

### Requirement: R-1 Enmienda de alcance documentada

**The system documentation MUST** seguir reflectindo que **CRUD invoices** e **CRUD billing** quedaron **fora da primeira entrega operativa do MVP v1** (2026-04-17). **The documentation SHALL** ademais enlazar o change archivado `2026-04-20-p1-invoices-billing-suppliers` como **programa P1** onde se implementou esa funcionalidade, co glosario: *invoices* = documentos **recibidos** polo taller; *billing* = documentos **emitidos** polo taller.

#### Scenario: Histórico MVP v1 intacto

- GIVEN documentación de peche MVP v1
- WHEN revísase Entra 1.1
- THEN o aplazamento original **MUST** seguir lexible

#### Scenario: Puente a P1

- GIVEN o change P1 archivado e os specs principais actualizados
- WHEN un PM abre os docs enlazados
- THEN **SHALL** ver alcance P1 acotado a invoices/billing/suppliers segundo eses specs

### Requirement: R-2 Tareas condicionais 3.2 / 3.3

**The OpenSpec tasks** do change archivado do roadmap **MUST** conservarse como **N/A** para a entrega MVP v1 histórica. **The tasks** do change `2026-04-20-p1-invoices-billing-suppliers` (archivado) **SHALL** reflejar la implementación HTTP/UI de P1 sen contradicer o peche 3.2–3.3 do arquivo antigo.

#### Scenario: Sen dobre conteo MVP

- GIVEN o arquivo de tarefas do roadmap archivado
- WHEN se audita cumprimento MVP v1
- THEN 3.2–3.3 **SHALL** seguir marcadas N/A para MVP v1

#### Scenario: Tareas P1 completadas

- GIVEN o change P1 archivado
- WHEN se revisan tasks.md no arquivo
- THEN **MUST** constar o traballo pechado para API/UI de invoices, billing e suppliers
