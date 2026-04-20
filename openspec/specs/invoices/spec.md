# Invoices (recibidas) — especificación P1

> **Promoción:** incorporado ao catálogo principal desde `openspec/changes/archive/2026-04-20-p1-invoices-billing-suppliers/` (2026-04-20).

## Purpose

**Invoice** en P1 designa **toda factura ou documento equivalente que o taller recibe** (provedores, servizos contratados, compras diversas). É o arquivo operativo de **contas por pagar / compras**, distinto de **billing** (facturas emitidas polo taller).

## Requirements

### Requirement: Alcance e non solapamento

**The system SHALL** modelar “invoices” só como documentos **recibidos** polo taller. **The system MUST NOT** usar este dominio para facturas emitidas a clientes, nóminas ou obrigas fiscais emitidas (iso é **billing**).

#### Scenario: Alta de factura recibida

- GIVEN un usuario staff (admin, manager ou employee segundo matriz `mvp-role-access`)
- WHEN rexistra unha factura recibida con datos mínimos acordados en deseño (p. ex. provedor opcional, importe, data, categoría)
- THEN o rexistro queda persistido e visible en listados staff

#### Scenario: Cliente sen acceso

- GIVEN un usuario con rol `client`
- WHEN intenta crear ou listar “invoices” de taller
- THEN **MUST** denegarse (403 ou equivalente documentado)

### Requirement: CRUD mínimo

**The system SHALL** expor operacións CRUD (crear, ler lista/detalle, actualizar, eliminar ou baixa lóxica) para invoices recibidas, con regras de autorización aliñadas a `mvp-role-access`.

#### Scenario: Eliminación idempotente ou segura

- GIVEN unha factura recibida existente
- WHEN staff autorizado solicita baixa
- THEN o sistema **MUST** aplicar a política de borrado definida en deseño (físico ou soft-delete) sen deixar referencias rotas obrigatorias
