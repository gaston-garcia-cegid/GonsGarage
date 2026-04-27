# Tasks: Next.js 16 + React 19 + Tailwind v4 + Zustand refactor

## Phase 1: Foundation — Next 16 + toolchain

- [x] 1.1 Add ADR under `docs/`: scope, GO/NO-GO, commands (`pnpm lint|typecheck|test|build`).
- [x] 1.2 Bump `next` + `eslint-config-next` to 16.x in `frontend/package.json`; `pnpm install`; lockfile updated.
- [x] 1.3 `cd frontend && pnpm typecheck && pnpm build` until green.
- [x] 1.4 Verify `frontend/next.config.ts` (`turbopack.root`, `output: "standalone"`); if Turbopack fails, try webpack path and note in ADR.

## Phase 2: Tailwind v4 + CSS

- [x] 2.1 `frontend/postcss.config.mjs` → Tailwind v4 PostCSS per upgrade guide.
- [x] 2.2 `frontend/src/app/globals.css` → v4 imports; preserve `src/styles/tokens.css`, `shadcn-theme.css`, `utilities.css` order.
- [x] 2.3 `frontend/tailwind.config.ts` + `frontend/src/styles/**`: dark `[data-theme="dark"]`, HSL tokens parity vs pre-change.
- [x] 2.4 `tailwindcss-animate` vs v4: upgrade, replace, or remove — row in `design.md` tabla.
- [x] 2.5 `pnpm lint` + `pnpm build`; fix utilities/`@apply` breaks.

## Phase 3: Zustand + auth

- [x] 3.1 RED: extend `frontend/src/stores/auth.client-role.test.ts` for unified auth + `apiClient` token (fails until refactor lands).
- [x] 3.2 Thin `frontend/src/contexts/AuthContext.tsx` delegating to `auth.store.ts` (no duplicate login/storage).
- [x] 3.3 Point `useAuth` imports to `@/stores`: `app/page.tsx`, `employees/page.tsx`, `auth/login/LoginForm.tsx`, `components/layouts/DashboardLayout.tsx`, `employees/page.test.tsx` mock.
- [x] 3.4 Grep `setToken`/`clearToken` in `frontend/src`; single sync from store.
- [x] 3.5 GREEN: `pnpm test` auth-related tests pass.
- [x] 3.6 `car.store.ts`, `appointment.store.ts`, `stores/index.ts`: client-only boundaries; document export changes in `design.md` if any.

## Phase 4: RSC / client islands

- [x] 4.1 Audit `frontend/src/app/**`: default server pages; `"use client"` only where store/Radix/router need it.
- [x] 4.2 Pilot one list route: server initial fetch + props; drop redundant data-only `useEffect` there.
- [x] 4.3 Pilot one detail route if server `fetch` matches client contract (`NEXT_PUBLIC_API_URL`).
- [x] 4.4 ADR: remaining data-only `useEffect` with “not candidate” reason each.

## Phase 5: React 19

- [x] 5.1 One mutation: `useOptimistic` + reconcile/error; ADR/design row.
- [x] 5.2 Optional: one `use()` in client; document vs `useEffect`.

## Phase 6: Verify + GO

- [x] 6.1 Local CI matrix in `frontend`: lint, typecheck, test `--passWithNoTests`, build.
- [x] 6.2 Manual smoke: auth, workshop, parts, accounting.
- [x] 6.3 ADR final: outcomes, **GO**, rollback SHA; refresh `design.md` before/after table.
- [x] 6.4 `.github/workflows/ci.yml` only if Node/steps must change; else ADR “CI unchanged”.
