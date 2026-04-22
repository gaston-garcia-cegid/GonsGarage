# Tasks: Taller — benchmark (recepción, día, reparaciones)

> Delta: `specs/workshop-repair-execution/spec.md` · `design.md` inexistente (inferido de código y delta).

## Phase 1: Backend — listado por día

- [x] 1.1 Fijar en código/Swag que el “día” usa ventana en **UTC** [00:00, 24:00) sobre `ServiceJob.OpenedAt` (o documentar otra regla en el mismo `GET`).
- [x] 1.2 Añadir `ListByOpenedOn` (o nombre equivalente) a `ServiceJobRepository` en `backend/internal/core/ports/repositories.go` e implementación en `backend/internal/repository/postgres/service_job_repository.go`.
- [x] 1.3 **TDD (RED)**: en `backend/internal/service/servicejob/service_job_service_test.go` casos: día sin jobs → `[]` sin error; día con apertura → al menos un ID (stub/fakes existentes).
- [x] 1.4 Método en `backend/internal/service/servicejob/` que delegue al repo; **(GREEN)** hasta verde.
- [x] 1.5 Nuevo handler (p. ej. `ListServiceJobsByOpenedOn`) con query `opened_on=YYYY-MM-DD` en `backend/internal/handler/service_job_handler_gin.go`; anotar `@Router` y `@Param`.
- [x] 1.6 En `backend/cmd/api/main.go`, registrar `GET ""` bajo `service-jobs` **antes** de `GET "/:id"`; mismo middleware `RequireWorkshopStaff` que hoy.
- [x] 1.7 `backend/internal/handler/mvp_role_access_test.go` (o prueba de integración equivalente): staff → 200/[] o 200+datos; fuera de staff según reglas actuales.

## Phase 2: Backend/Domain — enlace `repairs` (SHOULD en delta)

- [x] 2.1 Si aún no existe: migración y campo opcional `service_job_id` en `repairs` (`backend/internal/domain/repair.go` + GORM) alineada con criterio único de producto; sin backfill obligatorio.
- [x] 2.2 Repositorio: listar IDs de reparas por `service_job_id` (método acotado en `ports` o `RepairRepository` existente).
- [x] 2.3 `servicejob` al armar el detalle de `GET /service-jobs/:id`, incluir `repair_ids` (o array vacío) en la respuesta JSON; **MUST NOT** devolver obras inexistentes.
- [x] 2.4 Regenerar/actualizar Swagger en `backend` si el proyecto lo mantiene en el flujo de CI.

## Phase 3: API cliente + UI — “Hoy” y detalle

- [x] 3.1 `frontend/src/lib/api.ts`: tipo ampliado con `repair_ids?`; método `listServiceJobsByOpenedOn(date: string)`.
- [x] 3.2 `frontend/src/app/workshop/page.tsx`: sección “Hoy” que llame al `GET` con fecha actual (formateo UTC/ISO acorde a 1.1) y muestre lista vacía sin error; link a `workshop/{id}`.
- [x] 3.3 `frontend/src/app/workshop/[id]/page.tsx` (o componente compartido): sección reparas vinculadas (0..n) a partir de `repair_ids` / contrato de detalle unificado.

## Phase 4: Superficie recepción “app”

- [x] 4.1 Crear ruta dedicada, p. ej. `frontend/src/app/workshop/recepcion/page.tsx` o `recepcion/[jobId]`, reutilizando el mismo cuerpo `putServiceJobReception` y campos del detalle, sin API nueva.
- [x] 4.2 Navegación: entrada clara desde `workshop` (botón/ link “Recepção” / “Receção rápida`) y, si aplica, estilos `frontend/src/app/workshop/workshop.module.css` móvil-first.
- [x] 4.3 (Opcional) Ajuste responsive mínimo en el detalle actual para ancho móvil sin duplicar lógica de negocio.

## Phase 5: Cierre y verificación

- [x] 5.1 `go test` backend; `cd frontend && pnpm lint` + `typecheck` + `pnpm test -- --passWithNoTests` según `openspec/config.yaml`.
- [x] 5.2 Comprobar escenarios del delta: recepción vía ruta dedicada, listado día vacío, día con item, `repair_ids` vacío vs con IDs, 4xx recepción inalterados.
