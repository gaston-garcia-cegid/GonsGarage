# Design: Criação em modal na Contabilidade

## Technical approach

Mapear `openspec/changes/.../specs/ui-accounting-staff/spec.md`: cada uma das quatro `page.tsx` de listagem passa a ter estado `createOpen`, botões que abrem **Radix Dialog** já exposto em `@/components/ui/dialog.tsx` (overlay `z-50`, ESC, foco), com o **mesmo corpo de formulário** que hoje vive em `*/new/page.tsx`. Serviços (`supplierService`, `receivedInvoiceService`, etc.) **inalterados**. Sucesso: fechar dialog, `await load()` na lista. Rotas `*/new` fazem **`replace` para a lista com `?create=1`**; a lista lê o query uma vez, abre o modal e limpa o query (`router.replace` sem search) para não repetir ao refrescar.

**Rota de referência UX (success criteria da proposta)**: `frontend/src/app/employees/page.tsx` — padrão toolbar + modal in-page + `onSuccess` → refresh (`fetchEmployees`). Para **primitivas** preferir **Dialog shadcn** (accounting já usa `Button` de `components/ui/`); o overlay manual de `EmployeeModal` não é copiado literalmente para evitar duplicar z-index/scroll.

## Architecture decisions

| Decision | Choice | Alternatives | Rationale |
|----------|--------|--------------|-----------|
| Primitiva modal | `@/components/ui/dialog` (Radix) | Overlay inline estilo `EmployeeModal`; `Sheet` | ESC/foco/portal consistentes com `ui-component-system`; menos código custom. |
| Pós-criação | Fechar + `load()` na lista | `router.replace(/.../id)` como hoje em `new/page` | Spec exige lista a reflectir o novo registo sem depender de `*/new`; detalhe continua acessível pela linha da tabela. |
| Legacy `*/new` | `redirect`/`replace` → lista `?create=1` | Página `new` mínima que só monta modal | Bookmarks e links antigos abrem o mesmo fluxo modal sem página full-page órfã. |
| Organização de código | `CreateXxxForm` + `XxxCreateDialog` por domínio (ou form partilhado + wrapper) | Um mega-modal genérico | Quatro payloads distintos; extrair só **shell** partilhado se repetir header/footer. |

## Data flow

```
List page (AppShell)
  ├─ toolbar / empty CTA → setCreateOpen(true) ou ?create=1
  ├─ <Dialog open> → form fields → service.create(payload)
  ├─ success → onOpenChange(false) + load()
  └─ error → setError inline no corpo do dialog
```

## File changes

| File | Action | Description |
|------|--------|-------------|
| `frontend/src/app/accounting/suppliers/page.tsx` | Modify | Estado modal; `Dialog`; trocar `Link` primário por `Button`; empty state abre modal. |
| `frontend/src/app/accounting/suppliers/new/page.tsx` | Modify | Client mínimo: `useEffect` + `replace('/accounting/suppliers?create=1')` ou `redirect` em Server Component se migrar. |
| `frontend/src/app/accounting/received-invoices/page.tsx` (+ `new`) | Modify | Idem padrão. |
| `frontend/src/app/accounting/billing-documents/page.tsx` (+ `new`) | Modify | Idem. |
| `frontend/src/app/accounting/issued-invoices/page.tsx` (+ `new`) | Modify | Idem. |
| `frontend/src/app/accounting/components/AccountingCreateDialogShell.tsx` (opcional) | Create | Se os quatro dialogs repetirem header/footer/scroll; caso contrário inline por pasta. |
| `frontend/src/app/accounting/accounting.module.css` | Modify | Classes para corpo scrollável do dialog (`max-h`, `overflow-y`) alinhadas a tokens existentes. |

*(Forms podem ficar em `suppliers/SupplierCreateForm.tsx` etc., ou co-localizados na `page` até extrair.)*

## Interfaces / contracts

- Props típicas do wrapper: `{ open, onOpenChange, onCreated?: () => void }` onde `onCreated` dispara `load()`.
- Serviços: mesmas assinaturas actuais; nenhum contrato API novo.

## Testing strategy

| Layer | What | How |
|-------|------|-----|
| Integration | Abrir modal (toolbar + `?create=1`), cancelar, submit com serviço mockado | RTL + Vitest por lista (padrão `AppShell.test` / mocks de serviço). |
| E2E | — | Não requerido neste change (config `e2e: false`). |

## Migration / rollout

Sem migração de dados nem feature flags. Deploy único: listas passam a modal; `*/new` redireccionam.

## Open questions

- [ ] Confirmar com produto se, após criar, se mantém **só** na lista ou se se deseja **atalho** «Abrir detalhe» no toast/banner (fora do spec mínimo).
