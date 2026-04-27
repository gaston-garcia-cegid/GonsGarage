# Tasks: Stock de repuestos (manager / admin)

> `strict_tdd: true` en repo: en fases 2.x escribir **test antes** de código de producto (RED→GREEN→refactor breve en misma sesión).

## Phase 1: Foundation

- [x] 1.1 Create `backend/internal/domain/part_item.go` (GORM, soft delete, `Validate()`, `TableName()` `part_items`).
- [x] 1.2 Add `PartRepository` in `backend/internal/core/ports/repositories.go` (CRUD, `GetByBarcode`, `List` con filtros opcionales).
- [x] 1.3 Add `PartService` in `backend/internal/core/ports/services.go` (Create/Get/List/Update/Delete firmas alineadas a diseño).
- [x] 1.4 Create `backend/internal/repository/postgres/part_item_repository.go` implementando `PartRepository` (GORM; sin reglas negocio).
- [x] 1.5 Wire `PartItem` en `AutoMigrate` slice y `dropAllTables` en `backend/cmd/api/main.go`.

## Phase 2: Domain logic (TDD)

- [x] 2.1 (TDD) RED: `backend/internal/service/part/part_service_test.go` — cantidad negativa y UoM inválido fallan; GREEN: `part_service.go` validación + errores dominio.
- [x] 2.2 (TDD) RED: test barcode duplicado (mismo string no vacío); GREEN: comprobar repo antes de create/update en `part_service.go`.
- [x] 2.3 (TDD) RED: test `minimumQuantity` opcional coherente (≥0); GREEN: validar en servicio.

## Phase 3: HTTP + wiring

- [x] 3.1 Create `backend/internal/handler/part_handler.go` (DTOs camelCase, mapa errores→HTTP como `supplier_handler.go`).
- [x] 3.2 Register `protected.Group("/parts").Use(RequireStaffManagers())` con GET/POST/GET/:id/PATCH/DELETE en `backend/cmd/api/main.go` + inyección repo/svc/handler.
- [x] 3.3 (TDD) Extend `backend/internal/handler/mvp_role_access_test.go`: `client`/`employee` mutadores/list `parts` **MUST NOT** 2xx; `manager`/`admin` **MUST** pasar gate (403 solo por rol no aplica).

## Phase 4: Frontend

- [x] 4.1 Add `parts` API methods en `frontend/src/lib/api-client.ts` (`/parts`, query `barcode`/`search`).
- [x] 4.2 Extend `AppShellNavId` y nav condicional `canManageUsers` → `/admin/parts` en `frontend/src/components/layout/AppShell.tsx`.
- [x] 4.3 Create `frontend/src/app/admin/parts/layout.tsx` (guard `canManageUsers`, patrón `admin/users/layout.tsx`).
- [x] 4.4 Create `frontend/src/app/admin/parts/page.tsx` (lista, búsqueda por barcode/Enter, enlace crear/editar).
- [x] 4.5 Form create/edit (mismo change o sub-ruta) con campos diseño: `reference`, `brand`, `name`, `barcode`, `quantity`, `uom`, `minimumQuantity`.

## Phase 5: UI tests + matrix

- [x] 5.1 Update `frontend/src/components/layout/AppShell.test.tsx`: nav inventario para `manager`; ausente `client`/`employee` (spec `parts-inventory` + delta CI).
- [x] 5.2 Vitest mínimo: render `/admin/parts` con usuario manager mock y assert lista o empty state sin error.

## Phase 6: Spec / matrix follow-up (post-MVP código)

- [x] 6.1 Actualizar tabla resumen en `openspec/specs/mvp-role-access/spec.md` con fila Parts inventory (archivo principal; fuera del delta hasta archive).
