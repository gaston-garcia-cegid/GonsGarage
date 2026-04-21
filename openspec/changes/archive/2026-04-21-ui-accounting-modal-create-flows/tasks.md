# Tasks: Criação em modal na Contabilidade

> TDD (`strict_tdd` no repo): por domínio, **RED** (teste que falha) antes do **GREEN** (implementação que satisfaz o spec `ui-accounting-staff`).

## Phase 1: Foundation

- [x] 1.1 Em `frontend/src/app/accounting/accounting.module.css`, adicionar classes reutilizáveis para corpo do dialog (scroll, `max-height`) alinhadas a tokens existentes.
- [x] 1.2 Se o markup Dialog repetir 4× o mesmo header/footer, criar `frontend/src/app/accounting/components/AccountingCreateDialogShell.tsx`; senão manter inline até o REFACTOR. *(Decisão: mantido inline nesta entrega.)*

## Phase 2: Fornecedores (TDD)

- [x] 2.1 **RED**: `frontend/src/app/accounting/suppliers/page.test.tsx` — mock `useAuth` + `supplierService`; abrir modal via botão «Novo fornecedor» e assert `dialog` / título pt_PT; opcional: `?create=1` abre modal (mock `useSearchParams` ou wrapper).
- [x] 2.2 **GREEN**: Extrair formulário para `frontend/src/app/accounting/suppliers/SupplierCreateForm.tsx` (lógica de `new/page.tsx`); sucesso → `onCreated()` sem `router.replace` para detalhe; erros inline.
- [x] 2.3 **GREEN**: `frontend/src/app/accounting/suppliers/page.tsx` — `Dialog` de `@/components/ui/dialog`, estado `createOpen`, toolbar `Button` + empty state, `load()` após sucesso, consumir `create=1` e limpar query com `router.replace` sem search.
- [x] 2.4 **GREEN**: `frontend/src/app/accounting/suppliers/new/page.tsx` — `useEffect` + `router.replace('/accounting/suppliers?create=1')` (client mínimo).
- [x] 2.5 **REFACTOR**: Duplicação → usar shell da 1.2 só se já existir; senão deixar para fase final.

## Phase 3: Faturas recebidas (TDD)

- [x] 3.1 **RED**: `frontend/src/app/accounting/received-invoices/page.test.tsx` — abrir modal (toolbar + empty), mock `receivedInvoiceService`.
- [x] 3.2 **GREEN**: `ReceivedInvoiceCreateForm.tsx` + `page.tsx` com Dialog, `load()`, query `create=1`, pt_PT.
- [x] 3.3 **GREEN**: `received-invoices/new/page.tsx` → `replace` para `/accounting/received-invoices?create=1`.

## Phase 4: Documentos de faturação (TDD)

- [x] 4.1 **RED**: `frontend/src/app/accounting/billing-documents/page.test.tsx` — mock serviço de billing + abrir modal.
- [x] 4.2 **GREEN**: Form + `page.tsx` + `new/page.tsx` com o mesmo padrão da Phase 3 (URLs sob `/accounting/billing-documents`).

## Phase 5: Faturas emitidas (TDD)

- [x] 5.1 **RED**: `frontend/src/app/accounting/issued-invoices/page.test.tsx` — mock serviço + modal.
- [x] 5.2 **GREEN**: Form + `page.tsx` + `new/page.tsx` com o mesmo padrão (URLs sob `/accounting/issued-invoices`).

## Phase 6: Verificação

- [x] 6.1 `cd frontend && pnpm test -- --passWithNoTests` e `pnpm typecheck` e `pnpm build` — verde.
- [x] 6.2 Revisão manual: ESC/fechar, cancelar sem persistir, bookmark `*/new` abre fluxo modal (spec «Rotas legacy»). *(Checklist documentada; validação manual do autor da entrega.)*
