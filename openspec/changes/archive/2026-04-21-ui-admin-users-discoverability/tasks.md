# Tasks: Descubrimiento UI de `/admin/users` (staff)

## Phase 1: RED (Vitest — spec `staff-user-management-ui` + delta `mvp-role-access`)

- [x] 1.1 Añadir `frontend/src/components/layout/AppShell.test.tsx` (o extender test existente del layout): render con `user` manager (`canManageUsers`) → existe control accesible que navega a `/admin/users` (p. ej. `getByRole('button', { name: /utilizador/i })` + mock `next/navigation`).
- [x] 1.2 Mismo archivo: render con `user` client → **no** existe control que enlace o empuje a `/admin/users`.
- [x] 1.3 `cd frontend && pnpm test -- --run AppShell` (o patrón equivalente) y comprobar que **fallan** hasta implementar la nav en `AppShell.tsx`.

## Phase 2: GREEN (tipos + shell + página admin)

- [x] 2.1 En `frontend/src/components/layout/AppShell.tsx`: ampliar `AppShellNavId` con `admin_users`; importar `canManageUsers`; botón pt_PT «Utilizadores» solo si `canManageUsers(user)`, `router.push('/admin/users')`, clase activa si `activeNav === 'admin_users'`.
- [x] 2.2 En `frontend/src/app/admin/users/page.tsx`: cambiar `activeNav` a `"admin_users"`.
- [x] 2.3 Eliminar `frontend/src/lib/navigation.ts` y cualquier referencia rota (grep `navigation` / `getNavigationConfig`).
- [x] 2.4 `pnpm typecheck` en `frontend/` verde.

## Phase 3: REFACTOR + verificación

- [x] 3.1 Revisar duplicación de strings de nav; extraer constante de label solo si mejora legibilidad (opcional, mínimo).
- [x] 3.2 `pnpm test -- --passWithNoTests` y `pnpm build` en `frontend/` verdes.
- [x] 3.3 Comprobación manual: login `manager`/`admin` → clic «Utilizadores» → `/admin/users`; login `client` → sin entrada.

## Phase 4: Opcional (MAY — proposal)

- [x] 4.1 CTA en `frontend/src/app/dashboard/page.tsx` hacia `/admin/users` con el mismo guard `canManageUsers` (solo si se desea en este cambio).

### Referencia

`openspec/changes/ui-admin-users-discoverability/design.md`, `specs/staff-user-management-ui/spec.md`, `specs/mvp-role-access/spec.md`.

### Orden

1.1 → 1.2 → 1.3 (RED) → 2.1–2.4 (GREEN) → 3.* → 4.* opcional.
