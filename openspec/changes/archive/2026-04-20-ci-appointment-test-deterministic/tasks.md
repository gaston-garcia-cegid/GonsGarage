# Tasks: Test determinístico — UpdateAppointment repo error

## Phase 1: Fixture y test roto

- [x] 1.1 En `backend/internal/service/appointment/appointment_service_test.go`, en `TestAppointmentService_UpdateAppointment_RepoError`, sustituir `ScheduledAt: time.Now().UTC()` del `existing` por un instante fijo en `time.Local` dentro de 09:30–12:30 o 14:00–17:30 (p. ej. `time.Date(2026, 6, 15, 10, 0, 0, 0, time.Local)`).
- [x] 1.2 Garantizar que el `patch` no altere `ScheduledAt` a cero ni fuerce re-validación fuera de franja; mantener notas solo si `merged` sigue en ventana.

## Phase 2: Barrido y CI

- [x] 2.1 `rg "time\\.Now\\(" backend/internal/service/appointment/*_test.go` — si otro test del paquete mezcla `time.Now()` con `validateWorkshopClock`, alinear con hora fija o helper.
- [x] 2.2 (Opcional) Extraer `func workshopTimeLocal(y int, m time.Month, d, hh, mm int) time.Time` en el mismo `_test.go` o `_test_helpers.go` del paquete si se repite.

## Phase 3: Verificación

- [x] 3.1 `cd backend && go test ./internal/service/appointment/... -count=10 -race`
- [x] 3.2 `cd backend && go test ./... -count=1`
