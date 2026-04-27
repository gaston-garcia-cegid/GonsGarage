# Design: Stock de repuestos (manager / admin)

## Technical Approach

Añadir **dominio `PartItem`** (catálogo + stock en la misma fila MVP), **API REST** bajo `/api/v1/parts` con **`RequireStaffManagers()`** (misma superficie que `/employees` y `POST /admin/users`), y **UI** en `/admin/parts` con guard **`canManageUsers`** (equivale a manager+admin). **GORM `AutoMigrate`** en `cmd/api/main.go` siguiendo `Supplier`. **JSON camelCase** en handlers como `supplier_handler.go`. **Barcode**: campo de texto + `Enter` / búsqueda explícita; sin SDK.

## Architecture Decisions

| Decision | Alternativa | Elección | Rationale |
|----------|--------------|----------|-----------|
| Ruta API | `/inventory/items` | **`/parts`** | Corta, coherente con plural en repo (`/suppliers`, `/cars`). |
| Tabla GORM | SQL-only migration | **`AutoMigrate` + modelo** | Igual que `Supplier`, `ReceivedInvoice`; menos fricción en dev. |
| UoM | Tabla `units` | **`varchar` cerrado** (`unit`, `liter`, …) | MVP; ampliar enum en código. |
| Borrado | Hard delete | **Soft delete `deleted_at`** | Alineado a `Supplier`; rollback más simple. |
| Unicidad barcode | Nullable único global | **Unique index solo si `barcode` no vacío** | Postgres partial unique en fase apply o validación app + índice compuesto documentado. |
| Nav / UI guard | Nuevo helper rol | **`canManageUsers`** | Spec exige manager+admin; ya existe en frontend. |

## Data Flow

```
AppShell (canManageUsers) → /admin/parts
    → apiClient GET/POST/PATCH/DELETE /api/v1/parts[...]
        → Gin JWT → RequireStaffManagers → PartHandler → PartService → PartRepository (GORM)
```

Búsqueda por barcode: `GET /api/v1/parts?barcode=` (o query `search=`) → lista 0–1 ítems; formulario alta reutiliza respuesta vacía para pre-rellenar.

## File Changes

| Ruta | Acción | Descripción |
|------|--------|-------------|
| `backend/internal/domain/part_item.go` | Create | Struct GORM + `Validate()` + `TableName()`. |
| `backend/internal/core/ports/services.go` (+ repo si patrón separado) | Modify | `PartService` / repositorio según patrón existente (supplier usa repo dedicado). |
| `backend/internal/repository/postgres/part_item_repository.go` | Create | CRUD + `GetByBarcode` / list con filtros. |
| `backend/internal/service/part/` o `part_item/` | Create | Reglas negocio: qty ≥0, barcode duplicado, mínimos. |
| `backend/internal/handler/part_handler.go` | Create | DTOs request/response camelCase, mapa errores → HTTP. |
| `backend/cmd/api/main.go` | Modify | AutoMigrate modelo, wiring repo/svc/handler, `setupRoutes` grupo `/parts` + middleware. |
| `backend/cmd/api/main.go` `dropAllTables` | Modify | Incluir tabla nueva si reset dev. |
| `backend/internal/handler/mvp_role_access_test.go` | Modify | Matriz: client/employee 403/401; manager/admin 2xx en CRUD mínimo. |
| `frontend/src/lib/api-client.ts` | Modify | Métodos `parts.*`. |
| `frontend/src/app/admin/parts/*` | Create | `layout.tsx` guard + `page.tsx` lista; form create/edit (modal o página según patrón admin users). |
| `frontend/src/components/layout/AppShell.tsx` | Modify | `AppShellNavId` + botón nav condicional `canManageUsers`. |
| `frontend/src/components/layout/AppShell.test.tsx` | Modify | Nav ausente client/employee; presente manager. |

## Interfaces / Contracts

**REST (bajo `/api/v1/parts`, JWT, `RequireStaffManagers`):**

| Método | Ruta | Uso |
|--------|------|-----|
| GET | `/parts` | Lista; query `barcode`, `search`. |
| POST | `/parts` | Crear. |
| GET | `/parts/:id` | Detalle. |
| PATCH | `/parts/:id` | Actualizar (cantidad, UoM, mínimo, barcode, …). |
| DELETE | `/parts/:id` | Soft delete. |

**Cuerpo ejemplo (POST/PATCH)** — campos alineados al spec:

```json
{
  "reference": "NGK-123",
  "brand": "NGK",
  "name": "Bujía",
  "barcode": "5901234123457",
  "quantity": 12,
  "uom": "unit",
  "minimumQuantity": 4
}
```

**Respuesta lista ítem**: `id`, campos anteriores, `createdAt`, `updatedAt`, `deletedAt` opcional.

## Testing Strategy

| Capa | Qué | Cómo |
|------|-----|------|
| Unit | `PartService` validación / duplicados | `go test` tabla-driven en `internal/service/part`. |
| HTTP | Matriz roles + CRUD mínimo | Extender `mvp_role_access_test.go` con rutas `/parts`. |
| UI | Nav + guard | `AppShell.test.tsx` + smoke render `/admin/parts` con mock user. |

Con **`strict_tdd: true`**, tests nuevos **antes** de implementación en cada tarea que toque lógica (skill sdd-apply).

## Migration / Rollout

**AutoMigrate** al arranque API; sin feature flag. Rollback: deploy previo + opcional `DROP TABLE` en entornos dev con `RESET_DATABASE` (actualizar lista `dropAllTables`).

## Open Questions

- [ ] ¿Índice único parcial `barcode` vía `createIndexes` en Go o validación solo en servicio para MVP?
- [ ] ¿Paginación obligatoria en `GET /parts` desde v1 o lista completa hasta N ítems?
