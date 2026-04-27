# Design: Frontend ESLint warnings cleanup (React Hooks)

## Technical approach

Cumprir o delta **ui-brand-shell** (`pnpm lint` **0 warnings**) e os criterios da proposta, sen mudar contratos API nin rutas. Estado actual: **36 warnings** (0 errors), maioría **`react-hooks/set-state-in-effect`**; un **`react-hooks/exhaustive-deps`** en `CarsContainer.tsx`; regras **`immutability`** / **`purity`** seguen en `warn` en `eslint.config.mjs` pero non aparecen no lote actual.

1. **Inventario fixo**: o lint xa lista ~**24 ficheiros** (accounting, admin/parts, appointments, auth/login, cars, client, dashboard, employees, my-invoices, workshop, modais compartidos, `ThemeSwitcher`, `useAuthHydrationReady`).
2. **Patrón A — loaders `useCallback` + `useEffect(() => void load())`**: `load` chama **`setError(null)` (ou outro `setState`) antes do primeiro `await`**; o analizador trata iso como setState síncrono disparado dende o effect. **Corrección**: mover o reset de erro para **despois** do primeiro yield (`await Promise.resolve()`, inicio do `try` trala primeira await, ou só establecer erro cando falle a petición); ou **inlining** do fetch no `useEffect` cunha IIFE `async` que non chame `setState` síncrono antes de await. Manter cancelación (`let cancelled = false`) onde xa exista ou engadir en páxinas de detalle.
3. **Patrón B — UI de mount** (`ThemeSwitcher`, `useAuthHydrationReady`): evitar ramas que fagan **`setState` síncrono** no corpo do effect cando `hasHydrated()` é true; preferir **subscrición só** (`onFinishHydration`) + estado inicial derivado de `hasHydrated()` no `useState` inicial (xa existe no hook), ou **diferir** un tick (`queueMicrotask`) **documentado** na PR se React non ofrece API máis limpa sen duplicar lóxica.
4. **Patrón C — query params / modais** (`setCreateOpen`, `setShowCreateModal`, etc.): onde o effect só reacciona a `searchParams`/`router`, valorar **`startTransition`** para updates non urxentes ou reestruturar para que o primeiro `setState` non sexa síncrono respecto ao commit (seguir guía do aviso ESLint).
5. **Fase final**: con árbore a **0 warnings**, commit pequeno que pase **`set-state-in-effect`** (e opcionalmente **`immutability`** / **`purity`**) de `"warn"` a `"error"` en `frontend/eslint.config.mjs` **só se** `pnpm lint` segue verde.

## Architecture decisions

| Decision | Alternatives | Tradeoff | Choice |
|----------|----------------|----------|--------|
| Silenciar vs corrixir | `eslint-disable-next-line` masivo | Incumpre spec de waivers (≤3 ficheiros, acotado) | **Corrixir código**; disable só excepcional con motivo. |
| Fetch dende effect | React Query / Server Components | Alcance fóra da proposta | Manter **effects + servizos existentes**; só cambiar **orde/timing** de `setState`. |
| Theme / hydration | `useSyncExternalStore` vs microtask | máis API surface | Preferir **APIs de subscribe** (Zustand `onFinishHydration`, `matchMedia`) + estado inicial correcto; **microtask** só documentado. |
| Severidade ESLint | Manter `warn` para sempre | Non cumpre gate 0 warnings | **0 warnings** primeiro; **`error`** nun commit final opcional. |

## Data flow

```
useEffect([deps])
    → chama load() ou lóxica sync
         → ❌ setState antes de await  → ESLint set-state-in-effect
    → corrección: await primeiro tick OU setState só en .then / despois de await
         → ✅ mesmo resultado de datos, menos renders en cadea
```

## File changes

| Path | Action | Description |
|------|--------|-------------|
| `frontend/src/app/accounting/**/page.tsx`, `[id]/page.tsx` | Modify | Patrón A en listas e detalles (billing-documents, issued/received-invoices, suppliers). |
| `frontend/src/app/admin/parts/**` | Modify | PartCreateModal + páxinas list/detail. |
| `frontend/src/app/admin/users/ProvisionUserModal.tsx` | Modify | Effect + setState. |
| `frontend/src/app/appointments/**` | Modify | `AppointmentModal`, `page.tsx`. |
| `frontend/src/app/auth/login/LoginForm.tsx` | Modify | Effect inicial. |
| `frontend/src/app/cars/**` | Modify | `[id]/page`, `CarsContainer` (incl. deps `useCallback`). |
| `frontend/src/app/client/hooks/useClientData.ts` | Modify | Patrón A/B según fluxo. |
| `frontend/src/app/dashboard/page.tsx` | Modify | Effect. |
| `frontend/src/app/employees/page.tsx` | Modify | Effect. |
| `frontend/src/app/my-invoices/**` | Modify | Clients de lista/detalhe. |
| `frontend/src/app/workshop/**` | Modify | Várias liñas en `page.tsx` e recepción. |
| `frontend/src/components/appointments/NewAppointmentModal.tsx` | Modify | Múltiples effects. |
| `frontend/src/components/theme/ThemeSwitcher.tsx` | Modify | Patrón B. |
| `frontend/src/hooks/useAuthHydrationReady.ts` | Modify | Ramo sync `setReady(true)`. |
| `frontend/eslint.config.mjs` | Modify | Opcional: `warn` → `error` en regras `react-hooks/*` ao pechar. |

Sen ficheiros novos obrigatorios; **tests** só se o comportamento de mount/fetch cambie de xeito observable.

## Interfaces / contracts

Ningún cambio de contrato HTTP ou tipos públicos; só orde de execución e dependencias de hooks.

## Testing strategy

| Layer | What | How |
|-------|------|-----|
| Unit / RTL | Páxinas e hooks xa cubertos | `pnpm test -- --passWithNoTests` + executar suites que toquen ficheiros modificados. |
| Manual | Rotas accounting, cars `?addCar=1`, login, theme toggle | Checklist curto (ADR §6.2 onde aplique). |
| E2E | — | Fóra de alcance. |

## Migration / rollout

**No migration required** (sen datos nin feature flags). Rollout: PR(s) por fase (accounting → resto → ESLint severity).

## Open questions

- [ ] **Un PR ou varios**: equipo prefire un PR monolítico ou fases mergeables (proposta sugire fases).
- [ ] **`immutability` / `purity`**: permanecen en `warn` ata que reporten algo, ou subir a `error` xunto a `set-state-in-effect`?
