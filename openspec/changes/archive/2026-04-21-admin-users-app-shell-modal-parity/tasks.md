# Tasks: Admin Users — Shell & Modal Parity

## Phase 1: Tests first (Strict TDD — RED)

- [x] 1.1 Add `frontend/src/app/admin/users/page.test.tsx`: mock `@/stores` (`useAuth` admin or manager), `@/lib/api-client`; assert toolbar shows primary CTA (e.g. «Novo utilizador») and **no** `role('dialog')` until click — test **fails** on current inline-only page.
- [x] 1.2 Same file: after CTA click, expect `getByRole('dialog')` and heading inside dialog; test fails until modal exists.
- [x] 1.3 Same file: fill minimal valid fields, submit; mock `provisionUser` resolves success; expect dialog closes or success text visible per spec — fails until wired.

## Phase 2: Modal + page (GREEN)

- [x] 2.1 Create `frontend/src/app/admin/users/ProvisionUserModal.tsx`: `Dialog`/`DialogContent`/`DialogHeader`/`DialogTitle` from `@/components/ui/dialog`; move fields + submit from current `page.tsx`; props `open`, `onOpenChange`, `roleOptions`, `defaultRole`, `callerRole`; call `apiClient.provisionUser`; show inline error in modal on failure; `onOpenChange(false)` on success after optional `onProvisioned` callback.
- [x] 2.2 Refactor `frontend/src/app/admin/users/page.tsx`: pass `toolbar` to `AppShell` (fragment: `h1` + `Button` opening modal); state `createOpen`; remove inline form; keep `useMemo` role options on page, pass into modal; keep success `message` below toolbar or in page body (align with accounting empty-state pattern).
- [x] 2.3 Update `frontend/src/app/admin/users/admin-users.module.css`: toolbar row spacing; modal form field styles; remove unused full-page-only rules.

## Phase 3: Spec scenarios & edge cases (tests + polish)

- [x] 3.1 RTL: cancel/close dialog without submit — `provisionUser` **not** called; still on page (spec: close without success).
- [x] 3.2 RTL: `provisionUser` returns `success: false` — error text inside dialog; dialog stays open (spec: failed submit).
- [x] 3.3 Run `pnpm test` and `pnpm typecheck` in `frontend`; fix lint issues in touched files only.

## Phase 4: Cleanup

- [x] 4.1 Resolve design open question: pick one success surface (banner under toolbar vs text only) and document in `design.md` one line if behavior differs from draft.
- [x] 4.2 Manual smoke: manager vs admin role dropdown options in modal match current behavior.
