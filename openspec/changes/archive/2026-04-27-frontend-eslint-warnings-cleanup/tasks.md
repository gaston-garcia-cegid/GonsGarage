# Tasks: Frontend ESLint warnings cleanup (React Hooks)

## Nota de entrega (acordo de equipo)

- **Non abrir PRs**; traballar na rama acordada (**main** se é a rama activa).
- Acumular cambios localmente e facer **un único commit** ao pechar toda a fase 6 (código + `eslint.config.mjs`).
- **Subir a `error`**: `react-hooks/set-state-in-effect`, `react-hooks/immutability` e `react-hooks/purity` en `frontend/eslint.config.mjs` **no mesmo commit**, unha vez `pnpm lint` estea verde con regras en `error`.

## Phase 1: Baseline

- [x] 1.1 Executar `cd frontend && pnpm lint` e conservar mapa ficheiro→regra (referencia para verificar 0 problemas ao final).

## Phase 2: Contabilidade (patrón A do design)

- [x] 2.1 `frontend/src/app/accounting/billing-documents/page.tsx` e `[id]/page.tsx`: evitar `setState` síncrono antes do primeiro `await` nos loaders chamados dende `useEffect`.
- [x] 2.2 `frontend/src/app/accounting/issued-invoices/page.tsx` e `[id]/page.tsx`: igual.
- [x] 2.3 `frontend/src/app/accounting/received-invoices/page.tsx` e `[id]/page.tsx`: igual.
- [x] 2.4 `frontend/src/app/accounting/suppliers/page.tsx` e `[id]/page.tsx`: igual.

## Phase 3: Admin, citas, auth

- [x] 3.1 `frontend/src/app/admin/parts/page.tsx`, `[id]/page.tsx`, `components/PartCreateModal.tsx`: corrixir avisos `react-hooks/*` segundo design.
- [x] 3.2 `frontend/src/app/admin/users/ProvisionUserModal.tsx`: corrixir effect/setState.
- [x] 3.3 `frontend/src/app/appointments/page.tsx` e `components/AppointmentModal.tsx`: corrixir hooks afectados.
- [x] 3.4 `frontend/src/app/auth/login/LoginForm.tsx`: corrixir effect inicial.

## Phase 4: Cars, client, dashboard, employees, my-invoices

- [x] 4.1 `frontend/src/app/cars/[id]/page.tsx` e `cars/components/CarsContainer.tsx`: patrón A onde aplique; en `CarsContainer` corrixir `react-hooks/exhaustive-deps` (p.ex. dependencias de `useCallback` alineadas con `cars` / `cars.length` usados no corpo).
- [x] 4.2 `frontend/src/app/client/hooks/useClientData.ts`: patrón A ou B segundo fluxo.
- [x] 4.3 `frontend/src/app/dashboard/page.tsx`: corrixir effect.
- [x] 4.4 `frontend/src/app/employees/page.tsx`: corrixir effect.
- [x] 4.5 `frontend/src/app/my-invoices/MyInvoicesListClient.tsx` e `[id]/MyInvoiceDetailClient.tsx`: corrixir effects.

## Phase 5: Workshop, modais compartidos, tema, hidratación

- [x] 5.1 `frontend/src/app/workshop/page.tsx`, `[id]/page.tsx`, `recepcion/page.tsx`: corrixir todos os sitios sinalados polo lint.
- [x] 5.2 `frontend/src/components/appointments/NewAppointmentModal.tsx`: corrixir múltiples effects.
- [x] 5.3 `frontend/src/components/theme/ThemeSwitcher.tsx`: patrón B (mount / preferencia sen setState síncrono problemático no effect).
- [x] 5.4 `frontend/src/hooks/useAuthHydrationReady.ts`: patrón B (rama `hasHydrated()` sen setState síncrono no effect ou alternativa documentada).

## Phase 6: ESLint `error`, verificación, commit único

- [x] 6.1 `frontend/eslint.config.mjs`: poñer **`"error"`** en `react-hooks/set-state-in-effect`, `react-hooks/immutability`, `react-hooks/purity`.
- [x] 6.2 `cd frontend && pnpm lint && pnpm typecheck && pnpm test -- --passWithNoTests` — todo verde, ESLint **0 warnings e 0 errors**.
- [ ] 6.3 Smoke manual mínimo: accounting (unha lista + detalle), `/cars` con `?addCar=1` se aplica, login, toggle de tema (cumpre escenarios de calidade do delta **ui-brand-shell**). *(Pendente en browser por quen fusione.)*
- [x] 6.4 `git add` (incluír `openspec/changes/frontend-eslint-warnings-cleanup/*` se tocan na mesma sesión) e **un único `git commit`**; **non crear PR**.
