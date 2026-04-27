# ADR 0003 — Migración Next.js 16 + React 19 + Tailwind v4 (trabajo en `main`)

**Estado:** **GO** (cierre change SDD `nextjs-16-react19-migration`, Fase 6).  
**Fecha inicio:** 2026-04-27  
**Fecha cierre (verificación local):** 2026-04-27  
**Change SDD:** `openspec/changes/nextjs-16-react19-migration/` (proposal, specs, design, tasks).

**Relación:** amplía el [ADR 0001](./0001-next16-tailwind4-spike.md) (spike en rama, **NO-GO** por lint) con adopción real en `main` según `tasks.md`.

## Scope

- Subir **Next.js 16.x** + **`eslint-config-next`** alineado; mantener **React 19** ya declarado en `frontend/package.json`.
- **Tailwind CSS v4**, pipeline PostCSS/CSS, paridad de tokens y tema (Fase 2 del change).
- **Refactor Zustand / auth** una sola fuente de verdad; límites RSC/client (Fases 3–4).
- **React 19** `use` / `useOptimistic` donde se documente (Fase 5).
- Fuera de alcance: ver proposal (`openspec/changes/nextjs-16-react19-migration/proposal.md`).

## Comandos de calidad (barra oficial)

Ejecutar desde `frontend/`:

- `pnpm lint`
- `pnpm typecheck`
- `pnpm test -- --passWithNoTests`
- `pnpm build` (con `NEXT_PUBLIC_API_URL` si aplica en CI, ver `.github/workflows/ci.yml`)

## GO / NO-GO / DEFER

| Decisión | Criterio breve |
|----------|----------------|
| **GO** | Los cuatro comandos anteriores verdes en `main`; ADR actualizado con resultados; `package.json` refleja Next 16 + Tailwind v4 al cierre del change (ver specs delta `ui-brand-shell`). |
| **NO-GO** | Regresión bloqueante o CI rojo sin plan de revert acotado. |
| **DEFER** | Posponer solo parte documentada (p. ej. `use()` opcional) con fila en `design.md` del change. |

## Turbopack / webpack

- `frontend/next.config.ts` define `turbopack.root` al directorio del app (evita raíz incorrecta en monorepo).
- Si **Turbopack** falla en build/dev, probar `pnpm dev:webpack` o build con webpack según docs Next 16 y anotar aquí la evidencia.

## Evidencia (Phase 1 — toolchain)

- `frontend/package.json`: **Next 16.2.4**, **eslint-config-next 16.2.4** (2026-04-27).
- `pnpm typecheck`, `pnpm build` (Turbopack), `pnpm test`: verdes en Windows tras el bump.
- **ESLint:** `eslint.config.mjs` usa flat config nativo `eslint-config-next/core-web-vitals` (sin `FlatCompat`, evita error circular). Reglas nuevas `react-hooks/*` (set-state-in-effect, immutability, purity) en **`warn`** hasta refactors del change (Fase 4); correcciones puntuales: `AuthContext` (`checkAuthStatus` antes del effect), `CarModal` / `useCarValidation`, `dashboard` (`Date.now` en render).
- Cierre do change: ver **Phase 6** e `openspec/changes/nextjs-16-react19-migration/tasks.md` (todas as fases `[x]`).

## Phase 4 — RSC / client islands (2026-04-27)

### 4.1 Auditoría `frontend/src/app/**/page.tsx` (`use client`)

| Área | Patrón | Motivo típico de `"use client"` |
|------|--------|----------------------------------|
| Auth / Zustand | `page.tsx`, `dashboard`, `cars`, `workshop`, accounting lists, `employees`, `appointments`, `client` | `useAuth`, stores, `useRouter`, efectos de datos ou modais |
| Auth marketing | `app/page.tsx` (landing) | `useRouter`, redirect cando xa autenticado |
| Auth forms | `auth/register/page.tsx` | estado de formulario, `useAuth`, `useRouter` |
| Auth login | `auth/login/page.tsx` | **Sen** `"use client"` na ruta: **Server Component** que envolve `LoginForm` (illa cliente) |
| Detalle / CRUD pesado | `cars/[id]`, `workshop/[id]`, `admin/parts/[id]`, rutas accounting `.../[id]/page.tsx` | fetch + formularios + estado local |
| **As miñas facturas** | `my-invoices/page.tsx`, `my-invoices/[id]/page.tsx` | **Servidor** por defecto: fetch inicial opcional (`cookies` + Bearer); UI en `MyInvoices*Client.tsx` |

Regra aplicada: **default server** na ruta cando o layout non forza todo o subtree a cliente; as illas (`'use client'`) quedan en compoñentes fillos (`*Client.tsx`, `LoginForm`, modais, stores).

### 4.2 / 4.3 Piloto list + detalle (`/my-invoices`)

- **Lista**: `page.tsx` chama `fetchMyInvoicesInitialAuthenticated()` (`lib/server/my-invoices-initial.ts`, mesma orixe que `getPublicApiOrigin()`). Se existise cookie `auth_token` http, o HTML inicial podería levar filas sen `useEffect` de carga.
- **Hoxe**: a sesión JWT está en `localStorage` → o servidor devolve `[]` e a illa **`MyInvoicesListClient`** fai `listMine()` unha vez (comportamento equivalente ao anterior).
- **Detalle**: `my-invoices/[id]/page.tsx` servidor + `fetchMyInvoiceDetailAuthenticated(id)`; **`MyInvoiceDetailClient`** omite o primeiro GET se `initialRow` existe.

### 4.4 `useEffect` só-datos — inventario (non candidato a RSC nesta fase)

| Ficheiro (app) | Uso | Por que **non** candidato agora |
|----------------|-----|----------------------------------|
| `page.tsx` (landing) | redirect se autenticado | Depende de `useAuth` + hidratación cliente |
| `dashboard/page.tsx` | varios | Zustand + datos iniciais + auth |
| `cars/page.tsx`, `cars/[id]/page.tsx` | carga listas/detalles | Bearer só no cliente |
| `workshop/*.tsx` | listas / visita | idem + estado complexo |
| `appointments/page.tsx` | filtros + lista | store + interacción |
| `employees/page.tsx` | lista empregados | `apiClient` legacy + modais |
| `accounting/**/page.tsx` (listas e `[id]`) | fetch listas/detalles | JWT cliente; modais Radix |
| `admin/parts/**` | lista, detalle, modais | idem |
| `admin/users` + modais | datos + roles | idem |
| `my-invoices/layout.tsx` (migrado a gate) | redirect rol cliente | precisa `useRouter` + store |
| `auth/register/page.tsx` | redirect / mensaxes | estado formulario |
| `auth/login/LoginForm.tsx` | query `message` | só lectura de URL en cliente |
| `accounting/layout.tsx`, `workshop/layout.tsx`, `admin/*layout.tsx` | guards | auth + router |

**Candidatura futura:** mover fetch inicial a RSC cando exista **cookie httpOnly** ou **BFF** que reenvíe o Bearer ao Gin co mesmo contrato JSON.

## Phase 5 — React 19 (`useOptimistic`, `use()`) (2026-04-27)

### 5.1 Mutación optimista (detalle «As miñas facturas»)

- Ficheiro: `frontend/src/app/my-invoices/[id]/MyInvoiceDetailClient.tsx`.
- **Patrón:** `useOptimistic(row, (current, nextNotes) => ({ ...current, notes: nextNotes }))` + `<form action={async (formData) => { … }}>` que chama `mergeOptimisticNotes` **antes** do `await` a `issuedInvoiceService.patchNotes`.
- **Por que `action` e non `onSubmit` + `startTransition`:** co segundo patrón, nos tests (Vitest + RTL) o merge encolábase sen reflectirse no DOM durante o PATCH diferido; React 19 encadea a acción do formulario con **`startHostTransition`**, que mantiene o contexto correcto para updates optimistas (e alíña coa guía de «async transition / action»).
- **Reconciliación / erro:** éxito → `setRow`/`setNotes` coa resposta; erro → `setError` + `setRow(snapshot)` (snapshot capturado ao inicio da acción). `data-testid="invoice-notes-optimistic"` (sr-only) para asercións.
- Probas: `frontend/src/app/my-invoices/[id]/MyInvoiceDetailClient.test.tsx` (erro diferido + éxito).

### 5.2 `use()` (opcional)

- **Non adoptado** neste piloto: o fluxo usa props `initialRow` + servizo imperativo (`get` / `patchNotes`); non hai unha Promise estable exposta como recurso de lectura onde `use()` mellore claramente fronte a `useEffect` + estado.
- Cando exista lectura suspendible (p. ex. cache `react`/`use` sobre unha promise memoizada), documentar de novo na táboa do `design.md`.

## Phase 6 — Verificación + GO (2026-04-27)

### 6.1 Matriz local `frontend/` (alineada con `.github/workflows/ci.yml`)

| Paso | Comando | Resultado |
|------|---------|-------------|
| Lint | `pnpm lint` | **OK** (exit 0); **0 erros**, **36 warnings** (`react-hooks/set-state-in-effect` e similares en rutas accounting, `ThemeSwitcher`, `useAuthHydrationReady`, etc. — mesma política **warn** acordada na Fase 1). |
| Typecheck | `pnpm typecheck` | **OK** |
| Test | `pnpm test -- --passWithNoTests` | **OK** — 25 ficheiros, **111** tests |
| Build | `NEXT_PUBLIC_API_URL=http://localhost:8080 pnpm build` | **OK** (Next **16.2.4**, Turbopack); aviso Node sobre `tailwind.config.ts` sen `"type": "module"` (non bloqueante). |

### 6.2 Smoke manual (operador / pre-release)

Checklist recomendado (sin E2E automatizado no repo); marcar en despregue real:

| Área | Comprobar |
|------|-----------|
| **Auth** | Login (`/auth/login`), erro credenciais, logout dende shell; rexistro (`/auth/register`) carga. |
| **Workshop** | `/workshop` lista; abrir fluxo «Nova visita» ou navegar a `/workshop/[id]` se hai datos. |
| **Parts** | `/admin/parts` lista + filtros sen erro de consola. |
| **Accounting** | `/accounting` e unha lista con modais (p. ex. **facturas emitidas** `/accounting/issued-invoices`) abrir/pechar CTA. |

Cobertura **RTL** xa exercita varios fluxos de modais/listas; este checklist cubre navegación real + API.

### 6.3 Decisión **GO** e outcomes

- **GO:** a barra oficial de comandos (§Comandos de calidade) está verde coa evidencia da táboa 6.1; `package.json` mantén **Next 16.2.4** + **React 19** + Tailwind v4 (change completado segundo `tasks.md`).
- **Outcomes:** toolchain Next 16 + ESLint flat; Tailwind v4 + `tw-animate-css`; auth unificado en Zustand; piloto RSC `my-invoices`; notas de factura con `useOptimistic` + `form action`; inventario `useEffect` documentado (§4.4).

### 6.4 CI GitHub Actions

- **Sen cambios** en `.github/workflows/ci.yml`: xa usa **Node 22**, `pnpm install --frozen-lockfile`, `pnpm lint` / `typecheck` / `test -- --passWithNoTests` / `build` con `NEXT_PUBLIC_API_URL` no paso build — coincide coa matriz 6.1.

## Rollback

- `git revert` do rango de commits do change (ou revert puntual por commit).
- **SHA de referencia** (árbol ao pechar verificación documentada): `e636fd4263a0ed8f06170e2ee7495cf349c09f08`.
