# Tasks: Admin provision user roles (manager / employee / client)

## Phase 1: Tests first (RED)

- [ ] 1.1 Añadir `backend/internal/handler/admin_user_provision_test.go` (o extender `mvp_role_access_test.go`) con casos: sin JWT → no 2xx; `client`/`employee` → 403; `admin`+body `role=admin` → no 2xx; `admin`+`client` válido → 201 (mock repo si hace falta); `manager`+`manager` → no 2xx; `manager`+`employee` → 201.
- [ ] 1.2 `cd backend && go test ./internal/handler/... -count=1 -run Provision` y comprobar que **fallan** hasta existir la ruta.

## Phase 2: Contrato y dominio

- [ ] 2.1 En `backend/internal/core/ports/services.go` (o subpaquete ports): tipo request `ProvisionUserRequest` (email, password, firstName, lastName, role) y método `ProvisionUser(ctx, callerUserID, callerRole, req) (*domain.User, error)` en `AuthService` (o interfaz dedicada si preferís no hinchar auth).
- [ ] 2.2 En `backend/internal/service/auth/` (o servicio de usuarios que use `UserRepository`): implementar gates de `staff-user-provisioning` (admin vs manager vs target role); reutilizar hash/reglas de `Register`; rechazar `admin` y roles desconocidos; errores 409 email duplicado alineados a registro.

## Phase 3: HTTP (GREEN wiring)

- [ ] 3.1 Nuevo handler `backend/internal/handler/admin_user_handler.go` (nombre final acorde): `POST` JSON, lee `userID`+rol del JWT (`c.Get`), llama al servicio, responde 201 con usuario sin password.
- [ ] 3.2 En `backend/cmd/api/main.go`: grupo bajo `protected` (p. ej. `/admin/users`) con `middleware.GinBearerJWT` + `middleware.RequireStaffManagers()` (o middleware que refleje exactamente la matriz si `RequireStaffManagers` es demasiado permisivo para `manager`→`manager`).
- [ ] 3.3 Registrar handler en `setupRoutes` / inyección igual que `employeeHandler`.

## Phase 4: Verificación y documentación API

- [ ] 4.1 `go test ./... -count=1` en `backend/` verde; ajustar tests 1.1 hasta verde.
- [ ] 4.2 Anotaciones swag en el handler nuevo y `swag init` / `docs.go` si el repo regenera Swagger en CI.
- [ ] 4.3 Una línea en `deploy/README.md` o `docs/application-analysis.md` si hace falta mencionar el nuevo path (opcional).

## Phase 5: Frontend (opcional, misma change o follow-up)

- [ ] 5.1 `frontend/src/lib/api-client.ts`: método `provisionUser(body)` → `POST /api/v1/admin/users` (ajustar path al elegido en 3.2).
- [ ] 5.2 Página o secção bajo `AppShell` (solo `admin`/`manager` según UI): formulario pt_PT; `go`/`no-go` según `useAuth().user.role`.
- [ ] 5.3 `pnpm test` + `pnpm typecheck` en `frontend/`.

### Orden recomendado

1.1 → 1.2 (RED) → 2.1–2.2 (dominio) → 3.1–3.3 (ruta) → 4.1–4.2 (verde total) → 5.* si aplica.

### Referencia de escenarios

`openspec/changes/admin-provision-user-roles/specs/staff-user-provisioning/spec.md` y delta `specs/mvp-role-access/spec.md` (CI + matriz).
