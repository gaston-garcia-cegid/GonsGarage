# Design: P1 — Invoices recibidas, billing emitido y suppliers

## Technical Approach

Mapear los **tres specs** a **tres piezas de código** sin colisión de nombres: (1) **`domain.Invoice` + `InvoiceService`** se **mantienen** como factura **emitida al cliente** (encaja con “billing” y con el comentario actual en código); se amplían con repo Gin, CRUD staff y rutas HTTP. (2) Las facturas **recibidas** por el taller (proveedores, compras) son un **nuevo agregado** `ReceivedInvoice` (nombre de tipo; REST bajo **`/api/v1/received-invoices`** para no ocupar `/invoices` con semántica ambigua). (3) **Otros emitidos** (nóminas, IRS) viven en **`BillingDocument`** con `kind` discriminado y misma convención que `Appointment`/`Car` (handler + service + repo Postgres + sqlx opcional).

## Architecture Decisions

| Decisión | Opción elegida | Alternativas | Rationale |
|----------|------------------|--------------|-----------|
| Factura a cliente | Mantener `Invoice` / `InvoiceService` | Renombrar a `ClientInvoice` | Usuario acepta mantener orientación actual; menos churn y spec “billing” cubre emisión a cliente. |
| Facturas recibidas (spec `invoices`) | Nuevo `ReceivedInvoice` + tabla `received_invoices` | Reutilizar `Invoice` con flag | Evita mezclar `CustomerID` con proveedor y violar MUST NOT del spec. |
| Ruta REST recibidas | `/api/v1/received-invoices` | `/api/v1/invoices` | `/invoices` sugiere recurso único; el spec usa palabra “invoices” pero el contrato HTTP queda explícito en Swagger. |
| Billing nómina/IRS | Tabla `billing_documents` + `kind` enum | Tablas por tipo en v1 | Menos migraciones; filtrado por `kind` y RBAC en servicio. |
| Suppliers | `Supplier` + `suppliers` tabla + `/api/v1/suppliers` | Sin FK en v1 | FK opcional `received_invoices.supplier_id` nullable (spec MAY). |
| Persistencia | Patrón actual: GORM AutoMigrate + repo Postgres con sqlx cuando exista pool | Solo GORM | Alineado a `car_repository` / `appointment_repository`. |

## Data Flow

```
[Next.js staff] → Bearer JWT → Gin protected group
    ├→ POST/GET/PATCH/DELETE /received-invoices → ReceivedInvoiceService → postgres repo → received_invoices
    ├→ … /billing-documents → BillingDocumentService → billing_documents
    ├→ … /invoices (cliente) + staff CRUD extendido → InvoiceService → invoices (existente)
    └→ … /suppliers → SupplierService → suppliers
```

Cliente solo: `GET/PATCH` facturas propias vía rutas ya previstas por `InvoiceService` (lista “mis facturas”); staff usa endpoints de staff para emitidas y CRUD completo de recibidas/billing/suppliers.

## File Changes

| Ruta | Acción | Descripción |
|------|--------|---------------|
| `backend/internal/domain/received_invoice.go` (o nombre único) | Crear | Modelo + validación ligera |
| `backend/internal/domain/billing_document.go`, `supplier.go` | Crear | Tipos + `kind` billing |
| `backend/internal/repository/postgres/*_repository.go` | Crear | CRUD sqlx/GORM como cars |
| `backend/internal/handler/*_handler.go` | Crear | Gin + Swagger |
| `backend/internal/service/...` | Crear / modificar | Extender `InvoiceService`/ports con Create/List staff; nuevos servicios |
| `backend/cmd/api/main.go` | Modificar | AutoMigrate nuevos modelos; `protected` routes |
| `backend/migrations/*.sql` | Crear | DDL si no se confía solo a AutoMigrate en prod |
| `frontend/src/app/...` | Crear | Páginas taller bajo layout existente (dashboard/shell) |

## Interfaces / Contracts

- **Ports**: `ReceivedInvoiceRepository`, `BillingDocumentRepository`, `SupplierRepository`; ampliar `InvoiceRepository` con `Create`, `List` (filtros staff) si el spec staff-CRUD lo exige.
- **RBAC**: reutilizar `userRole` + `domain.User` (`IsClient`, `IsEmployee`, `CanManageUsers`) como en `car_service` / `appointment_service`.

## Testing Strategy

| Capa | Qué | Cómo |
|------|-----|------|
| Unit | Servicios recibidas/billing/supplier + reglas RBAC | Tablas stub como `car_service_test` |
| Integration | `httptest` rutas nuevas + JWT por rol | Patrón `mvp_role_access_test` / car integration |
| E2E | Opcional P1.1 | Playwright no requerido en v1 |

## Migration / Rollout

- Nuevas tablas: `received_invoices`, `billing_documents`, `suppliers`; migración idempotente o AutoMigrate en entornos dev.
- Tabla `invoices` si aún no existe en todas las bases: crear vía AutoMigrate del `domain.Invoice` ya definido.
- Sin feature flag obligatorio; rollback = revert deploy + `DROP TABLE` ordenado (hijos primero) si hiciera falta.

## Open Questions

- [ ] Catálogo cerrado de `BillingDocument.kind` en v1 (`client_invoice` | `payroll` | `irs` | `other`).
- [ ] ¿Unificar en UI “Billing” en una sola app section con pestañas por `kind` o rutas separadas?
- [ ] Numeración oficial / PDF: fuera de alcance P1 salvo decisión explícita en tasks.
