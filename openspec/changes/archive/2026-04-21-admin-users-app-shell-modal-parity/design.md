# Design: Admin Users — Shell & Modal Parity

## Technical approach

Refactor `frontend/src/app/admin/users/page.tsx` to mirror **`/appointments`**: pass **`toolbar`** to `AppShell` (fragment with `h1` + `Button`), hold `createOpen` state, render provisioning UI inside **Radix `Dialog`** (`@/components/ui/dialog`) like `accounting/suppliers/page.tsx`. Move form logic into **`ProvisionUserModal.tsx`** (props: `open`, `onOpenChange`, `roleOptions`, `defaultRole`, `onSuccess` callback). Keep **`apiClient.provisionUser`** and role matrix logic unchanged.

## Architecture decisions

| Decision | Choice | Alternatives | Rationale |
|----------|--------|----------------|-------------|
| Modal primitive | Shadcn `Dialog` | Custom div overlay | Same a11y and patterns as accounting lists |
| State location | Page owns `open` + success toast text | URL query `?create=1` | Matches appointments primary CTA; avoids extra routing |
| List body | Short intro or empty placeholder | Full user table | Out of scope; keep body minimal under toolbar |

## Data flow

```
AdminUsersPage
  ├─ useAuth → user, logout, roleOptions (useMemo)
  ├─ toolbar Button → setCreateOpen(true)
  └─ ProvisionUserModal(open)
        └─ onSubmit → apiClient.provisionUser → onSuccess → clear fields, onOpenChange(false), setMessage
```

## File changes

| File | Action | Description |
|------|--------|-------------|
| `frontend/src/app/admin/users/page.tsx` | Modify | Toolbar, modal state, render `ProvisionUserModal`; surface message/error outside or inside modal per UX consistency with accounting |
| `frontend/src/app/admin/users/ProvisionUserModal.tsx` | Create | Dialog + form; calls `provisionUser`; cancel/close |
| `frontend/src/app/admin/users/admin-users.module.css` | Modify | Toolbar layout; trim full-page form-only rules |
| `frontend/src/app/admin/users/page.test.tsx` | Create | RTL: mock `apiClient`, open CTA, `getByRole('dialog')`, submit success path (see `suppliers/page.test.tsx`) |

## Interfaces

- `ProvisionUserModalProps`: `{ open: boolean; onOpenChange: (v: boolean) => void; roleOptions: …; defaultRole: ProvisionRole; currentUserRole: UserRole; onProvisioned?: (email: string, role: string) => void }` — adjust to avoid prop drilling if context not needed.

## Testing strategy

| Layer | What | How |
|-------|------|-----|
| RTL | Toolbar → dialog → submit | `vi.mock('@/lib/api-client')`, `userEvent`, `waitFor` |
| Manual | Manager vs admin role options | Quick check in browser |

## Migration

No migration.

## Open questions

- [x] **Success feedback**: show **green message below the toolbar** in the main column (`styles.ok` on `page.tsx`), not only inside the modal — matches suppliers list pattern (status outside dialog).
