# Issue (copiar a GitHub)

**Título:** `test: GET /api/v1/employees accepts admin JWT`

**Labels sugeridos:** `test`, `backend`

---

## Contexto

Seguimiento P1 de [`docs/mvp-next-steps.md`](./mvp-next-steps.md). El middleware `RequireStaffManagers` debe permitir **admin** y **manager** en `GET /api/v1/employees`, y denegar **client** / **employee**.

## Spec

- [`openspec/specs/mvp-role-access/spec.md`](../openspec/specs/mvp-role-access/spec.md)

## Estado en repo

- Test: `backend/internal/handler/mvp_role_access_test.go` → `TestMVPAccess_EmployeesGET_AdminReachesHandler`
- Verificado en local: **PASS** (2026-06-02)

## Criterio de cierre

- [ ] Test sigue en verde en CI (`go test ./...` en `.github/workflows/ci.yml`)
- [ ] (Opcional) Añadir enlace a este issue en `mvp-next-steps.md` checkbox P1

## Notas

No requiere cambio de código salvo regresión futura. Cerrar como *done* / *verified* una vez confirmado en CI.
