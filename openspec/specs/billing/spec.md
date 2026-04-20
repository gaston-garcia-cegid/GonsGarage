# Billing (emitido) — especificación P1

> **Promoción:** incorporado ao catálogo principal desde `openspec/changes/archive/2026-04-20-p1-invoices-billing-suppliers/` (2026-04-20).

## Purpose

**Billing** en P1 designa **o que o taller emite**: facturas a clientes, **nóminas de soldos**, documentación **IRS** ou outras retencións/obrigas similares, e demais saídas contables-operativas acotadas nesta fase. É distinto de **invoices** (só documentos **recibidos**).

## Requirements

### Requirement: Alcance e non solapamento

**The system SHALL** modelar “billing” só como rexistros **emitidos** polo taller. **The system MUST NOT** mesturar neste dominio as facturas de compra recibidas de provedores (iso é **invoices**).

#### Scenario: Rexistro emitido staff-only

- GIVEN un usuario staff
- WHEN crea un rexistro de billing (p. ex. factura a cliente ou tipo “nómina” / “IRS” segundo catálogo P1)
- THEN queda almacenado con tipo/categoría distinguible e trazable

#### Scenario: Cliente ve o seu

- GIVEN un cliente autenticado
- WHEN lista ou consulta billing (no sentido de facturas emitidas a el)
- THEN **SHALL** ver só documentos emitidos **a el** ou vinculados á súa conta (sen ver nóminas de terceiros nin IRS global do taller)

### Requirement: CRUD mínimo

**The system SHALL** expor CRUD para entidades de billing acotadas en deseño, con matriz rol × operación coherente con `mvp-role-access`.

#### Scenario: Operación prohibida

- GIVEN un rol sen permiso para o subtipo (p. ex. cliente editando nómina)
- WHEN intenta mutación
- THEN **MUST** responder denegación acorde ao stack API existente
