# Proposal: P1 — Invoices, Billing y proveedores

## Intent

MVP v1 está cerrado; `p1-accounting-defer` documentó accounting HTTP/UI fuera de v1. Este change **planifica la primera entrega P1**: CRUD de facturas (invoices), CRUD de facturación/cobros (billing) y un **gestor de proveedores (suppliers)** integrado al taller, con API, persistencia y UI coherentes con el stack actual (Gin, Next, roles).

## Scope

### In Scope

- Especificar y luego implementar **CRUD invoices** (list/create/read/update/delete o equivalente acotado por negocio) con autorización por rol alineada a `mvp-role-access`.
- Especificar e implementar **CRUD billing** (entidad y reglas de negocio a definir en spec: p. ej. pagos, planes, estados; sin confundir con “solo notas” del `InvoiceService` actual).
- **Gestor de suppliers**: modelo de datos, API REST bajo `/api/v1`, repositorio Postgres, UI mínima (listado + alta/edición/baja o soft-delete).
- Migraciones, Swagger, tests de servicio/handler donde aplique; enlaces en docs/checklist P1.

### Out of Scope

- Contabilidad general, impuestos multi-país, integración con gateways de pago reales.
- Importación masiva CSV de proveedores (fase posterior).
- Cambios no relacionados a auth shell o brand.

## Capabilities

### New Capabilities

- `invoices`: CRUD HTTP + persistencia + UI; matriz rol × operación; compatibilidad con dominio/servicio invoice existente donde encaje.
- `billing`: CRUD billing (definir entidad en spec: cobro vs plan vs estado); permisos staff vs cliente según negocio.
- `suppliers`: CRUD proveedores; campos mínimos (nombre, contacto, NIF opcional, notas); visibilidad staff.

### Modified Capabilities

- `p1-accounting-defer`: delta — dejar explícito que la **implementación P1** queda cubierta por este change (sin borrar el histórico de aplazamiento MVP v1).

## Approach

1. **sdd-spec**: tres deltas (`invoices`, `billing`, `suppliers`) + delta a `p1-accounting-defer`; escenarios Given/When/Then y RFC 2119.
2. **sdd-design**: ER ligero (relaciones invoice ↔ customer/car/repair si aplica; billing; supplier sin acoplar a factura en v1 si reduce riesgo).
3. **sdd-tasks** → **sdd-apply**: repos Postgres + handlers Gin + rutas protegidas; frontend rutas bajo layout taller; reutilizar patrones `car`/`appointment`.
4. **sdd-verify**: `go test`, `pnpm` lint/test; smoke manual por rol.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `backend/internal/domain/` | New/Modified | Entidades billing/supplier; extensión invoice si hace falta |
| `backend/internal/repository/postgres/` | New | Repos + migraciones SQL |
| `backend/cmd/api/main.go` | Modified | Rutas `/api/v1/invoices`, billing, suppliers |
| `frontend/src/` | New | Páginas/componentes gestión P1 |
| `openspec/specs/` | New + delta | Tres specs nuevos + delta accounting-defer |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Solapamiento invoice vs billing | Med | Spec fija glosario y límites antes de código |
| Alcance billing vago | Med | Una entidad “MVP billing” acotada en spec v1 |
| Migración en prod compartida | Baja | Scripts idempotentes; backup documentado en deploy |

## Rollback Plan

- Revertir PR/commits del change; `docker compose` con imagen anterior; migraciones **down** manuales o restore snapshot DB si se aplicaron DDL destructivos. Mantener feature flags desactivados hasta verificación si se introducen.

## Dependencies

- Postgres y JWT actuales; roles `mvp-role-access`; ningún servicio externo obligatorio en v1.

## Success Criteria

- [ ] Specs delta aprobadas y trazables a tasks.
- [ ] CRUD invoices y billing expuestos en API con Swagger y tests mínimos verdes.
- [ ] CRUD suppliers usable por rol staff desde UI.
- [ ] Checklist/docs P1 actualizados; delta `p1-accounting-defer` refleja transición de “diferido” a “en implementación vía change X”.
