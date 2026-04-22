# Tasks: workshop-mechanic-vehicle-lifecycle (ServiceJob / taller)

## Phase 1: Domain and persistence

- [x] 1.1 Añadir tipos `ServiceJob`, `ServiceJobReception`, `ServiceJobHandover` y estados (open, in_progress, closed, cancelled) en `backend/internal/domain/`; `TableName` y validaciones básicas.
- [x] 1.2 Definir `ports` del repositorio (Create, GetByID, ListByCarID opcional, SaveReception, SaveHandover, transición estado) alineado con reparos.
- [x] 1.3 Implementar `internal/repository/postgres` para el agregado (1:1 recepción y cierre); índices por `car_id` y `service_job_id` según consultas.
- [x] 1.4 Incluir modelos en `AutoMigrate` y wiring en `backend/cmd/api/main.go` (y columnas/ensure como en `repairs` si aplica).
- [ ] 1.5 (Opcional este corte) Añadir `repairs.service_job_id` UUID NULL + FK: deixado comentado en `main.go` — aplicar cando se enlace execución a visita.

## Phase 2: Service and rules

- [x] 2.1 Implementar `internal/service/.../service_job` con permisos de coche como `RepairService` (empleado/manager/admin; cliente no muta).
- [x] 2.2 Reglas: crear job en *open*; requisitos de `PUT` recepción; *closed* solo tras *handover* validado; auditoría mínima quien/cuándo.
- [x] 2.3 Stub: si existen rutas OBD/estimate, responder 501 o error documentado sin 200 engañoso (spec *stub*).
- [x] 2.4 (TDD) `TestServiceJob_...` unitarios para validación recepción incompleta y cierre (escenarios *workshop-repair-execution*). Ver `openspec/config.yaml` *strict_tdd*.

## Phase 3: HTTP and wiring

- [x] 3.1 Añadir `ServiceJobHandler` (Gin): `POST /api/v1/service-jobs`, `GET` por id, `PUT` recepción y *handover*; cuerpos JSON con `schema_version` o campos fijos acordados en implementación.
- [x] 3.2 Registro de rutas en `main.go` bajo el mismo `protected` que otras; usar `RequireWorkshopStaff` o lógica en servicio según *design* (prohibir mutación a `client`).
- [x] 3.3 Registro *swagger* si el proyecto lo mantiene para reparos (misma convención). — *Comentarios swag añadidos; rexenerar con `swag` se pode en verify.*

## Phase 4: Role and flow tests (CI)

- [x] 4.1 En `internal/handler/mvp_role_access_test.go` (o *stub* router como *repair*): `POST /api/v1/service-jobs` con JWT *client* → no 2xx; con *employee* → 2xx con fakes; cumple delta *mvp-role-access*.
- [x] 4.2 Test de flujo: crear job → *PUT* recepción → *PUT* cierre; `GET` devuelve rastros; estado *closed* (o httptest en paquete handler dedicado).
- [x] 4.3 Ejecutar `cd backend && go test ./... -count=1 -race -timeout=2m` (comando CI *openspec/config*). — *Rexecutado con `go test ./...` (Windows: `-race` requiere CGO).*

## Phase 5: Frontend (MVP1 mínima)

- [x] 5.1 Cliente API en `frontend/src/lib/api.ts`: crear visita, leer, recepción, cierre; tipos TypeScript alineados con DTOs backend.
- [x] 5.2 Ruta bajo *staff*, p.ej. `frontend/src/app/workshop/` (lista y detalle) o integración con navegación *employee*; solo visible para personal (patrón *admin users* / *AppShell*).
- [x] 5.3 Formularios mínimos recepción (km, notas) y cierre; mensajes de error 4xx legibles; sin duplicar ruta cliente-only.

## Phase 6: Polish

- [ ] 6.1 *Smoke* manual: *seed* empleado, crear visita en coche de prueba, cierre; anotar en *verify* posterior.
- [ ] 6.2 (Opcional) *feature flag* env para desregistrar `service-jobs` vía `main.go` — sólo se aplica cando faga falla o rollback.
