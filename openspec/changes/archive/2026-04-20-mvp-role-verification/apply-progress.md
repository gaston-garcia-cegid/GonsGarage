# Apply progress — mvp-role-verification

**Mode**: Strict TDD (`openspec/config.yaml`).

## TDD cycle evidence

| Task | Test / evidence | RED | GREEN | Triangulate | Refactor |
|------|-----------------|-----|-------|-------------|----------|
| 1.1 | Code review `domain.User` + login | N/A (discovery) | Documented: JWT role from `users.role`; `employees` table not required for auth seed | — | — |
| 1.2 | `go run ./cmd/seed-mvp-users` (manual on dev DB) | N/A (new cmd) | `cmd/seed-mvp-users/main.go` idempotent loop | 3 roles | — |
| 1.3 | Env in `main.go` header | — | Defaults `*.gonsgarage.local` | — | — |
| 2.1–2.5 | `mvp_role_access_test.go` employees | Tests added first (would fail if routes wrong) | 403 client/employee, 200 manager | 3 roles | — |
| 3.1–3.3 | Same file + stubs `mvpUserRepo`/`mvpCarRepo`/`mvpRepairRepo` | POST client → not 2xx | 403 client, 201 employee | — | Interface stubs consolidated |
| 4.1–4.2 | Docs | — | `mvp-solo-checklist.md` + `application-analysis.md` | — | — |
| 5.1 | `go test ./internal/handler/...` | — | Pass (no `-race` on Windows agent) | — | — |
| 5.2 | Double seed | — | **Manual**: run twice against local Postgres when available | — | — |
| 5.3 | `go test ./...` | — | Pass | — | — |

## Files touched

- `backend/cmd/seed-mvp-users/main.go` (new)
- `backend/internal/handler/mvp_role_access_test.go` (new)
- `docs/mvp-solo-checklist.md`, `docs/application-analysis.md`

## Deviations

None from `design.md`.

## Issues

- **5.1 / 5.2**: `-race` and double `go run` seed not executed in this environment (Windows / no Postgres in agent); CI + local dev should run them.
