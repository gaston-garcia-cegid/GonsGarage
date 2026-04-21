# Proposal: Descubrimiento UI de `/admin/users` (staff)

## Intent

La página `frontend/src/app/admin/users/` y el aprovisionamiento ya existen; el layout usa `canManageUsers` (admin/manager). El problema es **descubrimiento**: `AppShell` (`frontend/src/components/layout/AppShell.tsx`) solo enlaza Painel, Viaturas, Marcações y luego facturas o contabilidad — **no hay ruta visible a `/admin/users`**. `frontend/src/lib/navigation.ts` lista `/admin/users` pero **no se importa** (config huérfana). Hay que exponer la función en la UI sin cambiar API ni roles.

## Scope

### In Scope

- Enlace en nav principal **solo** si `canManageUsers(user)`; etiquetas **pt_PT** alineadas al shell.
- Cambios en `AppShell` + `activeNav` en `frontend/src/app/admin/users/page.tsx` (hoy `activeNav="dashboard"`; conviene clave dedicada p. ej. `admin_users`).
- Decisión explícita sobre `navigation.ts`: **integrar** con el shell o **eliminar** para evitar duplicar mapas.
- Opcional: CTA en `frontend/src/app/dashboard/page.tsx` para el mismo guard.

### Out of Scope

- Backend, JWT, matriz de roles en servidor. i18n más allá del pt_PT existente. Rediseño grande del shell.

## Capabilities

### New Capabilities

- `staff-user-management-ui`: Shell autenticado con entrada descubrible a aprovisionamiento (`/admin/users`) para quien `canManageUsers`, sin exponerla a otros roles.

### Modified Capabilities

- `mvp-role-access`: Referencia o fila alineada para que la **UI** de gestión staff no dependa solo de URL manual, coherente con el aprovisionamiento ya documentado.

## Approach

1. **Nav:** Botón condicional (p. ej. “Utilizadores”) → `/admin/users`; estado activo cuando la ruta sea `/admin/users`.
2. **CTA dashboard (opcional):** Mismo guard.
3. **`navigation.ts`:** Unificar con `AppShell` o borrar si sigue muerto.
4. **Tests:** Cubrir presencia/ausencia del enlace por rol si el repo ya testea layout.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `frontend/src/components/layout/AppShell.tsx` | Modified | Ítem condicional; ampliar `AppShellNavId` si aplica. |
| `frontend/src/app/admin/users/page.tsx` | Modified | `activeNav` correcto. |
| `frontend/src/app/dashboard/page.tsx` | Optional | CTA. |
| `frontend/src/lib/navigation.ts` | Modified/Removed | Cablear o eliminar. |
| `openspec/specs/mvp-role-access/spec.md` | Modified (delta) | Coherencia documental. |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Nav saturada | Med | Label corto; revisar CSS responsive existente. |
| Duplicar config nav | Med | Una sola fuente (`navigation.ts` resuelto). |

## Rollback Plan

Revertir cambios en `AppShell`, `admin/users`, dashboard opcional y estado de `navigation.ts`. Sin datos que migrar.

## Dependencies

- Ninguna externa.

## Success Criteria

- [ ] Admin/manager **ve** enlace a gestión de usuarios sin escribir URL; client/employee sin permiso **no**.
- [ ] Sin config de nav **huérfana** sin decisión (integrar o borrar `navigation.ts`).
- [ ] Estado activo en `/admin/users` coherente con la ruta.
