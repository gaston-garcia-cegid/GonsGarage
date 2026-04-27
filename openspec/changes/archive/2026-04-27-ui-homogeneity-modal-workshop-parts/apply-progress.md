# Apply progress — ui-homogeneity-modal-workshop-parts

**Mode**: Strict TDD  
**Batches**: Phase 1 (RED) + Phase 2 (GREEN peças) + Phase 3 (GREEN taller + detalhe) + Phase 4 (qualidade)

## Completed tasks (cumulativo)

- [x] 1.1 — RTL `admin/parts/page.test.tsx`: `create=1` → dialog + heading + `replace`.
- [x] 1.2 — `workshop/[id]/page.test.tsx`: corpo inválido vs loading; happy path.
- [x] 2.1 — `PartCreateModal.tsx` (Dialog + `apiClient.createPart`).
- [x] 2.2–2.4 — `page.tsx`: `Suspense`, `useSearchParams`, ref, toolbar + empty-state botões, `handlePartCreated` → `load` + `router.push`.
- [x] 2.5 — `new/page.tsx` → `replace('/admin/parts?create=1')`.
- [x] 2.6 — `pnpm exec vitest run src/app/admin/parts/page.test.tsx` verde (5 testes).
- [x] 3.1–3.2 — `workshop/page.tsx`: `confirmOpen` + `Dialog` resumo viatura; `listLoading` vs `createSubmitting`; `createServiceJob` + `router.push`.
- [x] 3.3–3.5 — `workshop/[id]/page.tsx`: `normalizeJobId`, `loadState`, `isValidJobDetail`, mensagens erro/loading, `AppLoading` + link lista.
- [x] 3.6 — `workshop/[id]/page.test.tsx` verde (2 testes); asserções com `waitFor` para nós desligados pelo Strict Mode e estado estável pós-fetch.
- [x] 4.1 — `pnpm run lint` (0 erros; 8 avisos pré-existentes noutros ficheiros), `pnpm run typecheck`, `pnpm run test` — **22 ficheiros, 92 testes** (inclui `workshop/page.test.tsx` pós-verify).
- [x] Follow-up verify — `frontend/src/app/workshop/page.test.tsx`: RTL lista taller (`Dialog` «Nova visita», resumo viatura, `createServiceJob` + `router.push`).
- [x] 4.2 — `admin-parts.module.css`: removidas regras `.toolbar` / `.toolbar h1` (não referenciadas; toolbar na lista via `AppShell`). Restantes classes continuam em `page.tsx` e `[id]/page.tsx`.
- [ ] 4.3 — Smoke manual no browser (checklist abaixo; marcar `tasks.md` após validar).

## TDD Cycle Evidence

| Task | RED | GREEN | REFACTOR |
|------|-----|-------|----------|
| 1.1 | ✅ Written | ✅ Phase 2 — `PartCreateModal` + wiring; `vitest` 5/5 | ⏳ |
| 1.2 | ✅ Written | ✅ Phase 3 — detalhe visita + testes `waitFor`; `vitest` 2/2 | ⏳ |

## Files touched (cumulativo)

| File | Action |
|------|--------|
| `frontend/src/app/admin/parts/page.test.tsx` | Mock `useSearchParams` + testes; empty-state `button` |
| `frontend/src/app/workshop/[id]/page.test.tsx` | RED Phase 1; GREEN Phase 3 — `waitFor` + asserções estáveis |
| `frontend/src/app/admin/parts/components/PartCreateModal.tsx` | **Criado** (Phase 2) |
| `frontend/src/app/admin/parts/page.tsx` | **Reescrito** — `AdminPartsPageContent` + `Suspense` |
| `frontend/src/app/admin/parts/new/page.tsx` | Redirect só `?create=1` |
| `frontend/src/app/workshop/page.tsx` | Dialog confirmação + loading separado (Phase 3) |
| `frontend/src/app/workshop/[id]/page.tsx` | `loadState`, validação detalhe, UX erro/loading (Phase 3) |
| `frontend/src/app/workshop/page.test.tsx` | **Criado** — RTL confirmação Nova visita + navegação |
| `frontend/src/app/admin/parts/admin-parts.module.css` | Phase 4 — remoção `.toolbar` morto |

## Vitest (actual)

- Suite completa: **22 files, 92 passed** (`pnpm run test`).
- `src/app/admin/parts/page.test.tsx`: **5 passed**.
- `src/app/workshop/[id]/page.test.tsx`: **2 passed**.
- `src/app/workshop/page.test.tsx`: **2 passed**.

## Smoke manual (4.3) — checklist

1. **`/admin/parts/new`**: deve redireccionar para `/admin/parts?create=1` e abrir o modal de nova peça (comportamento coberto por testes; validar visualmente).
2. **Taller — Nova visita**: abrir diálogo de confirmação com resumo da viatura; confirmar cria e navega para detalhe.
3. **`/workshop/[id]`**: com visita válida ver dados; com resposta inválida ver mensagem de erro + link para lista (coberto por testes; validar copy/layout).

## Próximo

Marcar **4.3** `[x]` em `tasks.md` após correres o checklist no browser; depois **verify** / arquivo do change.

## Learned (tests)

`screen.findByText` pode resolver com um nó já **desligado** (`isConnected === false`) sob React Strict Mode; `waitFor(() => { expect(screen.getByText(...)).toBeInTheDocument(); })` re-consulta a árvore. Para “não loading + erro”, um único `waitFor` evita passar cedo e falhar no `queryByText` do loading após remount.

No diálogo **Nova visita**, matrícula/modelo repetem-se no `<select>` — usar `within(screen.getByRole('dialog'))` para asserções no resumo sem colidir com a opção.
