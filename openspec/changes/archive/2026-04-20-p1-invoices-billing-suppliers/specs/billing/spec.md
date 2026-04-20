# Billing (emitido) — especificación P1

## Purpose

**Billing** en P1 designa **lo que el taller emite**: facturas a clientes, **nóminas de sueldos**, documentación **IRS** u otras retenciones/obligaciones similares, y demás salidas contables-operativas acotadas en esta fase. Es distinto de **invoices** (solo documentos **recibidos**).

## Requirements

### Requirement: Alcance y no solapamiento

**The system SHALL** model “billing” solo como registros **emitidos** por el taller. **The system MUST NOT** mezclar en este dominio las facturas de compra recibidas de proveedores (eso es **invoices**).

#### Scenario: Registro emitido staff-only

- GIVEN un usuario staff
- WHEN crea un registro de billing (p. ej. factura a cliente o tipo “nómina” / “IRS” según catálogo P1)
- THEN queda almacenado con tipo/categoría distinguible y trazable

#### Scenario: Cliente ve lo suyo

- GIVEN un cliente autenticado
- WHEN lista o consulta billing
- THEN **SHALL** ver solo documentos emitidos **a él** o vinculados a su cuenta (sin ver nóminas de terceros ni IRS global del taller)

### Requirement: CRUD mínimo

**The system SHALL** exponer CRUD para entidades de billing acotadas en diseño, con matriz rol × operación coherente con `mvp-role-access`.

#### Scenario: Operación prohibida

- GIVEN un rol sin permiso para el subtipo (p. ej. cliente editando nómina)
- WHEN intenta mutación
- THEN **MUST** responder denegación acorde al stack API existente
