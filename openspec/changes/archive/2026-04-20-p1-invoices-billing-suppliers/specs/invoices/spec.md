# Invoices (recibidas) — especificación P1

## Purpose

**Invoice** en P1 designa **toda factura o documento equivalente que el taller recibe** (proveedores, servicios contratados, compras diversas). Es el archivo operativo de **cuentas por pagar / compras**, distinto de **billing** (facturas emitidas por el taller).

## Requirements

### Requirement: Alcance y no solapamiento

**The system SHALL** model “invoices” solo como documentos **recibidos** por el taller. **The system MUST NOT** usar este dominio para facturas emitidas a clientes, nóminas o obligaciones fiscales emitidas (eso es **billing**).

#### Scenario: Alta de factura recibida

- GIVEN un usuario staff (admin, manager o employee según matriz `mvp-role-access`)
- WHEN registra una factura recibida con datos mínimos acordados en diseño (p. ej. proveedor opcional, importe, fecha, categoría)
- THEN el registro queda persistido y visible en listados staff

#### Scenario: Cliente sin acceso

- GIVEN un usuario con rol `client`
- WHEN intenta crear o listar “invoices” de taller
- THEN **MUST** denegarse (403 o equivalente documentado)

### Requirement: CRUD mínimo

**The system SHALL** exponer operaciones CRUD (crear, leer lista/detalle, actualizar, eliminar o baja lógica) para invoices recibidas, con reglas de autorización alineadas a `mvp-role-access`.

#### Scenario: Eliminación idempotente o segura

- GIVEN una factura recibida existente
- WHEN staff autorizado solicita baja
- THEN el sistema **MUST** aplicar la política de borrado definida en diseño (físico o soft-delete) sin dejar referencias rotas obligatorias
