# Apply progress — ci-appointment-test-deterministic

**Mode**: Strict TDD (per `openspec/config.yaml`). This batch is **test-fix only** (no new production API).

## TDD cycle evidence

| Task | Test / evidence | Layer | Safety net | RED | GREEN | Triangulate | Refactor |
|------|-----------------|-------|------------|-----|-------|-------------|----------|
| 1.1 | `TestAppointmentService_UpdateAppointment_RepoError` | Unit | `go test ./internal/service/appointment/... -run … -count=20` baseline | Fixture could fail CI outside business hours (flaky) | `ScheduledAt` fixed `time.Date(..., time.Local)` 10:00 | Single path (repo err); window choice arbitrary but valid | Comment only |
| 1.2 | Same test + `UpdateAppointment` merge logic | — | Code review | N/A | `patch` has zero `ScheduledAt` → merged keeps fixture hour | — | — |
| 2.1 | `rg time.Now(` on `*_test.go` | — | No other hits | — | Only occurrence was 1.1 | — | — |
| 2.2 | Optional helper | — | N/A | — | Skipped: single site after 2.1 | — | — |
| 3.1 | Package stress | — | — | — | `go test ./internal/service/appointment/... -count=10` OK | — | — |
| 3.2 | Full backend | — | — | — | `go test ./... -count=1` OK | — | — |

**Local note (3.1)**: `-race` requires CGO/gcc; not available on this Windows agent. CI (Linux) should still run `-race` per `.github/workflows/ci.yml`.

## Deviations

- None for product code; **task 3.1** verified with `-count=10` without `-race` locally only.
