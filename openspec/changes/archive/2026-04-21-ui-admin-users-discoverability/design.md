# Design: Descubrimiento UI de `/admin/users` (staff)

## Technical approach

Extend `AppShell` with a **conditional** nav item using existing `canManageUsers` from `@/types/user`. Add `admin_users` to `AppShellNavId`; pass `activeNav="admin_users"` from `admin/users/page.tsx`. Remove dead `frontend/src/lib/navigation.ts` (unused `getNavigationConfig`) to satisfy single-source rule—restore only if a second consumer appears.

## Architecture decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Where to put nav | `AppShell.tsx` | Matches all staff pages; same `isClient` / `canManageUsers` patterns. |
| `navigation.ts` | Delete file | Zero imports; avoids dual config. |
| Optional dashboard CTA | Defer to apply if quick | Spec marks MAY; keeps first slice small. |
| Admin `activeNav` | New enum value | Avoid false «Painel» active on `/admin/users`. |

## Data flow

```
User (Zustand) → AppShell(user)
       → canManageUsers(user)? → show «Utilizadores» → router.push('/admin/users')
       → activeNav === 'admin_users' → active style
```

## File changes

| File | Action | Description |
|------|--------|-------------|
| `frontend/src/components/layout/AppShell.tsx` | Modify | Import `canManageUsers`; add nav button + `AppShellNavId` value `admin_users`. |
| `frontend/src/app/admin/users/page.tsx` | Modify | `activeNav="admin_users"`. |
| `frontend/src/lib/navigation.ts` | Delete | Unused; single source in `AppShell`. |
| Vitest (new or extend) | Add/Modify | Render `AppShell` with mock manager/client; assert link presence/absence per design. |
| `frontend/src/app/dashboard/page.tsx` | Optional | CTA block with same guard if product wants it in same change. |

## Interfaces / contracts

- `AppShellProps.activeNav` accepts `'admin_users'`.
- No API or store contract changes.

## Testing strategy

| Layer | What | How |
|-------|------|-----|
| Unit / component | Nav visibility by role | RTL: render `AppShell` with mock `User` (`manager` / `client`), `queryByRole('link', { name: … })` or button `href` via router mock |
| Regression | Employees unchanged | Existing tests untouched |

## Migration / rollout

No migration. Deploy front-only.

## Open questions

- [ ] Exact pt_PT label: «Utilizadores» vs «Conta de utilizadores» (product copy).
