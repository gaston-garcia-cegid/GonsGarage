# Suppliers — especificación P1

> **Promoción:** incorporado ao catálogo principal desde `openspec/changes/archive/2026-04-20-p1-invoices-billing-suppliers/` (2026-04-20).

## Purpose

**Suppliers** é o **xestor de provedores** do taller: máster de terceiros que abastecen bens ou servizos, usable para vincular ou contextualizar **invoices recibidas** (P1 **MAY** relacionar factura recibida con provedor; non é obrigatorio en v1).

## Requirements

### Requirement: CRUD provedores

**The system SHALL** permitir alta, listado, detalle, actualización e baixa (ou soft-delete) de provedores con campos mínimos: identificación comercial, contacto, identificador fiscal opcional, notas opcionais.

#### Scenario: Staff xestiona provedores

- GIVEN un usuario staff
- WHEN crea ou edita un provedor con datos válidos
- THEN o rexistro persístese e aparece en búsquedas/listados staff

#### Scenario: Cliente sen acceso ao máster

- GIVEN un usuario `client`
- WHEN intenta acceder ao CRUD de suppliers
- THEN **MUST** denegarse

### Requirement: Integración opcional con invoices

**The system MAY** asociar unha invoice recibida a un `supplier_id`; se non hai provedor, a factura **SHALL** poder existir igualmente.

#### Scenario: Factura sen provedor

- GIVEN staff crea invoice recibida sen elixir provedor
- WHEN garda
- THEN **MUST** aceptarse salvo regras de validación de negocio que o deseño engada explicitamente
