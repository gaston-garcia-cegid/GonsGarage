# Apply progress: admin-provision-user-roles

**Mode**: Strict TDD (per `openspec/config.yaml`)

## TDD cycle evidence

| Task / area | RED | GREEN | REFACTOR |
|-------------|-----|-------|----------|
| HTTP `POST /api/v1/admin/users` (middleware + matrix) | `admin_user_provision_test.go`: no JWT, client/employee 403, admin+admin body 400, admin+client 201, manager+manager 403, manager+employee 201 | `AdminUserHandler`, `setupRoutes`, `RequireStaffManagers` on `/admin` group | Handler uses shared error mapping; Gin test router helper `newProvisionTestRouter` |
| `ProvisionUser` service + ports | `auth_service_test.go`: admin→manager, manager≠manager, target admin, unknown role, caller employee, duplicate email, trim role | `ProvisionUserRequest` + `AuthService.ProvisionUser` in `services.go` / `auth_service.go` | Clear role matrix; `callerUserID` reserved for audit |
| Frontend API + page | N/A (optional phase); typecheck caught `user` null in `onSubmit` | `apiClient.provisionUser`, `/admin/users` layout + form | `canManageUsers` guard; role `<select>` by caller |

## Summary

- **What**: Staff provisioning endpoint and UI for admin/manager creating users with roles per matrix; Swagger regenerated.
- **Where**: `backend/internal/core/ports/services.go`, `auth_service.go`, `auth_service_test.go`, `admin_user_handler.go`, `admin_user_provision_test.go`, `cmd/api/main.go`, `docs/*`, `frontend/src/lib/api-client.ts`, `frontend/src/app/admin/users/*`.

## Status

All tasks in `tasks.md` marked complete for this change batch. Ready for `sdd-verify`.
