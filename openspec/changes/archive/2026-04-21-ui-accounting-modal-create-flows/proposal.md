# Proposal: Criação em modal na Contabilidade

## Intent

Em **Contabilidade**, os quatro fluxos (Fornecedores, Faturas recebidas, Documentos de faturação, Faturas emitidas) usam `*/new` em página cheia. Noutras áreas staff (**Funcionários**, **Marcações**), criação usa **modais no contexto da lista**. Unificar para o mesmo padrão reduz fragmentação UX e navegação desnecessária.

## Scope

### In Scope

- Quatro listas staff: modal de criação com a mesma lógica que `*/new/page.tsx` (extrair formulário partilhável).
- Toolbar, CTA em estado vazio, overlay, fecho (ESC/botão), **pt_PT**, alinhado a `employees/page.tsx` e `AppointmentContainer`.
- `*/new`: redirect para lista ou remoção com tratamento de bookmarks.

### Out of Scope

- API e regras P1 de `suppliers` / `invoices` / `billing`; edição em `/[id]` salvo sinergia trivial; billing cliente fora do existente.

## Capabilities

### New Capabilities

- `ui-accounting-staff`: Hub + quatro listas: criação **SHALL** ser **modal** na listagem, coerente com modais staff existentes.

### Modified Capabilities

- None (só superfície; CRUD nos specs de domínio inalterado).

## Approach

Extrair corpo de formulário de cada `new/page`, wrapper modal na `page` da lista, `onSuccess` → fechar + refresh. Reusar `components/ui` ou padrão estrutural de `EmployeeModal`. RTL: abrir, submit mock, cancelar.

## Affected Areas

| Area | Impact |
|------|--------|
| `frontend/src/app/accounting/{suppliers,received-invoices,billing-documents,issued-invoices}/` | Modified |
| `accounting.module.css` | Modified |
| `openspec/specs/ui-accounting-staff/spec.md` | New (delta) |

## Risks

| Risk | L | Mitigation |
|------|---|------------|
| Form longo | M | Corpo do modal com scroll; `max-height` como Employee modal. |
| Duplicação lista/new | M | Um componente de form; `new` só redirect. |

## Rollback Plan

Revert do PR: restaurar `Link` → `*/new` e páginas `new` como entrada única.

## Dependencies

Nenhuma.

## Success Criteria

- [ ] Criação nos quatro domínios **sem** navegar a `*/new`.
- [ ] Teclado/fecho e hierarquia visual alinhados a rota staff de referência (nomear no `design.md`).
- [ ] `pnpm test` + `pnpm build` verdes; `*/new` sem UX morta.
