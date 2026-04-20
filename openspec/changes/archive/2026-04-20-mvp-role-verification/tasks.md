# Tasks: MVP role verification (seeds + TDD HTTP)

## Phase 1: Discovery + seed command

- [x] 1.1 Confirmar en código que `users` basta para login JWT por rol (`domain.User` + `postgresRepo.NewPostgresUserRepository`); anotar si hace falta fila en `employees` para otro flujo (resolver open question del `design.md`).
- [x] 1.2 Crear `backend/cmd/seed-mvp-users/main.go`: mismo patrón que `cmd/seed-test-client/main.go` (GORM, sin AutoMigrate, `GetByEmail` → skip idempotente, `domain.NewUser` con `RoleAdmin` / `RoleManager` / `RoleEmployee`; opcional `SEED_CLIENT_*` o README que apunte a `seed-test-client`).
- [x] 1.3 Variables de entorno documentadas: `SEED_ADMIN_EMAIL|PASSWORD`, `SEED_MANAGER_*`, `SEED_EMPLOYEE_*`, defaults `*.local` (sin secretos en git).

## Phase 2: TDD — middleware employees (spec: Client forbidden / Manager allowed)

- [x] 2.1 **RED**: Añadir `backend/internal/handler/mvp_role_access_test.go` (o `middleware` si preferís solo gate) con helper JWT (`jwt.SigningMethodHS256`, claims `userID`, `role`, `email` alineados a `GinBearerJWT`).
- [x] 2.2 **GREEN**: Montar `gin.Engine` test con `api := /api/v1`, `GinBearerJWT`, grupo `/employees` + `RequireStaffManagers` + handler stub `GET ""` → 200 si llega al handler.
- [x] 2.3 Test: JWT `role=client`, `GET /api/v1/employees` → **403** (escenario *Client forbidden*).
- [x] 2.4 Test: JWT `role=employee`, mismo request → **403** (mismo requisito spec).
- [x] 2.5 Test: JWT `role=manager`, mismo request → **no** 403 por middleware (200 del stub; escenario *Manager allowed past role gate*).

## Phase 3: TDD — repairs mutation (spec: Client POST repair fails)

- [x] 3.1 En el mismo archivo de test (o split): montar grupo `/repairs` con `GinBearerJWT` + `RepairHandler` inyectando `RepairService` real y repos **stub** (usuario client en `userRepo`, coche válido opcional) según `design.md`.
- [x] 3.2 Test: JWT `client`, `POST /api/v1/repairs` con JSON mínimo válido → **403** (`ErrUnauthorizedAccess` mapeado en `GinCreateRepair`; escenario *Client POST repair fails*).
- [x] 3.3 Test: JWT `employee`, mismo POST con stubs que devuelvan usuario employee + car existente → **201 o 4xx distinto de 403 por rol** (escenario *Employee may mutate* — no debe fallar solo por `IsEmployee()`).

## Phase 4: Documentación (spec: Published matrix + invoices deferred)

- [x] 4.1 En `docs/mvp-solo-checklist.md`: sección o tabla con enlace a `openspec/changes/mvp-role-verification/specs/mvp-role-access/spec.md`, fila invoices **Deferred** y comandos `go run ./cmd/seed-mvp-users` + `seed-test-client`.
- [x] 4.2 Ajustar `docs/application-analysis.md` solo si la matriz de rutas difiere de lo probado (p. ej. nota “tests HTTP en `handler/mvp_role_access_test.go`”).

## Phase 5: Verificación

- [x] 5.1 `cd backend && go test ./internal/handler/... -count=1 -race` (o sin `-race` si CGO ausente; CI Linux con `-race`).
- [x] 5.2 Ejecutar dos veces `go run ./cmd/seed-mvp-users` contra Postgres dev: segunda sin error (escenario *Re-run is safe*).
- [x] 5.3 `cd backend && go test ./... -count=1`.
