# Lint baseline (Phase 1)

**When:** 2026-04-27T19:15:46.468Z
**Command:** `pnpm exec eslint . --format json` (cwd: `frontend/`)

**Totals:** 36 warnings, 0 errors

**Rexenerar mapa** (desde `frontend/`): `pnpm exec eslint . --format json --output-file ../openspec/changes/frontend-eslint-warnings-cleanup/lint-baseline.json` → `node ../openspec/changes/frontend-eslint-warnings-cleanup/scripts/build-lint-baseline-map.mjs` → podes borrar o `.json` despois se non o queres versionar.

## Map: file → rule (line)

### `frontend/src/app/accounting/billing-documents/[id]/page.tsx`
- react-hooks/set-state-in-effect @54

### `frontend/src/app/accounting/billing-documents/page.tsx`
- react-hooks/set-state-in-effect @44

### `frontend/src/app/accounting/issued-invoices/[id]/page.tsx`
- react-hooks/set-state-in-effect @41

### `frontend/src/app/accounting/issued-invoices/page.tsx`
- react-hooks/set-state-in-effect @37

### `frontend/src/app/accounting/received-invoices/[id]/page.tsx`
- react-hooks/set-state-in-effect @47

### `frontend/src/app/accounting/received-invoices/page.tsx`
- react-hooks/set-state-in-effect @37

### `frontend/src/app/accounting/suppliers/[id]/page.tsx`
- react-hooks/set-state-in-effect @47

### `frontend/src/app/accounting/suppliers/page.tsx`
- react-hooks/set-state-in-effect @37

### `frontend/src/app/admin/parts/[id]/page.tsx`
- react-hooks/set-state-in-effect @56

### `frontend/src/app/admin/parts/components/PartCreateModal.tsx`
- react-hooks/set-state-in-effect @47

### `frontend/src/app/admin/parts/page.tsx`
- react-hooks/set-state-in-effect @41

### `frontend/src/app/admin/users/ProvisionUserModal.tsx`
- react-hooks/set-state-in-effect @50

### `frontend/src/app/appointments/components/AppointmentModal.tsx`
- react-hooks/set-state-in-effect @61

### `frontend/src/app/appointments/page.tsx`
- react-hooks/set-state-in-effect @69

### `frontend/src/app/auth/login/LoginForm.tsx`
- react-hooks/set-state-in-effect @31

### `frontend/src/app/cars/[id]/page.tsx`
- react-hooks/set-state-in-effect @74

### `frontend/src/app/cars/components/CarsContainer.tsx`
- react-hooks/exhaustive-deps @110
- react-hooks/set-state-in-effect @206

### `frontend/src/app/client/hooks/useClientData.ts`
- react-hooks/set-state-in-effect @25

### `frontend/src/app/dashboard/page.tsx`
- react-hooks/set-state-in-effect @68

### `frontend/src/app/employees/page.tsx`
- react-hooks/set-state-in-effect @55

### `frontend/src/app/my-invoices/[id]/MyInvoiceDetailClient.tsx`
- react-hooks/set-state-in-effect @49

### `frontend/src/app/my-invoices/MyInvoicesListClient.tsx`
- react-hooks/set-state-in-effect @30

### `frontend/src/app/workshop/[id]/page.tsx`
- react-hooks/set-state-in-effect @73

### `frontend/src/app/workshop/page.tsx`
- react-hooks/set-state-in-effect @106
- react-hooks/set-state-in-effect @76
- react-hooks/set-state-in-effect @83

### `frontend/src/app/workshop/recepcion/page.tsx`
- react-hooks/set-state-in-effect @46
- react-hooks/set-state-in-effect @51
- react-hooks/set-state-in-effect @72
- react-hooks/set-state-in-effect @93

### `frontend/src/components/appointments/NewAppointmentModal.tsx`
- react-hooks/set-state-in-effect @67
- react-hooks/set-state-in-effect @78
- react-hooks/set-state-in-effect @85

### `frontend/src/components/theme/ThemeSwitcher.tsx`
- react-hooks/set-state-in-effect @12

### `frontend/src/hooks/useAuthHydrationReady.ts`
- react-hooks/set-state-in-effect @18

