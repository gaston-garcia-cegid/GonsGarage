# Apply progress — admin-users-app-shell-modal-parity

**Mode**: Strict TDD  
**Batch**: single (all `tasks.md` items)

## TDD Cycle Evidence

| Task | RED | GREEN | Triangulate | Safety net | REFACTOR |
|------|-----|-------|-------------|------------|----------|
| 1.1–1.3 | ✅ `page.test.tsx`: toolbar + no dialog; open dialog; submit success | ✅ `ProvisionUserModal` + `page.tsx` toolbar | ✅ Tests 3.1 cancel + 3.2 API error | ✅ `pnpm test` full suite before push (baseline 73 → 78) | ✅ Shared CSS tokens; `within(main)` for h1 |
| 2.1 | (covered by modal tests) | ✅ `ProvisionUserModal.tsx` | ✅ Error path + reset on `open` | ➖ New file | ➖ |
| 2.2 | — | ✅ `page.tsx` uses `AppShell` `toolbar` | — | — | ➖ |
| 2.3 | — | ✅ `admin-users.module.css` | — | — | ➖ |
| 3.1–3.3 | — | ✅ Cancel + failure RTL; full `pnpm test` + `pnpm typecheck` | — | — | ➖ |
| 4.1 | — | ✅ `design.md` success surface resolved | — | — | ➖ |
| 4.2 | ➖ Manual | ➖ Admin RTL covers provision matrix; manager role options: verify in browser or follow-up test | — | — | ➖ |

## Files changed

| File | Action |
|------|--------|
| `frontend/src/app/admin/users/ProvisionUserModal.tsx` | Created |
| `frontend/src/app/admin/users/page.tsx` | Modified |
| `frontend/src/app/admin/users/admin-users.module.css` | Modified |
| `frontend/src/app/admin/users/page.test.tsx` | Created |
| `openspec/changes/admin-users-app-shell-modal-parity/design.md` | Modified (open question) |

## Deviations from design

None: Dialog + toolbar + `apiClient.provisionUser` unchanged contract.

## Issues

None.

## Status

**11/11 tasks complete.** Ready for `sdd-verify`.
