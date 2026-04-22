# workshop-repair-execution Specification

> **Promoción** catálogo principal dende o change archivado `openspec/changes/archive/2026-04-22-workshop-mechanic-vehicle-lifecycle/` (2026-04-22).

## Purpose

Ciclo de servicio de taller por vehículo: la **visita (service job) es la unidad de trazabilidad** — recepción, cierre y, en fases posteriores, diagnóstico, presupuesto y entrega formal. **MVP1 mínima:** recepción + cierre; OBD/ presupuesto como **stubs** con contrato explícito.

| Fase | MVP1 | Post-MVP1 |
|------|------|-----------|
| A Recepción (km, fluidos, neumáticos, notas) | **SHALL** persistir | **MAY** fotos / tercero |
| B OBD / diagnóstico | **MAY** (stub) | códigos + historial |
| C Presupuesto / aprobación | **MAY** estado stub | flujo aprobado |
| D Ejecución (enlace a trabajo) | **SHOULD** enlazar a `repairs` | SLA y desvíos |
| E Entrega re-verificada | **SHALL** persistir cierre básico | firma / comprobante |

## Requirements

### Requirement: Agregado de visita (service job)

El sistema **SHALL** modelar una **visita de taller** (service job) como agregado propio con `car_id` **SHALL** obligatorio; `appointment_id` **MAY** nulo. El ciclo de visita **MUST NOT** depender solo de filas de `repairs` sin este agregado. Las reparaciones históricas sin visita **SHALL** seguir siendo legítimas (sin backfill forzado).

#### Scenario: Alta de visita

- GIVEN un `employee` con permiso de acceso al coche
- WHEN crea una visita para ese `car_id`
- THEN existe identificador de visita y estado mínimo **open** (nome concreto en implementación)

#### Scenario: Listado con legado

- GIVEN reparaciones antiguas sin service job
- WHEN se listan reparaciones del coche
- THEN se muestran como hoy; **MUST NOT** exigir migración automática a visita

### Requirement: Recepción y cierre mínimos

La visita **SHALL** aceptar **recepción** (p. ej. km, fluidos, neumáticos, notas) y **cierre/entrega** re-verificados. El persistido **SHALL** ser **estructurado y versionable** (tablas hijas o bloques con `schema_version`); **MUST NOT** aceptar JSON suelto sin esquema versionado en contrato de API.

#### Scenario: Recepción completa y válida

- GIVEN visita abierta
- WHEN el empleado envía un checklist de recepción **válido** según validación
- THEN la recepción queda persistida con al menos qué usuario y cuándo

#### Scenario: Recepción incompleta

- GIVEN la misma visita
- WHEN faltan campos obligatorios
- THEN la API **MUST** rechazar (no 2xx) con error interpretable

#### Scenario: Cierre

- GIVEN recepción persistida
- WHEN el empleado completa el checklist de cierre/entrega válido
- THEN el cierre queda persistido y el estado de visita **MUST** reflejar cierre (p. ej. **closed**)

### Requirement: OBD y presupuesto (stub MVP1+)

OBD, estimación o aprobación **MAY** exponerse como *stub* (p. ej. 501, o recurso vacío con semántica documentada) sin romper visitas v1.

#### Scenario: Lógica aún no implementada

- GIVEN MVP1 sin lógica de OBD/estimate
- WHEN se llama al sub-recurso correspondiente
- THEN la respuesta **MUST** ser predecible; **MUST NOT** 200 con cuerpo que implique un estado negocio falso

### Requirement: Trazas de transición mínima

Apertura, cierre o cancelación de la visita **SHALL** dejar rastro (quién/cuándo) al menos en API de lectura de la visita.

#### Scenario: Lectura tras cierre

- GIVEN cierre de visita completado
- WHEN se obtiene el recurso `GET` de la visita
- THEN la carga **SHALL** incluir marcas mínimas de apertura y cierre (actor/timestamp según criterio de cierre aprobado en implementación)

## Nota: roles

Resumen: las obligaciones **MUST/SHALL** de mutación e a matriz pública vive en `openspec/specs/mvp-role-access/spec.md`. O cliente **MUST NOT** en mutación de visita; o persoal *staff* (employee, manager, admin) na matriz; rotas bajo `RequireWorkshopStaff` / regras de servizo.
