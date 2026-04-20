# Tasks: P1 — received-invoices, billing-documents, suppliers (+ Invoice cliente)

## Phase 1 — Dominio, ports y esquema

- [x] 1.1 Cerrar enum `BillingDocumentKind` en `backend/internal/domain/billing_document.go` (`client_invoice`, `payroll`, `irs`, `other`) y campo `Kind` + metadatos mínimos acordados con spec.
- [x] 1.2 Crear `backend/internal/domain/received_invoice.go` (UUID, importe, fecha, categoría, `SupplierID` opcional, notas, timestamps, soft-delete si aplica).
- [x] 1.3 Crear `backend/internal/domain/supplier.go` (nombre, contacto, tax id opcional, notas, activo).
- [x] 1.4 Añadir interfaces en `backend/internal/core/ports/repositories.go` y `services.go`: `ReceivedInvoiceRepository`, `BillingDocumentRepository`, `SupplierRepository`; extender `InvoiceRepository` con `Create`, `List`/`ListForStaff` (firma en diseño) y total count si lista paginada.
- [x] 1.5 Añadir migración SQL idempotente en `backend/migrations/` para `suppliers`, `received_invoices`, `billing_documents` (FKs y índices; `received_invoices.supplier_id` nullable).

## Phase 2 — Repositorios y servicios (backend)

- [x] 2.1 **RED/GREEN**: `supplier_repository_test.go` + `supplier_repository.go` (`NewPostgresSupplierRepository`, CRUD + `List` + count); tests in-memory con `github.com/glebarez/sqlite` (sin CGO); rama `sqlx` en Postgres como `car_repository`.
- [x] 2.2 **RED/GREEN**: `received_invoice_repository_test.go` + `received_invoice_repository.go`: mismas operaciones; listados excluyen `deleted_at` no nulo; `supplier_id` opcional coherente con migración 010.
- [x] 2.3 **RED/GREEN**: `billing_document_repository_test.go` + `billing_document_repository.go`: CRUD + `List`; columnas `kind`, `title`, `customer_id` nullable alineadas a `domain.BillingDocument`.
- [x] 2.4 **RED/GREEN**: `invoice_repository_test.go` + `invoice_repository.go`: implementar `ports.InvoiceRepository` completo (`GetByID`, `Update`, `ListByCustomerID`, `Create`, `ListForStaff`, `Delete`); tabla `invoices` vía `AutoMigrate` en test.
- [x] 2.5 **RED/GREEN**: `backend/internal/service/supplier/supplier_service.go` + `supplier_service_test.go`: solo `IsEmployee()` para CRUD/list; cliente → `ErrUnauthorizedAccess` (spec `suppliers` + matriz `mvp-role-access`).
- [x] 2.6 Extender `backend/internal/service/invoice/invoice_service.go` y tests: staff `Create`/`List`/`Delete` o soft-delete; cliente sin nuevas rutas rotas (`ListMyInvoices` intacto).
- [x] 2.7 **RED/GREEN**: `backend/internal/service/billing_document/billing_document_service.go` + tests: staff CRUD/list; `doc.Validate()` antes de persistir; cliente denegado.
- [x] 2.8 **RED/GREEN**: `backend/internal/service/received_invoice/received_invoice_service.go` + tests: staff CRUD/list; `inv.Validate()` antes de crear/actualizar; escenarios “Cliente sin acceso” del spec `invoices`.

## Phase 3 — HTTP (Gin) y wiring

- [x] 3.1 Handlers `backend/internal/handler/received_invoice_handler.go`, `billing_document_handler.go`, `supplier_handler.go` + Swagger anotaciones; respuestas JSON camelCase como el resto del API.
- [x] 3.2 Handler `invoice_handler.go` (o extensión) para rutas staff de `Invoice` + rutas cliente existentes (`GET/PATCH` mis facturas).
- [x] 3.3 Registrar modelos en slice `AutoMigrate` y montar grupos en `backend/cmd/api/main.go` bajo `protected`: `/received-invoices`, `/billing-documents`, `/suppliers`, `/invoices` (matriz rol por spec `mvp-role-access`); inyectar repos/servicios Postgres nuevos.

## Phase 4 — Frontend (Next)

- [x] 4.1 Servicios/API client: `frontend/src/lib/services/` (o `api-client`) para los cuatro recursos; tipos en `frontend/src/types/` si aplica.
- [x] 4.2 Rutas staff: páginas listado/detalle/edición bajo layout taller (misma shell que dashboard/cars), navegación mínima en `AppShell` o menú existente.
- [x] 4.3 Vista cliente: lista / detalle de facturas propias (`Invoice`) solo lectura + notas si el spec lo mantiene.

## Phase 5 — Pruebas y documentación

- [x] 5.1 Ajustar o ampliar tests unitarios post-handlers si faltan casos RBAC (`invoices`, `billing`, `suppliers`); `invoice_service` ya cubierto en Fase 2 si no regresa.
- [x] 5.2 Integración `httptest`: nuevo archivo o extensión en `backend/internal/handler/` o `backend/tests/integration/` — staff 200 en CRUD recibidas; cliente 403 en recibidas; cliente 200 en sus `invoices`.
- [x] 5.3 `go vet ./...` y `go test ./...` (CI); `frontend`: `pnpm lint`, `pnpm typecheck`, `pnpm test`.
- [x] 5.4 Actualizar `docs/mvp-solo-checklist.md` (o doc P1 enlazado) con enlace a este change y decisión `BillingDocumentKind`; una línea en `README` si aplica rutas nuevas.
