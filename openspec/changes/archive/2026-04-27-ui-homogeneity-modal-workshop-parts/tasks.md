# Tasks: UI homogeneity — modal parity (peças + taller)

> `strict_tdd` no repo: fase 1 = testes que falham (RED); implementação nas fases 2–3 até verde; fase 4 = verificação global.

## Phase 1: Test-first (RED)

- [x] 1.1 Criar `frontend/src/app/admin/parts/page.test.tsx`: montar página com `useSearchParams` mock `create=1`; **assert** que o diálogo de criação fica visível (nome acessível / role). Deve **falhar** até existir `PartCreateModal` + wiring.
- [x] 1.2 Criar `frontend/src/app/workshop/[id]/page.test.tsx` (ou ficheiro colocado ao lado do page): mock `getServiceJob` a devolver `200` com corpo **sem** `job` válido; **assert** mensagem de erro ou estado terminal — **não** texto só "A carregar…". Deve **falhar** com o código actual.

## Phase 2: Peças — modal + redirect

- [x] 2.1 Criar `frontend/src/app/admin/parts/components/PartCreateModal.tsx` com `Dialog`, `Button`, `Input`, `Label`; props `open`, `onOpenChange`, `onSuccess(id: string)`; lógica `apiClient.createPart` espelhando `admin/parts/new/page.tsx`.
- [x] 2.2 Em `frontend/src/app/admin/parts/page.tsx`: `useSearchParams`, `useRouter`, ref anti-loop (padrão `appointments/page.tsx`); abrir modal com `?create=1` e `replace` para limpar query.
- [x] 2.3 Substituir `<Link href="/admin/parts/new">` na toolbar e na tabela vazia por acções que abrem modal ou `?create=1` (sem rota de página como fluxo primário).
- [x] 2.4 Embutir `<PartCreateModal />`; no sucesso: fechar, `void load()` para refrescar lista; opcional `router.push(/admin/parts/${id})` alinhado ao comportamento actual.
- [x] 2.5 Alterar `frontend/src/app/admin/parts/new/page.tsx` para `router.replace('/admin/parts?create=1')` em `useEffect` (client-only), retorno `null`.
- [x] 2.6 **Verificar** teste 1.1 a passar (`pnpm test` filtrado ao ficheiro).

## Phase 3: Taller — confirmação + detalhe

- [x] 3.1 Em `frontend/src/app/workshop/page.tsx`: estado `confirmOpen`; `Dialog` com resumo da viatura (`carId` → matrícula/modelo); botão **Nova visita** abre dialog; confirmar chama `createServiceJob` e depois `router.push`.
- [x] 3.2 No mesmo ficheiro: separar loading do POST do loading da lista (`loadJobs`) para não bloquear UI de forma ambígua.
- [x] 3.3 Em `frontend/src/app/workshop/[id]/page.tsx`: normalizar `id` de `useParams` (`string | string[]` → primeiro string); estado `loadState` (`idle` | `loading` | `success` | `error`).
- [x] 3.4 Após `getServiceJob`: se `error`, `loadState=error` + mensagem; se OK mas `!data?.job` ou `job.status` ausente, tratar como erro de produto ("Resposta inválida"); caso contrário `setDetail` + `success`.
- [x] 3.5 Render: `loading` → `AppLoading` ou cópia consistente; `error` → mensagem + `Link` para `/workshop`; nunca corpo vazio silencioso após fetch terminado.
- [x] 3.6 **Verificar** teste 1.2 a passar.

## Phase 4: Qualidade e limpeza

- [x] 4.1 `pnpm exec eslint` / `pnpm lint` em `frontend`; `pnpm typecheck`; `pnpm test` completo.
- [x] 4.2 Rever `frontend/src/app/admin/parts/admin-parts.module.css`: remover regras só usadas pelo form antigo se migradas para Tailwind no modal.
- [ ] 4.3 Confirmar manualmente: `/admin/parts/new` redireciona e abre modal; nova visita com confirmação; detalhe visita com dados ou erro legível.
