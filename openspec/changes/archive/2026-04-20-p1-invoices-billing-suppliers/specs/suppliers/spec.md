# Suppliers — especificación P1

## Purpose

**Suppliers** es el **gestor de proveedores** del taller: maestro de terceros que abastecen bienes o servicios, usable para vincular o contextualizar **invoices recibidas** (P1 **MAY** relacionar factura recibida con proveedor; no es obligatorio en v1).

## Requirements

### Requirement: CRUD proveedores

**The system SHALL** permitir alta, listado, detalle, actualización y baja (o soft-delete) de proveedores con campos mínimos: identificación comercial, contacto, identificador fiscal opcional, notas opcionales.

#### Scenario: Staff gestiona proveedores

- GIVEN un usuario staff
- WHEN crea o edita un proveedor con datos válidos
- THEN el registro se persiste y aparece en búsquedas/listados staff

#### Scenario: Cliente sin acceso al maestro

- GIVEN un usuario `client`
- WHEN intenta acceder al CRUD de suppliers
- THEN **MUST** denegarse

### Requirement: Integración opcional con invoices

**The system MAY** asociar una invoice recibida a un `supplier_id`; si no hay proveedor, la factura **SHALL** poder existir igualmente.

#### Scenario: Factura sin proveedor

- GIVEN staff crea invoice recibida sin elegir proveedor
- WHEN guarda
- THEN **MUST** aceptarse salvo reglas de validación de negocio que el diseño añada explícitamente
