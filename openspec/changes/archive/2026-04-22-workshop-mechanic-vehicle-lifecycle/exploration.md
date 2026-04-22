## Exploration: workshop-mechanic-vehicle-lifecycle (empleado / ciclo taller)

### Current State

- **Dominio `Repair`** (`backend/internal/domain/repair.go`): registro plano con `CarID`, `TechnicianID`, `Description`, `Status` (`pending|in_progress|completed|cancelled`), `Cost`, `StartedAt`, `CompletedAt`, soft delete. No hay checklist de recepción, OBD, presupuesto, aprobación, entrega ni historial de fases.
- **API** (`cmd/api/main.go`): `GET /api/v1/repairs/car/:carId`, `POST/GET/PUT/DELETE /api/v1/repairs/...`. Protegida con JWT; `RepairService` exige `user.IsEmployee()` para **crear**; clientes leen solo reparaciones de sus coches.
- **Servicio** (`internal/service/repair/repair_service.go`): validación mínima, duplicado por misma descripción + `StartedAt`. No modela flujo taller.
- **Frontend**: `cars/[id]/page.tsx` usa `apiClient.getRepairs/createRepair/updateRepair/deleteRepair` — CRUD en ficha coche; no hay área “taller” dedicada para mecánico.
- **Specs**: `mvp-role-access` ya matriz repares (cliente no muta; empleado puede según tests). La Fase C de `docs/mvp-minimum-phases.md` deja **POST/PATCH/DELETE + UI staff** como opcional MVP+.

### Affected Areas

- `backend/internal/domain/repair.go` — hoy insuficiente para el ciclo propuesto; o se extiende mucho o se añade agregado.
- `backend/internal/service/repair/repair_service.go` / `internal/handler/repair_handler_gin.go` — nuevas operaciones o nuevas rutas.
- `backend/internal/repository/postgres/repair_repository.go` — esquema y consultas.
- `backend/cmd/api/main.go` — registro de rutas y AutoMigrate.
- `frontend/src/lib/api.ts` (o `carApi`) — contratos; `frontend` nueva ruta staff o sección empleado.
- `openspec/specs/mvp-role-access/spec.md` — permisos finos por recurso/estado si aplica.

### Approaches

1. **Extender `repairs` con JSON / columnas nullable (fases en un solo registro)**  
   - Pros: menos tablas; migración incremental.  
   - Cons: mezcla recepción, OBD, presupuesto y cierre en un modelo que hoy es “una reparación”; riesgo de columnas gigantes o JSON opaco; difícil consultar por fase.  
   - Effort: **Medium** (corto plazo) / **High** (mantenimiento).

2. **Nuevo agregado `ServiceJob` / `WorkOrder` (1 vehículo en taller) con sub-entidades o estados**  
   - Pros: alinea con el lenguaje de negocio; `Repair` puede seguir como “línea de trabajo” o quedar referenciado; evolución clara por fases.  
   - Cons: más tablas, migración y UI nuevas; decidir relación con `Repair` existente y con `appointments`.  
   - Effort: **High** (correcto a medio plazo).

3. **Mantener `Repair` como “trabajo final” y tablas satélite** (`repair_intake`, `repair_obd_reading`, `repair_estimate`) ligadas por `repair_id`  
   - Pros: menores cambios de ruta si el work order = un `Repair` creado al inicio; satélites añaden datos por fase.  
   - Cons: orden de creación (¿crear repair al recibir o al aprobar?); riesgo de huérfanos si se cancela.  
   - Effort: **Medium–High**.

### Recommendation

**Enfoque 2 o 3 en versión reducida para MVP1:** introducir un **identificador de trabajo en taller** (agregado o `Repair` ampliado con **estado de workflow** explícito) y **al menos** persistir **recepción (checklist)** + **cierre (checklist)** como sub-recursos o filas hijas; dejar OBD y presupuesto como **campos opcionales o stub** hasta segundo corte. No acumular todo en un JSON monolítico sin contrato de spec.

Así se desbloquea el mecánico sin reescribir el historial de `repairs` ya listado por coche; la migración puede marcar filas antiguas como “legacy / sin intake”.

### Risks

- **Alcance épico** frente a un solo par de tablas y 3 endpoints.  
- **Duplicación** con la UI actual de reparaciones en `cars/[id]` si no se unifica el mensaje de producto.  
- **OBD real** no bloquea el diseño de dominio: el spec puede exigir “adjunto o payload estructurado” sin protocolo hardware.

### Ready for Proposal

**Yes** — el `proposal.md` existente sigue siendo válido como intención; esta exploración **acota la decisión de modelado** y justifica **sdd-spec** + **sdd-design** con elección explícita (agregado vs satélites). Falta decisión de producto: **¿un `Repair` por visita al taller o un work order que pueda generar varias “líneas” de coste?**

---

*Generado como artefacto SDD `exploration` para el change `workshop-mechanic-vehicle-lifecycle`.*
