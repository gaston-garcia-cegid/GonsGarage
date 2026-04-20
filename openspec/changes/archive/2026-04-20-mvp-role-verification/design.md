# Design: MVP role verification, seeds, and TDD

## Technical approach

Implement **`mvp-role-access`** by (1) **one idempotent seed command** mirroring `cmd/seed-test-client` (GORM + `postgresRepo.NewPostgresUserRepository`, no AutoMigrate in cmd), (2) **HTTP-level regression tests** with `gin` + `httptest` + JWTs signed like `gin_jwt_test.go`, mounting the **same route groups** as `setupRoutes` for `/employees` and `/repairs`, and (3) **docs** linking the spec matrix from `docs/mvp-solo-checklist.md`. **Repair authorization** stays in `RepairService` (`IsEmployee()` rejects client — `repair_service.go`); tests prove the **Gin stack** returns non-2xx for client mutations. **Employees** stay behind `middleware.RequireStaffManagers()` on `/api/v1/employees` (`main.go`).

## Architecture decisions

| Decision | Choice | Alternatives | Rationale |
|----------|--------|--------------|-----------|
| Seed layout | **Single** `cmd/seed-mvp-users` with env-driven emails/passwords per role | Four cmds; only docs | One `go run`, same idempotency pattern as `seed-test-client` |
| Keep `seed-test-client` | **Yes** | Remove | Backward compat; README/checklist can say “client also covered by legacy cmd” |
| Where to test HTTP auth | **New** `*_test.go` next to handlers or under `internal/handler` | E2E Playwright | Matches spec CI focus; `gin_jwt_test.go` already shows JWT signing |
| Handler deps | **Stubs/mocks** for repos + real `RepairService` where cheap | Full integration DB | Keeps CI fast; service layer already tested elsewhere |
| Frontend TDD | **Optional** follow-up in same change if time | Mandatory RTL | Spec prioritizes API; shell nav is role-agnostic today |

## Data flow

```
seed-mvp-users → PostgresUserRepository.Create (skip if GetByEmail exists)
       ↓
   users table (role column)

HTTP test → JWT (claims userID, role) → GinBearerJWT → RequireStaffManagers? → handler → service
```

## File changes

| File | Action | Description |
|------|--------|-------------|
| `backend/cmd/seed-mvp-users/main.go` | Create | Idempotent create for admin, manager, employee (+ optional client env or skip if using `seed-test-client`) |
| `backend/internal/handler/*_test.go` or `routes_auth_test.go` | Create | `httptest`: `GET /employees` 403 for client/employee; 200/401? for manager with mocks; `POST /repairs` non-2xx for client |
| `docs/mvp-solo-checklist.md` | Modify | Link to `openspec/.../mvp-role-access/spec.md` + env table for seeds |
| `docs/application-analysis.md` | Modify | Align matrix row with tests if any drift |

## Interfaces / contracts

- **Env vars** (defaults documented, `*.local` emails): e.g. `SEED_ADMIN_EMAIL`, `SEED_ADMIN_PASSWORD`, … `SEED_EMPLOYEE_*`, `SEED_MANAGER_*`; client optional `SEED_CLIENT_EMAIL` or defer to existing `seed-test-client`.
- **JWT claims** for tests: `userID` (string UUID), `role` (`client` | `employee` | `manager` | `admin`), `email` — same as `GinBearerJWT` expectations.

## Testing strategy

| Layer | What | How |
|-------|------|-----|
| Unit | `RepairService` client denied | Already implied by `IsEmployee`; extend if gaps |
| Integration | Gin routes + middleware | `httptest`, `gin.TestMode`, signed JWTs (`middleware/gin_jwt_test` pattern) |
| E2E | — | Out of scope per proposal |

## Migration / rollout

No production migration. Dev-only seeds; document “never run against prod DB”.

## Open questions

- [ ] Seeds for **staff** need **linked `employees` rows** or only `users`? (Employees domain may have separate table — verify before seed implementation.)
- [ ] Minimal **car** (+ owner) seed for staff repair POST happy path, or mock repos only for “employee allowed” assertion?
