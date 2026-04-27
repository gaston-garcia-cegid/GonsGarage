# Apply progress — `nextjs-16-react19-migration`

**Mode**: Strict TDD (`openspec/config.yaml`)  
**Batches completed**: Phases 2–6

## TDD Cycle Evidence (Phase 2) — cumulative

| Task | RED | GREEN | REFACTOR | Triangulation |
|------|-----|-------|----------|---------------|
| 2.1–2.5 | Contract Tailwind/PostCSS/globals | Config + tests verdes | Fix sintaxis contract | Ver `apply-progress` histórico Fase 2 |

## TDD Cycle Evidence (Phase 3) — cumulative

| Task | RED | GREEN | REFACTOR | Triangulation |
|------|-----|-------|----------|---------------|
| 3.1–3.6 | `auth.client-role` + apiClient | Store único + `AuthProvider` fino | Delegación `@/lib/api` → `api-client` | Ver evidencia Fase 3 previa |

## TDD Cycle Evidence (Phase 4) — cumulative

| Task | RED | GREEN | REFACTOR | Triangulation |
|------|-----|-------|----------|---------------|
| 4.1 Auditoría RSC | Contract: `auth/login/page.tsx` sen `'use client'`; `my-invoices/layout.tsx` servidor; list `page.tsx` servidor | `sdd-nextjs16-migration.contract.test.ts` (+3 tests) | `login/page.tsx` só Suspense + `LoginForm` | Táboa ADR 0003 §4.1 |
| 4.2 Lista piloto | `MyInvoicesListClient.test.tsx`: con `initialItems` non chama `listMine` | Split `my-invoices/page.tsx` + `MyInvoicesListClient.tsx` + layout servidor + `MyInvoicesAuthGate` | — | Caso baleiro chama `listMine` unha vez |
| 4.3 Detalle piloto | — | `my-invoices/[id]/page.tsx` servidor + `MyInvoiceDetailClient` + `fetchMyInvoiceDetailAuthenticated` | — | Mesmo contrato `/api/v1/invoices/:id` |
| 4.4 ADR useEffect | — | `docs/adr/0003-...` §4.4 inventario | — | — |

## Comandos (última verificación Phase 4)

```text
pnpm test
pnpm build
```

## Archivos tocados (Phase 4)

| Archivo | Acción |
|---------|--------|
| `frontend/src/lib/server/my-invoices-initial.ts` | Novo — fetch servidor opcional (cookie `auth_token`) |
| `frontend/src/app/my-invoices/page.tsx` | Servidor |
| `frontend/src/app/my-invoices/MyInvoicesListClient.tsx` | Novo — illa cliente |
| `frontend/src/app/my-invoices/MyInvoicesListClient.test.tsx` | Novo — RED/GREEN skip `listMine` |
| `frontend/src/app/my-invoices/[id]/page.tsx` | Servidor |
| `frontend/src/app/my-invoices/[id]/MyInvoiceDetailClient.tsx` | Novo — illa cliente |
| `frontend/src/app/my-invoices/layout.tsx` | Servidor |
| `frontend/src/app/my-invoices/MyInvoicesAuthGate.tsx` | Novo — lóxica antiga do layout cliente |
| `frontend/src/app/auth/login/page.tsx` | Sen `'use client'` |
| `frontend/src/lib/sdd-nextjs16-migration.contract.test.ts` | +describe Phase 4 |
| `docs/adr/0003-nextjs16-react19-migration-main.md` | §Phase 4 (4.1–4.4) |
| `openspec/.../design.md` | Fila RSC piloto |
| `openspec/.../tasks.md` | Phase 4 `[x]` |

## Deviations

None — JWT segue en `localStorage`; o piloto deixa lista/detalle preparadas para cookie/BFF sen cambiar o contrato JSON do Gin.

## TDD Cycle Evidence (Phase 5) — cumulative

| Task | RED | GREEN | REFACTOR | Triangulation |
|------|-----|-------|----------|---------------|
| 5.1 `useOptimistic` notas | `MyInvoiceDetailClient.test.tsx`: optimistic durante PATCH + revert en erro | `action={async (fd)=>…}` + merge antes do `await`; `waitFor` texto optimista | Eliminar `startTransition` manual; `await Promise.resolve()` intermedio xa non necesario | `patchNotes` recibe `borrador cliente`; segunda proba reconcilia con API |
| 5.2 `use()` | — | ADR §5.2 + fila `design.md` (non adoptado, motivo) | — | Spec: documentar fronte a `useEffect` |

## Comandos (última verificación Phase 5)

```text
pnpm test
pnpm typecheck
```

## Archivos tocados (Phase 5)

| Archivo | Acción |
|---------|--------|
| `frontend/src/app/my-invoices/[id]/MyInvoiceDetailClient.tsx` | `form action` async + `useOptimistic` |
| `frontend/src/app/my-invoices/[id]/MyInvoiceDetailClient.test.tsx` | RED/GREEN mutación notas |
| `frontend/src/app/my-invoices/MyInvoicesListClient.test.tsx` | Tipo `MockInstance` para `vi.spyOn` (tsc) |
| `docs/adr/0003-nextjs16-react19-migration-main.md` | §Phase 5 |
| `openspec/.../design.md` | Filas Phase 5 |
| `openspec/.../tasks.md` | Phase 5 `[x]` |

## TDD Cycle Evidence (Phase 6) — cumulative

| Task | RED | GREEN | REFACTOR | Triangulation |
|------|-----|-------|----------|---------------|
| 6.1 CI matrix | — | `pnpm lint`, `typecheck`, `test --passWithNoTests`, `build` (NEXT_PUBLIC_API_URL) verdes | — | Paridade con `ci.yml` |
| 6.2 Smoke manual | — | Checklist §ADR 6.2 (operador) | — | E2E fora de alcance; RTL cubre modais clave |
| 6.3 ADR GO | — | ADR 0003 estado **GO**, outcomes, SHA rollback | `design.md` filas verify/CI | — |
| 6.4 ci.yml | — | ADR «CI unchanged» (Node 22, pasos xa correctos) | — | Sen diff en workflow |

## Comandos (última verificación Phase 6)

```text
cd frontend && pnpm lint
cd frontend && pnpm typecheck
cd frontend && pnpm test -- --passWithNoTests
cd frontend && set NEXT_PUBLIC_API_URL=http://localhost:8080&& pnpm build
```

(En bash: `NEXT_PUBLIC_API_URL=http://localhost:8080 pnpm build`.)

## Archivos tocados (Phase 6)

| Archivo | Acción |
|---------|--------|
| `docs/adr/0003-nextjs16-react19-migration-main.md` | §Phase 6 + GO + rollback SHA |
| `openspec/.../design.md` | Filas verify / CI |
| `openspec/.../tasks.md` | Phase 6 `[x]` |
| `openspec/.../apply-progress.md` | Esta sección |

## Siguiente

Arquivar change (`sdd-archive`) ou seguir política de merge do repo.
