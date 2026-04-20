# Apply progress — mvp-post-v1-followup

**Change**: mvp-post-v1-followup  
**Mode**: Strict TDD (config `strict_tdd: true`; tarea 3.1 = solo test de regresión).

## TDD cycle evidence

| Task | RED | GREEN | REFACTOR |
|------|-----|-------|----------|
| 3.1 Admin `GET /employees` | Comportamiento ya en `RequireStaffManagers` (`RoleAdmin`); RED = brecha de cobertura respecto a tabla spec `mvp-role-access` (admin MUST). Test nuevo describe el contrato. | `go test ./internal/handler -run TestMVPAccess_EmployeesGET_Admin -count=1` PASS | Ninguno (test espejo del caso manager). |

## Archivos tocados

- `docs/mvp-next-steps.md` (nuevo)
- `docs/mvp-solo-checklist.md`, `docs/mvp-minimum-phases.md`, `docs/roadmap.md`
- `openspec/changes/archive/2026-04-20-readme-verified-refresh/**`, `.../2026-04-20-mvp-funcionando-plan/**`
- `openspec/changes/archive/2026-04-19-mvp-gap-roadmap-2026/proposal.md` (enlace)
- `backend/internal/handler/mvp_role_access_test.go`

## Desviaciones

Ninguna respecto a `proposal.md`.

## Issues

Ninguno.
