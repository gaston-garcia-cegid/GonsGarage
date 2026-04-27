# Proposal: Stock de repuestos (manager / admin)

## Intent

Inventario de repuestos con CRUD, cantidad + unidad (uds, L, etc.) y alta por **código de barras** (cualquier lector HID). Solo `manager` y `admin` en API y UI.

## Scope

### In Scope

- API REST CRUD repuesto/stock; UoM explícita; middleware `RequireStaffManagers`.
- UI sección dedicada (App Router) + flujo escanear → buscar o pre-rellenar (referencia, marca, etc.).
- Tests CI matriz/regresión; spec con **sugerencias típicas** de inventario acotadas al MVP (barcode único, mínimos, auditoría mínima).

### Out of Scope

- Consumo en órdenes de taller; multi-almacén; lotes; costes; app nativa; import CSV (MAY futuro).

## Capabilities

### New Capabilities

- `parts-inventory`: CRUD, stock numérico + UoM, barcode en web, reglas datos mínimos, UX sugerencias en spec.

### Modified Capabilities

- `mvp-role-access`: Fila **Stock repuestos** (HTTP+UI): manager/admin **MUST**; employee/client **MUST NOT**.

## Approach

Backend: tablas GORM + handlers bajo grupo con `RequireStaffManagers`. Frontend: patrón shell como `/admin/users`. Barcode: `input` + string + Enter, sin SDK. Spec: SHOULD mínimos/stock bajo; MAY movimientos/import (documentado defer).

## Affected Areas

| Area | Impact |
|------|--------|
| `backend/internal/*` | New domain/repo/service/handlers |
| `frontend/src/app/*` | New ruta + nav |
| `openspec/specs/mvp-role-access/spec.md` | Modified matrix |
| `openspec/specs/parts-inventory/spec.md` | New (post-archive) |

## Risks

| Risk | L | Mitigation |
|------|---|------------|
| Barcode duplicado | M | Validación + mensaje; fusión post-MVP |
| Prefijos pistola | M | Strip configurable o doc UI |
| Scope creep | M | Out of scope arriba |

## Rollback Plan

Revert migraciones + rutas + menú; redeploy anterior. Sin backup, se pierden datos nuevos.

## Dependencies

Ninguna bloqueante.

## Success Criteria

- [ ] API CRUD con 403 employee/client.
- [ ] UI solo manager/admin.
- [ ] Cantidad + UoM persistidas.
- [ ] Spec: escaneo HID + sugerencias comunes (único barcode, mínimos, auditoría ligera).
- [ ] Matriz `mvp-role-access` y tests CI actualizados.
