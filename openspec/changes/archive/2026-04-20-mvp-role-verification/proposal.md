# Proposal: Verificación MVP por rol, seeds y TDD

## Intent

El MVP se valida hoy sobre todo como **cliente** (`seed-test-client`). Faltan **usuarios seed** para **admin**, **manager** y **employee**, y una **matriz rol × funcionalidad** explícita (dashboard, coches, citas, repairs, employees, invoices). Sin eso no hay smoke reproducible ni regresiones detectables por rol. **Invoices**: hay servicio en código pero **sin HTTP/UI** en MVP v1 (`openspec/specs/p1-accounting-defer/spec.md`); la matriz debe reflejarlo, no prometer UI inexistente.

## Scope

### In Scope

- Matriz **rol × área** (doc + spec): acceso API/UI esperado; filas “diferido” donde aplique (invoices).
- **Seeds** idempotentes (env `SEED_*`), patrón `seed-test-client`, **≥1 usuario por rol** con credenciales documentadas fuera del repo.
- **TDD**: tests RED→GREEN en **autorización** (middleware/handlers/servicios) y, donde sea estable, **frontend** (RTL: nav/guards).

### Out of Scope

- CRUD invoices HTTP/UI completo; i18n; pagos; multi-tenant.

## Capabilities

### New Capabilities

- **`mvp-role-access`**: Requisitos y escenarios Given/When/Then para acceso por rol a MVP; seeds mínimos; tests que bloqueen regresiones de política.

### Modified Capabilities

- **None**

## Approach

1. Cruce `cmd/api/main.go` + `frontend/src/app/` vs `docs/application-analysis.md`; corregir doc si difiere.
2. Delta `specs/mvp-role-access/spec.md` con matriz y casos 403/forbidden.
3. Nuevo `cmd/seed-mvp-users` (o varios cmds) + checklist en `mvp-solo-checklist.md`.
4. TDD: priorizar `httptest` / tests de servicio; Vitest solo donde ya hay harness.

## Affected Areas

| Area | Impact |
|------|--------|
| `backend/cmd/` | Seed(s) |
| middleware, handlers, `*_test.go` | Tests por rol |
| `frontend/**/*.test.tsx` | Tests nav/guard opcionales |
| `docs/mvp-solo-checklist.md` | Enlace matriz + env vars |

## Risks

| Risk | L | Mitigation |
|------|---|------------|
| Seeds en BD compartida | M | emails `*.local`, idempotencia, doc “solo dev” |

## Rollback Plan

Revert commits; eliminar usuarios seed en dev por email (script o SQL documentado).

## Dependencies

- Postgres dev + `DATABASE_URL`.

## Success Criteria

- [ ] Spec `mvp-role-access` con matriz y negación.
- [ ] Seed por rol (cliente reutilizable o nuevo).
- [ ] Tests que fallen si se rompe política staff/cliente en repairs o employees.
- [ ] `go test ./...` y `pnpm test` verdes en CI.
