# Delta for p1-accounting-defer

## Context (unchanged narrative)

El roadmap previó P1 como rutas + UI para accounting o enmienda al checklist; el MVP v1 cerró con aplazamiento documentado. Este delta **mantiene** ese histórico y **añade** el programa P1 activo bajo el change `p1-invoices-billing-suppliers`, con glosario: **invoices** = facturas **recibidas** por el taller; **billing** = documentos **emitidos** por el taller (clientes, nóminas, IRS, etc.).

## ADDED Requirements

### Requirement: Programme P1 tracking

**The system documentation MUST** indicar que la implementación HTTP/UI de accounting en P1 (invoices recibidas, billing emitido, suppliers) se desarrolla según los specs del change `p1-invoices-billing-suppliers` (`invoices`, `billing`, `suppliers`), sin borrar la decisión de aplazamiento del MVP v1.

#### Scenario: Trazabilidad lector

- GIVEN un lector del spec `p1-accounting-defer` en catálogo principal
- WHEN busca “qué pasa después del aplazamiento MVP v1”
- THEN encuentra referencia explícita al change y a los tres specs P1

## MODIFIED Requirements

### R-1 Enmienda de alcance documentada

**The system documentation MUST** seguir reflejando que **CRUD invoices** y **CRUD billing** quedaron **fuera de la primera entrega operativa del MVP v1** (2026-04-17). **The documentation SHALL** además enlazar el change `p1-invoices-billing-suppliers` como **programa P1** donde se implementa esa funcionalidad con el glosario: *invoices* = documentos **recibidos** por el taller (proveedores, servicios, compras); *billing* = documentos **emitidos** por el taller (clientes, nóminas, IRS, etc.).

(Previously: solo backlog hasta change dedicado sin nombre de change ni glosario invoice/billing.)

#### Scenario: Histórico MVP v1 intacto

- GIVEN documentación de cierre MVP v1
- WHEN se revisa Entra 1.1
- THEN el aplazamiento original **MUST** seguir siendo legible

#### Scenario: Puente a P1

- GIVEN el change `p1-invoices-billing-suppliers` aprobado para ejecución
- WHEN un PM abre docs enlazados desde este delta
- THEN **SHALL** ver alcance P1 acotado a invoices/billing/suppliers según esos specs

### R-2 Tareas condicionales 3.2 / 3.3

**The OpenSpec tasks** del change archivado **MUST** conservarse como **N/A** para la entrega MVP v1 histórica. **The new tasks** del change `p1-invoices-billing-suppliers` **SHALL** cubrir implementación HTTP/UI de P1 sin contradecir el cierre 3.2–3.3 del archivo antiguo.

(Previously: tasks MUST ser N/A mientras solo aplicaba R-1 sin programa P1 nombrado.)

#### Scenario: Sin doble conteo MVP

- GIVEN el archivo de tareas del roadmap archivado
- WHEN se audita cumplimiento MVP v1
- THEN 3.2–3.3 **SHALL** seguir marcadas N/A para MVP v1

#### Scenario: Tareas P1 nuevas

- GIVEN el change `p1-invoices-billing-suppliers` entra en apply
- WHEN se generan tasks
- THEN **MUST** existir trabajo explícito para API/UI de invoices, billing y suppliers

## REMOVED Requirements

_None._
