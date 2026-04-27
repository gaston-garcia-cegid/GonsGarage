# Apply progress: spare-parts-inventory

**Mode**: Strict TDD (`strict_tdd: true`)

## Completed batches

- [x] **Phase 1** (tasks 1.1–1.5): domain, ports, `part_item_repository`, `main.go` AutoMigrate/drop, repo tests sqlite.
- [x] **Phase 2** (tasks 2.1–2.3): `PartService` + `part_service_test.go` (stubs), `ErrPartItemDuplicateBarcode`, gate `CanManageUsers()`.
- [x] **Phase 3** (tasks 3.1–3.3): `part_handler.go`, rutas `/api/v1/parts` en `main.go`, tests MVP (stub + handler real sqlite).
- [x] **Phase 4** (tasks 4.1–4.5): `api-client` parts, `AppShell` nav `admin_parts`, layout guard, lista `/admin/parts`, `/admin/parts/new` + `/admin/parts/[id]`, tipos `types/parts.ts`.
- [x] **Phase 5** (tasks 5.1–5.2): `AppShell.test.tsx` matriz nav Peças (stock); `admin/parts/page.test.tsx` empty state + filas + filtros.
- [x] **Phase 6** (task 6.1): `openspec/specs/mvp-role-access/spec.md` — fila matriz Parts inventory + requisito + escenarios CI (Go + Vitest) + nota promoción.

## TDD — Phase 1 (recap)

| Task | RED | GREEN |
|------|-----|-------|
| 1.4 | `part_item_repository_test.go` antes de repo estable | GORM repo + tests verdes |
| 1.1–1.3 | Tipos/interfaces consumidos por tests | `part_item.go`, ports |

## TDD — Phase 2

| Task | RED | GREEN | Triangulation / REFACTOR |
|------|-----|-------|---------------------------|
| 2.1 | Tests neg qty + invalid UOM fallan sin servicio completo | `PartService.Create` + `Validate()` vía dominio | Admin OK + employee 403 |
| 2.2 | Tests duplicate create/update | `ensureNoDuplicateBarcode` + `ErrPartItemDuplicateBarcode` | Update mismo barcode en misma fila permitido |
| 2.3 | Test minimum negativo | Mismo `Validate()` dominio | — |

## TDD — Phase 3

| Task | RED | GREEN | REFACTOR |
|------|-----|-------|----------|
| 3.3 | `mvp_role_access_test.go`: stub router GET/POST client+employee 403; manager/admin 2xx; fallos si falta middleware o handler | `part_handler.go` + `setupRoutes` `/parts` + tests con `PartHandler`+sqlite+`mvpUserRepo` (list vacía, POST 201, client/employee POST 403) | Reutilizar `partsRouterWithPartHandler` / `partsStaffStubRouter` |
| 3.1–3.2 | Cubiertos por integración 3.3 (handler + wiring); sin tests de unidad dedicados en `tasks.md` | — | — |

**Commands**: `go test ./internal/handler/...`, `go vet ./...`, `go build ./cmd/api`.

## TDD — Phase 4

| Task | RED | GREEN | REFACTOR |
|------|-----|-------|----------|
| 4.1–4.5 | Tests UI reservados en **Phase 5** (`AppShell.test.tsx`, smoke `/admin/parts`); esta fase = implementación + `pnpm exec tsc --noEmit` + eslint en archivos tocados | Tipos `parts`, métodos `listParts`/`getPart`/`createPart`/`updatePart`/`deletePart`, rutas App Router reales | Reutilizar `AppShell`, CSS módulo local `admin-parts.module.css` |

**Commands**: `pnpm exec eslint …`, `pnpm exec tsc --noEmit` (frontend).

## TDD — Phase 5

| Task | RED | GREEN | REFACTOR |
|------|-----|-------|----------|
| 5.1 | Nuevos casos en `AppShell.test.tsx`: botón Peças (stock) manager/admin; ausente client/employee; `activeNav` `admin_parts`; click → `/admin/parts` | UI ya entregada en fase 4 — tests verdes al añadir describe | — |
| 5.2 | `page.test.tsx`: mock `listParts` + manager; assert vacío / tabla / `listParts` con filtros | `AdminPartsPage` existente | — |

**Commands**: `pnpm exec vitest run src/components/layout/AppShell.test.tsx src/app/admin/parts/page.test.tsx`; suite completa `pnpm test`.

## Phase 6 (spec follow-up)

| Task | RED | GREEN | REFACTOR |
|------|-----|-------|----------|
| 6.1 | Catálogo principal sin fila Parts | Actualizar `mvp-role-access/spec.md` (matriz, requisito dedicado, CI regression, escenario matriz publicada) | Referencia delta `openspec/changes/spare-parts-inventory/specs/parts-inventory/spec.md` hasta archive |

## Deviations

- Fase 1: `Validate()` también en repo; fase 2 repite vía dominio en servicio (defensa en capas).
- Autorización en servicio con `CanManageUsers()` (manager+admin), no `IsEmployee()` como suppliers (diseño parts-inventory).

## Remaining

- Opcional: archivar change `spare-parts-inventory` (promover delta `parts-inventory` a `openspec/specs/` si aplica al workflow SDD).
