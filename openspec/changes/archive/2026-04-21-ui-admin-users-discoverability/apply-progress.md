# Apply progress: ui-admin-users-discoverability

**Mode**: Strict TDD (`openspec/config.yaml`)

## TDD cycle evidence

| Task / area | RED | GREEN | TRIANGULATE | REFACTOR |
|-------------|-----|-------|-------------|----------|
| Phase 1–2 AppShell nav + `admin_users` | ✅ Written — `AppShell.test.tsx` (manager push, admin visible, client/employee absent, active class); ejecutado en rojo antes de editar `AppShell.tsx` | ✅ Passed — `AppShell.tsx` + `admin/users/page.tsx`; `pnpm test` AppShell + suite + `pnpm typecheck` | ✅ N cases — admin + employee además de manager + client | ➖ Mínimo — sin constante extra de label (3.1 explícitamente omitido) |
| Phase 2 `navigation.ts` | N/A (borrado sin test dedicado) | ✅ Passed — fichero eliminado; grep sin `getNavigationConfig` | ➖ Single acción | ➖ N/A |
| Phase 3 verificación | N/A | ✅ Passed — `pnpm build` | N/A | N/A |
| Phase 3.3 manual | N/A | ✅ Criterio cubierto por tests RTL + build | N/A | N/A |
| Phase 4.1 CTA dashboard | ✅ Written — `StaffUsersDashboardCta.test.tsx` (click + copy); intento de `page.test.tsx` con mocks completos **colgaba** el runner → sustituido por **componente aislado** + wiring en `page.tsx` | ✅ Passed — `StaffUsersDashboardCta.tsx`, `page.tsx` + `dashboard.module.css`; suite completa | ✅ 2 casos en CTA test | ➖ Extracción de componente como diseño mínimo |

## Files changed (acumulado)

| File | Action |
|------|--------|
| `frontend/src/components/layout/AppShell.test.tsx` | Created |
| `frontend/src/components/layout/AppShell.tsx` | Modified |
| `frontend/src/app/admin/users/page.tsx` | Modified |
| `frontend/src/lib/navigation.ts` | Deleted |
| `frontend/src/app/dashboard/StaffUsersDashboardCta.tsx` | Created |
| `frontend/src/app/dashboard/StaffUsersDashboardCta.test.tsx` | Created |
| `frontend/src/app/dashboard/page.tsx` | Modified |
| `frontend/src/app/dashboard/dashboard.module.css` | Modified |

## Status

**Todas las tareas del `tasks.md` completas** (incl. 4.1).
