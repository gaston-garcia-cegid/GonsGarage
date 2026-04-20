# Apply progress — p1-invoices-billing-suppliers

**Mode**: Strict TDD (per `openspec/config.yaml`).

## TDD cycle evidence

| Task | Test / evidence | Layer | RED | GREEN | Refactor |
|------|-----------------|-------|-----|---------|----------|
| 1.1 | `billing_document_test.go` — kind + `Validate` | Unit | Tests before types | `billing_document.go` | — |
| 1.2 | `received_invoice_test.go` — `Validate` | Unit | — | `received_invoice.go` | — |
| 1.3 | `supplier_test.go` — `Validate` | Unit | — | `supplier.go` | — |
| 1.4 | Ports + `stubInvoiceRepo` extended | Unit | Compile until stubs | `repositories.go`, `services.go` | — |
| 1.5 | `010_p1_...sql` idempotent DDL | SQL | N/A | Migration file | — |
| 2.1 | `supplier_repository_test.go` + repo | Unit | Suite | `supplier_repository.go` | `clampRepoList` |
| 2.2 | `received_invoice_repository_test.go` + repo | Unit | Suite | `received_invoice_repository.go` (GORM) | — |
| 2.3 | `billing_document_repository_test.go` + repo | Unit | Suite | `billing_document_repository.go` | — |
| 2.4 | `invoice_repository_test.go` + repo | Unit | Suite | `invoice_repository.go` | — |
| 2.5 | `supplier_service_test.go` | Unit | RBAC + validate | `supplier_service.go` | — |
| 2.6 | `invoice_service_test.go` staff invoice | Unit | RBAC | `invoice_service.go` | Pagination clamp |
| 2.7 | `billing_document_service_test.go` | Unit | RBAC + validate | `billing_document_service.go` | — |
| 2.8 | `received_invoice_service_test.go` | Unit | RBAC + validate | `received_invoice_service.go` | — |
| 3.1 | `supplier_handler_test.go` List + `role_middleware_workshop_test.go` | HTTP/unit | Handler + staff middleware | Three resource handlers + `gin_context.go` helpers | — |
| 3.2 | `invoice_handler.go` + `ListMyInvoices` route `/invoices/me` | HTTP | N/A (new surface) | Staff POST/GET/DELETE under `RequireWorkshopStaff`; client GET+PATCH `/:id` | — |
| 3.3 | `main.go` compile + `AutoMigrate` slice + `dropAllTables` | Wiring | N/A | New repos/services wired; groups `/suppliers`, `/received-invoices`, `/billing-documents`, `/invoices` | — |
| 4.1 | `accounting.services.contract.test.ts` (Vitest) + servicios `supplier` / `received-invoice` / `billing-document` / `issued-invoice` + `types/accounting.ts` | Frontend unit | Contract tests (URLs/verbs) | Service modules + exports | — |
| 4.2 | Manual QA via rutas; layout `accounting` (gate staff) | UI | N/A | `/accounting/*` hub + list/detail/new por recurso; `AppShell` nav «Contabilidade» | — |
| 4.3 | Layout `my-invoices` (gate cliente) + `PATCH` notas | UI | N/A | `/my-invoices` lista + detalle notas; nav «As minhas faturas» | — |
| 5.1 | `p1_accounting_routes_test.go` — cliente 403 suppliers/billing + invoice staff POST | HTTP | ✅ Written | ✅ Passed | — |
| 5.2 | Mismo ficheiro — employee/manager CRUD recibidas 200/201 + cliente `GET /invoices/me` e `GET /invoices/:id` | HTTP | ✅ Written | ✅ Passed | — |
| 5.3 | `go vet` + `go test ./...`; `pnpm typecheck` + `pnpm test` (+ lint vía build) | CI | N/A | Executed en batch apply/verify | — |
| 5.4 | `docs/mvp-solo-checklist.md` bloque P1; `README.md` táboa doc | Doc | N/A | ✅ Done | — |

**Command**: `cd backend && go test ./... -count=1`

## Deviations

- **Middleware nuevo** `RequireWorkshopStaff` (admin/manager/employee): matriz spec “staff” incluye **employee**; `RequireStaffManagers` existente sigue solo para `/employees`.
- Repos **2.2–2.4**: solo GORM; **2.1** GORM + sqlx. Tests repo: glebarez/sqlite.
- Dominio: sin `default:gen_random_uuid()` en tags GORM; migración SQL 010 conserva defaults en Postgres.

## Issues

- `swag init` no ejecutado en este batch; anotaciones `@Summary`/`@Router` listas para regenerar `docs/` en CI o local.

## Remaining

- Ninguna tarea pendente neste change; opcional: `pnpm lint` standalone se se quere cero warnings no repo; `go test -cover` / cobertura frontend cando se instalen ferramentas.
