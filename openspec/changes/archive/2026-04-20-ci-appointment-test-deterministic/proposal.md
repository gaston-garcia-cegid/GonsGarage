# Proposal: Test determinístico — UpdateAppointment repo error

## Intent

`go test ./...` falla de forma **no determinista** en `TestAppointmentService_UpdateAppointment_RepoError`: el test espera `update failed` del stub, pero el servicio devuelve antes `ErrAppointmentOutsideBusinessHours` porque `validateWorkshopClock` usa **hora local** y el fixture usa `time.Now().UTC()` sin garantizar ventana 09:30–12:30 / 14:00–17:30. Eso bloquea CI y confunde diagnósticos (parece fallo de repo).

## Scope

### In Scope

- Ajustar el test (o el fixture `ScheduledAt`) para que el flujo llegue a `repo.Update` y reciba `updateErr` del stub.
- Opcional: helper compartido `mustWorkshopLocalTime(t)` en el paquete de tests para otros casos.

### Out of Scope

- Cambiar reglas de negocio de horario del taller.
- PgBouncer / Postgres / deploy Arnela (ya cubierto en docs y script `COMPOSE_OVERRIDE`).

## Capabilities

### New Capabilities

- None

### Modified Capabilities

- None

## Approach

En `appointment_service_test.go`, fijar `ScheduledAt` a un instante que, en `time.Local`, caiga **dentro** de franja permitida (p. ej. construir `time.Date(y, m, d, 10, 0, 0, 0, time.Local)` en día conocido). Mantener `t.Parallel()` compatible (cada test su propia hora fija). Re-ejecutar `go test ./... -race` como en CI.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `backend/internal/service/appointment/appointment_service_test.go` | Modified | Fixture de hora determinística. |
| `.github/workflows/ci.yml` | None | Ya ejecuta tests; pasará al verde local. |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Otros tests del paquete con `time.Now()` y mismo patrón | Med | Grep `time.Now` en `*_test.go` del paquete en la misma tanda. |

## Rollback Plan

Revertir el commit del test; no hay migración ni API.

## Dependencies

- Ninguna externa.

## Success Criteria

- [ ] `go test ./internal/service/appointment/... -count=10` sin fallos (estrés de orden/paralelismo).
- [ ] `go test ./...` en `backend/` con exit code 0 en entorno Linux típico (CI).
