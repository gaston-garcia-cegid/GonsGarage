# Design: UI homogeneity — modal parity (peças + taller)

## Technical Approach

Aplicar o **mesmo padrão que `appointments`**: estado local + `Dialog` (shadcn) + `useSearchParams` para `?create=1` com `router.replace` para limpar a URL. **Peças**: extrair o formulário actual de `admin/parts/new/page.tsx` para um modal na lista; **taller**: `Dialog` de confirmação antes do `POST` (resumo da viatura seleccionada) e feedback explícito; **detalhe visita**: corrigir máquina de estados do fetch (hoje `detail`/`err` permitem **“A carregar…”** indefinido ou corpo inválido sem mensagem).

## Architecture Decisions

| Decision | Alternatives | Rationale |
|----------|--------------|-----------|
| `?create=1` + modal em `admin/parts` | Só botão que abre modal sem query | Igual a `accounting/*/new` e `appointments?schedule=1`; bookmarks e testes podem forçar abertura. |
| Novo componente `PartCreateModal` (ou nome equivalente) sob `admin/parts/components/` | Form só inline na página | Reutiliza lógica `apiClient.createPart` + `PartItemWriteBody`; página mantém lista e filtros. |
| Taller continua com `@/lib/api` nesta fase | Migrar já todos os endpoints taller para `api-client` | `api-client` **não** expõe `service-jobs`; migrar seria diff maior; o bug de detalhe resolve-se no consumidor actual. |
| Confirmação **Nova visita** com `Dialog` | Drawer / página intermédia | Cumpre spec delta (lista-primário, primitivas); menor superfície que nova rota. |
| Detalhe: `loading` boolean + validar `data?.job` | Confiar só em `detail \|\| err` | Com `200` e corpo inesperado ou `data` ausente, hoje fica-se em “A carregar…” (`detail` e `err` ambos “vazios” úteis); `{}` truthy sem `job` quebra UI. |

## Data Flow

```
Lista peças — toolbar "Nova peça"
    → setCreateOpen(true) OU ?create=1 → replace /
    → Dialog com form → apiClient.createPart
    → sucesso: fechar modal, refresh lista, opcional router /admin/parts/{id}

Lista taller — "Nova visita"
    → Dialog (matrícula + modelo da carId seleccionada)
    → Confirmar → apiClient.createServiceJob(carId)  [@/lib/api]
    → sucesso: fechar, router.push(/workshop/{id}), lista pode refrescar

/workshop/[id]
    → jobId = normalizar useParams.id (string | string[] → string)
    → loading=true → getServiceJob(jobId)
    → error OU !data?.job → mensagem + link voltar
    → ok → setDetail(data); loading=false
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `frontend/src/app/admin/parts/components/PartCreateModal.tsx` | Create | Form criação (campos actuais + `Dialog`/`Button`/`Input`/`Label`); props `open`, `onOpenChange`, `onSuccess(partId)`. |
| `frontend/src/app/admin/parts/page.tsx` | Modify | `useSearchParams` + ref padrão appointments; toolbar `Button` abre modal; `PartCreateModal`; links “Nova peça”/“Criar a primeira” → `?create=1` ou handler; inputs pesquisa migrar a `Input` onde trivial. |
| `frontend/src/app/admin/parts/new/page.tsx` | Modify | `useEffect` → `router.replace('/admin/parts?create=1')` (espelho `accounting/suppliers/new`). |
| `frontend/src/app/workshop/page.tsx` | Modify | Estado `confirmOpen`; `Dialog` com texto da viatura; `onNewVisit` só após confirmar; `loading` só no POST ou estado dedicado. |
| `frontend/src/app/workshop/[id]/page.tsx` | Modify | `jobId` normalizado; `loadState`: `idle \| loading \| success \| error`; após fetch OK exigir `data.job`; mensagens para 403/404; `AppLoading` ou texto consistente com outras rotas. |
| `frontend/src/app/admin/parts/admin-parts.module.css` | Modify | Reduzir estilos só usados pelo form se passarem para Tailwind no modal. |

## Interfaces / Contracts

- **Sem alteração de contrato HTTP** previsto; resposta `GET /service-jobs/:id` já usa `json:"job"` (handler `serviceJobDetailResponse`).
- **Cliente**: tratar resposta como `ServiceJobDetail` só se `data && typeof data.job === 'object' && 'status' in data.job`; caso contrário erro de produto (“Resposta inválida”).

## Testing Strategy

| Layer | What | Approach |
|-------|------|------------|
| Unit/RTL | Abrir modal com `?create=1`, submeter mock `createPart`, fechar e lista refrescada | `admin/parts/page.test.tsx` (padrão accounting pages). |
| Unit/RTL | Detalhe visita com mock 200 sem `job` | Esperar mensagem de erro, não “A carregar…”. |
| Unit/RTL | `jobId` array (defensivo) | Mock `useParams` com array → um só UUID usado no fetch. |
| E2E | — | Fora de escopo (sem Playwright directo). |

## Migration / Rollout

No database migration. Deploy único: frontend + redirects; utilizadores com `/admin/parts/new` guardado passam pela lista com modal aberto.

## Open Questions

- [ ] Reprodução exacta do bug “não carrega” em produção (network 200 vazio vs 403 vs param) — confirmar na fase apply com DevTools.
- [ ] Após estabilizar, **opcional**: mover `service-jobs` para `api-client` para uma única fonte de token/interceptors.
