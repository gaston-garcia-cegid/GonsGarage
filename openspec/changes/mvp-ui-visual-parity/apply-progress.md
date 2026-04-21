# Apply progress — mvp-ui-visual-parity

**Mode:** Strict TDD (`openspec/config.yaml`).

## TDD cycle evidence

### Phase 1

| Task | Layer | RED | GREEN | Refactor |
|------|-------|-----|--------|----------|
| 1.1 | Doc | N/A | ✅ `docs/ui-audit-mvp-dark.md` | — |
| 1.2 | CSS tokens | N/A | ✅ `tokens.css` dark + loading vars | — |
| 1.3 | CSS | N/A | ✅ `utilities.css` `.app-loading` | — |
| 1.4 | React | ✅ `AppLoading.test.tsx` | ✅ `AppLoading.tsx` | — |
| 1.5 | Pages | N/A | ✅ dashboard, cars, car `[id]` | — |

### Phase 2

| Task | Layer | RED | GREEN | Refactor |
|------|-------|-----|--------|----------|
| 2.1 | employees | N/A (CSS inline → tokens) | ✅ hex → `var(--chip-*\|--color-error\|--surface-*)`, loading `AppLoading` | — |
| 2.2 | client shell | N/A | ✅ `DashboardLayout.module.css` + `client/page.tsx` loading | — |
| 2.3 | accounting | N/A | ✅ `layout.tsx` + 4× `[id]/page` `AppLoading`; `accounting.module.css` já tokens | — |
| 2.4 | dashboard CSS | N/A | ✅ `.spinner` removido de `dashboard.module.css` e `car-details.module.css` | — |
| 2.5 | inventory | N/A | ✅ `docs/ui-loader-exceptions.md`; `rg` 0 em `app/` | — |
| 2.6 | smoke | N/A | ✅ notas em `ui-audit-mvp-dark.md` | — |

### Phase 3 (spike Next 16 + Tailwind v4)

| Task | Layer | RED | GREEN | Refactor |
|------|-------|-----|--------|----------|
| 3.1 | ADR | N/A | ✅ `docs/adr/0001-next16-tailwind4-spike.md` | — |
| 3.2 | deps + CSS pipeline | N/A | ✅ rama `spike/next16-tailwind4`: `next@16.2.4`, `eslint-config-next@16.2.4`, `tailwindcss` + `@tailwindcss/postcss`, `postcss.config.mjs`, `@import "tailwindcss"` en `globals.css` | — |
| 3.3 | CI local | N/A | ⚠️ lint **falla** (26 errores hooks); typecheck / test / build **OK** — volcado en ADR | — |
| 3.4 | decisión | N/A | ✅ ADR: **NO-GO** merge a `main` hasta lint verde; pasos si GO | — |

## Files changed (cumulative)

**Fase 1:** ver histórico git.  
**Fase 2:** `employees/page.tsx`, `client/page.tsx`, `components/layouts/DashboardLayout.module.css`, `accounting/layout.tsx`, `accounting/**/[id]/page.tsx` (4), `my-invoices/layout.tsx`, `my-invoices/[id]/page.tsx`, `appointments/page.tsx`, `appointments/new/page.tsx`, `app/page.tsx`, `dashboard.module.css`, `car-details.module.css`, `docs/ui-loader-exceptions.md`, `docs/ui-audit-mvp-dark.md`.

**Estado:** Fases 1–3 completas (**15/26** tarefas). Fase 3: spike en rama `spike/next16-tailwind4`, **NO-GO** merge por ESLint. Siguiente: Fase 4 (shadcn) o PR de remediação de lint antes de subir Next en `main`.
