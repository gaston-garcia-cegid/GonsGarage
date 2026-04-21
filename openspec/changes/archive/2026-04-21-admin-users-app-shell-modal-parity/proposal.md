# Proposal: Admin Users — Shell & Modal Parity

## Intent

Match `/admin/users` to **Marcações** (`/appointments`): **AppShell** + **toolbar** (title + primary CTA) + **modal** for **«Criar utilizador»**, replacing the lone inline form so the tab matches other staff areas and accounting toolbar→modal create flows.

## Scope

### In Scope

- Toolbar (heading + button opening create flow); **Dialog** with existing fields; same **`apiClient.provisionUser`** and validation/error UX.
- Optional empty/success copy; spacing/typography aligned with appointments patterns where cheap.
- **Vitest RTL**: toolbar → modal → submit → mocked provision (accounting `*page.test.tsx` style).

### Out of Scope

- Backend/auth changes; full user list; role matrix; i18n overhaul.

## Capabilities

Specs skimmed: `staff-user-management-ui`, `mvp-role-access`, `staff-user-provisioning`. Toolbar vs inline is **not** a new normative requirement.

### New Capabilities

- None

### Modified Capabilities

- None

## Approach

- Copy **`frontend/src/app/appointments/page.tsx`**: `AppShell`, toolbar (`h1` + `Button`), modal state, list/empty body using shared list styles where fit.
- **ProvisionUserModal** (name flexible): extract current form; page holds `open` / `onSuccess` close.
- Reuse shared **Button** / **Input** / dialog; **ConfirmModal** only if destructive (unlikely).
- Tidy **`admin-users.module.css`**; prefer appointments **`styles.list`** spacing over full-page form layout.
- RTL beside admin users page.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `frontend/src/app/admin/users/page.tsx` | Modified | Toolbar + modal wiring |
| `frontend/src/app/admin/users/*.tsx` (new) | New | Modal + form |
| `frontend/src/app/admin/users/admin-users.module.css` | Modified | Toolbar/shell layout |
| `frontend/src/app/appointments/page.tsx` | Reference | Pattern only |
| `frontend/src/**/*.test.tsx` | New/Modified | Modal + mock submit |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Modal a11y | Low | Same Dialog as other modals |
| Visual drift vs Marcações | Med | Reuse toolbar/button variants |

## Rollback Plan

Revert frontend commits; inline form returns. No migrations.

## Dependencies

- None

## Success Criteria

- [ ] Toolbar + CTA; form **only** in modal.
- [ ] **`apiClient.provisionUser`** unchanged contract.
- [ ] RTL: open from toolbar, submit with mock.
- [ ] No primary lone inline form under `AppShell`.
