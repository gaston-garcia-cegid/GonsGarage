# Design: Next.js 16 + React 19 + Tailwind v4 + Zustand refactor

## Technical Approach

Execute on `main` in the order from the proposal, mapped to specs **ui-brand-shell** (stack + ADR + tabla “antes/después”) and **ui-component-system** (RSC/client, stores, React 19).

1. Bump **Next 16.x** + **`eslint-config-next`**; fix compile/runtime until `pnpm build` passes (repo already pins `turbopack.root` in `next.config.ts` for monorepo resolution — re-validate after bump).
2. **Tailwind v4**: migrate PostCSS/CSS entry (`globals.css`, `src/styles/*.css`) per official upgrade; align tokens with existing CSS variables (`--surface-page`, shadcn-style HSL vars).
3. **Zustand / auth**: eliminate parallel auth implementations (see decisions); keep store-only persistence keys (`gons-garage-auth`) and `merge`/`partialize` semantics unless intentionally changed (document in tabla).
4. **RSC**: route-by-route — move **initial** list/detail fetch to server only where the same public API contract can run on the server (env `NEXT_PUBLIC_API_URL` / origin helpers); keep **Zustand, Radix modals, router, `localStorage`** inside explicit `"use client"` islands.
5. **`use` / `useOptimistic`**: adopt only after (3)(4) stable; each use documented in this file or linked ADR (spec requirement).

## Architecture Decisions

| Decision | Alternatives | Tradeoff | Choice |
|----------|----------------|----------|--------|
| Dev/build bundler | Turbopack (default) vs `pnpm dev:webpack` | Rare plugin edge cases | Default **Turbopack**; keep existing **webpack** script as escape hatch until CI green. |
| Tailwind v4 shape | JS `tailwind.config` only vs CSS-first `@import "tailwindcss"` + `@theme` | v4 prefers CSS pipeline | **CSS-first v4**; shrink or adapt `tailwind.config.ts`; preserve `darkMode: selector` + `[data-theme="dark"]` behavior. |
| Auth sources of truth | Keep `AuthContext` + `auth.store` separate vs unify | Duplicate localStorage/login paths today | **Unify on `auth.store`**: `AuthProvider` becomes a thin `"use client"` wrapper (Context **or** direct re-exports) so `LoginForm` / layouts do not maintain a second auth state machine. |
| Store ↔ RSC | Import stores in server files | Build error / wrong runtime | **Stores only in client modules**; server passes **serializable props** where needed. |
| React `use()` | Wide adoption vs targeted | Suspense boundaries + client-only | **Targeted** client islands; document each adoption vs `useEffect`. |

## Data Flow

```
Request → app/layout.tsx (server) → fonts, globals.css
              ↓
        client: ThemeSwitcher, AuthProvider (thin) → Zustand rehydrate
              ↓
        page: RSC data (optional) + client children → stores / Radix / router
```

**Auth hydration**: preserve ordering so routes that read `isAuthenticated` do not flash wrong state (existing `useAuthHydrationReady` pattern or equivalent after refactor).

## File Changes

| Path | Action | Description |
|------|--------|-------------|
| `frontend/package.json`, `pnpm-lock.yaml` | Modify | Next 16, `eslint-config-next`, Tailwind 4, compatible tooling. |
| `frontend/postcss.config.mjs` | Modify | Tailwind v4 PostCSS plugin setup. |
| `frontend/tailwind.config.ts`, `frontend/src/app/globals.css`, `frontend/src/styles/**` | Modify | v4 pipeline; token/theme parity. |
| `frontend/next.config.ts` | Modify | Verify `output: "standalone"` + `turbopack.root` under Next 16. |
| `frontend/src/contexts/AuthContext.tsx`, auth consumers | Modify | Single auth path via store. |
| `frontend/src/stores/**` | Modify | Structure + boundaries; same observable UX unless tabla row says otherwise. |
| `frontend/src/app/**` | Modify | `"use client"` boundaries; incremental server data. |
| `docs/**` (ADR) + this file | Modify | **GO**, commands, outcomes, **before/after** rows for structural UI. |
| `.github/workflows/ci.yml` | Maybe | Node 22 already; adjust only if toolchain requires. |

## Interfaces / Contracts

- Public entry: `frontend/src/stores/index.ts` — preserve `useAuth`, `useCars`, `useAppointments`, etc., or document renames in the tabla.
- **`apiClient` token** must follow the same source as the store after unification (no second writer from Context).

## Testing Strategy

| Layer | What | How |
|-------|------|-----|
| Unit | `auth.store`, role helpers | Vitest (extend `auth.client-role.test.ts` patterns). |
| Integration | Login shell, key list pages | RTL + API mocks / MSW if already used. |
| E2E | — | Out of scope (manual checklist per proposal). |

## Migration / Rollout

Atomic commits on `main`; rollback: `git revert` commit range + ADR **ROLLBACK** with SHA (proposal).

## Before / after (structural UI — maintain during implementation)

| Área | Antes | Después | Motivo (testable) |
|------|--------|---------|-------------------|
| Next / ESLint | 15.5.x | 16.x + aligned `eslint-config-next` | Spec: official stack on `main`. |
| Tailwind | v3 + `@tailwind` directives | v4 CSS pipeline (`@import "tailwindcss"`, `@tailwindcss/postcss`) | Spec: v4 mandatory; evitar regresión visual documentada. |
| Animación utilidades | `tailwindcss-animate` (plugin v3) | **`tw-animate-css`** (`@import "tw-animate-css"` en `globals.css`) | Plugin legacy no compatible con pipeline v4; paquete oficial recomendado por shadcn/TW4. |
| Auth / `apiClient` | Context + store duplicaban sesión y token | **Store único** + `AuthProvider` que sólo chama `checkAuthStatus`; `useAuth` desde `@/stores`; `@/lib/api` **delega** `setToken`/`clearToken` ao singleton `@/lib/api-client` | Spec: unha fonte de verdade; token HTTP alineado para código legacy que aínda importa `@/lib/api`. |
| Datos iniciales | Muchos `useEffect` de fetch en páginas cliente | Server fetch donde contrato igual | Spec: sin `useEffect` solo-datos cuando el servidor basta. |
| RSC piloto facturas cliente | Lista/detalle `my-invoices` todo cliente | **`page.tsx` servidor** + `lib/server/my-invoices-initial.ts` (Bearer vía cookie cando exista) + illas `MyInvoices*Client.tsx` | JWT hoxe en `localStorage` → servidor devolve baleiro e o cliente fai un `listMine`/`get`; ADR 0003 §4.4 inventario do resto. |
| Notas factura (cliente) | `onSubmit` + `startTransition` manual | **`<form action={async (fd) => …}>`** + `useOptimistic` sobre `row` (`MyInvoiceDetailClient`) | React 19: `mergeOptimistic` debe correr nunha **acción / host transition**; `startTransition` desde `onSubmit` non enlazaba o entangled lane nos tests e o optimistic non pintaba; `action` usa `startHostTransition` e os tests RTL pasan. Erro: `setRow(snapshot)` + mensaxe API. |
| `use()` vs `useEffect` | — | **Sen `use()`** neste fluxo | Non hai unha **Promise estable** que suba por props/context para `use()`; carga inicial segue `initialRow` + `useEffect`/`load()`; ADR 0003 §Phase 5. |
| Verificación / GO | Criterio disperso | **Matriz 6.1** + checklist **6.2** no ADR 0003 §Phase 6; decisión **GO** | `pnpm lint|typecheck|test --passWithNoTests|build` alineados con `ci.yml`; smoke humano pre-release. |
| CI | — | **Sin cambios** en `ci.yml` | Node 22 + mesmos pasos; ADR 0003 §6.4. |

## Open Questions

- [x] **`tailwindcss-animate`**: **Reemplazado por `tw-animate-css`** (devDependency + import en `globals.css`); `tailwindcss-animate` eliminado del árbol de dependencias directas.
- [ ] **Radix + Next 16 + TW4**: first full `pnpm build` may surface peer warnings — resolve or pin per ADR row.
